package pegin_test

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/test/integration/regtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Bounds each attack loop; on the local stack each one ends after about 20 accepts.
	maxUnpaidAccepts = 100
)

// TestUnpaidQuotesDoNotBlockCommitFirstPegIns shows that the v2.6 liquidity DoS does not reach
// commit-first peg-ins.
//
// In the v2.6 quote flow, accepting a quote reserves its value until the deposit arrives or the
// quote expires, and accepting costs nothing. Accepting quotes and never paying them exhausts the
// liquidity the LPS offers, so honest quotes are refused; re-accepting as quotes expire keeps it
// exhausted. The 3.0 LPS still serves that flow, so the test runs the attack on the same LPS as a
// control, and keeps it going while a commit-first peg-in goes through.
//
// Commit-first has no step that reserves liquidity before payment: before paying, an institution
// can only read its deposit address. The LPS fronts RBTC only for a confirmed, registered deposit,
// and its claim checks the wallet balance, not the reserved quotes. Draining the LP with real
// deposits is not covered here: it costs the attacker the same amount in BTC plus the LP fee.
//
// The quotes accepted here stay reserved after they expire, until the LPS's deposit watcher sees
// the next BTC block (up to a minute later). On regtest blocks only come when a test mines them, so
// v2.6 quotes on the stack are refused until then; the next run of this test waits for it.
func TestUnpaidQuotesDoNotBlockCommitFirstPegIns(t *testing.T) {
	s := newPegInStack(t)
	attack := &quoteAttack{lps: s.lps, attacker: s.rsk.NewAccount(t, new(big.Int)).Address}

	runPhase(t, "quotes from an earlier attack have been released", func(t *testing.T) {
		// Left over from an earlier run, they would hold the lock this run's control must create.
		// The LPS checks quotes for expiry only when it sees a new BTC block, so each round mines one.
		s.driver.AdvanceUntil(t, "earlier quotes released", phaseTimeout, regtest.Step{BtcBlocks: 1},
			func(t *testing.T) (bool, string) {
				reserved := reservedQuoteLiquidity(t, s)
				return reserved.Sign() == 0, fmt.Sprintf("%s wei reserved by unpaid quotes", reserved)
			})
	})

	lpWalletBefore := s.rsk.Balance(t, s.lp)
	lpContractBalanceBefore := s.flyover.Balance(t, s.lp)

	runPhase(t, "unpaid v2.6 quotes lock the quote liquidity", func(t *testing.T) {
		accepted, available := attack.hold(t)
		t.Logf("[attack] %d unpaid quotes accepted, %s RBTC available (smallest quote %s RBTC)",
			accepted, rbtc(available), rbtc(smallestQuote()))
		require.Positive(t, accepted, "this attack's own quotes take the liquidity")
		assert.Equal(t, lpWalletBefore, s.rsk.Balance(t, s.lp), "accepting a quote moves no RBTC")
		assert.Equal(t, lpContractBalanceBefore, s.flyover.Balance(t, s.lp))
	})

	runPhase(t, "a commit-first peg-in settles while the attack goes on", func(t *testing.T) {
		// Every round of every wait below holds the lock first and fails if it cannot, so each
		// observation of the peg-in is made while the attack holds.
		attacked := *s
		attacked.driver = s.driver.WithEachRound(attack.holdEachRound)
		institution := s.rsk.NewAccount(t, new(big.Int)).Address
		// Registers the address the way the watchtower does; the contracts call it the registrant.
		fakeWatchtower := s.rsk.NewAccount(t, big.NewInt(1_000_000_000_000_000_000))
		institutionBalanceBefore := s.rsk.Balance(t, institution)

		paid := attacked.payDeposit(t, institution, depositAmount(t, s.flyover))
		s.flyover.RegisterAddress(t, fakeWatchtower, institution, paid.proof)
		attacked.waitWatchImported(t, institution)

		_, requested := attacked.waitClaimed(t, institution, paid)
		assert.Equal(t, s.lp, requested.Claimer)
		assert.Equal(t, new(big.Int).Add(institutionBalanceBefore, requested.NetToUser), s.rsk.Balance(t, institution),
			"the institution is fronted")

		_, resolved := attacked.waitResolved(t, institution, paid, requested)
		assert.Equal(t, new(big.Int).Add(lpContractBalanceBefore, resolved.ClaimerPayout), s.flyover.Balance(t, s.lp),
			"the LP gets its advance back")

		t.Logf("[attack] the lock held on all %d checks while the peg-in went through; %d quotes re-accepted as others expired",
			attack.checks, attack.reaccepted)
		require.Positive(t, attack.checks)
	})

	runPhase(t, "the attack still blocks honest v2.6 quotes after the peg-in settled", func(t *testing.T) {
		_, available := attack.hold(t)
		t.Logf("[attack] still running: %s RBTC available (smallest quote %s RBTC)", rbtc(available), rbtc(smallestQuote()))
		honestUser := s.rsk.NewAccount(t, new(big.Int)).Address
		quote := s.lps.PegInQuote(t, honestUser, smallestQuote())
		refused := s.lps.AcceptPegInQuote(t, quote.QuoteHash)
		assert.Equal(t, http.StatusConflict, refused.StatusCode, "an honest quote is refused: %s", refused.Body)
		assert.Contains(t, refused.Body, "not enough liquidity")
	})
}

