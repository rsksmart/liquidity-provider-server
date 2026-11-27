package mongo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDb "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

//nolint:funlen // One fixture shared across the round-trip scenarios keeps the event identity identical in each.
func TestPegInWatchMongoRepository(t *testing.T) {
	log.SetLevel(log.DebugLevel)
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
		logs := captureDbLog()
		require.NoError(t, repo.Upsert(context.Background(), entry))
		assert.Contains(t, logs.String(), dbLogLine(mongo.Upsert, entry))
		collection.AssertExpectations(t)
	})

	t.Run("returns nil on a duplicate-key upsert error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		collection.EXPECT().UpdateOne(mock.Anything, identity, mock.Anything, mock.Anything).
			Return(nil, mongoDb.WriteException{WriteErrors: []mongoDb.WriteError{{Code: 11000}}}).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		logs := captureDbLog()
		require.NoError(t, repo.Upsert(context.Background(), entry))
		assert.Empty(t, logs.String())
		collection.AssertExpectations(t)
	})

	t.Run("returns the UpdateOne error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		collection.EXPECT().UpdateOne(mock.Anything, identity, mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		logs := captureDbLog()
		require.ErrorIs(t, repo.Upsert(context.Background(), entry), assert.AnError)
		assert.Empty(t, logs.String())
		collection.AssertExpectations(t)
	})

	t.Run("returns the document matching the identity", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		collection.EXPECT().FindOne(mock.Anything, identity).
			Return(mongoDb.NewSingleResultFromDocument(entry, nil, nil)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		logs := captureDbLog()
		result, err := repo.Get(context.Background(), entry.RskAddress)
		require.NoError(t, err)
		assert.Equal(t, &entry, result)
		assert.Contains(t, logs.String(), dbLogLine(mongo.Read, entry))
		collection.AssertExpectations(t)
	})

	t.Run("returns nil when no document matches the identity", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		collection.EXPECT().FindOne(mock.Anything, identity).
			Return(mongoDb.NewSingleResultFromDocument(rootstock.PegInWatch{}, mongoDb.ErrNoDocuments, nil)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		logs := captureDbLog()
		result, err := repo.Get(context.Background(), entry.RskAddress)
		require.NoError(t, err)
		assert.Nil(t, result)
		assert.Empty(t, logs.String())
		collection.AssertExpectations(t)
	})

	t.Run("returns the FindOne error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		collection.EXPECT().FindOne(mock.Anything, identity).
			Return(mongoDb.NewSingleResultFromDocument(rootstock.PegInWatch{}, assert.AnError, nil)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		logs := captureDbLog()
		result, err := repo.Get(context.Background(), entry.RskAddress)
		require.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, result)
		assert.Empty(t, logs.String())
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
		logs := captureDbLog()
		result, err := repo.List(context.Background())
		require.NoError(t, err)
		assert.Equal(t, []rootstock.PegInWatch{entry, second}, result)
		assert.Contains(t, logs.String(), dbLogLine(mongo.Read, 2))
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
		logs := captureDbLog()
		result, err := repo.List(context.Background())
		require.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, result)
		assert.Empty(t, logs.String())
		collection.AssertExpectations(t)
	})

	t.Run("replaces the document when one row matches", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		imported := entry
		imported.State = rootstock.PegInWatchImported
		collection.EXPECT().ReplaceOne(mock.Anything, identity, imported).
			Return(&mongoDb.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		logs := captureDbLog()
		require.NoError(t, repo.Update(context.Background(), imported))
		assert.Contains(t, logs.String(), dbLogLine(mongo.Update, imported))
		collection.AssertExpectations(t)
	})

	t.Run("returns not-found when ReplaceOne matches nothing", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		imported := entry
		imported.State = rootstock.PegInWatchImported
		collection.EXPECT().ReplaceOne(mock.Anything, identity, imported).
			Return(&mongoDb.UpdateResult{MatchedCount: 0}, nil).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		logs := captureDbLog()
		require.ErrorContains(t, repo.Update(context.Background(), imported), "pegin watch not found")
		assert.Empty(t, logs.String())
		collection.AssertExpectations(t)
	})

	t.Run("returns the ReplaceOne error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		imported := entry
		imported.State = rootstock.PegInWatchImported
		collection.EXPECT().ReplaceOne(mock.Anything, identity, imported).
			Return(nil, assert.AnError).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		logs := captureDbLog()
		require.ErrorIs(t, repo.Update(context.Background(), imported), assert.AnError)
		assert.Empty(t, logs.String())
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
	log.SetLevel(log.DebugLevel)
	checkpointIdentity := bson.M{"_id": "verified"}
	t.Run("returns none when no checkpoint is stored", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCheckpointCollection)
		collection.EXPECT().FindOne(mock.Anything, checkpointIdentity).
			Return(mongoDb.NewSingleResultFromDocument(bson.M{}, mongoDb.ErrNoDocuments, nil)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		logs := captureDbLog()
		checkpoint, found, err := repo.GetCheckpoint(context.Background())

		require.NoError(t, err)
		assert.False(t, found)
		assert.Equal(t, rootstock.PegInWatchCheckpoint{}, checkpoint)
		assert.Empty(t, logs.String())
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
		logs := captureDbLog()
		checkpoint, found, err := repo.GetCheckpoint(context.Background())

		require.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, stored, checkpoint)
		assert.Contains(t, logs.String(), dbLogLine(mongo.Read, stored))
		collection.AssertExpectations(t)
	})
	t.Run("returns the read error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCheckpointCollection)
		collection.EXPECT().FindOne(mock.Anything, checkpointIdentity).
			Return(mongoDb.NewSingleResultFromDocument(bson.M{}, assert.AnError, nil)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		logs := captureDbLog()
		checkpoint, found, err := repo.GetCheckpoint(context.Background())

		require.ErrorIs(t, err, assert.AnError)
		assert.False(t, found)
		assert.Equal(t, rootstock.PegInWatchCheckpoint{}, checkpoint)
		assert.Empty(t, logs.String())
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
	firstSeenAt := time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)
	secondSeenAt := time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC)
	first := rootstock.PegInWatch{
		TxHash:           "0x01",
		LogIndex:         3,
		BlockNumber:      10,
		RskAddress:       "0xaaaa",
		Registrant:       "0xregistrant1",
		RegistrationRoot: [32]byte{1},
		Encoding:         1,
		BtcAddress:       "btc-first",
		State:            rootstock.PegInWatchDiscovered,
		LastSeenAt:       &firstSeenAt,
		LastError:        "first error",
		CreatedAt:        time.Date(2026, 11, 1, 8, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2026, 11, 1, 9, 0, 0, 0, time.UTC),
	}
	second := rootstock.PegInWatch{
		TxHash:           "0x02",
		LogIndex:         4,
		BlockNumber:      11,
		RskAddress:       "0xbbbb",
		Registrant:       "0xregistrant2",
		RegistrationRoot: [32]byte{2},
		Encoding:         2,
		BtcAddress:       "btc-second",
		State:            rootstock.PegInWatchImported,
		LastSeenAt:       &secondSeenAt,
		LastError:        "second error",
		CreatedAt:        time.Date(2026, 11, 2, 8, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2026, 11, 2, 9, 0, 0, 0, time.UTC),
	}
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
		logs := captureDbLog()
		require.NoError(t, repo.SaveVerifiedRegistrations(context.Background(), []rootstock.PegInWatch{first, second}, checkpoint))
		assert.Contains(t, logs.String(), registrationsLogLine([]rootstock.PegInWatch{first, second}, checkpoint))
		saveMocks.assertExpectations(t)
	})
	t.Run("writes only the checkpoint when there are no new watches", func(t *testing.T) {
		saveMocks := newCheckpointSaveMocks()
		saveMocks.expectTransaction(t, nil)
		saveMocks.checkpoints.EXPECT().UpdateOne(mock.Anything, checkpointIdentity, bson.M{"$set": checkpoint}, withUpdateUpsert()).
			Return(&mongoDb.UpdateResult{ModifiedCount: 1}, nil).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(saveMocks.client, time.Second))
		logs := captureDbLog()
		require.NoError(t, repo.SaveVerifiedRegistrations(context.Background(), []rootstock.PegInWatch{}, checkpoint))
		assert.Contains(t, logs.String(), registrationsLogLine([]rootstock.PegInWatch{}, checkpoint))
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
		logs := captureDbLog()
		err := repo.SaveVerifiedRegistrations(context.Background(), []rootstock.PegInWatch{first, second}, checkpoint)

		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, logs.String())
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
		logs := captureDbLog()
		err := repo.SaveVerifiedRegistrations(context.Background(), []rootstock.PegInWatch{first}, checkpoint)

		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, logs.String())
		saveMocks.assertExpectations(t)
	})
	t.Run("returns the session error without writing", func(t *testing.T) {
		saveMocks := newCheckpointSaveMocks()
		saveMocks.client.On("StartSession").Return(saveMocks.session, assert.AnError).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(saveMocks.client, time.Second))
		logs := captureDbLog()
		err := repo.SaveVerifiedRegistrations(context.Background(), []rootstock.PegInWatch{first}, checkpoint)

		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, logs.String())
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

