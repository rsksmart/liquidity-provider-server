package watcher_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	lp "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/require"
)

const testTransferTimeout = 60 * time.Second

func TestNewTransferColdWalletWatcher(t *testing.T) {
	ticker := &mocks.TickerMock{}
	useCase := &mockTransferUseCase{}

	w := watcher.NewTransferColdWalletWatcher(useCase, ticker, testTransferTimeout)

	test.AssertNonZeroValues(t, w)
}

func TestTransferColdWalletWatcher_Prepare(t *testing.T) {
	ticker := &mocks.TickerMock{}
	useCase := &mockTransferUseCase{}
	w := watcher.NewTransferColdWalletWatcher(useCase, ticker, testTransferTimeout)

	err := w.Prepare(context.Background())

	require.NoError(t, err)
}

func TestTransferColdWalletWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) watcher.Watcher {
		useCase := &mockTransferUseCase{}
		return watcher.NewTransferColdWalletWatcher(useCase, ticker, testTransferTimeout)
	})
}

func TestTransferColdWalletWatcher_Start_Error(t *testing.T) {
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop()

	expectedError := errors.New("cold wallet not configured")
	useCase := &mockTransferUseCase{err: expectedError}
	w := watcher.NewTransferColdWalletWatcher(useCase, ticker, testTransferTimeout)
	
	wg := sync.WaitGroup{}
	wg.Add(2)
	closeChannel := make(chan bool)
	defer test.AssertLogContains(t, "Error executing transfer to cold wallet")()
	
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
	ticker.AssertExpectations(t)
}

func TestTransferColdWalletWatcher_Start_BtcSuccess(t *testing.T) {
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop()

	result := &lp.TransferToColdWalletResult{
		BtcResult: lp.NetworkTransferResult{
			Status: lp.TransferStatusSuccess,
			TxHash: "btc_tx_hash_123",
			Amount: entities.NewWei(1000000),
			Fee:    entities.NewWei(5000),
		},
		RskResult: lp.NetworkTransferResult{
			Status: lp.TransferStatusSkippedNoExcess,
		},
	}
	useCase := &mockTransferUseCase{result: result}
	w := watcher.NewTransferColdWalletWatcher(useCase, ticker, testTransferTimeout)

	wg := sync.WaitGroup{}
	wg.Add(2)
	closeChannel := make(chan bool)
	defer test.AssertLogContains(t, "BTC transfer successful - TxHash: btc_tx_hash_123, Amount: 1000000, Fee: 5000")()
	
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
	ticker.AssertExpectations(t)
}

func TestTransferColdWalletWatcher_Start_RskSuccess(t *testing.T) {
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop()

	result := &lp.TransferToColdWalletResult{
		BtcResult: lp.NetworkTransferResult{
			Status: lp.TransferStatusSkippedNoExcess,
		},
		RskResult: lp.NetworkTransferResult{
			Status: lp.TransferStatusSuccess,
			TxHash: "rsk_tx_hash_456",
			Amount: entities.NewWei(2000000),
			Fee:    entities.NewWei(3000),
		},
	}
	useCase := &mockTransferUseCase{result: result}
	w := watcher.NewTransferColdWalletWatcher(useCase, ticker, testTransferTimeout)

	wg := sync.WaitGroup{}
	wg.Add(2)
	closeChannel := make(chan bool)
	defer test.AssertLogContains(t, "RSK transfer successful - TxHash: rsk_tx_hash_456, Amount: 2000000, Fee: 3000")()
	
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
	ticker.AssertExpectations(t)
}

func TestTransferColdWalletWatcher_Start_BothSkippedNoExcess(t *testing.T) {
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop()

	result := &lp.TransferToColdWalletResult{
		BtcResult: lp.NetworkTransferResult{
			Status: lp.TransferStatusSkippedNoExcess,
		},
		RskResult: lp.NetworkTransferResult{
			Status: lp.TransferStatusSkippedNoExcess,
		},
	}
	useCase := &mockTransferUseCase{result: result}
	w := watcher.NewTransferColdWalletWatcher(useCase, ticker, testTransferTimeout)

	wg := sync.WaitGroup{}
	wg.Add(2)
	closeChannel := make(chan bool)
	defer test.AssertLogContains(t, "transfer skipped - no excess liquidity")()
	
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
	ticker.AssertExpectations(t)
}

func TestTransferColdWalletWatcher_Start_BtcSkippedNotEconomical(t *testing.T) {
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop()

	result := &lp.TransferToColdWalletResult{
		BtcResult: lp.NetworkTransferResult{
			Status:  lp.TransferStatusSkippedNotEconomical,
			Message: "transfer amount too small",
		},
		RskResult: lp.NetworkTransferResult{
			Status: lp.TransferStatusSkippedNoExcess,
		},
	}
	useCase := &mockTransferUseCase{result: result}
	w := watcher.NewTransferColdWalletWatcher(useCase, ticker, testTransferTimeout)

	wg := sync.WaitGroup{}
	wg.Add(2)
	closeChannel := make(chan bool)
	defer test.AssertLogContains(t, "BTC transfer skipped - not economical: transfer amount too small")()
	
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
	ticker.AssertExpectations(t)
}

