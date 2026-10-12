package pegin_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test/integration/regtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// How long the institution waits between its two peg-ins. Nothing is mined meanwhile: only
// time passes, while the LPS watchers keep running.
const waitBetweenPegIns = 30 * time.Second

// TestTwoPegInsToTheSameAddress checks that an institution registers its address once and keeps
// using it: a second peg-in, sent after a wait, pays the same deposit address with no new
// registration, and both settle. The registration and the LPS watch outlive the wait.
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

	runPhase(t, "the registration outlives the wait", func(t *testing.T) {
		start := time.Now()
		s.driver.AdvanceUntil(t, "waiting", waitBetweenPegIns+phaseTimeout, regtest.Step{},
			func(t *testing.T) (bool, string) {
				requireAddressStillWatched(t, s, institution, watch)
				return time.Since(start) >= waitBetweenPegIns, "address registered and watch imported"
			})
		assert.Equal(t, first.address, s.flyover.PegInAddress(t, institution), "the deposit address does not change")
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
