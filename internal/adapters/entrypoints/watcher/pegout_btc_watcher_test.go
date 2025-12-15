package watcher_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPegoutBtcTransferWatcher_Start_SentPegout(t *testing.T) {
	testRetainedQuote := quote.RetainedPegoutQuote{QuoteHash: "010203", DepositAddress: test.AnyAddress, LpBtcTxHash: "040506", State: quote.PegoutStateSendPegoutSucceeded}
	testPegoutQuote := quote.PegoutQuote{Nonce: 5}
	rpc := blockchain.Rpc{}
	eventBus := &mocks.EventBusMock{}
	pegoutSentChannel := make(chan entities.Event)
	eventBus.On("Subscribe", quote.PegoutBtcSentEventId).Return((<-chan entities.Event)(pegoutSentChannel))
	ticker := &mocks.TickerMock{}
	ticker.EXPECT().C().Return(make(chan time.Time))
	ticker.EXPECT().Stop().Return()
	pegoutWatcher := watcher.NewPegoutBtcTransferWatcher(nil, nil, rpc, eventBus, ticker)

	go pegoutWatcher.Start()
	t.Run("handle quote without tx hash", func(t *testing.T) {
		checkFunction := test.AssertLogContains(t, "Quote 010203 doesn't have btc tx hash to watch")
		incomplete := testRetainedQuote
		incomplete.LpBtcTxHash = ""
		pegoutSentChannel <- quote.PegoutBtcSentToUserEvent{
			Event:         entities.NewBaseEvent(quote.PegoutBtcSentEventId),
			PegoutQuote:   testPegoutQuote,
			RetainedQuote: incomplete,
		}
		watchedQuote, ok := pegoutWatcher.GetWatchedQuote(test.AnyString)
		assert.False(t, ok)
		assert.Empty(t, watchedQuote)
		assert.Eventually(t, checkFunction, time.Second, 10*time.Millisecond)
	})
	t.Run("handle sent pegout", func(t *testing.T) {
		defer test.AssertNoLog(t)
		watchedQuote, ok := pegoutWatcher.GetWatchedQuote(test.AnyString)
		assert.False(t, ok)
		assert.Empty(t, watchedQuote)
		pegoutSentChannel <- quote.PegoutBtcSentToUserEvent{
			Event:         entities.NewBaseEvent(quote.PegoutBtcSentEventId),
			PegoutQuote:   testPegoutQuote,
			RetainedQuote: testRetainedQuote,
		}
		assert.Eventually(t, func() bool {
			watchedQuote, ok = pegoutWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			assert.True(t, ok)
			return assert.Equal(t, quote.WatchedPegoutQuote{PegoutQuote: testPegoutQuote, RetainedQuote: testRetainedQuote}, watchedQuote)
		}, time.Second, 10*time.Millisecond)
	})
	t.Run("handle already watched quote", func(t *testing.T) {
		checkFunction := test.AssertLogContains(t, "Quote 010203 is already watched")
		pegoutSentChannel <- quote.PegoutBtcSentToUserEvent{
			Event:       entities.NewBaseEvent(quote.PegoutBtcSentEventId),
			PegoutQuote: testPegoutQuote, RetainedQuote: testRetainedQuote,
		}
		assert.Eventually(t, checkFunction, time.Second, 10*time.Millisecond)
	})
	t.Run("handle incorrect event sent to bus", func(t *testing.T) {
		checkFunction := test.AssertLogContains(t, "Trying to parse wrong event in Pegout Bridge watcher")
		pegoutSentChannel <- quote.CallForUserCompletedEvent{Event: entities.NewBaseEvent(quote.CallForUserCompletedEventId)}
		assert.Eventually(t, checkFunction, time.Second, 10*time.Millisecond)
	})

	closeChannel := make(chan bool)
	go pegoutWatcher.Shutdown(closeChannel)
	<-closeChannel
	assert.Eventually(t, func() bool { return eventBus.AssertExpectations(t) && ticker.AssertExpectations(t) }, time.Second, 10*time.Millisecond)
}

