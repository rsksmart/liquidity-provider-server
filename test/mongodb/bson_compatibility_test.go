//go:build integration

package mongodb_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

// --- Layer 1: Frozen fixture dataset tests ---

func hasFixtureFiles() bool {
	dir := fixturesPath()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".json" {
			return true
		}
	}
	return false
}

func requireRestoredFixtures(t *testing.T) context.Context {
	t.Helper()

	if !hasFixtureFiles() {
		t.Skip("No fixture files found. Run generate-fixtures.sh first.")
	}
	cleanCollections(t)
	restoreFixtures(t)
	return context.Background()
}

func assertWeiEq(t *testing.T, expected int64, actual *entities.Wei) {
	t.Helper()

	require.NotNil(t, actual)
	assert.Zero(t, actual.Cmp(entities.NewWei(expected)))
}

func TestFixtures_PeginQuotes_Readable(t *testing.T) {
	ctx := requireRestoredFixtures(t)

	// Read all pegin quotes from the fixture collection via raw scan
	coll := rawCollection(mongo.PeginQuoteCollection)
	cursor, err := coll.Find(ctx, bson.M{})
	require.NoError(t, err)
	defer cursor.Close(ctx)

	type storedQuote struct {
		Hash string `bson:"hash"`
	}
	var hashes []storedQuote
	require.NoError(t, cursor.All(ctx, &hashes))
	require.NotEmpty(t, hashes, "fixture should have pegin quotes")

	for _, stored := range hashes {
		got, err := peginRepo.GetQuote(ctx, stored.Hash)
		require.NoError(t, err, "GetQuote failed for hash %s", stored.Hash)
		require.NotNil(t, got, "GetQuote returned nil for hash %s", stored.Hash)
		assert.NotEmpty(t, got.LbcAddress)
		assert.NotNil(t, got.CallFee)
		assert.NotNil(t, got.Value)
	}
}

func TestFixtures_RetainedPeginQuotes_Readable(t *testing.T) {
	ctx := requireRestoredFixtures(t)

	coll := rawCollection(mongo.RetainedPeginQuoteCollection)
	cursor, err := coll.Find(ctx, bson.M{})
	require.NoError(t, err)
	defer cursor.Close(ctx)

	type storedRetained struct {
		QuoteHash string `bson:"quote_hash"`
	}
	var items []storedRetained
	require.NoError(t, cursor.All(ctx, &items))
	require.NotEmpty(t, items)

	for _, item := range items {
		got, err := peginRepo.GetRetainedQuote(ctx, item.QuoteHash)
		require.NoError(t, err, "GetRetainedQuote failed for hash %s", item.QuoteHash)
		require.NotNil(t, got)
		assert.NotEmpty(t, got.State)
		assert.NotNil(t, got.RequiredLiquidity)
		assert.NotNil(t, got.CallForUserGasPrice)
		assert.NotNil(t, got.RegisterPeginGasPrice)
	}
}

func TestFixtures_PegoutQuotes_Readable(t *testing.T) {
	ctx := requireRestoredFixtures(t)

	coll := rawCollection(mongo.PegoutQuoteCollection)
	cursor, err := coll.Find(ctx, bson.M{})
	require.NoError(t, err)
	defer cursor.Close(ctx)

	type storedQuote struct {
		Hash string `bson:"hash"`
	}
	var hashes []storedQuote
	require.NoError(t, cursor.All(ctx, &hashes))
	require.NotEmpty(t, hashes)

	for _, stored := range hashes {
		got, err := pegoutRepo.GetQuote(ctx, stored.Hash)
		require.NoError(t, err, "GetQuote failed for hash %s", stored.Hash)
		require.NotNil(t, got)
		assert.NotNil(t, got.CallFee)
		assert.NotNil(t, got.Value)
	}
}

