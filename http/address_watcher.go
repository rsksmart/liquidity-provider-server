package http

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/connectors"
	"github.com/rsksmart/liquidity-provider/providers"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
	"math/big"
)

type BTCAddressWatcher struct {
	btc             connectors.BTCInterface
	rsk             connectors.RSKInterface
	lp              providers.LiquidityProvider
	calledForUser   bool
	registeredPegIn bool
	quote           *types.Quote
}

func NewBTCAddressWatcher(btc connectors.BTCInterface, rsk connectors.RSKInterface, provider providers.LiquidityProvider, q *types.Quote) (*BTCAddressWatcher, error) {
	watcher := BTCAddressWatcher{
		btc:             btc,
		rsk:             rsk,
		lp:              provider,
		quote:           q,
		calledForUser:   false,
		registeredPegIn: false,
	}
	return &watcher, nil
}

func (w *BTCAddressWatcher) RegisteredPegIn() bool {
	return w.registeredPegIn
}

func (w *BTCAddressWatcher) OnNewConfirmation(txHash string, confirmations int64, amount float64) {
	if !w.calledForUser && confirmations >= int64(w.quote.Confirmations) {
		_, err := w.performCallForUser()
		if err != nil {
			log.Errorf("error calling callForUser. value: %v. error: %v", txHash, err)
			return
		}
		w.calledForUser = true
	}

	if w.calledForUser && confirmations >= w.rsk.GetRequiredBridgeConfirmations() {
		_, err := w.performRegisterPegIn(txHash)
		if err != nil {
			log.Errorf("error calling registerPegIn. value: %v. error: %v", txHash, err)
			return
		}
		w.registeredPegIn = true
	}
}

func (w *BTCAddressWatcher) performCallForUser() (*gethTypes.Transaction, error) {
	q, err := w.rsk.ParseQuote(w.quote)
	opt := w.getTxOptions(q.GasLimit.Uint64()+100000, q.Value, q.LiquidityProviderRskAddress)
	if err != nil {
		return nil, err
	}
	res, err := w.rsk.CallForUser(opt, q)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (w *BTCAddressWatcher) performRegisterPegIn(txHash string) (*gethTypes.Transaction, error) {
	q, err := w.rsk.ParseQuote(w.quote)
	opt := w.getTxOptions(250000, nil, q.LiquidityProviderRskAddress)
	rawTx, err := w.btc.SerializeTx(txHash)
	if err != nil {
		return nil, err
	}
	pmt, err := w.btc.SerializePMT(txHash)
	if err != nil {
		return nil, err
	}
	h, err := w.rsk.HashQuote(w.quote)
	if err != nil {
		log.Errorf("error hashing quote: %v", err)
		return nil, err
	}
	hb, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	signature, err := w.lp.SignHash(hb)
	if err != nil {
		return nil, err
	}
	bh, err := w.btc.GetBlockNumberByTx(txHash)
	if err != nil {
		log.Errorf("error getting block height: %v", err)
		return nil, err
	}
	tx, err := w.rsk.RegisterPegIn(opt, q, signature, rawTx, pmt, big.NewInt(int64(bh)))
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (w *BTCAddressWatcher) getTxOptions(gasLimit uint64, value *big.Int, lpRskAddress common.Address) *bind.TransactOpts {
	opt := &bind.TransactOpts{
		GasLimit: gasLimit,
		Value:    value,
		From:     lpRskAddress,
		Signer:   w.lp.SignTx,
	}
	return opt
}
