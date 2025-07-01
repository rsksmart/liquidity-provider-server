package registry

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
)

var signingHashFunction = crypto.Keccak256

type UseCaseRegistry struct {
	getPeginQuoteUseCase          *pegin.GetQuoteUseCase
	registerProviderUseCase       *liquidity_provider.RegistrationUseCase
	callForUserUseCase            *pegin.CallForUserUseCase
	registerPeginUseCase          *pegin.RegisterPeginUseCase
	acceptPeginQuoteUseCase       *pegin.AcceptQuoteUseCase
	getWatchedPeginQuoteUseCase   *watcher.GetWatchedPeginQuoteUseCase
	expiredPeginQuoteUseCase      *pegin.ExpiredPeginQuoteUseCase
	cleanExpiredQuotesUseCase     *watcher.CleanExpiredQuotesUseCase
	getProviderDetailUseCase      *liquidity_provider.GetDetailUseCase
	getWatchedPegoutQuoteUseCase  *watcher.GetWatchedPegoutQuoteUseCase
	expiredPegoutUseCase          *pegout.ExpiredPegoutQuoteUseCase
	sendPegoutUseCase             *pegout.SendPegoutUseCase
	updatePegoutDepositUseCase    *watcher.UpdatePegoutQuoteDepositUseCase
	initPegoutDepositCacheUseCase *pegout.InitPegoutDepositCacheUseCase
	refundPegoutUseCase           *pegout.RefundPegoutUseCase
	getPegoutQuoteUseCase         *pegout.GetQuoteUseCase
	acceptPegoutQuoteUseCase      *pegout.AcceptQuoteUseCase
	getUserDepositsUseCase        *pegout.GetUserDepositsUseCase
	liquidityCheckUseCase         *liquidity_provider.CheckLiquidityUseCase
	penalizationAlertUseCase      *liquidity_provider.PenalizationAlertUseCase
	getProvidersUseCase           *liquidity_provider.GetProvidersUseCase
	getPeginCollateralUseCase     *pegin.GetCollateralUseCase
	getPegoutCollateralUseCase    *pegout.GetCollateralUseCase
	withdrawCollateralUseCase     *liquidity_provider.WithdrawCollateralUseCase
	healthUseCase                 *usecases.HealthUseCase
	resignUseCase                 *liquidity_provider.ResignUseCase
	changeStatusUseCase           *liquidity_provider.ChangeStatusUseCase
	addPeginCollateralUseCase     *pegin.AddCollateralUseCase
	addPegoutCollateralUseCase    *pegout.AddCollateralUseCase
	setPeginConfigUseCase         *liquidity_provider.SetPeginConfigUseCase
	setPegoutConfigUseCase        *liquidity_provider.SetPegoutConfigUseCase
	setGeneralConfigUseCase       *liquidity_provider.SetGeneralConfigUseCase
	getConfigurationUseCase       *liquidity_provider.GetConfigUseCase
	loginUseCase                  *liquidity_provider.LoginUseCase
	setCredentialsUseCase         *liquidity_provider.SetCredentialsUseCase
	defaultCredentialsUseCase     *liquidity_provider.GenerateDefaultCredentialsUseCase
	getManagementUiDataUseCase    *liquidity_provider.GetManagementUiDataUseCase
	bridgePegoutUseCase           *pegout.BridgePegoutUseCase
	peginStatusUseCase            *pegin.StatusUseCase
	pegoutStatusUseCase           *pegout.StatusUseCase
	availableLiquidityUseCase     *liquidity_provider.GetAvailableLiquidityUseCase
	updatePeginDepositUseCase     *watcher.UpdatePeginDepositUseCase
	getServerInfoUseCase          *liquidity_provider.ServerInfoUseCase
	getAssetReportUseCase         *reports.GetAssetReportUseCase
}

