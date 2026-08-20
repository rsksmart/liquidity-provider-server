package watcher

import (
	"cmp"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/datasets"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type sessionTicker struct {
	ticks   chan time.Time
	selects atomic.Int64
	stops   atomic.Int64
}

func newSessionTicker() *sessionTicker {
	return &sessionTicker{ticks: make(chan time.Time)}
}

func (ticker *sessionTicker) C() <-chan time.Time {
	ticker.selects.Add(1)
	return ticker.ticks
}

func (ticker *sessionTicker) Stop() {
	ticker.stops.Add(1)
}

type watcherSession struct {
	t            *testing.T
	watcher      Watcher
	ticker       *sessionTicker
	closeChannel chan bool
}

func startWatcherSession(t *testing.T, watcher Watcher, ticker *sessionTicker) *watcherSession {
	t.Helper()
	session := &watcherSession{
		t:            t,
		watcher:      watcher,
		ticker:       ticker,
		closeChannel: make(chan bool),
	}
	require.NoError(t, watcher.Prepare(context.Background()))
	go watcher.Start()
	session.waitUntilIdle(0)
	return session
}

func (session *watcherSession) poll() {
	session.t.Helper()
	idleBefore := session.ticker.selects.Load()
	select {
	case session.ticker.ticks <- time.Now():
	case <-time.After(time.Second):
		require.FailNow(session.t, "watcher did not accept the tick")
	}
	session.waitUntilIdle(idleBefore)
}

func (session *watcherSession) waitUntilIdle(idleBefore int64) {
	session.t.Helper()
	require.Eventually(
		session.t,
		func() bool { return session.ticker.selects.Load() > idleBefore },
		time.Second,
		time.Millisecond,
		"watcher poll did not return",
	)
}

func (session *watcherSession) stop() {
	session.t.Helper()
	go session.watcher.Shutdown(session.closeChannel)
	select {
	case <-session.closeChannel:
	case <-time.After(time.Second):
		require.FailNow(session.t, "watcher shutdown did not complete")
	}
	assert.Eventually(
		session.t,
		func() bool { return session.ticker.stops.Load() == 1 },
		time.Second,
		time.Millisecond,
		"shutdown must stop the watcher's ticker exactly once",
	)
}

// --- merged from pegin_address_registry_test_helpers_test.go ---

type depositAddress struct {
	payload []byte
	address string
}

type pinnedLBCRegistration struct {
	rskAddress string
	root       [32]byte
}

func pinnedLBCRegistrations(t *testing.T) []pinnedLBCRegistration {
	t.Helper()
	return []pinnedLBCRegistration{
		{
			rskAddress: "0x00000000000000000000000000000000000000a1",
			root:       registrationRoot(t, "5a16856e66cb2b1b463f7773c427085d55afdd19d778290b45fb959a6224877e"),
		},
		{
			rskAddress: "0x00000000000000000000000000000000000000a2",
			root:       registrationRoot(t, "95bdcbc864b6248af173e5feb03f5830479e972a492e705c3c2a2bddcb8ca643"),
		},
		{
			rskAddress: "0x00000000000000000000000000000000000000a3",
			root:       registrationRoot(t, "b80604f49f0685bb17dd0f5cc0f611383d724f10f53db8aebcec3e9541f552d8"),
		},
	}
}

func knownDepositAddress(index int) depositAddress {
	const checksumSize = 4
	decoded := datasets.Base58Addresses[index]
	payload := make([]byte, 0, len(decoded.Expected)+checksumSize)
	payload = append(payload, decoded.Expected...)
	payload = append(payload, chainhash.DoubleHashB(decoded.Expected)[:checksumSize]...)
	return depositAddress{payload: payload, address: decoded.Address}
}

func registrationRoot(t *testing.T, encoded string) [32]byte {
	t.Helper()
	decoded, err := hex.DecodeString(encoded)
	require.NoError(t, err)
	require.Len(t, decoded, 32)
	return [32]byte(decoded)
}

// --- merged from pegin_address_registry_watcher_test.go ---

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
	watcher    *PegInAddressRegistryWatcher
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
) *PegInAddressRegistryWatcher {
	replay := w.NewReplayRegisteredAddressesUseCase(
		repository,
		registry,
		rskRpc,
		eventBus,
		wallet,
		finalityDepth,
	)
	finalize := w.NewFinalizeRegisteredAddressImportUseCase(repository)
	return NewPegInAddressRegistryWatcher(
		replay,
		finalize,
		btcNetwork,
		wallet,
		eventBus,
		ticker,
		startBlock,
		pageSize,
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
		"checkpoint",
		"rescan",
		"update:first", "update:second", "update:third",
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

// --- merged from pegin_address_registry_resync_test.go ---

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

func newShortHeadCheckpointWatcher(t *testing.T) (
	*overlappingReplayRepository,
	*mocks.PegInAddressRegistryContractMock,
	*PegInAddressRegistryWatcher,
	rootstock.PegInAddressRegistryWatchCheckpoint,
) {
	t.Helper()
	original := rootstock.PegInAddressRegistryWatchCheckpoint{
		LocalRoot:          [32]byte{1},
		LastProcessedBlock: 105,
	}
	repository := &overlappingReplayRepository{checkpoint: original, found: true}
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(101), nil).Once()
	watcher := newRegistryWatcher(
		repository,
		registry,
		rskRpc,
		nil,
		mocks.NewBitcoinWalletMock(t),
		nil,
		nil,
		100,
		10,
		2,
	)
	return repository, registry, watcher, original
}

func TestPegInAddressRegistryWatcher_ScanKeepsCheckpointWhenHeadIsBelowStart(t *testing.T) {
	repository, registry, watcher, original := newShortHeadCheckpointWatcher(t)

	require.NoError(t, watcher.scan(context.Background()))

	checkpoint, found, setCalls, deleteCalls := repository.checkpointState()
	require.True(t, found)
	assert.Equal(t, original, checkpoint)
	assert.Zero(t, setCalls)
	assert.Zero(t, deleteCalls)
	registry.AssertNotCalled(t, "GetAddressRegisteredEvents", mock.Anything, mock.Anything, mock.Anything)
	registry.AssertNotCalled(t, "GetRegistrationRoot", mock.Anything, mock.Anything)
}

func TestPegInAddressRegistryWatcher_ResyncDeletesCheckpointWhenHeadIsBelowStart(t *testing.T) {
	repository, registry, watcher, _ := newShortHeadCheckpointWatcher(t)

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
	return newRegistryWatcher(repository, registry, rskRpc, nil, wallet, nil, nil, 0, 1, 0)
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
	watcher := newRegistryWatcher(
		repository,
		registry,
		rskRpc,
		nil,
		wallet,
		nil,
		nil,
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
	watcher := newRegistryWatcher(
		repository,
		registry,
		rskRpc,
		nil,
		mocks.NewBitcoinWalletMock(t),
		nil,
		nil,
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
		"PegIn address registry watcher scan failed: ReplayRegisteredAddresses: get PegIn address registry root: RSK registry root read failed",
		logEntries[0].Message(),
	)
	checkpoint, found, setCalls, deleteCalls := repository.checkpointState()
	require.True(t, found)
	assert.Equal(t, original, checkpoint)
	assert.Zero(t, setCalls)
	assert.Zero(t, deleteCalls)
}

// --- merged from pegin_address_registry_recovery_test.go ---

// A write that fails and a write lost to a dying process leave the same state behind, so a scenario
// injects this to stop a watcher at a boundary it cannot otherwise be stopped at.
var errProcessStopped = errors.New("process stopped")

// registryWatchStore stands in for the Mongo watch repository, keeping the behaviour a restarted
// watcher depends on: Upsert never overwrites an existing (tx_hash, log_index), Update needs the
// entry to exist, List answers in (block, log index) order, and the checkpoint reads as absent until
// it is first written.
type registryWatchStore struct {
	mutex         sync.Mutex
	entries       []rootstock.PegInAddressRegistryWatchEntry
	checkpoint    rootstock.PegInAddressRegistryWatchCheckpoint
	hasCheckpoint bool
	advances      []uint64
	clears        int

	// Each of these fails the next call to the method it names and is then cleared.
	failGet           error
	failUpdate        error
	failSetCheckpoint error
}

func (store *registryWatchStore) Upsert(_ context.Context, entry rootstock.PegInAddressRegistryWatchEntry) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if store.indexOf(entry.RskAddress) >= 0 {
		return nil
	}
	store.entries = append(store.entries, entry)
	return nil
}

func (store *registryWatchStore) Get(
	_ context.Context,
	rskAddress string,
) (*rootstock.PegInAddressRegistryWatchEntry, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if err := takeNextFailure(&store.failGet); err != nil {
		return nil, err
	}
	index := store.indexOf(rskAddress)
	if index < 0 {
		return nil, nil
	}
	found := store.entries[index]
	return &found, nil
}

func (store *registryWatchStore) List(context.Context) ([]rootstock.PegInAddressRegistryWatchEntry, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	return store.sortedEntries(), nil
}

func (store *registryWatchStore) Update(_ context.Context, entry rootstock.PegInAddressRegistryWatchEntry) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if err := takeNextFailure(&store.failUpdate); err != nil {
		return err
	}
	index := store.indexOf(entry.RskAddress)
	if index < 0 {
		return errors.New("pegin address registry watch entry not found")
	}
	store.entries[index] = entry
	return nil
}

