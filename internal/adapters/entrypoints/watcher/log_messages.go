package watcher

import (
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
)

// Log message templates for transfer_cold_wallet_watcher.go
// Messages whose only placeholder is the trailing error are kept as format
// constants used with logrus *f variants. Messages that require values are
// exposed as typed functions instead, so the compiler enforces the arguments
// the caller must provide.
// TODO FLY-2388: this is a pilot of the repo-wide log message centralization.
// Apply the same pattern (log_messages.go per package) to the remaining packages.
const (
	LogTransferError    = "TransferColdWalletWatcher: Error executing transfer to cold wallet: %v"
	LogTransferShutdown = "TransferColdWalletWatcher shut down"
)

func LogTransferSuccess(network, txHash, amount, fee string) string {
	return fmt.Sprintf("TransferColdWalletWatcher: %s transfer successful - TxHash: %s, Amount: %s, Fee: %s", network, txHash, amount, fee)
}

func LogTransferSkippedNoExcess(network string) string {
	return fmt.Sprintf("TransferColdWalletWatcher: %s transfer skipped - no excess liquidity", network)
}

func LogTransferSkippedNotEcon(network, message string) string {
	return fmt.Sprintf("TransferColdWalletWatcher: %s transfer skipped - not economical: %s", network, message)
}

func LogTransferSkippedCooldown(network string) string {
	return fmt.Sprintf("TransferColdWalletWatcher: %s transfer skipped - liquidity target cooldown active", network)
}

func LogTransferFailed(network, message string, err error) string {
	return fmt.Sprintf("TransferColdWalletWatcher: %s transfer failed - %s: %v", network, message, err)
}

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
)

func LogQuoteCleanerCleaned(count int) string {
	return fmt.Sprintf("Cleaned %d quotes:", count)
}

func LogQuoteCleanerQuote(quoteHash string) string {
	return fmt.Sprintf("Quote %s cleaned", quoteHash)
}

// Log message templates for liquidity_check.go
const (
	LogLiquidityCheckShutdown       = "LiquidityCheckWatcher shut down"
	LogLiquidityCheckError          = "Error checking liquidity inside watcher: %v"
	LogLiquidityCheckLowLiquidError = "Error checking low liquidity inside watcher: %v"
)

// Log message templates for btc_release_watcher.go
const (
	LogBtcReleaseError    = "Error checking BatchPegOuts in BtcReleaseWatcher: %v"
	LogBtcReleaseShutdown = "BtcReleaseWatcher shut down"
)

func LogBtcReleaseChecking(fromBlock, toBlock uint64, batchCount int) string {
	return fmt.Sprintf("Checking BatchPegOuts from block %d to %d, found %d batches", fromBlock, toBlock, batchCount)
}

func LogBtcReleaseProcessing(batch rootstock.BatchPegOut) string {
	return fmt.Sprintf("Processing BatchPegOut: %+v", batch)
}

func LogBtcReleaseNoQuotes(batchTxHash string) string {
	return fmt.Sprintf("No PegOuts to process in batch (%s)", batchTxHash)
}

func LogBtcReleaseProcessed(quoteCount uint, batchTxHash string) string {
	return fmt.Sprintf("Successfully processed %d quotes in BatchPegOut (%s)", quoteCount, batchTxHash)
}

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
	LogPegoutBtcWrongEvent     = LogPegoutBtcPrefix + "Trying to parse wrong event in Pegout Bridge watcher"
)

func LogPegoutBtcAlreadyWatched(quoteHash string) string {
	return fmt.Sprintf(LogPegoutBtcPrefix+"Quote %s is already watched", quoteHash)
}

func LogPegoutBtcNoTxHash(quoteHash string) string {
	return fmt.Sprintf(LogPegoutBtcPrefix+"Quote %s doesn't have btc tx hash to watch", quoteHash)
}

func LogPegoutBtcTxInfo(txHash string, err error) string {
	return fmt.Sprintf(LogPegoutBtcPrefix+"error getting Bitcoin transaction information (%s): %v", txHash, err)
}

