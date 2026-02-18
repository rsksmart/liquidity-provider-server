package mongo_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
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
}

var testRetainedPeginQuote = quote.RetainedPeginQuote{
	QuoteHash:             "8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819",
	DepositAddress:        "2N7Vw5f59V3o3bDcaJK5oA829LFTBYZHLoG",
	Signature:             "b24831aac7230910087d9818b378a31679be5e3991a7227cc160bc3add09e1645a26e9c740e3467f53953d7ec086c82bf8ef0eb03c118d0382ee6049a8f0119f1c",
	RequiredLiquidity:     entities.NewWei(100),
	State:                 quote.PeginStateCallForUserSucceeded,
	UserBtcTxHash:         "619c4d69ccaa5f78aaa2284817cf070609ac40af3792916ca3d0ef82b14af75f",
	CallForUserTxHash:     "0x2c73de184c80797c04a655217d121588e8d5c228d3e0cc26187cb249123aa7c3",
	RegisterPeginTxHash:   "0x3a0feaef4d803468ba5bfc1db78f4d2568de1b7cf002dec5991c469e6719db89",
	CallForUserGasUsed:    85000,
	CallForUserGasPrice:   entities.NewWei(21000000000),
	RegisterPeginGasUsed:  65000,
	RegisterPeginGasPrice: entities.NewWei(21000000000),
	OwnerAccountAddress:   "0x233845a26a4dA08E16218e7B401501D048670674",
}

var testPeginCreationData = quote.PeginCreationData{
	GasPrice:      entities.NewWei(55),
	FeePercentage: utils.NewBigFloat64(1.5),
	FixedFee:      entities.NewWei(100000),
}

func TestPeginMongoRepository_InsertQuote(t *testing.T) {
	t.Run("Insert pegin quote successfully", func(t *testing.T) {
		const expectedLog = "INSERT interaction with db: {PeginQuote:{FedBtcAddress:3LxPz39femVBL278mTiBvgzBNMVFqXssoH LbcAddress:0xAA9cAf1e3967600578727F975F283446A3Da6612 LpRskAddress:0x4202bac9919c3412fc7c8be4e678e26279386603 BtcRefundAddress:171gGjg8NeLUonNSrFmgwkgT1jgqzXR6QX RskRefundAddress:0xaD0DE1962ab903E06C725A1b343b7E8950a0Ff82 LpBtcAddress:17kksixYkbHeLy9okV16kr4eAxVhFkRhP CallFee:100000000000000 PenaltyFee:10000000000000 ContractAddress:0xaD0DE1962ab903E06C725A1b343b7E8950a0Ff82 Data:010203 GasLimit:21000 Nonce:8373381263192041574 Value:8000000000000000 AgreementTimestamp:1727298699 TimeForDeposit:3600 LpCallTime:7200 Confirmations:2 CallOnRegister:true GasFee:1341211956000} Hash:any value}"
		client, db := getClientAndDatabaseMocks()
		quoteCollection := &mocks.CollectionBindingMock{}
		creationDataCollection := &mocks.CollectionBindingMock{}
		db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(quoteCollection)
		db.EXPECT().Collection(mongo.PeginCreationDataCollection).Return(creationDataCollection)
		quoteCollection.On("InsertOne", mock.Anything, mock.MatchedBy(func(q mongo.StoredPeginQuote) bool {
			return q.Hash == test.AnyString && reflect.TypeOf(quote.PeginQuote{}).NumField() == test.CountNonZeroValues(q.PeginQuote)
		})).Return(nil, nil).Once()
		creationDataCollection.EXPECT().InsertOne(mock.Anything, mock.MatchedBy(func(q mongo.StoredPeginCreationData) bool {
			return q.Hash == test.AnyString && reflect.TypeOf(quote.PeginCreationData{}).NumField() == test.CountNonZeroValues(q.PeginCreationData)
		})).Return(nil, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)
		defer assertDbInteractionLog(t, expectedLog)()
		createdQuote := quote.CreatedPeginQuote{Hash: test.AnyString, Quote: testPeginQuote, CreationData: testPeginCreationData}
		err := repo.InsertQuote(context.Background(), createdQuote)
		quoteCollection.AssertExpectations(t)
		creationDataCollection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Db error when inserting pegin quote", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		quoteCollection := &mocks.CollectionBindingMock{}
		creationDataCollection := &mocks.CollectionBindingMock{}
		db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(quoteCollection)
		quoteCollection.On("InsertOne", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)
		createdQuote := quote.CreatedPeginQuote{Hash: test.AnyString, Quote: testPeginQuote, CreationData: testPeginCreationData}
		err := repo.InsertQuote(context.Background(), createdQuote)
		quoteCollection.AssertExpectations(t)
		creationDataCollection.AssertNotCalled(t, "InsertOne")
		require.Error(t, err)
	})

	t.Run("Db error when inserting pegin creation data", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		quoteCollection := &mocks.CollectionBindingMock{}
		creationDataCollection := &mocks.CollectionBindingMock{}
		db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(quoteCollection)
		db.EXPECT().Collection(mongo.PeginCreationDataCollection).Return(creationDataCollection)
		quoteCollection.On("InsertOne", mock.Anything, mock.Anything).Return(nil, nil).Once()
		creationDataCollection.On("InsertOne", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)
		createdQuote := quote.CreatedPeginQuote{
			Hash:         test.AnyString,
			Quote:        testPeginQuote,
			CreationData: testPeginCreationData,
		}
		err := repo.InsertQuote(context.Background(), createdQuote)
		quoteCollection.AssertExpectations(t)
		creationDataCollection.AssertExpectations(t)
		require.Error(t, err)
	})
}

