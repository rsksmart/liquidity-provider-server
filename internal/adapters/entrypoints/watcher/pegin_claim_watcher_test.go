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

func importedWatchEntry() rootstock.PegInWatch {
	return rootstock.PegInWatch{
		RskAddress: test.AnyRskAddress,
		BtcAddress: "bcrt1qimported",
		State:      rootstock.PegInWatchImported,
	}
}

func TestPegInClaimWatcher_EmptyWalletHistoryCreatesZeroClaims(t *testing.T) {
	watchRepo := mocks.NewPegInWatchRepositoryMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	runner := &countingClaimRunner{}
	watchRepo.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
	wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{}, nil).Once()

	ticker := &mocks.TickerMock{}
	ticks := make(chan time.Time)
	ticker.EXPECT().C().Return(ticks)
	ticker.EXPECT().Stop()

	claimWatcher := watcher.NewPegInClaimWatcher(runner, watchRepo, wallet, ticker)
	go claimWatcher.Start()
	ticks <- time.Now()
	go claimWatcher.Shutdown(make(chan bool, 1))

	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		assert.Equal(collect, int32(0), runner.runs.Load())
		wallet.AssertExpectations(newMockCollectT(collect))
		watchRepo.AssertExpectations(newMockCollectT(collect))
	}, time.Second, 10*time.Millisecond)
}

func TestPegInClaimWatcher_DiscoveredRowIsIgnored(t *testing.T) {
	watchRepo := mocks.NewPegInWatchRepositoryMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	runner := &countingClaimRunner{}
	discovered := importedWatchEntry()
	discovered.State = rootstock.PegInWatchDiscovered
	watchRepo.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{discovered}, nil).Once()

	ticker := &mocks.TickerMock{}
	ticks := make(chan time.Time)
	ticker.EXPECT().C().Return(ticks)
	ticker.EXPECT().Stop()

	claimWatcher := watcher.NewPegInClaimWatcher(runner, watchRepo, wallet, ticker)
	go claimWatcher.Start()
	ticks <- time.Now()
	go claimWatcher.Shutdown(make(chan bool, 1))

	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		assert.Equal(collect, int32(0), runner.runs.Load())
		wallet.AssertNotCalled(newMockCollectT(collect), "GetTransactions", mock.Anything)
	}, time.Second, 10*time.Millisecond)
}

func TestPegInClaimWatcher_PayingTransactionRunsClaim(t *testing.T) {
	watchRepo := mocks.NewPegInWatchRepositoryMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	runner := &countingClaimRunner{}
	watchRepo.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
	wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{{
		Hash: "aabbcc",
		Outputs: map[string][]*entities.Wei{
			"bcrt1qimported": {entities.NewWei(1)},
		},
	}}, nil).Once()

	ticker := &mocks.TickerMock{}
	ticks := make(chan time.Time)
	ticker.EXPECT().C().Return(ticks)
	ticker.EXPECT().Stop()

	claimWatcher := watcher.NewPegInClaimWatcher(runner, watchRepo, wallet, ticker)
	go claimWatcher.Start()
	ticks <- time.Now()
	go claimWatcher.Shutdown(make(chan bool, 1))

	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		assert.Equal(collect, int32(1), runner.runs.Load())
	}, time.Second, 10*time.Millisecond)
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
	claimWatcher := watcher.NewPegInClaimWatcher(
		&countingClaimRunner{},
		mocks.NewPegInWatchRepositoryMock(t),
		mocks.NewBitcoinWalletMock(t),
		&mocks.TickerMock{},
	)
	require.NoError(t, claimWatcher.Prepare(context.Background()))
}