func LogPegoutBtcRefundError(quoteHash string, err error) string {
	return fmt.Sprintf(LogPegoutBtcPrefix+"Error executing refund pegout on quote %s: %v", quoteHash, err)
}

// Log message templates for pegout_rsk_watcher.go
const (
	LogPegoutRskPrefix      = "PegoutRskDepositWatcher: "
	LogPegoutRskChainHeight = LogPegoutRskPrefix + "error getting Rootstock chain height: %v"
	LogPegoutRskShutdown    = LogPegoutRskPrefix + "shut down"
	LogPegoutRskWrongEvent  = LogPegoutRskPrefix + "Trying to parse wrong event in Pegout Rsk deposit watcher"
)

func LogPegoutRskStart(fromBlock uint64) string {
	return fmt.Sprintf(LogPegoutRskPrefix+"Starting to watch pegout deposits from block %d", fromBlock)
}

func LogPegoutRskAlreadyWatched(quoteHash string) string {
	return fmt.Sprintf(LogPegoutRskPrefix+"Quote %s is already watched", quoteHash)
}

func LogPegoutRskCheckingDeposit(txHash, quoteHash string) string {
	return fmt.Sprintf(LogPegoutRskPrefix+"Checking deposit of tx %s for quote %s", txHash, quoteHash)
}

func LogPegoutRskExpired(quoteHash string, expirationTime int64) string {
	return fmt.Sprintf(LogPegoutRskPrefix+"Quote %s expired at %d", quoteHash, expirationTime)
}

// Fragments composed at runtime into the reject reason in logRejectReason (LogPegoutRskPrefix is prepended at log time).
func LogPegoutRskRejectReason(quoteHash string) string {
	return fmt.Sprintf("Rejecting quote %s for the following reason: ", quoteHash)
}

func LogPegoutRskRejectExpired(expirationTime, secondsLate, confirmationTime int64) string {
	return fmt.Sprintf("quote expired at %d, %d seconds before its first confirmation at %d;", expirationTime, secondsLate, confirmationTime)
}

func LogPegoutRskRejectAmount(paidAmount, expectedAmount string) string {
	return fmt.Sprintf("transaction amount %s is less than expected %s;", paidAmount, expectedAmount)
}

func LogPegoutRskGetDepositsError(fromBlock, toBlock uint64, err error) string {
	return fmt.Sprintf(LogPegoutRskPrefix+"error executing getting deposits in range [%d, %d]: %v", fromBlock, toBlock, err)
}

func LogPegoutRskUpdateDepositError(quoteHash string, err error) string {
	return fmt.Sprintf(LogPegoutRskPrefix+"Error updating pegout deposit quote (%s): %v", quoteHash, err)
}

func LogPegoutRskUpdateExpiredError(quoteHash string, err error) string {
	return fmt.Sprintf(LogPegoutRskPrefix+"Error updating expired quote (%s): %v", quoteHash, err)
}

func LogPegoutRskReceiptError(quoteHash string, err error) string {
	return fmt.Sprintf(LogPegoutRskPrefix+"Error getting pegout deposit receipt of quote %s: %v", quoteHash, err)
}

func LogPegoutRskSendError(quoteHash string, err error) string {
	return fmt.Sprintf(LogPegoutRskPrefix+"Error sending pegout to the user (quote %s): %v", quoteHash, err)
}

func LogPegoutRskParseError(quoteHash string, err error) string {
	return fmt.Sprintf(LogPegoutRskPrefix+"Error parsing deposit event for quote %s: %v", quoteHash, err)
}

// Log message templates for pegin_bridge_watcher.go
const (
	LogPeginBridgePrefix         = "Pegin Bridge watcher: "
	LogPeginBridgeShutdown       = LogPeginBridgePrefix + "shut down"
	LogPeginBridgeBtcChainHeight = LogPeginBridgePrefix + "error getting Bitcoin chain height: %v"
	LogPeginBridgeWrongEvent     = LogPeginBridgePrefix + "Trying to parse wrong event"
)

func LogPeginBridgeAlreadyWatched(quoteHash string) string {
	return fmt.Sprintf(LogPeginBridgePrefix+"Quote %s is already watched", quoteHash)
}

