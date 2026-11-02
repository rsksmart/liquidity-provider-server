package watcher_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
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
	t.Run("error loading claimed quotes", func(t *testing.T) {
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		pegoutRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, quote.PegoutStateClaimed).Return(nil, assert.AnError).Once()
		rskRpc := &mocks.RootstockRpcServerMock{}
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(100), nil).Once()
		useCases := watcher.NewPegoutRskDepositWatcherUseCases(w.NewGetWatchedPegoutQuoteUseCase(pegoutRepository), nil)
		depositWatcher := watcher.NewPegoutRskDepositWatcher(useCases, blockchain.Rpc{Rsk: rskRpc}, nil, nil, time.Duration(1))
		err := depositWatcher.Prepare(context.Background())
		require.ErrorIs(t, err, assert.AnError)
		pegoutRepository.AssertExpectations(t)
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
		checkFunction := test.LogContains(t, watcher.LogPegoutRskAlreadyWatched(testRetainedQuote.QuoteHash))
		claimedChannel <- quote.ClaimedPegoutQuoteEvent{
			Event:         entities.NewBaseEvent(quote.ClaimedPegoutQuoteEventId),
			Quote:         testPegoutQuote,
			RetainedQuote: testRetainedQuote,
		}
		assert.Eventually(t, checkFunction, time.Second, 10*time.Millisecond)
	})
	t.Run("handle incorrect event", func(t *testing.T) {
		checkFunction := test.LogContains(t, watcher.LogPegoutRskWrongEvent)
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

func pegoutClaimedSendQuotes() (quote.PegoutQuote, quote.RetainedPegoutQuote) {
	return quote.PegoutQuote{Nonce: 1, Value: entities.NewWei(3), GasFee: entities.NewWei(1)},
		quote.RetainedPegoutQuote{
			QuoteHash: "0102030000000000000000000000000000000000000000000000000000000000",
			State:     quote.PegoutStateClaimed,
		}
}

func TestPegoutRskDepositWatcher_SendClaimed_Success(t *testing.T) {
	pegoutQuote, retainedQuote := pegoutClaimedSendQuotes()
	f := startPegoutClaimedSendWatcher(t, pegoutQuote, retainedQuote)
	rawTx := []byte{0x01, 0x02}
	f.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(1), nil).Once()
	f.pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil).Once()
	f.escrow.EXPECT().GetPegOutQuote(retainedQuote.QuoteHash).Return(pegoutQuote, nil).Once()
	f.escrow.EXPECT().GetPegOutState(retainedQuote.QuoteHash).Return(blockchain.EscrowedPegOutStateClaimed, nil).Once()
	f.quoteRepository.EXPECT().GetRetainedQuote(mock.Anything, retainedQuote.QuoteHash).Return(&retainedQuote, nil).Once()
	f.btcWallet.On("GetBalance").Return(entities.NewWei(10000), nil).Once()
	f.btcWallet.On("CreateUnfundedTransactionWithOpReturn", mock.Anything, pegoutQuote.Value, mock.Anything).Return(rawTx, nil).Once()
	f.pegoutContract.EXPECT().ValidatePegout(retainedQuote.QuoteHash, rawTx).Return(nil).Once()
	f.btcWallet.On("SendWithOpReturn", mock.Anything, pegoutQuote.Value, mock.Anything).
		Return(blockchain.BitcoinTransactionResult{Hash: test.AnyHash, Fee: entities.NewWei(1)}, nil).Once()
	f.quoteRepository.EXPECT().GetPegoutCreationData(mock.Anything, retainedQuote.QuoteHash).Return(quote.PegoutCreationDataZeroValue()).Once()
	f.eventBus.On("Publish", mock.AnythingOfType("quote.PegoutBtcSentToUserEvent")).Return().Once()
	f.quoteRepository.EXPECT().UpdateRetainedQuote(mock.Anything, mock.MatchedBy(func(q quote.RetainedPegoutQuote) bool {
		return q.State == quote.PegoutStateSendPegoutSucceeded && q.LpBtcTxHash == test.AnyHash
	})).Return(nil).Once()

	f.tickerChannel <- time.Now()

	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		_, ok := f.watcher.GetWatchedQuote(retainedQuote.QuoteHash)
		assert.False(collect, ok)
	}, time.Second, 10*time.Millisecond)
	f.assertExpectations(t)
}

func TestPegoutRskDepositWatcher_SendClaimed_RecoverableError(t *testing.T) {
	pegoutQuote, retainedQuote := pegoutClaimedSendQuotes()
	f := startPegoutClaimedSendWatcher(t, pegoutQuote, retainedQuote)
	f.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(1), nil).Once()
	f.pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil).Once()
	f.escrow.EXPECT().GetPegOutQuote(retainedQuote.QuoteHash).Return(quote.PegoutQuote{}, assert.AnError).Once()

	f.tickerChannel <- time.Now()

	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		f.escrow.AssertExpectations(newMockCollectT(collect))
	}, time.Second, 10*time.Millisecond)
	_, ok := f.watcher.GetWatchedQuote(retainedQuote.QuoteHash)
	assert.True(t, ok)
	f.assertExpectations(t)
}

