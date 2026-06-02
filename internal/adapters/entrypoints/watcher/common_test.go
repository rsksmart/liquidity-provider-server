package watcher_test

import (
	"sync"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// mockCollectT wraps assert.CollectT to satisfy mock.TestingT which also requires Logf.
type mockCollectT struct {
	*assert.CollectT
}

func (m mockCollectT) Logf(msg string, args ...interface{}) {
	log.Debugf(msg, args...)
}

func newMockCollectT(collect *assert.CollectT) mockCollectT {
	return mockCollectT{collect}
}

func createWatcherShutdownTest(t *testing.T, createFunc func(t utils.Ticker) watcher.Watcher) {
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
