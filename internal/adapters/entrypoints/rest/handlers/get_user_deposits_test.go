package handlers_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetUserQuotesHandlerHappyPath(t *testing.T) {
	address := "0x1234567890abcdef1234567890abcdef12345678"
	request := httptest.NewRequest(http.MethodGet, "/userQuotes?address="+address, nil)
	recorder := httptest.NewRecorder()

	timestamp := time.Now()
	expectedDeposits := []quote.PegoutDeposit{
		{
			TxHash:      "0xabc123",
			QuoteHash:   "0xdef456",
			Amount:      entities.NewWei(1000000000000000000),
			Timestamp:   timestamp,
			BlockNumber: 12345,
			From:        address,
		},
		{
			TxHash:      "0x789xyz",
			QuoteHash:   "0x012abc",
			Amount:      entities.NewWei(2000000000000000000),
			Timestamp:   timestamp.Add(time.Hour),
			BlockNumber: 12346,
			From:        address,
		},
	}

	mockUseCase := new(mocks.GetUserDepositsUseCaseMock)
	mockUseCase.On("Run", mock.Anything, address).Return(expectedDeposits, nil)

	handlerFunc := handlers.NewGetUserQuotesHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody []pkg.DepositEventDTO
	err := json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	assert.Len(t, responseBody, 2)
	assert.Equal(t, "0xdef456", responseBody[0].QuoteHash)
	assert.Equal(t, big.NewInt(1000000000000000000), responseBody[0].Amount)
	assert.Equal(t, address, responseBody[0].From)

	mockUseCase.AssertExpectations(t)
}

func TestGetUserQuotesHandlerEmptyList(t *testing.T) {
	address := "0x1234567890abcdef1234567890abcdef12345678"
	request := httptest.NewRequest(http.MethodGet, "/userQuotes?address="+address, nil)
	recorder := httptest.NewRecorder()

	mockUseCase := new(mocks.GetUserDepositsUseCaseMock)
	mockUseCase.On("Run", mock.Anything, address).Return([]quote.PegoutDeposit{}, nil)

	handlerFunc := handlers.NewGetUserQuotesHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var responseBody []pkg.DepositEventDTO
	err := json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	assert.Empty(t, responseBody)
	assert.NotNil(t, responseBody)

	mockUseCase.AssertExpectations(t)
}

// nolint:funlen
func TestGetUserQuotesHandlerErrorCases(t *testing.T) {

	t.Run("should return 400 when address parameter is missing", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/userQuotes", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetUserDepositsUseCaseMock)

		handlerFunc := handlers.NewGetUserQuotesHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "address parameter is required")

		mockUseCase.AssertNotCalled(t, "Run")
	})

	t.Run("should return 400 when address format is invalid", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/userQuotes?address=invalid-address", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetUserDepositsUseCaseMock)

		handlerFunc := handlers.NewGetUserQuotesHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "invalid request", errorResponse["message"])

		mockUseCase.AssertNotCalled(t, "Run")
	})

	t.Run("should return 500 on unexpected errors", func(t *testing.T) {
		address := "0x1234567890abcdef1234567890abcdef12345678"
		request := httptest.NewRequest(http.MethodGet, "/userQuotes?address="+address, nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetUserDepositsUseCaseMock)
		unexpectedError := errors.New("database error")
		mockUseCase.On("Run", mock.Anything, address).Return([]quote.PegoutDeposit(nil), unexpectedError)

		handlerFunc := handlers.NewGetUserQuotesHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "unknown error", errorResponse["message"])
	})
}

func TestGetUserQuotesHandlerErrorResponseFormat(t *testing.T) {

	t.Run("should set correct content type header on error", func(t *testing.T) {
		address := "0x1234567890abcdef1234567890abcdef12345678"
		request := httptest.NewRequest(http.MethodGet, "/userQuotes?address="+address, nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetUserDepositsUseCaseMock)
		mockUseCase.On("Run", mock.Anything, address).Return([]quote.PegoutDeposit(nil), errors.New("error"))

		handlerFunc := handlers.NewGetUserQuotesHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should include timestamp in error response", func(t *testing.T) {
		address := "0x1234567890abcdef1234567890abcdef12345678"
		request := httptest.NewRequest(http.MethodGet, "/userQuotes?address="+address, nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetUserDepositsUseCaseMock)
		mockUseCase.On("Run", mock.Anything, address).Return([]quote.PegoutDeposit(nil), errors.New("error"))

		handlerFunc := handlers.NewGetUserQuotesHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "timestamp")
		assert.NotZero(t, errorResponse["timestamp"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should set recoverable to true on invalid address error", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/userQuotes?address=invalid-address", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetUserDepositsUseCaseMock)

		handlerFunc := handlers.NewGetUserQuotesHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "recoverable")
		assert.Equal(t, true, errorResponse["recoverable"])
	})
}
