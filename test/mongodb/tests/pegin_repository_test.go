//go:build integration

package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/test/mongodb/support"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPegin_InsertAndGetQuote(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash := support.RandomHash()
	created := support.NewTestPeginQuote(hash)
	require.NoError(t, peginRepo.InsertQuote(ctx, created))

	got, err := peginRepo.GetQuote(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, created.Quote.FedBtcAddress, got.FedBtcAddress)
	assert.Equal(t, created.Quote.LbcAddress, got.LbcAddress)
	assert.Equal(t, created.Quote.LpRskAddress, got.LpRskAddress)
	assert.Equal(t, created.Quote.Nonce, got.Nonce)
	assert.Equal(t, created.Quote.GasLimit, got.GasLimit)
	assert.Equal(t, created.Quote.ChainId, got.ChainId)
	assertWeiEqual(t, created.Quote.CallFee, got.CallFee)
	assertWeiEqual(t, created.Quote.Value, got.Value)
	assertWeiEqual(t, created.Quote.GasFee, got.GasFee)
	assertWeiEqual(t, created.Quote.PenaltyFee, got.PenaltyFee)
}

func TestPegin_GetQuote_NotFound(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	got, err := peginRepo.GetQuote(ctx, support.RandomHash())
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestPegin_GetPeginCreationData(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash := support.RandomHash()
	created := support.NewTestPeginQuote(hash)
	require.NoError(t, peginRepo.InsertQuote(ctx, created))

	got := peginRepo.GetPeginCreationData(ctx, hash)
	require.NotNil(t, got.GasPrice)
	require.NotNil(t, got.FixedFee)
	assertWeiEqual(t, created.CreationData.GasPrice, got.GasPrice)
	assertWeiEqual(t, created.CreationData.FixedFee, got.FixedFee)
}

func TestPegin_GetPeginCreationData_NotFound(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	got := peginRepo.GetPeginCreationData(ctx, support.RandomHash())
	zeroVal := quote.PeginCreationDataZeroValue()
	assertWeiEqual(t, zeroVal.GasPrice, got.GasPrice)
	assertWeiEqual(t, zeroVal.FixedFee, got.FixedFee)
}

func TestPegin_GetQuotesByHashesAndDate(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	now := time.Now()
	hash1 := support.RandomHash()
	hash2 := support.RandomHash()
	hash3 := support.RandomHash()

	q1 := support.NewTestPeginQuote(hash1)
	q1.Quote.AgreementTimestamp = uint32(now.Unix())
	q2 := support.NewTestPeginQuote(hash2)
	q2.Quote.AgreementTimestamp = uint32(now.Add(-2 * time.Hour).Unix())
	q3 := support.NewTestPeginQuote(hash3)
	q3.Quote.AgreementTimestamp = uint32(now.Add(-48 * time.Hour).Unix())

	require.NoError(t, peginRepo.InsertQuote(ctx, q1))
	require.NoError(t, peginRepo.InsertQuote(ctx, q2))
	require.NoError(t, peginRepo.InsertQuote(ctx, q3))

	startDate := now.Add(-3 * time.Hour)
	endDate := now.Add(time.Hour)
	results, err := peginRepo.GetQuotesByHashesAndDate(ctx, []string{hash1, hash2, hash3}, startDate, endDate)
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestPegin_InsertAndGetRetainedQuote(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash := support.RandomHash()
	retained := support.NewTestRetainedPeginQuote(hash, quote.PeginStateWaitingForDeposit)
	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, retained))

	got, err := peginRepo.GetRetainedQuote(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, retained.QuoteHash, got.QuoteHash)
	assert.Equal(t, retained.DepositAddress, got.DepositAddress)
	assert.Equal(t, retained.Signature, got.Signature)
	assert.Equal(t, retained.State, got.State)
	assertWeiEqual(t, retained.RequiredLiquidity, got.RequiredLiquidity)
	assertWeiEqual(t, retained.CallForUserGasPrice, got.CallForUserGasPrice)
	assertWeiEqual(t, retained.RegisterPeginGasPrice, got.RegisterPeginGasPrice)
}

