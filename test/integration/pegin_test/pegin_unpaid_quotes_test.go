package pegin_test

import (
	"math/big"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/test/integration/regtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// The largest and smallest v2.6 peg-in quotes the local stack allows (maxValue and minValue in
	// docker-compose/local/lps/configurations/requests.js).
	maxQuoteValueWei = 10_000_000_000_000_000_000
	minQuoteValueWei = 600_000_000_000_000_000
	// Bounds each attack loop; on the local stack each one ends after about 20 accepts.
	maxUnpaidAccepts      = 100
	addressesAskedForFree = 20
)

// TestUnpaidQuotesDoNotBlockCommitFirstPegIns shows that the v2.6 liquidity DoS does not reach
// commit-first peg-ins.
//
// In the v2.6 quote flow, accepting a quote reserves its value until the deposit arrives or the
// quote expires (an hour on the local stack), and accepting costs nothing. Accepting quotes and
// never paying them exhausts the liquidity the LPS offers, so honest quotes are refused. The
// 3.0 LPS still serves that flow, so the test reproduces the attack on the same LPS as a control.
//
// Commit-first has no step that reserves liquidity before payment: before paying, an institution
// can only read its deposit address. The LPS fronts RBTC only for a confirmed, registered deposit,
// and its claim checks the wallet balance, not the reserved quotes. Draining the LP with real
// deposits is not covered here: it costs the attacker the same amount in BTC plus the LP fee.
//
// The quotes accepted here stay reserved for an hour, so v2.6 quotes on the stack are refused
// until they expire or the stack is reset. Commit-first peg-ins are not affected, which is the
// point of the test.
func TestUnpaidQuotesDoNotBlockCommitFirstPegIns(t *testing.T) {
	s := newPegInStack(t)
	largestQuote := new(big.Int).SetUint64(maxQuoteValueWei)
	smallestQuote := new(big.Int).SetUint64(minQuoteValueWei)
	lpWalletBefore := s.rsk.Balance(t, s.lp)
	lpContractBalanceBefore := s.flyover.Balance(t, s.lp)

	runPhase(t, "asking for commit-first deposit addresses reserves nothing", func(t *testing.T) {
		availableBefore := s.lps.AvailablePegInLiquidity(t)
		for range addressesAskedForFree {
			institution := s.rsk.NewAccount(t, new(big.Int)).Address
			require.NotEmpty(t, s.flyover.PegInAddress(t, institution))
			assert.Nil(t, s.db.Watch(t, institution), "the LPS does not learn about an address nobody paid")
		}
		assert.Equal(t, availableBefore, s.lps.AvailablePegInLiquidity(t))
		assert.Equal(t, lpWalletBefore, s.rsk.Balance(t, s.lp))
	})

	runPhase(t, "unpaid v2.6 quotes lock the quote liquidity", func(t *testing.T) {
		attacker := s.rsk.NewAccount(t, new(big.Int)).Address
		// The largest quotes take most of the liquidity, then the smallest ones take what is left,
		// so no quote the LP allows fits anymore.
		for _, value := range []*big.Int{largestQuote, smallestQuote} {
			accepted, rejected := acceptUnpaidQuotesUntilRefused(t, s.lps, attacker, value)
			t.Logf("%d unpaid quotes of %s wei accepted; quote liquidity left %s wei",
				accepted, value, s.lps.AvailablePegInLiquidity(t))
			require.Equal(t, http.StatusConflict, rejected.StatusCode, "an accept is refused: %s", rejected.Body)
			assert.Contains(t, rejected.Body, "not enough liquidity")
		}
		assert.Less(t, s.lps.AvailablePegInLiquidity(t).Cmp(smallestQuote), 0,
			"the liquidity left cannot cover even the smallest quote")
		assert.Equal(t, lpWalletBefore, s.rsk.Balance(t, s.lp), "accepting a quote moves no RBTC")
		assert.Equal(t, lpContractBalanceBefore, s.flyover.Balance(t, s.lp))
	})

	runPhase(t, "a commit-first peg-in settles while the quote liquidity is locked", func(t *testing.T) {
		institution := s.rsk.NewAccount(t, new(big.Int)).Address
		// Registers the address the way the watchtower does; the contracts call it the registrant.
		fakeWatchtower := s.rsk.NewAccount(t, big.NewInt(1_000_000_000_000_000_000))
		institutionBalanceBefore := s.rsk.Balance(t, institution)

		paid := s.payDeposit(t, institution, depositAmount(t, s.flyover))
		s.flyover.RegisterAddress(t, fakeWatchtower, institution, paid.proof)
		s.waitWatchImported(t, institution)

		_, requested := s.waitClaimed(t, institution, paid)
		assert.Equal(t, s.lp, requested.Claimer)
		assert.Equal(t, new(big.Int).Add(institutionBalanceBefore, requested.NetToUser), s.rsk.Balance(t, institution),
			"the institution is fronted")
		assert.Less(t, s.lps.AvailablePegInLiquidity(t).Cmp(smallestQuote), 0,
			"the quote liquidity was still locked when the LP fronted")

		_, resolved := s.waitResolved(t, institution, paid, requested)
		assert.Equal(t, new(big.Int).Add(lpContractBalanceBefore, resolved.ClaimerPayout), s.flyover.Balance(t, s.lp),
			"the LP gets its advance back")
	})
}

// acceptUnpaidQuotesUntilRefused accepts quotes of value and never pays them, until the LPS
// refuses one. It returns how many it accepted and the refusal.
func acceptUnpaidQuotesUntilRefused(
	t *testing.T,
	lps *regtest.Lps,
	attacker common.Address,
	value *big.Int,
) (int, regtest.AcceptResult) {
	t.Helper()
	for accepted := range maxUnpaidAccepts {
		quote := lps.PegInQuote(t, attacker, value)
		if result := lps.AcceptPegInQuote(t, quote.QuoteHash); result.StatusCode != http.StatusOK {
			return accepted, result
		}
	}
	require.FailNowf(t, "the quote liquidity never ran out", "after %d accepts of %s wei", maxUnpaidAccepts, value)
	return 0, regtest.AcceptResult{}
}
