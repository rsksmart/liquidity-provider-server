package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"net/http"
)

// NewResignationHandler
// @Title Provider resignation
// @Description Provider stops being a liquidity provider
// @Route /provider/resignation [post]
// @Success 204 object
func NewResignationHandler(useCase *liquidity_provider.ResignUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := useCase.Run()
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponse(w, http.StatusNoContent)
	}
}
