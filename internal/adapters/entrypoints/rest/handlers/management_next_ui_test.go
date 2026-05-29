package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const nextUiIndexTemplate = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="csrf-token" content="{{ .CsrfToken }}" />
    <title>LPS Management UI</title>
    <script nonce="{{ .ScriptNonce }}" type="module" crossorigin src="/management/next/assets/app.js"></script>
    <link rel="stylesheet" crossorigin href="/management/next/assets/app.css">
  </head>
  <body>
    <div id="root"></div>
    <script id="initial-data" type="application/json" nonce="{{ .ScriptNonce }}">{{ .InitialDataJSON }}</script>
  </body>
</html>`

const nextUiIndexMissingNonce = `<!doctype html>
<html>
  <head>
    <script type="module" src="/management/next/assets/app.js"></script>
  </head>
  <body>
    <div id="root"></div>
    <script>alert(1)</script>
  </body>
</html>`

func newNextUIHandlerTestFixtures(t *testing.T, indexHTML string, enableSecurityHeaders bool) (*ManagementNextUIHandler, *mocks.StoreMock, *mocks.GetManagementUiDataUseCaseMock) {
	t.Helper()

	dist := fstest.MapFS{
		"index.html":            &fstest.MapFile{Data: []byte(indexHTML)},
		"assets/app-abc123.js":  &fstest.MapFile{Data: []byte("console.log('ok')")},
		"assets/app-abc123.css": &fstest.MapFile{Data: []byte("body{}")},
	}

	env := environment.ManagementEnv{EnableSecurityHeaders: enableSecurityHeaders}
	mockStore := new(mocks.StoreMock)
	mockStore.On("Get", mock.Anything, "lp-session").Return(nil, errors.New("no session"))

	templateData := liquidity_provider.ManagementTemplateData{
		CredentialsSet: true,
		BaseUrl:        "http://localhost:8080",
		BtcAddress:     "tb1qexample",
		RskAddress:     "0xabc",
	}
	mockUseCase := new(mocks.GetManagementUiDataUseCaseMock)
	mockUseCase.On("Run", mock.Anything, false).Return(&liquidity_provider.ManagementTemplate{
		Name: liquidity_provider.ManagementUiTemplate,
		Data: templateData,
	}, nil)

	handler := NewManagementNextUIHandler(dist, env, mockStore, mockUseCase).(*ManagementNextUIHandler)
	return handler, mockStore, mockUseCase
}

func serveNextUIIndex(t *testing.T, handler http.Handler) *http.Response {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, "/management/next/", nil)
	request = mux.SetURLVars(request, map[string]string{"path": ""})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	t.Cleanup(func() { _ = response.Body.Close() })
	return response
}

func readResponseBody(t *testing.T, response *http.Response) string {
	t.Helper()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	return string(body)
}

func extractNonceFromCSP(csp string) string {
	re := regexp.MustCompile(`'nonce-([^']+)'`)
	matches := re.FindStringSubmatch(csp)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func normalizeCSP(csp string) string {
	nonce := extractNonceFromCSP(csp)
	if nonce == "" {
		return csp
	}
	return strings.ReplaceAll(csp, nonce, "NONCE")
}

func validateEveryScriptAndStyleHasNonce(html string, expectedNonce string) error {
	scriptTags := regexp.MustCompile(`(?i)<script[^>]*>`).FindAllString(html, -1)
	for _, tag := range scriptTags {
		if !strings.Contains(tag, `nonce="`+expectedNonce+`"`) {
			return errors.New("script tag missing expected nonce: " + tag)
		}
	}

	styleTags := regexp.MustCompile(`(?i)<style[^>]*>`).FindAllString(html, -1)
	for _, tag := range styleTags {
		if !strings.Contains(tag, `nonce="`+expectedNonce+`"`) {
			return errors.New("style tag missing expected nonce: " + tag)
		}
	}
	return nil
}