func TestTransferColdWalletWatcher_Start_RskSkippedNotEconomical(t *testing.T) {
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop()

	result := &lp.TransferToColdWalletResult{
		BtcResult: lp.NetworkTransferResult{
			Status: lp.TransferStatusSkippedNoExcess,
		},
		RskResult: lp.NetworkTransferResult{
			Status:  lp.TransferStatusSkippedNotEconomical,
			Message: "gas cost too high",
		},
	}
	useCase := &mockTransferUseCase{result: result}
	w := watcher.NewTransferColdWalletWatcher(useCase, ticker, testTransferTimeout)

	wg := sync.WaitGroup{}
	wg.Add(2)
	closeChannel := make(chan bool)
	defer test.AssertLogContains(t, "RSK transfer skipped - not economical: gas cost too high")()
	
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
	ticker.AssertExpectations(t)
}

func TestTransferColdWalletWatcher_Start_BtcFailed(t *testing.T) {
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop()

	transferError := errors.New("insufficient funds")
	result := &lp.TransferToColdWalletResult{
		BtcResult: lp.NetworkTransferResult{
			Status:  lp.TransferStatusFailed,
			Message: "transfer failed",
			Error:   transferError,
		},
		RskResult: lp.NetworkTransferResult{
			Status: lp.TransferStatusSkippedNoExcess,
		},
	}
	useCase := &mockTransferUseCase{result: result}
	w := watcher.NewTransferColdWalletWatcher(useCase, ticker, testTransferTimeout)

	wg := sync.WaitGroup{}
	wg.Add(2)
	closeChannel := make(chan bool)
	defer test.AssertLogContains(t, "BTC transfer failed - transfer failed: insufficient funds")()
	
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
	ticker.AssertExpectations(t)
}

func TestTransferColdWalletWatcher_Start_RskFailed(t *testing.T) {
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop()

	transferError := errors.New("gas price too low")
	result := &lp.TransferToColdWalletResult{
		BtcResult: lp.NetworkTransferResult{
			Status: lp.TransferStatusSkippedNoExcess,
		},
		RskResult: lp.NetworkTransferResult{
			Status:  lp.TransferStatusFailed,
			Message: "rsk transfer failed",
			Error:   transferError,
		},
	}
	useCase := &mockTransferUseCase{result: result}
	w := watcher.NewTransferColdWalletWatcher(useCase, ticker, testTransferTimeout)

	wg := sync.WaitGroup{}
	wg.Add(2)
	closeChannel := make(chan bool)
	defer test.AssertLogContains(t, "RSK transfer failed - rsk transfer failed: gas price too low")()
	
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
	ticker.AssertExpectations(t)
}

func TestTransferColdWalletWatcher_Start_BothTransfersSuccess(t *testing.T) {
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop()

	result := &lp.TransferToColdWalletResult{
		BtcResult: lp.NetworkTransferResult{
			Status: lp.TransferStatusSuccess,
			TxHash: "btc_tx_hash_789",
			Amount: entities.NewWei(5000000),
			Fee:    entities.NewWei(10000),
		},
		RskResult: lp.NetworkTransferResult{
			Status: lp.TransferStatusSuccess,
			TxHash: "rsk_tx_hash_012",
			Amount: entities.NewWei(3000000),
			Fee:    entities.NewWei(5000),
		},
	}
	useCase := &mockTransferUseCase{result: result}
	w := watcher.NewTransferColdWalletWatcher(useCase, ticker, testTransferTimeout)

	wg := sync.WaitGroup{}
	wg.Add(2)
	closeChannel := make(chan bool)
	defer test.AssertLogContains(t, "BTC transfer successful - TxHash: btc_tx_hash_789, Amount: 5000000, Fee: 10000")()
	
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
	ticker.AssertExpectations(t)
}

func TestTransferColdWalletWatcher_Start_BothFailed(t *testing.T) {
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop()

	result := &lp.TransferToColdWalletResult{
		BtcResult: lp.NetworkTransferResult{
			Status:  lp.TransferStatusFailed,
			Message: "btc error",
			Error:   errors.New("btc wallet unavailable"),
		},
		RskResult: lp.NetworkTransferResult{
			Status:  lp.TransferStatusFailed,
			Message: "rsk error",
			Error:   errors.New("rsk node disconnected"),
		},
	}
	useCase := &mockTransferUseCase{result: result}
	w := watcher.NewTransferColdWalletWatcher(useCase, ticker, testTransferTimeout)

	wg := sync.WaitGroup{}
	wg.Add(2)
	closeChannel := make(chan bool)
	defer test.AssertLogContains(t, "BTC transfer failed - btc error: btc wallet unavailable")()
	
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
	ticker.AssertExpectations(t)
}

func TestTransferColdWalletWatcher_Start_NilResult(t *testing.T) {
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop()

	useCase := &mockTransferUseCase{result: nil}
	w := watcher.NewTransferColdWalletWatcher(useCase, ticker, testTransferTimeout)

	wg := sync.WaitGroup{}
	wg.Add(2)
	closeChannel := make(chan bool)
	
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
	ticker.AssertExpectations(t)
}

// mockTransferUseCase implements watcher.TransferExcessToColdWalletUseCase interface
type mockTransferUseCase struct {
	result *lp.TransferToColdWalletResult
	err    error
}

func (m *mockTransferUseCase) Run(ctx context.Context) (*lp.TransferToColdWalletResult, error) {
	return m.result, m.err
}
