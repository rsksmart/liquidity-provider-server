package http

import (
	"context"
	"encoding/hex"
	"errors"
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

	unrecoverableError, err := w.performRefundPegout(txHash)
	if err != nil && unrecoverableError {
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		log.Error("Error refunding pegout: ", err)
		return
	} else if err != nil {
		log.Errorf("Error calling RefundPegout: %v. Retrying on next confirmation", err)
		return
	}

	err = w.rsk.SendRbtc(w.lp.SignTx, w.lp.Address(), w.rsk.GetBridgeAddress().Hex(), new(types.Wei).Add(w.quote.Value, w.quote.CallFee).Uint64())
	if err != nil {
		log.Errorf("Error sending RBTC to the bridge on pegout quote %s: %s", w.hash, err)
		_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInFailed)
		return
	}
	_ = w.closeAndUpdateQuoteState(types.RQStateRegisterPegInSucceeded)
}

func (w *BTCAddressPegOutWatcher) performRefundPegout(txHash string) (bool, error) {
	quote, err := w.rsk.ParsePegOutQuote(w.quote)
	if err != nil {
		return true, err
	}
	opt := &bind.TransactOpts{
		GasLimit: pegInGasLim,
		Value:    nil,
		From:     common.HexToAddress(w.quote.LPRSKAddr),
		Signer:   w.lp.SignTx,
	}

	mb, err := w.btc.BuildMerkleBranch(txHash)
	if err != nil {
		return true, err
	}
	bhh, err := w.btc.GetBlockHeaderHashByTx(txHash)
	if err != nil {
		return true, err
	}

	btcTxHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return true, err
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
		return false, err
	} else if err != nil {
		return true, err
	}
	s, err := w.rsk.GetTxStatus(context.Background(), tx)
	if err != nil || !s {
		return true, err
	}
	return false, err
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
		return fmt.Errorf("CallForUser transaction failed. hash: %v", tx.Hash())
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
		return fmt.Errorf("RegisterPegin transaction failed. hash: %v", tx.Hash())
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

func (b *BTCAddressPegOutWatcher) closeAndUpdateQuoteState(newState types.RQState) error {
	b.close()
	return b.updateQuoteState(newState)
}

