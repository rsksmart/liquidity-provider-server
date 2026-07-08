package routes

import (
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

func ConfigureRoutes(router Router, env environment.Environment, useCaseRegistry registry.UseCaseRegistry, endpointFactory EndpointFactory) {
	store, err := cookies.GetSessionCookieStore(env.Management)
	if err != nil {
		log.Fatal("Error registering routes: ", err)
	}

	registerPublicRoutes(router, env, endpointFactory.GetPublic(useCaseRegistry))

	if env.Management.EnableManagementApi {
		registerManagementRoutes(router, store, endpointFactory.GetPrivate(env, useCaseRegistry, store))
	}

	router.Handle(http.MethodOptions+" /", handlers.NewOptionsHandler())

	router.Use(middlewares.NewCorsMiddleware(env.AllowedOrigins))
}

func registerPublicRoutes(router Router, env environment.Environment, endpoints []PublicEndpoint) {
	captchaMiddleware := middlewares.NewCaptchaMiddleware(env.Captcha.Url, env.Captcha.Threshold, env.Captcha.Disabled, env.Captcha.SecretKey)
	for _, endpoint := range endpoints {
		handler := endpoint.Handler
		if endpoint.RequiresCaptcha {
			handler = captchaMiddleware(handler)
		}
		registerEndpoint(router, endpoint.Method, endpoint.Path, handler)
	}
}

func registerManagementRoutes(router Router, env environment.Environment, store sessions.Store, endpoints []Endpoint) {
	log.Warn(
		"Server is running with the management API exposed. This interface " +
			"includes endpoints that must remain private at all cost. Please shut down " +
			"the server if you haven't configured the WAF properly as explained in documentation.",
	)

	crossOriginMiddleware := middlewares.NewCrossOriginProtectionMiddleware()
	sessionValidator := middlewares.NewSessionMiddleware(store)
	var handler http.Handler
	for _, endpoint := range endpoints {
		if slices.Contains(AllowedPaths[:], endpoint.Path) {
			handler = useMiddlewares(endpoint.Handler, crossOriginMiddleware)
		} else {
			handler = useMiddlewares(endpoint.Handler, sessionValidator, crossOriginMiddleware)
		}
		registerEndpoint(router, endpoint.Method, endpoint.Path, handler)
	}
}

// registerEndpoint registers handler for the endpoint's method and path. For GET routes it also
// registers an explicit HEAD handler returning 405. http.ServeMux otherwise dispatches HEAD requests
// to the GET handler, whereas the previous gorilla router rejected any method not explicitly declared
// on a route; the explicit HEAD handler preserves that 405 behavior.
func registerEndpoint(router Router, method, path string, handler http.Handler) {
	router.Handle(method+" "+path, handler)
	if method == http.MethodGet {
		router.Handle(http.MethodHead+" "+path, methodNotAllowedHandler(method))
	}
}

// methodNotAllowedHandler responds with 405 and an Allow header, matching the previous gorilla router's
// response for a request whose method is not registered on an otherwise matching path.
func methodNotAllowedHandler(allowedMethod string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Allow", allowedMethod)
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
}

func useMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
