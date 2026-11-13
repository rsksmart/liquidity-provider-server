package watcher

import "fmt"

const LogRegistryCatchUpMismatch = "PegIn address registry catch-up failed, rebuilding the watch list: %v"

func LogRegistryWatchListRebuilt(fromBlock, toBlock uint64, registrations int) string {
	return fmt.Sprintf("Rebuilt the PegIn watch list from block %d to block %d with %d registrations", fromBlock, toBlock, registrations)
}
