package btcclient_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin/btcclient"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

const testString = "test"

var createWalletError = []byte("{\n    \"result\": null,\n    \"error\": {\n        \"code\": -4,\n        \"message\": \"Wallet file verification failed. Failed to create database path '/home/bitcoin/.bitcoin/regtest/wallets/postman5'. Database already exists.\"\n    },\n    \"id\": \"curltest\"\n}")
var createWalletResponse = []byte("{\n    \"result\": {\n        \"name\": \"postman\",\n        \"warning\": \"Wallet created successfully. The legacy wallet type is being deprecated and support for creating and opening legacy wallets will be removed in the future.\"\n    },\n    \"error\": null,\n    \"id\": \"curltest\"\n}")
var createWalletParams = btcclient.ReadonlyWalletRequest{
	WalletName:         test.AnyString,
	DisablePrivateKeys: true,
	Blank:              true,
	AvoidReuse:         true,
	Descriptors:        true,
}

func TestBtcSuiteClientAdapter_CreateReadonlyWallet(t *testing.T) {
	client := &mocks.RpcClientMock{}
	httpClient := &mocks.HttpClientMock{}
	adapter := btcclient.NewBtcSuiteClientAdapter(
		rpcclient.ConnConfig{DisableTLS: true, Host: testString + ":1234", User: testString, Pass: testString},
		client,
	)
	adapter.SetClient(httpClient)
	client.On("NextID").Return(uint64(1))
	reqBody, err := json.Marshal(btcclient.RpcRequestParamsObject[btcclient.ReadonlyWalletRequest]{
		Jsonrpc: btcjson.RpcVersion1,
		Method:  "createwallet",
		Params:  createWalletParams,
		ID:      1,
	})
	require.NoError(t, err)
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"http://test:1234",
		bytes.NewReader(reqBody),
	)
	require.NoError(t, err)
	req.SetBasicAuth(testString, testString)
	httpClient.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		receivedBody, readErr := io.ReadAll(r.Body)
		require.NoError(t, readErr)
		expectedBody, readErr := io.ReadAll(req.Body)
		require.NoError(t, readErr)
		return req.URL.String() == r.URL.String() && bytes.Equal(receivedBody, expectedBody) && r.Method == req.Method
	})).Return(&http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		Body:          io.NopCloser(bytes.NewReader(createWalletResponse)),
		ContentLength: int64(len(createWalletResponse)),
	}, nil)
	err = adapter.CreateReadonlyWallet(createWalletParams)
	require.NoError(t, err)
	client.AssertExpectations(t)
	httpClient.AssertExpectations(t)
}

func TestBtcSuiteClientAdapter_CreateReadonlyWallet_ErrorHanlding(t *testing.T) {
	t.Run("RPC server error", func(t *testing.T) {
		client := &mocks.RpcClientMock{}
		httpClient := &mocks.HttpClientMock{}
		adapter := btcclient.NewBtcSuiteClientAdapter(
			rpcclient.ConnConfig{DisableTLS: true, Host: testString + ":1234", User: testString, Pass: testString},
			client,
		)
		adapter.SetClient(httpClient)
		client.On("NextID").Return(uint64(1))

		httpClient.On("Do", mock.Anything).Return(&http.Response{
			Status:     "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, nil)
		err := adapter.CreateReadonlyWallet(createWalletParams)
		require.Error(t, err)
		client.AssertExpectations(t)
		httpClient.AssertExpectations(t)
	})
	t.Run("RPC client error", func(t *testing.T) {
		client := &mocks.RpcClientMock{}
		httpClient := &mocks.HttpClientMock{}
		adapter := btcclient.NewBtcSuiteClientAdapter(
			rpcclient.ConnConfig{DisableTLS: false, Host: testString + ":1234", User: testString, Pass: testString},
			client,
		)
		adapter.SetClient(httpClient)
		client.On("NextID").Return(uint64(1))

		httpClient.On("Do", mock.MatchedBy(func(r *http.Request) bool {
			return r.URL.String() == "https://test:1234"
		})).Return(&http.Response{
			Status:        "Bad request",
			StatusCode:    http.StatusBadRequest,
			Body:          io.NopCloser(bytes.NewReader(createWalletError)),
			ContentLength: int64(len(createWalletError)),
		}, nil)
		err := adapter.CreateReadonlyWallet(createWalletParams)
		require.Error(t, err)
		client.AssertExpectations(t)
		httpClient.AssertExpectations(t)
	})
}

func TestBtcSuiteClientAdapter_SignRawTransactionWithKey(t *testing.T) {
	keys := []string{"key"}
	client := &mocks.RpcClientMock{}
	receiveChannel := make(chan *rpcclient.Response, 1)
	client.On("SendCmd", &btcclient.SignRawTransactionWithKeyCmd{RawTx: rawTx, WifKeys: keys}).
		Return(receiveChannel)
	adapter := btcclient.NewBtcSuiteClientAdapter(
		rpcclient.ConnConfig{DisableTLS: true, Host: testString + ":1234", User: testString, Pass: testString},
		client,
	)
	tx := &wire.MsgTx{}
	txBytes, err := hex.DecodeString(rawTx)
	require.NoError(t, err)
	err = tx.DeserializeNoWitness(bytes.NewReader(txBytes))
	require.NoError(t, err)
	receiveChannel <- &rpcclient.Response{}
	_, _, _ = adapter.SignRawTransactionWithKey(tx, keys)
	client.AssertExpectations(t)
}
