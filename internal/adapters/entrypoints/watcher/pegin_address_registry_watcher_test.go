package watcher_test

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/crypto"
	watcherAdapter "github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test/datasets"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type pegInAddressRegistryWatcherFixture struct {
	t          *testing.T
	repository *mocks.PegInWatchRepositoryMock
	registry   *mocks.PegInAddressRegistryContractMock
	rskRpc     *mocks.RootstockRpcServerMock
	btcNetwork *mocks.BitcoinNetworkMock
	wallet     *mocks.BitcoinWalletMock
	ticker     *sessionTicker
	watcher    *watcherAdapter.PegInWatcher
}

//nolint:unparam // Fixed start/finality values make each scenario's cursor arithmetic explicit at the call site.
func newPegInWatcherFixture(
	t *testing.T,
	startBlock uint64,
	pageSize uint64,
	finalityDepth uint64,
) *pegInAddressRegistryWatcherFixture {
	t.Helper()
	fixture := &pegInAddressRegistryWatcherFixture{
		t:          t,
		repository: mocks.NewPegInWatchRepositoryMock(t),
		registry:   mocks.NewPegInAddressRegistryContractMock(t),
		rskRpc:     mocks.NewRootstockRpcServerMock(t),
		btcNetwork: mocks.NewBitcoinNetworkMock(t),
		wallet:     mocks.NewBitcoinWalletMock(t),
		ticker:     newSessionTicker(),
	}
	useCases := watcherAdapter.NewPegInWatcherUseCases(
		watcher.NewGetWatchedRegisteredAddressesUseCase(fixture.repository),
		watcher.NewCatchUpRegisteredAddressesUseCase(
			fixture.repository,
			fixture.registry,
			fixture.rskRpc,
			crypto.Keccak256,
			startBlock,
			pageSize,
			finalityDepth,
		),
		watcher.NewImportRegisteredAddressesUseCase(fixture.repository, fixture.registry, fixture.btcNetwork, fixture.wallet),
	)
	fixture.watcher = watcherAdapter.NewPegInWatcher(useCases, fixture.ticker)
	return fixture
}

func (fixture *pegInAddressRegistryWatcherFixture) expectBoundedRescan() {
	fixture.t.Helper()
	fixture.btcNetwork.EXPECT().GetHeight().Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{StartHeight: 100, StopHeight: 200}, nil).
		Once()
}

type catchUpExpectation struct {
	checkpoint rootstock.PegInWatchCheckpoint
	head       uint64
	fromBlock  uint64
	toBlock    uint64
	events     []blockchain.AddressRegistered
}

func (fixture *pegInAddressRegistryWatcherFixture) expectCatchUp(
	expected catchUpExpectation,
) *mocks.PegInWatchRepositoryMock_SaveVerifiedRegistrations_Call {
	fixture.t.Helper()
	fixture.repository.EXPECT().GetCheckpoint(mock.Anything).
		Return(expected.checkpoint, true, nil).Once()
	root := foldRegistrations(fixture.t, expected.checkpoint.Root, expected.events)
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(expected.head, nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, expected.fromBlock, uint64Pointer(expected.toBlock)).
		Return(expected.events, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, expected.toBlock).Return(root, nil).Once()
	saveCall := fixture.repository.EXPECT().SaveVerifiedRegistrations(
		mock.Anything,
		mock.Anything,
		rootstock.PegInWatchCheckpoint{Root: root, Block: expected.toBlock},
	).Return(nil)
	saveCall.Once()
	return saveCall
}

func foldRegistrations(t *testing.T, previousRoot [32]byte, events []blockchain.AddressRegistered) [32]byte {
	t.Helper()
	order := make([]int, len(events))
	for index := range order {
		order[index] = index
	}
	slices.SortFunc(order, func(first, second int) int {
		if byBlock := cmp.Compare(events[first].BlockNumber, events[second].BlockNumber); byBlock != 0 {
			return byBlock
		}
		return cmp.Compare(events[first].LogIndex, events[second].LogIndex)
	})
	root := previousRoot
	for _, index := range order {
		var err error
		root, err = blockchain.FoldPegInAddressRegistryRoot(crypto.Keccak256, root, events[index].RskAddress)
		require.NoError(t, err)
		events[index].RegistrationRoot = root
	}
	return root
}

