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
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSetCredentialsHandlerHappyPath(t *testing.T) {
	reqBody := pkg.CredentialsUpdateRequest{
		OldUsername: "admin",
		OldPassword: "oldpassword",
		NewUsername: "newadmin",
		NewPassword: "newpassword123!",
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/management/credentials", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	mockUseCase := new(mocks.SetCredentialsUseCaseMock)
	mockUseCase.On("Run", mock.Anything,
		lp.Credentials{Username: "admin", Password: "oldpassword"},
		lp.Credentials{Username: "newadmin", Password: "newpassword123!"},
	).Return(nil)

	mockSessionManager := new(mocks.SessionManagerMock)
	mockSessionManager.On("CloseSession", mock.Anything, mock.Anything).Return(nil)

	handlerFunc := handlers.NewSetCredentialsHandler(mockUseCase, mockSessionManager)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	mockUseCase.AssertExpectations(t)
	mockSessionManager.AssertExpectations(t)
}

// nolint:funlen
func TestSetCredentialsHandlerErrorCases(t *testing.T) {
	t.Run("should return 400 on malformed JSON", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/management/credentials", bytes.NewBufferString("{invalid json}"))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.SetCredentialsUseCaseMock)
		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewSetCredentialsHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		mockUseCase.AssertNotCalled(t, "Run")
		mockSessionManager.AssertNotCalled(t, "CloseSession")
	})

	t.Run("should return 400 on missing oldUsername", func(t *testing.T) {
		reqBody := pkg.CredentialsUpdateRequest{
			OldUsername: "",
			OldPassword: "oldpassword",
			NewUsername: "newadmin",
			NewPassword: "newpassword123!",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/management/credentials", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.SetCredentialsUseCaseMock)
		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewSetCredentialsHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var errorResponse rest.ErrorResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse.Details, "OldUsername")

		mockUseCase.AssertNotCalled(t, "Run")
	})

	t.Run("should return 400 on missing newPassword", func(t *testing.T) {
		reqBody := pkg.CredentialsUpdateRequest{
			OldUsername: "admin",
			OldPassword: "oldpassword",
			NewUsername: "newadmin",
			NewPassword: "",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/management/credentials", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.SetCredentialsUseCaseMock)
		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewSetCredentialsHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var errorResponse rest.ErrorResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse.Details, "NewPassword")

		mockUseCase.AssertNotCalled(t, "Run")
	})

	t.Run("should return 401 on BadLoginError", func(t *testing.T) {
		reqBody := pkg.CredentialsUpdateRequest{
			OldUsername: "admin",
			OldPassword: "wrongpassword",
			NewUsername: "newadmin",
			NewPassword: "newpassword123!",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/management/credentials", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.SetCredentialsUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(liquidity_provider.BadLoginError)

		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewSetCredentialsHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)

		var errorResponse rest.ErrorResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Equal(t, liquidity_provider.BadLoginError.Error(), errorResponse.Message)
		assert.True(t, errorResponse.Recoverable)

		mockUseCase.AssertExpectations(t)
		mockSessionManager.AssertNotCalled(t, "CloseSession")
	})

	t.Run("should return 400 on PasswordComplexityError", func(t *testing.T) {
		reqBody := pkg.CredentialsUpdateRequest{
			OldUsername: "admin",
			OldPassword: "oldpassword",
			NewUsername: "newadmin",
			NewPassword: "weak",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/management/credentials", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.SetCredentialsUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(utils.PasswordComplexityError)

		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewSetCredentialsHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var errorResponse rest.ErrorResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Equal(t, utils.PasswordComplexityError.Error(), errorResponse.Message)
		assert.True(t, errorResponse.Recoverable)

		mockUseCase.AssertExpectations(t)
		mockSessionManager.AssertNotCalled(t, "CloseSession")
	})

	t.Run("should return 500 on unexpected use case error", func(t *testing.T) {
		reqBody := pkg.CredentialsUpdateRequest{
			OldUsername: "admin",
			OldPassword: "oldpassword",
			NewUsername: "newadmin",
			NewPassword: "newpassword123!",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/management/credentials", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.SetCredentialsUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("database error"))

		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewSetCredentialsHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		var errorResponse rest.ErrorResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "unexpected login error", errorResponse.Message)
		assert.False(t, errorResponse.Recoverable)

		mockUseCase.AssertExpectations(t)
		mockSessionManager.AssertNotCalled(t, "CloseSession")
	})

	t.Run("should return error when CloseSession fails", func(t *testing.T) {
		reqBody := pkg.CredentialsUpdateRequest{
			OldUsername: "admin",
			OldPassword: "oldpassword",
			NewUsername: "newadmin",
			NewPassword: "newpassword123!",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/management/credentials", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.SetCredentialsUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		mockSessionManager := new(mocks.SessionManagerMock)
		mockSessionManager.On("CloseSession", mock.Anything, mock.Anything).Return(errors.New("session close error"))

		handlerFunc := handlers.NewSetCredentialsHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		// The error response is written by SessionManager.CloseSession
		mockUseCase.AssertExpectations(t)
		mockSessionManager.AssertExpectations(t)
	})
}

func TestSetCredentialsHandlerErrorResponseFormat(t *testing.T) {
	t.Run("error response should have correct Content-Type", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/management/credentials", bytes.NewBufferString("{invalid}"))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.SetCredentialsUseCaseMock)
		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewSetCredentialsHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	})

	t.Run("error response should have timestamp", func(t *testing.T) {
		reqBody := pkg.CredentialsUpdateRequest{
			OldUsername: "admin",
			OldPassword: "wrongpassword",
			NewUsername: "newadmin",
			NewPassword: "newpassword123!",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/management/credentials", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.SetCredentialsUseCaseMock)
		mockUseCase.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(liquidity_provider.BadLoginError)

		mockSessionManager := new(mocks.SessionManagerMock)

		handlerFunc := handlers.NewSetCredentialsHandler(mockUseCase, mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse rest.ErrorResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.NotZero(t, errorResponse.Timestamp)
	})
}
