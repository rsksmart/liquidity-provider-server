package pegin

import "fmt"

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