func TestPeginMongoRepository_GetQuote(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
	log.SetLevel(log.DebugLevel)
	t.Run("Get pegin quote successfully", func(t *testing.T) {
		const expectedLog = "READ interaction with db: {FedBtcAddress:3LxPz39femVBL278mTiBvgzBNMVFqXssoH LbcAddress:0xAA9cAf1e3967600578727F975F283446A3Da6612 LpRskAddress:0x4202bac9919c3412fc7c8be4e678e26279386603 BtcRefundAddress:171gGjg8NeLUonNSrFmgwkgT1jgqzXR6QX RskRefundAddress:0xaD0DE1962ab903E06C725A1b343b7E8950a0Ff82 LpBtcAddress:17kksixYkbHeLy9okV16kr4eAxVhFkRhP CallFee:100000000000000 PenaltyFee:10000000000000 ContractAddress:0xaD0DE1962ab903E06C725A1b343b7E8950a0Ff82 Data:010203 GasLimit:21000 Nonce:8373381263192041574 Value:8000000000000000 AgreementTimestamp:1727298699 TimeForDeposit:3600 LpCallTime:7200 Confirmations:2 CallOnRegister:true GasFee:1341211956000}"
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

//nolint:funlen
func TestPeginMongoRepository_GetRetainedQuote(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
	log.SetLevel(log.DebugLevel)
	t.Run("Get retained pegin quote successfully", func(t *testing.T) {
		const expectedLog = "READ interaction with db: {QuoteHash:8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819 DepositAddress:2N7Vw5f59V3o3bDcaJK5oA829LFTBYZHLoG Signature:b24831aac7230910087d9818b378a31679be5e3991a7227cc160bc3add09e1645a26e9c740e3467f53953d7ec086c82bf8ef0eb03c118d0382ee6049a8f0119f1c RequiredLiquidity:100 State:CallForUserSucceeded UserBtcTxHash:619c4d69ccaa5f78aaa2284817cf070609ac40af3792916ca3d0ef82b14af75f CallForUserTxHash:0x2c73de184c80797c04a655217d121588e8d5c228d3e0cc26187cb249123aa7c3 RegisterPeginTxHash:0x3a0feaef4d803468ba5bfc1db78f4d2568de1b7cf002dec5991c469e6719db89 CallForUserGasUsed:85000 CallForUserGasPrice:21000000000 RegisterPeginGasUsed:65000 RegisterPeginGasPrice:21000000000 OwnerAccountAddress:0x233845a26a4dA08E16218e7B401501D048670674}"
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
	t.Run("FillZeroValues is applied to retained pegin quote with missing gas fields", func(t *testing.T) {
		// Create a BSON document that represents what an old database record would look like
		// Mock strategy for similar tests did not work for mock limitations on unmarshalling into a struct
		oldBsonDocument := bson.D{
			{Key: "quote_hash", Value: testRetainedPeginQuote.QuoteHash},
			{Key: "deposit_address", Value: testRetainedPeginQuote.DepositAddress},
			{Key: "signature", Value: testRetainedPeginQuote.Signature},
			{Key: "required_liquidity", Value: testRetainedPeginQuote.RequiredLiquidity.String()},
			{Key: "state", Value: testRetainedPeginQuote.State},
			{Key: "user_btc_tx_hash", Value: testRetainedPeginQuote.UserBtcTxHash},
			{Key: "call_for_user_tx_hash", Value: testRetainedPeginQuote.CallForUserTxHash},
			{Key: "register_pegin_tx_hash", Value: testRetainedPeginQuote.RegisterPeginTxHash},
			{Key: "call_for_user_gas_used", Value: testRetainedPeginQuote.CallForUserGasUsed},
			{Key: "register_pegin_gas_used", Value: testRetainedPeginQuote.RegisterPeginGasUsed},
			{Key: "owner_account_address", Value: testRetainedPeginQuote.OwnerAccountAddress},
			// NOTE: call_for_user_gas_price and register_pegin_gas_price are MISSING
		}

		singleResult := mongoDb.NewSingleResultFromDocument(oldBsonDocument, nil, nil)

		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, bson.D{primitive.E{Key: "quote_hash", Value: test.AnyHash}}).
			Return(singleResult).Once()

		result, err := repo.GetRetainedQuote(context.Background(), test.AnyHash)

		collection.AssertExpectations(t)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.NotNil(t, result.CallForUserGasPrice, "CallForUserGasPrice should not be nil after normalization")
		assert.NotNil(t, result.RegisterPeginGasPrice, "RegisterPeginGasPrice should not be nil after normalization")
		assert.Equal(t, entities.NewWei(0), result.CallForUserGasPrice)
		assert.Equal(t, entities.NewWei(0), result.RegisterPeginGasPrice)

		assert.Equal(t, testRetainedPeginQuote.QuoteHash, result.QuoteHash)
		assert.Equal(t, testRetainedPeginQuote.CallForUserGasUsed, result.CallForUserGasUsed)
		assert.Equal(t, testRetainedPeginQuote.RegisterPeginGasUsed, result.RegisterPeginGasUsed)
	})
}

func TestPeginMongoRepository_InsertRetainedQuote(t *testing.T) {
	t.Run("Insert retained pegin quote successfully", func(t *testing.T) {
		const expectedLog = "INSERT interaction with db: {QuoteHash:8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819 DepositAddress:2N7Vw5f59V3o3bDcaJK5oA829LFTBYZHLoG Signature:b24831aac7230910087d9818b378a31679be5e3991a7227cc160bc3add09e1645a26e9c740e3467f53953d7ec086c82bf8ef0eb03c118d0382ee6049a8f0119f1c RequiredLiquidity:100 State:CallForUserSucceeded UserBtcTxHash:619c4d69ccaa5f78aaa2284817cf070609ac40af3792916ca3d0ef82b14af75f CallForUserTxHash:0x2c73de184c80797c04a655217d121588e8d5c228d3e0cc26187cb249123aa7c3 RegisterPeginTxHash:0x3a0feaef4d803468ba5bfc1db78f4d2568de1b7cf002dec5991c469e6719db89 CallForUserGasUsed:85000 CallForUserGasPrice:21000000000 RegisterPeginGasUsed:65000 RegisterPeginGasPrice:21000000000 OwnerAccountAddress:0x233845a26a4dA08E16218e7B401501D048670674}"
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.MatchedBy(func(q quote.RetainedPeginQuote) bool {
			test.AssertMaxZeroValues(t, q, 1)
			return q.QuoteHash == testRetainedPeginQuote.QuoteHash
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
		const expectedLog = "UPDATE interaction with db: {QuoteHash:8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819 DepositAddress:updated value Signature:updated value RequiredLiquidity:200 State:CallForUserFailed UserBtcTxHash:updated value CallForUserTxHash:updated value RegisterPeginTxHash:updated value CallForUserGasUsed:85000 CallForUserGasPrice:21000000000 RegisterPeginGasUsed:65000 RegisterPeginGasPrice:21000000000 OwnerAccountAddress:updated value}"
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		updatedQuote := testRetainedPeginQuote
		updatedQuote.State = quote.PeginStateCallForUserFailed
		updatedQuote.Signature = updated
		updatedQuote.RegisterPeginTxHash = updated
		updatedQuote.CallForUserTxHash = updated
		updatedQuote.DepositAddress = updated
		updatedQuote.UserBtcTxHash = updated
		updatedQuote.OwnerAccountAddress = updated
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

//nolint:funlen
func TestPeginMongoRepository_GetRetainedQuoteByState(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
	log.SetLevel(log.DebugLevel)
	states := []quote.PeginState{quote.PeginStateCallForUserSucceeded, quote.PeginStateCallForUserFailed}
	t.Run("Get retained pegin quotes by state successfully", func(t *testing.T) {
		const expectedLog = "READ interaction with db: [{QuoteHash:8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819 DepositAddress:2N7Vw5f59V3o3bDcaJK5oA829LFTBYZHLoG Signature:b24831aac7230910087d9818b378a31679be5e3991a7227cc160bc3add09e1645a26e9c740e3467f53953d7ec086c82bf8ef0eb03c118d0382ee6049a8f0119f1c RequiredLiquidity:100 State:CallForUserSucceeded UserBtcTxHash:619c4d69ccaa5f78aaa2284817cf070609ac40af3792916ca3d0ef82b14af75f CallForUserTxHash:0x2c73de184c80797c04a655217d121588e8d5c228d3e0cc26187cb249123aa7c3 RegisterPeginTxHash:0x3a0feaef4d803468ba5bfc1db78f4d2568de1b7cf002dec5991c469e6719db89 CallForUserGasUsed:85000 CallForUserGasPrice:21000000000 RegisterPeginGasUsed:65000 RegisterPeginGasPrice:21000000000 OwnerAccountAddress:0x233845a26a4dA08E16218e7B401501D048670674} {QuoteHash:second DepositAddress:2N7Vw5f59V3o3bDcaJK5oA829LFTBYZHLoG Signature:123 RequiredLiquidity:777 State:CallForUserSucceeded UserBtcTxHash:619c4d69ccaa5f78aaa2284817cf070609ac40af3792916ca3d0ef82b14af75f CallForUserTxHash:0x2c73de184c80797c04a655217d121588e8d5c228d3e0cc26187cb249123aa7c3 RegisterPeginTxHash:0x3a0feaef4d803468ba5bfc1db78f4d2568de1b7cf002dec5991c469e6719db89 CallForUserGasUsed:85000 CallForUserGasPrice:21000000000 RegisterPeginGasUsed:65000 RegisterPeginGasPrice:21000000000 OwnerAccountAddress:0x233845a26a4dA08E16218e7B401501D048670674}]"
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
	t.Run("FillZeroValues is applied to retained pegin quotes with missing gas fields", func(t *testing.T) {
		// Mock strategy for similar tests did not work for mock limitations on unmarshalling into a struct
		firstOldDocument := bson.D{
			{Key: "quote_hash", Value: "first"},
			{Key: "deposit_address", Value: testRetainedPeginQuote.DepositAddress},
			{Key: "signature", Value: testRetainedPeginQuote.Signature},
			{Key: "required_liquidity", Value: testRetainedPeginQuote.RequiredLiquidity.String()},
			{Key: "state", Value: testRetainedPeginQuote.State},
			{Key: "user_btc_tx_hash", Value: testRetainedPeginQuote.UserBtcTxHash},
			{Key: "call_for_user_tx_hash", Value: testRetainedPeginQuote.CallForUserTxHash},
			{Key: "register_pegin_tx_hash", Value: testRetainedPeginQuote.RegisterPeginTxHash},
			{Key: "call_for_user_gas_used", Value: testRetainedPeginQuote.CallForUserGasUsed},
			{Key: "register_pegin_gas_used", Value: testRetainedPeginQuote.RegisterPeginGasUsed},
			{Key: "owner_account_address", Value: testRetainedPeginQuote.OwnerAccountAddress},
			// NOTE: call_for_user_gas_price and register_pegin_gas_price are MISSING
		}

		secondOldDocument := bson.D{
			{Key: "quote_hash", Value: "second"},
			{Key: "deposit_address", Value: testRetainedPeginQuote.DepositAddress},
			{Key: "signature", Value: "different_signature"},
			{Key: "required_liquidity", Value: entities.NewWei(500).String()},
			{Key: "state", Value: quote.PeginStateCallForUserFailed},
			{Key: "user_btc_tx_hash", Value: testRetainedPeginQuote.UserBtcTxHash},
			{Key: "call_for_user_tx_hash", Value: testRetainedPeginQuote.CallForUserTxHash},
			{Key: "register_pegin_tx_hash", Value: testRetainedPeginQuote.RegisterPeginTxHash},
			{Key: "call_for_user_gas_used", Value: uint64(75000)},
			{Key: "register_pegin_gas_used", Value: uint64(55000)},
			{Key: "owner_account_address", Value: testRetainedPeginQuote.OwnerAccountAddress},
			// NOTE: call_for_user_gas_price and register_pegin_gas_price are MISSING
		}

		cursor, err := mongoDb.NewCursorFromDocuments([]any{firstOldDocument, secondOldDocument}, nil, nil)
		require.NoError(t, err)

		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("Find", mock.Anything,
			bson.D{primitive.E{Key: "state", Value: bson.D{primitive.E{Key: "$in", Value: states}}}},
		).Return(cursor, nil).Once()

		result, err := repo.GetRetainedQuoteByState(context.Background(), states...)

		collection.AssertExpectations(t)
		require.NoError(t, err)
		require.Len(t, result, 2)

		// Verify normalization applied to first document
		assert.NotNil(t, result[0].CallForUserGasPrice, "CallForUserGasPrice should not be nil after normalization")
		assert.NotNil(t, result[0].RegisterPeginGasPrice, "RegisterPeginGasPrice should not be nil after normalization")
		assert.Equal(t, entities.NewWei(0), result[0].CallForUserGasPrice)
		assert.Equal(t, entities.NewWei(0), result[0].RegisterPeginGasPrice)
		assert.Equal(t, "first", result[0].QuoteHash)
		assert.Equal(t, testRetainedPeginQuote.CallForUserGasUsed, result[0].CallForUserGasUsed)

		// Verify normalization applied to second document
		assert.NotNil(t, result[1].CallForUserGasPrice, "CallForUserGasPrice should not be nil after normalization")
		assert.NotNil(t, result[1].RegisterPeginGasPrice, "RegisterPeginGasPrice should not be nil after normalization")
		assert.Equal(t, entities.NewWei(0), result[1].CallForUserGasPrice)
		assert.Equal(t, entities.NewWei(0), result[1].RegisterPeginGasPrice)
		assert.Equal(t, "second", result[1].QuoteHash)
		assert.Equal(t, uint64(75000), result[1].CallForUserGasUsed)
		assert.Equal(t, uint64(55000), result[1].RegisterPeginGasUsed)
	})
}

// nolint:funlen
func TestPeginMongoRepository_DeleteQuotes(t *testing.T) {
	var hashes = []string{"pegin1", "pegin2", "pegin3"}
	log.SetLevel(log.DebugLevel)
	t.Run("Delete quotes successfully", func(t *testing.T) {
		client, quoteCollection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
		retainedCollection := &mocks.CollectionBindingMock{}
		creationDataCollection := &mocks.CollectionBindingMock{}
		parsedClientMock, ok := client.Database(mongo.DbName).(*mocks.DbBindingMock)
		require.True(t, ok)
		parsedClientMock.On("Collection", mongo.RetainedPeginQuoteCollection).Return(retainedCollection)
		parsedClientMock.On("Collection", mongo.PeginCreationDataCollection).Return(creationDataCollection)
		quoteCollection.On("DeleteMany", mock.Anything, bson.D{primitive.E{Key: "hash", Value: bson.D{primitive.E{Key: "$in", Value: hashes}}}}).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything,
			bson.D{primitive.E{Key: "quote_hash", Value: bson.D{primitive.E{Key: "$in", Value: hashes}}}},
		).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		creationDataCollection.On("DeleteMany", mock.Anything,
			bson.D{primitive.E{Key: "hash", Value: bson.D{primitive.E{Key: "$in", Value: hashes}}}},
		).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)
		defer assertDbInteractionLog(t, mongo.Delete)()
		count, err := repo.DeleteQuotes(context.Background(), hashes)
		quoteCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
		creationDataCollection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, uint(9), count)
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
	t.Run("Db error when deleting pegin creation data", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
		retainedCollection := &mocks.CollectionBindingMock{}
		creationDataCollection := &mocks.CollectionBindingMock{}
		parsedClient, ok := client.Database(mongo.DbName).(*mocks.DbBindingMock)
		require.True(t, ok)
		parsedClient.On("Collection", mongo.RetainedPeginQuoteCollection).Return(retainedCollection)
		parsedClient.On("Collection", mongo.PeginCreationDataCollection).Return(creationDataCollection)
		collection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		creationDataCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		count, err := repo.DeleteQuotes(context.Background(), []string{test.AnyString})
		collection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
		creationDataCollection.AssertExpectations(t)
		require.Error(t, err)
		assert.Zero(t, count)
	})
	t.Run("Error when deletion count missmatch", func(t *testing.T) {
		client, quoteCollection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)
		retainedCollection := &mocks.CollectionBindingMock{}
		creationDataCollection := &mocks.CollectionBindingMock{}
		parsedClientMock, ok := client.Database(mongo.DbName).(*mocks.DbBindingMock)
		require.True(t, ok)
		parsedClientMock.On("Collection", mongo.RetainedPeginQuoteCollection).Return(retainedCollection)
		parsedClientMock.On("Collection", mongo.PeginCreationDataCollection).Return(creationDataCollection)
		quoteCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 4}, nil).Once()
		creationDataCollection.EXPECT().DeleteMany(mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 1}, nil).Once()
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		count, err := repo.DeleteQuotes(context.Background(), hashes)
		quoteCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
		creationDataCollection.AssertExpectations(t)
		require.ErrorContains(t, err, "pegin quote collections didn't match")
		assert.Zero(t, count)
	})
}

func TestPeginMongoRepository_GetPeginCreationData(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	t.Run("read pegin creation data properly", func(t *testing.T) {
		const (
			expectedLog = "READ interaction with db: {GasPrice:55 FeePercentage:1.5 FixedFee:100000}"
			hash        = "8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819"
		)
		client, collection := getClientAndCollectionMocks(mongo.PeginCreationDataCollection)
		collection.EXPECT().FindOne(mock.Anything, bson.D{primitive.E{Key: "hash", Value: hash}}).
			Return(mongoDb.NewSingleResultFromDocument(testPeginCreationData, nil, nil)).Once()
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		defer assertDbInteractionLog(t, expectedLog)()
		result := repo.GetPeginCreationData(context.Background(), hash)
		collection.AssertExpectations(t)
		assert.Equal(t, testPeginCreationData, result)
	})
	t.Run("return zero value on invalid hash", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PeginCreationDataCollection)
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		result := repo.GetPeginCreationData(context.Background(), test.AnyString)
		collection.AssertNotCalled(t, "FindOne")
		assert.Equal(t, quote.PeginCreationDataZeroValue(), result)
	})
	t.Run("return zero value on db error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PeginCreationDataCollection)
		collection.EXPECT().FindOne(mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(nil, nil, nil)).Once()
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		result := repo.GetPeginCreationData(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		assert.Equal(t, quote.PeginCreationDataZeroValue(), result)
	})
}

func TestPeginMongoRepository_GetQuotes(t *testing.T) {
	t.Run("Get quotes with hash filters and timestamp filters", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		peginCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(peginCollection)

		hashList := []string{"27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f"}
		expectedQuotes := []quote.PeginQuote{testPeginQuote}
		peginCollection.On("Find", mock.Anything, mock.MatchedBy(func(filter bson.M) bool {
			return true
		}), mock.Anything).Return(mongoDb.NewCursorFromDocuments([]any{testPeginQuote}, nil, nil))
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)

		startDateTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		endDateTime := time.Date(2025, 1, 1, 23, 59, 59, 0, time.UTC)

		result, err := repo.GetQuotesByHashesAndDate(context.Background(), hashList, startDateTime, endDateTime)

		require.NoError(t, err)
		assert.Equal(t, expectedQuotes, result)

		peginCollection.AssertExpectations(t)
		peginCollection.AssertExpectations(t)
	})

	t.Run("error reading quotes from DB", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PeginQuoteCollection)

		collection.On("Find", mock.Anything, mock.Anything).Return(nil, mongoDb.ErrNoDocuments).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)

		hashList := []string{"27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f"}
		startDateTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		endDateTime := time.Date(2025, 1, 1, 23, 59, 59, 0, time.UTC)

		quotes, err := repo.GetQuotesByHashesAndDate(context.Background(), hashList, startDateTime, endDateTime)
		require.Error(t, err)
		assert.Equal(t, "mongo: no documents in result", err.Error())
		assert.Nil(t, quotes)
	})
}