func (store *registryWatchStore) GetCheckpoint(
	context.Context,
) (rootstock.PegInAddressRegistryWatchCheckpoint, bool, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	return store.checkpoint, store.hasCheckpoint, nil
}

func (store *registryWatchStore) SetCheckpoint(
	_ context.Context,
	checkpoint rootstock.PegInAddressRegistryWatchCheckpoint,
) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if err := takeNextFailure(&store.failSetCheckpoint); err != nil {
		return err
	}
	store.checkpoint = checkpoint
	store.hasCheckpoint = true
	store.advances = append(store.advances, checkpoint.LastProcessedBlock)
	return nil
}

func (store *registryWatchStore) DeleteCheckpoint(context.Context) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.checkpoint = rootstock.PegInAddressRegistryWatchCheckpoint{}
	store.hasCheckpoint = false
	store.clears++
	return nil
}

func (store *registryWatchStore) checkpointClears() int {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	return store.clears
}

func (store *registryWatchStore) checkpointFound() bool {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	return store.hasCheckpoint
}

// checkpointAdvances is every block the checkpoint was moved to, in order, so a restart can be
// asserted not to move it backwards or skip a range.
func (store *registryWatchStore) checkpointAdvances() []uint64 {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	return append([]uint64(nil), store.advances...)
}

func (store *registryWatchStore) indexOf(rskAddress string) int {
	for index := range store.entries {
		if store.entries[index].RskAddress == rskAddress {
			return index
		}
	}
	return -1
}

func (store *registryWatchStore) sortedEntries() []rootstock.PegInAddressRegistryWatchEntry {
	entries := slices.Clone(store.entries)
	slices.SortFunc(entries, func(first, second rootstock.PegInAddressRegistryWatchEntry) int {
		return cmp.Or(
			cmp.Compare(first.BlockNumber, second.BlockNumber),
			cmp.Compare(first.LogIndex, second.LogIndex),
		)
	})
	return entries
}

// takeNextFailure clears the failure as it returns it, so an injected crash stops one write and the
// process that restarts after it gets through.
func takeNextFailure(failure *error) error {
	err := *failure
	*failure = nil
	return err
}

// rskChain is the log set the RSK node answers range queries from. A scenario mutates it to model a
// reorg or a node that under-reported a block, so a test states what the chain holds rather than
// what the node returns on call N.
type rskChain struct {
	mutex       sync.Mutex
	head        uint64
	logs        []blockchain.AddressRegistered
	requested   [][2]uint64
	omitOnce    map[string]int
	includeOnce map[string]bool
}

func (chain *rskChain) setHead(head uint64) {
	chain.mutex.Lock()
	defer chain.mutex.Unlock()
	chain.head = head
}

func (chain *rskChain) currentHead() uint64 {
	chain.mutex.Lock()
	defer chain.mutex.Unlock()
	return chain.head
}

func (chain *rskChain) add(event blockchain.AddressRegistered) {
	chain.mutex.Lock()
	defer chain.mutex.Unlock()
	chain.logs = append(chain.logs, event)
}

func (chain *rskChain) omitOnNextQuery(event blockchain.AddressRegistered) {
	chain.omitOnNextQueries(event, 1)
}

func (chain *rskChain) omitOnNextQueries(event blockchain.AddressRegistered, count int) {
	chain.mutex.Lock()
	defer chain.mutex.Unlock()
	if chain.omitOnce == nil {
		chain.omitOnce = make(map[string]int)
	}
	chain.omitOnce[eventIdentity(event)] += count
}

func (chain *rskChain) includeOnNextQuery(event blockchain.AddressRegistered) {
	chain.mutex.Lock()
	defer chain.mutex.Unlock()
	if chain.includeOnce == nil {
		chain.includeOnce = make(map[string]bool)
	}
	chain.includeOnce[eventIdentity(event)] = true
}

func (chain *rskChain) drop(event blockchain.AddressRegistered) {
	chain.mutex.Lock()
	defer chain.mutex.Unlock()
	remaining := make([]blockchain.AddressRegistered, 0, len(chain.logs))
	for _, log := range chain.logs {
		if log.TxHash != event.TxHash || log.LogIndex != event.LogIndex {
			remaining = append(remaining, log)
		}
	}
	chain.logs = remaining
}

func (chain *rskChain) logsIn(fromBlock, toBlock uint64) []blockchain.AddressRegistered {
	chain.mutex.Lock()
	defer chain.mutex.Unlock()
	chain.requested = append(chain.requested, [2]uint64{fromBlock, toBlock})
	events := make([]blockchain.AddressRegistered, 0)
	for _, log := range chain.logs {
		identity := eventIdentity(log)
		inRequestedRange := log.BlockNumber >= fromBlock && log.BlockNumber <= toBlock
		if inRequestedRange || chain.includeOnce[identity] {
			if chain.omitOnce[eventIdentity(log)] > 0 {
				chain.omitOnce[eventIdentity(log)]--
				continue
			}
			delete(chain.includeOnce, identity)
			events = append(events, log)
		}
	}
	return events
}

