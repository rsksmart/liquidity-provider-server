package http

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/common"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	mongoDB "github.com/rsksmart/liquidity-provider-server/mongo"
	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider-server/pegout"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rsksmart/liquidity-provider-server/connectors"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
)

type BTCAddressWatcher struct {
	hash         string
	btc          connectors.BTCConnector
	rsk          connectors.RSKConnector
	lp           pegin.LiquidityProvider
	dbMongo      mongoDB.DBConnector
	state        types.RQState
	quote        *pegin.Quote
	done         chan struct{}
	closed       bool
	signature    []byte
	sharedLocker sync.Locker
}

type BTCAddressPegOutWatcher struct {
	hash                 string
	derivationAddress    string
	addressDecryptionKey string
	btc                  connectors.BTCConnector
	rsk                  connectors.RSKConnector
	lp                   pegout.LiquidityProvider
	dbMongo              mongoDB.DBConnector
	state                types.RQState
	quote                *pegout.Quote
	done                 chan struct{}
	closed               bool
	signature            []byte
	sharedLocker         sync.Locker
}

type RegisterPegoutWatcher struct {
	hash              string
	btc               connectors.BTCConnector
	rsk               connectors.RSKConnector
	lp                pegout.LiquidityProvider
	dbMongo           mongoDB.DBConnector
	state             types.RQState
	quote             *pegout.Quote
	done              chan struct{}
	closed            bool
	signature         []byte
	sharedLocker      sync.Locker
	derivationAddress common.Address
}

const (
	pegInGasLim           = 1500000
	CFUExtraGas           = 150000
	WatcherClosedError    = "watcher is closed; cannot handle OnNewConfirmation; hash: %v"
	WatcherOnExpireError  = "watcher is closed; cannot handle OnExpire; hash: %v"
	TimeExpiredError      = "time has expired for quote: %v"
	UpdateQuoteStateError = "error updating quote state; hash: %v; error: %v"
)

func NewBTCAddressWatcher(hash string,
	btc connectors.BTCConnector, rsk connectors.RSKConnector, provider pegin.LiquidityProvider,
	dbMongo mongoDB.DBConnector, q *pegin.Quote, signature []byte, state types.RQState, sharedLocker sync.Locker) *BTCAddressWatcher {
	watcher := BTCAddressWatcher{
		hash:         hash,
		btc:          btc,
		rsk:          rsk,
		lp:           provider,
		dbMongo:      dbMongo,
		quote:        q,
		state:        state,
		signature:    signature,
		done:         make(chan struct{}),
		sharedLocker: sharedLocker,
	}
	return &watcher
}

func (w *BTCAddressWatcher) OnNewConfirmation(txHash string, confirmations int64, amount btcutil.Amount) {
	if w.closed {
		log.Errorf(WatcherClosedError, w.hash)
		return
	}
	log.Debugf("processing OnNewConfirmation event for tx %v; confirmations: %v; received amount: %v", txHash, confirmations, amount)

	if w.state == types.RQStateWaitingForDeposit && confirmations >= int64(w.quote.Confirmations) {
		err := w.performCallForUser()
		if err != nil {
			log.Errorf("error calling callForUser. value: %v. error: %v", txHash, err)
			return
		}
		log.Debugf("registered callforuser for tx %v", txHash)
	}

	if w.state == types.RQStateCallForUserSucceeded && confirmations >= w.rsk.GetRequiredBridgeConfirmations() {
		err := w.performRegisterPegIn(txHash)
		if err != nil {
			log.Errorf("error calling registerPegIn. value: %v. error: %v", txHash, err)
		}
	}
}

func (w *BTCAddressWatcher) OnExpire() {
	if w.closed {
		log.Errorf(WatcherOnExpireError, w.hash)
		return
	}
	log.Debugf(TimeExpiredError, w.hash)
	_ = w.closeAndUpdateQuoteState(types.RQStateTimeForDepositElapsed)
}

func (w *BTCAddressWatcher) Done() <-chan struct{} {
	return w.done
}

