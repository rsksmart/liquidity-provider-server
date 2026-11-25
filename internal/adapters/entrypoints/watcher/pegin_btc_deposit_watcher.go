package watcher

import (
	"context"
	"fmt"
	"maps"
	"math/big"
	"slices"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	log "github.com/sirupsen/logrus"
)

type ClaimPegInUseCase interface {
	Run(ctx context.Context, watch rootstock.PegInWatch, depositTxID string) error
}

type PeginDepositAddressWatcher struct {
	watches             map[string]rootstock.PegInWatch
	getWatchedUseCase   *w.GetWatchedRegisteredAddressesUseCase
	claimUseCase        ClaimPegInUseCase
	btcWallet           blockchain.BitcoinWallet
	rpc                 blockchain.Rpc
	ticker              utils.Ticker
	eventBus            entities.EventBus
	importedEvents      <-chan entities.Event
	hasUncheckedWatches bool
	watcherStopChannel  chan bool
	currentBlock        *big.Int
	currentBlockMutex   sync.RWMutex
}

type PeginDepositAddressWatcherUseCases struct {
	claimUseCase      ClaimPegInUseCase
	getWatchedUseCase *w.GetWatchedRegisteredAddressesUseCase
}

func NewPeginDepositAddressWatcherUseCases(
	claimUseCase ClaimPegInUseCase,
	getWatchedUseCase *w.GetWatchedRegisteredAddressesUseCase,
) *PeginDepositAddressWatcherUseCases {
	return &PeginDepositAddressWatcherUseCases{
		claimUseCase:      claimUseCase,
		getWatchedUseCase: getWatchedUseCase,
	}
}

func NewPeginDepositAddressWatcher(
	useCases *PeginDepositAddressWatcherUseCases,
	btcWallet blockchain.BitcoinWallet,
	rpc blockchain.Rpc,
	eventBus entities.EventBus,
	ticker utils.Ticker,
) *PeginDepositAddressWatcher {
	watcherStopChannel := make(chan bool, 1)
	return &PeginDepositAddressWatcher{
		watches:            make(map[string]rootstock.PegInWatch),
		claimUseCase:       useCases.claimUseCase,
		getWatchedUseCase:  useCases.getWatchedUseCase,
		btcWallet:          btcWallet,
		eventBus:           eventBus,
		watcherStopChannel: watcherStopChannel,
		rpc:                rpc,
		ticker:             ticker,
		currentBlockMutex:  sync.RWMutex{},
	}
}

func (watcher *PeginDepositAddressWatcher) Prepare(ctx context.Context) error {
	watcher.currentBlockMutex.Lock()
	defer watcher.currentBlockMutex.Unlock()
	watcher.currentBlock = big.NewInt(0)

	watcher.importedEvents = watcher.eventBus.Subscribe(blockchain.PegInAddressesImportedEventId)
	watches, err := watcher.getWatchedUseCase.Run(ctx)
	if err != nil {
		return fmt.Errorf("load the imported PegIn addresses: %w", err)
	}
	for _, watch := range watches {
		watcher.rememberWatch(watch)
	}
	return nil
}

func (watcher *PeginDepositAddressWatcher) Start() {
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.checkDepositsWhenDue(context.Background())
		case event, open := <-watcher.importedEvents:
			if open {
				watcher.rememberImportedWatches(event)
			} else {
				watcher.importedEvents = nil
			}
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *PeginDepositAddressWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- true
	closeChannel <- true
	log.Debug(LogPeginBtcShutdown)
}

func (watcher *PeginDepositAddressWatcher) checkDepositsWhenDue(ctx context.Context) {
	watcher.currentBlockMutex.Lock()
	defer watcher.currentBlockMutex.Unlock()
	height, err := watcher.rpc.Btc.GetHeight()
	if err != nil {
		log.Errorf(LogPeginBtcChainHeight, err)
		return
	}
	if !watcher.hasUncheckedWatches && height.Cmp(watcher.currentBlock) <= 0 {
		return
	}
	if watcher.checkDeposits(ctx) {
		watcher.currentBlock = height
		watcher.hasUncheckedWatches = false
	}
}

func (watcher *PeginDepositAddressWatcher) rememberWatch(watch rootstock.PegInWatch) {
	if watch.State == rootstock.PegInWatchImported {
		watcher.watches[watch.BtcAddress] = watch
	}
}

func (watcher *PeginDepositAddressWatcher) rememberImportedWatches(event entities.Event) {
	imported, ok := event.(blockchain.PegInAddressesImportedEvent)
	if !ok {
		log.Error(LogPeginBtcWrongEvent)
		return
	}
	for _, watch := range imported.Watches {
		watcher.rememberWatch(watch)
	}
	watcher.hasUncheckedWatches = true
}

type pegInDeposit struct {
	btcAddress string
	txID       string
}

func (watcher *PeginDepositAddressWatcher) checkDeposits(ctx context.Context) (listed bool) {
	outputs, err := watcher.btcWallet.GetUnspentOutputs(slices.Sorted(maps.Keys(watcher.watches)))
	if err != nil {
		log.Errorf(LogPeginBtcListDepositsError, err)
		return false
	}
	checked := make(map[pegInDeposit]bool)
	for _, output := range outputs {
		deposit := pegInDeposit{btcAddress: output.Address, txID: output.TxID}
		watch, isWatched := watcher.watches[output.Address]
		if isWatched && !checked[deposit] {
			checked[deposit] = true
			watcher.claim(ctx, watch, output.TxID)
		}
	}
	return true
}

func (watcher *PeginDepositAddressWatcher) claim(ctx context.Context, watch rootstock.PegInWatch, depositTxID string) {
	if err := watcher.claimUseCase.Run(ctx, watch, depositTxID); err != nil {
		log.Error(LogPeginBtcClaimError(watch.RskAddress, depositTxID, err))
	}
}

func (watcher *PeginDepositAddressWatcher) GetCurrentBlock() *big.Int {
	watcher.currentBlockMutex.RLock()
	defer watcher.currentBlockMutex.RUnlock()
	return watcher.currentBlock
}
