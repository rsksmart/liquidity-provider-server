package http

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/rsksmart/liquidity-provider-server/connectors"

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

func (lp LiquidityProviderMock) SignTx(_ common.Address, _ *gethTypes.Transaction) (*gethTypes.Transaction, error) {
	return nil, nil
}

func (lp LiquidityProviderMock) Address() string {
	return lp.address
}

func (lp LiquidityProviderMock) GetQuote(quote *types.Quote, _ uint64, _ *types.Wei) (*types.Quote, error) {
	res := *quote
	res.CallFee = types.NewWei(0)
	res.PenaltyFee = types.NewWei(0)
	return &res, nil
}

func (lp LiquidityProviderMock) SignQuote(_ []byte, _ string, _ *types.Wei) ([]byte, error) {
	return nil, nil
}

var providerMocks = []LiquidityProviderMock{
	{address: "123"},
	{address: "0x00d80aA033fb51F191563B08Dc035fA128e942C5"},
}

var cfgData = ConfigData{
	MaxQuoteValue: 600000000000000000,
	RSK: LiquidityProviderList{
		Endpoint:                    "",
		LBCAddr:                     "",
		BridgeAddr:                  "",
		RequiredBridgeConfirmations: 10,
		MaxQuoteValue:               600000000000000000,
	},
}