// nolint:funlen, maintidx
func TestPeginMongoRepository_ListQuotesByDateRange(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// Test data setup
	testHash1 := "8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819"
	testHash2 := "9e2cb3dc668b7fcf52f20242713578f2e1a4793e762c3b82f66c97ed775b7920"

	testStoredQuote1 := mongo.StoredPeginQuote{
		PeginQuote: testPeginQuote,
		Hash:       testHash1,
	}
	testStoredQuote2 := mongo.StoredPeginQuote{
		PeginQuote: testPeginQuote,
		Hash:       testHash2,
	}

	testRetainedQuote2 := testRetainedPeginQuote
	testRetainedQuote2.QuoteHash = testHash2
	testRetainedQuote2.State = quote.PeginStateWaitingForDeposit

	startDate := time.Date(2023, 9, 25, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 9, 26, 23, 59, 59, 0, time.UTC)

	t.Run("Successfully list quotes with pagination and retained quotes", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		peginCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(peginCollection).Times(1)
		db.EXPECT().Collection(mongo.RetainedPeginQuoteCollection).Return(retainedCollection).Times(1)

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		peginCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(mongoDb.NewCursorFromDocuments([]any{testStoredQuote1, testStoredQuote2}, nil, nil)).Once()

		retainedFilter := bson.D{{Key: "quote_hash", Value: bson.D{{Key: "$in", Value: []string{testHash1, testHash2}}}}}
		retainedCollection.On("Find", mock.Anything, retainedFilter).
			Return(mongoDb.NewCursorFromDocuments([]any{testRetainedPeginQuote, testRetainedQuote2}, nil, nil)).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)

		defer assertDbInteractionLog(t, "READ interaction with db: 2")()

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, 2, count)
		require.Len(t, result, 2)

		// Verify first quote
		assert.Equal(t, testPeginQuote, result[0].Quote)
		assert.Equal(t, testRetainedPeginQuote, result[0].RetainedQuote)

		// Verify second quote
		assert.Equal(t, testPeginQuote, result[1].Quote)
		assert.Equal(t, testRetainedQuote2, result[1].RetainedQuote)

		peginCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
	})

	t.Run("Successfully list quotes without pagination", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		peginCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(peginCollection).Times(1)
		db.EXPECT().Collection(mongo.RetainedPeginQuoteCollection).Return(retainedCollection).Times(1)

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		// When page=0 and perPage=0, no pagination should be applied
		peginCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(mongoDb.NewCursorFromDocuments([]any{testStoredQuote1}, nil, nil)).Once()

		retainedFilter := bson.D{{Key: "quote_hash", Value: bson.D{{Key: "$in", Value: []string{testHash1}}}}}
		retainedCollection.On("Find", mock.Anything, retainedFilter).
			Return(mongoDb.NewCursorFromDocuments([]any{testRetainedPeginQuote}, nil, nil)).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)

		defer assertDbInteractionLog(t, "READ interaction with db: 1")()

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 0, 0)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
		require.Len(t, result, 1)
		assert.Equal(t, testPeginQuote, result[0].Quote)
		assert.Equal(t, testRetainedPeginQuote, result[0].RetainedQuote)

		peginCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
	})

	t.Run("Successfully return empty result", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		peginCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(peginCollection).Times(1)

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		peginCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(mongoDb.NewCursorFromDocuments([]any{}, nil, nil)).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)

		defer assertDbInteractionLog(t, "READ interaction with db: []")()

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, 0, count)
		assert.Empty(t, result)

		peginCollection.AssertExpectations(t)
	})

	t.Run("Successfully list quotes without retained quotes", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		peginCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(peginCollection).Times(1)
		db.EXPECT().Collection(mongo.RetainedPeginQuoteCollection).Return(retainedCollection).Times(1)

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		peginCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(mongoDb.NewCursorFromDocuments([]any{testStoredQuote1}, nil, nil)).Once()

		retainedFilter := bson.D{{Key: "quote_hash", Value: bson.D{{Key: "$in", Value: []string{testHash1}}}}}
		// Return empty cursor for retained quotes
		retainedCollection.On("Find", mock.Anything, retainedFilter).
			Return(mongoDb.NewCursorFromDocuments([]any{}, nil, nil)).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)

		defer assertDbInteractionLog(t, "READ interaction with db: 1")()

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
		require.Len(t, result, 1)
		assert.Equal(t, testPeginQuote, result[0].Quote)
		// RetainedQuote should be empty struct since no retained quote was found
		assert.Equal(t, quote.RetainedPeginQuote{}, result[0].RetainedQuote)

		peginCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
	})

	t.Run("Error when fetching quotes from database", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		peginCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(peginCollection).Times(1)

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		peginCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(nil, assert.AnError).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 1, 10)

		require.Error(t, err)
		assert.Equal(t, 0, count)
		assert.Nil(t, result)

		peginCollection.AssertExpectations(t)
	})

	t.Run("Error when fetching retained quotes", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		peginCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(peginCollection).Times(1)
		db.EXPECT().Collection(mongo.RetainedPeginQuoteCollection).Return(retainedCollection).Times(1)

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		peginCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(mongoDb.NewCursorFromDocuments([]any{testStoredQuote1}, nil, nil)).Once()

		retainedFilter := bson.D{{Key: "quote_hash", Value: bson.D{{Key: "$in", Value: []string{testHash1}}}}}
		retainedCollection.On("Find", mock.Anything, retainedFilter).
			Return(nil, assert.AnError).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 1, 10)

		require.Error(t, err)
		assert.Equal(t, 1, count) // Should still return the count from quotes even if retained quotes fail
		require.Len(t, result, 1) // Should return the partial result

		peginCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
	})

	t.Run("Successfully handle pagination edge cases", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		peginCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(peginCollection).Times(1)
		db.EXPECT().Collection(mongo.RetainedPeginQuoteCollection).Return(retainedCollection).Times(1)

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		// Test page 2 with perPage 1 (should skip 1 and limit 1)
		peginCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(mongoDb.NewCursorFromDocuments([]any{testStoredQuote2}, nil, nil)).Once()

		retainedFilter := bson.D{{Key: "quote_hash", Value: bson.D{{Key: "$in", Value: []string{testHash2}}}}}
		retainedCollection.On("Find", mock.Anything, retainedFilter).
			Return(mongoDb.NewCursorFromDocuments([]any{}, nil, nil)).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)

		defer assertDbInteractionLog(t, "READ interaction with db: 1")()

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 2, 1)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
		require.Len(t, result, 1)

		peginCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
	})

	t.Run("Should fill zero values for retained pegin quotes with missing gas fields", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		peginCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(peginCollection).Times(1)
		db.EXPECT().Collection(mongo.RetainedPeginQuoteCollection).Return(retainedCollection).Times(1)

		// Create old database record with missing gas fields (represented as BSON document)
		oldRetainedDocument := bson.D{
			{Key: "quote_hash", Value: testHash1},
			{Key: "deposit_address", Value: "test_deposit_address"},
			{Key: "signature", Value: "test_signature"},
			{Key: "required_liquidity", Value: uint64(1000000)},
			{Key: "state", Value: "WaitingForDeposit"},
			{Key: "call_for_user_gas_used", Value: uint64(21000)},
			{Key: "register_pegin_gas_used", Value: uint64(22000)},
			{Key: "owner_account_address", Value: "0x123"},
			// Note: CallForUserGasPrice and RegisterPeginGasPrice are missing (nil)
		}

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		peginCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(mongoDb.NewCursorFromDocuments([]any{testStoredQuote1}, nil, nil)).Once()

		retainedFilter := bson.D{{Key: "quote_hash", Value: bson.D{{Key: "$in", Value: []string{testHash1}}}}}
		retainedCollection.On("Find", mock.Anything, retainedFilter).
			Return(mongoDb.NewCursorFromDocuments([]any{oldRetainedDocument}, nil, nil)).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPeginMongoRepository(conn)

		defer assertDbInteractionLog(t, "READ interaction with db: 1")()

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
		require.Len(t, result, 1)
		assert.Equal(t, testPeginQuote, result[0].Quote)

		// Verify that FillZeroValues() was applied - gas prices should be zero Wei instead of nil
		assert.NotNil(t, result[0].RetainedQuote.CallForUserGasPrice)
		assert.NotNil(t, result[0].RetainedQuote.RegisterPeginGasPrice)
		assert.Equal(t, entities.NewWei(0), result[0].RetainedQuote.CallForUserGasPrice)
		assert.Equal(t, entities.NewWei(0), result[0].RetainedQuote.RegisterPeginGasPrice)

		peginCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
	})
}

