package handlers_test

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
	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/routes"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type nextUiInitialData struct {
	LoggedIn bool                                      `json:"loggedIn"`
	Data     liquidity_provider.ManagementTemplateData `json:"data"`
}

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

var nextUIInitialDataScriptRe = regexp.MustCompile(`(?s)<script id="initial-data"[^>]*>(.*?)</script>`)

func parseNextUIInitialDataFromHTML(t *testing.T, html string) nextUiInitialData {
	t.Helper()

	matches := nextUIInitialDataScriptRe.FindStringSubmatch(html)
	require.Len(t, matches, 2)

	var parsed nextUiInitialData
	require.NoError(t, json.Unmarshal([]byte(strings.TrimSpace(matches[1])), &parsed))
	return parsed
}

func newNextUIHandlerWithLoggedInSession(
	t *testing.T,
	templateData liquidity_provider.ManagementTemplateData,
) (*handlers.ManagementNextUIHandler, *mocks.StoreMock, *mocks.GetManagementUiDataUseCaseMock) {
	t.Helper()

	dist := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte(nextUiIndexTemplate)},
	}
	env := environment.ManagementEnv{EnableSecurityHeaders: false}

	mockStore := new(mocks.StoreMock)
	validSession := sessions.NewSession(mockStore, "lp-session")
	validSession.IsNew = false
	mockStore.On("Get", mock.Anything, "lp-session").Return(validSession, nil)

	mockUseCase := new(mocks.GetManagementUiDataUseCaseMock)
	mockUseCase.On("Run", mock.Anything, true).Return(&liquidity_provider.ManagementTemplate{
		Name: liquidity_provider.ManagementUiTemplate,
		Data: templateData,
	}, nil)

	handler, ok := handlers.NewManagementNextUIHandler(dist, env, mockStore, mockUseCase).(*handlers.ManagementNextUIHandler)
	require.True(t, ok)
	return handler, mockStore, mockUseCase
}

func newNextUIHandlerTestFixtures(t *testing.T, indexHTML string) (*handlers.ManagementNextUIHandler, *mocks.StoreMock, *mocks.GetManagementUiDataUseCaseMock) {
	t.Helper()

	dist := fstest.MapFS{
		"index.html":            &fstest.MapFile{Data: []byte(indexHTML)},
		"assets/app-abc123.js":  &fstest.MapFile{Data: []byte("console.log('ok')")},
		"assets/app-abc123.css": &fstest.MapFile{Data: []byte("body{}")},
	}

	env := environment.ManagementEnv{EnableSecurityHeaders: false}
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

	handler, ok := handlers.NewManagementNextUIHandler(dist, env, mockStore, mockUseCase).(*handlers.ManagementNextUIHandler)
	require.True(t, ok)
	return handler, mockStore, mockUseCase
}

func serveNextUIIndex(t *testing.T, handler http.Handler) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, "/management/next/", nil)
	request = mux.SetURLVars(request, map[string]string{"path": ""})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func readResponseBody(t *testing.T, recorder *httptest.ResponseRecorder) string {
	t.Helper()
	body, err := io.ReadAll(recorder.Body)
	require.NoError(t, err)
	return string(body)
}

func newNextUIRouterTestServer(t *testing.T, handler http.Handler) (string, *http.Client) {
	t.Helper()

	router := mux.NewRouter()
	router.Path(routes.NextUiPath).Methods(http.MethodGet).Handler(handler)
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	return server.URL, &http.Client{}
}

