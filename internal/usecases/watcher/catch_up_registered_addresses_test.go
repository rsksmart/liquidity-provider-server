package watcher_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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

//nolint:funlen
func TestCatchUpRegisteredAddressesUseCase_Run(t *testing.T) {
	t.Run("starts from the start block when no point is stored", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(nil)
		events := chainedRegistrations(t, [32]byte{},
			blockchain.AddressRegistered{TxHash: "0x01", RskAddress: catchUpAddress(1), Registrant: "0xr", BlockNumber: 101},
			blockchain.AddressRegistered{TxHash: "0x02", RskAddress: catchUpAddress(2), Registrant: "0xr", BlockNumber: 104, LogIndex: 3},
		)
		root := events[1].RegistrationRoot
		fixture.withPage(200, 100, 109, events)
		fixture.withContractRoot(109, root)
		fixture.repository.EXPECT().SaveVerifiedRegistrations(
			mock.Anything,
			mock.MatchedBy(func(watches []rootstock.PegInWatch) bool {
				return len(watches) == 2 &&
					watches[0].TxHash == "0x01" && watches[0].RskAddress == catchUpAddress(1) && watches[0].Registrant == "0xr" &&
					watches[0].State == rootstock.PegInWatchDiscovered &&
					watches[1].BlockNumber == 104 && watches[1].LogIndex == 3 &&
					watches[1].RegistrationRoot == root
			}),
			rootstock.PegInWatchCheckpoint{Root: root, Block: 109},
		).Return(nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assert.Equal(t, rootstock.PegInWatchCheckpoint{Root: root, Block: 109}, result.Checkpoint)
		assert.Equal(t, []string{"0x01", "0x02"}, savedTxHashes(result.Watches))
	})
	t.Run("continues after the stored block from the stored root", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{9}, Block: 150}
		fixture.withCheckpoint(&stored)
		events := chainedRegistrations(t, stored.Root,
			blockchain.AddressRegistered{TxHash: "0x03", RskAddress: catchUpAddress(3), BlockNumber: 152, LogIndex: 1},
		)
		fixture.withPage(155, 151, 153, events)
		fixture.withContractRoot(153, events[0].RegistrationRoot)
		fixture.repository.EXPECT().SaveVerifiedRegistrations(
			mock.Anything,
			mock.Anything,
			rootstock.PegInWatchCheckpoint{Root: events[0].RegistrationRoot, Block: 153},
		).Return(nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assert.Equal(t, uint64(153), result.Checkpoint.Block)
		assert.Equal(t, []string{"0x03"}, savedTxHashes(result.Watches))
	})
	t.Run("reads one page at a time", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(&rootstock.PegInWatchCheckpoint{Root: [32]byte{9}, Block: 120})
		fixture.withPage(200, 121, 130, []blockchain.AddressRegistered{})
		fixture.withContractRoot(130, [32]byte{9})
		fixture.repository.EXPECT().SaveVerifiedRegistrations(mock.Anything, mock.Anything, rootstock.PegInWatchCheckpoint{Root: [32]byte{9}, Block: 130}).
			Return(nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assert.Equal(t, uint64(130), result.Checkpoint.Block)
	})
	t.Run("keeps the stored root when the page has no registrations", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{9}, Block: 150}
		fixture.withCheckpoint(&stored)
		fixture.withPage(160, 151, 158, []blockchain.AddressRegistered{})
		fixture.withContractRoot(158, stored.Root)
		fixture.repository.EXPECT().SaveVerifiedRegistrations(
			mock.Anything,
			[]rootstock.PegInWatch{},
			rootstock.PegInWatchCheckpoint{Root: stored.Root, Block: 158},
		).Return(nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assert.Empty(t, result.Watches)
		assert.Equal(t, stored.Root, result.Checkpoint.Root)
	})
	t.Run("folds the registrations in chain order", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(nil)
		chain := chainedRegistrations(t, [32]byte{},
			blockchain.AddressRegistered{TxHash: "0xfirst", RskAddress: catchUpAddress(1), BlockNumber: 103, LogIndex: 9},
			blockchain.AddressRegistered{TxHash: "0xmiddle", RskAddress: catchUpAddress(2), BlockNumber: 105, LogIndex: 1},
			blockchain.AddressRegistered{TxHash: "0xlast", RskAddress: catchUpAddress(3), BlockNumber: 105, LogIndex: 4},
		)
		fixture.withPage(107, 100, 105, []blockchain.AddressRegistered{chain[2], chain[0], chain[1]})
		fixture.withContractRoot(105, chain[2].RegistrationRoot)
		fixture.repository.EXPECT().SaveVerifiedRegistrations(
			mock.Anything,
			mock.Anything,
			rootstock.PegInWatchCheckpoint{Root: chain[2].RegistrationRoot, Block: 105},
		).Return(nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assert.Equal(t, []string{"0xfirst", "0xmiddle", "0xlast"}, savedTxHashes(result.Watches))
	})
	t.Run("hashes nothing when there are no new registrations", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{9}, Block: 150}
		fixture.withCheckpoint(&stored)
		fixture.withPage(160, 151, 158, []blockchain.AddressRegistered{})
		fixture.withContractRoot(158, stored.Root)
		fixture.repository.EXPECT().SaveVerifiedRegistrations(mock.Anything, mock.Anything, rootstock.PegInWatchCheckpoint{Root: stored.Root, Block: 158}).
			Return(nil).Once()

		_, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assert.Zero(t, fixture.hashCalls)
		fixture.repository.AssertNotCalled(t, "List", mock.Anything)
	})
	t.Run("folds registrations in one block by log index, whatever their transaction hashes", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(nil)
		lower := blockchain.AddressRegistered{TxHash: "0xffff", RskAddress: "0x00000000000000000000000000000000000000a1", BlockNumber: 105, LogIndex: 1}
		higher := blockchain.AddressRegistered{TxHash: "0x0000", RskAddress: "0x00000000000000000000000000000000000000a2", BlockNumber: 105, LogIndex: 4}
		pinned := [32]byte(common.HexToHash("0x95bdcbc864b6248af173e5feb03f5830479e972a492e705c3c2a2bddcb8ca643"))
		chain := chainedRegistrations(t, [32]byte{}, lower, higher)
		fixture.withPage(107, 100, 105, []blockchain.AddressRegistered{chain[1], chain[0]})
		fixture.withContractRoot(105, pinned)
		fixture.repository.EXPECT().SaveVerifiedRegistrations(mock.Anything, mock.Anything, rootstock.PegInWatchCheckpoint{Root: pinned, Block: 105}).
			Return(nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assert.Equal(t, pinned, result.Checkpoint.Root)
		assert.Equal(t, []string{"0xffff", "0x0000"}, savedTxHashes(result.Watches))
	})
	t.Run("folds registrations in one transaction by log index", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(nil)
		chain := chainedRegistrations(t, [32]byte{},
			blockchain.AddressRegistered{TxHash: "0xaaaa", RskAddress: catchUpAddress(1), BlockNumber: 103, LogIndex: 2},
			blockchain.AddressRegistered{TxHash: "0xaaaa", RskAddress: catchUpAddress(2), BlockNumber: 103, LogIndex: 7},
		)
		fixture.withPage(105, 100, 103, []blockchain.AddressRegistered{chain[1], chain[0]})
		fixture.withContractRoot(103, chain[1].RegistrationRoot)
		fixture.repository.EXPECT().SaveVerifiedRegistrations(
			mock.Anything,
			mock.Anything,
			rootstock.PegInWatchCheckpoint{Root: chain[1].RegistrationRoot, Block: 103},
		).Return(nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assert.Equal(t, chain[1].RegistrationRoot, result.Checkpoint.Root)
	})
	t.Run("rejects an event with an invalid address and saves nothing", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(nil)
		fixture.withPage(105, 100, 103, []blockchain.AddressRegistered{{TxHash: "0x05", RskAddress: "0x1234", BlockNumber: 101}})

		_, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, blockchain.InvalidAddressError)
	})
	t.Run("does nothing while the chain is shorter than the finality depth", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(nil)
		fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(1), nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assert.Empty(t, result.Watches)
	})
	t.Run("does nothing while the start block is above the finalized head", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(nil)
		fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(50), nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assert.Empty(t, result.Watches)
		assert.Equal(t, rootstock.PegInWatchCheckpoint{}, result.Checkpoint)
	})
}

