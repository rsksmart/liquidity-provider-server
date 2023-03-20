package http

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/connectors/bindings"
	"github.com/rsksmart/liquidity-provider-server/storage"

	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider-server/pegout"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/rsksmart/liquidity-provider-server/connectors"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/http/testmocks"
	"github.com/rsksmart/liquidity-provider/types"
	"github.com/stretchr/testify/assert"
	http2 "github.com/stretchr/testify/http"
	"github.com/stretchr/testify/mock"
)

type basicTestCase struct {
	caseName   string
	request    string
	assertions func(response *http.Response)
}

type LiquidityProviderMock struct {
	address string
}

type LiquidityPegOutProviderMock struct {
	address string
}

func (lp LiquidityPegOutProviderMock) GetQuote(quote *pegout.Quote, test uint64, gas uint64, gasPrice *big.Int) (*pegout.Quote, error) {
	return quote, nil
}

func (lp LiquidityPegOutProviderMock) SignTx(address common.Address, transaction *gethTypes.Transaction) (*gethTypes.Transaction, error) {
	return &gethTypes.Transaction{}, nil
}

func (lp LiquidityProviderMock) SignTx(_ common.Address, _ *gethTypes.Transaction) (*gethTypes.Transaction, error) {
	return nil, nil
}

func (lp LiquidityProviderMock) Address() string {
	return lp.address
}

func (lp LiquidityProviderMock) GetQuote(quote *pegin.Quote, _ uint64, _ *types.Wei) (*pegin.Quote, error) {
	res := *quote
	res.CallFee = types.NewWei(0)
	res.PenaltyFee = types.NewWei(0)
	return &res, nil
}

func (lp LiquidityProviderMock) SignQuote(_ []byte, _ string, _ *types.Wei) ([]byte, error) {
	return []byte("fb4a3e40390dee7db6e861e10e5e3b39a0cf546eeccc8c0902249419140d9f29335023e3a83deee747f4987e9cd32773d2afa5176295dc2042255b57a30300201c"), nil
}

func (lp LiquidityPegOutProviderMock) SignQuote(hash []byte, depositAddr string, satoshis uint64) ([]byte, error) {
	return hex.DecodeString("fb4a3e40390dee7db6e861e10e5e3b39a0cf546eeccc8c0902249419140d9f29335023e3a83deee747f4987e9cd32773d2afa5176295dc2042255b57a30300201c")
}

func (lp LiquidityPegOutProviderMock) Address() string {
	return lp.address
}

var providerMocks = []LiquidityProviderMock{
	{address: "123"},
	{address: "0x00d80aA033fb51F191563B08Dc035fA128e942C5"},
}

var providerPegOutMocks = []LiquidityPegOutProviderMock{
	{address: "456"},
	{address: "0xa554d96413FF72E93437C4072438302C38350EE3"},
}

var providerCfgData = pegin.ProviderConfig{}

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

var testQuotes = []*pegin.Quote{
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
		LpCallTime:         3600,
		Confirmations:      10,
	},
}

var testPegOutQuotes = []*pegout.Quote{
	{
		LBCAddr:               "2ff74F841b95E000625b3A77fed03714874C4fEa",
		LPRSKAddr:             "0xa554d96413FF72E93437C4072438302C38350EE3",
		RSKRefundAddr:         "0x5F3b836CA64DA03e613887B46f71D168FC8B5Bdf",
		CallFee:               250,
		PenaltyFee:            5000,
		Nonce:                 int64(rand.Int()),
		Value:                 250,
		AgreementTimestamp:    0,
		DepositDateLimit:      0,
		DepositConfirmations:  0,
		TransferConfirmations: 0,
		TransferTime:          0,
		ExpireDate:            0,
		ExpireBlocks:          0,
	},
}

func testGetProviderByAddress(t *testing.T) {
	var liquidityProviders []pegin.LiquidityProvider
	for _, providerMock := range providerMocks {
		liquidityProviders = append(liquidityProviders, providerMock)
	}

	for _, tt := range liquidityProviders {
		result := pegin.GetPeginProviderByAddress(liquidityProviders, tt.Address())
		assert.EqualValues(t, tt.Address(), result.Address())
	}
}

