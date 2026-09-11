package mongo_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDb "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

//nolint:funlen // One fixture shared across the round-trip scenarios keeps the event identity identical in each.
func TestPegInWatchMongoRepository(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond)
	entry := rootstock.PegInWatch{
		TxHash:      "0x1234",
		LogIndex:    7,
		BlockNumber: 100,
		RskAddress:  "0xabcd",
		State:       rootstock.PegInWatchDiscovered,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	identity := bson.M{"rsk_address": entry.RskAddress}

	t.Run("upserts with set-on-insert", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		collection.EXPECT().UpdateOne(
			mock.Anything,
			identity,
			bson.M{"$setOnInsert": entry},
			withUpdateUpsert(),
		).Return(&mongoDb.UpdateResult{UpsertedCount: 1}, nil).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		require.NoError(t, repo.Upsert(context.Background(), entry))
		collection.AssertExpectations(t)
	})

	t.Run("returns nil on a duplicate-key upsert error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		collection.EXPECT().UpdateOne(mock.Anything, identity, mock.Anything, mock.Anything).
			Return(nil, mongoDb.WriteException{WriteErrors: []mongoDb.WriteError{{Code: 11000}}}).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		require.NoError(t, repo.Upsert(context.Background(), entry))
		collection.AssertExpectations(t)
	})

	t.Run("returns the UpdateOne error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		collection.EXPECT().UpdateOne(mock.Anything, identity, mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		require.ErrorIs(t, repo.Upsert(context.Background(), entry), assert.AnError)
		collection.AssertExpectations(t)
	})

	t.Run("returns the document matching the identity", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		collection.EXPECT().FindOne(mock.Anything, identity).
			Return(mongoDb.NewSingleResultFromDocument(entry, nil, nil)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		result, err := repo.Get(context.Background(), entry.RskAddress)
		require.NoError(t, err)
		assert.Equal(t, &entry, result)
		collection.AssertExpectations(t)
	})

	t.Run("returns nil when no document matches the identity", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		collection.EXPECT().FindOne(mock.Anything, identity).
			Return(mongoDb.NewSingleResultFromDocument(rootstock.PegInWatch{}, mongoDb.ErrNoDocuments, nil)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		result, err := repo.Get(context.Background(), entry.RskAddress)
		require.NoError(t, err)
		assert.Nil(t, result)
		collection.AssertExpectations(t)
	})

	t.Run("returns the FindOne error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		collection.EXPECT().FindOne(mock.Anything, identity).
			Return(mongoDb.NewSingleResultFromDocument(rootstock.PegInWatch{}, assert.AnError, nil)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		result, err := repo.Get(context.Background(), entry.RskAddress)
		require.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, result)
		collection.AssertExpectations(t)
	})

	t.Run("lists documents sorted by block number and log index", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		second := entry
		second.LogIndex = 9
		collection.EXPECT().Find(
			mock.Anything,
			bson.M{"rsk_address": bson.M{"$exists": true}},
			sortedBy(bson.D{{Key: "block_number", Value: 1}, {Key: "log_index", Value: 1}}),
		).Return(mongoDb.NewCursorFromDocuments([]any{entry, second}, nil, nil)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		result, err := repo.List(context.Background())
		require.NoError(t, err)
		assert.Equal(t, []rootstock.PegInWatch{entry, second}, result)
		collection.AssertExpectations(t)
	})

	t.Run("returns the Find error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		collection.EXPECT().Find(
			mock.Anything,
			bson.M{"rsk_address": bson.M{"$exists": true}},
			sortedBy(bson.D{{Key: "block_number", Value: 1}, {Key: "log_index", Value: 1}}),
		).Return(nil, assert.AnError).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		result, err := repo.List(context.Background())
		require.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, result)
		collection.AssertExpectations(t)
	})

	t.Run("replaces the document when one row matches", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		imported := entry
		imported.State = rootstock.PegInWatchImported
		collection.EXPECT().ReplaceOne(mock.Anything, identity, imported).
			Return(&mongoDb.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		require.NoError(t, repo.Update(context.Background(), imported))
		collection.AssertExpectations(t)
	})

	t.Run("returns not-found when ReplaceOne matches nothing", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		imported := entry
		imported.State = rootstock.PegInWatchImported
		collection.EXPECT().ReplaceOne(mock.Anything, identity, imported).
			Return(&mongoDb.UpdateResult{MatchedCount: 0}, nil).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		require.ErrorContains(t, repo.Update(context.Background(), imported), "pegin watch not found")
		collection.AssertExpectations(t)
	})

	t.Run("returns the ReplaceOne error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		imported := entry
		imported.State = rootstock.PegInWatchImported
		collection.EXPECT().ReplaceOne(mock.Anything, identity, imported).
			Return(nil, assert.AnError).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		require.ErrorIs(t, repo.Update(context.Background(), imported), assert.AnError)
		collection.AssertExpectations(t)
	})

	t.Run("replaces a suffix", testPegInWatchMongoRepositoryReplaceFromBlock)

	t.Run("updates and reads a document in unsupported encoding state", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		unsupported := entry
		unsupported.State = rootstock.PegInWatchUnsupportedEncoding
		unsupported.Encoding = 1
		collection.EXPECT().ReplaceOne(mock.Anything, identity, unsupported).
			Return(&mongoDb.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil).Once()
		collection.EXPECT().FindOne(mock.Anything, identity).
			Return(mongoDb.NewSingleResultFromDocument(unsupported, nil, nil)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		require.NoError(t, repo.Update(context.Background(), unsupported))
		result, err := repo.Get(context.Background(), unsupported.RskAddress)
		require.NoError(t, err)
		assert.Equal(t, &unsupported, result)
		collection.AssertExpectations(t)
	})
}

