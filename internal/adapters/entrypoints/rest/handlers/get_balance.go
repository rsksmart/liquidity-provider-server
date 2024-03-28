package handlers

import (
	"errors"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

// NewGetLiquidityStatusHandler creates a new handler for fetching liquidity status
// @Title Fetch Liquidity Status
// @Description Fetches the available liquidity for both pegin & pegout
// @Success 200 {object} pkg.LiquidityStatus
// @Route /liquidity/status [get]
func NewGetLiquidityStatusHandler(useCase *liquidity_provider.LiquidityStatusUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		status, err := useCase.Run(req.Context())
		if err != nil {
			if errors.Is(err, usecases.PublicLiquidityCheckDisabledError) {
				jsonErr := rest.NewErrorResponseWithDetails(usecases.PublicLiquidityCheckDisabledError.Error(), rest.DetailsFromError(err), false)
				rest.JsonErrorResponse(w, http.StatusForbidden, jsonErr)
			} else {
				jsonErr := rest.NewErrorResponseWithDetails(usecases.PublicLiquidityCheckError.Error(), rest.DetailsFromError(err), false)
				rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			}
			return
		}
		response := pkg.LiquidityStatus{
			Available: pkg.Available{
				Pegin:  status.Available.Pegin,
				Pegout: status.Available.Pegout,
			},
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
