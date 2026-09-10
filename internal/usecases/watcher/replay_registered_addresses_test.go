package watcher_test

import (
	"context"
	"math"
	"sort"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	addressA = "0x00000000000000000000000000000000000000a1"
	addressB = "0x00000000000000000000000000000000000000b2"
	addressC = "0x00000000000000000000000000000000000000c3"
	addressD = "0x00000000000000000000000000000000000000d4"
)

type memoryWatchRepository struct {
	mu            sync.Mutex
	rows          []rootstock.PegInWatch
	deletes       []uint64
	listCalls     int
	listErrAt     map[int]error
	corruptListAt int
}

func (repository *memoryWatchRepository) Upsert(_ context.Context, watch rootstock.PegInWatch) error {
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

func (repository *memoryWatchRepository) Get(
	_ context.Context,
	rskAddress string,
) (*rootstock.PegInWatch, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	for index := range repository.rows {
		if repository.rows[index].RskAddress == rskAddress {
			entry := repository.rows[index]
			return &entry, nil
		}
	}
	return nil, nil
}

func (repository *memoryWatchRepository) List(context.Context) ([]rootstock.PegInWatch, error) {
	repository.mu.Lock()
	repository.listCalls++
	call := repository.listCalls
	err := repository.listErrAt[call]
	rows := append([]rootstock.PegInWatch(nil), repository.rows...)
	if repository.corruptListAt == call && len(rows) != 0 {
		rows = rows[:len(rows)-1]
	}
	repository.mu.Unlock()
	return rows, err
}

func (repository *memoryWatchRepository) Update(_ context.Context, watch rootstock.PegInWatch) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	for index := range repository.rows {
		if repository.rows[index].RskAddress == watch.RskAddress {
			repository.rows[index] = watch
			return nil
		}
	}
	return nil
}

func (repository *memoryWatchRepository) DeleteFromBlock(_ context.Context, fromBlock uint64) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	repository.deletes = append(repository.deletes, fromBlock)
	kept := repository.rows[:0]
	for _, row := range repository.rows {
		if row.BlockNumber < fromBlock {
			kept = append(kept, row)
		}
	}
	repository.rows = kept
	return nil
}

type memoryCheckpointRepository struct {
	checkpoint rootstock.PegInWatchCheckpoint
	found      bool
	sets       []rootstock.PegInWatchCheckpoint
}

type fakeRegistry struct {
	blockchain.PegInAddressRegistryContract
	events              []blockchain.AddressRegistered
	roots               map[uint64][32]byte
	rootErrors          map[uint64]error
	fetchError          error
	rootReads           []uint64
	eventReads          [][2]uint64
	deleteCountsAtFetch []int
	repository          *memoryWatchRepository
}

func (registry *fakeRegistry) GetRegistrationRoot(
	_ context.Context,
	blockNumber uint64,
) ([32]byte, error) {
	registry.rootReads = append(registry.rootReads, blockNumber)
	if err := registry.rootErrors[blockNumber]; err != nil {
		return [32]byte{}, err
	}
	return registry.roots[blockNumber], nil
}

func (registry *fakeRegistry) GetAddressRegisteredEvents(
	_ context.Context,
	fromBlock uint64,
	toBlock *uint64,
) ([]blockchain.AddressRegistered, error) {
	to := ^uint64(0)
	if toBlock != nil {
		to = *toBlock
	}
	registry.eventReads = append(registry.eventReads, [2]uint64{fromBlock, to})
	registry.deleteCountsAtFetch = append(registry.deleteCountsAtFetch, len(registry.repository.deletes))
	if registry.fetchError != nil {
		return nil, registry.fetchError
	}
	result := make([]blockchain.AddressRegistered, 0)
	for _, event := range registry.events {
		if event.BlockNumber >= fromBlock && event.BlockNumber <= to {
			result = append(result, event)
		}
	}
	return result, nil
}

func (registry *fakeRegistry) GetPegInAddress(string) (blockchain.PegInAddress, error) {
	return blockchain.PegInAddress{Encoding: blockchain.PegInAddressRegistryEncodingBech32}, nil
}

type fakeRootstockRPC struct {
	blockchain.RootstockRpcServer
	head uint64
	err  error
}

func (rpc *fakeRootstockRPC) GetHeight(context.Context) (uint64, error) {
	return rpc.head, rpc.err
}

type replayScenario struct {
	repository  *memoryWatchRepository
	checkpoints *memoryCheckpointRepository
	registry    *fakeRegistry
	rpc         *fakeRootstockRPC
	useCase     *watcher.ReplayRegisteredAddressesUseCase
}

