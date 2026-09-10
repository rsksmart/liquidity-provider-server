package watcher

import (
	"fmt"
	"sort"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
)

// RegistrationRoot is intentionally absent. A stored row cannot prove its own integrity.
type localRegistration struct {
	blockNumber uint64
	logIndex    uint
	txHash      string
	rskAddress  string
}

type localRootPoint struct {
	blockNumber uint64
	root        [32]byte
}

// localRootTimeline holds the local registration root at the end of every block that registered
// an address. It answers the root at any height without another database read.
type localRootTimeline struct {
	points []localRootPoint
}

// NewLocalRootTimeline folds the stored addresses in [startBlock, head] into a root per block.
// A row outside that range cannot be compared against the chain, so the timeline refuses to
// describe it and returns an error instead.
func NewLocalRootTimeline(
	entries []rootstock.PegInWatch,
	startBlock uint64,
	head uint64,
	hashFunction entities.HashFunction,
) (localRootTimeline, error) {
	registrations, err := collectLocalRegistrations(entries, startBlock, head)
	if err != nil {
		return localRootTimeline{}, err
	}
	sort.Slice(registrations, func(firstIndex, secondIndex int) bool {
		return localRegistrationBefore(registrations[firstIndex], registrations[secondIndex])
	})

	timeline := localRootTimeline{}
	root := [32]byte{}
	for _, registration := range registrations {
		root, err = blockchain.FoldPegInAddressRegistryRoot(
			hashFunction,
			root,
			registration.rskAddress,
		)
		if err != nil {
			return localRootTimeline{}, fmt.Errorf(
				"fold stored AddressRegistered event %s/%d: %w",
				registration.txHash,
				registration.logIndex,
				err,
			)
		}
		timeline.points = recordRoot(timeline.points, registration.blockNumber, root)
	}
	return timeline, nil
}

func collectLocalRegistrations(
	entries []rootstock.PegInWatch,
	startBlock uint64,
	head uint64,
) ([]localRegistration, error) {
	registrations := make([]localRegistration, 0, len(entries))
	for _, entry := range entries {
		if entry.BlockNumber < startBlock {
			return nil, fmt.Errorf("PegIn watch exists below configured start block %d", startBlock)
		}
		if entry.BlockNumber > head {
			return nil, fmt.Errorf("PegIn watch exists above current head %d", head)
		}
		registrations = append(registrations, localRegistration{
			blockNumber: entry.BlockNumber,
			logIndex:    entry.LogIndex,
			txHash:      entry.TxHash,
			rskAddress:  entry.RskAddress,
		})
	}
	return registrations, nil
}

func localRegistrationBefore(first localRegistration, second localRegistration) bool {
	if first.blockNumber != second.blockNumber {
		return first.blockNumber < second.blockNumber
	}
	if first.logIndex != second.logIndex {
		return first.logIndex < second.logIndex
	}
	if first.txHash != second.txHash {
		return first.txHash < second.txHash
	}
	return first.rskAddress < second.rskAddress
}

// recordRoot keeps one point per block: the root after the last registration of that block.
func recordRoot(points []localRootPoint, blockNumber uint64, root [32]byte) []localRootPoint {
	point := localRootPoint{blockNumber: blockNumber, root: root}
	lastPoint := len(points) - 1
	if lastPoint >= 0 && points[lastPoint].blockNumber == blockNumber {
		points[lastPoint] = point
		return points
	}
	return append(points, point)
}

// RootAt returns the local root at the given height. It returns the last root at or below the
// height, and the zero root before the first registration.
func (timeline localRootTimeline) RootAt(height uint64) [32]byte {
	index := sort.Search(len(timeline.points), func(index int) bool {
		return timeline.points[index].blockNumber > height
	})
	if index == 0 {
		return [32]byte{}
	}
	return timeline.points[index-1].root
}
