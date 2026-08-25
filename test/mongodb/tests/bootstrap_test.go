//go:build integration

package mongodb_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/test/mongodb/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodrv "go.mongodb.org/mongo-driver/v2/mongo"
)

func TestBootstrap_IndexesCreated(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type expectedIndex struct {
		collection string
		field      string
		unique     bool
	}
	expectedIndexes := []expectedIndex{
		{collection: mongo.DepositEventsCollection, field: "tx_hash", unique: true},
		{collection: mongo.TrustedAccountCollection, field: "address", unique: true},
		{collection: mongo.BatchPegOutEventsCollection, field: "transaction_hash", unique: true},
		{collection: mongo.RetainedPegoutQuoteCollection, field: "bridge_rebalances.tx_hash", unique: false},
		{collection: mongo.RetainedPeginQuoteCollection, field: "state", unique: false},
		{collection: mongo.RetainedPegoutQuoteCollection, field: "state", unique: false},
		{collection: mongo.PeginQuoteCollection, field: "agreement_timestamp", unique: false},
		{collection: mongo.PegoutQuoteCollection, field: "agreement_timestamp", unique: false},
		{collection: mongo.RetainedPeginQuoteCollection, field: "quote_hash", unique: false},
		{collection: mongo.RetainedPegoutQuoteCollection, field: "quote_hash", unique: false},
		{collection: mongo.PeginQuoteCollection, field: "hash", unique: false},
		{collection: mongo.PegoutQuoteCollection, field: "hash", unique: false},
		{collection: mongo.PegInWatchCollection, field: "rsk_address", unique: true},
		{collection: mongo.PegInWatchCollection, field: "state", unique: false},
	}

	for _, expected := range expectedIndexes {
		t.Run(expected.collection+"."+expected.field, func(t *testing.T) {
			db := mongoClient.Database(testDbName)
			cursor, err := db.Collection(expected.collection).Indexes().List(ctx)
			require.NoError(t, err, "listing indexes for %s", expected.collection)
			defer func() { _ = cursor.Close(ctx) }()

			var indexes []map[string]any
			require.NoError(t, cursor.All(ctx, &indexes))

			i := slices.IndexFunc(indexes, func(idx map[string]any) bool {
				return utils.IndexKeysContainField(idx["key"], expected.field)
			})
			require.GreaterOrEqual(t, i, 0, "index on %s.%s not found", expected.collection, expected.field)

			idx := indexes[i]
			assert.True(t,
				utils.IndexKeysIsSingleFieldAscending(idx["key"], expected.field),
				"index on %s.%s must be single-field ascending, got key=%v", expected.collection, expected.field, idx["key"])
			unique, _ := idx["unique"].(bool)
			assert.Equal(t, expected.unique, unique, "unique flag on %s.%s", expected.collection, expected.field)
		})
	}
}

func TestBootstrap_RegistryConstruction(t *testing.T) {
	db := registry.NewDatabaseRegistry(conn)
	assert.NotNil(t, db.PeginRepository)
	assert.NotNil(t, db.PegoutRepository)
	assert.NotNil(t, db.LiquidityProviderRepository)
	assert.NotNil(t, db.TrustedAccountRepository)
	assert.NotNil(t, db.PenalizedEventRepository)
	assert.NotNil(t, db.BatchPegOutRepository)
	assert.NotNil(t, db.PegInWatchRepository)
	assert.NotNil(t, db.Connection)
}

func TestBootstrap_UniqueIndex_DepositEvents(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	coll := rawCollection(mongo.DepositEventsCollection)
	doc := map[string]any{"tx_hash": "duplicate_tx_hash", "quote_hash": "qh1"}
	_, err := coll.InsertOne(ctx, doc)
	require.NoError(t, err)

	_, err = coll.InsertOne(ctx, doc)
	assert.Error(t, err, "second insert with same tx_hash should fail due to unique index")
	require.True(t, mongodrv.IsDuplicateKeyError(err), "expected duplicate key error, got: %v", err)

	count, countErr := mongoClient.Database(testDbName).
		Collection(mongo.DepositEventsCollection).
		CountDocuments(ctx, bson.M{"tx_hash": "duplicate_tx_hash"})
	require.NoError(t, countErr)
	assert.Equal(t, int64(1), count)
}

func TestBootstrap_UniqueIndex_TrustedAccounts(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	coll := rawCollection(mongo.TrustedAccountCollection)
	doc := map[string]any{"address": "0xduplicateaddress", "name": "test"}
	_, err := coll.InsertOne(ctx, doc)
	require.NoError(t, err)

	_, err = coll.InsertOne(ctx, doc)
	assert.Error(t, err, "second insert with same address should fail due to unique index")
	require.True(t, mongodrv.IsDuplicateKeyError(err), "expected duplicate key error, got: %v", err)

	count, countErr := mongoClient.Database(testDbName).
		Collection(mongo.TrustedAccountCollection).
		CountDocuments(ctx, bson.M{"address": "0xduplicateaddress"})
	require.NoError(t, countErr)
	assert.Equal(t, int64(1), count)
}

func TestBootstrap_UniqueIndex_BatchPegOutEvents(t *testing.T) {
	cleanCollections(t)
	ctx := context.Background()

	coll := rawCollection(mongo.BatchPegOutEventsCollection)
	doc := map[string]any{"transaction_hash": "duplicate_batch_tx", "block_number": 1}
	_, err := coll.InsertOne(ctx, doc)
	require.NoError(t, err)

	_, err = coll.InsertOne(ctx, doc)
	assert.Error(t, err, "second insert with same transaction_hash should fail due to unique index")
	require.True(t, mongodrv.IsDuplicateKeyError(err), "expected duplicate key error, got: %v", err)

	count, countErr := mongoClient.Database(testDbName).
		Collection(mongo.BatchPegOutEventsCollection).
		CountDocuments(ctx, bson.M{"transaction_hash": "duplicate_batch_tx"})
	require.NoError(t, countErr)
	assert.Equal(t, int64(1), count)
}
