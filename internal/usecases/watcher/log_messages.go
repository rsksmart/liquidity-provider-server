package watcher

import "fmt"

const LogRegistryCatchUpMismatch = "PegIn address registry catch-up failed, rebuilding the watch list: %v"

func LogRegistryCheckpointAboveFinalizedHead(block, finalizedHead uint64) string {
	return fmt.Sprintf(
		"Stored PegIn address registry checkpoint at block %d is above the finalized RSK head %d, waiting for the node to catch up",
		block,
		finalizedHead,
	)
}

func LogRegistryDocumentsDoNotMatchCheckpoint(block, startBlock uint64) string {
	return fmt.Sprintf("Stored PegIn watch documents do not match the checkpoint at block %d, rebuilding from the start block %d", block, startBlock)
}

func LogRegistryStartBlockAboveFinalizedHead(startBlock, finalizedHead uint64) string {
	return fmt.Sprintf(
		"PegIn address registry start block %d is above the finalized RSK head %d, waiting until the head reaches it",
		startBlock,
		finalizedHead,
	)
}

func LogRegistryWatchListRebuilt(fromBlock, toBlock uint64, registrations int) string {
	return fmt.Sprintf("Rebuilt the PegIn watch list from block %d to block %d with %d registrations", fromBlock, toBlock, registrations)
}