// NewUseCaseRegistry
// nolint:funlen
func NewUseCaseRegistry(
	env environment.Environment,
	rskRegistry *Rootstock,
	btcRegistry *Bitcoin,
	databaseRegistry *Database,
	liquidityProvider *dataproviders.LocalLiquidityProvider,
	messaging *Messaging,
	mutexes entities.ApplicationMutexes,
) *UseCaseRegistry {
	return &UseCaseRegistry{
		getPeginQuoteUseCase: pegin.NewGetQuoteUseCase(
			messaging.Rpc,
			rskRegistry.Contracts,
			databaseRegistry.PeginRepository,
			liquidityProvider,
			liquidityProvider,
			env.Rsk.FeeCollectorAddress,
		),
		registerProviderUseCase: liquidity_provider.NewRegistrationUseCase(
			rskRegistry.Contracts,
			liquidityProvider,
		),
		callForUserUseCase: pegin.NewCallForUserUseCase(
			rskRegistry.Contracts,
			databaseRegistry.PeginRepository,
			messaging.Rpc,
			liquidityProvider,
			messaging.EventBus,
			mutexes.RskWalletMutex(),
		),
		registerPeginUseCase: pegin.NewRegisterPeginUseCase(
			rskRegistry.Contracts,
			databaseRegistry.PeginRepository,
			messaging.EventBus,
			messaging.Rpc,
			mutexes.RskWalletMutex(),
		),
		acceptPeginQuoteUseCase: pegin.NewAcceptQuoteUseCase(
			databaseRegistry.PeginRepository,
			rskRegistry.Contracts,
			messaging.Rpc,
			liquidityProvider,
			liquidityProvider,
			messaging.EventBus,
			mutexes.PeginLiquidityMutex(),
		),
		getWatchedPeginQuoteUseCase: watcher.NewGetWatchedPeginQuoteUseCase(databaseRegistry.PeginRepository),
		expiredPeginQuoteUseCase:    pegin.NewExpiredPeginQuoteUseCase(databaseRegistry.PeginRepository),
		cleanExpiredQuotesUseCase: watcher.NewCleanExpiredQuotesUseCase(
			databaseRegistry.PeginRepository,
			databaseRegistry.PegoutRepository,
		),
		getProviderDetailUseCase: liquidity_provider.NewGetDetailUseCase(
			env.Captcha.SiteKey,
			env.Captcha.Disabled,
			liquidityProvider,
			liquidityProvider,
			liquidityProvider,
		),
		getWatchedPegoutQuoteUseCase: watcher.NewGetWatchedPegoutQuoteUseCase(
			databaseRegistry.PegoutRepository,
		),
		expiredPegoutUseCase:       pegout.NewExpiredPegoutQuoteUseCase(databaseRegistry.PegoutRepository),
		updatePegoutDepositUseCase: watcher.NewUpdatePegoutQuoteDepositUseCase(databaseRegistry.PegoutRepository),
		initPegoutDepositCacheUseCase: pegout.NewInitPegoutDepositCacheUseCase(
			databaseRegistry.PegoutRepository,
			rskRegistry.Contracts,
			messaging.Rpc,
		),
		refundPegoutUseCase: pegout.NewRefundPegoutUseCase(
			databaseRegistry.PegoutRepository,
			rskRegistry.Contracts,
			messaging.EventBus,
			messaging.Rpc,
			mutexes.RskWalletMutex(),
		),
		getPegoutQuoteUseCase: pegout.NewGetQuoteUseCase(
			messaging.Rpc,
			rskRegistry.Contracts,
			databaseRegistry.PegoutRepository,
			liquidityProvider,
			liquidityProvider,
			btcRegistry.PaymentWallet,
			env.Rsk.FeeCollectorAddress,
		),
		acceptPegoutQuoteUseCase: pegout.NewAcceptQuoteUseCase(
			databaseRegistry.PegoutRepository,
			rskRegistry.Contracts,
			liquidityProvider,
			liquidityProvider,
			messaging.EventBus,
			mutexes.PegoutLiquidityMutex(),
		),
		sendPegoutUseCase: pegout.NewSendPegoutUseCase(
			btcRegistry.PaymentWallet,
			databaseRegistry.PegoutRepository,
			messaging.Rpc,
			messaging.EventBus,
			rskRegistry.Contracts,
			mutexes.BtcWalletMutex(),
		),
		getUserDepositsUseCase: pegout.NewGetUserDepositsUseCase(databaseRegistry.PegoutRepository),
		liquidityCheckUseCase: liquidity_provider.NewCheckLiquidityUseCase(
			liquidityProvider,
			liquidityProvider,
			rskRegistry.Contracts,
			messaging.AlertSender,
			env.Provider.AlertRecipientEmail,
		),
		penalizationAlertUseCase: liquidity_provider.NewPenalizationAlertUseCase(
			rskRegistry.Contracts,
			messaging.AlertSender,
			env.Provider.AlertRecipientEmail,
		),
		addPeginCollateralUseCase:  pegin.NewAddCollateralUseCase(rskRegistry.Contracts, liquidityProvider),
		addPegoutCollateralUseCase: pegout.NewAddCollateralUseCase(rskRegistry.Contracts, liquidityProvider),
		changeStatusUseCase:        liquidity_provider.NewChangeStatusUseCase(rskRegistry.Contracts, liquidityProvider),
		resignUseCase:              liquidity_provider.NewResignUseCase(rskRegistry.Contracts, liquidityProvider),
		getProvidersUseCase:        liquidity_provider.NewGetProvidersUseCase(rskRegistry.Contracts),
		getPeginCollateralUseCase:  pegin.NewGetCollateralUseCase(rskRegistry.Contracts, liquidityProvider),
		getPegoutCollateralUseCase: pegout.NewGetCollateralUseCase(rskRegistry.Contracts, liquidityProvider),
		withdrawCollateralUseCase:  liquidity_provider.NewWithdrawCollateralUseCase(rskRegistry.Contracts),
		healthUseCase:              usecases.NewHealthUseCase(rskRegistry.Client, btcRegistry.RpcConnection, databaseRegistry.Connection),
		setGeneralConfigUseCase: liquidity_provider.NewSetGeneralConfigUseCase(
			databaseRegistry.LiquidityProviderRepository,
			rskRegistry.Wallet,
			signingHashFunction,
		),
		setPeginConfigUseCase: liquidity_provider.NewSetPeginConfigUseCase(
			databaseRegistry.LiquidityProviderRepository,
			rskRegistry.Wallet,
			signingHashFunction,
		),
		setPegoutConfigUseCase: liquidity_provider.NewSetPegoutConfigUseCase(
			databaseRegistry.LiquidityProviderRepository,
			rskRegistry.Wallet,
			signingHashFunction,
			rskRegistry.Contracts,
		),
		getConfigurationUseCase: liquidity_provider.NewGetConfigUseCase(liquidityProvider, liquidityProvider, liquidityProvider),
		loginUseCase:            liquidity_provider.NewLoginUseCase(databaseRegistry.LiquidityProviderRepository, messaging.EventBus),
		setCredentialsUseCase: liquidity_provider.NewSetCredentialsUseCase(
			databaseRegistry.LiquidityProviderRepository,
			rskRegistry.Wallet,
			signingHashFunction,
			messaging.EventBus,
		),
		defaultCredentialsUseCase: liquidity_provider.NewGenerateDefaultCredentialsUseCase(
			databaseRegistry.LiquidityProviderRepository,
			messaging.EventBus,
		),
		getManagementUiDataUseCase: liquidity_provider.NewGetManagementUiDataUseCase(
			databaseRegistry.LiquidityProviderRepository,
			liquidityProvider,
			liquidityProvider,
			liquidityProvider,
			rskRegistry.Contracts,
			env.Provider.ApiBaseUrl,
		),
		bridgePegoutUseCase: pegout.NewBridgePegoutUseCase(
			databaseRegistry.PegoutRepository,
			liquidityProvider,
			rskRegistry.Wallet,
			rskRegistry.Contracts,
			mutexes.RskWalletMutex(),
		),
		peginStatusUseCase:        pegin.NewStatusUseCase(databaseRegistry.PeginRepository),
		pegoutStatusUseCase:       pegout.NewStatusUseCase(databaseRegistry.PegoutRepository),
		availableLiquidityUseCase: liquidity_provider.NewGetAvailableLiquidityUseCase(liquidityProvider, liquidityProvider, liquidityProvider),
		updatePeginDepositUseCase: watcher.NewUpdatePeginDepositUseCase(databaseRegistry.PeginRepository),
		getServerInfoUseCase:      liquidity_provider.NewServerInfoUseCase(),
		getAssetReportUseCase: reports.NewGetAssetReportUseCase(
			btcRegistry.MonitoringWallet,
			messaging.Rpc,
			liquidityProvider,
			liquidityProvider,
			liquidityProvider,
			databaseRegistry.PeginRepository,
			databaseRegistry.PegoutRepository,
			rskRegistry.Contracts,
		),
	}
}

