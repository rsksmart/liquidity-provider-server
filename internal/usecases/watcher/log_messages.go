package watcher

import "fmt"

const LogRegistryCatchUpMismatch = "PegIn address registry catch-up failed, rebuilding the watch list: %v"

func LogRegistryPointAboveFinalizedHead(block, finalizedHead uint64) string {
	return fmt.Sprintf(
		"Stored PegIn address registry point at block %d is above the finalized RSK head %d, waiting for the node to catch up",
		block,
		finalizedHead,
	)
}

func LogRegistryRowsDoNotMatchPoint(block, startBlock uint64) string {
	return fmt.Sprintf("Stored PegIn watch rows do not match the stored point at block %d, rebuilding from the start block %d", block, startBlock)
}

func LogRegistryWatchListRebuilt(fromBlock, toBlock uint64, registrations int) string {
	return fmt.Sprintf("Rebuilt the PegIn watch list from block %d to block %d with %d registrations", fromBlock, toBlock, registrations)
}
