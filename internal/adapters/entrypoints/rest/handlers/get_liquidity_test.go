package handlers_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	lpUseCase "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetAvailableLiquidityHandlerHappyPath(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/providers/liquidity", nil)
	recorder := httptest.NewRecorder()

	peginLiquidity := entities.NewWei(5000000000000000000)
	pegoutLiquidity := entities.NewWei(3000000000000000000)

	expectedResult := liquidity_provider.AvailableLiquidity{
		PeginLiquidity:  peginLiquidity,
		PegoutLiquidity: pegoutLiquidity,
	}

	mockUseCase := new(mocks.GetAvailableLiquidityUseCaseMock)
	mockUseCase.On("Run", mock.Anything).Return(expectedResult, nil)

	handlerFunc := handlers.NewGetAvailableLiquidityHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody pkg.AvailableLiquidityDTO
	err := json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	assert.Equal(t, big.NewInt(5000000000000000000), responseBody.PeginLiquidityAmount)
	assert.Equal(t, big.NewInt(3000000000000000000), responseBody.PegoutLiquidityAmount)

	mockUseCase.AssertExpectations(t)
}

func TestGetAvailableLiquidityHandlerErrorCases(t *testing.T) {

	t.Run("should return 403 when liquidity check is not enabled", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/providers/liquidity", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetAvailableLiquidityUseCaseMock)
		mockUseCase.On("Run", mock.Anything).Return(liquidity_provider.AvailableLiquidity{}, lpUseCase.LiquidityCheckNotEnabledError)

		handlerFunc := handlers.NewGetAvailableLiquidityHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusForbidden, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "public liquidity check is not enabled", errorResponse["message"])
	})

	t.Run("should return 500 on unexpected errors", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/providers/liquidity", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetAvailableLiquidityUseCaseMock)
		unexpectedError := errors.New("unexpected blockchain error")
		mockUseCase.On("Run", mock.Anything).Return(liquidity_provider.AvailableLiquidity{}, unexpectedError)

		handlerFunc := handlers.NewGetAvailableLiquidityHandler(mockUseCase)
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
func TestGetAvailableLiquidityHandlerErrorResponseFormat(t *testing.T) {

	t.Run("should set correct content type header on error", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/providers/liquidity", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetAvailableLiquidityUseCaseMock)
		mockUseCase.On("Run", mock.Anything).Return(liquidity_provider.AvailableLiquidity{}, errors.New("error"))

		handlerFunc := handlers.NewGetAvailableLiquidityHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should include timestamp in error response", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/providers/liquidity", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetAvailableLiquidityUseCaseMock)
		mockUseCase.On("Run", mock.Anything).Return(liquidity_provider.AvailableLiquidity{}, errors.New("error"))

		handlerFunc := handlers.NewGetAvailableLiquidityHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "timestamp")
		assert.NotZero(t, errorResponse["timestamp"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should set recoverable to false on errors", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/providers/liquidity", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetAvailableLiquidityUseCaseMock)
		mockUseCase.On("Run", mock.Anything).Return(liquidity_provider.AvailableLiquidity{}, errors.New("error"))

		handlerFunc := handlers.NewGetAvailableLiquidityHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "recoverable")
		assert.Equal(t, false, errorResponse["recoverable"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should set recoverable to false on forbidden error", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/providers/liquidity", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetAvailableLiquidityUseCaseMock)
		mockUseCase.On("Run", mock.Anything).Return(liquidity_provider.AvailableLiquidity{}, lpUseCase.LiquidityCheckNotEnabledError)

		handlerFunc := handlers.NewGetAvailableLiquidityHandler(mockUseCase)
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