func testGetProviderByAddressWhenNotFoundShouldReturnNull(t *testing.T) {
	var liquidityProviders []pegin.LiquidityProvider
	for _, providerMock := range providerMocks {
		liquidityProviders = append(liquidityProviders, providerMock)
	}

	var nonLiquidityProviderAddress = "0xa554d96413FF72E93437C4072438302C38350EE3"
	result := pegin.GetPeginProviderByAddress(liquidityProviders, nonLiquidityProviderAddress)
	assert.Empty(t, result)
}

func testCheckHealth(t *testing.T) {
	rsk := new(testmocks.RskMock)
	btc := new(testmocks.BtcMock)
	lp := new(storage.LPRepository)
	mongoDb, _ := testmocks.NewDbMock("", testQuotes[0], nil)

	srv := New(rsk, btc, mongoDb, cfgData, lp, providerCfgData)

	w := http2.TestResponseWriter{}
	req, err := http.NewRequest("GET", "health", bytes.NewReader([]byte{}))
	if err != nil {
		t.Fatalf("couldn't instantiate request. error: %v", err)
	}
	mongoDb.On("CheckConnection").Return(nil).Times(1)
	rsk.On("CheckConnection").Return(nil).Times(1)
	btc.On("CheckConnection").Return(nil).Times(1)
	srv.checkHealthHandler(&w, req)
	mongoDb.AssertExpectations(t)
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
	mongoDb.On("CheckConnection").Return(errors.New("db error")).Times(1)
	rsk.On("CheckConnection").Return(nil).Times(1)
	btc.On("CheckConnection").Return(nil).Times(1)
	srv.checkHealthHandler(&w, req)
	mongoDb.AssertExpectations(t)
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
	mongoDb.On("CheckConnection").Return(errors.New("db error")).Times(1)
	rsk.On("CheckConnection").Return(errors.New("rsk error")).Times(1)
	btc.On("CheckConnection").Return(errors.New("btc error")).Times(1)
	srv.checkHealthHandler(&w, req)
	mongoDb.AssertExpectations(t)
	rsk.AssertExpectations(t)
	btc.AssertExpectations(t)
	assert.EqualValues(t, 200, w.StatusCode)
	assert.EqualValues(t, "application/json", w.Header().Get("Content-Type"))
	assert.EqualValues(t, "{\"status\":\"degraded\",\"services\":{\"db\":\"unreachable\",\"rsk\":\"unreachable\",\"btc\":\"unreachable\"}}\n", w.Output)
}

