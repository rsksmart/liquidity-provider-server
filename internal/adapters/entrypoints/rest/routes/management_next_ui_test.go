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

	js := findFirstAsset(t, dist, ".js")
	css := findFirstAsset(t, dist, ".css")

	require.NotEmpty(t, js)
	require.NotEmpty(t, css)
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

	js := findFirstAsset(t, dist, ".js")
	css := findFirstAsset(t, dist, ".css")

	// Hashed assets should be immutable cached + correct MIME.
	for _, tc := range []struct {
		name           string
		path           string
		contentTypeSub string
	}{
		{name: "js", path: js, contentTypeSub: "javascript"},
		{name: "css", path: css, contentTypeSub: "text/css"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, testServer.URL+"/management/next/"+tc.path, nil)
			require.NoError(t, err)
			response, err := client.Do(request)
			require.NoError(t, err)
			defer response.Body.Close()
			require.Equal(t, http.StatusOK, response.StatusCode)
			require.Equal(t, "public, max-age=31536000, immutable", response.Header.Get("Cache-Control"))
			require.Contains(t, strings.ToLower(response.Header.Get("Content-Type")), tc.contentTypeSub)
		})
	}
}

func findFirstAsset(t *testing.T, dist fs.FS, ext string) string {
	t.Helper()

	var out string
	require.NoError(t, fs.WalkDir(dist, "assets", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(p, ext) {
			out = p
			return fs.SkipAll
		}
		return nil
	}))
	require.NotEmpty(t, out, "expected at least one %s under dist/assets/", ext)
	return out
}
