package watcher_test

import (
	w "github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
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
		assert.Eventually(t, func() bool {
			return ticker.AssertExpectations(t) && useCase.AssertExpectations(t)
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
		assert.Eventually(t, func() bool {
			return ticker.AssertExpectations(t) && useCase.AssertExpectations(t)
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
		assert.Eventually(t, func() bool {
			return ticker.AssertExpectations(t) && useCase.AssertExpectations(t)
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
		assert.Eventually(t, func() bool {
			return ticker.AssertExpectations(t) && useCase.AssertExpectations(t)
		}, time.Second*1, time.Millisecond*100)
	})
}
