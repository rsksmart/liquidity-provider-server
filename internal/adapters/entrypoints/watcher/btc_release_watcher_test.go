package watcher_test

import (
	"context"
	w "github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBtcReleaseWatcher_Prepare(t *testing.T) {
	t.Run("should start from latest if start block is 0", func(t *testing.T) {
		const (
			startBlock = 0
			pageSize   = 10
			timeout    = 30 * 1000
		)
		contracts := blockchain.RskContracts{Bridge: &mocks.BridgeMock{}}
		rskRpc := &mocks.RootstockRpcServerMock{}
		rpc := blockchain.Rpc{Rsk: rskRpc}
		useCase := &mocks.UpdateBtcReleaseUseCaseMock{}
		watcher := w.NewBtcReleaseWatcher(contracts, rpc, useCase, &mocks.TickerMock{}, startBlock, pageSize, timeout)

		rskRpc.EXPECT().GetHeight(mock.Anything).Return(50, nil).Once()
		ctx := context.Background()
		err := watcher.Prepare(ctx)
		require.NoError(t, err)
		assert.Equal(t, uint64(50), watcher.CurrentBlock())
		rskRpc.AssertExpectations(t)
	})
	t.Run("should start from given block if start block is set", func(t *testing.T) {
		const (
			startBlock = 15
			pageSize   = 10
			timeout    = 30 * 1000
		)
		contracts := blockchain.RskContracts{
			Bridge: &mocks.BridgeMock{},
		}
		rskRpc := &mocks.RootstockRpcServerMock{}
		rpc := blockchain.Rpc{Rsk: rskRpc}
		useCase := &mocks.UpdateBtcReleaseUseCaseMock{}
		watcher := w.NewBtcReleaseWatcher(contracts, rpc, useCase, &mocks.TickerMock{}, startBlock, pageSize, timeout)

		ctx := context.Background()
		err := watcher.Prepare(ctx)
		require.NoError(t, err)
		assert.Equal(t, uint64(15), watcher.CurrentBlock())
		rskRpc.AssertNotCalled(t, "GetHeight")
	})
	t.Run("should handle error when getting height", func(t *testing.T) {
		const (
			startBlock = 0
			pageSize   = 10
			timeout    = 30 * 1000
		)
		contracts := blockchain.RskContracts{Bridge: &mocks.BridgeMock{}}
		rskRpc := &mocks.RootstockRpcServerMock{}
		rpc := blockchain.Rpc{Rsk: rskRpc}
		useCase := &mocks.UpdateBtcReleaseUseCaseMock{}
		watcher := w.NewBtcReleaseWatcher(contracts, rpc, useCase, &mocks.TickerMock{}, startBlock, pageSize, timeout)

		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(0), assert.AnError).Once()
		ctx := context.Background()
		err := watcher.Prepare(ctx)
		require.Error(t, err)
		rskRpc.AssertExpectations(t)
	})
}

