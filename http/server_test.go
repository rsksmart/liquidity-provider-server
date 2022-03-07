package http

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil"
	"math/big"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/http/testmocks"
	"github.com/rsksmart/liquidity-provider/providers"
	"github.com/rsksmart/liquidity-provider/types"
	"github.com/stretchr/testify/assert"
	http2 "github.com/stretchr/testify/http"
	"github.com/stretchr/testify/mock"
)

type LiquidityProviderMock struct {
	address string
}

func (lp LiquidityProviderMock) RefundLiquidity(_ []byte) error {
	return nil
}

func (lp LiquidityProviderMock) SignTx(common.Address, *gethTypes.Transaction) (*gethTypes.Transaction, error) {
	return nil, nil
}

func (lp LiquidityProviderMock) Address() string {
	return lp.address
}

func (lp LiquidityProviderMock) GetQuote(q types.Quote, _ uint64, _ uint64) *types.Quote {
	return &q
}

func (lp LiquidityProviderMock) SignQuote(_ []byte, _ string, _ *big.Int) ([]byte, error) {
	return nil, nil
}

func (lp LiquidityProviderMock) SetLiquidity(_ *big.Int) {
}

var providerMocks = []LiquidityProviderMock{
	{address: "123"},
	{address: "0x00d80aA033fb51F191563B08Dc035fA128e942C5"},
}

var testQuotes = []*types.Quote{
	{
		FedBTCAddr:         "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		LBCAddr:            "2ff74F841b95E000625b3A77fed03714874C4fEa",
		LPRSKAddr:          "0x00d80aA033fb51F191563B08Dc035fA128e942C5",
		BTCRefundAddr:      "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		RSKRefundAddr:      "0x5F3b836CA64DA03e613887B46f71D168FC8B5Bdf",
		LPBTCAddr:          "2NDjJznHgtH1rzq63eeFG3SiDi5wxE25FSz",
		CallFee:            250,
		PenaltyFee:         5000,
		ContractAddr:       "0x87136cf829edaF7c46Eb943063369a1C8D4f9085",
		Data:               "",
		GasLimit:           6000000,
		Nonce:              int64(rand.Int()),
		Value:              250,
		AgreementTimestamp: 0,
		TimeForDeposit:     3600,
		CallTime:           3600,
		Confirmations:      10,
	},
}

func testGetProviderByAddress(t *testing.T) {
	var liquidityProviders []providers.LiquidityProvider
	for _, providerMock := range providerMocks {
		liquidityProviders = append(liquidityProviders, providerMock)
	}

	for _, tt := range liquidityProviders {
		result := getProviderByAddress(liquidityProviders, tt.Address())
		assert.EqualValues(t, tt.Address(), result.Address())
	}
}

func testCheckHealth(t *testing.T) {
	rsk := new(testmocks.RskMock)
	btc := new(testmocks.BtcMock)
	db := testmocks.NewDbMock("", testQuotes[0])
	srv := New(rsk, btc, db)

	w := http2.TestResponseWriter{}
	req, err := http.NewRequest("GET", "health", bytes.NewReader([]byte{}))
	if err != nil {
		t.Fatalf("couldn't instantiate request. error: %v", err)
	}
	db.On("CheckConnection").Return(nil).Times(1)
	rsk.On("CheckConnection").Return(nil).Times(1)
	btc.On("CheckConnection").Return(nil).Times(1)
	srv.checkHealthHandler(&w, req)
	db.AssertExpectations(t)
	rsk.AssertExpectations(t)
	btc.AssertExpectations(t)
	assert.EqualValues(t, 200, w.StatusCode)
	assert.EqualValues(t, "application/json", w.Header().Get("Content-Type"))
	assert.EqualValues(t, "{\"status\":\"ok\",\"services\":{\"db\":\"ok\",\"rsk\":\"ok\",\"btc\":\"ok\"}}\n", w.Output)

	w = http2.TestResponseWriter{}
	req, err = http.NewRequest("GET", "health", bytes.NewReader([]byte{}))
	if err != nil {
		t.Fatalf("couldn't instantiate request. error: %v", err)
	}
	db.On("CheckConnection").Return(errors.New("db error")).Times(1)
	rsk.On("CheckConnection").Return(nil).Times(1)
	btc.On("CheckConnection").Return(nil).Times(1)
	srv.checkHealthHandler(&w, req)
	db.AssertExpectations(t)
	rsk.AssertExpectations(t)
	btc.AssertExpectations(t)
	assert.EqualValues(t, 200, w.StatusCode)
	assert.EqualValues(t, "application/json", w.Header().Get("Content-Type"))
	assert.EqualValues(t, "{\"status\":\"degraded\",\"services\":{\"db\":\"unreachable\",\"rsk\":\"ok\",\"btc\":\"ok\"}}\n", w.Output)

	w = http2.TestResponseWriter{}
	req, err = http.NewRequest("GET", "health", bytes.NewReader([]byte{}))
	if err != nil {
		t.Fatalf("couldn't instantiate request. error: %v", err)
	}
	db.On("CheckConnection").Return(errors.New("db error")).Times(1)
	rsk.On("CheckConnection").Return(errors.New("rsk error")).Times(1)
	btc.On("CheckConnection").Return(errors.New("btc error")).Times(1)
	srv.checkHealthHandler(&w, req)
	db.AssertExpectations(t)
	rsk.AssertExpectations(t)
	btc.AssertExpectations(t)
	assert.EqualValues(t, 200, w.StatusCode)
	assert.EqualValues(t, "application/json", w.Header().Get("Content-Type"))
	assert.EqualValues(t, "{\"status\":\"degraded\",\"services\":{\"db\":\"unreachable\",\"rsk\":\"unreachable\",\"btc\":\"unreachable\"}}\n", w.Output)
}

