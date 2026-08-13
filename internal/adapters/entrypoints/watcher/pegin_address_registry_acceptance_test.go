package watcher_test

import (
	"fmt"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