func testGetQuoteComplete(t *testing.T) {
	quote := testQuotes[0]
	callContractArgumentsField := `"callContractArguments":"%v",`
	basicQuoteFields := `"callEoaOrContractAddress":"%v","valueToTransfer":%v,` +
		`"rskRefundAddress":"%v","lpAddress":"%v", "bitcoinRefundAddress":"%v"`

	rsk := new(testmocks.RskMock)
	btc := new(testmocks.BtcMock)
	lpRepo := new(storage.LPRepository)
	mongoDb, _ := testmocks.NewDbMock("", quote, nil)

	srv := New(rsk, btc, mongoDb, cfgData, lpRepo, providerCfgData)

	detailMock := types.ProviderRegisterRequest{}
	for _, lp := range providerMocks {
		rsk.On("GetCollateral", lp.address).Return(big.NewInt(10), big.NewInt(10), nil)
		rsk.On("RegisterProvider", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
		mongoDb.On("InsertProvider", mock.Anything, mock.Anything).Return(nil)
		err := srv.AddProvider(lp, detailMock)
		if err != nil {
			t.Fatalf("couldn't add provider. error: %v", err)
		}
	}

	destAddrEOA := "0x63C46fBf3183B0a230833a7076128bdf3D5Bc03F"
	destAddrSC := "0x63C46fBf3183B0a230833a7076128bdf3D5Bc04F"
	callArgs := "0x"
	value := quote.Value
	rskRefAddrEOA := "0x9d93929a9099be4355fc2389fbf253982f9df47c"
	rskRefAddrSC := "0x1eD614cd3443EFd9c70F04b6d777aed947A4b0c4"
	btcRefAddr := "myCqdohiF3cvopyoPMB2rGTrJZx9jJ2ihT"

	eoaQuote := pegin.Quote{
		FedBTCAddr:         "",
		LBCAddr:            "",
		LPRSKAddr:          "",
		BTCRefundAddr:      btcRefAddr,
		RSKRefundAddr:      rskRefAddrEOA,
		LPBTCAddr:          "",
		CallFee:            types.NewWei(0),
		PenaltyFee:         types.NewWei(0),
		Nonce:              0,
		Value:              value.Copy(),
		AgreementTimestamp: 0,
		TimeForDeposit:     0,
		LpCallTime:         0,
		Confirmations:      0,
		GasLimit:           10000,
		ContractAddr:       rskRefAddrEOA,
	}

	scQuoteCallContractEoa := eoaQuote
	scQuoteCallContractEoa.RSKRefundAddr = rskRefAddrSC
	scQuoteCallContractEoa.ContractAddr = destAddrEOA

	scQuoteCallContractSc := eoaQuote
	scQuoteCallContractSc.RSKRefundAddr = rskRefAddrSC
	scQuoteCallContractSc.ContractAddr = destAddrSC
	scQuoteCallContractSc.Data = callArgs

	defaultMocks := func(rskMock *testmocks.RskMock, btcMock *testmocks.BtcMock, dbMock *testmocks.DbMock) {
		rskMock.On("EstimateGas", mock.Anything, value.AsBigInt(), mock.Anything).Times(1)
		rskMock.On("GasPrice").Times(1)
		rskMock.On("GetFedAddress").Times(1)
		rskMock.On("GetLBCAddress").Times(1)
		rskMock.On("IsEOA", "").Return(false, errors.New("invalid address"))
		rskMock.On("IsEOA", rskRefAddrEOA).Return(true, nil)
		rskMock.On("IsEOA", destAddrEOA).Return(true, nil)
		rskMock.On("IsEOA", destAddrSC).Return(false, nil)
		rskMock.On("IsEOA", rskRefAddrSC).Return(false, nil)
		rskMock.On("GetMinimumLockTxValue").Return(big.NewInt(0), nil).Times(1)
		rskMock.On("HashQuote", mock.Anything).Times(len(providerMocks)).Return("", nil)
		rskMock.On("HashQuote", mock.Anything).Times(len(providerMocks)).Return("", nil)
		dbMock.On("InsertQuote", "", mock.Anything).Times(len(providerMocks)).Return(quote)
	}

	testCases := []*struct {
		basicTestCase
		customMocks func(rskMock *testmocks.RskMock, btcMock *testmocks.BtcMock, dbMock *testmocks.DbMock)
	}{
		{
			basicTestCase: basicTestCase{
				caseName: "Return error when requested amount below bridge's min pegin tx value",
				request: fmt.Sprintf("{"+basicQuoteFields+"}",
					destAddrEOA, value, rskRefAddrEOA, rskRefAddrEOA, btcRefAddr,
				),
				assertions: func(res *http.Response) {
					response := &ErrorBody{}
					json.NewDecoder(res.Body).Decode(response)
					assert.EqualValues(t, "application/json", res.Header.Get("Content-Type"))
					assert.EqualValues(t, 400, res.StatusCode)
					assert.EqualValues(t, "requested amount below bridge's min pegin tx value", response.Message)
				},
			},
			customMocks: func(rskMock *testmocks.RskMock, btcMock *testmocks.BtcMock, dbMock *testmocks.DbMock) {
				rsk.On("GetMinimumLockTxValue").Return(new(big.Int).Add(big.NewInt(-1), new(big.Int).Add(quote.Value.AsBigInt(), quote.CallFee.AsBigInt())), nil).Times(1)
			},
		},
		{
			basicTestCase: basicTestCase{
				caseName: "Return quote successfully for EOA origin",
				request: fmt.Sprintf("{"+basicQuoteFields+"}",
					rskRefAddrEOA, value, rskRefAddrEOA, rskRefAddrEOA, btcRefAddr,
				),
				assertions: func(res *http.Response) {
					var response []*QuoteReturn
					json.NewDecoder(res.Body).Decode(&response)
					assert.EqualValues(t, "application/json", res.Header.Get("Content-Type"))
					assert.EqualValues(t, 200, res.StatusCode)
					assert.EqualValues(t, eoaQuote, *response[0].Quote)
				},
			},
		},
		{
			basicTestCase: basicTestCase{
				caseName: "Return quote successfully for SC origin and SC call contract address",
				request: fmt.Sprintf("{"+callContractArgumentsField+basicQuoteFields+"}",
					callArgs, destAddrSC, value, rskRefAddrSC, rskRefAddrSC, btcRefAddr,
				),
				assertions: func(res *http.Response) {
					var response []*QuoteReturn
					json.NewDecoder(res.Body).Decode(&response)
					assert.EqualValues(t, "application/json", res.Header.Get("Content-Type"))
					assert.EqualValues(t, 200, res.StatusCode)
					assert.EqualValues(t, scQuoteCallContractSc, *response[0].Quote)
				},
			},
		},
		{
			basicTestCase: basicTestCase{
				caseName: "Return quote successfully for SC origin and EOA call contract address",
				request: fmt.Sprintf("{"+basicQuoteFields+"}",
					destAddrEOA, value, rskRefAddrSC, rskRefAddrSC, btcRefAddr,
				),
				assertions: func(res *http.Response) {
					var response []*QuoteReturn
					json.NewDecoder(res.Body).Decode(&response)
					assert.EqualValues(t, "application/json", res.Header.Get("Content-Type"))
					assert.EqualValues(t, 200, res.StatusCode)
					assert.EqualValues(t, scQuoteCallContractEoa, *response[0].Quote)
				},
			},
		},
		{
			basicTestCase: basicTestCase{
				caseName: "Return error when transfer value is too high",
				request: fmt.Sprintf("{"+basicQuoteFields+"}",
					destAddrEOA, 600000000000000001, rskRefAddrEOA, rskRefAddrEOA, btcRefAddr,
				),
				assertions: func(res *http.Response) {
					response := &ErrorBody{}
					json.NewDecoder(res.Body).Decode(&response)
					assert.EqualValues(t, "application/json", res.Header.Get("Content-Type"))
					assert.EqualValues(t, 400, res.StatusCode)
					assert.EqualValues(t, "value to transfer is higher than max allowed", response.Message)
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.caseName, func(t *testing.T) {
			req, err := http.NewRequest("POST", "getQuote", bytes.NewReader([]byte(test.request)))
			if err != nil {
				t.Fatalf("couldn't instantiate request. error: %v", err)
			}
			if test.customMocks != nil {
				test.customMocks(rsk, btc, mongoDb)
			}
			defaultMocks(rsk, btc, mongoDb)
			rr := httptest.NewRecorder()

			srv.getQuoteHandler(rr, req)
			test.assertions(rr.Result())

			rsk.Calls = []mock.Call{}
			btc.Calls = []mock.Call{}
			mongoDb.Calls = []mock.Call{}
		})
	}

}

func testAcceptQuoteComplete(t *testing.T) {
	for _, quote := range testQuotes {
		hash := "555c9cfba7638a40a71a17a34fef0c3e192c1fbf4b311ad6e2ae288e97794228"
		rsk := new(testmocks.RskMock)
		btc := new(testmocks.BtcMock)
		lpRepo := new(storage.LPRepository)
		mongoDb, _ := testmocks.NewDbMock("", quote, nil)
		sat, _ := new(types.Wei).Add(quote.Value, quote.CallFee).ToSatoshi().Float64()
		minAmount := btcutil.Amount(uint64(math.Ceil(sat)))
		expTime := time.Unix(int64(quote.AgreementTimestamp+quote.TimeForDeposit), 0)
		fedInfo := &connectors.FedInfo{}

		srv := newServer(rsk, btc, mongoDb, func() time.Time {
			return time.Unix(0, 0)
		}, cfgData, lpRepo, providerCfgData)

		detailMock := types.ProviderRegisterRequest{}
		for _, lp := range providerMocks {
			rsk.On("GetCollateral", lp.address).Times(1).Return(big.NewInt(10), big.NewInt(10), nil)
			rsk.On("GetCollateral", lp.address).Return(big.NewInt(10), big.NewInt(10), nil)
			rsk.On("RegisterProvider", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
			mongoDb.On("InsertProvider", mock.Anything, mock.Anything).Return(nil)
			err := srv.AddProvider(lp, detailMock)
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

		mongoDb.On("GetQuote", hash).Times(1).Return(quote, nil)
		mongoDb.On("GetRetainedQuote", hash).Times(1).Return(nil, nil)
		rsk.On("GasPrice").Times(1)
		rsk.On("FetchFederationInfo").Times(1).Return(fedInfo, nil)
		btc.On("GetParams")
		rsk.On("GetDerivedBitcoinAddress", fedInfo, nil, btcRefAddr, lbcAddr, lpBTCAddr, hashBytes)
		btc.On("AddAddressWatcher", "", minAmount, time.Minute, expTime, mock.AnythingOfType("*http.BTCAddressWatcher"), mock.AnythingOfType("func(connectors.AddressWatcher)"))
		srv.acceptQuoteHandler(&w, req)
		mongoDb.AssertExpectations(t)
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
	lp := new(storage.LPRepository)
	mongoDb, _ := testmocks.NewDbMock(hash, quote, nil)
	sat, _ := new(types.Wei).Add(quote.Value, quote.CallFee).ToSatoshi().Float64()
	minAmount := btcutil.Amount(uint64(math.Ceil(sat)))
	expTime := time.Unix(int64(quote.AgreementTimestamp+quote.TimeForDeposit), 0)

	srv := newServer(rsk, btc, mongoDb, func() time.Time {
		return time.Unix(0, 0)
	}, cfgData, lp, providerCfgData)

	detailMock := types.ProviderRegisterRequest{}
	for _, lp := range providerMocks {
		rsk.On("GetCollateral", lp.address).Times(1).Return(big.NewInt(10), big.NewInt(10), nil)
		rsk.On("GetCollateral", lp.address).Return(big.NewInt(10), big.NewInt(10), nil)
		rsk.On("RegisterProvider", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
		mongoDb.On("InsertProvider", mock.Anything, mock.Anything).Return(nil)
		err := srv.AddProvider(lp, detailMock)
		if err != nil {
			t.Errorf("couldn't add provider. error: %v", err)
		}
	}

	mongoDb.On("GetRetainedQuotes", []types.RQState{types.RQStateWaitingForDeposit, types.RQStateCallForUserSucceeded}).Times(1).Return([]*types.RetainedQuote{{QuoteHash: hash}})
	mongoDb.On("GetQuote", hash).Times(1).Return(quote)
	btc.On("AddAddressWatcher", "", minAmount, time.Minute, expTime, mock.AnythingOfType("*http.BTCAddressWatcher"), mock.AnythingOfType("func(connectors.AddressWatcher)")).Times(1).Return("")
	err := srv.initBtcWatchers()
	if err != nil {
		t.Errorf("couldn't init BTC watchers. error: %v", err)
	}
	mongoDb.AssertExpectations(t)
	btc.AssertExpectations(t)
	rsk.AssertExpectations(t)
}

func testGetQuoteExpTime(t *testing.T) {
	quote := pegin.Quote{AgreementTimestamp: 2, TimeForDeposit: 3}
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

func testGetProviders(t *testing.T) {
	rsk := new(testmocks.RskMock)
	btc := new(testmocks.BtcMock)
	lp := new(storage.LPRepository)

	mongoDb, _ := testmocks.NewDbMock("", testQuotes[0], nil)

	srv := New(rsk, btc, mongoDb, cfgData, lp, providerCfgData)
	req, err := http.NewRequest("GET", "getProviders", bytes.NewReader([]byte("")))
	w := http2.TestResponseWriter{}

	if err != nil {
		t.Fatalf("couldn't instantiate request. error: %v", err)
	}

	mongoDb.On("GetProviders").Return([]int64{}, nil)
	rsk.On("GetProviders", mock.Anything).Return([]bindings.LiquidityBridgeContractProvider{}, nil)
	srv.getProvidersHandler(&w, req)

	assert.EqualValues(t, "application/json", w.Header().Get("Content-Type"))
	assert.EqualValues(t, "[]\n", w.Output)
}

func testcAcceptQuotePegoutComplete(t *testing.T) {
	for _, quote := range testPegOutQuotes {
		hash := "555c9cfba7638a40a71a17a34fef0c3e192c1fbf4b311ad6e2ae288e97794228"
		derivationAddress := "2NFwPDdXtAmGijQPbpK7s1z9bRGRx2SkB6D"
		rsk := new(testmocks.RskMock)
		btc := new(testmocks.BtcMock)
		lp := new(storage.LPRepository)

		mongoDb, _ := testmocks.NewDbMock("", nil, quote)
		minAmount := quote.Value + quote.CallFee
		expTime := time.Unix(int64(quote.AgreementTimestamp+quote.DepositDateLimit), 0)

		srv := newServer(rsk, btc, mongoDb, func() time.Time {
			return time.Unix(0, 0)
		}, cfgData, lp, providerCfgData)

		detailMock := types.ProviderRegisterRequest{}
		for _, lp := range providerPegOutMocks {
			rsk.On("GetCollateral", lp.address).Return(big.NewInt(10), big.NewInt(10), nil)
			rsk.On("GetCollateral", lp.address).Return(big.NewInt(10), big.NewInt(10), nil)
			rsk.On("RegisterProvider", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
			mongoDb.On("InsertProvider", mock.Anything, mock.Anything).Return(nil)
			err := srv.AddPegOutProvider(lp, detailMock)
			if err != nil {
				t.Fatalf("couldn't add provider. error: %v", err)
			}
		}

		w := http2.TestResponseWriter{}
		body := fmt.Sprintf("{\"quoteHash\":\"%v\", \"derivationAddress\":\"%v\"}", hash, derivationAddress)

		req, err := http.NewRequest("POST", "pegout/acceptQuote", bytes.NewReader([]byte(body)))
		if err != nil {
			t.Errorf("couldn't instantiate request. error: %v", err)
		}

		mongoDb.On("GetPegOutQuote", hash).Times(1).Return(quote, nil)
		rsk.On("AddQuoteToWatch", "555c9cfba7638a40a71a17a34fef0c3e192c1fbf4b311ad6e2ae288e97794228", time.Minute, expTime, mock.AnythingOfType("*http.RegisterPegoutWatcher"), mock.AnythingOfType("func(connectors.QuotePegOutWatcher)")).Times(1).Return("")

		btc.On("AddAddressWatcher", "2NFwPDdXtAmGijQPbpK7s1z9bRGRx2SkB6D", minAmount, time.Minute, expTime, mock.AnythingOfType("*http.BTCAddressWatcher"), mock.AnythingOfType("func(connectors.AddressWatcher)")).Times(1).Return("")
		srv.acceptQuotePegOutHandler(&w, req)
		response := AcceptResPegOut{}
		json.Unmarshal([]byte(w.Output), &response)
		assert.Equal(t, "fb4a3e40390dee7db6e861e10e5e3b39a0cf546eeccc8c0902249419140d9f29335023e3a83deee747f4987e9cd32773d2afa5176295dc2042255b57a30300201c", response.Signature)
	}
}

func testAddCollateral(t *testing.T) {
	rsk := new(testmocks.RskMock)
	srv := New(rsk, nil, nil, cfgData, nil, providerCfgData)

	for _, provider := range providerMocks {
		srv.providers = append(srv.providers, provider)
	}

	rsk.On("AddCollateral", mock.Anything).Return(nil).Once()
	rsk.On("GetCollateral", mock.Anything).Return(big.NewInt(2), big.NewInt(4), nil).Once()
	rsk.On("GetCollateral", providerMocks[1].address).Return(big.NewInt(7), big.NewInt(4), nil)

	testCases := []*basicTestCase{
		{
			caseName: "Returns 400 on missing address",
			request:  fmt.Sprintf(`{"amount": %v}`, 5),
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(&body)
				assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)
				assert.EqualValues(t, "LpRskAddress is required", body.Message)
			},
		},
		{
			caseName: "Returns 400 on missing amount",
			request:  fmt.Sprintf(`{"lpRskAddress": "%v"}`, providerMocks[1].address),
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(&body)
				assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)
				assert.EqualValues(t, "Amount is required", body.Message)
			},
		},
		{
			caseName: "Returns 400 on decoding error",
			request:  fmt.Sprint(""),
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(&body)
				assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)
				assert.Contains(t, body.Message, "Unable to deserialize payload")
			},
		},
		{
			caseName: "Returns 400 on invalid address",
			request:  fmt.Sprintf(`{"lpRskAddress": "%v", "amount": %v}`, providerMocks[0].address, 5),
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(&body)
				assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)
				assert.EqualValues(t, "LpRskAddress is eth_addr", body.Message)
			},
		},
		{
			caseName: "Returns 409 on non registered provider",
			request:  fmt.Sprintf(`{"lpRskAddress": "%v", "amount": %v}`, "0x9D93929A9099be4355fC2389FbF253982F9dF47c", 5),
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(&body)
				assert.EqualValues(t, http.StatusNotFound, res.StatusCode)
				assert.EqualValues(t, "missing liquidity provider", body.Message)
			},
		},
		{
			caseName: "Returns 409 on when provided collateral is lower than minimal",
			request:  fmt.Sprintf(`{"lpRskAddress": "%v", "amount": %v}`, providerMocks[1].address, 1),
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(&body)
				assert.EqualValues(t, http.StatusConflict, res.StatusCode)
				assert.EqualValues(t, "Amount is lower than min collateral", body.Message)
			},
		},
		{
			caseName: "Returns 200 on successful add",
			request:  fmt.Sprintf(`{"lpRskAddress": "%v", "amount": %v}`, providerMocks[1].address, 5),
			assertions: func(res *http.Response) {
				body := &AddCollateralResponse{}
				json.NewDecoder(res.Body).Decode(&body)
				assert.EqualValues(t, http.StatusOK, res.StatusCode)
				assert.EqualValues(t, uint64(7), body.NewCollateralBalance)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.caseName, func(t *testing.T) {
			req, err := http.NewRequest("POST", "addCollateral", bytes.NewReader([]byte(test.request)))
			if err != nil {
				t.Fatalf("couldn't instantiate request. error: %v", err)
			}
			rr := httptest.NewRecorder()
			srv.addCollateral(rr, req)
			test.assertions(rr.Result())
			rsk.Calls = []mock.Call{}
		})
	}
}

