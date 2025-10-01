package watcher_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
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

func TestNewPegoutRskDepositWatcher(t *testing.T) {
	ticker := &mocks.TickerMock{}
	providerMock := &mocks.ProviderMock{}
	contracts := blockchain.RskContracts{Lbc: &mocks.LiquidityBridgeContractMock{}}
	rpc := blockchain.Rpc{Btc: &mocks.BtcRpcMock{}, Rsk: &mocks.RootstockRpcServerMock{}}
	eventBus := &mocks.EventBusMock{}
	useCases := watcher.NewPegoutRskDepositWatcherUseCases(
		&w.GetWatchedPegoutQuoteUseCase{},
		&pegout.ExpiredPegoutQuoteUseCase{},
		&pegout.SendPegoutUseCase{},
		&w.UpdatePegoutQuoteDepositUseCase{},
		&pegout.InitPegoutDepositCacheUseCase{},
	)
	depositWatcher := watcher.NewPegoutRskDepositWatcher(useCases, providerMock, rpc, contracts, eventBus, 1, ticker, time.Duration(1))
	// the strcut has 17, but we need the mutexes to have the zero value
	assert.Equal(t, 15, test.CountNonZeroValues(depositWatcher))
}

// nolint:funlen
func TestPegoutRskDepositWatcher_Prepare(t *testing.T) {
	t.Run("should handle error during cache initialization", func(t *testing.T) {
		contracts := blockchain.RskContracts{}
		rskRpc := &mocks.RootstockRpcServerMock{}
		rpc := blockchain.Rpc{Rsk: rskRpc}
		initCacheUseCase := pegout.NewInitPegoutDepositCacheUseCase(&mocks.PegoutQuoteRepositoryMock{}, contracts, rpc)
		useCases := watcher.NewPegoutRskDepositWatcherUseCases(nil, nil, nil, nil, initCacheUseCase)
		depositWatcher := watcher.NewPegoutRskDepositWatcher(useCases, nil, rpc, contracts, nil, 1, nil, time.Duration(1))
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(0), assert.AnError)
		err := depositWatcher.Prepare(context.Background())
		require.Error(t, err)
		rskRpc.AssertExpectations(t)
	})
	t.Run("should initialize quote cache", func(t *testing.T) {
		testRetainedQuotes := []quote.RetainedPegoutQuote{
			{QuoteHash: "0102", State: quote.PegoutStateWaitingForDeposit},
			{QuoteHash: "0203", State: quote.PegoutStateWaitingForDepositConfirmations},
			{QuoteHash: "0304", State: quote.PegoutStateWaitingForDeposit},
		}
		lbc := &mocks.LiquidityBridgeContractMock{}
		lbc.On("GetDepositEvents", mock.Anything, mock.Anything, mock.Anything).Return([]quote.PegoutDeposit{}, nil)
		providerMock := &mocks.ProviderMock{}
		providerMock.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.DefaultPegoutConfiguration()).Once()
		contracts := blockchain.RskContracts{Lbc: lbc}
		rskRpc := &mocks.RootstockRpcServerMock{}
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(8000), nil).Once()
		rpc := blockchain.Rpc{Rsk: rskRpc}
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		pegoutRepository.EXPECT().UpsertPegoutDeposits(mock.Anything, mock.Anything).Return(nil).Once()
		pegoutRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, quote.PegoutStateWaitingForDeposit).Return([]quote.RetainedPegoutQuote{testRetainedQuotes[0], testRetainedQuotes[2]}, nil).Once()
		pegoutRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, quote.PegoutStateWaitingForDepositConfirmations).Return([]quote.RetainedPegoutQuote{testRetainedQuotes[1]}, nil).Once()
		for i, q := range testRetainedQuotes {
			pegoutRepository.EXPECT().GetQuote(mock.Anything, q.QuoteHash).Return(&quote.PegoutQuote{Nonce: int64(i + 1), ExpireBlock: uint32((i + 1) * 1000)}, nil).Once()
			pegoutRepository.EXPECT().GetPegoutCreationData(mock.Anything, q.QuoteHash).Return(quote.PegoutCreationData{GasPrice: entities.NewWei(int64(i))}).Once()
		}

		initCacheUseCase := pegout.NewInitPegoutDepositCacheUseCase(pegoutRepository, contracts, rpc)
		getWatchedQuotesUseCase := w.NewGetWatchedPegoutQuoteUseCase(pegoutRepository)
		useCases := watcher.NewPegoutRskDepositWatcherUseCases(getWatchedQuotesUseCase, nil, nil, nil, initCacheUseCase)
		depositWatcher := watcher.NewPegoutRskDepositWatcher(useCases, providerMock, rpc, contracts, nil, 3000, nil, time.Duration(1))
		err := depositWatcher.Prepare(context.Background())
		require.NoError(t, err)
		t.Run("should initialize cache successfully", func(t *testing.T) {
			for _, q := range testRetainedQuotes {
				watchedQuote, ok := depositWatcher.GetWatchedQuote(q.QuoteHash)
				assert.True(t, ok)
				assert.NotEmpty(t, watchedQuote)
			}
			providerMock.AssertExpectations(t)
			rskRpc.AssertExpectations(t)
			lbc.AssertExpectations(t)
			pegoutRepository.AssertExpectations(t)
		})
		t.Run("current block should be the oldest of the cache", func(t *testing.T) {
			assert.Equal(t, uint64(500), depositWatcher.GetCurrentBlock())
		})
	})

	t.Run("should start from the current block if cacheStartBlock is not provided", func(t *testing.T) {
		latestBlock := uint64(567)
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		lbc := &mocks.LiquidityBridgeContractMock{}
		contracts := blockchain.RskContracts{Lbc: lbc}
		rskRpc := &mocks.RootstockRpcServerMock{}
		rpc := blockchain.Rpc{Rsk: rskRpc}
		providerMock := &mocks.ProviderMock{}
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(latestBlock, nil).Once()
		pegoutRepository.EXPECT().
			GetRetainedQuoteByState(mock.Anything, quote.PegoutStateWaitingForDeposit).
			Return([]quote.RetainedPegoutQuote{}, nil).
			Once()
		pegoutRepository.EXPECT().
			GetRetainedQuoteByState(mock.Anything, quote.PegoutStateWaitingForDepositConfirmations).
			Return([]quote.RetainedPegoutQuote{}, nil).
			Once()
		providerMock.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.DefaultPegoutConfiguration()).Once()
		initCacheUseCase := pegout.NewInitPegoutDepositCacheUseCase(pegoutRepository, contracts, rpc)
		getWatchedQuotesUseCase := w.NewGetWatchedPegoutQuoteUseCase(pegoutRepository)
		useCases := watcher.NewPegoutRskDepositWatcherUseCases(getWatchedQuotesUseCase, nil, nil, nil, initCacheUseCase)
		depositWatcher := watcher.NewPegoutRskDepositWatcher(useCases, providerMock, rpc, contracts, nil, 0, nil, time.Duration(1))
		err := depositWatcher.Prepare(context.Background())
		require.NoError(t, err)
		assert.Equal(t, latestBlock, depositWatcher.GetCurrentBlock())
		rskRpc.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
		pegoutRepository.AssertNotCalled(t, "UpsertPegoutDeposits")
		lbc.AssertNotCalled(t, "GetDepositEvents")
	})

	t.Run("should start from the oldest quote if cacheStartBlock is not provided but some quotes are not processed yet", func(t *testing.T) {
		const (
			quoteHash1 = "0a0b"
			quoteHash2 = "0c0d"
			quoteHash3 = "0e0f"
		)
		testRetainedQuotes := []quote.RetainedPegoutQuote{
			{QuoteHash: quoteHash1, State: quote.PegoutStateWaitingForDeposit},
			{QuoteHash: quoteHash2, State: quote.PegoutStateWaitingForDepositConfirmations},
			{QuoteHash: quoteHash3, State: quote.PegoutStateWaitingForDeposit},
		}
		testQuotes := []quote.PegoutQuote{
			{Nonce: 1, ExpireBlock: 3123},
			{Nonce: 2, ExpireBlock: 1234},
			{Nonce: 3, ExpireBlock: 6241},
		}
		quotesCreationData := []quote.PegoutCreationData{
			{GasPrice: entities.NewWei(1)},
			{GasPrice: entities.NewWei(2)},
			{GasPrice: entities.NewWei(3)},
		}
		latestBlock := uint64(7000)
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		lbc := &mocks.LiquidityBridgeContractMock{}
		contracts := blockchain.RskContracts{Lbc: lbc}
		rskRpc := &mocks.RootstockRpcServerMock{}
		rpc := blockchain.Rpc{Rsk: rskRpc}
		providerMock := &mocks.ProviderMock{}
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(latestBlock, nil).Once()
		pegoutRepository.EXPECT().
			GetRetainedQuoteByState(mock.Anything, quote.PegoutStateWaitingForDeposit).
			Return([]quote.RetainedPegoutQuote{testRetainedQuotes[0], testRetainedQuotes[2]}, nil).
			Once()
		pegoutRepository.EXPECT().
			GetRetainedQuoteByState(mock.Anything, quote.PegoutStateWaitingForDepositConfirmations).
			Return([]quote.RetainedPegoutQuote{testRetainedQuotes[1]}, nil).
			Once()
		pegoutRepository.EXPECT().GetQuote(mock.Anything, quoteHash1).Return(&testQuotes[0], nil).Once()
		pegoutRepository.EXPECT().GetQuote(mock.Anything, quoteHash2).Return(&testQuotes[1], nil).Once()
		pegoutRepository.EXPECT().GetQuote(mock.Anything, quoteHash3).Return(&testQuotes[2], nil).Once()
		pegoutRepository.EXPECT().GetPegoutCreationData(mock.Anything, quoteHash1).Return(quotesCreationData[0]).Once()
		pegoutRepository.EXPECT().GetPegoutCreationData(mock.Anything, quoteHash2).Return(quotesCreationData[1]).Once()
		pegoutRepository.EXPECT().GetPegoutCreationData(mock.Anything, quoteHash3).Return(quotesCreationData[2]).Once()

		providerMock.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.DefaultPegoutConfiguration()).Once()
		initCacheUseCase := pegout.NewInitPegoutDepositCacheUseCase(pegoutRepository, contracts, rpc)
		getWatchedQuotesUseCase := w.NewGetWatchedPegoutQuoteUseCase(pegoutRepository)
		useCases := watcher.NewPegoutRskDepositWatcherUseCases(getWatchedQuotesUseCase, nil, nil, nil, initCacheUseCase)
		depositWatcher := watcher.NewPegoutRskDepositWatcher(useCases, providerMock, rpc, contracts, nil, 0, nil, time.Duration(1))
		err := depositWatcher.Prepare(context.Background())
		require.NoError(t, err)
		assert.Equal(t, uint64(734), depositWatcher.GetCurrentBlock())
		rskRpc.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
		pegoutRepository.AssertNotCalled(t, "UpsertPegoutDeposits")
		lbc.AssertNotCalled(t, "GetDepositEvents")
	})

	t.Run("should handle error getting height if cacheStartBlock is not provided", func(t *testing.T) {
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		lbc := &mocks.LiquidityBridgeContractMock{}
		contracts := blockchain.RskContracts{Lbc: lbc}
		rskRpc := &mocks.RootstockRpcServerMock{}
		rpc := blockchain.Rpc{Rsk: rskRpc}
		providerMock := &mocks.ProviderMock{}
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(0), assert.AnError).Once()
		initCacheUseCase := pegout.NewInitPegoutDepositCacheUseCase(pegoutRepository, contracts, rpc)
		getWatchedQuotesUseCase := w.NewGetWatchedPegoutQuoteUseCase(pegoutRepository)
		useCases := watcher.NewPegoutRskDepositWatcherUseCases(getWatchedQuotesUseCase, nil, nil, nil, initCacheUseCase)
		depositWatcher := watcher.NewPegoutRskDepositWatcher(useCases, providerMock, rpc, contracts, nil, 0, nil, time.Duration(1))
		err := depositWatcher.Prepare(context.Background())
		require.Error(t, err)
		rskRpc.AssertExpectations(t)
		providerMock.AssertNotCalled(t, "PegoutConfiguration")
		pegoutRepository.AssertNotCalled(t, "GetQuote")
		pegoutRepository.AssertNotCalled(t, "GetRetainedQuoteByState")
		pegoutRepository.AssertNotCalled(t, "UpsertPegoutDeposits")
		lbc.AssertNotCalled(t, "GetDepositEvents")
	})
}

