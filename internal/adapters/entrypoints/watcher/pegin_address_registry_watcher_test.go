package watcher_test

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
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
	repository *mocks.PegInAddressRegistryWatchRepositoryMock
	registry   *mocks.PegInAddressRegistryContractMock
	rskRpc     *mocks.RootstockRpcServerMock
	btcNetwork *mocks.BtcRpcMock
	wallet     *mocks.BitcoinWalletMock
	ticker     *sessionTicker
	watcher    *watcherAdapter.PegInAddressRegistryWatcher
}

//nolint:unparam // Fixed start/finality values make each scenario's cursor arithmetic explicit at the call site.
func newPegInAddressRegistryWatcherFixture(
	t *testing.T,
	startBlock uint64,
	pageSize uint64,
	finalityDepth uint64,
) *pegInAddressRegistryWatcherFixture {
	t.Helper()
	fixture := &pegInAddressRegistryWatcherFixture{
		t:          t,
		repository: mocks.NewPegInAddressRegistryWatchRepositoryMock(t),
		registry:   mocks.NewPegInAddressRegistryContractMock(t),
		rskRpc:     mocks.NewRootstockRpcServerMock(t),
		btcNetwork: &mocks.BtcRpcMock{},
		wallet:     mocks.NewBitcoinWalletMock(t),
		ticker:     newSessionTicker(),
	}
	t.Cleanup(func() { fixture.btcNetwork.AssertExpectations(t) })
	useCases := watcherAdapter.NewPegInAddressRegistryWatcherUseCases(
		watcher.NewGetWatchedRegisteredAddressesUseCase(fixture.repository),
		watcher.NewGetRegistryWatchCursorUseCase(fixture.repository),
		watcher.NewSetRegistryWatchCursorUseCase(fixture.repository),
		watcher.NewDiscoverRegisteredAddressUseCase(fixture.repository, fixture.registry, fixture.wallet),
		watcher.NewMarkRegisteredAddressImportedUseCase(fixture.repository),
		watcher.NewRecordRegisteredAddressWatchErrorUseCase(fixture.repository),
	)
	fixture.watcher = watcherAdapter.NewPegInAddressRegistryWatcher(
		useCases,
		fixture.registry,
		fixture.rskRpc,
		fixture.btcNetwork,
		fixture.wallet,
		fixture.ticker,
		startBlock,
		pageSize,
		finalityDepth,
	)
	return fixture
}

func (fixture *pegInAddressRegistryWatcherFixture) expectBoundedRescan() {
	fixture.t.Helper()
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().RescanBlockchain(int64(100)).
		Return(blockchain.BitcoinRescanResult{StartHeight: 100, StopHeight: 200}, nil).
		Once()
}

func (fixture *pegInAddressRegistryWatcherFixture) expectCursorAdvance(toBlock uint64) {
	fixture.t.Helper()
	fixture.repository.On("SetCursor", mock.Anything, toBlock).Return(nil).Once()
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

func discoveredWatchEntry(event blockchain.AddressRegistered) rootstock.PegInAddressRegistryWatch {
	return rootstock.PegInAddressRegistryWatch{
		TxHash:      event.TxHash,
		LogIndex:    event.LogIndex,
		BlockNumber: event.BlockNumber,
		RskAddress:  event.RskAddress,
		State:       rootstock.PegInAddressRegistryWatchDiscovered,
	}
}

func TestPegInAddressRegistryWatcher_PrepareResumesWithOverlap(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 3, 2)
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInAddressRegistryWatch{{
		TxHash: "persisted",
		State:  rootstock.PegInAddressRegistryWatchImported,
	}}, nil).Once()
	fixture.repository.EXPECT().GetCursor(mock.Anything).Return(uint64(105), true, nil).Once()
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(112), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(104), uint64Pointer(108)).
		Return([]blockchain.AddressRegistered{}, nil).Once()
	fixture.expectCursorAdvance(108)

	fixture.runScan()
}

