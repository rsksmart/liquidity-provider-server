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
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetPegoutQuoteHandlerHappyPath(t *testing.T) {
	reqBody := pkg.PegoutQuoteRequest{
		To:               test.AnyBtcAddress,
		ValueToTransfer:  big.NewInt(1000000000000000000),
		RskRefundAddress: test.AnyRskAddress,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/pegout/getQuotes", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	expectedQuote := pegout.GetPegoutQuoteResult{
		PegoutQuote: createTestPegoutQuote(),
		Hash:        test.AnyHash,
	}

	mockUseCase := new(mocks.GetPegoutQuoteUseCaseMock)
	mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegout.QuoteRequest")).Return(expectedQuote, nil)

	handlerFunc := handlers.NewGetPegoutQuoteHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody []pkg.GetPegoutQuoteResponse
	err = json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	require.Len(t, responseBody, 1)
	assert.Equal(t, test.AnyHash, responseBody[0].QuoteHash)
	assert.Equal(t, expectedQuote.PegoutQuote.LbcAddress, responseBody[0].Quote.LBCAddr)
	assert.Equal(t, expectedQuote.PegoutQuote.LpRskAddress, responseBody[0].Quote.LPRSKAddr)

	mockUseCase.AssertExpectations(t)
}

// nolint:funlen,maintidx
func TestGetPegoutQuoteHandlerErrorCases(t *testing.T) {

	t.Run("should handle malformed JSON in request body", func(t *testing.T) {
		malformedJSON := []byte(`{"to": "some-btc-address"`)
		request := httptest.NewRequest(http.MethodPost, "/pegout/getQuotes", bytes.NewBuffer(malformedJSON))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPegoutQuoteUseCaseMock)

		handlerFunc := handlers.NewGetPegoutQuoteHandler(mockUseCase)
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
		reqBody := pkg.PegoutQuoteRequest{
			To:               "", // Missing required field
			ValueToTransfer:  big.NewInt(1000),
			RskRefundAddress: test.AnyRskAddress,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/getQuotes", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPegoutQuoteUseCaseMock)

		handlerFunc := handlers.NewGetPegoutQuoteHandler(mockUseCase)
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

	t.Run("should handle invalid rskRefundAddress format", func(t *testing.T) {
		reqBody := pkg.PegoutQuoteRequest{
			To:               test.AnyBtcAddress,
			ValueToTransfer:  big.NewInt(1000),
			RskRefundAddress: "invalid-address", // Invalid ETH address
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/getQuotes", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPegoutQuoteUseCaseMock)

		handlerFunc := handlers.NewGetPegoutQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertNotCalled(t, "Run")
	})

	t.Run("should return 400 on BtcAddressNotSupportedError", func(t *testing.T) {
		reqBody := createValidPegoutQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/getQuotes", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPegoutQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegout.QuoteRequest")).
			Return(pegout.GetPegoutQuoteResult{}, blockchain.BtcAddressNotSupportedError)

		handlerFunc := handlers.NewGetPegoutQuoteHandler(mockUseCase)
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
		reqBody := createValidPegoutQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/getQuotes", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPegoutQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegout.QuoteRequest")).
			Return(pegout.GetPegoutQuoteResult{}, blockchain.BtcAddressInvalidNetworkError)

		handlerFunc := handlers.NewGetPegoutQuoteHandler(mockUseCase)
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
		reqBody := createValidPegoutQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/getQuotes", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPegoutQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegout.QuoteRequest")).
			Return(pegout.GetPegoutQuoteResult{}, usecases.RskAddressNotSupportedError)

		handlerFunc := handlers.NewGetPegoutQuoteHandler(mockUseCase)
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
		reqBody := createValidPegoutQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/getQuotes", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPegoutQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegout.QuoteRequest")).
			Return(pegout.GetPegoutQuoteResult{}, usecases.TxBelowMinimumError)

		handlerFunc := handlers.NewGetPegoutQuoteHandler(mockUseCase)
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
		reqBody := createValidPegoutQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/getQuotes", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPegoutQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegout.QuoteRequest")).
			Return(pegout.GetPegoutQuoteResult{}, liquidity_provider.AmountOutOfRangeError)

		handlerFunc := handlers.NewGetPegoutQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "invalid request", errorResponse["message"])
	})

	t.Run("should return 409 on NoLiquidityError", func(t *testing.T) {
		reqBody := createValidPegoutQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/getQuotes", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPegoutQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegout.QuoteRequest")).
			Return(pegout.GetPegoutQuoteResult{}, usecases.NoLiquidityError)

		handlerFunc := handlers.NewGetPegoutQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusConflict, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "no enough liquidity", errorResponse["message"])
	})

	t.Run("should return 500 on unexpected errors", func(t *testing.T) {
		reqBody := createValidPegoutQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/getQuotes", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPegoutQuoteUseCaseMock)
		unexpectedError := errors.New("unexpected database error")
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegout.QuoteRequest")).
			Return(pegout.GetPegoutQuoteResult{}, unexpectedError)

		handlerFunc := handlers.NewGetPegoutQuoteHandler(mockUseCase)
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
func TestGetPegoutQuoteHandlerErrorResponseFormat(t *testing.T) {

	t.Run("should set correct content type header on error", func(t *testing.T) {
		reqBody := createValidPegoutQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/getQuotes", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPegoutQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegout.QuoteRequest")).
			Return(pegout.GetPegoutQuoteResult{}, usecases.TxBelowMinimumError)

		handlerFunc := handlers.NewGetPegoutQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should include timestamp in error response", func(t *testing.T) {
		reqBody := createValidPegoutQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/getQuotes", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPegoutQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegout.QuoteRequest")).
			Return(pegout.GetPegoutQuoteResult{}, usecases.TxBelowMinimumError)

		handlerFunc := handlers.NewGetPegoutQuoteHandler(mockUseCase)
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
		reqBody := createValidPegoutQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/getQuotes", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPegoutQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegout.QuoteRequest")).
			Return(pegout.GetPegoutQuoteResult{}, usecases.NoLiquidityError)

		handlerFunc := handlers.NewGetPegoutQuoteHandler(mockUseCase)
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
		reqBody := createValidPegoutQuoteRequest()
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/getQuotes", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPegoutQuoteUseCaseMock)
		unexpectedError := errors.New("unexpected error")
		mockUseCase.On("Run", mock.Anything, mock.AnythingOfType("pegout.QuoteRequest")).
			Return(pegout.GetPegoutQuoteResult{}, unexpectedError)

		handlerFunc := handlers.NewGetPegoutQuoteHandler(mockUseCase)
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

