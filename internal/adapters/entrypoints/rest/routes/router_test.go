package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/routes"
	"github.com/stretchr/testify/assert"
)

func TestRouter_Handle(t *testing.T) {
	router := routes.NewRouter()
	called := false
	router.Handle("GET /test", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	handler := router.BuildHandler()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestRouter_Handler(t *testing.T) {
	router := routes.NewRouter()
	router.Handle("GET /test", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	_, pattern := router.Handler(req)
	assert.Equal(t, "GET /test", pattern)
}

func TestRouter_Use_AppliesMiddlewareInRegistrationOrder(t *testing.T) {
	router := routes.NewRouter()
	var order []string
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "first")
			next.ServeHTTP(w, r)
		})
	})
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "second")
			next.ServeHTTP(w, r)
		})
	})
	router.Handle("GET /test", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		order = append(order, "handler")
	}))

	handler := router.BuildHandler()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, []string{"first", "second", "handler"}, order)
}
