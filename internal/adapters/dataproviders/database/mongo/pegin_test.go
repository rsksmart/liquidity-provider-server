package mongo_test

import (
	"context"
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
	"reflect"
	"testing"
)

var testPeginQuote = quote.PeginQuote{
	FedBtcAddress:      test.AnyAddress,
	LbcAddress:         test.AnyAddress,
	LpRskAddress:       test.AnyAddress,
	BtcRefundAddress:   test.AnyAddress,
	RskRefundAddress:   test.AnyAddress,
	LpBtcAddress:       test.AnyAddress,
	CallFee:            entities.NewWei(1),
	PenaltyFee:         entities.NewWei(2),
	ContractAddress:    test.AnyAddress,
	Data:               test.AnyString,
	GasLimit:           1,
	Nonce:              2,
	Value:              entities.NewWei(3),
	AgreementTimestamp: 4,
	TimeForDeposit:     5,
	LpCallTime:         6,
	Confirmations:      7,
	CallOnRegister:     true,
	GasFee:             entities.NewWei(4),
	ProductFeeAmount:   8,
}

var testRetainedPeginQuote = quote.RetainedPeginQuote{
	QuoteHash:           test.AnyString,
	DepositAddress:      test.AnyAddress,
	Signature:           test.AnyString,
	RequiredLiquidity:   entities.NewWei(100),
	State:               quote.PeginStateCallForUserSucceeded,
	UserBtcTxHash:       test.AnyString,
	CallForUserTxHash:   test.AnyString,
	RegisterPeginTxHash: test.AnyString,
}

func TestPeginMongoRepository_InsertQuote(t *testing.T) {
	t.Run("Insert pegin quote successfully", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.MatchedBy(func(q mongo.StoredPeginQuote) bool {
			return q.Hash == test.AnyString && reflect.TypeOf(quote.PeginQuote{}).NumField() == test.CountNonZeroValues(q.PeginQuote)
		})).Return(nil, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPeginMongoRepository(conn)
		defer assertDbInteractionLog(t, mongo.Insert)()
		err := repo.InsertQuote(context.Background(), test.AnyString, testPeginQuote)
		collection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Db error when inserting pegin quote", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPeginMongoRepository(conn)
		err := repo.InsertQuote(context.Background(), test.AnyString, testPeginQuote)
		collection.AssertExpectations(t)
		require.Error(t, err)
	})
}

func TestPeginMongoRepository_GetQuote(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
	log.SetLevel(log.DebugLevel)
	t.Run("Get pegin quote successfully", func(t *testing.T) {
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, bson.D{primitive.E{Key: "hash", Value: test.AnyString}}).
			Return(mongoDb.NewSingleResultFromDocument(mongo.StoredPeginQuote{
				PeginQuote: testPeginQuote,
				Hash:       test.AnyString,
			}, nil, nil)).Once()
		defer assertDbInteractionLog(t, mongo.Read)()
		result, err := repo.GetQuote(context.Background(), test.AnyString)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, testPeginQuote, *result)
	})
	t.Run("Db error when getting pegin quote", func(t *testing.T) {
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(mongo.StoredPeginQuote{}, assert.AnError, nil)).Once()
		result, err := repo.GetQuote(context.Background(), test.AnyString)
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Pegin quote not present in db", func(t *testing.T) {
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(mongo.StoredPeginQuote{}, mongoDb.ErrNoDocuments, nil)).Once()
		result, err := repo.GetQuote(context.Background(), test.AnyString)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestPeginMongoRepository_GetRetainedQuote(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
	log.SetLevel(log.DebugLevel)
	t.Run("Get retained pegin quote successfully", func(t *testing.T) {
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, bson.D{primitive.E{Key: "quote_hash", Value: test.AnyString}}).
			Return(mongoDb.NewSingleResultFromDocument(testRetainedPeginQuote, nil, nil)).Once()
		defer assertDbInteractionLog(t, mongo.Read)()
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyString)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, testRetainedPeginQuote, *result)
	})
	t.Run("Db error when getting retained pegin quote", func(t *testing.T) {
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(quote.RetainedPeginQuote{}, assert.AnError, nil)).Once()
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyString)
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Retained pegin quote not present in db", func(t *testing.T) {
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(quote.RetainedPeginQuote{}, mongoDb.ErrNoDocuments, nil)).Once()
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyString)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestPeginMongoRepository_InsertRetainedQuote(t *testing.T) {
	t.Run("Insert retained pegin quote successfully", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.MatchedBy(func(q quote.RetainedPeginQuote) bool {
			return q.QuoteHash == test.AnyString && reflect.TypeOf(quote.RetainedPeginQuote{}).NumField() == test.CountNonZeroValues(q)
		})).Return(nil, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPeginMongoRepository(conn)
		defer assertDbInteractionLog(t, mongo.Insert)()
		err := repo.InsertRetainedQuote(context.Background(), testRetainedPeginQuote)
		collection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Db error when inserting retained pegin quote", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPeginMongoRepository(conn)
		err := repo.InsertRetainedQuote(context.Background(), testRetainedPeginQuote)
		collection.AssertExpectations(t)
		require.Error(t, err)
	})
}

