package mongo_test

import (
	"context"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"testing"
	"time"
)

var testPegoutQuote = quote.PegoutQuote{
	LbcAddress:            test.AnyAddress,
	LpRskAddress:          test.AnyAddress,
	BtcRefundAddress:      test.AnyAddress,
	RskRefundAddress:      test.AnyAddress,
	LpBtcAddress:          test.AnyAddress,
	CallFee:               entities.NewWei(1),
	PenaltyFee:            2,
	Nonce:                 3,
	DepositAddress:        test.AnyAddress,
	Value:                 entities.NewWei(4),
	AgreementTimestamp:    5,
	DepositDateLimit:      6,
	DepositConfirmations:  7,
	TransferConfirmations: 8,
	TransferTime:          9,
	ExpireDate:            10,
	ExpireBlock:           11,
	GasFee:                entities.NewWei(12),
	ProductFeeAmount:      13,
}

var testRetainedPegoutQuote = quote.RetainedPegoutQuote{
	QuoteHash:          test.AnyString,
	DepositAddress:     test.AnyAddress,
	Signature:          test.AnyString,
	RequiredLiquidity:  entities.NewWei(55),
	State:              quote.PegoutStateWaitingForDepositConfirmations,
	UserRskTxHash:      test.AnyString,
	LpBtcTxHash:        test.AnyString,
	RefundPegoutTxHash: test.AnyString,
	BridgeRefundTxHash: test.AnyString,
}

var testPegoutDeposit = quote.PegoutDeposit{
	TxHash:      test.AnyString,
	QuoteHash:   test.AnyString,
	Amount:      entities.NewWei(999),
	Timestamp:   time.Unix(1715001146288, 0).UTC(),
	BlockNumber: 789,
	From:        test.AnyAddress,
}

func TestPegoutMongoRepository_InsertQuote(t *testing.T) {
	t.Run("Insert pegout quote successfully", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegoutQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.MatchedBy(func(q mongo.StoredPegoutQuote) bool {
			return q.Hash == test.AnyString && reflect.TypeOf(quote.PegoutQuote{}).NumField() == test.CountNonZeroValues(q.PegoutQuote)
		})).Return(nil, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		defer assertDbInteractionLog(t, mongo.Insert)()
		err := repo.InsertQuote(context.Background(), test.AnyString, testPegoutQuote)
		collection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Db error when inserting pegout quote", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegoutQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		err := repo.InsertQuote(context.Background(), test.AnyString, testPegoutQuote)
		collection.AssertExpectations(t)
		require.Error(t, err)
	})
}

func TestPegoutMongoRepository_GetQuote(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.PegoutQuoteCollection)
	log.SetLevel(log.DebugLevel)
	t.Run("Get pegout quote successfully", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, bson.D{primitive.E{Key: "hash", Value: test.AnyString}}).
			Return(mongoDb.NewSingleResultFromDocument(mongo.StoredPegoutQuote{
				PegoutQuote: testPegoutQuote,
				Hash:        test.AnyString,
			}, nil, nil)).Once()
		defer assertDbInteractionLog(t, mongo.Read)()
		result, err := repo.GetQuote(context.Background(), test.AnyString)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, testPegoutQuote, *result)
	})
	t.Run("Db error when getting pegout quote", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(mongo.StoredPegoutQuote{}, assert.AnError, nil)).Once()
		result, err := repo.GetQuote(context.Background(), test.AnyString)
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Pegout quote not present in db", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(mongo.StoredPegoutQuote{}, mongoDb.ErrNoDocuments, nil)).Once()
		result, err := repo.GetQuote(context.Background(), test.AnyString)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestPegoutMongoRepository_GetRetainedQuote(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
	log.SetLevel(log.DebugLevel)
	t.Run("Get retained pegout quote successfully", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, bson.D{primitive.E{Key: "quote_hash", Value: test.AnyString}}).
			Return(mongoDb.NewSingleResultFromDocument(testRetainedPegoutQuote, nil, nil)).Once()
		defer assertDbInteractionLog(t, mongo.Read)()
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyString)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, testRetainedPegoutQuote, *result)
	})
	t.Run("Db error when getting retained pegout quote", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(quote.RetainedPegoutQuote{}, assert.AnError, nil)).Once()
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyString)
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Retained pegout quote not present in db", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(quote.RetainedPegoutQuote{}, mongoDb.ErrNoDocuments, nil)).Once()
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyString)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestPegoutMongoRepository_InsertRetainedQuote(t *testing.T) {
	t.Run("Insert retained pegout quote successfully", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.MatchedBy(func(q quote.RetainedPegoutQuote) bool {
			return q.QuoteHash == test.AnyString && reflect.TypeOf(quote.RetainedPegoutQuote{}).NumField() == test.CountNonZeroValues(q)
		})).Return(nil, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		defer assertDbInteractionLog(t, mongo.Insert)()
		err := repo.InsertRetainedQuote(context.Background(), testRetainedPegoutQuote)
		collection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Db error when inserting retained pegout quote", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		err := repo.InsertRetainedQuote(context.Background(), testRetainedPegoutQuote)
		collection.AssertExpectations(t)
		require.Error(t, err)
	})
}

