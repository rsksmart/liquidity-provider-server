package watcher_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"sync"
	"testing"
	"time"
)

func createWatcherShutdownTest(t *testing.T, createFunc func(t watcher.Ticker) watcher.Watcher) {
	tickerChannel := make(chan time.Time, 1)
	ticker := &mocks.TickerMock{}
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop().Return()
	w := createFunc(ticker)
	closeChannel := make(chan bool, 1)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		w.Start()
		wg.Done()
	}()
	go func() {
		w.Shutdown(closeChannel)
		<-closeChannel
		wg.Done()
	}()
	wg.Wait()
	ticker.AssertExpectations(t)
}