func TestPegInWatchMongoRepository_GetRootAt(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	atOrBelow := bson.M{"block_number": bson.M{"$lte": uint64(150)}}
	newestFirst := mock.MatchedBy(func(opt options.Lister[options.FindOneOptions]) bool {
		resolved := &options.FindOneOptions{}
		for _, fn := range opt.List() {
			if err := fn(resolved); err != nil {
				return false
			}
		}
		return assert.ObjectsAreEqual(bson.D{{Key: "block_number", Value: -1}, {Key: "log_index", Value: -1}}, resolved.Sort)
	})
	t.Run("returns the root of the last row at or below the block", func(t *testing.T) {
		lastRow := rootstock.PegInWatch{RskAddress: "0xaaaa", BlockNumber: 149, RegistrationRoot: [32]byte{4, 5}}
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		collection.EXPECT().FindOne(mock.Anything, atOrBelow, newestFirst).
			Return(mongoDb.NewSingleResultFromDocument(lastRow, nil, nil)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		logs := captureDbLog()
		root, err := repo.GetRootAt(context.Background(), 150)

		require.NoError(t, err)
		assert.Equal(t, [32]byte{4, 5}, root)
		assert.Contains(t, logs.String(), dbLogLine(mongo.Read, lastRow))
		collection.AssertExpectations(t)
	})
	t.Run("returns the zero root when no row is at or below the block", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		collection.EXPECT().FindOne(mock.Anything, atOrBelow, newestFirst).
			Return(mongoDb.NewSingleResultFromDocument(bson.M{}, mongoDb.ErrNoDocuments, nil)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		logs := captureDbLog()
		root, err := repo.GetRootAt(context.Background(), 150)

		require.NoError(t, err)
		assert.Equal(t, [32]byte{}, root)
		assert.Empty(t, logs.String())
		collection.AssertExpectations(t)
	})
	t.Run("returns the read error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInWatchCollection)
		collection.EXPECT().FindOne(mock.Anything, atOrBelow, newestFirst).
			Return(mongoDb.NewSingleResultFromDocument(bson.M{}, assert.AnError, nil)).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(client, time.Second))
		logs := captureDbLog()
		_, err := repo.GetRootAt(context.Background(), 150)

		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, logs.String())
		collection.AssertExpectations(t)
	})
}

