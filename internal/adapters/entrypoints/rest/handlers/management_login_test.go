package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestManagementLoginHandlerHappyPath(t *testing.T) {
	reqBody := pkg.LoginRequest{
		Username: "admin",
		Password: "secretpassword",
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/management/login", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	mockUseCase := new(mocks.LoginUseCaseMock)
	mockUseCase.On("Run", mock.Anything, lp.Credentials{
		Username: "admin",
		Password: "secretpassword",
	}).Return(nil)

	mockSessionManager := new(mocks.SessionManagerMock)
	mockSessionManager.On("CloseSession", mock.Anything, mock.Anything).Return(nil)
	mockSessionManager.On("CreateSession", mock.Anything, mock.Anything).Return(nil)

	handlerFunc := handlers.NewManagementLoginHandler(mockUseCase, mockSessionManager)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	mockUseCase.AssertExpectations(t)
	mockSessionManager.AssertExpectations(t)
}

// nolint:funlen
func TestManagementLoginHandlerErrorCases(t *testing.T) {
	t.Run("should return 400 on malformed JSON", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/management/login", bytes.NewBufferString("{invalid json}"))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.LoginUseCaseMock)
		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewManagementLoginHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		mockUseCase.AssertNotCalled(t, "Run")
		mockSessionManager.AssertNotCalled(t, "CloseSession")
		mockSessionManager.AssertNotCalled(t, "CreateSession")
	})

	t.Run("should return 400 on missing username", func(t *testing.T) {
		reqBody := pkg.LoginRequest{
			Username: "",
			Password: "secretpassword",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/management/login", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.LoginUseCaseMock)
		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewManagementLoginHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var errorResponse rest.ErrorResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse.Details, "Username")

		mockUseCase.AssertNotCalled(t, "Run")
	})

	t.Run("should return 400 on missing password", func(t *testing.T) {
		reqBody := pkg.LoginRequest{
			Username: "admin",
			Password: "",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/management/login", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.LoginUseCaseMock)
		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewManagementLoginHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var errorResponse rest.ErrorResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse.Details, "Password")

		mockUseCase.AssertNotCalled(t, "Run")
	})

	t.Run("should return 401 on BadLoginError", func(t *testing.T) {
		reqBody := pkg.LoginRequest{
			Username: "admin",
			Password: "wrongpassword",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/management/login", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.LoginUseCaseMock)
		mockUseCase.On("Run", mock.Anything, lp.Credentials{
			Username: "admin",
			Password: "wrongpassword",
		}).Return(liquidity_provider.BadLoginError)

		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewManagementLoginHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)

		var errorResponse rest.ErrorResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Equal(t, liquidity_provider.BadLoginError.Error(), errorResponse.Message)
		assert.False(t, errorResponse.Recoverable)

		mockUseCase.AssertExpectations(t)
		mockSessionManager.AssertNotCalled(t, "CloseSession")
		mockSessionManager.AssertNotCalled(t, "CreateSession")
	})

	t.Run("should return 500 on unexpected use case error", func(t *testing.T) {
		reqBody := pkg.LoginRequest{
			Username: "admin",
			Password: "secretpassword",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/management/login", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.LoginUseCaseMock)
		mockUseCase.On("Run", mock.Anything, lp.Credentials{
			Username: "admin",
			Password: "secretpassword",
		}).Return(errors.New("database connection error"))

		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewManagementLoginHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		var errorResponse rest.ErrorResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "session creation error", errorResponse.Message)
		assert.False(t, errorResponse.Recoverable)

		mockUseCase.AssertExpectations(t)
		mockSessionManager.AssertNotCalled(t, "CloseSession")
		mockSessionManager.AssertNotCalled(t, "CreateSession")
	})

	t.Run("should return 500 when CloseSession fails", func(t *testing.T) {
		reqBody := pkg.LoginRequest{
			Username: "admin",
			Password: "secretpassword",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/management/login", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.LoginUseCaseMock)
		mockUseCase.On("Run", mock.Anything, lp.Credentials{
			Username: "admin",
			Password: "secretpassword",
		}).Return(nil)

		mockSessionManager := new(mocks.SessionManagerMock)
		mockSessionManager.On("CloseSession", mock.Anything, mock.Anything).Return(errors.New("session close error"))

		handlerFunc := handlers.NewManagementLoginHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		// The error response is written by SessionManager.CloseSession, so we just verify the handler returned
		mockUseCase.AssertExpectations(t)
		mockSessionManager.AssertExpectations(t)
		mockSessionManager.AssertNotCalled(t, "CreateSession")
	})

	t.Run("should return 500 when CreateSession fails", func(t *testing.T) {
		reqBody := pkg.LoginRequest{
			Username: "admin",
			Password: "secretpassword",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/management/login", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.LoginUseCaseMock)
		mockUseCase.On("Run", mock.Anything, lp.Credentials{
			Username: "admin",
			Password: "secretpassword",
		}).Return(nil)

		mockSessionManager := new(mocks.SessionManagerMock)
		mockSessionManager.On("CloseSession", mock.Anything, mock.Anything).Return(nil)
		mockSessionManager.On("CreateSession", mock.Anything, mock.Anything).Return(errors.New("session creation error"))

		handlerFunc := handlers.NewManagementLoginHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		// The error response is written by SessionManager.CreateSession, so we just verify the handler returned
		mockUseCase.AssertExpectations(t)
		mockSessionManager.AssertExpectations(t)
	})
}

// nolint:funlen
func TestManagementLoginHandlerErrorResponseFormat(t *testing.T) {
	t.Run("error response should have correct Content-Type", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/management/login", bytes.NewBufferString("{invalid}"))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.LoginUseCaseMock)
		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewManagementLoginHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	})

	t.Run("error response should have timestamp", func(t *testing.T) {
		reqBody := pkg.LoginRequest{
			Username: "admin",
			Password: "wrongpassword",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/management/login", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.LoginUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.Anything).Return(liquidity_provider.BadLoginError)

		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewManagementLoginHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse rest.ErrorResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.NotEmpty(t, errorResponse.Timestamp)
	})

	t.Run("401 error should have recoverable flag set to false", func(t *testing.T) {
		reqBody := pkg.LoginRequest{
			Username: "admin",
			Password: "wrongpassword",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/management/login", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.LoginUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.Anything).Return(liquidity_provider.BadLoginError)

		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewManagementLoginHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse rest.ErrorResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.False(t, errorResponse.Recoverable)
	})
}
