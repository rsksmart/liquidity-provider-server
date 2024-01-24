package routes

import (
	"github.com/gorilla/mux"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/middlewares"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"net/http"
)

func ConfigureRoutes(router *mux.Router, env environment.Environment, useCaseRegistry registry.UseCaseRegistry) {
	router.Use(middlewares.NewCorsMiddleware())
	captchaMiddleware := middlewares.NewCaptchaMiddleware(env.Captcha.Threshold, env.Captcha.Disabled, env.Captcha.SecretKey)

	router.Path("/health").Methods(http.MethodGet).HandlerFunc(handlers.NewHealthCheckHandler(useCaseRegistry.HealthUseCase()))
	router.Path("/getProviders").Methods(http.MethodGet).HandlerFunc(handlers.NewGetProvidersHandler(useCaseRegistry.GetProvidersUseCase()))
	router.Path("/pegin/getQuote").Methods(http.MethodPost).HandlerFunc(handlers.NewGetPeginQuoteHandler(useCaseRegistry.GetPeginQuoteUseCase()))
	router.Path("/pegin/acceptQuote").Methods(http.MethodPost).Handler(
		captchaMiddleware(
			handlers.NewAcceptPeginQuoteHandler(useCaseRegistry.GetAcceptPeginQuoteUseCase()),
		),
	)
	router.Path("/pegout/getQuotes").Methods(http.MethodPost).HandlerFunc(handlers.NewGetPegoutQuoteHandler(useCaseRegistry.GetPegoutQuoteUseCase()))
	router.Path("/pegout/acceptQuote").Methods(http.MethodPost).Handler(
		captchaMiddleware(
			handlers.NewAcceptPegoutQuoteHandler(useCaseRegistry.GetAcceptPegoutQuoteUseCase()),
		),
	)
	router.Path("/userQuotes").Methods(http.MethodGet).HandlerFunc(handlers.NewGetUserQuotesHandler(useCaseRegistry.GetUserDepositsUseCase()))
	router.Path("/providers/details").Methods(http.MethodGet).HandlerFunc(handlers.NewProviderDetailsHandler(useCaseRegistry.GetProviderDetailUseCase()))

	if env.EnableManagementApi {
		router.Path("/pegin/collateral").Methods(http.MethodGet).
			HandlerFunc(handlers.NewGetPeginCollateralHandler(useCaseRegistry.GetPeginCollateralUseCase()))
		router.Path("/pegin/addCollateral").Methods(http.MethodPost).
			HandlerFunc(handlers.NewAddPeginCollateralHandler(useCaseRegistry.AddPeginCollateralUseCase()))
		router.Path("/pegin/withdrawCollateral").Methods(http.MethodPost).
			HandlerFunc(handlers.NewWithdrawPeginCollateralHandler(useCaseRegistry.WithdrawPeginCollateralUseCase()))
		router.Path("/pegout/collateral").Methods(http.MethodGet).
			HandlerFunc(handlers.NewGetPegoutCollateralHandler(useCaseRegistry.GetPegoutCollateralUseCase()))
		router.Path("/pegout/addCollateral").Methods(http.MethodPost).
			HandlerFunc(handlers.NewAddPegoutCollateralHandler(useCaseRegistry.AddPegoutCollateralUseCase()))
		router.Path("/pegout/withdrawCollateral").Methods(http.MethodPost).
			HandlerFunc(handlers.NewWithdrawPegoutCollateralHandler(useCaseRegistry.WithdrawPegoutCollateralUseCase()))
		router.Path("/providers/changeStatus").Methods(http.MethodPost).
			HandlerFunc(handlers.NewChangeStatusHandler(useCaseRegistry.ChangeStatusUseCase()))
		router.Path("/providers/resignation").Methods(http.MethodPost).
			HandlerFunc(handlers.NewResignationHandler(useCaseRegistry.ResignationUseCase()))
	}

	router.Methods(http.MethodOptions).HandlerFunc(handlers.NewOptionsHandler())
}
