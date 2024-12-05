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
	LbcAddress:            "0xc2A630c053D12D63d32b025082f6Ba268db18300",
	LpRskAddress:          "0x7c4890a0f1d4bbf2c669ac2d1effa185c505359b",
	BtcRefundAddress:      "n2Ge4xMVQKp5Hzzf8xTBJBLppRgjRZYYyq",
	RskRefundAddress:      "0x79568C2989232dcA1840087d73d403602364c0D4",
	LpBtcAddress:          "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe",
	CallFee:               entities.NewWei(100000000000000),
	PenaltyFee:            10000000000000,
	Nonce:                 6410832321595034747,
	DepositAddress:        "n2Ge4xMVQKp5Hzzf8xTBJBLppRgjRZYYyq",
	Value:                 entities.NewWei(5000000000000000),
	AgreementTimestamp:    1721944367,
	DepositDateLimit:      1721951567,
	DepositConfirmations:  4,
	TransferConfirmations: 2,
	TransferTime:          7200,
	ExpireDate:            1721958767,
	ExpireBlock:           5366409,
	GasFee:                entities.NewWei(4170000000000),
	ProductFeeAmount:      13,
}

var testRetainedPegoutQuote = quote.RetainedPegoutQuote{
	QuoteHash:          "27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f",
	DepositAddress:     "mkE1WWdiu5VgjfugomDk8GxV6JdEEEJR9s",
	Signature:          "5c9eab91c753355f87c19d09ea88b2fd02773981e513bc2821fed5ceba0d452a0a3d21e2252cb35348ce5c6803117e3abb62837beb8f5866a375ce66587d004b1c",
	RequiredLiquidity:  entities.NewWei(55),
	State:              quote.PegoutStateWaitingForDepositConfirmations,
	UserRskTxHash:      "0x6b2e1e4daf8cf00c5c3534b72cdeec3526e8a38f70c11e44888b6e4ae1ee7d38",
	LpBtcTxHash:        "6ac3779dc33ad52f3409cbb909bcd458745995496a2a3954406206f6e5d4cb0e",
	RefundPegoutTxHash: "0x8e773a2826e73f8e5792304379a7e46dff38f17089c6d344335e03537b31c2bc",
	BridgeRefundTxHash: "0x4f3f6f0664a732e4c907971e75c1e3fd8671461dcb53f566660432fc47255d8b",
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
		const expectedLog = "INSERT interaction with db: {PegoutQuote:{LbcAddress:0xc2A630c053D12D63d32b025082f6Ba268db18300 LpRskAddress:0x7c4890a0f1d4bbf2c669ac2d1effa185c505359b BtcRefundAddress:n2Ge4xMVQKp5Hzzf8xTBJBLppRgjRZYYyq RskRefundAddress:0x79568C2989232dcA1840087d73d403602364c0D4 LpBtcAddress:mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe CallFee:100000000000000 PenaltyFee:10000000000000 Nonce:6410832321595034747 DepositAddress:n2Ge4xMVQKp5Hzzf8xTBJBLppRgjRZYYyq Value:5000000000000000 AgreementTimestamp:1721944367 DepositDateLimit:1721951567 DepositConfirmations:4 TransferConfirmations:2 TransferTime:7200 ExpireDate:1721958767 ExpireBlock:5366409 GasFee:4170000000000 ProductFeeAmount:13} Hash:any value}"
		client, collection := getClientAndCollectionMocks(mongo.PegoutQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.MatchedBy(func(q mongo.StoredPegoutQuote) bool {
			return q.Hash == test.AnyString && reflect.TypeOf(quote.PegoutQuote{}).NumField() == test.CountNonZeroValues(q.PegoutQuote)
		})).Return(nil, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		defer assertDbInteractionLog(t, expectedLog)()
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
		const expectedLog = "READ interaction with db: {LbcAddress:0xc2A630c053D12D63d32b025082f6Ba268db18300 LpRskAddress:0x7c4890a0f1d4bbf2c669ac2d1effa185c505359b BtcRefundAddress:n2Ge4xMVQKp5Hzzf8xTBJBLppRgjRZYYyq RskRefundAddress:0x79568C2989232dcA1840087d73d403602364c0D4 LpBtcAddress:mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe CallFee:100000000000000 PenaltyFee:10000000000000 Nonce:6410832321595034747 DepositAddress:n2Ge4xMVQKp5Hzzf8xTBJBLppRgjRZYYyq Value:5000000000000000 AgreementTimestamp:1721944367 DepositDateLimit:1721951567 DepositConfirmations:4 TransferConfirmations:2 TransferTime:7200 ExpireDate:1721958767 ExpireBlock:5366409 GasFee:4170000000000 ProductFeeAmount:13}"
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, bson.D{primitive.E{Key: "hash", Value: test.AnyHash}}).
			Return(mongoDb.NewSingleResultFromDocument(mongo.StoredPegoutQuote{
				PegoutQuote: testPegoutQuote,
				Hash:        test.AnyString,
			}, nil, nil)).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, testPegoutQuote, *result)
	})
	t.Run("Db error when getting pegout quote", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(mongo.StoredPegoutQuote{}, assert.AnError, nil)).Once()
		result, err := repo.GetQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Pegout quote not present in db", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(mongo.StoredPegoutQuote{}, mongoDb.ErrNoDocuments, nil)).Once()
		result, err := repo.GetQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
	t.Run("Fail on invalid pegout quote hash", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		result, err := repo.GetQuote(context.Background(), test.AnyString)
		collection.AssertNotCalled(t, "FindOne")
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestPegoutMongoRepository_GetRetainedQuote(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
	log.SetLevel(log.DebugLevel)
	t.Run("Get retained pegout quote successfully", func(t *testing.T) {
		const expectedLog = "READ interaction with db: {QuoteHash:27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f DepositAddress:mkE1WWdiu5VgjfugomDk8GxV6JdEEEJR9s Signature:5c9eab91c753355f87c19d09ea88b2fd02773981e513bc2821fed5ceba0d452a0a3d21e2252cb35348ce5c6803117e3abb62837beb8f5866a375ce66587d004b1c RequiredLiquidity:55 State:WaitingForDepositConfirmations UserRskTxHash:0x6b2e1e4daf8cf00c5c3534b72cdeec3526e8a38f70c11e44888b6e4ae1ee7d38 LpBtcTxHash:6ac3779dc33ad52f3409cbb909bcd458745995496a2a3954406206f6e5d4cb0e RefundPegoutTxHash:0x8e773a2826e73f8e5792304379a7e46dff38f17089c6d344335e03537b31c2bc BridgeRefundTxHash:0x4f3f6f0664a732e4c907971e75c1e3fd8671461dcb53f566660432fc47255d8b}"
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, bson.D{primitive.E{Key: "quote_hash", Value: test.AnyHash}}).
			Return(mongoDb.NewSingleResultFromDocument(testRetainedPegoutQuote, nil, nil)).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, testRetainedPegoutQuote, *result)
	})
	t.Run("Db error when getting retained pegout quote", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(quote.RetainedPegoutQuote{}, assert.AnError, nil)).Once()
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Retained pegout quote not present in db", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(quote.RetainedPegoutQuote{}, mongoDb.ErrNoDocuments, nil)).Once()
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
	t.Run("Fail on invalid pegout quote hash", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyString)
		collection.AssertNotCalled(t, "FindOne")
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestPegoutMongoRepository_InsertRetainedQuote(t *testing.T) {
	t.Run("Insert retained pegout quote successfully", func(t *testing.T) {
		const expectedLog = "INSERT interaction with db: {QuoteHash:27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f DepositAddress:mkE1WWdiu5VgjfugomDk8GxV6JdEEEJR9s Signature:5c9eab91c753355f87c19d09ea88b2fd02773981e513bc2821fed5ceba0d452a0a3d21e2252cb35348ce5c6803117e3abb62837beb8f5866a375ce66587d004b1c RequiredLiquidity:55 State:WaitingForDepositConfirmations UserRskTxHash:0x6b2e1e4daf8cf00c5c3534b72cdeec3526e8a38f70c11e44888b6e4ae1ee7d38 LpBtcTxHash:6ac3779dc33ad52f3409cbb909bcd458745995496a2a3954406206f6e5d4cb0e RefundPegoutTxHash:0x8e773a2826e73f8e5792304379a7e46dff38f17089c6d344335e03537b31c2bc BridgeRefundTxHash:0x4f3f6f0664a732e4c907971e75c1e3fd8671461dcb53f566660432fc47255d8b}"
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.MatchedBy(func(q quote.RetainedPegoutQuote) bool {
			return q.QuoteHash == testRetainedPegoutQuote.QuoteHash && reflect.TypeOf(quote.RetainedPegoutQuote{}).NumField() == test.CountNonZeroValues(q)
		})).Return(nil, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		defer assertDbInteractionLog(t, expectedLog)()
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
		const expectedLog = "UPDATE interaction with db: {QuoteHash:27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f DepositAddress:updated value Signature:updated value RequiredLiquidity:200 State:SendPegoutFailed UserRskTxHash:updated value LpBtcTxHash:updated value RefundPegoutTxHash:updated value BridgeRefundTxHash:updated value}"
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		updatedQuote := testRetainedPegoutQuote
		updatedQuote.State = quote.PegoutStateSendPegoutFailed
		updatedQuote.Signature = updated
		updatedQuote.RefundPegoutTxHash = updated
		updatedQuote.LpBtcTxHash = updated
		updatedQuote.DepositAddress = updated
		updatedQuote.UserRskTxHash = updated
		updatedQuote.BridgeRefundTxHash = updated
		updatedQuote.RequiredLiquidity = entities.NewWei(200)
		collection.On("UpdateOne", mock.Anything,
			bson.D{primitive.E{Key: "quote_hash", Value: testRetainedPegoutQuote.QuoteHash}},
			bson.D{primitive.E{Key: "$set", Value: updatedQuote}},
		).Return(&mongoDb.UpdateResult{ModifiedCount: 1}, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		defer assertDbInteractionLog(t, expectedLog)()
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
			bson.D{primitive.E{Key: "hash", Value: bson.D{primitive.E{Key: "$in", Value: hashes}}}},
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
		const expectedLog = "READ interaction with db: [{QuoteHash:27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f DepositAddress:mkE1WWdiu5VgjfugomDk8GxV6JdEEEJR9s Signature:5c9eab91c753355f87c19d09ea88b2fd02773981e513bc2821fed5ceba0d452a0a3d21e2252cb35348ce5c6803117e3abb62837beb8f5866a375ce66587d004b1c RequiredLiquidity:55 State:WaitingForDepositConfirmations UserRskTxHash:0x6b2e1e4daf8cf00c5c3534b72cdeec3526e8a38f70c11e44888b6e4ae1ee7d38 LpBtcTxHash:6ac3779dc33ad52f3409cbb909bcd458745995496a2a3954406206f6e5d4cb0e RefundPegoutTxHash:0x8e773a2826e73f8e5792304379a7e46dff38f17089c6d344335e03537b31c2bc BridgeRefundTxHash:0x4f3f6f0664a732e4c907971e75c1e3fd8671461dcb53f566660432fc47255d8b} {QuoteHash:other hash DepositAddress:mkE1WWdiu5VgjfugomDk8GxV6JdEEEJR9s Signature:456 RequiredLiquidity:777 State:WaitingForDepositConfirmations UserRskTxHash:0x6b2e1e4daf8cf00c5c3534b72cdeec3526e8a38f70c11e44888b6e4ae1ee7d38 LpBtcTxHash:6ac3779dc33ad52f3409cbb909bcd458745995496a2a3954406206f6e5d4cb0e RefundPegoutTxHash:0x8e773a2826e73f8e5792304379a7e46dff38f17089c6d344335e03537b31c2bc BridgeRefundTxHash:0x4f3f6f0664a732e4c907971e75c1e3fd8671461dcb53f566660432fc47255d8b}]"
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		secondQuote := testRetainedPegoutQuote
		secondQuote.QuoteHash = "other hash"
		secondQuote.Signature = "456"
		secondQuote.RequiredLiquidity = entities.NewWei(777)
		collection.On("Find", mock.Anything,
			bson.D{primitive.E{Key: "state", Value: bson.D{primitive.E{Key: "$in", Value: states}}}},
		).Return(mongoDb.NewCursorFromDocuments([]any{testRetainedPegoutQuote, secondQuote}, nil, nil)).Once()
		defer assertDbInteractionLog(t, expectedLog)()
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
	t.Run("Should sanitize address properly", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client))
		collection.On("Find", mock.Anything,
			bson.M{"from": bson.M{"$regex": "0x1234567890abcdef1234567890abcdef12345678\\(a\\+\\)\\+", "$options": "i"}},
			options.Find().SetSort(bson.M{"timestamp": -1}),
		).Return(mongoDb.NewCursorFromDocuments([]any{testPegoutDeposit}, nil, nil)).Once()
		result, err := repo.ListPegoutDepositsByAddress(context.Background(), "0x1234567890abcdef1234567890abcdef12345678(a+)+")
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, []quote.PegoutDeposit{testPegoutDeposit}[0], result[0])
	})
}