func newReplayScenario(t *testing.T, head uint64) *replayScenario {
	t.Helper()
	repository := &memoryWatchRepository{listErrAt: make(map[int]error)}
	checkpoints := &memoryCheckpointRepository{}
	registry := &fakeRegistry{
		roots:      make(map[uint64][32]byte),
		rootErrors: make(map[uint64]error),
		repository: repository,
	}
	rpc := &fakeRootstockRPC{head: head}
	scenario := &replayScenario{
		repository:  repository,
		checkpoints: checkpoints,
		registry:    registry,
		rpc:         rpc,
	}
	var err error
	scenario.useCase, err = watcher.NewReplayRegisteredAddressesUseCase(
		repository,
		registry,
		rpc,
		crypto.Keccak256,
		100,
		10,
	)
	require.NoError(t, err)
	return scenario
}

func (scenario *replayScenario) run(
	ctx context.Context,
	startBlock uint64,
	pageSize uint64,
) (watcher.ReplayResult, error) {
	useCase := scenario.useCase
	if startBlock != 100 || pageSize != 10 {
		var err error
		useCase, err = watcher.NewReplayRegisteredAddressesUseCase(
			scenario.repository,
			scenario.registry,
			scenario.rpc,
			crypto.Keccak256,
			startBlock,
			pageSize,
		)
		if err != nil {
			return watcher.ReplayResult{}, err
		}
	}
	var checkpoint *rootstock.PegInWatchCheckpoint
	if scenario.checkpoints.found {
		checkpointCopy := scenario.checkpoints.checkpoint
		checkpoint = &checkpointCopy
	}
	result, err := useCase.Run(ctx, checkpoint)
	if err == nil && result.Checkpoint != nil {
		scenario.checkpoints.checkpoint = *result.Checkpoint
		scenario.checkpoints.found = true
		scenario.checkpoints.sets = append(scenario.checkpoints.sets, *result.Checkpoint)
	}
	return result, err
}

func watchRow(block uint64, logIndex uint, txHash string, address string) rootstock.PegInWatch {
	return rootstock.PegInWatch{
		BlockNumber: block,
		LogIndex:    logIndex,
		TxHash:      txHash,
		RskAddress:  address,
		State:       rootstock.PegInWatchImported,
	}
}

func chainEvents(
	t *testing.T,
	registrations ...rootstock.PegInWatch,
) []blockchain.AddressRegistered {
	t.Helper()
	ordered := append([]rootstock.PegInWatch(nil), registrations...)
	sort.SliceStable(ordered, func(first, second int) bool {
		if ordered[first].BlockNumber == ordered[second].BlockNumber {
			return ordered[first].LogIndex < ordered[second].LogIndex
		}
		return ordered[first].BlockNumber < ordered[second].BlockNumber
	})
	root := [32]byte{}
	events := make([]blockchain.AddressRegistered, 0, len(ordered))
	for _, registration := range ordered {
		var err error
		root, err = blockchain.FoldPegInAddressRegistryRoot(
			crypto.Keccak256,
			root,
			registration.RskAddress,
		)
		require.NoError(t, err)
		events = append(events, blockchain.AddressRegistered{
			BlockNumber:      registration.BlockNumber,
			LogIndex:         registration.LogIndex,
			TxHash:           registration.TxHash,
			RskAddress:       registration.RskAddress,
			RegistrationRoot: root,
		})
	}
	return events
}

func rootAt(t *testing.T, events []blockchain.AddressRegistered, height uint64) [32]byte {
	t.Helper()
	root := [32]byte{}
	for _, event := range events {
		if event.BlockNumber > height {
			break
		}
		var err error
		root, err = blockchain.FoldPegInAddressRegistryRoot(crypto.Keccak256, root, event.RskAddress)
		require.NoError(t, err)
	}
	return root
}

func configureChain(t *testing.T, scenario *replayScenario, events []blockchain.AddressRegistered) {
	t.Helper()
	scenario.registry.events = events
	for height := uint64(0); height <= scenario.rpc.head; height++ {
		scenario.registry.roots[height] = rootAt(t, events, height)
		if height == scenario.rpc.head {
			break
		}
	}
}

func TestNewReplayRegisteredAddressesUseCase(t *testing.T) {
	require.NotNil(t, newReplayScenario(t, 0).useCase)
}

func TestNewReplayRegisteredAddressesUseCase_RejectsZeroPageSize(t *testing.T) {
	useCase, err := watcher.NewReplayRegisteredAddressesUseCase(
		nil,
		nil,
		nil,
		crypto.Keccak256,
		100,
		0,
	)

	require.ErrorContains(t, err, "replay page size must be greater than zero")
	assert.Nil(t, useCase)
}

