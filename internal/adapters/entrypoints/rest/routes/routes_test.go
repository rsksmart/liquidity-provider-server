package routes_test

import (
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/routes"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// nolint:gosec // Linter is assuming the header name is a password
const csrfTokenHeaderName = "X-Csrf-Token"

var testAllowedDomains = []string{"https://allowed.com", "https://another-allowed.com"}

type openApiSpecification struct {
	// Path - Verb
	Paths map[string]map[string]any `yaml:"paths"`
}

// nolint:funlen
func TestConfigureRoutes_Public(t *testing.T) {
	onlyPublicRouter := mux.NewRouter()
	useCaseRegistry := &mocks.UseCaseRegistryMock{}
	setupRegistryMock(useCaseRegistry)
	onlyPublicEnv := environment.Environment{
		Management: environment.ManagementEnv{
			EnableManagementApi:  false,
			SessionAuthKey:       hex.EncodeToString(make([]byte, 32)),
			SessionEncryptionKey: hex.EncodeToString(make([]byte, 32)),
			SessionTokenAuthKey:  hex.EncodeToString(make([]byte, 32)),
			UseHttps:             false,
		},
		AllowedOrigins: testAllowedDomains,
	}

	routes.ConfigureRoutes(onlyPublicRouter, onlyPublicEnv, useCaseRegistry, newBlockedEndpointFactory())
	onlyPublicRoutes := make([]*mux.Route, 0)

	err := onlyPublicRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		onlyPublicRoutes = append(onlyPublicRoutes, route)
		return nil
	})
	require.NoError(t, err)

	t.Run("should configure cors middleware if origin is present", func(t *testing.T) {
		testCorsMiddleware(t, onlyPublicRoutes, onlyPublicRouter, testAllowedDomains[0], true)
	})

	t.Run("should not allow any domain if origin is not present", func(t *testing.T) {
		testCorsMiddleware(t, onlyPublicRoutes, onlyPublicRouter, "", false)
	})

	t.Run("should not allow any domain if origin is not allowed", func(t *testing.T) {
		testCorsMiddleware(t, onlyPublicRoutes, onlyPublicRouter, "https://not-allowed.com", false)
	})

	t.Run("should configure options handler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/aPath", nil)
		ok := slices.ContainsFunc(onlyPublicRoutes, func(r *mux.Route) bool {
			return r.Match(req, &mux.RouteMatch{})
		})
		assert.True(t, ok)
	})

	t.Run("should register public routes", func(t *testing.T) {
		testPublicRoutesRegistration(t, useCaseRegistry, onlyPublicRoutes, onlyPublicRouter)
	})

	t.Run("should register management routes only if Management API is enabled", func(t *testing.T) {
		managementRoutes := routes.GetManagementEndpoints(onlyPublicEnv, useCaseRegistry, &mocks.StoreMock{})
		for _, endpoint := range managementRoutes {
			req := httptest.NewRequest(endpoint.Method, endpoint.Path, nil)
			result := slices.ContainsFunc(onlyPublicRoutes, func(r *mux.Route) bool {
				return r.Match(req, &mux.RouteMatch{})
			})
			assert.False(t, result)
		}
	})
}

