package connectors

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider/providers"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
	"math/big"
)

type AddressWatcher interface {
	OnNewConfirmation(txHash string, confirmations int64, amount float64)
}

type BTCAddressWatcher struct {
	btc                    BTCInterface
	rsk                    RSKInterface
	lp                     providers.LiquidityProvider
	calledForUser          bool
	quote                  *types.Quote
	reqBridgeConfirmations int64
}

func NewBTCAddressWatcher(btc BTCInterface, rsk RSKInterface, provider providers.LiquidityProvider, q *types.Quote, reqBridgeConfirms int64) (*BTCAddressWatcher, error) {
	watcher := BTCAddressWatcher{
		btc:                    btc,
		rsk:                    rsk,
		lp:                     provider,
		calledForUser:          false,
		quote:                  q,
		reqBridgeConfirmations: reqBridgeConfirms,
	}
	return &watcher, nil
}

func (w BTCAddressWatcher) OnNewConfirmation(txHash string, confirmations int64, amount float64) {
	quoteConfirmations := new(big.Int).SetUint64(uint64(w.quote.Confirmations))
	if !w.calledForUser && confirmations >= quoteConfirmations.Int64() {
		_, err := w.performCallForUser()
		if err != nil {
			log.Errorf("error calling callForUser. value: %v. error: %v", txHash, err)
			return
		}
		w.calledForUser = true
	}

	if w.calledForUser && confirmations <= w.reqBridgeConfirmations {
		_, err := w.performRegisterPegIn(txHash)
		if err != nil {
			log.Errorf("error calling registerPegIn. value: %v. error: %v", txHash, err)
			return
		}
	}
}

func getChainIdFromNetwork(params chaincfg.Params) *big.Int {
	switch params.Name {
	case "mainnet":
		return big.NewInt(30)
	default:
		return big.NewInt(31)

	}
}

func (w BTCAddressWatcher) performCallForUser() (*gethTypes.Transaction, error) {
	q, err := w.rsk.ParseQuote(w.quote)
	opt := w.getTxOptions(q.GasLimit.Uint64(), q.Value, q.LiquidityProviderRskAddress)
	if err != nil {
		return nil, err
	}
	res, err := w.rsk.CallForUser(opt, q)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (w BTCAddressWatcher) performRegisterPegIn(txHash string) (*gethTypes.Transaction, error) {
	q, err := w.rsk.ParseQuote(w.quote)
	opt := w.getTxOptions(q.GasLimit.Uint64(), q.Value, q.LiquidityProviderRskAddress)
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
	bh, err := w.rsk.GetBlockHeight()
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

func (w BTCAddressWatcher) getTxOptions(gasLimit uint64, value *big.Int, lpRskAddress common.Address) *bind.TransactOpts {
	tx := gethTypes.NewTx(nil) // TODO: empty transaction is ok at this point?
	chainId := getChainIdFromNetwork(w.btc.GetParams())
	// add 10% to the gas limit, so it's enough for everything
	limit := gasLimit + (10 * gasLimit / 100)

	signer := func(common.Address, *gethTypes.Transaction) (*gethTypes.Transaction, error) {
		return w.lp.SignTx(tx, chainId)
	}
	opt := &bind.TransactOpts{
		GasLimit: limit,
		Value:    value,
		From:     lpRskAddress,
		Signer:   signer, // TODO: get signer fn. and assign to signer
	}
	return opt
}