func replacedModels(t *testing.T, fromBlock uint64, watches ...rootstock.PegInWatch) interface{} {
	t.Helper()
	return mock.MatchedBy(func(models []mongoDb.WriteModel) bool {
		if len(models) != len(watches)+1 {
			return false
		}
		addresses := make([]string, 0, len(watches))
		for index, watch := range watches {
			addresses = append(addresses, watch.RskAddress)
			if !isRebuiltWatchUpsert(models[index+1], watch) {
				return false
			}
		}
		deleteModel, ok := models[0].(*mongoDb.DeleteManyModel)
		return ok && assert.ObjectsAreEqual(bson.M{
			"block_number": bson.M{"$gte": fromBlock},
			"rsk_address":  bson.M{"$nin": addresses},
		}, deleteModel.Filter)
	})
}

func isRebuiltWatchUpsert(model mongoDb.WriteModel, watch rootstock.PegInWatch) bool {
	upsert, ok := model.(*mongoDb.UpdateOneModel)
	if !ok || upsert.Upsert == nil || !*upsert.Upsert || !assert.ObjectsAreEqual(bson.M{"rsk_address": watch.RskAddress}, upsert.Filter) {
		return false
	}
	update, ok := upsert.Update.(bson.M)
	return ok && keepsImportState(update, watch)
}

