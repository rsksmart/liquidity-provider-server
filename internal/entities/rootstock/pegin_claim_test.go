package rootstock_test

import (
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculatePegInClaimPayableValue_FeeBelowAmount(t *testing.T) {
	amount := entities.NewWei(1000)
	fee := entities.NewWei(250)

	got, err := rootstock.CalculatePegInClaimPayableValue(amount, fee)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 0, got.Cmp(entities.NewWei(750)))
}

func TestCalculatePegInClaimPayableValue_FeeEqualToAmount(t *testing.T) {
	amount := entities.NewWei(1000)
	fee := entities.NewWei(1000)

	got, err := rootstock.CalculatePegInClaimPayableValue(amount, fee)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 0, got.Cmp(entities.NewWei(0)))
}

func TestCalculatePegInClaimPayableValue_FeeAboveAmount(t *testing.T) {
	amount := entities.NewWei(100)
	fee := entities.NewWei(250)

	got, err := rootstock.CalculatePegInClaimPayableValue(amount, fee)

	assert.Nil(t, got)
	require.ErrorIs(t, err, blockchain.ErrIncorrectFronting)
}

func TestCalculatePegInClaimPayableValue_NilInputs(t *testing.T) {
	amount := entities.NewWei(1000)
	fee := entities.NewWei(250)

	t.Run("nil amount", func(t *testing.T) {
		got, err := rootstock.CalculatePegInClaimPayableValue(nil, fee)
		assert.Nil(t, got)
		require.ErrorIs(t, err, blockchain.ErrIncorrectFronting)
	})
	t.Run("nil fee", func(t *testing.T) {
		got, err := rootstock.CalculatePegInClaimPayableValue(amount, nil)
		assert.Nil(t, got)
		require.ErrorIs(t, err, blockchain.ErrIncorrectFronting)
	})
	t.Run("both nil", func(t *testing.T) {
		got, err := rootstock.CalculatePegInClaimPayableValue(nil, nil)
		assert.Nil(t, got)
		require.ErrorIs(t, err, blockchain.ErrIncorrectFronting)
	})
}

func TestCalculatePegInClaimPayableValue_DoesNotMutateInputs(t *testing.T) {
	amount := entities.NewWei(1000)
	fee := entities.NewWei(250)

	_, err := rootstock.CalculatePegInClaimPayableValue(amount, fee)

	require.NoError(t, err)
	assert.Equal(t, 0, amount.Cmp(entities.NewWei(1000)))
	assert.Equal(t, 0, fee.Cmp(entities.NewWei(250)))
}

func TestPegInClaim_IsTerminal(t *testing.T) {
	cases := []struct {
		state    rootstock.PegInClaimState
		terminal bool
	}{
		{rootstock.PegInClaimCandidate, false},
		{rootstock.PegInClaimSubmitting, false},
		{rootstock.PegInClaimClaimed, false},
		{rootstock.PegInClaimRaceLost, true},
		{rootstock.PegInClaimRetryableFailure, false},
		{rootstock.PegInClaimResolved, true},
		{rootstock.PegInClaimResolveFailed, true},
		{rootstock.PegInClaimState(""), false},
	}
	for _, tc := range cases {
		t.Run(string(tc.state), func(t *testing.T) {
			claim := &rootstock.PegInClaim{State: tc.state}
			assert.Equal(t, tc.terminal, claim.IsTerminal())
		})
	}
}

func TestPegInClaim_IsTerminal_NilReceiver(t *testing.T) {
	var claim *rootstock.PegInClaim
	assert.False(t, claim.IsTerminal())
}

func TestPegInClaim_IsClaimed(t *testing.T) {
	cases := []struct {
		state   rootstock.PegInClaimState
		claimed bool
	}{
		{rootstock.PegInClaimCandidate, false},
		{rootstock.PegInClaimSubmitting, false},
		{rootstock.PegInClaimClaimed, true},
		{rootstock.PegInClaimRaceLost, false},
		{rootstock.PegInClaimRetryableFailure, false},
		{rootstock.PegInClaimResolved, false},
		{rootstock.PegInClaimResolveFailed, false},
		{rootstock.PegInClaimState(""), false},
	}
	for _, tc := range cases {
		t.Run(string(tc.state), func(t *testing.T) {
			claim := &rootstock.PegInClaim{State: tc.state}
			assert.Equal(t, tc.claimed, claim.IsClaimed())
		})
	}
}

