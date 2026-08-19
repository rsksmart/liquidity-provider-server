package watcher_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	watcherAdapter "github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	watcherUseCase "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type pegInAddressRegistryWatchRepositorySetMock struct {
	*mocks.PegInAddressRegistryWatchRepositoryMock
	checkpoints mock.Mock
}

func newPegInAddressRegistryWatchRepositorySetMock(
	t *testing.T,
) *pegInAddressRegistryWatchRepositorySetMock {
	t.Helper()
	repository := &pegInAddressRegistryWatchRepositorySetMock{
		PegInAddressRegistryWatchRepositoryMock: mocks.NewPegInAddressRegistryWatchRepositoryMock(t),
	}
	repository.checkpoints.Test(t)
	t.Cleanup(func() { repository.checkpoints.AssertExpectations(t) })
	return repository
}

func (repository *pegInAddressRegistryWatchRepositorySetMock) GetCheckpoint(
	ctx context.Context,
) (rootstock.PegInAddressRegistryWatchCheckpoint, bool, error) {
	args := repository.checkpoints.MethodCalled("GetCheckpoint", ctx)
	checkpoint, ok := args.Get(0).(rootstock.PegInAddressRegistryWatchCheckpoint)
	if !ok {
		return rootstock.PegInAddressRegistryWatchCheckpoint{}, false, errors.New("invalid checkpoint mock result")
	}
	return checkpoint, args.Bool(1), args.Error(2)
}

func (repository *pegInAddressRegistryWatchRepositorySetMock) SetCheckpoint(
	ctx context.Context,
	checkpoint rootstock.PegInAddressRegistryWatchCheckpoint,
) error {
	return repository.checkpoints.MethodCalled("SetCheckpoint", ctx, checkpoint).Error(0)
}

func (repository *pegInAddressRegistryWatchRepositorySetMock) DeleteCheckpoint(ctx context.Context) error {
	return repository.checkpoints.MethodCalled("DeleteCheckpoint", ctx).Error(0)
}

type pegInAddressRegistryWatcherFixture struct {
	t          *testing.T
	repository *pegInAddressRegistryWatchRepositorySetMock
	registry   *mocks.PegInAddressRegistryContractMock
	rskRpc     *mocks.RootstockRpcServerMock
	btcNetwork *mocks.BtcRpcMock
	wallet     *mocks.BitcoinWalletMock
	ticker     *sessionTicker
	watcher    *watcherAdapter.PegInAddressRegistryWatcher
}

//nolint:unparam // Fixed start/finality values make each scenario's checkpoint arithmetic explicit at the call site.
func newPegInAddressRegistryWatcherFixture(
	t *testing.T,
	startBlock uint64,
	pageSize uint64,
	finalityDepth uint64,
) *pegInAddressRegistryWatcherFixture {
	t.Helper()
	fixture := &pegInAddressRegistryWatcherFixture{
		t:          t,
		repository: newPegInAddressRegistryWatchRepositorySetMock(t),
		registry:   mocks.NewPegInAddressRegistryContractMock(t),
		rskRpc:     mocks.NewRootstockRpcServerMock(t),
		btcNetwork: &mocks.BtcRpcMock{},
		wallet:     mocks.NewBitcoinWalletMock(t),
		ticker:     newSessionTicker(),
	}
	t.Cleanup(func() { fixture.btcNetwork.AssertExpectations(t) })
	fixture.watcher = newRegistryWatcher(
		fixture.repository,
		fixture.registry,
		fixture.rskRpc,
		fixture.btcNetwork,
		fixture.wallet,
		nil,
		fixture.ticker,
		startBlock,
		pageSize,
		finalityDepth,
	)
	return fixture
}

