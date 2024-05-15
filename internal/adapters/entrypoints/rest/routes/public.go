package routes

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/registry"
	"net/http"
)

type PublicEndpoint struct {
	Endpoint
	RequiresCaptcha bool
}

// nolint:funlen
func getPublicEndpoints(useCaseRegistry registry.UseCaseRegistry) []PublicEndpoint {
	return []PublicEndpoint{
		{
			Endpoint: Endpoint{
				Path:    "/health",
				Method:  http.MethodGet,
				Handler: handlers.NewHealthCheckHandler(useCaseRegistry.HealthUseCase()),
			},
		},
		{
			Endpoint: Endpoint{
				Path:    "/getProviders",
				Method:  http.MethodGet,
				Handler: handlers.NewGetProvidersHandler(useCaseRegistry.GetProvidersUseCase()),
			},
		},
		{
			Endpoint: Endpoint{
				Path:    "/pegin/getQuote",
				Method:  http.MethodPost,
				Handler: handlers.NewGetPeginQuoteHandler(useCaseRegistry.GetPeginQuoteUseCase()),
			},
		},
		{
			Endpoint: Endpoint{
				Path:    "/pegin/acceptQuote",
				Method:  http.MethodPost,
				Handler: handlers.NewAcceptPeginQuoteHandler(useCaseRegistry.GetAcceptPeginQuoteUseCase()),
			},
			RequiresCaptcha: true,
		},
		{
			Endpoint: Endpoint{
				Path:    "/pegout/getQuotes",
				Method:  http.MethodPost,
				Handler: handlers.NewGetPegoutQuoteHandler(useCaseRegistry.GetPegoutQuoteUseCase()),
			},
		},
		{
			Endpoint: Endpoint{
				Path:    "/pegout/acceptQuote",
				Method:  http.MethodPost,
				Handler: handlers.NewAcceptPegoutQuoteHandler(useCaseRegistry.GetAcceptPegoutQuoteUseCase()),
			},
			RequiresCaptcha: true,
		},
		{
			Endpoint: Endpoint{
				Path:    "/userQuotes",
				Method:  http.MethodGet,
				Handler: handlers.NewGetUserQuotesHandler(useCaseRegistry.GetUserDepositsUseCase()),
			},
		},
		{
			Endpoint: Endpoint{
				Path:    "/providers/details",
				Method:  http.MethodGet,
				Handler: handlers.NewProviderDetailsHandler(useCaseRegistry.GetProviderDetailUseCase()),
			},
		},
		{
			Endpoint: Endpoint{
				Path:    "/pegin/status",
				Method:  http.MethodGet,
				Handler: handlers.NewGetPeginQuoteStatusHandler(useCaseRegistry.GetPeginStatusUseCase()),
			},
		},
		{
			Endpoint: Endpoint{
				Path:    "/pegout/status",
				Method:  http.MethodGet,
				Handler: handlers.NewGetPegoutQuoteStatusHandler(useCaseRegistry.GetPegoutStatusUseCase()),
			},
		},
	}
}
