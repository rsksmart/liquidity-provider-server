package handlers

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
)

// SessionManager handles session creation and closure for the management API
type SessionManager interface {
	CloseSession(req *http.Request, w http.ResponseWriter) error
	CreateSession(req *http.Request, w http.ResponseWriter) error
}

// CookieSessionManager is the default implementation of SessionManager using cookies
type CookieSessionManager struct {
	env environment.ManagementEnv
}

// NewCookieSessionManager creates a new CookieSessionManager
func NewCookieSessionManager(env environment.ManagementEnv) *CookieSessionManager {
	return &CookieSessionManager{env: env}
}

// CloseSession closes the current management session
func (m *CookieSessionManager) CloseSession(req *http.Request, w http.ResponseWriter) error {
	const errorMsg = "error closing session"
	cookieStore, err := cookies.GetSessionCookieStore(m.env)
	if err != nil {
		jsonErr := rest.NewErrorResponseWithDetails(errorMsg, rest.DetailsFromError(err), false)
		rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
		return err
	}

	err = cookies.CloseManagementSession(&cookies.CloseSessionArgs{
		Store:   cookieStore,
		Request: req,
		Writer:  w,
	})
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
	cookieStore, err := cookies.GetSessionCookieStore(m.env)
	if err != nil {
		jsonErr := rest.NewErrorResponseWithDetails(errorMsg, rest.DetailsFromError(err), false)
		rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
		return err
	}

	err = cookies.CreateManagementSession(&cookies.CreateSessionArgs{
		Store:   cookieStore,
		Env:     m.env,
		Request: req,
		Writer:  w,
	})
	if err != nil {
		jsonErr := rest.NewErrorResponseWithDetails(errorMsg, rest.DetailsFromError(err), false)
		rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
		return err
	}
	return nil
}
