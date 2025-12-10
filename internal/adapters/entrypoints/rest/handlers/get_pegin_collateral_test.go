package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPeginCollateralHandlerHappyPath(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/pegin/collateral", nil)
	recorder := httptest.NewRecorder()

	expectedCollateral := entities.NewWei(5000000000000000000)

	mockUseCase := new(mocks.GetPeginCollateralUseCaseMock)
	mockUseCase.On("Run").Return(expectedCollateral, nil)

	handlerFunc := handlers.NewGetPeginCollateralHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody pkg.GetCollateralResponse
	err := json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	assert.Equal(t, expectedCollateral.AsBigInt(), responseBody.Collateral)

	mockUseCase.AssertExpectations(t)
}

func TestGetPeginCollateralHandlerErrorCases(t *testing.T) {

	t.Run("should return 500 on unexpected errors", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/pegin/collateral", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginCollateralUseCaseMock)
		unexpectedError := errors.New("unexpected blockchain error")
		mockUseCase.On("Run").Return((*entities.Wei)(nil), unexpectedError)

		handlerFunc := handlers.NewGetPeginCollateralHandler(mockUseCase)
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

func TestGetPeginCollateralHandlerErrorResponseFormat(t *testing.T) {

	t.Run("should set correct content type header on error", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/pegin/collateral", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginCollateralUseCaseMock)
		mockUseCase.On("Run").Return((*entities.Wei)(nil), errors.New("error"))

		handlerFunc := handlers.NewGetPeginCollateralHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should include timestamp in error response", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/pegin/collateral", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginCollateralUseCaseMock)
		mockUseCase.On("Run").Return((*entities.Wei)(nil), errors.New("error"))

		handlerFunc := handlers.NewGetPeginCollateralHandler(mockUseCase)
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
		request := httptest.NewRequest(http.MethodGet, "/pegin/collateral", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetPeginCollateralUseCaseMock)
		mockUseCase.On("Run").Return((*entities.Wei)(nil), errors.New("error"))

		handlerFunc := handlers.NewGetPeginCollateralHandler(mockUseCase)
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