func testGetCollateral(t *testing.T) {
	rsk := new(testmocks.RskMock)
	srv := New(rsk, nil, nil, cfgData, nil, providerCfgData)

	for _, provider := range providerMocks {
		srv.providers = append(srv.providers, provider)
	}

	rsk.On("GetCollateral", providerMocks[0].address).Return(big.NewInt(0), big.NewInt(4), nil).Once()
	rsk.On("GetCollateral", providerMocks[0].address).Return(nil, nil, errors.New("some error"))
	rsk.On("GetCollateral", providerMocks[1].address).Return(big.NewInt(300), big.NewInt(4), nil)
	rsk.On("GetCollateral", "anything").Return(nil, nil, connectors.NewInvalidAddressError("anything"))

	testCases := []*basicTestCase{
		{
			caseName: "Fail on invalid address",
			request:  "anything",
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(body)
				assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)
				assert.Contains(t, body.Message, "invalid address")
			},
		},
		{
			caseName: "Return 404 when no collateral is found",
			request:  providerMocks[0].address,
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(body)
				assert.EqualValues(t, http.StatusNotFound, res.StatusCode)
				assert.EqualValues(t, body.Message, "no collateral found")
			},
		},
		{
			caseName: "Return collateral successfully",
			request:  providerMocks[1].address,
			assertions: func(res *http.Response) {
				body := &GetCollateralResponse{}
				json.NewDecoder(res.Body).Decode(body)
				assert.EqualValues(t, http.StatusOK, res.StatusCode)
				assert.EqualValues(t, uint64(300), body.Collateral)
			},
		},
		{
			caseName: "Return 500 on get collateral error",
			request:  providerMocks[0].address,
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(body)
				assert.EqualValues(t, http.StatusInternalServerError, res.StatusCode)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.caseName, func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("collateral?address=%s", test.request), nil)
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}
			rr := httptest.NewRecorder()
			srv.getCollateralHandler(rr, req)
			test.assertions(rr.Result())
		})
	}
}

