package http

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
	"strings"
	"time"

	"github.com/rsksmart/liquidity-provider-server/connectors"
	"github.com/rsksmart/liquidity-provider/providers"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
)

type BTCAddressWatcher struct {
	btc           connectors.BTCConnector
	rsk           connectors.RSKConnector
	lp            providers.LiquidityProvider
	calledForUser bool
	quote         *types.Quote
	done          chan struct{}
	signature     []byte
}

const (
	pegInGasLim = 250000
	CFUExtraGas = 150000
)

func NewBTCAddressWatcher(btc connectors.BTCConnector, rsk connectors.RSKConnector, provider providers.LiquidityProvider, q *types.Quote, signature []byte) (*BTCAddressWatcher, error) {
	watcher := BTCAddressWatcher{
		btc:           btc,
		rsk:           rsk,
		lp:            provider,
		quote:         q,
		calledForUser: false,
		signature:     signature,
		done:          make(chan struct{}),
	}
	return &watcher, nil
}

func (w *BTCAddressWatcher) OnNewConfirmation(txHash string, confirmations int64, _ float64) {
	if !w.calledForUser && confirmations >= int64(w.quote.Confirmations) {
		err := w.performCallForUser()
		if err != nil {
			log.Errorf("error calling callForUser. value: %v. error: %v", txHash, err)
			close(w.done)
			return
		}
		w.calledForUser = true
	}

	if w.calledForUser && confirmations >= w.rsk.GetRequiredBridgeConfirmations() {
		err := w.performRegisterPegIn(txHash)
		if err != nil {
			log.Errorf("error calling registerPegIn. value: %v. error: %v", txHash, err)
		}
		close(w.done)
	}
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
			return nil // allow retrying in case the bridge didn't acknowledge all required confirmations have occurred
		}
	}
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
