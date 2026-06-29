package watcher

// Log message templates for transfer_cold_wallet_watcher.go
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

// Log message templates for bitcoin_peer_watcher.go
const (
	LogBitcoinPeerError    = "BitcoinPeerWatcher: error running peer check: %v"
	LogBitcoinPeerShutdown = "BitcoinPeerWatcher shut down"
)

// Log message templates for bitcoin_reorg_watcher.go
const (
	LogBitcoinReorgError    = "BitcoinReorgWatcher: error running reorg check: %v"
	LogBitcoinReorgShutdown = "BitcoinReorgWatcher shut down"
)

// Log message templates for rootstock_peer_watcher.go
const (
	LogRootstockPeerError    = "RootstockPeerWatcher: error running peer check: %v"
	LogRootstockPeerShutdown = "RootstockPeerWatcher shut down"
)

// Log message templates for rootstock_reorg_watcher.go
const (
	LogRootstockReorgError    = "RootstockReorgWatcher: error running reorg check: %v"
	LogRootstockReorgShutdown = "RootstockReorgWatcher shut down"
)

// Log message templates for eclipse.go
const (
	LogEclipseError    = "Error executing eclipse check: %v"
	LogEclipseShutdown = "EclipseWatcher shut down"
)

// Log message templates for penalization_alert.go
const (
	LogPenalizationError    = "Error checking penalization events inside watcher: %v"
	LogPenalizationShutdown = "PenalizationAlertWatcher shut down"
)

// Log message templates for quote_cleaner.go
const (
	LogQuoteCleanerShutdown = "QuoteCleanerWatcher shut down"
	LogQuoteCleanerError    = "Error cleaning quotes: %v"
	LogQuoteCleanerCleaned  = "Cleaned %d quotes:"
	LogQuoteCleanerQuote    = "Quote %s cleaned"
)

// Log message templates for liquidity_check.go
const (
	LogLiquidityCheckShutdown       = "LiquidityCheckWatcher shut down"
	LogLiquidityCheckError          = "Error checking liquidity inside watcher: %v"
	LogLiquidityCheckLowLiquidError = "Error checking low liquidity inside watcher: %v"
)

// Log message templates for btc_release_watcher.go
const (
	LogBtcReleaseError      = "Error checking BatchPegOuts in BtcReleaseWatcher: %v"
	LogBtcReleaseShutdown   = "BtcReleaseWatcher shut down"
	LogBtcReleaseChecking   = "Checking BatchPegOuts from block %d to %d, found %d batches"
	LogBtcReleaseProcessing = "Processing BatchPegOut: %+v"
	LogBtcReleaseNoQuotes   = "No PegOuts to process in batch (%s)"
	LogBtcReleaseProcessed  = "Successfully processed %d quotes in BatchPegOut (%s)"
)

// Log message templates for pegout_bridge_watcher.go
const (
	LogPegoutBridgePrefix         = "PegoutBridgeWatcher: "
	LogPegoutBridgeShutdown       = LogPegoutBridgePrefix + "shut down"
	LogPegoutBridgeGetQuotesError = LogPegoutBridgePrefix + "error getting pegout quotes: %v"
	LogPegoutBridgeSendError      = LogPegoutBridgePrefix + "error sending pegout to bridge: %v"
)

// Log message templates for pegout_btc_watcher.go
const (
	LogPegoutBtcPrefix         = "PegoutBtcTransferWatcher: "
	LogPegoutBtcShutdown       = LogPegoutBtcPrefix + "shut down"
	LogPegoutBtcBtcChainHeight = LogPegoutBtcPrefix + "error getting Bitcoin chain height: %v"
	LogPegoutBtcTxInfo         = LogPegoutBtcPrefix + "error getting Bitcoin transaction information (%s): %v"
	LogPegoutBtcRefundError    = LogPegoutBtcPrefix + "Error executing refund pegout on quote %s: %v"
	LogPegoutBtcWrongEvent     = LogPegoutBtcPrefix + "Trying to parse wrong event in Pegout Bridge watcher"
	LogPegoutBtcAlreadyWatched = LogPegoutBtcPrefix + "Quote %s is already watched"
	LogPegoutBtcNoTxHash       = LogPegoutBtcPrefix + "Quote %s doesn't have btc tx hash to watch"
)

