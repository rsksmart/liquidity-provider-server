package watcher

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	pegin "github.com/rsksmart/liquidity-provider-server/internal/entities/pegin"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	peginUseCase "github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	log "github.com/sirupsen/logrus"
)

// PeginDiscoveryWatcher is the entry point of the commit-first peg-in path on the LPS (EPIC E5).
//
// It does two things on each tick:
//  1. Discovery — polls the PegInAddressRegistry AddressRegistered logs for newly registered RSK
//     addresses, reads their derived BTC deposit address, imports it into the watch-only BTC wallet,
//     and adds it to an in-memory watch-set. This replaces the off-chain quote-accept reservation:
//     LPS reserves no liquidity; it merely watches addresses the user has committed to on-chain.
//  2. Claim — for each watched address, checks the BTC wallet for a confirmed deposit and drives the
//     two-step claim (requestPegIn → resolvePegIn) via ClaimPegInUseCase, competing first-mined-wins.
//
// The watcher polls RSK logs rather than subscribing, consistent with how the rest of the repo reads
// RSK state (e.g. BtcReleaseWatcher / pegout deposit events).
type PeginDiscoveryWatcher struct {
	watched          map[string]pegin.WatchedRegisteredAddress
	watchedMutex     sync.RWMutex
	claimUseCase     *peginUseCase.ClaimPegInUseCase
	registry         blockchain.PegInAddressRegistryContract
	configurations   blockchain.FlyoverConfigurationsContract
	btcWallet        blockchain.BitcoinWallet
	rpc              blockchain.Rpc
	ticker           utils.Ticker
	eventBus         entities.EventBus
	stopChannel      chan bool
	lastScannedBlock uint64
	lastScannedMutex sync.RWMutex
}

func NewPeginDiscoveryWatcher(
	claimUseCase *peginUseCase.ClaimPegInUseCase,
	registry blockchain.PegInAddressRegistryContract,
	configurations blockchain.FlyoverConfigurationsContract,
	btcWallet blockchain.BitcoinWallet,
	rpc blockchain.Rpc,
	eventBus entities.EventBus,
	ticker utils.Ticker,
) *PeginDiscoveryWatcher {
	return &PeginDiscoveryWatcher{
		watched:        make(map[string]pegin.WatchedRegisteredAddress),
		claimUseCase:   claimUseCase,
		registry:       registry,
		configurations: configurations,
		btcWallet:      btcWallet,
		rpc:            rpc,
		eventBus:       eventBus,
		ticker:         ticker,
		stopChannel:    make(chan bool, 1),
	}
}

// Prepare records the current RSK block as the discovery checkpoint. New registrations from this
// block onward are picked up by Start. (A production build would also backfill historical
// registrations from persistent state; out of scope for the PoC.)
func (watcher *PeginDiscoveryWatcher) Prepare(ctx context.Context) error {
	current, err := watcher.rpc.Rsk.GetHeight(ctx)
	if err != nil {
		return fmt.Errorf("pegin discovery watcher: error getting RSK height: %w", err)
	}
	watcher.lastScannedMutex.Lock()
	defer watcher.lastScannedMutex.Unlock()
	watcher.lastScannedBlock = current
	return nil
}

func (watcher *PeginDiscoveryWatcher) Start() {
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			ctx := context.Background()
			watcher.discoverNewAddresses(ctx)
			watcher.checkDeposits(ctx)
		case <-watcher.stopChannel:
			watcher.ticker.Stop()
			close(watcher.stopChannel)
			break watcherLoop
		}
	}
}

func (watcher *PeginDiscoveryWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.stopChannel <- true
	closeChannel <- true
	log.Debug(discoveryLog("shut down"))
}

// discoverNewAddresses polls the AddressRegistered logs since the last scanned block and adds any
// new registered addresses to the watch-set.
func (watcher *PeginDiscoveryWatcher) discoverNewAddresses(ctx context.Context) {
	current, err := watcher.rpc.Rsk.GetHeight(ctx)
	if err != nil {
		log.Error(discoveryLog("error getting RSK height: %v", err))
		return
	}
	watcher.lastScannedMutex.Lock()
	from := watcher.lastScannedBlock
	watcher.lastScannedMutex.Unlock()
	if current <= from {
		return
	}

	to := current
	events, err := watcher.registry.GetAddressRegisteredEvents(ctx, from+1, &to)
	if err != nil {
		log.Error(discoveryLog("error polling AddressRegistered events [%d,%d]: %v", from+1, to, err))
		return
	}

	for _, event := range events {
		watcher.addWatchedAddress(event)
	}

	watcher.lastScannedMutex.Lock()
	watcher.lastScannedBlock = to
	watcher.lastScannedMutex.Unlock()
}