func testWithdrawCollateral(t *testing.T) {
	request := fmt.Sprintf(`{ "lpRskAddress": "%s" }`, providerMocks[1].address)

	rsk := new(testmocks.RskMock)
	srv := New(rsk, nil, nil, cfgData, nil, providerCfgData)
	for _, provider := range providerMocks {
		srv.providers = append(srv.providers, provider)
	}

	rsk.On("WithdrawCollateral", mock.Anything).Return(connectors.WithdrawCollateralError).Once()
	rsk.On("WithdrawCollateral", mock.Anything).Return(errors.New("some error")).Once()
	rsk.On("WithdrawCollateral", mock.Anything).Return(nil).Once()

	testCases := []*basicTestCase{
		{
			caseName: "Fail on invalid address",
			request:  `{ "lpRskAddress": "anything" }`,
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(body)
				assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)
				assert.EqualValues(t, body.Message, "LpRskAddress is eth_addr")
			},
		},
		{
			caseName: "Fail on non registered provider",
			request:  `{ "lpRskAddress": "0xa554d96413FF72E93437C4072438302C38350EE3" }`,
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(body)
				assert.EqualValues(t, http.StatusNotFound, res.StatusCode)
				assert.EqualValues(t, body.Message, "missing liquidity provider")
			},
		},
		{
			caseName: "Fail when provider didn't resigned",
			request:  request,
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(body)
				assert.EqualValues(t, http.StatusConflict, res.StatusCode)
				assert.Contains(t, body.Message, "withdraw collateral error")
			},
		},
		{
			caseName: "Fail on transaction error",
			request:  request,
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(body)
				assert.EqualValues(t, http.StatusInternalServerError, res.StatusCode)
			},
		},
		{
			caseName: "Return 204 on successful update",
			request:  request,
			assertions: func(res *http.Response) {
				assert.EqualValues(t, http.StatusNoContent, res.StatusCode)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.caseName, func(t *testing.T) {
			req, err := http.NewRequest("POST", "withdrawCollateral", bytes.NewReader([]byte(test.request)))
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}
			rr := httptest.NewRecorder()
			srv.withdrawCollateral(rr, req)
			test.assertions(rr.Result())
		})
	}
}