//nolint:funlen // One scenario records the complete persist/import/cursor ordering across three events.
func TestPegInAddressRegistryWatcher_SortsEventsAndClampsToFinalizedHead(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 20, 2)
	events := []blockchain.AddressRegistered{
		{TxHash: "third", RskAddress: "third-rsk", BlockNumber: 102, LogIndex: 5},
		{TxHash: "first", RskAddress: "first-rsk", BlockNumber: 101, LogIndex: 9},
		{TxHash: "second", RskAddress: "second-rsk", BlockNumber: 102, LogIndex: 1},
	}
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInAddressRegistryWatch{}, nil).Once()
	fixture.repository.EXPECT().GetCursor(mock.Anything).Return(uint64(0), false, nil).Once()
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(105), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(100), uint64Pointer(103)).
		Return(events, nil).Once()

	var operations []string
	var operationsMutex sync.Mutex
	fixture.repository.On("Upsert", mock.Anything, mock.Anything).
		Run(func(arguments mock.Arguments) {
			entry, ok := arguments.Get(1).(rootstock.PegInAddressRegistryWatch)
			require.True(t, ok)
			operationsMutex.Lock()
			operations = append(operations, fmt.Sprintf("upsert:%d/%d/%s", entry.BlockNumber, entry.LogIndex, entry.TxHash))
			operationsMutex.Unlock()
		}).
		Return(nil).
		Times(3)
	for index, event := range events {
		deposit := depositAddressFixture(index)
		persisted := discoveredWatchEntry(event)
		fixture.repository.EXPECT().Get(mock.Anything, event.TxHash, event.LogIndex).Return(nil, nil).Once()
		fixture.repository.EXPECT().Get(mock.Anything, event.TxHash, event.LogIndex).Return(&persisted, nil).Once()
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
		fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(entry rootstock.PegInAddressRegistryWatch) bool {
			return entry.TxHash == event.TxHash &&
				entry.BtcAddress == deposit.address &&
				entry.DepositTxID == "" &&
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
	fixture.repository.On("SetCursor", mock.Anything, uint64(103)).
		Run(func(mock.Arguments) {
			operationsMutex.Lock()
			operations = append(operations, "cursor")
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
		"cursor",
	}, operations)
	operationsMutex.Unlock()
}

func TestPegInAddressRegistryWatcher_SkipsUnsupportedEncodingsAndContinues(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 1, 2)
	events := []blockchain.AddressRegistered{
		{TxHash: "bech32", RskAddress: "bech32-rsk", BlockNumber: 100, LogIndex: 1},
		{TxHash: "bech32m", RskAddress: "bech32m-rsk", BlockNumber: 100, LogIndex: 2},
		{TxHash: "base58", RskAddress: "base58-rsk", BlockNumber: 100, LogIndex: 3},
	}
	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInAddressRegistryWatch{}, nil).Once()
	fixture.repository.EXPECT().GetCursor(mock.Anything).Return(uint64(0), false, nil).Once()
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(102), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(100), uint64Pointer(100)).
		Return(events, nil).Once()
	fixture.repository.EXPECT().Upsert(mock.Anything, mock.Anything).Return(nil).Times(3)

	encodings := []blockchain.PegInAddressRegistryEncoding{
		blockchain.PegInAddressRegistryEncodingBech32,
		blockchain.PegInAddressRegistryEncodingBech32M,
		blockchain.PegInAddressRegistryEncodingBase58,
	}
	for index, event := range events {
		deposit := depositAddressFixture(index)
		persisted := discoveredWatchEntry(event)
		fixture.repository.EXPECT().Get(mock.Anything, event.TxHash, event.LogIndex).Return(nil, nil).Once()
		fixture.repository.EXPECT().Get(mock.Anything, event.TxHash, event.LogIndex).Return(&persisted, nil).Once()
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
		fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(entry rootstock.PegInAddressRegistryWatch) bool {
			return entry.TxHash == event.TxHash &&
				entry.State == expectedState &&
				entry.BtcAddress == expectedAddress &&
				(expectedState != rootstock.PegInAddressRegistryWatchImported ||
					entry.DepositTxID == "")
		})).Return(nil).Once()
	}
	fixture.expectCursorAdvance(100)

	fixture.runScan()
}

func TestPegInAddressRegistryWatcher_DoesNotReimportPersistedEntry(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 1, 2)
	event := blockchain.AddressRegistered{
		TxHash: "duplicate", RskAddress: "duplicate-rsk", BlockNumber: 100, LogIndex: 1,
	}
	persisted := rootstock.PegInAddressRegistryWatch{
		TxHash: "duplicate", LogIndex: 1, State: rootstock.PegInAddressRegistryWatchImported,
	}

	fixture.repository.EXPECT().List(mock.Anything).Return([]rootstock.PegInAddressRegistryWatch{persisted}, nil).Once()
	fixture.repository.EXPECT().GetCursor(mock.Anything).Return(uint64(99), true, nil).Once()
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(102), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(100), uint64Pointer(100)).
		Return([]blockchain.AddressRegistered{event}, nil).Once()
	fixture.repository.EXPECT().Get(mock.Anything, event.TxHash, event.LogIndex).Return(&persisted, nil).Once()
	fixture.expectCursorAdvance(100)

	fixture.runScan()
}

