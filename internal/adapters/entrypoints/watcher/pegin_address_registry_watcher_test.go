package watcher_test

import (
	"context"
	"math/big"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/crypto"
	entrypoint "github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	usecase "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test/datasets"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	watcherAddressA = "0x00000000000000000000000000000000000000a1"
	watcherAddressB = "0x00000000000000000000000000000000000000b2"
	// watcherStartBlock is the configured lower scan bound shared by every fixture.
	watcherStartBlock = uint64(100)
)

type memoryPegInWatchRepository struct {
	mu        sync.Mutex
	rows      []rootstock.PegInWatch
	updateErr error
}

func (repository *memoryPegInWatchRepository) Upsert(
	_ context.Context,
	watch rootstock.PegInWatch,
) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	for index := range repository.rows {
		if repository.rows[index].RskAddress == watch.RskAddress {
			return nil
		}
	}
	repository.rows = append(repository.rows, watch)
	return nil
}

func (repository *memoryPegInWatchRepository) Get(
	_ context.Context,
	rskAddress string,
) (*rootstock.PegInWatch, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	for index := range repository.rows {
		if repository.rows[index].RskAddress == rskAddress {
			watch := repository.rows[index]
			return &watch, nil
		}
	}
	return nil, nil
}

func (repository *memoryPegInWatchRepository) List(context.Context) ([]rootstock.PegInWatch, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	rows := append([]rootstock.PegInWatch(nil), repository.rows...)
	sort.Slice(rows, func(first, second int) bool {
		if rows[first].BlockNumber == rows[second].BlockNumber {
			return rows[first].LogIndex < rows[second].LogIndex
		}
		return rows[first].BlockNumber < rows[second].BlockNumber
	})
	return rows, nil
}

func (repository *memoryPegInWatchRepository) Update(
	_ context.Context,
	watch rootstock.PegInWatch,
) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	if repository.updateErr != nil {
		return repository.updateErr
	}
	for index := range repository.rows {
		if repository.rows[index].RskAddress == watch.RskAddress {
			repository.rows[index] = watch
			return nil
		}
	}
	return nil
}

func (repository *memoryPegInWatchRepository) DeleteFromBlock(
	_ context.Context,
	fromBlock uint64,
) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	kept := repository.rows[:0]
	for _, watch := range repository.rows {
		if watch.BlockNumber < fromBlock {
			kept = append(kept, watch)
		}
	}
	repository.rows = kept
	return nil
}

func (repository *memoryPegInWatchRepository) seed(watches ...rootstock.PegInWatch) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	repository.rows = append(repository.rows, watches...)
}

func (repository *memoryPegInWatchRepository) setUpdateError(err error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	repository.updateErr = err
}

func (repository *memoryPegInWatchRepository) state(rskAddress string) rootstock.PegInWatchState {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	for _, watch := range repository.rows {
		if watch.RskAddress == rskAddress {
			return watch.State
		}
	}
	return ""
}

func (repository *memoryPegInWatchRepository) rowsSnapshot(t *testing.T) []rootstock.PegInWatch {
	t.Helper()
	rows, err := repository.List(context.Background())
	require.NoError(t, err)
	return rows
}

func (repository *memoryPegInWatchRepository) clear() {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	repository.rows = nil
}

type watcherChainState struct {
	mu              sync.RWMutex
	head            uint64
	heightErr       error
	roots           map[uint64][32]byte
	events          []blockchain.AddressRegistered
	requested       [][2]uint64
	omit            map[string]int
	heightHook      func(uint64)
	rootReadStarted chan struct{}
	releaseRootRead chan struct{}
	rootReads       atomic.Int64
	activeRootReads atomic.Int64
	maxRootReads    atomic.Int64
}

func (state *watcherChainState) height(context.Context) (uint64, error) {
	state.mu.RLock()
	head, err, hook := state.head, state.heightErr, state.heightHook
	state.mu.RUnlock()
	if hook != nil {
		hook(head)
	}
	return head, err
}

func (state *watcherChainState) root(
	_ context.Context,
	blockNumber uint64,
) ([32]byte, error) {
	state.mu.RLock()
	root := state.roots[blockNumber]
	started, release := state.rootReadStarted, state.releaseRootRead
	state.mu.RUnlock()

	call := state.rootReads.Add(1)
	active := state.activeRootReads.Add(1)
	for {
		maximum := state.maxRootReads.Load()
		if active <= maximum || state.maxRootReads.CompareAndSwap(maximum, active) {
			break
		}
	}
	defer state.activeRootReads.Add(-1)
	if call == 1 && started != nil {
		close(started)
		<-release
	}
	return root, nil
}

func (state *watcherChainState) registeredEvents(
	_ context.Context,
	fromBlock uint64,
	toBlock *uint64,
) ([]blockchain.AddressRegistered, error) {
	state.mu.Lock()
	defer state.mu.Unlock()
	to := state.head
	if toBlock != nil {
		to = *toBlock
	}
	state.requested = append(state.requested, [2]uint64{fromBlock, to})
	result := make([]blockchain.AddressRegistered, 0)
	for _, event := range state.events {
		if event.BlockNumber >= fromBlock && event.BlockNumber <= to {
			identity := event.TxHash + "/" + strconv.FormatUint(uint64(event.LogIndex), 10)
			if state.omit[identity] > 0 {
				state.omit[identity]--
				continue
			}
			result = append(result, event)
		}
	}
	return result, nil
}

