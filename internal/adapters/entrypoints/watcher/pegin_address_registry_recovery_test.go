package watcher_test

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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
			name:    "stopped after the imported state was persisted, before the checkpoint advanced",
			stop:    func(store *registryWatchStore) { store.failSetCheckpoint = errProcessStopped },
			imports: 1,
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
