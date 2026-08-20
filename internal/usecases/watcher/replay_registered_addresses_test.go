package watcher_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type replayCheckpointRepository struct {
	*mocks.PegInAddressRegistryWatchRepositoryMock
	checkpoints mock.Mock
}

func newReplayCheckpointRepository(t *testing.T) *replayCheckpointRepository {
	t.Helper()
	repository := &replayCheckpointRepository{
		PegInAddressRegistryWatchRepositoryMock: mocks.NewPegInAddressRegistryWatchRepositoryMock(t),
	}
	repository.checkpoints.Test(t)
	t.Cleanup(func() { repository.checkpoints.AssertExpectations(t) })
	return repository
}

func (repository *replayCheckpointRepository) GetCheckpoint(
	ctx context.Context,
) (rootstock.PegInAddressRegistryWatchCheckpoint, bool, error) {
	args := repository.checkpoints.MethodCalled("GetCheckpoint", ctx)
	checkpoint, ok := args.Get(0).(rootstock.PegInAddressRegistryWatchCheckpoint)
	if !ok {
		return rootstock.PegInAddressRegistryWatchCheckpoint{}, false, errors.New("invalid checkpoint mock result")
	}
	return checkpoint, args.Bool(1), args.Error(2)
}

func (repository *replayCheckpointRepository) SetCheckpoint(
	ctx context.Context,
	checkpoint rootstock.PegInAddressRegistryWatchCheckpoint,
) error {
	return repository.checkpoints.MethodCalled("SetCheckpoint", ctx, checkpoint).Error(0)
}

func (repository *replayCheckpointRepository) DeleteCheckpoint(ctx context.Context) error {
	return repository.checkpoints.MethodCalled("DeleteCheckpoint", ctx).Error(0)
}

type replayHarness struct {
	t          *testing.T
	repository *replayCheckpointRepository
	registry   *mocks.PegInAddressRegistryContractMock
	rskRpc     *mocks.RootstockRpcServerMock
}

func newReplayHarness(t *testing.T) *replayHarness {
	t.Helper()
	return &replayHarness{
		t:          t,
		repository: newReplayCheckpointRepository(t),
		registry:   mocks.NewPegInAddressRegistryContractMock(t),
		rskRpc:     mocks.NewRootstockRpcServerMock(t),
	}
}

func (harness *replayHarness) useCase(
	finalityDepth uint64,
) *watcher.ReplayRegisteredAddressesUseCase {
	harness.t.Helper()
	return newReplayUseCase(
		harness.t,
		harness.repository,
		harness.registry,
		harness.rskRpc,
		finalityDepth,
	)
}

func newReplayUseCase(
	t *testing.T,
	repository rootstock.PegInAddressRegistryWatchRepositorySet,
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	finalityDepth uint64,
) *watcher.ReplayRegisteredAddressesUseCase {
	t.Helper()
	return watcher.NewReplayRegisteredAddressesUseCase(
		repository,
		registry,
		rskRpc,
		nil,
		mocks.NewBitcoinWalletMock(t),
		finalityDepth,
	)
}

func TestNewReplayRegisteredAddressesUseCase(t *testing.T) {
	require.NotNil(t, newReplayHarness(t).useCase(2))
}

func TestReplayRegisteredAddressesUseCase_Run_DoesNotQueryBeforeConfiguredStart(t *testing.T) {
	const (
		startBlock    = uint64(100)
		finalityDepth = uint64(2)
		head          = uint64(101)
	)
	require.Less(t, head-finalityDepth, startBlock)
	harness := newReplayHarness(t)
	harness.rskRpc.EXPECT().GetHeight(mock.Anything).Return(head, nil).Once()

	pending, err := harness.useCase(finalityDepth).Run(context.Background(), false, startBlock, 10)

	require.NoError(t, err)
	assert.Empty(t, pending)
	harness.registry.AssertNotCalled(t, "GetAddressRegisteredEvents", mock.Anything, mock.Anything, mock.Anything)
	harness.registry.AssertNotCalled(t, "GetRegistrationRoot", mock.Anything, mock.Anything)
	harness.repository.checkpoints.AssertNotCalled(t, "SetCheckpoint", mock.Anything, mock.Anything)
	harness.repository.checkpoints.AssertNotCalled(t, "DeleteCheckpoint", mock.Anything)
}

func TestReplayRegisteredAddressesUseCase_Run_WrapsErrors(t *testing.T) {
	harness := newReplayHarness(t)
	harness.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(0), assert.AnError).Once()

	_, err := harness.useCase(0).Run(context.Background(), false, 100, 10)

	require.ErrorIs(t, err, assert.AnError)
	assert.ErrorContains(t, err, string(usecases.ReplayRegisteredAddressesId))
}

