package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"net/http"
)

// NewManagementLogoutHandler
// @Title Management Logout
// @Description Logout from the Management API session
// @Success 204
// @Route /management/logout [post]
func NewManagementLogoutHandler(env environment.ManagementEnv) http.HandlerFunc {
	// this handler doesn't use a use case because it doesn't involve any business logic
	return func(w http.ResponseWriter, req *http.Request) {
		const errorMsg = "logout error"
		store, err := cookies.GetSessionCookieStore(env)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(errorMsg, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}

		err = cookies.CloseManagementSession(&cookies.CloseSessionArgs{
			Store:   store,
			Request: req,
			Writer:  w,
		})
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(errorMsg, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponse(w, http.StatusNoContent)
	}
}
