package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type SetCredentialsUseCase interface {
	Run(ctx context.Context, oldCredentials, newCredentials lp.Credentials) error
}

// NewSetCredentialsHandler
// @Title Set Login Credentials
// @Description Set new credentials to log into the Management API
// @Success 200 object
// @Route /management/credentials [post]
func NewSetCredentialsHandler(useCase SetCredentialsUseCase, sessionManager SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		request := pkg.CredentialsUpdateRequest{}
		if err = rest.DecodeRequest(w, req, &request); err != nil {
			return
		} else if err = rest.ValidateRequest(w, &request); err != nil {
			return
		}

		oldCredentials := lp.Credentials{Username: request.OldUsername, Password: request.OldPassword}
		newCredentials := lp.Credentials{Username: request.NewUsername, Password: request.NewPassword}

		err = useCase.Run(req.Context(), oldCredentials, newCredentials)
		if errors.Is(err, liquidity_provider.BadLoginError) {
			jsonErr := rest.NewErrorResponse(err.Error(), true)
			rest.JsonErrorResponse(w, http.StatusUnauthorized, jsonErr)
			return
		} else if errors.Is(err, utils.PasswordComplexityError) {
			jsonErr := rest.NewErrorResponse(err.Error(), true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		} else if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("unexpected login error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}

		if err = sessionManager.CloseSession(req, w); err != nil {
			return
		}
		rest.JsonResponse(w, http.StatusOK)
	}
}
