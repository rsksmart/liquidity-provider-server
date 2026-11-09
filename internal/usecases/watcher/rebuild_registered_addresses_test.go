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
	"github.com/rsksmart/liquidity-provider-server/test"
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
	fixture.withStoredRoots(rootTimeline(rootPoint{110, shared}, rootPoint{130, oldFork}, rootPoint{140, oldForkTip}), 8)
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

func TestCatchUpRegisteredAddressesUseCase_RebuildsToTheChainRootWhenRootsMatchAgainBelowTheMismatch(t *testing.T) {
	fixture := newCatchUpFixture(t)
	rootA := foldedRoot(t, [32]byte{}, catchUpAddress(1))
	rootAB := foldedRoot(t, rootA, catchUpAddress(2))
	storedTip := foldedRoot(t, rootAB, catchUpAddress(3))
	canonical := chainedRegistrations(t, rootAB,
		blockchain.AddressRegistered{TxHash: "0x0d", RskAddress: catchUpAddress(4), BlockNumber: 140},
	)
	chainTip := canonical[0].RegistrationRoot
	stored := rootstock.PegInWatchCheckpoint{Root: storedTip, Block: 150}
	fixture.withCheckpoint(&stored)
	fixture.withStoredRoots(rootTimeline(rootPoint{110, rootA}, rootPoint{120, rootAB}, rootPoint{140, storedTip}), 8)
	fixture.withChainRoots(rootTimeline(rootPoint{115, rootA}, rootPoint{120, rootAB}, rootPoint{140, chainTip}), 9)
	fixture.withPage(160, 151, 158, []blockchain.AddressRegistered{})
	fixture.withEvents(140, 149, canonical)
	fixture.withEvents(150, 158, []blockchain.AddressRegistered{})
	next := rootstock.PegInWatchCheckpoint{Root: chainTip, Block: 158}
	fixture.repository.EXPECT().ReplaceVerifiedRegistrations(mock.Anything, uint64(140), mock.Anything, next).Return(nil).Once()

	result, err := fixture.useCase.Run(context.Background())

	require.NoError(t, err)
	assertCatchUpResult(t, expectedCatchUpResult{txHashes: []string{"0x0d"}, checkpoint: next, rebuilt: true}, result)
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

func TestCatchUpRegisteredAddressesUseCase_RebuildsFromBlockZeroWithoutReadingARootBeforeIt(t *testing.T) {
	fixture := newCatchUpFixtureWithStartBlockAndPageSize(t, 0, 0)
	stale := foldedRoot(t, [32]byte{}, catchUpAddress(1))
	events := chainedRegistrations(t, [32]byte{},
		blockchain.AddressRegistered{TxHash: "0x02", RskAddress: catchUpAddress(2), BlockNumber: 0},
	)
	chainRoot := events[0].RegistrationRoot
	stored := rootstock.PegInWatchCheckpoint{Root: stale, Block: 10}
	fixture.withCheckpoint(&stored)
	fixture.withPage(15, 11, 13, []blockchain.AddressRegistered{})
	var storedReads []uint64
	fixture.repository.EXPECT().GetRootAt(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, block uint64) ([32]byte, error) {
			storedReads = append(storedReads, block)
			return rootTimeline(rootPoint{0, stale})(block), nil
		})
	fixture.withChainRoots(rootTimeline(rootPoint{0, chainRoot}), 7)
	fixture.withEvents(0, 13, events)
	next := rootstock.PegInWatchCheckpoint{Root: chainRoot, Block: 13}
	fixture.repository.EXPECT().ReplaceVerifiedRegistrations(mock.Anything, uint64(0), mock.Anything, next).Return(nil).Once()

	result, err := fixture.useCase.Run(context.Background())

	require.NoError(t, err)
	assertCatchUpResult(t, expectedCatchUpResult{txHashes: []string{"0x02"}, checkpoint: next, rebuilt: true}, result)
	for _, block := range storedReads {
		assert.LessOrEqual(t, block, stored.Block)
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
		assertCatchUpResult(t, expectedCatchUpResult{checkpoint: stored}, result)
		fixture.repository.AssertNotCalled(t, "ReplaceVerifiedRegistrations", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
	t.Run("changes no data when the node cannot read a searched block", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		mismatchedPage(fixture)
		fixture.withContractRoot(150, [32]byte{3})
		fixture.repository.EXPECT().GetRootAt(mock.Anything, uint64(150)).Return(stored.Root, nil).Once()
		fixture.repository.EXPECT().GetRootAt(mock.Anything, uint64(125)).Return([32]byte{}, nil).Once()
		fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, uint64(125)).Return([32]byte{}, assert.AnError).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, watcher.RegistryOldBlockReadError)
		assert.Contains(t, err.Error(), "old block 125")
		assertCatchUpResult(t, expectedCatchUpResult{checkpoint: stored}, result)
		fixture.repository.AssertNotCalled(t, "ReplaceVerifiedRegistrations", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
	t.Run("changes no data when the stored root read fails", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		mismatchedPage(fixture)
		fixture.withContractRoot(150, [32]byte{3})
		fixture.repository.EXPECT().GetRootAt(mock.Anything, uint64(150)).Return(stored.Root, nil).Once()
		fixture.repository.EXPECT().GetRootAt(mock.Anything, uint64(125)).Return([32]byte{}, assert.AnError).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, assert.AnError)
		assertCatchUpResult(t, expectedCatchUpResult{checkpoint: stored}, result)
		fixture.repository.AssertNotCalled(t, "ReplaceVerifiedRegistrations", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
	t.Run("changes no data when the stored root before the first bad block cannot be read", func(t *testing.T) {
		fixture := newCatchUpFixtureWithPageSize(t, 0)
		stale := foldedRoot(t, [32]byte{}, catchUpAddress(1))
		chainRoot := foldedRoot(t, [32]byte{}, catchUpAddress(2))
		stale150 := rootstock.PegInWatchCheckpoint{Root: stale, Block: 150}
		fixture.withCheckpoint(&stale150)
		fixture.withPage(155, 151, 153, []blockchain.AddressRegistered{})
		readsOfBlock100 := 0
		fixture.repository.EXPECT().GetRootAt(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, block uint64) ([32]byte, error) {
				if block == 100 {
					readsOfBlock100++
					if readsOfBlock100 == 2 {
						return [32]byte{}, assert.AnError
					}
				}
				return rootTimeline(rootPoint{101, stale})(block), nil
			})
		fixture.registry.EXPECT().GetRegistrationRoot(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, block uint64) ([32]byte, error) {
				return rootTimeline(rootPoint{101, chainRoot})(block), nil
			})

		result, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, assert.AnError)
		assert.Contains(t, err.Error(), "stored PegIn watch root at block 100")
		assertCatchUpResult(t, expectedCatchUpResult{checkpoint: stale150}, result)
		assert.Equal(t, 2, readsOfBlock100)
		fixture.repository.AssertNotCalled(t, "ReplaceVerifiedRegistrations", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
	t.Run("changes no data when the stored rows cannot be read at the stored block", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		mismatchedPage(fixture)
		fixture.withContractRoot(150, [32]byte{3})
		fixture.repository.EXPECT().GetRootAt(mock.Anything, uint64(150)).Return([32]byte{}, assert.AnError).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, assert.AnError)
		assertCatchUpResult(t, expectedCatchUpResult{checkpoint: stored}, result)
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
		assertCatchUpResult(t, expectedCatchUpResult{checkpoint: stored}, result)
	})
}

func TestCatchUpRegisteredAddressesUseCase_RebuildsFromTheStartBlockOnAFirstStart(t *testing.T) {
	t.Run("gets every earlier registration and reaches the head in one run", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(nil)
		events := chainedRegistrations(t, [32]byte{},
			blockchain.AddressRegistered{TxHash: "0x01", RskAddress: catchUpAddress(1), BlockNumber: 100},
			blockchain.AddressRegistered{TxHash: "0x02", RskAddress: catchUpAddress(2), BlockNumber: 115},
			blockchain.AddressRegistered{TxHash: "0x03", RskAddress: catchUpAddress(3), BlockNumber: 129, LogIndex: 4},
			blockchain.AddressRegistered{TxHash: "0x04", RskAddress: catchUpAddress(4), BlockNumber: 130},
		)
		head := rootstock.PegInWatchCheckpoint{Root: events[3].RegistrationRoot, Block: 130}
		fixture.repository.EXPECT().GetRootAt(mock.Anything, uint64(99)).Return([32]byte{}, nil).Once()
		fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(132), nil).Once()
		fixture.withEvents(100, 109, events[:1])
		fixture.withEvents(110, 119, events[1:2])
		fixture.withEvents(120, 129, events[2:3])
		fixture.withEvents(130, 130, events[3:])
		fixture.withContractRoot(130, head.Root)
		fixture.repository.EXPECT().ReplaceVerifiedRegistrations(mock.Anything, uint64(100), mock.Anything, head).Return(nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assertCatchUpResult(t, expectedCatchUpResult{txHashes: []string{"0x01", "0x02", "0x03", "0x04"}, checkpoint: head, rebuilt: true}, result)
	})
	t.Run("folds from the stored rows below the start block", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(nil)
		below := foldedRoot(t, [32]byte{}, catchUpAddress(1))
		events := chainedRegistrations(t, below,
			blockchain.AddressRegistered{TxHash: "0x02", RskAddress: catchUpAddress(2), BlockNumber: 103},
		)
		next := rootstock.PegInWatchCheckpoint{Root: events[0].RegistrationRoot, Block: 105}
		fixture.repository.EXPECT().GetRootAt(mock.Anything, uint64(99)).Return(below, nil).Once()
		fixture.withPage(107, 100, 105, events)
		fixture.withContractRoot(105, next.Root)
		fixture.repository.EXPECT().ReplaceVerifiedRegistrations(mock.Anything, uint64(100), mock.Anything, next).Return(nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assertCatchUpResult(t, expectedCatchUpResult{txHashes: []string{"0x02"}, checkpoint: next, rebuilt: true}, result)
	})
}