func LogPeginBridgeBtcTxInfo(txHash string, err error) string {
	return fmt.Sprintf(LogPeginBridgePrefix+"error getting Bitcoin transaction information (%s): %v", txHash, err)
}

func LogPeginBridgeRegisterError(quoteHash string, err error) string {
	return fmt.Sprintf(LogPeginBridgePrefix+"Error executing register pegin on quote %s: %v", quoteHash, err)
}

// Log message templates for pegin_btc_deposit_watcher.go
const (
	LogPeginBtcPrefix      = "PeginDepositAddressWatcher: "
	LogPeginBtcChainHeight = LogPeginBtcPrefix + "error getting Bitcoin chain height: %v"
	LogPeginBtcShutdown    = LogPeginBtcPrefix + "shut down"
	LogPeginBtcWrongEvent  = LogPeginBtcPrefix + "trying to parse wrong event"
)

func LogPeginBtcAlreadyWatched(quoteHash string) string {
	return fmt.Sprintf(LogPeginBtcPrefix+"Quote %s is already watched", quoteHash)
}

func LogPeginBtcExpired(quoteHash string, expirationTime int64) string {
	return fmt.Sprintf(LogPeginBtcPrefix+"Quote %s expired at %d", quoteHash, expirationTime)
}

func LogPeginBtcCheckingTx(txHash, quoteHash string) string {
	return fmt.Sprintf(LogPeginBtcPrefix+"Checking transaction %s for quote %s", txHash, quoteHash)
}

// Fragments composed at runtime into the reject reason in logRejectReason (LogPeginBtcPrefix is prepended at log time).
func LogPeginBtcRejectReason(quoteHash string) string {
	return fmt.Sprintf("Rejecting quote %s for the following reason: ", quoteHash)
}

func LogPeginBtcRejectExpired(expirationTime, secondsLate, blockTime int64) string {
	return fmt.Sprintf("quote expired at %d, %d seconds before its first confirmation at %d;", expirationTime, secondsLate, blockTime)
}

func LogPeginBtcRejectAmount(paidAmount, expectedAmount string) string {
	return fmt.Sprintf("transaction amount %s is less than expected %s;", paidAmount, expectedAmount)
}

func LogPeginBtcImportAddress(depositAddress string, err error) string {
	return fmt.Sprintf(LogPeginBtcPrefix+"error while importing deposit address (%s): %v", depositAddress, err)
}

func LogPeginBtcCallForUserError(quoteHash string, err error) string {
	return fmt.Sprintf(LogPeginBtcPrefix+"Error executing call for user on quote %s: %v", quoteHash, err)
}

func LogPeginBtcUpdateExpiredError(quoteHash string, err error) string {
	return fmt.Sprintf(LogPeginBtcPrefix+"Error updating expired quote (%s): %v", quoteHash, err)
}

const (
	LogPegoutEscrowPrefix   = "PegoutEscrowWatcher: "
	LogPegoutEscrowError    = LogPegoutEscrowPrefix + "error checking escrowed peg-outs: %v"
	LogPegoutEscrowShutdown = LogPegoutEscrowPrefix + "shut down"
)

func LogPegoutEscrowStart(fromBlock uint64) string {
	return fmt.Sprintf(LogPegoutEscrowPrefix+"watching escrowed peg-outs from block %d", fromBlock)
}

func LogPegoutEscrowChecking(fromBlock, toBlock uint64, requestedCount int) string {
	return fmt.Sprintf(LogPegoutEscrowPrefix+"checking escrowed peg-outs from block %d to %d, found %d requests", fromBlock, toBlock, requestedCount)
}

func LogPegoutEscrowStateError(requestHash string, err error) string {
	return fmt.Sprintf(LogPegoutEscrowPrefix+"error reading peg-out state for %s: %v", requestHash, err)
}

func LogPegoutEscrowClaimError(requestHash string, err error) string {
	return fmt.Sprintf(LogPegoutEscrowPrefix+"error claiming peg-out %s: %v", requestHash, err)
}