func getNextUIAt(t *testing.T, client *http.Client, baseURL, path string) *http.Response {
	t.Helper()

	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseURL+path, nil)
	require.NoError(t, err)
	resp, err := client.Do(request)
	require.NoError(t, err)
	return resp
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
	handler, mockStore, mockUseCase := newNextUIHandlerTestFixtures(t, nextUiIndexTemplate)
	recorder := serveNextUIIndex(t, handler)
	body := readResponseBody(t, recorder)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Header().Get("Content-Type"), "text/html")
	require.Contains(t, recorder.Header().Get("Cache-Control"), "no-store")

	csrfToken := csrf.Token(httptest.NewRequest(http.MethodGet, "/", nil))
	require.Contains(t, body, `<meta name="csrf-token" content="`+csrfToken+`"`)

	scriptNonceRe := regexp.MustCompile(`nonce="([a-f0-9]+)"`)
	nonceMatches := scriptNonceRe.FindAllStringSubmatch(body, -1)
	require.NotEmpty(t, nonceMatches)
	expectedNonce := nonceMatches[0][1]
	assertEveryScriptAndStyleHasNonce(t, body, expectedNonce)

	parsed := parseNextUIInitialDataFromHTML(t, body)
	assert.False(t, parsed.LoggedIn)
	assert.Equal(t, "http://localhost:8080", parsed.Data.BaseUrl)
	assert.True(t, parsed.Data.CredentialsSet)
	assert.Equal(t, "tb1qexample", parsed.Data.BtcAddress)
	assert.Equal(t, "0xabc", parsed.Data.RskAddress)

	mockStore.AssertExpectations(t)
	mockUseCase.AssertExpectations(t)
}

func TestManagementNextUIHandler_TemplatedIndexLoggedInInitialData(t *testing.T) {
	templateData := liquidity_provider.ManagementTemplateData{
		CredentialsSet: true,
		BaseUrl:        "http://localhost:8080",
		BtcAddress:     "tb1qloggedin",
		RskAddress:     "0xloggedin",
	}
	handler, mockStore, mockUseCase := newNextUIHandlerWithLoggedInSession(t, templateData)
	recorder := serveNextUIIndex(t, handler)
	body := readResponseBody(t, recorder)

	require.Equal(t, http.StatusOK, recorder.Code)

	parsed := parseNextUIInitialDataFromHTML(t, body)
	assert.True(t, parsed.LoggedIn)
	assert.Equal(t, "tb1qloggedin", parsed.Data.BtcAddress)

	mockStore.AssertExpectations(t)
	mockUseCase.AssertExpectations(t)
}

func TestManagementNextUIHandler_SecurityHeadersUseNextUiPolicy(t *testing.T) {
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
	legacyHandler := handlers.NewManagementInterfaceHandler(env, mockStore, mockUseCase)
	legacyHandler(legacyRecorder, request)

	nextDist := fstest.MapFS{"index.html": &fstest.MapFile{Data: []byte(nextUiIndexTemplate)}}
	nextHandler := handlers.NewManagementNextUIHandler(nextDist, env, mockStore, mockUseCase)
	nextRecorder := serveNextUIIndex(t, nextHandler)

	require.Equal(t, http.StatusOK, nextRecorder.Code)

	nextCSP := nextRecorder.Header().Get("Content-Security-Policy")
	require.Contains(t, nextCSP, "style-src 'self' 'nonce-")
	require.Contains(t, nextCSP, "script-src 'self' 'nonce-")
	require.NotContains(t, nextCSP, "sha256-yr5DcAJJmu0m4Rv1KfUyA8AJj1t0kAJ1D2JuSBIT1DU=")

	assert.Equal(t, legacyRecorder.Header().Get("Strict-Transport-Security"), nextRecorder.Header().Get("Strict-Transport-Security"))
	assert.Equal(t, legacyRecorder.Header().Get("X-Frame-Options"), nextRecorder.Header().Get("X-Frame-Options"))
	assert.Equal(t, legacyRecorder.Header().Get("X-Content-Type-Options"), nextRecorder.Header().Get("X-Content-Type-Options"))
}

func TestManagementNextUIHandler_SecurityHeadersDisabled(t *testing.T) {
	handler, _, _ := newNextUIHandlerTestFixtures(t, nextUiIndexTemplate)
	recorder := serveNextUIIndex(t, handler)

	assert.Empty(t, recorder.Header().Get("Content-Security-Policy"))
	assert.Empty(t, recorder.Header().Get("Strict-Transport-Security"))
	assert.Empty(t, recorder.Header().Get("X-Frame-Options"))
	assert.Empty(t, recorder.Header().Get("X-Content-Type-Options"))
}