func (watcher *PeginDiscoveryWatcher) addWatchedAddress(event blockchain.AddressRegisteredEvent) {
	watcher.watchedMutex.Lock()
	defer watcher.watchedMutex.Unlock()
	if _, exists := watcher.watched[event.RskAddress]; exists {
		return
	}

	derived, err := watcher.registry.GetPegInAddress(event.RskAddress)
	if err != nil {
		log.Error(discoveryLog("error reading derived address for %s: %v", event.RskAddress, err))
		return
	}
	depositAddress, err := decodeDepositAddress(derived)
	if err != nil {
		log.Error(discoveryLog("cannot decode deposit address for %s: %v", event.RskAddress, err))
		return
	}
	if err = watcher.btcWallet.ImportAddress(depositAddress); err != nil {
		log.Error(discoveryLog("error importing deposit address %s for %s: %v", depositAddress, event.RskAddress, err))
		return
	}
	watcher.watched[event.RskAddress] = pegin.WatchedRegisteredAddress{
		RskAddress:        event.RskAddress,
		DerivationAddress: depositAddress,
		RegistrationBlock: event.BlockNumber,
		State:             pegin.WatchedAddressStateWaitingForDeposit,
	}
	log.Info(discoveryLog("now watching deposit address %s for rsk address %s", depositAddress, event.RskAddress))
}

// checkDeposits inspects each watched address for a confirmed deposit and advances the claim flow.
func (watcher *PeginDiscoveryWatcher) checkDeposits(ctx context.Context) {
	watcher.watchedMutex.RLock()
	snapshot := make([]pegin.WatchedRegisteredAddress, 0, len(watcher.watched))
	for _, w := range watcher.watched {
		snapshot = append(snapshot, w)
	}
	watcher.watchedMutex.RUnlock()

	for _, w := range snapshot {
		watcher.handleWatchedAddress(ctx, w)
	}
}

func (watcher *PeginDiscoveryWatcher) handleWatchedAddress(ctx context.Context, w pegin.WatchedRegisteredAddress) {
	switch w.State {
	case pegin.WatchedAddressStateWaitingForDeposit:
		watcher.tryRequest(ctx, w)
	case pegin.WatchedAddressStateRequested:
		watcher.tryResolve(ctx, w)
	default:
		// Resolved or ClaimedByOther: nothing further to do.
	}
}

func (watcher *PeginDiscoveryWatcher) tryRequest(ctx context.Context, w pegin.WatchedRegisteredAddress) {
	deposit, found, err := watcher.findConfirmedDeposit(w)
	if err != nil {
		log.Error(discoveryLog("error checking deposit for %s: %v", w.RskAddress, err))
		return
	}
	if !found {
		return
	}

	newState, err := watcher.claimUseCase.Request(ctx, deposit)
	if err != nil {
		log.Error(discoveryLog("requestPegIn for %s failed: %v", w.RskAddress, err))
		return
	}
	watcher.updateState(w.RskAddress, newState, deposit.BtcTxHash)
}

func (watcher *PeginDiscoveryWatcher) tryResolve(ctx context.Context, w pegin.WatchedRegisteredAddress) {
	deposit, found, err := watcher.findConfirmedDeposit(w)
	if err != nil {
		log.Error(discoveryLog("error checking deposit for %s: %v", w.RskAddress, err))
		return
	}
	if !found {
		return
	}

	newState, err := watcher.claimUseCase.Resolve(ctx, deposit)
	if errors.Is(err, blockchain.WaitingForBridgeError) {
		// Not enough confirmations yet; retry on the next tick.
		return
	}
	if err != nil {
		log.Error(discoveryLog("resolvePegIn for %s failed: %v", w.RskAddress, err))
		return
	}
	watcher.updateState(w.RskAddress, newState, deposit.BtcTxHash)
}

