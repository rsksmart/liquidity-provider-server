package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

// NewGetPegoutCollateralHandler
// @Title Get PegOut Collateral
// @Description Get PegOut Collateral
// @Success 200 object pkg.GetCollateralResponse
// @Route /pegout/collateral [get]
func NewGetPegoutCollateralHandler(useCase *pegout.GetCollateralUseCase) http.HandlerFunc {
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