func TestPegoutMongoRepository_UpdateRetainedQuote(t *testing.T) {
	const updated = "updated value"
	t.Run("Update retained pegout quote successfully", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		updatedQuote := testRetainedPegoutQuote
		updatedQuote.State = quote.PegoutStateSendPegoutFailed
		updatedQuote.Signature = updated
		updatedQuote.RefundPegoutTxHash = updated
		updatedQuote.LpBtcTxHash = updated
		updatedQuote.DepositAddress = updated
		updatedQuote.UserRskTxHash = updated
		updatedQuote.RequiredLiquidity = entities.NewWei(200)
		collection.On("UpdateOne", mock.Anything,
			bson.D{primitive.E{Key: "quote_hash", Value: testRetainedPegoutQuote.QuoteHash}},
			bson.D{primitive.E{Key: "$set", Value: updatedQuote}},
		).Return(&mongoDb.UpdateResult{ModifiedCount: 1}, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		defer assertDbInteractionLog(t, mongo.Update)()
		err := repo.UpdateRetainedQuote(context.Background(), updatedQuote)
		collection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Db error when updating retained pegout quote", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		err := repo.UpdateRetainedQuote(context.Background(), testRetainedPegoutQuote)
		collection.AssertExpectations(t)
		require.Error(t, err)
	})
	t.Run("Retained pegout quote to update not found", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).
			Return(&mongoDb.UpdateResult{ModifiedCount: 0}, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		err := repo.UpdateRetainedQuote(context.Background(), testRetainedPegoutQuote)
		collection.AssertExpectations(t)
		require.ErrorIs(t, err, usecases.QuoteNotFoundError)
	})
	t.Run("Update more than one retained pegout quote", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).
			Return(&mongoDb.UpdateResult{ModifiedCount: 2}, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		err := repo.UpdateRetainedQuote(context.Background(), testRetainedPegoutQuote)
		collection.AssertExpectations(t)
		require.ErrorContains(t, err, "multiple documents updated")
	})
}

func TestPegoutMongoRepository_DeleteQuotes(t *testing.T) {
	var hashes = []string{"pegout1", "pegout2", "pegout3"}
	log.SetLevel(log.DebugLevel)
	t.Run("Delete quotes successfully", func(t *testing.T) {
		client, quoteCollection := getClientAndCollectionMocks(mongo.PegoutQuoteCollection)
		retainedCollection := &mocks.CollectionBindingMock{}
		client.Database(mongo.DbName).(*mocks.DbBindingMock).On("Collection", mongo.RetainedPegoutQuoteCollection).Return(retainedCollection)
		quoteCollection.On("DeleteMany", mock.Anything,
			bson.D{primitive.E{Key: "quote_hash", Value: bson.D{primitive.E{Key: "$in", Value: hashes}}}},
		).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything,
			bson.D{primitive.E{Key: "quote_hash", Value: bson.D{primitive.E{Key: "$in", Value: hashes}}}},
		).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		defer assertDbInteractionLog(t, mongo.Delete)()
		count, err := repo.DeleteQuotes(context.Background(), hashes)
		quoteCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, uint(6), count)
	})
	t.Run("Db error when deleting pegout quotes", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegoutQuoteCollection)
		collection.On("DeleteMany", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		count, err := repo.DeleteQuotes(context.Background(), []string{test.AnyString})
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Zero(t, count)
	})
	t.Run("Db error when deleting retained pegout quotes", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegoutQuoteCollection)
		retainedCollection := &mocks.CollectionBindingMock{}
		client.Database(mongo.DbName).(*mocks.DbBindingMock).On("Collection", mongo.RetainedPegoutQuoteCollection).Return(retainedCollection)
		collection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		count, err := repo.DeleteQuotes(context.Background(), []string{test.AnyString})
		collection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
		require.Error(t, err)
		assert.Zero(t, count)
	})
	t.Run("Error when deletion count missmatch", func(t *testing.T) {
		client, quoteCollection := getClientAndCollectionMocks(mongo.PegoutQuoteCollection)
		retainedCollection := &mocks.CollectionBindingMock{}
		client.Database(mongo.DbName).(*mocks.DbBindingMock).On("Collection", mongo.RetainedPegoutQuoteCollection).Return(retainedCollection)
		quoteCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 4}, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		count, err := repo.DeleteQuotes(context.Background(), hashes)
		quoteCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
		require.ErrorContains(t, err, "pegout quote collections didn't match")
		assert.Zero(t, count)
	})
}

