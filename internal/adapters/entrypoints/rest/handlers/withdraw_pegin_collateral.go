package handlers

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"net/http"
)

// NewWithdrawPeginCollateralHandler
// @Title Withdraw PegIn Collateral
// @Description Withdraw PegIn collateral of a resigned LP
// @Route /pegin/withdrawCollateral [post]
// @Success 204 object
func NewWithdrawPeginCollateralHandler(useCase *pegin.WithdrawCollateralUseCase) http.HandlerFunc {
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
