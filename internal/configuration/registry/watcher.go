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
	PegoutBridgeWatcher        *watcher.PegoutBridgeWatcher
}

// nolint:funlen
func NewWatcherRegistry(
	env environment.Environment,
	useCaseRegistry *UseCaseRegistry,
	rskRegistry *Rootstock,
	btcRegistry *Bitcoin,
	liquidityProvider *dataproviders.LocalLiquidityProvider,
	messaging *Messaging,
	tickers *watcher.ApplicationTickers,
) *WatcherRegistry {
	return &WatcherRegistry{
		PeginDepositAddressWatcher: watcher.NewPeginDepositAddressWatcher(
			watcher.NewPeginDepositAddressWatcherUseCases(
				useCaseRegistry.callForUserUseCase,
				useCaseRegistry.getWatchedPeginQuoteUseCase,
				useCaseRegistry.updatePeginDepositUseCase,
				useCaseRegistry.expiredPeginQuoteUseCase,
			),
			btcRegistry.MonitoringWallet,
			messaging.Rpc,
			messaging.EventBus,
			tickers.PeginDepositWatcherTicker,
		),
		PeginBridgeWatcher: watcher.NewPeginBridgeWatcher(
			useCaseRegistry.registerPeginUseCase,
			useCaseRegistry.getWatchedPeginQuoteUseCase,
			rskRegistry.Contracts,
			messaging.Rpc,
			messaging.EventBus,
			tickers.PeginBridgeWatcherTicker,
		),
		QuoteCleanerWatcher: watcher.NewQuoteCleanerWatcher(
			useCaseRegistry.cleanExpiredQuotesUseCase,
			tickers.QuoteCleanerTicker,
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
			tickers.PegoutDepositWatcherTicker,
		),
		PegoutBtcTransferWatcher: watcher.NewPegoutBtcTransferWatcher(
			useCaseRegistry.getWatchedPegoutQuoteUseCase,
			useCaseRegistry.refundPegoutUseCase,
			messaging.Rpc,
			messaging.EventBus,
			tickers.PegoutBtcTransferWatcherTicker,
		),
		LiquidityCheckWatcher: watcher.NewLiquidityCheckWatcher(useCaseRegistry.liquidityCheckUseCase, tickers.LiquidityCheckTicker),
		PenalizationAlertWatcher: watcher.NewPenalizationAlertWatcher(
			messaging.Rpc,
			useCaseRegistry.penalizationAlertUseCase,
			tickers.PenalizationCheckTicker,
		),
		PegoutBridgeWatcher: watcher.NewPegoutBridgeWatcher(
			useCaseRegistry.getWatchedPegoutQuoteUseCase,
			useCaseRegistry.bridgePegoutUseCase,
			tickers.PegoutBridgeWatcherTicker,
		),
	}
}