func TestPeginMongoRepository_UpdateRetainedQuote(t *testing.T) {
	const updated = "updated value"
	t.Run("Update retained pegin quote successfully", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		updatedQuote := testRetainedPeginQuote
		updatedQuote.State = quote.PeginStateCallForUserFailed
		updatedQuote.Signature = updated
		updatedQuote.RegisterPeginTxHash = updated
		updatedQuote.CallForUserTxHash = updated
		updatedQuote.DepositAddress = updated
		updatedQuote.UserBtcTxHash = updated
		updatedQuote.RequiredLiquidity = entities.NewWei(200)
		collection.On("UpdateOne", mock.Anything,
			bson.D{primitive.E{Key: "quote_hash", Value: testRetainedPeginQuote.QuoteHash}},
			bson.D{primitive.E{Key: "$set", Value: updatedQuote}},
		).Return(&mongoDb.UpdateResult{ModifiedCount: 1}, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPeginMongoRepository(conn)
		defer assertDbInteractionLog(t, mongo.Update)()
		err := repo.UpdateRetainedQuote(context.Background(), updatedQuote)
		collection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Db error when updating retained pegin quote", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPeginMongoRepository(conn)
		err := repo.UpdateRetainedQuote(context.Background(), testRetainedPeginQuote)
		collection.AssertExpectations(t)
		require.Error(t, err)
	})
	t.Run("Retained pegin quote to update not found", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).
			Return(&mongoDb.UpdateResult{ModifiedCount: 0}, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPeginMongoRepository(conn)
		err := repo.UpdateRetainedQuote(context.Background(), testRetainedPeginQuote)
		collection.AssertExpectations(t)
		require.ErrorIs(t, err, usecases.QuoteNotFoundError)
	})
	t.Run("Update more than one retained pegin quote", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).
			Return(&mongoDb.UpdateResult{ModifiedCount: 2}, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPeginMongoRepository(conn)
		err := repo.UpdateRetainedQuote(context.Background(), testRetainedPeginQuote)
		collection.AssertExpectations(t)
		require.ErrorContains(t, err, "multiple documents updated")
	})
}

func TestPeginMongoRepository_GetRetainedQuoteByState(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
	log.SetLevel(log.DebugLevel)
	states := []quote.PeginState{quote.PeginStateCallForUserSucceeded, quote.PeginStateCallForUserFailed}
	t.Run("Get retained pegin quotes by state successfully", func(t *testing.T) {
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client))
		secondQuote := testRetainedPeginQuote
		secondQuote.QuoteHash = "second"
		secondQuote.Signature = "123"
		secondQuote.RequiredLiquidity = entities.NewWei(777)
		collection.On("Find", mock.Anything,
			bson.D{primitive.E{Key: "state", Value: bson.D{primitive.E{Key: "$in", Value: states}}}},
		).Return(mongoDb.NewCursorFromDocuments([]any{testRetainedPeginQuote, secondQuote}, nil, nil)).Once()
		defer assertDbInteractionLog(t, mongo.Read)()
		result, err := repo.GetRetainedQuoteByState(context.Background(), states...)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, []quote.RetainedPeginQuote{testRetainedPeginQuote, secondQuote}, result)
	})
	t.Run("Db error when getting retained pegin quotes by state", func(t *testing.T) {
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client))
		collection.On("Find", mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()
		result, err := repo.GetRetainedQuoteByState(context.Background(), states...)
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestPeginMongoRepository_DeleteQuotes(t *testing.T) {
	var hashes = []string{"pegin1", "pegin2", "pegin3"}
	log.SetLevel(log.DebugLevel)
	t.Run("Delete quotes successfully", func(t *testing.T) {
		client, quoteCollection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
		retainedCollection := &mocks.CollectionBindingMock{}
		client.Database(mongo.DbName).(*mocks.DbBindingMock).On("Collection", mongo.RetainedPeginQuoteCollection).
			Return(retainedCollection)
		quoteCollection.On("DeleteMany", mock.Anything,
			bson.D{primitive.E{Key: "hash", Value: bson.D{primitive.E{Key: "$in", Value: hashes}}}},
		).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything,
			bson.D{primitive.E{Key: "quote_hash", Value: bson.D{primitive.E{Key: "$in", Value: hashes}}}},
		).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPeginMongoRepository(conn)
		defer assertDbInteractionLog(t, mongo.Delete)()
		count, err := repo.DeleteQuotes(context.Background(), hashes)
		quoteCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, uint(6), count)
	})
	t.Run("Db error when deleting pegin quotes", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
		collection.On("DeleteMany", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPeginMongoRepository(conn)
		count, err := repo.DeleteQuotes(context.Background(), []string{test.AnyString})
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Zero(t, count)
	})
	t.Run("Db error when deleting retained pegin quotes", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
		retainedCollection := &mocks.CollectionBindingMock{}
		client.Database(mongo.DbName).(*mocks.DbBindingMock).On("Collection", mongo.RetainedPeginQuoteCollection).Return(retainedCollection)
		collection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPeginMongoRepository(conn)
		count, err := repo.DeleteQuotes(context.Background(), []string{test.AnyString})
		collection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
		require.Error(t, err)
		assert.Zero(t, count)
	})
	t.Run("Error when deletion count missmatch", func(t *testing.T) {
		client, quoteCollection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
		retainedCollection := &mocks.CollectionBindingMock{}
		client.Database(mongo.DbName).(*mocks.DbBindingMock).On("Collection", mongo.RetainedPeginQuoteCollection).Return(retainedCollection)
		quoteCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 4}, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPeginMongoRepository(conn)
		count, err := repo.DeleteQuotes(context.Background(), hashes)
		quoteCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
		require.ErrorContains(t, err, "pegin quote collections didn't match")
		assert.Zero(t, count)
	})
}
