package watcher

import (
	"cmp"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/crypto"
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

func newRegistryWatcher(
	repository rootstock.PegInWatchRepository,
	checkpoints rootstock.PegInWatchCheckpointRepository,
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	btcNetwork blockchain.BitcoinNetwork,
	wallet blockchain.BitcoinWallet,
	eventBus entities.EventBus,
	ticker utils.Ticker,
	startBlock uint64,
	pageSize uint64,
) *PegInWatcher {
	replay := w.NewReplayRegisteredAddressesUseCase(
		repository,
		checkpoints,
		registry,
		rskRpc,
		eventBus,
		wallet,
		crypto.Keccak256,
	)
	finalize := w.NewFinalizeRegisteredAddressImportUseCase(repository)
	return NewPegInWatcher(
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

func uint64Pointer(expected uint64) interface{} {
	return mock.MatchedBy(func(actual *uint64) bool {
		return actual != nil && *actual == expected
	})
}

type overlappingReplayRepository struct {
	mutex      sync.Mutex
	checkpoint rootstock.PegInWatchCheckpoint
	found      bool
	setCalls   int
	pruneCalls int
}

func (repository *overlappingReplayRepository) Upsert(
	context.Context,
	rootstock.PegInWatch,
) error {
	return nil
}

func (repository *overlappingReplayRepository) Get(
	context.Context,
	string,
) (*rootstock.PegInWatch, error) {
	return nil, nil
}

func (repository *overlappingReplayRepository) List(
	context.Context,
) ([]rootstock.PegInWatch, error) {
	return nil, nil
}

func (repository *overlappingReplayRepository) Update(
	context.Context,
	rootstock.PegInWatch,
) error {
	return nil
}

func (repository *overlappingReplayRepository) GetCheckpoint(
	context.Context,
) (*rootstock.PegInWatchCheckpoint, error) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()
	if !repository.found {
		return nil, nil
	}
	checkpoint := repository.checkpoint
	return &checkpoint, nil
}

func (repository *overlappingReplayRepository) SetCheckpoint(
	_ context.Context,
	checkpoint rootstock.PegInWatchCheckpoint,
) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()
	repository.setCalls++
	repository.checkpoint = checkpoint
	repository.found = true
	return nil
}

func (repository *overlappingReplayRepository) DeleteFromBlock(context.Context, uint64) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()
	repository.pruneCalls++
	return nil
}

func (repository *overlappingReplayRepository) checkpointState() (
	rootstock.PegInWatchCheckpoint,
	bool,
	int,
	int,
) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()
	return repository.checkpoint, repository.found, repository.setCalls, repository.pruneCalls
}

func newShortHeadCheckpointWatcher(t *testing.T) (
	*overlappingReplayRepository,
	*mocks.PegInAddressRegistryContractMock,
	*PegInWatcher,
	rootstock.PegInWatchCheckpoint,
) {
	t.Helper()
	original := rootstock.PegInWatchCheckpoint{
		LocalRoot:          [32]byte{1},
		LastProcessedBlock: 105,
	}
	repository := &overlappingReplayRepository{checkpoint: original, found: true}
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(99), nil).Once()
	watcher := newRegistryWatcher(
		repository,
		repository,
		registry,
		rskRpc,
		nil,
		mocks.NewBitcoinWalletMock(t),
		nil,
		nil,
		100,
		10,
	)
	return repository, registry, watcher, original
}

func TestPegInWatcher_ScanKeepsCheckpointWhenHeadIsBelowStart(t *testing.T) {
	repository, registry, watcher, original := newShortHeadCheckpointWatcher(t)

	require.NoError(t, watcher.scan(context.Background()))

	checkpoint, found, setCalls, pruneCalls := repository.checkpointState()
	require.True(t, found)
	assert.Equal(t, original, checkpoint)
	assert.Zero(t, setCalls)
	assert.Zero(t, pruneCalls)
	registry.AssertNotCalled(t, "GetAddressRegisteredEvents", mock.Anything, mock.Anything, mock.Anything)
	registry.AssertNotCalled(t, "GetRegistrationRoot", mock.Anything, mock.Anything)
}