func rskAddress(suffix int) string {
	return fmt.Sprintf("0x%040x", suffix)
}

// runScan is one boot, one tick, one shutdown. The session returns from the tick when the scan has
// finished, so a scenario does not have to nominate a repository call as the end of the poll.
func (fixture *pegInAddressRegistryWatcherFixture) runScan() {
	fixture.t.Helper()
	session := startWatcherSession(fixture.t, fixture.watcher, fixture.ticker)
	session.poll()
	session.stop()
}

// depositAddress pairs a real base58 address with the payload the registry returns for it, so a
// scenario can assert that the scanner imports the address rather than the bytes. The address is a
// known-good constant taken from the shared dataset, not something this file re-derives.
type depositAddress struct {
	payload []byte
	address string
}

// depositAddressFixture builds the registry's return value for one of the dataset addresses:
// version ++ hash ++ four-byte double-SHA256 checksum, exactly as getPegInAddress emits it.
func depositAddressFixture(index int) depositAddress {
	const checksumSize = 4
	decoded := datasets.Base58Addresses[index]
	payload := make([]byte, 0, len(decoded.Expected)+checksumSize)
	payload = append(payload, decoded.Expected...)
	payload = append(payload, chainhash.DoubleHashB(decoded.Expected)[:checksumSize]...)
	return depositAddress{payload: payload, address: decoded.Address}
}

func discoveredWatchEntry(event blockchain.AddressRegistered) rootstock.PegInWatch {
	return rootstock.PegInWatch{
		TxHash:      event.TxHash,
		LogIndex:    event.LogIndex,
		BlockNumber: event.BlockNumber,
		RskAddress:  event.RskAddress,
		State:       rootstock.PegInWatchDiscovered,
	}
}

func TestPegInWatcher_ContinuesFromStoredCheckpoint(t *testing.T) {
	fixture := newPegInWatcherFixture(t, 100, 3, 2)
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{{
		TxHash: "persisted",
		State:  rootstock.PegInWatchImported,
	}}, nil).Once()
	fixture.expectCatchUp(catchUpExpectation{
		checkpoint: rootstock.PegInWatchCheckpoint{Root: [32]byte{1}, Block: 105},
		head:       110,
		fromBlock:  106,
		toBlock:    108,
	})

	fixture.runScan()
}

//nolint:funlen // One scenario records the complete persist/import/cursor ordering across three events.
func TestPegInWatcher_SortsEventsAndClampsToFinalizedHead(t *testing.T) {
	fixture := newPegInWatcherFixture(t, 100, 20, 2)
	events := []blockchain.AddressRegistered{
		{TxHash: "third", RskAddress: rskAddress(3), BlockNumber: 102, LogIndex: 5},
		{TxHash: "first", RskAddress: rskAddress(1), BlockNumber: 101, LogIndex: 9},
		{TxHash: "second", RskAddress: rskAddress(2), BlockNumber: 102, LogIndex: 1},
	}
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{}, nil).Once()
	var operations []string
	var operationsMutex sync.Mutex
	fixture.expectCatchUp(catchUpExpectation{checkpoint: rootstock.PegInWatchCheckpoint{Block: 99}, head: 105, fromBlock: 100, toBlock: 103, events: events}).
		Run(func(_ context.Context, saved []rootstock.PegInWatch, _ rootstock.PegInWatchCheckpoint) {
			operationsMutex.Lock()
			for _, entry := range saved {
				operations = append(operations, fmt.Sprintf("save:%d/%d/%s", entry.BlockNumber, entry.LogIndex, entry.TxHash))
			}
			operationsMutex.Unlock()
		})
	for index, event := range events {
		deposit := depositAddressFixture(index)
		persisted := discoveredWatchEntry(event)
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
		fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(entry rootstock.PegInWatch) bool {
			return entry.TxHash == event.TxHash &&
				entry.BtcAddress == deposit.address &&
				entry.Encoding == uint8(blockchain.PegInAddressRegistryEncodingBase58) &&
				entry.State == rootstock.PegInWatchImported
		})).Run(func(mock.Arguments) {
			operationsMutex.Lock()
			operations = append(operations, "update:"+event.TxHash)
			operationsMutex.Unlock()
		}).Return(nil).Once()
	}
	fixture.btcNetwork.EXPECT().GetHeight().Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().RescanBlockchain(int64(100)).
		Run(func(int64) {
			operationsMutex.Lock()
			operations = append(operations, "rescan")
			operationsMutex.Unlock()
		}).
		Return(blockchain.BitcoinRescanResult{StartHeight: 100, StopHeight: 200}, nil).
		Once()
	fixture.runScan()

	operationsMutex.Lock()
	assert.Equal(t, []string{
		"save:101/9/first", "save:102/1/second", "save:102/5/third",
		"import:first", "import:second", "import:third",
		"rescan",
		"update:first", "update:second", "update:third",
	}, operations)
	operationsMutex.Unlock()
}