func testGetQuoteComplete(t *testing.T) {
	for _, quote := range testQuotes {
		rsk := new(testmocks.RskMock)
		btc := new(testmocks.BtcMock)
		db := testmocks.NewDbMock("", quote)

		srv := New(rsk, btc, db)

		for _, lp := range providerMocks {
			rsk.On("GetCollateral", lp.address).Return(nil)
			rsk.On("GetAvailableLiquidity", lp.address).Times(1).Return()
			err := srv.AddProvider(lp)
			if err != nil {
				t.Fatalf("couldn't add provider. error: %v", err)
			}
		}
		w := http2.TestResponseWriter{}
		destAddr := "0x63C46fBf3183B0a230833a7076128bdf3D5Bc03F"
		callArgs := ""
		value := quote.Value
		gasLim := 500000
		rskRefAddr := "0x2428E03389e9db669698E0Ffa16FD66DC8156b3c"
		btcRefAddr := "myCqdohiF3cvopyoPMB2rGTrJZx9jJ2ihT"
		body := fmt.Sprintf(
			"{\"callContractAddress\":\"%v\","+
				"\"callContractArguments\":\"%v\","+
				"\"valueToTransfer\":%v,"+
				"\"gaslimit\":%v,"+
				"\"RskRefundAddress\":\"%v\","+
				"\"bitcoinRefundAddress\":\"%v\"}",
			destAddr, callArgs, value, gasLim, rskRefAddr, btcRefAddr)

		tq := types.Quote{
			FedBTCAddr:         "",
			LBCAddr:            "",
			LPRSKAddr:          "",
			BTCRefundAddr:      btcRefAddr,
			RSKRefundAddr:      rskRefAddr,
			LPBTCAddr:          "",
			CallFee:            0,
			PenaltyFee:         0,
			ContractAddr:       destAddr,
			Data:               callArgs,
			GasLimit:           500000,
			Nonce:              0,
			Value:              uint64(value),
			AgreementTimestamp: 0,
			TimeForDeposit:     0,
			CallTime:           0,
			Confirmations:      0,
		}
		req, err := http.NewRequest("POST", "getQuote", bytes.NewReader([]byte(body)))
		if err != nil {
			t.Fatalf("couldn't instantiate request. error: %v", err)
		}
		rsk.On("EstimateGas", destAddr, uint64(value), []byte(callArgs)).Times(1)
		rsk.On("GasPrice").Times(1)
		rsk.On("GetFedAddress").Times(1)
		rsk.On("GetLBCAddress").Times(1)
		rsk.On("GetMinimumLockTxValue").Return(big.NewInt(0), nil).Times(1)
		rsk.On("HashQuote", &tq).Times(len(providerMocks)).Return("", nil)
		db.On("InsertQuote", "", &tq).Times(len(providerMocks)).Return(quote)

		srv.getQuoteHandler(&w, req)
		db.AssertExpectations(t)
		rsk.AssertExpectations(t)
		btc.AssertExpectations(t)
		assert.EqualValues(t, "application/json", w.Header().Get("Content-Type"))

		req, err = http.NewRequest("POST", "getQuote", bytes.NewReader([]byte(body)))
		if err != nil {
			t.Fatalf("couldn't instantiate request. error: %v", err)
		}
		w = http2.TestResponseWriter{}
		rsk.On("EstimateGas", destAddr, uint64(value), []byte(callArgs)).Times(1)
		rsk.On("GasPrice").Times(1)
		rsk.On("GetFedAddress").Times(1)
		rsk.On("GetLBCAddress").Times(1)
		rsk.On("GetMinimumLockTxValue").Return(big.NewInt(int64(quote.Value+quote.CallFee-1)), nil).Times(1)
		rsk.On("HashQuote", &tq).Times(len(providerMocks)).Return("", nil)
		db.On("InsertQuote", "", &tq).Times(len(providerMocks)).Return(quote)
		srv.getQuoteHandler(&w, req)
		assert.EqualValues(t, "bad request; requested amount below bridge's min pegin tx value\n", w.Output)
	}
}

