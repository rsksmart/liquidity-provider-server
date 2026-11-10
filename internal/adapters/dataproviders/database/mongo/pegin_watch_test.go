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

func TestPegInWatchMongoRepository_GetCheckpoint(t *testing.T) {
	checkpointIdentity := bson.M{"_id": "verified"}
	t.Run("returns none when no checkpoint is stored", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCheckpointCollection)
		collection.EXPECT().FindOne(mock.Anything, checkpointIdentity).
			Return(mongoDb.NewSingleResultFromDocument(bson.M{}, mongoDb.ErrNoDocuments, nil)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		checkpoint, found, err := repo.GetCheckpoint(context.Background())

		require.NoError(t, err)
		assert.False(t, found)
		assert.Equal(t, rootstock.PegInWatchCheckpoint{}, checkpoint)
		collection.AssertExpectations(t)
	})
	t.Run("returns the stored checkpoint", func(t *testing.T) {
		stored := rootstock.PegInWatchCheckpoint{Root: [32]byte{1, 2, 3}, Block: 123}
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCheckpointCollection)
		collection.EXPECT().FindOne(mock.Anything, checkpointIdentity).
			Return(mongoDb.NewSingleResultFromDocument(
				bson.M{"_id": "verified", "root": stored.Root, "block": stored.Block},
				nil,
				nil,
			)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		checkpoint, found, err := repo.GetCheckpoint(context.Background())

		require.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, stored, checkpoint)
		collection.AssertExpectations(t)
	})
	t.Run("returns the read error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCheckpointCollection)
		collection.EXPECT().FindOne(mock.Anything, checkpointIdentity).
			Return(mongoDb.NewSingleResultFromDocument(bson.M{}, assert.AnError, nil)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		checkpoint, found, err := repo.GetCheckpoint(context.Background())

		require.ErrorIs(t, err, assert.AnError)
		assert.False(t, found)
		assert.Equal(t, rootstock.PegInWatchCheckpoint{}, checkpoint)
		collection.AssertExpectations(t)
	})
}

type checkpointSaveMocks struct {
	client      *mocks.DbClientBindingMock
	watches     *mocks.CollectionBindingMock
	checkpoints *mocks.CollectionBindingMock
	session     *mocks.SessionBindingMock
}

func newCheckpointSaveMocks() checkpointSaveMocks {
	client, db := getClientAndDatabaseMocks()
	saveMocks := checkpointSaveMocks{
		client:      client,
		watches:     &mocks.CollectionBindingMock{},
		checkpoints: &mocks.CollectionBindingMock{},
		session:     &mocks.SessionBindingMock{},
	}
	db.On("Collection", mongo.PegInWatchCollection).Return(saveMocks.watches)
	db.On("Collection", mongo.PegInWatchCheckpointCollection).Return(saveMocks.checkpoints)
	saveMocks.session.On("EndSession", mock.Anything).Return().Once()
	return saveMocks
}

func (saveMocks checkpointSaveMocks) expectTransaction(t *testing.T, expectedErr error) {
	saveMocks.client.On("StartSession").Return(saveMocks.session, nil).Once()
	saveMocks.session.On("WithTransaction", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn, ok := args.Get(1).(func(context.Context) (any, error))
			require.True(t, ok)
			result, cbErr := fn(context.Background())
			require.ErrorIs(t, cbErr, expectedErr)
			assert.Nil(t, result)
		}).
		Return(nil, expectedErr).Once()
}

func (saveMocks checkpointSaveMocks) assertExpectations(t *testing.T) {
	saveMocks.client.AssertExpectations(t)
	saveMocks.watches.AssertExpectations(t)
	saveMocks.checkpoints.AssertExpectations(t)
	saveMocks.session.AssertExpectations(t)
}

func verifiedRegistrationsFixture() (rootstock.PegInWatch, rootstock.PegInWatch, rootstock.PegInWatchCheckpoint) {
	first := rootstock.PegInWatch{TxHash: "0x01", RskAddress: "0xaaaa", BlockNumber: 10, State: rootstock.PegInWatchDiscovered}
	second := rootstock.PegInWatch{TxHash: "0x02", RskAddress: "0xbbbb", BlockNumber: 11, State: rootstock.PegInWatchDiscovered}
	return first, second, rootstock.PegInWatchCheckpoint{Root: [32]byte{7}, Block: 11}
}

