package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type LoginUseCase interface {
	Run(ctx context.Context, credentials lp.Credentials) error
}

// NewManagementLoginHandler
// @Title Management Login
// @Description Authenticate to start a Management API session
// @Success 200 object
// @Route /management/login [post]
func NewManagementLoginHandler(useCase LoginUseCase, sessionManager SessionManager) http.HandlerFunc {
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

		if err = sessionManager.CloseSession(req, w); err != nil {
			return
		}

		if err = sessionManager.CreateSession(req, w); err != nil {
			return
		}

		rest.JsonResponse(w, http.StatusOK)
	}
}