func TestPegInWatcher_SkipsUnsupportedEncodingsAndContinues(t *testing.T) {
	fixture := newPegInWatcherFixture(t, 100, 1, 2)
	events := []blockchain.AddressRegistered{
		{TxHash: "bech32", RskAddress: rskAddress(1), BlockNumber: 100, LogIndex: 1},
		{TxHash: "bech32m", RskAddress: rskAddress(2), BlockNumber: 100, LogIndex: 2},
		{TxHash: "base58", RskAddress: rskAddress(3), BlockNumber: 100, LogIndex: 3},
	}
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{}, nil).Once()
	fixture.expectCatchUp(catchUpExpectation{checkpoint: rootstock.PegInWatchCheckpoint{Block: 99}, head: 102, fromBlock: 100, toBlock: 100, events: events})

	encodings := []blockchain.PegInAddressRegistryEncoding{
		blockchain.PegInAddressRegistryEncodingBech32,
		blockchain.PegInAddressRegistryEncodingBech32M,
		blockchain.PegInAddressRegistryEncodingBase58,
	}
	for index, event := range events {
		deposit := depositAddressFixture(index)
		persisted := discoveredWatchEntry(event)
		fixture.repository.EXPECT().Get(mock.Anything, event.RskAddress).Return(&persisted, nil).Once()
		fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).
			Return(blockchain.PegInAddress{
				Payload:  deposit.payload,
				Encoding: encodings[index],
			}, nil).
			Once()
		expectedState := rootstock.PegInWatchUnsupportedEncoding
		// An encoding the scanner cannot render must leave BtcAddress empty rather than storing the
		// raw payload, which is not a valid address in any encoding.
		expectedAddress := ""
		if encodings[index] == blockchain.PegInAddressRegistryEncodingBase58 {
			expectedState = rootstock.PegInWatchImported
			expectedAddress = deposit.address
			fixture.wallet.EXPECT().ImportAddress(deposit.address).Return(nil).Once()
			fixture.expectBoundedRescan()
		}
		fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(entry rootstock.PegInWatch) bool {
			return entry.TxHash == event.TxHash &&
				entry.State == expectedState &&
				entry.BtcAddress == expectedAddress
		})).Return(nil).Once()
	}

	fixture.runScan()
}

func TestPegInWatcher_DoesNotReimportPersistedEntry(t *testing.T) {
	fixture := newPegInWatcherFixture(t, 100, 1, 2)
	event := blockchain.AddressRegistered{
		TxHash: "duplicate", RskAddress: rskAddress(1), BlockNumber: 100, LogIndex: 1,
	}
	persisted := rootstock.PegInWatch{
		TxHash: "duplicate", LogIndex: 1, RskAddress: event.RskAddress, State: rootstock.PegInWatchImported,
	}

	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{persisted}, nil).Once()
	fixture.expectCatchUp(catchUpExpectation{
		checkpoint: rootstock.PegInWatchCheckpoint{Block: 99},
		head:       102,
		fromBlock:  100,
		toBlock:    100,
		events:     []blockchain.AddressRegistered{event},
	})
	fixture.repository.EXPECT().Get(mock.Anything, event.RskAddress).Return(&persisted, nil).Once()

	fixture.runScan()
}

