package cold_wallet

// Log message templates for cold-wallet transfer operations.
// Format arguments are filled at the call site with logrus *f variants.
const (
	LogTransferError           = "TransferColdWalletWatcher: Error executing transfer to cold wallet: "
	LogTransferShutdown        = "TransferColdWalletWatcher shut down"
	LogTransferSuccess         = "TransferColdWalletWatcher: %s transfer successful - TxHash: %s, Amount: %s, Fee: %s"
	LogTransferSkippedNoExcess = "TransferColdWalletWatcher: %s transfer skipped - no excess liquidity"
	LogTransferSkippedNotEcon  = "TransferColdWalletWatcher: %s transfer skipped - not economical: %s"
	LogTransferSkippedCooldown = "TransferColdWalletWatcher: %s transfer skipped - liquidity target cooldown active"
	LogTransferFailed          = "TransferColdWalletWatcher: %s transfer failed - %s: %v"
)