func TestCatchUpRegisteredAddressesUseCase_FirstStartErrors(t *testing.T) {
	t.Run("changes no data when the first start does not match the chain", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(nil)
		fixture.repository.EXPECT().GetRootAt(mock.Anything, uint64(99)).Return([32]byte{}, nil).Once()
		fixture.withPage(107, 100, 105, []blockchain.AddressRegistered{})
		fixture.withContractRoot(105, [32]byte{7})

		result, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, watcher.RegistryRootMismatchError)
		assert.Contains(t, err.Error(), "from block 100")
		assertCatchUpResult(t, expectedCatchUpResult{}, result)
		fixture.repository.AssertNotCalled(t, "ReplaceVerifiedRegistrations", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
	t.Run("changes no data when the stored root read fails", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(nil)
		fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(107), nil).Once()
		fixture.repository.EXPECT().GetRootAt(mock.Anything, uint64(99)).Return([32]byte{}, assert.AnError).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.ErrorIs(t, err, assert.AnError)
		assertCatchUpResult(t, expectedCatchUpResult{}, result)
		fixture.repository.AssertNotCalled(t, "ReplaceVerifiedRegistrations", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestCatchUpRegisteredAddressesUseCase_HandlesAnInvalidStoredPoint(t *testing.T) {
	t.Run("rebuilds from the start block when the stored point is below it", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		fixture.withCheckpoint(&rootstock.PegInWatchCheckpoint{Root: [32]byte{3}, Block: 50})
		events := chainedRegistrations(t, [32]byte{},
			blockchain.AddressRegistered{TxHash: "0x01", RskAddress: catchUpAddress(1), BlockNumber: 102},
		)
		next := rootstock.PegInWatchCheckpoint{Root: events[0].RegistrationRoot, Block: 105}
		fixture.repository.EXPECT().GetRootAt(mock.Anything, uint64(99)).Return([32]byte{}, nil).Once()
		fixture.withPage(107, 100, 105, events)
		fixture.withContractRoot(105, next.Root)
		fixture.repository.EXPECT().ReplaceVerifiedRegistrations(mock.Anything, uint64(100), mock.Anything, next).Return(nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assertCatchUpResult(t, expectedCatchUpResult{txHashes: []string{"0x01"}, checkpoint: next, rebuilt: true}, result)
	})
	t.Run("waits without reading or writing while the stored point is above the finalized head", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{3}, Block: 500}
		fixture.withCheckpoint(&stored)
		fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(300), nil).Once()
		logged := test.AssertLogContains(t, watcher.LogRegistryPointAboveFinalizedHead(500, 298))

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assert.True(t, logged())
		assertCatchUpResult(t, expectedCatchUpResult{checkpoint: stored}, result)
		fixture.registry.AssertNotCalled(t, "GetAddressRegisteredEvents", mock.Anything, mock.Anything, mock.Anything)
		fixture.repository.AssertNotCalled(t, "SaveVerifiedRegistrations", mock.Anything, mock.Anything, mock.Anything)
		fixture.repository.AssertNotCalled(t, "ReplaceVerifiedRegistrations", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
	t.Run("keeps the stored point while it is below the start block and the start block is above the finalized head", func(t *testing.T) {
		fixture := newCatchUpFixture(t)
		stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{3}, Block: 50}
		fixture.withCheckpoint(&stored)
		fixture.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(60), nil).Once()

		result, err := fixture.useCase.Run(context.Background())

		require.NoError(t, err)
		assertCatchUpResult(t, expectedCatchUpResult{checkpoint: stored}, result)
		fixture.registry.AssertNotCalled(t, "GetAddressRegisteredEvents", mock.Anything, mock.Anything, mock.Anything)
		fixture.repository.AssertNotCalled(t, "ReplaceVerifiedRegistrations", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestCatchUpRegisteredAddressesUseCase_RebuildsFromTheStartBlockWhenStoredRowsDoNotMatchTheStoredPoint(t *testing.T) {
	fixture := newCatchUpFixture(t)
	stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{3}, Block: 150}
	fixture.withCheckpoint(&stored)
	events := chainedRegistrations(t, [32]byte{},
		blockchain.AddressRegistered{TxHash: "0x01", RskAddress: catchUpAddress(1), BlockNumber: 120},
	)
	next := rootstock.PegInWatchCheckpoint{Root: events[0].RegistrationRoot, Block: 153}
	fixture.withPage(155, 151, 153, []blockchain.AddressRegistered{})
	fixture.withContractRoot(153, next.Root)
	fixture.withContractRoot(150, next.Root)
	fixture.repository.EXPECT().GetRootAt(mock.Anything, uint64(150)).Return([32]byte{4}, nil).Once()
	fixture.repository.EXPECT().GetRootAt(mock.Anything, uint64(99)).Return([32]byte{}, nil).Once()
	fixture.withEvents(100, 109, []blockchain.AddressRegistered{})
	fixture.withEvents(110, 119, []blockchain.AddressRegistered{})
	fixture.withEvents(120, 129, events)
	fixture.withEvents(130, 139, []blockchain.AddressRegistered{})
	fixture.withEvents(140, 149, []blockchain.AddressRegistered{})
	fixture.withEvents(150, 153, []blockchain.AddressRegistered{})
	fixture.withContractRoot(153, next.Root)
	fixture.repository.EXPECT().ReplaceVerifiedRegistrations(mock.Anything, uint64(100), mock.Anything, next).Return(nil).Once()
	logged := test.AssertLogContains(t, watcher.LogRegistryRowsDoNotMatchPoint(150, 100))

	result, err := fixture.useCase.Run(context.Background())

	require.NoError(t, err)
	assert.True(t, logged())
	assertCatchUpResult(t, expectedCatchUpResult{txHashes: []string{"0x01"}, checkpoint: next, rebuilt: true}, result)
}