func (w *BTCAddressPegOutWatcher) OnNewConfirmation(txHash string, confirmations int64, amount btcutil.Amount) {
	if w.closed {
		log.Errorf(WatcherClosedError, w.hash)
		return
	}
	log.Debugf("processing OnNewConfirmation event for tx %v; confirmations: %v; received amount: %v", txHash, confirmations, amount)

	if !(w.state == types.RQStateCallForUserSucceeded && confirmations >= int64(w.quote.TransferConfirmations)) {
		return
	}

	quote, err := w.rsk.ParsePegOutQuote(w.quote)
	if err != nil {
		log.Error("Error parsing pegout quote: ", err)
		_ = w.closeAndUpdateQuoteState(types.RQStateCallForUserFailed)
	}
	opt := &bind.TransactOpts{
		GasLimit: pegInGasLim,
		Value:    nil,
		From:     common.HexToAddress(w.quote.LPRSKAddr),
		Signer:   w.lp.SignTx,
	}

	mb, err := w.btc.BuildMerkleBranch(txHash)
	if err != nil {
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		log.Error("Error refunding pegout: ", err)
	}
	bhh, err := w.btc.GetBlockHeaderHashByTx(txHash)
	if err != nil {
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		log.Error("Error refunding pegout: ", err)
	}

	btcTxHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		log.Error("Error refunding pegout: ", err)
	}

	var btcTxHashBytes [32]byte
	copy(btcTxHashBytes[:], btcTxHash[:])

	var mbHashes [][32]byte
	var mbHash [32]byte
	for _, hash := range mb.Hashes {
		copy(mbHash[:], hash[:])
		mbHashes = append(mbHashes, mbHash)
	}

	w.sharedLocker.Lock()
	defer w.sharedLocker.Unlock()

	tx, err := w.rsk.RefundPegOut(opt, quote, btcTxHashBytes, bhh, big.NewInt(int64(mb.Path)), mbHashes)
	if err != nil && strings.Contains(err.Error(), "LBC: Don't have required confirmations") {
		log.Errorf("Error calling RefundPegout: %v. Retrying on next confirmation", err)
		return
	} else if err != nil {
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		log.Error("Error refunding pegout: ", err)
		return
	}
	s, err := w.rsk.GetTxStatus(context.Background(), tx)
	if err != nil || !s {
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		log.Error("Error refunding pegout: ", err)
	}

	keys, err := w.dbMongo.GetAddressKeys(w.hash)
	if err != nil {
		log.Errorf("Error sending RBTC to the bridge on pegout quote %s: %s", w.hash, err)
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		return
	}
	privateKey, err := decrypt(keys.PrivateKey, []byte(w.addressDecryptionKey))
	if err != nil {
		log.Errorf("Error sending RBTC to the bridge on pegout quote %s: %s", w.hash, err)
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		return
	}
	err = w.rsk.SendRbtc(string(privateKey), w.rsk.GetBridgeAddress().Hex(), new(types.Wei).Add(w.quote.Value, w.quote.CallFee).Uint64())
	if err != nil {
		log.Errorf("Error sending RBTC to the bridge on pegout quote %s: %s", w.hash, err)
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		return
	}
	_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInSucceeded)
}

func (w *BTCAddressPegOutWatcher) OnExpire() {
	if w.closed {
		log.Errorf(WatcherOnExpireError, w.hash)
		return
	}
	log.Debugf(TimeExpiredError, w.hash)
	_ = w.closeAndUpdateQuoteState(types.RQStateTimeForDepositElapsed)
}

func (w *BTCAddressPegOutWatcher) Done() <-chan struct{} {
	return w.done
}

