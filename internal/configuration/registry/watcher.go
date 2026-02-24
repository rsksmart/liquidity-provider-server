package registry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher/monitoring"
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
	PegoutBridgeWatcher        *watcher.PegoutBridgeWatcher
	BitcoinEclipseWatcher      *watcher.EclipseWatcher
	RskEclipseWatcher          *watcher.EclipseWatcher
	BtcReleaseWatcher          *watcher.BtcReleaseWatcher
	QuoteMetricsWatcher        *monitoring.QuoteMetricsWatcher
	AssetReportWatcher         *monitoring.AssetReportWatcher
	TransferColdWalletWatcher  *watcher.TransferColdWalletWatcher
	ColdWalletMetricsWatcher   *monitoring.ColdWalletMetricsWatcher
}

// nolint:funlen
func NewWatcherRegistry(
	env environment.Environment,
	useCaseRegistry *UseCaseRegistry,
	rskRegistry *Rootstock,
	btcRegistry *Bitcoin,
	lpRegistry *LiquidityProvider,
	messaging *Messaging,
	tickers *watcher.ApplicationTickers,
	timeouts environment.ApplicationTimeouts,
) *WatcherRegistry {
	appMetrics := monitoring.NewMetrics(prometheus.DefaultRegisterer)

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
			lpRegistry.LiquidityProvider,
			messaging.Rpc,
			rskRegistry.Contracts,
			messaging.EventBus,
			env.Pegout.DepositCacheStartBlock,
			tickers.PegoutDepositWatcherTicker,
			timeouts.PegoutDepositCheck.Seconds(),
		),
		PegoutBtcTransferWatcher: watcher.NewPegoutBtcTransferWatcher(
			useCaseRegistry.getWatchedPegoutQuoteUseCase,
			useCaseRegistry.refundPegoutUseCase,
			messaging.Rpc,
			messaging.EventBus,
			tickers.PegoutBtcTransferWatcherTicker,
		),
		LiquidityCheckWatcher: watcher.NewLiquidityCheckWatcher(
			useCaseRegistry.liquidityCheckUseCase,
			useCaseRegistry.lowLiquidityAlertUseCase,
			tickers.LiquidityCheckTicker,
			timeouts.WatcherValidation.Seconds(),
		),
		PenalizationAlertWatcher: watcher.NewPenalizationAlertWatcher(
			messaging.Rpc,
			useCaseRegistry.penalizationAlertUseCase,
			tickers.PenalizationCheckTicker,
			timeouts.WatcherValidation.Seconds(),
		),
		PegoutBridgeWatcher: watcher.NewPegoutBridgeWatcher(
			useCaseRegistry.getWatchedPegoutQuoteUseCase,
			useCaseRegistry.bridgePegoutUseCase,
			tickers.PegoutBridgeWatcherTicker,
		),
		BitcoinEclipseWatcher: watcher.NewEclipseWatcher(
			useCaseRegistry.btcEclipseCheckUseCase,
			entities.NodeTypeBitcoin,
			env.Eclipse.FillWithDefaults().AlertCooldownSeconds,
			tickers.BitcoinEclipseCheckTicker,
		),
		RskEclipseWatcher: watcher.NewEclipseWatcher(
			useCaseRegistry.rskEclipseCheckUseCase,
			entities.NodeTypeRootstock,
			env.Eclipse.FillWithDefaults().AlertCooldownSeconds,
			tickers.RskEclipseCheckTicker,
		),
		BtcReleaseWatcher: watcher.NewBtcReleaseWatcher(
			rskRegistry.Contracts,
			messaging.Rpc,
			useCaseRegistry.updateBtcReleaseUseCase,
			tickers.BtcReleaseCheckTicker,
			env.Pegout.BtcReleaseWatcherStartBlock,
			env.Pegout.BtcReleaseWatcherPageSize,
			timeouts.BtcReleaseCheck.Seconds(),
		),
		QuoteMetricsWatcher: monitoring.NewQuoteMetricsWatcher(
			appMetrics,
			messaging.EventBus,
			useCaseRegistry.GetServerInfoUseCase(),
		),
		AssetReportWatcher: monitoring.NewAssetReportWatcher(
			appMetrics,
			useCaseRegistry.GetAssetsReportUseCase(),
			tickers.AssetReportTicker,
		),
		TransferColdWalletWatcher: watcher.NewTransferColdWalletWatcher(
			useCaseRegistry.TransferExcessToColdWalletUseCase(),
			tickers.TransferColdWalletTicker,
			timeouts.WatcherValidation.Seconds(),
		),
		ColdWalletMetricsWatcher: monitoring.NewColdWalletMetricsWatcher(
			appMetrics,
			messaging.EventBus,
		),
	}
}
