package registry

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
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
	eventBus entities.EventBus,
) *WatcherRegistry {
	return &WatcherRegistry{
		PeginDepositAddressWatcher: watcher.NewPeginDepositAddressWatcher(
			useCaseRegistry.callForUserUseCase,
			useCaseRegistry.getWatchedPeginQuoteUseCase,
			useCaseRegistry.expiredPeginQuoteUseCase,
			btcRegistry.Wallet,
			btcRegistry.RpcServer,
			eventBus,
		),
		PeginBridgeWatcher: watcher.NewPeginBridgeWatcher(
			useCaseRegistry.registerPeginUseCase,
			useCaseRegistry.getWatchedPeginQuoteUseCase,
			rskRegistry.Bridge,
			btcRegistry.RpcServer,
			eventBus,
		),
		QuoteCleanerWatcher: watcher.NewQuoteCleanerWatcher(
			useCaseRegistry.cleanExpiredQuotesUseCase,
		),
		PegoutRskDepositWatcher: watcher.NewPegoutRskDepositWatcher(
			useCaseRegistry.getWatchedPegoutQuoteUseCase,
			useCaseRegistry.expiredPegoutUseCase,
			useCaseRegistry.sendPegoutUseCase,
			useCaseRegistry.updatePegoutDepositUseCase,
			useCaseRegistry.initPegoutDepositCacheUseCase,
			liquidityProvider,
			rskRegistry.RpcServer,
			rskRegistry.Lbc,
			eventBus,
			env.Pegout.DepositCacheStartBlock,
		),
		PegoutBtcTransferWatcher: watcher.NewPegoutBtcTransferWatcher(
			useCaseRegistry.getWatchedPegoutQuoteUseCase,
			useCaseRegistry.refundPegoutUseCase,
			btcRegistry.RpcServer,
			eventBus,
		),
		LiquidityCheckWatcher: watcher.NewLiquidityCheckWatcher(useCaseRegistry.liquidityCheckUseCase),
		PenalizationAlertWatcher: watcher.NewPenalizationAlertWatcher(
			rskRegistry.RpcServer,
			useCaseRegistry.penalizationAlertUseCase,
		),
	}
}
