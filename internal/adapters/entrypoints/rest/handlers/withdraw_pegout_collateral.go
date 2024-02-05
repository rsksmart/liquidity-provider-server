package handlers

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"net/http"
)

func NewWithdrawPegoutCollateralHandler(useCase *pegout.WithdrawCollateralUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := useCase.Run()
		if errors.Is(err, usecases.ProviderNotResignedError) {
			jsonErr := rest.NewErrorResponseWithDetails("unknown error", rest.DetailsFromError(err), true)
			rest.JsonErrorResponse(w, http.StatusConflict, jsonErr)
		} else if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("unknown error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
		} else {
			rest.JsonResponse(w, http.StatusNoContent)
		}
	}
}