func TestPegInWatchMongoRepository_SaveVerifiedRegistrations(t *testing.T) {
	first, second, checkpoint := verifiedRegistrationsFixture()
	checkpointIdentity := bson.M{"_id": "verified"}

	t.Run("writes the watches and the checkpoint in one transaction", func(t *testing.T) {
		saveMocks := newCheckpointSaveMocks()
		saveMocks.expectTransaction(t, nil)
		for _, watch := range []rootstock.PegInWatch{first, second} {
			saveMocks.watches.EXPECT().UpdateOne(
				mock.Anything,
				bson.M{"rsk_address": watch.RskAddress},
				bson.M{"$setOnInsert": watch},
				withUpdateUpsert(),
			).Return(&mongoDb.UpdateResult{UpsertedCount: 1}, nil).Once()
		}
		saveMocks.checkpoints.EXPECT().UpdateOne(
			mock.Anything,
			checkpointIdentity,
			bson.M{"$set": checkpoint},
			withUpdateUpsert(),
		).Return(&mongoDb.UpdateResult{UpsertedCount: 1}, nil).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(saveMocks.client, time.Second))
		require.NoError(t, repo.SaveVerifiedRegistrations(context.Background(), []rootstock.PegInWatch{first, second}, checkpoint))
		saveMocks.assertExpectations(t)
	})
	t.Run("writes only the checkpoint when there are no new watches", func(t *testing.T) {
		saveMocks := newCheckpointSaveMocks()
		saveMocks.expectTransaction(t, nil)
		saveMocks.checkpoints.EXPECT().UpdateOne(mock.Anything, checkpointIdentity, bson.M{"$set": checkpoint}, withUpdateUpsert()).
			Return(&mongoDb.UpdateResult{ModifiedCount: 1}, nil).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(saveMocks.client, time.Second))
		require.NoError(t, repo.SaveVerifiedRegistrations(context.Background(), []rootstock.PegInWatch{}, checkpoint))
		saveMocks.assertExpectations(t)
	})
}

func TestPegInWatchMongoRepository_SaveVerifiedRegistrationsErrors(t *testing.T) {
	first, second, checkpoint := verifiedRegistrationsFixture()
	checkpointIdentity := bson.M{"_id": "verified"}

	t.Run("stops before the checkpoint when a watch write fails", func(t *testing.T) {
		saveMocks := newCheckpointSaveMocks()
		saveMocks.expectTransaction(t, assert.AnError)
		saveMocks.watches.EXPECT().UpdateOne(mock.Anything, bson.M{"rsk_address": first.RskAddress}, mock.Anything, withUpdateUpsert()).
			Return(nil, assert.AnError).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(saveMocks.client, time.Second))
		err := repo.SaveVerifiedRegistrations(context.Background(), []rootstock.PegInWatch{first, second}, checkpoint)

		require.ErrorIs(t, err, assert.AnError)
		saveMocks.assertExpectations(t)
	})
	t.Run("returns the checkpoint write error", func(t *testing.T) {
		saveMocks := newCheckpointSaveMocks()
		saveMocks.expectTransaction(t, assert.AnError)
		saveMocks.watches.EXPECT().UpdateOne(mock.Anything, mock.Anything, mock.Anything, withUpdateUpsert()).
			Return(&mongoDb.UpdateResult{UpsertedCount: 1}, nil).Once()
		saveMocks.checkpoints.EXPECT().UpdateOne(mock.Anything, checkpointIdentity, bson.M{"$set": checkpoint}, withUpdateUpsert()).
			Return(nil, assert.AnError).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(saveMocks.client, time.Second))
		err := repo.SaveVerifiedRegistrations(context.Background(), []rootstock.PegInWatch{first}, checkpoint)

		require.ErrorIs(t, err, assert.AnError)
		saveMocks.assertExpectations(t)
	})
	t.Run("returns the session error without writing", func(t *testing.T) {
		saveMocks := newCheckpointSaveMocks()
		saveMocks.client.On("StartSession").Return(saveMocks.session, assert.AnError).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(saveMocks.client, time.Second))
		err := repo.SaveVerifiedRegistrations(context.Background(), []rootstock.PegInWatch{first}, checkpoint)

		require.ErrorIs(t, err, assert.AnError)
		saveMocks.assertExpectations(t)
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