func TestPegoutRskDepositWatcher_SendClaimed_NonRecoverableError(t *testing.T) {
	pegoutQuote, retainedQuote := pegoutClaimedSendQuotes()
	f := startPegoutClaimedSendWatcher(t, pegoutQuote, retainedQuote)
	invalidQuote := pegoutQuote
	invalidQuote.DepositAddress = "not-hex"
	f.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(1), nil).Once()
	f.pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil).Once()
	f.escrow.EXPECT().GetPegOutQuote(retainedQuote.QuoteHash).Return(invalidQuote, nil).Once()
	f.quoteRepository.EXPECT().UpdateRetainedQuote(mock.Anything, mock.MatchedBy(func(q quote.RetainedPegoutQuote) bool {
		return q.State == quote.PegoutStateSendPegoutFailed
	})).Return(nil).Once()
	f.eventBus.On("Publish", mock.AnythingOfType("quote.PegoutBtcSentToUserEvent")).Return().Once()

	f.tickerChannel <- time.Now()

	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		_, ok := f.watcher.GetWatchedQuote(retainedQuote.QuoteHash)
		assert.False(collect, ok)
	}, time.Second, 10*time.Millisecond)
	f.assertExpectations(t)
	f.btcWallet.AssertNotCalled(t, "SendWithOpReturn", mock.Anything, mock.Anything, mock.Anything)
}

func TestPegoutRskDepositWatcher_SendClaimed_NoNewBlock(t *testing.T) {
	pegoutQuote, retainedQuote := pegoutClaimedSendQuotes()
	f := startPegoutClaimedSendWatcher(t, pegoutQuote, retainedQuote)
	f.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(0), nil).Once()

	f.tickerChannel <- time.Now()

	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		f.rskRpc.AssertExpectations(newMockCollectT(collect))
	}, time.Second, 10*time.Millisecond)
	_, ok := f.watcher.GetWatchedQuote(retainedQuote.QuoteHash)
	assert.True(t, ok)
	f.escrow.AssertNotCalled(t, "GetPegOutQuote", mock.Anything)
}

func TestPegoutRskDepositWatcher_SendClaimed_ChainHeightError(t *testing.T) {
	pegoutQuote, retainedQuote := pegoutClaimedSendQuotes()
	f := startPegoutClaimedSendWatcher(t, pegoutQuote, retainedQuote)
	checkFunction := test.LogContains(t, fmt.Sprintf(watcher.LogPegoutRskChainHeight, assert.AnError))
	f.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(0), assert.AnError).Once()

	f.tickerChannel <- time.Now()

	assert.Eventually(t, checkFunction, time.Second, 10*time.Millisecond)
	_, ok := f.watcher.GetWatchedQuote(retainedQuote.QuoteHash)
	assert.True(t, ok)
	f.rskRpc.AssertExpectations(t)
	f.escrow.AssertNotCalled(t, "GetPegOutQuote", mock.Anything)
}

type pegoutClaimedSendFixture struct {
	watcher         *watcher.PegoutRskDepositWatcher
	tickerChannel   chan time.Time
	rskRpc          *mocks.RootstockRpcServerMock
	pegoutContract  *mocks.PegoutContractMock
	escrow          *mocks.PegOutEscrowContractMock
	quoteRepository *mocks.PegoutQuoteRepositoryMock
	btcWallet       *mocks.BitcoinWalletMock
	eventBus        *mocks.EventBusMock
}

func startPegoutClaimedSendWatcher(
	t *testing.T,
	pegoutQuote quote.PegoutQuote,
	retainedQuote quote.RetainedPegoutQuote,
) *pegoutClaimedSendFixture {
	t.Helper()
	f := &pegoutClaimedSendFixture{
		tickerChannel:   make(chan time.Time),
		rskRpc:          &mocks.RootstockRpcServerMock{},
		pegoutContract:  &mocks.PegoutContractMock{},
		escrow:          &mocks.PegOutEscrowContractMock{},
		quoteRepository: &mocks.PegoutQuoteRepositoryMock{},
		btcWallet:       &mocks.BitcoinWalletMock{},
		eventBus:        &mocks.EventBusMock{},
	}
	ticker := &mocks.TickerMock{}
	ticker.EXPECT().C().Return(f.tickerChannel)
	ticker.EXPECT().Stop().Return()
	claimedChannel := make(chan entities.Event)
	f.eventBus.On("Subscribe", quote.ClaimedPegoutQuoteEventId).Return((<-chan entities.Event)(claimedChannel))

	rpc := blockchain.Rpc{Rsk: f.rskRpc}
	contracts := blockchain.RskContracts{PegOut: f.pegoutContract, PegOutEscrow: f.escrow}
	sendPegoutUseCase := pegout.NewSendPegoutUseCase(
		f.btcWallet,
		f.quoteRepository,
		rpc,
		f.eventBus,
		contracts,
		environment.NewApplicationMutexes().BtcWalletMutex(),
		rootstock.ParseDepositEventByQuoteHash,
	)
	useCases := watcher.NewPegoutRskDepositWatcherUseCases(nil, sendPegoutUseCase)
	f.watcher = watcher.NewPegoutRskDepositWatcher(useCases, rpc, f.eventBus, ticker, time.Duration(1))

	go f.watcher.Start()
	t.Cleanup(func() {
		closeChannel := make(chan bool)
		go f.watcher.Shutdown(closeChannel)
		<-closeChannel
	})

	claimedChannel <- quote.ClaimedPegoutQuoteEvent{
		Event:         entities.NewBaseEvent(quote.ClaimedPegoutQuoteEventId),
		Quote:         pegoutQuote,
		RetainedQuote: retainedQuote,
	}
	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		_, ok := f.watcher.GetWatchedQuote(retainedQuote.QuoteHash)
		assert.True(collect, ok)
	}, time.Second, 10*time.Millisecond)
	return f
}

func (f *pegoutClaimedSendFixture) assertExpectations(t *testing.T) {
	f.rskRpc.AssertExpectations(t)
	f.pegoutContract.AssertExpectations(t)
	f.escrow.AssertExpectations(t)
	f.quoteRepository.AssertExpectations(t)
	f.btcWallet.AssertExpectations(t)
	f.eventBus.AssertExpectations(t)
}