// nolint:funlen
func TestPeginMongoRepository_GetRetainedQuotesForAddress(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	const address = "0xAA9cAf1e3967600578727F975F283446A3Da6612"

	t.Run("Get retained pegin quotes for address with specific state", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		const expectedLog = "READ interaction with db: [{QuoteHash:8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819 DepositAddress:2N7Vw5f59V3o3bDcaJK5oA829LFTBYZHLoG Signature:b24831aac7230910087d9818b378a31679be5e3991a7227cc160bc3add09e1645a26e9c740e3467f53953d7ec086c82bf8ef0eb03c118d0382ee6049a8f0119f1c RequiredLiquidity:100 State:WaitingForDeposit UserBtcTxHash:619c4d69ccaa5f78aaa2284817cf070609ac40af3792916ca3d0ef82b14af75f CallForUserTxHash:0x2c73de184c80797c04a655217d121588e8d5c228d3e0cc26187cb249123aa7c3 RegisterPeginTxHash:0x3a0feaef4d803468ba5bfc1db78f4d2568de1b7cf002dec5991c469e6719db89 CallForUserGasUsed:85000 CallForUserGasPrice:21000000000 RegisterPeginGasUsed:65000 RegisterPeginGasPrice:21000000000 OwnerAccountAddress:0xAA9cAf1e3967600578727F975F283446A3Da6612}]"

		mockQuote := testRetainedPeginQuote
		mockQuote.State = quote.PeginStateWaitingForDeposit
		mockQuote.OwnerAccountAddress = address

		collection.On("Find", mock.Anything, mock.MatchedBy(func(filter bson.D) bool {
			// Assert that the filter structure matches what we expect
			assert.Len(t, filter, 2)
			assert.Equal(t, "owner_account_address", filter[0].Key)
			assert.Equal(t, address, filter[0].Value)
			assert.Equal(t, "state", filter[1].Key)
			stateFilter, ok := filter[1].Value.(bson.D)
			assert.True(t, ok)
			assert.Len(t, stateFilter, 1)
			assert.Equal(t, "$in", stateFilter[0].Key)
			stateValues, ok := stateFilter[0].Value.([]quote.PeginState)
			assert.True(t, ok)
			assert.Len(t, stateValues, 1)
			assert.Contains(t, stateValues, quote.PeginStateWaitingForDeposit)

			return true
		})).Return(mongoDb.NewCursorFromDocuments([]any{mockQuote}, nil, nil)).Once()

		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetRetainedQuotesForAddress(context.Background(), address, quote.PeginStateWaitingForDeposit)

		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mockQuote, result[0])
	})

	t.Run("Get retained pegin quotes for address with multiple specific states", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		const expectedLog = "READ interaction with db: [{QuoteHash:8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819 DepositAddress:2N7Vw5f59V3o3bDcaJK5oA829LFTBYZHLoG Signature:b24831aac7230910087d9818b378a31679be5e3991a7227cc160bc3add09e1645a26e9c740e3467f53953d7ec086c82bf8ef0eb03c118d0382ee6049a8f0119f1c RequiredLiquidity:100 State:CallForUserSucceeded UserBtcTxHash:619c4d69ccaa5f78aaa2284817cf070609ac40af3792916ca3d0ef82b14af75f CallForUserTxHash:0x2c73de184c80797c04a655217d121588e8d5c228d3e0cc26187cb249123aa7c3 RegisterPeginTxHash:0x3a0feaef4d803468ba5bfc1db78f4d2568de1b7cf002dec5991c469e6719db89 CallForUserGasUsed:85000 CallForUserGasPrice:21000000000 RegisterPeginGasUsed:65000 RegisterPeginGasPrice:21000000000 OwnerAccountAddress:0xAA9cAf1e3967600578727F975F283446A3Da6612} {QuoteHash:second DepositAddress:2N7Vw5f59V3o3bDcaJK5oA829LFTBYZHLoG Signature:123 RequiredLiquidity:777 State:RegisterPegInSucceeded UserBtcTxHash:619c4d69ccaa5f78aaa2284817cf070609ac40af3792916ca3d0ef82b14af75f CallForUserTxHash:0x2c73de184c80797c04a655217d121588e8d5c228d3e0cc26187cb249123aa7c3 RegisterPeginTxHash:0x3a0feaef4d803468ba5bfc1db78f4d2568de1b7cf002dec5991c469e6719db89 CallForUserGasUsed:85000 CallForUserGasPrice:21000000000 RegisterPeginGasUsed:65000 RegisterPeginGasPrice:21000000000 OwnerAccountAddress:0xAA9cAf1e3967600578727F975F283446A3Da6612}]"

		firstQuote := testRetainedPeginQuote
		firstQuote.State = quote.PeginStateCallForUserSucceeded
		firstQuote.OwnerAccountAddress = address

		secondQuote := testRetainedPeginQuote
		secondQuote.QuoteHash = "second"
		secondQuote.Signature = "123"
		secondQuote.RequiredLiquidity = entities.NewWei(777)
		secondQuote.State = quote.PeginStateRegisterPegInSucceeded
		secondQuote.OwnerAccountAddress = address

		collection.On("Find", mock.Anything, mock.MatchedBy(func(filter bson.D) bool {
			// Assert that the filter structure matches what we expect
			assert.Len(t, filter, 2)
			assert.Equal(t, "owner_account_address", filter[0].Key)
			assert.Equal(t, address, filter[0].Value)
			assert.Equal(t, "state", filter[1].Key)
			stateFilter, ok := filter[1].Value.(bson.D)
			assert.True(t, ok)
			assert.Len(t, stateFilter, 1)
			assert.Equal(t, "$in", stateFilter[0].Key)
			stateValues, ok := stateFilter[0].Value.([]quote.PeginState)
			assert.True(t, ok)
			assert.Len(t, stateValues, 2)
			assert.Contains(t, stateValues, quote.PeginStateCallForUserSucceeded)
			assert.Contains(t, stateValues, quote.PeginStateRegisterPegInSucceeded)

			return true
		})).Return(mongoDb.NewCursorFromDocuments([]any{firstQuote, secondQuote}, nil, nil)).Once()

		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetRetainedQuotesForAddress(context.Background(), address, quote.PeginStateCallForUserSucceeded, quote.PeginStateRegisterPegInSucceeded)

		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, firstQuote, result[0])
		assert.Equal(t, secondQuote, result[1])
	})

	t.Run("Empty result with no matching quotes", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		const expectedLog = "READ interaction with db: []"

		collection.On("Find", mock.Anything, mock.MatchedBy(func(filter bson.D) bool {
			// Assert that the filter structure matches what we expect
			assert.Len(t, filter, 2)
			assert.Equal(t, "owner_account_address", filter[0].Key)
			assert.Equal(t, address, filter[0].Value)
			assert.Equal(t, "state", filter[1].Key)
			stateFilter, ok := filter[1].Value.(bson.D)
			assert.True(t, ok)
			assert.Len(t, stateFilter, 1)
			assert.Equal(t, "$in", stateFilter[0].Key)
			stateValues, ok := stateFilter[0].Value.([]quote.PeginState)
			assert.True(t, ok)
			assert.Len(t, stateValues, 2)
			assert.Contains(t, stateValues, quote.PeginStateWaitingForDeposit)
			assert.Contains(t, stateValues, quote.PeginStateWaitingForDepositConfirmations)

			return true
		})).Return(mongoDb.NewCursorFromDocuments([]any{}, nil, nil)).Once()

		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetRetainedQuotesForAddress(context.Background(), address, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations)

		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Db error when getting retained pegin quotes for address", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))

		collection.On("Find", mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()

		result, err := repo.GetRetainedQuotesForAddress(context.Background(), address)

		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("FillZeroValues is applied to retained pegin quotes with missing gas fields", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPeginQuoteCollection)
		repo := mongo.NewPeginMongoRepository(mongo.NewConnection(client, time.Duration(1)))

		// Mock strategy for similar tests did not work for mock limitations on unmarshalling into a struct
		firstOldDocument := bson.D{
			{Key: "quote_hash", Value: "address_first"},
			{Key: "deposit_address", Value: testRetainedPeginQuote.DepositAddress},
			{Key: "signature", Value: testRetainedPeginQuote.Signature},
			{Key: "required_liquidity", Value: testRetainedPeginQuote.RequiredLiquidity.String()},
			{Key: "state", Value: quote.PeginStateWaitingForDeposit},
			{Key: "user_btc_tx_hash", Value: testRetainedPeginQuote.UserBtcTxHash},
			{Key: "call_for_user_tx_hash", Value: testRetainedPeginQuote.CallForUserTxHash},
			{Key: "register_pegin_tx_hash", Value: testRetainedPeginQuote.RegisterPeginTxHash},
			{Key: "call_for_user_gas_used", Value: testRetainedPeginQuote.CallForUserGasUsed},
			{Key: "register_pegin_gas_used", Value: testRetainedPeginQuote.RegisterPeginGasUsed},
			{Key: "owner_account_address", Value: address},
			// NOTE: call_for_user_gas_price and register_pegin_gas_price are MISSING
		}

		secondOldDocument := bson.D{
			{Key: "quote_hash", Value: "address_second"},
			{Key: "deposit_address", Value: testRetainedPeginQuote.DepositAddress},
			{Key: "signature", Value: "address_signature"},
			{Key: "required_liquidity", Value: entities.NewWei(300).String()},
			{Key: "state", Value: quote.PeginStateWaitingForDepositConfirmations},
			{Key: "user_btc_tx_hash", Value: testRetainedPeginQuote.UserBtcTxHash},
			{Key: "call_for_user_tx_hash", Value: testRetainedPeginQuote.CallForUserTxHash},
			{Key: "register_pegin_tx_hash", Value: testRetainedPeginQuote.RegisterPeginTxHash},
			{Key: "call_for_user_gas_used", Value: uint64(90000)},
			{Key: "register_pegin_gas_used", Value: uint64(70000)},
			{Key: "owner_account_address", Value: address},
			// NOTE: call_for_user_gas_price and register_pegin_gas_price are MISSING
		}

		cursor, err := mongoDb.NewCursorFromDocuments([]any{firstOldDocument, secondOldDocument}, nil, nil)
		require.NoError(t, err)

		collection.On("Find", mock.Anything, mock.Anything).
			Return(cursor, nil).Once()

		result, err := repo.GetRetainedQuotesForAddress(context.Background(), address, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations)

		collection.AssertExpectations(t)
		require.NoError(t, err)
		require.Len(t, result, 2)

		// Verify normalization applied to first document
		assert.NotNil(t, result[0].CallForUserGasPrice, "CallForUserGasPrice should not be nil after normalization")
		assert.NotNil(t, result[0].RegisterPeginGasPrice, "RegisterPeginGasPrice should not be nil after normalization")
		assert.Equal(t, entities.NewWei(0), result[0].CallForUserGasPrice)
		assert.Equal(t, entities.NewWei(0), result[0].RegisterPeginGasPrice)
		assert.Equal(t, "address_first", result[0].QuoteHash)
		assert.Equal(t, testRetainedPeginQuote.CallForUserGasUsed, result[0].CallForUserGasUsed)

		// Verify normalization applied to second document
		assert.NotNil(t, result[1].CallForUserGasPrice, "CallForUserGasPrice should not be nil after normalization")
		assert.NotNil(t, result[1].RegisterPeginGasPrice, "RegisterPeginGasPrice should not be nil after normalization")
		assert.Equal(t, entities.NewWei(0), result[1].CallForUserGasPrice)
		assert.Equal(t, entities.NewWei(0), result[1].RegisterPeginGasPrice)
		assert.Equal(t, "address_second", result[1].QuoteHash)
		assert.Equal(t, uint64(90000), result[1].CallForUserGasUsed)
		assert.Equal(t, uint64(70000), result[1].RegisterPeginGasUsed)
	})
}