func TestPegInWatcher_SerializesOverlappingScans(t *testing.T) {
	repository := &overlappingReplayRepository{}
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)

	firstRootReadStarted := make(chan struct{})
	releaseFirstRootRead := make(chan struct{})
	var rootReads atomic.Int64
	var activeRootReads atomic.Int64
	var maximumActiveRootReads atomic.Int64
	registry.EXPECT().
		GetRegistrationRoot(mock.Anything, uint64(0)).
		RunAndReturn(func(context.Context, uint64) ([32]byte, error) {
			call := rootReads.Add(1)
			active := activeRootReads.Add(1)
			for {
				maximum := maximumActiveRootReads.Load()
				if active <= maximum || maximumActiveRootReads.CompareAndSwap(maximum, active) {
					break
				}
			}
			defer activeRootReads.Add(-1)
			if call == 1 {
				close(firstRootReadStarted)
				<-releaseFirstRootRead
			}
			return [32]byte{}, nil
		}).
		Twice()
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(0), nil).Twice()

	watcher := newOverlappingReplayWatcher(repository, registry, rskRpc, wallet)

	firstResult := make(chan error, 1)
	go func() {
		firstResult <- watcher.scan(context.Background())
	}()
	select {
	case <-firstRootReadStarted:
	case <-time.After(time.Second):
		require.FailNow(t, "first replay did not reach the root read")
	}

	secondStarted := make(chan struct{})
	secondResult := make(chan error, 1)
	go func() {
		close(secondStarted)
		secondResult <- watcher.scan(context.Background())
	}()
	<-secondStarted
	assert.Never(t, func() bool {
		return rootReads.Load() > 1
	}, 50*time.Millisecond, time.Millisecond, "a second replay entered while the first was active")

	close(releaseFirstRootRead)
	require.NoError(t, <-firstResult)
	require.NoError(t, <-secondResult)
	assert.Equal(t, int64(2), rootReads.Load())
	assert.Equal(t, int64(1), maximumActiveRootReads.Load())
}

func newOverlappingReplayWatcher(
	repository *overlappingReplayRepository,
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	wallet blockchain.BitcoinWallet,
) *PegInWatcher {
	return newRegistryWatcher(repository, repository, registry, rskRpc, nil, wallet, nil, nil, 0, 1)
}

func TestPegInWatcher_LogsRootReadFailureWithoutPublishingCheckpoint(t *testing.T) {
	rootErr := errors.New("RSK registry root read failed")
	original := rootstock.PegInWatchCheckpoint{
		LocalRoot:          [32]byte{1},
		LastProcessedBlock: 100,
	}
	repository := &overlappingReplayRepository{checkpoint: original, found: true}
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	registry.EXPECT().
		GetRegistrationRoot(mock.Anything, uint64(101)).
		Return([32]byte{}, rootErr).
		Once()
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(101), nil).Once()
	watcher := newRegistryWatcher(
		repository,
		repository,
		registry,
		rskRpc,
		nil,
		mocks.NewBitcoinWalletMock(t),
		nil,
		nil,
		0,
		10,
	)
	require.NoError(t, watcher.Prepare(context.Background()))
	capturedLogs := test.CaptureStructuredLogs(t)

	watcher.scanAndLog()

	logEntries := capturedLogs()
	require.Len(t, logEntries, 1)
	assert.Equal(t, "error", logEntries[0].Level())
	assert.Equal(
		t,
		"PegIn address registry watcher scan failed: ReplayRegisteredAddresses: get PegIn address registry root at block 101: RSK registry root read failed",
		logEntries[0].Message(),
	)
	checkpoint, found, setCalls, pruneCalls := repository.checkpointState()
	require.True(t, found)
	assert.Equal(t, original, checkpoint)
	assert.Zero(t, setCalls)
	assert.Zero(t, pruneCalls)
}

// A write that fails and a write lost to a dying process leave the same state behind, so a scenario
// injects this to stop a watcher at a boundary it cannot otherwise be stopped at.
var errProcessStopped = errors.New("process stopped")

