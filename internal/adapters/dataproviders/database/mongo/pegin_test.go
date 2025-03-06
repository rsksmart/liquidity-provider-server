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
	"time"
)

var testPeginQuote = quote.PeginQuote{
	FedBtcAddress:      "3LxPz39femVBL278mTiBvgzBNMVFqXssoH",
	LbcAddress:         "0xAA9cAf1e3967600578727F975F283446A3Da6612",
	LpRskAddress:       "0x4202bac9919c3412fc7c8be4e678e26279386603",
	BtcRefundAddress:   "171gGjg8NeLUonNSrFmgwkgT1jgqzXR6QX",
	RskRefundAddress:   "0xaD0DE1962ab903E06C725A1b343b7E8950a0Ff82",
	LpBtcAddress:       "17kksixYkbHeLy9okV16kr4eAxVhFkRhP",
	CallFee:            entities.NewWei(100000000000000),
	PenaltyFee:         entities.NewWei(10000000000000),
	ContractAddress:    "0xaD0DE1962ab903E06C725A1b343b7E8950a0Ff82",
	Data:               "010203",
	GasLimit:           21000,
	Nonce:              8373381263192041574,
	Value:              entities.NewWei(8000000000000000),
	AgreementTimestamp: 1727298699,
	TimeForDeposit:     3600,
	LpCallTime:         7200,
	Confirmations:      2,
	CallOnRegister:     true,
	GasFee:             entities.NewWei(1341211956000),
	ProductFeeAmount:   1,
}

var testRetainedPeginQuote = quote.RetainedPeginQuote{
	QuoteHash:           "8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819",
	DepositAddress:      "2N7Vw5f59V3o3bDcaJK5oA829LFTBYZHLoG",
	Signature:           "b24831aac7230910087d9818b378a31679be5e3991a7227cc160bc3add09e1645a26e9c740e3467f53953d7ec086c82bf8ef0eb03c118d0382ee6049a8f0119f1c",
	RequiredLiquidity:   entities.NewWei(100),
	State:               quote.PeginStateCallForUserSucceeded,
	UserBtcTxHash:       "619c4d69ccaa5f78aaa2284817cf070609ac40af3792916ca3d0ef82b14af75f",
	CallForUserTxHash:   "0x2c73de184c80797c04a655217d121588e8d5c228d3e0cc26187cb249123aa7c3",
	RegisterPeginTxHash: "0x3a0feaef4d803468ba5bfc1db78f4d2568de1b7cf002dec5991c469e6719db89",
}

