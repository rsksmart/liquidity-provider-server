package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProvidersHandlerHappyPath(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/getProviders", nil)
	recorder := httptest.NewRecorder()

	expectedProviders := []liquidity_provider.RegisteredLiquidityProvider{
		{
			Id:           1,
			Address:      "0x1234567890abcdef1234567890abcdef12345678",
			Name:         "Provider One",
			ApiBaseUrl:   "https://provider1.example.com",
			Status:       true,
			ProviderType: liquidity_provider.FullProvider,
		},
		{
			Id:           2,
			Address:      "0xabcdef1234567890abcdef1234567890abcdef12",
			Name:         "Provider Two",
			ApiBaseUrl:   "https://provider2.example.com",
			Status:       false,
			ProviderType: liquidity_provider.PeginProvider,
		},
	}

	mockUseCase := new(mocks.GetProvidersUseCaseMock)
	mockUseCase.On("Run").Return(expectedProviders, nil)

	handlerFunc := handlers.NewGetProvidersHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody []pkg.LiquidityProvider
	err := json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	assert.Len(t, responseBody, 2)

	assert.Equal(t, uint64(1), responseBody[0].Id)
	assert.Equal(t, "0x1234567890abcdef1234567890abcdef12345678", responseBody[0].Provider)
	assert.Equal(t, "Provider One", responseBody[0].Name)
	assert.Equal(t, "https://provider1.example.com", responseBody[0].ApiBaseUrl)
	assert.True(t, responseBody[0].Status)
	assert.Equal(t, "both", responseBody[0].ProviderType)

	assert.Equal(t, uint64(2), responseBody[1].Id)
	assert.Equal(t, "0xabcdef1234567890abcdef1234567890abcdef12", responseBody[1].Provider)
	assert.Equal(t, "Provider Two", responseBody[1].Name)
	assert.Equal(t, "https://provider2.example.com", responseBody[1].ApiBaseUrl)
	assert.False(t, responseBody[1].Status)
	assert.Equal(t, "pegin", responseBody[1].ProviderType)

	mockUseCase.AssertExpectations(t)
}

func TestGetProvidersHandlerEmptyList(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/getProviders", nil)
	recorder := httptest.NewRecorder()

	expectedProviders := []liquidity_provider.RegisteredLiquidityProvider{}

	mockUseCase := new(mocks.GetProvidersUseCaseMock)
	mockUseCase.On("Run").Return(expectedProviders, nil)

	handlerFunc := handlers.NewGetProvidersHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody []pkg.LiquidityProvider
	err := json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	assert.Empty(t, responseBody)
	assert.NotNil(t, responseBody)

	mockUseCase.AssertExpectations(t)
}

func TestGetProvidersHandlerErrorCases(t *testing.T) {

	t.Run("should return 500 on unexpected errors", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/getProviders", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetProvidersUseCaseMock)
		unexpectedError := errors.New("unexpected blockchain error")
		mockUseCase.On("Run").Return([]liquidity_provider.RegisteredLiquidityProvider(nil), unexpectedError)

		handlerFunc := handlers.NewGetProvidersHandler(mockUseCase)
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

func TestGetProvidersHandlerErrorResponseFormat(t *testing.T) {

	t.Run("should set correct content type header on error", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/getProviders", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetProvidersUseCaseMock)
		mockUseCase.On("Run").Return([]liquidity_provider.RegisteredLiquidityProvider(nil), errors.New("error"))

		handlerFunc := handlers.NewGetProvidersHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should include timestamp in error response", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/getProviders", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetProvidersUseCaseMock)
		mockUseCase.On("Run").Return([]liquidity_provider.RegisteredLiquidityProvider(nil), errors.New("error"))

		handlerFunc := handlers.NewGetProvidersHandler(mockUseCase)
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
		request := httptest.NewRequest(http.MethodGet, "/getProviders", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetProvidersUseCaseMock)
		mockUseCase.On("Run").Return([]liquidity_provider.RegisteredLiquidityProvider(nil), errors.New("error"))

		handlerFunc := handlers.NewGetProvidersHandler(mockUseCase)
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
