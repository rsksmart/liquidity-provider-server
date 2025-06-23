package utils

import "time"

// Ticker is an interface to be able to mock time.Ticker in the unit test
type Ticker interface {
	Stop()
	C() <-chan time.Time
}

type TickerWrapper struct {
	ticker *time.Ticker
}

func NewTickerWrapper(d time.Duration) Ticker {
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
