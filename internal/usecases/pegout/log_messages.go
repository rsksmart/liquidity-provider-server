package pegout

import "fmt"

const (
	LogClaimPegoutPrefix = "ClaimPegOut: "
)

func LogClaimPegoutLostRace(requestHash string) string {
	return fmt.Sprintf(LogClaimPegoutPrefix+"lost race for %s", requestHash)
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

func LogClaimPegoutSuccess(requestHash, txHash string) string {
	return fmt.Sprintf(LogClaimPegoutPrefix+"claimed %s in tx %s", requestHash, txHash)
}
