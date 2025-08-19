package watcher_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPeginDepositAddressWatcher_Shutdown(t *testing.T) {
	eventBus := &mocks.EventBusMock{}
	eventBus.On("Subscribe", mock.Anything).Return(make(<-chan entities.Event))
	createWatcherShutdownTest(t, func(ticker watcher.Ticker) watcher.Watcher {
		useCases := watcher.NewPeginDepositAddressWatcherUseCases(nil, nil, nil, nil)
		return watcher.NewPeginDepositAddressWatcher(useCases, nil, blockchain.Rpc{}, eventBus, ticker)
	})
}

func TestPeginDepositAddressWatcher_Prepare(t *testing.T) {
	t.Run("prepare successfully", func(t *testing.T) {
		retainedQuotes := []quote.RetainedPeginQuote{
			{QuoteHash: "q1", DepositAddress: "addr1"},
			{QuoteHash: "q2", DepositAddress: "addr2"},
			{QuoteHash: "q3", DepositAddress: "addr4"},
		}
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		peginRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, quote.PeginStateWaitingForDeposit).Return(retainedQuotes[0:1], nil).Once()
		peginRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, quote.PeginStateWaitingForDepositConfirmations).Return(retainedQuotes[1:], nil).Once()
		btcWallet := &mocks.BitcoinWalletMock{}
		for i, q := range retainedQuotes {
			peginRepository.EXPECT().GetQuote(mock.Anything, q.QuoteHash).Return(&quote.PeginQuote{Nonce: int64(i)}, nil).Once()
			peginRepository.EXPECT().GetPeginCreationData(mock.Anything, q.QuoteHash).Return(quote.PeginCreationData{GasPrice: entities.NewWei(int64(i))}).Once()
			btcWallet.EXPECT().ImportAddress(q.DepositAddress).Return(nil).Once()
		}
		getWatchedUseCase := w.NewGetWatchedPeginQuoteUseCase(peginRepository)
		useCases := watcher.NewPeginDepositAddressWatcherUseCases(nil, getWatchedUseCase, nil, nil)
		addressWatcher := watcher.NewPeginDepositAddressWatcher(useCases, btcWallet, blockchain.Rpc{}, nil, nil)
		err := addressWatcher.Prepare(context.Background())
		require.NoError(t, err)
		peginRepository.AssertExpectations(t)
		btcWallet.AssertExpectations(t)
		for _, q := range retainedQuotes {
			watchedQuote, ok := addressWatcher.GetWatchedQuote(q.QuoteHash)
			require.True(t, ok)
			require.Equal(t, q.QuoteHash, watchedQuote.RetainedQuote.QuoteHash)
		}
	})
	t.Run("handle error getting retained quotes", func(t *testing.T) {
		peginRepository := &mocks.PeginQuoteRepositoryMock{}
		peginRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		getWatchedUseCase := w.NewGetWatchedPeginQuoteUseCase(peginRepository)
		useCases := watcher.NewPeginDepositAddressWatcherUseCases(nil, getWatchedUseCase, nil, nil)
		addressWatcher := watcher.NewPeginDepositAddressWatcher(useCases, nil, blockchain.Rpc{}, nil, nil)
		err := addressWatcher.Prepare(context.Background())
		require.Error(t, err)
		peginRepository.AssertExpectations(t)
	})
}