func TestReplayRegisteredAddressesUseCase_RunReturnsPerRunCheckpoint(t *testing.T) {
	useCase, err := watcher.NewReplayRegisteredAddressesUseCase(
		nil,
		nil,
		&fakeRootstockRPC{head: 99},
		crypto.Keccak256,
		100,
		10,
	)
	require.NoError(t, err)

	result, err := useCase.Run(context.Background(), nil)

	require.NoError(t, err)
	assert.Nil(t, result.Checkpoint)
}

func TestReplayRegisteredAddressesUseCase_Run_StopsBeforeStateReadsWhenHeadIsBelowStart(t *testing.T) {
	scenario := newReplayScenario(t, 99)

	result, err := scenario.run(context.Background(), 100, 10)

	require.NoError(t, err)
	assert.Empty(t, result.Reason)
	assert.Nil(t, result.Checkpoint)
	assert.Nil(t, result.Mismatch)
	assert.Equal(t, 0, scenario.repository.listCalls)
	assert.Empty(t, scenario.registry.rootReads)
}

func TestReplayRegisteredAddressesUseCase_Run_WrapsHeadErrorWithoutMutation(t *testing.T) {
	scenario := newReplayScenario(t, 100)
	scenario.rpc.err = assert.AnError

	_, err := scenario.run(context.Background(), 100, 10)

	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, string(usecases.ReplayRegisteredAddressesId))
	assert.Empty(t, scenario.repository.deletes)
	assert.Empty(t, scenario.checkpoints.sets)
}

func TestReplayRegisteredAddressesUseCase_Run_HealthyTimelineIgnoresStoredRootsAndUsesDeterministicOrder(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	rows := []rootstock.PegInWatch{
		watchRow(102, 1, "tx-b", addressD),
		watchRow(100, 0, "tx-a", addressA),
		watchRow(102, 1, "tx-a", addressC),
		watchRow(102, 0, "tx-z", addressB),
	}
	for index := range rows {
		rows[index].RegistrationRoot = [32]byte{byte(index + 1)}
	}
	scenario.repository.rows = rows
	expectedOrder := chainEvents(t,
		watchRow(100, 0, "tx-a", addressA),
		watchRow(102, 0, "tx-z", addressB),
		watchRow(102, 1, "tx-a", addressC),
		watchRow(102, 1, "tx-b", addressD),
	)
	scenario.registry.roots[104] = rootAt(t, expectedOrder, 104)

	result, err := scenario.run(context.Background(), 100, 10)

	require.NoError(t, err)
	assert.Empty(t, result.Reason)
	assert.Nil(t, result.Mismatch)
	assert.Equal(t, []uint64{104}, scenario.registry.rootReads)
	assert.Empty(t, scenario.registry.eventReads)
	assert.Empty(t, scenario.repository.deletes)
	require.Len(t, scenario.checkpoints.sets, 1)
	assert.Equal(t, scenario.registry.roots[104], scenario.checkpoints.sets[0].LocalRoot)
}

func TestReplayRegisteredAddressesUseCase_Run_EmptyDatabaseWithZeroRootCheckpointsWithoutLogs(t *testing.T) {
	scenario := newReplayScenario(t, 100)

	_, err := scenario.run(context.Background(), 100, 10)

	require.NoError(t, err)
	assert.Empty(t, scenario.registry.eventReads)
	assert.Empty(t, scenario.repository.deletes)
	assert.Equal(t, []rootstock.PegInWatchCheckpoint{{
		VerifiedThroughBlock: 100,
	}}, scenario.checkpoints.sets)
}

func TestReplayRegisteredAddressesUseCase_Run_HealthyStaleCheckpointAdvances(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	events := chainEvents(t, watchRow(100, 0, "a", addressA))
	configureChain(t, scenario, events)
	scenario.repository.rows = []rootstock.PegInWatch{watchRow(100, 0, "a", addressA)}
	scenario.checkpoints.checkpoint = rootstock.PegInWatchCheckpoint{
		LocalRoot:            rootAt(t, events, 100),
		VerifiedThroughBlock: 100,
	}
	scenario.checkpoints.found = true

	_, err := scenario.run(context.Background(), 100, 10)

	require.NoError(t, err)
	assert.Empty(t, scenario.repository.deletes)
	assert.Equal(t, uint64(104), scenario.checkpoints.checkpoint.VerifiedThroughBlock)
	assert.Empty(t, scenario.registry.eventReads)
}