func TestPegInWatcher_RecordsEntryErrorAndContinues(t *testing.T) {
	fixture := newPegInWatcherFixture(t, 100, 1, 2)
	events := []blockchain.AddressRegistered{
		{TxHash: "broken", RskAddress: rskAddress(1), BlockNumber: 100, LogIndex: 1},
		{TxHash: "valid", RskAddress: rskAddress(2), BlockNumber: 100, LogIndex: 2},
	}
	fixture.repository.EXPECT().List(mock.Anything).Return(nil, nil).Once()
	fixture.expectCatchUp(catchUpExpectation{checkpoint: rootstock.PegInWatchCheckpoint{Block: 99}, head: 102, fromBlock: 100, toBlock: 100, events: events})
	broken, valid := events[0], events[1]

	brokenEntry := discoveredWatchEntry(broken)
	fixture.repository.EXPECT().Get(mock.Anything, broken.RskAddress).
		Return(&brokenEntry, nil).
		Once()
	fixture.registry.EXPECT().GetPegInAddress(broken.RskAddress).
		Return(blockchain.PegInAddress{}, assert.AnError).
		Once()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(entry rootstock.PegInWatch) bool {
		return entry.TxHash == broken.TxHash &&
			entry.State == rootstock.PegInWatchDiscovered &&
			entry.LastError != ""
	})).Return(nil).Once()

	validEntry := discoveredWatchEntry(valid)
	fixture.repository.EXPECT().Get(mock.Anything, valid.RskAddress).
		Return(&validEntry, nil).
		Once()
	validDeposit := depositAddressFixture(0)
	fixture.registry.EXPECT().GetPegInAddress(valid.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  validDeposit.payload,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.wallet.EXPECT().ImportAddress(validDeposit.address).Return(nil).Once()
	fixture.expectBoundedRescan()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(entry rootstock.PegInWatch) bool {
		return entry.TxHash == valid.TxHash &&
			entry.State == rootstock.PegInWatchImported &&
			entry.LastError == ""
	})).Return(nil).Once()

	fixture.runScan()
}

func TestPegInWatcher_RetriesPersistedDiscoveredEntry(t *testing.T) {
	fixture := newPegInWatcherFixture(t, 100, 1, 2)
	entry := rootstock.PegInWatch{
		TxHash:      "retry-registration",
		LogIndex:    1,
		BlockNumber: 90,
		RskAddress:  "retry-rsk",
		State:       rootstock.PegInWatchDiscovered,
		LastError:   "previous transient failure",
	}
	fixture.repository.EXPECT().List(mock.Anything).
		Return([]rootstock.PegInWatch{entry}, nil).
		Once()
	fixture.expectCatchUp(catchUpExpectation{
		checkpoint: rootstock.PegInWatchCheckpoint{Block: 105},
		head:       108,
		fromBlock:  106,
		toBlock:    106,
	})
	fixture.repository.EXPECT().Get(mock.Anything, entry.RskAddress).Return(&entry, nil).Once()
	deposit := depositAddressFixture(0)
	fixture.registry.EXPECT().GetPegInAddress(entry.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  deposit.payload,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.wallet.EXPECT().ImportAddress(deposit.address).Return(nil).Once()
	fixture.expectBoundedRescan()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(updated rootstock.PegInWatch) bool {
		return updated.TxHash == entry.TxHash &&
			updated.BtcAddress == deposit.address &&
			updated.State == rootstock.PegInWatchImported &&
			updated.LastError == ""
	})).Return(nil).Once()

	fixture.runScan()
}

