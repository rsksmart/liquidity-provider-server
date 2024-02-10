package registry

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
)

type UseCaseRegistry struct {
	getPeginQuoteUseCase            *pegin.GetQuoteUseCase
	registerProviderUseCase         *liquidity_provider.RegistrationUseCase
	callForUserUseCase              *pegin.CallForUserUseCase
	registerPeginUseCase            *pegin.RegisterPeginUseCase
	acceptPeginQuoteUseCase         *pegin.AcceptQuoteUseCase
	getWatchedPeginQuoteUseCase     *watcher.GetWatchedPeginQuoteUseCase
	expiredPeginQuoteUseCase        *pegin.ExpiredPeginQuoteUseCase
	cleanExpiredQuotesUseCase       *watcher.CleanExpiredQuotesUseCase
	getProviderDetailUseCase        *liquidity_provider.GetDetailUseCase
	getWatchedPegoutQuoteUseCase    *watcher.GetWatchedPegoutQuoteUseCase
	expiredPegoutUseCase            *pegout.ExpiredPegoutQuoteUseCase
	sendPegoutUseCase               *pegout.SendPegoutUseCase
	updatePegoutDepositUseCase      *watcher.UpdatePegoutQuoteDepositUseCase
	initPegoutDepositCacheUseCase   *pegout.InitPegoutDepositCacheUseCase
	refundPegoutUseCase             *pegout.RefundPegoutUseCase
	getPegoutQuoteUseCase           *pegout.GetQuoteUseCase
	acceptPegoutQuoteUseCase        *pegout.AcceptQuoteUseCase
	getUserDepositsUseCase          *pegout.GetUserDepositsUseCase
	liquidityCheckUseCase           *liquidity_provider.CheckLiquidityUseCase
	penalizationAlertUseCase        *liquidity_provider.PenalizationAlertUseCase
	getProvidersUseCase             *liquidity_provider.GetProvidersUseCase
	getPeginCollateralUseCase       *pegin.GetCollateralUseCase
	getPegoutCollateralUseCase      *pegout.GetCollateralUseCase
	withdrawPeginCollateralUseCase  *pegin.WithdrawCollateralUseCase
	withdrawPegoutCollateralUseCase *pegout.WithdrawCollateralUseCase
	healthUseCase                   *usecases.HealthUseCase
	resignUseCase                   *liquidity_provider.ResignUseCase
	changeStatusUseCase             *liquidity_provider.ChangeStatusUseCase
	addPeginCollateralUseCase       *pegin.AddCollateralUseCase
	addPegoutCollateralUseCase      *pegout.AddCollateralUseCase
}