// nolint:funlen,cyclop
func TestPegoutBtcTransferWatcher_Start_BlockchainCheck(t *testing.T) {
	testRetainedQuote := quote.RetainedPegoutQuote{QuoteHash: "070809", DepositAddress: test.AnyAddress, LpBtcTxHash: "030201", State: quote.PegoutStateSendPegoutSucceeded}
	testPegoutQuote := quote.PegoutQuote{Nonce: 5, TransferConfirmations: 5}
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	btcRpc := &mocks.BtcRpcMock{}
	rpc := blockchain.Rpc{Btc: btcRpc}
	eventBus := &mocks.EventBusMock{}
	pegoutSentChannel := make(chan entities.Event)
	eventBus.On("Subscribe", quote.PegoutBtcSentEventId).Return((<-chan entities.Event)(pegoutSentChannel))
	eventBus.On("Publish", mock.Anything).Return(nil)
	tickerChannel := make(chan time.Time)
	ticker := &mocks.TickerMock{}
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop().Return()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return(nil)
	mutex.On("Unlock").Return()
	lbc := &mocks.LiquidityBridgeContractMock{}
	refundPegoutReceipt := blockchain.TransactionReceipt{
		TransactionHash:   test.AnyHash,
		BlockHash:         "0xblock123",
		BlockNumber:       uint64(1000),
		From:              "0x123",
		To:                "0x456",
		CumulativeGasUsed: big.NewInt(21000),
		GasUsed:           big.NewInt(21000),
		Value:             entities.NewWei(0),
		GasPrice:          entities.NewWei(1000000000),
	}
	lbc.On("RefundPegout", mock.Anything, mock.Anything).Return(refundPegoutReceipt, nil).Once()
	refundUseCase := pegout.NewRefundPegoutUseCase(pegoutRepository, blockchain.RskContracts{Lbc: lbc}, eventBus, rpc, mutex)
	pegoutWatcher := watcher.NewPegoutBtcTransferWatcher(nil, refundUseCase, rpc, eventBus, ticker)
	resetMocks := func() {
		btcRpc.Calls = []mock.Call{}
		btcRpc.ExpectedCalls = []*mock.Call{}
		pegoutRepository.Calls = []mock.Call{}
		pegoutRepository.ExpectedCalls = []*mock.Call{}
	}
	go pegoutWatcher.Start()
	t.Run("should only update block upwards", func(t *testing.T) {
		resetMocks()
		btcRpc.On("GetHeight").Return(big.NewInt(5), nil).Once()
		btcRpc.On("GetHeight").Return(big.NewInt(4), nil).Once()
		btcRpc.On("GetHeight").Return(big.NewInt(6), nil).Once()

		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool { return assert.Equal(t, big.NewInt(5), pegoutWatcher.GetCurrentBlock()) }, time.Second, 10*time.Millisecond)

		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool { return assert.Equal(t, big.NewInt(5), pegoutWatcher.GetCurrentBlock()) }, time.Second, 10*time.Millisecond)

		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool { return assert.Equal(t, big.NewInt(6), pegoutWatcher.GetCurrentBlock()) }, time.Second, 10*time.Millisecond)

		btcRpc.AssertExpectations(t)
	})
	t.Run("should handle error getting height", func(t *testing.T) {
		resetMocks()
		checkFunction := test.AssertLogContains(t, "error getting Bitcoin chain height")
		btcRpc.On("GetHeight").Return(nil, assert.AnError).Once()
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool { return checkFunction() && btcRpc.AssertExpectations(t) }, time.Second, 10*time.Millisecond)
	})
	pegoutSentChannel <- quote.PegoutBtcSentToUserEvent{
		Event:       entities.NewBaseEvent(quote.PegoutBtcSentEventId),
		PegoutQuote: testPegoutQuote, RetainedQuote: testRetainedQuote,
	}
	t.Run("should handle error getting tx info", func(t *testing.T) {
		resetMocks()
		checkFunction := test.AssertLogContains(t, "error getting Bitcoin transaction information (030201)")
		btcRpc.On("GetHeight").Return(big.NewInt(8), nil).Once()
		btcRpc.On("GetTransactionInfo", testRetainedQuote.LpBtcTxHash).Return(blockchain.BitcoinTransactionInformation{}, assert.AnError).Once()
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool { return checkFunction() && btcRpc.AssertExpectations(t) }, time.Second, 10*time.Millisecond)
	})
	t.Run("shouldn't refund pegout if transaction is not mature enough", func(t *testing.T) {
		resetMocks()
		btcRpc.On("GetHeight").Return(big.NewInt(9), nil).Once()
		btcRpc.On("GetTransactionInfo", testRetainedQuote.LpBtcTxHash).Return(blockchain.BitcoinTransactionInformation{Confirmations: 1}, nil).Once()
		watchedQuote, ok := pegoutWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
		assert.True(t, ok)
		assert.Equal(t, quote.WatchedPegoutQuote{PegoutQuote: testPegoutQuote, RetainedQuote: testRetainedQuote}, watchedQuote)
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			watchedQuote, ok = pegoutWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			return btcRpc.AssertExpectations(t) && assert.True(t, ok) &&
				assert.Equal(t, quote.WatchedPegoutQuote{PegoutQuote: testPegoutQuote, RetainedQuote: testRetainedQuote}, watchedQuote)
		}, time.Second, 10*time.Millisecond)
	})
	const errorMsg = "Error executing refund pegout on quote"
	t.Run("shouldn't stop tracking quote on recoverable error", func(t *testing.T) {
		resetMocks()
		checkFunction := test.AssertLogContains(t, errorMsg)
		btcRpc.On("GetHeight").Return(big.NewInt(10), nil).Once()
		btcRpc.On("GetTransactionInfo", testRetainedQuote.LpBtcTxHash).Return(blockchain.BitcoinTransactionInformation{Confirmations: 10}, nil).Once()
		pegoutRepository.EXPECT().GetQuote(mock.Anything, testRetainedQuote.QuoteHash).Return(nil, assert.AnError).Once()
		watchedQuote, ok := pegoutWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
		assert.True(t, ok)
		assert.Equal(t, quote.WatchedPegoutQuote{PegoutQuote: testPegoutQuote, RetainedQuote: testRetainedQuote}, watchedQuote)
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			watchedQuote, ok = pegoutWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			return btcRpc.AssertExpectations(t) && assert.True(t, ok) && pegoutRepository.AssertExpectations(t) && checkFunction() &&
				assert.Equal(t, quote.WatchedPegoutQuote{PegoutQuote: testPegoutQuote, RetainedQuote: testRetainedQuote}, watchedQuote)
		}, time.Second, 10*time.Millisecond)
	})
	t.Run("should stop tracking quote on non-recoverable error", func(t *testing.T) {
		resetMocks()
		checkFunction := test.AssertLogContains(t, errorMsg)
		btcRpc.On("GetHeight").Return(big.NewInt(11), nil).Once()
		btcRpc.On("GetTransactionInfo", testRetainedQuote.LpBtcTxHash).Return(blockchain.BitcoinTransactionInformation{Confirmations: 10}, nil).Once()
		pegoutRepository.EXPECT().GetQuote(mock.Anything, testRetainedQuote.QuoteHash).Return(nil, errors.Join(assert.AnError, usecases.NonRecoverableError)).Once()
		watchedQuote, ok := pegoutWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
		assert.True(t, ok)
		assert.Equal(t, quote.WatchedPegoutQuote{PegoutQuote: testPegoutQuote, RetainedQuote: testRetainedQuote}, watchedQuote)
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			watchedQuote, ok = pegoutWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			return btcRpc.AssertExpectations(t) && assert.False(t, ok) && pegoutRepository.AssertExpectations(t) && checkFunction() &&
				assert.Empty(t, watchedQuote)
		}, time.Second, 10*time.Millisecond)
	})
	pegoutSentChannel <- quote.PegoutBtcSentToUserEvent{
		Event:       entities.NewBaseEvent(quote.PegoutBtcSentEventId),
		PegoutQuote: testPegoutQuote, RetainedQuote: testRetainedQuote,
	}
	t.Run("should stop tracking quote on successful refund", func(t *testing.T) {
		resetMocks()
		btcRpc.On("GetHeight").Return(big.NewInt(12), nil).Once()
		btcRpc.On("GetTransactionInfo", testRetainedQuote.LpBtcTxHash).Return(blockchain.BitcoinTransactionInformation{Confirmations: 10}, nil).Twice()
		btcRpc.On("GetTransactionBlockInfo", mock.Anything).Return(blockchain.BitcoinBlockInformation{}, nil).Once()
		btcRpc.On("BuildMerkleBranch", mock.Anything).Return(blockchain.MerkleBranch{}, nil).Once()
		btcRpc.On("GetRawTransaction", mock.Anything).Return([]byte{1, 2, 3}, nil).Once()
		pegoutRepository.EXPECT().GetQuote(mock.Anything, testRetainedQuote.QuoteHash).Return(&testPegoutQuote, nil).Once()
		pegoutRepository.EXPECT().UpdateRetainedQuote(mock.Anything, mock.Anything).Return(nil).Once()
		assert.Eventually(t, func() bool {
			watchedQuote, ok := pegoutWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			return assert.True(t, ok) && assert.Equal(t, quote.WatchedPegoutQuote{PegoutQuote: testPegoutQuote, RetainedQuote: testRetainedQuote}, watchedQuote)
		}, time.Second, 10*time.Millisecond)

		tickerChannel <- time.Now()

		assert.Eventually(t, func() bool {
			watchedQuote, ok := pegoutWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			return assert.False(t, ok) && assert.Empty(t, watchedQuote) && btcRpc.AssertExpectations(t) && pegoutRepository.AssertExpectations(t) && lbc.AssertExpectations(t)
		}, time.Second, 10*time.Millisecond)
	})
	closeChannel := make(chan bool)
	go pegoutWatcher.Shutdown(closeChannel)
	<-closeChannel
	assert.Eventually(t, func() bool { return eventBus.AssertExpectations(t) && ticker.AssertExpectations(t) }, time.Second, 10*time.Millisecond)
}

