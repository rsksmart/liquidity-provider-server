package middlewares_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/middlewares"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSessionMiddleware_HappyPath(t *testing.T) {
	mockStore := mocks.NewSessionStoreMock(t)
	mockHandler := mocks.NewHandlerMock(t)

	mockStore.On("Refresh", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Return(nil)
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

	handler := middlewares.NewSessionMiddleware(mockStore)(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "success", rr.Body.String())

	mockStore.AssertExpectations(t)
	mockHandler.AssertExpectations(t)
	mockStore.AssertCalled(t, "Refresh", rr, req)
	mockHandler.AssertCalled(t, "ServeHTTP", rr, req)
}

func TestSessionMiddleware_RefreshError_ReturnsForbidden(t *testing.T) {
	mockStore := mocks.NewSessionStoreMock(t)
	mockHandler := mocks.NewHandlerMock(t)

	mockStore.On("Refresh", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Return(cookies.ErrSessionNotRecognized)

	handler := middlewares.NewSessionMiddleware(mockStore)(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)

	var errorResponse rest.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.Equal(t, "session not recognized", errorResponse.Message)
	assert.True(t, errorResponse.Recoverable)

	mockStore.AssertExpectations(t)
	mockHandler.AssertNotCalled(t, "ServeHTTP")
}