// An entry that keeps failing the same way every tick must not rewrite the same error to Mongo,
// which only holds if the persisted error survives the successful steps that precede the failure.
func TestPegInWatcher_SuppressesRepeatedIdenticalEntryErrors(t *testing.T) {
	fixture := newPegInWatcherFixture(t, 100, 1, 2)
	entry := rootstock.PegInWatch{
		TxHash:      "stuck",
		LogIndex:    1,
		BlockNumber: 90,
		RskAddress:  "stuck-rsk",
		State:       rootstock.PegInWatchDiscovered,
		LastError:   fmt.Sprintf("import PegIn address for event stuck/1: %v", assert.AnError),
	}
	fixture.repository.EXPECT().List(mock.Anything).
		Return([]rootstock.PegInWatch{entry}, nil).
		Once()
	fixture.expectCatchUp(catchUpExpectation{
		checkpoint: rootstock.PegInWatchCheckpoint{Block: 105},
		head:       108,
		fromBlock:  106,
		toBlock:    106,
	})
	fixture.repository.EXPECT().Get(mock.Anything, entry.RskAddress).Return(&entry, nil).Once()
	deposit := depositAddressFixture(0)
	fixture.registry.EXPECT().GetPegInAddress(entry.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  deposit.payload,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.wallet.EXPECT().ImportAddress(deposit.address).Return(assert.AnError).Once()

	fixture.runScan()

	fixture.repository.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestPegInWatcher_LeavesDiscoveredWhenRescanFails(t *testing.T) {
	fixture := newPegInWatcherFixture(t, 100, 1, 2)
	events := []blockchain.AddressRegistered{
		{TxHash: "needs-rescan", RskAddress: rskAddress(1), BlockNumber: 100, LogIndex: 1},
	}
	fixture.repository.EXPECT().List(mock.Anything).Return(nil, nil).Once()
	fixture.expectCatchUp(catchUpExpectation{checkpoint: rootstock.PegInWatchCheckpoint{Block: 99}, head: 102, fromBlock: 100, toBlock: 100, events: events})
	event := events[0]
	persisted := discoveredWatchEntry(event)
	fixture.repository.EXPECT().Get(mock.Anything, event.RskAddress).Return(&persisted, nil).Once()
	deposit := depositAddressFixture(0)
	fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  deposit.payload,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.wallet.EXPECT().ImportAddress(deposit.address).Return(nil).Once()
	fixture.btcNetwork.EXPECT().GetHeight().Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().RescanBlockchain(int64(100)).Return(blockchain.BitcoinRescanResult{}, assert.AnError).Once()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(entry rootstock.PegInWatch) bool {
		return entry.TxHash == event.TxHash &&
			entry.State == rootstock.PegInWatchDiscovered &&
			entry.LastError != ""
	})).Return(nil).Once()

	fixture.runScan()
}

// The registry returns a payload rather than an address, so a payload that cannot be encoded must
// park the entry with the reason instead of importing something a node would reject.
func TestPegInWatcher_RecordsUnencodablePayload(t *testing.T) {
	fixture := newPegInWatcherFixture(t, 100, 1, 2)
	entry := rootstock.PegInWatch{
		TxHash:      "truncated-payload",
		LogIndex:    1,
		BlockNumber: 90,
		RskAddress:  "truncated-rsk",
		State:       rootstock.PegInWatchDiscovered,
	}
	fixture.repository.EXPECT().List(mock.Anything).
		Return([]rootstock.PegInWatch{entry}, nil).
		Once()
	fixture.expectCatchUp(catchUpExpectation{
		checkpoint: rootstock.PegInWatchCheckpoint{Block: 105},
		head:       108,
		fromBlock:  106,
		toBlock:    106,
	})
	fixture.repository.EXPECT().Get(mock.Anything, entry.RskAddress).Return(&entry, nil).Once()
	truncated := depositAddressFixture(0).payload[:20]
	fixture.registry.EXPECT().GetPegInAddress(entry.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  truncated,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(updated rootstock.PegInWatch) bool {
		return updated.TxHash == entry.TxHash &&
			updated.State == rootstock.PegInWatchDiscovered &&
			updated.BtcAddress == "" &&
			strings.Contains(updated.LastError, "encode PegIn address for event truncated-payload/1")
	})).Return(nil).Once()

	fixture.runScan()

	fixture.wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
	fixture.wallet.AssertNotCalled(t, "RescanBlockchain", mock.Anything)
}