func newRegistryWatcher(
	repository rootstock.PegInAddressRegistryWatchRepositorySet,
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	btcNetwork blockchain.BitcoinNetwork,
	wallet blockchain.BitcoinWallet,
	eventBus entities.EventBus,
	ticker utils.Ticker,
	startBlock uint64,
	pageSize uint64,
	finalityDepth uint64,
) *watcherAdapter.PegInAddressRegistryWatcher {
	discover := watcherUseCase.NewDiscoverRegisteredAddressUseCase(repository, registry, wallet)
	getWatched := watcherUseCase.NewGetWatchedRegisteredAddressesUseCase(repository)
	replay := watcherUseCase.NewReplayRegisteredAddressesUseCase(
		discover,
		getWatched,
		repository,
		registry,
		rskRpc,
		eventBus,
		startBlock,
		pageSize,
		finalityDepth,
	)
	return watcherAdapter.NewPegInAddressRegistryWatcher(
		watcherAdapter.NewPegInAddressRegistryWatcherUseCases(discover, replay),
		btcNetwork,
		wallet,
		eventBus,
		ticker,
	)
}

func (fixture *pegInAddressRegistryWatcherFixture) expectBoundedRescan() {
	fixture.t.Helper()
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{StartHeight: 100, StopHeight: 200}, nil).
		Once()
}

func discoveredWatchEntry(event blockchain.AddressRegistered) rootstock.PegInAddressRegistryWatchEntry {
	return rootstock.PegInAddressRegistryWatchEntry{
		TxHash:      event.TxHash,
		LogIndex:    event.LogIndex,
		BlockNumber: event.BlockNumber,
		RskAddress:  event.RskAddress,
		State:       rootstock.PegInAddressRegistryWatchDiscovered,
	}
}

func (fixture *pegInAddressRegistryWatcherFixture) expectDiscover(event blockchain.AddressRegistered) {
	fixture.t.Helper()
	persisted := discoveredWatchEntry(event)
	fixture.repository.EXPECT().Get(mock.Anything, event.RskAddress).Return(nil, nil).Once()
	fixture.repository.EXPECT().Upsert(mock.Anything, mock.Anything).Return(nil).Once()
	fixture.repository.EXPECT().Get(mock.Anything, event.RskAddress).Return(&persisted, nil).Once()
}

func (fixture *pegInAddressRegistryWatcherFixture) expectCheckpoint(
	checkpoint rootstock.PegInAddressRegistryWatchCheckpoint,
	found bool,
) {
	fixture.t.Helper()
	fixture.repository.checkpoints.On("GetCheckpoint", mock.Anything).
		Return(checkpoint, found, nil).
		Once()
}

func (fixture *pegInAddressRegistryWatcherFixture) expectCheckpointAdvance(toBlock uint64) {
	fixture.t.Helper()
	fixture.repository.checkpoints.On(
		"SetCheckpoint",
		mock.Anything,
		mock.MatchedBy(func(checkpoint rootstock.PegInAddressRegistryWatchCheckpoint) bool {
			return checkpoint.LastProcessedBlock == toBlock
		}),
	).Return(nil).Once()
}

// runScan is one boot, one tick, one shutdown. The session returns from the tick when the scan has
// finished, so a scenario does not have to nominate a repository call as the end of the poll.
func (fixture *pegInAddressRegistryWatcherFixture) runScan() {
	fixture.t.Helper()
	session := startWatcherSession(fixture.t, fixture.watcher, fixture.ticker)
	session.poll()
	session.stop()
}

func TestPegInAddressRegistryWatcher_DoesNotQueryOrCheckpointBeforeConfiguredStart(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 10, 2)
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(101), nil).Once()

	fixture.runScan()

	fixture.registry.AssertNotCalled(t, "GetPegInAddress", mock.Anything)
	fixture.registry.AssertNotCalled(t, "GetAddressRegisteredEvents", mock.Anything, mock.Anything, mock.Anything)
	fixture.registry.AssertNotCalled(t, "GetRegistrationRoot", mock.Anything, mock.Anything)
	fixture.repository.checkpoints.AssertNotCalled(t, "SetCheckpoint", mock.Anything, mock.Anything)
	fixture.repository.checkpoints.AssertNotCalled(t, "DeleteCheckpoint", mock.Anything)
}

