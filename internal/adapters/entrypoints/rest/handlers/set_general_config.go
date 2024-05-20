package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

// NewSetGeneralConfigHandler
// @Title Set General Config
// @Description Set general configurations of the server. Included in the management API.
// @Param GeneralConfigurationRequest  body pkg.GeneralConfigurationRequest true "General parameters for the quote computation"
// @Success 204 object
// @Route /configuration [post]
func NewSetGeneralConfigHandler(useCase *liquidity_provider.SetGeneralConfigUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		request := &pkg.GeneralConfigurationRequest{}
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