// registryWatchStore stands in for the Mongo watch repository, keeping the behaviour a restarted
// watcher depends on: Upsert never overwrites an existing (tx_hash, log_index), Update needs the
// entry to exist, List answers in (block, log index) order, and the checkpoint reads as absent until
// it is first written.
type registryWatchStore struct {
	mutex         sync.Mutex
	entries       []rootstock.PegInWatch
	checkpoint    rootstock.PegInWatchCheckpoint
	hasCheckpoint bool
	advances      []uint64

	// Each of these fails the next call to the method it names and is then cleared.
	failGet           error
	failUpdate        error
	failSetCheckpoint error
}

func (store *registryWatchStore) Upsert(_ context.Context, entry rootstock.PegInWatch) error {
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
) (*rootstock.PegInWatch, error) {
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

func (store *registryWatchStore) List(context.Context) ([]rootstock.PegInWatch, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	return store.sortedEntries(), nil
}

func (store *registryWatchStore) Update(_ context.Context, entry rootstock.PegInWatch) error {
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
) (*rootstock.PegInWatchCheckpoint, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if !store.hasCheckpoint {
		return nil, nil
	}
	checkpoint := store.checkpoint
	return &checkpoint, nil
}

func (store *registryWatchStore) SetCheckpoint(
	_ context.Context,
	checkpoint rootstock.PegInWatchCheckpoint,
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

func (store *registryWatchStore) DeleteFromBlock(_ context.Context, fromBlock uint64) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.entries = pruneWatchEntriesFromBlock(store.entries, fromBlock)
	return nil
}

func pruneWatchEntriesFromBlock(
	entries []rootstock.PegInWatch,
	fromBlock uint64,
) []rootstock.PegInWatch {
	return slices.DeleteFunc(entries, func(entry rootstock.PegInWatch) bool {
		return entry.BlockNumber >= fromBlock
	})
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

func (store *registryWatchStore) sortedEntries() []rootstock.PegInWatch {
	entries := slices.Clone(store.entries)
	slices.SortFunc(entries, func(first, second rootstock.PegInWatch) int {
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
	mutex     sync.Mutex
	head      uint64
	logs      []blockchain.AddressRegistered
	requested [][2]uint64
	omitOnce  map[string]int
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

func (chain *rskChain) logsIn(fromBlock, toBlock uint64) []blockchain.AddressRegistered {
	chain.mutex.Lock()
	defer chain.mutex.Unlock()
	chain.requested = append(chain.requested, [2]uint64{fromBlock, toBlock})
	events := make([]blockchain.AddressRegistered, 0)
	for _, log := range chain.logs {
		identity := eventIdentity(log)
		inRequestedRange := log.BlockNumber >= fromBlock && log.BlockNumber <= toBlock
		if inRequestedRange {
			if chain.omitOnce[identity] > 0 {
				chain.omitOnce[identity]--
				continue
			}
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

	mutex     sync.Mutex
	deposits  map[string]depositAddress
	imports   []string
	rootReads []uint64
	rescanErr error

	startBlock uint64
	pageSize   uint64
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

func (bus *registryReorgEventBus) clearPublished() {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()
	bus.published = nil
}

const registryRestartStartBlock = uint64(100)

func newRegistryRestartFixture(t *testing.T, pageSize uint64) *registryRestartFixture {
	t.Helper()
	fixture := &registryRestartFixture{
		t:          t,
		store:      &registryWatchStore{},
		chain:      &rskChain{},
		registry:   mocks.NewPegInAddressRegistryContractMock(t),
		rskRpc:     mocks.NewRootstockRpcServerMock(t),
		btcNetwork: &mocks.BtcRpcMock{},
		wallet:     mocks.NewBitcoinWalletMock(t),
		eventBus:   newRegistryReorgEventBus(),
		deposits:   make(map[string]depositAddress),
		startBlock: registryRestartStartBlock,
		pageSize:   pageSize,
	}

	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).RunAndReturn(func(context.Context) (uint64, error) {
		return fixture.chain.currentHead(), nil
	})
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, fromBlock uint64, toBlock *uint64) ([]blockchain.AddressRegistered, error) {
			require.NotNil(t, toBlock, "the scanner must always bound its range at the captured head")
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
		RunAndReturn(func(int64) (blockchain.BitcoinRescanResult, error) {
			return blockchain.BitcoinRescanResult{StartHeight: 100, StopHeight: 200}, fixture.rescanErr
		}).
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

func (fixture *registryRestartFixture) watchSet() []rootstock.PegInWatch {
	fixture.t.Helper()
	entries, err := fixture.store.List(context.Background())
	require.NoError(fixture.t, err)
	return entries
}

func (fixture *registryRestartFixture) checkpointSnapshot() (
	rootstock.PegInWatchCheckpoint,
	bool,
) {
	fixture.store.mutex.Lock()
	defer fixture.store.mutex.Unlock()
	return fixture.store.checkpoint, fixture.store.hasCheckpoint
}

func (fixture *registryRestartFixture) checkpointAdvances() []uint64 {
	return fixture.store.checkpointAdvances()
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
		fixture.store,
		fixture.registry,
		fixture.rskRpc,
		fixture.btcNetwork,
		fixture.wallet,
		fixture.eventBus,
		ticker,
		fixture.startBlock,
		fixture.pageSize,
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
func TestPegInWatcher_ConvergesAfterRestartAtEveryBoundary(t *testing.T) {
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
			fixture := newRegistryRestartFixture(t, 3)
			event := fixture.chainRegisters("registration", registrationBlock, 4)
			fixture.chain.setHead(104)

			testCase.stop(fixture.store)
			fixture.scanOnce()
			fixture.scanOnce()

			requireImportedEntry(t, fixture, event)
			assert.Len(t, fixture.importedAddresses(), testCase.imports)
			assert.Equal(t, []uint64{104}, fixture.checkpointAdvances())
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
	assert.Equal(t, rootstock.PegInWatchImported, entries[0].State)
	assert.Equal(t, fixture.depositAddressOf(event), entries[0].BtcAddress)
	assert.Empty(t, entries[0].LastError)
}

// A node that delivers the same log twice, or delivers it a poll later than it should have, must
// not produce a second entry, a second import, a state regression, or a checkpoint that moves
// anywhere but forward.
func TestPegInWatcher_AbsorbsDuplicateAndLateEventDelivery(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 3)
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
		{102, 104}, {105, 107}, {108, 110},
	}, fixture.chain.requestedRanges())
	entries := fixture.watchSet()
	require.Len(t, entries, 2)
	assert.Equal(t, late.TxHash, entries[0].TxHash)
	assert.Equal(t, repeated.TxHash, entries[1].TxHash)
	for _, entry := range entries {
		assert.Equal(t, rootstock.PegInWatchImported, entry.State)
	}
	assert.Equal(
		t,
		[]string{fixture.depositAddressOf(late), fixture.depositAddressOf(repeated)},
		fixture.importedAddresses(),
	)
	assert.Equal(t, fixture.chain.rootAt(110), fixture.store.checkpoint.LocalRoot)
	assert.Equal(t, []uint64{106, 110}, fixture.checkpointAdvances())
}

func TestPegInWatcher_ImportsThroughCapturedHead(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 2)
	atHead := fixture.chainRegisters("captured-head", 103, 0)
	fixture.chain.setHead(103)

	fixture.scanOnce()

	requireImportedEntry(t, fixture, atHead)
	requireMatchingCheckpoint(t, fixture, atHead.RegistrationRoot, 103)
	assert.Equal(t, [][2]uint64{{103, 103}}, fixture.chain.requestedRanges())
}

func TestPegInWatcher_TimerAndReorgUseSameReplayPath(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 2)
	fixture.chainRegisters("same-path", 103, 0)
	fixture.chain.setHead(103)
	session := fixture.startScanner()
	defer session.stop()

	session.poll()
	firstRanges := fixture.chain.requestedRanges()
	require.Equal(t, [][2]uint64{{103, 103}}, firstRanges)

	fixture.store.mutex.Lock()
	fixture.store.entries = nil
	fixture.store.checkpoint = rootstock.PegInWatchCheckpoint{}
	fixture.store.hasCheckpoint = false
	fixture.store.mutex.Unlock()
	fixture.eventBus.Publish(blockchain.NodeReorgCheckEvent{
		BaseEvent:    entities.NewBaseEvent(blockchain.NodeReorgCheckEventId),
		NodeType:     entities.NodeTypeRootstock,
		CurrentDepth: 1,
	})

	require.Eventually(t, func() bool {
		return len(fixture.chain.requestedRanges()) == 2*len(firstRanges)
	}, time.Second, time.Millisecond)
	assert.Equal(t, append(firstRanges, firstRanges...), fixture.chain.requestedRanges())
}

func TestPegInWatcher_IgnoresInvalidReorgSignalsAndClosedSubscription(t *testing.T) {
	eventBus := newRegistryReorgEventBus()
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	ticker := newSessionTicker()
	repository := &overlappingReplayRepository{}
	watcher := newRegistryWatcher(
		repository,
		repository,
		mocks.NewPegInAddressRegistryContractMock(t),
		rskRpc,
		nil,
		mocks.NewBitcoinWalletMock(t),
		eventBus,
		ticker,
		100,
		2,
	)
	session := startWatcherSession(t, watcher, ticker)
	defer session.stop()

	eventBus.Publish(entities.NewBaseEvent(blockchain.NodeReorgCheckEventId))
	eventBus.Publish(blockchain.NodeReorgCheckEvent{
		BaseEvent:    entities.NewBaseEvent(blockchain.NodeReorgCheckEventId),
		NodeType:     entities.NodeTypeBitcoin,
		CurrentDepth: 1,
	})
	eventBus.Publish(blockchain.NodeReorgCheckEvent{
		BaseEvent:    entities.NewBaseEvent(blockchain.NodeReorgCheckEventId),
		NodeType:     entities.NodeTypeRootstock,
		CurrentDepth: 0,
	})
	selectsBeforeClose := ticker.selects.Load()
	close(eventBus.events)

	require.Eventually(t, func() bool {
		return ticker.selects.Load() > selectsBeforeClose
	}, time.Second, time.Millisecond)
	rskRpc.AssertNotCalled(t, "GetHeight", mock.Anything)
}

func TestPegInWatcher_RescanFailureLeavesEntryPending(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 2)
	fixture.rescanErr = assert.AnError
	event := fixture.chainRegisters("rescan-failure", 100, 0)
	fixture.chain.setHead(100)

	fixture.scanOnce()

	entries := fixture.watchSet()
	require.Len(t, entries, 1)
	assert.Equal(t, event.RskAddress, entries[0].RskAddress)
	assert.Equal(t, rootstock.PegInWatchDiscovered, entries[0].State)
}

type checkpointRestartStore struct {
	mutex             sync.Mutex
	entries           []rootstock.PegInWatch
	checkpoint        rootstock.PegInWatchCheckpoint
	hasCheckpoint     bool
	failSetCheckpoint error
	operations        []string
}

func (store *checkpointRestartStore) Upsert(_ context.Context, entry rootstock.PegInWatch) error {
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
) (*rootstock.PegInWatch, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	index := store.entryIndex(rskAddress)
	if index < 0 {
		return nil, nil
	}
	entry := store.entries[index]
	return &entry, nil
}

func (store *checkpointRestartStore) List(context.Context) ([]rootstock.PegInWatch, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	return append([]rootstock.PegInWatch(nil), store.entries...), nil
}

func (store *checkpointRestartStore) Update(_ context.Context, entry rootstock.PegInWatch) error {
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
) (*rootstock.PegInWatchCheckpoint, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if !store.hasCheckpoint {
		return nil, nil
	}
	checkpoint := store.checkpoint
	return &checkpoint, nil
}

func (store *checkpointRestartStore) SetCheckpoint(
	_ context.Context,
	checkpoint rootstock.PegInWatchCheckpoint,
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

func (store *checkpointRestartStore) DeleteFromBlock(_ context.Context, fromBlock uint64) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.entries = pruneWatchEntriesFromBlock(store.entries, fromBlock)
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
	[]rootstock.PegInWatch,
	rootstock.PegInWatchCheckpoint,
	[]string,
) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	return append([]rootstock.PegInWatch(nil), store.entries...),
		store.checkpoint,
		append([]string(nil), store.operations...)
}

func TestPegInWatcher_RestartsAfterEntryWriteBeforeCheckpoint(t *testing.T) {
	scenario := newCheckpointRestartScenario(t)

	scenario.scanOnce(t)
	entries, checkpoint, _ := scenario.store.snapshot()
	require.Len(t, entries, 2)
	assert.Equal(t, rootstock.PegInWatchDiscovered, entries[1].State)
	assert.Equal(t, uint64(100), checkpoint.LastProcessedBlock)
	assert.Equal(t, scenario.previousRoot, checkpoint.LocalRoot)

	scenario.scanOnce(t)
	entries, checkpoint, operations := scenario.store.snapshot()
	require.Len(t, entries, 2)
	assert.Equal(t, rootstock.PegInWatchImported, entries[1].State)
	assert.Equal(t, rootstock.PegInWatchCheckpoint{
		LocalRoot:          scenario.expectedRoot,
		LastProcessedBlock: 103,
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

const (
	checkpointRestartStartBlock = uint64(100)
	checkpointRestartEventBlock = uint64(101)
)

func newCheckpointRestartScenario(t *testing.T) *checkpointRestartScenario {
	t.Helper()
	vectors := pinnedLBCRegistrations(t)
	prior := vectors[0]
	registered := vectors[1]
	event := blockchain.AddressRegistered{
		TxHash:           "registration",
		LogIndex:         1,
		BlockNumber:      checkpointRestartEventBlock,
		RskAddress:       registered.rskAddress,
		RegistrationRoot: registered.root,
	}
	deposit := knownDepositAddress(0)
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(103), nil).Twice()
	wallet := mocks.NewBitcoinWalletMock(t)
	wallet.EXPECT().ImportAddress(deposit.address).Return(nil).Twice()
	wallet.EXPECT().RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{StartHeight: 100, StopHeight: 200}, nil).
		Once()
	btcNetwork := &mocks.BtcRpcMock{}
	t.Cleanup(func() { btcNetwork.AssertExpectations(t) })
	btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()

	return &checkpointRestartScenario{
		store:        newCheckpointRestartStore(prior),
		registry:     newCheckpointRestartRegistry(t, event, prior.root, deposit),
		rskRpc:       rskRpc,
		btcNetwork:   btcNetwork,
		wallet:       wallet,
		previousRoot: prior.root,
		expectedRoot: registered.root,
	}
}

func newCheckpointRestartStore(prior pinnedLBCRegistration) *checkpointRestartStore {
	return &checkpointRestartStore{
		entries: []rootstock.PegInWatch{{
			TxHash:           "prior-registration",
			BlockNumber:      checkpointRestartStartBlock,
			RskAddress:       prior.rskAddress,
			RegistrationRoot: prior.root,
			State:            rootstock.PegInWatchImported,
		}},
		checkpoint: rootstock.PegInWatchCheckpoint{
			LocalRoot:          prior.root,
			LastProcessedBlock: checkpointRestartStartBlock,
		},
		hasCheckpoint:     true,
		failSetCheckpoint: errProcessStopped,
	}
}

func newCheckpointRestartRegistry(
	t *testing.T,
	event blockchain.AddressRegistered,
	previousRoot [32]byte,
	deposit depositAddress,
) *mocks.PegInAddressRegistryContractMock {
	t.Helper()
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(101), uint64Pointer(103)).
		Return([]blockchain.AddressRegistered{event}, nil).
		Once()
	registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(103)).
		Return(event.RegistrationRoot, nil).
		Twice()
	registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(100)).Return(previousRoot, nil).Once()
	registry.EXPECT().GetPegInAddress(event.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  deposit.payload,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Twice()
	return registry
}

func (scenario *checkpointRestartScenario) scanOnce(t *testing.T) {
	t.Helper()
	ticker := newSessionTicker()
	watcher := newRegistryWatcher(
		scenario.store,
		scenario.store,
		scenario.registry,
		scenario.rskRpc,
		scenario.btcNetwork,
		scenario.wallet,
		nil,
		ticker,
		100,
		3,
	)
	session := startWatcherSession(t, watcher, ticker)
	session.poll()
	session.stop()
}

func TestFirstBootBackfillsRegistrationsAndStoresMatchingCheckpoint(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 2)
	first := fixture.chainRegisters("first", 100, 0)
	second := fixture.chainRegisters("second", 101, 0)
	fixture.chain.setHead(103)

	session := fixture.startScanner()
	defer session.stop()
	session.poll()

	requireCompleteImportedWatchSet(t, fixture, first, second)
	requireMatchingCheckpoint(t, fixture, second.RegistrationRoot, 103)
	initialRanges := fixture.chain.requestedRanges()
	require.NotEmpty(t, initialRanges)
	assert.Equal(t, fixture.startBlock, initialRanges[0][0],
		"the first mismatch is at the deployment block in this fixture")
	assert.Equal(t, []uint64{103, 101, 100}, fixture.rootReadBlocks(),
		"first boot must check the head and find the first mismatching block")

	requestCountBeforeIncrementalPoll := len(fixture.chain.requestedRanges())
	third := fixture.chainRegisters("third", 103, 0)
	fixture.chain.setHead(104)
	session.poll()

	requireCompleteImportedWatchSet(t, fixture, first, second, third)
	requireMatchingCheckpoint(t, fixture, third.RegistrationRoot, 104)
	incrementalRanges := fixture.chain.requestedRanges()[requestCountBeforeIncrementalPoll:]
	require.NotEmpty(t, incrementalRanges)
	assert.Equal(t, uint64(103), incrementalRanges[0][0], "normal operation must continue after the saved checkpoint")
	assert.Equal(t, []uint64{103, 104}, fixture.checkpointAdvances())
	rootReads := fixture.rootReadBlocks()
	assert.Contains(t, rootReads, uint64(104))
}

