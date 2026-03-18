package watcher_test

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	w "github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher/monitoring"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBitcoinPeerWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) w.Watcher {
		btcMock := &mocks.BtcRpcMock{}
		useCase := &mocks.NodePeerAlertUseCaseMock{}
		rpc := blockchain.Rpc{Btc: btcMock}
		appMetrics := monitoring.NewMetrics(prometheus.NewRegistry())
		return w.NewBitcoinPeerWatcher(rpc, useCase, ticker, 3, 15*time.Second, 30*time.Minute, appMetrics)
	})
}

// nolint:funlen
func TestBitcoinPeerWatcher_Start(t *testing.T) {
	t.Run("should alert when peer count is below threshold", func(t *testing.T) {
		runBitcoinPeerSubtest(t, int64(1), nil, 3, true, nil)
	})
	t.Run("should not alert when peer count is at or above threshold", func(t *testing.T) {
		runBitcoinPeerSubtest(t, int64(5), nil, 3, false, nil)
	})
	t.Run("should not alert when peer count equals threshold", func(t *testing.T) {
		runBitcoinPeerSubtest(t, int64(3), nil, 3, false, nil)
	})
	t.Run("should not alert when minPeers is zero", func(t *testing.T) {
		runBitcoinPeerSubtest(t, int64(0), nil, 0, false, nil)
	})
	t.Run("should continue running on RPC error", func(t *testing.T) {
		runBitcoinPeerSubtest(t, int64(0), assert.AnError, 3, false, nil)
	})
	t.Run("should continue running on alert send error", func(t *testing.T) {
		runBitcoinPeerSubtest(t, int64(1), nil, 3, true, assert.AnError)
	})
	t.Run("should suppress alert during cooldown period", func(t *testing.T) {
		btcMock := &mocks.BtcRpcMock{}
		useCase := &mocks.NodePeerAlertUseCaseMock{}
		ticker := &mocks.TickerMock{}
		tickerChannel := make(chan time.Time)
		closeChannel := make(chan bool)
		ticker.EXPECT().C().Return(tickerChannel)
		ticker.EXPECT().Stop().Return().Once()
		btcMock.On("GetConnectionCount").Return(int64(1), nil)
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeBitcoin, int64(1), uint64(3)).Return(nil).Once()
		rpc := blockchain.Rpc{Btc: btcMock}
		appMetrics := monitoring.NewMetrics(prometheus.NewRegistry())
		watcher := w.NewBitcoinPeerWatcher(rpc, useCase, ticker, 3, 15*time.Second, time.Hour, appMetrics)
		go watcher.Start()
		tickerChannel <- time.Now()
		tickerChannel <- time.Now()
		go watcher.Shutdown(closeChannel)
		<-closeChannel
		assert.Eventually(t, func() bool {
			return ticker.AssertExpectations(t) && useCase.AssertExpectations(t)
		}, time.Second, 100*time.Millisecond)
	})
}

func runBitcoinPeerSubtest(t *testing.T, peerCount int64, rpcErr error, minPeers uint64, expectAlert bool, alertErr error) {
	t.Helper()
	btcMock := &mocks.BtcRpcMock{}
	useCase := &mocks.NodePeerAlertUseCaseMock{}
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	closeChannel := make(chan bool)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop().Return().Once()
	btcMock.On("GetConnectionCount").Return(peerCount, rpcErr).Once()
	if expectAlert {
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeBitcoin, peerCount, minPeers).Return(alertErr).Once()
	}
	rpc := blockchain.Rpc{Btc: btcMock}
	appMetrics := monitoring.NewMetrics(prometheus.NewRegistry())
	watcher := w.NewBitcoinPeerWatcher(rpc, useCase, ticker, minPeers, 15*time.Second, 30*time.Minute, appMetrics)
	go watcher.Start()
	tickerChannel <- time.Now()
	go watcher.Shutdown(closeChannel)
	<-closeChannel
	assert.Eventually(t, func() bool {
		return ticker.AssertExpectations(t) && btcMock.AssertExpectations(t)
	}, time.Second, 100*time.Millisecond)
	if !expectAlert {
		useCase.AssertNotCalled(t, "Run", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	}
}
