package handlers

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

// NewManagementLoginHandler
// @Title Management Login
// @Description Authenticate to start a Management API session
// @Success 200 object
// @Route /management/login [post]
func NewManagementLoginHandler(env environment.ManagementEnv, useCase *liquidity_provider.LoginUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		const errorMsg = "session creation error"
		var err error
		request := pkg.LoginRequest{}
		if err = rest.DecodeRequest(w, req, &request); err != nil {
			return
		} else if err = rest.ValidateRequest(w, &request); err != nil {
			return
		}

		err = useCase.Run(req.Context(), lp.Credentials{
			Username: request.Username,
			Password: request.Password,
		})
		if errors.Is(err, liquidity_provider.BadLoginError) {
			jsonErr := rest.NewErrorResponse(err.Error(), false)
			rest.JsonErrorResponse(w, http.StatusUnauthorized, jsonErr)
			return
		} else if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(errorMsg, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		if err = closeManagementSession(req, w, env); err != nil {
			return
		}

		cookieStore, err := cookies.GetSessionCookieStore(env)
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
			jsonErr := rest.NewErrorResponseWithDetails(errorMsg, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponse(w, http.StatusOK)
	}
}
