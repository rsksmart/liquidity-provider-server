package routes

import (
	"context"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/managementnextui"
	"github.com/stretchr/testify/require"
)

func TestManagementNextUI_EmbeddedFSShape(t *testing.T) {
	dist := managementnextui.Dist()

	_, err := fs.Stat(dist, "index.html")
	require.NoError(t, err, "expected dist/index.html to be embedded (run `cd ui && pnpm run build` before go test/build)")

	assets := listEmbeddedAssets(t, dist)
	require.NotEmpty(t, findFirstAssetPath(assets, ".js"), "expected at least one .js under dist/assets/")
}

func TestManagementNextUI_Headers_IndexAndAssets(t *testing.T) {
	dist := managementnextui.Dist()
	handler := handlers.NewManagementNextUIHandler(dist)

	router := mux.NewRouter()
	router.Path(NextUiPath).Methods(http.MethodGet).Handler(handler)

	testServer := httptest.NewServer(router)
	t.Cleanup(testServer.Close)

	client := &http.Client{}

	// Index should never be immutable cached.
	{
		request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, testServer.URL+"/management/next/", nil)
		require.NoError(t, err)
		response, err := client.Do(request)
		require.NoError(t, err)
		defer response.Body.Close()
		require.Equal(t, http.StatusOK, response.StatusCode)
		require.Contains(t, response.Header.Get("Cache-Control"), "no-store")
		require.Contains(t, response.Header.Get("Content-Type"), "text/html")
	}

	for _, assetPath := range listEmbeddedAssets(t, dist) {
		contentTypeSub, ok := assetContentTypeSub(assetPath)
		if !ok {
			continue
		}
		t.Run(assetPath, func(t *testing.T) {
			request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, testServer.URL+"/management/next/"+assetPath, nil)
			require.NoError(t, err)
			response, err := client.Do(request)
			require.NoError(t, err)
			defer response.Body.Close()
			require.Equal(t, http.StatusOK, response.StatusCode)
			require.Equal(t, "public, max-age=31536000, immutable", response.Header.Get("Cache-Control"))
			require.Contains(t, strings.ToLower(response.Header.Get("Content-Type")), contentTypeSub)
		})
	}
}

func listEmbeddedAssets(t *testing.T, dist fs.FS) []string {
	t.Helper()

	var paths []string
	require.NoError(t, fs.WalkDir(dist, "assets", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		paths = append(paths, p)
		return nil
	}))
	return paths
}

func findFirstAssetPath(paths []string, ext string) string {
	for _, p := range paths {
		if strings.HasSuffix(p, ext) {
			return p
		}
	}
	return ""
}

func assetContentTypeSub(path string) (string, bool) {
	switch {
	case strings.HasSuffix(path, ".js"):
		return "javascript", true
	case strings.HasSuffix(path, ".css"):
		return "text/css", true
	default:
		return "", false
	}
}