func (state *watcherChainState) set(
	head uint64,
	roots map[uint64][32]byte,
	events []blockchain.AddressRegistered,
) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.head = head
	state.heightErr = nil
	state.roots = roots
	state.events = events
}

func (state *watcherChainState) setHeightHook(hook func(uint64)) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.heightHook = hook
}

func (state *watcherChainState) blockFirstRootRead() (<-chan struct{}, chan<- struct{}) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.rootReadStarted = make(chan struct{})
	state.releaseRootRead = make(chan struct{})
	return state.rootReadStarted, state.releaseRootRead
}

func (state *watcherChainState) omitNext(event blockchain.AddressRegistered) {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.omit == nil {
		state.omit = make(map[string]int)
	}
	state.omit[event.TxHash+"/"+strconv.FormatUint(uint64(event.LogIndex), 10)]++
}

func (state *watcherChainState) requestedRanges() [][2]uint64 {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return append([][2]uint64(nil), state.requested...)
}

func (state *watcherChainState) clearRequestedRanges() {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.requested = nil
}

func (state *watcherChainState) setHeightError(err error) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.heightErr = err
}

type pegInWatcherFixture struct {
	target     *entrypoint.PegInWatcher
	repository *memoryPegInWatchRepository
	chain      *watcherChainState
	registry   *mocks.PegInAddressRegistryContractMock
	rskRpc     *mocks.RootstockRpcServerMock
	btcNetwork *mocks.BtcRpcMock
	wallet     *mocks.BitcoinWalletMock
	eventBus   *mocks.EventBusMock
	ticker     *mocks.TickerMock
	ticks      chan time.Time
	reorgs     chan entities.Event
	startBlock uint64
	pageSize   uint64
}

func newPegInWatcherFixture(t *testing.T, withEventBus bool) *pegInWatcherFixture {
	return newPegInWatcherFixtureWithPaging(t, withEventBus, 10)
}