func keepsImportState(update bson.M, watch rootstock.PegInWatch) bool {
	return assert.ObjectsAreEqual(bson.M{
		"$set": bson.M{
			"tx_hash":           watch.TxHash,
			"log_index":         watch.LogIndex,
			"block_number":      watch.BlockNumber,
			"registrant":        watch.Registrant,
			"registration_root": watch.RegistrationRoot,
			"updated_at":        watch.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"encoding":     watch.Encoding,
			"btc_address":  watch.BtcAddress,
			"state":        watch.State,
			"last_seen_at": watch.LastSeenAt,
			"last_error":   watch.LastError,
			"created_at":   watch.CreatedAt,
		},
	}, update)
}

func orderedBulkWrite() interface{} {
	return mock.MatchedBy(func(opt options.Lister[options.BulkWriteOptions]) bool {
		resolved := &options.BulkWriteOptions{}
		for _, fn := range opt.List() {
			if err := fn(resolved); err != nil {
				return false
			}
		}
		return resolved.Ordered != nil && *resolved.Ordered
	})
}

func TestPegInWatchMongoRepository_ReplaceVerifiedRegistrations(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	first, second, checkpoint := verifiedRegistrationsFixture()
	checkpointIdentity := bson.M{"_id": "verified"}

	t.Run("removes, upserts and saves the checkpoint in one transaction with one bulk write", func(t *testing.T) {
		saveMocks := newCheckpointSaveMocks()
		saveMocks.expectTransaction(t, nil)
		saveMocks.watches.EXPECT().BulkWrite(mock.Anything, replacedModels(t, 10, first, second), orderedBulkWrite()).
			Return(&mongoDb.BulkWriteResult{}, nil).Once()
		saveMocks.checkpoints.EXPECT().UpdateOne(mock.Anything, checkpointIdentity, bson.M{"$set": checkpoint}, withUpdateUpsert()).
			Return(&mongoDb.UpdateResult{ModifiedCount: 1}, nil).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(saveMocks.client, time.Second))
		logs := captureDbLog()
		require.NoError(t, repo.ReplaceVerifiedRegistrations(context.Background(), 10, []rootstock.PegInWatch{first, second}, checkpoint))
		assert.Contains(t, logs.String(), dbLogLine(mongo.Delete, "removed 0 registrations at or after block 10"))
		assert.Contains(t, logs.String(), registrationsLogLine([]rootstock.PegInWatch{first, second}, checkpoint))
		saveMocks.assertExpectations(t)
	})
	t.Run("removes every row from the block when the chain has none there", func(t *testing.T) {
		saveMocks := newCheckpointSaveMocks()
		saveMocks.expectTransaction(t, nil)
		saveMocks.watches.EXPECT().BulkWrite(mock.Anything, replacedModels(t, 10), orderedBulkWrite()).
			Return(&mongoDb.BulkWriteResult{DeletedCount: 2}, nil).Once()
		saveMocks.checkpoints.EXPECT().UpdateOne(mock.Anything, checkpointIdentity, bson.M{"$set": checkpoint}, withUpdateUpsert()).
			Return(&mongoDb.UpdateResult{ModifiedCount: 1}, nil).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(saveMocks.client, time.Second))
		logs := captureDbLog()
		require.NoError(t, repo.ReplaceVerifiedRegistrations(context.Background(), 10, []rootstock.PegInWatch{}, checkpoint))
		assert.Contains(t, logs.String(), dbLogLine(mongo.Delete, "removed 2 registrations at or after block 10"))
		assert.Contains(t, logs.String(), registrationsLogLine([]rootstock.PegInWatch{}, checkpoint))
		saveMocks.assertExpectations(t)
	})
}

