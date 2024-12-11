package middlewares

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type SessionMiddlewares struct {
	Csrf             func(next http.Handler) http.Handler
	SessionValidator func(next http.Handler) http.Handler
}

func NewSessionMiddlewares(env environment.ManagementEnv, store sessions.Store) SessionMiddlewares {
	csrfStep := csrfMiddleware(env)
	sessionStep := sessionMiddleware(store)
	return SessionMiddlewares{
		Csrf:             csrfStep,
		SessionValidator: sessionStep,
	}
}

func sessionMiddleware(store sessions.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, cookies.ManagementSessionCookieName)
			if err != nil {
				jsonErr := rest.NewErrorResponseWithDetails("session validation error", rest.DetailsFromError(err), false)
				rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
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

func csrfMiddleware(env environment.ManagementEnv) func(next http.Handler) http.Handler {
	authKey, err := utils.DecodeKey(env.SessionTokenAuthKey, cookies.KeysBytesLength)
	if err != nil {
		log.Fatalf("error decoding session token auth key: %v", err)
	}
	return csrf.Protect(
		authKey,
		csrf.MaxAge(cookies.SessionMaxSeconds),
		csrf.CookieName(cookies.CsrfCookieName),
		csrf.Path("/"),
		csrf.Secure(env.UseHttps),
		csrf.HttpOnly(true),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			details := rest.DetailsFromError(csrf.FailureReason(r))
			jsonErr := rest.NewErrorResponseWithDetails("CRSF token validation error", details, true)
			rest.JsonErrorResponse(w, http.StatusForbidden, jsonErr)
		})),
	)
}
