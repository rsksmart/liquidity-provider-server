package watcher

import "github.com/rsksmart/liquidity-provider-server/internal/entities/utils"

type ApplicationTickers struct {
	LiquidityCheckTicker           utils.Ticker
	PeginBridgeWatcherTicker       utils.Ticker
	QuoteCleanerTicker             utils.Ticker
	PeginDepositWatcherTicker      utils.Ticker
	PenalizationCheckTicker        utils.Ticker
	PegoutDepositWatcherTicker     utils.Ticker
	PegoutBtcTransferWatcherTicker utils.Ticker
	PegoutBridgeWatcherTicker      utils.Ticker
}

func NewApplicationTickers() *ApplicationTickers {
	return &ApplicationTickers{
		LiquidityCheckTicker:           utils.NewTickerWrapper(liquidityCheckInterval),
		PeginBridgeWatcherTicker:       utils.NewTickerWrapper(peginBridgeWatcherInterval),
		QuoteCleanerTicker:             utils.NewTickerWrapper(quoteCleanInterval),
		PeginDepositWatcherTicker:      utils.NewTickerWrapper(peginDepositWatcherInterval),
		PenalizationCheckTicker:        utils.NewTickerWrapper(penalizationCheckInterval),
		PegoutDepositWatcherTicker:     utils.NewTickerWrapper(pegoutDepositWatcherInterval),
		PegoutBtcTransferWatcherTicker: utils.NewTickerWrapper(pegoutBtcTransferWatcherInterval),
		PegoutBridgeWatcherTicker:      utils.NewTickerWrapper(pegoutBridgeWatcherInterval),
	}
}