func TestPegInWatchMongoRepository_ReplaceVerifiedRegistrationsErrors(t *testing.T) {
	first, _, checkpoint := verifiedRegistrationsFixture()
	checkpointIdentity := bson.M{"_id": "verified"}

	t.Run("stops before the checkpoint when the bulk write fails", func(t *testing.T) {
		saveMocks := newCheckpointSaveMocks()
		saveMocks.expectTransaction(t, assert.AnError)
		saveMocks.watches.EXPECT().BulkWrite(mock.Anything, mock.Anything, orderedBulkWrite()).Return(nil, assert.AnError).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(saveMocks.client, time.Second))
		logs := captureDbLog()
		err := repo.ReplaceVerifiedRegistrations(context.Background(), 10, []rootstock.PegInWatch{first}, checkpoint)

		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, logs.String())
		saveMocks.assertExpectations(t)
	})
	t.Run("returns the checkpoint write error", func(t *testing.T) {
		saveMocks := newCheckpointSaveMocks()
		saveMocks.expectTransaction(t, assert.AnError)
		saveMocks.watches.EXPECT().BulkWrite(mock.Anything, mock.Anything, orderedBulkWrite()).Return(&mongoDb.BulkWriteResult{}, nil).Once()
		saveMocks.checkpoints.EXPECT().UpdateOne(mock.Anything, checkpointIdentity, bson.M{"$set": checkpoint}, withUpdateUpsert()).
			Return(nil, assert.AnError).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(saveMocks.client, time.Second))
		logs := captureDbLog()
		err := repo.ReplaceVerifiedRegistrations(context.Background(), 10, []rootstock.PegInWatch{first}, checkpoint)

		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, logs.String())
		saveMocks.assertExpectations(t)
	})
	t.Run("returns the session error without writing", func(t *testing.T) {
		saveMocks := newCheckpointSaveMocks()
		saveMocks.client.On("StartSession").Return(saveMocks.session, assert.AnError).Once()

		repo := mongo.NewPegInWatchMongoRepository(mongo.NewConnection(saveMocks.client, time.Second))
		logs := captureDbLog()
		err := repo.ReplaceVerifiedRegistrations(context.Background(), 10, []rootstock.PegInWatch{first}, checkpoint)

		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, logs.String())
		saveMocks.assertExpectations(t)
	})
}

func captureDbLog() *test.ThreadSafeBuffer {
	buffer := new(test.ThreadSafeBuffer)
	log.SetOutput(buffer)
	return buffer
}

func dbLogLine(interaction mongo.DbInteraction, value any) string {
	return fmt.Sprintf("%s interaction with db: %+v", interaction, value)
}

func registrationsLogLine(watches []rootstock.PegInWatch, checkpoint rootstock.PegInWatchCheckpoint) string {
	return fmt.Sprintf("%s interaction with db: {Watches:%+v Checkpoint:%+v}", mongo.Upsert, watches, checkpoint)
}
