package handlers

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type ChangeStatusUseCase interface {
	Run(newStatus bool) error
}

// NewChangeStatusHandler
// @Title Change Provider Status
// @Description Changes the status of the provider
// @Param ChangeStatusRequest body pkg.ChangeStatusRequest true "Change Provider Status Request"
// @Success 204 object
// @Route /providers/changeStatus [post]
func NewChangeStatusHandler(useCase ChangeStatusUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		request := pkg.ChangeStatusRequest{}
		if err = rest.DecodeRequest(w, req, &request); err != nil {
			return
		} else if err = rest.ValidateRequest(w, &request); err != nil {
			return
		}

		err = useCase.Run(*request.Status)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponse(w, http.StatusNoContent)
	}
}