func TestReplayRegisteredAddressesUseCase_Run_TrustedCheckpointReplaysOnlyFollowingRange(t *testing.T) {
	logHook := logtest.NewGlobal()
	defer logHook.Reset()
	scenario := newReplayScenario(t, 104)
	events := chainEvents(t,
		watchRow(100, 0, "a", addressA),
		watchRow(104, 0, "b", addressB),
	)
	configureChain(t, scenario, events)
	scenario.repository.rows = []rootstock.PegInWatch{watchRow(100, 0, "a", addressA)}
	scenario.checkpoints.checkpoint = rootstock.PegInWatchCheckpoint{
		LocalRoot:            rootAt(t, events, 100),
		VerifiedThroughBlock: 100,
	}
	scenario.checkpoints.found = true

	result, err := scenario.run(context.Background(), 100, 10)

	require.NoError(t, err)
	assert.Equal(t, blockchain.PegInAddressRegistryRecoveryCatchUp, result.Reason)
	assert.Nil(t, result.Mismatch)
	assert.Equal(t, []uint64{101}, scenario.repository.deletes)
	assert.Equal(t, [][2]uint64{{101, 104}}, scenario.registry.eventReads)
	assert.Equal(t, []uint64{104, 100}, scenario.registry.rootReads)
	for _, entry := range logHook.AllEntries() {
		assert.NotEqual(t, "PegIn address registry root mismatch", entry.Message)
	}
}

func TestReplayRegisteredAddressesUseCase_Run_EmptyRowsCatchUpFromTrustedZeroRootCheckpoint(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	events := chainEvents(t, watchRow(104, 0, "a", addressA))
	configureChain(t, scenario, events)
	scenario.checkpoints.checkpoint = rootstock.PegInWatchCheckpoint{
		VerifiedThroughBlock: 100,
	}
	scenario.checkpoints.found = true

	_, err := scenario.run(context.Background(), 100, 10)

	require.NoError(t, err)
	assert.Equal(t, []uint64{101}, scenario.repository.deletes)
	assert.Equal(t, [][2]uint64{{101, 104}}, scenario.registry.eventReads)
	assert.Equal(t, []uint64{104, 100}, scenario.registry.rootReads)
}

func TestReplayRegisteredAddressesUseCase_Run_CatchUpPersistsDiscoveredWithoutImporting(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	events := chainEvents(t, watchRow(104, 0, "a", addressA))
	configureChain(t, scenario, events)
	scenario.checkpoints.checkpoint = rootstock.PegInWatchCheckpoint{
		VerifiedThroughBlock: 100,
	}
	scenario.checkpoints.found = true

	_, err := scenario.run(context.Background(), 100, 10)

	require.NoError(t, err)
	require.Len(t, scenario.repository.rows, 1)
	assert.Equal(t, rootstock.PegInWatchDiscovered, scenario.repository.rows[0].State)
	assert.Equal(t, 2, scenario.repository.listCalls,
		"Replay must only list for initial and persisted-state verification; import belongs to Discover")
}

func TestReplayRegisteredAddressesUseCase_Run_TrustedCatchUpReportsEventRootMismatch(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	events := chainEvents(t,
		watchRow(100, 0, "a", addressA),
		watchRow(104, 0, "b", addressB),
	)
	configureChain(t, scenario, events)
	scenario.repository.rows = []rootstock.PegInWatch{watchRow(100, 0, "a", addressA)}
	scenario.checkpoints.checkpoint = rootstock.PegInWatchCheckpoint{
		LocalRoot:            rootAt(t, events, 100),
		VerifiedThroughBlock: 100,
	}
	scenario.checkpoints.found = true
	scenario.registry.events[1].RegistrationRoot = [32]byte{9}

	_, err := scenario.run(context.Background(), 100, 10)

	require.Error(t, err)
}

func TestReplayRegisteredAddressesUseCase_Run_RootAtBetweenEventsFindsHeadBoundary(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	events := chainEvents(t,
		watchRow(100, 0, "a", addressA),
		watchRow(102, 0, "b", addressB),
		watchRow(104, 0, "c", addressC),
	)
	configureChain(t, scenario, events)
	scenario.repository.rows = []rootstock.PegInWatch{
		watchRow(100, 0, "a", addressA),
		watchRow(102, 0, "b", addressB),
	}

	_, err := scenario.run(context.Background(), 100, 10)

	require.NoError(t, err)
	assert.Equal(t, []uint64{104}, scenario.repository.deletes)
	assert.Contains(t, scenario.registry.rootReads, uint64(103))
	assert.Equal(t, [][2]uint64{{104, 104}}, scenario.registry.eventReads)
}

type mismatchBoundaryCase struct {
	name       string
	startBlock uint64
	head       uint64
	chainRows  []rootstock.PegInWatch
	localRows  []rootstock.PegInWatch
	wantDelete uint64
}

