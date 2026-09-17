package watcher_test

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func depositTxIDArg(args mock.Arguments) string {
	claim, ok := args.Get(1).(rootstock.PegInClaim)
	if !ok {
		return ""
	}
	return claim.DepositTxID
}

func submittingClaimRow(depositTxID string) rootstock.PegInClaim {
	return rootstock.PegInClaim{
		RskAddress:  test.AnyRskAddress,
		DepositTxID: depositTxID,
		State:       rootstock.PegInClaimSubmitting,
	}
}

func newClaimWatcherForTest(
	t *testing.T,
	runner watcher.PegInClaimRunner,
	settler watcher.PegInClaimSettler,
	claims *mocks.PegInClaimRepositoryMock,
	watches *mocks.PegInWatchRepositoryMock,
	wallet *mocks.BitcoinWalletMock,
	ticker utils.Ticker,
) *watcher.PegInClaimWatcher {
	t.Helper()
	return watcher.NewPegInClaimWatcher(runner, settler, claims, watches, wallet, ticker)
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
	runner *mocks.PegInClaimRunnerMock,
	setup func(*mocks.PegInWatchRepositoryMock, *mocks.BitcoinWalletMock),
	assertTick func(*assert.CollectT, *mocks.PegInWatchRepositoryMock, *mocks.BitcoinWalletMock, *mocks.PegInClaimRunnerMock),
) {
	t.Helper()
	watchRepo := mocks.NewPegInWatchRepositoryMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	claims := mocks.NewPegInClaimRepositoryMock(t)
	settler := mocks.NewPegInClaimSettlerMock(t)
	claims.EXPECT().ListByStates(mock.Anything, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{}, nil).Once()
	setup(watchRepo, wallet)

	ticker := &mocks.TickerMock{}
	ticks := make(chan time.Time)
	ticker.EXPECT().C().Return(ticks)
	ticker.EXPECT().Stop()

	claimWatcher := watcher.NewPegInClaimWatcher(
		runner,
		settler,
		claims,
		watchRepo,
		wallet,
		ticker,
	)
	go claimWatcher.Start()
	ticks <- time.Now()
	go claimWatcher.Shutdown(make(chan bool, 1))

	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		assertTick(collect, watchRepo, wallet, runner)
	}, time.Second, 10*time.Millisecond)
}

func TestPegInClaimWatcher_EmptyWalletHistoryCreatesZeroClaims(t *testing.T) {
	runner := mocks.NewPegInClaimRunnerMock(t)
	runClaimWatcherTick(t, runner, func(watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		watchRepo.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
		wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{}, nil).Once()
	}, func(collect *assert.CollectT, watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock, runner *mocks.PegInClaimRunnerMock) {
		mt := newMockCollectT(collect)
		wallet.AssertNumberOfCalls(mt, "GetTransactions", 1)
		runner.AssertNotCalled(mt, "Run", mock.Anything, mock.Anything, mock.Anything)
		watchRepo.AssertExpectations(mt)
	})
}