// nolint: cyclop, funlen, gocognit, gocyclo
func validatePeginPipelineStructure(pipeline mongoDb.Pipeline, states []quote.PeginState, startDate, endDate time.Time) bool {
	// Verify the pipeline structure
	if len(pipeline) != 5 {
		return false
	}

	// Stage 1: $match by date - verify dates are correct
	matchStage := pipeline[0]
	if len(matchStage) == 0 || matchStage[0].Key != "$match" {
		return false
	}
	matchDateFilter, ok := matchStage[0].Value.(bson.M)
	if !ok {
		return false
	}
	timestampFilter, ok := matchDateFilter["agreement_timestamp"].(bson.M)
	if !ok {
		return false
	}
	gteValue, hasGte := timestampFilter["$gte"]
	lteValue, hasLte := timestampFilter["$lte"]
	if !hasGte || !hasLte {
		return false
	}
	// Verify the dates match
	if gteValue != startDate.Unix() || lteValue != endDate.Unix() {
		return false
	}

	// Stage 2: $lookup
	lookupStage := pipeline[1]
	if len(lookupStage) == 0 || lookupStage[0].Key != "$lookup" {
		return false
	}

	// Stage 3: $unwind
	unwindStage := pipeline[2]
	if len(unwindStage) == 0 || unwindStage[0].Key != "$unwind" {
		return false
	}
	unwindConfig, ok := unwindStage[0].Value.(bson.M)
	if !ok {
		return false
	}
	preserveNullAndEmptyArrays, ok := unwindConfig["preserveNullAndEmptyArrays"].(bool)
	if !ok || !preserveNullAndEmptyArrays {
		return false
	}

	// Stage 4: $match by state with $or to include non-accepted quotes
	matchStateStage := pipeline[3]
	if len(matchStateStage) == 0 || matchStateStage[0].Key != "$match" {
		return false
	}
	matchStateFilter, ok := matchStateStage[0].Value.(bson.M)
	if !ok {
		return false
	}
	orConditions, ok := matchStateFilter["$or"].([]bson.M)
	if !ok || len(orConditions) != 2 {
		return false
	}
	// First condition: retained.state in states
	stateFilter, ok := orConditions[0]["retained.state"].(bson.M)
	if !ok {
		return false
	}
	inStates, ok := stateFilter["$in"].([]quote.PeginState)
	if !ok {
		return false
	}
	// Verify the states match (must have same length and contain same elements)
	if len(inStates) != len(states) {
		return false
	}
	stateMap := make(map[quote.PeginState]bool)
	for _, s := range states {
		stateMap[s] = true
	}
	for _, s := range inStates {
		if !stateMap[s] {
			return false
		}
	}
	// Second condition: retained.state is nil (for non-accepted quotes)
	if orConditions[1]["retained.state"] != nil {
		return false
	}

	// Stage 5: $limit - verify limit is set to 501
	limitStage := pipeline[4]
	if len(limitStage) == 0 || limitStage[0].Key != "$limit" {
		return false
	}
	limitValue, ok := limitStage[0].Value.(int)
	if !ok || limitValue != 501 {
		return false
	}

	return true
}

