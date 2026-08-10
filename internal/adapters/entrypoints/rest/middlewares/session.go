package middlewares

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	log "github.com/sirupsen/logrus"
)

// NewSessionMiddleware validates the management session cookie, rejecting
// requests without a recognized session with a 403 JSON error.
func NewSessionMiddleware(store sessions.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, cookies.ManagementSessionCookieName)
			if err != nil {
				jsonErr := rest.NewErrorResponseWithDetails("session validation error", rest.DetailsFromError(err), false)
				rest.JsonErrorResponse(w, http.StatusForbidden, jsonErr)
				return
			} else if session.IsNew {
				jsonErr := rest.NewErrorResponse("session not recognized", true)
				rest.JsonErrorResponse(w, http.StatusForbidden, jsonErr)
				return
			}

			next.ServeHTTP(w, r)

			if err = session.Save(r, w); err != nil {
				log.Error("Error saving session: ", err)
			}
		})
	}
}
