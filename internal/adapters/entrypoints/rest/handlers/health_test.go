package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHealthCheckHandlerAllServicesOk(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	expectedResult := usecases.HealthStatus{
		Status: usecases.SvcStatusOk,
		Services: usecases.Services{
			Db:  usecases.SvcStatusOk,
			Rsk: usecases.SvcStatusOk,
			Btc: usecases.SvcStatusOk,
		},
	}

	mockUseCase := new(mocks.HealthUseCaseMock)
	mockUseCase.On("Run", mock.Anything).Return(expectedResult)

	handlerFunc := handlers.NewHealthCheckHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody pkg.HealthResponse
	err := json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	assert.Equal(t, "ok", responseBody.Status)
	assert.Equal(t, "ok", responseBody.Services.Db)
	assert.Equal(t, "ok", responseBody.Services.Rsk)
	assert.Equal(t, "ok", responseBody.Services.Btc)

	mockUseCase.AssertExpectations(t)
}

func TestHealthCheckHandlerDegradedStatus(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	expectedResult := usecases.HealthStatus{
		Status: usecases.SvcStatusDegraded,
		Services: usecases.Services{
			Db:  usecases.SvcStatusUnreachable,
			Rsk: usecases.SvcStatusOk,
			Btc: usecases.SvcStatusOk,
		},
	}

	mockUseCase := new(mocks.HealthUseCaseMock)
	mockUseCase.On("Run", mock.Anything).Return(expectedResult)

	handlerFunc := handlers.NewHealthCheckHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var responseBody pkg.HealthResponse
	err := json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	assert.Equal(t, "degraded", responseBody.Status)
	assert.Equal(t, "unreachable", responseBody.Services.Db)
	assert.Equal(t, "ok", responseBody.Services.Rsk)
	assert.Equal(t, "ok", responseBody.Services.Btc)

	mockUseCase.AssertExpectations(t)
}

func TestHealthCheckHandlerMultipleServicesUnreachable(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	expectedResult := usecases.HealthStatus{
		Status: usecases.SvcStatusDegraded,
		Services: usecases.Services{
			Db:  usecases.SvcStatusUnreachable,
			Rsk: usecases.SvcStatusUnreachable,
			Btc: usecases.SvcStatusUnreachable,
		},
	}

	mockUseCase := new(mocks.HealthUseCaseMock)
	mockUseCase.On("Run", mock.Anything).Return(expectedResult)

	handlerFunc := handlers.NewHealthCheckHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var responseBody pkg.HealthResponse
	err := json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	assert.Equal(t, "degraded", responseBody.Status)
	assert.Equal(t, "unreachable", responseBody.Services.Db)
	assert.Equal(t, "unreachable", responseBody.Services.Rsk)
	assert.Equal(t, "unreachable", responseBody.Services.Btc)

	mockUseCase.AssertExpectations(t)
}

func TestHealthCheckHandlerResponseStructure(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	expectedResult := usecases.HealthStatus{
		Status: usecases.SvcStatusOk,
		Services: usecases.Services{
			Db:  usecases.SvcStatusOk,
			Rsk: usecases.SvcStatusOk,
			Btc: usecases.SvcStatusOk,
		},
	}

	mockUseCase := new(mocks.HealthUseCaseMock)
	mockUseCase.On("Run", mock.Anything).Return(expectedResult)

	handlerFunc := handlers.NewHealthCheckHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	var rawResponse map[string]interface{}
	err := json.NewDecoder(recorder.Body).Decode(&rawResponse)
	require.NoError(t, err)

	assert.Len(t, rawResponse, 2)
	assert.Contains(t, rawResponse, "status")
	assert.Contains(t, rawResponse, "services")

	services, ok := rawResponse["services"].(map[string]interface{})
	require.True(t, ok)
	assert.Len(t, services, 3)
	assert.Contains(t, services, "db")
	assert.Contains(t, services, "rsk")
	assert.Contains(t, services, "btc")

	mockUseCase.AssertExpectations(t)
}
