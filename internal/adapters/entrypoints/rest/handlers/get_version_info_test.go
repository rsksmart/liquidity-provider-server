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

func TestVersionInfoHandlerHappyPath(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/version", nil)
	recorder := httptest.NewRecorder()

	expectedResult := liquidity_provider.ServerInfo{
		Version:  "v1.2.3",
		Revision: "abc123def456",
	}

	mockUseCase := new(mocks.ServerInfoUseCaseMock)
	mockUseCase.On("Run").Return(expectedResult, nil)

	handlerFunc := handlers.NewVersionInfoHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody pkg.ServerInfoDTO
	err := json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	assert.Equal(t, "v1.2.3", responseBody.Version)
	assert.Equal(t, "abc123def456", responseBody.Revision)

	mockUseCase.AssertExpectations(t)
}

func TestVersionInfoHandlerErrorCases(t *testing.T) {

	t.Run("should return 500 when unable to read build info", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/version", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.ServerInfoUseCaseMock)
		mockUseCase.On("Run").Return(liquidity_provider.ServerInfo{}, errors.New("unable to read build info"))

		handlerFunc := handlers.NewVersionInfoHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "unable to read build info", errorResponse["message"])
	})
}

func TestVersionInfoHandlerErrorResponseFormat(t *testing.T) {

	t.Run("should set correct content type header on error", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/version", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.ServerInfoUseCaseMock)
		mockUseCase.On("Run").Return(liquidity_provider.ServerInfo{}, errors.New("error"))

		handlerFunc := handlers.NewVersionInfoHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should set recoverable to false on errors", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/version", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.ServerInfoUseCaseMock)
		mockUseCase.On("Run").Return(liquidity_provider.ServerInfo{}, errors.New("error"))

		handlerFunc := handlers.NewVersionInfoHandler(mockUseCase)
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
