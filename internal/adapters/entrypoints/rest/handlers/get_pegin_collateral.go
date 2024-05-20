package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

// NewGetPeginCollateralHandler
// @Title Get PegIn Collateral
// @Description Get PegIn Collateral
// @Success 200 object pkg.GetCollateralResponse
// @Route /pegin/collateral [get]
func NewGetPeginCollateralHandler(useCase *pegin.GetCollateralUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		collateral, err := useCase.Run()
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &pkg.GetCollateralResponse{Collateral: collateral.Uint64()})
	}
}
