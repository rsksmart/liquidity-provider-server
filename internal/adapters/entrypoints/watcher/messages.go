package watcher

// Log message templates for the cold-wallet transfer watcher.
// Format arguments are filled at the call site with logrus *f variants.
// TODO FLY-2388: this is a pilot of the repo-wide log message centralization.
// Apply the same pattern (messages.go per package) to the remaining packages.
const (
	LogTransferError           = "TransferColdWalletWatcher: Error executing transfer to cold wallet: %v"
	LogTransferShutdown        = "TransferColdWalletWatcher shut down"
	LogTransferSuccess         = "TransferColdWalletWatcher: %s transfer successful - TxHash: %s, Amount: %s, Fee: %s"
	LogTransferSkippedNoExcess = "TransferColdWalletWatcher: %s transfer skipped - no excess liquidity"
	LogTransferSkippedNotEcon  = "TransferColdWalletWatcher: %s transfer skipped - not economical: %s"
	LogTransferSkippedCooldown = "TransferColdWalletWatcher: %s transfer skipped - liquidity target cooldown active"
	LogTransferFailed          = "TransferColdWalletWatcher: %s transfer failed - %s: %v"
)
