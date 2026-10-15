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

func submittingClaim(depositTxID string) rootstock.PegInClaim {
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
	submittingClaims *mocks.PegInClaimsGetterMock,
	importedWatches *mocks.PegInWatchesGetterMock,
	wallet *mocks.BitcoinWalletMock,
	ticker utils.Ticker,
) *watcher.PegInClaimWatcher {
	t.Helper()
	return watcher.NewPegInClaimWatcher(
		watcher.NewPegInClaimWatcherUseCases(runner, settler, submittingClaims, importedWatches),
		wallet,
		ticker,
	)
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
	setup func(*mocks.PegInWatchesGetterMock, *mocks.BitcoinWalletMock),
	assertTick func(*assert.CollectT, *mocks.PegInWatchesGetterMock, *mocks.BitcoinWalletMock, *mocks.PegInClaimRunnerMock),
) {
	t.Helper()
	importedWatches := mocks.NewPegInWatchesGetterMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	submittingClaims := mocks.NewPegInClaimsGetterMock(t)
	settler := mocks.NewPegInClaimSettlerMock(t)
	submittingClaims.EXPECT().Run(mock.Anything, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{}, nil).Once()
	setup(importedWatches, wallet)

	ticker := &mocks.TickerMock{}
	ticks := make(chan time.Time)
	ticker.EXPECT().C().Return(ticks)
	ticker.EXPECT().Stop()

	claimWatcher := newClaimWatcherForTest(
		t,
		runner,
		settler,
		submittingClaims,
		importedWatches,
		wallet,
		ticker,
	)
	go claimWatcher.Start()
	ticks <- time.Now()
	go claimWatcher.Shutdown(make(chan bool, 1))

	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		assertTick(collect, importedWatches, wallet, runner)
	}, time.Second, 10*time.Millisecond)
}

func TestPegInClaimWatcher_EmptyWalletHistoryCreatesZeroClaims(t *testing.T) {
	runner := mocks.NewPegInClaimRunnerMock(t)
	runClaimWatcherTick(t, runner, func(importedWatches *mocks.PegInWatchesGetterMock, wallet *mocks.BitcoinWalletMock) {
		importedWatches.EXPECT().Run(mock.Anything, rootstock.PegInWatchImported).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
		wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{}, nil).Once()
	}, func(collect *assert.CollectT, importedWatches *mocks.PegInWatchesGetterMock, wallet *mocks.BitcoinWalletMock, runner *mocks.PegInClaimRunnerMock) {
		mt := newMockCollectT(collect)
		wallet.AssertNumberOfCalls(mt, "GetTransactions", 1)
		runner.AssertNotCalled(mt, "Run", mock.Anything, mock.Anything, mock.Anything)
		importedWatches.AssertExpectations(mt)
	})
}

