package watcher

import "fmt"

func LogRegistryEventRootMismatch(block uint64, localRoot, eventRoot [32]byte) string {
	return fmt.Sprintf("PegIn address registry root mismatch at block %d: local root 0x%x, event root 0x%x", block, localRoot, eventRoot)
}

func LogRegistryContractRootMismatch(block uint64, localRoot, contractRoot [32]byte) string {
	return fmt.Sprintf("PegIn address registry root mismatch at block %d: local root 0x%x, contract root 0x%x", block, localRoot, contractRoot)
}