func TestPegInAddressRegistryWatcher_VerifiesRootAtCapturedHead(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 10, 2)
	fixture.repository.EXPECT().List(mock.Anything).
		Return([]rootstock.PegInAddressRegistryWatchEntry{}, nil).
		Once()
	fixture.expectCheckpoint(rootstock.PegInAddressRegistryWatchCheckpoint{}, false)
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(104), nil).Once()
	fixture.registry.EXPECT().
		GetAddressRegisteredEvents(mock.Anything, uint64(100), uint64Pointer(102)).
		Return([]blockchain.AddressRegistered{}, nil).
		Once()
	fixture.expectCheckpointAdvance(102)

	var latestHead atomic.Uint64
	latestHead.Store(104)
	fixture.registry.EXPECT().
		GetAddressRegisteredEvents(mock.Anything, uint64(103), uint64Pointer(104)).
		Run(func(context.Context, uint64, *uint64) {
			latestHead.Store(105)
		}).
		Return([]blockchain.AddressRegistered{}, nil).
		Once()
	fixture.registry.EXPECT().
		GetRegistrationRoot(mock.Anything, uint64(104)).
		RunAndReturn(func(context.Context, uint64) ([32]byte, error) {
			require.Equal(t, uint64(105), latestHead.Load(), "latest state must advance before the pinned root read")
			return [32]byte{}, nil
		}).
		Once()

	session := startWatcherSession(t, fixture.watcher, fixture.ticker)
	session.poll()
	session.stop()

	fixture.repository.checkpoints.AssertNotCalled(t, "DeleteCheckpoint", mock.Anything)
}

//nolint:funlen // One scenario records the complete persist/import/checkpoint ordering across three events.
func TestPegInAddressRegistryWatcher_SortsEventsAndClampsToFinalizedHead(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 20, 2)
	events := []blockchain.AddressRegistered{
		{
			TxHash:           "third",
			RskAddress:       "0x00000000000000000000000000000000000000a3",
			BlockNumber:      102,
			LogIndex:         5,
			RegistrationRoot: registrationRoot(t, "b80604f49f0685bb17dd0f5cc0f611383d724f10f53db8aebcec3e9541f552d8"),
		},
		{
			TxHash:           "first",
			RskAddress:       "0x00000000000000000000000000000000000000a1",
			BlockNumber:      101,
			LogIndex:         9,
			RegistrationRoot: registrationRoot(t, "5a16856e66cb2b1b463f7773c427085d55afdd19d778290b45fb959a6224877e"),
		},
		{
			TxHash:           "second",
			RskAddress:       "0x00000000000000000000000000000000000000a2",
			BlockNumber:      102,
			LogIndex:         1,
			RegistrationRoot: registrationRoot(t, "95bdcbc864b6248af173e5feb03f5830479e972a492e705c3c2a2bddcb8ca643"),
		},
	}
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInAddressRegistryWatchEntry{}, nil).Once()
	fixture.expectCheckpoint(rootstock.PegInAddressRegistryWatchCheckpoint{}, false)
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(105), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(100), uint64Pointer(103)).
		Return(events, nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(104), uint64Pointer(105)).
		Return([]blockchain.AddressRegistered{}, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(105)).
		Return(events[0].RegistrationRoot, nil).Once()

	var operations []string
	var operationsMutex sync.Mutex
	ordered := []blockchain.AddressRegistered{events[1], events[2], events[0]}
	for index, event := range ordered {
		deposit := knownDepositAddress(index)
		persisted := discoveredWatchEntry(event)
		fixture.repository.EXPECT().Get(mock.Anything, event.RskAddress).Return(nil, nil).Once()
		fixture.repository.On("Upsert", mock.Anything, mock.MatchedBy(func(entry rootstock.PegInAddressRegistryWatchEntry) bool {
			return entry.RskAddress == event.RskAddress
		})).
			Run(func(arguments mock.Arguments) {
				entry, ok := arguments.Get(1).(rootstock.PegInAddressRegistryWatchEntry)
				require.True(t, ok)
				operationsMutex.Lock()
				operations = append(operations, fmt.Sprintf("upsert:%d/%d/%s", entry.BlockNumber, entry.LogIndex, entry.TxHash))
				operationsMutex.Unlock()
			}).
			Return(nil).
			Once()
		fixture.repository.EXPECT().Get(mock.Anything, event.RskAddress).Return(&persisted, nil).Once()
		fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).
			Return(blockchain.PegInAddress{
				Payload:  deposit.payload,
				Encoding: blockchain.PegInAddressRegistryEncodingBase58,
			}, nil).
			Once()
		var importErr error
		if event.TxHash == "second" {
			importErr = errors.New("address already imported")
		}
		fixture.wallet.EXPECT().ImportAddress(deposit.address).
			Run(func(string) {
				operationsMutex.Lock()
				operations = append(operations, "import:"+event.TxHash)
				operationsMutex.Unlock()
			}).
			Return(importErr).
			Once()
		fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(entry rootstock.PegInAddressRegistryWatchEntry) bool {
			return entry.RskAddress == event.RskAddress &&
				entry.BtcAddress == deposit.address &&
				entry.Encoding == uint8(blockchain.PegInAddressRegistryEncodingBase58) &&
				entry.State == rootstock.PegInAddressRegistryWatchImported
		})).Run(func(mock.Arguments) {
			operationsMutex.Lock()
			operations = append(operations, "update:"+event.TxHash)
			operationsMutex.Unlock()
		}).Return(nil).Once()
	}
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().RescanBlockchain(int64(100)).
		Run(func(int64) {
			operationsMutex.Lock()
			operations = append(operations, "rescan")
			operationsMutex.Unlock()
		}).
		Return(blockchain.BitcoinRescanResult{StartHeight: 100, StopHeight: 200}, nil).
		Once()
	fixture.repository.checkpoints.On(
		"SetCheckpoint",
		mock.Anything,
		mock.MatchedBy(func(checkpoint rootstock.PegInAddressRegistryWatchCheckpoint) bool {
			return checkpoint.LastProcessedBlock == 103
		}),
	).
		Run(func(mock.Arguments) {
			operationsMutex.Lock()
			operations = append(operations, "checkpoint")
			operationsMutex.Unlock()
		}).
		Return(nil).
		Once()
	fixture.runScan()

	operationsMutex.Lock()
	assert.Equal(t, []string{
		"upsert:101/9/first", "import:first",
		"upsert:102/1/second", "import:second",
		"upsert:102/5/third", "import:third",
		"rescan",
		"update:first", "update:second", "update:third",
		"checkpoint",
	}, operations)
	operationsMutex.Unlock()
}

