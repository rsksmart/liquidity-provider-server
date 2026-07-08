package routes_test

import (
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/routes"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigureRoutes_StaticAssets_DoesNotListDirectory(t *testing.T) {
	handler := managementAPIHandler(t)

	for _, path := range []string{routes.StaticRootPath, routes.StaticPath} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		responseRecorder := httptest.NewRecorder()
		handler.ServeHTTP(responseRecorder, req)

		body, err := io.ReadAll(responseRecorder.Body)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
		assert.NotContains(t, string(body), "<pre>")
		assert.NotContains(t, strings.ToLower(string(body)), "index of")
	}
}

func TestConfigureRoutes_StaticAssets_ServesFile(t *testing.T) {
	handler := managementAPIHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/static/management.js", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func managementAPIHandler(t *testing.T) http.Handler {
	t.Helper()
	router := routes.NewRouter()
	useCaseRegistry := &mocks.UseCaseRegistryMock{}
	setupRegistryMock(useCaseRegistry)
	env := environment.Environment{
		Management: environment.ManagementEnv{
			EnableManagementApi:  true,
			SessionAuthKey:       hex.EncodeToString(make([]byte, 32)),
			SessionEncryptionKey: hex.EncodeToString(make([]byte, 32)),
			SessionTokenAuthKey:  hex.EncodeToString(make([]byte, 32)),
			UseHttps:             false,
		},
		AllowedOrigins: testAllowedDomains,
	}
	routes.ConfigureRoutes(router, env, useCaseRegistry, routes.NewEndpointFactory())
	return router.BuildHandler()
}
