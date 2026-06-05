package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	log "github.com/sirupsen/logrus"
)

type nextUiIndexData struct {
	CsrfToken       string
	ScriptNonce     string
	InitialDataJSON template.JS
}

type nextUiInitialData struct {
	LoggedIn bool                                      `json:"loggedIn"`
	Data     liquidity_provider.ManagementTemplateData `json:"data"`
}

type ManagementNextUIHandler struct {
	dist      fs.FS
	env       environment.ManagementEnv
	store     sessions.Store
	useCase   GetManagementUiDataUseCase
	indexTmpl *template.Template
}

func NewManagementNextUIHandler(
	dist fs.FS,
	env environment.ManagementEnv,
	store sessions.Store,
	useCase GetManagementUiDataUseCase,
) http.Handler {
	indexBytes, err := fs.ReadFile(dist, "index.html")
	if err != nil {
		panic("management next UI: missing embedded index.html: " + err.Error())
	}
	indexTmpl := template.Must(template.New("index.html").Parse(string(indexBytes)))

	return &ManagementNextUIHandler{
		dist:      dist,
		env:       env,
		store:     store,
		useCase:   useCase,
		indexTmpl: indexTmpl,
	}
}

func (handler *ManagementNextUIHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	routeVars := mux.Vars(request)
	requestedPath := routeVars["path"]

	// Normalize and block traversal attempts.
	cleanPath := path.Clean("/" + requestedPath)
	cleanPath = strings.TrimPrefix(cleanPath, "/")
	if cleanPath == ".." || strings.HasPrefix(cleanPath, "../") {
		http.Error(responseWriter, "invalid path", http.StatusBadRequest)
		return
	}

	// Default doc (and SPA entrypoint fallback). index.html is a Go template and must not be served raw.
	if cleanPath == "" || cleanPath == "." || cleanPath == "index.html" {
		handler.serveIndex(responseWriter, request)
		return
	}

	// Hash-named Vite assets live under dist/assets/.
	if strings.HasPrefix(cleanPath, "assets/") {
		handler.serveFile(responseWriter, request, cleanPath, true)
		return
	}

	// For client-side routes like /management/next/providers, serve index.html.
	// For any other file, attempt to serve it, otherwise fall back to index.html.
	if _, err := fs.Stat(handler.dist, cleanPath); err == nil {
		handler.serveFile(responseWriter, request, cleanPath, false)
		return
	}
	handler.serveIndex(responseWriter, request)
}

func (handler *ManagementNextUIHandler) serveIndex(responseWriter http.ResponseWriter, request *http.Request) {
	const errorGeneratingTemplate = "Error generating template: %v"

	session, err := handler.store.Get(request, cookies.ManagementSessionCookieName)
	loggedIn := err == nil && !session.IsNew
	result, err := handler.useCase.Run(request.Context(), loggedIn)
	if err != nil {
		log.Errorf(errorGeneratingTemplate, err)
		sendErrorTemplate(responseWriter)
		return
	}

	nonceRaw, err := utils.GetRandomBytes(nonceBytes)
	if err != nil {
		log.Errorf(errorGeneratingTemplate, err)
		sendErrorTemplate(responseWriter)
		return
	}
	nonce := hex.EncodeToString(nonceRaw)

	var jsonBuf bytes.Buffer
	encoder := json.NewEncoder(&jsonBuf)
	encoder.SetEscapeHTML(true)
	if err := encoder.Encode(nextUiInitialData{LoggedIn: loggedIn, Data: result.Data}); err != nil {
		log.Errorf(errorGeneratingTemplate, err)
		sendErrorTemplate(responseWriter)
		return
	}

	responseWriter.Header().Set("Cache-Control", "no-store")
	responseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	if handler.env.EnableSecurityHeaders {
		htmlTemplateSecurityHeaders(responseWriter, nonce)
	}

	data := nextUiIndexData{
		CsrfToken:   csrf.Token(request),
		ScriptNonce: nonce,
		// JSON is server-generated with SetEscapeHTML; template.JS marks safe for application/json script block.
		InitialDataJSON: template.JS(strings.TrimSpace(jsonBuf.String())), //nolint:gosec // G203
	}

	if err := handler.indexTmpl.Execute(responseWriter, data); err != nil {
		log.Errorf("Error executing next UI index template: %s", err.Error())
	}
}

func (handler *ManagementNextUIHandler) serveFile(responseWriter http.ResponseWriter, request *http.Request, filePath string, immutable bool) {
	ext := path.Ext(filePath)
	if ext != "" {
		if ct := mime.TypeByExtension(ext); ct != "" {
			responseWriter.Header().Set("Content-Type", ct)
		}
	}

	if immutable {
		responseWriter.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	} else {
		responseWriter.Header().Set("Cache-Control", "no-store")
	}

	handler.serveFileBody(responseWriter, request, filePath)
}

func (handler *ManagementNextUIHandler) serveFileBody(responseWriter http.ResponseWriter, request *http.Request, filePath string) {
	file, err := handler.dist.Open(filePath)
	if err != nil {
		http.NotFound(responseWriter, request)
		return
	}
	defer file.Close()

	if _, err := io.Copy(responseWriter, file); err != nil {
		http.Error(responseWriter, "failed to write response", http.StatusInternalServerError)
		return
	}
}