func testAcceptQuoteComplete(t *testing.T) {
	for _, quote := range testQuotes {
		hash := "555c9cfba7638a40a71a17a34fef0c3e192c1fbf4b311ad6e2ae288e97794228"
		rsk := new(testmocks.RskMock)
		btc := new(testmocks.BtcMock)
		db := testmocks.NewDbMock(hash, quote)
		minAmount := btcutil.Amount(quote.Value + quote.CallFee)
		expTime := time.Unix(int64(quote.AgreementTimestamp+quote.TimeForDeposit), 0)

		srv := newServer(rsk, btc, db, func() time.Time {
			return time.Unix(0, 0)
		})
		for _, lp := range providerMocks {
			rsk.On("GetCollateral", lp.address).Times(1).Return(big.NewInt(10), big.NewInt(10))
			rsk.On("GetAvailableLiquidity", lp.address).Times(1).Return()
			err := srv.AddProvider(lp)
			if err != nil {
				t.Errorf("couldn't add provider. error: %v", err)
			}
		}
		w := http2.TestResponseWriter{}
		body := fmt.Sprintf("{\"quoteHash\":\"%v\"}", hash)

		btcRefAddr, lpBTCAddr, lbcAddr, err := decodeAddresses(quote.BTCRefundAddr, quote.LPBTCAddr, quote.LBCAddr)
		if err != nil {
			t.Errorf("couldn't decode addresses. error: %v", err)
		}
		req, err := http.NewRequest("POST", "acceptQuote", bytes.NewReader([]byte(body)))
		if err != nil {
			t.Errorf("couldn't instantiate request. error: %v", err)
		}
		hashBytes, err := hex.DecodeString(hash)
		if err != nil {
			t.Errorf("couldn't decode hash. error: %v", err)
		}

		db.On("GetQuote", hash).Times(1).Return(quote)
		rsk.On("GasPrice").Times(1)
		btc.On("GetDerivedBitcoinAddress", btcRefAddr, lbcAddr, lpBTCAddr, hashBytes).Times(1).Return("")
		btc.On("AddAddressWatcher", "", minAmount, time.Minute, expTime, mock.AnythingOfType("*http.BTCAddressWatcher")).Times(1).Return("")
		srv.acceptQuoteHandler(&w, req)
		db.AssertExpectations(t)
		btc.AssertExpectations(t)
		rsk.AssertExpectations(t)
		assert.EqualValues(t, "application/json", w.Header().Get("Content-Type"))
	}
}

func testInitBtcWatchers(t *testing.T) {
	hash := "555c9cfba7638a40a71a17a34fef0c3e192c1fbf4b311ad6e2ae288e97794228"
	quote := testQuotes[0]
	rsk := new(testmocks.RskMock)
	btc := new(testmocks.BtcMock)
	db := testmocks.NewDbMock(hash, quote)
	minAmount := btcutil.Amount(quote.Value + quote.CallFee)
	expTime := time.Unix(int64(quote.AgreementTimestamp+quote.TimeForDeposit), 0)

	srv := newServer(rsk, btc, db, func() time.Time {
		return time.Unix(0, 0)
	})
	for _, lp := range providerMocks {
		rsk.On("GetCollateral", lp.address).Times(1).Return(big.NewInt(10), big.NewInt(10))
		rsk.On("GetAvailableLiquidity", lp.address).Times(1).Return()
		err := srv.AddProvider(lp)
		if err != nil {
			t.Errorf("couldn't add provider. error: %v", err)
		}
	}

	db.On("GetRetainedQuotes", []types.RQState{types.RQStateWaitingForDeposit, types.RQStateCallForUserSucceeded, types.RQStateRegisterPegInSucceeded}).Times(1).Return([]*types.RetainedQuote{{QuoteHash: hash}})
	db.On("GetQuote", hash).Times(1).Return(quote)
	btc.On("AddAddressWatcher", "", minAmount, time.Minute, expTime, mock.AnythingOfType("*http.BTCAddressWatcher")).Times(1).Return("")
	err := srv.initBtcWatchers()
	if err != nil {
		t.Errorf("couldn't init BTC watchers. error: %v", err)
	}
	db.AssertExpectations(t)
	btc.AssertExpectations(t)
	rsk.AssertExpectations(t)
}

func testGetQuoteExpTime(t *testing.T) {
	quote := types.Quote{AgreementTimestamp: 2, TimeForDeposit: 3}
	expTime := getQuoteExpTime(&quote)
	assert.Equal(t, time.Unix(5, 0), expTime)
}

func TestLiquidityProviderServer(t *testing.T) {
	t.Run("get provider by address", testGetProviderByAddress)
	t.Run("check health", testCheckHealth)
	t.Run("get quote", testGetQuoteComplete)
	t.Run("accept quote", testAcceptQuoteComplete)
	t.Run("init BTC watchers", testInitBtcWatchers)
	t.Run("get quote exp time", testGetQuoteExpTime)
}
