package handlers

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type GetPegoutCollateralUseCase interface {
	Run() (*entities.Wei, error)
}

// NewGetPegoutCollateralHandler
// @Title Get PegOut Collateral
// @Description Get PegOut Collateral
// @Success 200 object pkg.GetCollateralResponse
// @Route /pegout/collateral [get]
func NewGetPegoutCollateralHandler(useCase GetPegoutCollateralUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		collateral, err := useCase.Run()
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &pkg.GetCollateralResponse{Collateral: collateral.AsBigInt()})
	}
}
