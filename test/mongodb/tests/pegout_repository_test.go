//go:build integration

package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/test/mongodb/support"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPegout_InsertAndGetQuote(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash := support.RandomHash()
	created := support.NewTestPegoutQuote(hash)
	require.NoError(t, pegoutRepo.InsertQuote(ctx, created))

	got, err := pegoutRepo.GetQuote(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, created.Quote.LbcAddress, got.LbcAddress)
	assert.Equal(t, created.Quote.LpRskAddress, got.LpRskAddress)
	assert.Equal(t, created.Quote.Nonce, got.Nonce)
	assert.Equal(t, created.Quote.ChainId, got.ChainId)
	assertWeiEqual(t, created.Quote.CallFee, got.CallFee)
	assertWeiEqual(t, created.Quote.Value, got.Value)
	assertWeiEqual(t, created.Quote.GasFee, got.GasFee)
	assertWeiEqual(t, created.Quote.PenaltyFee, got.PenaltyFee)
}

func TestPegout_GetQuote_NotFound(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	got, err := pegoutRepo.GetQuote(ctx, support.RandomHash())
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestPegout_GetPegoutCreationData(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash := support.RandomHash()
	created := support.NewTestPegoutQuote(hash)
	require.NoError(t, pegoutRepo.InsertQuote(ctx, created))

	got := pegoutRepo.GetPegoutCreationData(ctx, hash)
	require.NotNil(t, got.GasPrice)
	require.NotNil(t, got.FixedFee)
	assertWeiEqual(t, created.CreationData.GasPrice, got.GasPrice)
	assertWeiEqual(t, created.CreationData.FixedFee, got.FixedFee)
}

func TestPegout_GetPegoutCreationData_NotFound(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	got := pegoutRepo.GetPegoutCreationData(ctx, support.RandomHash())
	zeroVal := quote.PegoutCreationDataZeroValue()
	assertWeiEqual(t, zeroVal.GasPrice, got.GasPrice)
	assertWeiEqual(t, zeroVal.FixedFee, got.FixedFee)
}

func TestPegout_GetQuotesByHashesAndDate(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	h1 := support.RandomHash()
	h2 := support.RandomHash()
	h3 := support.RandomHash()

	q1 := support.NewTestPegoutQuote(h1)
	q1.Quote.Nonce = int64(1)
	q1.Quote.AgreementTimestamp = uint32(now.Unix())
	q2 := support.NewTestPegoutQuote(h2)
	q2.Quote.Nonce = int64(2)
	q2.Quote.AgreementTimestamp = uint32(now.Add(-2 * time.Hour).Unix())
	q3 := support.NewTestPegoutQuote(h3)
	q3.Quote.Nonce = int64(3)
	q3.Quote.AgreementTimestamp = uint32(now.Add(-48 * time.Hour).Unix())

	require.NoError(t, pegoutRepo.InsertQuote(ctx, q1))
	require.NoError(t, pegoutRepo.InsertQuote(ctx, q2))
	require.NoError(t, pegoutRepo.InsertQuote(ctx, q3))

	startDate := now.Add(-3 * time.Hour)
	endDate := now.Add(time.Hour)
	results, err := pegoutRepo.GetQuotesByHashesAndDate(ctx, []string{h1, h2, h3}, startDate, endDate)
	require.NoError(t, err)
	assert.Len(t, results, 2)
	nonces := []int64{results[0].Nonce, results[1].Nonce}
	assert.ElementsMatch(t, []int64{1, 2}, nonces)
}

func TestPegout_InsertAndGetRetainedQuote(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash := support.RandomHash()
	retained := support.NewTestRetainedPegoutQuote(hash, quote.PegoutStateWaitingForDeposit)
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, retained))

	got, err := pegoutRepo.GetRetainedQuote(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, retained.QuoteHash, got.QuoteHash)
	assert.Equal(t, retained.DepositAddress, got.DepositAddress)
	assert.Equal(t, retained.State, got.State)
	assertWeiEqual(t, retained.RequiredLiquidity, got.RequiredLiquidity)
	assertWeiEqual(t, retained.BridgeRefundGasPrice, got.BridgeRefundGasPrice)
	assertWeiEqual(t, retained.RefundPegoutGasPrice, got.RefundPegoutGasPrice)
	assertWeiEqual(t, retained.SendPegoutBtcFee, got.SendPegoutBtcFee)
}

