package handlers

import (
	"errors"

	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

// NewSetPegoutConfigHandler
// @Title Set Pegout Config
// @Description Set the configuration for the Pegout service. Included in the management API.
// @Param PegoutConfigurationRequest  body pkg.PegoutConfigurationRequest true "Specification of the thresholds for the PegOut service"
// @Success 204 object
// @Route /pegout/configuration [post]
func NewSetPegoutConfigHandler(useCase *liquidity_provider.SetPegoutConfigUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		request := &pkg.PegoutConfigurationRequest{}
		if err = rest.DecodeRequest(w, req, request); err != nil {
			return
		} else if err = rest.ValidateRequest(w, request); err != nil {
			return
		}

		err = useCase.Run(req.Context(), pkg.FromPegoutConfigurationDTO(request.Configuration))
		if err != nil {
			if errors.Is(err, usecases.TxBelowMinimumError) || errors.Is(err, usecases.NonPositiveWeiError) {
				jsonErr := rest.NewErrorResponseWithDetails("Validation error", rest.DetailsFromError(err), true)
				rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			} else {
				jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
				rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			}
			return
		}
		rest.JsonResponse(w, http.StatusNoContent)
	}
}