// nolint:funlen
func TestConfigureRoutes_Management(t *testing.T) {
	managementRouter := mux.NewRouter()
	useCaseRegistry := &mocks.UseCaseRegistryMock{}
	setupRegistryMock(useCaseRegistry)
	managementEnv := environment.Environment{
		Management: environment.ManagementEnv{
			EnableManagementApi:  true,
			SessionAuthKey:       hex.EncodeToString(make([]byte, 32)),
			SessionEncryptionKey: hex.EncodeToString(make([]byte, 32)),
			SessionTokenAuthKey:  hex.EncodeToString(make([]byte, 32)),
			UseHttps:             false,
		},
		AllowedOrigins: testAllowedDomains,
	}
	routes.ConfigureRoutes(managementRouter, managementEnv, useCaseRegistry, newBlockedEndpointFactory())
	managementAndPublicRoutes := make([]*mux.Route, 0)

	err := managementRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		managementAndPublicRoutes = append(managementAndPublicRoutes, route)
		return nil
	})
	require.NoError(t, err)

	t.Run("should configure cors middleware", func(t *testing.T) {
		testCorsMiddleware(t, managementAndPublicRoutes, managementRouter, testAllowedDomains[1], true)
	})

	t.Run("should not allow any domain if origin is not present", func(t *testing.T) {
		testCorsMiddleware(t, managementAndPublicRoutes, managementRouter, "", false)
	})

	t.Run("should not allow any domain if origin is not allowed", func(t *testing.T) {
		testCorsMiddleware(t, managementAndPublicRoutes, managementRouter, "https://not-allowed.com", false)
	})

	t.Run("should configure options handler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/aPath", nil)
		ok := slices.ContainsFunc(managementAndPublicRoutes, func(r *mux.Route) bool {
			return r.Match(req, &mux.RouteMatch{})
		})
		assert.True(t, ok)
	})

	t.Run("should register public routes", func(t *testing.T) {
		testPublicRoutesRegistration(t, useCaseRegistry, managementAndPublicRoutes, managementRouter)
	})

	t.Run("should register management routes only if Management API is enabled", func(t *testing.T) {
		managementRoutes := routes.GetManagementEndpoints(managementEnv, useCaseRegistry, &mocks.StoreMock{})
		for _, endpoint := range managementRoutes {
			req := httptest.NewRequest(endpoint.Method, endpoint.Path, nil)
			ok := slices.ContainsFunc(managementAndPublicRoutes, func(r *mux.Route) bool {
				return r.Match(req, &mux.RouteMatch{})
			})
			assert.True(t, ok)
		}
		t.Run("should register management routes with proper middlewares", func(t *testing.T) {
			for _, endpoint := range managementRoutes {
				if slices.Contains(routes.AllowedPaths[:], endpoint.Path) {
					assertHasCsrfMiddleware(t, managementRouter, endpoint)
				} else {
					req := httptest.NewRequest(http.MethodGet, routes.UiPath, nil)
					responseRecorder := httptest.NewRecorder()
					managementRouter.ServeHTTP(responseRecorder, req)
					assertHasCsrfMiddleware(t, managementRouter, endpoint)
					// nolint:bodyclose
					assertHasSessionMiddleware(t, managementRouter, endpoint, responseRecorder.Result().Cookies()[0], responseRecorder.Header().Get(csrfTokenHeaderName))
					require.NoError(t, responseRecorder.Result().Body.Close())
				}
			}
		})
	})
}

func testCorsMiddleware(
	t *testing.T,
	routesToTest []*mux.Route,
	routerToTest *mux.Router,
	origin string,
	isOriginAllowed bool,
) {
	for _, route := range routesToTest {
		methods, methodsErr := route.GetMethods()
		require.NoError(t, methodsErr)
		for _, method := range methods {
			if method != http.MethodOptions {
				path, pathErr := route.GetPathTemplate()
				require.NoError(t, pathErr)
				req := httptest.NewRequest(method, path, nil)
				req.Header.Set("Origin", origin)
				responseRecorder := httptest.NewRecorder()
				routerToTest.ServeHTTP(responseRecorder, req)
				if isOriginAllowed {
					assertHasCorsHeadersAllowed(t, responseRecorder, origin)
				} else {
					assertHasCorsHeadersNotAllowed(t, responseRecorder)
				}
			}
		}
	}
}

func testPublicRoutesRegistration(t *testing.T, useCaseRegistry registry.UseCaseRegistry, routesToTest []*mux.Route, routerToTest *mux.Router) {
	publicRoutes := routes.GetPublicEndpoints(useCaseRegistry)
	for _, endpoint := range publicRoutes {
		req := httptest.NewRequest(endpoint.Method, endpoint.Path, nil)
		publicRoutesOk := slices.ContainsFunc(routesToTest, func(r *mux.Route) bool {
			return r.Match(req, &mux.RouteMatch{})
		})
		assert.True(t, publicRoutesOk)
	}
	t.Run("should use captcha middleware in proper routes", func(t *testing.T) {
		for _, endpoint := range publicRoutes {
			if endpoint.RequiresCaptcha {
				req := httptest.NewRequest(endpoint.Method, endpoint.Path, nil)
				responseRecorder := httptest.NewRecorder()
				routerToTest.ServeHTTP(responseRecorder, req)
				assert.Contains(t, responseRecorder.Body.String(), "missing X-Captcha-Token header")
			}
		}
	})
}

