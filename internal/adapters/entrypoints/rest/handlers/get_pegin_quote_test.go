package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	anyContractAddress = "0x79568c2989232dCa1840087D73d403602364c0D4"
	anyRefundAddress   = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
)

func validPeginQuoteRequest() pkg.PeginQuoteRequest {
	return pkg.PeginQuoteRequest{
		CallEoaOrContractAddress: anyContractAddress,
		CallContractArguments:    "0xabcdef",
		ValueToTransfer:          big.NewInt(1000000),
		RskRefundAddress:         anyRefundAddress,
	}
}

func anyPeginQuoteResult() pegin.GetPeginQuoteResult {
	return pegin.GetPeginQuoteResult{
		PeginQuote: quote.PeginQuote{
			FedBtcAddress:    "2N5muMepJizJE1gR7FbHJU6CD18V3BpNF9p",
			LbcAddress:       anyContractAddress,
			LpRskAddress:     anyContractAddress,
			BtcRefundAddress: "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe",
			RskRefundAddress: anyRefundAddress,
			LpBtcAddress:     "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe",
			CallFee:          entities.NewWei(1000),
			PenaltyFee:       entities.NewWei(500),
			ContractAddress:  anyContractAddress,
			Data:             "abcdef",
			GasLimit:         21000,
			Nonce:            42,
			Value:            entities.NewWei(1000000),
			TimeForDeposit:   3600,
			LpCallTime:       7200,
			Confirmations:    6,
			GasFee:           entities.NewWei(200),
			ProductFeeAmount: entities.NewWei(100),
		},
		Hash: "abc123def456abc123def456abc123def456abc123def456abc123def456abc1",
	}
}