func TestCatchUpRegisteredAddressesUseCase_RootChecks(t *testing.T) {
	t.Run("reports a mismatch with an event root and saves nothing", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{1}, Block: 150}
		fixture.withCheckpoint(&stored)
		events := chainedRegistrations(t, [32]byte{2},
			blockchain.AddressRegistered{TxHash: "0x04", RskAddress: catchUpAddress(4), BlockNumber: 151},
		)
		fixture.withPage(155, 151, 153, events)

		result, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, watcher.RegistryRootMismatchError)
		assert.Contains(t, err.Error(), "at block 151")
		assert.Empty(t, result.Watches)
	})
	t.Run("reports a mismatch when the contract root moved without events and saves nothing", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{1}, Block: 150}
		fixture.withCheckpoint(&stored)
		fixture.withPage(155, 151, 153, []blockchain.AddressRegistered{})
		fixture.withContractRoot(153, [32]byte{2})

		result, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, watcher.RegistryRootMismatchError)
		assert.Contains(t, err.Error(), "at block 153")
		assert.Empty(t, result.Watches)
	})
}

func TestCatchUpRegisteredAddressesUseCase_RunErrors(t *testing.T) {
	stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{1}, Block: 150}
	t.Run("returns the checkpoint read error", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.repository.EXPECT().GetCheckpoint(mock.Anything).Return(rootstock.PegInWatchCheckpoint{}, false, assert.AnError).Once()

		_, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, assert.AnError)
	})
	t.Run("changes no data when the height read fails", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(&stored)
		fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(0), assert.AnError).Once()

		_, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, assert.AnError)
	})
	t.Run("changes no data when the event read fails", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(&stored)
		fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(155), nil).Once()
		fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(151), blockPointer(153)).Return(nil, assert.AnError).Once()

		_, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, assert.AnError)
	})
	t.Run("changes no data when the contract root read fails", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(&stored)
		fixture.withPage(155, 151, 153, []blockchain.AddressRegistered{})
		fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(153)).Return([32]byte{}, assert.AnError).Once()

		_, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, assert.AnError)
	})
	t.Run("returns the save error", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(&stored)
		fixture.withPage(155, 151, 153, []blockchain.AddressRegistered{})
		fixture.withContractRoot(153, stored.Root)
		fixture.repository.EXPECT().SaveVerifiedRegistrations(mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, result.Watches)
	})
}