func newPegInWatcherFixtureWithPaging(
	t *testing.T,
	withEventBus bool,
	pageSize uint64,
) *pegInWatcherFixture {
	t.Helper()
	repository := &memoryPegInWatchRepository{}
	chain := &watcherChainState{roots: make(map[uint64][32]byte)}
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	btcNetwork := &mocks.BtcRpcMock{}
	wallet := mocks.NewBitcoinWalletMock(t)
	ticker := mocks.NewTickerMock(t)
	rskRpc.EXPECT().GetHeight(mock.Anything).RunAndReturn(chain.height).Maybe()
	registry.EXPECT().
		GetRegistrationRoot(mock.Anything, mock.Anything).
		RunAndReturn(chain.root).
		Maybe()
	registry.EXPECT().
		GetAddressRegisteredEvents(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(chain.registeredEvents).
		Maybe()

	var eventBus *mocks.EventBusMock
	var bus entities.EventBus
	if withEventBus {
		eventBus = &mocks.EventBusMock{}
		bus = eventBus
	}
	fixture := &pegInWatcherFixture{
		repository: repository,
		chain:      chain,
		registry:   registry,
		rskRpc:     rskRpc,
		btcNetwork: btcNetwork,
		wallet:     wallet,
		eventBus:   eventBus,
		ticker:     ticker,
		ticks:      make(chan time.Time, 1),
		reorgs:     make(chan entities.Event, 10),
		startBlock: watcherStartBlock,
		pageSize:   pageSize,
	}
	fixture.resetTarget(t, bus)
	return fixture
}

func (fixture *pegInWatcherFixture) resetTarget(t *testing.T, bus entities.EventBus) {
	t.Helper()
	replay, err := usecase.NewReplayRegisteredAddressesUseCase(
		fixture.repository,
		fixture.registry,
		fixture.rskRpc,
		crypto.Keccak256,
		fixture.startBlock,
		fixture.pageSize,
	)
	require.NoError(t, err)
	fixture.target = entrypoint.NewPegInWatcher(
		replay,
		usecase.NewDiscoverRegisteredAddressUseCase(fixture.repository, fixture.registry, fixture.wallet),
		usecase.NewGetPendingRegisteredAddressImportsUseCase(fixture.repository),
		usecase.NewFinalizeRegisteredAddressImportUseCase(fixture.repository),
		fixture.btcNetwork,
		fixture.wallet,
		bus,
		fixture.ticker,
	)
}

func (fixture *pegInWatcherFixture) start(t *testing.T) {
	t.Helper()
	fixture.ticker.EXPECT().C().Return(fixture.ticks).Maybe()
	fixture.ticker.EXPECT().Stop().Return().Once()
	if fixture.eventBus != nil {
		fixture.eventBus.
			On("Subscribe", blockchain.NodeReorgCheckEventId).
			Return((<-chan entities.Event)(fixture.reorgs)).
			Once()
	}
	go fixture.target.Start()
}

func (fixture *pegInWatcherFixture) stop(t *testing.T) {
	t.Helper()
	closed := make(chan bool, 1)
	fixture.target.Shutdown(closed)
	<-closed
	assert.Eventually(
		t,
		func() bool {
			return countCallsSince(fixture.ticker.Calls, 0, "Stop", nil) == 1
		},
		time.Second,
		10*time.Millisecond,
	)
}

func (fixture *pegInWatcherFixture) tick() {
	fixture.ticks <- time.Now()
}

func (fixture *pegInWatcherFixture) setChain(
	head uint64,
	events ...blockchain.AddressRegistered,
) {
	ordered := append([]blockchain.AddressRegistered(nil), events...)
	sort.Slice(ordered, func(first, second int) bool {
		if ordered[first].BlockNumber == ordered[second].BlockNumber {
			return ordered[first].LogIndex < ordered[second].LogIndex
		}
		return ordered[first].BlockNumber < ordered[second].BlockNumber
	})
	roots := make(map[uint64][32]byte)
	var root [32]byte
	eventIndex := 0
	for block := fixture.startBlock; block <= head; block++ {
		for eventIndex < len(ordered) && ordered[eventIndex].BlockNumber == block {
			root = ordered[eventIndex].RegistrationRoot
			eventIndex++
		}
		roots[block] = root
		if block == ^uint64(0) {
			break
		}
	}
	fixture.chain.set(head, roots, ordered)
}

func (fixture *pegInWatcherFixture) expectSupportedImport(
	t *testing.T,
	event blockchain.AddressRegistered,
	addressIndex int,
) {
	t.Helper()
	const checksumSize = 4
	decoded := datasets.Base58Addresses[addressIndex]
	payload := make([]byte, 0, len(decoded.Expected)+checksumSize)
	payload = append(payload, decoded.Expected...)
	payload = append(payload, chainhash.DoubleHashB(decoded.Expected)[:checksumSize]...)
	fixture.registry.EXPECT().
		GetPegInAddress(event.RskAddress).
		Return(blockchain.PegInAddress{
			Payload: payload, Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.wallet.EXPECT().ImportAddress(decoded.Address).Return(nil).Once()
}

func watcherRoot(t *testing.T, previous [32]byte, rskAddress string) [32]byte {
	t.Helper()
	root, err := blockchain.FoldPegInAddressRegistryRoot(crypto.Keccak256, previous, rskAddress)
	require.NoError(t, err)
	return root
}

func watcherEvent(
	t *testing.T,
	blockNumber uint64,
	logIndex uint,
	rskAddress string,
	previousRoot [32]byte,
) blockchain.AddressRegistered {
	t.Helper()
	return blockchain.AddressRegistered{
		TxHash:           rskAddress,
		LogIndex:         logIndex,
		BlockNumber:      blockNumber,
		RskAddress:       rskAddress,
		RegistrationRoot: watcherRoot(t, previousRoot, rskAddress),
	}
}

func watcherDepositPayload() ([]byte, string) {
	const checksumSize = 4
	decoded := datasets.Base58Addresses[0]
	payload := make([]byte, 0, len(decoded.Expected)+checksumSize)
	payload = append(payload, decoded.Expected...)
	payload = append(payload, chainhash.DoubleHashB(decoded.Expected)[:checksumSize]...)
	return payload, decoded.Address
}

func countCallsSince(calls []mock.Call, first int, method string, matches func(mock.Arguments) bool) int {
	count := 0
	for _, call := range calls[first:] {
		if call.Method == method && (matches == nil || matches(call.Arguments)) {
			count++
		}
	}
	return count
}

func publishedEvents(t *testing.T, bus *mocks.EventBusMock) []entities.Event {
	t.Helper()
	result := make([]entities.Event, 0)
	for _, call := range bus.Calls {
		if call.Method != "Publish" {
			continue
		}
		event, ok := call.Arguments.Get(0).(entities.Event)
		require.True(t, ok)
		result = append(result, event)
	}
	return result
}

func TestPegInWatcher_PrepareReplaysWithoutStoredCheckpoint(t *testing.T) {
	fixture := newPegInWatcherFixture(t, false)
	fixture.chain.set(100, map[uint64][32]byte{100: {}}, nil)
	require.NoError(t, fixture.target.Prepare(context.Background()))

	event := watcherEvent(t, 101, 1, watcherAddressA, [32]byte{})
	fixture.registry.EXPECT().
		GetPegInAddress(watcherAddressA).
		Return(blockchain.PegInAddress{Encoding: blockchain.PegInAddressRegistryEncodingBech32}, nil).
		Once()
	firstCall := len(fixture.registry.Calls)
	fixture.chain.set(102, map[uint64][32]byte{
		100: {},
		101: event.RegistrationRoot,
		102: event.RegistrationRoot,
	}, []blockchain.AddressRegistered{event})

	require.NoError(t, fixture.target.Prepare(context.Background()))

	assert.Equal(t, 1, countCallsSince(fixture.registry.Calls, firstCall, "GetRegistrationRoot",
		func(arguments mock.Arguments) bool { return arguments.Get(1) == uint64(101) }))
}

func TestPegInWatcher_NextScanUsesSuccessfulReplayCheckpoint(t *testing.T) {
	fixture := newPegInWatcherFixture(t, false)
	fixture.chain.set(100, map[uint64][32]byte{100: {}}, nil)
	require.NoError(t, fixture.target.Prepare(context.Background()))

	event := watcherEvent(t, 101, 1, watcherAddressA, [32]byte{})
	fixture.registry.EXPECT().
		GetPegInAddress(watcherAddressA).
		Return(blockchain.PegInAddress{Encoding: blockchain.PegInAddressRegistryEncodingBech32}, nil).
		Once()
	firstCall := len(fixture.registry.Calls)
	fixture.chain.set(102, map[uint64][32]byte{
		100: {},
		101: event.RegistrationRoot,
		102: event.RegistrationRoot,
	}, []blockchain.AddressRegistered{event})
	fixture.start(t)
	defer fixture.stop(t)

	fixture.tick()

	assert.Eventually(t, func() bool {
		return fixture.repository.state(watcherAddressA) == rootstock.PegInWatchUnsupportedEncoding
	}, time.Second, 10*time.Millisecond)
	assert.Equal(t, 1, countCallsSince(fixture.registry.Calls, firstCall, "GetRegistrationRoot",
		func(arguments mock.Arguments) bool { return arguments.Get(1) == uint64(100) }))
	assert.Zero(t, countCallsSince(fixture.registry.Calls, firstCall, "GetRegistrationRoot",
		func(arguments mock.Arguments) bool { return arguments.Get(1) == uint64(101) }))
	assert.Equal(t, 1, countCallsSince(fixture.registry.Calls, firstCall, "GetAddressRegisteredEvents",
		func(arguments mock.Arguments) bool { return arguments.Get(1) == uint64(101) }))
}

func TestPegInWatcher_PendingRowsAreDiscoveredRescannedOnceAndFinalized(t *testing.T) {
	fixture := newPegInWatcherFixture(t, false)
	fixture.chain.set(99, nil, nil)
	require.NoError(t, fixture.target.Prepare(context.Background()))

	first := watcherEvent(t, 100, 1, watcherAddressA, [32]byte{})
	second := watcherEvent(t, 100, 2, watcherAddressB, first.RegistrationRoot)
	fixture.repository.seed(
		rootstock.NewPegInWatch(first.TxHash, first.LogIndex, first.BlockNumber, first.RskAddress, "", first.RegistrationRoot),
		rootstock.NewPegInWatch(second.TxHash, second.LogIndex, second.BlockNumber, second.RskAddress, "", second.RegistrationRoot),
	)
	fixture.chain.set(100, map[uint64][32]byte{100: second.RegistrationRoot}, nil)
	payload, btcAddress := watcherDepositPayload()
	fixture.registry.EXPECT().
		GetPegInAddress(watcherAddressA).
		Return(blockchain.PegInAddress{
			Payload: payload, Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.registry.EXPECT().
		GetPegInAddress(watcherAddressB).
		Return(blockchain.PegInAddress{
			Payload: payload, Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.wallet.EXPECT().ImportAddress(btcAddress).Return(nil).Twice()
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(500), nil).Once()
	fixture.wallet.EXPECT().
		RescanBlockchain(int64(400)).
		Return(blockchain.BitcoinRescanResult{}, nil).
		Once()
	fixture.start(t)
	defer fixture.stop(t)

	fixture.tick()

	assert.Eventually(t, func() bool {
		return fixture.repository.state(watcherAddressA) == rootstock.PegInWatchImported &&
			fixture.repository.state(watcherAddressB) == rootstock.PegInWatchImported
	}, time.Second, 10*time.Millisecond)
	fixture.registry.AssertNumberOfCalls(t, "GetPegInAddress", 2)
	fixture.wallet.AssertNumberOfCalls(t, "ImportAddress", 2)
	fixture.wallet.AssertNumberOfCalls(t, "RescanBlockchain", 1)
}

func TestPegInWatcher_ReplayErrorSkipsDiscoverAndKeepsCheckpoint(t *testing.T) {
	fixture := newPegInWatcherFixture(t, false)
	fixture.chain.set(100, map[uint64][32]byte{100: {}}, nil)
	require.NoError(t, fixture.target.Prepare(context.Background()))
	fixture.start(t)
	defer fixture.stop(t)

	rpcCalls := len(fixture.rskRpc.Calls)
	fixture.chain.setHeightError(assert.AnError)
	fixture.tick()
	assert.Eventually(t, func() bool {
		return countCallsSince(fixture.rskRpc.Calls, rpcCalls, "GetHeight", nil) == 1
	}, time.Second, 10*time.Millisecond)
	fixture.registry.AssertNotCalled(t, "GetPegInAddress", mock.Anything)

	event := watcherEvent(t, 101, 1, watcherAddressA, [32]byte{})
	fixture.registry.EXPECT().
		GetPegInAddress(watcherAddressA).
		Return(blockchain.PegInAddress{Encoding: blockchain.PegInAddressRegistryEncodingBech32}, nil).
		Once()
	firstCall := len(fixture.registry.Calls)
	fixture.chain.set(102, map[uint64][32]byte{
		100: {},
		101: event.RegistrationRoot,
		102: event.RegistrationRoot,
	}, []blockchain.AddressRegistered{event})
	fixture.tick()

	assert.Eventually(t, func() bool {
		return fixture.repository.state(watcherAddressA) == rootstock.PegInWatchUnsupportedEncoding
	}, time.Second, 10*time.Millisecond)
	assert.Equal(t, 1, countCallsSince(fixture.registry.Calls, firstCall, "GetRegistrationRoot",
		func(arguments mock.Arguments) bool { return arguments.Get(1) == uint64(100) }))
	assert.Zero(t, countCallsSince(fixture.registry.Calls, firstCall, "GetRegistrationRoot",
		func(arguments mock.Arguments) bool { return arguments.Get(1) == uint64(101) }))
}

func TestPegInWatcher_DiscoverErrorKeepsSuccessfulReplayCheckpoint(t *testing.T) {
	fixture := newPegInWatcherFixture(t, false)
	fixture.chain.set(100, map[uint64][32]byte{100: {}}, nil)
	require.NoError(t, fixture.target.Prepare(context.Background()))
	fixture.start(t)
	defer fixture.stop(t)

	first := watcherEvent(t, 101, 1, watcherAddressA, [32]byte{})
	fixture.chain.set(101, map[uint64][32]byte{
		100: {},
		101: first.RegistrationRoot,
	}, []blockchain.AddressRegistered{first})
	fixture.registry.EXPECT().
		GetPegInAddress(watcherAddressA).
		Return(blockchain.PegInAddress{}, assert.AnError).
		Once()
	fixture.repository.setUpdateError(assert.AnError)
	fixture.tick()
	assert.Eventually(t, func() bool {
		return countCallsSince(fixture.registry.Calls, 0, "GetPegInAddress", nil) == 1
	}, time.Second, 10*time.Millisecond)

	fixture.repository.setUpdateError(nil)
	second := watcherEvent(t, 102, 2, watcherAddressB, first.RegistrationRoot)
	fixture.registry.EXPECT().
		GetPegInAddress(watcherAddressA).
		Return(blockchain.PegInAddress{Encoding: blockchain.PegInAddressRegistryEncodingBech32}, nil).
		Once()
	fixture.registry.EXPECT().
		GetPegInAddress(watcherAddressB).
		Return(blockchain.PegInAddress{Encoding: blockchain.PegInAddressRegistryEncodingBech32}, nil).
		Once()
	firstCall := len(fixture.registry.Calls)
	fixture.chain.set(103, map[uint64][32]byte{
		100: {},
		101: first.RegistrationRoot,
		102: second.RegistrationRoot,
		103: second.RegistrationRoot,
	}, []blockchain.AddressRegistered{first, second})
	fixture.tick()

	assert.Eventually(t, func() bool {
		return fixture.repository.state(watcherAddressA) == rootstock.PegInWatchUnsupportedEncoding &&
			fixture.repository.state(watcherAddressB) == rootstock.PegInWatchUnsupportedEncoding
	}, time.Second, 10*time.Millisecond)
	assert.Equal(t, 1, countCallsSince(fixture.registry.Calls, firstCall, "GetRegistrationRoot",
		func(arguments mock.Arguments) bool { return arguments.Get(1) == uint64(101) }))
	assert.Zero(t, countCallsSince(fixture.registry.Calls, firstCall, "GetRegistrationRoot",
		func(arguments mock.Arguments) bool { return arguments.Get(1) == uint64(102) }))
	assert.Equal(t, 1, countCallsSince(fixture.registry.Calls, firstCall, "GetAddressRegisteredEvents",
		func(arguments mock.Arguments) bool { return arguments.Get(1) == uint64(102) }))
}

func TestPegInWatcher_PublishesMismatchAndResyncOnlyForSuccessfulRootMismatch(t *testing.T) {
	fixture := newPegInWatcherFixture(t, true)
	event := watcherEvent(t, 100, 1, watcherAddressA, [32]byte{})
	fixture.chain.set(100, map[uint64][32]byte{100: event.RegistrationRoot}, []blockchain.AddressRegistered{event})
	fixture.registry.EXPECT().
		GetPegInAddress(watcherAddressA).
		Return(blockchain.PegInAddress{Encoding: blockchain.PegInAddressRegistryEncodingBech32}, nil).
		Once()
	fixture.eventBus.On("Publish", mock.Anything).Return().Twice()

	require.NoError(t, fixture.target.Prepare(context.Background()))

	events := publishedEvents(t, fixture.eventBus)
	require.Len(t, events, 2)
	mismatch, ok := events[0].(blockchain.PegInAddressRegistryRootMismatchEvent)
	require.True(t, ok)
	assert.Equal(t, uint64(100), mismatch.BlockNumber)
	assert.Equal(t, [32]byte{}, mismatch.LocalRoot)
	assert.Equal(t, event.RegistrationRoot, mismatch.ChainRoot)
	resync, ok := events[1].(blockchain.PegInAddressRegistryResyncStartedEvent)
	require.True(t, ok)
	assert.Equal(t, string(blockchain.PegInAddressRegistryRecoveryRootMismatch), resync.Reason)
}

func TestPegInWatcher_PublishesOnlyResyncForSuccessfulCatchUp(t *testing.T) {
	fixture := newPegInWatcherFixture(t, true)
	fixture.chain.set(100, map[uint64][32]byte{100: {}}, nil)
	require.NoError(t, fixture.target.Prepare(context.Background()))

	event := watcherEvent(t, 101, 1, watcherAddressA, [32]byte{})
	fixture.chain.set(101, map[uint64][32]byte{
		100: {},
		101: event.RegistrationRoot,
	}, []blockchain.AddressRegistered{event})
	fixture.registry.EXPECT().
		GetPegInAddress(watcherAddressA).
		Return(blockchain.PegInAddress{Encoding: blockchain.PegInAddressRegistryEncodingBech32}, nil).
		Once()
	fixture.eventBus.On("Publish", mock.Anything).Return().Once()
	fixture.start(t)
	defer fixture.stop(t)

	fixture.tick()

	assert.Eventually(t, func() bool {
		return len(publishedEvents(t, fixture.eventBus)) == 1
	}, time.Second, 10*time.Millisecond)
	events := publishedEvents(t, fixture.eventBus)
	require.Len(t, events, 1)
	resync, ok := events[0].(blockchain.PegInAddressRegistryResyncStartedEvent)
	require.True(t, ok)
	assert.Equal(t, string(blockchain.PegInAddressRegistryRecoveryCatchUp), resync.Reason)
}

func TestPegInWatcher_PublishesNothingOnReplayError(t *testing.T) {
	fixture := newPegInWatcherFixture(t, true)
	fixture.chain.setHeightError(assert.AnError)

	err := fixture.target.Prepare(context.Background())

	require.ErrorIs(t, err, assert.AnError)
	fixture.eventBus.AssertNotCalled(t, "Publish", mock.Anything)
}

func TestPegInWatcher_SerializesOverlappingScans(t *testing.T) {
	fixture := newPegInWatcherFixture(t, false)
	fixture.setChain(100)
	firstRootReadStarted, releaseFirstRootRead := fixture.chain.blockFirstRootRead()

	firstResult := make(chan error, 1)
	go func() {
		firstResult <- fixture.target.Prepare(context.Background())
	}()
	select {
	case <-firstRootReadStarted:
	case <-time.After(time.Second):
		require.FailNow(t, "first replay did not reach the root read")
	}

	secondResult := make(chan error, 1)
	go func() {
		secondResult <- fixture.target.Prepare(context.Background())
	}()
	assert.Never(t, func() bool {
		return fixture.chain.rootReads.Load() > 1
	}, 50*time.Millisecond, time.Millisecond, "a second replay entered while the first was active")

	close(releaseFirstRootRead)
	require.NoError(t, <-firstResult)
	require.NoError(t, <-secondResult)
	assert.Equal(t, int64(2), fixture.chain.rootReads.Load())
	assert.Equal(t, int64(1), fixture.chain.maxRootReads.Load())
}

// A node that delivers the same log twice, or delivers it a poll later than it should have, must
// not produce a second entry, a second import, or a state regression.
func TestPegInWatcher_AbsorbsDuplicateAndLateEventDelivery(t *testing.T) {
	fixture := newPegInWatcherFixtureWithPaging(t, false, 3)
	fixture.setChain(106)
	require.NoError(t, fixture.target.Prepare(context.Background()))
	fixture.chain.clearRequestedRanges()

	late := watcherEvent(t, 102, 1, watcherAddressA, [32]byte{})
	repeated := watcherEvent(t, 105, 2, watcherAddressB, late.RegistrationRoot)
	fixture.expectSupportedImport(t, late, 0)
	fixture.expectSupportedImport(t, repeated, 1)
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().
		RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{}, nil).
		Once()
	fixture.setChain(110, late, late, repeated)

	require.NoError(t, fixture.target.Prepare(context.Background()))
	require.NoError(t, fixture.target.Prepare(context.Background()))

	rows := fixture.repository.rowsSnapshot(t)
	require.Len(t, rows, 2)
	assert.Equal(t, late.TxHash, rows[0].TxHash)
	assert.Equal(t, repeated.TxHash, rows[1].TxHash)
	assert.Equal(t, rootstock.PegInWatchImported, rows[0].State)
	assert.Equal(t, rootstock.PegInWatchImported, rows[1].State)
	assert.Equal(t, [][2]uint64{{102, 104}, {105, 107}, {108, 110}}, fixture.chain.requestedRanges())
}

func TestPegInWatcher_ImportsThroughCapturedHead(t *testing.T) {
	fixture := newPegInWatcherFixtureWithPaging(t, false, 2)
	atHead := watcherEvent(t, 103, 0, watcherAddressA, [32]byte{})
	afterHead := watcherEvent(t, 104, 0, watcherAddressB, atHead.RegistrationRoot)
	fixture.setChain(103, atHead)
	fixture.chain.setHeightHook(func(uint64) {
		fixture.setChain(104, atHead, afterHead)
	})
	fixture.expectSupportedImport(t, atHead, 0)
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().
		RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{}, nil).
		Once()

	require.NoError(t, fixture.target.Prepare(context.Background()))

	rows := fixture.repository.rowsSnapshot(t)
	require.Len(t, rows, 1)
	assert.Equal(t, atHead.RskAddress, rows[0].RskAddress)
	assert.Equal(t, [][2]uint64{{103, 103}}, fixture.chain.requestedRanges())
}

func TestPegInWatcher_TimerAndReorgUseSameReplayPath(t *testing.T) {
	fixture := newPegInWatcherFixtureWithPaging(t, true, 2)
	event := watcherEvent(t, 103, 0, watcherAddressA, [32]byte{})
	fixture.setChain(103, event)
	fixture.registry.EXPECT().
		GetPegInAddress(watcherAddressA).
		Return(blockchain.PegInAddress{Encoding: blockchain.PegInAddressRegistryEncodingBech32}, nil).
		Times(3)
	fixture.eventBus.On("Publish", mock.Anything).Return().Times(6)

	require.NoError(t, fixture.target.Prepare(context.Background()))
	fixture.repository.clear()
	fixture.chain.clearRequestedRanges()
	fixture.start(t)
	defer fixture.stop(t)

	fixture.tick()
	assert.Eventually(t, func() bool {
		return fixture.repository.state(watcherAddressA) == rootstock.PegInWatchUnsupportedEncoding
	}, time.Second, time.Millisecond)
	timerRanges := fixture.chain.requestedRanges()
	require.NotEmpty(t, timerRanges)

	fixture.repository.clear()
	fixture.chain.clearRequestedRanges()
	fixture.reorgs <- blockchain.NodeReorgCheckEvent{
		BaseEvent:    entities.NewBaseEvent(blockchain.NodeReorgCheckEventId),
		NodeType:     entities.NodeTypeRootstock,
		CurrentDepth: 1,
	}
	assert.Eventually(t, func() bool {
		return fixture.repository.state(watcherAddressA) == rootstock.PegInWatchUnsupportedEncoding
	}, time.Second, time.Millisecond)
	assert.Equal(t, timerRanges, fixture.chain.requestedRanges())
}

func TestPegInWatcher_IgnoresInvalidReorgSignalsAndClosedSubscription(t *testing.T) {
	fixture := newPegInWatcherFixture(t, true)
	fixture.start(t)
	defer fixture.stop(t)
	assert.Eventually(t, func() bool {
		return countCallsSince(fixture.ticker.Calls, 0, "C", nil) > 0
	}, time.Second, time.Millisecond)

	fixture.reorgs <- entities.NewBaseEvent(blockchain.NodeReorgCheckEventId)
	fixture.reorgs <- blockchain.NodeReorgCheckEvent{
		BaseEvent:    entities.NewBaseEvent(blockchain.NodeReorgCheckEventId),
		NodeType:     entities.NodeTypeBitcoin,
		CurrentDepth: 1,
	}
	fixture.reorgs <- blockchain.NodeReorgCheckEvent{
		BaseEvent:    entities.NewBaseEvent(blockchain.NodeReorgCheckEventId),
		NodeType:     entities.NodeTypeRootstock,
		CurrentDepth: 0,
	}
	close(fixture.reorgs)

	assert.Never(t, func() bool {
		return countCallsSince(fixture.rskRpc.Calls, 0, "GetHeight", nil) > 0
	}, 50*time.Millisecond, time.Millisecond)
}

func TestPegInWatcher_RescanFailureLeavesEntryPending(t *testing.T) {
	fixture := newPegInWatcherFixture(t, false)
	event := watcherEvent(t, 100, 0, watcherAddressA, [32]byte{})
	fixture.setChain(100, event)
	fixture.expectSupportedImport(t, event, 0)
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().
		RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{}, assert.AnError).
		Once()

	require.NoError(t, fixture.target.Prepare(context.Background()))

	rows := fixture.repository.rowsSnapshot(t)
	require.Len(t, rows, 1)
	assert.Equal(t, rootstock.PegInWatchDiscovered, rows[0].State)
	assert.Contains(t, rows[0].LastError, "rescan PegIn addresses")
}

func TestFirstBootBackfillsRegistrationsAndContinuesIncrementally(t *testing.T) {
	fixture := newPegInWatcherFixtureWithPaging(t, false, 2)
	first := watcherEvent(t, 100, 0, watcherAddressA, [32]byte{})
	second := watcherEvent(t, 101, 0, watcherAddressB, first.RegistrationRoot)
	fixture.setChain(103, first, second)
	fixture.expectSupportedImport(t, first, 0)
	fixture.expectSupportedImport(t, second, 1)
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().
		RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{}, nil).
		Once()

	require.NoError(t, fixture.target.Prepare(context.Background()))

	rows := fixture.repository.rowsSnapshot(t)
	require.Len(t, rows, 2)
	assert.Equal(t, rootstock.PegInWatchImported, rows[0].State)
	assert.Equal(t, rootstock.PegInWatchImported, rows[1].State)
	require.NotEmpty(t, fixture.chain.requestedRanges())
	assert.Equal(t, uint64(100), fixture.chain.requestedRanges()[0][0])

	fixture.chain.clearRequestedRanges()
	third := watcherEvent(t, 103, 1, "0x00000000000000000000000000000000000000c3", second.RegistrationRoot)
	fixture.setChain(104, first, second, third)
	fixture.expectSupportedImport(t, third, 2)
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().
		RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{}, nil).
		Once()

	fixture.start(t)
	defer fixture.stop(t)
	fixture.tick()
	assert.Eventually(t, func() bool {
		return fixture.repository.state(third.RskAddress) == rootstock.PegInWatchImported
	}, time.Second, time.Millisecond)
	require.NotEmpty(t, fixture.chain.requestedRanges())
	assert.Equal(t, uint64(104), fixture.chain.requestedRanges()[0][1])
}

func TestFirstBootSkipsIdleDeploymentRange(t *testing.T) {
	fixture := newPegInWatcherFixtureWithPaging(t, false, 2)
	first := watcherEvent(t, 102, 0, watcherAddressA, [32]byte{})
	fixture.setChain(103, first)
	fixture.expectSupportedImport(t, first, 0)
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().
		RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{}, nil).
		Once()

	require.NoError(t, fixture.target.Prepare(context.Background()))

	require.Len(t, fixture.repository.rowsSnapshot(t), 1)
	assert.Equal(t, [][2]uint64{{102, 103}}, fixture.chain.requestedRanges())
}