func TestPegoutMongoRepository_GetRetainedQuoteByState(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
	log.SetLevel(log.DebugLevel)
	states := []quote.PegoutState{quote.PegoutStateSendPegoutSucceeded, quote.PegoutStateSendPegoutFailed}
	t.Run("Get retained pegout quotes by state successfully", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		secondQuote := testRetainedPegoutQuote
		secondQuote.QuoteHash = "other hash"
		secondQuote.Signature = "456"
		secondQuote.RequiredLiquidity = entities.NewWei(777)
		collection.On("Find", mock.Anything,
			bson.D{primitive.E{Key: "state", Value: bson.D{primitive.E{Key: "$in", Value: states}}}},
		).Return(mongoDb.NewCursorFromDocuments([]any{testRetainedPegoutQuote, secondQuote}, nil, nil)).Once()
		defer assertDbInteractionLog(t, mongo.Read)()
		result, err := repo.GetRetainedQuoteByState(context.Background(), states...)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, []quote.RetainedPegoutQuote{testRetainedPegoutQuote, secondQuote}, result)
	})
	t.Run("Db error when getting retained pegout quotes by state", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("Find", mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()
		result, err := repo.GetRetainedQuoteByState(context.Background(), states...)
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestPegoutMongoRepository_ListPegoutDepositsByAddress(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.DepositEventsCollection)
	log.SetLevel(log.DebugLevel)
	t.Run("List pegout deposits by address successfully", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("Find", mock.Anything,
			bson.M{"from": bson.M{"$regex": test.AnyAddress, "$options": "i"}},
			options.Find().SetSort(bson.M{"timestamp": -1}),
		).Return(mongoDb.NewCursorFromDocuments([]any{testPegoutDeposit}, nil, nil)).Once()
		defer assertDbInteractionLog(t, mongo.Read)()
		result, err := repo.ListPegoutDepositsByAddress(context.Background(), test.AnyAddress)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, []quote.PegoutDeposit{testPegoutDeposit}[0], result[0])
	})
	t.Run("Db error when listing pegout deposits by address", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("Find", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()
		result, err := repo.ListPegoutDepositsByAddress(context.Background(), test.AnyAddress)
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestPegoutMongoRepository_UpsertPegoutDeposit(t *testing.T) {
	t.Run("Upsert pegout deposit successfully", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.DepositEventsCollection)
		newDeposit := testPegoutDeposit
		newDeposit.Amount = entities.NewWei(1000)
		newDeposit.Timestamp = time.Now().UTC()
		newDeposit.BlockNumber = 790
		collection.On("ReplaceOne", mock.Anything,
			bson.M{"tx_hash": testPegoutDeposit.TxHash},
			newDeposit,
			options.Replace().SetUpsert(true),
		).Return(&mongoDb.UpdateResult{ModifiedCount: 1}, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		defer assertDbInteractionLog(t, mongo.Upsert)()
		err := repo.UpsertPegoutDeposit(context.Background(), newDeposit)
		collection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Db error when upserting pegout deposit", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.DepositEventsCollection)
		collection.On("ReplaceOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		err := repo.UpsertPegoutDeposit(context.Background(), testPegoutDeposit)
		collection.AssertExpectations(t)
		require.Error(t, err)
	})
	t.Run("Error when upserting more than one pegout deposit", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.DepositEventsCollection)
		collection.On("ReplaceOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(&mongoDb.UpdateResult{ModifiedCount: 2}, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		err := repo.UpsertPegoutDeposit(context.Background(), testPegoutDeposit)
		collection.AssertExpectations(t)
		require.ErrorContains(t, err, "multiple deposits updated")
	})
}

func TestPegoutMongoRepository_UpsertPegoutDeposits(t *testing.T) {
	t.Run("Upsert pegout deposits successfully", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.DepositEventsCollection)
		deposits := []quote.PegoutDeposit{
			{TxHash: "tx1", QuoteHash: "quote1", Amount: entities.NewWei(1000), Timestamp: time.Now().UTC(), BlockNumber: 790, From: test.AnyAddress},
			{TxHash: "tx2", QuoteHash: "quote2", Amount: entities.NewWei(2000), Timestamp: time.Now().UTC(), BlockNumber: 791, From: test.AnyAddress},
		}
		collection.On("BulkWrite", mock.Anything,
			[]mongoDb.WriteModel{
				&mongoDb.ReplaceOneModel{
					Upsert:      btcjson.Bool(true),
					Filter:      bson.M{"tx_hash": deposits[0].TxHash},
					Replacement: deposits[0],
				},
				&mongoDb.ReplaceOneModel{
					Upsert:      btcjson.Bool(true),
					Filter:      bson.M{"tx_hash": deposits[1].TxHash},
					Replacement: deposits[1],
				},
			}).
			Return(nil, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		defer assertDbInteractionLog(t, mongo.Upsert)()
		err := repo.UpsertPegoutDeposits(context.Background(), deposits)
		collection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Db error when upserting pegout deposits", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.DepositEventsCollection)
		collection.On("BulkWrite", mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		err := repo.UpsertPegoutDeposits(context.Background(), []quote.PegoutDeposit{testPegoutDeposit})
		collection.AssertExpectations(t)
		require.Error(t, err)
	})
}

// nolint:funlen
func TestPegoutMongoRepository_UpdateRetainedQuotes(t *testing.T) {
	retainedQuotes := []quote.RetainedPegoutQuote{
		{QuoteHash: "quote1", DepositAddress: test.AnyAddress, Signature: test.AnyString, RequiredLiquidity: entities.NewWei(1000), State: quote.PegoutStateSendPegoutSucceeded},
		{QuoteHash: "quote2", DepositAddress: test.AnyAddress, Signature: test.AnyString, RequiredLiquidity: entities.NewWei(2000), State: quote.PegoutStateSendPegoutFailed},
	}
	t.Run("Update retained quotes successfully", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		session := &mocks.SessionBindingMock{}
		client.On("StartSession").Return(session, nil).Once()
		session.On("EndSession", mock.Anything).Return().Once()
		session.On("WithTransaction", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				fn := args.Get(1).(func(mongoDb.SessionContext) (any, error))
				count, err := fn(mongoDb.NewSessionContext(context.Background(), mongoDb.SessionFromContext(context.Background())))
				require.NoError(t, err)
				assert.Equal(t, int64(len(retainedQuotes)), count)
			}).
			Return(any(int64(len(retainedQuotes))), nil)
		for _, q := range retainedQuotes {
			collection.On("UpdateOne", mock.Anything,
				bson.D{primitive.E{Key: "quote_hash", Value: q.QuoteHash}},
				bson.D{primitive.E{Key: "$set", Value: q}},
			).Return(&mongoDb.UpdateResult{ModifiedCount: 1}, nil).Once()
		}
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		defer assertDbInteractionLog(t, mongo.Update)()
		err := repo.UpdateRetainedQuotes(context.Background(), retainedQuotes)
		collection.AssertExpectations(t)
		client.AssertExpectations(t)
		session.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Error creating session", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		session := &mocks.SessionBindingMock{}
		client.On("StartSession").Return(session, assert.AnError).Once()
		session.On("EndSession", mock.Anything).Return().Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		defer test.AssertNoLog(t)()
		err := repo.UpdateRetainedQuotes(context.Background(), retainedQuotes)
		collection.AssertExpectations(t)
		client.AssertExpectations(t)
		session.AssertExpectations(t)
		require.Error(t, err)
	})
	t.Run("Error updating one quote", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		session := &mocks.SessionBindingMock{}
		client.On("StartSession").Return(session, nil).Once()
		session.On("EndSession", mock.Anything).Return().Once()
		session.On("WithTransaction", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				fn := args.Get(1).(func(mongoDb.SessionContext) (any, error))
				count, err := fn(mongoDb.NewSessionContext(context.Background(), mongoDb.SessionFromContext(context.Background())))
				require.Error(t, err)
				assert.Equal(t, int64(0), count)
			}).
			Return(int64(0), assert.AnError)

		collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongoDb.UpdateResult{ModifiedCount: 1}, nil).Once()
		collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		defer test.AssertNoLog(t)()
		err := repo.UpdateRetainedQuotes(context.Background(), retainedQuotes)
		collection.AssertExpectations(t)
		client.AssertExpectations(t)
		session.AssertExpectations(t)
		require.Error(t, err)
	})
	t.Run("Update count mismatch", func(t *testing.T) {
		client, _ := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		session := &mocks.SessionBindingMock{}
		client.On("StartSession").Return(session, nil).Once()
		session.On("EndSession", mock.Anything).Return().Once()
		session.On("WithTransaction", mock.Anything, mock.Anything).Return(any(int64(len(retainedQuotes)-1)), nil)
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		defer test.AssertNoLog(t)()
		err := repo.UpdateRetainedQuotes(context.Background(), retainedQuotes)
		client.AssertExpectations(t)
		session.AssertExpectations(t)
		require.ErrorContains(t, err, "mismatch on updated documents. Expected 2, updated 1")
	})
}
