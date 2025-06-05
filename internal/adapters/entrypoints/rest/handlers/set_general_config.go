package handlers

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

type SetGeneralConfigUseCase interface {
	Run(ctx context.Context, config liquidity_provider.GeneralConfiguration) error
}

// NewSetGeneralConfigHandler
// @Title Set General Config
// @Description Set general configurations of the server. Included in the management API.
// @Param GeneralConfigurationRequest  body pkg.GeneralConfigurationRequest true "General parameters for the quote computation"
// @Success 204 object
// @Route /configuration [post]
func NewSetGeneralConfigHandler(useCase SetGeneralConfigUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		request := &pkg.GeneralConfigurationRequest{}
		if err = rest.DecodeRequest(w, req, request); err != nil {
			return
		} else if err = rest.ValidateRequest(w, request); err != nil {
			return
		}

		config, err := pkg.FromGeneralConfigurationDTO(request.Configuration)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("Invalid configuration", rest.DetailsFromError(err), true)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		err = useCase.Run(req.Context(), config)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		} else {
			rest.JsonResponse(w, http.StatusNoContent)
		}
	}
}