// nolint:funlen,maintidx
func TestNewGetPeginQuoteHandler(t *testing.T) {
	const path = "/pegin/getQuote"

	tests := []struct {
		name           string
		buildBody      func() ([]byte, error)
		setupMock      func(useCase *mocks.GetPeginQuoteUseCaseMock)
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "should return 200 with quote on valid request",
			buildBody: func() ([]byte, error) {
				return json.Marshal(validPeginQuoteRequest())
			},
			setupMock: func(useCase *mocks.GetPeginQuoteUseCaseMock) {
				result := anyPeginQuoteResult()
				useCase.EXPECT().Run(mock.Anything, mock.Anything).Return(result, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response []pkg.GetPeginQuoteResponse
				require.NoError(t, json.Unmarshal(body, &response))
				require.Len(t, response, 1)
				result := anyPeginQuoteResult()
				assert.Equal(t, result.Hash, response[0].QuoteHash)
				assert.Equal(t, result.PeginQuote.ContractAddress, response[0].Quote.ContractAddr)
				assert.Equal(t, result.PeginQuote.RskRefundAddress, response[0].Quote.RSKRefundAddr)
				assert.Equal(t, result.PeginQuote.Value.AsBigInt(), response[0].Quote.Value)
			},
		},
		{
			name: "should return 200 with empty callContractArguments",
			buildBody: func() ([]byte, error) {
				req := validPeginQuoteRequest()
				req.CallContractArguments = ""
				return json.Marshal(req)
			},
			setupMock: func(useCase *mocks.GetPeginQuoteUseCaseMock) {
				result := anyPeginQuoteResult()
				useCase.EXPECT().Run(mock.Anything, mock.Anything).Return(result, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse:  nil,
		},
		{
			name: "should return 200 with callContractArguments without 0x prefix",
			buildBody: func() ([]byte, error) {
				req := validPeginQuoteRequest()
				req.CallContractArguments = "abcdef"
				return json.Marshal(req)
			},
			setupMock: func(useCase *mocks.GetPeginQuoteUseCaseMock) {
				result := anyPeginQuoteResult()
				useCase.EXPECT().Run(mock.Anything, mock.Anything).Return(result, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse:  nil,
		},
		{
			name: "should return 400 on malformed JSON body",
			buildBody: func() ([]byte, error) {
				return []byte(`{"callEoaOrContractAddress": "0x79568c2989232dCa1840087D73d403602364c0D4"`), nil
			},
			setupMock:      func(useCase *mocks.GetPeginQuoteUseCaseMock) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.Equal(t, "Error decoding request", errorResponse.Message)
				assert.True(t, errorResponse.Recoverable)
			},
		},
		{
			name: "should return 400 when callEoaOrContractAddress is missing",
			buildBody: func() ([]byte, error) {
				req := validPeginQuoteRequest()
				req.CallEoaOrContractAddress = ""
				return json.Marshal(req)
			},
			setupMock:      func(useCase *mocks.GetPeginQuoteUseCaseMock) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.Equal(t, "validation error", errorResponse.Message)
				assert.True(t, errorResponse.Recoverable)
				assert.Contains(t, errorResponse.Details, "CallEoaOrContractAddress")
			},
		},
		{
			name: "should return 400 when callEoaOrContractAddress is not a valid eth address",
			buildBody: func() ([]byte, error) {
				req := validPeginQuoteRequest()
				req.CallEoaOrContractAddress = "not-an-address"
				return json.Marshal(req)
			},
			setupMock:      func(useCase *mocks.GetPeginQuoteUseCaseMock) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.Equal(t, "validation error", errorResponse.Message)
				assert.True(t, errorResponse.Recoverable)
				assert.Contains(t, errorResponse.Details, "CallEoaOrContractAddress")
			},
		},
		{
			name: "should return 400 when rskRefundAddress is missing",
			buildBody: func() ([]byte, error) {
				req := validPeginQuoteRequest()
				req.RskRefundAddress = ""
				return json.Marshal(req)
			},
			setupMock:      func(useCase *mocks.GetPeginQuoteUseCaseMock) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.Equal(t, "validation error", errorResponse.Message)
				assert.True(t, errorResponse.Recoverable)
				assert.Contains(t, errorResponse.Details, "RskRefundAddress")
			},
		},
		{
			name: "should return 400 when rskRefundAddress is not a valid eth address",
			buildBody: func() ([]byte, error) {
				req := validPeginQuoteRequest()
				req.RskRefundAddress = "not-an-address"
				return json.Marshal(req)
			},
			setupMock:      func(useCase *mocks.GetPeginQuoteUseCaseMock) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.Equal(t, "validation error", errorResponse.Message)
				assert.True(t, errorResponse.Recoverable)
				assert.Contains(t, errorResponse.Details, "RskRefundAddress")
			},
		},
		{
			name: "should return 400 when valueToTransfer is missing",
			buildBody: func() ([]byte, error) {
				req := validPeginQuoteRequest()
				req.ValueToTransfer = nil
				return json.Marshal(req)
			},
			setupMock:      func(useCase *mocks.GetPeginQuoteUseCaseMock) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.Equal(t, "validation error", errorResponse.Message)
				assert.True(t, errorResponse.Recoverable)
				assert.Contains(t, errorResponse.Details, "ValueToTransfer")
			},
		},
		{
			name: "should return 400 when callContractArguments exceeds max length",
			buildBody: func() ([]byte, error) {
				req := validPeginQuoteRequest()
				oversized := make([]byte, 8195)
				for i := range oversized {
					oversized[i] = 'a'
				}
				req.CallContractArguments = string(oversized)
				return json.Marshal(req)
			},
			setupMock:      func(useCase *mocks.GetPeginQuoteUseCaseMock) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.Equal(t, "validation error", errorResponse.Message)
				assert.True(t, errorResponse.Recoverable)
				assert.Contains(t, errorResponse.Details, "CallContractArguments")
			},
		},
		{
			name: "should return 400 when callContractArguments is not valid hex",
			buildBody: func() ([]byte, error) {
				req := validPeginQuoteRequest()
				req.CallContractArguments = "not-hex-data!"
				return json.Marshal(req)
			},
			setupMock:      func(useCase *mocks.GetPeginQuoteUseCaseMock) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.True(t, errorResponse.Recoverable)
			},
		},
		{
			name: "should return 400 when use case returns BtcAddressNotSupportedError",
			buildBody: func() ([]byte, error) {
				return json.Marshal(validPeginQuoteRequest())
			},
			setupMock: func(useCase *mocks.GetPeginQuoteUseCaseMock) {
				useCase.EXPECT().Run(mock.Anything, mock.Anything).
					Return(pegin.GetPeginQuoteResult{}, blockchain.BtcAddressNotSupportedError)
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.Equal(t, "invalid request", errorResponse.Message)
				assert.True(t, errorResponse.Recoverable)
				assert.Contains(t, errorResponse.Details, "error")
			},
		},
		{
			name: "should return 400 when use case returns BtcAddressInvalidNetworkError",
			buildBody: func() ([]byte, error) {
				return json.Marshal(validPeginQuoteRequest())
			},
			setupMock: func(useCase *mocks.GetPeginQuoteUseCaseMock) {
				useCase.EXPECT().Run(mock.Anything, mock.Anything).
					Return(pegin.GetPeginQuoteResult{}, blockchain.BtcAddressInvalidNetworkError)
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.Equal(t, "invalid request", errorResponse.Message)
				assert.True(t, errorResponse.Recoverable)
				assert.Contains(t, errorResponse.Details, "error")
			},
		},
		{
			name: "should return 400 when use case returns RskAddressNotSupportedError",
			buildBody: func() ([]byte, error) {
				return json.Marshal(validPeginQuoteRequest())
			},
			setupMock: func(useCase *mocks.GetPeginQuoteUseCaseMock) {
				useCase.EXPECT().Run(mock.Anything, mock.Anything).
					Return(pegin.GetPeginQuoteResult{}, usecases.RskAddressNotSupportedError)
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.Equal(t, "invalid request", errorResponse.Message)
				assert.True(t, errorResponse.Recoverable)
				assert.Contains(t, errorResponse.Details, "error")
			},
		},
		{
			name: "should return 400 when use case returns TxBelowMinimumError",
			buildBody: func() ([]byte, error) {
				return json.Marshal(validPeginQuoteRequest())
			},
			setupMock: func(useCase *mocks.GetPeginQuoteUseCaseMock) {
				useCase.EXPECT().Run(mock.Anything, mock.Anything).
					Return(pegin.GetPeginQuoteResult{}, usecases.TxBelowMinimumError)
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.Equal(t, "invalid request", errorResponse.Message)
				assert.True(t, errorResponse.Recoverable)
				assert.Contains(t, errorResponse.Details, "error")
			},
		},
		{
			name: "should return 400 when use case returns DataCapExceededError",
			buildBody: func() ([]byte, error) {
				return json.Marshal(validPeginQuoteRequest())
			},
			setupMock: func(useCase *mocks.GetPeginQuoteUseCaseMock) {
				useCase.EXPECT().Run(mock.Anything, mock.Anything).
					Return(pegin.GetPeginQuoteResult{}, pegin.DataCapExceededError)
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.Equal(t, "invalid request", errorResponse.Message)
				assert.True(t, errorResponse.Recoverable)
				assert.Contains(t, errorResponse.Details, "error")
			},
		},
		{
			name: "should return 400 when use case returns AmountOutOfRangeError",
			buildBody: func() ([]byte, error) {
				return json.Marshal(validPeginQuoteRequest())
			},
			setupMock: func(useCase *mocks.GetPeginQuoteUseCaseMock) {
				useCase.EXPECT().Run(mock.Anything, mock.Anything).
					Return(pegin.GetPeginQuoteResult{}, liquidity_provider.AmountOutOfRangeError)
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.Equal(t, "invalid request", errorResponse.Message)
				assert.True(t, errorResponse.Recoverable)
				assert.Contains(t, errorResponse.Details, "error")
			},
		},
		{
			name: "should return 400 when use case returns a wrapped known error",
			buildBody: func() ([]byte, error) {
				return json.Marshal(validPeginQuoteRequest())
			},
			setupMock: func(useCase *mocks.GetPeginQuoteUseCaseMock) {
				wrappedErr := usecases.WrapUseCaseError(usecases.GetPeginQuoteId, usecases.TxBelowMinimumError)
				useCase.EXPECT().Run(mock.Anything, mock.Anything).
					Return(pegin.GetPeginQuoteResult{}, wrappedErr)
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.Equal(t, "invalid request", errorResponse.Message)
				assert.True(t, errorResponse.Recoverable)
			},
		},
		{
			name: "should return 500 when use case returns an unknown error",
			buildBody: func() ([]byte, error) {
				return json.Marshal(validPeginQuoteRequest())
			},
			setupMock: func(useCase *mocks.GetPeginQuoteUseCaseMock) {
				useCase.EXPECT().Run(mock.Anything, mock.Anything).
					Return(pegin.GetPeginQuoteResult{}, errors.New("unexpected internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body []byte) {
				var errorResponse rest.ErrorResponse
				require.NoError(t, json.Unmarshal(body, &errorResponse))
				assert.Equal(t, handlers.UnknownErrorMessage, errorResponse.Message)
				assert.False(t, errorResponse.Recoverable)
				assert.Contains(t, errorResponse.Details, "error")
			},
		},
		{
			name: "should not call use case when request decoding fails",
			buildBody: func() ([]byte, error) {
				return []byte(`{invalid json`), nil
			},
			setupMock: func(useCase *mocks.GetPeginQuoteUseCaseMock) {
				useCase.AssertNotCalled(t, "Run")
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "should not call use case when request validation fails",
			buildBody: func() ([]byte, error) {
				req := validPeginQuoteRequest()
				req.CallEoaOrContractAddress = ""
				return json.Marshal(req)
			},
			setupMock: func(useCase *mocks.GetPeginQuoteUseCaseMock) {
				useCase.AssertNotCalled(t, "Run")
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := tt.buildBody()
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			useCase := new(mocks.GetPeginQuoteUseCaseMock)
			tt.setupMock(useCase)

			handler := handlers.NewGetPeginQuoteHandler(useCase)
			handler.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)
			assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder.Body.Bytes())
			}
		})
	}
}