func assertEveryScriptAndStyleHasNonce(t *testing.T, html string, expectedNonce string) {
	t.Helper()
	require.NoError(t, validateEveryScriptAndStyleHasNonce(html, expectedNonce))
}

func TestManagementNextUIHandler_TemplatedIndexNonceCsrfAndInitialData(t *testing.T) {
	handler, mockStore, mockUseCase := newNextUIHandlerTestFixtures(t, nextUiIndexTemplate, false)
	response := serveNextUIIndex(t, handler)
	body := readResponseBody(t, response)

	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Contains(t, response.Header.Get("Content-Type"), "text/html")
	require.Contains(t, response.Header.Get("Cache-Control"), "no-store")

	csrfToken := csrf.Token(httptest.NewRequest(http.MethodGet, "/", nil))
	require.Contains(t, body, `<meta name="csrf-token" content="`+csrfToken+`"`)

	scriptNonceRe := regexp.MustCompile(`nonce="([a-f0-9]+)"`)
	nonceMatches := scriptNonceRe.FindAllStringSubmatch(body, -1)
	require.NotEmpty(t, nonceMatches)
	expectedNonce := nonceMatches[0][1]
	assertEveryScriptAndStyleHasNonce(t, body, expectedNonce)

	initialDataRe := regexp.MustCompile(`(?s)<script id="initial-data"[^>]*>(.*?)</script>`)
	initialMatches := initialDataRe.FindStringSubmatch(body)
	require.Len(t, initialMatches, 2)

	var parsed liquidity_provider.ManagementTemplateData
	require.NoError(t, json.Unmarshal([]byte(strings.TrimSpace(initialMatches[1])), &parsed))
	assert.Equal(t, "http://localhost:8080", parsed.BaseUrl)
	assert.True(t, parsed.CredentialsSet)
	assert.Equal(t, "tb1qexample", parsed.BtcAddress)
	assert.Equal(t, "0xabc", parsed.RskAddress)

	mockStore.AssertExpectations(t)
	mockUseCase.AssertExpectations(t)
}

func TestManagementNextUIHandler_SecurityHeadersMatchLegacy(t *testing.T) {
	mockStore := new(mocks.StoreMock)
	mockStore.On("Get", mock.Anything, "lp-session").Return(nil, errors.New("no session"))

	templateData := liquidity_provider.ManagementTemplateData{
		CredentialsSet: false,
		BaseUrl:        "http://localhost:8080",
	}
	mockUseCase := new(mocks.GetManagementUiDataUseCaseMock)
	mockUseCase.On("Run", mock.Anything, false).Return(&liquidity_provider.ManagementTemplate{
		Name: liquidity_provider.ManagementLoginTemplate,
		Data: templateData,
	}, nil)

	env := environment.ManagementEnv{EnableSecurityHeaders: true}
	request := httptest.NewRequest(http.MethodGet, "/management", nil)
	legacyRecorder := httptest.NewRecorder()
	legacyHandler := NewManagementInterfaceHandler(env, mockStore, mockUseCase)
	legacyHandler(legacyRecorder, request)

	nextDist := fstest.MapFS{"index.html": &fstest.MapFile{Data: []byte(nextUiIndexTemplate)}}
	nextHandler := NewManagementNextUIHandler(nextDist, env, mockStore, mockUseCase)
	nextResponse := serveNextUIIndex(t, nextHandler)

	require.Equal(t, http.StatusOK, nextResponse.StatusCode)

	legacyCSP := normalizeCSP(legacyRecorder.Header().Get("Content-Security-Policy"))
	nextCSP := normalizeCSP(nextResponse.Header.Get("Content-Security-Policy"))
	assert.Equal(t, legacyCSP, nextCSP)
	assert.Equal(t, legacyRecorder.Header().Get("Strict-Transport-Security"), nextResponse.Header.Get("Strict-Transport-Security"))
	assert.Equal(t, legacyRecorder.Header().Get("X-Frame-Options"), nextResponse.Header.Get("X-Frame-Options"))
	assert.Equal(t, legacyRecorder.Header().Get("X-Content-Type-Options"), nextResponse.Header.Get("X-Content-Type-Options"))
}