// nolint:funlen
func TestBtcReleaseWatcher_Start(t *testing.T) {
	const (
		startBlock = 10
		pageSize   = 15
		timeout    = 30 * 1000
	)
	mockEvents := []rootstock.BatchPegOut{
		{
			TransactionHash:    test.AnyHash,
			BlockHash:          test.AnyString,
			BlockNumber:        5,
			BtcTxHash:          test.AnyUrl,
			ReleaseRskTxHashes: []string{"1", "2", "3"},
		},
		{
			TransactionHash:    test.AnyString,
			BlockHash:          test.AnyHash,
			BlockNumber:        8,
			BtcTxHash:          test.AnyUrl,
			ReleaseRskTxHashes: []string{"4", "5", "6"},
		},
	}
	t.Run("should run tick with a full page", func(t *testing.T) {
		bridge := &mocks.BridgeMock{}
		contracts := blockchain.RskContracts{Bridge: bridge}
		rskRpc := &mocks.RootstockRpcServerMock{}
		rpc := blockchain.Rpc{Rsk: rskRpc}
		useCase := &mocks.UpdateBtcReleaseUseCaseMock{}
		ticker := &mocks.TickerMock{}
		watcher := w.NewBtcReleaseWatcher(contracts, rpc, useCase, ticker, startBlock, pageSize, timeout)

		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(50), nil).Once()
		bridge.On("GetBatchPegOutCreatedEvent", mock.Anything, uint64(10), mock.MatchedBy(func(toBlock *uint64) bool {
			return toBlock != nil && *toBlock == 25
		})).Return(mockEvents, nil).Once()
		useCase.EXPECT().Run(mock.Anything, mockEvents[0]).Return(uint(3), nil).Once()
		useCase.EXPECT().Run(mock.Anything, mockEvents[1]).Return(uint(0), nil).Once()
		tickerChannel := make(chan time.Time)
		ticker.EXPECT().C().Return(tickerChannel)
		assertQuotesLog := test.AssertLogContains(t, "Successfully processed 3 quotes in BatchPegOut (d8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb)")

		err := watcher.Prepare(context.Background())
		require.NoError(t, err)
		go watcher.Start()

		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			return useCase.AssertExpectations(t) &&
				bridge.AssertExpectations(t) &&
				rskRpc.AssertExpectations(t) &&
				assertQuotesLog()
		}, time.Second*3, time.Millisecond*100)
	})
	t.Run("should run tick with a with reduce page if its too close to the latest block", func(t *testing.T) {
		bridge := &mocks.BridgeMock{}
		contracts := blockchain.RskContracts{Bridge: bridge}
		rskRpc := &mocks.RootstockRpcServerMock{}
		rpc := blockchain.Rpc{Rsk: rskRpc}
		useCase := &mocks.UpdateBtcReleaseUseCaseMock{}
		ticker := &mocks.TickerMock{}
		watcher := w.NewBtcReleaseWatcher(contracts, rpc, useCase, ticker, startBlock, pageSize, timeout)

		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(20), nil).Once()
		bridge.On("GetBatchPegOutCreatedEvent", mock.Anything, uint64(10), mock.MatchedBy(func(toBlock *uint64) bool {
			return toBlock != nil && *toBlock == 20
		})).Return(mockEvents, nil).Once()
		useCase.EXPECT().Run(mock.Anything, mockEvents[0]).Return(uint(0), nil).Once()
		useCase.EXPECT().Run(mock.Anything, mockEvents[1]).Return(uint(0), nil).Once()
		tickerChannel := make(chan time.Time)
		ticker.EXPECT().C().Return(tickerChannel)
		assertNoQuotesLog := test.AssertLogContains(t, "No PegOuts to process in batch (any value)")

		err := watcher.Prepare(context.Background())
		require.NoError(t, err)
		go watcher.Start()

		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			return useCase.AssertExpectations(t) &&
				bridge.AssertExpectations(t) &&
				rskRpc.AssertExpectations(t) &&
				assertNoQuotesLog()
		}, time.Second*3, time.Millisecond*100)
	})
}