var testQuotes = []*types.Quote{
	{
		FedBTCAddr:         "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		LBCAddr:            "2ff74F841b95E000625b3A77fed03714874C4fEa",
		LPRSKAddr:          "0x00d80aA033fb51F191563B08Dc035fA128e942C5",
		BTCRefundAddr:      "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		RSKRefundAddr:      "0x5F3b836CA64DA03e613887B46f71D168FC8B5Bdf",
		LPBTCAddr:          "2NDjJznHgtH1rzq63eeFG3SiDi5wxE25FSz",
		CallFee:            types.NewWei(250),
		PenaltyFee:         types.NewWei(5000),
		ContractAddr:       "0x87136cf829edaF7c46Eb943063369a1C8D4f9085",
		Data:               "",
		GasLimit:           6000000,
		Nonce:              int64(rand.Int()),
		Value:              types.NewWei(250),
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

func testGetProviderByAddressWhenNotFoundShouldReturnNull(t *testing.T) {
	var liquidityProviders []providers.LiquidityProvider
	for _, providerMock := range providerMocks {
		liquidityProviders = append(liquidityProviders, providerMock)
	}

	var nonLiquidityProviderAddress = "0xa554d96413FF72E93437C4072438302C38350EE3"
	result := getProviderByAddress(liquidityProviders, nonLiquidityProviderAddress)
	assert.Empty(t, result)
}

func testCheckHealth(t *testing.T) {
	rsk := new(testmocks.RskMock)
	btc := new(testmocks.BtcMock)
	db := testmocks.NewDbMock("", testQuotes[0])

	srv := New(rsk, btc, db, cfgData)

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

		srv := New(rsk, btc, db, cfgData)

		for _, lp := range providerMocks {
			rsk.On("GetCollateral", lp.address).Return(nil)
			err := srv.AddProvider(lp)
			if err != nil {
				t.Fatalf("couldn't add provider. error: %v", err)
			}
		}
		w := http2.TestResponseWriter{}
		destAddr := "0x63C46fBf3183B0a230833a7076128bdf3D5Bc03F"
		callArgs := ""
		value := quote.Value
		rskRefAddr := "0x2428E03389e9db669698E0Ffa16FD66DC8156b3c"
		btcRefAddr := "myCqdohiF3cvopyoPMB2rGTrJZx9jJ2ihT"
		body := fmt.Sprintf(
			"{\"callContractAddress\":\"%v\","+
				"\"callContractArguments\":\"%v\","+
				"\"valueToTransfer\":%v,"+
				"\"RskRefundAddress\":\"%v\","+
				"\"bitcoinRefundAddress\":\"%v\"}",
			destAddr, callArgs, value, rskRefAddr, btcRefAddr)

		tq := types.Quote{
			FedBTCAddr:         "",
			LBCAddr:            "",
			LPRSKAddr:          "",
			BTCRefundAddr:      btcRefAddr,
			RSKRefundAddr:      rskRefAddr,
			LPBTCAddr:          "",
			CallFee:            types.NewWei(0),
			PenaltyFee:         types.NewWei(0),
			ContractAddr:       destAddr,
			Data:               callArgs,
			GasLimit:           10000,
			Nonce:              0,
			Value:              value.Copy(),
			AgreementTimestamp: 0,
			TimeForDeposit:     0,
			CallTime:           0,
			Confirmations:      0,
		}
		req, err := http.NewRequest("POST", "getQuote", bytes.NewReader([]byte(body)))
		if err != nil {
			t.Fatalf("couldn't instantiate request. error: %v", err)
		}
		rsk.On("EstimateGas", destAddr, value.AsBigInt(), []byte(callArgs)).Times(1)
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
		rsk.On("EstimateGas", destAddr, value.AsBigInt(), []byte(callArgs)).Times(1)
		rsk.On("GasPrice").Times(1)
		rsk.On("GetFedAddress").Times(1)
		rsk.On("GetLBCAddress").Times(1)
		rsk.On("GetMinimumLockTxValue").Return(new(big.Int).Add(big.NewInt(-1), new(big.Int).Add(quote.Value.AsBigInt(), quote.CallFee.AsBigInt())), nil).Times(1)
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
		sat, _ := new(types.Wei).Add(quote.Value, quote.CallFee).ToSatoshi().Float64()
		minAmount := btcutil.Amount(uint64(math.Ceil(sat)))
		expTime := time.Unix(int64(quote.AgreementTimestamp+quote.TimeForDeposit), 0)
		fedInfo := &connectors.FedInfo{}

		srv := newServer(rsk, btc, db, func() time.Time {
			return time.Unix(0, 0)
		}, cfgData)
		for _, lp := range providerMocks {
			rsk.On("GetCollateral", lp.address).Times(1).Return(big.NewInt(10), big.NewInt(10))
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

		db.On("GetQuote", hash).Times(1).Return(quote, nil)
		db.On("GetRetainedQuote", hash).Times(1).Return(nil, nil)
		rsk.On("GasPrice").Times(1)
		rsk.On("FetchFederationInfo").Times(1).Return(fedInfo, nil)
		btc.On("GetDerivedBitcoinAddress", fedInfo, btcRefAddr, lbcAddr, lpBTCAddr, hashBytes).Times(1).Return("")
		btc.On("AddAddressWatcher", "", minAmount, time.Minute, expTime, mock.AnythingOfType("*http.BTCAddressWatcher"), mock.AnythingOfType("func(connectors.AddressWatcher)")).Times(1).Return("")
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
	sat, _ := new(types.Wei).Add(quote.Value, quote.CallFee).ToSatoshi().Float64()
	minAmount := btcutil.Amount(uint64(math.Ceil(sat)))
	expTime := time.Unix(int64(quote.AgreementTimestamp+quote.TimeForDeposit), 0)

	srv := newServer(rsk, btc, db, func() time.Time {
		return time.Unix(0, 0)
	}, cfgData)
	for _, lp := range providerMocks {
		rsk.On("GetCollateral", lp.address).Times(1).Return(big.NewInt(10), big.NewInt(10))
		err := srv.AddProvider(lp)
		if err != nil {
			t.Errorf("couldn't add provider. error: %v", err)
		}
	}

	db.On("GetRetainedQuotes", []types.RQState{types.RQStateWaitingForDeposit, types.RQStateCallForUserSucceeded}).Times(1).Return([]*types.RetainedQuote{{QuoteHash: hash}})
	db.On("GetQuote", hash).Times(1).Return(quote)
	btc.On("AddAddressWatcher", "", minAmount, time.Minute, expTime, mock.AnythingOfType("*http.BTCAddressWatcher"), mock.AnythingOfType("func(connectors.AddressWatcher)")).Times(1).Return("")
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

func testDecodeAddress(t *testing.T) {
	_, _, _, err := decodeAddresses("1PRTTaJesdNovgne6Ehcdu1fpEdX7913CK", "1JRRmhqTc87SmLjSHaiJjHyuJfDUc8AQDF", "0xa554d96413FF72E93437C4072438302C38350EE3")
	assert.Empty(t, err)
}

func testDecodeAddressWithAnInvalidBtcRefundAddr(t *testing.T) {
	_, _, _, err := decodeAddresses("0xa554d96413FF72E93437C4072438302C38350EE3", "1JRRmhqTc87SmLjSHaiJjHyuJfDUc8AQDF", "0xa554d96413FF72E93437C4072438302C38350EE3")
	assert.Equal(t, "the provider address is not a valid base58 encoded address. address: 0xa554d96413FF72E93437C4072438302C38350EE3", err.Error())
}

func testDecodeAddressWithAnInvalidLpBTCAddrB(t *testing.T) {
	_, _, _, err := decodeAddresses("1PRTTaJesdNovgne6Ehcdu1fpEdX7913CK", "0xa554d96413FF72E93437C4072438302C38350EE3", "0xa554d96413FF72E93437C4072438302C38350EE3")
	assert.Equal(t, "the provider address is not a valid base58 encoded address. address: 0xa554d96413FF72E93437C4072438302C38350EE3", err.Error())
}

func testDecodeAddressWithAnInvalidLbcAddrB(t *testing.T) {
	_, _, _, err := decodeAddresses("1PRTTaJesdNovgne6Ehcdu1fpEdX7913CK", "1JRRmhqTc87SmLjSHaiJjHyuJfDUc8AQDF", "1JRRmhqTc87SmLjSHaiJjHyuJfDUc8AQDF")
	assert.Equal(t, "invalid address: 1JRRmhqTc87SmLjSHaiJjHyuJfDUc8AQDF", err.Error())
}

func testInvalidQuoteValue(t *testing.T) {
	for _, quote := range testQuotes {
		rsk := new(testmocks.RskMock)
		btc := new(testmocks.BtcMock)
		db := testmocks.NewDbMock("", quote)

		srv := New(rsk, btc, db, cfgData)

		for _, lp := range providerMocks {
			rsk.On("GetCollateral", lp.address).Return(nil)
			err := srv.AddProvider(lp)
			if err != nil {
				t.Fatalf("couldn't add provider. error: %v", err)
			}
		}
		w := http2.TestResponseWriter{}
		destAddr := "0x63C46fBf3183B0a230833a7076128bdf3D5Bc03F"
		callArgs := ""
		rskRefAddr := "0x2428E03389e9db669698E0Ffa16FD66DC8156b3c"
		btcRefAddr := "myCqdohiF3cvopyoPMB2rGTrJZx9jJ2ihT"
		body := fmt.Sprintf(
			"{\"callContractAddress\":\"%v\","+
				"\"callContractArguments\":\"%v\","+
				"\"valueToTransfer\":%v,"+
				"\"RskRefundAddress\":\"%v\","+
				"\"bitcoinRefundAddress\":\"%v\"}",
			destAddr, callArgs, 600000000000000001, rskRefAddr, btcRefAddr)

		req, err := http.NewRequest("POST", "getQuote", bytes.NewReader([]byte(body)))
		if err != nil {
			t.Log("Tewst")
			t.Fatalf("couldn't instantiate request. error: %v", err)
		}

		srv.getQuoteHandler(&w, req)
		assert.EqualValues(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
		assert.EqualValues(t, "internal server error\n", w.Output)
	}
}

func testGetProviders(t *testing.T) {
	rsk := new(testmocks.RskMock)
	btc := new(testmocks.BtcMock)
	db := testmocks.NewDbMock("", nil)

	srv := New(rsk, btc, db, cfgData)
	req, err := http.NewRequest("GET", "getProviders", bytes.NewReader([]byte("")))
	w := http2.TestResponseWriter{}

	if err != nil {
		t.Fatalf("couldn't instantiate request. error: %v", err)
	}

	rsk.On("GetProviders").Return(nil)
	srv.getProvidersHandler(&w, req)

	assert.EqualValues(t, "application/json", w.Header().Get("Content-Type"))
	assert.EqualValues(t, "null\n", w.Output)
}

func TestLiquidityProviderServer(t *testing.T) {
	t.Run("get provider by address", testGetProviderByAddress)
	t.Run("check health", testCheckHealth)
	t.Run("get provider should return null when provider not found", testGetProviderByAddressWhenNotFoundShouldReturnNull)
	t.Run("get quote", testGetQuoteComplete)
	t.Run("get quote invalid quote value", testInvalidQuoteValue)
	t.Run("accept quote", testAcceptQuoteComplete)
	t.Run("init BTC watchers", testInitBtcWatchers)
	t.Run("get quote exp time", testGetQuoteExpTime)
	t.Run("decode address", testDecodeAddress)
	t.Run("decode address with an invalid btcRefundAddr", testDecodeAddressWithAnInvalidBtcRefundAddr)
	t.Run("decode address with an invalid lpBTCAddrB", testDecodeAddressWithAnInvalidLpBTCAddrB)
	t.Run("decode address with an invalid lbcAddrB", testDecodeAddressWithAnInvalidLbcAddrB)
	t.Run("get registered providers", testGetProviders)
}