func buildPeginAggregationDocument(q quote.PeginQuote, retained quote.RetainedPeginQuote) bson.D {
	return bson.D{
		{Key: "fed_btc_address", Value: q.FedBtcAddress},
		{Key: "lbc_address", Value: q.LbcAddress},
		{Key: "lp_rsk_address", Value: q.LpRskAddress},
		{Key: "btc_refund_address", Value: q.BtcRefundAddress},
		{Key: "rsk_refund_address", Value: q.RskRefundAddress},
		{Key: "lp_btc_address", Value: q.LpBtcAddress},
		{Key: "call_fee", Value: q.CallFee.String()},
		{Key: "penalty_fee", Value: q.PenaltyFee.String()},
		{Key: "contract_address", Value: q.ContractAddress},
		{Key: "data", Value: q.Data},
		{Key: "gas_limit", Value: q.GasLimit},
		{Key: "nonce", Value: q.Nonce},
		{Key: "value", Value: q.Value.String()},
		{Key: "agreement_timestamp", Value: q.AgreementTimestamp},
		{Key: "time_for_deposit", Value: q.TimeForDeposit},
		{Key: "lp_call_time", Value: q.LpCallTime},
		{Key: "confirmations", Value: q.Confirmations},
		{Key: "call_on_register", Value: q.CallOnRegister},
		{Key: "gas_fee", Value: q.GasFee.String()},
		{Key: "hash", Value: retained.QuoteHash},
		{Key: "retained", Value: bson.D{
			{Key: "quote_hash", Value: retained.QuoteHash},
			{Key: "deposit_address", Value: retained.DepositAddress},
			{Key: "signature", Value: retained.Signature},
			{Key: "required_liquidity", Value: retained.RequiredLiquidity.String()},
			{Key: "state", Value: retained.State},
			{Key: "user_btc_tx_hash", Value: retained.UserBtcTxHash},
			{Key: "call_for_user_tx_hash", Value: retained.CallForUserTxHash},
			{Key: "register_pegin_tx_hash", Value: retained.RegisterPeginTxHash},
			{Key: "call_for_user_gas_used", Value: retained.CallForUserGasUsed},
			{Key: "call_for_user_gas_price", Value: retained.CallForUserGasPrice.String()},
			{Key: "register_pegin_gas_used", Value: retained.RegisterPeginGasUsed},
			{Key: "register_pegin_gas_price", Value: retained.RegisterPeginGasPrice.String()},
			{Key: "owner_account_address", Value: retained.OwnerAccountAddress},
		}},
	}
}