func (b *BTCAddressPegOutWatcher) close() {
	b.closed = true
	close(b.done)
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

type DepositEventWatcher interface {
	Init(waitingForDepositQuotes, waitingForConfirmationQuotes map[string]*WatchedQuote)
	WatchNewQuote(quoteHash, signature string, quote *pegout.Quote) error
	EndChannel() chan<- bool
}

type DepositEventWatcherImpl struct {
	lastCheckedBlock     uint64
	nonDepositedQuotes   map[string]*WatchedQuote
	depositedQuotes      map[string]*WatchedQuote
	checkInterval        time.Duration
	endChannel           chan bool
	addLocker            sync.Locker
	rsk                  connectors.RSKConnector
	btc                  connectors.BTCConnector
	db                   mongoDB.DBConnector
	pegoutLocker         sync.Locker
	liquidityProvider    pegout.LiquidityProvider
	finalizationCallback func(hash string, quote *WatchedQuote, endState types.RQState)
}

func NewDepositEventWatcher(checkInterval time.Duration, liquidityProvider pegout.LiquidityProvider,
	addLocker sync.Locker, pegoutLocker sync.Locker, endChannel chan bool,
	rsk connectors.RSKConnector, btc connectors.BTCConnector, db mongoDB.DBConnector,
	finalizationCallback func(hash string, quote *WatchedQuote, endState types.RQState)) DepositEventWatcher {
	return &DepositEventWatcherImpl{
		checkInterval:        checkInterval,
		endChannel:           endChannel,
		addLocker:            addLocker,
		rsk:                  rsk,
		btc:                  btc,
		db:                   db,
		pegoutLocker:         pegoutLocker,
		liquidityProvider:    liquidityProvider,
		finalizationCallback: finalizationCallback,
	}
}

type WatchedQuote struct {
	Data         *pegout.Quote
	Signature    string
	DepositBlock uint64
}

func (watcher *DepositEventWatcherImpl) Init(waitingForDepositQuotes, waitingForConfirmationQuotes map[string]*WatchedQuote) {
	if waitingForDepositQuotes == nil || waitingForConfirmationQuotes == nil {
		log.Fatal("invalid initial pegout quote map")
	}
	var oldestBlock uint32
	for _, quote := range waitingForDepositQuotes {
		watcher.updateOldestBlock(quote, &oldestBlock)
	}
	for _, quote := range waitingForConfirmationQuotes {
		watcher.updateOldestBlock(quote, &oldestBlock)
	}
	watcher.lastCheckedBlock = uint64(oldestBlock)
	watcher.nonDepositedQuotes = waitingForDepositQuotes
	watcher.depositedQuotes = waitingForConfirmationQuotes
	watcher.watchDepositEvent()
}

func (watcher *DepositEventWatcherImpl) updateOldestBlock(quote *WatchedQuote, oldestBlock *uint32) {
	creationBlock := watcher.liquidityProvider.GetCreationBlock(quote.Data)
	if *oldestBlock == 0 || *oldestBlock > creationBlock {
		*oldestBlock = creationBlock
	}
}

func (watcher *DepositEventWatcherImpl) WatchNewQuote(quoteHash, signature string, quote *pegout.Quote) error {
	if watcher.nonDepositedQuotes == nil {
		return errors.New("not initialized")
	}
	watcher.addLocker.Lock()
	defer watcher.addLocker.Unlock()
	_, existsOnNonDeposited := watcher.nonDepositedQuotes[quoteHash]
	_, existsOnDeposited := watcher.depositedQuotes[quoteHash]
	if !existsOnNonDeposited && !existsOnDeposited {
		watcher.nonDepositedQuotes[quoteHash] = &WatchedQuote{Data: quote, Signature: signature}
		return nil
	} else {
		return errors.New("already watched")
	}
}

func (watcher *DepositEventWatcherImpl) watchDepositEvent() {
	ticker := time.NewTicker(watcher.checkInterval)
	for {
		select {
		case <-watcher.endChannel:
			ticker.Stop()
			return
		case <-ticker.C:
			height, err := watcher.rsk.GetRskHeight()
			if err != nil {
				log.Error("Error getting rsk height: ", err)
				break
			}
			err = watcher.checkDeposits(height)
			if err != nil {
				log.Error("Error getting pegout deposit events: ", err)
				break
			}
			quotes := watcher.getConfirmedQuotes(height)
			watcher.cleanExpiredQuotes()
			watcher.handleDepositedQuotes(quotes)
		}
	}
}

func (watcher *DepositEventWatcherImpl) checkDeposits(height uint64) error {
	if height == watcher.lastCheckedBlock {
		return nil
	}
	events, err := watcher.rsk.GetDepositEvents(watcher.lastCheckedBlock-1, height)
	if err != nil {
		return err
	}
	log.Debugf("Checking block interval %d-%d for deposits", watcher.lastCheckedBlock-1, height)
	for _, event := range events {
		quote, exists := watcher.nonDepositedQuotes[event.QuoteHash]
		if exists && event.IsValidForQuote(quote.Data) {
			quote.DepositBlock = event.BlockNumber
			_ = watcher.db.UpdateDepositedPegOutQuote(event.QuoteHash, quote.DepositBlock)
			watcher.depositedQuotes[event.QuoteHash] = quote
			delete(watcher.nonDepositedQuotes, event.QuoteHash)
		}
	}
	return nil
}

func (watcher *DepositEventWatcherImpl) getConfirmedQuotes(height uint64) map[string]*WatchedQuote {
	confirmedQuotes := make(map[string]*WatchedQuote, 0)
	for hash, quote := range watcher.depositedQuotes {
		if uint64(quote.Data.DepositConfirmations)+quote.DepositBlock < height {
			confirmedQuotes[hash] = quote
			delete(watcher.depositedQuotes, hash)
		}
	}
	return confirmedQuotes
}

func (watcher *DepositEventWatcherImpl) cleanExpiredQuotes() {
	now := time.Now()
	for hash, quote := range watcher.nonDepositedQuotes {
		if now.After(quote.Data.GetExpirationTime()) {
			log.Debugf(TimeExpiredError, hash)
			if err := watcher.updateQuoteState(hash, types.RQStateWaitingForDeposit, types.RQStateTimeForDepositElapsed); err == nil {
				delete(watcher.nonDepositedQuotes, hash)
			}
		}
	}
}

func (watcher *DepositEventWatcherImpl) updateQuoteState(hash string, oldState, newState types.RQState) error {
	err := watcher.db.UpdateRetainedPegOutQuoteState(hash, oldState, newState)
	if err != nil {
		log.Errorf(UpdateQuoteStateError, hash, err)
		return err
	}
	return nil
}

func (watcher *DepositEventWatcherImpl) handleDepositedQuotes(quotes map[string]*WatchedQuote) {
	var newState types.RQState
	for hash, quote := range quotes {
		err := watcher.handleDepositedQuote(quote)
		if err == nil {
			newState = types.RQStateCallForUserSucceeded
		} else {
			newState = types.RQStateCallForUserFailed
		}
		if err = watcher.updateQuoteState(hash, types.RQStateWaitingForDepositConfirmations, newState); err != nil {
			log.Errorf("Error updating quote %s: %v", hash, err)
		} else {
			log.Debug("registered pegout quote: ", hash)
			watcher.finalizationCallback(hash, quote, newState)
		}
	}
}

func (watcher *DepositEventWatcherImpl) handleDepositedQuote(quote *WatchedQuote) error {
	paredQuote, err := watcher.rsk.ParsePegOutQuote(quote.Data)
	if err != nil {
		return err
	}

	satoshi, _ := quote.Data.Value.ToSatoshi().Float64()
	watcher.pegoutLocker.Lock()
	defer watcher.pegoutLocker.Unlock()
	err = watcher.btc.LockBtc(satoshi)
	if err != nil {
		return err
	}

	opt := &bind.TransactOpts{
		From:   paredQuote.LpRskAddress,
		Signer: watcher.liquidityProvider.SignTx,
	}

	signatureBytes, err := hex.DecodeString(quote.Signature)
	if err != nil {
		return err
	}
	tx, err := watcher.rsk.RegisterPegOut(opt, paredQuote, signatureBytes)
	if err != nil {
		return err
	}

	if status, err := watcher.rsk.GetTxStatus(context.Background(), tx); err != nil || !status {
		_ = watcher.btc.UnlockBtc(satoshi)
		return err
	}

	err = watcher.btc.UnlockBtc(satoshi)
	if err != nil {
		return err
	}

	_, err = watcher.btc.SendBtc(quote.Data.DepositAddr, uint64(math.Ceil(satoshi)))
	if err != nil {
		return err
	}

	return nil
}

func (watcher *DepositEventWatcherImpl) EndChannel() chan<- bool {
	return watcher.endChannel
}
