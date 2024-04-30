package registry

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
)

type WatcherRegistry struct {
	PeginDepositAddressWatcher *watcher.PeginDepositAddressWatcher
	PeginBridgeWatcher         *watcher.PeginBridgeWatcher
	QuoteCleanerWatcher        *watcher.QuoteCleanerWatcher
	PegoutRskDepositWatcher    *watcher.PegoutRskDepositWatcher
	PegoutBtcTransferWatcher   *watcher.PegoutBtcTransferWatcher
	LiquidityCheckWatcher      *watcher.LiquidityCheckWatcher
	PenalizationAlertWatcher   *watcher.PenalizationAlertWatcher
}

func NewWatcherRegistry(
	env environment.Environment,
	useCaseRegistry *UseCaseRegistry,
	rskRegistry *Rootstock,
	btcRegistry *Bitcoin,
	liquidityProvider *dataproviders.LocalLiquidityProvider,
	messaging *Messaging,
) *WatcherRegistry {
	return &WatcherRegistry{
		PeginDepositAddressWatcher: watcher.NewPeginDepositAddressWatcher(
			useCaseRegistry.callForUserUseCase,
			useCaseRegistry.getWatchedPeginQuoteUseCase,
			useCaseRegistry.expiredPeginQuoteUseCase,
			btcRegistry.MonitoringWallet,
			messaging.Rpc,
			messaging.EventBus,
		),
		PeginBridgeWatcher: watcher.NewPeginBridgeWatcher(
			useCaseRegistry.registerPeginUseCase,
			useCaseRegistry.getWatchedPeginQuoteUseCase,
			rskRegistry.Contracts,
			messaging.Rpc,
			messaging.EventBus,
		),
		QuoteCleanerWatcher: watcher.NewQuoteCleanerWatcher(
			useCaseRegistry.cleanExpiredQuotesUseCase,
		),
		PegoutRskDepositWatcher: watcher.NewPegoutRskDepositWatcher(
			watcher.NewPegoutRskDepositWatcherUseCases(
				useCaseRegistry.getWatchedPegoutQuoteUseCase,
				useCaseRegistry.expiredPegoutUseCase,
				useCaseRegistry.sendPegoutUseCase,
				useCaseRegistry.updatePegoutDepositUseCase,
				useCaseRegistry.initPegoutDepositCacheUseCase,
			),
			liquidityProvider,
			messaging.Rpc,
			rskRegistry.Contracts,
			messaging.EventBus,
			env.Pegout.DepositCacheStartBlock,
		),
		PegoutBtcTransferWatcher: watcher.NewPegoutBtcTransferWatcher(
			useCaseRegistry.getWatchedPegoutQuoteUseCase,
			useCaseRegistry.refundPegoutUseCase,
			messaging.Rpc,
			messaging.EventBus,
		),
		LiquidityCheckWatcher: watcher.NewLiquidityCheckWatcher(useCaseRegistry.liquidityCheckUseCase),
		PenalizationAlertWatcher: watcher.NewPenalizationAlertWatcher(
			messaging.Rpc,
			useCaseRegistry.penalizationAlertUseCase,
		),
	}
}