// NewUseCaseRegistry
// nolint:funlen
func NewUseCaseRegistry(
	env environment.Environment,
	rskRegistry *Rootstock,
	btcRegistry *Bitcoin,
	databaseRegistry *Database,
	liquidityProvider *dataproviders.LocalLiquidityProvider,
	eventBus entities.EventBus,
	alertSender entities.AlertSender,
	mutexes entities.ApplicationMutexes,
) *UseCaseRegistry {
	return &UseCaseRegistry{
		getPeginQuoteUseCase: pegin.NewGetQuoteUseCase(
			rskRegistry.RpcServer,
			rskRegistry.FeeCollector,
			rskRegistry.Bridge,
			rskRegistry.Lbc,
			databaseRegistry.PeginRepository,
			liquidityProvider,
			liquidityProvider,
			env.Rsk.FeeCollectorAddress,
		),
		registerProviderUseCase: liquidity_provider.NewRegistrationUseCase(
			rskRegistry.Lbc,
			liquidityProvider,
		),
		callForUserUseCase: pegin.NewCallForUserUseCase(
			rskRegistry.Lbc,
			databaseRegistry.PeginRepository,
			btcRegistry.RpcServer,
			liquidityProvider,
			eventBus,
			mutexes.RskWalletMutex(),
		),
		registerPeginUseCase: pegin.NewRegisterPeginUseCase(
			rskRegistry.Lbc,
			databaseRegistry.PeginRepository,
			eventBus,
			rskRegistry.Bridge,
			btcRegistry.RpcServer,
			mutexes.RskWalletMutex(),
		),
		acceptPeginQuoteUseCase: pegin.NewAcceptQuoteUseCase(
			databaseRegistry.PeginRepository,
			rskRegistry.Bridge,
			btcRegistry.RpcServer,
			rskRegistry.RpcServer,
			liquidityProvider,
			liquidityProvider,
			eventBus,
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
			rskRegistry.Lbc,
			rskRegistry.RpcServer,
		),
		refundPegoutUseCase: pegout.NewRefundPegoutUseCase(
			databaseRegistry.PegoutRepository,
			rskRegistry.Lbc,
			eventBus,
			btcRegistry.RpcServer,
			rskRegistry.Wallet,
			rskRegistry.Bridge,
			mutexes.RskWalletMutex(),
		),
		getPegoutQuoteUseCase: pegout.NewGetQuoteUseCase(
			rskRegistry.RpcServer,
			rskRegistry.FeeCollector,
			rskRegistry.Bridge,
			rskRegistry.Lbc,
			databaseRegistry.PegoutRepository,
			liquidityProvider,
			liquidityProvider,
			btcRegistry.Wallet,
			env.Rsk.FeeCollectorAddress,
		),
		acceptPegoutQuoteUseCase: pegout.NewAcceptQuoteUseCase(
			databaseRegistry.PegoutRepository,
			rskRegistry.Lbc,
			liquidityProvider,
			liquidityProvider,
			eventBus,
			mutexes.PegoutLiquidityMutex(),
		),
		sendPegoutUseCase: pegout.NewSendPegoutUseCase(
			btcRegistry.Wallet,
			databaseRegistry.PegoutRepository,
			rskRegistry.RpcServer,
			eventBus,
			mutexes.BtcWalletMutex(),
		),
		getUserDepositsUseCase: pegout.NewGetUserDepositsUseCase(databaseRegistry.PegoutRepository),
		liquidityCheckUseCase: liquidity_provider.NewCheckLiquidityUseCase(
			liquidityProvider,
			liquidityProvider,
			rskRegistry.Bridge,
			alertSender,
			env.Provider.AlertRecipientEmail,
		),
		penalizationAlertUseCase: liquidity_provider.NewPenalizationAlertUseCase(
			rskRegistry.Lbc,
			alertSender,
			env.Provider.AlertRecipientEmail,
		),
		addPegoutCollateralUseCase:      pegout.NewAddCollateralUseCase(rskRegistry.Lbc, liquidityProvider),
		changeStatusUseCase:             liquidity_provider.NewChangeStatusUseCase(rskRegistry.Lbc, liquidityProvider),
		resignUseCase:                   liquidity_provider.NewResignUseCase(rskRegistry.Lbc, liquidityProvider),
		getProvidersUseCase:             liquidity_provider.NewGetProvidersUseCase(rskRegistry.Lbc),
		getPeginCollateralUseCase:       pegin.NewGetCollateralUseCase(rskRegistry.Lbc, liquidityProvider),
		getPegoutCollateralUseCase:      pegout.NewGetCollateralUseCase(rskRegistry.Lbc, liquidityProvider),
		withdrawPeginCollateralUseCase:  pegin.NewWithdrawCollateralUseCase(rskRegistry.Lbc),
		withdrawPegoutCollateralUseCase: pegout.NewWithdrawCollateralUseCase(rskRegistry.Lbc),
		healthUseCase:                   usecases.NewHealthUseCase(rskRegistry.Client, btcRegistry.Connection, databaseRegistry.Connection),
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

func (registry *UseCaseRegistry) WithdrawPeginCollateralUseCase() *pegin.WithdrawCollateralUseCase {
	return registry.withdrawPeginCollateralUseCase
}

func (registry *UseCaseRegistry) WithdrawPegoutCollateralUseCase() *pegout.WithdrawCollateralUseCase {
	return registry.withdrawPegoutCollateralUseCase
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
