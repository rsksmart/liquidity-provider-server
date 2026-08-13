package watcher

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type overlappingReplayRepository struct {
	mutex       sync.Mutex
	checkpoint  rootstock.PegInAddressRegistryWatchCheckpoint
	found       bool
	setCalls    int
	deleteCalls int
}

func (repository *overlappingReplayRepository) Upsert(
	context.Context,
	rootstock.PegInAddressRegistryWatchEntry,
) error {
	return nil
}

func (repository *overlappingReplayRepository) Get(
	context.Context,
	string,
) (*rootstock.PegInAddressRegistryWatchEntry, error) {
	return nil, nil
}

func (repository *overlappingReplayRepository) List(
	context.Context,
) ([]rootstock.PegInAddressRegistryWatchEntry, error) {
	return nil, nil
}

func (repository *overlappingReplayRepository) Update(
	context.Context,
	rootstock.PegInAddressRegistryWatchEntry,
) error {
	return nil
}

func (repository *overlappingReplayRepository) GetCheckpoint(
	context.Context,
) (rootstock.PegInAddressRegistryWatchCheckpoint, bool, error) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()
	return repository.checkpoint, repository.found, nil
}

func (repository *overlappingReplayRepository) SetCheckpoint(
	_ context.Context,
	checkpoint rootstock.PegInAddressRegistryWatchCheckpoint,
) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()
	repository.setCalls++
	repository.checkpoint = checkpoint
	repository.found = true
	return nil
}

func (repository *overlappingReplayRepository) DeleteCheckpoint(context.Context) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()
	repository.deleteCalls++
	repository.checkpoint = rootstock.PegInAddressRegistryWatchCheckpoint{}
	repository.found = false
	return nil
}

func (repository *overlappingReplayRepository) checkpointState() (
	rootstock.PegInAddressRegistryWatchCheckpoint,
	bool,
	int,
	int,
) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()
	return repository.checkpoint, repository.found, repository.setCalls, repository.deleteCalls
}

func TestPegInAddressRegistryWatcher_TrustedCheckpointWaitsForConfiguredStart(t *testing.T) {
	original := rootstock.PegInAddressRegistryWatchCheckpoint{
		LocalRoot:          [32]byte{1},
		LastProcessedBlock: 105,
	}
	repository := &overlappingReplayRepository{checkpoint: original, found: true}
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(101), nil).Once()
	watcher := newTestPegInAddressRegistryWatcher(
		repository,
		registry,
		rskRpc,
		mocks.NewBitcoinWalletMock(t),
		100,
		10,
		2,
	)
	require.NoError(t, watcher.Prepare(context.Background()))

	require.NoError(t, watcher.scan(context.Background()))

	checkpoint, found, setCalls, deleteCalls := repository.checkpointState()
	require.True(t, found)
	assert.Equal(t, original, checkpoint)
	assert.Zero(t, setCalls)
	assert.Zero(t, deleteCalls)
	registry.AssertNotCalled(t, "GetAddressRegisteredEvents", mock.Anything, mock.Anything, mock.Anything)
	registry.AssertNotCalled(t, "GetRegistrationRoot", mock.Anything, mock.Anything)
}

func TestPegInAddressRegistryWatcher_ReorgBelowConfiguredStartRemainsCold(t *testing.T) {
	repository := &overlappingReplayRepository{
		checkpoint: rootstock.PegInAddressRegistryWatchCheckpoint{
			LocalRoot:          [32]byte{1},
			LastProcessedBlock: 105,
		},
		found: true,
	}
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(101), nil).Once()
	watcher := newTestPegInAddressRegistryWatcher(
		repository,
		registry,
		rskRpc,
		mocks.NewBitcoinWalletMock(t),
		100,
		10,
		2,
	)
	require.NoError(t, watcher.Prepare(context.Background()))

	require.NoError(t, watcher.resync(context.Background()))

	checkpoint, found, setCalls, deleteCalls := repository.checkpointState()
	assert.False(t, found)
	assert.Zero(t, checkpoint)
	assert.Zero(t, setCalls)
	assert.Equal(t, 1, deleteCalls)
	registry.AssertNotCalled(t, "GetAddressRegisteredEvents", mock.Anything, mock.Anything, mock.Anything)
	registry.AssertNotCalled(t, "GetRegistrationRoot", mock.Anything, mock.Anything)
}