func (w *BTCAddressWatcher) performCallForUser() error {
	q, err := w.rsk.ParseQuote(w.quote)
	if err != nil {
		_ = w.closeAndUpdateQuoteState(types.RQStateCallForUserFailed)
		return err
	}

	w.sharedLocker.Lock()
	defer w.sharedLocker.Unlock()

	lbcBalance, err := w.rsk.GetLbcBalance(w.lp.Address())
	if err != nil {
		return err
	}

	var val *big.Int
	if lbcBalance.Cmp(q.Value) >= 0 { // lbc balance is sufficient, no need to transfer any value
		val = nil
	} else { // lbc balance is not sufficient, calc delta to transfer
		val = new(big.Int).Sub(q.Value, lbcBalance)
	}

	opt := &bind.TransactOpts{
		GasLimit: uint64(q.GasLimit + CFUExtraGas),
		Value:    val,
		From:     q.LiquidityProviderRskAddress,
		Signer:   w.lp.SignTx,
	}
	tx, err := w.rsk.CallForUser(opt, q)
	if err != nil {
		_ = w.closeAndUpdateQuoteState(types.RQStateCallForUserFailed)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*8760) // timeout is a year
	defer cancel()
	s, err := w.rsk.GetTxStatus(ctx, tx)
	if err != nil || !s {
		_ = w.closeAndUpdateQuoteState(types.RQStateCallForUserFailed)
		return fmt.Errorf("transaction failed. hash: %v", tx.Hash())
	}

	err = w.updateQuoteState(types.RQStateCallForUserSucceeded)
	if err != nil {
		w.close()
		return err
	}
	return nil
}

func (w *BTCAddressWatcher) performRegisterPegIn(txHash string) error {
	q, err := w.rsk.ParseQuote(w.quote)
	if err != nil {
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		return err
	}
	opt := &bind.TransactOpts{
		GasLimit: pegInGasLim,
		Value:    nil,
		From:     q.LiquidityProviderRskAddress,
		Signer:   w.lp.SignTx,
	}
	rawTx, err := w.btc.SerializeTx(txHash)
	if err != nil {
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		return err
	}
	pmt, err := w.btc.SerializePMT(txHash)
	if err != nil {
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		return err
	}
	bh, err := w.btc.GetBlockNumberByTx(txHash)
	if err != nil {
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		return err
	}

	err = w.rsk.RegisterPegInWithoutTx(q, w.signature, rawTx, pmt, big.NewInt(bh))
	if err != nil {
		if strings.Contains(err.Error(), "Failed to validate BTC transaction") {
			log.Debugf("bridge failed to validate BTC transaction. retrying on next confirmation. tx: %v", txHash)
			return nil // allow retrying in case the bridge didn't acknowledge all required confirmations have occurred
		}
	}

	log.Debugf("calling pegin for tx %v", txHash)
	tx, err := w.rsk.RegisterPegIn(opt, q, w.signature, rawTx, pmt, big.NewInt(bh))
	if err != nil {
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*8760) // timeout is a year
	defer cancel()
	s, err := w.rsk.GetTxStatus(ctx, tx)
	if err != nil || !s {
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		return fmt.Errorf("transaction failed. hash: %v", tx.Hash())
	}

	err = w.updateQuoteState(types.RQStateRegisterPegInSucceeded)
	if err != nil {
		w.close()
		return err
	}
	log.Debugf("registered pegin for tx %v", txHash)

	w.close()
	return nil
}

func (w *BTCAddressWatcher) updateQuoteState(newState types.RQState) error {
	err := w.dbMongo.UpdateRetainedQuoteState(w.hash, w.state, newState)
	if err != nil {
		log.Errorf(UpdateQuoteStateError, w.hash, err)
		return err
	}

	w.state = newState
	return nil
}

func (w *BTCAddressWatcher) closeAndUpdateQuoteState(newState types.RQState) error {
	w.close()
	return w.updateQuoteState(newState)
}

func (w *BTCAddressWatcher) close() {
	w.closed = true
	close(w.done)
}

func (watcher *RegisterPegoutWatcher) GetQuote() *pegout.Quote {
	return watcher.quote
}

func (watcher *RegisterPegoutWatcher) GetState() types.RQState {
	return watcher.state
}

func (watcher *RegisterPegoutWatcher) GetWatchedAddress() common.Address {
	return watcher.derivationAddress
}

