package pegin_test

import (
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test/integration/regtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Waited past the v2.6 quote's expiry, so the LPS has certainly closed its window.
	expiryMargin = 5 * time.Second
	// Longer deposit windows would make the test wait that long; the local stack uses 30 seconds.
	maxDepositWindow = 10 * time.Minute
)

// TestTwoPegInsToTheSameAddress checks that an institution registers its address once and keeps
// using it: a second peg-in, sent after a wait, pays the same deposit address with no new
// registration, and both settle. The registration and the LPS watch outlive the wait.
//
// As a control, the wait outlasts the deposit window of a v2.6 quote taken at its start: once it
// has passed, the LPS refuses that quote as expired, so a v2.6 user would need a new quote and a
// new deposit address, while the institution keeps its commit-first address.
//
//nolint:funlen // one function keeps the scenario's phases in order and lets them share its state
func TestTwoPegInsToTheSameAddress(t *testing.T) {
	s := newPegInStack(t)
	institution := s.rsk.NewAccount(t, new(big.Int)).Address
	// Registers the address the way the watchtower does; the contracts call it the registrant.
	fakeWatchtower := s.rsk.NewAccount(t, big.NewInt(1_000_000_000_000_000_000))
	amount := depositAmount(t, s.flyover)
	institutionBalanceBefore := s.rsk.Balance(t, institution)
	lpBalanceBefore := s.flyover.Balance(t, s.lp)
	fakeWatchtowerBalanceBefore := s.flyover.Balance(t, fakeWatchtower.Address)
	t.Logf("institution %s, fake watchtower %s, deposits of %s", institution, fakeWatchtower.Address, amount)

	var (
		first, second                   deposit
		firstRequested, secondRequested *regtest.PegInRequested
		firstResolved, secondResolved   *regtest.PegInResolved
		watch                           *rootstock.PegInWatch
	)

	runPhase(t, "the first peg-in registers the address and settles", func(t *testing.T) {
		first = s.payDeposit(t, institution, amount)
		s.flyover.RegisterAddress(t, fakeWatchtower, institution, first.proof)
		watch = s.waitWatchImported(t, institution)
		_, firstRequested = s.waitClaimed(t, institution, first)
		_, firstResolved = s.waitResolved(t, institution, first, firstRequested)
		assert.Equal(t, fakeWatchtower.Address, firstResolved.Registrant, "the first peg-in pays the registrant fee")
		assert.Positive(t, firstResolved.RegistrantFee.Sign())
	})

	runPhase(t, "the registration outlives a v2.6 quote's deposit window", func(t *testing.T) {
		quote := s.lps.PegInQuote(t, institution, new(big.Int).SetUint64(minQuoteValueWei))
		window := time.Duration(quote.Quote.TimeForDeposit) * time.Second
		require.LessOrEqualf(t, window, maxDepositWindow,
			"the LPS gives v2.6 quotes %s to deposit; lower timeForDeposit in docker-compose/local/lps/configurations/requests.js", window)
		waitUntil := regtest.PegInQuoteExpiry(quote).Add(expiryMargin)
		t.Logf("v2.6 deposit window %s, waiting until %s", window, waitUntil.Format(time.TimeOnly))

		s.driver.AdvanceUntil(t, "waiting", time.Until(waitUntil)+phaseTimeout, regtest.Step{},
			func(t *testing.T) (bool, string) {
				requireAddressStillWatched(t, s, institution, watch)
				return time.Now().After(waitUntil), "address registered and watch imported"
			})
		assert.Equal(t, first.address, s.flyover.PegInAddress(t, institution), "the deposit address does not change")

		expired := s.lps.AcceptPegInQuote(t, quote.QuoteHash)
		assert.Equal(t, http.StatusGone, expired.StatusCode, "the v2.6 quote's window is closed: %s", expired.Body)
		assert.Contains(t, expired.Body, "expired quote")
	})

	runPhase(t, "the second deposit reuses the address without a new registration", func(t *testing.T) {
		second = s.payDeposit(t, institution, amount)
		assert.Equal(t, first.address, second.address)
		var claim *rootstock.PegInClaim
		claim, secondRequested = s.waitClaimed(t, institution, second)
		assert.NotEqual(t, firstRequested.PegInId, secondRequested.PegInId, "each deposit is its own peg-in")
		assert.Equal(t, s.lp, secondRequested.Claimer)
		assert.Equal(t, second.amountWei, secondRequested.Amount)
		assert.Equal(t, second.address, claim.BtcAddress)
		assert.Equal(t,
			new(big.Int).Add(institutionBalanceBefore, new(big.Int).Add(firstRequested.NetToUser, secondRequested.NetToUser)),
			s.rsk.Balance(t, institution), "the institution receives both net amounts")
		assert.Len(t, s.flyover.AddressRegisteredEvents(t, institution), 1, "the address is registered once")
	})

	runPhase(t, "the second peg-in settles without a second registrant fee", func(t *testing.T) {
		_, secondResolved = s.waitResolved(t, institution, second, secondRequested)
		feeAtClaim := new(big.Int).Sub(secondRequested.Amount, secondRequested.NetToUser)
		// The registrant fee is paid once per address, so the LP keeps its whole fee.
		assert.Equal(t, common.Address{}, secondResolved.Registrant)
		assert.Zero(t, secondResolved.RegistrantFee.Sign())
		assert.Equal(t, new(big.Int).Add(secondRequested.NetToUser, feeAtClaim), secondResolved.ClaimerPayout,
			"fronted + fee")

		assert.Equal(t,
			new(big.Int).Add(lpBalanceBefore, new(big.Int).Add(firstResolved.ClaimerPayout, secondResolved.ClaimerPayout)),
			s.flyover.Balance(t, s.lp))
		assert.Equal(t, new(big.Int).Add(fakeWatchtowerBalanceBefore, firstResolved.RegistrantFee),
			s.flyover.Balance(t, fakeWatchtower.Address), "the fake watchtower is paid only for the first peg-in")
	})
}

// requireAddressStillWatched fails as soon as the registry forgets the address or the LPS
// changes its watch of it.
func requireAddressStillWatched(t *testing.T, s *pegInStack, institution common.Address, imported *rootstock.PegInWatch) {
	t.Helper()
	require.True(t, s.flyover.IsRegistered(t, institution), "the registry dropped the address")
	watch := s.db.Watch(t, institution)
	require.NotNil(t, watch, "the LPS deleted the watch")
	require.Equal(t, rootstock.PegInWatchImported, watch.State, "the LPS stopped watching the address")
	require.Equal(t, imported.BtcAddress, watch.BtcAddress)
	require.Equal(t, imported.TxHash, watch.TxHash, "the watch still comes from the same registration")
}
