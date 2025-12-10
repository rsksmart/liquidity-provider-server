package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResignationHandlerHappyPath(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/providers/resignation", nil)
	recorder := httptest.NewRecorder()

	mockUseCase := new(mocks.ResignUseCaseMock)
	mockUseCase.On("Run").Return(nil)

	handlerFunc := handlers.NewResignationHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	assert.Empty(t, recorder.Body.String())

	mockUseCase.AssertExpectations(t)
}

func TestResignationHandlerErrorCases(t *testing.T) {

	t.Run("should return 500 on unexpected errors", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/providers/resignation", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.ResignUseCaseMock)
		unexpectedError := errors.New("unexpected blockchain error")
		mockUseCase.On("Run").Return(unexpectedError)

		handlerFunc := handlers.NewResignationHandler(mockUseCase)
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

func TestResignationHandlerErrorResponseFormat(t *testing.T) {

	t.Run("should set correct content type header on error", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/providers/resignation", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.ResignUseCaseMock)
		mockUseCase.On("Run").Return(errors.New("error"))

		handlerFunc := handlers.NewResignationHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should include timestamp in error response", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/providers/resignation", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.ResignUseCaseMock)
		mockUseCase.On("Run").Return(errors.New("error"))

		handlerFunc := handlers.NewResignationHandler(mockUseCase)
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
		request := httptest.NewRequest(http.MethodPost, "/providers/resignation", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.ResignUseCaseMock)
		mockUseCase.On("Run").Return(errors.New("error"))

		handlerFunc := handlers.NewResignationHandler(mockUseCase)
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
