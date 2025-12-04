package handlers

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type GetProvidersUseCase interface {
	Run() ([]liquidity_provider.RegisteredLiquidityProvider, error)
}

// NewGetProvidersHandler
// @Title Get Providers
// @Description Returns a list of providers.
// @Success 200 array pkg.LiquidityProvider
// @Route /getProviders [get]
func NewGetProvidersHandler(useCase GetProvidersUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		providers, err := useCase.Run()
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}

		result := make([]pkg.LiquidityProvider, 0)
		for _, provider := range providers {
			result = append(result,
				pkg.LiquidityProvider{
					Id:           provider.Id,
					Provider:     provider.Address,
					Name:         provider.Name,
					ApiBaseUrl:   provider.ApiBaseUrl,
					Status:       provider.Status,
					ProviderType: string(provider.ProviderType),
				},
			)
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &result)
	}
}
