package handlers

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"net/http"
)

// NewManagementLoginHandler
// @Title Management Login
// @Description Authenticate to start a Management API session
// @Success 200
// @Route /management/login [post]
func NewManagementLoginHandler(env environment.ManagementEnv, useCase *liquidity_provider.LoginUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		const errorMsg = "error closing old session"
		err := useCase.Run()
		if errors.Is(err, liquidity_provider.BadLoginError) {
			jsonErr := rest.NewErrorResponse(err.Error(), false)
			rest.JsonErrorResponse(w, http.StatusUnauthorized, jsonErr)
			return
		} else if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("unexpected login error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}

		cookieStore, err := cookies.GetSessionCookieStore(env)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(errorMsg, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}

		err = cookies.CloseManagementSession(&cookies.CloseSessionArgs{
			Store:   cookieStore,
			Request: req,
			Writer:  w,
		})
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(errorMsg, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}

		err = cookies.CreateManagementSession(&cookies.CreateSessionArgs{
			Store:   cookieStore,
			Env:     env,
			Request: req,
			Writer:  w,
		})
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("session creation error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponse(w, http.StatusOK)
	}
}