func TestPegout_GetRetainedQuote_NotFound(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	got, err := pegoutRepo.GetRetainedQuote(ctx, support.RandomHash())
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestPegout_UpdateRetainedQuote(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash := support.RandomHash()
	retained := support.NewTestRetainedPegoutQuote(hash, quote.PegoutStateWaitingForDeposit)
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, retained))

	retained.State = quote.PegoutStateSendPegoutSucceeded
	retained.LpBtcTxHash = "updated_btc_tx"
	require.NoError(t, pegoutRepo.UpdateRetainedQuote(ctx, retained))

	got, err := pegoutRepo.GetRetainedQuote(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, quote.PegoutStateSendPegoutSucceeded, got.State)
	assert.Equal(t, "updated_btc_tx", got.LpBtcTxHash)
}

func TestPegout_UpdateRetainedQuote_NotFound(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	retained := support.NewTestRetainedPegoutQuote(support.RandomHash(), quote.PegoutStateSendPegoutSucceeded)
	err := pegoutRepo.UpdateRetainedQuote(ctx, retained)
	require.ErrorIs(t, err, usecases.QuoteNotFoundError)
}

func TestPegout_UpdateRetainedQuotes_HappyPath(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	h1 := support.RandomHash()
	h2 := support.RandomHash()
	r1 := support.NewTestRetainedPegoutQuote(h1, quote.PegoutStateWaitingForDeposit)
	r2 := support.NewTestRetainedPegoutQuote(h2, quote.PegoutStateWaitingForDeposit)
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, r1))
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, r2))

	r1.State = quote.PegoutStateSendPegoutSucceeded
	r2.State = quote.PegoutStateSendPegoutSucceeded
	err := pegoutRepo.UpdateRetainedQuotes(ctx, []quote.RetainedPegoutQuote{r1, r2})
	require.NoError(t, err)

	got1, err := pegoutRepo.GetRetainedQuote(ctx, h1)
	require.NoError(t, err)
	assert.Equal(t, quote.PegoutStateSendPegoutSucceeded, got1.State)

	got2, err := pegoutRepo.GetRetainedQuote(ctx, h2)
	require.NoError(t, err)
	assert.Equal(t, quote.PegoutStateSendPegoutSucceeded, got2.State)
}

func TestPegout_UpdateRetainedQuotes_PartialFailure(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	h1 := support.RandomHash()
	h2 := support.RandomHash()
	hNonExistent := support.RandomHash()
	r1 := support.NewTestRetainedPegoutQuote(h1, quote.PegoutStateWaitingForDeposit)
	r2 := support.NewTestRetainedPegoutQuote(h2, quote.PegoutStateWaitingForDeposit)
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, r1))
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, r2))

	r1.State = quote.PegoutStateSendPegoutSucceeded
	r2.State = quote.PegoutStateSendPegoutSucceeded
	rNonExistent := support.NewTestRetainedPegoutQuote(hNonExistent, quote.PegoutStateSendPegoutSucceeded)

	err := pegoutRepo.UpdateRetainedQuotes(ctx, []quote.RetainedPegoutQuote{r1, r2, rNonExistent})
	require.Error(t, err, "batch update with non-existent quote should return an error")
	assert.ErrorContains(t, err, "mismatch")
}

func TestPegout_GetRetainedQuoteByState(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	h1 := support.RandomHash()
	h2 := support.RandomHash()
	h3 := support.RandomHash()
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPegoutQuote(h1, quote.PegoutStateWaitingForDeposit)))
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPegoutQuote(h2, quote.PegoutStateSendPegoutSucceeded)))
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPegoutQuote(h3, quote.PegoutStateBtcReleased)))

	results, err := pegoutRepo.GetRetainedQuoteByState(ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateBtcReleased)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	hashes := []string{results[0].QuoteHash, results[1].QuoteHash}
	assert.ElementsMatch(t, []string{h1, h3}, hashes)
	for _, r := range results {
		assert.Contains(t, []quote.PegoutState{quote.PegoutStateWaitingForDeposit, quote.PegoutStateBtcReleased}, r.State)
	}
}