func testProviderResign(t *testing.T) {
	request := fmt.Sprintf(`{ "lpRskAddress": "%s" }`, providerMocks[1].address)

	rsk := new(testmocks.RskMock)
	srv := New(rsk, nil, nil, cfgData, nil, providerCfgData)
	for _, provider := range providerMocks {
		srv.providers = append(srv.providers, provider)
	}

	rsk.On("Resign", mock.Anything).Return(connectors.ProviderResignError).Once()
	rsk.On("Resign", mock.Anything).Return(errors.New("some error")).Once()
	rsk.On("Resign", mock.Anything).Return(nil).Once()

	testCases := []*basicTestCase{
		{
			caseName: "Fail on invalid address",
			request:  `{ "lpRskAddress": "dsadasda" }`,
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(body)
				assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)
				assert.EqualValues(t, body.Message, "LpRskAddress is eth_addr")
			},
		},
		{
			caseName: "Fail on non registered provider",
			request:  `{ "lpRskAddress": "0x1eD614cd3443EFd9c70F04b6d777aed947A4b0c4" }`,
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(body)
				assert.EqualValues(t, http.StatusNotFound, res.StatusCode)
				assert.EqualValues(t, body.Message, "missing liquidity provider")
			},
		},
		{
			caseName: "Fail when provider has resigned before",
			request:  request,
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(body)
				assert.EqualValues(t, http.StatusConflict, res.StatusCode)
				assert.Contains(t, body.Message, "provider has already resigned")
			},
		},
		{
			caseName: "Fail on transaction error",
			request:  request,
			assertions: func(res *http.Response) {
				body := &ErrorBody{}
				json.NewDecoder(res.Body).Decode(body)
				assert.EqualValues(t, http.StatusInternalServerError, res.StatusCode)
			},
		},
		{
			caseName: "Return 204 on successful resign",
			request:  request,
			assertions: func(res *http.Response) {
				assert.EqualValues(t, http.StatusNoContent, res.StatusCode)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.caseName, func(t *testing.T) {
			req, err := http.NewRequest("POST", "provider/resignation", bytes.NewReader([]byte(test.request)))
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}
			rr := httptest.NewRecorder()
			srv.providerResignHandler(rr, req)
			test.assertions(rr.Result())
		})
	}
}