func TestFirstBootSkipsIdleDeploymentRange(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 2)
	first := fixture.chainRegisters("first-after-idle", 102, 0)
	fixture.chain.setHead(103)

	fixture.scanOnce()

	requireCompleteImportedWatchSet(t, fixture, first)
	requireMatchingCheckpoint(t, fixture, first.RegistrationRoot, 103)
	assert.Equal(t, [][2]uint64{{102, 103}}, fixture.chain.requestedRanges())
	assert.Equal(t, []uint64{103, 101, 102}, fixture.rootReadBlocks())
}

func TestRestartResumesFromTrustedCheckpointWithMatchingRoot(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 2)
	first := fixture.chainRegisters("persisted-first", 100, 0)
	second := fixture.chainRegisters("persisted-second", 101, 0)
	fixture.chain.setHead(102)
	fixture.scanOnce()

	requireCompleteImportedWatchSet(t, fixture, first, second)
	requireMatchingCheckpoint(t, fixture, second.RegistrationRoot, 102)
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
	requireMatchingCheckpoint(t, fixture, third.RegistrationRoot, 103)
	assert.Equal(t, []uint64{102, 103}, fixture.checkpointAdvances())
	rootReads := fixture.rootReadBlocks()
	require.NotEmpty(t, rootReads)
	assert.Contains(t, rootReads, uint64(103), "restart must verify against the captured head")
}