func TestManagementNextUIHandler_NonceRegressionGuard(t *testing.T) {
	handler, _, _ := newNextUIHandlerTestFixtures(t, nextUiIndexMissingNonce)
	recorder := serveNextUIIndex(t, handler)
	body := readResponseBody(t, recorder)

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

	handler := handlers.NewManagementNextUIHandler(dist, env, mockStore, mockUseCase)
	recorder := serveNextUIIndex(t, handler)
	body := readResponseBody(t, recorder)

	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, body, "Error opening management UI")
}

func TestManagementNextUIHandler_ServesTemplatedIndexForExplicitIndexHTMLPath(t *testing.T) {
	handler, mockStore, mockUseCase := newNextUIHandlerTestFixtures(t, nextUiIndexTemplate)

	request := httptest.NewRequest(http.MethodGet, "/management/next/index.html", nil)
	request = mux.SetURLVars(request, map[string]string{"path": "index.html"})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	body := readResponseBody(t, recorder)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Header().Get("Content-Type"), "text/html")
	require.Contains(t, body, `<meta name="csrf-token"`)
	require.Contains(t, body, `<script id="initial-data"`)
	mockStore.AssertExpectations(t)
	mockUseCase.AssertExpectations(t)
}

func TestManagementNextUIHandler_NormalizesTraversalToSafePath(t *testing.T) {
	handler, _, _ := newNextUIHandlerTestFixtures(t, nextUiIndexTemplate)

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
	handler, _, _ := newNextUIHandlerTestFixtures(t, nextUiIndexTemplate)
	baseURL, client := newNextUIRouterTestServer(t, handler)

	t.Run("index without trailing slash", func(t *testing.T) {
		resp := getNextUIAt(t, client, baseURL, "/management/next")
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Contains(t, resp.Header.Get("Cache-Control"), "no-store")
	})

	t.Run("index", func(t *testing.T) {
		resp := getNextUIAt(t, client, baseURL, "/management/next/")
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Contains(t, resp.Header.Get("Cache-Control"), "no-store")
	})

	t.Run("hashed asset immutable", func(t *testing.T) {
		resp := getNextUIAt(t, client, baseURL, "/management/next/assets/app-abc123.js")
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "public, max-age=31536000, immutable", resp.Header.Get("Cache-Control"))
	})

	t.Run("existing non-asset file no-store", func(t *testing.T) {
		resp := getNextUIAt(t, client, baseURL, "/management/next/robots.txt")
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Contains(t, resp.Header.Get("Cache-Control"), "no-store")
	})

	t.Run("missing file falls back to index", func(t *testing.T) {
		resp := getNextUIAt(t, client, baseURL, "/management/next/some/client/route")
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Contains(t, resp.Header.Get("Content-Type"), "text/html")
	})

	t.Run("directory path falls back to index", func(t *testing.T) {
		resp := getNextUIAt(t, client, baseURL, "/management/next/assets")
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
	env := environment.ManagementEnv{EnableSecurityHeaders: false}
	mockStore := new(mocks.StoreMock)
	mockStore.On("Get", mock.Anything, "lp-session").Return(nil, errors.New("no session"))
	mockUseCase := new(mocks.GetManagementUiDataUseCaseMock)
	mockUseCase.On("Run", mock.Anything, false).Return(&liquidity_provider.ManagementTemplate{
		Name: liquidity_provider.ManagementUiTemplate,
		Data: liquidity_provider.ManagementTemplateData{
			CredentialsSet: true,
			BaseUrl:        "http://localhost:8080",
		},
	}, nil)
	handler := handlers.NewManagementNextUIHandler(dist, env, mockStore, mockUseCase)

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