func TestPegout_GetQuotesByState(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	h1 := support.RandomHash()
	h2 := support.RandomHash()
	q1 := support.NewTestPegoutQuote(h1)
	q1.Quote.Nonce = int64(1)
	q2 := support.NewTestPegoutQuote(h2)
	q2.Quote.Nonce = int64(2)
	require.NoError(t, pegoutRepo.InsertQuote(ctx, q1))
	require.NoError(t, pegoutRepo.InsertQuote(ctx, q2))
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPegoutQuote(h1, quote.PegoutStateWaitingForDeposit)))
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPegoutQuote(h2, quote.PegoutStateSendPegoutFailed)))

	results, err := pegoutRepo.GetQuotesByState(ctx, quote.PegoutStateWaitingForDeposit)
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, int64(1), results[0].Nonce)
}

func TestPegout_DeleteQuotes(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	h1 := support.RandomHash()
	h2 := support.RandomHash()
	require.NoError(t, pegoutRepo.InsertQuote(ctx, support.NewTestPegoutQuote(h1)))
	require.NoError(t, pegoutRepo.InsertQuote(ctx, support.NewTestPegoutQuote(h2)))
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPegoutQuote(h1, quote.PegoutStateWaitingForDeposit)))
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPegoutQuote(h2, quote.PegoutStateSendPegoutSucceeded)))

	count, err := pegoutRepo.DeleteQuotes(ctx, []string{h1, h2})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, uint(4))

	got, err := pegoutRepo.GetQuote(ctx, h1)
	require.NoError(t, err)
	assert.Nil(t, got)

	got, err = pegoutRepo.GetQuote(ctx, h2)
	require.NoError(t, err)
	assert.Nil(t, got)

	retained, err := pegoutRepo.GetRetainedQuote(ctx, h1)
	require.NoError(t, err)
	assert.Nil(t, retained)

	retained, err = pegoutRepo.GetRetainedQuote(ctx, h2)
	require.NoError(t, err)
	assert.Nil(t, retained)
}

func TestPegout_UpsertPegoutDeposit(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	txHash := "0x" + support.RandomHash()
	quoteHash := support.RandomHash()
	deposit := support.NewTestPegoutDeposit(txHash, quoteHash)

	err := pegoutRepo.UpsertPegoutDeposit(ctx, deposit)
	require.NoError(t, err)

	deposits, err := pegoutRepo.ListPegoutDepositsByAddress(ctx, deposit.From)
	require.NoError(t, err)
	require.Len(t, deposits, 1)
	assert.Equal(t, txHash, deposits[0].TxHash)
	assertWeiEqual(t, deposit.Amount, deposits[0].Amount)
}

func TestPegout_UpsertPegoutDeposit_UniqueIndex(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	txHash := "0x" + support.RandomHash()
	quoteHash := support.RandomHash()
	deposit := support.NewTestPegoutDeposit(txHash, quoteHash)
	require.NoError(t, pegoutRepo.UpsertPegoutDeposit(ctx, deposit))

	deposit.BlockNumber = 99999
	require.NoError(t, pegoutRepo.UpsertPegoutDeposit(ctx, deposit))

	deposits, err := pegoutRepo.ListPegoutDepositsByAddress(ctx, deposit.From)
	require.NoError(t, err)
	require.Len(t, deposits, 1)
	assert.Equal(t, uint64(99999), deposits[0].BlockNumber)
}

func TestPegout_UpsertPegoutDeposits(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	d1 := support.NewTestPegoutDeposit("0x"+support.RandomHash(), support.RandomHash())
	d2 := support.NewTestPegoutDeposit("0x"+support.RandomHash(), support.RandomHash())
	// Make the behavior under test explicit (both deposits belong to same sender).
	commonFrom := "0x1234567890abcdef1234567890abcdef12345678"
	d1.From = commonFrom
	d2.From = commonFrom

	err := pegoutRepo.UpsertPegoutDeposits(ctx, []quote.PegoutDeposit{d1, d2})
	require.NoError(t, err)

	deposits, err := pegoutRepo.ListPegoutDepositsByAddress(ctx, commonFrom)
	require.NoError(t, err)
	assert.Len(t, deposits, 2)
	txHashes := []string{deposits[0].TxHash, deposits[1].TxHash}
	assert.ElementsMatch(t, []string{d1.TxHash, d2.TxHash}, txHashes)
}

