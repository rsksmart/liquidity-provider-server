package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetPeginQuoteHandlerHappyPath(t *testing.T) {
	reqBody := pkg.PeginQuoteRequest{
		CallEoaOrContractAddress: test.AnyRskAddress,
		CallContractArguments:    "0x1234",
		ValueToTransfer:          big.NewInt(1000000000000000000),
		RskRefundAddress:         test.AnyRskAddress,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/pegin/getQuote", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	expectedQuote := pegin.GetPeginQuoteResult{
		PeginQuote: createTestPeginQuote(),
		Hash:       test.AnyHash,
	}

	mockUseCase := new(mocks.GetPeginQuoteUseCaseMock)
	mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegin.QuoteRequest")).Return(expectedQuote, nil)

	handlerFunc := handlers.NewGetPeginQuoteHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody []pkg.GetPeginQuoteResponse
	err = json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	require.Len(t, responseBody, 1)
	assert.Equal(t, test.AnyHash, responseBody[0].QuoteHash)
	assert.Equal(t, expectedQuote.PeginQuote.FedBtcAddress, responseBody[0].Quote.FedBTCAddr)
	assert.Equal(t, expectedQuote.PeginQuote.LbcAddress, responseBody[0].Quote.LBCAddr)

	mockUseCase.AssertExpectations(t)
}

// nolint:funlen,maintidx
func TestGetPeginQuoteHandlerErrorCases(t *testing.T) {

	t.Run("should handle malformed JSON in request body", func(t *testing.T) {
		malformedJSON := []byte(`{"callEoaOrContractAddress": "0x123"`)
		request := httptest.NewRequest(http.MethodPost, "/pegin/getQuote", bytes.NewBuffer(malformedJSON))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginQuoteUseCaseMock)

		handlerFunc := handlers.NewGetPeginQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertNotCalled(t, "Run")

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Equal(t, "Error decoding request", errorResponse["message"])
	})

	t.Run("should handle request validation failure - missing required fields", func(t *testing.T) {
		reqBody := pkg.PeginQuoteRequest{
			CallEoaOrContractAddress: "", // Missing required field
			ValueToTransfer:          big.NewInt(1000),
			RskRefundAddress:         test.AnyRskAddress,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/getQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginQuoteUseCaseMock)

		handlerFunc := handlers.NewGetPeginQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertNotCalled(t, "Run")

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Contains(t, errorResponse["message"], "validation error")
	})

	t.Run("should handle invalid eth_addr format", func(t *testing.T) {
		reqBody := pkg.PeginQuoteRequest{
			CallEoaOrContractAddress: "invalid-address", // Invalid ETH address
			CallContractArguments:    "0x1234",
			ValueToTransfer:          big.NewInt(1000),
			RskRefundAddress:         test.AnyRskAddress,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/getQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginQuoteUseCaseMock)

		handlerFunc := handlers.NewGetPeginQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertNotCalled(t, "Run")
	})

	t.Run("should return 400 on invalid callContractArguments hex", func(t *testing.T) {
		reqBody := pkg.PeginQuoteRequest{
			CallEoaOrContractAddress: test.AnyRskAddress,
			CallContractArguments:    "0xGGGG", // Invalid hex
			ValueToTransfer:          big.NewInt(1000),
			RskRefundAddress:         test.AnyRskAddress,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/getQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginQuoteUseCaseMock)

		handlerFunc := handlers.NewGetPeginQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertNotCalled(t, "Run")

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
	})

	t.Run("should return 400 on BtcAddressNotSupportedError", func(t *testing.T) {
		reqBody := createValidPeginQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/getQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegin.QuoteRequest")).
			Return(pegin.GetPeginQuoteResult{}, blockchain.BtcAddressNotSupportedError)

		handlerFunc := handlers.NewGetPeginQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "invalid request", errorResponse["message"])
	})

	t.Run("should return 400 on BtcAddressInvalidNetworkError", func(t *testing.T) {
		reqBody := createValidPeginQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/getQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegin.QuoteRequest")).
			Return(pegin.GetPeginQuoteResult{}, blockchain.BtcAddressInvalidNetworkError)

		handlerFunc := handlers.NewGetPeginQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "invalid request", errorResponse["message"])
	})

	t.Run("should return 400 on RskAddressNotSupportedError", func(t *testing.T) {
		reqBody := createValidPeginQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/getQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegin.QuoteRequest")).
			Return(pegin.GetPeginQuoteResult{}, usecases.RskAddressNotSupportedError)

		handlerFunc := handlers.NewGetPeginQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "invalid request", errorResponse["message"])
	})

	t.Run("should return 400 on TxBelowMinimumError", func(t *testing.T) {
		reqBody := createValidPeginQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/getQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegin.QuoteRequest")).
			Return(pegin.GetPeginQuoteResult{}, usecases.TxBelowMinimumError)

		handlerFunc := handlers.NewGetPeginQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "invalid request", errorResponse["message"])
	})

	t.Run("should return 400 on AmountOutOfRangeError", func(t *testing.T) {
		reqBody := createValidPeginQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/getQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegin.QuoteRequest")).
			Return(pegin.GetPeginQuoteResult{}, liquidity_provider.AmountOutOfRangeError)

		handlerFunc := handlers.NewGetPeginQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "invalid request", errorResponse["message"])
	})

	t.Run("should return 500 on unexpected errors", func(t *testing.T) {
		reqBody := createValidPeginQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/getQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginQuoteUseCaseMock)
		unexpectedError := errors.New("unexpected database error")
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegin.QuoteRequest")).
			Return(pegin.GetPeginQuoteResult{}, unexpectedError)

		handlerFunc := handlers.NewGetPeginQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "unknown error", errorResponse["message"])
	})
}

