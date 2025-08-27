package middlewares_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/middlewares"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSessionMiddleware_HappyPath(t *testing.T) {
	mockStore := mocks.NewStoreMock(t)
	mockHandler := mocks.NewHandlerMock(t)

	// Create a valid, existing session (not new) using New method first to get proper session
	// This ensures the session has the proper store reference internally
	validSession := sessions.NewSession(mockStore, cookies.ManagementSessionCookieName)
	validSession.IsNew = false // Mark as existing session

	mockStore.On("Get", mock.AnythingOfType("*http.Request"), cookies.ManagementSessionCookieName).Return(validSession, nil)
	mockStore.On("Save", mock.AnythingOfType("*http.Request"), mock.AnythingOfType("*httptest.ResponseRecorder"), validSession).Return(nil)
	// Setup handler mock to write a success response when called
	mockHandler.On("ServeHTTP", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		w, ok := args.Get(0).(http.ResponseWriter)
		if !ok {
			t.Errorf("Expected http.ResponseWriter, got %T", args.Get(0))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("success"))
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	})

	// Create the middleware
	middleware := middlewares.NewSessionMiddlewares(mockManagementEnv(), mockStore)
	handler := middleware.SessionValidator(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "success", rr.Body.String())

	mockStore.AssertExpectations(t)
	mockHandler.AssertExpectations(t)

	mockStore.AssertCalled(t, "Get", req, cookies.ManagementSessionCookieName)
	mockStore.AssertCalled(t, "Save", req, rr, validSession)
	mockHandler.AssertCalled(t, "ServeHTTP", rr, req)
}

func TestSessionMiddleware_StoreGetError_ReturnsForbidden(t *testing.T) {
	mockStore := mocks.NewStoreMock(t)
	mockHandler := mocks.NewHandlerMock(t)

	// Setup mock to return an error from store.Get
	mockStore.On("Get", mock.AnythingOfType("*http.Request"), cookies.ManagementSessionCookieName).Return((*sessions.Session)(nil), assert.AnError)

	// Handler should NOT be called in this case, so no expectations set

	middleware := middlewares.NewSessionMiddlewares(mockManagementEnv(), mockStore)
	handler := middleware.SessionValidator(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)

	// Verify the error response structure
	var errorResponse rest.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.Equal(t, "session validation error", errorResponse.Message)
	assert.False(t, errorResponse.Recoverable)
	assert.Contains(t, errorResponse.Details["error"], "assert.AnError")

	mockStore.AssertExpectations(t)
	mockHandler.AssertExpectations(t)

	mockHandler.AssertNotCalled(t, "ServeHTTP")
}

func TestSessionMiddleware_NewSession_ReturnsForbidden(t *testing.T) {
	mockStore := mocks.NewStoreMock(t)
	mockHandler := mocks.NewHandlerMock(t)

	// Create a new session (IsNew = true) using proper method
	newSession := sessions.NewSession(mockStore, cookies.ManagementSessionCookieName)
	newSession.IsNew = true // This is the default, but being explicit

	mockStore.On("Get", mock.AnythingOfType("*http.Request"), cookies.ManagementSessionCookieName).Return(newSession, nil)

	// Handler should NOT be called in this case, so no expectations set

	middleware := middlewares.NewSessionMiddlewares(mockManagementEnv(), mockStore)
	handler := middleware.SessionValidator(mockHandler)

	// Create test request and response
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)

	// Verify the error response structure
	var errorResponse rest.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.Equal(t, "session not recognized", errorResponse.Message)
	assert.True(t, errorResponse.Recoverable)
	assert.Empty(t, errorResponse.Details) // This error doesn't include details

	mockStore.AssertExpectations(t)
	mockHandler.AssertExpectations(t)

	mockHandler.AssertNotCalled(t, "ServeHTTP")
}

func TestSessionMiddleware_SessionSaveError_StillProcessesRequest(t *testing.T) {
	mockStore := mocks.NewStoreMock(t)
	mockHandler := mocks.NewHandlerMock(t)

	// Create a valid, existing session
	validSession := sessions.NewSession(mockStore, cookies.ManagementSessionCookieName)
	validSession.IsNew = false

	// Save returns an error but request should still succeed
	mockStore.On("Get", mock.AnythingOfType("*http.Request"), cookies.ManagementSessionCookieName).Return(validSession, nil)
	mockStore.On("Save", mock.AnythingOfType("*http.Request"), mock.AnythingOfType("*httptest.ResponseRecorder"), validSession).Return(assert.AnError)
	mockHandler.On("ServeHTTP", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		w, ok := args.Get(0).(http.ResponseWriter)
		if !ok {
			t.Errorf("Expected http.ResponseWriter, got %T", args.Get(0))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("success"))
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	})

	middleware := middlewares.NewSessionMiddlewares(mockManagementEnv(), mockStore)
	handler := middleware.SessionValidator(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Even though session.Save() failed, the request should still succeed
	// because the error is only logged, not returned
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "success", rr.Body.String())

	mockStore.AssertExpectations(t)
	mockHandler.AssertExpectations(t)

	mockStore.AssertCalled(t, "Get", req, cookies.ManagementSessionCookieName)
	mockStore.AssertCalled(t, "Save", req, rr, validSession)
	mockHandler.AssertCalled(t, "ServeHTTP", rr, req)
}

// Helper function to create a mock management environment
func mockManagementEnv() environment.ManagementEnv {
	return environment.ManagementEnv{
		EnableManagementApi:   true,
		SessionAuthKey:        "01fbac02d66202e8468d2a4f1deba4fa5c2491f592e0e22e32fe1e6acac25923",
		SessionEncryptionKey:  "02fbac02d66202e8468d2a4f1deba4fa5c2491f592e0e22e32fe1e6acac25923",
		SessionTokenAuthKey:   "03fbac02d66202e8468d2a4f1deba4fa5c2491f592e0e22e32fe1e6acac25923",
		UseHttps:              false,
		EnableSecurityHeaders: false,
	}
}
