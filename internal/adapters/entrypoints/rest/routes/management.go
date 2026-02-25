package routes

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/assets"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
)

const (
	LoginPath          = "/management/login"
	UiPath             = "/management"
	ManualApprovalPath = "/management/manual-approval"
	StaticPath         = "/static/{file}"
	IconPath           = "/favicon.ico"
)

var AllowedPaths = [...]string{LoginPath, UiPath, StaticPath, IconPath}

// nolint:funlen
func GetManagementEndpoints(env environment.Environment, useCaseRegistry registry.UseCaseRegistry, store sessions.Store) []Endpoint {
	sessionManager := handlers.NewCookieSessionManager(env.Management)
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
			Path:   "/reports/summaries",
			Method: http.MethodGet,
			Handler: handlers.NewGetReportSummariesHandler(
				handlers.SingleFlightGroup,
				handlers.SummariesReportSingleFlightKey,
				useCaseRegistry.SummariesUseCase(),
			),
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
			Path:   "/reports/pegin",
			Method: http.MethodGet,
			Handler: handlers.NewGetReportsPeginHandler(
				handlers.SingleFlightGroup,
				handlers.PegInReportSingleFlightKey,
				useCaseRegistry.GetPeginReportUseCase(),
			),
		},
		{
			Path:   "/reports/pegout",
			Method: http.MethodGet,
			Handler: handlers.NewGetReportsPegoutHandler(
				handlers.SingleFlightGroup,
				handlers.PegOutReportSingleFlightKey,
				useCaseRegistry.GetPegoutReportUseCase(),
			),
		},
		{
			Path:   "/reports/revenue",
			Method: http.MethodGet,
			Handler: handlers.NewGetReportsRevenueHandler(
				handlers.SingleFlightGroup,
				handlers.RevenueReportSingleFlightKey,
				useCaseRegistry.GetRevenueReportUseCase(),
			),
		},
		{
			Path:    "/reports/assets",
			Method:  http.MethodGet,
			Handler: handlers.NewGetReportsAssetsHandler(useCaseRegistry.GetAssetsReportUseCase()),
		},
		{
			Path:    "/reports/transactions",
			Method:  http.MethodGet,
			Handler: handlers.NewGetReportsTransactionHandler(useCaseRegistry.GetTransactionsReportUseCase()),
		},
		{
			Path:    LoginPath,
			Method:  http.MethodPost,
			Handler: handlers.NewManagementLoginHandler(useCaseRegistry.LoginUseCase(), sessionManager),
		},
		{
			Path:    "/management/logout",
			Method:  http.MethodPost,
			Handler: handlers.NewManagementLogoutHandler(sessionManager),
		},
		{
			Path:    "/management/credentials",
			Method:  http.MethodPost,
			Handler: handlers.NewSetCredentialsHandler(useCaseRegistry.SetCredentialsUseCase(), sessionManager),
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
		{
			Path:    "/management/trusted-accounts",
			Method:  http.MethodGet,
			Handler: handlers.NewGetTrustedAccountsHandler(useCaseRegistry.GetTrustedAccountsUseCase()),
		},
		{
			Path:    "/management/trusted-accounts",
			Method:  http.MethodPost,
			Handler: handlers.NewAddTrustedAccountHandler(useCaseRegistry.AddTrustedAccountUseCase()),
		},
		{
			Path:    "/management/trusted-accounts",
			Method:  http.MethodPut,
			Handler: handlers.NewUpdateTrustedAccountHandler(useCaseRegistry.UpdateTrustedAccountUseCase()),
		},
		{
			Path:    "/management/trusted-accounts",
			Method:  http.MethodDelete,
			Handler: handlers.NewDeleteTrustedAccountHandler(useCaseRegistry.DeleteTrustedAccountUseCase()),
		},
		{
			Path:    ManualApprovalPath,
			Method:  http.MethodGet,
			Handler: handlers.NewManagementInterfaceHandler(env.Management, store, useCaseRegistry.GetManagementUiDataUseCase(), liquidity_provider.ManagementManualApprovalTemplate),
		},
		{
			Path:    ManualApprovalPath + "/pending",
			Method:  http.MethodGet,
			Handler: handlers.NewGetPendingTransactionsHandler(),
		},
		{
			Path:    ManualApprovalPath + "/history",
			Method:  http.MethodGet,
			Handler: handlers.NewGetHistoryHandler(),
		},
		{
			Path:    ManualApprovalPath + "/approve",
			Method:  http.MethodPost,
			Handler: handlers.NewApproveTransactionsHandler(),
		},
		{
			Path:    ManualApprovalPath + "/deny",
			Method:  http.MethodPost,
			Handler: handlers.NewDenyTransactionsHandler(),
		},
	}
}