func (registry *UseCaseRegistry) GetPeginQuoteUseCase() *pegin.GetQuoteUseCase {
	return registry.getPeginQuoteUseCase
}

func (registry *UseCaseRegistry) GetRegistrationUseCase() *liquidity_provider.RegistrationUseCase {
	return registry.registerProviderUseCase
}

func (registry *UseCaseRegistry) GetAcceptPeginQuoteUseCase() *pegin.AcceptQuoteUseCase {
	return registry.acceptPeginQuoteUseCase
}

func (registry *UseCaseRegistry) GetProviderDetailUseCase() *liquidity_provider.GetDetailUseCase {
	return registry.getProviderDetailUseCase
}

func (registry *UseCaseRegistry) GetPegoutQuoteUseCase() *pegout.GetQuoteUseCase {
	return registry.getPegoutQuoteUseCase
}

func (registry *UseCaseRegistry) GetAcceptPegoutQuoteUseCase() *pegout.AcceptQuoteUseCase {
	return registry.acceptPegoutQuoteUseCase
}

func (registry *UseCaseRegistry) GetUserDepositsUseCase() *pegout.GetUserDepositsUseCase {
	return registry.getUserDepositsUseCase
}

func (registry *UseCaseRegistry) GetProvidersUseCase() *liquidity_provider.GetProvidersUseCase {
	return registry.getProvidersUseCase
}