func TestMissedEventSignalsMismatchAndRepairsCompleteWatchSet(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 2)
	capturedLogs := test.CaptureStructuredLogs(t)
	first := fixture.chainRegisters("first", 100, 0)
	fixture.chain.setHead(101)
	fixture.scanOnce()
	requireCompleteImportedWatchSet(t, fixture, first)

	skipped := fixture.chainRegisters("skipped", 101, 0)
	last := fixture.chainRegisters("last", 102, 0)
	fixture.chain.omitOnNextQuery(skipped)
	fixture.chain.setHead(103)
	fixture.scanOnce()
	fixture.eventBus.clearPublished()
	logOffset := len(capturedLogs())
	fixture.scanOnce()

	requireCompleteWatchSet(t, fixture, first, skipped, last)
	assert.Equal(t, []string{
		fixture.depositAddressOf(first),
		fixture.depositAddressOf(skipped),
		fixture.depositAddressOf(last),
	}, fixture.importedAddresses())
	requireMatchingCheckpoint(t, fixture, last.RegistrationRoot, 103)
	assert.Equal(t, 1, replayCountFromDeployment(fixture))
	assert.Len(t, fixture.importedAddresses(), 3)
	mismatchSignal := requireOneIntegritySignal(t, fixture, 103, last.RegistrationRoot)
	logEntry := requireSingleMismatchLog(t, capturedLogs()[logOffset:])
	assert.Equal(t, "error", logEntry.Level())
	assert.InDelta(t, float64(103), logEntry.Field(test.LogKey("block_number")), 0)
	assert.Equal(t, fmt.Sprintf("0x%x", mismatchSignal.LocalRoot), logEntry.Field(test.LogKey("local_root")))
	assert.Equal(t, fmt.Sprintf("0x%x", last.RegistrationRoot), logEntry.Field(test.LogKey("chain_root")))
	assert.NotEqual(t, logEntry.Field(test.LogKey("chain_root")), logEntry.Field(test.LogKey("local_root")))
	assert.Equal(t, "captured_head", logEntry.Field(test.LogKey("source")))
}

