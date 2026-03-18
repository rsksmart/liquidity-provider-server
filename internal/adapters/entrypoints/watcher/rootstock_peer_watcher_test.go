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

func TestRootstockPeerWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) w.Watcher {
		rskMock := &mocks.RootstockRpcServerMock{}
		useCase := &mocks.NodePeerAlertUseCaseMock{}
		rpc := blockchain.Rpc{Rsk: rskMock}
		appMetrics := monitoring.NewMetrics(prometheus.NewRegistry())
		return w.NewRootstockPeerWatcher(rpc, useCase, ticker, 3, 15*time.Second, 30*time.Minute, appMetrics)
	})
}

// nolint:funlen
func TestRootstockPeerWatcher_Start(t *testing.T) {
	t.Run("should alert when peer count is below threshold", func(t *testing.T) {
		runRootstockPeerSubtest(t, uint64(1), nil, 3, true, nil)
	})
	t.Run("should not alert when peer count is at or above threshold", func(t *testing.T) {
		runRootstockPeerSubtest(t, uint64(5), nil, 3, false, nil)
	})
	t.Run("should not alert when peer count equals threshold", func(t *testing.T) {
		runRootstockPeerSubtest(t, uint64(3), nil, 3, false, nil)
	})
	t.Run("should not alert when minPeers is zero", func(t *testing.T) {
		runRootstockPeerSubtest(t, uint64(0), nil, 0, false, nil)
	})
	t.Run("should continue running on RPC error", func(t *testing.T) {
		runRootstockPeerSubtest(t, uint64(0), assert.AnError, 3, false, nil)
	})
	t.Run("should continue running on alert send error", func(t *testing.T) {
		runRootstockPeerSubtest(t, uint64(1), nil, 3, true, assert.AnError)
	})
	t.Run("should suppress alert during cooldown period", func(t *testing.T) {
		rskMock := &mocks.RootstockRpcServerMock{}
		useCase := &mocks.NodePeerAlertUseCaseMock{}
		ticker := &mocks.TickerMock{}
		tickerChannel := make(chan time.Time)
		closeChannel := make(chan bool)
		ticker.EXPECT().C().Return(tickerChannel)
		ticker.EXPECT().Stop().Return().Once()
		rskMock.On("PeerCount", mock.Anything).Return(uint64(1), nil)
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeRootstock, int64(1), uint64(3)).Return(nil).Once()
		rpc := blockchain.Rpc{Rsk: rskMock}
		appMetrics := monitoring.NewMetrics(prometheus.NewRegistry())
		watcher := w.NewRootstockPeerWatcher(rpc, useCase, ticker, 3, 15*time.Second, time.Hour, appMetrics)
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

func runRootstockPeerSubtest(t *testing.T, peerCount uint64, rpcErr error, minPeers uint64, expectAlert bool, alertErr error) {
	t.Helper()
	rskMock := &mocks.RootstockRpcServerMock{}
	useCase := &mocks.NodePeerAlertUseCaseMock{}
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	closeChannel := make(chan bool)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop().Return().Once()
	rskMock.On("PeerCount", mock.Anything).Return(peerCount, rpcErr).Once()
	if expectAlert {
		useCase.EXPECT().Run(mock.Anything, entities.NodeTypeRootstock, int64(peerCount), minPeers).Return(alertErr).Once()
	}
	rpc := blockchain.Rpc{Rsk: rskMock}
	appMetrics := monitoring.NewMetrics(prometheus.NewRegistry())
	watcher := w.NewRootstockPeerWatcher(rpc, useCase, ticker, minPeers, 15*time.Second, 30*time.Minute, appMetrics)
	go watcher.Start()
	tickerChannel <- time.Now()
	go watcher.Shutdown(closeChannel)
	<-closeChannel
	assert.Eventually(t, func() bool {
		return ticker.AssertExpectations(t) && rskMock.AssertExpectations(t)
	}, time.Second, 100*time.Millisecond)
	if !expectAlert {
		useCase.AssertNotCalled(t, "Run", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	}
}