func TestPegout_ListPegoutDepositsByAddress(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	addr := "0xABCdef1234567890abcdef1234567890abcdef12"
	otherAddr := "0x9999999999999999999999999999999999999999"

	d1 := support.NewTestPegoutDeposit("0x"+support.RandomHash(), support.RandomHash())
	d1.From = addr
	d2 := support.NewTestPegoutDeposit("0x"+support.RandomHash(), support.RandomHash())
	d2.From = otherAddr

	require.NoError(t, pegoutRepo.UpsertPegoutDeposit(ctx, d1))
	require.NoError(t, pegoutRepo.UpsertPegoutDeposit(ctx, d2))

	// Case-insensitive search
	deposits, err := pegoutRepo.ListPegoutDepositsByAddress(ctx, "0xabcdef1234567890ABCDEF1234567890abcdef12")
	require.NoError(t, err)
	assert.Len(t, deposits, 1)
	assert.Equal(t, addr, deposits[0].From)
}

func TestPegout_ListQuotesByDateRange(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	h1 := support.RandomHash()
	h2 := support.RandomHash()
	h3 := support.RandomHash()

	q1 := support.NewTestPegoutQuote(h1)
	q1.Quote.Nonce = int64(1)
	q1.Quote.AgreementTimestamp = uint32(now.Unix())
	q2 := support.NewTestPegoutQuote(h2)
	q2.Quote.Nonce = int64(2)
	q2.Quote.AgreementTimestamp = uint32(now.Add(-1 * time.Hour).Unix())
	q3 := support.NewTestPegoutQuote(h3)
	q3.Quote.Nonce = int64(3)
	q3.Quote.AgreementTimestamp = uint32(now.Add(-48 * time.Hour).Unix())

	require.NoError(t, pegoutRepo.InsertQuote(ctx, q1))
	require.NoError(t, pegoutRepo.InsertQuote(ctx, q2))
	require.NoError(t, pegoutRepo.InsertQuote(ctx, q3))

	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPegoutQuote(h1, quote.PegoutStateWaitingForDeposit)))
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPegoutQuote(h2, quote.PegoutStateSendPegoutSucceeded)))

	startDate := now.Add(-3 * time.Hour)
	endDate := now.Add(time.Hour)

	results, count, err := pegoutRepo.ListQuotesByDateRange(ctx, startDate, endDate, 0, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, results, 2)
	nonces := []int64{results[0].Quote.Nonce, results[1].Quote.Nonce}
	assert.ElementsMatch(t, []int64{1, 2}, nonces)

	// With pagination
	results, count, err = pegoutRepo.ListQuotesByDateRange(ctx, startDate, endDate, 1, 1)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
	require.Len(t, results, 1)
	// Sorted by agreement_timestamp ascending, so the oldest in-range quote should be returned first.
	assert.Equal(t, int64(2), results[0].Quote.Nonce)
}

func TestPegout_GetRetainedQuotesForAddress(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	addr := "0x1234567890abcdef1234567890abcdef12345678"
	otherAddr := "0xaaaa567890abcdef1234567890abcdef12345678"

	h1 := support.RandomHash()
	h2 := support.RandomHash()
	h3 := support.RandomHash()

	r1 := support.NewTestRetainedPegoutQuote(h1, quote.PegoutStateWaitingForDeposit)
	r1.OwnerAccountAddress = addr
	r2 := support.NewTestRetainedPegoutQuote(h2, quote.PegoutStateSendPegoutSucceeded)
	r2.OwnerAccountAddress = addr
	r3 := support.NewTestRetainedPegoutQuote(h3, quote.PegoutStateWaitingForDeposit)
	r3.OwnerAccountAddress = otherAddr

	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, r1))
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, r2))
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, r3))

	results, err := pegoutRepo.GetRetainedQuotesForAddress(ctx, addr, quote.PegoutStateWaitingForDeposit, quote.PegoutStateSendPegoutSucceeded)
	require.NoError(t, err)
	assert.Len(t, results, 2)
	hashes := []string{results[0].QuoteHash, results[1].QuoteHash}
	assert.ElementsMatch(t, []string{h1, h2}, hashes)
	for _, r := range results {
		assert.Equal(t, addr, r.OwnerAccountAddress)
	}
}

