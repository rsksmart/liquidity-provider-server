package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	lp "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type GetLiquidityRatioUseCase interface {
	Run(ctx context.Context, proposedBtcPercentage uint64) (lp.LiquidityRatioDetail, error)
}

// NewGetLiquidityRatioHandler
// @Title Get Liquidity Ratio
// @Description Get the current liquidity ratio status. Optionally pass btcPercentage query param to preview a change. Included in the management API.
// @Param btcPercentage query int false "Proposed BTC percentage for preview"
// @Success 200 {object} pkg.LiquidityRatioResponse
// @Route /management/liquidity-ratio [get]
func NewGetLiquidityRatioHandler(useCase GetLiquidityRatioUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var proposedPercentage uint64
		if param := req.URL.Query().Get("btcPercentage"); param != "" {
			parsed, err := strconv.ParseUint(param, 10, 64)
			if err != nil {
				rest.ValidateRequestError(w, err)
				return
			}
			proposedPercentage = parsed
		}

		result, err := useCase.Run(req.Context(), proposedPercentage)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}

		response := pkg.ToLiquidityRatioResponse(result)
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