func TestPegoutBtcTransferWatcher_Prepare(t *testing.T) {
	t.Run("prepare watcher successfully", func(t *testing.T) {
		quotes := []quote.RetainedPegoutQuote{
			{QuoteHash: "pegout1", RequiredLiquidity: entities.NewWei(utils.MustGetRandomInt())},
			{QuoteHash: "pegout2", RequiredLiquidity: entities.NewWei(utils.MustGetRandomInt())},
			{QuoteHash: "pegout3", RequiredLiquidity: entities.NewWei(utils.MustGetRandomInt())},
		}
		quoteRepository := &mocks.PegoutQuoteRepositoryMock{}
		quoteRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, quote.PegoutStateSendPegoutSucceeded).Return(quotes, nil)
		for i, q := range quotes {
			quoteRepository.EXPECT().GetQuote(mock.Anything, q.QuoteHash).
				Return(&quote.PegoutQuote{Value: q.RequiredLiquidity}, nil)
			quoteRepository.EXPECT().GetPegoutCreationData(mock.Anything, mock.Anything).Return(quote.PegoutCreationData{GasPrice: entities.NewWei(int64(i))}).Once()
		}
		useCase := w.NewGetWatchedPegoutQuoteUseCase(quoteRepository)
		pegoutWatcher := watcher.NewPegoutBtcTransferWatcher(useCase, nil, blockchain.Rpc{}, nil, nil)
		err := pegoutWatcher.Prepare(context.Background())
		require.NoError(t, err)
		for i, q := range quotes {
			watchedQuote, ok := pegoutWatcher.GetWatchedQuote(q.QuoteHash)
			require.True(t, ok)
			assert.Equal(t, quote.WatchedPegoutQuote{
				PegoutQuote:   quote.PegoutQuote{Value: q.RequiredLiquidity},
				RetainedQuote: q,
				CreationData:  quote.PegoutCreationData{GasPrice: entities.NewWei(int64(i))},
			}, watchedQuote)
		}
		quoteRepository.AssertExpectations(t)
	})
	t.Run("handle error preparing watcher", func(t *testing.T) {
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		pegoutRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		useCase := w.NewGetWatchedPegoutQuoteUseCase(pegoutRepository)
		addressWatcher := watcher.NewPegoutBtcTransferWatcher(useCase, nil, blockchain.Rpc{}, nil, nil)
		err := addressWatcher.Prepare(context.Background())
		require.Error(t, err)
		pegoutRepository.AssertExpectations(t)
	})
}