func TestPegInWatcher_RetriesDiscoveredEntriesWhenCatchUpFails(t *testing.T) {
	fixture := newPegInWatcherFixture(t, 100, 1, 2)
	entry := rootstock.PegInWatch{
		TxHash:      "pending-import",
		LogIndex:    1,
		BlockNumber: 90,
		RskAddress:  rskAddress(7),
		State:       rootstock.PegInWatchDiscovered,
	}
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{entry}, nil).Once()
	fixture.repository.EXPECT().GetCheckpoint(mock.Anything).
		Return(rootstock.PegInWatchCheckpoint{Root: [32]byte{1}, Block: 105}, true, nil).Once()
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(0), assert.AnError).Once()
	fixture.repository.EXPECT().Get(mock.Anything, entry.RskAddress).Return(&entry, nil).Once()
	deposit := depositAddressFixture(0)
	fixture.registry.EXPECT().GetPegInAddress(entry.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  deposit.payload,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.wallet.EXPECT().ImportAddress(deposit.address).Return(nil).Once()
	fixture.expectBoundedRescan()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(updated rootstock.PegInWatch) bool {
		return updated.TxHash == entry.TxHash && updated.State == rootstock.PegInWatchImported
	})).Return(nil).Once()

	fixture.runScan()

	fixture.repository.AssertNotCalled(t, "SaveVerifiedRegistrations", mock.Anything, mock.Anything, mock.Anything)
}

func TestPegInWatcher_ImportsTheSavedPagesWhenALaterPageFails(t *testing.T) {
	fixture := newPegInWatcherFixture(t, 100, 1, 2)
	stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{1}, Block: 105}
	events := []blockchain.AddressRegistered{{TxHash: "saved", BlockNumber: 106, LogIndex: 1, RskAddress: rskAddress(7)}}
	root := foldRegistrations(t, stored.Root, events)
	saved := discoveredWatchEntry(events[0])
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{}, nil).Once()
	fixture.repository.EXPECT().GetCheckpoint(mock.Anything).Return(stored, true, nil).Once()
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(109), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(106), uint64Pointer(106)).Return(events, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(106)).Return(root, nil).Once()
	fixture.repository.EXPECT().SaveVerifiedRegistrations(mock.Anything, mock.Anything, rootstock.PegInWatchCheckpoint{Root: root, Block: 106}).
		Return(nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(107), uint64Pointer(107)).Return(nil, assert.AnError).Once()
	fixture.repository.EXPECT().Get(mock.Anything, saved.RskAddress).Return(&saved, nil).Once()
	deposit := depositAddressFixture(0)
	fixture.registry.EXPECT().GetPegInAddress(saved.RskAddress).
		Return(blockchain.PegInAddress{Payload: deposit.payload, Encoding: blockchain.PegInAddressRegistryEncodingBase58}, nil).Once()
	fixture.wallet.EXPECT().ImportAddress(deposit.address).Return(nil).Once()
	fixture.expectBoundedRescan()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(updated rootstock.PegInWatch) bool {
		return updated.RskAddress == saved.RskAddress && updated.State == rootstock.PegInWatchImported
	})).Return(nil).Once()

	fixture.runScan()
}

func TestPegInWatcher_RetriesDiscoveredEntriesWhenRebuildFails(t *testing.T) {
	fixture := newPegInWatcherFixture(t, 100, 1, 2)
	entry := rootstock.PegInWatch{
		TxHash:      "pending-import",
		LogIndex:    1,
		BlockNumber: 90,
		RskAddress:  rskAddress(7),
		State:       rootstock.PegInWatchDiscovered,
	}
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{entry}, nil).Once()
	fixture.repository.EXPECT().GetCheckpoint(mock.Anything).
		Return(rootstock.PegInWatchCheckpoint{Root: [32]byte{1}, Block: 105}, true, nil).Once()
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(108), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(106), uint64Pointer(106)).
		Return([]blockchain.AddressRegistered{}, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(106)).Return([32]byte{2}, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(105)).Return([32]byte{}, assert.AnError).Once()
	fixture.repository.EXPECT().Get(mock.Anything, entry.RskAddress).Return(&entry, nil).Once()
	deposit := depositAddressFixture(0)
	fixture.registry.EXPECT().GetPegInAddress(entry.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  deposit.payload,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.wallet.EXPECT().ImportAddress(deposit.address).Return(nil).Once()
	fixture.expectBoundedRescan()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(updated rootstock.PegInWatch) bool {
		return updated.TxHash == entry.TxHash && updated.State == rootstock.PegInWatchImported
	})).Return(nil).Once()

	fixture.runScan()

	fixture.repository.AssertNotCalled(t, "SaveVerifiedRegistrations", mock.Anything, mock.Anything, mock.Anything)
}