func mismatchBoundaryCases() []mismatchBoundaryCase {
	return []mismatchBoundaryCase{
		{
			name:       "first",
			startBlock: 100,
			head:       104,
			chainRows:  []rootstock.PegInWatch{watchRow(100, 0, "a", addressA)},
			localRows:  []rootstock.PegInWatch{watchRow(100, 0, "x", addressD)},
			wantDelete: 100,
		},
		{
			name:       "middle",
			startBlock: 100,
			head:       104,
			chainRows: []rootstock.PegInWatch{
				watchRow(100, 0, "a", addressA),
				watchRow(102, 0, "b", addressB),
			},
			localRows: []rootstock.PegInWatch{
				watchRow(100, 0, "a", addressA),
				watchRow(102, 0, "x", addressD),
			},
			wantDelete: 102,
		},
		{
			name:       "head",
			startBlock: 100,
			head:       104,
			chainRows: []rootstock.PegInWatch{
				watchRow(100, 0, "a", addressA),
				watchRow(104, 0, "b", addressB),
			},
			localRows:  []rootstock.PegInWatch{watchRow(100, 0, "a", addressA)},
			wantDelete: 104,
		},
		{
			name:       "start zero",
			startBlock: 0,
			head:       1,
			chainRows:  []rootstock.PegInWatch{watchRow(0, 0, "a", addressA)},
			localRows:  []rootstock.PegInWatch{watchRow(0, 0, "x", addressD)},
			wantDelete: 0,
		},
	}
}

func TestReplayRegisteredAddressesUseCase_Run_FindsMismatchBoundaries(t *testing.T) {
	for _, testCase := range mismatchBoundaryCases() {
		t.Run(testCase.name, func(t *testing.T) {
			scenario := newReplayScenario(t, testCase.head)
			events := chainEvents(t, testCase.chainRows...)
			configureChain(t, scenario, events)
			scenario.repository.rows = append([]rootstock.PegInWatch(nil), testCase.localRows...)
			localRootAtHead := rootAt(t, chainEvents(t, testCase.localRows...), testCase.head)

			result, err := scenario.run(context.Background(), testCase.startBlock, 10)

			require.NoError(t, err)
			assert.Equal(t, blockchain.PegInAddressRegistryRecoveryRootMismatch, result.Reason)
			assert.Equal(t, &watcher.ReplayRootMismatch{
				BlockNumber: testCase.head,
				LocalRoot:   localRootAtHead,
				ChainRoot:   scenario.registry.roots[testCase.head],
				Source:      "captured_head",
			}, result.Mismatch)
			assert.Equal(t, []uint64{testCase.wantDelete}, scenario.repository.deletes)
		})
	}
}

func TestReplayRegisteredAddressesUseCase_Run_NonMonotoneSearchReturnsFalseBoundaryAfterTruePredecessor(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	events := chainEvents(t,
		watchRow(100, 0, "a", addressA),
		watchRow(104, 0, "b", addressB),
	)
	configureChain(t, scenario, events)
	scenario.repository.rows = []rootstock.PegInWatch{watchRow(100, 0, "a", addressA)}
	localRoot := rootAt(t, chainEvents(t, watchRow(100, 0, "a", addressA)), 104)
	scenario.registry.roots[100] = localRoot
	scenario.registry.roots[101] = [32]byte{1}
	scenario.registry.roots[102] = localRoot
	scenario.registry.roots[103] = localRoot

	_, err := scenario.run(context.Background(), 100, 10)

	require.NoError(t, err)
	require.Equal(t, []uint64{104}, scenario.repository.deletes)
	assert.Equal(t, localRoot, scenario.registry.roots[103])
	assert.NotEqual(t, localRoot, scenario.registry.roots[104])
	assert.Equal(t, []uint64{104, 102, 103}, scenario.registry.rootReads)
	assert.Equal(t, [][2]uint64{{104, 104}}, scenario.registry.eventReads)
}

func TestReplayRegisteredAddressesUseCase_Run_BoundsSearchRootReadsAndFetchesLogsAfterSearch(t *testing.T) {
	scenario := newReplayScenario(t, 199)
	events := chainEvents(t,
		watchRow(100, 0, "a", addressA),
		watchRow(199, 0, "b", addressB),
	)
	configureChain(t, scenario, events)
	scenario.repository.rows = []rootstock.PegInWatch{watchRow(100, 0, "a", addressA)}

	_, err := scenario.run(context.Background(), 100, 10)

	require.NoError(t, err)
	searchReads := len(scenario.registry.rootReads) - 1
	assert.LessOrEqual(t, searchReads, 8)
	assert.Equal(t, []uint64{199}, scenario.repository.deletes)
	assert.Equal(t, [][2]uint64{{199, 199}}, scenario.registry.eventReads)
}