func TestPegInClaimWatcher_ListErrorCreatesZeroClaims(t *testing.T) {
	runner := mocks.NewPegInClaimRunnerMock(t)
	runClaimWatcherTick(t, runner, func(importedWatches *mocks.PegInWatchesGetterMock, _ *mocks.BitcoinWalletMock) {
		importedWatches.EXPECT().Run(mock.Anything, rootstock.PegInWatchImported).Return(nil, assert.AnError).Once()
	}, func(collect *assert.CollectT, importedWatches *mocks.PegInWatchesGetterMock, wallet *mocks.BitcoinWalletMock, runner *mocks.PegInClaimRunnerMock) {
		mt := newMockCollectT(collect)
		importedWatches.AssertNumberOfCalls(mt, "Run", 1)
		wallet.AssertNotCalled(mt, "GetTransactions", mock.Anything)
		runner.AssertNotCalled(mt, "Run", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestPegInClaimWatcher_WalletErrorCreatesZeroClaims(t *testing.T) {
	runner := mocks.NewPegInClaimRunnerMock(t)
	runClaimWatcherTick(t, runner, func(importedWatches *mocks.PegInWatchesGetterMock, wallet *mocks.BitcoinWalletMock) {
		importedWatches.EXPECT().Run(mock.Anything, rootstock.PegInWatchImported).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
		wallet.EXPECT().GetTransactions("bcrt1qimported").Return(nil, assert.AnError).Once()
	}, func(collect *assert.CollectT, _ *mocks.PegInWatchesGetterMock, wallet *mocks.BitcoinWalletMock, runner *mocks.PegInClaimRunnerMock) {
		mt := newMockCollectT(collect)
		wallet.AssertNumberOfCalls(mt, "GetTransactions", 1)
		runner.AssertNotCalled(mt, "Run", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestPegInClaimWatcher_ZeroFirstOutputIsSkipped(t *testing.T) {
	runner := mocks.NewPegInClaimRunnerMock(t)
	runClaimWatcherTick(t, runner, func(importedWatches *mocks.PegInWatchesGetterMock, wallet *mocks.BitcoinWalletMock) {
		importedWatches.EXPECT().Run(mock.Anything, rootstock.PegInWatchImported).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
		wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{{
			Hash: "aabbcc",
			Outputs: map[string][]*entities.Wei{
				"other": {entities.NewWei(1)},
			},
		}}, nil).Once()
	}, func(collect *assert.CollectT, _ *mocks.PegInWatchesGetterMock, wallet *mocks.BitcoinWalletMock, runner *mocks.PegInClaimRunnerMock) {
		mt := newMockCollectT(collect)
		wallet.AssertNumberOfCalls(mt, "GetTransactions", 1)
		runner.AssertNotCalled(mt, "Run", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestPegInClaimWatcher_PayingTransactionRunsClaim(t *testing.T) {
	runner := mocks.NewPegInClaimRunnerMock(t)
	runner.On("Run", mock.Anything, mock.Anything, "aabbcc").Return(nil).Once()
	runClaimWatcherTick(t, runner, func(importedWatches *mocks.PegInWatchesGetterMock, wallet *mocks.BitcoinWalletMock) {
		importedWatches.EXPECT().Run(mock.Anything, rootstock.PegInWatchImported).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
		wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{
			payingWatchTx("bcrt1qimported"),
		}, nil).Once()
	}, func(collect *assert.CollectT, _ *mocks.PegInWatchesGetterMock, _ *mocks.BitcoinWalletMock, runner *mocks.PegInClaimRunnerMock) {
		runner.AssertNumberOfCalls(newMockCollectT(collect), "Run", 1)
	})
}

func TestPegInClaimWatcher_RunErrorDoesNotStopTick(t *testing.T) {
	runner := mocks.NewPegInClaimRunnerMock(t)
	runner.On("Run", mock.Anything, mock.Anything, "aabbcc").Return(assert.AnError).Once()
	runClaimWatcherTick(t, runner, func(importedWatches *mocks.PegInWatchesGetterMock, wallet *mocks.BitcoinWalletMock) {
		importedWatches.EXPECT().Run(mock.Anything, rootstock.PegInWatchImported).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
		wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{
			payingWatchTx("bcrt1qimported"),
		}, nil).Once()
	}, func(collect *assert.CollectT, _ *mocks.PegInWatchesGetterMock, _ *mocks.BitcoinWalletMock, runner *mocks.PegInClaimRunnerMock) {
		runner.AssertNumberOfCalls(newMockCollectT(collect), "Run", 1)
	})
}

func TestPegInClaimWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) watcher.Watcher {
		return newClaimWatcherForTest(
			t,
			mocks.NewPegInClaimRunnerMock(t),
			mocks.NewPegInClaimSettlerMock(t),
			mocks.NewPegInClaimsGetterMock(t),
			mocks.NewPegInWatchesGetterMock(t),
			mocks.NewBitcoinWalletMock(t),
			ticker,
		)
	})
}

func TestPegInClaimWatcher_Prepare(t *testing.T) {
	runner := mocks.NewPegInClaimRunnerMock(t)
	settler := mocks.NewPegInClaimSettlerMock(t)
	settler.On("Run", mock.Anything, mock.Anything).Return(nil).Once()
	submittingClaims := mocks.NewPegInClaimsGetterMock(t)
	submittingClaims.EXPECT().Run(mock.Anything, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{{
			RskAddress:  test.AnyRskAddress,
			DepositTxID: "aa",
			State:       rootstock.PegInClaimSubmitting,
		}}, nil).Once()
	claimWatcher := newClaimWatcherForTest(
		t,
		runner,
		settler,
		submittingClaims,
		mocks.NewPegInWatchesGetterMock(t),
		mocks.NewBitcoinWalletMock(t),
		&mocks.TickerMock{},
	)
	require.NoError(t, claimWatcher.Prepare(context.Background()))
	settler.AssertNumberOfCalls(t, "Run", 1)
	runner.AssertNotCalled(t, "Run", mock.Anything, mock.Anything, mock.Anything)
}

func TestPegInClaimWatcher_PrepareContinuesAfterClaimSpecificError(t *testing.T) {
	settler := mocks.NewPegInClaimSettlerMock(t)
	var ids []string
	settler.On("Run", mock.Anything, submittingClaim("aa")).
		Run(func(args mock.Arguments) {
			ids = append(ids, depositTxIDArg(args))
		}).
		Return(assert.AnError).Once()
	settler.On("Run", mock.Anything, submittingClaim("bb")).
		Run(func(args mock.Arguments) {
			ids = append(ids, depositTxIDArg(args))
		}).
		Return(nil).Once()
	submittingClaims := mocks.NewPegInClaimsGetterMock(t)
	submittingClaims.EXPECT().Run(mock.Anything, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{submittingClaim("aa"), submittingClaim("bb")}, nil).Once()
	claimWatcher := newClaimWatcherForTest(
		t,
		mocks.NewPegInClaimRunnerMock(t),
		settler,
		submittingClaims,
		mocks.NewPegInWatchesGetterMock(t),
		mocks.NewBitcoinWalletMock(t),
		&mocks.TickerMock{},
	)
	require.NoError(t, claimWatcher.Prepare(context.Background()))
	assert.Equal(t, []string{"aa", "bb"}, ids)
}

func TestPegInClaimWatcher_PrepareAbortsOnInfrastructureUnavailable(t *testing.T) {
	settler := mocks.NewPegInClaimSettlerMock(t)
	var ids []string
	settler.On("Run", mock.Anything, submittingClaim("aa")).
		Run(func(args mock.Arguments) {
			ids = append(ids, depositTxIDArg(args))
		}).
		Return(errors.Join(assert.AnError, usecases.InfrastructureUnavailableError)).Once()
	submittingClaims := mocks.NewPegInClaimsGetterMock(t)
	submittingClaims.EXPECT().Run(mock.Anything, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{submittingClaim("aa"), submittingClaim("bb")}, nil).Once()
	claimWatcher := newClaimWatcherForTest(
		t,
		mocks.NewPegInClaimRunnerMock(t),
		settler,
		submittingClaims,
		mocks.NewPegInWatchesGetterMock(t),
		mocks.NewBitcoinWalletMock(t),
		&mocks.TickerMock{},
	)
	require.NoError(t, claimWatcher.Prepare(context.Background()))
	assert.Equal(t, []string{"aa"}, ids)
	settler.AssertNotCalled(t, "Run", mock.Anything, submittingClaim("bb"))
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
	submittingClaims := mocks.NewPegInClaimsGetterMock(t)
	importedWatches := mocks.NewPegInWatchesGetterMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	submittingClaims.EXPECT().Run(mock.Anything, rootstock.PegInClaimSubmitting).
		Run(func(_ context.Context, _ ...rootstock.PegInClaimState) { record("list-submitting") }).
		Return([]rootstock.PegInClaim{submittingClaim("aa")}, nil).Once()
	importedWatches.EXPECT().Run(mock.Anything, rootstock.PegInWatchImported).
		Run(func(_ context.Context, _ ...rootstock.PegInWatchState) { record("list-watches") }).
		Return([]rootstock.PegInWatch{}, nil).Once()

	ticker := &mocks.TickerMock{}
	ticks := make(chan time.Time)
	ticker.EXPECT().C().Return(ticks)
	ticker.EXPECT().Stop()
	claimWatcher := newClaimWatcherForTest(t, runner, settler, submittingClaims, importedWatches, wallet, ticker)
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

func TestPegInClaimWatcher_TickContinuesAfterClaimSpecificError(t *testing.T) {
	settler := mocks.NewPegInClaimSettlerMock(t)
	var mu sync.Mutex
	var ids []string
	settler.On("Run", mock.Anything, submittingClaim("aa")).
		Run(func(args mock.Arguments) {
			mu.Lock()
			ids = append(ids, depositTxIDArg(args))
			mu.Unlock()
		}).
		Return(assert.AnError).Once()
	settler.On("Run", mock.Anything, submittingClaim("bb")).
		Run(func(args mock.Arguments) {
			mu.Lock()
			ids = append(ids, depositTxIDArg(args))
			mu.Unlock()
		}).
		Return(nil).Once()
	runner := mocks.NewPegInClaimRunnerMock(t)
	runner.On("Run", mock.Anything, mock.Anything, "aabbcc").Return(nil).Once()
	submittingClaims := mocks.NewPegInClaimsGetterMock(t)
	importedWatches := mocks.NewPegInWatchesGetterMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	submittingClaims.EXPECT().Run(mock.Anything, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{submittingClaim("aa"), submittingClaim("bb")}, nil).Once()
	importedWatches.EXPECT().Run(mock.Anything, rootstock.PegInWatchImported).Return([]rootstock.PegInWatch{importedWatchEntry()}, nil).Once()
	wallet.EXPECT().GetTransactions("bcrt1qimported").Return([]blockchain.BitcoinTransactionInformation{
		payingWatchTx("bcrt1qimported"),
	}, nil).Once()

	ticker := &mocks.TickerMock{}
	ticks := make(chan time.Time)
	ticker.EXPECT().C().Return(ticks)
	ticker.EXPECT().Stop()
	claimWatcher := newClaimWatcherForTest(t, runner, settler, submittingClaims, importedWatches, wallet, ticker)
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
	settler.On("Run", mock.Anything, submittingClaim("aa")).
		Run(func(args mock.Arguments) {
			mu.Lock()
			ids = append(ids, depositTxIDArg(args))
			mu.Unlock()
		}).
		Return(errors.Join(assert.AnError, usecases.InfrastructureUnavailableError)).Once()
	runner := mocks.NewPegInClaimRunnerMock(t)
	submittingClaims := mocks.NewPegInClaimsGetterMock(t)
	importedWatches := mocks.NewPegInWatchesGetterMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	submittingClaims.EXPECT().Run(mock.Anything, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{submittingClaim("aa"), submittingClaim("bb")}, nil).Once()

	ticker := &mocks.TickerMock{}
	ticks := make(chan time.Time)
	ticker.EXPECT().C().Return(ticks)
	ticker.EXPECT().Stop()
	claimWatcher := newClaimWatcherForTest(t, runner, settler, submittingClaims, importedWatches, wallet, ticker)
	go claimWatcher.Start()
	ticks <- time.Now()
	go claimWatcher.Shutdown(make(chan bool, 1))

	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		mu.Lock()
		defer mu.Unlock()
		assert.Equal(collect, []string{"aa"}, ids)
		runner.AssertNotCalled(newMockCollectT(collect), "Run", mock.Anything, mock.Anything, mock.Anything)
		importedWatches.AssertNotCalled(newMockCollectT(collect), "Run", mock.Anything)
		wallet.AssertNotCalled(newMockCollectT(collect), "GetTransactions", mock.Anything)
	}, time.Second, 10*time.Millisecond)
}

func TestPegInClaimWatcher_StopChannelIsStruct(t *testing.T) {
	claimWatcher := newClaimWatcherForTest(
		t,
		mocks.NewPegInClaimRunnerMock(t),
		mocks.NewPegInClaimSettlerMock(t),
		mocks.NewPegInClaimsGetterMock(t),
		mocks.NewPegInWatchesGetterMock(t),
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
		mocks.NewPegInClaimsGetterMock(t),
		mocks.NewPegInWatchesGetterMock(t),
		mocks.NewBitcoinWalletMock(t),
		ticker,
	)
	closeChannel := make(chan bool, 1)
	go claimWatcher.Start()
	claimWatcher.Shutdown(closeChannel)
	assert.True(t, <-closeChannel)
}
