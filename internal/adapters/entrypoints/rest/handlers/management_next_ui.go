package handlers

import (
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/mux"
)

type ManagementNextUIHandler struct {
	dist fs.FS
}

func NewManagementNextUIHandler(dist fs.FS) http.Handler {
	return &ManagementNextUIHandler{dist: dist}
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

	// Default doc (and SPA entrypoint fallback).
	if cleanPath == "" || cleanPath == "." {
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
	responseWriter.Header().Set("Cache-Control", "no-store")
	responseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	handler.serveFileBody(responseWriter, request, "index.html")
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