func TestPegInAddressRegistryWatcher_RecordsEntryErrorAndContinues(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 1, 2)
	broken := blockchain.AddressRegistered{
		TxHash: "broken", RskAddress: "broken-rsk", BlockNumber: 100, LogIndex: 1,
	}
	valid := blockchain.AddressRegistered{
		TxHash: "valid", RskAddress: "valid-rsk", BlockNumber: 100, LogIndex: 2,
	}
	fixture.repository.EXPECT().List(mock.Anything).Return(nil, nil).Once()
	fixture.repository.EXPECT().GetCursor(mock.Anything).Return(uint64(0), false, nil).Once()
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(102), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(100), uint64Pointer(100)).
		Return([]blockchain.AddressRegistered{broken, valid}, nil).
		Once()
	fixture.repository.EXPECT().Upsert(mock.Anything, mock.Anything).Return(nil).Twice()

	brokenEntry := discoveredWatchEntry(broken)
	fixture.repository.EXPECT().Get(mock.Anything, broken.TxHash, broken.LogIndex).
		Return(nil, nil).
		Once()
	fixture.repository.EXPECT().Get(mock.Anything, broken.TxHash, broken.LogIndex).
		Return(&brokenEntry, nil).
		Once()
	fixture.registry.EXPECT().GetPegInAddress(broken.RskAddress).
		Return(blockchain.PegInAddress{}, assert.AnError).
		Once()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(entry rootstock.PegInAddressRegistryWatch) bool {
		return entry.TxHash == broken.TxHash &&
			entry.State == rootstock.PegInAddressRegistryWatchDiscovered &&
			entry.LastError != ""
	})).Return(nil).Once()

	validEntry := discoveredWatchEntry(valid)
	fixture.repository.EXPECT().Get(mock.Anything, valid.TxHash, valid.LogIndex).
		Return(nil, nil).
		Once()
	fixture.repository.EXPECT().Get(mock.Anything, valid.TxHash, valid.LogIndex).
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
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(entry rootstock.PegInAddressRegistryWatch) bool {
		return entry.TxHash == valid.TxHash &&
			entry.State == rootstock.PegInAddressRegistryWatchImported &&
			entry.DepositTxID == "" &&
			entry.LastError == ""
	})).Return(nil).Once()
	fixture.expectCursorAdvance(100)

	fixture.runScan()
}

func TestPegInAddressRegistryWatcher_RetriesPersistedDiscoveredEntryOutsideOverlap(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 1, 2)
	entry := rootstock.PegInAddressRegistryWatch{
		TxHash:      "retry-registration",
		LogIndex:    1,
		BlockNumber: 90,
		RskAddress:  "retry-rsk",
		State:       rootstock.PegInAddressRegistryWatchDiscovered,
		LastError:   "previous transient failure",
	}
	fixture.repository.EXPECT().List(mock.Anything).
		Return([]rootstock.PegInAddressRegistryWatch{entry}, nil).
		Once()
	fixture.repository.EXPECT().GetCursor(mock.Anything).Return(uint64(105), true, nil).Once()
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(108), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(104), uint64Pointer(106)).
		Return(nil, nil).
		Once()
	fixture.repository.EXPECT().Get(mock.Anything, entry.TxHash, entry.LogIndex).Return(&entry, nil).Once()
	deposit := depositAddressFixture(0)
	fixture.registry.EXPECT().GetPegInAddress(entry.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  deposit.payload,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.wallet.EXPECT().ImportAddress(deposit.address).Return(nil).Once()
	fixture.expectBoundedRescan()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(updated rootstock.PegInAddressRegistryWatch) bool {
		return updated.TxHash == entry.TxHash &&
			updated.BtcAddress == deposit.address &&
			updated.DepositTxID == "" &&
			updated.State == rootstock.PegInAddressRegistryWatchImported &&
			updated.LastError == ""
	})).Return(nil).Once()
	fixture.expectCursorAdvance(106)

	fixture.runScan()
}

