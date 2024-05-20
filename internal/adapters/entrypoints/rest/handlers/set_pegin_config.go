package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
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

		err = useCase.Run(req.Context(), *request.Configuration)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		} else {
			rest.JsonResponse(w, http.StatusNoContent)
		}
	}
}
