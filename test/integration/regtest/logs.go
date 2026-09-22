package regtest

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/stretchr/testify/require"
)

const (
	lpsContainer = "lps01"
	lpsLogFile   = "/home/lps/logs/lps.log"
)

func LPSLogLen(t *testing.T) int {
	t.Helper()
	return len(readLPSLog(t))
}

func WaitSubmittingEmptyTxHashLog(t *testing.T, before int, user common.Address, depositTxID string) {
	t.Helper()
	snippet := pegin.LogPegInClaimSubmittingEmptyTxHash(user.Hex(), depositTxID)
	require.Eventually(t, func() bool {
		out, err := exec.Command("docker", "exec", lpsContainer, "cat", lpsLogFile).Output()
		if err != nil || len(out) < before {
			return false
		}
		return strings.Contains(string(out[before:]), snippet)
	}, WaitTimeout, PollInterval, "LPS did not log the empty-TxHash incident for %s / %s after the marked offset", user.Hex(), depositTxID)
}

func readLPSLog(t *testing.T) []byte {
	t.Helper()
	out, err := exec.Command("docker", "exec", lpsContainer, "cat", lpsLogFile).Output()
	require.NoError(t, err)
	return out
}