func TestLiquidityProviderServer(t *testing.T) {
	t.Run("get provider by address", testGetProviderByAddress)
	t.Run("check health", testCheckHealth)
	t.Run("get provider should return null when provider not found", testGetProviderByAddressWhenNotFoundShouldReturnNull)
	t.Run("get quote", testGetQuoteComplete)
	t.Run("accept quote", testAcceptQuoteComplete)
	t.Run("init BTC watchers", testInitBtcWatchers)
	t.Run("get quote exp time", testGetQuoteExpTime)
	t.Run("decode address", testDecodeAddress)
	t.Run("decode address with an invalid btcRefundAddr", testDecodeAddressWithAnInvalidBtcRefundAddr)
	t.Run("decode address with an invalid lpBTCAddrB", testDecodeAddressWithAnInvalidLpBTCAddrB)
	t.Run("decode address with an invalid lbcAddrB", testDecodeAddressWithAnInvalidLbcAddrB)
	t.Run("get registered providers", testGetProviders)
	t.Run("accept quote pegout", testcAcceptQuotePegoutComplete)
	t.Run("add collateral", testAddCollateral)
	t.Run("get collateral", testGetCollateral)
	t.Run("withdraw collateral", testWithdrawCollateral)
	t.Run("test provider resign", testProviderResign)
}
