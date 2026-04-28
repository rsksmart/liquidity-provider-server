package handlers

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
)

// NewManagementLogoutHandler
// @Title Management Logout
// @Description Logout from the Management API session
// @Success 204 object
// @Route /management/logout [post]
func NewManagementLogoutHandler(sessionManager SessionManager) http.HandlerFunc {
	// this handler doesn't use a use case because it doesn't involve any business logic
	return func(w http.ResponseWriter, req *http.Request) {
		if err := sessionManager.CloseSession(req, w); err != nil {
			return
		}
		rest.JsonResponse(w, http.StatusNoContent)
	}
}
