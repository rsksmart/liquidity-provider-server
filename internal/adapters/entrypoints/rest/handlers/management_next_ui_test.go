package handlers

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

type errReadFile struct{}

func (errReadFile) Stat() (fs.FileInfo, error) { return nil, errors.New("stat error") }
func (errReadFile) Read([]byte) (int, error)  { return 0, errors.New("read error") }
func (errReadFile) Close() error              { return nil }

type errReadFS struct{}

func (errReadFS) Open(name string) (fs.File, error) {
	if name == "index.html" {
		return errReadFile{}, nil
	}
	return nil, fs.ErrNotExist
}

func TestManagementNextUIHandler_NormalizesTraversalToSafePath(t *testing.T) {
	dist := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("<html/>")},
	}
	handler := NewManagementNextUIHandler(dist)

	request := httptest.NewRequest(http.MethodGet, "/management/next/../secrets.txt", nil)
	request = mux.SetURLVars(request, map[string]string{"path": "../secrets.txt"})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer func() { _ = response.Body.Close() }()
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Contains(t, response.Header.Get("Cache-Control"), "no-store")
}

func TestManagementNextUIHandler_ServesEmbeddedAssetsAndSpaFallback(t *testing.T) {
	dist := fstest.MapFS{
		"index.html":             &fstest.MapFile{Data: []byte("<html/>")},
		"assets/app-abc123.js":   &fstest.MapFile{Data: []byte("console.log('ok')")},
		"assets/app-abc123.css":  &fstest.MapFile{Data: []byte("body{}")},
		"robots.txt":             &fstest.MapFile{Data: []byte("User-agent: *")},
		"assets/nested/file.txt": &fstest.MapFile{Data: []byte("x")},
	}
	handler := NewManagementNextUIHandler(dist)

	router := mux.NewRouter()
	router.Path("/management/next/{path:.*}").Methods(http.MethodGet).Handler(handler)

	server := httptest.NewServer(router)
	t.Cleanup(server.Close)

	client := &http.Client{}

	t.Run("index", func(t *testing.T) {
		request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+"/management/next/", nil)
		require.NoError(t, err)
		resp, err := client.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Contains(t, resp.Header.Get("Cache-Control"), "no-store")
	})

	t.Run("hashed asset immutable", func(t *testing.T) {
		request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+"/management/next/assets/app-abc123.js", nil)
		require.NoError(t, err)
		resp, err := client.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "public, max-age=31536000, immutable", resp.Header.Get("Cache-Control"))
	})

	t.Run("existing non-asset file no-store", func(t *testing.T) {
		request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+"/management/next/robots.txt", nil)
		require.NoError(t, err)
		resp, err := client.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Contains(t, resp.Header.Get("Cache-Control"), "no-store")
	})

	t.Run("missing file falls back to index", func(t *testing.T) {
		request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+"/management/next/some/client/route", nil)
		require.NoError(t, err)
		resp, err := client.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Contains(t, resp.Header.Get("Content-Type"), "text/html")
	})
}

func TestManagementNextUIHandler_PropagatesCopyError(t *testing.T) {
	handler := NewManagementNextUIHandler(errReadFS{})

	request := httptest.NewRequest(http.MethodGet, "/management/next/", nil)
	request = mux.SetURLVars(request, map[string]string{"path": ""})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer func() { _ = response.Body.Close() }()

	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "failed to write response")
}