func TestReplayRegisteredAddressesUseCase_Run_SearchReadErrorDoesNotMutateOrFetchLogs(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	scenario.repository.rows = []rootstock.PegInWatch{watchRow(100, 0, "a", addressA)}
	scenario.registry.roots[104] = [32]byte{9}
	scenario.registry.rootErrors[102] = assert.AnError

	_, err := scenario.run(context.Background(), 100, 10)

	require.ErrorIs(t, err, assert.AnError)
	assert.Empty(t, scenario.repository.deletes)
	assert.Empty(t, scenario.checkpoints.sets)
	assert.Empty(t, scenario.registry.eventReads)
}

func TestReplayRegisteredAddressesUseCase_Run_FirstStartFindsFirstRegistrationBoundary(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	events := chainEvents(t, watchRow(102, 0, "a", addressA))
	configureChain(t, scenario, events)

	_, err := scenario.run(context.Background(), 100, 2)

	require.NoError(t, err)
	assert.Equal(t, []uint64{102}, scenario.repository.deletes)
	assert.Equal(t, [][2]uint64{{102, 103}, {104, 104}}, scenario.registry.eventReads)
	assert.Equal(t, scenario.registry.roots[104], scenario.checkpoints.checkpoint.LocalRoot)
}

func TestReplayRegisteredAddressesUseCase_Run_EmptyDatabaseSkipsIdleDeploymentRange(t *testing.T) {
	scenario := newReplayScenario(t, 10_000)
	events := chainEvents(t, watchRow(9_000, 0, "a", addressA))
	configureChain(t, scenario, events)

	_, err := scenario.run(context.Background(), 100, 1_000)

	require.NoError(t, err)
	assert.Equal(t, []uint64{9_000}, scenario.repository.deletes)
	assert.Equal(t, [][2]uint64{{9_000, 9_999}, {10_000, 10_000}}, scenario.registry.eventReads)
	assert.Contains(t, scenario.registry.rootReads, uint64(10_000))
	assert.Contains(t, scenario.registry.rootReads, uint64(5_050),
		"boundary search must probe the idle range before fetching events")
}

func TestReplayRegisteredAddressesUseCase_Run_RepairsProjectionDifferences(t *testing.T) {
	chainRows := []rootstock.PegInWatch{
		watchRow(100, 0, "a", addressA),
		watchRow(102, 0, "b", addressB),
		watchRow(104, 0, "c", addressC),
	}
	tests := []struct {
		name      string
		localRows []rootstock.PegInWatch
	}{
		{"missing middle", []rootstock.PegInWatch{
			watchRow(100, 0, "a", addressA),
			watchRow(104, 0, "c", addressC),
		}},
		{"missing tail", chainRows[:2]},
		{"extra row", append(append([]rootstock.PegInWatch(nil), chainRows...),
			watchRow(103, 0, "x", addressD))},
		{"changed address", []rootstock.PegInWatch{
			watchRow(100, 0, "a", addressA),
			watchRow(102, 0, "b", addressD),
			watchRow(104, 0, "c", addressC),
		}},
		{"reordered sequence", []rootstock.PegInWatch{
			watchRow(100, 0, "a", addressA),
			watchRow(102, 1, "b", addressB),
			watchRow(102, 0, "c", addressC),
		}},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			scenario := newReplayScenario(t, 104)
			events := chainEvents(t, chainRows...)
			configureChain(t, scenario, events)
			scenario.repository.rows = append([]rootstock.PegInWatch(nil), testCase.localRows...)

			_, err := scenario.run(context.Background(), 100, 10)

			require.NoError(t, err)
			assert.Equal(t, scenario.registry.roots[104], scenario.checkpoints.checkpoint.LocalRoot)
			assert.NotEmpty(t, scenario.repository.deletes)
		})
	}
}

func TestReplayRegisteredAddressesUseCase_Run_PreservesHeightShiftThatKeepsSequence(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	events := chainEvents(t,
		watchRow(100, 0, "a", addressA),
		watchRow(102, 0, "x", addressB),
		watchRow(104, 0, "y", addressC),
	)
	configureChain(t, scenario, events)
	scenario.repository.rows = []rootstock.PegInWatch{
		watchRow(100, 0, "a", addressA),
		watchRow(101, 0, "x", addressB),
	}

	_, err := scenario.run(context.Background(), 100, 10)

	require.NoError(t, err)
	assert.Equal(t, []uint64{104}, scenario.repository.deletes)
	require.Len(t, scenario.repository.rows, 3)
	assert.Equal(t, uint64(101), scenario.repository.rows[1].BlockNumber)
}

func TestReplayRegisteredAddressesUseCase_Run_MultipleEventsInBlockUseLogOrder(t *testing.T) {
	scenario := newReplayScenario(t, 100)
	events := chainEvents(t,
		watchRow(100, 1, "b", addressB),
		watchRow(100, 0, "a", addressA),
	)
	configureChain(t, scenario, events)
	scenario.registry.events = []blockchain.AddressRegistered{events[1], events[0]}

	_, err := scenario.run(context.Background(), 100, 10)

	require.NoError(t, err)
	assert.Equal(t, scenario.registry.roots[100], scenario.checkpoints.checkpoint.LocalRoot)
}

