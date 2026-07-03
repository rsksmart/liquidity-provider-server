package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/middlewares"
	"github.com/stretchr/testify/assert"
)

func okHandler(t *testing.T) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("ok"))
		assert.NoError(t, err)
	})
}

func TestNewCrossOriginProtectionMiddleware(t *testing.T) {
	handler := middlewares.NewCrossOriginProtectionMiddleware()(okHandler(t))

	t.Run("allows unsafe request with no Origin nor Sec-Fetch-Site (same-origin or non-browser)", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/test", nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "ok", response.Body.String())
	})

	t.Run("rejects unsafe cross-site request", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/test", nil)
		request.Header.Set("Origin", "https://evil.example.com")
		request.Header.Set("Sec-Fetch-Site", "cross-site")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
		assert.Contains(t, response.Body.String(), "cross-origin request rejected")
	})

	t.Run("allows safe-method cross-site request", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/test", nil)
		request.Header.Set("Origin", "https://evil.example.com")
		request.Header.Set("Sec-Fetch-Site", "cross-site")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("rejects unsafe cross-site request from any origin", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/test", nil)
		request.Header.Set("Origin", "https://allowed.com")
		request.Header.Set("Sec-Fetch-Site", "cross-site")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
		assert.Contains(t, response.Body.String(), "cross-origin request rejected")
	})
}
