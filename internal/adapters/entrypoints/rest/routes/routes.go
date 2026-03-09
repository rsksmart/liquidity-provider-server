package routes

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/middlewares"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	log "github.com/sirupsen/logrus"
	"net/http"
	"slices"
)

type Endpoint struct {
	Path    string
	Method  string
	Handler http.Handler
}

// EndpointFactory abstraction to be able to mock the endpoints in tests
type EndpointFactory interface {
	GetPublic(useCaseRegistry registry.UseCaseRegistry) []PublicEndpoint
	GetPrivate(env environment.Environment, useCaseRegistry registry.UseCaseRegistry, store sessions.Store) []Endpoint
}

type endpointFactoryImpl struct{}

func NewEndpointFactory() EndpointFactory {
	return &endpointFactoryImpl{}
}

func (f *endpointFactoryImpl) GetPublic(useCaseRegistry registry.UseCaseRegistry) []PublicEndpoint {
	return GetPublicEndpoints(useCaseRegistry)
}

func (f *endpointFactoryImpl) GetPrivate(env environment.Environment, useCaseRegistry registry.UseCaseRegistry, store sessions.Store) []Endpoint {
	return GetManagementEndpoints(env, useCaseRegistry, store)
}

func ConfigureRoutes(router *mux.Router, env environment.Environment, useCaseRegistry registry.UseCaseRegistry, endpointFactory EndpointFactory) {
	router.Use(middlewares.NewCorsMiddleware(env.AllowedOrigins))

	store, err := cookies.GetSessionCookieStore(env.Management)
	if err != nil {
		log.Fatal("Error registering routes: ", err)
	}

	registerPublicRoutes(router, env, endpointFactory.GetPublic(useCaseRegistry))

	if env.Management.EnableManagementApi {
		registerManagementRoutes(router, env, store, endpointFactory.GetPrivate(env, useCaseRegistry, store))
	}

	router.Methods(http.MethodOptions).HandlerFunc(handlers.NewOptionsHandler())
}

func registerPublicRoutes(router *mux.Router, env environment.Environment, endpoints []PublicEndpoint) {
	captchaMiddleware := middlewares.NewCaptchaMiddleware(env.Captcha.Url, env.Captcha.Threshold, env.Captcha.Disabled, env.Captcha.SecretKey)
	for _, endpoint := range endpoints {
		handler := endpoint.Handler
		if endpoint.RequiresCaptcha {
			handler = useMiddlewares(handler, captchaMiddleware)
		}
		router.Path(endpoint.Path).Methods(endpoint.Method).Handler(handler)
	}
}

func registerManagementRoutes(router *mux.Router, env environment.Environment, store sessions.Store, endpoints []Endpoint) {
	log.Warn(
		"Server is running with the management API exposed. This interface " +
			"includes endpoints that must remain private at all cost. Please shut down " +
			"the server if you haven't configured the WAF properly as explained in documentation.",
	)

	sessionMiddlewares := middlewares.NewSessionMiddlewares(env.Management, store)
	var handler http.Handler
	for _, endpoint := range endpoints {
		if slices.Contains(AllowedPaths[:], endpoint.Path) {
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
