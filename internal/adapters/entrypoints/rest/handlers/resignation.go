package handlers

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"net/http"
)

type ResignUseCase interface {
	Run() error
}

// NewResignationHandler
// @Title Provider resignation
// @Description Provider stops being a liquidity provider
// @Route /providers/resignation [post]
// @Success 204 object
func NewResignationHandler(useCase ResignUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := useCase.Run()
		if errors.Is(err, blockchain.ContractPausedError) {
			jsonErr := rest.NewErrorResponseWithDetails("protocol is paused", rest.DetailsFromError(err), true)
			rest.JsonErrorResponse(w, http.StatusServiceUnavailable, jsonErr)
			return
		} else if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponse(w, http.StatusNoContent)
	}
}