func TestPegInClaim_IsClaimed_NilReceiver(t *testing.T) {
	var claim *rootstock.PegInClaim
	assert.False(t, claim.IsClaimed())
}

func TestPegInClaim_IsResolvable(t *testing.T) {
	cases := []struct {
		name       string
		state      rootstock.PegInClaimState
		pegInID    string
		resolvable bool
	}{
		{"claimed with peg-in id", rootstock.PegInClaimClaimed, "0a1b", true},
		{"claimed without peg-in id", rootstock.PegInClaimClaimed, "", false},
		{"candidate", rootstock.PegInClaimCandidate, "0a1b", false},
		{"submitting", rootstock.PegInClaimSubmitting, "0a1b", false},
		{"race lost", rootstock.PegInClaimRaceLost, "0a1b", false},
		{"retryable failure", rootstock.PegInClaimRetryableFailure, "0a1b", false},
		{"resolved", rootstock.PegInClaimResolved, "0a1b", false},
		{"resolve failed", rootstock.PegInClaimResolveFailed, "0a1b", false},
		{"empty state", rootstock.PegInClaimState(""), "0a1b", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			claim := &rootstock.PegInClaim{State: tc.state, PegInID: tc.pegInID}
			assert.Equal(t, tc.resolvable, claim.IsResolvable())
		})
	}
}

func TestPegInClaim_IsResolvable_NilReceiver(t *testing.T) {
	var claim *rootstock.PegInClaim
	assert.False(t, claim.IsResolvable())
}

func TestNewCandidatePegInClaim_UsesUTCTimestamps(t *testing.T) {
	before := time.Now().UTC().Add(-time.Second)
	entry := rootstock.PegInWatch{
		RskAddress: "0xrsk",
		BtcAddress: "btc-addr",
	}

	claim := rootstock.NewCandidatePegInClaim(entry, "deposit-txid", nil)
	after := time.Now().UTC().Add(time.Second)

	assert.Equal(t, "0xrsk", claim.RskAddress)
	assert.Equal(t, "deposit-txid", claim.DepositTxID)
	assert.Equal(t, "btc-addr", claim.BtcAddress)
	assert.Equal(t, rootstock.PegInClaimCandidate, claim.State)
	assert.Equal(t, time.UTC, claim.CreatedAt.Location())
	assert.Equal(t, time.UTC, claim.UpdatedAt.Location())
	assert.True(t, !claim.CreatedAt.Before(before) && !claim.CreatedAt.After(after))
	assert.Equal(t, claim.CreatedAt, claim.UpdatedAt)
}

func TestNewCandidatePegInClaim_PreservesExistingCreatedAt(t *testing.T) {
	createdAt := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)
	existing := &rootstock.PegInClaim{CreatedAt: createdAt}
	entry := rootstock.PegInWatch{RskAddress: "0xrsk", BtcAddress: "btc-addr"}

	before := time.Now().UTC().Add(-time.Second)
	claim := rootstock.NewCandidatePegInClaim(entry, "deposit-txid", existing)
	after := time.Now().UTC().Add(time.Second)

	assert.Equal(t, createdAt, claim.CreatedAt)
	assert.True(t, !claim.UpdatedAt.Before(before) && !claim.UpdatedAt.After(after))
	assert.NotEqual(t, claim.CreatedAt, claim.UpdatedAt)
}

func TestNewPegInClaimCompletedEvent(t *testing.T) {
	claim := rootstock.PegInClaim{RskAddress: "0xabc", DepositTxID: "aa", State: rootstock.PegInClaimClaimed}

	event := rootstock.NewPegInClaimCompletedEvent(claim)

	assert.Equal(t, rootstock.PegInClaimCompletedEventId, event.Id())
	assert.False(t, event.CreationTimestamp().IsZero())
	assert.Equal(t, claim, event.Claim)
}
