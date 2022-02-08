package http

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rsksmart/liquidity-provider-server/storage"
	"math/big"
	"strings"
	"time"

	"github.com/rsksmart/liquidity-provider-server/connectors"
	"github.com/rsksmart/liquidity-provider/providers"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
)

type BTCAddressWatcher struct {
	hash          string
	btc           connectors.BTCConnector
	rsk           connectors.RSKConnector
	lp            providers.LiquidityProvider
	db            storage.DBConnector
	calledForUser bool
	quote         *types.Quote
	done          chan struct{}
	closed        bool
	signature     []byte
}

const (
	pegInGasLim = 250000
	CFUExtraGas = 150000
)

func NewBTCAddressWatcher(hash string,
	btc connectors.BTCConnector, rsk connectors.RSKConnector, provider providers.LiquidityProvider, db storage.DBConnector,
	q *types.Quote, signature []byte, calledForUser bool) *BTCAddressWatcher {
	watcher := BTCAddressWatcher{
		hash:          hash,
		btc:           btc,
		rsk:           rsk,
		lp:            provider,
		db:            db,
		quote:         q,
		calledForUser: calledForUser,
		signature:     signature,
		done:          make(chan struct{}),
	}
	return &watcher
}

func (w *BTCAddressWatcher) OnNewConfirmation(txHash string, confirmations int64, amount btcutil.Amount) {
	if w.closed {
		log.Errorf("watcher is closed; cannot handle OnNewConfirmation; hash: %v", w.hash)
		return
	}
	log.Debugf("processing OnNewConfirmation event for tx %v; confirmations: %v; received amount: %v", txHash, confirmations, amount)

	if !w.calledForUser && confirmations >= int64(w.quote.Confirmations) {
		err := w.performCallForUser()
		if err != nil {
			log.Errorf("error calling callForUser. value: %v. error: %v", txHash, err)
			err = w.db.DeleteRetainedQuote(w.hash)
			if err != nil {
				log.Errorf("error deleting retained quote from db. hash: %v", w.hash)
			}

			w.close()
			return
		}

		err = w.db.SetRetainedQuoteCalledForUserFlag(w.hash)
		if err != nil {
			log.Errorf("error setting CalledForUser flag for retained quote. hash: %v", w.hash)
			err = w.db.DeleteRetainedQuote(w.hash)
			if err != nil {
				log.Errorf("error deleting retained quote from db. hash: %v", w.hash)
			}

			w.close()
			return
		}
		w.calledForUser = true
		log.Debugf("registered callforuser for tx %v", txHash)
	}

	if w.calledForUser && confirmations >= w.rsk.GetRequiredBridgeConfirmations() {
		err := w.performRegisterPegIn(txHash)
		if err != nil {
			log.Errorf("error calling registerPegIn. value: %v. error: %v", txHash, err)
		}
		log.Debugf("performed RegisterPegIn for tx %v", txHash)
	}
}

func (w *BTCAddressWatcher) OnExpire() {
	if w.closed {
		log.Errorf("watcher is closed; cannot handle OnExpire; hash: %v", w.hash)
		return
	}
	log.Debugf("time has expired for quote: %v", w.hash)
	err := w.db.DeleteRetainedQuote(w.hash)
	if err != nil {
		log.Errorf("error deleting retained quote from db. hash: %v", w.hash)
	}
	w.close()
}

func (w *BTCAddressWatcher) Done() <-chan struct{} {
	return w.done
}

func (w *BTCAddressWatcher) performCallForUser() error {
	q, err := w.rsk.ParseQuote(w.quote)
	if err != nil {
		return err
	}
	opt := &bind.TransactOpts{
		GasLimit: uint64(q.GasLimit + CFUExtraGas),
		Value:    big.NewInt(int64(q.Value)),
		From:     q.LiquidityProviderRskAddress,
		Signer:   w.lp.SignTx,
	}
	tx, err := w.rsk.CallForUser(opt, q)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*8760) // timeout is a year
	defer cancel()
	s, err := w.rsk.GetTxStatus(ctx, tx)
	if err != nil || !s {
		return fmt.Errorf("transaction failed. hash: %v", tx.Hash())
	}
	return nil
}

func (w *BTCAddressWatcher) performRegisterPegIn(txHash string) error {
	q, err := w.rsk.ParseQuote(w.quote)
	if err != nil {
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
		return err
	}
	pmt, err := w.btc.SerializePMT(txHash)
	if err != nil {
		return err
	}
	bh, err := w.btc.GetBlockNumberByTx(txHash)
	if err != nil {
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
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*8760) // timeout is a year
	defer cancel()
	s, err := w.rsk.GetTxStatus(ctx, tx)
	if err != nil || !s {
		return fmt.Errorf("transaction failed. hash: %v", tx.Hash())
	}
	log.Debugf("registered pegin for tx %v", txHash)
	err = w.notifyProvider()
	if err != nil {
		log.Errorf("error refunding provider. value: %v. error: %v", txHash, err)
	}
	w.close()
	return nil
}

func (w *BTCAddressWatcher) notifyProvider() error {
	h, err := w.rsk.HashQuote(w.quote)
	if err != nil {
		return err
	}
	hb, err := hex.DecodeString(h)
	if err != nil {
		return err
	}
	err = w.lp.RefundLiquidity(hb)
	if err != nil {
		return fmt.Errorf("failed to refund to liquidity provider. quote: %v. error: %v", h, err)
	}
	return nil
}

func (w *BTCAddressWatcher) close() {
	w.closed = true
	close(w.done)
}
