package routes_test

import (
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/gorilla/csrf"
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
	onlyPublicMux := http.NewServeMux()
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

	routes.ConfigureRoutes(onlyPublicMux, onlyPublicEnv, useCaseRegistry, newBlockedEndpointFactory())
	publicEndpoints := toEndpoints(routes.GetPublicEndpoints(useCaseRegistry))

	t.Run("should configure cors middleware if origin is present", func(t *testing.T) {
		testCorsMiddleware(t, publicEndpoints, onlyPublicMux, testAllowedDomains[0], true)
	})

	t.Run("should not allow any domain if origin is not present", func(t *testing.T) {
		testCorsMiddleware(t, publicEndpoints, onlyPublicMux, "", false)
	})

	t.Run("should not allow any domain if origin is not allowed", func(t *testing.T) {
		testCorsMiddleware(t, publicEndpoints, onlyPublicMux, "https://not-allowed.com", false)
	})

	t.Run("should configure options handler", func(t *testing.T) {
		assertConfiguresOptionsHandler(t, onlyPublicMux)
	})

	t.Run("should register public routes", func(t *testing.T) {
		testPublicRoutesRegistration(t, useCaseRegistry, onlyPublicMux)
	})

	t.Run("should register management routes only if Management API is enabled", func(t *testing.T) {
		managementRoutes := routes.GetManagementEndpoints(onlyPublicEnv, useCaseRegistry, &mocks.StoreMock{})
		for _, endpoint := range managementRoutes {
			req := httptest.NewRequest(endpoint.Method, requestPath(endpoint.Path), nil)
			_, pattern := onlyPublicMux.Handler(req)
			assert.Empty(t, pattern, "management route %s %s should not be registered", endpoint.Method, endpoint.Path)
		}
	})
}

// nolint:funlen
func TestConfigureRoutes_Management(t *testing.T) {
	managementMux := http.NewServeMux()
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
	routes.ConfigureRoutes(managementMux, managementEnv, useCaseRegistry, newBlockedEndpointFactory())

	managementEndpoints := routes.GetManagementEndpoints(managementEnv, useCaseRegistry, &mocks.StoreMock{})
	allEndpoints := append(toEndpoints(routes.GetPublicEndpoints(useCaseRegistry)), managementEndpoints...)

	t.Run("should configure cors middleware", func(t *testing.T) {
		testCorsMiddleware(t, allEndpoints, managementMux, testAllowedDomains[1], true)
	})

	t.Run("should not allow any domain if origin is not present", func(t *testing.T) {
		testCorsMiddleware(t, allEndpoints, managementMux, "", false)
	})

	t.Run("should not allow any domain if origin is not allowed", func(t *testing.T) {
		testCorsMiddleware(t, allEndpoints, managementMux, "https://not-allowed.com", false)
	})

	t.Run("should configure options handler", func(t *testing.T) {
		assertConfiguresOptionsHandler(t, managementMux)
	})

	t.Run("should register public routes", func(t *testing.T) {
		testPublicRoutesRegistration(t, useCaseRegistry, managementMux)
	})

	t.Run("should register management routes only if Management API is enabled", func(t *testing.T) {
		for _, endpoint := range managementEndpoints {
			req := httptest.NewRequest(endpoint.Method, requestPath(endpoint.Path), nil)
			_, pattern := managementMux.Handler(req)
			assert.NotEmpty(t, pattern, "management route %s %s not registered", endpoint.Method, endpoint.Path)
		}
		t.Run("should register management routes with proper middlewares", func(t *testing.T) {
			for _, endpoint := range managementEndpoints {
				if slices.Contains(routes.AllowedPaths[:], endpoint.Path) {
					assertHasCsrfMiddleware(t, managementMux, endpoint)
				} else {
					req := httptest.NewRequest(http.MethodGet, routes.UiPath, nil)
					responseRecorder := httptest.NewRecorder()
					managementMux.ServeHTTP(responseRecorder, req)
					assertHasCsrfMiddleware(t, managementMux, endpoint)
					// nolint:bodyclose
					assertHasSessionMiddleware(t, managementMux, endpoint, responseRecorder.Result().Cookies()[0], responseRecorder.Header().Get(csrfTokenHeaderName))
					require.NoError(t, responseRecorder.Result().Body.Close())
				}
			}
		})
	})
}