func TestPegout_GetRetainedQuotesInBatch(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	releaseTx1 := "0xrelease_abc123"
	releaseTx2 := "0xrelease_def456"

	h1 := support.RandomHash()
	h2 := support.RandomHash()
	h3 := support.RandomHash()

	r1 := support.NewTestRetainedPegoutQuote(h1, quote.PegoutStateBridgeTxSucceeded)
	r1.BridgeRefundTxHash = releaseTx1
	r2 := support.NewTestRetainedPegoutQuote(h2, quote.PegoutStateBridgeTxSucceeded)
	r2.BridgeRefundTxHash = releaseTx2
	r3 := support.NewTestRetainedPegoutQuote(h3, quote.PegoutStateBridgeTxSucceeded)
	r3.BridgeRefundTxHash = "0xother_release"

	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, r1))
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, r2))
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, r3))

	batch := rootstock.BatchPegOut{
		TransactionHash:    "0x" + support.RandomHash(),
		BlockHash:          "0xblockhash",
		BlockNumber:        100,
		BtcTxHash:          "btctx123",
		ReleaseRskTxHashes: []string{releaseTx1, releaseTx2},
	}

	results, err := pegoutRepo.GetRetainedQuotesInBatch(ctx, batch)
	require.NoError(t, err)
	assert.Len(t, results, 2)
	releaseHashes := []string{results[0].BridgeRefundTxHash, results[1].BridgeRefundTxHash}
	assert.ElementsMatch(t, []string{releaseTx1, releaseTx2}, releaseHashes)
}

func TestPegout_GetQuotesWithRetainedByStateAndDate(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	h1 := support.RandomHash()
	h2 := support.RandomHash()
	h3 := support.RandomHash()

	q1 := support.NewTestPegoutQuote(h1)
	q1.Quote.Nonce = int64(1)
	q1.Quote.AgreementTimestamp = uint32(now.Unix())
	q2 := support.NewTestPegoutQuote(h2)
	q2.Quote.Nonce = int64(2)
	q2.Quote.AgreementTimestamp = uint32(now.Add(-1 * time.Hour).Unix())
	q3 := support.NewTestPegoutQuote(h3)
	q3.Quote.Nonce = int64(3)
	q3.Quote.AgreementTimestamp = uint32(now.Unix())

	require.NoError(t, pegoutRepo.InsertQuote(ctx, q1))
	require.NoError(t, pegoutRepo.InsertQuote(ctx, q2))
	require.NoError(t, pegoutRepo.InsertQuote(ctx, q3))

	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPegoutQuote(h1, quote.PegoutStateWaitingForDeposit)))
	require.NoError(t, pegoutRepo.InsertRetainedQuote(ctx, support.NewTestRetainedPegoutQuote(h2, quote.PegoutStateSendPegoutFailed)))
	// h3 has no retained quote

	startDate := now.Add(-3 * time.Hour)
	endDate := now.Add(time.Hour)

	results, err := pegoutRepo.GetQuotesWithRetainedByStateAndDate(ctx, []quote.PegoutState{quote.PegoutStateWaitingForDeposit}, startDate, endDate)
	require.NoError(t, err)
	// Should include h1 (matching state) and h3 (no retained = included)
	assert.Len(t, results, 2)
	nonces := []int64{results[0].Quote.Nonce, results[1].Quote.Nonce}
	assert.ElementsMatch(t, []int64{1, 3}, nonces)

	for _, r := range results {
		assert.NotNil(t, r.RetainedQuote.BridgeRefundGasPrice)
		assert.NotNil(t, r.RetainedQuote.RefundPegoutGasPrice)
		assert.NotNil(t, r.RetainedQuote.SendPegoutBtcFee)

		// When there is no retained quote, the repository normalizes it to a zero-valued retained struct.
		if r.Quote.Nonce == int64(3) {
			assert.Empty(t, r.RetainedQuote.QuoteHash)
			assert.Empty(t, r.RetainedQuote.State)
		}
		if r.Quote.Nonce == int64(1) {
			assert.Equal(t, h1, r.RetainedQuote.QuoteHash)
			assert.Equal(t, quote.PegoutStateWaitingForDeposit, r.RetainedQuote.State)
		}
	}
}
