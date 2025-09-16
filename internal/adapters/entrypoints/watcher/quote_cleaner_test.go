package watcher_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQuoteCleanerWatcher_Start(t *testing.T) {
	t.Run("handle use case error", func(t *testing.T) {
		checkFunction := test.AssertLogContains(t, "Error cleaning quotes")
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		peginRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, mock.Anything).Return(nil, assert.AnError)
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		ticker := &mocks.TickerMock{}
		tickerChannel := make(chan time.Time)
		ticker.EXPECT().C().Return(tickerChannel)
		ticker.EXPECT().Stop()
		useCase := w.NewCleanExpiredQuotesUseCase(peginRepository, pegoutRepository)
		quoteCleaner := watcher.NewQuoteCleanerWatcher(useCase, ticker)
		go quoteCleaner.Start()
		tickerChannel <- time.Now()
		go quoteCleaner.Shutdown(make(chan bool))
		assert.Eventually(t, func() bool {
			return checkFunction() && ticker.AssertExpectations(t) && peginRepository.AssertExpectations(t)
		}, time.Second, 10*time.Millisecond)
	})
	t.Run("clean quotes successfully", func(t *testing.T) {
		checkFunction := test.AssertLogContains(t, "Cleaned 3 quotes")
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		peginRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, mock.Anything).
			Return([]quote.RetainedPeginQuote{{QuoteHash: "pegin1"}, {QuoteHash: "pegin2"}}, nil)
		peginRepository.EXPECT().DeleteQuotes(mock.Anything, []string{"pegin1", "pegin2"}).Return(2, nil)
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		pegoutRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, mock.Anything).
			Return([]quote.RetainedPegoutQuote{{QuoteHash: "pegout1"}}, nil)
		pegoutRepository.EXPECT().DeleteQuotes(mock.Anything, []string{"pegout1"}).Return(1, nil)
		ticker := &mocks.TickerMock{}
		tickerChannel := make(chan time.Time)
		ticker.EXPECT().C().Return(tickerChannel)
		ticker.EXPECT().Stop()
		useCase := w.NewCleanExpiredQuotesUseCase(peginRepository, pegoutRepository)
		quoteCleaner := watcher.NewQuoteCleanerWatcher(useCase, ticker)
		go quoteCleaner.Start()
		tickerChannel <- time.Now()
		go quoteCleaner.Shutdown(make(chan bool))
		assert.Eventually(t, func() bool {
			return checkFunction() && ticker.AssertExpectations(t) && peginRepository.AssertExpectations(t) && pegoutRepository.AssertExpectations(t)
		}, time.Second, 10*time.Millisecond)
	})
}

func TestQuoteCleanerWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) watcher.Watcher {
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		useCase := w.NewCleanExpiredQuotesUseCase(peginRepository, pegoutRepository)
		return watcher.NewQuoteCleanerWatcher(useCase, ticker)
	})
}

func TestQuoteCleanerWatcher_Prepare(t *testing.T) {
	peginRepository := &mocks.PeginQuoteRepositoryMock{}
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	ticker := &mocks.TickerMock{}
	useCase := w.NewCleanExpiredQuotesUseCase(peginRepository, pegoutRepository)
	quoteCleaner := watcher.NewQuoteCleanerWatcher(useCase, ticker)
	err := quoteCleaner.Prepare(context.Background())
	require.NoError(t, err)
}
