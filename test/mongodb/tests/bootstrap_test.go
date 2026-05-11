//go:build integration

package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/test/mongodb/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	mongodrv "go.mongodb.org/mongo-driver/mongo"
)

func TestBootstrap_IndexesCreated(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexedCollections := map[string]string{
		mongo.DepositEventsCollection:     "tx_hash",
		mongo.TrustedAccountCollection:    "address",
		mongo.BatchPegOutEventsCollection: "transaction_hash",
	}

	for collName, field := range indexedCollections {
		t.Run(collName, func(t *testing.T) {
			db := mongoClient.Database(testDbName)
			cursor, err := db.Collection(collName).Indexes().List(ctx)
			require.NoError(t, err, "listing indexes for %s", collName)
			defer func() { _ = cursor.Close(ctx) }()

			var indexes []map[string]any
			require.NoError(t, cursor.All(ctx, &indexes))

			found := false
			for _, idx := range indexes {
				if utils.IndexKeysContainField(idx["key"], field) {
					unique, _ := idx["unique"].(bool)
					assert.True(t, unique, "index on %s.%s should be unique", collName, field)
					found = true
				}
			}
			require.True(t, found, "unique index on %s.%s not found", collName, field)
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
