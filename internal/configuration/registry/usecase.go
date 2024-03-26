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
	setPeginConfigUseCase           *liquidity_provider.SetPeginConfigUseCase
	setPegoutConfigUseCase          *liquidity_provider.SetPegoutConfigUseCase
	setGeneralConfigUseCase         *liquidity_provider.SetGeneralConfigUseCase
	getConfigurationUseCase         *liquidity_provider.GetConfigUseCase
	liquidityStatusUseCase			*liquidity_provider.LiquidityStatusUseCase
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
			rskRegistry.Wallet,
			mutexes.RskWalletMutex(),
		),
		getPegoutQuoteUseCase: pegout.NewGetQuoteUseCase(
			messaging.Rpc,
			rskRegistry.Contracts,
			databaseRegistry.PegoutRepository,
			liquidityProvider,
			liquidityProvider,
			btcRegistry.Wallet,
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
			btcRegistry.Wallet,
			databaseRegistry.PegoutRepository,
			messaging.Rpc,
			messaging.EventBus,
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
		addPeginCollateralUseCase:       pegin.NewAddCollateralUseCase(rskRegistry.Contracts, liquidityProvider),
		addPegoutCollateralUseCase:      pegout.NewAddCollateralUseCase(rskRegistry.Contracts, liquidityProvider),
		changeStatusUseCase:             liquidity_provider.NewChangeStatusUseCase(rskRegistry.Contracts, liquidityProvider),
		resignUseCase:                   liquidity_provider.NewResignUseCase(rskRegistry.Contracts, liquidityProvider),
		getProvidersUseCase:             liquidity_provider.NewGetProvidersUseCase(rskRegistry.Contracts),
		getPeginCollateralUseCase:       pegin.NewGetCollateralUseCase(rskRegistry.Contracts, liquidityProvider),
		getPegoutCollateralUseCase:      pegout.NewGetCollateralUseCase(rskRegistry.Contracts, liquidityProvider),
		withdrawPeginCollateralUseCase:  pegin.NewWithdrawCollateralUseCase(rskRegistry.Contracts),
		withdrawPegoutCollateralUseCase: pegout.NewWithdrawCollateralUseCase(rskRegistry.Contracts),
		healthUseCase:                   usecases.NewHealthUseCase(rskRegistry.Client, btcRegistry.Connection, databaseRegistry.Connection),
		liquidityStatusUseCase:		 	 liquidity_provider.NewLiquidityStatusUseCase(rskRegistry.Contracts, liquidityProvider, messaging.Rpc,btcRegistry.Wallet, liquidityProvider),
		setGeneralConfigUseCase: liquidity_provider.NewSetGeneralConfigUseCase(
			databaseRegistry.LiquidityProviderRepository,
			rskRegistry.Wallet,
			crypto.Keccak256,
		),
		setPeginConfigUseCase: liquidity_provider.NewSetPeginConfigUseCase(
			databaseRegistry.LiquidityProviderRepository,
			rskRegistry.Wallet,
			crypto.Keccak256,
		),
		setPegoutConfigUseCase: liquidity_provider.NewSetPegoutConfigUseCase(
			databaseRegistry.LiquidityProviderRepository,
			rskRegistry.Wallet,
			crypto.Keccak256,
		),
		getConfigurationUseCase: liquidity_provider.NewGetConfigUseCase(liquidityProvider, liquidityProvider, liquidityProvider),
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
func (registry *UseCaseRegistry) GetLiquidityStatusUseCase() *liquidity_provider.LiquidityStatusUseCase {
    return registry.liquidityStatusUseCase
}