func TestReplayRegisteredAddressesUseCase_Run_DoesNotTrustChangedOrFutureCheckpoint(t *testing.T) {
	for _, checkpoint := range []rootstock.PegInWatchCheckpoint{
		{LocalRoot: [32]byte{7}, VerifiedThroughBlock: 100},
		{VerifiedThroughBlock: 105},
	} {
		scenario := newReplayScenario(t, 104)
		events := chainEvents(t,
			watchRow(100, 0, "a", addressA),
			watchRow(104, 0, "b", addressB),
		)
		configureChain(t, scenario, events)
		scenario.repository.rows = []rootstock.PegInWatch{watchRow(100, 0, "a", addressA)}
		scenario.checkpoints.checkpoint = checkpoint
		scenario.checkpoints.found = true

		_, err := scenario.run(context.Background(), 100, 10)

		require.NoError(t, err)
		assert.Equal(t, []uint64{104}, scenario.repository.deletes)
	}
}

// A checkpoint sitting on the captured head cannot be trusted for catch-up: the run only reaches
// this path because the local and chain roots differ at the head, so a checkpoint that matched both
// there would contradict that reading. Resuming from it would ask for the block after the head.
func TestReplayRegisteredAddressesUseCase_Run_DoesNotCatchUpFromCheckpointAtCapturedHead(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	events := chainEvents(t,
		watchRow(100, 0, "a", addressA),
		watchRow(104, 0, "b", addressB),
	)
	configureChain(t, scenario, events)
	scenario.repository.rows = []rootstock.PegInWatch{watchRow(100, 0, "a", addressA)}
	scenario.checkpoints.checkpoint = rootstock.PegInWatchCheckpoint{
		LocalRoot:            rootAt(t, events, 100),
		VerifiedThroughBlock: 104,
	}
	scenario.checkpoints.found = true

	_, err := scenario.run(context.Background(), 100, 10)

	require.NoError(t, err)
	assert.Equal(t, []uint64{104}, scenario.repository.deletes)
	assert.Equal(t, [][2]uint64{{104, 104}}, scenario.registry.eventReads)
	assert.Equal(t, []uint64{104, 102, 103}, scenario.registry.rootReads,
		"the boundary search must run instead of the checkpoint being trusted at the head")
	assert.Equal(t, scenario.registry.roots[104], scenario.checkpoints.checkpoint.LocalRoot)
}

// No stored row can sit above the highest representable block, so a head there is reconciled
// normally rather than treated as an overflow hazard.
func TestReplayRegisteredAddressesUseCase_Run_ReconcilesHighestRepresentableHead(t *testing.T) {
	scenario := newReplayScenario(t, math.MaxUint64)
	rows := []rootstock.PegInWatch{watchRow(math.MaxUint64, 0, "a", addressA)}
	scenario.repository.rows = append([]rootstock.PegInWatch(nil), rows...)
	scenario.registry.roots[math.MaxUint64] = rootAt(t, chainEvents(t, rows...), math.MaxUint64)

	_, err := scenario.run(context.Background(), 100, 10)

	require.NoError(t, err)
	assert.Empty(t, scenario.repository.deletes, "no row can be above the highest block")
	assert.Empty(t, scenario.registry.eventReads)
	assert.Equal(t, []rootstock.PegInWatchCheckpoint{{
		LocalRoot:            scenario.registry.roots[math.MaxUint64],
		VerifiedThroughBlock: math.MaxUint64,
	}}, scenario.checkpoints.sets)
}

func TestReplayRegisteredAddressesUseCase_Run_RowBelowStartFailsWithoutMutation(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	scenario.repository.rows = []rootstock.PegInWatch{watchRow(99, 0, "a", addressA)}

	_, err := scenario.run(context.Background(), 100, 10)

	require.ErrorContains(t, err, "PegIn watch exists below configured start block 100")
	assert.Equal(t, []uint64{104}, scenario.registry.rootReads)
	assert.Empty(t, scenario.repository.deletes)
	assert.Empty(t, scenario.checkpoints.sets)
}

func TestReplayRegisteredAddressesUseCase_Run_RowAboveHeadFailsWithoutMutation(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	scenario.repository.rows = []rootstock.PegInWatch{watchRow(105, 0, "a", addressA)}

	_, err := scenario.run(context.Background(), 100, 10)

	require.ErrorContains(t, err, "PegIn watch exists above current head 104")
	assert.Equal(t, []uint64{104}, scenario.registry.rootReads)
	assert.Empty(t, scenario.repository.deletes)
	require.Len(t, scenario.repository.rows, 1)
	assert.Empty(t, scenario.checkpoints.sets)
}