func TestManagementNextUIHandler_SecurityHeadersDisabled(t *testing.T) {
	handler, _, _ := newNextUIHandlerTestFixtures(t, nextUiIndexTemplate, false)
	response := serveNextUIIndex(t, handler)

	assert.Empty(t, response.Header.Get("Content-Security-Policy"))
	assert.Empty(t, response.Header.Get("Strict-Transport-Security"))
	assert.Empty(t, response.Header.Get("X-Frame-Options"))
	assert.Empty(t, response.Header.Get("X-Content-Type-Options"))
}

func TestManagementNextUIHandler_NonceRegressionGuard(t *testing.T) {
	handler, _, _ := newNextUIHandlerTestFixtures(t, nextUiIndexMissingNonce, false)
	response := serveNextUIIndex(t, handler)
	body := readResponseBody(t, response)

	scriptNonceRe := regexp.MustCompile(`nonce="([a-f0-9]+)"`)
	nonceMatches := scriptNonceRe.FindAllStringSubmatch(body, -1)
	require.Empty(t, nonceMatches, "bad fixture should render scripts without nonce attributes")

	require.Contains(t, body, "<script>alert(1)</script>")
	require.Error(t, validateEveryScriptAndStyleHasNonce(body, "any-nonce-would-fail"))
}

func TestManagementNextUIHandler_UseCaseErrorUsesErrorTemplate(t *testing.T) {
	dist := fstest.MapFS{"index.html": &fstest.MapFile{Data: []byte(nextUiIndexTemplate)}}
	env := environment.ManagementEnv{EnableSecurityHeaders: false}
	mockStore := new(mocks.StoreMock)
	mockStore.On("Get", mock.Anything, "lp-session").Return(nil, errors.New("no session"))
	mockUseCase := new(mocks.GetManagementUiDataUseCaseMock)
	mockUseCase.On("Run", mock.Anything, false).Return(nil, errors.New("database error"))

	handler := NewManagementNextUIHandler(dist, env, mockStore, mockUseCase)
	response := serveNextUIIndex(t, handler)
	body := readResponseBody(t, response)

	require.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "Error opening management UI")
}

func TestManagementNextUIHandler_NormalizesTraversalToSafePath(t *testing.T) {
	handler, _, _ := newNextUIHandlerTestFixtures(t, nextUiIndexTemplate, false)

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
	handler, _, _ := newNextUIHandlerTestFixtures(t, nextUiIndexTemplate, false)

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

type errReadFile struct{}

func (errReadFile) Stat() (fs.FileInfo, error) { return nil, errors.New("stat error") }
func (errReadFile) Read([]byte) (int, error)   { return 0, errors.New("read error") }
func (errReadFile) Close() error               { return nil }

type failingReadFS struct {
	inner    fs.FS
	failPath string
}

func (f failingReadFS) Open(name string) (fs.File, error) {
	if name == f.failPath {
		return errReadFile{}, nil
	}
	return f.inner.Open(name)
}

func TestManagementNextUIHandler_PropagatesCopyError(t *testing.T) {
	inner := fstest.MapFS{
		"index.html":       &fstest.MapFile{Data: []byte(nextUiIndexTemplate)},
		"assets/broken.js": &fstest.MapFile{Data: []byte("x")},
	}
	dist := failingReadFS{inner: inner, failPath: "assets/broken.js"}
	handler, _, _ := newNextUIHandlerTestFixtures(t, nextUiIndexTemplate, false)
	handler.dist = dist

	request := httptest.NewRequest(http.MethodGet, "/management/next/assets/broken.js", nil)
	request = mux.SetURLVars(request, map[string]string{"path": "assets/broken.js"})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer func() { _ = response.Body.Close() }()
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "failed to write response")
}