func TestPegoutRskDepositWatcher_Shutdown(t *testing.T) {
	eventBus := &mocks.EventBusMock{}
	eventBus.On("Subscribe", mock.Anything).Return(make(<-chan entities.Event))
	createWatcherShutdownTest(t, func(ticker utils.Ticker) watcher.Watcher {
		return watcher.NewPegoutRskDepositWatcher(&watcher.PegoutRskDepositWatcherUseCases{}, nil, blockchain.Rpc{}, blockchain.RskContracts{}, eventBus, 0, ticker, time.Duration(1))
	})
}

func TestPegoutRskDepositWatcher_Start_QuoteAccepted(t *testing.T) {
	ticker := &mocks.TickerMock{}
	ticker.EXPECT().C().Return(make(chan time.Time))
	ticker.EXPECT().Stop().Return()
	lbc := &mocks.LiquidityBridgeContractMock{}
	providerMock := &mocks.ProviderMock{}
	contracts := blockchain.RskContracts{Lbc: lbc}
	rskRpc := &mocks.RootstockRpcServerMock{}
	rpc := blockchain.Rpc{Rsk: rskRpc}
	eventBus := &mocks.EventBusMock{}
	acceptPegoutChannel := make(chan entities.Event)
	eventBus.On("Subscribe", quote.AcceptedPegoutQuoteEventId).Return((<-chan entities.Event)(acceptPegoutChannel))

	testPegoutQuote := quote.PegoutQuote{Nonce: 1}
	testRetainedQuote := quote.RetainedPegoutQuote{QuoteHash: "010203"}

	useCases := watcher.NewPegoutRskDepositWatcherUseCases(nil, nil, nil, nil, nil)
	depositWatcher := watcher.NewPegoutRskDepositWatcher(useCases, providerMock, rpc, contracts, eventBus, 3000, ticker, time.Duration(1))

	go depositWatcher.Start()

	t.Run("handle accepted pegin quote", func(t *testing.T) {
		defer test.AssertNoLog(t)
		watchedQuote, ok := depositWatcher.GetWatchedQuote(test.AnyString)
		assert.False(t, ok)
		assert.Empty(t, watchedQuote)
		acceptPegoutChannel <- quote.AcceptedPegoutQuoteEvent{
			Event: entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
			Quote: testPegoutQuote, RetainedQuote: testRetainedQuote,
		}
		assert.Eventually(t, func() bool {
			watchedQuote, ok = depositWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			assert.True(t, ok)
			return assert.Equal(t, quote.WatchedPegoutQuote{
				PegoutQuote:   testPegoutQuote,
				RetainedQuote: testRetainedQuote,
			}, watchedQuote)
		}, time.Second, 10*time.Millisecond)
	})
	t.Run("handle already watched quote", func(t *testing.T) {
		checkFunction := test.AssertLogContains(t, "Quote 010203 is already watched")
		acceptPegoutChannel <- quote.AcceptedPegoutQuoteEvent{
			Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
			Quote:         testPegoutQuote,
			RetainedQuote: testRetainedQuote,
		}
		assert.Eventually(t, checkFunction, time.Second, 10*time.Millisecond)
	})
	t.Run("handle incorrect event sent to bus", func(t *testing.T) {
		checkFunction := test.AssertLogContains(t, "Trying to parse wrong event in Pegout Rsk deposit watcher")
		acceptPegoutChannel <- quote.AcceptedPeginQuoteEvent{Event: entities.NewBaseEvent(quote.PegoutQuoteCompletedEventId)}
		assert.Eventually(t, checkFunction, time.Second, 10*time.Millisecond)
	})

	closeChannel := make(chan bool)
	go depositWatcher.Shutdown(closeChannel)
	<-closeChannel
	assert.Eventually(t, func() bool { return eventBus.AssertExpectations(t) && ticker.AssertExpectations(t) }, time.Second, 10*time.Millisecond)
}