func TestPegin_GetRetainedQuote_NotFound(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	got, err := peginRepo.GetRetainedQuote(ctx, support.RandomHash())
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestPegin_UpdateRetainedQuote(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash := support.RandomHash()
	retained := support.NewTestRetainedPeginQuote(hash, quote.PeginStateWaitingForDeposit)
	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, retained))

	retained.State = quote.PeginStateCallForUserSucceeded
	retained.CallForUserTxHash = "0xupdatedcalltx"
	retained.CallForUserGasUsed = 75000
	require.NoError(t, peginRepo.UpdateRetainedQuote(ctx, retained))

	got, err := peginRepo.GetRetainedQuote(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, quote.PeginStateCallForUserSucceeded, got.State)
	assert.Equal(t, "0xupdatedcalltx", got.CallForUserTxHash)
	assert.Equal(t, uint64(75000), got.CallForUserGasUsed)
}

func TestPegin_UpdateRetainedQuote_NotFound(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	retained := support.NewTestRetainedPeginQuote(support.RandomHash(), quote.PeginStateCallForUserSucceeded)
	err := peginRepo.UpdateRetainedQuote(ctx, retained)
	require.ErrorIs(t, err, usecases.QuoteNotFoundError)
}

func TestPegin_GetRetainedQuoteByState(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	h1 := support.RandomHash()
	h2 := support.RandomHash()
	h3 := support.RandomHash()
	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPeginQuote(h1, quote.PeginStateWaitingForDeposit)))
	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPeginQuote(h2, quote.PeginStateCallForUserSucceeded)))
	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPeginQuote(h3, quote.PeginStateRegisterPegInSucceeded)))

	results, err := peginRepo.GetRetainedQuoteByState(ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateCallForUserSucceeded)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	for _, r := range results {
		assert.Contains(t, []quote.PeginState{quote.PeginStateWaitingForDeposit, quote.PeginStateCallForUserSucceeded}, r.State)
		assert.NotNil(t, r.CallForUserGasPrice)
		assert.NotNil(t, r.RegisterPeginGasPrice)
	}
}

func TestPegin_GetQuotesByState(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	h1 := support.RandomHash()
	h2 := support.RandomHash()
	require.NoError(t, peginRepo.InsertQuote(ctx, support.NewTestPeginQuote(h1)))
	require.NoError(t, peginRepo.InsertQuote(ctx, support.NewTestPeginQuote(h2)))
	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPeginQuote(h1, quote.PeginStateWaitingForDeposit)))
	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPeginQuote(h2, quote.PeginStateCallForUserFailed)))

	results, err := peginRepo.GetQuotesByState(ctx, quote.PeginStateWaitingForDeposit)
	require.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestPegin_DeleteQuotes(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	h1 := support.RandomHash()
	h2 := support.RandomHash()
	require.NoError(t, peginRepo.InsertQuote(ctx, support.NewTestPeginQuote(h1)))
	require.NoError(t, peginRepo.InsertQuote(ctx, support.NewTestPeginQuote(h2)))
	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPeginQuote(h1, quote.PeginStateWaitingForDeposit)))
	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPeginQuote(h2, quote.PeginStateCallForUserSucceeded)))

	count, err := peginRepo.DeleteQuotes(ctx, []string{h1, h2})
	require.NoError(t, err)
	assert.Greater(t, count, uint(0))

	got, err := peginRepo.GetQuote(ctx, h1)
	require.NoError(t, err)
	assert.Nil(t, got)

	got2, err := peginRepo.GetRetainedQuote(ctx, h1)
	require.NoError(t, err)
	assert.Nil(t, got2)
}

func TestPegin_ListQuotesByDateRange(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	now := time.Now()
	h1 := support.RandomHash()
	h2 := support.RandomHash()
	h3 := support.RandomHash()

	q1 := support.NewTestPeginQuote(h1)
	q1.Quote.AgreementTimestamp = uint32(now.Unix())
	q2 := support.NewTestPeginQuote(h2)
	q2.Quote.AgreementTimestamp = uint32(now.Add(-1 * time.Hour).Unix())
	q3 := support.NewTestPeginQuote(h3)
	q3.Quote.AgreementTimestamp = uint32(now.Add(-48 * time.Hour).Unix())

	require.NoError(t, peginRepo.InsertQuote(ctx, q1))
	require.NoError(t, peginRepo.InsertQuote(ctx, q2))
	require.NoError(t, peginRepo.InsertQuote(ctx, q3))

	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPeginQuote(h1, quote.PeginStateWaitingForDeposit)))
	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPeginQuote(h2, quote.PeginStateCallForUserSucceeded)))

	startDate := now.Add(-3 * time.Hour)
	endDate := now.Add(time.Hour)

	// No pagination
	results, count, err := peginRepo.ListQuotesByDateRange(ctx, startDate, endDate, 0, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, results, 2)

	// With pagination
	results, count, err = peginRepo.ListQuotesByDateRange(ctx, startDate, endDate, 1, 1)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.Len(t, results, 1)
}