func TestPegInAddressRegistryWatcher_SkipsUnsupportedEncodingsAndContinues(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 1, 2)
	events := unsupportedEncodingEvents(t)
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInAddressRegistryWatchEntry{}, nil).Once()
	fixture.expectCheckpoint(rootstock.PegInAddressRegistryWatchCheckpoint{}, false)
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(102), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(100), uint64Pointer(100)).
		Return(events, nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(101), uint64Pointer(102)).
		Return([]blockchain.AddressRegistered{}, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(102)).
		Return(events[2].RegistrationRoot, nil).Once()

	expectEncodingOutcomes(fixture, events)
	fixture.expectCheckpointAdvance(100)

	fixture.runScan()
}

func unsupportedEncodingEvents(t *testing.T) []blockchain.AddressRegistered {
	t.Helper()
	return []blockchain.AddressRegistered{
		{
			TxHash:           "bech32",
			RskAddress:       "0x00000000000000000000000000000000000000a1",
			BlockNumber:      100,
			LogIndex:         1,
			RegistrationRoot: registrationRoot(t, "5a16856e66cb2b1b463f7773c427085d55afdd19d778290b45fb959a6224877e"),
		},
		{
			TxHash:           "bech32m",
			RskAddress:       "0x00000000000000000000000000000000000000a2",
			BlockNumber:      100,
			LogIndex:         2,
			RegistrationRoot: registrationRoot(t, "95bdcbc864b6248af173e5feb03f5830479e972a492e705c3c2a2bddcb8ca643"),
		},
		{
			TxHash:           "base58",
			RskAddress:       "0x00000000000000000000000000000000000000a3",
			BlockNumber:      100,
			LogIndex:         3,
			RegistrationRoot: registrationRoot(t, "b80604f49f0685bb17dd0f5cc0f611383d724f10f53db8aebcec3e9541f552d8"),
		},
	}
}

