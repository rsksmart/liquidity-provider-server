package pegout

import (
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

const (
	LogClaimPegoutPrefix = "ClaimPegOut: "
)

func LogClaimPegoutLostRace(requestHash string) string {
	return fmt.Sprintf(LogClaimPegoutPrefix+"lost race for %s", requestHash)
}

func LogClaimPegoutWindowClosed(requestHash string) string {
	return fmt.Sprintf(LogClaimPegoutPrefix+"claim window closed for %s", requestHash)
}

func LogClaimPegoutCapacitySkip(requestHash string) string {
	return fmt.Sprintf(LogClaimPegoutPrefix+"capacity gate failed for %s", requestHash)
}

func LogClaimPegoutProfitabilitySkip(requestHash string) string {
	return fmt.Sprintf(LogClaimPegoutPrefix+"profitability gate failed for %s", requestHash)
}

func LogClaimPegoutRestrictedSkip(requestHash string, until uint64) string {
	return fmt.Sprintf(LogClaimPegoutPrefix+"LP restricted until %d, skipping %s", until, requestHash)
}

func LogClaimPegoutAlreadyClaimed(requestHash string) string {
	return fmt.Sprintf(LogClaimPegoutPrefix+"already claimed locally %s", requestHash)
}

func LogClaimPegoutSuccess(requestHash, quoteHash, txHash string) string {
	return fmt.Sprintf(LogClaimPegoutPrefix+"claimed %s as %s in tx %s", requestHash, quoteHash, txHash)
}

func LogClaimPegoutPendingPromoted(quoteHash string) string {
	return fmt.Sprintf(LogClaimPegoutPrefix+"pending claim %s was mined, promoted to claimed", quoteHash)
}

func LogClaimPegoutPendingDropped(quoteHash string) string {
	return fmt.Sprintf(LogClaimPegoutPrefix+"pending claim %s can no longer be mined, dropped", quoteHash)
}

func LogClaimPegoutPendingSettled(quoteHash string, state blockchain.EscrowedPegOutState) string {
	return fmt.Sprintf(LogClaimPegoutPrefix+"pending claim %s was mined but the escrow already settled it (state %d), dropped", quoteHash, state)
}
