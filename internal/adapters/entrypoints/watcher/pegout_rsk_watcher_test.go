package watcher_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewPegoutRskDepositWatcher(t *testing.T) {
	ticker := &mocks.TickerMock{}
	eventBus := &mocks.EventBusMock{}
	rpc := blockchain.Rpc{}
	useCases := watcher.NewPegoutRskDepositWatcherUseCases(
		&w.GetWatchedPegoutQuoteUseCase{},
		&pegout.SendPegoutUseCase{},
	)
	depositWatcher := watcher.NewPegoutRskDepositWatcher(useCases, rpc, eventBus, ticker, time.Duration(1))
	require.NotNil(t, depositWatcher)
}

func TestPegoutRskDepositWatcher_Prepare(t *testing.T) {
	t.Run("loads claimed quotes", func(t *testing.T) {
		claimed := []quote.RetainedPegoutQuote{
			{QuoteHash: "aa", State: quote.PegoutStateClaimed, RequiredLiquidity: entities.NewWei(1)},
			{QuoteHash: "bb", State: quote.PegoutStateClaimed, RequiredLiquidity: entities.NewWei(2)},
		}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		pegoutRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, quote.PegoutStateClaimed).Return(claimed, nil).Once()
		for _, q := range claimed {
			pegoutRepository.EXPECT().GetQuote(mock.Anything, q.QuoteHash).Return(&quote.PegoutQuote{Nonce: 1}, nil).Once()
			pegoutRepository.EXPECT().GetPegoutCreationData(mock.Anything, q.QuoteHash).Return(quote.PegoutCreationData{}).Once()
		}
		rskRpc := &mocks.RootstockRpcServerMock{}
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(100), nil).Once()
		rpc := blockchain.Rpc{Rsk: rskRpc}
		getWatchedQuotesUseCase := w.NewGetWatchedPegoutQuoteUseCase(pegoutRepository)
		useCases := watcher.NewPegoutRskDepositWatcherUseCases(getWatchedQuotesUseCase, nil)
		depositWatcher := watcher.NewPegoutRskDepositWatcher(useCases, rpc, nil, nil, time.Duration(1))
		err := depositWatcher.Prepare(context.Background())
		require.NoError(t, err)
		assert.Equal(t, uint64(100), depositWatcher.GetCurrentBlock())
		for _, q := range claimed {
			watched, ok := depositWatcher.GetWatchedQuote(q.QuoteHash)
			assert.True(t, ok)
			assert.Equal(t, q, watched.RetainedQuote)
		}
		pegoutRepository.AssertExpectations(t)
		rskRpc.AssertExpectations(t)
	})
	t.Run("error getting height", func(t *testing.T) {
		rskRpc := &mocks.RootstockRpcServerMock{}
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(0), assert.AnError).Once()
		useCases := watcher.NewPegoutRskDepositWatcherUseCases(nil, nil)
		depositWatcher := watcher.NewPegoutRskDepositWatcher(useCases, blockchain.Rpc{Rsk: rskRpc}, nil, nil, time.Duration(1))
		err := depositWatcher.Prepare(context.Background())
		require.Error(t, err)
		rskRpc.AssertExpectations(t)
	})
}

func TestPegoutRskDepositWatcher_Shutdown(t *testing.T) {
	eventBus := &mocks.EventBusMock{}
	eventBus.On("Subscribe", mock.Anything).Return(make(<-chan entities.Event))
	createWatcherShutdownTest(t, func(ticker utils.Ticker) watcher.Watcher {
		return watcher.NewPegoutRskDepositWatcher(
			watcher.NewPegoutRskDepositWatcherUseCases(nil, nil),
			blockchain.Rpc{},
			eventBus,
			ticker,
			time.Duration(1),
		)
	})
}

func TestPegoutRskDepositWatcher_Start_Claimed(t *testing.T) {
	ticker := &mocks.TickerMock{}
	ticker.EXPECT().C().Return(make(chan time.Time))
	ticker.EXPECT().Stop().Return()
	rskRpc := &mocks.RootstockRpcServerMock{}
	rpc := blockchain.Rpc{Rsk: rskRpc}
	eventBus := &mocks.EventBusMock{}
	claimedChannel := make(chan entities.Event)
	eventBus.On("Subscribe", quote.ClaimedPegoutQuoteEventId).Return((<-chan entities.Event)(claimedChannel))

	testPegoutQuote := quote.PegoutQuote{Nonce: 1}
	testRetainedQuote := quote.RetainedPegoutQuote{QuoteHash: "010203", State: quote.PegoutStateClaimed}

	useCases := watcher.NewPegoutRskDepositWatcherUseCases(nil, nil)
	depositWatcher := watcher.NewPegoutRskDepositWatcher(useCases, rpc, eventBus, ticker, time.Duration(1))
	go depositWatcher.Start()

	t.Run("handle claimed pegout quote", func(t *testing.T) {
		defer test.AssertNoLog(t)()
		claimedChannel <- quote.ClaimedPegoutQuoteEvent{
			Event:         entities.NewBaseEvent(quote.ClaimedPegoutQuoteEventId),
			Quote:         testPegoutQuote,
			RetainedQuote: testRetainedQuote,
		}
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			watchedQuote, ok := depositWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			assert.True(collect, ok)
			assert.Equal(collect, testPegoutQuote, watchedQuote.PegoutQuote)
			assert.Equal(collect, testRetainedQuote, watchedQuote.RetainedQuote)
		}, time.Second, 10*time.Millisecond)
	})
	t.Run("handle already watched quote", func(t *testing.T) {
		checkFunction := test.AssertLogContains(t, watcher.LogPegoutRskAlreadyWatched(testRetainedQuote.QuoteHash))
		claimedChannel <- quote.ClaimedPegoutQuoteEvent{
			Event:         entities.NewBaseEvent(quote.ClaimedPegoutQuoteEventId),
			Quote:         testPegoutQuote,
			RetainedQuote: testRetainedQuote,
		}
		assert.Eventually(t, checkFunction, time.Second, 10*time.Millisecond)
	})
	t.Run("handle incorrect event", func(t *testing.T) {
		checkFunction := test.AssertLogContains(t, watcher.LogPegoutRskWrongEvent)
		claimedChannel <- quote.AcceptedPeginQuoteEvent{Event: entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId)}
		assert.Eventually(t, checkFunction, time.Second, 10*time.Millisecond)
	})

	closeChannel := make(chan bool)
	go depositWatcher.Shutdown(closeChannel)
	<-closeChannel
	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		mt := newMockCollectT(collect)
		eventBus.AssertExpectations(mt)
		ticker.AssertExpectations(mt)
	}, time.Second, 10*time.Millisecond)
}
