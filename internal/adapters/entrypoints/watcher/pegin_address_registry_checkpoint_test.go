package watcher_test

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"

	watcherAdapter "github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	watcherUseCase "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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
	assert.Equal(t, rootstock.PegInAddressRegistryWatchImported, entries[0].State)
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
		"entry-update",
		"checkpoint",
		"checkpoint",
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
		Once()
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(103), nil).Twice()
	wallet := mocks.NewBitcoinWalletMock(t)
	wallet.EXPECT().ImportAddress(deposit.address).Return(nil).Once()
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
	discover := watcherUseCase.NewDiscoverRegisteredAddressUseCase(
		scenario.store,
		scenario.registry,
		scenario.wallet,
	)
	getWatched := watcherUseCase.NewGetWatchedRegisteredAddressesUseCase(scenario.store)
	watcher := watcherAdapter.NewPegInAddressRegistryWatcher(
		watcherAdapter.NewPegInAddressRegistryWatcherUseCases(discover, getWatched),
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
