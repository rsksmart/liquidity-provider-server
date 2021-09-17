package http

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
	"strings"
  
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
}

const (
	pegInGasLim = 250000
	CFUExtraGas = 100000
)

func NewBTCAddressWatcher(btc connectors.BTCConnector, rsk connectors.RSKConnector, provider providers.LiquidityProvider, q *types.Quote) (*BTCAddressWatcher, error) {
	watcher := BTCAddressWatcher{
		btc:           btc,
		rsk:           rsk,
		lp:            provider,
		quote:         q,
		calledForUser: false,
		done:          make(chan struct{}),
	}
	return &watcher, nil
}

func (w *BTCAddressWatcher) OnNewConfirmation(txHash string, confirmations int64, amount float64) {
	if !w.calledForUser && confirmations >= int64(w.quote.Confirmations) {
		err := w.performCallForUser()
		if err != nil {
			log.Errorf("error calling callForUser. value: %v. error: %v", txHash, err)
			return
		}
	}

	if w.calledForUser && confirmations >= w.rsk.GetRequiredBridgeConfirmations() {
		err := w.performRegisterPegIn(txHash)
		if err != nil {
			log.Errorf("error calling registerPegIn. value: %v. error: %v", txHash, err)
			return
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
		GasLimit: q.GasLimit.Uint64() + CFUExtraGas,
		Value:    q.Value,
		From:     q.LiquidityProviderRskAddress,
		Signer:   w.lp.SignTx,
	}
	_, err = w.rsk.CallForUser(opt, q)
	if err != nil {
		return err
	}
	w.calledForUser = true
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
	h, err := w.rsk.HashQuote(w.quote)
	if err != nil {
		return err
	}
	hb, err := hex.DecodeString(h)
	if err != nil {
		return err
	}
	signature, err := w.lp.SignHash(hb)
	if err != nil {
		return err
	}
	bh, err := w.btc.GetBlockNumberByTx(txHash)
	if err != nil {
		return err
	}
	err = w.rsk.RegisterPegInWithoutTx(q, signature, rawTx, pmt, big.NewInt(bh))
	if err != nil {
		if strings.Contains(err.Error(), "Failed to validate BTC transaction") {
			return nil // allow retrying in case the bridge didn't acknowledge all required confirmations have occurred
		}
	}
	_, err = w.rsk.RegisterPegIn(opt, q, signature, rawTx, pmt, big.NewInt(bh))
	if err != nil {
		return err
	}
	w.registeredPegIn = true

	return nil
}