// nolint:funlen
func TestBtcReleaseWatcher_Start_ErrorCases(t *testing.T) {
	const (
		startBlock = 10
		pageSize   = 15
		timeout    = 30 * 1000
	)
	mockEvents := []rootstock.BatchPegOut{
		{
			TransactionHash:    test.AnyHash,
			BlockHash:          test.AnyString,
			BlockNumber:        5,
			BtcTxHash:          test.AnyUrl,
			ReleaseRskTxHashes: []string{"1", "2", "3"},
		},
		{
			TransactionHash:    test.AnyString,
			BlockHash:          test.AnyHash,
			BlockNumber:        8,
			BtcTxHash:          test.AnyUrl,
			ReleaseRskTxHashes: []string{"4", "5", "6"},
		},
	}
	t.Run("should handle error getting height", func(t *testing.T) {
		bridge := &mocks.BridgeMock{}
		contracts := blockchain.RskContracts{Bridge: bridge}
		rskRpc := &mocks.RootstockRpcServerMock{}
		rpc := blockchain.Rpc{Rsk: rskRpc}
		useCase := &mocks.UpdateBtcReleaseUseCaseMock{}
		ticker := &mocks.TickerMock{}
		watcher := w.NewBtcReleaseWatcher(contracts, rpc, useCase, ticker, startBlock, pageSize, timeout)

		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(0), assert.AnError).Once()
		tickerChannel := make(chan time.Time)
		ticker.EXPECT().C().Return(tickerChannel)
		assertErrorLog := test.AssertLogContains(t, "error getting RSK height in BtcReleaseWatcher")

		err := watcher.Prepare(context.Background())
		require.NoError(t, err)
		go watcher.Start()

		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			return useCase.AssertNotCalled(t, "Run") &&
				bridge.AssertNotCalled(t, "GetBatchPegOutCreatedEvent") &&
				rskRpc.AssertExpectations(t) &&
				assertErrorLog()
		}, time.Second*3, time.Millisecond*100)
	})
	t.Run("should handle error getting event", func(t *testing.T) {
		bridge := &mocks.BridgeMock{}
		contracts := blockchain.RskContracts{Bridge: bridge}
		rskRpc := &mocks.RootstockRpcServerMock{}
		rpc := blockchain.Rpc{Rsk: rskRpc}
		useCase := &mocks.UpdateBtcReleaseUseCaseMock{}
		ticker := &mocks.TickerMock{}
		watcher := w.NewBtcReleaseWatcher(contracts, rpc, useCase, ticker, startBlock, pageSize, timeout)

		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(20), nil).Once()
		bridge.On("GetBatchPegOutCreatedEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		tickerChannel := make(chan time.Time)
		ticker.EXPECT().C().Return(tickerChannel)
		assertNoQuotesLog := test.AssertLogContains(t, "error fetching BatchPegOutCreated events in BtcReleaseWatcher")

		err := watcher.Prepare(context.Background())
		require.NoError(t, err)
		go watcher.Start()

		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			return useCase.AssertNotCalled(t, "Run") &&
				bridge.AssertExpectations(t) &&
				rskRpc.AssertExpectations(t) &&
				assertNoQuotesLog()
		}, time.Second*3, time.Millisecond*100)
	})
	t.Run("should handle error in use case", func(t *testing.T) {
		bridge := &mocks.BridgeMock{}
		contracts := blockchain.RskContracts{Bridge: bridge}
		rskRpc := &mocks.RootstockRpcServerMock{}
		rpc := blockchain.Rpc{Rsk: rskRpc}
		useCase := &mocks.UpdateBtcReleaseUseCaseMock{}
		ticker := &mocks.TickerMock{}
		watcher := w.NewBtcReleaseWatcher(contracts, rpc, useCase, ticker, startBlock, pageSize, timeout)

		rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(20), nil).Once()
		bridge.On("GetBatchPegOutCreatedEvent", mock.Anything, mock.Anything, mock.Anything).Return(mockEvents, nil).Once()
		useCase.EXPECT().Run(mock.Anything, mock.Anything).Return(uint(0), assert.AnError).Once()
		tickerChannel := make(chan time.Time)
		ticker.EXPECT().C().Return(tickerChannel)
		assertNoQuotesLog := test.AssertLogContains(t, "error processing BatchPegOut")

		err := watcher.Prepare(context.Background())
		require.NoError(t, err)
		go watcher.Start()

		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			return useCase.AssertExpectations(t) && bridge.AssertExpectations(t) && rskRpc.AssertExpectations(t) && assertNoQuotesLog()
		}, time.Second*3, time.Millisecond*100)
	})
}

func TestBtcReleaseWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) w.Watcher {
		contracts := blockchain.RskContracts{
			Bridge: &mocks.BridgeMock{},
		}
		rpc := blockchain.Rpc{
			Rsk: &mocks.RootstockRpcServerMock{},
		}
		useCase := &mocks.UpdateBtcReleaseUseCaseMock{}
		return w.NewBtcReleaseWatcher(contracts, rpc, useCase, ticker, 0, 10, 30)
	})
}