func TestFixtures_RetainedPegoutQuotes_Readable(t *testing.T) {
	ctx := requireRestoredFixtures(t)

	coll := rawCollection(mongo.RetainedPegoutQuoteCollection)
	cursor, err := coll.Find(ctx, bson.M{})
	require.NoError(t, err)
	defer cursor.Close(ctx)

	type storedRetained struct {
		QuoteHash string `bson:"quote_hash"`
	}
	var items []storedRetained
	require.NoError(t, cursor.All(ctx, &items))
	require.NotEmpty(t, items)

	for _, item := range items {
		got, err := pegoutRepo.GetRetainedQuote(ctx, item.QuoteHash)
		require.NoError(t, err, "GetRetainedQuote failed for hash %s", item.QuoteHash)
		require.NotNil(t, got)
		assert.NotEmpty(t, got.State)
		assert.NotNil(t, got.BridgeRefundGasPrice)
		assert.NotNil(t, got.RefundPegoutGasPrice)
		assert.NotNil(t, got.SendPegoutBtcFee)
	}
}

func TestFixtures_LPConfigurations_Readable(t *testing.T) {
	ctx := requireRestoredFixtures(t)

	peginCfg, err := lpRepo.GetPeginConfiguration(ctx)
	require.NoError(t, err)
	require.NotNil(t, peginCfg)
	assert.NotNil(t, peginCfg.Value.PenaltyFee)

	pegoutCfg, err := lpRepo.GetPegoutConfiguration(ctx)
	require.NoError(t, err)
	require.NotNil(t, pegoutCfg)
	assert.NotNil(t, pegoutCfg.Value.PenaltyFee)

	generalCfg, err := lpRepo.GetGeneralConfiguration(ctx)
	require.NoError(t, err)
	require.NotNil(t, generalCfg)
	assert.NotEmpty(t, generalCfg.Value.RskConfirmations)

	creds, err := lpRepo.GetCredentials(ctx)
	require.NoError(t, err)
	require.NotNil(t, creds)
	assert.NotEmpty(t, creds.Value.HashedUsername)
}

func TestFixtures_TrustedAccounts_Readable(t *testing.T) {
	ctx := requireRestoredFixtures(t)

	all, err := trustedRepo.GetAllTrustedAccounts(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, all)

	for _, a := range all {
		assert.NotEmpty(t, a.Value.Address)
		assert.NotNil(t, a.Value.BtcLockingCap)
		assert.NotNil(t, a.Value.RbtcLockingCap)
	}
}

func TestFixtures_Deposits_Readable(t *testing.T) {
	ctx := requireRestoredFixtures(t)

	deposits, err := pegoutRepo.ListPegoutDepositsByAddress(ctx, "0x1234567890abcdef1234567890abcdef12345678")
	require.NoError(t, err)
	require.NotEmpty(t, deposits)

	for _, d := range deposits {
		assert.NotEmpty(t, d.TxHash)
		assert.NotNil(t, d.Amount)
	}
}

func TestFixtures_Penalizations_Readable(t *testing.T) {
	ctx := requireRestoredFixtures(t)

	// Get all penalization hashes from raw collection
	coll := rawCollection(mongo.PenalizedEventCollection)
	cursor, err := coll.Find(ctx, bson.M{})
	require.NoError(t, err)
	defer cursor.Close(ctx)

	type storedPenalization struct {
		QuoteHash string `bson:"quote_hash"`
	}
	var items []storedPenalization
	require.NoError(t, cursor.All(ctx, &items))
	require.NotEmpty(t, items)

	hashes := make([]string, len(items))
	for i, item := range items {
		hashes[i] = item.QuoteHash
	}

	results, err := penaltyRepo.GetPenalizationsByQuoteHashes(ctx, hashes)
	require.NoError(t, err)
	require.Len(t, results, len(items))

	want := make(map[string]struct{}, len(hashes))
	for _, h := range hashes {
		want[h] = struct{}{}
	}
	for _, r := range results {
		_, ok := want[r.QuoteHash]
		assert.True(t, ok, "unexpected penalization quote hash returned: %s", r.QuoteHash)
		delete(want, r.QuoteHash)
	}
	assert.Empty(t, want, "missing penalizations for some requested quote hashes")
}

// --- Layer 2: Hand-crafted raw BSON legacy shape tests ---