// nolint:funlen
func TestPeginDepositAddressWatcher_Start_QuoteAccepted(t *testing.T) {
	testRetainedQuote := quote.RetainedPeginQuote{QuoteHash: "010203", DepositAddress: test.AnyAddress}
	testPeginQuote := quote.PeginQuote{Nonce: 5}
	btcWallet := &mocks.BtcWalletMock{}
	rpc := blockchain.Rpc{}
	eventBus := &mocks.EventBusMock{}
	acceptPeginChannel := make(chan entities.Event)
	eventBus.On("Subscribe", quote.AcceptedPeginQuoteEventId).Return((<-chan entities.Event)(acceptPeginChannel))
	ticker := &mocks.TickerMock{}
	ticker.EXPECT().C().Return(make(chan time.Time))
	ticker.EXPECT().Stop().Return()
	useCases := watcher.NewPeginDepositAddressWatcherUseCases(nil, nil, nil, nil)
	peginWatcher := watcher.NewPeginDepositAddressWatcher(useCases, btcWallet, rpc, eventBus, ticker)

	go peginWatcher.Start()

	t.Run("handle error importing address", func(t *testing.T) {
		checkFunction := test.AssertLogContains(t, "error while importing deposit address (any address)")
		btcWallet.On("ImportAddress", mock.Anything).Return(assert.AnError).Once()
		acceptPeginChannel <- quote.AcceptedPeginQuoteEvent{
			Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
			Quote:         testPeginQuote,
			RetainedQuote: testRetainedQuote,
		}
		watchedQuote, ok := peginWatcher.GetWatchedQuote(test.AnyString)
		assert.False(t, ok)
		assert.Empty(t, watchedQuote)
		assert.Eventually(t, func() bool { return btcWallet.AssertExpectations(t) && checkFunction() }, time.Second, 10*time.Millisecond)
	})
	t.Run("handle accepted pegin quote", func(t *testing.T) {
		defer test.AssertNoLog(t)
		btcWallet.Calls = []mock.Call{}
		btcWallet.ExpectedCalls = []*mock.Call{}
		btcWallet.On("ImportAddress", test.AnyAddress).Return(nil).Once()
		watchedQuote, ok := peginWatcher.GetWatchedQuote(test.AnyString)
		assert.False(t, ok)
		assert.Empty(t, watchedQuote)
		acceptPeginChannel <- quote.AcceptedPeginQuoteEvent{
			Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
			Quote:         testPeginQuote,
			RetainedQuote: testRetainedQuote,
		}
		assert.Eventually(t, func() bool {
			watchedQuote, ok = peginWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			assert.True(t, ok)
			return assert.Equal(t, quote.WatchedPeginQuote{
				PeginQuote:    testPeginQuote,
				RetainedQuote: testRetainedQuote,
			}, watchedQuote) && btcWallet.AssertExpectations(t)
		}, time.Second, 10*time.Millisecond)
	})
	t.Run("handle already watched quote", func(t *testing.T) {
		checkFunction := test.AssertLogContains(t, "Quote 010203 is already watched")
		acceptPeginChannel <- quote.AcceptedPeginQuoteEvent{
			Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
			Quote:         testPeginQuote,
			RetainedQuote: testRetainedQuote,
		}
		assert.Eventually(t, checkFunction, time.Second, 10*time.Millisecond)
	})
	t.Run("handle incorrect event sent to bus", func(t *testing.T) {
		checkFunction := test.AssertLogContains(t, "trying to parse wrong event")
		acceptPeginChannel <- quote.PegoutQuoteCompletedEvent{
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

// nolint:funlen,cyclop,maintidx,gocyclo
func TestPeginDepositAddressWatcher_Start_BlockchainCheck(t *testing.T) {
	peginRepository := &mocks.PeginQuoteRepositoryMock{}
	peginRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, quote.PeginStateWaitingForDeposit).Return([]quote.RetainedPeginQuote{}, nil).Once()
	peginRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, quote.PeginStateWaitingForDepositConfirmations).Return([]quote.RetainedPeginQuote{}, nil).Once()
	btcWallet := &mocks.BtcWalletMock{}
	btcRpc := &mocks.BtcRpcMock{}
	rpc := blockchain.Rpc{Btc: btcRpc}
	eventBus := &mocks.EventBusMock{}
	eventBus.On("Publish", mock.Anything).Return(nil).Twice()
	acceptPeginChannel := make(chan entities.Event)
	eventBus.On("Subscribe", quote.AcceptedPeginQuoteEventId).Return((<-chan entities.Event)(acceptPeginChannel))
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop().Return()
	lbc := &mocks.LiquidityBridgeContractMock{}
	lbc.On("GetBalance", mock.Anything).Return(entities.NewWei(1000), nil)
	lbc.On("CallForUser", mock.Anything, mock.Anything).Return(test.AnyHash, nil)
	bridge := &mocks.BridgeMock{}
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1), nil)
	peginProvider := &mocks.ProviderMock{}
	peginProvider.On("RskAddress").Return(test.AnyAddress)
	appMutexes := environment.NewApplicationMutexes()
	getUseCase := w.NewGetWatchedPeginQuoteUseCase(peginRepository)
	expireUseCase := pegin.NewExpiredPeginQuoteUseCase(peginRepository)
	updateUseCase := w.NewUpdatePeginDepositUseCase(peginRepository)
	cfuUseCase := pegin.NewCallForUserUseCase(blockchain.RskContracts{Lbc: lbc, Bridge: bridge}, peginRepository, rpc, peginProvider, eventBus, appMutexes.RskWalletMutex())
	useCases := watcher.NewPeginDepositAddressWatcherUseCases(cfuUseCase, getUseCase, updateUseCase, expireUseCase)
	peginWatcher := watcher.NewPeginDepositAddressWatcher(useCases, btcWallet, rpc, eventBus, ticker)

	resetMocks := func() {
		btcRpc.ExpectedCalls = []*mock.Call{}
		btcRpc.Calls = []mock.Call{}
		btcWallet.ExpectedCalls = []*mock.Call{}
		btcWallet.Calls = []mock.Call{}
		peginRepository.ExpectedCalls = []*mock.Call{}
		peginRepository.Calls = []mock.Call{}
	}

	prepareErr := peginWatcher.Prepare(context.Background())
	require.NoError(t, prepareErr)
	go peginWatcher.Start()

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
	t.Run("should handle error getting block height", func(t *testing.T) {
		resetMocks()
		checkFunction := test.AssertLogContains(t, assert.AnError.Error())
		btcRpc.On("GetHeight").Return(nil, assert.AnError).Once()
		tickerChannel <- time.Now()
		assert.Eventually(t, checkFunction, time.Second, 10*time.Millisecond)
		btcRpc.AssertExpectations(t)
	})
	t.Run("should call expire use case on expired quotes", func(t *testing.T) {
		expiredRetained := quote.RetainedPeginQuote{QuoteHash: test.AnyHash, DepositAddress: test.AnyAddress, State: quote.PeginStateWaitingForDeposit}
		expiredQuote := quote.PeginQuote{Nonce: 6, AgreementTimestamp: 1}
		t.Run("should handle error when expiring quotes", func(t *testing.T) {
			resetMocks()
			checkFunction := test.AssertLogContains(t, "Error updating expired quote (d8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb)")
			peginRepository.EXPECT().UpdateRetainedQuote(mock.Anything, mock.Anything).Return(assert.AnError).Once()
			btcRpc.On("GetHeight").Return(big.NewInt(9), nil).Once()
			btcWallet.On("ImportAddress", test.AnyAddress).Return(nil).Once()
			btcWallet.On("GetTransactions", test.AnyAddress).Return([]blockchain.BitcoinTransactionInformation{}, nil).Once()
			acceptPeginChannel <- quote.AcceptedPeginQuoteEvent{
				Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
				Quote:         expiredQuote,
				RetainedQuote: expiredRetained,
			}
			assert.Eventually(t, func() bool {
				watchedQuote, ok := peginWatcher.GetWatchedQuote(test.AnyHash)
				assert.True(t, ok)
				return assert.Equal(t, quote.WatchedPeginQuote{
					PeginQuote:    expiredQuote,
					RetainedQuote: expiredRetained,
				}, watchedQuote)
			}, time.Second, 10*time.Millisecond)

			tickerChannel <- time.Now()

			assert.Eventually(t, func() bool {
				watchedQuote, ok := peginWatcher.GetWatchedQuote(test.AnyHash)
				return checkFunction() && assert.NotEmpty(t, watchedQuote) && assert.True(t, ok) &&
					peginRepository.AssertExpectations(t) && btcRpc.AssertExpectations(t) && btcWallet.AssertExpectations(t)
			}, time.Second, 10*time.Millisecond)
		})
		t.Run("should stop tracking quotes after expiring them", func(t *testing.T) {
			resetMocks()
			peginRepository.EXPECT().UpdateRetainedQuote(mock.Anything, mock.Anything).Return(nil).Once()
			btcRpc.On("GetHeight").Return(big.NewInt(10), nil).Once()
			btcWallet.On("GetTransactions", test.AnyAddress).Return([]blockchain.BitcoinTransactionInformation{}, nil).Once()
			assert.Eventually(t, func() bool {
				watchedQuote, ok := peginWatcher.GetWatchedQuote(test.AnyHash)
				assert.True(t, ok)
				return assert.Equal(t, quote.WatchedPeginQuote{
					PeginQuote:    expiredQuote,
					RetainedQuote: expiredRetained,
				}, watchedQuote)
			}, time.Second, 10*time.Millisecond)

			tickerChannel <- time.Now()

			assert.Eventually(t, func() bool {
				watchedQuote, ok := peginWatcher.GetWatchedQuote(test.AnyHash)
				return assert.Empty(t, watchedQuote) && assert.False(t, ok) && peginRepository.AssertExpectations(t) && btcRpc.AssertExpectations(t)
			}, time.Second, 10*time.Millisecond)
		})
	})
	t.Run("should update state to WaitingForDepositConfirmations after first confirmation", func(t *testing.T) {
		btcWallet.On("ImportAddress", test.AnyAddress).Return(nil).Once()
		testRetainedQuote := quote.RetainedPeginQuote{QuoteHash: test.AnyHash, DepositAddress: test.AnyAddress, State: quote.PeginStateWaitingForDeposit}
		testQuote := quote.PeginQuote{Nonce: 8, AgreementTimestamp: uint32(time.Now().Unix()), TimeForDeposit: 6000, Confirmations: 10, Value: entities.NewWei(1)}
		acceptPeginChannel <- quote.AcceptedPeginQuoteEvent{
			Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
			Quote:         testQuote,
			RetainedQuote: testRetainedQuote,
		}
		assert.Eventually(t, func() bool { return btcWallet.AssertExpectations(t) }, time.Second, 10*time.Millisecond)
		t.Run("should handle error getting transactions from wallet", func(t *testing.T) {
			resetMocks()
			const errorMsg = "error getting tx"
			btcRpc.On("GetHeight").Return(big.NewInt(11), nil).Once()
			checkFunction := test.AssertLogContains(t, errorMsg)
			btcWallet.On("GetTransactions", testRetainedQuote.DepositAddress).Return(nil, errors.New(errorMsg)).Once()
			tickerChannel <- time.Now()
			assert.Eventually(t, func() bool {
				return checkFunction() && btcWallet.AssertExpectations(t)
			}, time.Second, 10*time.Millisecond)
		})
		t.Run("should handle error getting transaction block", func(t *testing.T) {
			resetMocks()
			const errorMsg = "error getting block"
			btcRpc.On("GetHeight").Return(big.NewInt(12), nil).Once()
			checkFunction := test.AssertLogContains(t, errorMsg)
			btcWallet.On("GetTransactions", testRetainedQuote.DepositAddress).Return([]blockchain.BitcoinTransactionInformation{
				{Hash: test.AnyHash, Confirmations: 5, Outputs: map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(100)}}},
			}, nil).Once()
			btcRpc.On("GetTransactionBlockInfo", test.AnyHash).Return(blockchain.BitcoinBlockInformation{}, errors.New(errorMsg)).Once()
			tickerChannel <- time.Now()
			assert.Eventually(t, func() bool {
				return checkFunction() && btcWallet.AssertExpectations(t) && btcRpc.AssertExpectations(t)
			}, time.Second, 10*time.Millisecond)
		})
		t.Run("should not update quote if doesn't meet the conditions", func(t *testing.T) {
			resetMocks()
			// incorrect amount
			btcRpc.On("GetHeight").Return(big.NewInt(13), nil).Once()
			btcWallet.On("GetTransactions", testRetainedQuote.DepositAddress).Return([]blockchain.BitcoinTransactionInformation{
				{Hash: test.AnyHash, Confirmations: 1, Outputs: map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(0)}}},
			}, nil).Once()
			btcRpc.On("GetTransactionBlockInfo", test.AnyHash).Return(blockchain.BitcoinBlockInformation{Time: time.Now()}, nil).Once()
			tickerChannel <- time.Now()

			// incorrect time
			btcRpc.On("GetHeight").Return(big.NewInt(14), nil).Once()
			btcWallet.On("GetTransactions", testRetainedQuote.DepositAddress).Return([]blockchain.BitcoinTransactionInformation{
				{Hash: test.AnyHash, Confirmations: 1, Outputs: map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(100)}}},
			}, nil).Once()
			btcRpc.On("GetTransactionBlockInfo", test.AnyHash).Return(blockchain.BitcoinBlockInformation{Time: time.Now().Add(8000 * time.Second)}, nil).Once()
			tickerChannel <- time.Now()

			assert.Eventually(t, func() bool {
				return btcWallet.AssertExpectations(t) && btcRpc.AssertExpectations(t) && peginRepository.AssertNotCalled(t, "UpdateRetainedQuote", mock.Anything, mock.Anything)
			}, time.Second, 10*time.Millisecond)
		})
		t.Run("should update quote successfully and continue tracking it", func(t *testing.T) {
			resetMocks()
			btcRpc.On("GetHeight").Return(big.NewInt(15), nil).Once()
			tx := blockchain.BitcoinTransactionInformation{Hash: test.AnyHash, Confirmations: 1, Outputs: map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(100)}}}
			btcWallet.On("GetTransactions", testRetainedQuote.DepositAddress).Return([]blockchain.BitcoinTransactionInformation{tx}, nil).Once()
			btcRpc.On("GetTransactionInfo", test.AnyHash).Return(tx, nil).Once()
			btcRpc.On("GetTransactionBlockInfo", test.AnyHash).Return(blockchain.BitcoinBlockInformation{Time: time.Now()}, nil).Once()
			peginRepository.EXPECT().UpdateRetainedQuote(mock.Anything, mock.Anything).Return(nil).Once()
			tickerChannel <- time.Now()
			assert.Eventually(t, func() bool {
				watchedQuote, _ := peginWatcher.GetWatchedQuote(test.AnyHash)
				return btcWallet.AssertExpectations(t) && btcRpc.AssertExpectations(t) &&
					peginRepository.AssertExpectations(t) && assert.Equal(t, quote.PeginStateWaitingForDepositConfirmations, watchedQuote.RetainedQuote.State)
			}, time.Second, 10*time.Millisecond)
		})
	})
	t.Run("should execute call for user use case after confirmations have passed", func(t *testing.T) {
		t.Run("should not execute call for user if confirmations aren't enough", func(t *testing.T) {
			resetMocks()
			btcRpc.On("GetHeight").Return(big.NewInt(16), nil).Once()
			btcRpc.On("GetTransactionInfo", test.AnyHash).Return(blockchain.BitcoinTransactionInformation{Confirmations: 5}, nil).Once()
			tickerChannel <- time.Now()
			assert.Eventually(t, func() bool {
				return btcRpc.AssertExpectations(t) && peginRepository.AssertNotCalled(t, "UpdateRetainedQuote", mock.Anything, mock.Anything)
			}, time.Second, 10*time.Millisecond)
		})
		t.Run("shouldn't stop tracking quote on recoverable error", func(t *testing.T) {
			resetMocks()
			btcRpc.On("GetHeight").Return(big.NewInt(17), nil).Once()
			btcRpc.On("GetTransactionInfo", test.AnyHash).Return(blockchain.BitcoinTransactionInformation{Confirmations: 10}, nil).Once()
			peginRepository.EXPECT().GetQuote(mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
			tickerChannel <- time.Now()
			assert.Eventually(t, func() bool {
				watchedQuote, ok := peginWatcher.GetWatchedQuote(test.AnyHash)
				return btcRpc.AssertExpectations(t) && peginRepository.AssertExpectations(t) && assert.True(t, ok) && assert.NotEmpty(t, watchedQuote)
			}, time.Second, 10*time.Millisecond)
		})
		t.Run("should stop tracking quote on non-recoverable error", func(t *testing.T) {
			resetMocks()
			btcRpc.On("GetHeight").Return(big.NewInt(18), nil).Once()
			btcRpc.On("GetTransactionInfo", test.AnyHash).Return(blockchain.BitcoinTransactionInformation{Confirmations: 10}, nil).Once()
			peginRepository.EXPECT().GetQuote(mock.Anything, mock.Anything).Return(nil, errors.Join(assert.AnError, usecases.NonRecoverableError)).Once()
			tickerChannel <- time.Now()
			assert.Eventually(t, func() bool {
				watchedQuote, ok := peginWatcher.GetWatchedQuote(test.AnyHash)
				return btcRpc.AssertExpectations(t) && peginRepository.AssertExpectations(t) && assert.False(t, ok) && assert.Empty(t, watchedQuote)
			}, time.Second, 10*time.Millisecond)
		})
		t.Run("should stop tracking quote on successful call for user", func(t *testing.T) {
			resetMocks()
			testRetainedQuote := quote.RetainedPeginQuote{QuoteHash: test.AnyHash, DepositAddress: test.AnyAddress, State: quote.PeginStateWaitingForDepositConfirmations}
			testQuote := quote.PeginQuote{Nonce: 8, AgreementTimestamp: uint32(time.Now().Unix()), TimeForDeposit: 6000, Confirmations: 10, Value: entities.NewWei(1)}

			btcWallet.On("ImportAddress", test.AnyAddress).Return(nil).Once()
			btcRpc.On("GetTransactionInfo", mock.Anything).Return(blockchain.BitcoinTransactionInformation{
				Hash: test.AnyHash, Confirmations: 15, Outputs: map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(100)}},
			}, nil).Twice()
			btcRpc.On("GetHeight").Return(big.NewInt(19), nil).Once()
			btcRpc.On("GetTransactionBlockInfo", mock.Anything).Return(blockchain.BitcoinBlockInformation{Time: time.Now()}, nil).Once()
			peginRepository.EXPECT().GetQuote(mock.Anything, mock.Anything).Return(&testQuote, nil).Once()
			peginRepository.EXPECT().UpdateRetainedQuote(mock.Anything, mock.Anything).Return(nil).Once()
			peginRepository.EXPECT().GetPeginCreationData(mock.Anything, mock.Anything).Return(quote.PeginCreationData{GasPrice: entities.NewWei(1)}).Once()

			acceptPeginChannel <- quote.AcceptedPeginQuoteEvent{
				Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
				Quote:         testQuote,
				RetainedQuote: testRetainedQuote,
			}
			tickerChannel <- time.Now()
			assert.Eventually(t, func() bool {
				watchedQuote, ok := peginWatcher.GetWatchedQuote(test.AnyHash)
				return btcWallet.AssertExpectations(t) && btcRpc.AssertExpectations(t) && peginRepository.AssertExpectations(t) &&
					assert.False(t, ok) && assert.Empty(t, watchedQuote)
			}, time.Second, 10*time.Millisecond)
		})
	})
	t.Run("should update expired quote if block was mined before expiration", func(t *testing.T) {
		const (
			otherHash = "quote-hash-1"
			txHash    = "tx-hash-1"
		)
		resetMocks()
		now := time.Now().Unix()
		btcWallet.On("ImportAddress", test.AnyAddress).Return(nil).Once()
		testRetainedQuote := quote.RetainedPeginQuote{QuoteHash: otherHash, DepositAddress: test.AnyAddress, State: quote.PeginStateWaitingForDeposit}
		testQuote := quote.PeginQuote{Nonce: 123, AgreementTimestamp: uint32(now - 6000), TimeForDeposit: 5000, Confirmations: 10, Value: entities.NewWei(1)}
		acceptPeginChannel <- quote.AcceptedPeginQuoteEvent{
			Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
			Quote:         testQuote,
			RetainedQuote: testRetainedQuote,
		}
		assert.Eventually(t, func() bool { return btcWallet.AssertExpectations(t) }, time.Second, 10*time.Millisecond)
		btcRpc.On("GetHeight").Return(big.NewInt(20), nil).Once()
		tx := blockchain.BitcoinTransactionInformation{Hash: txHash, Confirmations: 1, Outputs: map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(1)}}}
		btcWallet.On("GetTransactions", testRetainedQuote.DepositAddress).Return([]blockchain.BitcoinTransactionInformation{tx}, nil).Once()
		btcRpc.On("GetTransactionBlockInfo", txHash).Return(blockchain.BitcoinBlockInformation{Time: time.Unix(now-2000, 0)}, nil).Once()
		btcRpc.On("GetTransactionInfo", txHash).Return(tx, nil).Once()
		peginRepository.EXPECT().UpdateRetainedQuote(mock.Anything, mock.Anything).Return(nil).Once()
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			watchedQuote, _ := peginWatcher.GetWatchedQuote(otherHash)
			return btcWallet.AssertExpectations(t) && btcRpc.AssertExpectations(t) &&
				peginRepository.AssertExpectations(t) && assert.Equal(t, quote.PeginStateWaitingForDepositConfirmations, watchedQuote.RetainedQuote.State)
		}, time.Second, 10*time.Millisecond)
	})
	t.Run("should update expired quote and call for user if block was mined before expiration and confirmations already passed", func(t *testing.T) {
		const (
			otherHash = "quote-hash-2"
			txHash    = "tx-hash-2"
		)
		resetMocks()
		now := time.Now().Unix()
		btcWallet.On("ImportAddress", test.AnyAddress).Return(nil).Once()
		testRetainedQuote := quote.RetainedPeginQuote{QuoteHash: otherHash, DepositAddress: test.AnyAddress, State: quote.PeginStateWaitingForDeposit}
		testQuote := quote.PeginQuote{Nonce: 123, AgreementTimestamp: uint32(now - 6000), TimeForDeposit: 5000, Confirmations: 10, Value: entities.NewWei(1)}
		acceptPeginChannel <- quote.AcceptedPeginQuoteEvent{
			Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
			Quote:         testQuote,
			RetainedQuote: testRetainedQuote,
		}
		assert.Eventually(t, func() bool { return btcWallet.AssertExpectations(t) }, time.Second, 10*time.Millisecond)
		btcRpc.On("GetHeight").Return(big.NewInt(21), nil).Once()
		confirmedTx := blockchain.BitcoinTransactionInformation{Hash: txHash, Confirmations: 100, Outputs: map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(1)}}}
		unconfirmedTx := confirmedTx
		unconfirmedTx.Confirmations = 1
		btcWallet.On("GetTransactions", testRetainedQuote.DepositAddress).Return([]blockchain.BitcoinTransactionInformation{confirmedTx}, nil).Once()
		btcRpc.On("GetTransactionBlockInfo", txHash).Return(blockchain.BitcoinBlockInformation{Time: time.Unix(now-2000, 0)}, nil).Twice()
		btcRpc.On("GetTransactionInfo", txHash).Return(confirmedTx, nil).Twice()
		btcRpc.On("GetTransactionInfo", mock.Anything).Return(unconfirmedTx, nil).Once() // the quote from the old test still in the watcher
		peginRepository.EXPECT().GetQuote(mock.Anything, mock.Anything).Return(&testQuote, nil).Once()
		peginRepository.EXPECT().UpdateRetainedQuote(mock.Anything, mock.Anything).Return(nil).Twice()
		peginRepository.EXPECT().GetPeginCreationData(mock.Anything, mock.Anything).Return(quote.PeginCreationData{GasPrice: entities.NewWei(1)}).Once()
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			watchedQuote, ok := peginWatcher.GetWatchedQuote(test.AnyHash)
			return btcWallet.AssertExpectations(t) && btcRpc.AssertExpectations(t) && peginRepository.AssertExpectations(t) &&
				assert.False(t, ok) && assert.Empty(t, watchedQuote)
		}, time.Second, 10*time.Millisecond)
	})
}
