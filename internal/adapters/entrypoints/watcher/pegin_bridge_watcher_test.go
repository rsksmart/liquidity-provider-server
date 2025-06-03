package watcher_test

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

func TestPeginBridgeWatcher_Prepare(t *testing.T) {
	quotes := []quote.RetainedPeginQuote{
		{QuoteHash: "pegin1", RequiredLiquidity: entities.NewWei(utils.MustGetRandomInt())},
		{QuoteHash: "pegin2", RequiredLiquidity: entities.NewWei(utils.MustGetRandomInt())},
		{QuoteHash: "pegin3", RequiredLiquidity: entities.NewWei(utils.MustGetRandomInt())},
	}
	quoteRepository := &mocks.PeginQuoteRepositoryMock{}
	quoteRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, quote.PeginStateCallForUserSucceeded).Return(quotes, nil)
	for i, q := range quotes {
		quoteRepository.EXPECT().GetQuote(mock.Anything, q.QuoteHash).
			Return(&quote.PeginQuote{Value: q.RequiredLiquidity}, nil)
		quoteRepository.EXPECT().GetPeginCreationData(mock.Anything, q.QuoteHash).
			Return(quote.PeginCreationData{GasPrice: entities.NewWei(int64(i))})
	}
	useCase := w.NewGetWatchedPeginQuoteUseCase(quoteRepository)
	peginWatcher := watcher.NewPeginBridgeWatcher(nil, useCase, blockchain.RskContracts{}, blockchain.Rpc{}, nil, nil)
	err := peginWatcher.Prepare(context.Background())
	require.NoError(t, err)
	for i, q := range quotes {
		watchedQuote, ok := peginWatcher.GetWatchedQuote(q.QuoteHash)
		require.True(t, ok)
		assert.Equal(t, quote.WatchedPeginQuote{
			PeginQuote:    quote.PeginQuote{Value: q.RequiredLiquidity},
			RetainedQuote: q,
			CreationData:  quote.PeginCreationData{GasPrice: entities.NewWei(int64(i))},
		}, watchedQuote)
	}
	quoteRepository.AssertExpectations(t)
}

func TestPeginBridgeWatcher_Prepare_ErrorHandling(t *testing.T) {
	qupteRepository := &mocks.PeginQuoteRepositoryMock{}
	qupteRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, mock.Anything).Return(nil, assert.AnError)
	useCase := w.NewGetWatchedPeginQuoteUseCase(qupteRepository)
	peginWatcher := watcher.NewPeginBridgeWatcher(nil, useCase, blockchain.RskContracts{}, blockchain.Rpc{}, nil, nil)
	err := peginWatcher.Prepare(context.Background())
	require.Error(t, err)
	qupteRepository.AssertExpectations(t)
}

func TestPeginBridgeWatcher_Shutdown(t *testing.T) {
	eventBus := &mocks.EventBusMock{}
	eventBus.On("Subscribe", mock.Anything).Return(make(<-chan entities.Event))
	createWatcherShutdownTest(t, func(ticker watcher.Ticker) watcher.Watcher {
		return watcher.NewPeginBridgeWatcher(nil, nil, blockchain.RskContracts{}, blockchain.Rpc{}, eventBus, ticker)
	})
	eventBus.AssertExpectations(t)
}

