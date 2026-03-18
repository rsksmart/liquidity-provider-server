package handlers

import (
	"context"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type SetLiquidityRatioUseCase interface {
	Run(ctx context.Context, btcPercentage uint64) error
}

// NewSetLiquidityRatioHandler
// @Title Set Liquidity Ratio
// @Description Set the BTC/RBTC liquidity ratio. Included in the management API.
// @Param SetLiquidityRatioRequest body pkg.SetLiquidityRatioRequest true "New BTC percentage"
// @Success 204 object
// @Route /management/liquidity-ratio [post]
func NewSetLiquidityRatioHandler(useCase SetLiquidityRatioUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		request := &pkg.SetLiquidityRatioRequest{}
		if err = rest.DecodeRequest(w, req, request); err != nil {
			return
		} else if err = rest.ValidateRequest(w, request); err != nil {
			return
		}

		err = useCase.Run(req.Context(), request.BtcPercentage)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}

		rest.JsonResponse(w, http.StatusNoContent)
	}
}
