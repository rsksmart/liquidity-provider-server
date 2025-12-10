package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestManagementLogoutHandlerHappyPath(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/management/logout", nil)
	recorder := httptest.NewRecorder()

	mockSessionManager := new(mocks.SessionManagerMock)
	mockSessionManager.On("CloseSession", mock.Anything, mock.Anything).Return(nil)

	handlerFunc := handlers.NewManagementLogoutHandler(mockSessionManager)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	assert.Empty(t, recorder.Body.String())

	mockSessionManager.AssertExpectations(t)
}

func TestManagementLogoutHandlerErrorCases(t *testing.T) {
	t.Run("should return error when CloseSession fails", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/management/logout", nil)
		recorder := httptest.NewRecorder()

		mockSessionManager := new(mocks.SessionManagerMock)
		mockSessionManager.On("CloseSession", mock.Anything, mock.Anything).Return(errors.New("session close error"))

		handlerFunc := handlers.NewManagementLogoutHandler(mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		// The error response is written by SessionManager.CloseSession
		// We verify the handler returned early and CloseSession was called
		mockSessionManager.AssertExpectations(t)
	})
}

func TestManagementLogoutHandlerResponseFormat(t *testing.T) {
	t.Run("successful logout should return empty body", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/management/logout", nil)
		recorder := httptest.NewRecorder()

		mockSessionManager := new(mocks.SessionManagerMock)
		mockSessionManager.On("CloseSession", mock.Anything, mock.Anything).Return(nil)

		handlerFunc := handlers.NewManagementLogoutHandler(mockSessionManager)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusNoContent, recorder.Code)
		assert.Empty(t, recorder.Body.String())
	})
}