// nolint: funlen
func TestPeginMongoRepository_GetQuotesWithRetainedByStateAndDate_Success(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	states := []quote.PeginState{quote.PeginStateCallForUserSucceeded, quote.PeginStateRegisterPegInSucceeded, quote.PeginStateCallForUserFailed}

	client, db := getClientAndDatabaseMocks()
	collection := &mocks.CollectionBindingMock{}
	db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(collection)

	// Create test documents that simulate aggregation pipeline results
	firstQuote := testPeginQuote
	firstQuote.AgreementTimestamp = uint32(startDate.Add(24 * time.Hour).Unix())
	firstRetained := testRetainedPeginQuote
	firstRetained.QuoteHash = "test-hash-1"
	firstRetained.State = quote.PeginStateCallForUserSucceeded

	secondQuote := testPeginQuote
	secondQuote.AgreementTimestamp = uint32(startDate.Add(48 * time.Hour).Unix())
	secondRetained := testRetainedPeginQuote
	secondRetained.QuoteHash = "test-hash-2"
	secondRetained.State = quote.PeginStateRegisterPegInSucceeded

	thirdQuote := testPeginQuote
	thirdQuote.AgreementTimestamp = uint32(startDate.Add(72 * time.Hour).Unix())
	thirdRetained := testRetainedPeginQuote
	thirdRetained.QuoteHash = "test-hash-3"
	thirdRetained.State = quote.PeginStateCallForUserFailed

	// Simulate aggregation result documents
	doc1 := buildPeginAggregationDocument(firstQuote, firstRetained)
	doc2 := buildPeginAggregationDocument(secondQuote, secondRetained)
	doc3 := buildPeginAggregationDocument(thirdQuote, thirdRetained)

	cursor, err := mongoDb.NewCursorFromDocuments([]any{doc1, doc2, doc3}, nil, nil)
	require.NoError(t, err)

	collection.On("Aggregate", mock.Anything, mock.MatchedBy(func(pipeline mongoDb.Pipeline) bool {
		return validatePeginPipelineStructure(pipeline, states, startDate, endDate)
	})).Return(cursor, nil).Once()

	conn := mongo.NewConnection(client, time.Duration(1))
	repo := mongo.NewPeginMongoRepository(conn)

	result, err := repo.GetQuotesWithRetainedByStateAndDate(context.Background(), states, startDate, endDate)

	collection.AssertExpectations(t)
	require.NoError(t, err)
	require.Len(t, result, 3)

	assert.Equal(t, firstRetained.QuoteHash, result[0].RetainedQuote.QuoteHash)
	assert.Equal(t, firstRetained.State, result[0].RetainedQuote.State)
	assert.Equal(t, firstQuote.CallFee, result[0].Quote.CallFee)

	assert.Equal(t, secondRetained.QuoteHash, result[1].RetainedQuote.QuoteHash)
	assert.Equal(t, secondRetained.State, result[1].RetainedQuote.State)
	assert.Equal(t, secondQuote.CallFee, result[1].Quote.CallFee)

	assert.Equal(t, thirdRetained.QuoteHash, result[2].RetainedQuote.QuoteHash)
	assert.Equal(t, thirdRetained.State, result[2].RetainedQuote.State)
	assert.Equal(t, thirdQuote.CallFee, result[2].Quote.CallFee)
}