// A restart has no durable watcher checkpoint. It must derive the matching boundary from rows,
// replay the missing suffix, and finish pending imports through the public Prepare API.
func TestRestartConvergesFromRowsWithNilCheckpoint(t *testing.T) {
	tests := []struct {
		name       string
		firstState rootstock.PegInWatchState
	}{
		{name: "after the first entry was persisted", firstState: rootstock.PegInWatchDiscovered},
		{name: "after the first entry was imported", firstState: rootstock.PegInWatchImported},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newPegInWatcherFixtureWithPaging(t, false, 2)
			first := watcherEvent(t, 100, 0, watcherAddressA, [32]byte{})
			second := watcherEvent(t, 101, 0, watcherAddressB, first.RegistrationRoot)
			fixture.setChain(103, first, second)
			firstRow := rootstock.NewPegInWatch(
				first.TxHash,
				first.LogIndex,
				first.BlockNumber,
				first.RskAddress,
				"",
				first.RegistrationRoot,
			)
			firstRow.State = testCase.firstState
			if testCase.firstState == rootstock.PegInWatchImported {
				firstRow.BtcAddress = datasets.Base58Addresses[0].Address
			} else {
				fixture.expectSupportedImport(t, first, 0)
			}
			fixture.repository.seed(firstRow)
			fixture.expectSupportedImport(t, second, 1)
			fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
			fixture.wallet.EXPECT().
				RescanBlockchain(int64(100)).
				Return(blockchain.BitcoinRescanResult{}, nil).
				Once()
			fixture.resetTarget(t, nil)

			require.NoError(t, fixture.target.Prepare(context.Background()))

			rows := fixture.repository.rowsSnapshot(t)
			require.Len(t, rows, 2)
			assert.Equal(t, rootstock.PegInWatchImported, rows[0].State)
			assert.Equal(t, rootstock.PegInWatchImported, rows[1].State)
			require.NotEmpty(t, fixture.chain.requestedRanges())
			assert.Equal(t, uint64(101), fixture.chain.requestedRanges()[0][0],
				"the persisted row must establish the matching replay boundary")
		})
	}
}

