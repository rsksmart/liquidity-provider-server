package watcher_test

import (
	"context"
	"math/bits"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type rootPoint struct {
	block uint64
	root  [32]byte
}

func rootTimeline(points ...rootPoint) func(uint64) [32]byte {
	sort.Slice(points, func(first, second int) bool { return points[first].block < points[second].block })
	return func(block uint64) [32]byte {
		root := [32]byte{}
		for _, point := range points {
			if point.block > block {
				break
			}
			root = point.root
		}
		return root
	}
}

func foldedRoot(t *testing.T, previous [32]byte, address string) [32]byte {
	t.Helper()
	root, err := blockchain.FoldPegInAddressRegistryRoot(crypto.Keccak256, previous, address)
	require.NoError(t, err)
	return root
}

func (fixture *catchUpFixture) withStoredRoots(timeline func(uint64) [32]byte, calls int) {
	fixture.repository.EXPECT().GetRootAt(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, block uint64) ([32]byte, error) { return timeline(block), nil }).Times(calls)
}

func (fixture *catchUpFixture) withChainRoots(timeline func(uint64) [32]byte, calls int) {
	fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, block uint64) ([32]byte, error) { return timeline(block), nil }).Times(calls)
}

func TestCatchUpRegisteredAddressesUseCase_RebuildsAfterASkippedEvent(t *testing.T) {
	fixture := newCatchUpFixture(t)
	stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{9}, Block: 150}
	fixture.withCheckpoint(&stored)
	events := chainedRegistrations(t, stored.Root,
		blockchain.AddressRegistered{TxHash: "0x01", RskAddress: catchUpAddress(1), BlockNumber: 152},
		blockchain.AddressRegistered{TxHash: "0x02", RskAddress: catchUpAddress(2), BlockNumber: 155},
	)
	finalRoot := events[1].RegistrationRoot
	fixture.withPage(160, 151, 158, events[:1])
	fixture.withContractRoot(158, finalRoot)
	fixture.withContractRoot(150, stored.Root)
	fixture.withEvents(151, 158, events)
	fixture.withContractRoot(158, finalRoot)
	next := rootstock.PegInWatchCheckpoint{Root: finalRoot, Block: 158}
	fixture.repository.EXPECT().ReplaceVerifiedRegistrations(
		mock.Anything,
		uint64(151),
		mock.MatchedBy(func(watches []rootstock.PegInWatch) bool {
			return len(watches) == 2 && watches[0].TxHash == "0x01" && watches[1].TxHash == "0x02"
		}),
		next,
	).Return(nil).Once()

	result, err := fixture.useCase.Run(context.Background())

	require.NoError(t, err)
	assertCatchUpResult(t, expectedCatchUpResult{txHashes: []string{"0x01", "0x02"}, checkpoint: next, rebuilt: true}, result)
	fixture.repository.AssertNotCalled(t, "SaveVerifiedRegistrations", mock.Anything, mock.Anything, mock.Anything)
	fixture.repository.AssertNotCalled(t, "GetRootAt", mock.Anything, mock.Anything)
}

//nolint:funlen
func TestCatchUpRegisteredAddressesUseCase_RebuildsAfterAReorgBelowTheStoredBlock(t *testing.T) {
	fixture := newCatchUpFixture(t)
	shared := foldedRoot(t, [32]byte{}, catchUpAddress(1))
	oldFork := foldedRoot(t, shared, catchUpAddress(2))
	oldForkTip := foldedRoot(t, oldFork, catchUpAddress(3))
	reorged := chainedRegistrations(t, shared,
		blockchain.AddressRegistered{TxHash: "0x0b", RskAddress: catchUpAddress(4), BlockNumber: 135},
		blockchain.AddressRegistered{TxHash: "0x0c", RskAddress: catchUpAddress(3), BlockNumber: 145, LogIndex: 2},
	)
	newTip := reorged[1].RegistrationRoot
	stored := rootstock.PegInWatchCheckpoint{Root: oldForkTip, Block: 150}
	fixture.withCheckpoint(&stored)
	fixture.withStoredRoots(rootTimeline(rootPoint{110, shared}, rootPoint{130, oldFork}, rootPoint{140, oldForkTip}), 7)
	fixture.withChainRoots(rootTimeline(rootPoint{110, shared}, rootPoint{135, reorged[0].RegistrationRoot}, rootPoint{145, newTip}), 9)
	fixture.withPage(160, 151, 158, []blockchain.AddressRegistered{})
	fixture.withEvents(130, 139, reorged[:1])
	fixture.withEvents(140, 149, reorged[1:])
	fixture.withEvents(150, 158, []blockchain.AddressRegistered{})
	next := rootstock.PegInWatchCheckpoint{Root: newTip, Block: 158}
	fixture.repository.EXPECT().ReplaceVerifiedRegistrations(
		mock.Anything,
		uint64(130),
		mock.MatchedBy(func(watches []rootstock.PegInWatch) bool {
			return len(watches) == 2 && watches[0].RskAddress == catchUpAddress(4) &&
				watches[1].RskAddress == catchUpAddress(3) && watches[1].BlockNumber == 145
		}),
		next,
	).Return(nil).Once()

	result, err := fixture.useCase.Run(context.Background())

	require.NoError(t, err)
	assertCatchUpResult(t, expectedCatchUpResult{txHashes: []string{"0x0b", "0x0c"}, checkpoint: next, rebuilt: true}, result)

	fixture.withCheckpoint(&next)
	fixture.withPage(165, 159, 163, []blockchain.AddressRegistered{})
	fixture.withContractRoot(163, newTip)
	afterRebuild := rootstock.PegInWatchCheckpoint{Root: newTip, Block: 163}
	fixture.repository.EXPECT().SaveVerifiedRegistrations(mock.Anything, []rootstock.PegInWatch{}, afterRebuild).Return(nil).Once()

	nextResult, err := fixture.useCase.Run(context.Background())

	require.NoError(t, err)
	assertCatchUpResult(t, expectedCatchUpResult{checkpoint: afterRebuild}, nextResult)
	fixture.repository.AssertNumberOfCalls(t, "ReplaceVerifiedRegistrations", 1)
}