func eventIdentity(event blockchain.AddressRegistered) string {
	return fmt.Sprintf("%s/%d", event.TxHash, event.LogIndex)
}

func (chain *rskChain) rootAt(toBlock uint64) [32]byte {
	chain.mutex.Lock()
	defer chain.mutex.Unlock()
	events := chain.orderedUniqueLogs()
	var root [32]byte
	for _, event := range events {
		if event.BlockNumber > toBlock {
			break
		}
		root = event.RegistrationRoot
	}
	return root
}

func (chain *rskChain) registrationCount() int {
	chain.mutex.Lock()
	defer chain.mutex.Unlock()
	return len(chain.orderedUniqueLogs())
}

func (chain *rskChain) orderedUniqueLogs() []blockchain.AddressRegistered {
	events := make([]blockchain.AddressRegistered, 0, len(chain.logs))
	seen := make(map[string]struct{})
	for _, event := range chain.logs {
		identity := eventIdentity(event)
		if _, duplicate := seen[identity]; duplicate {
			continue
		}
		seen[identity] = struct{}{}
		events = append(events, event)
	}
	slices.SortFunc(events, func(first, second blockchain.AddressRegistered) int {
		return cmp.Or(
			cmp.Compare(first.BlockNumber, second.BlockNumber),
			cmp.Compare(first.LogIndex, second.LogIndex),
		)
	})
	return events
}

func (chain *rskChain) requestedRanges() [][2]uint64 {
	chain.mutex.Lock()
	defer chain.mutex.Unlock()
	return append([][2]uint64(nil), chain.requested...)
}

// registryRestartFixture owns everything that outlives a restart: the watch set, the RSK chain, and
// the Bitcoin node's wallet.
type registryRestartFixture struct {
	t          *testing.T
	store      *registryWatchStore
	chain      *rskChain
	registry   *mocks.PegInAddressRegistryContractMock
	rskRpc     *mocks.RootstockRpcServerMock
	btcNetwork *mocks.BtcRpcMock
	wallet     *mocks.BitcoinWalletMock
	eventBus   *registryReorgEventBus

	mutex       sync.Mutex
	deposits    map[string]depositAddress
	imports     []string
	heightCalls int
	rootReads   []uint64

	startBlock    uint64
	pageSize      uint64
	finalityDepth uint64
}

type registryReorgEventBus struct {
	events    chan entities.Event
	mutex     sync.Mutex
	published []entities.Event
}

func newRegistryReorgEventBus() *registryReorgEventBus {
	return &registryReorgEventBus{events: make(chan entities.Event, 10)}
}

func (bus *registryReorgEventBus) Publish(event entities.Event) {
	bus.mutex.Lock()
	bus.published = append(bus.published, event)
	bus.mutex.Unlock()
	if event.Id() == blockchain.NodeReorgCheckEventId {
		bus.events <- event
	}
}

func (bus *registryReorgEventBus) Subscribe(id entities.EventId) <-chan entities.Event {
	if id != blockchain.NodeReorgCheckEventId {
		return nil
	}
	return bus.events
}

func (bus *registryReorgEventBus) Shutdown(closeChannel chan<- bool) {
	closeChannel <- true
}

func (bus *registryReorgEventBus) publishedCount(id entities.EventId) int {
	return len(bus.publishedEvents(id))
}

func (bus *registryReorgEventBus) publishedEvents(id entities.EventId) []entities.Event {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()
	events := make([]entities.Event, 0)
	for _, event := range bus.published {
		if event.Id() == id {
			events = append(events, event)
		}
	}
	return events
}

//nolint:unparam // Fixed start/page/finality values keep each scenario's checkpoint arithmetic explicit at the call site.
func newRegistryRestartFixture(
	t *testing.T,
	startBlock uint64,
	pageSize uint64,
	finalityDepth uint64,
) *registryRestartFixture {
	t.Helper()
	fixture := &registryRestartFixture{
		t:             t,
		store:         &registryWatchStore{},
		chain:         &rskChain{},
		registry:      mocks.NewPegInAddressRegistryContractMock(t),
		rskRpc:        mocks.NewRootstockRpcServerMock(t),
		btcNetwork:    &mocks.BtcRpcMock{},
		wallet:        mocks.NewBitcoinWalletMock(t),
		eventBus:      newRegistryReorgEventBus(),
		deposits:      make(map[string]depositAddress),
		startBlock:    startBlock,
		pageSize:      pageSize,
		finalityDepth: finalityDepth,
	}

	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).RunAndReturn(func(context.Context) (uint64, error) {
		fixture.mutex.Lock()
		fixture.heightCalls++
		fixture.mutex.Unlock()
		return fixture.chain.currentHead(), nil
	})
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, fromBlock uint64, toBlock *uint64) ([]blockchain.AddressRegistered, error) {
			require.NotNil(t, toBlock, "the scanner must always bound its range at the finalized head")
			return fixture.chain.logsIn(fromBlock, *toBlock), nil
		})
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, blockNumber uint64) ([32]byte, error) {
			fixture.mutex.Lock()
			fixture.rootReads = append(fixture.rootReads, blockNumber)
			fixture.mutex.Unlock()
			return fixture.chain.rootAt(blockNumber), nil
		}).
		Maybe()
	fixture.registry.EXPECT().GetPegInAddress(mock.Anything).
		RunAndReturn(func(rskAddress string) (blockchain.PegInAddress, error) {
			fixture.mutex.Lock()
			defer fixture.mutex.Unlock()
			deposit, registered := fixture.deposits[rskAddress]
			if !registered {
				return blockchain.PegInAddress{}, fmt.Errorf("no registration for %s", rskAddress)
			}
			return blockchain.PegInAddress{
				Payload:  deposit.payload,
				Encoding: blockchain.PegInAddressRegistryEncodingBase58,
			}, nil
		})
	fixture.expectAddressImports()
	return fixture
}

func (fixture *registryRestartFixture) expectAddressImports() {
	fixture.t.Helper()
	// A node refuses a repeated import, and what the scanner does with that refusal is what several
	// of the restart cases turn on.
	fixture.wallet.EXPECT().ImportAddress(mock.Anything).RunAndReturn(func(address string) error {
		fixture.mutex.Lock()
		defer fixture.mutex.Unlock()
		alreadyImported := false
		for _, imported := range fixture.imports {
			alreadyImported = alreadyImported || imported == address
		}
		fixture.imports = append(fixture.imports, address)
		if alreadyImported {
			return errors.New("address already imported")
		}
		return nil
	})
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Maybe()
	fixture.wallet.EXPECT().RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{StartHeight: 100, StopHeight: 200}, nil).
		Maybe()
}