type sessionContextKey struct{}

func testPegInWatchMongoRepositoryReplaceFromBlock(t *testing.T) {
	t.Run("replaces the suffix in one transaction", testReplaceFromBlockTransaction)
	t.Run("returns an insert error so the transaction can roll back", testReplaceFromBlockRollback)
}

func testReplaceFromBlockTransaction(t *testing.T) {
	filter := bson.M{
		"rsk_address":  bson.M{"$exists": true},
		"block_number": bson.M{"$gte": uint64(101)},
	}
	replacement := rootstock.PegInWatch{
		BlockNumber: 101,
		RskAddress:  "0xreplacement",
	}
	client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
	session := &mocks.SessionBindingMock{}
	sessionContext := context.WithValue(context.Background(), sessionContextKey{}, "session")
	client.On("StartSession").Return(session, nil).Once()
	session.On("EndSession", mock.Anything).Return().Once()
	session.On("WithTransaction", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn, ok := args.Get(1).(func(context.Context) (any, error))
			require.True(t, ok)
			result, callbackErr := fn(sessionContext)
			require.NoError(t, callbackErr)
			assert.Nil(t, result)
		}).
		Return(nil, nil).Once()
	collection.EXPECT().DeleteMany(sessionContext, filter).
		Return(&mongoDb.DeleteResult{DeletedCount: 2}, nil).Once()
	collection.EXPECT().UpdateOne(
		sessionContext,
		rskAddressIdentity(replacement.RskAddress),
		bson.M{"$setOnInsert": replacement},
		withUpdateUpsert(),
	).Return(&mongoDb.UpdateResult{UpsertedCount: 1}, nil).Once()

	repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
	require.NoError(t, repo.ReplaceFromBlock(context.Background(), 101, []rootstock.PegInWatch{replacement}))
	collection.AssertExpectations(t)
	client.AssertExpectations(t)
	session.AssertExpectations(t)
}

func testReplaceFromBlockRollback(t *testing.T) {
	filter := bson.M{
		"rsk_address":  bson.M{"$exists": true},
		"block_number": bson.M{"$gte": uint64(101)},
	}
	replacement := rootstock.PegInWatch{
		BlockNumber: 101,
		RskAddress:  "0xreplacement",
	}
	client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
	session := &mocks.SessionBindingMock{}
	client.On("StartSession").Return(session, nil).Once()
	session.On("EndSession", mock.Anything).Return().Once()
	session.On("WithTransaction", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn, ok := args.Get(1).(func(context.Context) (any, error))
			require.True(t, ok)
			result, callbackErr := fn(context.Background())
			require.ErrorIs(t, callbackErr, assert.AnError)
			assert.Nil(t, result)
		}).
		Return(nil, assert.AnError).Once()
	collection.EXPECT().DeleteMany(mock.Anything, filter).
		Return(&mongoDb.DeleteResult{DeletedCount: 2}, nil).Once()
	collection.EXPECT().UpdateOne(
		mock.Anything,
		rskAddressIdentity(replacement.RskAddress),
		bson.M{"$setOnInsert": replacement},
		withUpdateUpsert(),
	).Return(nil, assert.AnError).Once()

	repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
	require.ErrorIs(t,
		repo.ReplaceFromBlock(context.Background(), 101, []rootstock.PegInWatch{replacement}),
		assert.AnError,
	)
	collection.AssertExpectations(t)
	client.AssertExpectations(t)
	session.AssertExpectations(t)
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

func rskAddressIdentity(rskAddress string) bson.M {
	return bson.M{"rsk_address": rskAddress}
}
