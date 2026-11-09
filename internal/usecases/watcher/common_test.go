package watcher_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	catchUpStartBlock    uint64 = 100
	catchUpPageSize      uint64 = 10
	catchUpFinalityDepth uint64 = 2
)

type catchUpFixture struct {
	repository *mocks.PegInWatchRepositoryMock
	registry   *mocks.PegInAddressRegistryContractMock
	rskRpc     *mocks.RootstockRpcServerMock
	hashCalls  int
	useCase    *watcher.CatchUpRegisteredAddressesUseCase
}

func newCatchUpFixture(t *testing.T) *catchUpFixture {
	t.Helper()
	return newCatchUpFixtureWithPageSize(t, catchUpPageSize)
}

func newCatchUpFixtureWithPageSize(t *testing.T, pageSize uint64) *catchUpFixture {
	t.Helper()
	fixture := &catchUpFixture{
		repository: mocks.NewPegInWatchRepositoryMock(t),
		registry:   mocks.NewPegInAddressRegistryContractMock(t),
		rskRpc:     mocks.NewRootstockRpcServerMock(t),
	}
	countingHash := func(data ...[]byte) []byte {
		fixture.hashCalls++
		return crypto.Keccak256(data...)
	}
	fixture.useCase = watcher.NewCatchUpRegisteredAddressesUseCase(
		fixture.repository,
		fixture.registry,
		fixture.rskRpc,
		countingHash,
		catchUpStartBlock,
		pageSize,
		catchUpFinalityDepth,
	)
	return fixture
}

func (fixture *catchUpFixture) withCheckpoint(checkpoint *rootstock.PegInWatchCheckpoint) {
	if checkpoint == nil {
		fixture.repository.EXPECT().GetCheckpoint(mock.Anything).Return(rootstock.PegInWatchCheckpoint{}, false, nil).Once()
		return
	}
	fixture.repository.EXPECT().GetCheckpoint(mock.Anything).Return(*checkpoint, true, nil).Once()
}

func (fixture *catchUpFixture) withPage(head, fromBlock, toBlock uint64, events []blockchain.AddressRegistered) {
	fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(head, nil).Once()
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, fromBlock, blockPointer(toBlock)).Return(events, nil).Once()
}

func catchUpAddress(suffix int) string {
	return fmt.Sprintf("0x%040x", suffix)
}

func chainedRegistrations(t *testing.T, previousRoot [32]byte, events ...blockchain.AddressRegistered) []blockchain.AddressRegistered {
	t.Helper()
	root := previousRoot
	for index := range events {
		var err error
		root, err = blockchain.FoldPegInAddressRegistryRoot(crypto.Keccak256, root, events[index].RskAddress)
		require.NoError(t, err)
		events[index].RegistrationRoot = root
	}
	return events
}

func (fixture *catchUpFixture) withContractRoot(block uint64, root [32]byte) {
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, block).Return(root, nil).Once()
}

func blockPointer(expected uint64) interface{} {
	return mock.MatchedBy(func(actual *uint64) bool {
		return actual != nil && *actual == expected
	})
}

func savedTxHashes(watches []rootstock.PegInWatch) []string {
	hashes := make([]string, 0, len(watches))
	for _, watch := range watches {
		hashes = append(hashes, watch.TxHash)
	}
	return hashes
}

func (fixture *catchUpFixture) withEvents(fromBlock, toBlock uint64, events []blockchain.AddressRegistered) {
	fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, fromBlock, blockPointer(toBlock)).Return(events, nil).Once()
}

type expectedCatchUpResult struct {
	txHashes   []string
	checkpoint rootstock.PegInWatchCheckpoint
	rebuilt    bool
}

func assertCatchUpResult(t *testing.T, expected expectedCatchUpResult, result watcher.CatchUpRegisteredAddressesResult) {
	t.Helper()
	txHashes := expected.txHashes
	if txHashes == nil {
		txHashes = []string{}
	}
	assert.Equal(t, txHashes, savedTxHashes(result.Watches))
	assert.Equal(t, expected.checkpoint, result.Checkpoint)
	assert.Equal(t, expected.rebuilt, result.Rebuilt)
}
