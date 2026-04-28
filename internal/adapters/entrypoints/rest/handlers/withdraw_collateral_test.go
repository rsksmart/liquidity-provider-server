package handlers_test

import (
	"encoding/json"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithdrawCollateralHandlerHappyPath(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/providers/withdrawCollateral", nil)
	recorder := httptest.NewRecorder()

	mockUseCase := new(mocks.WithdrawCollateralUseCaseMock)
	mockUseCase.On("Run").Return(nil)

	handlerFunc := handlers.NewWithdrawCollateralHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	assert.Empty(t, recorder.Body.String())

	mockUseCase.AssertExpectations(t)
}

// nolint:funlen
func TestWithdrawCollateralHandlerErrorCases(t *testing.T) {

	t.Run("should return 409 on ProviderNotResignedError", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/providers/withdrawCollateral", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.WithdrawCollateralUseCaseMock)
		mockUseCase.On("Run").Return(liquidity_provider.ProviderNotResignedError)

		handlerFunc := handlers.NewWithdrawCollateralHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusConflict, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "unknown error", errorResponse["message"])
	})

	t.Run("should return 500 on unexpected errors", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/providers/withdrawCollateral", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.WithdrawCollateralUseCaseMock)
		unexpectedError := errors.New("unexpected blockchain error")
		mockUseCase.On("Run").Return(unexpectedError)

		handlerFunc := handlers.NewWithdrawCollateralHandler(mockUseCase)
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

// nolint:funlen
func TestWithdrawCollateralHandlerErrorResponseFormat(t *testing.T) {

	t.Run("should set correct content type header on error", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/providers/withdrawCollateral", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.WithdrawCollateralUseCaseMock)
		mockUseCase.On("Run").Return(liquidity_provider.ProviderNotResignedError)

		handlerFunc := handlers.NewWithdrawCollateralHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should include timestamp in error response", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/providers/withdrawCollateral", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.WithdrawCollateralUseCaseMock)
		mockUseCase.On("Run").Return(liquidity_provider.ProviderNotResignedError)

		handlerFunc := handlers.NewWithdrawCollateralHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "timestamp")
		assert.NotZero(t, errorResponse["timestamp"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should set recoverable to true for ProviderNotResignedError", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/providers/withdrawCollateral", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.WithdrawCollateralUseCaseMock)
		mockUseCase.On("Run").Return(liquidity_provider.ProviderNotResignedError)

		handlerFunc := handlers.NewWithdrawCollateralHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "recoverable")
		assert.Equal(t, true, errorResponse["recoverable"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should set recoverable to false on unexpected errors", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/providers/withdrawCollateral", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.WithdrawCollateralUseCaseMock)
		unexpectedError := errors.New("unexpected error")
		mockUseCase.On("Run").Return(unexpectedError)

		handlerFunc := handlers.NewWithdrawCollateralHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "recoverable")
		assert.Equal(t, false, errorResponse["recoverable"])
		mockUseCase.AssertExpectations(t)
	})
}