func TestPegInAddressRegistryWatcher_SerializesOverlappingResyncs(t *testing.T) {
	repository := &overlappingReplayRepository{}
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)

	firstQueryStarted := make(chan struct{})
	releaseFirstQuery := make(chan struct{})
	var queryCalls atomic.Int64
	var activeQueries atomic.Int64
	var maximumActiveQueries atomic.Int64
	registry.EXPECT().
		GetAddressRegisteredEvents(mock.Anything, uint64(0), mock.Anything).
		RunAndReturn(func(context.Context, uint64, *uint64) ([]blockchain.AddressRegistered, error) {
			call := queryCalls.Add(1)
			active := activeQueries.Add(1)
			for {
				maximum := maximumActiveQueries.Load()
				if active <= maximum || maximumActiveQueries.CompareAndSwap(maximum, active) {
					break
				}
			}
			defer activeQueries.Add(-1)
			if call == 1 {
				close(firstQueryStarted)
				<-releaseFirstQuery
			}
			return nil, nil
		}).
		Twice()
	registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(0)).Return([32]byte{}, nil).Twice()
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(0), nil).Twice()

	watcher := newOverlappingReplayWatcher(repository, registry, rskRpc, wallet)

	firstResult := make(chan error, 1)
	go func() {
		firstResult <- watcher.resync(context.Background())
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
		secondResult <- watcher.resync(context.Background())
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

func newOverlappingReplayWatcher(
	repository rootstock.PegInAddressRegistryWatchRepositorySet,
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	wallet blockchain.BitcoinWallet,
) *PegInAddressRegistryWatcher {
	return newTestPegInAddressRegistryWatcher(repository, registry, rskRpc, wallet, 0, 1, 0)
}

func newTestPegInAddressRegistryWatcher(
	repository rootstock.PegInAddressRegistryWatchRepositorySet,
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	wallet blockchain.BitcoinWallet,
	startBlock uint64,
	pageSize uint64,
	finalityDepth uint64,
) *PegInAddressRegistryWatcher {
	useCases := NewPegInAddressRegistryWatcherUseCases(
		w.NewDiscoverRegisteredAddressUseCase(repository, registry, wallet),
		w.NewGetWatchedRegisteredAddressesUseCase(repository),
	)
	return NewPegInAddressRegistryWatcher(
		useCases,
		repository,
		registry,
		rskRpc,
		nil,
		wallet,
		nil,
		nil,
		startBlock,
		pageSize,
		finalityDepth,
	)
}

func TestPegInAddressRegistryWatcher_DoesNotPublishIncrementalReplayBeforeRootVerification(t *testing.T) {
	rootErr := errors.New("root read failed")
	original := rootstock.PegInAddressRegistryWatchCheckpoint{
		LastProcessedBlock: 100,
	}
	repository := &overlappingReplayRepository{checkpoint: original, found: true}
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	registry.EXPECT().
		GetAddressRegisteredEvents(mock.Anything, uint64(101), mock.Anything).
		Return(nil, nil).
		Once()
	registry.EXPECT().
		GetRegistrationRoot(mock.Anything, uint64(101)).
		Return([32]byte{}, rootErr).
		Once()
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(101), nil).Once()
	watcher := newTestPegInAddressRegistryWatcher(
		repository,
		registry,
		rskRpc,
		wallet,
		0,
		10,
		0,
	)
	require.NoError(t, watcher.Prepare(context.Background()))

	err := watcher.scan(context.Background())

	require.ErrorIs(t, err, rootErr)
	checkpoint, found, checkpointErr := repository.GetCheckpoint(context.Background())
	require.NoError(t, checkpointErr)
	require.True(t, found)
	assert.Equal(t, original, checkpoint)
}

func TestPegInAddressRegistryWatcher_LogsRootReadFailureWithoutPublishingCheckpoint(t *testing.T) {
	rootErr := errors.New("RSK registry root read failed")
	original := rootstock.PegInAddressRegistryWatchCheckpoint{
		LocalRoot:          [32]byte{1},
		LastProcessedBlock: 100,
	}
	repository := &overlappingReplayRepository{checkpoint: original, found: true}
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	registry.EXPECT().
		GetAddressRegisteredEvents(mock.Anything, uint64(101), mock.Anything).
		Return(nil, nil).
		Once()
	registry.EXPECT().
		GetRegistrationRoot(mock.Anything, uint64(101)).
		Return([32]byte{}, rootErr).
		Once()
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(101), nil).Once()
	watcher := newTestPegInAddressRegistryWatcher(
		repository,
		registry,
		rskRpc,
		mocks.NewBitcoinWalletMock(t),
		0,
		10,
		0,
	)
	require.NoError(t, watcher.Prepare(context.Background()))
	capturedLogs := test.CaptureStructuredLogs(t)

	watcher.scanAndLog()

	logEntries := capturedLogs()
	require.Len(t, logEntries, 1)
	assert.Equal(t, "error", logEntries[0].Level())
	assert.Equal(
		t,
		"PegIn address registry watcher scan failed: get PegIn address registry root: RSK registry root read failed",
		logEntries[0].Message(),
	)
	checkpoint, found, setCalls, deleteCalls := repository.checkpointState()
	require.True(t, found)
	assert.Equal(t, original, checkpoint)
	assert.Zero(t, setCalls)
	assert.Zero(t, deleteCalls)
}
