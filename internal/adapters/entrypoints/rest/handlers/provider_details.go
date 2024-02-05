package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
)

func NewProviderDetailsHandler(useCase *liquidity_provider.GetDetailUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		result, err := useCase.Run()
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("unknown error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		response := pkg.ProviderDetailResponse{
			SiteKey: result.SiteKey,
			Pegin: pkg.ProviderDetail{
				Fee:                   result.Pegin.Fee.Uint64(),
				MinTransactionValue:   result.Pegin.MinTransactionValue.Uint64(),
				MaxTransactionValue:   result.Pegin.MaxTransactionValue.Uint64(),
				RequiredConfirmations: result.Pegin.RequiredConfirmations,
			},
			Pegout: pkg.ProviderDetail{
				Fee:                   result.Pegout.Fee.Uint64(),
				MinTransactionValue:   result.Pegout.MinTransactionValue.Uint64(),
				MaxTransactionValue:   result.Pegout.MaxTransactionValue.Uint64(),
				RequiredConfirmations: result.Pegout.RequiredConfirmations,
			},
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
