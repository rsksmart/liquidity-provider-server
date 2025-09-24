package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
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
		peginFeePercentage, _ := result.Pegin.FeePercentage.Native().Float64()
		pegoutFeePercentage, _ := result.Pegout.FeePercentage.Native().Float64()
		response := pkg.ProviderDetailResponse{
			SiteKey:               result.SiteKey,
			LiquidityCheckEnabled: result.LiquidityCheckEnabled,
			UsingSegwitFederation: result.UsingSegwitFederation,
			Pegin: pkg.ProviderDetail{
				Fee:                   result.Pegin.FixedFee.AsBigInt(),
				FixedFee:              result.Pegin.FixedFee.AsBigInt(),
				FeePercentage:         peginFeePercentage,
				MinTransactionValue:   result.Pegin.MinTransactionValue.AsBigInt(),
				MaxTransactionValue:   result.Pegin.MaxTransactionValue.AsBigInt(),
				RequiredConfirmations: result.Pegin.RequiredConfirmations,
			},
			Pegout: pkg.ProviderDetail{
				Fee:                   result.Pegout.FixedFee.AsBigInt(),
				FixedFee:              result.Pegout.FixedFee.AsBigInt(),
				FeePercentage:         pegoutFeePercentage,
				MinTransactionValue:   result.Pegout.MinTransactionValue.AsBigInt(),
				MaxTransactionValue:   result.Pegout.MaxTransactionValue.AsBigInt(),
				RequiredConfirmations: result.Pegout.RequiredConfirmations,
			},
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
