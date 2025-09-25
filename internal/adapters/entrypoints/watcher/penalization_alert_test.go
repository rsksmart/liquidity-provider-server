package watcher_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewPenalizationAlertWatcher(t *testing.T) {
	rpc := blockchain.Rpc{
		Btc: &mocks.BtcRpcMock{},
		Rsk: &mocks.RootstockRpcServerMock{},
	}
	penalizationWatcher := watcher.NewPenalizationAlertWatcher(rpc, &liquidity_provider.PenalizationAlertUseCase{}, &mocks.TickerMock{}, time.Duration(1))
	assert.Equal(t, 5, test.CountNonZeroValues(penalizationWatcher))
}

func TestPenalizationAlertWatcher_Start(t *testing.T) {
	t.Run("shouldn't update block if use case had an error", func(t *testing.T) {
		rskRpc := &mocks.RootstockRpcServerMock{}
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(555, nil).Once()
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(600, nil).Once()
		collateral := &mocks.CollateralManagementContractMock{}
		collateral.EXPECT().GetPenalizedEvents(mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)
		useCase := liquidity_provider.NewPenalizationAlertUseCase(
			blockchain.RskContracts{CollateralManagement: collateral},
			&mocks.AlertSenderMock{},
			test.AnyString,
			mocks.NewPenalizedEventRepositoryMock(t),
		)
		ticker := &mocks.TickerMock{}
		tickerChannel := make(chan time.Time)
		ticker.EXPECT().C().Return(tickerChannel)
		penalizationWatcher := watcher.NewPenalizationAlertWatcher(blockchain.Rpc{Rsk: rskRpc}, useCase, ticker, time.Duration(1))
		err := penalizationWatcher.Prepare(context.Background())
		require.NoError(t, err)
		go penalizationWatcher.Start()
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			return assert.Equal(t, uint64(555), penalizationWatcher.GetCurrentBlock()) && rskRpc.AssertExpectations(t) && collateral.AssertExpectations(t)
		}, time.Second, 10*time.Millisecond)
	})
	t.Run("should update block if use case executed successfully", func(t *testing.T) {
		rskRpc := &mocks.RootstockRpcServerMock{}
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(555, nil).Once()
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(600, nil).Once()
		collateral := &mocks.CollateralManagementContractMock{}
		collateral.EXPECT().GetPenalizedEvents(mock.Anything, mock.Anything, mock.Anything).Return([]penalization.PenalizedEvent{}, nil)
		useCase := liquidity_provider.NewPenalizationAlertUseCase(
			blockchain.RskContracts{CollateralManagement: collateral},
			&mocks.AlertSenderMock{},
			test.AnyString,
			mocks.NewPenalizedEventRepositoryMock(t),
		)
		ticker := &mocks.TickerMock{}
		tickerChannel := make(chan time.Time)
		ticker.EXPECT().C().Return(tickerChannel)
		penalizationWatcher := watcher.NewPenalizationAlertWatcher(blockchain.Rpc{Rsk: rskRpc}, useCase, ticker, time.Duration(1))
		err := penalizationWatcher.Prepare(context.Background())
		require.NoError(t, err)
		go penalizationWatcher.Start()
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			return assert.Equal(t, uint64(599), penalizationWatcher.GetCurrentBlock()) && rskRpc.AssertExpectations(t) && collateral.AssertExpectations(t)
		}, time.Second, 10*time.Millisecond)
	})
}

func TestPenalizationAlertWatcher_Prepare(t *testing.T) {
	t.Run("prepare successfully", func(t *testing.T) {
		rskRpc := &mocks.RootstockRpcServerMock{}
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(555, nil).Once()
		penalizationWatcher := watcher.NewPenalizationAlertWatcher(blockchain.Rpc{Rsk: rskRpc}, nil, nil, time.Duration(1))
		err := penalizationWatcher.Prepare(context.Background())
		require.NoError(t, err)
		assert.Equal(t, uint64(555), penalizationWatcher.GetCurrentBlock())
		rskRpc.AssertExpectations(t)
	})
	t.Run("handle get height error", func(t *testing.T) {
		rskRpc := &mocks.RootstockRpcServerMock{}
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(0, assert.AnError).Once()
		penalizationWatcher := watcher.NewPenalizationAlertWatcher(blockchain.Rpc{Rsk: rskRpc}, nil, nil, time.Duration(1))
		err := penalizationWatcher.Prepare(context.Background())
		require.Error(t, err)
		assert.Zero(t, penalizationWatcher.GetCurrentBlock())
		rskRpc.AssertExpectations(t)
	})
}

func TestPenalizationAlertWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) watcher.Watcher {
		return watcher.NewPenalizationAlertWatcher(blockchain.Rpc{}, nil, ticker, time.Duration(1))
	})
}