func TestReplayRegisteredAddressesUseCase_Run_RowAboveHeadSurvivesRootReadError(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	scenario.repository.rows = []rootstock.PegInWatch{watchRow(105, 0, "a", addressA)}
	scenario.registry.rootErrors[104] = assert.AnError

	_, err := scenario.run(context.Background(), 100, 10)

	require.ErrorIs(t, err, assert.AnError)
	assert.Empty(t, scenario.repository.deletes)
	require.Len(t, scenario.repository.rows, 1)
}

func TestReplayRegisteredAddressesUseCase_Run_FetchFailureDoesNotDeleteSuffix(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	scenario.registry.roots[104] = [32]byte{9}
	scenario.registry.fetchError = assert.AnError
	original := rootstock.PegInWatchCheckpoint{LocalRoot: [32]byte{3}, VerifiedThroughBlock: 99}
	scenario.checkpoints.checkpoint = original
	scenario.checkpoints.found = true

	_, err := scenario.run(context.Background(), 100, 10)

	require.ErrorIs(t, err, assert.AnError)
	assert.Empty(t, scenario.repository.deletes)
	assert.Equal(t, original, scenario.checkpoints.checkpoint)
	assert.Empty(t, scenario.checkpoints.sets)
}

func TestReplayRegisteredAddressesUseCase_Run_FoldFailureDoesNotDeleteSuffix(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	scenario.registry.roots[104] = [32]byte{9}
	scenario.registry.events = []blockchain.AddressRegistered{{
		BlockNumber: 104,
		LogIndex:    0,
		TxHash:      "bad-address",
		RskAddress:  "not-an-address",
	}}

	_, err := scenario.run(context.Background(), 100, 10)

	require.ErrorContains(t, err, "fold AddressRegistered event bad-address/0")
	assert.Empty(t, scenario.repository.deletes)
}

func TestReplayRegisteredAddressesUseCase_Run_EventMismatchDoesNotDeleteOrLog(t *testing.T) {
	logHook := logtest.NewGlobal()
	defer logHook.Reset()
	scenario := newReplayScenario(t, 104)
	events := chainEvents(t, watchRow(104, 0, "bad-root", addressA))
	configureChain(t, scenario, events)
	scenario.registry.events[0].RegistrationRoot = [32]byte{9}

	result, err := scenario.run(context.Background(), 100, 10)

	require.ErrorContains(t, err, "event_bad-root_0")
	assert.Empty(t, result.Reason)
	assert.Nil(t, result.Checkpoint)
	assert.Nil(t, result.Mismatch)
	assert.Empty(t, scenario.repository.deletes)
	for _, entry := range logHook.AllEntries() {
		assert.NotEqual(t, "PegIn address registry root mismatch", entry.Message)
	}
}

func TestReplayRegisteredAddressesUseCase_Run_DeletesOnlyAfterFetchingAndFoldingSuffix(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	events := chainEvents(t, watchRow(104, 0, "a", addressA))
	configureChain(t, scenario, events)
	deleteCountsAtFold := make([]int, 0)
	hashFunction := func(data ...[]byte) []byte {
		deleteCountsAtFold = append(deleteCountsAtFold, len(scenario.repository.deletes))
		return crypto.Keccak256(data...)
	}
	var err error
	scenario.useCase, err = watcher.NewReplayRegisteredAddressesUseCase(
		scenario.repository,
		scenario.registry,
		scenario.rpc,
		hashFunction,
		100,
		10,
	)
	require.NoError(t, err)

	_, err = scenario.run(context.Background(), 100, 10)

	require.NoError(t, err)
	assert.Equal(t, []int{0}, scenario.registry.deleteCountsAtFetch)
	require.NotEmpty(t, deleteCountsAtFold)
	assert.Zero(t, deleteCountsAtFold[0], "the event must be folded before suffix deletion")
	assert.Equal(t, []uint64{104}, scenario.repository.deletes)
}

func TestReplayRegisteredAddressesUseCase_Run_FinalPersistedMismatchDoesNotCheckpointOrRetry(t *testing.T) {
	scenario := newReplayScenario(t, 104)
	events := chainEvents(t, watchRow(104, 0, "a", addressA))
	configureChain(t, scenario, events)
	scenario.repository.corruptListAt = 2

	_, err := scenario.run(context.Background(), 100, 10)

	require.Error(t, err)
	assert.Equal(t, []uint64{104}, scenario.repository.deletes)
	assert.Empty(t, scenario.checkpoints.sets)
	assert.Equal(t, 2, scenario.repository.listCalls)
}
