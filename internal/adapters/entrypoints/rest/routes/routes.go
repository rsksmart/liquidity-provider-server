package routes

import (
	"github.com/gorilla/mux"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/middlewares"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Endpoint struct {
	Path    string
	Method  string
	Handler http.Handler
}

func ConfigureRoutes(router *mux.Router, env environment.Environment, useCaseRegistry registry.UseCaseRegistry) {
	router.Use(middlewares.NewCorsMiddleware())

	// TODO this handler is temporary, it is only to provide the CSRF token to the client, in further tasks, this
	// token should be provided inside the login template
	// ---
	router.Path("/csrf").Methods(http.MethodGet).
		Handler(
			middlewares.NewDummyCsrfTokenHandler(env.Management),
		)
	// ---

	registerPublicRoutes(router, env, useCaseRegistry)

	if env.Management.EnableManagementApi {
		registerManagementRoutes(router, env, useCaseRegistry)
	}

	router.Methods(http.MethodOptions).HandlerFunc(handlers.NewOptionsHandler())
}

func registerPublicRoutes(router *mux.Router, env environment.Environment, useCaseRegistry registry.UseCaseRegistry) {
	captchaMiddleware := middlewares.NewCaptchaMiddleware(env.Captcha.Url, env.Captcha.Threshold, env.Captcha.Disabled, env.Captcha.SecretKey)
	for _, endpoint := range getPublicEndpoints(useCaseRegistry) {
		handler := endpoint.Handler
		if endpoint.RequiresCaptcha {
			handler = useMiddlewares(endpoint.Handler, captchaMiddleware)
		}
		router.Path(endpoint.Path).Methods(endpoint.Method).Handler(handler)
	}
}

func registerManagementRoutes(router *mux.Router, env environment.Environment, useCaseRegistry registry.UseCaseRegistry) {
	log.Warn(
		"Server is running with the management API exposed. This interface " +
			"includes endpoints that must remain private at all cost. Please shut down " +
			"the server if you haven't configured the WAF properly as explained in documentation.",
	)

	store, err := cookies.GetSessionCookieStore(env.Management)
	if err != nil {
		log.Fatal("Error registering management routes: ", err)
	}
	sessionMiddlewares := middlewares.NewSessionMiddlewares(env.Management, store)
	managementEndpoints := getManagementEndpoints(env, useCaseRegistry)
	var handler http.Handler
	for _, endpoint := range managementEndpoints {
		if endpoint.Path == LOGIN_PATH {
			handler = useMiddlewares(endpoint.Handler, sessionMiddlewares.Csrf)
		} else {
			handler = useMiddlewares(endpoint.Handler, sessionMiddlewares.SessionValidator, sessionMiddlewares.Csrf)
		}
		router.Path(endpoint.Path).Methods(endpoint.Method).Handler(handler)
	}
}

func useMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