// The deposit address comes from the shared dataset, so assertions compare against a known-good
// constant instead of a value this test re-derives with the code under test.
func (fixture *registryRestartFixture) chainRegisters(
	txHash string,
	blockNumber uint64,
	logIndex uint,
) blockchain.AddressRegistered {
	fixture.t.Helper()
	registrationIndex := fixture.chain.registrationCount()
	vectors := pinnedLBCRegistrations(fixture.t)
	require.Less(
		fixture.t,
		registrationIndex,
		len(vectors),
		"acceptance fixtures must use only pinned LBC registration vectors",
	)
	registration := vectors[registrationIndex]
	fixture.mutex.Lock()
	if _, exists := fixture.deposits[registration.rskAddress]; !exists {
		fixture.deposits[registration.rskAddress] = knownDepositAddress(registrationIndex)
	}
	fixture.mutex.Unlock()
	event := blockchain.AddressRegistered{
		TxHash:           txHash,
		RskAddress:       registration.rskAddress,
		RegistrationRoot: registration.root,
		BlockNumber:      blockNumber,
		LogIndex:         logIndex,
	}
	fixture.chain.add(event)
	return event
}

func (fixture *registryRestartFixture) depositAddressOf(event blockchain.AddressRegistered) string {
	fixture.mutex.Lock()
	defer fixture.mutex.Unlock()
	return fixture.deposits[event.RskAddress].address
}

func (fixture *registryRestartFixture) importedAddresses() []string {
	fixture.mutex.Lock()
	defer fixture.mutex.Unlock()
	return append([]string(nil), fixture.imports...)
}

func (fixture *registryRestartFixture) watchSet() []rootstock.PegInAddressRegistryWatchEntry {
	fixture.t.Helper()
	entries, err := fixture.store.List(context.Background())
	require.NoError(fixture.t, err)
	return entries
}

func (fixture *registryRestartFixture) checkpointSnapshot() (
	rootstock.PegInAddressRegistryWatchCheckpoint,
	bool,
) {
	fixture.store.mutex.Lock()
	defer fixture.store.mutex.Unlock()
	return fixture.store.checkpoint, fixture.store.hasCheckpoint
}

func (fixture *registryRestartFixture) checkpointAdvances() []uint64 {
	return fixture.store.checkpointAdvances()
}

func (fixture *registryRestartFixture) rootstockHeightCalls() int {
	fixture.mutex.Lock()
	defer fixture.mutex.Unlock()
	return fixture.heightCalls
}

func (fixture *registryRestartFixture) rootReadBlocks() []uint64 {
	fixture.mutex.Lock()
	defer fixture.mutex.Unlock()
	return append([]uint64(nil), fixture.rootReads...)
}

func (fixture *registryRestartFixture) startScanner() *watcherSession {
	fixture.t.Helper()
	ticker := newSessionTicker()
	return startWatcherSession(fixture.t, newRegistryWatcher(
		fixture.store,
		fixture.registry,
		fixture.rskRpc,
		fixture.btcNetwork,
		fixture.wallet,
		fixture.eventBus,
		ticker,
		fixture.startBlock,
		fixture.pageSize,
		fixture.finalityDepth,
	), ticker)
}

// Each call is a boot, one poll and a shutdown, so calling scanOnce twice is a restart.
func (fixture *registryRestartFixture) scanOnce() {
	fixture.t.Helper()
	session := fixture.startScanner()
	session.poll()
	session.stop()
}

// Wherever the process stops, one registration must converge on one imported watch entry.
//
//nolint:funlen // The three scanner boundaries share one world and one convergence assertion.
func TestPegInAddressRegistryWatcher_ConvergesAfterRestartAtEveryBoundary(t *testing.T) {
	const registrationBlock = 101
	tests := []struct {
		name string
		stop func(store *registryWatchStore)
		// imports counts every import the node saw across both lifetimes.
		imports int
	}{
		{
			name:    "stopped after the entry was persisted, before the import",
			stop:    func(store *registryWatchStore) { store.failGet = errProcessStopped },
			imports: 1,
		},
		{
			name: "stopped after the import, before the imported state was persisted",
			stop: func(store *registryWatchStore) { store.failUpdate = errProcessStopped },
			// The restart re-imports, and the node's refusal is what the scanner reads as success.
			imports: 2,
		},
		{
			name: "stopped after verify, before the checkpoint advanced",
			stop: func(store *registryWatchStore) { store.failSetCheckpoint = errProcessStopped },
			// Checkpoint publish now precedes finalize, so a failed checkpoint leaves the
			// entry discovered and the restart re-imports.
			imports: 2,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newRegistryRestartFixture(t, 100, 3, 2)
			event := fixture.chainRegisters("registration", registrationBlock, 4)
			fixture.chain.setHead(104)

			testCase.stop(fixture.store)
			fixture.scanOnce()
			fixture.scanOnce()

			requireImportedEntry(t, fixture, event)
			assert.Len(t, fixture.importedAddresses(), testCase.imports)
			assert.Equal(t, []uint64{102}, fixture.checkpointAdvances())
		})
	}
}

func requireImportedEntry(
	t *testing.T,
	fixture *registryRestartFixture,
	event blockchain.AddressRegistered,
) {
	t.Helper()
	entries := fixture.watchSet()
	require.Len(t, entries, 1, "a replayed registration must not produce a second watch entry")
	assert.Equal(t, event.TxHash, entries[0].TxHash)
	assert.Equal(t, event.LogIndex, entries[0].LogIndex)
	assert.Equal(t, rootstock.PegInAddressRegistryWatchImported, entries[0].State)
	assert.Equal(t, fixture.depositAddressOf(event), entries[0].BtcAddress)
	assert.Empty(t, entries[0].LastError)
}

// A node that delivers the same log twice, or delivers it a poll later than it should have, must
// not produce a second entry, a second import, a state regression, or a checkpoint that moves
// anywhere but forward.
func TestPegInAddressRegistryWatcher_AbsorbsDuplicateAndLateEventDelivery(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 100, 3, 2)
	fixture.chain.setHead(106)

	session := fixture.startScanner()
	defer session.stop()
	session.poll()
	assert.Empty(t, fixture.watchSet())

	// The node under-reported block 102 on the first poll and now returns its log twice; the second
	// registration is placed where the following poll's overlap window will deliver it again.
	late := fixture.chainRegisters("late-delivery", 102, 1)
	fixture.chain.add(late)
	repeated := fixture.chainRegisters("repeated-delivery", 105, 2)
	fixture.chain.setHead(110)
	session.poll()

	assert.Len(t, fixture.watchSet(), 2)
	session.poll()

	assert.Equal(t, [][2]uint64{
		{100, 102}, {103, 104}, {105, 106},
		{105, 107},
		{100, 102}, {103, 105}, {106, 108},
		{109, 110},
		{109, 110},
	}, fixture.chain.requestedRanges())
	entries := fixture.watchSet()
	require.Len(t, entries, 2)
	assert.Equal(t, late.TxHash, entries[0].TxHash)
	assert.Equal(t, repeated.TxHash, entries[1].TxHash)
	for _, entry := range entries {
		assert.Equal(t, rootstock.PegInAddressRegistryWatchImported, entry.State)
	}
	assert.Equal(
		t,
		[]string{fixture.depositAddressOf(late), fixture.depositAddressOf(repeated)},
		fixture.importedAddresses(),
	)
	assert.Equal(t, fixture.chain.rootAt(108), fixture.store.checkpoint.LocalRoot)
	assert.Equal(t, []uint64{104, 108}, fixture.checkpointAdvances())
}