func TestBSON_Wei_StoredAsInt64(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash := randomHash()
	doc := bson.D{
		{Key: "hash", Value: hash},
		{Key: "fed_address", Value: "2N1234567890abcdef"},
		{Key: "lbc_address", Value: "0xaabbccdd11223344556677889900aabbccddeeff"},
		{Key: "lp_rsk_address", Value: "0x1234567890abcdef1234567890abcdef12345678"},
		{Key: "btc_refund_address", Value: "mzBc4XEFSdzCDcTxAgf6EZXgsZWpztRhef"},
		{Key: "rsk_refund_address", Value: "0xabcdef1234567890abcdef1234567890abcdef12"},
		{Key: "lp_btc_address", Value: "mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6"},
		{Key: "call_fee", Value: int64(1000000)},
		{Key: "penalty_fee", Value: int64(5000000)},
		{Key: "contract_address", Value: "0x1234567890abcdef1234567890abcdef12345678"},
		{Key: "data", Value: "0x"},
		{Key: "gas_limit", Value: int32(100000)},
		{Key: "nonce", Value: int64(42)},
		{Key: "value", Value: int64(500000000)},
		{Key: "agreement_timestamp", Value: int32(1700000000)},
		{Key: "time_for_deposit", Value: int32(3600)},
		{Key: "lp_call_time", Value: int32(7200)},
		{Key: "confirmations", Value: int32(10)},
		{Key: "call_on_register", Value: false},
		{Key: "gas_fee", Value: int64(21000)},
		{Key: "chain_id", Value: int64(31)},
	}

	coll := rawCollection(mongo.PeginQuoteCollection)
	_, err := coll.InsertOne(ctx, doc)
	require.NoError(t, err)

	got, err := peginRepo.GetQuote(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, got)
	assertWeiEq(t, 1000000, got.CallFee)
	assertWeiEq(t, 500000000, got.Value)
	assertWeiEq(t, 21000, got.GasFee)
}

func TestBSON_Wei_StoredAsString(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash := randomHash()
	doc := bson.D{
		{Key: "hash", Value: hash},
		{Key: "fed_address", Value: "2N1234567890abcdef"},
		{Key: "lbc_address", Value: "0xaabbccdd11223344556677889900aabbccddeeff"},
		{Key: "lp_rsk_address", Value: "0x1234567890abcdef1234567890abcdef12345678"},
		{Key: "btc_refund_address", Value: "mzBc4XEFSdzCDcTxAgf6EZXgsZWpztRhef"},
		{Key: "rsk_refund_address", Value: "0xabcdef1234567890abcdef1234567890abcdef12"},
		{Key: "lp_btc_address", Value: "mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6"},
		{Key: "call_fee", Value: "1000000"},
		{Key: "penalty_fee", Value: "5000000"},
		{Key: "contract_address", Value: "0x1234567890abcdef1234567890abcdef12345678"},
		{Key: "data", Value: "0x"},
		{Key: "gas_limit", Value: int32(100000)},
		{Key: "nonce", Value: int64(42)},
		{Key: "value", Value: "500000000000000000"},
		{Key: "agreement_timestamp", Value: int32(1700000000)},
		{Key: "time_for_deposit", Value: int32(3600)},
		{Key: "lp_call_time", Value: int32(7200)},
		{Key: "confirmations", Value: int32(10)},
		{Key: "call_on_register", Value: false},
		{Key: "gas_fee", Value: "21000"},
		{Key: "chain_id", Value: int64(31)},
	}

	coll := rawCollection(mongo.PeginQuoteCollection)
	_, err := coll.InsertOne(ctx, doc)
	require.NoError(t, err)

	got, err := peginRepo.GetQuote(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, got)
	assertWeiEq(t, 1000000, got.CallFee)
	require.NotNil(t, got.Value)
	assert.Zero(t, got.Value.Cmp(entities.NewWei(500000000000000000)))
}

