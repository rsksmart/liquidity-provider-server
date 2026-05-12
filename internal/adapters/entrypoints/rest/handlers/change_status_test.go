package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangeStatusHandlerHappyPath(t *testing.T) {
	status := true
	reqBody := pkg.ChangeStatusRequest{
		Status: &status,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/providers/changeStatus", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	mockUseCase := new(mocks.ChangeStatusUseCaseMock)
	mockUseCase.On("Run", true).Return(nil)

	handlerFunc := handlers.NewChangeStatusHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	assert.Empty(t, recorder.Body.String())

	mockUseCase.AssertExpectations(t)
}

func TestChangeStatusHandlerHappyPathSetToFalse(t *testing.T) {
	status := false
	reqBody := pkg.ChangeStatusRequest{
		Status: &status,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/providers/changeStatus", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	mockUseCase := new(mocks.ChangeStatusUseCaseMock)
	mockUseCase.On("Run", false).Return(nil)

	handlerFunc := handlers.NewChangeStatusHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	assert.Empty(t, recorder.Body.String())

	mockUseCase.AssertExpectations(t)
}

// nolint:funlen
func TestChangeStatusHandlerErrorCases(t *testing.T) {

	t.Run("should handle malformed JSON in request body", func(t *testing.T) {
		malformedJSON := []byte(`{"status": "not-a-boolean"`)
		request := httptest.NewRequest(http.MethodPost, "/providers/changeStatus", bytes.NewBuffer(malformedJSON))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.ChangeStatusUseCaseMock)

		handlerFunc := handlers.NewChangeStatusHandler(mockUseCase)
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

	t.Run("should return 500 on unexpected errors", func(t *testing.T) {
		status := true
		reqBody := pkg.ChangeStatusRequest{
			Status: &status,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/providers/changeStatus", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.ChangeStatusUseCaseMock)
		unexpectedError := errors.New("unexpected blockchain error")
		mockUseCase.On("Run", true).Return(unexpectedError)

		handlerFunc := handlers.NewChangeStatusHandler(mockUseCase)
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
func TestChangeStatusHandlerErrorResponseFormat(t *testing.T) {

	t.Run("should set correct content type header on error", func(t *testing.T) {
		status := true
		reqBody := pkg.ChangeStatusRequest{
			Status: &status,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/providers/changeStatus", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.ChangeStatusUseCaseMock)
		mockUseCase.On("Run", true).Return(errors.New("error"))

		handlerFunc := handlers.NewChangeStatusHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should include timestamp in error response", func(t *testing.T) {
		status := true
		reqBody := pkg.ChangeStatusRequest{
			Status: &status,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/providers/changeStatus", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.ChangeStatusUseCaseMock)
		mockUseCase.On("Run", true).Return(errors.New("error"))

		handlerFunc := handlers.NewChangeStatusHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "timestamp")
		assert.NotZero(t, errorResponse["timestamp"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should set recoverable to false on errors", func(t *testing.T) {
		status := true
		reqBody := pkg.ChangeStatusRequest{
			Status: &status,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/providers/changeStatus", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.ChangeStatusUseCaseMock)
		mockUseCase.On("Run", true).Return(errors.New("error"))

		handlerFunc := handlers.NewChangeStatusHandler(mockUseCase)
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