func TestMissedEventSignalsMismatchAndRepairsCompleteWatchSet(t *testing.T) {
	fixture := newPegInWatcherFixtureWithPaging(t, true, 2)
	fixture.eventBus.On("Publish", mock.Anything).Return().Maybe()
	first := watcherEvent(t, 100, 0, watcherAddressA, [32]byte{})
	fixture.setChain(101, first)
	fixture.expectSupportedImport(t, first, 0)
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().
		RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{}, nil).
		Once()
	require.NoError(t, fixture.target.Prepare(context.Background()))

	skipped := watcherEvent(t, 101, 0, watcherAddressB, first.RegistrationRoot)
	last := watcherEvent(t, 102, 0, "0x00000000000000000000000000000000000000c3", skipped.RegistrationRoot)
	fixture.setChain(103, first, skipped, last)
	fixture.chain.omitNext(skipped)
	publishedBeforeFailure := len(publishedEvents(t, fixture.eventBus))

	require.Error(t, fixture.target.Prepare(context.Background()))
	assert.Len(t, publishedEvents(t, fixture.eventBus), publishedBeforeFailure,
		"Replay must not publish when rebuild validation fails")

	fixture.expectSupportedImport(t, skipped, 1)
	fixture.expectSupportedImport(t, last, 2)
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().
		RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{}, nil).
		Once()
	require.NoError(t, fixture.target.Prepare(context.Background()))

	rows := fixture.repository.rowsSnapshot(t)
	require.Len(t, rows, 3)
	for _, row := range rows {
		assert.Equal(t, rootstock.PegInWatchImported, row.State)
	}
	assert.Len(t, publishedEvents(t, fixture.eventBus), publishedBeforeFailure+2,
		"the watcher must publish mismatch and resync only after a successful rebuild")
}