func TestPeginBridgeWatcher_Start_CfuCompleted(t *testing.T) {
	quoteRepository := &mocks.PeginQuoteRepositoryMock{}
	contracts := blockchain.RskContracts{}
	rpc := blockchain.Rpc{}
	eventBus := &mocks.EventBusMock{}
	cfuChannel := make(chan entities.Event)
	eventBus.On("Subscribe", quote.CallForUserCompletedEventId).Return((<-chan entities.Event)(cfuChannel))
	appMutexes := environment.NewApplicationMutexes()
	getUseCase := w.NewGetWatchedPeginQuoteUseCase(quoteRepository)
	registerUseCase := pegin.NewRegisterPeginUseCase(contracts, quoteRepository, eventBus, rpc, appMutexes.RskWalletMutex())
	ticker := &mocks.TickerMock{}
	ticker.EXPECT().C().Return(make(chan time.Time))
	ticker.EXPECT().Stop().Return()
	peginWatcher := watcher.NewPeginBridgeWatcher(registerUseCase, getUseCase, contracts, rpc, eventBus, ticker)

	go peginWatcher.Start()

	testPeginQuote := quote.PeginQuote{Nonce: quote.NewNonce(1)}
	testRetainedQuote := quote.RetainedPeginQuote{QuoteHash: test.AnyString, State: quote.PeginStateCallForUserSucceeded}

	t.Run("handle call for user performed", func(t *testing.T) {
		watchedQuote, ok := peginWatcher.GetWatchedQuote(test.AnyString)
		assert.False(t, ok)
		assert.Empty(t, watchedQuote)
		cfuChannel <- quote.CallForUserCompletedEvent{
			Event:         entities.NewBaseEvent(quote.CallForUserCompletedEventId),
			PeginQuote:    testPeginQuote,
			RetainedQuote: testRetainedQuote,
		}
		assert.Eventually(t, func() bool {
			watchedQuote, ok = peginWatcher.GetWatchedQuote(test.AnyString)
			assert.True(t, ok)
			return assert.Equal(t, quote.WatchedPeginQuote{
				PeginQuote:    testPeginQuote,
				RetainedQuote: testRetainedQuote,
			}, watchedQuote)
		}, time.Second, 10*time.Millisecond)
	})
	t.Run("handle already watched quote", func(t *testing.T) {
		checkFunction := test.AssertLogContains(t, "Quote any value is already watched")
		cfuChannel <- quote.CallForUserCompletedEvent{
			Event:         entities.NewBaseEvent(quote.CallForUserCompletedEventId),
			PeginQuote:    testPeginQuote,
			RetainedQuote: testRetainedQuote,
		}
		assert.Eventually(t, checkFunction, time.Second, 10*time.Millisecond)
	})
	t.Run("handle incorrect event sent to bus", func(t *testing.T) {
		checkFunction := test.AssertLogContains(t, "Trying to parse wrong event")
		cfuChannel <- quote.PegoutQuoteCompletedEvent{
			Event: entities.NewBaseEvent(quote.PegoutQuoteCompletedEventId),
		}
		assert.Eventually(t, checkFunction, time.Second, 10*time.Millisecond)
	})
	closeChannel := make(chan bool)
	go peginWatcher.Shutdown(closeChannel)
	<-closeChannel
	assert.Eventually(t, func() bool {
		return eventBus.AssertExpectations(t) && ticker.AssertExpectations(t)
	}, time.Second, 10*time.Millisecond)
}