func TestPegin_GetRetainedQuotesForAddress(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	addr := "0x1234567890abcdef1234567890abcdef12345678"
	otherAddr := "0xaaaa567890abcdef1234567890abcdef12345678"

	h1 := support.RandomHash()
	h2 := support.RandomHash()
	h3 := support.RandomHash()

	r1 := support.NewTestRetainedPeginQuote(h1, quote.PeginStateWaitingForDeposit)
	r1.OwnerAccountAddress = addr
	r2 := support.NewTestRetainedPeginQuote(h2, quote.PeginStateCallForUserSucceeded)
	r2.OwnerAccountAddress = addr
	r3 := support.NewTestRetainedPeginQuote(h3, quote.PeginStateWaitingForDeposit)
	r3.OwnerAccountAddress = otherAddr

	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, r1))
	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, r2))
	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, r3))

	results, err := peginRepo.GetRetainedQuotesForAddress(ctx, addr, quote.PeginStateWaitingForDeposit, quote.PeginStateCallForUserSucceeded)
	require.NoError(t, err)
	assert.Len(t, results, 2)
	for _, r := range results {
		assert.Equal(t, addr, r.OwnerAccountAddress)
	}
}

func TestPegin_GetQuotesWithRetainedByStateAndDate(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	h1 := support.RandomHash()
	h2 := support.RandomHash()
	h3 := support.RandomHash()

	q1 := support.NewTestPeginQuote(h1)
	q1.Quote.Nonce = int64(1)
	q1.Quote.AgreementTimestamp = uint32(now.Unix())
	q2 := support.NewTestPeginQuote(h2)
	q2.Quote.Nonce = int64(2)
	q2.Quote.AgreementTimestamp = uint32(now.Add(-1 * time.Hour).Unix())
	q3 := support.NewTestPeginQuote(h3)
	q3.Quote.Nonce = int64(3)
	q3.Quote.AgreementTimestamp = uint32(now.Unix())

	require.NoError(t, peginRepo.InsertQuote(ctx, q1))
	require.NoError(t, peginRepo.InsertQuote(ctx, q2))
	require.NoError(t, peginRepo.InsertQuote(ctx, q3))

	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPeginQuote(h1, quote.PeginStateWaitingForDeposit)))
	require.NoError(t, peginRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPeginQuote(h2, quote.PeginStateCallForUserFailed)))
	// h3 has no retained quote

	startDate := now.Add(-3 * time.Hour)
	endDate := now.Add(time.Hour)

	results, err := peginRepo.GetQuotesWithRetainedByStateAndDate(ctx, []quote.PeginState{quote.PeginStateWaitingForDeposit}, startDate, endDate)
	require.NoError(t, err)
	// Should include h1 (matching state) and h3 (no retained = included)
	assert.Len(t, results, 2)
	nonces := []int64{results[0].Quote.Nonce, results[1].Quote.Nonce}
	assert.ElementsMatch(t, []int64{1, 3}, nonces)

	for _, r := range results {
		assert.NotNil(t, r.RetainedQuote.CallForUserGasPrice)
		assert.NotNil(t, r.RetainedQuote.RegisterPeginGasPrice)

		// When there is no retained quote, the repository normalizes it to a zero-valued retained struct.
		if r.Quote.Nonce == int64(3) {
			assert.Empty(t, r.RetainedQuote.QuoteHash)
			assert.Empty(t, r.RetainedQuote.State)
		}
		if r.Quote.Nonce == int64(1) {
			assert.Equal(t, h1, r.RetainedQuote.QuoteHash)
			assert.Equal(t, quote.PeginStateWaitingForDeposit, r.RetainedQuote.State)
		}
	}
}
