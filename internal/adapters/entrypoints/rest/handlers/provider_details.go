package handlers

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

// NewProviderDetailsHandler
// @Title Provider detail
// @Description Returns the details of the provider that manages this instance of LPS
// @Success 200 object pkg.ProviderDetailResponse "Detail of the provider that manges this instance"
// @Route /providers/details [get]
func NewProviderDetailsHandler(useCase *liquidity_provider.GetDetailUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		result, err := useCase.Run(req.Context())
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		response := pkg.ProviderDetailResponse{
			SiteKey:               result.SiteKey,
			LiquidityCheckEnabled: result.LiquidityCheckEnabled,
			Pegin: pkg.ProviderDetail{
				FixedFee:              result.Pegin.FixedFee.Uint64(),
				PercentageFee:         result.Pegin.PercentageFee,
				MinTransactionValue:   result.Pegin.MinTransactionValue.Uint64(),
				MaxTransactionValue:   result.Pegin.MaxTransactionValue.Uint64(),
				RequiredConfirmations: result.Pegin.RequiredConfirmations,
			},
			Pegout: pkg.ProviderDetail{
				FixedFee:              result.Pegout.FixedFee.Uint64(),
				PercentageFee:         result.Pegout.PercentageFee,
				MinTransactionValue:   result.Pegout.MinTransactionValue.Uint64(),
				MaxTransactionValue:   result.Pegout.MaxTransactionValue.Uint64(),
				RequiredConfirmations: result.Pegout.RequiredConfirmations,
			},
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