func assertHasSessionMiddleware(t *testing.T, router *mux.Router, endpoint routes.Endpoint, cookie *http.Cookie, token string) {
	request := httptest.NewRequest(endpoint.Method, endpoint.Path, nil)
	request.AddCookie(cookie)
	request.Header.Set(csrfTokenHeaderName, token)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	// nolint:bodyclose
	assert.Equal(t, http.StatusForbidden, response.Result().StatusCode)
	assert.Contains(t, response.Body.String(), "session not recognized")
	require.NoError(t, response.Result().Body.Close())
}

func assertHasCsrfMiddleware(t *testing.T, router *mux.Router, endpoint routes.Endpoint) {
	req := httptest.NewRequest(endpoint.Method, endpoint.Path, nil)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	// nolint:bodyclose
	i := slices.IndexFunc(responseRecorder.Result().Cookies(), func(cookie *http.Cookie) bool {
		return cookie.Name == cookies.CsrfCookieName
	})
	require.NoError(t, responseRecorder.Result().Body.Close())
	assert.NotEqual(t, -1, i, "response does not have CSRF cookie")
}

func setupRegistryMock(registryMock *mocks.UseCaseRegistryMock) {
	acceptQuoteUseCase := &pegin.AcceptQuoteUseCase{}

	registryMock.EXPECT().HealthUseCase().Return(&usecases.HealthUseCase{})
	registryMock.EXPECT().GetProvidersUseCase().Return(&liquidity_provider.GetProvidersUseCase{})
	registryMock.EXPECT().GetPeginQuoteUseCase().Return(&pegin.GetQuoteUseCase{})
	registryMock.EXPECT().GetAcceptPeginQuoteUseCase().Return(acceptQuoteUseCase)
	registryMock.EXPECT().GetPegoutQuoteUseCase().Return(&pegout.GetQuoteUseCase{})
	registryMock.EXPECT().GetAcceptPegoutQuoteUseCase().Return(&pegout.AcceptQuoteUseCase{})
	registryMock.EXPECT().GetUserDepositsUseCase().Return(&pegout.GetUserDepositsUseCase{})
	registryMock.EXPECT().GetProviderDetailUseCase().Return(&liquidity_provider.GetDetailUseCase{})
	registryMock.EXPECT().GetPeginStatusUseCase().Return(&pegin.StatusUseCase{})
	registryMock.EXPECT().GetPegoutStatusUseCase().Return(&pegout.StatusUseCase{})
	registryMock.EXPECT().GetAvailableLiquidityUseCase().Return(&liquidity_provider.GetAvailableLiquidityUseCase{})
	registryMock.EXPECT().SummariesUseCase().Return(&reports.SummariesUseCase{})
	registryMock.EXPECT().GetServerInfoUseCase().Return(&liquidity_provider.ServerInfoUseCase{})

	registryMock.EXPECT().GetPeginCollateralUseCase().Return(&pegin.GetCollateralUseCase{})
	registryMock.EXPECT().AddPeginCollateralUseCase().Return(&pegin.AddCollateralUseCase{})
	registryMock.EXPECT().GetPegoutCollateralUseCase().Return(&pegout.GetCollateralUseCase{})
	registryMock.EXPECT().AddPegoutCollateralUseCase().Return(&pegout.AddCollateralUseCase{})
	registryMock.EXPECT().ChangeStatusUseCase().Return(&liquidity_provider.ChangeStatusUseCase{})
	registryMock.EXPECT().ResignationUseCase().Return(&liquidity_provider.ResignUseCase{})
	registryMock.EXPECT().WithdrawCollateralUseCase().Return(&liquidity_provider.WithdrawCollateralUseCase{})
	registryMock.EXPECT().GetConfigurationUseCase().Return(&liquidity_provider.GetConfigUseCase{})
	registryMock.EXPECT().SetGeneralConfigUseCase().Return(&liquidity_provider.SetGeneralConfigUseCase{})
	registryMock.EXPECT().SetPeginConfigUseCase().Return(&liquidity_provider.SetPeginConfigUseCase{})
	registryMock.EXPECT().SetPegoutConfigUseCase().Return(&liquidity_provider.SetPegoutConfigUseCase{})
	registryMock.EXPECT().SetCredentialsUseCase().Return(&liquidity_provider.SetCredentialsUseCase{})
	registryMock.EXPECT().LoginUseCase().Return(&liquidity_provider.LoginUseCase{})
	registryMock.EXPECT().GetManagementUiDataUseCase().Return(&liquidity_provider.GetManagementUiDataUseCase{})
	registryMock.EXPECT().GetServerInfoUseCase().Return(&liquidity_provider.ServerInfoUseCase{})
	registryMock.EXPECT().GetPeginReportUseCase().Return(&reports.GetPeginReportUseCase{})
	registryMock.EXPECT().GetPegoutReportUseCase().Return(&reports.GetPegoutReportUseCase{})
	registryMock.EXPECT().GetRevenueReportUseCase().Return(&reports.GetRevenueReportUseCase{})
	registryMock.EXPECT().GetAssetsReportUseCase().Return(&reports.GetAssetsReportUseCase{})
	registryMock.EXPECT().GetTransactionsReportUseCase().Return(&reports.GetTransactionsUseCase{})
	registryMock.EXPECT().GetTrustedAccountsUseCase().Return(&liquidity_provider.GetTrustedAccountsUseCase{})
	registryMock.EXPECT().UpdateTrustedAccountUseCase().Return(&liquidity_provider.UpdateTrustedAccountUseCase{})
	registryMock.EXPECT().AddTrustedAccountUseCase().Return(&liquidity_provider.AddTrustedAccountUseCase{})
	registryMock.EXPECT().DeleteTrustedAccountUseCase().Return(&liquidity_provider.DeleteTrustedAccountUseCase{})
	registryMock.EXPECT().RecommendedPegoutUseCase().Return(&pegout.RecommendedPegoutUseCase{})
	registryMock.EXPECT().RecommendedPeginUseCase().Return(&pegin.RecommendedPeginUseCase{})
}