//nolint:funlen
func TestCatchUpRegisteredAddressesUseCase_FindsTheFirstMismatchingBlock(t *testing.T) {
	cases := []struct {
		name        string
		storedBlock uint64
		firstBad    uint64
	}{
		{name: "the start block is the first bad block", storedBlock: 150, firstBad: 100},
		{name: "the stored block is the first bad block", storedBlock: 150, firstBad: 150},
		{name: "a block in the middle is the first bad block", storedBlock: 150, firstBad: 131},
		{name: "a stored block at the start block needs no search", storedBlock: 100, firstBad: 100},
		{name: "a million blocks need about log2 reads", storedBlock: 1_000_100, firstBad: 777_877},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newCatchUpFixtureWithPageSize(t, 0)
			stale := foldedRoot(t, [32]byte{}, catchUpAddress(1))
			events := chainedRegistrations(t, [32]byte{},
				blockchain.AddressRegistered{TxHash: "0x02", RskAddress: catchUpAddress(2), BlockNumber: testCase.firstBad},
			)
			chainRoot := events[0].RegistrationRoot
			finalizedHead := testCase.storedBlock + 3
			stored := rootstock.PegInWatchCheckpoint{Root: stale, Block: testCase.storedBlock}
			fixture.withCheckpoint(&stored)
			fixture.withPage(finalizedHead+catchUpFinalityDepth, testCase.storedBlock+1, finalizedHead, []blockchain.AddressRegistered{})
			fixture.repository.EXPECT().GetRootAt(mock.Anything, mock.Anything).
				RunAndReturn(func(_ context.Context, block uint64) ([32]byte, error) {
					return rootTimeline(rootPoint{testCase.firstBad, stale})(block), nil
				})
			chainReads := 0
			fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, mock.Anything).
				RunAndReturn(func(_ context.Context, block uint64) ([32]byte, error) {
					chainReads++
					return rootTimeline(rootPoint{testCase.firstBad, chainRoot})(block), nil
				})
			fixture.withEvents(testCase.firstBad, finalizedHead, events)
			next := rootstock.PegInWatchCheckpoint{Root: chainRoot, Block: finalizedHead}
			fixture.repository.EXPECT().ReplaceVerifiedRegistrations(mock.Anything, testCase.firstBad, mock.Anything, next).Return(nil).Once()

			result, err := fixture.useCase.Run(context.Background())

			require.NoError(t, err)
			assertCatchUpResult(t, expectedCatchUpResult{txHashes: []string{"0x02"}, checkpoint: next, rebuilt: true}, result)
			searchReads := chainReads - 3
			assert.LessOrEqual(t, searchReads, bits.Len64(testCase.storedBlock-catchUpStartBlock))
		})
	}
}

//nolint:funlen
func TestCatchUpRegisteredAddressesUseCase_RebuildErrors(t *testing.T) {
	stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{1}, Block: 150}
	mismatchedPage := func(fixture *catchUpFixture) {
		fixture.withCheckpoint(&stored)
		fixture.withPage(155, 151, 153, []blockchain.AddressRegistered{})
		fixture.withContractRoot(153, [32]byte{2})
	}
	t.Run("changes no data when the node cannot read the stored block", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		mismatchedPage(fixture)
		fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(150)).Return([32]byte{}, assert.AnError).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, watcher.RegistryOldBlockReadError)
		require.ErrorIs(t, err, assert.AnError)
		assert.Contains(t, err.Error(), "old block 150")
		assertCatchUpResult(t, expectedCatchUpResult{}, result)
		fixture.repository.AssertNotCalled(t, "ReplaceVerifiedRegistrations", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
	t.Run("changes no data when the node cannot read a searched block", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		mismatchedPage(fixture)
		fixture.withContractRoot(150, [32]byte{3})
		fixture.repository.EXPECT().GetRootAt(mock.Anything, uint64(125)).Return([32]byte{}, nil).Once()
		fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(125)).Return([32]byte{}, assert.AnError).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, watcher.RegistryOldBlockReadError)
		assert.Contains(t, err.Error(), "old block 125")
		assertCatchUpResult(t, expectedCatchUpResult{}, result)
		fixture.repository.AssertNotCalled(t, "ReplaceVerifiedRegistrations", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
	t.Run("changes no data when the stored root read fails", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		mismatchedPage(fixture)
		fixture.withContractRoot(150, [32]byte{3})
		fixture.repository.EXPECT().GetRootAt(mock.Anything, uint64(125)).Return([32]byte{}, assert.AnError).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, assert.AnError)
		assertCatchUpResult(t, expectedCatchUpResult{}, result)
		fixture.repository.AssertNotCalled(t, "ReplaceVerifiedRegistrations", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
	t.Run("returns the replace error", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		mismatchedPage(fixture)
		fixture.withContractRoot(150, stored.Root)
		fixture.withEvents(151, 153, []blockchain.AddressRegistered{})
		fixture.withContractRoot(153, stored.Root)
		fixture.repository.EXPECT().ReplaceVerifiedRegistrations(mock.Anything, uint64(151), mock.Anything, mock.Anything).
			Return(assert.AnError).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, assert.AnError)
		assertCatchUpResult(t, expectedCatchUpResult{}, result)
	})
}
