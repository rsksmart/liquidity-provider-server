package watcher_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestLiquidityCheckWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) watcher.Watcher {
		return watcher.NewLiquidityCheckWatcher(nil, ticker, time.Duration(1))
	})
}

func TestNewLiquidityCheckWatcher(t *testing.T) {
	ticker := &mocks.TickerMock{}
	providerMock := &mocks.ProviderMock{}
	useCase := liquidity_provider.NewCheckLiquidityUseCase(providerMock, providerMock, blockchain.RskContracts{}, &mocks.AlertSenderMock{}, test.AnyString)
	test.AssertNonZeroValues(t, watcher.NewLiquidityCheckWatcher(useCase, ticker, time.Duration(1)))
}

func TestLiquidityCheckWatcher_Start(t *testing.T) {
	tickerChannel := make(chan time.Time)
	ticker := &mocks.TickerMock{}
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop().Return()
	providerMock := &mocks.ProviderMock{}
	providerMock.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(nil)
	providerMock.On("HasPegoutLiquidity", mock.Anything, mock.Anything).Return(nil)
	bridgeMock := &mocks.BridgeMock{}
	bridgeMock.On("GetMinimumLockTxValue").Return(entities.NewWei(5), nil)
	useCase := liquidity_provider.NewCheckLiquidityUseCase(providerMock, providerMock, blockchain.RskContracts{Bridge: bridgeMock}, &mocks.AlertSenderMock{}, test.AnyString)
	w := watcher.NewLiquidityCheckWatcher(useCase, ticker, time.Duration(1))
	wg := sync.WaitGroup{}
	wg.Add(2)
	closeChannel := make(chan bool)
	defer test.AssertNoLog(t)()
	go func() {
		defer wg.Done()
		<-closeChannel
	}()
	go func() {
		defer wg.Done()
		w.Start()
	}()
	tickerChannel <- time.Now()
	w.Shutdown(closeChannel)
	wg.Wait()
	providerMock.AssertExpectations(t)
	bridgeMock.AssertExpectations(t)
}

func TestLiquidityCheckWatcher_Start_ErrorHandling(t *testing.T) {
	tickerChannel := make(chan time.Time)
	closeChannel := make(chan bool)
	ticker := &mocks.TickerMock{}
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop().Return()
	providerMock := &mocks.ProviderMock{}
	bridgeMock := &mocks.BridgeMock{}
	bridgeMock.On("GetMinimumLockTxValue").Return(nil, assert.AnError)
	useCase := liquidity_provider.NewCheckLiquidityUseCase(providerMock, providerMock, blockchain.RskContracts{Bridge: bridgeMock}, &mocks.AlertSenderMock{}, test.AnyString)
	w := watcher.NewLiquidityCheckWatcher(useCase, ticker, time.Duration(1))
	wg := sync.WaitGroup{}
	wg.Add(2)
	defer test.AssertLogContains(t, assert.AnError.Error())
	go func() {
		defer wg.Done()
		w.Start()
	}()
	go func() {
		defer wg.Done()
		<-closeChannel
	}()
	tickerChannel <- time.Now()
	w.Shutdown(closeChannel)
	wg.Wait()
	bridgeMock.AssertExpectations(t)
}

func TestLiquidityCheckWatcher_Prepare(t *testing.T) {
	w := watcher.NewLiquidityCheckWatcher(nil, &mocks.TickerMock{}, time.Duration(1))
	require.NoError(t, w.Prepare(context.Background()))
}
