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
		{rootstock.PegInClaimClaimed, true},
		{rootstock.PegInClaimRaceLost, true},
		{rootstock.PegInClaimRetryableFailure, false},
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

func TestNewCandidatePegInClaim_UsesUTCTimestamps(t *testing.T) {
	before := time.Now().UTC().Add(-time.Second)
	entry := rootstock.PegInWatch{
		RskAddress: "0xrsk",
		BtcAddress: "btc-addr",
	}
	reserved := entities.NewWei(500)

	claim := rootstock.NewCandidatePegInClaim(entry, "deposit-txid", reserved, nil)
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
	claim := rootstock.NewCandidatePegInClaim(entry, "deposit-txid", entities.NewWei(1), existing)
	after := time.Now().UTC().Add(time.Second)

	assert.Equal(t, createdAt, claim.CreatedAt)
	assert.True(t, !claim.UpdatedAt.Before(before) && !claim.UpdatedAt.After(after))
	assert.NotEqual(t, claim.CreatedAt, claim.UpdatedAt)
}

func TestNewCandidatePegInClaim_CopiesReservedWei(t *testing.T) {
	reserved := entities.NewWei(500)
	entry := rootstock.PegInWatch{RskAddress: "0xrsk", BtcAddress: "btc-addr"}

	claim := rootstock.NewCandidatePegInClaim(entry, "deposit-txid", reserved, nil)

	require.NotNil(t, claim.ReservedWei)
	assert.NotSame(t, reserved, claim.ReservedWei)
	assert.Equal(t, 0, claim.ReservedWei.Cmp(entities.NewWei(500)))

	reserved.Add(reserved, entities.NewWei(100))
	assert.Equal(t, 0, claim.ReservedWei.Cmp(entities.NewWei(500)))
}

func TestNewCandidatePegInClaim_NilReservedIsZeroWei(t *testing.T) {
	entry := rootstock.PegInWatch{RskAddress: "0xrsk", BtcAddress: "btc-addr"}

	claim := rootstock.NewCandidatePegInClaim(entry, "deposit-txid", nil, nil)

	require.NotNil(t, claim.ReservedWei)
	assert.Equal(t, 0, claim.ReservedWei.Cmp(entities.NewWei(0)))
}