func TestPeginMongoRepository_GetQuotesWithRetainedByStateAndDate_EmptyResult(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	states := []quote.PeginState{quote.PeginStateCallForUserSucceeded}

	client, db := getClientAndDatabaseMocks()
	collection := &mocks.CollectionBindingMock{}
	db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(collection)

	cursor, err := mongoDb.NewCursorFromDocuments([]any{}, nil, nil)
	require.NoError(t, err)

	collection.On("Aggregate", mock.Anything, mock.Anything).Return(cursor, nil).Once()

	conn := mongo.NewConnection(client, time.Duration(1))
	repo := mongo.NewPeginMongoRepository(conn)

	result, err := repo.GetQuotesWithRetainedByStateAndDate(context.Background(), states, startDate, endDate)

	collection.AssertExpectations(t)
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestPeginMongoRepository_GetQuotesWithRetainedByStateAndDate_AggregationError(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	states := []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}

	client, db := getClientAndDatabaseMocks()
	collection := &mocks.CollectionBindingMock{}
	db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(collection)

	expectedError := assert.AnError
	collection.On("Aggregate", mock.Anything, mock.Anything).Return(nil, expectedError).Once()

	conn := mongo.NewConnection(client, time.Duration(1))
	repo := mongo.NewPeginMongoRepository(conn)

	result, err := repo.GetQuotesWithRetainedByStateAndDate(context.Background(), states, startDate, endDate)

	collection.AssertExpectations(t)
	require.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
}

func TestPeginMongoRepository_GetQuotesWithRetainedByStateAndDate_DecodeError(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	states := []quote.PeginState{quote.PeginStateCallForUserSucceeded}

	client, db := getClientAndDatabaseMocks()
	collection := &mocks.CollectionBindingMock{}
	db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(collection)

	// Create a document with invalid data that will fail decode
	invalidDoc := bson.D{
		{Key: "call_fee", Value: "invalid"}, // Invalid type for Wei
	}

	cursor, err := mongoDb.NewCursorFromDocuments([]any{invalidDoc}, nil, nil)
	require.NoError(t, err)

	collection.On("Aggregate", mock.Anything, mock.Anything).Return(cursor, nil).Once()

	conn := mongo.NewConnection(client, time.Duration(1))
	repo := mongo.NewPeginMongoRepository(conn)

	result, err := repo.GetQuotesWithRetainedByStateAndDate(context.Background(), states, startDate, endDate)

	collection.AssertExpectations(t)
	require.Error(t, err)
	assert.Nil(t, result)
}

// nolint: funlen
func TestPeginMongoRepository_GetQuotesWithRetainedByStateAndDate_ZeroValueNormalization(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	states := []quote.PeginState{quote.PeginStateCallForUserSucceeded}

	client, db := getClientAndDatabaseMocks()
	collection := &mocks.CollectionBindingMock{}
	db.EXPECT().Collection(mongo.PeginQuoteCollection).Return(collection)

	testQuote := testPeginQuote
	testQuote.AgreementTimestamp = uint32(startDate.Add(24 * time.Hour).Unix())

	// Create document with missing gas price fields
	doc := bson.D{
		{Key: "fed_btc_address", Value: testQuote.FedBtcAddress},
		{Key: "lbc_address", Value: testQuote.LbcAddress},
		{Key: "lp_rsk_address", Value: testQuote.LpRskAddress},
		{Key: "btc_refund_address", Value: testQuote.BtcRefundAddress},
		{Key: "rsk_refund_address", Value: testQuote.RskRefundAddress},
		{Key: "lp_btc_address", Value: testQuote.LpBtcAddress},
		{Key: "call_fee", Value: testQuote.CallFee.String()},
		{Key: "penalty_fee", Value: testQuote.PenaltyFee.String()},
		{Key: "contract_address", Value: testQuote.ContractAddress},
		{Key: "data", Value: testQuote.Data},
		{Key: "gas_limit", Value: testQuote.GasLimit},
		{Key: "nonce", Value: testQuote.Nonce},
		{Key: "value", Value: testQuote.Value.String()},
		{Key: "agreement_timestamp", Value: testQuote.AgreementTimestamp},
		{Key: "time_for_deposit", Value: testQuote.TimeForDeposit},
		{Key: "lp_call_time", Value: testQuote.LpCallTime},
		{Key: "confirmations", Value: testQuote.Confirmations},
		{Key: "call_on_register", Value: testQuote.CallOnRegister},
		{Key: "gas_fee", Value: testQuote.GasFee.String()},
		{Key: "hash", Value: "test-hash"},
		{Key: "retained", Value: bson.D{
			{Key: "quote_hash", Value: "test-hash"},
			{Key: "deposit_address", Value: testRetainedPeginQuote.DepositAddress},
			{Key: "signature", Value: testRetainedPeginQuote.Signature},
			{Key: "required_liquidity", Value: testRetainedPeginQuote.RequiredLiquidity.String()},
			{Key: "state", Value: quote.PeginStateCallForUserSucceeded},
			{Key: "user_btc_tx_hash", Value: testRetainedPeginQuote.UserBtcTxHash},
			{Key: "call_for_user_tx_hash", Value: testRetainedPeginQuote.CallForUserTxHash},
			{Key: "register_pegin_tx_hash", Value: testRetainedPeginQuote.RegisterPeginTxHash},
			{Key: "call_for_user_gas_used", Value: uint64(85000)},
			{Key: "register_pegin_gas_used", Value: uint64(65000)},
			{Key: "owner_account_address", Value: testRetainedPeginQuote.OwnerAccountAddress},
			// NOTE: call_for_user_gas_price and register_pegin_gas_price are MISSING
		}},
	}

	cursor, err := mongoDb.NewCursorFromDocuments([]any{doc}, nil, nil)
	require.NoError(t, err)

	collection.On("Aggregate", mock.Anything, mock.Anything).Return(cursor, nil).Once()

	conn := mongo.NewConnection(client, time.Duration(1))
	repo := mongo.NewPeginMongoRepository(conn)

	result, err := repo.GetQuotesWithRetainedByStateAndDate(context.Background(), states, startDate, endDate)

	collection.AssertExpectations(t)
	require.NoError(t, err)
	require.Len(t, result, 1)

	// Verify normalization applied
	assert.NotNil(t, result[0].RetainedQuote.CallForUserGasPrice, "CallForUserGasPrice should not be nil after normalization")
	assert.NotNil(t, result[0].RetainedQuote.RegisterPeginGasPrice, "RegisterPeginGasPrice should not be nil after normalization")
	assert.Equal(t, entities.NewWei(0), result[0].RetainedQuote.CallForUserGasPrice)
	assert.Equal(t, entities.NewWei(0), result[0].RetainedQuote.RegisterPeginGasPrice)
}