func TestPeginMongoRepository_InsertQuote(t *testing.T) {
	t.Run("Insert pegin quote successfully", func(t *testing.T) {
		const expectedLog = "INSERT interaction with db: {PeginQuote:{FedBtcAddress:3LxPz39femVBL278mTiBvgzBNMVFqXssoH LbcAddress:0xAA9cAf1e3967600578727F975F283446A3Da6612 LpRskAddress:0x4202bac9919c3412fc7c8be4e678e26279386603 BtcRefundAddress:171gGjg8NeLUonNSrFmgwkgT1jgqzXR6QX RskRefundAddress:0xaD0DE1962ab903E06C725A1b343b7E8950a0Ff82 LpBtcAddress:17kksixYkbHeLy9okV16kr4eAxVhFkRhP CallFee:100000000000000 PenaltyFee:10000000000000 ContractAddress:0xaD0DE1962ab903E06C725A1b343b7E8950a0Ff82 Data:010203 GasLimit:21000 Nonce:8373381263192041574 Value:8000000000000000 AgreementTimestamp:1727298699 TimeForDeposit:3600 LpCallTime:7200 Confirmations:2 CallOnRegister:true GasFee:1341211956000 ProductFeeAmount:1} Hash:any value}"
		client, collection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.MatchedBy(func(q mongo.StoredPeginQuote) bool {
			return q.Hash == test.AnyString && reflect.TypeOf(quote.PeginQuote{}).NumField() == test.CountNonZeroValues(q.PeginQuote)
		})).Return(nil, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)
		defer assertDbInteractionLog(t, expectedLog)()
		err := repo.InsertQuote(context.Background(), test.AnyString, testPeginQuote)
		collection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Db error when inserting pegin quote", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
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
		const expectedLog = "READ interaction with db: {FedBtcAddress:3LxPz39femVBL278mTiBvgzBNMVFqXssoH LbcAddress:0xAA9cAf1e3967600578727F975F283446A3Da6612 LpRskAddress:0x4202bac9919c3412fc7c8be4e678e26279386603 BtcRefundAddress:171gGjg8NeLUonNSrFmgwkgT1jgqzXR6QX RskRefundAddress:0xaD0DE1962ab903E06C725A1b343b7E8950a0Ff82 LpBtcAddress:17kksixYkbHeLy9okV16kr4eAxVhFkRhP CallFee:100000000000000 PenaltyFee:10000000000000 ContractAddress:0xaD0DE1962ab903E06C725A1b343b7E8950a0Ff82 Data:010203 GasLimit:21000 Nonce:8373381263192041574 Value:8000000000000000 AgreementTimestamp:1727298699 TimeForDeposit:3600 LpCallTime:7200 Confirmations:2 CallOnRegister:true GasFee:1341211956000 ProductFeeAmount:1}"
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, bson.D{primitive.E{Key: "hash", Value: test.AnyHash}}).
			Return(mongoDb.NewSingleResultFromDocument(mongo.StoredPeginQuote{
				PeginQuote: testPeginQuote,
				Hash:       test.AnyString,
			}, nil, nil)).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, testPeginQuote, *result)
	})
	t.Run("Db error when getting pegin quote", func(t *testing.T) {
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(mongo.StoredPeginQuote{}, assert.AnError, nil)).Once()
		result, err := repo.GetQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Pegin quote not present in db", func(t *testing.T) {
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(mongo.StoredPeginQuote{}, mongoDb.ErrNoDocuments, nil)).Once()
		result, err := repo.GetQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
	t.Run("Fail on invalid pegin quote hash", func(t *testing.T) {
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		result, err := repo.GetQuote(context.Background(), test.AnyString)
		collection.AssertNotCalled(t, "FindOne")
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestPeginMongoRepository_GetRetainedQuote(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
	log.SetLevel(log.DebugLevel)
	t.Run("Get retained pegin quote successfully", func(t *testing.T) {
		const expectedLog = "READ interaction with db: {QuoteHash:8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819 DepositAddress:2N7Vw5f59V3o3bDcaJK5oA829LFTBYZHLoG Signature:b24831aac7230910087d9818b378a31679be5e3991a7227cc160bc3add09e1645a26e9c740e3467f53953d7ec086c82bf8ef0eb03c118d0382ee6049a8f0119f1c RequiredLiquidity:100 State:CallForUserSucceeded UserBtcTxHash:619c4d69ccaa5f78aaa2284817cf070609ac40af3792916ca3d0ef82b14af75f CallForUserTxHash:0x2c73de184c80797c04a655217d121588e8d5c228d3e0cc26187cb249123aa7c3 RegisterPeginTxHash:0x3a0feaef4d803468ba5bfc1db78f4d2568de1b7cf002dec5991c469e6719db89}"
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, bson.D{primitive.E{Key: "quote_hash", Value: test.AnyHash}}).
			Return(mongoDb.NewSingleResultFromDocument(testRetainedPeginQuote, nil, nil)).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, testRetainedPeginQuote, *result)
	})
	t.Run("Db error when getting retained pegin quote", func(t *testing.T) {
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(quote.RetainedPeginQuote{}, assert.AnError, nil)).Once()
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Retained pegin quote not present in db", func(t *testing.T) {
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(quote.RetainedPeginQuote{}, mongoDb.ErrNoDocuments, nil)).Once()
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
	t.Run("Fail on invalid pegin quote hash", func(t *testing.T) {
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyString)
		collection.AssertNotCalled(t, "FindOne")
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestPeginMongoRepository_InsertRetainedQuote(t *testing.T) {
	t.Run("Insert retained pegin quote successfully", func(t *testing.T) {
		const expectedLog = "INSERT interaction with db: {QuoteHash:8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819 DepositAddress:2N7Vw5f59V3o3bDcaJK5oA829LFTBYZHLoG Signature:b24831aac7230910087d9818b378a31679be5e3991a7227cc160bc3add09e1645a26e9c740e3467f53953d7ec086c82bf8ef0eb03c118d0382ee6049a8f0119f1c RequiredLiquidity:100 State:CallForUserSucceeded UserBtcTxHash:619c4d69ccaa5f78aaa2284817cf070609ac40af3792916ca3d0ef82b14af75f CallForUserTxHash:0x2c73de184c80797c04a655217d121588e8d5c228d3e0cc26187cb249123aa7c3 RegisterPeginTxHash:0x3a0feaef4d803468ba5bfc1db78f4d2568de1b7cf002dec5991c469e6719db89}"
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.MatchedBy(func(q quote.RetainedPeginQuote) bool {
			return q.QuoteHash == testRetainedPeginQuote.QuoteHash && reflect.TypeOf(quote.RetainedPeginQuote{}).NumField() == test.CountNonZeroValues(q)
		})).Return(nil, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)
		defer assertDbInteractionLog(t, expectedLog)()
		err := repo.InsertRetainedQuote(context.Background(), testRetainedPeginQuote)
		collection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Db error when inserting retained pegin quote", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)
		err := repo.InsertRetainedQuote(context.Background(), testRetainedPeginQuote)
		collection.AssertExpectations(t)
		require.Error(t, err)
	})
}

func TestPeginMongoRepository_UpdateRetainedQuote(t *testing.T) {
	const updated = "updated value"
	t.Run("Update retained pegin quote successfully", func(t *testing.T) {
		const expectedLog = "UPDATE interaction with db: {QuoteHash:8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819 DepositAddress:updated value Signature:updated value RequiredLiquidity:200 State:CallForUserFailed UserBtcTxHash:updated value CallForUserTxHash:updated value RegisterPeginTxHash:updated value}"
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
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)
		defer assertDbInteractionLog(t, expectedLog)()
		err := repo.UpdateRetainedQuote(context.Background(), updatedQuote)
		collection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Db error when updating retained pegin quote", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)
		err := repo.UpdateRetainedQuote(context.Background(), testRetainedPeginQuote)
		collection.AssertExpectations(t)
		require.Error(t, err)
	})
	t.Run("Retained pegin quote to update not found", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).
			Return(&mongoDb.UpdateResult{ModifiedCount: 0}, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)
		err := repo.UpdateRetainedQuote(context.Background(), testRetainedPeginQuote)
		collection.AssertExpectations(t)
		require.ErrorIs(t, err, usecases.QuoteNotFoundError)
	})
	t.Run("Update more than one retained pegin quote", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).
			Return(&mongoDb.UpdateResult{ModifiedCount: 2}, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
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
		const expectedLog = "READ interaction with db: [{QuoteHash:8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819 DepositAddress:2N7Vw5f59V3o3bDcaJK5oA829LFTBYZHLoG Signature:b24831aac7230910087d9818b378a31679be5e3991a7227cc160bc3add09e1645a26e9c740e3467f53953d7ec086c82bf8ef0eb03c118d0382ee6049a8f0119f1c RequiredLiquidity:100 State:CallForUserSucceeded UserBtcTxHash:619c4d69ccaa5f78aaa2284817cf070609ac40af3792916ca3d0ef82b14af75f CallForUserTxHash:0x2c73de184c80797c04a655217d121588e8d5c228d3e0cc26187cb249123aa7c3 RegisterPeginTxHash:0x3a0feaef4d803468ba5bfc1db78f4d2568de1b7cf002dec5991c469e6719db89} {QuoteHash:second DepositAddress:2N7Vw5f59V3o3bDcaJK5oA829LFTBYZHLoG Signature:123 RequiredLiquidity:777 State:CallForUserSucceeded UserBtcTxHash:619c4d69ccaa5f78aaa2284817cf070609ac40af3792916ca3d0ef82b14af75f CallForUserTxHash:0x2c73de184c80797c04a655217d121588e8d5c228d3e0cc26187cb249123aa7c3 RegisterPeginTxHash:0x3a0feaef4d803468ba5bfc1db78f4d2568de1b7cf002dec5991c469e6719db89}]"
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		secondQuote := testRetainedPeginQuote
		secondQuote.QuoteHash = "second"
		secondQuote.Signature = "123"
		secondQuote.RequiredLiquidity = entities.NewWei(777)
		collection.On("Find", mock.Anything,
			bson.D{primitive.E{Key: "state", Value: bson.D{primitive.E{Key: "$in", Value: states}}}},
		).Return(mongoDb.NewCursorFromDocuments([]any{testRetainedPeginQuote, secondQuote}, nil, nil)).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetRetainedQuoteByState(context.Background(), states...)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, []quote.RetainedPeginQuote{testRetainedPeginQuote, secondQuote}, result)
	})
	t.Run("Db error when getting retained pegin quotes by state", func(t *testing.T) {
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
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
		parsedClientMock, ok := client.Database(mongo.DbName).(*mocks.DbBindingMock)
		require.True(t, ok)
		parsedClientMock.On("Collection", mongo.RetainedPeginQuoteCollection).Return(retainedCollection)
		quoteCollection.On("DeleteMany", mock.Anything, bson.D{primitive.E{Key: "hash", Value: bson.D{primitive.E{Key: "$in", Value: hashes}}}}).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything,
			bson.D{primitive.E{Key: "quote_hash", Value: bson.D{primitive.E{Key: "$in", Value: hashes}}}},
		).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
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
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		count, err := repo.DeleteQuotes(context.Background(), []string{test.AnyString})
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Zero(t, count)
	})
	t.Run("Db error when deleting retained pegin quotes", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
		retainedCollection := &mocks.CollectionBindingMock{}
		parsedClient, ok := client.Database(mongo.DbName).(*mocks.DbBindingMock)
		require.True(t, ok)
		parsedClient.On("Collection", mongo.RetainedPeginQuoteCollection).Return(retainedCollection)
		collection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		count, err := repo.DeleteQuotes(context.Background(), []string{test.AnyString})
		collection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
		require.Error(t, err)
		assert.Zero(t, count)
	})
	t.Run("Error when deletion count missmatch", func(t *testing.T) {
		client, quoteCollection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
		retainedCollection := &mocks.CollectionBindingMock{}
		parsedClientMock, ok := client.Database(mongo.DbName).(*mocks.DbBindingMock)
		require.True(t, ok)
		parsedClientMock.On("Collection", mongo.RetainedPeginQuoteCollection).Return(retainedCollection)
		quoteCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 4}, nil).Once()
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		count, err := repo.DeleteQuotes(context.Background(), hashes)
		quoteCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
		require.ErrorContains(t, err, "pegin quote collections didn't match")
		assert.Zero(t, count)
	})
}

