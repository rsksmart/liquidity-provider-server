package watcher_test

import (
	"context"
	"fmt"
	"testing"

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
	useCase    *watcher.CatchUpRegisteredAddressesUseCase
}

func newCatchUpFixture(t *testing.T) *catchUpFixture {
	t.Helper()
	fixture := &catchUpFixture{
		repository: mocks.NewPegInWatchRepositoryMock(t),
		registry:   mocks.NewPegInAddressRegistryContractMock(t),
		rskRpc:     mocks.NewRootstockRpcServerMock(t),
	}
	fixture.useCase = watcher.NewCatchUpRegisteredAddressesUseCase(
		fixture.repository,
		fixture.registry,
		fixture.rskRpc,
		catchUpStartBlock,
		catchUpPageSize,
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
		events := []blockchain.AddressRegistered{
			{TxHash: "0x01", RskAddress: "0xa1", Registrant: "0xr", RegistrationRoot: [32]byte{1}, BlockNumber: 101},
			{TxHash: "0x02", RskAddress: "0xa2", Registrant: "0xr", RegistrationRoot: [32]byte{2}, BlockNumber: 104, LogIndex: 3},
		}
		fixture.withPage(200, 100, 109, events)
		fixture.repository.EXPECT().SaveVerifiedRegistrations(
			mock.Anything,
			mock.MatchedBy(func(watches []rootstock.PegInWatch) bool {
				return len(watches) == 2 &&
					watches[0].TxHash == "0x01" && watches[0].RskAddress == "0xa1" && watches[0].Registrant == "0xr" &&
					watches[0].State == rootstock.PegInWatchDiscovered &&
					watches[1].BlockNumber == 104 && watches[1].LogIndex == 3 &&
					watches[1].RegistrationRoot == [32]byte{2}
			}),
			rootstock.PegInWatchCheckpoint{Root: [32]byte{2}, Block: 109},
		).Return(nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assert.Equal(t, rootstock.PegInWatchCheckpoint{Root: [32]byte{2}, Block: 109}, result.Checkpoint)
		assert.Equal(t, []string{"0x01", "0x02"}, savedTxHashes(result.Watches))
	})
	t.Run("continues from the stored point, rereading the finality window", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(&rootstock.PegInWatchCheckpoint{Root: [32]byte{9}, Block: 150})
		events := []blockchain.AddressRegistered{
			{TxHash: "0x03", RskAddress: "0xa3", RegistrationRoot: [32]byte{3}, BlockNumber: 152, LogIndex: 1},
		}
		fixture.withPage(155, 149, 153, events)
		fixture.repository.EXPECT().SaveVerifiedRegistrations(mock.Anything, mock.Anything, rootstock.PegInWatchCheckpoint{Root: [32]byte{3}, Block: 153}).
			Return(nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assert.Equal(t, uint64(153), result.Checkpoint.Block)
		assert.Equal(t, []string{"0x03"}, savedTxHashes(result.Watches))
	})
	t.Run("reads one page at a time", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(&rootstock.PegInWatchCheckpoint{Root: [32]byte{9}, Block: 120})
		fixture.withPage(200, 119, 130, []blockchain.AddressRegistered{})
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
		fixture.withPage(160, 149, 158, []blockchain.AddressRegistered{})
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
	t.Run("stores the root of the last registration in chain order", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(nil)
		events := []blockchain.AddressRegistered{
			{TxHash: "0xlast", RegistrationRoot: [32]byte{3}, BlockNumber: 105, LogIndex: 4},
			{TxHash: "0xfirst", RegistrationRoot: [32]byte{1}, BlockNumber: 103, LogIndex: 9},
			{TxHash: "0xmiddle", RegistrationRoot: [32]byte{2}, BlockNumber: 105, LogIndex: 1},
		}
		fixture.withPage(107, 100, 105, events)
		fixture.repository.EXPECT().SaveVerifiedRegistrations(mock.Anything, mock.Anything, rootstock.PegInWatchCheckpoint{Root: [32]byte{3}, Block: 105}).
			Return(nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assert.Equal(t, []string{"0xfirst", "0xmiddle", "0xlast"}, savedTxHashes(result.Watches))
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
		fixture.registry.EXPECT().GetAddressRegisteredEvents(mock.Anything, uint64(149), blockPointer(153)).Return(nil, assert.AnError).Once()

		_, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, assert.AnError)
	})
	t.Run("returns the save error", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(&stored)
		fixture.withPage(155, 149, 153, []blockchain.AddressRegistered{})
		fixture.repository.EXPECT().SaveVerifiedRegistrations(mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, result.Watches)
	})
}

func TestSortAddressRegisteredEvents(t *testing.T) {
	events := []blockchain.AddressRegistered{
		{TxHash: "0x01", BlockNumber: 12, LogIndex: 0},
		{TxHash: "0x02", BlockNumber: 11, LogIndex: 9},
		{TxHash: "0x03", BlockNumber: 11, LogIndex: 2},
		{TxHash: "0x04", BlockNumber: 10, LogIndex: 5},
	}

	watcher.SortAddressRegisteredEvents(events)

	ordered := make([]string, 0, len(events))
	for _, event := range events {
		ordered = append(ordered, fmt.Sprintf("%d/%d", event.BlockNumber, event.LogIndex))
	}
	assert.Equal(t, []string{"10/5", "11/2", "11/9", "12/0"}, ordered)
}
