package watcher

import "time"

// Ticker is an interface to be able to mock time.Ticker in the unit test
type Ticker interface {
	Stop()
	C() <-chan time.Time
}

type TickerWrapper struct {
	ticker *time.Ticker
}

func NewTickerWrapper(d time.Duration) *TickerWrapper {
	return &TickerWrapper{
		ticker: time.NewTicker(d),
	}
}

func (t *TickerWrapper) Stop() {
	t.ticker.Stop()
}

func (t *TickerWrapper) C() <-chan time.Time {
	return t.ticker.C
}

type ApplicationTickers struct {
	LiquidityCheckTicker           Ticker
	PeginBridgeWatcherTicker       Ticker
	QuoteCleanerTicker             Ticker
	PeginDepositWatcherTicker      Ticker
	PenalizationCheckTicker        Ticker
	PegoutDepositWatcherTicker     Ticker
	PegoutBtcTransferWatcherTicker Ticker
	PegoutBridgeWatcherTicker      Ticker
}

func NewApplicationTickers() *ApplicationTickers {
	return &ApplicationTickers{
		LiquidityCheckTicker:           NewTickerWrapper(liquidityCheckInterval),
		PeginBridgeWatcherTicker:       NewTickerWrapper(peginBridgeWatcherInterval),
		QuoteCleanerTicker:             NewTickerWrapper(quoteCleanInterval),
		PeginDepositWatcherTicker:      NewTickerWrapper(peginDepositWatcherInterval),
		PenalizationCheckTicker:        NewTickerWrapper(penalizationCheckInterval),
		PegoutDepositWatcherTicker:     NewTickerWrapper(pegoutDepositWatcherInterval),
		PegoutBtcTransferWatcherTicker: NewTickerWrapper(pegoutBtcTransferWatcherInterval),
		PegoutBridgeWatcherTicker:      NewTickerWrapper(pegoutBridgeWatcherInterval),
	}
}
