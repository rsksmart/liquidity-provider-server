package watcher_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type countingClaimRunner struct {
	runs atomic.Int32
}

func (runner *countingClaimRunner) Run(
	_ context.Context,
	_ rootstock.PegInWatch,
	_ string,
) error {
	runner.runs.Add(1)
	return nil
}

func (runner *countingClaimRunner) ReconcileSubmitting(context.Context) error {
	return nil
}

type recordingClaimRunner struct {
	countingClaimRunner
	reconciles atomic.Int32
}

func (runner *recordingClaimRunner) ReconcileSubmitting(context.Context) error {
	runner.reconciles.Add(1)
	return nil
}

type errorClaimRunner struct {
	countingClaimRunner
}

func (runner *errorClaimRunner) Run(
	_ context.Context,
	_ rootstock.PegInWatch,
	_ string,
) error {
	runner.runs.Add(1)
	return assert.AnError
}

func importedWatchEntry() rootstock.PegInWatch {
	return rootstock.PegInWatch{
		RskAddress: test.AnyRskAddress,
		BtcAddress: "bcrt1qimported",
		State:      rootstock.PegInWatchImported,
	}
}

func payingWatchTx(btcAddress string) blockchain.BitcoinTransactionInformation {
	return blockchain.BitcoinTransactionInformation{
		Hash: "aabbcc",
		Outputs: map[string][]*entities.Wei{
			btcAddress: {entities.NewWei(1)},
		},
	}
}

func runClaimWatcherTick(
	t *testing.T,
	runner watcher.PegInClaimRunner,
	setup func(*mocks.PegInWatchRepositoryMock, *mocks.BitcoinWalletMock),
	assertTick func(*assert.CollectT, *mocks.PegInWatchRepositoryMock, *mocks.BitcoinWalletMock),
) {
	t.Helper()
	watchRepo := mocks.NewPegInWatchRepositoryMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	setup(watchRepo, wallet)

	ticker := &mocks.TickerMock{}
	ticks := make(chan time.Time)
	ticker.EXPECT().C().Return(ticks)
	ticker.EXPECT().Stop()

	claimWatcher := watcher.NewPegInClaimWatcher(runner, watchRepo, wallet, ticker)
	go claimWatcher.Start()
	ticks <- time.Now()
	go claimWatcher.Shutdown(make(chan bool, 1))

	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		assertTick(collect, watchRepo, wallet)
	}, time.Second, 10*time.Millisecond)
}

func TestPegInClaimWatcher_EmptyWalletHistoryCreatesZeroClaims(t *testing.T) {
	runner := &countingClaimRunner{}
	runClaimWatcherTick(t, runner, func(watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		watchRepo.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
		wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{}, nil).Once()
	}, func(collect *assert.CollectT, watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		assert.Equal(collect, int32(0), runner.runs.Load())
		wallet.AssertExpectations(newMockCollectT(collect))
		watchRepo.AssertExpectations(newMockCollectT(collect))
	})
}

func TestPegInClaimWatcher_DiscoveredRowIsIgnored(t *testing.T) {
	runner := &countingClaimRunner{}
	discovered := importedWatchEntry()
	discovered.State = rootstock.PegInWatchDiscovered
	runClaimWatcherTick(t, runner, func(watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		watchRepo.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{discovered}, nil).Once()
	}, func(collect *assert.CollectT, _ *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		assert.Equal(collect, int32(0), runner.runs.Load())
		wallet.AssertNotCalled(newMockCollectT(collect), "GetTransactions", mock.Anything)
	})
}

func TestPegInClaimWatcher_ListErrorCreatesZeroClaims(t *testing.T) {
	runner := &countingClaimRunner{}
	runClaimWatcherTick(t, runner, func(watchRepo *mocks.PegInWatchRepositoryMock, _ *mocks.BitcoinWalletMock) {
		watchRepo.EXPECT().List(mock.Anything).Return(nil, assert.AnError).Once()
	}, func(collect *assert.CollectT, _ *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		assert.Equal(collect, int32(0), runner.runs.Load())
		wallet.AssertNotCalled(newMockCollectT(collect), "GetTransactions", mock.Anything)
	})
}

func TestPegInClaimWatcher_WalletErrorCreatesZeroClaims(t *testing.T) {
	runner := &countingClaimRunner{}
	runClaimWatcherTick(t, runner, func(watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		watchRepo.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
		wallet.EXPECT().GetTransactions("bcrt1qimported").Return(nil, assert.AnError).Once()
	}, func(collect *assert.CollectT, _ *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		assert.Equal(collect, int32(0), runner.runs.Load())
		wallet.AssertExpectations(newMockCollectT(collect))
	})
}

func TestPegInClaimWatcher_ZeroFirstOutputIsSkipped(t *testing.T) {
	runner := &countingClaimRunner{}
	runClaimWatcherTick(t, runner, func(watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		watchRepo.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
		wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{{
			Hash: "aabbcc",
			Outputs: map[string][]*entities.Wei{
				"other": {entities.NewWei(1)},
			},
		}}, nil).Once()
	}, func(collect *assert.CollectT, _ *mocks.PegInWatchRepositoryMock, _ *mocks.BitcoinWalletMock) {
		assert.Equal(collect, int32(0), runner.runs.Load())
	})
}

func TestPegInClaimWatcher_PayingTransactionRunsClaim(t *testing.T) {
	runner := &countingClaimRunner{}
	runClaimWatcherTick(t, runner, func(watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		watchRepo.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
		wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{
			payingWatchTx("bcrt1qimported"),
		}, nil).Once()
	}, func(collect *assert.CollectT, _ *mocks.PegInWatchRepositoryMock, _ *mocks.BitcoinWalletMock) {
		assert.Equal(collect, int32(1), runner.runs.Load())
	})
}

func TestPegInClaimWatcher_RunErrorDoesNotStopTick(t *testing.T) {
	runner := &errorClaimRunner{}
	runClaimWatcherTick(t, runner, func(watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		watchRepo.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
		wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{
			payingWatchTx("bcrt1qimported"),
		}, nil).Once()
	}, func(collect *assert.CollectT, _ *mocks.PegInWatchRepositoryMock, _ *mocks.BitcoinWalletMock) {
		assert.Equal(collect, int32(1), runner.runs.Load())
	})
}

func TestPegInClaimWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) watcher.Watcher {
		return watcher.NewPegInClaimWatcher(
			&countingClaimRunner{},
			mocks.NewPegInWatchRepositoryMock(t),
			mocks.NewBitcoinWalletMock(t),
			ticker,
		)
	})
}

func TestPegInClaimWatcher_Prepare(t *testing.T) {
	runner := &recordingClaimRunner{}
	claimWatcher := watcher.NewPegInClaimWatcher(
		runner,
		mocks.NewPegInWatchRepositoryMock(t),
		mocks.NewBitcoinWalletMock(t),
		&mocks.TickerMock{},
	)
	require.NoError(t, claimWatcher.Prepare(context.Background()))
	assert.Equal(t, int32(1), runner.reconciles.Load())
	assert.Equal(t, int32(0), runner.runs.Load())
}