//nolint:funlen // One scenario covers the rebuild, the reload and the import of the only new address.
func TestPegInWatcher_ReloadsTheWatchListAfterARebuild(t *testing.T) {
	fixture := newPegInWatcherFixture(t, 100, 1, 2)
	stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{1}, Block: 105}
	dropped := rootstock.PegInWatch{TxHash: "dropped", BlockNumber: 106, RskAddress: rskAddress(5), State: rootstock.PegInWatchDiscovered}
	kept := rootstock.PegInWatch{TxHash: "kept", BlockNumber: 100, RskAddress: rskAddress(6), State: rootstock.PegInWatchImported}
	rebuilt := []blockchain.AddressRegistered{
		{TxHash: "kept-again", BlockNumber: 106, LogIndex: 1, RskAddress: kept.RskAddress},
		{TxHash: "added", BlockNumber: 106, LogIndex: 2, RskAddress: rskAddress(7)},
	}
	root := foldRegistrations(t, stored.Root, rebuilt)
	added := discoveredWatchEntry(rebuilt[1])
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{kept, dropped}, nil).Once()
	fixture.repository.EXPECT().GetCheckpoint(mock.Anything).Return(stored, true, nil).Once()
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(108), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(106), uint64Pointer(106)).
		Return([]blockchain.AddressRegistered{}, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(106)).Return(root, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(105)).Return(stored.Root, nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(106), uint64Pointer(106)).
		Return(rebuilt, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(106)).Return(root, nil).Once()
	fixture.repository.EXPECT().ReplaceVerifiedRegistrations(mock.Anything, uint64(106), mock.Anything, rootstock.PegInWatchCheckpoint{Root: root, Block: 106}).
		Return(nil).Once()
	keptAfterRebuild := kept
	keptAfterRebuild.TxHash, keptAfterRebuild.BlockNumber = "kept-again", 106
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{keptAfterRebuild, added}, nil).Once()
	fixture.repository.EXPECT().Get(mock.Anything, added.RskAddress).Return(&added, nil).Once()
	deposit := depositAddressFixture(0)
	fixture.registry.EXPECT().GetPegInAddress(added.RskAddress).
		Return(blockchain.PegInAddress{Payload: deposit.payload, Encoding: blockchain.PegInAddressRegistryEncodingBase58}, nil).Once()
	fixture.wallet.EXPECT().ImportAddress(deposit.address).Return(nil).Once()
	fixture.expectBoundedRescan()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(updated rootstock.PegInWatch) bool {
		return updated.RskAddress == added.RskAddress && updated.State == rootstock.PegInWatchImported
	})).Return(nil).Once()

	fixture.runScan()

	fixture.repository.AssertNotCalled(t, "Get", mock.Anything, dropped.RskAddress)
	fixture.repository.AssertNotCalled(t, "Get", mock.Anything, kept.RskAddress)
}