// An entry that keeps failing the same way every tick must not rewrite the same error to Mongo,
// which only holds if the persisted error survives the successful steps that precede the failure.
func TestPegInAddressRegistryWatcher_SuppressesRepeatedIdenticalEntryErrors(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 1, 2)
	entry := rootstock.PegInAddressRegistryWatch{
		TxHash:      "stuck",
		LogIndex:    1,
		BlockNumber: 90,
		RskAddress:  "stuck-rsk",
		State:       rootstock.PegInAddressRegistryWatchDiscovered,
		LastError:   fmt.Sprintf("import PegIn address for event stuck/1: %v", assert.AnError),
	}
	fixture.repository.EXPECT().List(mock.Anything).
		Return([]rootstock.PegInAddressRegistryWatch{entry}, nil).
		Once()
	fixture.repository.EXPECT().GetCursor(mock.Anything).Return(uint64(105), true, nil).Once()
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(108), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(104), uint64Pointer(106)).
		Return(nil, nil).
		Once()
	fixture.repository.EXPECT().Get(mock.Anything, entry.TxHash, entry.LogIndex).Return(&entry, nil).Once()
	deposit := depositAddressFixture(0)
	fixture.registry.EXPECT().GetPegInAddress(entry.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  deposit.payload,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.wallet.EXPECT().ImportAddress(deposit.address).Return(assert.AnError).Once()
	fixture.expectCursorAdvance(106)

	fixture.runScan()

	fixture.repository.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestPegInAddressRegistryWatcher_LeavesDiscoveredWhenRescanFails(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 1, 2)
	event := blockchain.AddressRegistered{
		TxHash: "needs-rescan", RskAddress: "needs-rescan-rsk", BlockNumber: 100, LogIndex: 1,
	}
	fixture.repository.EXPECT().List(mock.Anything).Return(nil, nil).Once()
	fixture.repository.EXPECT().GetCursor(mock.Anything).Return(uint64(0), false, nil).Once()
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(102), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(100), uint64Pointer(100)).
		Return([]blockchain.AddressRegistered{event}, nil).Once()
	fixture.repository.EXPECT().Upsert(mock.Anything, mock.Anything).Return(nil).Once()
	persisted := discoveredWatchEntry(event)
	fixture.repository.EXPECT().Get(mock.Anything, event.TxHash, event.LogIndex).Return(nil, nil).Once()
	fixture.repository.EXPECT().Get(mock.Anything, event.TxHash, event.LogIndex).Return(&persisted, nil).Once()
	deposit := depositAddressFixture(0)
	fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  deposit.payload,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.wallet.EXPECT().ImportAddress(deposit.address).Return(nil).Once()
	fixture.btcNetwork.On("GetHeight").Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().RescanBlockchain(int64(100)).Return(blockchain.BitcoinRescanResult{}, assert.AnError).Once()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(entry rootstock.PegInAddressRegistryWatch) bool {
		return entry.TxHash == event.TxHash &&
			entry.State == rootstock.PegInAddressRegistryWatchDiscovered &&
			entry.LastError != ""
	})).Return(nil).Once()
	fixture.expectCursorAdvance(100)

	fixture.runScan()
}

// The registry returns a payload rather than an address, so a payload that cannot be encoded must
// park the entry with the reason instead of importing something a node would reject.
func TestPegInAddressRegistryWatcher_RecordsUnencodablePayload(t *testing.T) {
	fixture := newPegInAddressRegistryWatcherFixture(t, 100, 1, 2)
	entry := rootstock.PegInAddressRegistryWatch{
		TxHash:      "truncated-payload",
		LogIndex:    1,
		BlockNumber: 90,
		RskAddress:  "truncated-rsk",
		State:       rootstock.PegInAddressRegistryWatchDiscovered,
	}
	fixture.repository.EXPECT().List(mock.Anything).
		Return([]rootstock.PegInAddressRegistryWatch{entry}, nil).
		Once()
	fixture.repository.EXPECT().GetCursor(mock.Anything).Return(uint64(105), true, nil).Once()
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(108), nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(104), uint64Pointer(106)).
		Return(nil, nil).
		Once()
	fixture.repository.EXPECT().Get(mock.Anything, entry.TxHash, entry.LogIndex).Return(&entry, nil).Once()
	truncated := depositAddressFixture(0).payload[:20]
	fixture.registry.EXPECT().GetPegInAddress(entry.RskAddress).
		Return(blockchain.PegInAddress{
			Payload:  truncated,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, nil).
		Once()
	fixture.repository.On("Update", mock.Anything, mock.MatchedBy(func(updated rootstock.PegInAddressRegistryWatch) bool {
		return updated.TxHash == entry.TxHash &&
			updated.State == rootstock.PegInAddressRegistryWatchDiscovered &&
			updated.BtcAddress == "" &&
			updated.DepositTxID == "" &&
			strings.Contains(updated.LastError, "encode PegIn address for event truncated-payload/1")
	})).Return(nil).Once()
	fixture.expectCursorAdvance(106)

	fixture.runScan()

	fixture.wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
	fixture.wallet.AssertNotCalled(t, "RescanBlockchain", mock.Anything)
}

func uint64Pointer(expected uint64) interface{} {
	return mock.MatchedBy(func(actual *uint64) bool {
		return actual != nil && *actual == expected
	})
}