func expectEncodingOutcomes(
	fixture *pegInAddressRegistryWatcherFixture,
	events []blockchain.AddressRegistered,
) {
	fixture.t.Helper()
	encodings := []blockchain.PegInAddressRegistryEncoding{
		blockchain.PegInAddressRegistryEncodingBech32,
		blockchain.PegInAddressRegistryEncodingBech32M,
		blockchain.PegInAddressRegistryEncodingBase58,
	}
	for index, event := range events {
		deposit := knownDepositAddress(index)
		fixture.expectDiscover(event)
		fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).
			Return(blockchain.PegInAddress{
				Payload:  deposit.payload,
				Encoding: encodings[index],
			}, nil).
			Once()
		expectedState := rootstock.PegInAddressRegistryWatchUnsupportedEncoding
		// An encoding the scanner cannot render must leave BtcAddress empty rather than storing the
		// raw payload, which is not a valid address in any encoding.
		expectedAddress := ""
		if encodings[index] == blockchain.PegInAddressRegistryEncodingBase58 {
			expectedState = rootstock.PegInAddressRegistryWatchImported
			expectedAddress = deposit.address
			fixture.wallet.EXPECT().ImportAddress(deposit.address).Return(nil).Once()
			fixture.expectBoundedRescan()
		}
		fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(entry rootstock.PegInAddressRegistryWatchEntry) bool {
			return entry.TxHash == event.TxHash &&
				entry.State == expectedState &&
				entry.BtcAddress == expectedAddress
		})).Return(nil).Once()
	}
}

func TestPegInAddressRegistryWatcher_RecordsEntryErrorAndContinues(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 1, 2)
	broken, valid := entryErrorEvents(t)
	fixture.repository.EXPECT().List(mock.Anything).Return(nil, nil).Once()
	fixture.expectCheckpoint(rootstock.PegInAddressRegistryWatchCheckpoint{}, false)
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(102), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(100), uint64Pointer(100)).
		Return([]blockchain.AddressRegistered{broken, valid}, nil).
		Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(101), uint64Pointer(102)).
		Return([]blockchain.AddressRegistered{}, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(102)).
		Return(valid.RegistrationRoot, nil).Once()

	expectBrokenEntryUpdate(fixture, broken)
	expectValidEntryUpdate(fixture, valid)
	fixture.expectCheckpointAdvance(100)

	fixture.runScan()
}

func entryErrorEvents(t *testing.T) (blockchain.AddressRegistered, blockchain.AddressRegistered) {
	t.Helper()
	broken := blockchain.AddressRegistered{
		TxHash:           "broken",
		RskAddress:       "0x00000000000000000000000000000000000000a1",
		BlockNumber:      100,
		LogIndex:         1,
		RegistrationRoot: registrationRoot(t, "5a16856e66cb2b1b463f7773c427085d55afdd19d778290b45fb959a6224877e"),
	}
	valid := blockchain.AddressRegistered{
		TxHash:           "valid",
		RskAddress:       "0x00000000000000000000000000000000000000a2",
		BlockNumber:      100,
		LogIndex:         2,
		RegistrationRoot: registrationRoot(t, "95bdcbc864b6248af173e5feb03f5830479e972a492e705c3c2a2bddcb8ca643"),
	}
	return broken, valid
}

func expectBrokenEntryUpdate(
	fixture *pegInAddressRegistryWatcherFixture,
	broken blockchain.AddressRegistered,
) {
	fixture.t.Helper()
	fixture.expectDiscover(broken)
	fixture.registry.EXPECT().GetPegInAddress(broken.RskAddress).
		Return(blockchain.PegInAddress{}, assert.AnError).
		Once()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(entry rootstock.PegInAddressRegistryWatchEntry) bool {
		return entry.TxHash == broken.TxHash &&
			entry.State == rootstock.PegInAddressRegistryWatchDiscovered &&
			entry.LastError != ""
	})).Return(nil).Once()
}

func expectValidEntryUpdate(
	fixture *pegInAddressRegistryWatcherFixture,
	valid blockchain.AddressRegistered,
) {
	fixture.t.Helper()
	fixture.expectDiscover(valid)
	validDeposit := knownDepositAddress(0)
	fixture.registry.EXPECT().GetPegInAddress(valid.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  validDeposit.payload,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.wallet.EXPECT().ImportAddress(validDeposit.address).Return(nil).Once()
	fixture.expectBoundedRescan()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(entry rootstock.PegInAddressRegistryWatchEntry) bool {
		return entry.TxHash == valid.TxHash &&
			entry.State == rootstock.PegInAddressRegistryWatchImported &&
			entry.LastError == ""
	})).Return(nil).Once()
}