func assertHasCorsHeadersAllowed(t *testing.T, recorder *httptest.ResponseRecorder, origin string) {
	assert.Equal(t, origin, recorder.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Content-Type, Origin, Accept, token, X-Captcha-Token, X-Csrf-Token", recorder.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", recorder.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "true", recorder.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "Origin", recorder.Header().Get("Vary"))
}

func assertHasCorsHeadersNotAllowed(t *testing.T, recorder *httptest.ResponseRecorder) {
	assert.Empty(t, recorder.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Content-Type, Origin, Accept, token, X-Captcha-Token, X-Csrf-Token", recorder.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", recorder.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Origin", recorder.Header().Get("Vary"))
}

type blockedEndpointFactory struct {
	realFactory routes.EndpointFactory
}

func newBlockedEndpointFactory() routes.EndpointFactory {
	return &blockedEndpointFactory{
		realFactory: routes.NewEndpointFactory(),
	}
}

func (f *blockedEndpointFactory) GetPublic(useCaseRegistry registry.UseCaseRegistry) []routes.PublicEndpoint {
	dummyEndpoints := make([]routes.PublicEndpoint, 0)
	endpoints := f.realFactory.GetPublic(useCaseRegistry)
	for _, endpoint := range endpoints {
		dummyEndpoints = append(dummyEndpoints, routes.PublicEndpoint{
			Endpoint: routes.Endpoint{
				Path:    endpoint.Path,
				Method:  endpoint.Method,
				Handler: teapotHandler(),
			},
			RequiresCaptcha: endpoint.RequiresCaptcha,
		})
	}
	return dummyEndpoints
}

func (f *blockedEndpointFactory) GetPrivate(env environment.Environment, useCaseRegistry registry.UseCaseRegistry, store sessions.Store) []routes.Endpoint {
	dummyEndpoints := make([]routes.Endpoint, 0)
	endpoints := f.realFactory.GetPrivate(env, useCaseRegistry, store)
	for _, endpoint := range endpoints {
		dummyEndpoints = append(dummyEndpoints, routes.Endpoint{
			Path:    endpoint.Path,
			Method:  endpoint.Method,
			Handler: teapotHandler(),
		})
	}
	return dummyEndpoints
}

func teapotHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set(csrfTokenHeaderName, csrf.Token(req))
		w.WriteHeader(http.StatusTeapot)
	}
}
