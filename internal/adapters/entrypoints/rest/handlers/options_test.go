package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/stretchr/testify/assert"
)

func TestOptionsHandlerHappyPath(t *testing.T) {
	request := httptest.NewRequest(http.MethodOptions, "/any-endpoint", nil)
	recorder := httptest.NewRecorder()

	handlerFunc := handlers.NewOptionsHandler()
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Empty(t, recorder.Body.String())
}

func TestOptionsHandlerDifferentPaths(t *testing.T) {
	paths := []string{
		"/pegin/getQuote",
		"/pegout/getQuote",
		"/providers/liquidity",
		"/health",
		"/version",
	}

	for _, path := range paths {
		t.Run("path "+path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodOptions, path, nil)
			recorder := httptest.NewRecorder()

			handlerFunc := handlers.NewOptionsHandler()
			handler := http.HandlerFunc(handlerFunc)

			handler.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
		})
	}
}