// quoteAttack keeps the v2.6 quote liquidity locked with unpaid quotes, and counts how often it
// confirmed the lock and how many quotes it had to accept again as earlier ones expired.
type quoteAttack struct {
	lps        *regtest.Lps
	attacker   common.Address
	checks     int
	reaccepted int
}

// hold accepts unpaid quotes, the largest first and then the smallest, until the LPS refuses
// them, and fails unless no quote the LP allows fits anymore. While earlier quotes are still
// reserved it only tops the lock up. It returns how many quotes it accepted and the liquidity
// left.
func (attack *quoteAttack) hold(t *testing.T) (int, *big.Int) {
	t.Helper()
	total := 0
	for _, value := range []*big.Int{largestQuote(), smallestQuote()} {
		accepted, refused := acceptUnpaidQuotesUntilRefused(t, attack.lps, attack.attacker, value)
		require.Equal(t, http.StatusConflict, refused.StatusCode, "an accept is refused: %s", refused.Body)
		require.Contains(t, refused.Body, "not enough liquidity")
		total += accepted
	}
	available := attack.lps.AvailablePegInLiquidity(t)
	require.Negative(t, available.Cmp(smallestQuote()),
		"the liquidity left (%s RBTC) cannot cover even the smallest quote", rbtc(available))
	return total, available
}

// holdEachRound is the driver action that keeps the attack going while the stack advances.
func (attack *quoteAttack) holdEachRound(t *testing.T) {
	t.Helper()
	accepted, available := attack.hold(t)
	attack.checks++
	attack.reaccepted += accepted
	t.Logf("[attack] check %d: lock holds, %s RBTC available (smallest quote %s RBTC), %d quotes re-accepted",
		attack.checks, rbtc(available), rbtc(smallestQuote()), accepted)
}

func largestQuote() *big.Int {
	return new(big.Int).SetUint64(maxQuoteValueWei)
}

func smallestQuote() *big.Int {
	return new(big.Int).SetUint64(minQuoteValueWei)
}

// rbtc formats a wei amount in RBTC with four decimals, for the logs.
func rbtc(wei *big.Int) string {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18)).Text('f', 4)
}

// reservedQuoteLiquidity is what the LPS holds back for accepted v2.6 quotes: the LP's RBTC, in its
// wallet and in PegInContract, minus the liquidity the LPS still offers.
func reservedQuoteLiquidity(t *testing.T, s *pegInStack) *big.Int {
	t.Helper()
	total := new(big.Int).Add(s.rsk.Balance(t, s.lp), s.flyover.Balance(t, s.lp))
	return total.Sub(total, s.lps.AvailablePegInLiquidity(t))
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
