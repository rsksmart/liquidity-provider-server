package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	lpUseCase "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type GetAvailableLiquidityUseCase interface {
	Run(ctx context.Context) (liquidity_provider.AvailableLiquidity, error)
}

// NewGetAvailableLiquidityHandler
// @Title Fetch Available Liquidity
// @Description Fetches the available liquidity for both PegIn and PegOut operations.
// This might be disabled by the liquidity provider for privacy reasons.
// @Success 200 {object} pkg.AvailableLiquidityDTO
// @Route /providers/liquidity [get]
func NewGetAvailableLiquidityHandler(useCase GetAvailableLiquidityUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		result, err := useCase.Run(req.Context())
		if errors.Is(err, lpUseCase.LiquidityCheckNotEnabledError) {
			rest.JsonErrorResponse(w, http.StatusForbidden, rest.NewErrorResponse(err.Error(), false))
			return
		} else if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		response := pkg.ToAvailableLiquidityDTO(result)
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
