package mongo_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDb "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

//nolint:funlen // One fixture shared across the round-trip scenarios keeps the event identity identical in each.
func TestPegInAddressRegistryWatchMongoRepository(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond)
	entry := rootstock.PegInAddressRegistryWatchEntry{
		TxHash:      "0x1234",
		LogIndex:    7,
		BlockNumber: 100,
		RskAddress:  "0xabcd",
		State:       rootstock.PegInAddressRegistryWatchDiscovered,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	identity := bson.M{"rsk_address": entry.RskAddress}

	t.Run("upserts an event once", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInAddressRegistryWatchCollection)
		collection.EXPECT().UpdateOne(
			mock.Anything,
			identity,
			bson.M{"$setOnInsert": entry},
			withUpdateUpsert(),
		).Return(&mongoDb.UpdateResult{UpsertedCount: 1}, nil).Once()

		repo := mongo.NewPegInAddressRegistryWatchMongoRepository(mongo.NewConnection(client, time.Second))
		require.NoError(t, repo.Upsert(context.Background(), entry))
		collection.AssertExpectations(t)
	})

	t.Run("treats a duplicate-key race as an idempotent upsert", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInAddressRegistryWatchCollection)
		collection.EXPECT().UpdateOne(mock.Anything, identity, mock.Anything, mock.Anything).
			Return(nil, mongoDb.WriteException{WriteErrors: []mongoDb.WriteError{{Code: 11000}}}).Once()

		repo := mongo.NewPegInAddressRegistryWatchMongoRepository(mongo.NewConnection(client, time.Second))
		require.NoError(t, repo.Upsert(context.Background(), entry))
		collection.AssertExpectations(t)
	})

	t.Run("reads an event by stable identity", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInAddressRegistryWatchCollection)
		collection.EXPECT().FindOne(mock.Anything, identity).
			Return(mongoDb.NewSingleResultFromDocument(entry, nil, nil)).Once()

		repo := mongo.NewPegInAddressRegistryWatchMongoRepository(mongo.NewConnection(client, time.Second))
		result, err := repo.Get(context.Background(), entry.RskAddress)
		require.NoError(t, err)
		assert.Equal(t, &entry, result)
		collection.AssertExpectations(t)
	})

	t.Run("lists events in processing order", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInAddressRegistryWatchCollection)
		second := entry
		second.LogIndex = 9
		collection.EXPECT().Find(
			mock.Anything,
			bson.M{"rsk_address": bson.M{"$exists": true}},
			sortedBy(bson.D{{Key: "block_number", Value: 1}, {Key: "log_index", Value: 1}}),
		).Return(mongoDb.NewCursorFromDocuments([]any{entry, second}, nil, nil)).Once()

		repo := mongo.NewPegInAddressRegistryWatchMongoRepository(mongo.NewConnection(client, time.Second))
		result, err := repo.List(context.Background())
		require.NoError(t, err)
		assert.Equal(t, []rootstock.PegInAddressRegistryWatchEntry{entry, second}, result)
		collection.AssertExpectations(t)
	})

	t.Run("persists a state transition", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInAddressRegistryWatchCollection)
		imported := entry
		imported.State = rootstock.PegInAddressRegistryWatchImported
		collection.EXPECT().ReplaceOne(mock.Anything, identity, imported).
			Return(&mongoDb.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil).Once()

		repo := mongo.NewPegInAddressRegistryWatchMongoRepository(mongo.NewConnection(client, time.Second))
		require.NoError(t, repo.Update(context.Background(), imported))
		collection.AssertExpectations(t)
	})
}

func TestPegInAddressRegistryWatchMongoRepository_Cursor(t *testing.T) {
	t.Run("returns an absent cursor for an empty watch set", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInAddressRegistryWatchCollection)
		collection.EXPECT().FindOne(mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(bson.M{}, mongoDb.ErrNoDocuments, nil)).Once()

		repo := mongo.NewPegInAddressRegistryWatchMongoRepository(mongo.NewConnection(client, time.Second))
		block, found, err := repo.GetCursor(context.Background())
		require.NoError(t, err)
		assert.Zero(t, block)
		assert.False(t, found)
	})

	t.Run("writes and reads the scan cursor", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInAddressRegistryWatchCollection)
		collection.EXPECT().UpdateOne(
			mock.Anything,
			bson.M{"_id": "scanCursor"},
			bson.M{"$set": bson.M{"last_scanned_block": uint64(123)}},
			withUpdateUpsert(),
		).Return(&mongoDb.UpdateResult{UpsertedCount: 1}, nil).Once()
		collection.EXPECT().FindOne(mock.Anything, bson.M{"_id": "scanCursor"}).
			Return(mongoDb.NewSingleResultFromDocument(
				bson.M{"_id": "scanCursor", "last_scanned_block": uint64(123)},
				nil,
				nil,
			)).Once()

		repo := mongo.NewPegInAddressRegistryWatchMongoRepository(mongo.NewConnection(client, time.Second))
		require.NoError(t, repo.SetCursor(context.Background(), 123))
		block, found, err := repo.GetCursor(context.Background())
		require.NoError(t, err)
		assert.Equal(t, uint64(123), block)
		assert.True(t, found)
		collection.AssertExpectations(t)
	})
}

func withUpdateUpsert() interface{} {
	return mock.MatchedBy(func(opt options.Lister[options.UpdateOneOptions]) bool {
		resolved := &options.UpdateOneOptions{}
		for _, fn := range opt.List() {
			if err := fn(resolved); err != nil {
				return false
			}
		}
		return resolved.Upsert != nil && *resolved.Upsert
	})
}