func (watcher *RegisterPegoutWatcher) OnDepositConfirmationsReached() bool {
	if watcher.state != types.RQStateWaitingForDeposit {
		return false
	}

	quote, err := watcher.rsk.ParsePegOutQuote(watcher.quote)
	if err != nil {
		log.Error("Error parsing pegout quote: ", err)
		_ = watcher.closeAndUpdateQuoteState(types.RQStateCallForUserFailed)
		return false
	}

	satoshi, _ := watcher.quote.Value.ToSatoshi().Float64()
	watcher.sharedLocker.Lock()
	defer watcher.sharedLocker.Unlock()
	err = watcher.btc.LockBtc(satoshi)
	if err != nil {
		log.Error("Error locking btc: ", err)
		_ = watcher.closeAndUpdateQuoteState(types.RQStateCallForUserFailed)
		return false
	}

	opt := &bind.TransactOpts{
		From:   quote.LpRskAddress,
		Signer: watcher.lp.SignTx,
	}
	tx, err := watcher.rsk.RegisterPegOut(opt, quote, watcher.signature)
	if err != nil {
		log.Error("Error registering pegout: ", err)
		_ = watcher.closeAndUpdateQuoteState(types.RQStateCallForUserFailed)
		return false
	}

	if status, err := watcher.rsk.GetTxStatus(context.Background(), tx); err != nil || !status {
		_ = watcher.btc.UnlockBtc(satoshi)
		_ = watcher.closeAndUpdateQuoteState(types.RQStateCallForUserFailed)
		log.Errorf("transaction failed. hash: %v", tx.Hash())
		return false
	}
	err = watcher.btc.UnlockBtc(satoshi)
	if err != nil {
		log.Error("Error unlocking BTC before sending to destination: ", err)
		_ = watcher.closeAndUpdateQuoteState(types.RQStateCallForUserFailed)
		return false
	}

	_, err = watcher.btc.SendBTC(watcher.quote.DepositAddr, uint64(math.Ceil(satoshi)))
	if err != nil {
		log.Errorf("Error sending BTC to address %s on pegout quote %s: %s", watcher.quote.DepositAddr, watcher.hash, err)
		_ = watcher.closeAndUpdateQuoteState(types.RQStateCallForUserFailed)
		return false
	}

	err = watcher.updateQuoteState(types.RQStateCallForUserSucceeded)
	watcher.close()
	if err != nil {
		log.Debugf("error updating quote state for quote %s", watcher.hash)
	}
	log.Debugf("registered pegout for tx %v", tx.Hash())
	return true
}

func (watcher *RegisterPegoutWatcher) OnExpire() {
	if watcher.closed {
		log.Errorf(WatcherOnExpireError, watcher.hash)
		return
	}
	log.Debugf(TimeExpiredError, watcher.hash)
	_ = watcher.closeAndUpdateQuoteState(types.RQStateTimeForDepositElapsed)
}

func (watcher *RegisterPegoutWatcher) Done() <-chan struct{} {
	return watcher.done
}

func (watcher *RegisterPegoutWatcher) closeAndUpdateQuoteState(newState types.RQState) error {
	watcher.close()
	return watcher.updateQuoteState(newState)
}

func (b *BTCAddressPegOutWatcher) closeAndUpdateQuoteState(newState types.RQState) error {
	b.close()
	return b.updateQuoteState(newState)
}

func (watcher *RegisterPegoutWatcher) close() {
	watcher.closed = true
	close(watcher.done)
}

func (b *BTCAddressPegOutWatcher) close() {
	b.closed = true
	close(b.done)
}

func (watcher *RegisterPegoutWatcher) updateQuoteState(newState types.RQState) error {
	err := watcher.dbMongo.UpdateRetainedPegOutQuoteState(watcher.hash, watcher.state, newState)
	if err != nil {
		log.Errorf(UpdateQuoteStateError, watcher.hash, err)
		return err
	}

	watcher.state = newState
	return nil
}

func (r *BTCAddressPegOutWatcher) updateQuoteState(newState types.RQState) error {
	err := r.dbMongo.UpdateRetainedPegOutQuoteState(r.hash, r.state, newState)
	if err != nil {
		log.Errorf(UpdateQuoteStateError, r.hash, err)
		return err
	}

	r.state = newState
	return nil
}