func TestBSON_Wei_StoredAsNull(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash := randomHash()
	doc := bson.D{
		{Key: "quote_hash", Value: hash},
		{Key: "deposit_address", Value: "2N1234567890abcdef"},
		{Key: "signature", Value: "0xsig"},
		{Key: "required_liquidity", Value: "600000000"},
		{Key: "state", Value: "WaitingForDeposit"},
		{Key: "user_btc_tx_hash", Value: ""},
		{Key: "call_for_user_tx_hash", Value: ""},
		{Key: "register_pegin_tx_hash", Value: ""},
		{Key: "call_for_user_gas_used", Value: int64(0)},
		{Key: "call_for_user_gas_price", Value: nil},
		{Key: "register_pegin_gas_used", Value: int64(0)},
		{Key: "register_pegin_gas_price", Value: nil},
		{Key: "owner_account_address", Value: "0x1234567890abcdef1234567890abcdef12345678"},
	}

	coll := rawCollection(mongo.RetainedPeginQuoteCollection)
	_, err := coll.InsertOne(ctx, doc)
	require.NoError(t, err)

	got, err := peginRepo.GetRetainedQuote(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, got)
	// FillZeroValues should set nil gas prices to NewWei(0)
	assert.NotNil(t, got.CallForUserGasPrice)
	assert.NotNil(t, got.RegisterPeginGasPrice)
	assertWeiEq(t, 0, got.CallForUserGasPrice)
	assertWeiEq(t, 0, got.RegisterPeginGasPrice)
}

func TestBSON_Wei_StoredAsNilString(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash := randomHash()
	doc := bson.D{
		{Key: "hash", Value: hash},
		{Key: "fed_address", Value: "2N1234567890abcdef"},
		{Key: "lbc_address", Value: "0xaabbccdd11223344556677889900aabbccddeeff"},
		{Key: "lp_rsk_address", Value: "0x1234567890abcdef1234567890abcdef12345678"},
		{Key: "btc_refund_address", Value: "mzBc4XEFSdzCDcTxAgf6EZXgsZWpztRhef"},
		{Key: "rsk_refund_address", Value: "0xabcdef1234567890abcdef1234567890abcdef12"},
		{Key: "lp_btc_address", Value: "mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6"},
		{Key: "call_fee", Value: "<nil>"},
		{Key: "penalty_fee", Value: "5000000"},
		{Key: "contract_address", Value: "0x1234567890abcdef1234567890abcdef12345678"},
		{Key: "data", Value: "0x"},
		{Key: "gas_limit", Value: int32(100000)},
		{Key: "nonce", Value: int64(42)},
		{Key: "value", Value: "500000000"},
		{Key: "agreement_timestamp", Value: int32(1700000000)},
		{Key: "time_for_deposit", Value: int32(3600)},
		{Key: "lp_call_time", Value: int32(7200)},
		{Key: "confirmations", Value: int32(10)},
		{Key: "call_on_register", Value: false},
		{Key: "gas_fee", Value: "21000"},
		{Key: "chain_id", Value: int64(31)},
	}

	coll := rawCollection(mongo.PeginQuoteCollection)
	_, err := coll.InsertOne(ctx, doc)
	require.NoError(t, err)

	got, err := peginRepo.GetQuote(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, got)
	// "<nil>" string should be handled gracefully — CallFee should be zero-valued
	assertWeiEq(t, 0, got.CallFee)
}

func TestBSON_BigFloat_StoredAsDouble(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash := randomHash()
	doc := bson.D{
		{Key: "hash", Value: hash},
		{Key: "gas_price", Value: "60000000"},
		{Key: "fee_percentage", Value: 0.5},
		{Key: "fixed_fee", Value: "1000000"},
	}

	coll := rawCollection(mongo.PeginCreationDataCollection)
	_, err := coll.InsertOne(ctx, doc)
	require.NoError(t, err)

	got := peginRepo.GetPeginCreationData(ctx, hash)
	assertWeiEq(t, 60000000, got.GasPrice)
	assertWeiEq(t, 1000000, got.FixedFee)
	val, _ := got.FeePercentage.Native().Float64()
	assert.InDelta(t, 0.5, val, 0.001)
}

