package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"net/http"
)

func closeManagementSession(req *http.Request, w http.ResponseWriter, env environment.ManagementEnv) error {
	const errorMsg = "error closing session"
	cookieStore, err := cookies.GetSessionCookieStore(env)
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
