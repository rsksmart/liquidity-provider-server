package middlewares

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
)

// NewSessionMiddleware validates the management session cookie, rejecting
// requests without a recognized session with a 403 JSON error.
func NewSessionMiddleware(store cookies.SessionStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := store.Refresh(w, r); err != nil {
				jsonErr := rest.NewErrorResponse("session not recognized", true)
				rest.JsonErrorResponse(w, http.StatusForbidden, jsonErr)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
