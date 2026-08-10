package handlers

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
)

// SessionManager handles session creation and closure for the management API
type SessionManager interface {
	CloseSession(req *http.Request, w http.ResponseWriter) error
	CreateSession(req *http.Request, w http.ResponseWriter) error
}

// CookieSessionManager is the default implementation of SessionManager using cookies
type CookieSessionManager struct {
	store cookies.SessionStore
}

// NewCookieSessionManager creates a new CookieSessionManager
func NewCookieSessionManager(store cookies.SessionStore) *CookieSessionManager {
	return &CookieSessionManager{store: store}
}

// CloseSession closes the current management session
func (m *CookieSessionManager) CloseSession(req *http.Request, w http.ResponseWriter) error {
	const errorMsg = "error closing session"
	err := m.store.Close(w, req)
	if err != nil {
		jsonErr := rest.NewErrorResponseWithDetails(errorMsg, rest.DetailsFromError(err), false)
		rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
		return err
	}
	return nil
}

// CreateSession creates a new management session
func (m *CookieSessionManager) CreateSession(req *http.Request, w http.ResponseWriter) error {
	const errorMsg = "session creation error"
	err := m.store.Create(w, req)
	if err != nil {
		jsonErr := rest.NewErrorResponseWithDetails(errorMsg, rest.DetailsFromError(err), false)
		rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
		return err
	}
	return nil
}
