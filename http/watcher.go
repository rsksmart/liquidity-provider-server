package http

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	mongoDB "github.com/rsksmart/liquidity-provider-server/mongo"
	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider-server/pegout"

	"github.com/btcsuite/btcutil"
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
	hash              string
	derivationAddress string
	btc               connectors.BTCConnector
	rsk               connectors.RSKConnector
	lp                pegout.LiquidityProvider
	dbMongo           mongoDB.DB
	state             types.RQState
	quote             *pegout.Quote
	done              chan struct{}
	closed            bool
	signature         []byte
	sharedLocker      sync.Locker
}

type RegisterPegoutWatcher struct {
	hash              string
	btc               connectors.BTCConnector
	rsk               connectors.RSKConnector
	lp                pegout.LiquidityProvider
	dbMongo           mongoDB.DB
	state             types.RQState
	quote             *pegout.Quote
	done              chan struct{}
	closed            bool
	signature         []byte
	sharedLocker      sync.Locker
	derivationAddress string
}

const (
	pegInGasLim           = 250000
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
	w.closeAndUpdateQuoteState(types.RQStateRegisterPegInSucceeded)
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

func (r *RegisterPegoutWatcher) OnRegisterPegOut(newState types.RQState) {
	if r.closed {
		log.Errorf(WatcherClosedError, r.hash)
		return
	}

	if newState == types.RQStateRegisterPegInFailed {
		err := r.closeAndUpdateQuoteState(newState)
		if err != nil {
			log.Errorf("error calling registerPegOut. value: %v. error: %v", r.hash, err)
			return
		}
	}

	if newState == types.RQStateCallForUserSucceeded {
		if newState != r.state {
			txHash, err := r.btc.SendBTC(r.derivationAddress, uint(r.quote.Value))
			if err != nil {
				log.Errorf("Error to send %v BTC to %v of quote hash %v", r.derivationAddress, r.quote.Value, r.hash)
				log.Errorf("Error: %v", err)
				return
			}

			log.Infof("it was sent %v BTC to %v of quote hash %v (transaction hash: %v)", r.derivationAddress, r.quote.Value, r.hash, txHash)

			r.updateQuoteState(newState)
			r.close()
		}
	}
}

func (r *RegisterPegoutWatcher) OnExpire() {
	if r.closed {
		log.Errorf(WatcherOnExpireError, r.hash)
		return
	}
	log.Debugf(TimeExpiredError, r.hash)
	_ = r.closeAndUpdateQuoteState(types.RQStateTimeForDepositElapsed)
}

func (r *RegisterPegoutWatcher) Done() <-chan struct{} {
	return r.done
}

func (r *RegisterPegoutWatcher) closeAndUpdateQuoteState(newState types.RQState) error {
	r.close()
	return r.updateQuoteState(newState)
}

func (b *BTCAddressPegOutWatcher) closeAndUpdateQuoteState(newState types.RQState) error {
	b.close()
	return b.updateQuoteState(newState)
}

func (r *RegisterPegoutWatcher) close() {
	r.closed = true
	close(r.done)
}

func (b *BTCAddressPegOutWatcher) close() {
	b.closed = true
	close(b.done)
}

func (r *RegisterPegoutWatcher) updateQuoteState(newState types.RQState) error {
	err := r.dbMongo.UpdateRetainedPegOutQuoteState(r.hash, r.state, newState)
	if err != nil {
		log.Errorf(UpdateQuoteStateError, r.hash, err)
		return err
	}

	r.state = newState
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
