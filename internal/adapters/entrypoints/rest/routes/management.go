package routes

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"net/http"
)

const LOGIN_PATH = "/management/login"

// nolint:funlen
func getManagementEndpoints(env environment.Environment, useCaseRegistry registry.UseCaseRegistry) []Endpoint {
	return []Endpoint{
		{
			Path:    "/pegin/collateral",
			Method:  http.MethodGet,
			Handler: handlers.NewGetPeginCollateralHandler(useCaseRegistry.GetPeginCollateralUseCase()),
		},
		{
			Path:    "/pegin/addCollateral",
			Method:  http.MethodPost,
			Handler: handlers.NewAddPeginCollateralHandler(useCaseRegistry.AddPeginCollateralUseCase()),
		},
		{
			Path:    "/pegin/withdrawCollateral",
			Method:  http.MethodPost,
			Handler: handlers.NewWithdrawPeginCollateralHandler(useCaseRegistry.WithdrawPeginCollateralUseCase()),
		},
		{
			Path:    "/pegout/collateral",
			Method:  http.MethodGet,
			Handler: handlers.NewGetPegoutCollateralHandler(useCaseRegistry.GetPegoutCollateralUseCase()),
		},
		{
			Path:    "/pegout/addCollateral",
			Method:  http.MethodPost,
			Handler: handlers.NewAddPegoutCollateralHandler(useCaseRegistry.AddPegoutCollateralUseCase()),
		},
		{
			Path:    "/pegout/withdrawCollateral",
			Method:  http.MethodPost,
			Handler: handlers.NewWithdrawPegoutCollateralHandler(useCaseRegistry.WithdrawPegoutCollateralUseCase()),
		},
		{
			Path:    "/providers/changeStatus",
			Method:  http.MethodPost,
			Handler: handlers.NewChangeStatusHandler(useCaseRegistry.ChangeStatusUseCase()),
		},
		{
			Path:    "/providers/resignation",
			Method:  http.MethodPost,
			Handler: handlers.NewResignationHandler(useCaseRegistry.ResignationUseCase()),
		},
		{
			Path:    "/configuration",
			Method:  http.MethodPost,
			Handler: handlers.NewSetGeneralConfigHandler(useCaseRegistry.SetGeneralConfigUseCase()),
		},
		{
			Path:    "/configuration",
			Method:  http.MethodGet,
			Handler: handlers.NewGetConfigurationHandler(useCaseRegistry.GetConfigurationUseCase()),
		},
		{
			Path:    "/pegin/configuration",
			Method:  http.MethodPost,
			Handler: handlers.NewSetPeginConfigHandler(useCaseRegistry.SetPeginConfigUseCase()),
		},
		{
			Path:    "/pegout/configuration",
			Method:  http.MethodPost,
			Handler: handlers.NewSetPegoutConfigHandler(useCaseRegistry.SetPegoutConfigUseCase()),
		},
		{
			Path:    LOGIN_PATH,
			Method:  http.MethodPost,
			Handler: handlers.NewManagementLoginHandler(env.Management, useCaseRegistry.LoginUseCase()),
		},
		{
			Path:    "/management/logout",
			Method:  http.MethodPost,
			Handler: handlers.NewManagementLogoutHandler(env.Management),
		},
		{
			Path:    "/management/credentials",
			Method:  http.MethodPost,
			Handler: handlers.NewSetCredentialsHandler(env.Management, useCaseRegistry.SetCredentialsUseCase()),
		},
	}
}