func TestConfigureRoutes_HeadMatchesGetRoute(t *testing.T) {
	publicMux := http.NewServeMux()
	useCaseRegistry := &mocks.UseCaseRegistryMock{}
	setupRegistryMock(useCaseRegistry)
	env := environment.Environment{
		Management: environment.ManagementEnv{
			EnableManagementApi:  false,
			SessionAuthKey:       hex.EncodeToString(make([]byte, 32)),
			SessionEncryptionKey: hex.EncodeToString(make([]byte, 32)),
			SessionTokenAuthKey:  hex.EncodeToString(make([]byte, 32)),
			UseHttps:             false,
		},
		AllowedOrigins: testAllowedDomains,
	}
	routes.ConfigureRoutes(publicMux, env, useCaseRegistry, newBlockedEndpointFactory())

	// ServeMux routes HEAD to the GET handler (net/http strips the body), so a HEAD request reaches
	// the handler instead of returning 405. This is an intentional deviation from the gorilla router.
	req := httptest.NewRequest(http.MethodHead, "/health", nil)
	responseRecorder := httptest.NewRecorder()
	publicMux.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusTeapot, responseRecorder.Code)
}

// toEndpoints flattens the public endpoint wrappers to plain endpoints so the registration and CORS
// helpers can iterate public and management routes uniformly.
func toEndpoints(publicEndpoints []routes.PublicEndpoint) []routes.Endpoint {
	result := make([]routes.Endpoint, 0, len(publicEndpoints))
	for _, endpoint := range publicEndpoints {
		result = append(result, endpoint.Endpoint)
	}
	return result
}

// requestPath turns a registered pattern into a concrete request path, substituting the {file}
// wildcard of the static assets route.
func requestPath(pattern string) string {
	if pattern == routes.StaticPath {
		return "/static/test"
	}
	return pattern
}

func assertConfiguresOptionsHandler(t *testing.T, handler http.Handler) {
	req := httptest.NewRequest(http.MethodOptions, "/aPath", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func testCorsMiddleware(
	t *testing.T,
	endpoints []routes.Endpoint,
	handler http.Handler,
	origin string,
	isOriginAllowed bool,
) {
	for _, endpoint := range endpoints {
		req := httptest.NewRequest(endpoint.Method, requestPath(endpoint.Path), nil)
		req.Header.Set("Origin", origin)
		responseRecorder := httptest.NewRecorder()
		handler.ServeHTTP(responseRecorder, req)
		if isOriginAllowed {
			assertHasCorsHeadersAllowed(t, responseRecorder, origin)
		} else {
			assertHasCorsHeadersNotAllowed(t, responseRecorder)
		}
	}
}

func testPublicRoutesRegistration(t *testing.T, useCaseRegistry registry.UseCaseRegistry, mux *http.ServeMux) {
	publicRoutes := routes.GetPublicEndpoints(useCaseRegistry)
	for _, endpoint := range publicRoutes {
		req := httptest.NewRequest(endpoint.Method, requestPath(endpoint.Path), nil)
		_, pattern := mux.Handler(req)
		assert.NotEmpty(t, pattern, "public route %s %s not registered", endpoint.Method, endpoint.Path)
	}
	t.Run("should use captcha middleware in proper routes", func(t *testing.T) {
		for _, endpoint := range publicRoutes {
			if endpoint.RequiresCaptcha {
				req := httptest.NewRequest(endpoint.Method, requestPath(endpoint.Path), nil)
				responseRecorder := httptest.NewRecorder()
				mux.ServeHTTP(responseRecorder, req)
				assert.Contains(t, responseRecorder.Body.String(), "missing X-Captcha-Token header")
			}
		}
	})
}

func assertHasSessionMiddleware(t *testing.T, handler http.Handler, endpoint routes.Endpoint, cookie *http.Cookie, token string) {
	request := httptest.NewRequest(endpoint.Method, requestPath(endpoint.Path), nil)
	request.AddCookie(cookie)
	request.Header.Set(csrfTokenHeaderName, token)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	// nolint:bodyclose
	assert.Equal(t, http.StatusForbidden, response.Result().StatusCode)
	assert.Contains(t, response.Body.String(), "session not recognized")
	require.NoError(t, response.Result().Body.Close())
}

func assertHasCsrfMiddleware(t *testing.T, handler http.Handler, endpoint routes.Endpoint) {
	req := httptest.NewRequest(endpoint.Method, requestPath(endpoint.Path), nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, req)
	// nolint:bodyclose
	i := slices.IndexFunc(responseRecorder.Result().Cookies(), func(cookie *http.Cookie) bool {
		return cookie.Name == cookies.CsrfCookieName
	})
	require.NoError(t, responseRecorder.Result().Body.Close())
	assert.NotEqual(t, -1, i, "response does not have CSRF cookie")
}

// nolint:funlen
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
	registryMock.EXPECT().GetLiquidityRatioUseCase().Return(&liquidity_provider.GetLiquidityRatioUseCase{})
	registryMock.EXPECT().SetLiquidityRatioUseCase().Return(&liquidity_provider.SetLiquidityRatioUseCase{})
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