func TestPegInAddressRegistryWatcher_DoesNotAcceptCheckpointAfterSkippedTailEvent(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 100, 2, 1)
	first := fixture.chainRegisters("first", 100, 0)
	skippedTail := fixture.chainRegisters("skipped-tail", 101, 0)
	fixture.chain.omitOnNextQuery(skippedTail)
	fixture.chain.setHead(102)

	fixture.scanOnce()

	requireImportedEntryCount(t, fixture, first, skippedTail)
	assert.Equal(t, fixture.chain.rootAt(101), fixture.store.checkpoint.LocalRoot)
}

func TestPegInAddressRegistryWatcher_SecondColdMismatchLeavesNoTrustedCheckpoint(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 100, 2, 0)
	fixture.chainRegisters("baseline", 100, 0)
	omitted := fixture.chainRegisters("omitted-twice", 100, 1)
	fixture.chain.omitOnNextQueries(omitted, 2)
	fixture.chain.setHead(100)

	fixture.scanOnce()

	assert.Equal(t, 2, fixture.store.checkpointClears(), "cold verification mismatch must invalidate its page checkpoint")
	assert.False(t, fixture.store.checkpointFound(), "restart must not trust the failed cold replay")

	fixture.chain.drop(omitted)
	fixture.scanOnce()
	replaysFromDeployment := 0
	for _, requestedRange := range fixture.chain.requestedRanges() {
		if requestedRange[0] == 100 {
			replaysFromDeployment++
		}
	}
	assert.Equal(t, 3, replaysFromDeployment, "restart must begin at the configured deployment block")
}

func TestPegInAddressRegistryWatcher_DoesNotCheckpointUnconfirmedLog(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 100, 10, 2)
	confirmed := fixture.chainRegisters("confirmed", 100, 0)
	unconfirmed := fixture.chainRegisters("unconfirmed", 103, 0)
	fixture.chain.includeOnNextQuery(unconfirmed)
	fixture.chain.setHead(104)

	fixture.scanOnce()

	requireImportedEntry(t, fixture, confirmed)
	assert.Equal(t, []string{fixture.depositAddressOf(confirmed)}, fixture.importedAddresses())
	assert.Equal(t, fixture.chain.rootAt(102), fixture.store.checkpoint.LocalRoot)
	assert.Equal(t, uint64(102), fixture.store.checkpoint.LastProcessedBlock)

	fixture.chain.drop(unconfirmed)
	fixture.chain.setHead(108)
	fixture.scanOnce()

	requireImportedEntry(t, fixture, confirmed)
	assert.Equal(t, fixture.chain.rootAt(106), fixture.store.checkpoint.LocalRoot)
}

func TestPegInAddressRegistryWatcher_UnconfirmedTailOmissionTriggersAuthoritativeReplay(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 100, 10, 2)
	confirmed := fixture.chainRegisters("confirmed", 100, 0)
	fixture.chainRegisters("unconfirmed-first", 103, 0)
	unconfirmedTail := fixture.chainRegisters("unconfirmed-tail", 104, 0)
	fixture.chain.omitOnNextQuery(unconfirmedTail)
	fixture.chain.setHead(104)

	fixture.scanOnce()

	requireImportedEntry(t, fixture, confirmed)
	assert.Equal(t, 1, fixture.store.checkpointClears())
	replaysFromDeployment := 0
	for _, requestedRange := range fixture.chain.requestedRanges() {
		if requestedRange[0] == 100 {
			replaysFromDeployment++
		}
	}
	assert.Equal(t, 2, replaysFromDeployment)
}

func TestPegInAddressRegistryWatcher_ReorgSignalDiscardsCheckpointAndReplaysFromDeployment(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 100, 2, 1)
	fixture.chainRegisters("existing", 100, 0)
	fixture.chain.setHead(102)
	session := fixture.startScanner()
	defer session.stop()
	session.poll()

	replaysFromDeployment := func() int {
		count := 0
		for _, requestedRange := range fixture.chain.requestedRanges() {
			if requestedRange[0] == 100 {
				count++
			}
		}
		return count
	}
	require.Equal(t, 1, replaysFromDeployment())
	fixture.eventBus.Publish(blockchain.NodeReorgCheckEvent{
		BaseEvent:    entities.NewBaseEvent(blockchain.NodeReorgCheckEventId),
		NodeType:     entities.NodeTypeRootstock,
		CurrentDepth: 1,
	})

	assert.Eventually(t, func() bool {
		return replaysFromDeployment() == 2 && fixture.store.checkpointClears() == 1
	}, time.Second, time.Millisecond)
}

func TestPegInAddressRegistryWatcher_ColdReplayWaitsForConfiguredStart(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 100, 10, 2)
	boundary := fixture.chainRegisters("boundary", 100, 0)

	fixture.chain.setHead(100)
	fixture.scanOnce()
	checkpoint, found := fixture.checkpointSnapshot()
	assert.False(t, found)
	assert.Zero(t, checkpoint)
	assert.Empty(t, fixture.chain.requestedRanges())

	fixture.chain.setHead(101)
	fixture.scanOnce()
	_, found = fixture.checkpointSnapshot()
	assert.False(t, found)
	assert.Empty(t, fixture.chain.requestedRanges())

	fixture.chain.setHead(102)
	fixture.scanOnce()
	checkpoint, found = fixture.checkpointSnapshot()
	require.True(t, found)
	assert.Equal(t, uint64(100), checkpoint.LastProcessedBlock)
	assert.Equal(t, [][2]uint64{{100, 100}, {101, 102}}, fixture.chain.requestedRanges())
	requireImportedEntry(t, fixture, boundary)
}

func TestPegInAddressRegistryWatcher_ReorgInvalidatesCheckpointBeforeShortHeadReturn(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 100, 2, 2)
	fixture.chainRegisters("existing", 100, 0)
	fixture.chain.setHead(104)
	session := fixture.startScanner()
	defer session.stop()
	session.poll()

	fixture.chain.setHead(1)
	fixture.eventBus.Publish(blockchain.NodeReorgCheckEvent{
		BaseEvent:    entities.NewBaseEvent(blockchain.NodeReorgCheckEventId),
		NodeType:     entities.NodeTypeRootstock,
		CurrentDepth: 1,
	})
	require.Eventually(t, func() bool {
		return fixture.rootstockHeightCalls() == 2
	}, time.Second, time.Millisecond, "reorg signal was not consumed")
	assert.Equal(t, 1, fixture.store.checkpointClears(), "accepted reorg must invalidate before the short-head return")

	fixture.chain.setHead(104)
	session.poll()
	replaysFromDeployment := 0
	for _, requestedRange := range fixture.chain.requestedRanges() {
		if requestedRange[0] == 100 {
			replaysFromDeployment++
		}
	}
	assert.Equal(t, 2, replaysFromDeployment, "the next run must remain forced into cold replay")
}