func TestPegInAddressRegistryWatcher_RetriesPersistedDiscoveredEntryOutsideOverlap(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 1, 2)
	entry := rootstock.PegInAddressRegistryWatchEntry{
		TxHash:      "retry-registration",
		LogIndex:    1,
		BlockNumber: 90,
		RskAddress:  "retry-rsk",
		State:       rootstock.PegInAddressRegistryWatchDiscovered,
		LastError:   "previous transient failure",
	}
	fixture.repository.EXPECT().List(mock.Anything).
		Return([]rootstock.PegInAddressRegistryWatchEntry{entry}, nil).
		Once()
	fixture.expectCheckpoint(rootstock.PegInAddressRegistryWatchCheckpoint{LastProcessedBlock: 105}, true)
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(108), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(106), uint64Pointer(106)).
		Return(nil, nil).
		Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(107), uint64Pointer(108)).
		Return(nil, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(108)).Return([32]byte{}, nil).Once()
	fixture.repository.EXPECT().Get(mock.Anything, entry.RskAddress).Return(&entry, nil).Once()
	deposit := knownDepositAddress(0)
	fixture.registry.EXPECT().GetPegInAddress(entry.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  deposit.payload,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.wallet.EXPECT().ImportAddress(deposit.address).Return(nil).Once()
	fixture.expectBoundedRescan()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(updated rootstock.PegInAddressRegistryWatchEntry) bool {
		return updated.TxHash == entry.TxHash &&
			updated.BtcAddress == deposit.address &&
			updated.State == rootstock.PegInAddressRegistryWatchImported &&
			updated.LastError == ""
	})).Return(nil).Once()
	fixture.expectCheckpointAdvance(106)

	fixture.runScan()
}

// An entry that keeps failing the same way every tick must not rewrite the same error to Mongo,
// which only holds if the persisted error survives the successful steps that precede the failure.
func TestPegInAddressRegistryWatcher_SuppressesRepeatedIdenticalEntryErrors(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 1, 2)
	entry := rootstock.PegInAddressRegistryWatchEntry{
		TxHash:      "stuck",
		LogIndex:    1,
		BlockNumber: 90,
		RskAddress:  "stuck-rsk",
		State:       rootstock.PegInAddressRegistryWatchDiscovered,
		LastError:   fmt.Sprintf("import PegIn address for event stuck/1: %v", assert.AnError),
	}
	fixture.repository.EXPECT().List(mock.Anything).
		Return([]rootstock.PegInAddressRegistryWatchEntry{entry}, nil).
		Once()
	fixture.expectCheckpoint(rootstock.PegInAddressRegistryWatchCheckpoint{LastProcessedBlock: 105}, true)
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(108), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(106), uint64Pointer(106)).
		Return(nil, nil).
		Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(107), uint64Pointer(108)).
		Return(nil, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(108)).Return([32]byte{}, nil).Once()
	fixture.repository.EXPECT().Get(mock.Anything, entry.RskAddress).Return(&entry, nil).Once()
	deposit := knownDepositAddress(0)
	fixture.registry.EXPECT().GetPegInAddress(entry.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  deposit.payload,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.wallet.EXPECT().ImportAddress(deposit.address).Return(assert.AnError).Once()
	fixture.expectCheckpointAdvance(106)

	fixture.runScan()

	fixture.repository.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// The registry returns a payload rather than an address, so a payload that cannot be encoded must
// park the entry with the reason instead of importing something a node would reject.
func TestPegInAddressRegistryWatcher_RecordsUnencodablePayload(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 1, 2)
	entry := rootstock.PegInAddressRegistryWatchEntry{
		TxHash:      "truncated-payload",
		LogIndex:    1,
		BlockNumber: 90,
		RskAddress:  "truncated-rsk",
		State:       rootstock.PegInAddressRegistryWatchDiscovered,
	}
	fixture.repository.EXPECT().List(mock.Anything).
		Return([]rootstock.PegInAddressRegistryWatchEntry{entry}, nil).
		Once()
	fixture.expectCheckpoint(rootstock.PegInAddressRegistryWatchCheckpoint{LastProcessedBlock: 105}, true)
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(108), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(106), uint64Pointer(106)).
		Return(nil, nil).
		Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(107), uint64Pointer(108)).
		Return(nil, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(108)).Return([32]byte{}, nil).Once()
	fixture.repository.EXPECT().Get(mock.Anything, entry.RskAddress).Return(&entry, nil).Once()
	truncated := knownDepositAddress(0).payload[:20]
	fixture.registry.EXPECT().GetPegInAddress(entry.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  truncated,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(updated rootstock.PegInAddressRegistryWatchEntry) bool {
		return updated.TxHash == entry.TxHash &&
			updated.State == rootstock.PegInAddressRegistryWatchDiscovered &&
			updated.BtcAddress == "" &&
			strings.Contains(updated.LastError, "encode PegIn address for event truncated-payload/1")
	})).Return(nil).Once()
	fixture.expectCheckpointAdvance(106)

	fixture.runScan()

	fixture.wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
	fixture.wallet.AssertNotCalled(t, "RescanBlockchain", mock.Anything)
}

