package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
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

		err = useCase.Run(req.Context(), *request.Configuration)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("unknown error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		} else {
			rest.JsonResponse(w, http.StatusNoContent)
		}
	}
}