func requireImportedEntryCount(
	t *testing.T,
	fixture *registryRestartFixture,
	events ...blockchain.AddressRegistered,
) {
	t.Helper()
	entries := fixture.watchSet()
	require.Len(t, entries, len(events))
	for index, event := range events {
		assert.Equal(t, event.TxHash, entries[index].TxHash)
		assert.Equal(t, rootstock.PegInAddressRegistryWatchImported, entries[index].State)
	}
}

// The checkpoint may never advance over a log that can still be reorganised away.
func TestPegInAddressRegistryWatcher_DoesNotCompleteLogRemovedBeforeFinality(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 100, 10, 2)
	fixture.chain.setHead(104)
	reorgedOut := fixture.chainRegisters("reorged-out", 103, 0)

	session := fixture.startScanner()
	defer session.stop()
	session.poll()

	assert.Equal(t, [][2]uint64{{100, 102}, {103, 104}}, fixture.chain.requestedRanges())
	assert.Empty(t, fixture.watchSet(), "a log above the finalized head must not reach the watch set")
	assert.Empty(t, fixture.importedAddresses())

	fixture.chain.drop(reorgedOut)
	replacement := fixture.chainRegisters("replacement", 103, 0)
	fixture.chain.setHead(108)
	session.poll()

	entries := fixture.watchSet()
	require.Len(t, entries, 1)
	assert.Equal(t, replacement.TxHash, entries[0].TxHash)
	assert.Equal(t, rootstock.PegInAddressRegistryWatchImported, entries[0].State)
	assert.Equal(t, []string{fixture.depositAddressOf(replacement)}, fixture.importedAddresses())
	assert.Equal(
		t,
		[][2]uint64{{100, 102}, {103, 104}, {103, 106}, {107, 108}},
		fixture.chain.requestedRanges(),
	)
	assert.Equal(t, []uint64{102, 106}, fixture.checkpointAdvances())
}

// --- merged from pegin_address_registry_checkpoint_test.go ---

type checkpointRestartStore struct {
	mutex             sync.Mutex
	entries           []rootstock.PegInAddressRegistryWatchEntry
	checkpoint        rootstock.PegInAddressRegistryWatchCheckpoint
	hasCheckpoint     bool
	failSetCheckpoint error
	operations        []string
}

func (store *checkpointRestartStore) Upsert(_ context.Context, entry rootstock.PegInAddressRegistryWatchEntry) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.operations = append(store.operations, "entry-upsert")
	if store.entryIndex(entry.RskAddress) < 0 {
		store.entries = append(store.entries, entry)
	}
	return nil
}

func (store *checkpointRestartStore) Get(
	_ context.Context,
	rskAddress string,
) (*rootstock.PegInAddressRegistryWatchEntry, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	index := store.entryIndex(rskAddress)
	if index < 0 {
		return nil, nil
	}
	entry := store.entries[index]
	return &entry, nil
}

func (store *checkpointRestartStore) List(context.Context) ([]rootstock.PegInAddressRegistryWatchEntry, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	return append([]rootstock.PegInAddressRegistryWatchEntry(nil), store.entries...), nil
}

func (store *checkpointRestartStore) Update(_ context.Context, entry rootstock.PegInAddressRegistryWatchEntry) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	index := store.entryIndex(entry.RskAddress)
	if index < 0 {
		return errors.New("pegin address registry watch entry not found")
	}
	store.entries[index] = entry
	store.operations = append(store.operations, "entry-update")
	return nil
}

func (store *checkpointRestartStore) GetCheckpoint(
	context.Context,
) (rootstock.PegInAddressRegistryWatchCheckpoint, bool, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	return store.checkpoint, store.hasCheckpoint, nil
}

func (store *checkpointRestartStore) SetCheckpoint(
	_ context.Context,
	checkpoint rootstock.PegInAddressRegistryWatchCheckpoint,
) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.operations = append(store.operations, "checkpoint")
	if store.failSetCheckpoint != nil {
		err := store.failSetCheckpoint
		store.failSetCheckpoint = nil
		return err
	}
	store.checkpoint = checkpoint
	store.hasCheckpoint = true
	return nil
}

func (store *checkpointRestartStore) DeleteCheckpoint(context.Context) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.checkpoint = rootstock.PegInAddressRegistryWatchCheckpoint{}
	store.hasCheckpoint = false
	return nil
}

func (store *checkpointRestartStore) entryIndex(rskAddress string) int {
	for index := range store.entries {
		if store.entries[index].RskAddress == rskAddress {
			return index
		}
	}
	return -1
}

func (store *checkpointRestartStore) snapshot() (
	[]rootstock.PegInAddressRegistryWatchEntry,
	rootstock.PegInAddressRegistryWatchCheckpoint,
	[]string,
) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	return append([]rootstock.PegInAddressRegistryWatchEntry(nil), store.entries...),
		store.checkpoint,
		append([]string(nil), store.operations...)
}

func TestPegInAddressRegistryWatcher_RestartsAfterEntryWriteBeforeCheckpoint(t *testing.T) {
	scenario := newCheckpointRestartScenario(t)

	scenario.scanOnce(t)
	entries, checkpoint, _ := scenario.store.snapshot()
	require.Len(t, entries, 1)
	assert.Equal(t, rootstock.PegInAddressRegistryWatchDiscovered, entries[0].State)
	assert.Equal(t, uint64(100), checkpoint.LastProcessedBlock)
	assert.Equal(t, scenario.previousRoot, checkpoint.LocalRoot)

	scenario.scanOnce(t)
	entries, checkpoint, operations := scenario.store.snapshot()
	require.Len(t, entries, 1)
	assert.Equal(t, rootstock.PegInAddressRegistryWatchImported, entries[0].State)
	assert.Equal(t, rootstock.PegInAddressRegistryWatchCheckpoint{
		LocalRoot:          scenario.expectedRoot,
		LastProcessedBlock: 102,
	}, checkpoint)
	assert.Equal(t, []string{
		"entry-upsert",
		"checkpoint",
		"checkpoint",
		"entry-update",
	}, operations)
}

type checkpointRestartScenario struct {
	store        *checkpointRestartStore
	registry     *mocks.PegInAddressRegistryContractMock
	rskRpc       *mocks.RootstockRpcServerMock
	btcNetwork   *mocks.BtcRpcMock
	wallet       *mocks.BitcoinWalletMock
	previousRoot [32]byte
	expectedRoot [32]byte
}