// nolint:funlen
func TestPeginBridgeWatcher_Start_BlockchainCheck(t *testing.T) {
	quoteHash := hex.EncodeToString([]byte{0x20})
	userTx := hex.EncodeToString([]byte{0x12})
	testPeginQuote := quote.PeginQuote{Nonce: quote.NewNonce(1)}
	testRetainedQuote := quote.RetainedPeginQuote{
		QuoteHash:     quoteHash,
		State:         quote.PeginStateCallForUserSucceeded,
		UserBtcTxHash: userTx,
	}

	quoteRepository := &mocks.PeginQuoteRepositoryMock{}
	quoteRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, quote.PeginStateCallForUserSucceeded).Return([]quote.RetainedPeginQuote{}, nil)
	bridge := &mocks.BridgeMock{}
	lbc := &mocks.LbcMock{}
	lbc.On("RegisterPegin", mock.Anything).Return(test.AnyHash, nil)
	contracts := blockchain.RskContracts{Bridge: bridge, Lbc: lbc}
	btcRpc := &mocks.BtcRpcMock{}
	rpc := blockchain.Rpc{Btc: btcRpc}
	eventBus := &mocks.EventBusMock{}
	cfuChannel := make(chan entities.Event)
	eventBus.On("Subscribe", quote.CallForUserCompletedEventId).Return((<-chan entities.Event)(cfuChannel))
	eventBus.On("Publish", mock.Anything).Return()
	appMutexes := environment.NewApplicationMutexes()
	getUseCase := w.NewGetWatchedPeginQuoteUseCase(quoteRepository)
	registerUseCase := pegin.NewRegisterPeginUseCase(contracts, quoteRepository, eventBus, rpc, appMutexes.RskWalletMutex())
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop().Return()
	peginWatcher := watcher.NewPeginBridgeWatcher(registerUseCase, getUseCase, contracts, rpc, eventBus, ticker)
	resetMocks := func() {
		btcRpc.ExpectedCalls = []*mock.Call{}
		btcRpc.Calls = []mock.Call{}
		bridge.ExpectedCalls = []*mock.Call{}
		bridge.Calls = []mock.Call{}
		quoteRepository.ExpectedCalls = []*mock.Call{}
		quoteRepository.Calls = []mock.Call{}
	}

	prepareErr := peginWatcher.Prepare(context.Background())
	require.NoError(t, prepareErr)
	go peginWatcher.Start()

	quoteRepository.AssertExpectations(t)
	t.Run("should only update current block upwards", func(t *testing.T) {
		resetMocks()
		btcRpc.On("GetHeight").Return(big.NewInt(5), nil).Once()
		btcRpc.On("GetHeight").Return(big.NewInt(4), nil).Once()
		btcRpc.On("GetHeight").Return(big.NewInt(7), nil).Once()

		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			return assert.Equal(t, big.NewInt(5), peginWatcher.GetCurrentBlock())
		}, time.Second, 10*time.Millisecond)

		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			return assert.Equal(t, big.NewInt(5), peginWatcher.GetCurrentBlock())
		}, time.Second, 10*time.Millisecond)

		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			return assert.Equal(t, big.NewInt(7), peginWatcher.GetCurrentBlock())
		}, time.Second, 10*time.Millisecond)

		btcRpc.AssertExpectations(t)
	})

	t.Run("should not run register pegin on an unconfirmed quote", func(t *testing.T) {
		resetMocks()
		bridge.On("GetRequiredTxConfirmations").Return(uint64(10)).Once()
		btcRpc.On("GetHeight").Return(big.NewInt(10), nil).Once()
		btcRpc.On("GetTransactionInfo", testRetainedQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
			Hash:          userTx,
			Confirmations: 9,
		}, nil).Once()
		cfuChannel <- quote.CallForUserCompletedEvent{
			Event:         entities.NewBaseEvent(quote.CallForUserCompletedEventId),
			PeginQuote:    testPeginQuote,
			RetainedQuote: testRetainedQuote,
		}
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			return btcRpc.AssertExpectations(t) && bridge.AssertExpectations(t)
		}, time.Second, 10*time.Millisecond)
	})

	t.Run("should run register pegin on a confirmed quote", func(t *testing.T) {
		const errorMsg = "Error executing register pegin on quote 20:"
		t.Run("should continue watching quote on recoverable error", func(t *testing.T) {
			resetMocks()
			defer test.AssertLogContains(t, errorMsg)()
			watchedQuote, ok := peginWatcher.GetWatchedQuote(quoteHash)
			assert.True(t, ok)
			assert.Equal(t, quote.WatchedPeginQuote{
				PeginQuote:    testPeginQuote,
				RetainedQuote: testRetainedQuote,
			}, watchedQuote)
			bridge.On("GetRequiredTxConfirmations").Return(uint64(10)).Once()
			btcRpc.On("GetHeight").Return(big.NewInt(12), nil).Once()
			btcRpc.On("GetTransactionInfo", testRetainedQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
				Hash:          userTx,
				Confirmations: 10,
			}, nil).Once()
			quoteRepository.EXPECT().GetQuote(mock.Anything, quoteHash).Return(nil, assert.AnError).Once()
			tickerChannel <- time.Now()
			assert.Eventually(t, func() bool {
				return btcRpc.AssertExpectations(t) && bridge.AssertExpectations(t) && quoteRepository.AssertExpectations(t)
			}, time.Second, 10*time.Millisecond)
			watchedQuote, ok = peginWatcher.GetWatchedQuote(quoteHash)
			assert.True(t, ok)
			assert.Equal(t, quote.WatchedPeginQuote{
				PeginQuote:    testPeginQuote,
				RetainedQuote: testRetainedQuote,
			}, watchedQuote)
		})
		t.Run("should stop watching quote on unrecoverable error", func(t *testing.T) {
			resetMocks()
			defer test.AssertLogContains(t, errorMsg)()
			watchedQuote, ok := peginWatcher.GetWatchedQuote(quoteHash)
			assert.True(t, ok)
			assert.Equal(t, quote.WatchedPeginQuote{
				PeginQuote:    testPeginQuote,
				RetainedQuote: testRetainedQuote,
			}, watchedQuote)
			bridge.On("GetRequiredTxConfirmations").Return(uint64(10)).Once()
			btcRpc.On("GetHeight").Return(big.NewInt(13), nil).Once()
			btcRpc.On("GetTransactionInfo", testRetainedQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
				Hash:          test.AnyHash,
				Confirmations: 10,
			}, nil).Once()
			quoteRepository.EXPECT().GetQuote(mock.Anything, quoteHash).Return(nil, errors.Join(assert.AnError, usecases.NonRecoverableError)).Once()
			tickerChannel <- time.Now()
			assert.Eventually(t, func() bool {
				return btcRpc.AssertExpectations(t) && bridge.AssertExpectations(t) && quoteRepository.AssertExpectations(t)
			}, time.Second, 10*time.Millisecond)
			watchedQuote, ok = peginWatcher.GetWatchedQuote(quoteHash)
			assert.False(t, ok)
			assert.Empty(t, watchedQuote)
		})
		t.Run("should stop watching quote on successful register", func(t *testing.T) {
			resetMocks()
			cfuChannel <- quote.CallForUserCompletedEvent{
				Event:         entities.NewBaseEvent(quote.CallForUserCompletedEventId),
				PeginQuote:    testPeginQuote,
				RetainedQuote: testRetainedQuote,
			}
			time.Sleep(time.Second)
			watchedQuote, ok := peginWatcher.GetWatchedQuote(quoteHash)
			assert.True(t, ok)
			assert.Equal(t, quote.WatchedPeginQuote{
				PeginQuote:    testPeginQuote,
				RetainedQuote: testRetainedQuote,
			}, watchedQuote)
			bridge.On("GetRequiredTxConfirmations").Return(uint64(10)).Twice()
			btcRpc.On("GetHeight").Return(big.NewInt(14), nil).Once()
			btcRpc.On("GetTransactionInfo", testRetainedQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{
				Hash:          userTx,
				Confirmations: 10,
			}, nil).Twice()
			btcRpc.On("GetRawTransaction", mock.Anything).Return([]byte{0x01}, nil).Once()
			btcRpc.On("GetPartialMerkleTree", mock.Anything).Return([]byte{0x01}, nil).Once()
			btcRpc.On("GetTransactionBlockInfo", mock.Anything).Return(blockchain.BitcoinBlockInformation{}, nil).Once()
			quoteRepository.EXPECT().GetQuote(mock.Anything, quoteHash).Return(&testPeginQuote, nil).Once()
			quoteRepository.EXPECT().UpdateRetainedQuote(mock.Anything, mock.Anything).Return(nil).Once()
			tickerChannel <- time.Now()
			assert.Eventually(t, func() bool {
				return btcRpc.AssertExpectations(t) && bridge.AssertExpectations(t) && quoteRepository.AssertExpectations(t)
			}, time.Second, 10*time.Millisecond)
			watchedQuote, ok = peginWatcher.GetWatchedQuote(quoteHash)
			assert.False(t, ok)
			assert.Empty(t, watchedQuote)
		})
	})

	t.Run("should handle error getting height", func(t *testing.T) {
		resetMocks()
		checkFunction := test.AssertLogContains(t, assert.AnError.Error())
		btcRpc.On("GetHeight").Return(nil, assert.AnError).Once()
		tickerChannel <- time.Now()
		assert.Eventually(t, checkFunction, time.Second, 10*time.Millisecond)
		btcRpc.AssertExpectations(t)
		bridge.AssertNotCalled(t, "GetRequiredTxConfirmations")
	})

	t.Run("should handle error getting tx information", func(t *testing.T) {
		resetMocks()
		checkFunction := test.AssertLogContains(t, assert.AnError.Error())
		btcRpc.On("GetHeight").Return(big.NewInt(15), nil).Once()
		btcRpc.On("GetTransactionInfo", testRetainedQuote.UserBtcTxHash).Return(blockchain.BitcoinTransactionInformation{}, assert.AnError).Once()
		failRetainedQuote := testRetainedQuote
		failRetainedQuote.QuoteHash = "fail"
		cfuChannel <- quote.CallForUserCompletedEvent{
			Event:         entities.NewBaseEvent(quote.CallForUserCompletedEventId),
			PeginQuote:    testPeginQuote,
			RetainedQuote: failRetainedQuote,
		}
		tickerChannel <- time.Now()
		assert.Eventually(t, checkFunction, time.Second, 10*time.Millisecond)
		btcRpc.AssertExpectations(t)
		bridge.AssertNotCalled(t, "GetRequiredTxConfirmations")
	})
}