func TestPegInAddressRegistryWatcher_DoesNotReimportPersistedEntry(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 1, 2)
	event := unsupportedEncodingEvents(t)[0]
	persisted := rootstock.PegInAddressRegistryWatchEntry{
		TxHash:     event.TxHash,
		LogIndex:   event.LogIndex,
		RskAddress: event.RskAddress,
		State:      rootstock.PegInAddressRegistryWatchImported,
		BtcAddress: knownDepositAddress(0).address,
	}
	fixture.repository.EXPECT().List(mock.Anything).
		Return([]rootstock.PegInAddressRegistryWatchEntry{persisted}, nil).
		Once()
	fixture.expectCheckpoint(rootstock.PegInAddressRegistryWatchCheckpoint{}, false)
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(102), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(100), uint64Pointer(100)).
		Return([]blockchain.AddressRegistered{event}, nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(101), uint64Pointer(102)).
		Return([]blockchain.AddressRegistered{}, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(102)).
		Return(event.RegistrationRoot, nil).Once()
	fixture.repository.EXPECT().Get(mock.Anything, event.RskAddress).Return(&persisted, nil).Once()
	fixture.expectCheckpointAdvance(100)

	fixture.runScan()

	fixture.wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
	fixture.repository.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}

func TestPegInAddressRegistryWatcher_DuplicateRskAddressKeepsOneWatchRow(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 1, 2)
	first := unsupportedEncodingEvents(t)[0]
	replayRoot, err := blockchain.FoldPegInAddressRegistryRoot(first.RegistrationRoot, first.RskAddress)
	require.NoError(t, err)
	replay := blockchain.AddressRegistered{
		TxHash:           "replay-reg",
		RskAddress:       first.RskAddress,
		BlockNumber:      100,
		LogIndex:         2,
		RegistrationRoot: replayRoot,
	}
	imported := rootstock.PegInAddressRegistryWatchEntry{
		TxHash:     first.TxHash,
		LogIndex:   first.LogIndex,
		RskAddress: first.RskAddress,
		State:      rootstock.PegInAddressRegistryWatchImported,
		BtcAddress: knownDepositAddress(0).address,
	}
	fixture.repository.EXPECT().List(mock.Anything).
		Return([]rootstock.PegInAddressRegistryWatchEntry{imported}, nil).
		Once()
	fixture.expectCheckpoint(rootstock.PegInAddressRegistryWatchCheckpoint{}, false)
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(102), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(100), uint64Pointer(100)).
		Return([]blockchain.AddressRegistered{first, replay}, nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(101), uint64Pointer(102)).
		Return([]blockchain.AddressRegistered{}, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(102)).
		Return(replayRoot, nil).Once()
	fixture.repository.EXPECT().Get(mock.Anything, first.RskAddress).Return(&imported, nil).Twice()
	fixture.expectCheckpointAdvance(100)

	fixture.runScan()

	fixture.wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
	fixture.repository.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}

func uint64Pointer(expected uint64) interface{} {
	return mock.MatchedBy(func(actual *uint64) bool {
		return actual != nil && *actual == expected
	})
}
