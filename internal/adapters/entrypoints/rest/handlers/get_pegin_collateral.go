package handlers

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type GetPeginCollateralUseCase interface {
	Run() (*entities.Wei, error)
}

// NewGetPeginCollateralHandler
// @Title Get PegIn Collateral
// @Description Get PegIn Collateral
// @Success 200 object pkg.GetCollateralResponse
// @Route /pegin/collateral [get]
func NewGetPeginCollateralHandler(useCase GetPeginCollateralUseCase) http.HandlerFunc {
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