// nolint:funlen
func TestGetPeginQuoteHandlerErrorResponseFormat(t *testing.T) {

	t.Run("should set correct content type header on error", func(t *testing.T) {
		reqBody := createValidPeginQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/getQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegin.QuoteRequest")).
			Return(pegin.GetPeginQuoteResult{}, usecases.TxBelowMinimumError)

		handlerFunc := handlers.NewGetPeginQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should include timestamp in error response", func(t *testing.T) {
		reqBody := createValidPeginQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/getQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegin.QuoteRequest")).
			Return(pegin.GetPeginQuoteResult{}, usecases.TxBelowMinimumError)

		handlerFunc := handlers.NewGetPeginQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "timestamp")
		assert.NotZero(t, errorResponse["timestamp"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should include recoverable flag as true for domain errors", func(t *testing.T) {
		reqBody := createValidPeginQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/getQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegin.QuoteRequest")).
			Return(pegin.GetPeginQuoteResult{}, liquidity_provider.AmountOutOfRangeError)

		handlerFunc := handlers.NewGetPeginQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "recoverable")
		assert.Equal(t, true, errorResponse["recoverable"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should set recoverable to false on unexpected errors", func(t *testing.T) {
		reqBody := createValidPeginQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/getQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginQuoteUseCaseMock)
		unexpectedError := errors.New("unexpected error")
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegin.QuoteRequest")).
			Return(pegin.GetPeginQuoteResult{}, unexpectedError)

		handlerFunc := handlers.NewGetPeginQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "recoverable")
		assert.Equal(t, false, errorResponse["recoverable"])
		mockUseCase.AssertExpectations(t)
	})
}

// Helper functions

func createValidPeginQuoteRequest() pkg.PeginQuoteRequest {
	return pkg.PeginQuoteRequest{
		CallEoaOrContractAddress: test.AnyRskAddress,
		CallContractArguments:    "0x1234",
		ValueToTransfer:          big.NewInt(1000000000000000000),
		RskRefundAddress:         test.AnyRskAddress,
	}
}

func createTestPeginQuote() quote.PeginQuote {
	return quote.PeginQuote{
		FedBtcAddress:      "2N5W5MxrGKMNNRzoBMN2hKKUNxEJUUuGcLp",
		LbcAddress:         "0x85FaB18a0d06fb14651c8F5EE9C7f4b00D80d70c",
		LpRskAddress:       "0x9D93929A9099be4355fC2389FbF253982F9dF47c",
		BtcRefundAddress:   "2MvMxL8KLzw4R8Y9wQP8QNNpYQqGKSUJe6J",
		RskRefundAddress:   test.AnyRskAddress,
		LpBtcAddress:       "2N7Vw5f59V3o3bDcaJK5oA829LFTBYZHLoG",
		CallFee:            entities.NewWei(100),
		PenaltyFee:         entities.NewWei(200),
		ContractAddress:    test.AnyRskAddress,
		Data:               "0x1234",
		GasLimit:           21000,
		Nonce:              1,
		Value:              entities.NewWei(1000000000000000000),
		AgreementTimestamp: 1640995200,
		TimeForDeposit:     3600,
		LpCallTime:         1800,
		Confirmations:      6,
		CallOnRegister:     true,
		GasFee:             entities.NewWei(50),
	}
}