// nolint:funlen
func TestPegoutRskDepositWatcher_Start_BlockchainCheck_CheckDeposits(t *testing.T) {
	ticker := &mocks.TickerMock{}
	btcWallet := &mocks.BtcWalletMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop().Return()
	lbc := &mocks.LiquidityBridgeContractMock{}
	providerMock := &mocks.ProviderMock{}
	contracts := blockchain.RskContracts{Lbc: lbc}
	rskRpc := &mocks.RootstockRpcServerMock{}
	rpc := blockchain.Rpc{Rsk: rskRpc}
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	eventBus := &mocks.EventBusMock{}
	acceptPegoutChannel := make(chan entities.Event)
	eventBus.On("Subscribe", quote.AcceptedPegoutQuoteEventId).Return((<-chan entities.Event)(acceptPegoutChannel))

	testPegoutQuote := quote.PegoutQuote{Nonce: 1, Value: entities.NewWei(3), ExpireBlock: 100, ExpireDate: uint32(time.Now().Unix() + 600)}
	testRetainedQuote := quote.RetainedPegoutQuote{QuoteHash: "010203", State: quote.PegoutStateWaitingForDeposit}

	updatePegoutDeposit := w.NewUpdatePegoutQuoteDepositUseCase(pegoutRepository)
	useCases := watcher.NewPegoutRskDepositWatcherUseCases(nil, nil, nil, updatePegoutDeposit, nil)
	depositWatcher := watcher.NewPegoutRskDepositWatcher(useCases, providerMock, rpc, contracts, eventBus, 0, ticker, time.Duration(1))

	go depositWatcher.Start()
	t.Run("should handle error getting deposits", func(t *testing.T) {
		checkFunction := test.AssertLogContains(t, "error executing getting deposits in range [0, 5]")
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(5), nil).Once()
		lbc.On("GetDepositEvents", mock.Anything, uint64(0), mock.MatchedBy(matchUinPtr(5))).Return(nil, assert.AnError)
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool { return rskRpc.AssertExpectations(t) && lbc.AssertExpectations(t) && checkFunction() }, time.Second, 10*time.Millisecond)
	})
	t.Run("shouldn't update quote if deposit is not valid", func(t *testing.T) {
		rskRpc.Calls = []mock.Call{}
		rskRpc.ExpectedCalls = []*mock.Call{}
		lbc.Calls = []mock.Call{}
		lbc.ExpectedCalls = []*mock.Call{}
		acceptPegoutChannel <- quote.AcceptedPegoutQuoteEvent{
			Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
			Quote:         testPegoutQuote,
			RetainedQuote: testRetainedQuote,
		}

		// incorrect amount
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(6), nil).Once()
		lbc.On("GetDepositEvents", mock.Anything, uint64(5), mock.MatchedBy(matchUinPtr(6))).Return([]quote.PegoutDeposit{{
			TxHash:      test.AnyHash,
			QuoteHash:   testRetainedQuote.QuoteHash,
			Amount:      entities.NewWei(1),
			Timestamp:   time.Now(),
			BlockNumber: 6,
		}}, nil).Once()
		tickerChannel <- time.Now()

		// expired in time
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(7), nil).Once()
		lbc.On("GetDepositEvents", mock.Anything, uint64(6), mock.MatchedBy(matchUinPtr(7))).Return([]quote.PegoutDeposit{{
			TxHash:      test.AnyHash,
			QuoteHash:   testRetainedQuote.QuoteHash,
			Amount:      entities.NewWei(10),
			Timestamp:   time.Now().Add(time.Second * 1000),
			BlockNumber: 6,
		}}, nil).Once()
		tickerChannel <- time.Now()

		// expired in blocks
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(8), nil).Once()
		lbc.On("GetDepositEvents", mock.Anything, uint64(7), mock.MatchedBy(matchUinPtr(8))).Return([]quote.PegoutDeposit{{
			TxHash:      test.AnyHash,
			QuoteHash:   testRetainedQuote.QuoteHash,
			Amount:      entities.NewWei(10),
			Timestamp:   time.Now(),
			BlockNumber: 500,
		}}, nil).Once()
		tickerChannel <- time.Now()

		assert.Eventually(t, func() bool {
			watchedQuote, _ := depositWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			return lbc.AssertExpectations(t) && rskRpc.AssertExpectations(t) && assert.Equal(t, quote.PegoutStateWaitingForDeposit, watchedQuote.RetainedQuote.State)
		}, time.Second, 10*time.Millisecond)
	})
	t.Run("should update state to WaitingForDepositConfirmations after checking a valid deposit", func(t *testing.T) {
		rskRpc.Calls = []mock.Call{}
		rskRpc.ExpectedCalls = []*mock.Call{}
		lbc.Calls = []mock.Call{}
		lbc.ExpectedCalls = []*mock.Call{}
		acceptPegoutChannel <- quote.AcceptedPegoutQuoteEvent{
			Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
			Quote:         testPegoutQuote,
			RetainedQuote: testRetainedQuote,
		}

		validDeposit := quote.PegoutDeposit{
			TxHash:      test.AnyHash,
			QuoteHash:   testRetainedQuote.QuoteHash,
			Amount:      entities.NewWei(10),
			Timestamp:   time.Now(),
			BlockNumber: 6,
		}

		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(9), nil).Once()
		lbc.On("GetDepositEvents", mock.Anything, uint64(8), mock.MatchedBy(matchUinPtr(9))).Return([]quote.PegoutDeposit{validDeposit}, nil).Once()
		updatedRetained := testRetainedQuote
		updatedRetained.UserRskTxHash = validDeposit.TxHash
		updatedRetained.State = quote.PegoutStateWaitingForDepositConfirmations
		pegoutRepository.EXPECT().UpdateRetainedQuote(mock.Anything, updatedRetained).Return(nil).Once()
		pegoutRepository.EXPECT().UpsertPegoutDeposit(mock.Anything, validDeposit).Return(nil).Once()
		// not mature enough yet
		rskRpc.EXPECT().GetTransactionReceipt(mock.Anything, validDeposit.TxHash).Return(blockchain.TransactionReceipt{BlockNumber: 10}, nil).Once()
		tickerChannel <- time.Now()

		assert.Eventually(t, func() bool {
			watchedQuote, _ := depositWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			return lbc.AssertExpectations(t) && rskRpc.AssertExpectations(t) &&
				pegoutRepository.AssertExpectations(t) && btcWallet.AssertNotCalled(t, "SendWithOpReturn") &&
				assert.Equal(t, quote.PegoutStateWaitingForDepositConfirmations, watchedQuote.RetainedQuote.State)
		}, time.Second, 10*time.Millisecond)
	})
}

