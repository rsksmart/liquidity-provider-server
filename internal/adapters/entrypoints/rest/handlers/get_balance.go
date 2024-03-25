package handlers

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
)

// NewGetLiquidityStatusHandler creates a new handler for fetching liquidity status
// @Title Fetch Liquidity Status
// @Description Fetches the available liquidity for both pegin & pegout
// @Success 200 {object} LiquidityStatus
// @Route /liquidity/status [get]
func NewGetLiquidityStatusHandler(useCase *liquidity_provider.LiquidityStatusUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		response, err := useCase.Run(req.Context())
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("failed to fetch liquidity status", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}