func TestPegInClaimWatcher_DiscoveredRowIsIgnored(t *testing.T) {
	runner := mocks.NewPegInClaimRunnerMock(t)
	discovered := importedWatchEntry()
	discovered.State = rootstock.PegInWatchDiscovered
	runClaimWatcherTick(t, runner, func(watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		watchRepo.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{discovered}, nil).Once()
	}, func(collect *assert.CollectT, watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock, runner *mocks.PegInClaimRunnerMock) {
		mt := newMockCollectT(collect)
		watchRepo.AssertNumberOfCalls(mt, "List", 1)
		wallet.AssertNotCalled(mt, "GetTransactions", mock.Anything)
		runner.AssertNotCalled(mt, "Run", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestPegInClaimWatcher_ListErrorCreatesZeroClaims(t *testing.T) {
	runner := mocks.NewPegInClaimRunnerMock(t)
	runClaimWatcherTick(t, runner, func(watchRepo *mocks.PegInWatchRepositoryMock, _ *mocks.BitcoinWalletMock) {
		watchRepo.EXPECT().List(mock.Anything).Return(nil, assert.AnError).Once()
	}, func(collect *assert.CollectT, watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock, runner *mocks.PegInClaimRunnerMock) {
		mt := newMockCollectT(collect)
		watchRepo.AssertNumberOfCalls(mt, "List", 1)
		wallet.AssertNotCalled(mt, "GetTransactions", mock.Anything)
		runner.AssertNotCalled(mt, "Run", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestPegInClaimWatcher_WalletErrorCreatesZeroClaims(t *testing.T) {
	runner := mocks.NewPegInClaimRunnerMock(t)
	runClaimWatcherTick(t, runner, func(watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		watchRepo.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
		wallet.EXPECT().GetTransactions("bcrt1qimported").Return(nil, assert.AnError).Once()
	}, func(collect *assert.CollectT, _ *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock, runner *mocks.PegInClaimRunnerMock) {
		mt := newMockCollectT(collect)
		wallet.AssertNumberOfCalls(mt, "GetTransactions", 1)
		runner.AssertNotCalled(mt, "Run", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestPegInClaimWatcher_ZeroFirstOutputIsSkipped(t *testing.T) {
	runner := mocks.NewPegInClaimRunnerMock(t)
	runClaimWatcherTick(t, runner, func(watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		watchRepo.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
		wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{{
			Hash: "aabbcc",
			Outputs: map[string][]*entities.Wei{
				"other": {entities.NewWei(1)},
			},
		}}, nil).Once()
	}, func(collect *assert.CollectT, _ *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock, runner *mocks.PegInClaimRunnerMock) {
		mt := newMockCollectT(collect)
		wallet.AssertNumberOfCalls(mt, "GetTransactions", 1)
		runner.AssertNotCalled(mt, "Run", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestPegInClaimWatcher_PayingTransactionRunsClaim(t *testing.T) {
	runner := mocks.NewPegInClaimRunnerMock(t)
	runner.On("Run", mock.Anything, mock.Anything, "aabbcc").Return(nil).Once()
	runClaimWatcherTick(t, runner, func(watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		watchRepo.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
		wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{
			payingWatchTx("bcrt1qimported"),
		}, nil).Once()
	}, func(collect *assert.CollectT, _ *mocks.PegInWatchRepositoryMock, _ *mocks.BitcoinWalletMock, runner *mocks.PegInClaimRunnerMock) {
		runner.AssertNumberOfCalls(newMockCollectT(collect), "Run", 1)
	})
}

func TestPegInClaimWatcher_RunErrorDoesNotStopTick(t *testing.T) {
	runner := mocks.NewPegInClaimRunnerMock(t)
	runner.On("Run", mock.Anything, mock.Anything, "aabbcc").Return(assert.AnError).Once()
	runClaimWatcherTick(t, runner, func(watchRepo *mocks.PegInWatchRepositoryMock, wallet *mocks.BitcoinWalletMock) {
		watchRepo.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
		wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{
			payingWatchTx("bcrt1qimported"),
		}, nil).Once()
	}, func(collect *assert.CollectT, _ *mocks.PegInWatchRepositoryMock, _ *mocks.BitcoinWalletMock, runner *mocks.PegInClaimRunnerMock) {
		runner.AssertNumberOfCalls(newMockCollectT(collect), "Run", 1)
	})
}

func TestPegInClaimWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) watcher.Watcher {
		return watcher.NewPegInClaimWatcher(
			mocks.NewPegInClaimRunnerMock(t),
			mocks.NewPegInClaimSettlerMock(t),
			mocks.NewPegInClaimRepositoryMock(t),
			mocks.NewPegInWatchRepositoryMock(t),
			mocks.NewBitcoinWalletMock(t),
			ticker,
		)
	})
}

func TestPegInClaimWatcher_Prepare(t *testing.T) {
	runner := mocks.NewPegInClaimRunnerMock(t)
	settler := mocks.NewPegInClaimSettlerMock(t)
	settler.On("Run", mock.Anything, mock.Anything).Return(nil).Once()
	claims := mocks.NewPegInClaimRepositoryMock(t)
	claims.EXPECT().ListByStates(mock.Anything, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{{
			RskAddress:  test.AnyRskAddress,
			DepositTxID: "aa",
			State:       rootstock.PegInClaimSubmitting,
		}}, nil).Once()
	claimWatcher := watcher.NewPegInClaimWatcher(
		runner,
		settler,
		claims,
		mocks.NewPegInWatchRepositoryMock(t),
		mocks.NewBitcoinWalletMock(t),
		&mocks.TickerMock{},
	)
	require.NoError(t, claimWatcher.Prepare(context.Background()))
	settler.AssertNumberOfCalls(t, "Run", 1)
	runner.AssertNotCalled(t, "Run", mock.Anything, mock.Anything, mock.Anything)
}

func TestPegInClaimWatcher_PrepareContinuesAfterRowSpecificError(t *testing.T) {
	settler := mocks.NewPegInClaimSettlerMock(t)
	var ids []string
	settler.On("Run", mock.Anything, submittingClaimRow("aa")).
		Run(func(args mock.Arguments) {
			ids = append(ids, depositTxIDArg(args))
		}).
		Return(assert.AnError).Once()
	settler.On("Run", mock.Anything, submittingClaimRow("bb")).
		Run(func(args mock.Arguments) {
			ids = append(ids, depositTxIDArg(args))
		}).
		Return(nil).Once()
	claims := mocks.NewPegInClaimRepositoryMock(t)
	claims.EXPECT().ListByStates(mock.Anything, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{submittingClaimRow("aa"), submittingClaimRow("bb")}, nil).Once()
	claimWatcher := newClaimWatcherForTest(
		t,
		mocks.NewPegInClaimRunnerMock(t),
		settler,
		claims,
		mocks.NewPegInWatchRepositoryMock(t),
		mocks.NewBitcoinWalletMock(t),
		&mocks.TickerMock{},
	)
	require.NoError(t, claimWatcher.Prepare(context.Background()))
	assert.Equal(t, []string{"aa", "bb"}, ids)
}

func TestPegInClaimWatcher_PrepareAbortsOnInfrastructureUnavailable(t *testing.T) {
	settler := mocks.NewPegInClaimSettlerMock(t)
	var ids []string
	settler.On("Run", mock.Anything, submittingClaimRow("aa")).
		Run(func(args mock.Arguments) {
			ids = append(ids, depositTxIDArg(args))
		}).
		Return(errors.Join(assert.AnError, usecases.InfrastructureUnavailableError)).Once()
	claims := mocks.NewPegInClaimRepositoryMock(t)
	claims.EXPECT().ListByStates(mock.Anything, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{submittingClaimRow("aa"), submittingClaimRow("bb")}, nil).Once()
	claimWatcher := newClaimWatcherForTest(
		t,
		mocks.NewPegInClaimRunnerMock(t),
		settler,
		claims,
		mocks.NewPegInWatchRepositoryMock(t),
		mocks.NewBitcoinWalletMock(t),
		&mocks.TickerMock{},
	)
	require.NoError(t, claimWatcher.Prepare(context.Background()))
	assert.Equal(t, []string{"aa"}, ids)
	settler.AssertNotCalled(t, "Run", mock.Anything, submittingClaimRow("bb"))
}

func TestPegInClaimWatcher_TickSettlesBeforeNewClaims(t *testing.T) {
	var mu sync.Mutex
	var steps []string
	record := func(step string) {
		mu.Lock()
		steps = append(steps, step)
		mu.Unlock()
	}
	settler := mocks.NewPegInClaimSettlerMock(t)
	settler.On("Run", mock.Anything, mock.Anything).Return(nil).Once()
	runner := mocks.NewPegInClaimRunnerMock(t)
	claims := mocks.NewPegInClaimRepositoryMock(t)
	watches := mocks.NewPegInWatchRepositoryMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	claims.EXPECT().ListByStates(mock.Anything, rootstock.PegInClaimSubmitting).
		Run(func(_ context.Context, _ ...rootstock.PegInClaimState) { record("list-submitting") }).
		Return([]rootstock.PegInClaim{submittingClaimRow("aa")}, nil).Once()
	watches.EXPECT().List(mock.Anything).
		Run(func(_ context.Context) { record("list-watches") }).
		Return([]rootstock.PegInWatch{}, nil).Once()

	ticker := &mocks.TickerMock{}
	ticks := make(chan time.Time)
	ticker.EXPECT().C().Return(ticks)
	ticker.EXPECT().Stop()
	claimWatcher := newClaimWatcherForTest(t, runner, settler, claims, watches, wallet, ticker)
	go claimWatcher.Start()
	ticks <- time.Now()
	go claimWatcher.Shutdown(make(chan bool, 1))

	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		settler.AssertNumberOfCalls(newMockCollectT(collect), "Run", 1)
		mu.Lock()
		defer mu.Unlock()
		assert.Equal(collect, []string{"list-submitting", "list-watches"}, steps)
	}, time.Second, 10*time.Millisecond)
}

func TestPegInClaimWatcher_TickContinuesAfterRowSpecificError(t *testing.T) {
	settler := mocks.NewPegInClaimSettlerMock(t)
	var mu sync.Mutex
	var ids []string
	settler.On("Run", mock.Anything, submittingClaimRow("aa")).
		Run(func(args mock.Arguments) {
			mu.Lock()
			ids = append(ids, depositTxIDArg(args))
			mu.Unlock()
		}).
		Return(assert.AnError).Once()
	settler.On("Run", mock.Anything, submittingClaimRow("bb")).
		Run(func(args mock.Arguments) {
			mu.Lock()
			ids = append(ids, depositTxIDArg(args))
			mu.Unlock()
		}).
		Return(nil).Once()
	runner := mocks.NewPegInClaimRunnerMock(t)
	runner.On("Run", mock.Anything, mock.Anything, "aabbcc").Return(nil).Once()
	claims := mocks.NewPegInClaimRepositoryMock(t)
	watches := mocks.NewPegInWatchRepositoryMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	claims.EXPECT().ListByStates(mock.Anything, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{submittingClaimRow("aa"), submittingClaimRow("bb")}, nil).Once()
	watches.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
	wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{
		payingWatchTx("bcrt1qimported"),
	}, nil).Once()

	ticker := &mocks.TickerMock{}
	ticks := make(chan time.Time)
	ticker.EXPECT().C().Return(ticks)
	ticker.EXPECT().Stop()
	claimWatcher := newClaimWatcherForTest(t, runner, settler, claims, watches, wallet, ticker)
	go claimWatcher.Start()
	ticks <- time.Now()
	go claimWatcher.Shutdown(make(chan bool, 1))

	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		mu.Lock()
		defer mu.Unlock()
		assert.Equal(collect, []string{"aa", "bb"}, ids)
		runner.AssertNumberOfCalls(newMockCollectT(collect), "Run", 1)
	}, time.Second, 10*time.Millisecond)
}

func TestPegInClaimWatcher_TickAbortsOnInfrastructureUnavailable(t *testing.T) {
	settler := mocks.NewPegInClaimSettlerMock(t)
	var mu sync.Mutex
	var ids []string
	settler.On("Run", mock.Anything, submittingClaimRow("aa")).
		Run(func(args mock.Arguments) {
			mu.Lock()
			ids = append(ids, depositTxIDArg(args))
			mu.Unlock()
		}).
		Return(errors.Join(assert.AnError, usecases.InfrastructureUnavailableError)).Once()
	runner := mocks.NewPegInClaimRunnerMock(t)
	claims := mocks.NewPegInClaimRepositoryMock(t)
	watches := mocks.NewPegInWatchRepositoryMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	claims.EXPECT().ListByStates(mock.Anything, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{submittingClaimRow("aa"), submittingClaimRow("bb")}, nil).Once()

	ticker := &mocks.TickerMock{}
	ticks := make(chan time.Time)
	ticker.EXPECT().C().Return(ticks)
	ticker.EXPECT().Stop()
	claimWatcher := newClaimWatcherForTest(t, runner, settler, claims, watches, wallet, ticker)
	go claimWatcher.Start()
	ticks <- time.Now()
	go claimWatcher.Shutdown(make(chan bool, 1))

	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		mu.Lock()
		defer mu.Unlock()
		assert.Equal(collect, []string{"aa"}, ids)
		runner.AssertNotCalled(newMockCollectT(collect), "Run", mock.Anything, mock.Anything, mock.Anything)
		watches.AssertNotCalled(newMockCollectT(collect), "List", mock.Anything)
		wallet.AssertNotCalled(newMockCollectT(collect), "GetTransactions", mock.Anything)
	}, time.Second, 10*time.Millisecond)
}

func TestPegInClaimWatcher_StopChannelIsStruct(t *testing.T) {
	claimWatcher := newClaimWatcherForTest(
		t,
		mocks.NewPegInClaimRunnerMock(t),
		mocks.NewPegInClaimSettlerMock(t),
		mocks.NewPegInClaimRepositoryMock(t),
		mocks.NewPegInWatchRepositoryMock(t),
		mocks.NewBitcoinWalletMock(t),
		&mocks.TickerMock{},
	)
	field, ok := reflect.TypeOf(claimWatcher).Elem().FieldByName("watcherStopChannel")
	require.True(t, ok)
	assert.Equal(t, reflect.ChanOf(reflect.BothDir, reflect.TypeOf(struct{}{})), field.Type)
}

func TestPegInClaimWatcher_ShutdownSendsTrueOnCloseChannel(t *testing.T) {
	ticker := &mocks.TickerMock{}
	ticker.EXPECT().C().Return(make(chan time.Time))
	ticker.EXPECT().Stop()
	claimWatcher := newClaimWatcherForTest(
		t,
		mocks.NewPegInClaimRunnerMock(t),
		mocks.NewPegInClaimSettlerMock(t),
		mocks.NewPegInClaimRepositoryMock(t),
		mocks.NewPegInWatchRepositoryMock(t),
		mocks.NewBitcoinWalletMock(t),
		ticker,
	)
	closeChannel := make(chan bool, 1)
	go claimWatcher.Start()
	claimWatcher.Shutdown(closeChannel)
	assert.True(t, <-closeChannel)
}