// nolint:funlen,cyclop
func TestPegoutRskDepositWatcher_Start_BlockchainCheck_CheckQuotes(t *testing.T) {
	mutexes := environment.NewApplicationMutexes()
	ticker := &mocks.TickerMock{}
	btcWallet := &mocks.BtcWalletMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop().Return()
	lbc := &mocks.LiquidityBridgeContractMock{}
	providerMock := &mocks.ProviderMock{}
	contracts := blockchain.RskContracts{Lbc: lbc}
	rskRpc := &mocks.RootstockRpcServerMock{}
	rpc := blockchain.Rpc{Rsk: rskRpc}
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	eventBus := &mocks.EventBusMock{}
	acceptPegoutChannel := make(chan entities.Event)
	eventBus.On("Subscribe", quote.AcceptedPegoutQuoteEventId).Return((<-chan entities.Event)(acceptPegoutChannel))
	eventBus.On("Publish", mock.Anything).Return(make(<-chan entities.Event))

	testPegoutQuote := quote.PegoutQuote{Nonce: 1, Value: entities.NewWei(3), ExpireBlock: 100, ExpireDate: uint32(time.Now().Unix() + 600), DepositConfirmations: 5}
	testRetainedQuote := quote.RetainedPegoutQuote{QuoteHash: "0102030000000000000000000000000000000000000000000000000000000000", State: quote.PegoutStateWaitingForDepositConfirmations, UserRskTxHash: test.AnyHash}

	expireUseCase := pegout.NewExpiredPegoutQuoteUseCase(pegoutRepository)
	sendPegoutUseCase := pegout.NewSendPegoutUseCase(btcWallet, pegoutRepository, rpc, eventBus, contracts, mutexes.BtcWalletMutex(), rootstock.ParseDepositEvent)
	useCases := watcher.NewPegoutRskDepositWatcherUseCases(nil, expireUseCase, sendPegoutUseCase, nil, nil)
	depositWatcher := watcher.NewPegoutRskDepositWatcher(useCases, providerMock, rpc, contracts, eventBus, 0, ticker, time.Duration(1))

	go depositWatcher.Start()
	t.Run("should stop tracking after cleaning expired quote", func(t *testing.T) {
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(10), nil).Once()
		lbc.On("GetDepositEvents", mock.Anything, uint64(0), mock.MatchedBy(matchUinPtr(10))).Return([]quote.PegoutDeposit{}, nil).Once()
		expired := testPegoutQuote
		expired.ExpireDate = uint32(time.Now().Unix() - 600)
		expiredRetained := testRetainedQuote
		expiredRetained.State = quote.PegoutStateWaitingForDeposit
		acceptPegoutChannel <- quote.AcceptedPegoutQuoteEvent{
			Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
			Quote:         expired,
			RetainedQuote: expiredRetained,
		}
		assert.Eventually(t, func() bool {
			q, ok := depositWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			return assert.True(t, ok) && assert.NotEmpty(t, q)
		}, time.Second, 10*time.Millisecond)
		pegoutRepository.EXPECT().UpdateRetainedQuote(mock.Anything, mock.Anything).Return(nil).Once()
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			q, ok := depositWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			return assert.False(t, ok) && assert.Empty(t, q) && pegoutRepository.AssertExpectations(t)
		}, time.Second, 10*time.Millisecond)
	})
	acceptPegoutChannel <- quote.AcceptedPegoutQuoteEvent{
		Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
		Quote:         testPegoutQuote,
		RetainedQuote: testRetainedQuote,
	}
	t.Run("shouldn't stop tracking on recoverable error when sending pegout", func(t *testing.T) {
		rskRpc.Calls = []mock.Call{}
		rskRpc.ExpectedCalls = []*mock.Call{}
		lbc.Calls = []mock.Call{}
		lbc.ExpectedCalls = []*mock.Call{}
		pegoutRepository.Calls = []mock.Call{}
		pegoutRepository.ExpectedCalls = []*mock.Call{}

		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(20), nil).Once()
		receipt := &blockchain.TransactionReceipt{
			BlockNumber: 10,
			Value:       entities.NewWei(3),
		}
		receipt = test.AddDepositLogFromQuote(t, receipt, testPegoutQuote, testRetainedQuote)
		rskRpc.EXPECT().GetTransactionReceipt(mock.Anything, testRetainedQuote.UserRskTxHash).
			Return(*receipt, nil).Once()
		lbc.On("GetDepositEvents", mock.Anything, uint64(10), mock.MatchedBy(matchUinPtr(20))).Return([]quote.PegoutDeposit{}, nil).Once()
		pegoutRepository.EXPECT().GetQuote(mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		assert.Eventually(t, func() bool {
			q, ok := depositWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			return assert.True(t, ok) && assert.NotEmpty(t, q)
		}, time.Second, 10*time.Millisecond)

		tickerChannel <- time.Now()

		assert.Eventually(t, func() bool {
			q, ok := depositWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			return assert.True(t, ok) && assert.NotEmpty(t, q) && pegoutRepository.AssertExpectations(t) &&
				lbc.AssertExpectations(t) && rskRpc.AssertExpectations(t)
		}, time.Second, 10*time.Millisecond)
	})
	t.Run("should stop tracking on non-recoverable error when sending pegout", func(t *testing.T) {
		rskRpc.Calls = []mock.Call{}
		rskRpc.ExpectedCalls = []*mock.Call{}
		lbc.Calls = []mock.Call{}
		lbc.ExpectedCalls = []*mock.Call{}
		pegoutRepository.Calls = []mock.Call{}
		pegoutRepository.ExpectedCalls = []*mock.Call{}

		receipt := &blockchain.TransactionReceipt{
			BlockNumber: 10,
			Value:       entities.NewWei(3),
		}

		receipt = test.AddDepositLogFromQuote(t, receipt, testPegoutQuote, testRetainedQuote)

		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(21), nil).Once()
		rskRpc.EXPECT().GetTransactionReceipt(mock.Anything, testRetainedQuote.UserRskTxHash).
			Return(*receipt, nil).Once()
		lbc.On("GetDepositEvents", mock.Anything, uint64(20), mock.MatchedBy(matchUinPtr(21))).Return([]quote.PegoutDeposit{}, nil).Once()
		pegoutRepository.EXPECT().GetQuote(mock.Anything, mock.Anything).Return(nil, errors.Join(assert.AnError, usecases.NonRecoverableError)).Once()

		assert.Eventually(t, func() bool {
			q, ok := depositWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			return assert.True(t, ok) && assert.NotEmpty(t, q)
		}, time.Second, 10*time.Millisecond)

		tickerChannel <- time.Now()

		assert.Eventually(t, func() bool {
			q, ok := depositWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			return assert.False(t, ok) && assert.Empty(t, q) && pegoutRepository.AssertExpectations(t) &&
				lbc.AssertExpectations(t) && rskRpc.AssertExpectations(t)
		}, time.Second, 10*time.Millisecond)
	})
	acceptPegoutChannel <- quote.AcceptedPegoutQuoteEvent{
		Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
		Quote:         testPegoutQuote,
		RetainedQuote: testRetainedQuote,
	}
	t.Run("should stop tracking after send pegout successfully", func(t *testing.T) {
		rskRpc.Calls = []mock.Call{}
		rskRpc.ExpectedCalls = []*mock.Call{}
		lbc.Calls = []mock.Call{}
		lbc.ExpectedCalls = []*mock.Call{}
		pegoutRepository.Calls = []mock.Call{}
		pegoutRepository.ExpectedCalls = []*mock.Call{}
		btcWallet.Calls = []mock.Call{}
		btcWallet.ExpectedCalls = []*mock.Call{}

		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(22), nil).Twice()
		receipt := &blockchain.TransactionReceipt{
			BlockNumber: 10,
			Value:       entities.NewWei(3),
		}
		receipt = test.AddDepositLogFromQuote(t, receipt, testPegoutQuote, testRetainedQuote)
		receipt.Value = entities.NewWei(0)
		receipt.To = "0xaabb"
		rskRpc.EXPECT().GetTransactionReceipt(mock.Anything, testRetainedQuote.UserRskTxHash).
			Return(*receipt, nil).Twice()
		lbc.On("GetDepositEvents", mock.Anything, uint64(21), mock.MatchedBy(matchUinPtr(22))).Return([]quote.PegoutDeposit{}, nil).Once()
		pegoutRepository.EXPECT().GetQuote(mock.Anything, mock.Anything).Return(&testPegoutQuote, nil).Once()
		pegoutRepository.EXPECT().GetPegoutCreationData(mock.Anything, testRetainedQuote.QuoteHash).Return(quote.PegoutCreationDataZeroValue()).Once()
		pegoutRepository.EXPECT().UpdateRetainedQuote(mock.Anything, mock.Anything).Return(nil).Once()
		rskRpc.EXPECT().GetBlockByHash(mock.Anything, mock.Anything).Return(blockchain.BlockInfo{Timestamp: time.Now()}, nil).Once()
		lbc.On("IsPegOutQuoteCompleted", testRetainedQuote.QuoteHash).Return(false, nil).Once()
		btcWallet.On("GetBalance").Return(entities.NewWei(10000), nil).Once()
		btcWallet.On("SendWithOpReturn", mock.Anything, mock.Anything, mock.Anything).Return(test.AnyHash, nil).Once()

		assert.Eventually(t, func() bool {
			q, ok := depositWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			return assert.True(t, ok) && assert.NotEmpty(t, q)
		}, time.Second, 10*time.Millisecond)

		tickerChannel <- time.Now()

		assert.Eventually(t, func() bool {
			q, ok := depositWatcher.GetWatchedQuote(testRetainedQuote.QuoteHash)
			return assert.False(t, ok) && assert.Empty(t, q) && pegoutRepository.AssertExpectations(t) &&
				lbc.AssertExpectations(t) && rskRpc.AssertExpectations(t) && btcWallet.AssertExpectations(t)
		}, time.Second, 10*time.Millisecond)
	})
}

func matchUinPtr(target uint64) func(uin *uint64) bool {
	return func(uin *uint64) bool {
		return *uin == target
	}
}