func TestSilentEventStreamHealthCheckSignalsMismatchAndReplays(t *testing.T) {
	fixture := newPegInWatcherFixtureWithPaging(t, true, 2)
	fixture.eventBus.On("Publish", mock.Anything).Return().Maybe()
	first := watcherEvent(t, 100, 0, watcherAddressA, [32]byte{})
	fixture.setChain(101, first)
	fixture.expectSupportedImport(t, first, 0)
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().
		RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{}, nil).
		Once()
	require.NoError(t, fixture.target.Prepare(context.Background()))
	publishedBeforeHealthCheck := len(publishedEvents(t, fixture.eventBus))

	silentlyOmitted := watcherEvent(t, 100, 1, watcherAddressB, first.RegistrationRoot)
	fixture.setChain(101, first, silentlyOmitted)
	fixture.registry.EXPECT().
		GetPegInAddress(first.RskAddress).
		Return(blockchain.PegInAddress{
			Payload: func() []byte {
				payload, _ := watcherDepositPayload()
				return payload
			}(),
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.wallet.EXPECT().ImportAddress(datasets.Base58Addresses[0].Address).Return(nil).Once()
	fixture.expectSupportedImport(t, silentlyOmitted, 1)
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().
		RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{}, nil).
		Once()

	require.NoError(t, fixture.target.Prepare(context.Background()))

	rows := fixture.repository.rowsSnapshot(t)
	require.Len(t, rows, 2)
	assert.Equal(t, rootstock.PegInWatchImported, rows[0].State)
	assert.Equal(t, rootstock.PegInWatchImported, rows[1].State)
	assert.Len(t, publishedEvents(t, fixture.eventBus), publishedBeforeHealthCheck+2)
}
