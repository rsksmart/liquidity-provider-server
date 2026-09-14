package watcher

import "context"

// findFirstMismatchingBlock binary searches the lowest block in [startBlock, head]
// where matchesAt no longer reports a match. It assumes matchesAt is monotonic:
// once a block mismatches, every higher block mismatches too. When every probed
// block matches, it returns head.
func findFirstMismatchingBlock(
	ctx context.Context,
	matchesAt func(context.Context, uint64) (bool, error),
	startBlock uint64,
	head uint64,
) (uint64, error) {
	low, high := startBlock, head
	for low < high {
		middle := low + (high-low)/2
		matches, err := matchesAt(ctx, middle)
		if err != nil {
			return 0, err
		}
		if matches {
			low = middle + 1
		} else {
			high = middle
		}
	}
	return low, nil
}