func TestCatchUpRegisteredAddressesUseCase_SavesRegistrationsInChainOrder(t *testing.T) {
	fixture := newCatchUpFixture(t)
	fixture.withCheckpoint(nil)
	chain := chainedRegistrations(t, [32]byte{},
		blockchain.AddressRegistered{TxHash: "0x04", RskAddress: catchUpAddress(4), BlockNumber: 101, LogIndex: 5},
		blockchain.AddressRegistered{TxHash: "0x03", RskAddress: catchUpAddress(3), BlockNumber: 102, LogIndex: 2},
		blockchain.AddressRegistered{TxHash: "0x02", RskAddress: catchUpAddress(2), BlockNumber: 102, LogIndex: 9},
		blockchain.AddressRegistered{TxHash: "0x01", RskAddress: catchUpAddress(1), BlockNumber: 103, LogIndex: 0},
	)
	fixture.withPage(105, 100, 103, []blockchain.AddressRegistered{chain[3], chain[2], chain[1], chain[0]})
	fixture.withContractRoot(103, chain[3].RegistrationRoot)
	fixture.repository.EXPECT().SaveVerifiedRegistrations(mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	result, err := fixture.useCase.Run(context.Background())

	require.NoError(t, err)
	assert.Equal(t, []string{"0x04", "0x03", "0x02", "0x01"}, savedTxHashes(result.Watches))
}

func TestCatchUpRegisteredAddressesUseCase_PageRange(t *testing.T) {
	stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{9}, Block: 109}
	cases := []struct {
		name          string
		storedBlock   uint64
		finalizedHead uint64
		pageSize      uint64
		toBlock       uint64
		hasPage       bool
	}{
		{name: "reads nothing past the finalized head", storedBlock: 110, finalizedHead: 110, pageSize: 5},
		{name: "reads a single block", storedBlock: 109, finalizedHead: 110, pageSize: 5, toBlock: 110, hasPage: true},
		{name: "stops at the finalized head", storedBlock: 109, finalizedHead: 113, pageSize: 5, toBlock: 113, hasPage: true},
		{name: "ends a full page at the finalized head", storedBlock: 109, finalizedHead: 114, pageSize: 5, toBlock: 114, hasPage: true},
		{name: "stops at the page size", storedBlock: 109, finalizedHead: 200, pageSize: 5, toBlock: 114, hasPage: true},
		{name: "reads to the finalized head with a zero page size", storedBlock: 109, finalizedHead: 200, toBlock: 200, hasPage: true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newCatchUpFixtureWithPageSize(t, testCase.pageSize)
			checkpoint := rootstock.PegInWatchCheckpoint{Root: stored.Root, Block: testCase.storedBlock}
			fixture.withCheckpoint(&checkpoint)
			expected := checkpoint
			if testCase.hasPage {
				fixture.withPage(testCase.finalizedHead+catchUpFinalityDepth, testCase.storedBlock+1, testCase.toBlock, []blockchain.AddressRegistered{})
				fixture.withContractRoot(testCase.toBlock, stored.Root)
				expected = rootstock.PegInWatchCheckpoint{Root: stored.Root, Block: testCase.toBlock}
				fixture.repository.EXPECT().SaveVerifiedRegistrations(mock.Anything, []rootstock.PegInWatch{}, expected).Return(nil).Once()
			} else {
				fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(testCase.finalizedHead+catchUpFinalityDepth, nil).Once()
			}

			result, err := fixture.useCase.Run(context.Background())

			require.NoError(t, err)
			assert.Equal(t, expected, result.Checkpoint)
			assert.Empty(t, result.Watches)
		})
	}
}