func (registry *UseCaseRegistry) GetPeginCollateralUseCase() *pegin.GetCollateralUseCase {
	return registry.getPeginCollateralUseCase
}

func (registry *UseCaseRegistry) GetPegoutCollateralUseCase() *pegout.GetCollateralUseCase {
	return registry.getPegoutCollateralUseCase
}

func (registry *UseCaseRegistry) WithdrawCollateralUseCase() *liquidity_provider.WithdrawCollateralUseCase {
	return registry.withdrawCollateralUseCase
}

func (registry *UseCaseRegistry) HealthUseCase() *usecases.HealthUseCase {
	return registry.healthUseCase
}

func (registry *UseCaseRegistry) ResignationUseCase() *liquidity_provider.ResignUseCase {
	return registry.resignUseCase
}

func (registry *UseCaseRegistry) ChangeStatusUseCase() *liquidity_provider.ChangeStatusUseCase {
	return registry.changeStatusUseCase
}

func (registry *UseCaseRegistry) AddPeginCollateralUseCase() *pegin.AddCollateralUseCase {
	return registry.addPeginCollateralUseCase
}

func (registry *UseCaseRegistry) AddPegoutCollateralUseCase() *pegout.AddCollateralUseCase {
	return registry.addPegoutCollateralUseCase
}

func (registry *UseCaseRegistry) SetPeginConfigUseCase() *liquidity_provider.SetPeginConfigUseCase {
	return registry.setPeginConfigUseCase
}

func (registry *UseCaseRegistry) SetPegoutConfigUseCase() *liquidity_provider.SetPegoutConfigUseCase {
	return registry.setPegoutConfigUseCase
}

func (registry *UseCaseRegistry) SetGeneralConfigUseCase() *liquidity_provider.SetGeneralConfigUseCase {
	return registry.setGeneralConfigUseCase
}

func (registry *UseCaseRegistry) GetConfigurationUseCase() *liquidity_provider.GetConfigUseCase {
	return registry.getConfigurationUseCase
}

func (registry *UseCaseRegistry) LoginUseCase() *liquidity_provider.LoginUseCase {
	return registry.loginUseCase
}

func (registry *UseCaseRegistry) SetCredentialsUseCase() *liquidity_provider.SetCredentialsUseCase {
	return registry.setCredentialsUseCase
}

func (registry *UseCaseRegistry) GenerateDefaultCredentialsUseCase() *liquidity_provider.GenerateDefaultCredentialsUseCase {
	return registry.defaultCredentialsUseCase
}

func (registry *UseCaseRegistry) GetManagementUiDataUseCase() *liquidity_provider.GetManagementUiDataUseCase {
	return registry.getManagementUiDataUseCase
}

func (registry *UseCaseRegistry) GetPeginStatusUseCase() *pegin.StatusUseCase {
	return registry.peginStatusUseCase
}

func (registry *UseCaseRegistry) GetPegoutStatusUseCase() *pegout.StatusUseCase {
	return registry.pegoutStatusUseCase
}

func (registry *UseCaseRegistry) GetAvailableLiquidityUseCase() *liquidity_provider.GetAvailableLiquidityUseCase {
	return registry.availableLiquidityUseCase
}

func (registry *UseCaseRegistry) GetServerInfoUseCase() *liquidity_provider.ServerInfoUseCase {
	return registry.getServerInfoUseCase
}

func (registry *UseCaseRegistry) GetAssetReportUseCase() *reports.GetAssetReportUseCase {
	return registry.getAssetReportUseCase
}