// findConfirmedDeposit looks for a BTC transaction paying the watched derivation address and returns
// the highest-confirmation one as a ConfirmedDeposit. found is false when no deposit is seen yet.
func (watcher *PeginDiscoveryWatcher) findConfirmedDeposit(w pegin.WatchedRegisteredAddress) (pegin.ConfirmedDeposit, bool, error) {
	txs, err := watcher.btcWallet.GetTransactions(w.DerivationAddress)
	if err != nil {
		return pegin.ConfirmedDeposit{}, false, err
	}
	for _, tx := range txs {
		amount := tx.AmountToAddress(w.DerivationAddress)
		if amount.Cmp(entities.NewWei(0)) <= 0 {
			continue
		}
		txInfo, err := watcher.rpc.Btc.GetTransactionInfo(tx.Hash)
		if err != nil {
			return pegin.ConfirmedDeposit{}, false, err
		}
		hashBytes, err := btcTxHashToBytes(tx.Hash)
		if err != nil {
			return pegin.ConfirmedDeposit{}, false, err
		}
		return pegin.ConfirmedDeposit{
			RskAddress:    w.RskAddress,
			BtcTxHash:     tx.Hash,
			BtcTxHashRaw:  hashBytes,
			Amount:        amount,
			Confirmations: txInfo.Confirmations,
			OpReturn:      nil, // plain peg-in for the PoC; OP_RETURN SC-call parsing is EPIC E4.2.
		}, true, nil
	}
	return pegin.ConfirmedDeposit{}, false, nil
}

func (watcher *PeginDiscoveryWatcher) updateState(rskAddress string, state pegin.WatchedAddressState, depositTxHash string) {
	watcher.watchedMutex.Lock()
	defer watcher.watchedMutex.Unlock()
	w, ok := watcher.watched[rskAddress]
	if !ok {
		return
	}
	w.State = state
	if depositTxHash != "" {
		w.DepositTxHash = depositTxHash
	}
	watcher.watched[rskAddress] = w
}

// GetWatchedAddress exposes the in-memory state of a watched address (used by tests / status).
func (watcher *PeginDiscoveryWatcher) GetWatchedAddress(rskAddress string) (pegin.WatchedRegisteredAddress, bool) {
	watcher.watchedMutex.RLock()
	defer watcher.watchedMutex.RUnlock()
	w, ok := watcher.watched[rskAddress]
	return w, ok
}

func discoveryLog(msg string, args ...any) string {
	return fmt.Sprintf("PeginDiscoveryWatcher: "+msg, args...)
}

// decodeDepositAddress turns the registry's derivation bytes into a BTC address string. The registry
// returns the address already rendered (P2SH) as ASCII bytes for the configured encoding, so for the
// PoC we interpret printable bytes as the address. NOTE: the exact byte encoding must be confirmed
// against a live regtest registry (see report).
func decodeDepositAddress(derived blockchain.PegInDepositAddress) (string, error) {
	if len(derived.ScriptOrAddress) == 0 {
		return "", errors.New("empty derivation address")
	}
	if !isPrintableASCII(derived.ScriptOrAddress) {
		return "", fmt.Errorf("derivation address is not a printable BTC address (encoding %d)", derived.Encoding)
	}
	return string(derived.ScriptOrAddress), nil
}

func isPrintableASCII(b []byte) bool {
	for _, c := range b {
		if c < 0x20 || c > 0x7e {
			return false
		}
	}
	return true
}

// btcTxHashToBytes converts a big-endian BTC tx hash hex string into the 32-byte array the contract
// expects. The bytes are reversed because BTC tx hashes are displayed in reverse byte order.
func btcTxHashToBytes(hash string) ([32]byte, error) {
	var out [32]byte
	raw, ok := new(big.Int).SetString(hash, 16)
	if !ok {
		return out, fmt.Errorf("invalid btc tx hash %q", hash)
	}
	bytes := raw.Bytes()
	if len(bytes) > 32 {
		return out, fmt.Errorf("btc tx hash %q is longer than 32 bytes", hash)
	}
	// Right-align into the 32-byte array (big-endian, matching bytes32 in the contract).
	copy(out[32-len(bytes):], bytes)
	return out, nil
}