func TestPeginMongoRepository_GetQuotes(t *testing.T) {
	t.Run("Successfully retrieves quotes", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
		log.SetLevel(log.DebugLevel)
		hashes := []string{testRetainedPegoutQuote.QuoteHash}
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("Find", mock.Anything,
			bson.M{"hash": bson.M{"$in": hashes}},
		).Return(mongoDb.NewCursorFromDocuments([]any{testPeginQuote}, nil, nil)).Once()
		result, err := repo.GetQuotes(context.Background(), hashes)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, []quote.PeginQuote{testPeginQuote}, result)
	})

	t.Run("Fails validation for hashes", func(t *testing.T) {
		client, _ := getClientAndCollectionMocks(mongo.PeginQuoteCollection)

		invalidHashes := []string{"invalidHash"}
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)

		_, err := repo.GetQuotes(context.Background(), invalidHashes)
		require.Error(t, err)
		assert.Equal(t, "invalid quote hash length: expected 64 characters, got 11", err.Error())
	})

	t.Run("No quotes found", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)

		expectedHashes := []string{testRetainedPegoutQuote.QuoteHash}
		collection.On("Find", mock.Anything, bson.M{"hash": bson.M{"$in": expectedHashes}}).Return(nil, mongoDb.ErrNoDocuments).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)

		quotes, err := repo.GetQuotes(context.Background(), expectedHashes)
		require.NoError(t, err)
		assert.Nil(t, quotes)
	})
}
