package pegin

import "fmt"

func LogPegInClaimEmptyTxHash(rskAddress, depositTxID string) string {
	return fmt.Sprintf(
		"PegInClaimWatcher: submitting claim %s/%s has empty TxHash; follow incident-recovery; not resubmitting",
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

func LogPegInClaimSubmittingEmptyTxHash(rskAddress, depositTxID string) string {
	return fmt.Sprintf(
		"PegInClaim: submitting claim %s/%s has empty TxHash; follow incident-recovery; not resubmitting",
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