func TestSilentEventStreamHealthCheckSignalsMismatchAndReplays(t *testing.T) {
	fixture := newRegistryRestartFixture(t, 2)
	capturedLogs := test.CaptureStructuredLogs(t)
	first := fixture.chainRegisters("first", 100, 0)
	fixture.chain.setHead(101)

	session := fixture.startScanner()
	defer session.stop()
	session.poll()
	requireCompleteImportedWatchSet(t, fixture, first)
	requireMatchingCheckpoint(t, fixture, first.RegistrationRoot, 101)

	session.poll()
	requireCompleteImportedWatchSet(t, fixture, first)
	requireMatchingCheckpoint(t, fixture, first.RegistrationRoot, 101)
	assert.Equal(t, 1, replayCountFromDeployment(fixture), "an equal-root tick must not start a cold replay")
	assert.Equal(t, []uint64{101}, fixture.checkpointAdvances(), "an equal-root tick must not republish progress")
	fixture.eventBus.clearPublished()
	logOffset := len(capturedLogs())

	// The node's first reading of block 100 omitted this registration. The head and
	// trusted checkpoint cannot advance, so only the same-head root check can expose the gap.
	silentlyOmitted := fixture.chainRegisters("silently-omitted", 100, 1)
	session.poll()

	requireCompleteWatchSet(t, fixture, first, silentlyOmitted)
	assert.Equal(t, []string{
		fixture.depositAddressOf(first),
		fixture.depositAddressOf(first),
		fixture.depositAddressOf(silentlyOmitted),
	}, fixture.importedAddresses())
	requireMatchingCheckpoint(t, fixture, silentlyOmitted.RegistrationRoot, 101)
	assert.Equal(t, 2, replayCountFromDeployment(fixture),
		"the health mismatch must use the cold replay procedure")
	assert.Len(t, fixture.importedAddresses(), 3, "health repair must reimport the pruned registration")
	mismatchSignal := requireOneIntegritySignal(t, fixture, 101, silentlyOmitted.RegistrationRoot)
	assert.Equal(t, first.RegistrationRoot, mismatchSignal.LocalRoot)
	logEntry := requireSingleMismatchLog(t, capturedLogs()[logOffset:])
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
) []rootstock.PegInWatch {
	t.Helper()
	entries := requireCompleteWatchSet(t, fixture, events...)
	expectedImports := make([]string, 0, len(events))
	for _, event := range events {
		expectedImports = append(expectedImports, fixture.depositAddressOf(event))
	}
	assert.Equal(t, expectedImports, fixture.importedAddresses())
	return entries
}

func requireCompleteWatchSet(
	t *testing.T,
	fixture *registryRestartFixture,
	events ...blockchain.AddressRegistered,
) []rootstock.PegInWatch {
	t.Helper()
	entries := fixture.watchSet()
	require.Len(t, entries, len(events))
	for index, event := range events {
		assert.Equal(t, event.TxHash, entries[index].TxHash)
		assert.Equal(t, event.LogIndex, entries[index].LogIndex)
		assert.Equal(t, event.BlockNumber, entries[index].BlockNumber)
		assert.Equal(t, event.RskAddress, entries[index].RskAddress)
		assert.Equal(t, event.RegistrationRoot, entries[index].RegistrationRoot)
		assert.Equal(t, rootstock.PegInWatchImported, entries[index].State)
		assert.Equal(t, fixture.depositAddressOf(event), entries[index].BtcAddress)
		assert.Empty(t, entries[index].LastError)
	}
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
