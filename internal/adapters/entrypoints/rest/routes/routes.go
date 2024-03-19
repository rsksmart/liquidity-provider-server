package routes

import (
	"github.com/gorilla/mux"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/middlewares"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func ConfigureRoutes(router *mux.Router, env environment.Environment, useCaseRegistry registry.UseCaseRegistry) {
	router.Use(middlewares.NewCorsMiddleware())
	captchaMiddleware := middlewares.NewCaptchaMiddleware(env.Captcha.Url, env.Captcha.Threshold, env.Captcha.Disabled, env.Captcha.SecretKey)

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
		log.Warn(
			"Server is running with the management API exposed. This interface " +
				"includes endpoints that must remain private at all cost. Please shut down " +
				"the server if you haven't configured the WAF properly as explained in documentation.",
		)
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
		router.Path("/configuration").Methods(http.MethodPost).
			HandlerFunc(handlers.NewSetGeneralConfigHandler(useCaseRegistry.SetGeneralConfigUseCase()))
		router.Path("/configuration").Methods(http.MethodGet).
			HandlerFunc(handlers.NewGetConfigurationHandler(useCaseRegistry.GetConfigurationUseCase()))
		router.Path("/pegin/configuration").Methods(http.MethodPost).
			HandlerFunc(handlers.NewSetPeginConfigHandler(useCaseRegistry.SetPeginConfigUseCase()))
		router.Path("/pegout/configuration").Methods(http.MethodPost).
			HandlerFunc(handlers.NewSetPegoutConfigHandler(useCaseRegistry.SetPegoutConfigUseCase()))
	}

	router.Methods(http.MethodOptions).HandlerFunc(handlers.NewOptionsHandler())
}