func TestPegoutMongoRepository_UpsertPegoutDeposit(t *testing.T) {
	t.Run("Upsert pegout deposit successfully", func(t *testing.T) {
		now := time.Now().UTC()
		expectedLog := "UPSERT interaction with db: {TxHash:any value QuoteHash:any value Amount:1000 Timestamp:" + now.String() + " BlockNumber:790 From:any address}"
		client, collection := getClientAndCollectionMocks(mongo.DepositEventsCollection)
		newDeposit := testPegoutDeposit
		newDeposit.Amount = entities.NewWei(1000)
		newDeposit.Timestamp = now
		newDeposit.BlockNumber = 790
		collection.On("ReplaceOne", mock.Anything,
			bson.M{"tx_hash": testPegoutDeposit.TxHash},
			newDeposit,
			options.Replace().SetUpsert(true),
		).Return(&mongoDb.UpdateResult{ModifiedCount: 1}, nil).Once()
		conn := mongo.NewConnection(client)
		repo := mongo.NewPegoutMongoRepository(conn)
		defer test.AssertLogContains(t, expectedLog)()
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
		now := time.Now().UTC()
		expectedLog := "UPSERT interaction with db: [{TxHash:tx1 QuoteHash:quote1 Amount:1000 Timestamp:" + now.String() + " BlockNumber:790 From:any address} {TxHash:tx2 QuoteHash:quote2 Amount:2000 Timestamp:" + now.String() + " BlockNumber:791 From:any address}]"
		client, collection := getClientAndCollectionMocks(mongo.DepositEventsCollection)
		deposits := []quote.PegoutDeposit{
			{TxHash: "tx1", QuoteHash: "quote1", Amount: entities.NewWei(1000), Timestamp: now, BlockNumber: 790, From: test.AnyAddress},
			{TxHash: "tx2", QuoteHash: "quote2", Amount: entities.NewWei(2000), Timestamp: now, BlockNumber: 791, From: test.AnyAddress},
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
		defer test.AssertLogContains(t, expectedLog)()
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
		const expectedLog = "UPDATE interaction with db: [{QuoteHash:quote1 DepositAddress:any address Signature:any value RequiredLiquidity:1000 State:SendPegoutSucceeded UserRskTxHash: LpBtcTxHash: RefundPegoutTxHash: BridgeRefundTxHash:} {QuoteHash:quote2 DepositAddress:any address Signature:any value RequiredLiquidity:2000 State:SendPegoutFailed UserRskTxHash: LpBtcTxHash: RefundPegoutTxHash: BridgeRefundTxHash:}]"
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
		defer assertDbInteractionLog(t, expectedLog)()
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
