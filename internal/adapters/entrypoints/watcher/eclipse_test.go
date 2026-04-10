package watcher_test

import (
	"testing"
	"time"

	w "github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// nolint:funlen
func TestEclipseWatcher_Start(t *testing.T) {
	const cooldownSeconds = 2
	t.Run("should start and run without errors", func(t *testing.T) {
		ticker := &mocks.TickerMock{}
		useCase := &mocks.EclipseCheckUseCaseMock{}
		tickerChannel := make(chan time.Time)
		closeChannel := make(chan bool)
		ticker.EXPECT().C().Return(tickerChannel)
		ticker.EXPECT().Stop().Return().Once()
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeBitcoin).Return(nil).Once()
		eclipseWatcher := w.NewEclipseWatcher(useCase, entities.NodeTypeBitcoin, cooldownSeconds, ticker)
		go eclipseWatcher.Start()
		tickerChannel <- time.Now()
		go eclipseWatcher.Shutdown(closeChannel)
		<-closeChannel
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			mt := mockCollectT{collect}
			ticker.AssertExpectations(mt)
			useCase.AssertExpectations(mt)
		}, time.Second*1, time.Millisecond*100)
	})
	t.Run("should not perform eclipse check during the cooldown", func(t *testing.T) {
		ticker := &mocks.TickerMock{}
		useCase := &mocks.EclipseCheckUseCaseMock{}
		tickerChannel := make(chan time.Time)
		closeChannel := make(chan bool)
		ticker.EXPECT().C().Return(tickerChannel)
		ticker.EXPECT().Stop().Return().Once()
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeBitcoin).Return(watcher.NodeEclipseDetectedError).Once()
		eclipseWatcher := w.NewEclipseWatcher(useCase, entities.NodeTypeBitcoin, cooldownSeconds, ticker)
		go eclipseWatcher.Start()
		tickerChannel <- time.Now()
		tickerChannel <- time.Now()
		go eclipseWatcher.Shutdown(closeChannel)
		<-closeChannel
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			mt := mockCollectT{collect}
			ticker.AssertExpectations(mt)
			useCase.AssertExpectations(mt)
		}, time.Second*1, time.Millisecond*100)
	})
	t.Run("should run again the eclipse check after an error", func(t *testing.T) {
		ticker := &mocks.TickerMock{}
		useCase := &mocks.EclipseCheckUseCaseMock{}
		tickerChannel := make(chan time.Time)
		closeChannel := make(chan bool)
		ticker.EXPECT().C().Return(tickerChannel)
		ticker.EXPECT().Stop().Return().Once()
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeBitcoin).Return(assert.AnError).Once()
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeBitcoin).Return(nil).Once()
		eclipseWatcher := w.NewEclipseWatcher(useCase, entities.NodeTypeBitcoin, cooldownSeconds, ticker)
		go eclipseWatcher.Start()
		tickerChannel <- time.Now()
		tickerChannel <- time.Now()
		go eclipseWatcher.Shutdown(closeChannel)
		<-closeChannel
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			mt := mockCollectT{collect}
			ticker.AssertExpectations(mt)
			useCase.AssertExpectations(mt)
		}, time.Second*1, time.Millisecond*100)
	})
	t.Run("should run again the eclipse check after the cooldown", func(t *testing.T) {
		ticker := &mocks.TickerMock{}
		useCase := &mocks.EclipseCheckUseCaseMock{}
		tickerChannel := make(chan time.Time)
		closeChannel := make(chan bool)
		ticker.EXPECT().C().Return(tickerChannel)
		ticker.EXPECT().Stop().Return().Once()
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeBitcoin).Return(watcher.NodeEclipseDetectedError).Once()
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeBitcoin).Return(nil).Once()
		eclipseWatcher := w.NewEclipseWatcher(useCase, entities.NodeTypeBitcoin, cooldownSeconds, ticker)
		go eclipseWatcher.Start()
		tickerChannel <- time.Now()
		time.Sleep(time.Second * time.Duration(cooldownSeconds+1))
		tickerChannel <- time.Now()
		go eclipseWatcher.Shutdown(closeChannel)
		<-closeChannel
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			mt := mockCollectT{collect}
			ticker.AssertExpectations(mt)
			useCase.AssertExpectations(mt)
		}, time.Second*1, time.Millisecond*100)
	})
}

func TestEclipseWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) w.Watcher {
		useCase := &mocks.EclipseCheckUseCaseMock{}
		return w.NewEclipseWatcher(useCase, entities.NodeTypeBitcoin, 1, ticker)
	})
}