func createValidPegoutQuoteRequest() pkg.PegoutQuoteRequest {
	return pkg.PegoutQuoteRequest{
		To:               test.AnyBtcAddress,
		ValueToTransfer:  big.NewInt(1000000000000000000),
		RskRefundAddress: test.AnyRskAddress,
	}
}

func createTestPegoutQuote() quote.PegoutQuote {
	return quote.PegoutQuote{
		LbcAddress:            "0x85FaB18a0d06fb14651c8F5EE9C7f4b00D80d70c",
		LpRskAddress:          "0x9D93929A9099be4355fC2389FbF253982F9dF47c",
		BtcRefundAddress:      test.AnyBtcAddress,
		RskRefundAddress:      test.AnyRskAddress,
		LpBtcAddress:          test.AnyBtcAddress,
		CallFee:               entities.NewWei(100),
		PenaltyFee:            entities.NewWei(200),
		Nonce:                 1,
		DepositAddress:        "0xDepositAddress",
		Value:                 entities.NewWei(1000000000000000000),
		AgreementTimestamp:    1640995200,
		DepositDateLimit:      1641081600,
		DepositConfirmations:  6,
		TransferConfirmations: 10,
		TransferTime:          3600,
		ExpireDate:            1641168000,
		ExpireBlock:           1000,
		GasFee:                entities.NewWei(50),
		ProductFeeAmount:      entities.NewWei(25),
	}
}
