package watcher_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPenalizationAlertWatcher_Start(t *testing.T) {
	t.Run("shouldn't update block if use case had an error", func(t *testing.T) {
		rskRpc := &mocks.RootstockRpcServerMock{}
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(555, nil).Once()
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(600, nil).Once()
		lbc := &mocks.LbcMock{}
		lbc.On("GetPeginPunishmentEvents", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)
		useCase := liquidity_provider.NewPenalizationAlertUseCase(blockchain.RskContracts{Lbc: lbc}, &mocks.AlertSenderMock{}, test.AnyString)
		ticker := &mocks.TickerMock{}
		tickerChannel := make(chan time.Time)
		ticker.EXPECT().C().Return(tickerChannel)
		penalizationWatcher := watcher.NewPenalizationAlertWatcher(blockchain.Rpc{Rsk: rskRpc}, useCase, ticker)
		err := penalizationWatcher.Prepare(context.Background())
		require.NoError(t, err)
		go penalizationWatcher.Start()
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			return assert.Equal(t, uint64(555), penalizationWatcher.GetCurrentBlock()) && rskRpc.AssertExpectations(t) && lbc.AssertExpectations(t)
		}, time.Second, 10*time.Millisecond)
	})
	t.Run("should update block if use case executed successfully", func(t *testing.T) {
		rskRpc := &mocks.RootstockRpcServerMock{}
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(555, nil).Once()
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(600, nil).Once()
		lbc := &mocks.LbcMock{}
		lbc.On("GetPeginPunishmentEvents", mock.Anything, mock.Anything, mock.Anything).Return([]lp.PunishmentEvent{}, nil)
		useCase := liquidity_provider.NewPenalizationAlertUseCase(blockchain.RskContracts{Lbc: lbc}, &mocks.AlertSenderMock{}, test.AnyString)
		ticker := &mocks.TickerMock{}
		tickerChannel := make(chan time.Time)
		ticker.EXPECT().C().Return(tickerChannel)
		penalizationWatcher := watcher.NewPenalizationAlertWatcher(blockchain.Rpc{Rsk: rskRpc}, useCase, ticker)
		err := penalizationWatcher.Prepare(context.Background())
		require.NoError(t, err)
		go penalizationWatcher.Start()
		tickerChannel <- time.Now()
		assert.Eventually(t, func() bool {
			return assert.Equal(t, uint64(599), penalizationWatcher.GetCurrentBlock()) && rskRpc.AssertExpectations(t) && lbc.AssertExpectations(t)
		}, time.Second, 10*time.Millisecond)
	})
}

func TestPenalizationAlertWatcher_Prepare(t *testing.T) {
	t.Run("prepare successfully", func(t *testing.T) {
		rskRpc := &mocks.RootstockRpcServerMock{}
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(555, nil).Once()
		penalizationWatcher := watcher.NewPenalizationAlertWatcher(blockchain.Rpc{Rsk: rskRpc}, nil, nil)
		err := penalizationWatcher.Prepare(context.Background())
		require.NoError(t, err)
		assert.Equal(t, uint64(555), penalizationWatcher.GetCurrentBlock())
		rskRpc.AssertExpectations(t)
	})
	t.Run("handle get height error", func(t *testing.T) {
		rskRpc := &mocks.RootstockRpcServerMock{}
		rskRpc.EXPECT().GetHeight(mock.Anything).Return(0, assert.AnError).Once()
		penalizationWatcher := watcher.NewPenalizationAlertWatcher(blockchain.Rpc{Rsk: rskRpc}, nil, nil)
		err := penalizationWatcher.Prepare(context.Background())
		require.Error(t, err)
		assert.Zero(t, penalizationWatcher.GetCurrentBlock())
		rskRpc.AssertExpectations(t)
	})
}

func TestPenalizationAlertWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker watcher.Ticker) watcher.Watcher {
		return watcher.NewPenalizationAlertWatcher(blockchain.Rpc{}, nil, ticker)
	})
}