// Log message templates for pegout_rsk_watcher.go
const (
	LogPegoutRskPrefix             = "PegoutRskDepositWatcher: "
	LogPegoutRskStart              = LogPegoutRskPrefix + "Starting to watch pegout deposits from block %d"
	LogPegoutRskChainHeight        = LogPegoutRskPrefix + "error getting Rootstock chain height: %v"
	LogPegoutRskShutdown           = LogPegoutRskPrefix + "shut down"
	LogPegoutRskWrongEvent         = LogPegoutRskPrefix + "Trying to parse wrong event in Pegout Rsk deposit watcher"
	LogPegoutRskAlreadyWatched     = LogPegoutRskPrefix + "Quote %s is already watched"
	LogPegoutRskGetDepositsError   = LogPegoutRskPrefix + "error executing getting deposits in range [%d, %d]: %v"
	LogPegoutRskCheckingDeposit    = LogPegoutRskPrefix + "Checking deposit of tx %s for quote %s"
	LogPegoutRskUpdateDepositError = LogPegoutRskPrefix + "Error updating pegout deposit quote (%s): %v"
	LogPegoutRskUpdateExpiredError = LogPegoutRskPrefix + "Error updating expired quote (%s): %v"
	LogPegoutRskExpired            = LogPegoutRskPrefix + "Quote %s expired at %d"
	LogPegoutRskReceiptError       = LogPegoutRskPrefix + "Error getting pegout deposit receipt of quote %s: %v"
	LogPegoutRskSendError          = LogPegoutRskPrefix + "Error sending pegout to the user (quote %s): %v"
	LogPegoutRskParseError         = LogPegoutRskPrefix + "Error parsing deposit event for quote %s: %v"
	// Fragments composed at runtime into the reject reason in logRejectReason (LogPegoutRskPrefix is prepended at log time).
	LogPegoutRskRejectReason  = "Rejecting quote %s for the following reason: "
	LogPegoutRskRejectExpired = "quote expired at %d, %d seconds before its first confirmation at %d;"
	LogPegoutRskRejectAmount  = "transaction amount %s is less than expected %s;"
)

// Log message templates for pegin_bridge_watcher.go
const (
	LogPeginBridgePrefix         = "Pegin Bridge watcher: "
	LogPeginBridgeShutdown       = LogPeginBridgePrefix + "shut down"
	LogPeginBridgeBtcChainHeight = LogPeginBridgePrefix + "error getting Bitcoin chain height: %v"
	LogPeginBridgeBtcTxInfo      = LogPeginBridgePrefix + "error getting Bitcoin transaction information (%s): %v"
	LogPeginBridgeWrongEvent     = LogPeginBridgePrefix + "Trying to parse wrong event"
	LogPeginBridgeAlreadyWatched = LogPeginBridgePrefix + "Quote %s is already watched"
	LogPeginBridgeRegisterError  = LogPeginBridgePrefix + "Error executing register pegin on quote %s: %v"
)

// Log message templates for pegin_btc_deposit_watcher.go
const (
	LogPeginBtcPrefix             = "PeginDepositAddressWatcher: "
	LogPeginBtcChainHeight        = LogPeginBtcPrefix + "error getting Bitcoin chain height: %v"
	LogPeginBtcShutdown           = LogPeginBtcPrefix + "shut down"
	LogPeginBtcWrongEvent         = LogPeginBtcPrefix + "trying to parse wrong event"
	LogPeginBtcAlreadyWatched     = LogPeginBtcPrefix + "Quote %s is already watched"
	LogPeginBtcImportAddress      = LogPeginBtcPrefix + "error while importing deposit address (%s): %v"
	LogPeginBtcCallForUserError   = LogPeginBtcPrefix + "Error executing call for user on quote %s: %v"
	LogPeginBtcUpdateExpiredError = LogPeginBtcPrefix + "Error updating expired quote (%s): %v"
	LogPeginBtcExpired            = LogPeginBtcPrefix + "Quote %s expired at %d"
	LogPeginBtcCheckingTx         = LogPeginBtcPrefix + "Checking transaction %s for quote %s"
	// Fragments composed at runtime into the reject reason in logRejectReason (LogPeginBtcPrefix is prepended at log time).
	LogPeginBtcRejectReason  = "Rejecting quote %s for the following reason: "
	LogPeginBtcRejectExpired = "quote expired at %d, %d seconds before its first confirmation at %d;"
	LogPeginBtcRejectAmount  = "transaction amount %s is less than expected %s;"
)
