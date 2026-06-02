package watcher_test

import (
	"testing"
	"time"

	w "github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRootstockPeerWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) w.Watcher {
		useCase := &mocks.NodePeerCheckUseCaseMock{}
		return w.NewRootstockPeerWatcher(useCase, ticker, 15*time.Second)
	})
}

func TestRootstockPeerWatcher_Start(t *testing.T) {
	t.Run("should call use case on tick with rootstock node type", func(t *testing.T) {
		useCase := &mocks.NodePeerCheckUseCaseMock{}
		ticker := &mocks.TickerMock{}
		tickerChannel := make(chan time.Time)
		closeChannel := make(chan bool)
		ticker.EXPECT().C().Return(tickerChannel)
		ticker.EXPECT().Stop().Return().Once()
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeRootstock).Return(nil).Once()
		watcher := w.NewRootstockPeerWatcher(useCase, ticker, 15*time.Second)
		go watcher.Start()
		tickerChannel <- time.Now()
		go watcher.Shutdown(closeChannel)
		<-closeChannel
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			mt := newMockCollectT(collect)
			useCase.AssertExpectations(mt)
			ticker.AssertExpectations(mt)
		}, time.Second, 10*time.Millisecond)
	})

	t.Run("should continue running on use case error", func(t *testing.T) {
		useCase := &mocks.NodePeerCheckUseCaseMock{}
		ticker := &mocks.TickerMock{}
		tickerChannel := make(chan time.Time)
		closeChannel := make(chan bool)
		ticker.EXPECT().C().Return(tickerChannel)
		ticker.EXPECT().Stop().Return().Once()
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeRootstock).Return(assert.AnError).Once()
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeRootstock).Return(nil).Once()
		watcher := w.NewRootstockPeerWatcher(useCase, ticker, 15*time.Second)
		go watcher.Start()
		tickerChannel <- time.Now()
		tickerChannel <- time.Now()
		go watcher.Shutdown(closeChannel)
		<-closeChannel
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			mt := newMockCollectT(collect)
			ticker.AssertExpectations(mt)
			useCase.AssertExpectations(mt)
		}, time.Second, 10*time.Millisecond)
	})
}
