package watcher_test

import (
	"context"
	"testing"
	"time"

	w "github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRootstockReorgWatcher_Prepare(t *testing.T) {
	watcher := w.NewRootstockReorgWatcher(nil, nil, time.Second)
	assert.NoError(t, watcher.Prepare(context.Background()))
}

func TestRootstockReorgWatcher_Start(t *testing.T) {
	t.Run("should call reorg check on tick", func(t *testing.T) {
		ticker := &mocks.TickerMock{}
		useCase := mocks.NewNodeReorgCheckUseCaseMock(t)
		tickerChannel := make(chan time.Time)
		closeChannel := make(chan bool)
		ticker.EXPECT().C().Return(tickerChannel)
		ticker.EXPECT().Stop().Return().Once()
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeRootstock).Return(nil).Once()
		rw := w.NewRootstockReorgWatcher(useCase, ticker, time.Second)
		go rw.Start()
		tickerChannel <- time.Now()
		go rw.Shutdown(closeChannel)
		<-closeChannel
		assert.Eventually(t, func() bool {
			return ticker.AssertExpectations(t) && useCase.AssertExpectations(t)
		}, time.Second, time.Millisecond*50)
	})
	t.Run("should continue after reorg check error", func(t *testing.T) {
		ticker := &mocks.TickerMock{}
		useCase := mocks.NewNodeReorgCheckUseCaseMock(t)
		tickerChannel := make(chan time.Time)
		closeChannel := make(chan bool)
		ticker.EXPECT().C().Return(tickerChannel)
		ticker.EXPECT().Stop().Return().Once()
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeRootstock).Return(assert.AnError).Once()
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeRootstock).Return(nil).Once()
		rw := w.NewRootstockReorgWatcher(useCase, ticker, time.Second)
		go rw.Start()
		tickerChannel <- time.Now()
		tickerChannel <- time.Now()
		go rw.Shutdown(closeChannel)
		<-closeChannel
		assert.Eventually(t, func() bool {
			return ticker.AssertExpectations(t) && useCase.AssertExpectations(t)
		}, time.Second, time.Millisecond*50)
	})
}

func TestRootstockReorgWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) w.Watcher {
		useCase := mocks.NewNodeReorgCheckUseCaseMock(t)
		return w.NewRootstockReorgWatcher(useCase, ticker, time.Second)
	})
}