func TestBSON_RetainedPeginQuote_MissingGasFields(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash := randomHash()
	doc := bson.D{
		{Key: "quote_hash", Value: hash},
		{Key: "deposit_address", Value: "2N1234567890abcdef"},
		{Key: "signature", Value: "0xsig"},
		{Key: "required_liquidity", Value: "600000000"},
		{Key: "state", Value: "WaitingForDeposit"},
		{Key: "user_btc_tx_hash", Value: "btctx123"},
		{Key: "call_for_user_tx_hash", Value: ""},
		{Key: "register_pegin_tx_hash", Value: ""},
		{Key: "call_for_user_gas_used", Value: int64(0)},
		// call_for_user_gas_price is MISSING
		{Key: "register_pegin_gas_used", Value: int64(0)},
		// register_pegin_gas_price is MISSING
		{Key: "owner_account_address", Value: "0x1234567890abcdef1234567890abcdef12345678"},
	}

	coll := rawCollection(mongo.RetainedPeginQuoteCollection)
	_, err := coll.InsertOne(ctx, doc)
	require.NoError(t, err)

	got, err := peginRepo.GetRetainedQuote(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.NotNil(t, got.CallForUserGasPrice)
	assert.NotNil(t, got.RegisterPeginGasPrice)
	assertWeiEq(t, 0, got.CallForUserGasPrice)
	assertWeiEq(t, 0, got.RegisterPeginGasPrice)
}

func TestBSON_RetainedPegoutQuote_MissingGasFields(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	hash := randomHash()
	doc := bson.D{
		{Key: "quote_hash", Value: hash},
		{Key: "deposit_address", Value: "0xdeposit123"},
		{Key: "signature", Value: "0xsig"},
		{Key: "required_liquidity", Value: "1200000000"},
		{Key: "state", Value: "WaitingForDeposit"},
		{Key: "user_rsk_tx_hash", Value: ""},
		{Key: "lp_btc_tx_hash", Value: ""},
		{Key: "refund_pegout_tx_hash", Value: ""},
		{Key: "bridge_refund_tx_hash", Value: ""},
		{Key: "bridge_refund_gas_used", Value: int64(0)},
		// bridge_refund_gas_price MISSING
		{Key: "refund_pegout_gas_used", Value: int64(0)},
		// refund_pegout_gas_price MISSING
		// send_pegout_btc_fee MISSING
		{Key: "btc_release_tx_hash", Value: ""},
		{Key: "owner_account_address", Value: "0x1234567890abcdef1234567890abcdef12345678"},
	}

	coll := rawCollection(mongo.RetainedPegoutQuoteCollection)
	_, err := coll.InsertOne(ctx, doc)
	require.NoError(t, err)

	got, err := pegoutRepo.GetRetainedQuote(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.NotNil(t, got.BridgeRefundGasPrice)
	assert.NotNil(t, got.RefundPegoutGasPrice)
	assert.NotNil(t, got.SendPegoutBtcFee)
	assertWeiEq(t, 0, got.BridgeRefundGasPrice)
	assertWeiEq(t, 0, got.RefundPegoutGasPrice)
	assertWeiEq(t, 0, got.SendPegoutBtcFee)
}

func TestBSON_PegoutDeposit_WeiAmount(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	doc := bson.D{
		{Key: "tx_hash", Value: "0x" + randomHash()},
		{Key: "quote_hash", Value: randomHash()},
		{Key: "amount", Value: int64(1000000000)},
		{Key: "timestamp", Value: "2024-01-15T10:30:00Z"},
		{Key: "block_number", Value: int64(50000)},
		{Key: "from", Value: "0xTestAddress1234567890abcdef1234567890abcd"},
	}

	coll := rawCollection(mongo.DepositEventsCollection)
	_, err := coll.InsertOne(ctx, doc)
	require.NoError(t, err)

	deposits, err := pegoutRepo.ListPegoutDepositsByAddress(ctx, "0xTestAddress1234567890abcdef1234567890abcd")
	require.NoError(t, err)
	require.Len(t, deposits, 1)
	assert.NotNil(t, deposits[0].Amount)
	assertWeiEq(t, 1000000000, deposits[0].Amount)
}
