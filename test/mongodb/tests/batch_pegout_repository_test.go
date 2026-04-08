//go:build integration

package mongodb_test

import (
	"context"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test/mongodb/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

// countDocuments returns the number of documents matching filter.
// mongo.CollectionBinding does not expose CountDocuments; Aggregate is available, so we use
// $match + $count instead of reaching past the binding to the underlying *mongo.Collection.
func countDocuments(t *testing.T, ctx context.Context, coll mongo.CollectionBinding, filter bson.M) int64 {
	t.Helper()

	cursor, err := coll.Aggregate(ctx, bson.A{
		bson.M{"$match": filter},
		bson.M{"$count": "count"},
	})
	require.NoError(t, err)
	defer cursor.Close(ctx)

	var result struct {
		Count int64 `bson:"count"`
	}
	if cursor.Next(ctx) {
		require.NoError(t, cursor.Decode(&result))
	}
	require.NoError(t, cursor.Err())
	return result.Count
}

func TestBatch_Upsert_InsertThenUpdate(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	txHash := "0x" + utils.RandomHash()
	batch := utils.NewTestBatchPegOut(txHash)
	err := batchRepo.UpsertBatch(ctx, batch)
	require.NoError(t, err)

	// Verify via raw collection read
	coll := rawCollection(mongo.BatchPegOutEventsCollection)
	var inserted rootstock.BatchPegOut
	err = coll.FindOne(ctx, bson.M{"transaction_hash": txHash}).Decode(&inserted)
	require.NoError(t, err)
	assert.Equal(t, batch.TransactionHash, inserted.TransactionHash)
	assert.Equal(t, batch.BlockHash, inserted.BlockHash)
	assert.Equal(t, batch.BlockNumber, inserted.BlockNumber)
	assert.Equal(t, batch.BtcTxHash, inserted.BtcTxHash)
	assert.Equal(t, batch.ReleaseRskTxHashes, inserted.ReleaseRskTxHashes)

	// Upsert with updated fields
	batch.BlockNumber = 100001
	batch.BtcTxHash = "updated_btc_tx"
	err = batchRepo.UpsertBatch(ctx, batch)
	require.NoError(t, err)

	count := countDocuments(t, ctx, coll, bson.M{"transaction_hash": txHash})
	assert.Equal(t, int64(1), count)

	var updated rootstock.BatchPegOut
	err = coll.FindOne(ctx, bson.M{"transaction_hash": txHash}).Decode(&updated)
	require.NoError(t, err)
	assert.Equal(t, uint64(100001), updated.BlockNumber)
	assert.Equal(t, "updated_btc_tx", updated.BtcTxHash)
}

func TestBatch_UpsertIdempotent(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	txHash := "0x" + utils.RandomHash()
	batch := utils.NewTestBatchPegOut(txHash)

	err := batchRepo.UpsertBatch(ctx, batch)
	require.NoError(t, err)

	err = batchRepo.UpsertBatch(ctx, batch)
	require.NoError(t, err)

	coll := rawCollection(mongo.BatchPegOutEventsCollection)
	count := countDocuments(t, ctx, coll, bson.M{"transaction_hash": txHash})
	assert.Equal(t, int64(1), count)

	var result rootstock.BatchPegOut
	err = coll.FindOne(ctx, bson.M{"transaction_hash": txHash}).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, batch.TransactionHash, result.TransactionHash)
	assert.Equal(t, batch.BlockHash, result.BlockHash)
	assert.Equal(t, batch.BlockNumber, result.BlockNumber)
	assert.Equal(t, batch.BtcTxHash, result.BtcTxHash)
	assert.Equal(t, batch.ReleaseRskTxHashes, result.ReleaseRskTxHashes)
}

func TestBatch_Upsert_SameTransactionHashUpdatesExistingDocument(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	txHash := "0x" + utils.RandomHash()
	batch1 := utils.NewTestBatchPegOut(txHash)
	batch1.BlockNumber = 1

	batch2 := utils.NewTestBatchPegOut(txHash)
	batch2.BlockNumber = 2

	require.NoError(t, batchRepo.UpsertBatch(ctx, batch1))
	require.NoError(t, batchRepo.UpsertBatch(ctx, batch2))

	var result rootstock.BatchPegOut
	coll := rawCollection(mongo.BatchPegOutEventsCollection)
	count := countDocuments(t, ctx, coll, bson.M{"transaction_hash": txHash})
	assert.Equal(t, int64(1), count)

	err := coll.FindOne(ctx, bson.M{"transaction_hash": txHash}).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), result.BlockNumber)
}
