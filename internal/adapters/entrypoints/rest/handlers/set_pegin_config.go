package handlers

import (
	"errors"

	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

// NewSetPeginConfigHandler
// @Title Set Pegin Config
// @Description Set the configuration for the Pegin service. Included in the management API.
// @Param PeginConfigurationRequest  body pkg.PeginConfigurationRequest true "Specification of the thresholds for the PegIn service"
// @Success 204 object
// @Route /pegin/configuration [post]
func NewSetPeginConfigHandler(useCase *liquidity_provider.SetPeginConfigUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		request := &pkg.PeginConfigurationRequest{}
		if err = rest.DecodeRequest(w, req, request); err != nil {
			return
		} else if err = rest.ValidateRequest(w, request); err != nil {
			return
		}

		err = useCase.Run(req.Context(), pkg.FromPeginConfigurationDTO(request.Configuration))
		if err != nil {
			// Check if this is a validation error
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