func TestReplayRegisteredAddressesUseCase_Run_DeletesCheckpointWhenFinalizedHeadIsBelowStart(t *testing.T) {
	const (
		startBlock    = uint64(100)
		finalityDepth = uint64(2)
		head          = uint64(101)
	)
	require.Less(t, head-finalityDepth, startBlock)
	harness := newReplayHarness(t)
	harness.repository.checkpoints.On("DeleteCheckpoint", mock.Anything).Return(nil).Once()
	harness.rskRpc.EXPECT().GetHeight(mock.Anything).Return(head, nil).Once()

	pending, err := harness.useCase(finalityDepth).Run(context.Background(), true, startBlock, 10)

	require.NoError(t, err)
	assert.Empty(t, pending)
	harness.registry.AssertNotCalled(t, "GetAddressRegisteredEvents", mock.Anything, mock.Anything, mock.Anything)
}

func exclusiveQueryGate(
	started, release chan struct{},
	calls, active, maximum *atomic.Int64,
) func(context.Context, uint64, *uint64) ([]blockchain.AddressRegistered, error) {
	return func(context.Context, uint64, *uint64) ([]blockchain.AddressRegistered, error) {
		call := calls.Add(1)
		n := active.Add(1)
		for {
			current := maximum.Load()
			if n <= current || maximum.CompareAndSwap(current, n) {
				break
			}
		}
		defer active.Add(-1)
		if call == 1 {
			close(started)
			<-release
		}
		return nil, nil
	}
}

func TestReplayRegisteredAddressesUseCase_Run_SecondRunWaitsUntilFirstRunFinishes(t *testing.T) {
	harness := newReplayHarness(t)
	harness.repository.EXPECT().List(mock.Anything).Return(nil, nil).Twice()
	harness.repository.checkpoints.On("GetCheckpoint", mock.Anything).
		Return(rootstock.PegInAddressRegistryWatchCheckpoint{}, false, nil).
		Twice()
	harness.repository.checkpoints.On("SetCheckpoint", mock.Anything, mock.Anything).Return(nil).Twice()
	firstQueryStarted := make(chan struct{})
	releaseFirstQuery := make(chan struct{})
	var queryCalls, activeQueries, maximumActiveQueries atomic.Int64
	harness.registry.EXPECT().
		GetAddressRegisteredEvents(mock.Anything, uint64(0), mock.Anything).
		RunAndReturn(exclusiveQueryGate(firstQueryStarted, releaseFirstQuery, &queryCalls, &activeQueries, &maximumActiveQueries)).
		Twice()
	harness.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(0)).Return([32]byte{}, nil).Twice()
	harness.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(0), nil).Twice()
	useCase := harness.useCase(0)

	firstResult := make(chan error, 1)
	go func() {
		_, err := useCase.Run(context.Background(), false, 0, 1)
		firstResult <- err
	}()
	select {
	case <-firstQueryStarted:
	case <-time.After(time.Second):
		require.FailNow(t, "first replay did not reach the page query")
	}

	secondStarted := make(chan struct{})
	secondResult := make(chan error, 1)
	go func() {
		close(secondStarted)
		_, err := useCase.Run(context.Background(), false, 0, 1)
		secondResult <- err
	}()
	<-secondStarted
	assert.Never(t, func() bool {
		return queryCalls.Load() > 1
	}, 50*time.Millisecond, time.Millisecond, "a second replay entered while the first was active")

	close(releaseFirstQuery)
	require.NoError(t, <-firstResult)
	require.NoError(t, <-secondResult)
	assert.Equal(t, int64(2), queryCalls.Load())
	assert.Equal(t, int64(1), maximumActiveQueries.Load())
}

func TestReplayRegisteredAddressesUseCase_Run_DoesNotPublishIncrementalReplayBeforeRootVerification(t *testing.T) {
	rootErr := errors.New("root read failed")
	original := rootstock.PegInAddressRegistryWatchCheckpoint{LastProcessedBlock: 100}
	harness := newReplayHarness(t)
	harness.repository.EXPECT().List(mock.Anything).Return(nil, nil).Once()
	harness.repository.checkpoints.On("GetCheckpoint", mock.Anything).Return(original, true, nil).Once()
	harness.registry.EXPECT().
		GetAddressRegisteredEvents(mock.Anything, uint64(101), mock.Anything).
		Return(nil, nil).
		Once()
	harness.registry.EXPECT().
		GetRegistrationRoot(mock.Anything, uint64(101)).
		Return([32]byte{}, rootErr).
		Once()
	harness.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(101), nil).Once()

	_, err := harness.useCase(0).Run(context.Background(), false, 0, 10)

	require.ErrorIs(t, err, rootErr)
	harness.repository.checkpoints.AssertNotCalled(t, "SetCheckpoint", mock.Anything, mock.Anything)
	harness.repository.checkpoints.AssertNotCalled(t, "DeleteCheckpoint", mock.Anything)
}
