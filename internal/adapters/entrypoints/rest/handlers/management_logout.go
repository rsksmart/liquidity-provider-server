package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"net/http"
)

// NewManagementLogoutHandler
// @Title Management Logout
// @Description Logout from the Management API session
// @Success 204 object
// @Route /management/logout [post]
func NewManagementLogoutHandler(env environment.ManagementEnv) http.HandlerFunc {
	// this handler doesn't use a use case because it doesn't involve any business logic
	return func(w http.ResponseWriter, req *http.Request) {
		if err := closeManagementSession(req, w, env); err != nil {
			return
		}
		rest.JsonResponse(w, http.StatusNoContent)
	}
}
