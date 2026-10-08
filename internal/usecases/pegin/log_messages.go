package pegin

import (
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
)

func LogPegInClaimEmptyRequestTxHash(rskAddress, depositTxID string) string {
	return fmt.Sprintf(
		"PegInClaimWatcher: submitting claim %s/%s has empty RequestTxHash; follow incident-recovery; not resubmitting",
		rskAddress,
		depositTxID,
	)
}

func LogPegInClaimMissingReceipt(txHash, rskAddress, depositTxID string) string {
	return fmt.Sprintf(
		"PegInClaimWatcher: submitting tx %s for %s/%s has no recoverable receipt; follow incident-recovery; not resubmitting",
		txHash,
		rskAddress,
		depositTxID,
	)
}

func LogPegInClaimSubmittingEmptyRequestTxHash(rskAddress, depositTxID string) string {
	return fmt.Sprintf(
		"PegInClaim: submitting claim %s/%s has empty RequestTxHash; follow incident-recovery; not resubmitting",
		rskAddress,
		depositTxID,
	)
}

func LogPegInClaimInsufficientWalletLiquidity(rskAddress, depositTxID string, available, required fmt.Stringer) string {
	return fmt.Sprintf(
		"PegInClaim: not enough wallet liquidity for %s/%s (available %s wei, required %s wei); will retry",
		rskAddress,
		depositTxID,
		available,
		required,
	)
}

func LogPegInClaimBelowMinimum(rskAddress, depositTxID string) string {
	return fmt.Sprintf(
		"PegInClaim: deposit %s/%s is below the Flyover minimum; not claiming",
		rskAddress,
		depositTxID,
	)
}

func LogPegInClaimMissingEvent(txHash, rskAddress, depositTxID string, err error) string {
	return fmt.Sprintf(
		"PegInClaimWatcher: receipt %s for %s/%s is missing PegInRequested; follow incident-recovery; not resubmitting: %v",
		txHash,
		rskAddress,
		depositTxID,
		err,
	)
}

func LogPegInResolved(claim rootstock.PegInClaim) string {
	return fmt.Sprintf(
		"PegInResolve: resolved %s/%s (pegInId %s, resolveTxHash %s, claimerPayout %v wei, registrantFee %v wei)",
		claim.RskAddress,
		claim.DepositTxID,
		claim.PegInID,
		claim.ResolveTxHash,
		claim.ClaimerPayout,
		claim.RegistrantFee,
	)
}

func LogPegInResolveAlreadyProcessed(claim rootstock.PegInClaim) string {
	return fmt.Sprintf(
		"PegInResolve: %s/%s (pegInId %s) was already resolved on-chain; marking it resolved",
		claim.RskAddress,
		claim.DepositTxID,
		claim.PegInID,
	)
}

func LogPegInResolveWaitingForBridge(claim rootstock.PegInClaim) string {
	return fmt.Sprintf(
		"PegInResolve: bridge rejected %s/%s (pegInId %s) with code -303 after the confirmation check passed; will retry",
		claim.RskAddress,
		claim.DepositTxID,
		claim.PegInID,
	)
}

func LogPegInResolveFailed(claim rootstock.PegInClaim, err error) string {
	return fmt.Sprintf(
		"PegInResolve: resolve failed for %s/%s (pegInId %s); not retrying: %v",
		claim.RskAddress,
		claim.DepositTxID,
		claim.PegInID,
		err,
	)
}

func LogPegInResolveTxNotResolved(claim rootstock.PegInClaim, err error) string {
	return fmt.Sprintf(
		"PegInResolve: resolve tx %s for %s/%s (pegInId %s) did not resolve the peg-in; will retry: %v",
		claim.ResolveTxHash,
		claim.RskAddress,
		claim.DepositTxID,
		claim.PegInID,
		err,
	)
}