func newCheckpointRestartScenario(t *testing.T) *checkpointRestartScenario {
	t.Helper()
	const (
		startBlock = uint64(100)
		eventBlock = uint64(101)
	)
	vectors := pinnedLBCRegistrations(t)
	previousRoot := vectors[0].root
	rskAddress := vectors[1].rskAddress
	expectedRoot := vectors[1].root
	event := blockchain.AddressRegistered{
		TxHash:           "registration",
		LogIndex:         1,
		BlockNumber:      eventBlock,
		RskAddress:       rskAddress,
		RegistrationRoot: expectedRoot,
	}
	deposit := knownDepositAddress(0)
	store := &checkpointRestartStore{
		checkpoint: rootstock.PegInAddressRegistryWatchCheckpoint{
			LocalRoot:          previousRoot,
			LastProcessedBlock: startBlock,
		},
		hasCheckpoint:     true,
		failSetCheckpoint: errProcessStopped,
	}
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(101), uint64Pointer(102)).
		Return([]blockchain.AddressRegistered{event}, nil).
		Twice()
	registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(103), uint64Pointer(103)).
		Return([]blockchain.AddressRegistered{}, nil).
		Twice()
	registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(103)).Return(expectedRoot, nil).Twice()
	registry.EXPECT().GetPegInAddress(rskAddress).
		Return(blockchain.PegInAddress{
			Payload:  deposit.payload,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Twice()
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(103), nil).Twice()
	wallet := mocks.NewBitcoinWalletMock(t)
	wallet.EXPECT().ImportAddress(deposit.address).Return(nil).Twice()
	btcNetwork := &mocks.BtcRpcMock{}
	t.Cleanup(func() { btcNetwork.AssertExpectations(t) })
	btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
	wallet.EXPECT().RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{StartHeight: 100, StopHeight: 200}, nil).
		Once()

	return &checkpointRestartScenario{
		store:        store,
		registry:     registry,
		rskRpc:       rskRpc,
		btcNetwork:   btcNetwork,
		wallet:       wallet,
		previousRoot: previousRoot,
		expectedRoot: expectedRoot,
	}
}

func (scenario *checkpointRestartScenario) scanOnce(t *testing.T) {
	t.Helper()
	ticker := newSessionTicker()
	watcher := newRegistryWatcher(
		scenario.store,
		scenario.registry,
		scenario.rskRpc,
		scenario.btcNetwork,
		scenario.wallet,
		nil,
		ticker,
		100,
		3,
		1,
	)
	session := startWatcherSession(t, watcher, ticker)
	session.poll()
	session.stop()
}

// --- merged from pegin_address_registry_acceptance_test.go ---

func TestFirstBootBackfillsRegistrationsAndStoresMatchingCheckpoint(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 100, 2, 1)
	capturedLogs := test.CaptureStructuredLogs(t)
	first := fixture.chainRegisters("first", 100, 0)
	second := fixture.chainRegisters("second", 101, 0)
	fixture.chain.setHead(103)

	session := fixture.startScanner()
	defer session.stop()
	session.poll()

	requireCompleteImportedWatchSet(t, fixture, first, second)
	requireMatchingCheckpoint(t, fixture, second.RegistrationRoot, 102)
	initialRanges := fixture.chain.requestedRanges()
	require.NotEmpty(t, initialRanges)
	assert.Equal(t, fixture.startBlock, initialRanges[0][0],
		"first boot must replay from the configured deployment block")
	assert.Equal(t, []uint64{103}, fixture.rootReadBlocks(), "the local root must be checked at the captured head")
	requireNoIntegritySignal(t, fixture, capturedLogs())

	requestCountBeforeIncrementalPoll := len(fixture.chain.requestedRanges())
	third := fixture.chainRegisters("third", 103, 0)
	fixture.chain.setHead(104)
	session.poll()

	requireCompleteImportedWatchSet(t, fixture, first, second, third)
	requireMatchingCheckpoint(t, fixture, third.RegistrationRoot, 103)
	incrementalRanges := fixture.chain.requestedRanges()[requestCountBeforeIncrementalPoll:]
	require.NotEmpty(t, incrementalRanges)
	assert.Equal(t, uint64(103), incrementalRanges[0][0], "normal operation must continue after the saved checkpoint")
	assert.Equal(t, []uint64{102, 103}, fixture.checkpointAdvances())
	assert.Equal(t, []uint64{103, 104}, fixture.rootReadBlocks())
	requireNoIntegritySignal(t, fixture, capturedLogs())
}

func TestRestartResumesFromTrustedCheckpointWithMatchingRoot(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 100, 2, 1)
	capturedLogs := test.CaptureStructuredLogs(t)
	first := fixture.chainRegisters("persisted-first", 100, 0)
	second := fixture.chainRegisters("persisted-second", 101, 0)
	fixture.chain.setHead(102)
	fixture.scanOnce()

	requireCompleteImportedWatchSet(t, fixture, first, second)
	requireMatchingCheckpoint(t, fixture, second.RegistrationRoot, 101)
	entriesBeforeRestart := fixture.watchSet()
	requestCountBeforeRestart := len(fixture.chain.requestedRanges())

	third := fixture.chainRegisters("after-restart", 102, 0)
	fixture.chain.setHead(103)
	fixture.scanOnce()

	entriesAfterRestart := requireCompleteImportedWatchSet(t, fixture, first, second, third)
	assert.Equal(t, entriesBeforeRestart, entriesAfterRestart[:len(entriesBeforeRestart)],
		"persisted imported entries must not regress during restart")
	assert.Len(t, fixture.importedAddresses(), 3, "restart must import only the new registration")
	restartRanges := fixture.chain.requestedRanges()[requestCountBeforeRestart:]
	require.NotEmpty(t, restartRanges)
	assert.Equal(t, uint64(102), restartRanges[0][0], "restart must resume after the trusted checkpoint")
	for _, requestedRange := range restartRanges {
		assert.NotEqual(t, uint64(100), requestedRange[0], "restart must not replay from the deployment block")
	}
	requireMatchingCheckpoint(t, fixture, third.RegistrationRoot, 102)
	assert.Equal(t, []uint64{101, 102}, fixture.checkpointAdvances())
	rootReads := fixture.rootReadBlocks()
	require.NotEmpty(t, rootReads)
	assert.Equal(t, uint64(103), rootReads[len(rootReads)-1], "restart must verify against the captured head")
	requireNoIntegritySignal(t, fixture, capturedLogs())
}

func TestMissedEventSignalsMismatchAndRepairsCompleteWatchSet(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 100, 2, 1)
	capturedLogs := test.CaptureStructuredLogs(t)
	first := fixture.chainRegisters("first", 100, 0)
	fixture.chain.setHead(101)
	fixture.scanOnce()
	requireCompleteImportedWatchSet(t, fixture, first)
	requireNoIntegritySignal(t, fixture, capturedLogs())

	skipped := fixture.chainRegisters("skipped", 101, 0)
	last := fixture.chainRegisters("last", 102, 0)
	fixture.chain.omitOnNextQuery(skipped)
	fixture.chain.setHead(103)
	fixture.scanOnce()

	requireCompleteImportedWatchSet(t, fixture, first, skipped, last)
	requireMatchingCheckpoint(t, fixture, last.RegistrationRoot, 102)
	assert.Equal(t, 2, replayCountFromDeployment(fixture),
		"the mismatch repair must add one cold replay from the deployment block")
	assert.Len(t, fixture.importedAddresses(), 3, "repair must not repeat the first registration import")
	mismatchSignal := requireOneIntegritySignal(t, fixture, last.BlockNumber, last.RegistrationRoot)
	logEntry := requireSingleMismatchLog(t, capturedLogs())
	assert.Equal(t, "error", logEntry.Level())
	assert.InDelta(t, float64(last.BlockNumber), logEntry.Field(test.LogKey("block_number")), 0)
	assert.Equal(t, fmt.Sprintf("0x%x", mismatchSignal.LocalRoot), logEntry.Field(test.LogKey("local_root")))
	assert.Equal(t, fmt.Sprintf("0x%x", last.RegistrationRoot), logEntry.Field(test.LogKey("chain_root")))
	assert.NotEqual(t, logEntry.Field(test.LogKey("chain_root")), logEntry.Field(test.LogKey("local_root")))
	assert.Equal(t, "event_last_0", logEntry.Field(test.LogKey("source")))
}