//nolint:funlen // One scenario covers the failed reload after a rebuild and the reload on the next tick.
func TestPegInWatcher_ReloadsOnTheNextTickWhenTheReloadAfterARebuildFails(t *testing.T) {
	fixture := newPegInWatcherFixture(t, 100, 1, 2)
	stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{1}, Block: 105}
	dropped := rootstock.PegInWatch{TxHash: "dropped", BlockNumber: 106, RskAddress: rskAddress(5), State: rootstock.PegInWatchDiscovered}
	rebuiltEvents := []blockchain.AddressRegistered{{TxHash: "added", BlockNumber: 106, LogIndex: 1, RskAddress: rskAddress(7)}}
	root := foldRegistrations(t, stored.Root, rebuiltEvents)
	rebuiltPoint := rootstock.PegInWatchCheckpoint{Root: root, Block: 106}
	added := discoveredWatchEntry(rebuiltEvents[0])
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{dropped}, nil).Once()
	fixture.repository.EXPECT().GetCheckpoint(mock.Anything).Return(stored, true, nil).Once()
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(108), nil).Times(3)
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(106), uint64Pointer(106)).
		Return([]blockchain.AddressRegistered{}, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(106)).Return(root, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(105)).Return(stored.Root, nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(106), uint64Pointer(106)).
		Return(rebuiltEvents, nil).Once()
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(106)).Return(root, nil).Once()
	fixture.repository.EXPECT().ReplaceVerifiedRegistrations(mock.Anything, uint64(106), mock.Anything, rebuiltPoint).Return(nil).Once()
	fixture.repository.EXPECT().List(mock.Anything).Return(nil, assert.AnError).Once()
	fixture.repository.EXPECT().GetCheckpoint(mock.Anything).Return(rebuiltPoint, true, nil).Twice()
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{added}, nil).Once()
	fixture.repository.EXPECT().Get(mock.Anything, added.RskAddress).Return(&added, nil).Once()
	deposit := depositAddressFixture(0)
	fixture.registry.EXPECT().GetPegInAddress(added.RskAddress).
		Return(blockchain.PegInAddress{Payload: deposit.payload, Encoding: blockchain.PegInAddressRegistryEncodingBase58}, nil).Once()
	fixture.wallet.EXPECT().ImportAddress(deposit.address).Return(nil).Once()
	fixture.expectBoundedRescan()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(updated rootstock.PegInWatch) bool {
		return updated.RskAddress == added.RskAddress && updated.State == rootstock.PegInWatchImported
	})).Return(nil).Once()

	session := startWatcherSession(t, fixture.watcher, fixture.ticker)
	session.poll()
	fixture.repository.AssertNotCalled(t, "Get", mock.Anything, dropped.RskAddress)
	fixture.wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
	session.poll()
	session.poll()
	session.stop()

	fixture.repository.AssertNotCalled(t, "Get", mock.Anything, dropped.RskAddress)
	fixture.repository.AssertNumberOfCalls(t, "List", 3)
}

func TestPegInWatcher_KeepsTheImportResultsWhenTheImportFails(t *testing.T) {
	fixture := newPegInWatcherFixture(t, 100, 1, 2)
	stale := rootstock.PegInWatch{TxHash: "stale", BlockNumber: 90, RskAddress: rskAddress(1), State: rootstock.PegInWatchDiscovered}
	unreadable := rootstock.PegInWatch{TxHash: "unreadable", BlockNumber: 91, RskAddress: rskAddress(2), State: rootstock.PegInWatchDiscovered}
	storedStale, storedUnreadable := stale, unreadable
	storedStale.State, storedUnreadable.State = rootstock.PegInWatchImported, rootstock.PegInWatchImported
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInWatch{stale, unreadable}, nil).Once()
	fixture.repository.EXPECT().GetCheckpoint(mock.Anything).
		Return(rootstock.PegInWatchCheckpoint{Root: [32]byte{1}, Block: 105}, true, nil).Twice()
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(107), nil).Twice()
	fixture.repository.EXPECT().Get(mock.Anything, stale.RskAddress).Return(&storedStale, nil).Once()
	fixture.repository.EXPECT().Get(mock.Anything, unreadable.RskAddress).Return(nil, assert.AnError).Once()
	fixture.repository.EXPECT().Get(mock.Anything, unreadable.RskAddress).Return(&storedUnreadable, nil).Once()

	session := startWatcherSession(t, fixture.watcher, fixture.ticker)
	session.poll()
	session.poll()
	session.stop()

	fixture.repository.AssertNumberOfCalls(t, "Get", 3)
	fixture.wallet.AssertNotCalled(t, "RescanBlockchain", mock.Anything)
}

func uint64Pointer(expected uint64) interface{} {
	return mock.MatchedBy(func(actual *uint64) bool {
		return actual != nil && *actual == expected
	})
}
