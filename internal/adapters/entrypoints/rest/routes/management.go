package routes

import (
	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/assets"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"net/http"
)

const (
	LoginPath  = "/management/login"
	UiPath     = "/management"
	StaticPath = "/static/{file}"
	IconPath   = "/favicon.ico"
)

var AllowedPaths = [...]string{LoginPath, UiPath, StaticPath, IconPath}

// nolint:funlen
func getManagementEndpoints(env environment.Environment, useCaseRegistry registry.UseCaseRegistry, store sessions.Store) []Endpoint {
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
			Path:    "/providers/withdrawCollateral",
			Method:  http.MethodPost,
			Handler: handlers.NewWithdrawCollateralHandler(useCaseRegistry.WithdrawCollateralUseCase()),
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
			Path:    LoginPath,
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
		{
			Path:    UiPath,
			Method:  http.MethodGet,
			Handler: handlers.NewManagementInterfaceHandler(env.Management, store, useCaseRegistry.GetManagementUiDataUseCase()),
		},
		{
			Path:    StaticPath,
			Method:  http.MethodGet,
			Handler: http.FileServer(http.FS(assets.FileSystem)),
		},
		{
			Path:    IconPath,
			Method:  http.MethodGet,
			Handler: http.FileServer(http.FS(assets.FileSystem)),
		},
	}
}