func TestSilentEventStreamHealthCheckSignalsMismatchAndReplays(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 100, 2, 1)
	capturedLogs := test.CaptureStructuredLogs(t)
	first := fixture.chainRegisters("first", 100, 0)
	fixture.chain.setHead(101)

	session := fixture.startScanner()
	defer session.stop()
	session.poll()
	requireCompleteImportedWatchSet(t, fixture, first)
	requireMatchingCheckpoint(t, fixture, first.RegistrationRoot, 100)

	session.poll()
	requireCompleteImportedWatchSet(t, fixture, first)
	requireMatchingCheckpoint(t, fixture, first.RegistrationRoot, 100)
	assert.Equal(t, 1, replayCountFromDeployment(fixture), "an equal-root tick must not start a cold replay")
	assert.Equal(t, []uint64{100}, fixture.checkpointAdvances(), "an equal-root tick must not republish progress")
	requireNoIntegritySignal(t, fixture, capturedLogs())

	// The node's first reading of finalized block 100 omitted this registration. The head and
	// trusted checkpoint cannot advance, so only the same-head root check can expose the gap.
	silentlyOmitted := fixture.chainRegisters("silently-omitted", 100, 1)
	session.poll()

	requireCompleteImportedWatchSet(t, fixture, first, silentlyOmitted)
	requireMatchingCheckpoint(t, fixture, silentlyOmitted.RegistrationRoot, 100)
	assert.Equal(t, 2, replayCountFromDeployment(fixture),
		"the health mismatch must use the cold replay procedure")
	assert.Len(t, fixture.importedAddresses(), 2, "health repair must not repeat the first registration import")
	mismatchSignal := requireOneIntegritySignal(t, fixture, 101, silentlyOmitted.RegistrationRoot)
	assert.Equal(t, first.RegistrationRoot, mismatchSignal.LocalRoot)
	logEntry := requireSingleMismatchLog(t, capturedLogs())
	assert.Equal(t, "error", logEntry.Level())
	assert.InDelta(t, float64(101), logEntry.Field(test.LogKey("block_number")), 0)
	assert.Equal(t, fmt.Sprintf("0x%x", first.RegistrationRoot), logEntry.Field(test.LogKey("local_root")))
	assert.Equal(t, fmt.Sprintf("0x%x", silentlyOmitted.RegistrationRoot), logEntry.Field(test.LogKey("chain_root")))
	assert.Equal(t, "captured_head", logEntry.Field(test.LogKey("source")))
}

func requireCompleteImportedWatchSet(
	t *testing.T,
	fixture *registryRestartFixture,
	events ...blockchain.AddressRegistered,
) []rootstock.PegInAddressRegistryWatchEntry {
	t.Helper()
	entries := fixture.watchSet()
	require.Len(t, entries, len(events))
	expectedImports := make([]string, 0, len(events))
	for index, event := range events {
		assert.Equal(t, event.TxHash, entries[index].TxHash)
		assert.Equal(t, event.LogIndex, entries[index].LogIndex)
		assert.Equal(t, event.BlockNumber, entries[index].BlockNumber)
		assert.Equal(t, event.RskAddress, entries[index].RskAddress)
		assert.Equal(t, event.RegistrationRoot, entries[index].RegistrationRoot)
		assert.Equal(t, rootstock.PegInAddressRegistryWatchImported, entries[index].State)
		assert.Equal(t, fixture.depositAddressOf(event), entries[index].BtcAddress)
		assert.Empty(t, entries[index].LastError)
		expectedImports = append(expectedImports, fixture.depositAddressOf(event))
	}
	assert.Equal(t, expectedImports, fixture.importedAddresses())
	return entries
}

func requireMatchingCheckpoint(
	t *testing.T,
	fixture *registryRestartFixture,
	expectedRoot [32]byte,
	expectedBlock uint64,
) {
	t.Helper()
	checkpoint, found := fixture.checkpointSnapshot()
	require.True(t, found)
	assert.Equal(t, expectedBlock, checkpoint.LastProcessedBlock)
	assert.Equal(t, expectedRoot, checkpoint.LocalRoot)
	assert.Equal(t, fixture.chain.rootAt(expectedBlock), checkpoint.LocalRoot,
		"the checkpoint and chain roots must describe the same block")
}

func replayCountFromDeployment(fixture *registryRestartFixture) int {
	count := 0
	for _, requestedRange := range fixture.chain.requestedRanges() {
		if requestedRange[0] == fixture.startBlock {
			count++
		}
	}
	return count
}

func requireNoIntegritySignal(
	t *testing.T,
	fixture *registryRestartFixture,
	logEntries []test.LogEntry,
) {
	t.Helper()
	assert.Zero(t, fixture.eventBus.publishedCount(blockchain.PegInAddressRegistryRootMismatchEventId))
	assert.Zero(t, fixture.eventBus.publishedCount(blockchain.PegInAddressRegistryResyncStartedEventId))
	assert.Empty(t, mismatchLogs(logEntries))
}

func requireOneIntegritySignal(
	t *testing.T,
	fixture *registryRestartFixture,
	expectedBlock uint64,
	expectedChainRoot [32]byte,
) blockchain.PegInAddressRegistryRootMismatchEvent {
	t.Helper()
	mismatchEvents := fixture.eventBus.publishedEvents(blockchain.PegInAddressRegistryRootMismatchEventId)
	require.Len(t, mismatchEvents, 1)
	mismatch, ok := mismatchEvents[0].(blockchain.PegInAddressRegistryRootMismatchEvent)
	require.True(t, ok, "the metrics watcher requires the concrete root-mismatch event type")
	assert.Equal(t, expectedBlock, mismatch.BlockNumber)
	assert.Equal(t, expectedChainRoot, mismatch.ChainRoot)
	assert.NotEqual(t, mismatch.LocalRoot, mismatch.ChainRoot)

	resyncEvents := fixture.eventBus.publishedEvents(blockchain.PegInAddressRegistryResyncStartedEventId)
	require.Len(t, resyncEvents, 1)
	resync, ok := resyncEvents[0].(blockchain.PegInAddressRegistryResyncStartedEvent)
	require.True(t, ok, "the metrics watcher requires the concrete resync event type")
	assert.Equal(t, "root_mismatch", resync.Reason)
	return mismatch
}

func requireSingleMismatchLog(t *testing.T, logEntries []test.LogEntry) test.LogEntry {
	t.Helper()
	entries := mismatchLogs(logEntries)
	require.Len(t, entries, 1, "one mismatch must produce one structured log")
	return entries[0]
}

func mismatchLogs(logEntries []test.LogEntry) []test.LogEntry {
	entries := make([]test.LogEntry, 0)
	for _, entry := range logEntries {
		if entry.Message() == "PegIn address registry root mismatch" {
			entries = append(entries, entry)
		}
	}
	return entries
}
