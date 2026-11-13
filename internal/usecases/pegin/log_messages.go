package pegin

import (
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
)

func LogPegInClaimAlreadyStored(rskAddress, depositTxID string, state rootstock.PegInClaimState) string {
	return fmt.Sprintf("PegIn claim of deposit %s for %s is already %s", depositTxID, rskAddress, state)
}

func LogPegInClaimWaitingForConfirmations(rskAddress, depositTxID string, confirmations, required uint64) string {
	return fmt.Sprintf("PegIn deposit %s for %s has %d of %d confirmations", depositTxID, rskAddress, confirmations, required)
}

func LogPegInClaimHardPaused(rskAddress, depositTxID string) string {
	return fmt.Sprintf("Skipping the claim of PegIn deposit %s for %s: the PegIn contract is hard paused", depositTxID, rskAddress)
}

func LogPegInClaimBelowMinimum(rskAddress, depositTxID string, amount, minimum *entities.Wei) string {
	return fmt.Sprintf("PegIn deposit %s for %s of %v wei is below the minimum of %v wei", depositTxID, rskAddress, amount, minimum)
}

func LogPegInClaimDeductionAboveAmount(rskAddress, depositTxID string, amount, deduction *entities.Wei) string {
	return fmt.Sprintf("PegIn deposit %s for %s of %v wei does not cover the fee and gas fee of %v wei", depositTxID, rskAddress, amount, deduction)
}

func LogPegInClaimNotClaimable(rskAddress, depositTxID string, reason error) string {
	return fmt.Sprintf("PegIn deposit %s for %s can't be claimed yet: %v", depositTxID, rskAddress, reason)
}

func LogPegInClaimInsufficientBalance(rskAddress, depositTxID string, balance, required *entities.Wei) string {
	return fmt.Sprintf("Wallet balance of %v wei can't claim PegIn deposit %s for %s, which needs %v wei", balance, depositTxID, rskAddress, required)
}

func LogPegInClaimRaceLost(rskAddress, depositTxID string) string {
	return fmt.Sprintf("PegIn deposit %s for %s was claimed by another liquidity provider", depositTxID, rskAddress)
}

func LogPegInClaimStoredMeanwhile(rskAddress, depositTxID string) string {
	return fmt.Sprintf("PegIn claim of deposit %s for %s was stored by another claim attempt", depositTxID, rskAddress)
}

func LogPegInClaimed(rskAddress, depositTxID, txHash, pegInID string) string {
	return fmt.Sprintf("Claimed PegIn deposit %s for %s in transaction %s, pegInId %s", depositTxID, rskAddress, txHash, pegInID)
}

func LogPegInClaimIdMismatch(depositTxID, storedPegInID, eventPegInID string) string {
	return fmt.Sprintf("PegIn claim of deposit %s stored pegInId %s but the PegInRequested event has %s", depositTxID, storedPegInID, eventPegInID)
}

func LogPegInClaimUserPaymentFailed(rskAddress, depositTxID, txHash string) string {
	return fmt.Sprintf("PegIn claim of deposit %s in transaction %s could not pay %s; the fronted amount went to the PegIn contract balance", depositTxID, txHash, rskAddress)
}
