package pegin_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test/integration/regtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Rounds driven after the settlement, so the resolve watcher passes over the claim again.
const extraRounds = 5

// TestCommitFirstPegIn drives one commit-first peg-in on the local regtest stack, from the
// institution's BTC deposit to the LP getting its advance back, and checks each step on both
// chains and in the LPS database. The test acts as the institution and the watchtower, the LPS
// acts on its own.
func TestCommitFirstPegIn(t *testing.T) {
	s := newPegInStack(t)
	institution := s.rsk.NewAccount(t, new(big.Int)).Address
	// Registers the address the way the watchtower does; the contracts call it the registrant.
	fakeWatchtower := s.rsk.NewAccount(t, big.NewInt(1_000_000_000_000_000_000))
	amount := depositAmount(t, s.flyover)
	institutionBalanceBefore := s.rsk.Balance(t, institution)
	lpBalanceBefore := s.flyover.Balance(t, s.lp)
	fakeWatchtowerBalanceBefore := s.flyover.Balance(t, fakeWatchtower.Address)
	t.Logf("institution %s, fake watchtower %s, deposit %s", institution, fakeWatchtower.Address, amount)

	var (
		paid      deposit
		requested *regtest.PegInRequested
		claim     *rootstock.PegInClaim
	)

	runPhase(t, "the institution's deposit reaches the bridge", func(t *testing.T) {
		paid = s.payDeposit(t, institution, amount)
	})

	runPhase(t, "the fake watchtower registers the address", func(t *testing.T) {
		s.flyover.RegisterAddress(t, fakeWatchtower, institution, paid.proof)
		assert.True(t, s.flyover.IsRegistered(t, institution))
		events := s.flyover.AddressRegisteredEvents(t, institution)
		require.Len(t, events, 1)
		assert.Equal(t, institution, events[0].RskAddr)
		assert.Equal(t, fakeWatchtower.Address, events[0].Registrant)
	})

	runPhase(t, "the LPS discovers and imports the address", func(t *testing.T) {
		watch := s.waitWatchImported(t, institution)
		assert.Equal(t, paid.address, watch.BtcAddress)
		assert.Equal(t, fakeWatchtower.Address.Hex(), common.HexToAddress(watch.Registrant).Hex())
	})

	runPhase(t, "the LPS fronts the RBTC and settles its claim", func(t *testing.T) {
		claim, requested = s.waitClaimed(t, institution, paid)
		fee := s.flyover.PegInFee(t, paid.amountWei)
		assert.Equal(t, s.lp, requested.Claimer)
		assert.Equal(t, paid.amountWei, requested.Amount)
		assert.Equal(t, new(big.Int).Sub(paid.amountWei, fee), requested.NetToUser)
		assert.Equal(t, new(big.Int).Add(institutionBalanceBefore, requested.NetToUser), s.rsk.Balance(t, institution),
			"the institution receives the net amount")
		assert.Equal(t, paid.address, claim.BtcAddress)
		assert.NotEmpty(t, claim.RequestTxHash)
	})

	runPhase(t, "the LPS resolves the peg-in and gets its advance back", func(t *testing.T) {
		var resolved *regtest.PegInResolved
		claim, resolved = s.waitResolved(t, institution, paid, requested)
		feeAtClaim := new(big.Int).Sub(requested.Amount, requested.NetToUser)
		// The first peg-in to an address pays the registrant fee to the fake watchtower, capped at
		// the LP's fee.
		registrantFee := minBigInt(s.flyover.PegInConfiguration(t).RegistrantFee, feeAtClaim)
		assert.Equal(t, s.lp, resolved.Claimer)
		assert.Equal(t, fakeWatchtower.Address, resolved.Registrant)
		assert.Equal(t, registrantFee, resolved.RegistrantFee)
		assert.Equal(t, new(big.Int).Sub(new(big.Int).Add(requested.NetToUser, feeAtClaim), registrantFee),
			resolved.ClaimerPayout, "fronted + fee - registrant fee")
		assert.Zero(t, resolved.UserPayout.Sign())

		assert.Equal(t, new(big.Int).Add(lpBalanceBefore, resolved.ClaimerPayout), s.flyover.Balance(t, s.lp))
		assert.Equal(t, new(big.Int).Add(fakeWatchtowerBalanceBefore, resolved.RegistrantFee),
			s.flyover.Balance(t, fakeWatchtower.Address))

		assert.NotEmpty(t, claim.ResolveTxHash)
		assert.Zero(t, entities.NewBigWei(resolved.ClaimerPayout).Cmp(claim.ClaimerPayout))
		assert.Zero(t, entities.NewBigWei(resolved.RegistrantFee).Cmp(claim.RegistrantFee))
	})

	runPhase(t, "the peg-in is settled only once", func(t *testing.T) {
		resolveTxHash := claim.ResolveTxHash
		rounds := 0
		s.driver.AdvanceUntil(t, "extra resolve passes", phaseTimeout, regtest.Step{BtcBlocks: 1, RskBlocks: 1},
			func(t *testing.T) (bool, string) {
				rounds++
				return rounds > extraRounds, fmt.Sprintf("round %d of %d", rounds, extraRounds)
			})
		assert.Len(t, s.flyover.PegInResolvedEvents(t, requested.PegInId), 1)
		claim = s.db.Claim(t, institution, paid.txID)
		assert.Equal(t, rootstock.PegInClaimResolved, claim.State)
		assert.Equal(t, resolveTxHash, claim.ResolveTxHash)
	})
}