func TestPegoutBtcTransferWatcher_Shutdown(t *testing.T) {
	eventBus := &mocks.EventBusMock{}
	eventBus.On("Subscribe", mock.Anything).Return(make(<-chan entities.Event))
	createWatcherShutdownTest(t, func(ticker utils.Ticker) watcher.Watcher {
		return watcher.NewPegoutBtcTransferWatcher(nil, nil, blockchain.Rpc{}, eventBus, ticker)
	})
}

func TestPegoutBtcTransferWatcher(t *testing.T) {
	t.Run("watcher doesn't run into a deadlock", func(t *testing.T) {
		eventBus := &mocks.EventBusMock{}
		eventBus.On("Subscribe", quote.PegoutBtcSentEventId).Return((<-chan entities.Event)(make(chan entities.Event)))
		ticker := &mocks.TickerMock{}
		tickerChannel := make(chan time.Time)
		shutdownChannel := make(chan bool)
		ticker.EXPECT().C().Return(tickerChannel)
		ticker.EXPECT().Stop().Return()
		btcRpc := &mocks.BtcRpcMock{}
		rpc := blockchain.Rpc{Btc: btcRpc}
		btcRpc.On("GetHeight").Return(big.NewInt(0), nil).Once()
		quoteRepository := &mocks.PegoutQuoteRepositoryMock{}
		quoteRepository.EXPECT().
			GetRetainedQuoteByState(mock.Anything, quote.PegoutStateSendPegoutSucceeded).
			After(time.Second*2).
			Return([]quote.RetainedPegoutQuote{}, nil)
		useCase := w.NewGetWatchedPegoutQuoteUseCase(quoteRepository)
		pegoutWatcher := watcher.NewPegoutBtcTransferWatcher(useCase, nil, rpc, eventBus, ticker)

		prepareDoneChannel := make(chan bool, 1)
		startDoneChannel := make(chan bool, 1)
		go assert.NotPanics(t, func() {
			err := pegoutWatcher.Prepare(context.Background())
			require.NoError(t, err)
			prepareDoneChannel <- true
		})
		go assert.NotPanics(t, func() {
			pegoutWatcher.Start()
			startDoneChannel <- true
		})

		tickerChannel <- time.Now()
		go pegoutWatcher.Shutdown(shutdownChannel)
		<-shutdownChannel
		require.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.NotEmpty(c, prepareDoneChannel)
			assert.NotEmpty(c, startDoneChannel)
		}, time.Second*5, time.Millisecond*100)
	})
}
