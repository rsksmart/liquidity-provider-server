package mongo_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"

	"github.com/btcsuite/btcd/btcjson"
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
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testPegoutQuote = quote.PegoutQuote{
	LbcAddress:            "0xc2A630c053D12D63d32b025082f6Ba268db18300",
	LpRskAddress:          "0x7c4890a0f1d4bbf2c669ac2d1effa185c505359b",
	BtcRefundAddress:      "n2Ge4xMVQKp5Hzzf8xTBJBLppRgjRZYYyq",
	RskRefundAddress:      "0x79568C2989232dcA1840087d73d403602364c0D4",
	LpBtcAddress:          "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe",
	CallFee:               entities.NewWei(100000000000000),
	PenaltyFee:            entities.NewWei(10000000000000),
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
	ProductFeeAmount:      entities.NewWei(13),
}

var testRetainedPegoutQuote = quote.RetainedPegoutQuote{
	QuoteHash:            "27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f",
	DepositAddress:       "mkE1WWdiu5VgjfugomDk8GxV6JdEEEJR9s",
	Signature:            "5c9eab91c753355f87c19d09ea88b2fd02773981e513bc2821fed5ceba0d452a0a3d21e2252cb35348ce5c6803117e3abb62837beb8f5866a375ce66587d004b1c",
	RequiredLiquidity:    entities.NewWei(55),
	State:                quote.PegoutStateWaitingForDepositConfirmations,
	UserRskTxHash:        "0x6b2e1e4daf8cf00c5c3534b72cdeec3526e8a38f70c11e44888b6e4ae1ee7d38",
	LpBtcTxHash:          "6ac3779dc33ad52f3409cbb909bcd458745995496a2a3954406206f6e5d4cb0e",
	RefundPegoutTxHash:   "0x8e773a2826e73f8e5792304379a7e46dff38f17089c6d344335e03537b31c2bc",
	BridgeRefundTxHash:   "0x4f3f6f0664a732e4c907971e75c1e3fd8671461dcb53f566660432fc47255d8b",
	BridgeRefundGasUsed:  21000,
	BridgeRefundGasPrice: entities.NewWei(20000000000),
	RefundPegoutGasUsed:  22000,
	RefundPegoutGasPrice: entities.NewWei(25000000000),
	SendPegoutBtcFee:     entities.NewWei(15000),
	BtcReleaseTxHash:     "0xd8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb",
	OwnerAccountAddress:  "0x233845a26a4dA08E16218e7B401501D048670674",
}

var testPegoutDeposit = quote.PegoutDeposit{
	TxHash:      test.AnyString,
	QuoteHash:   test.AnyString,
	Amount:      entities.NewWei(999),
	Timestamp:   time.Unix(1715001146288, 0).UTC(),
	BlockNumber: 789,
	From:        test.AnyAddress,
}

var testPegoutCreationData = quote.PegoutCreationData{
	FeeRate:       utils.NewBigFloat64(1.55),
	FeePercentage: utils.NewBigFloat64(2.41),
	GasPrice:      entities.NewWei(123),
	FixedFee:      entities.NewWei(456),
}

func TestPegoutMongoRepository_InsertQuote(t *testing.T) {
	t.Run("Insert pegout quote successfully", func(t *testing.T) {
		const expectedLog = "INSERT interaction with db: {PegoutQuote:{LbcAddress:0xc2A630c053D12D63d32b025082f6Ba268db18300 LpRskAddress:0x7c4890a0f1d4bbf2c669ac2d1effa185c505359b BtcRefundAddress:n2Ge4xMVQKp5Hzzf8xTBJBLppRgjRZYYyq RskRefundAddress:0x79568C2989232dcA1840087d73d403602364c0D4 LpBtcAddress:mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe CallFee:100000000000000 PenaltyFee:10000000000000 Nonce:6410832321595034747 DepositAddress:n2Ge4xMVQKp5Hzzf8xTBJBLppRgjRZYYyq Value:5000000000000000 AgreementTimestamp:1721944367 DepositDateLimit:1721951567 DepositConfirmations:4 TransferConfirmations:2 TransferTime:7200 ExpireDate:1721958767 ExpireBlock:5366409 GasFee:4170000000000 ProductFeeAmount:13} Hash:any value}"
		client, db := getClientAndDatabaseMocks()
		quoteCollection := &mocks.CollectionBindingMock{}
		creationDataCollection := &mocks.CollectionBindingMock{}
		db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(quoteCollection)
		db.EXPECT().Collection(mongo.PegoutCreationDataCollection).Return(creationDataCollection)
		quoteCollection.On("InsertOne", mock.Anything, mock.MatchedBy(func(q mongo.StoredPegoutQuote) bool {
			return q.Hash == test.AnyString && reflect.TypeOf(quote.PegoutQuote{}).NumField() == test.CountNonZeroValues(q.PegoutQuote)
		})).Return(nil, nil).Once()
		creationDataCollection.EXPECT().InsertOne(mock.Anything, mock.MatchedBy(func(q mongo.StoredPegoutCreationData) bool {
			return q.Hash == test.AnyString && reflect.TypeOf(quote.PegoutCreationData{}).NumField() == test.CountNonZeroValues(q.PegoutCreationData)
		})).Return(nil, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)
		defer assertDbInteractionLog(t, expectedLog)()
		createdQuote := quote.CreatedPegoutQuote{Hash: test.AnyString, Quote: testPegoutQuote, CreationData: testPegoutCreationData}
		err := repo.InsertQuote(context.Background(), createdQuote)
		quoteCollection.AssertExpectations(t)
		creationDataCollection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Db error when inserting pegout quote", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		quoteCollection := &mocks.CollectionBindingMock{}
		creationDataCollection := &mocks.CollectionBindingMock{}
		db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(quoteCollection)
		db.EXPECT().Collection(mongo.PegoutCreationDataCollection).Return(creationDataCollection)
		quoteCollection.On("InsertOne", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)
		createdQuote := quote.CreatedPegoutQuote{Hash: test.AnyString, Quote: testPegoutQuote, CreationData: testPegoutCreationData}
		err := repo.InsertQuote(context.Background(), createdQuote)
		quoteCollection.AssertExpectations(t)
		creationDataCollection.AssertNotCalled(t, "InsertOne")
		require.Error(t, err)
	})
	t.Run("Db error when inserting pegout creation data", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		quoteCollection := &mocks.CollectionBindingMock{}
		creationDataCollection := &mocks.CollectionBindingMock{}
		db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(quoteCollection)
		db.EXPECT().Collection(mongo.PegoutCreationDataCollection).Return(creationDataCollection)
		quoteCollection.On("InsertOne", mock.Anything, mock.Anything).Return(nil, nil).Once()
		creationDataCollection.On("InsertOne", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)
		createdQuote := quote.CreatedPegoutQuote{
			Hash:         test.AnyString,
			Quote:        testPegoutQuote,
			CreationData: testPegoutCreationData,
		}
		err := repo.InsertQuote(context.Background(), createdQuote)
		quoteCollection.AssertExpectations(t)
		creationDataCollection.AssertExpectations(t)
		require.Error(t, err)
	})
}

func TestPegoutMongoRepository_GetQuote(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.PegoutQuoteCollection)
	log.SetLevel(log.DebugLevel)
	t.Run("Get pegout quote successfully", func(t *testing.T) {
		const expectedLog = "READ interaction with db: {LbcAddress:0xc2A630c053D12D63d32b025082f6Ba268db18300 LpRskAddress:0x7c4890a0f1d4bbf2c669ac2d1effa185c505359b BtcRefundAddress:n2Ge4xMVQKp5Hzzf8xTBJBLppRgjRZYYyq RskRefundAddress:0x79568C2989232dcA1840087d73d403602364c0D4 LpBtcAddress:mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe CallFee:100000000000000 PenaltyFee:10000000000000 Nonce:6410832321595034747 DepositAddress:n2Ge4xMVQKp5Hzzf8xTBJBLppRgjRZYYyq Value:5000000000000000 AgreementTimestamp:1721944367 DepositDateLimit:1721951567 DepositConfirmations:4 TransferConfirmations:2 TransferTime:7200 ExpireDate:1721958767 ExpireBlock:5366409 GasFee:4170000000000 ProductFeeAmount:13}"
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
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
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(mongo.StoredPegoutQuote{}, assert.AnError, nil)).Once()
		result, err := repo.GetQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Pegout quote not present in db", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(mongo.StoredPegoutQuote{}, mongoDb.ErrNoDocuments, nil)).Once()
		result, err := repo.GetQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
	t.Run("Fail on invalid pegout quote hash", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		result, err := repo.GetQuote(context.Background(), test.AnyString)
		collection.AssertNotCalled(t, "FindOne")
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

//nolint:funlen
func TestPegoutMongoRepository_GetRetainedQuote(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
	log.SetLevel(log.DebugLevel)
	t.Run("Get retained pegout quote successfully", func(t *testing.T) {
		const expectedLog = "READ interaction with db: {QuoteHash:27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f DepositAddress:mkE1WWdiu5VgjfugomDk8GxV6JdEEEJR9s Signature:5c9eab91c753355f87c19d09ea88b2fd02773981e513bc2821fed5ceba0d452a0a3d21e2252cb35348ce5c6803117e3abb62837beb8f5866a375ce66587d004b1c RequiredLiquidity:55 State:WaitingForDepositConfirmations UserRskTxHash:0x6b2e1e4daf8cf00c5c3534b72cdeec3526e8a38f70c11e44888b6e4ae1ee7d38 LpBtcTxHash:6ac3779dc33ad52f3409cbb909bcd458745995496a2a3954406206f6e5d4cb0e RefundPegoutTxHash:0x8e773a2826e73f8e5792304379a7e46dff38f17089c6d344335e03537b31c2bc BridgeRefundTxHash:0x4f3f6f0664a732e4c907971e75c1e3fd8671461dcb53f566660432fc47255d8b BridgeRefundGasUsed:21000 BridgeRefundGasPrice:20000000000 RefundPegoutGasUsed:22000 RefundPegoutGasPrice:25000000000 SendPegoutBtcFee:15000 BtcReleaseTxHash:0xd8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb OwnerAccountAddress:0x233845a26a4dA08E16218e7B401501D048670674}"
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, bson.D{primitive.E{Key: "quote_hash", Value: test.AnyHash}}).
			Return(mongoDb.NewSingleResultFromDocument(testRetainedPegoutQuote, nil, nil)).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, testRetainedPegoutQuote, *result)
	})
	t.Run("Db error when getting retained pegout quote", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(quote.RetainedPegoutQuote{}, assert.AnError, nil)).Once()
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Retained pegout quote not present in db", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(quote.RetainedPegoutQuote{}, mongoDb.ErrNoDocuments, nil)).Once()
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
	t.Run("Fail on invalid pegout quote hash", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		result, err := repo.GetRetainedQuote(context.Background(), test.AnyString)
		collection.AssertNotCalled(t, "FindOne")
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("FillZeroValues is applied to retained pegout quote with missing gas fields", func(t *testing.T) {
		// Mock strategy for similar tests did not work for mock limitations on unmarshalling into a struct
		oldBsonDocument := bson.D{
			{Key: "quote_hash", Value: testRetainedPegoutQuote.QuoteHash},
			{Key: "deposit_address", Value: testRetainedPegoutQuote.DepositAddress},
			{Key: "signature", Value: testRetainedPegoutQuote.Signature},
			{Key: "required_liquidity", Value: testRetainedPegoutQuote.RequiredLiquidity.String()},
			{Key: "state", Value: testRetainedPegoutQuote.State},
			{Key: "user_rsk_tx_hash", Value: testRetainedPegoutQuote.UserRskTxHash},
			{Key: "lp_btc_tx_hash", Value: testRetainedPegoutQuote.LpBtcTxHash},
			{Key: "refund_pegout_tx_hash", Value: testRetainedPegoutQuote.RefundPegoutTxHash},
			{Key: "bridge_refund_tx_hash", Value: testRetainedPegoutQuote.BridgeRefundTxHash},
			{Key: "bridge_refund_gas_used", Value: testRetainedPegoutQuote.BridgeRefundGasUsed},
			{Key: "refund_pegout_gas_used", Value: testRetainedPegoutQuote.RefundPegoutGasUsed},
			{Key: "btc_release_tx_hash", Value: testRetainedPegoutQuote.BtcReleaseTxHash},
			{Key: "owner_account_address", Value: testRetainedPegoutQuote.OwnerAccountAddress},
			// NOTE: bridge_refund_gas_price, refund_pegout_gas_price, and send_pegout_btc_fee are MISSING
		}

		singleResult := mongoDb.NewSingleResultFromDocument(oldBsonDocument, nil, nil)

		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, bson.D{primitive.E{Key: "quote_hash", Value: test.AnyHash}}).
			Return(singleResult).Once()

		result, err := repo.GetRetainedQuote(context.Background(), test.AnyHash)

		collection.AssertExpectations(t)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.NotNil(t, result.BridgeRefundGasPrice, "BridgeRefundGasPrice should not be nil after normalization")
		assert.NotNil(t, result.RefundPegoutGasPrice, "RefundPegoutGasPrice should not be nil after normalization")
		assert.NotNil(t, result.SendPegoutBtcFee, "SendPegoutBtcFee should not be nil after normalization")
		assert.Equal(t, entities.NewWei(0), result.BridgeRefundGasPrice)
		assert.Equal(t, entities.NewWei(0), result.RefundPegoutGasPrice)
		assert.Equal(t, entities.NewWei(0), result.SendPegoutBtcFee)

		assert.Equal(t, testRetainedPegoutQuote.QuoteHash, result.QuoteHash)
		assert.Equal(t, testRetainedPegoutQuote.BridgeRefundGasUsed, result.BridgeRefundGasUsed)
		assert.Equal(t, testRetainedPegoutQuote.RefundPegoutGasUsed, result.RefundPegoutGasUsed)
	})
}

func TestPegoutMongoRepository_InsertRetainedQuote(t *testing.T) {
	t.Run("Insert retained pegout quote successfully", func(t *testing.T) {
		const expectedLog = "INSERT interaction with db: {QuoteHash:27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f DepositAddress:mkE1WWdiu5VgjfugomDk8GxV6JdEEEJR9s Signature:5c9eab91c753355f87c19d09ea88b2fd02773981e513bc2821fed5ceba0d452a0a3d21e2252cb35348ce5c6803117e3abb62837beb8f5866a375ce66587d004b1c RequiredLiquidity:55 State:WaitingForDepositConfirmations UserRskTxHash:0x6b2e1e4daf8cf00c5c3534b72cdeec3526e8a38f70c11e44888b6e4ae1ee7d38 LpBtcTxHash:6ac3779dc33ad52f3409cbb909bcd458745995496a2a3954406206f6e5d4cb0e RefundPegoutTxHash:0x8e773a2826e73f8e5792304379a7e46dff38f17089c6d344335e03537b31c2bc BridgeRefundTxHash:0x4f3f6f0664a732e4c907971e75c1e3fd8671461dcb53f566660432fc47255d8b BridgeRefundGasUsed:21000 BridgeRefundGasPrice:20000000000 RefundPegoutGasUsed:22000 RefundPegoutGasPrice:25000000000 SendPegoutBtcFee:15000 BtcReleaseTxHash:0xd8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb OwnerAccountAddress:0x233845a26a4dA08E16218e7B401501D048670674}"
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.MatchedBy(func(q quote.RetainedPegoutQuote) bool {
			return q.QuoteHash == testRetainedPegoutQuote.QuoteHash && reflect.TypeOf(quote.RetainedPegoutQuote{}).NumField() == test.CountNonZeroValues(q)
		})).Return(nil, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)
		defer assertDbInteractionLog(t, expectedLog)()
		err := repo.InsertRetainedQuote(context.Background(), testRetainedPegoutQuote)
		collection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Db error when inserting retained pegout quote", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		collection.On("InsertOne", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)
		err := repo.InsertRetainedQuote(context.Background(), testRetainedPegoutQuote)
		collection.AssertExpectations(t)
		require.Error(t, err)
	})
}

func TestPegoutMongoRepository_UpdateRetainedQuote(t *testing.T) {
	const updated = "updated value"
	t.Run("Update retained pegout quote successfully", func(t *testing.T) {
		const expectedLog = "UPDATE interaction with db: {QuoteHash:27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f DepositAddress:updated value Signature:updated value RequiredLiquidity:200 State:SendPegoutFailed UserRskTxHash:updated value LpBtcTxHash:updated value RefundPegoutTxHash:updated value BridgeRefundTxHash:updated value BridgeRefundGasUsed:21000 BridgeRefundGasPrice:20000000000 RefundPegoutGasUsed:22000 RefundPegoutGasPrice:25000000000 SendPegoutBtcFee:15000 BtcReleaseTxHash:0xd8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb OwnerAccountAddress:0x233845a26a4dA08E16218e7B401501D048670674}"
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
		conn := mongo.NewConnection(client, time.Duration(1))
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
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)
		err := repo.UpdateRetainedQuote(context.Background(), testRetainedPegoutQuote)
		collection.AssertExpectations(t)
		require.Error(t, err)
	})
	t.Run("Retained pegout quote to update not found", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).
			Return(&mongoDb.UpdateResult{ModifiedCount: 0}, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)
		err := repo.UpdateRetainedQuote(context.Background(), testRetainedPegoutQuote)
		collection.AssertExpectations(t)
		require.ErrorIs(t, err, usecases.QuoteNotFoundError)
	})
	t.Run("Update more than one retained pegout quote", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).
			Return(&mongoDb.UpdateResult{ModifiedCount: 2}, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)
		err := repo.UpdateRetainedQuote(context.Background(), testRetainedPegoutQuote)
		collection.AssertExpectations(t)
		require.ErrorContains(t, err, "multiple documents updated")
	})
}

// nolint:funlen
func TestPegoutMongoRepository_DeleteQuotes(t *testing.T) {
	var hashes = []string{"pegout1", "pegout2", "pegout3"}
	log.SetLevel(log.DebugLevel)
	t.Run("Delete quotes successfully", func(t *testing.T) {
		client, quoteCollection := getClientAndCollectionMocks(mongo.PegoutQuoteCollection)
		retainedCollection := &mocks.CollectionBindingMock{}
		creationDataCollection := &mocks.CollectionBindingMock{}
		parsedClientMock, ok := client.Database(mongo.DbName).(*mocks.DbBindingMock)
		require.True(t, ok)
		parsedClientMock.On("Collection", mongo.RetainedPegoutQuoteCollection).Return(retainedCollection)
		parsedClientMock.On("Collection", mongo.PegoutCreationDataCollection).Return(creationDataCollection)
		quoteCollection.On("DeleteMany", mock.Anything,
			bson.D{primitive.E{Key: "hash", Value: bson.D{primitive.E{Key: "$in", Value: hashes}}}},
		).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything,
			bson.D{primitive.E{Key: "quote_hash", Value: bson.D{primitive.E{Key: "$in", Value: hashes}}}},
		).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		creationDataCollection.On("DeleteMany", mock.Anything,
			bson.D{primitive.E{Key: "hash", Value: bson.D{primitive.E{Key: "$in", Value: hashes}}}},
		).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)
		defer assertDbInteractionLog(t, mongo.Delete)()
		count, err := repo.DeleteQuotes(context.Background(), hashes)
		quoteCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
		creationDataCollection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, uint(9), count)
	})
	t.Run("Db error when deleting pegout quotes", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegoutQuoteCollection)
		collection.On("DeleteMany", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)
		count, err := repo.DeleteQuotes(context.Background(), []string{test.AnyString})
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Zero(t, count)
	})
	t.Run("Db error when deleting retained pegout quotes", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegoutQuoteCollection)
		retainedCollection := &mocks.CollectionBindingMock{}
		parsedClientMock, ok := client.Database(mongo.DbName).(*mocks.DbBindingMock)
		require.True(t, ok)
		parsedClientMock.On("Collection", mongo.RetainedPegoutQuoteCollection).Return(retainedCollection)
		collection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)
		count, err := repo.DeleteQuotes(context.Background(), []string{test.AnyString})
		collection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
		require.Error(t, err)
		assert.Zero(t, count)
	})
	t.Run("Db error when deleting pegout creation data", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegoutQuoteCollection)
		retainedCollection := &mocks.CollectionBindingMock{}
		creationDataCollection := &mocks.CollectionBindingMock{}
		parsedClientMock, ok := client.Database(mongo.DbName).(*mocks.DbBindingMock)
		require.True(t, ok)
		parsedClientMock.On("Collection", mongo.RetainedPegoutQuoteCollection).Return(retainedCollection)
		parsedClientMock.On("Collection", mongo.PegoutCreationDataCollection).Return(creationDataCollection)
		collection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		creationDataCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)
		count, err := repo.DeleteQuotes(context.Background(), []string{test.AnyString})
		collection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
		creationDataCollection.AssertExpectations(t)
		require.Error(t, err)
		assert.Zero(t, count)
	})
	t.Run("Error when deletion count missmatch", func(t *testing.T) {
		client, quoteCollection := getClientAndCollectionMocks(mongo.PegoutQuoteCollection)
		retainedCollection := &mocks.CollectionBindingMock{}
		creationDataCollection := &mocks.CollectionBindingMock{}
		parsedClientMock, ok := client.Database(mongo.DbName).(*mocks.DbBindingMock)
		require.True(t, ok)
		parsedClientMock.On("Collection", mongo.RetainedPegoutQuoteCollection).Return(retainedCollection)
		parsedClientMock.On("Collection", mongo.PegoutCreationDataCollection).Return(creationDataCollection)
		quoteCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		retainedCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 4}, nil).Once()
		creationDataCollection.On("DeleteMany", mock.Anything, mock.Anything).Return(&mongoDb.DeleteResult{DeletedCount: 3}, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)
		count, err := repo.DeleteQuotes(context.Background(), hashes)
		quoteCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
		creationDataCollection.AssertExpectations(t)
		require.ErrorContains(t, err, "pegout quote collections didn't match")
		assert.Zero(t, count)
	})
}

//nolint:funlen
func TestPegoutMongoRepository_GetRetainedQuoteByState(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
	log.SetLevel(log.DebugLevel)
	states := []quote.PegoutState{quote.PegoutStateSendPegoutSucceeded, quote.PegoutStateSendPegoutFailed}
	t.Run("Get retained pegout quotes by state successfully", func(t *testing.T) {
		const expectedLog = "READ interaction with db: [{QuoteHash:27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f DepositAddress:mkE1WWdiu5VgjfugomDk8GxV6JdEEEJR9s Signature:5c9eab91c753355f87c19d09ea88b2fd02773981e513bc2821fed5ceba0d452a0a3d21e2252cb35348ce5c6803117e3abb62837beb8f5866a375ce66587d004b1c RequiredLiquidity:55 State:WaitingForDepositConfirmations UserRskTxHash:0x6b2e1e4daf8cf00c5c3534b72cdeec3526e8a38f70c11e44888b6e4ae1ee7d38 LpBtcTxHash:6ac3779dc33ad52f3409cbb909bcd458745995496a2a3954406206f6e5d4cb0e RefundPegoutTxHash:0x8e773a2826e73f8e5792304379a7e46dff38f17089c6d344335e03537b31c2bc BridgeRefundTxHash:0x4f3f6f0664a732e4c907971e75c1e3fd8671461dcb53f566660432fc47255d8b BridgeRefundGasUsed:21000 BridgeRefundGasPrice:20000000000 RefundPegoutGasUsed:22000 RefundPegoutGasPrice:25000000000 SendPegoutBtcFee:15000 BtcReleaseTxHash:0xd8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb OwnerAccountAddress:0x233845a26a4dA08E16218e7B401501D048670674} {QuoteHash:other hash DepositAddress:mkE1WWdiu5VgjfugomDk8GxV6JdEEEJR9s Signature:456 RequiredLiquidity:777 State:WaitingForDepositConfirmations UserRskTxHash:0x6b2e1e4daf8cf00c5c3534b72cdeec3526e8a38f70c11e44888b6e4ae1ee7d38 LpBtcTxHash:6ac3779dc33ad52f3409cbb909bcd458745995496a2a3954406206f6e5d4cb0e RefundPegoutTxHash:0x8e773a2826e73f8e5792304379a7e46dff38f17089c6d344335e03537b31c2bc BridgeRefundTxHash:0x4f3f6f0664a732e4c907971e75c1e3fd8671461dcb53f566660432fc47255d8b BridgeRefundGasUsed:21000 BridgeRefundGasPrice:20000000000 RefundPegoutGasUsed:22000 RefundPegoutGasPrice:25000000000 SendPegoutBtcFee:15000 BtcReleaseTxHash:0xd8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb OwnerAccountAddress:0x233845a26a4dA08E16218e7B401501D048670674}]"
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
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
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("Find", mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()
		result, err := repo.GetRetainedQuoteByState(context.Background(), states...)
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("FillZeroValues is applied to retained pegout quotes with missing gas fields", func(t *testing.T) {
		// Mock strategy for similar tests did not work for mock limitations on unmarshalling into a struct
		firstOldDocument := bson.D{
			{Key: "quote_hash", Value: "state_first"},
			{Key: "deposit_address", Value: testRetainedPegoutQuote.DepositAddress},
			{Key: "signature", Value: testRetainedPegoutQuote.Signature},
			{Key: "required_liquidity", Value: testRetainedPegoutQuote.RequiredLiquidity.String()},
			{Key: "state", Value: quote.PegoutStateSendPegoutSucceeded},
			{Key: "user_rsk_tx_hash", Value: testRetainedPegoutQuote.UserRskTxHash},
			{Key: "lp_btc_tx_hash", Value: testRetainedPegoutQuote.LpBtcTxHash},
			{Key: "refund_pegout_tx_hash", Value: testRetainedPegoutQuote.RefundPegoutTxHash},
			{Key: "bridge_refund_tx_hash", Value: testRetainedPegoutQuote.BridgeRefundTxHash},
			{Key: "bridge_refund_gas_used", Value: testRetainedPegoutQuote.BridgeRefundGasUsed},
			{Key: "refund_pegout_gas_used", Value: testRetainedPegoutQuote.RefundPegoutGasUsed},
			{Key: "btc_release_tx_hash", Value: testRetainedPegoutQuote.BtcReleaseTxHash},
			{Key: "owner_account_address", Value: testRetainedPegoutQuote.OwnerAccountAddress},
			// NOTE: bridge_refund_gas_price, refund_pegout_gas_price, and send_pegout_btc_fee are MISSING
		}

		secondOldDocument := bson.D{
			{Key: "quote_hash", Value: "state_second"},
			{Key: "deposit_address", Value: testRetainedPegoutQuote.DepositAddress},
			{Key: "signature", Value: "state_signature"},
			{Key: "required_liquidity", Value: entities.NewWei(400).String()},
			{Key: "state", Value: quote.PegoutStateSendPegoutFailed},
			{Key: "user_rsk_tx_hash", Value: testRetainedPegoutQuote.UserRskTxHash},
			{Key: "lp_btc_tx_hash", Value: testRetainedPegoutQuote.LpBtcTxHash},
			{Key: "refund_pegout_tx_hash", Value: testRetainedPegoutQuote.RefundPegoutTxHash},
			{Key: "bridge_refund_tx_hash", Value: testRetainedPegoutQuote.BridgeRefundTxHash},
			{Key: "bridge_refund_gas_used", Value: uint64(25000)},
			{Key: "refund_pegout_gas_used", Value: uint64(30000)},
			{Key: "btc_release_tx_hash", Value: testRetainedPegoutQuote.BtcReleaseTxHash},
			{Key: "owner_account_address", Value: testRetainedPegoutQuote.OwnerAccountAddress},
			// NOTE: bridge_refund_gas_price, refund_pegout_gas_price, and send_pegout_btc_fee are MISSING
		}

		cursor, err := mongoDb.NewCursorFromDocuments([]any{firstOldDocument, secondOldDocument}, nil, nil)
		require.NoError(t, err)

		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("Find", mock.Anything,
			bson.D{primitive.E{Key: "state", Value: bson.D{primitive.E{Key: "$in", Value: states}}}},
		).Return(cursor, nil).Once()

		result, err := repo.GetRetainedQuoteByState(context.Background(), states...)

		collection.AssertExpectations(t)
		require.NoError(t, err)
		require.Len(t, result, 2)

		// Verify normalization applied to first document
		assert.NotNil(t, result[0].BridgeRefundGasPrice, "BridgeRefundGasPrice should not be nil after normalization")
		assert.NotNil(t, result[0].RefundPegoutGasPrice, "RefundPegoutGasPrice should not be nil after normalization")
		assert.NotNil(t, result[0].SendPegoutBtcFee, "SendPegoutBtcFee should not be nil after normalization")
		assert.Equal(t, entities.NewWei(0), result[0].BridgeRefundGasPrice)
		assert.Equal(t, entities.NewWei(0), result[0].RefundPegoutGasPrice)
		assert.Equal(t, entities.NewWei(0), result[0].SendPegoutBtcFee)
		assert.Equal(t, "state_first", result[0].QuoteHash)

		// Verify normalization applied to second document
		assert.NotNil(t, result[1].BridgeRefundGasPrice, "BridgeRefundGasPrice should not be nil after normalization")
		assert.NotNil(t, result[1].RefundPegoutGasPrice, "RefundPegoutGasPrice should not be nil after normalization")
		assert.NotNil(t, result[1].SendPegoutBtcFee, "SendPegoutBtcFee should not be nil after normalization")
		assert.Equal(t, entities.NewWei(0), result[1].BridgeRefundGasPrice)
		assert.Equal(t, entities.NewWei(0), result[1].RefundPegoutGasPrice)
		assert.Equal(t, entities.NewWei(0), result[1].SendPegoutBtcFee)
		assert.Equal(t, "state_second", result[1].QuoteHash)
		assert.Equal(t, uint64(25000), result[1].BridgeRefundGasUsed)
		assert.Equal(t, uint64(30000), result[1].RefundPegoutGasUsed)
	})
}

func TestPegoutMongoRepository_ListPegoutDepositsByAddress(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.DepositEventsCollection)
	log.SetLevel(log.DebugLevel)
	t.Run("List pegout deposits by address successfully", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
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
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("Find", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()
		result, err := repo.ListPegoutDepositsByAddress(context.Background(), test.AnyAddress)
		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Should sanitize address properly", func(t *testing.T) {
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
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
		conn := mongo.NewConnection(client, time.Duration(1))
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
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)
		err := repo.UpsertPegoutDeposit(context.Background(), testPegoutDeposit)
		collection.AssertExpectations(t)
		require.Error(t, err)
	})
	t.Run("Error when upserting more than one pegout deposit", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.DepositEventsCollection)
		collection.On("ReplaceOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(&mongoDb.UpdateResult{ModifiedCount: 2}, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
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
		conn := mongo.NewConnection(client, time.Duration(1))
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
		conn := mongo.NewConnection(client, time.Duration(1))
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
		const expectedLog = "UPDATE interaction with db: [{QuoteHash:quote1 DepositAddress:any address Signature:any value RequiredLiquidity:1000 State:SendPegoutSucceeded UserRskTxHash: LpBtcTxHash: RefundPegoutTxHash: BridgeRefundTxHash: BridgeRefundGasUsed:0 BridgeRefundGasPrice:<nil> RefundPegoutGasUsed:0 RefundPegoutGasPrice:<nil> SendPegoutBtcFee:<nil> BtcReleaseTxHash: OwnerAccountAddress:} {QuoteHash:quote2 DepositAddress:any address Signature:any value RequiredLiquidity:2000 State:SendPegoutFailed UserRskTxHash: LpBtcTxHash: RefundPegoutTxHash: BridgeRefundTxHash: BridgeRefundGasUsed:0 BridgeRefundGasPrice:<nil> RefundPegoutGasUsed:0 RefundPegoutGasPrice:<nil> SendPegoutBtcFee:<nil> BtcReleaseTxHash: OwnerAccountAddress:}]"
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		session := &mocks.SessionBindingMock{}
		client.On("StartSession").Return(session, nil).Once()
		session.On("EndSession", mock.Anything).Return().Once()
		session.On("WithTransaction", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				fn, ok := args.Get(1).(func(mongoDb.SessionContext) (any, error))
				require.True(t, ok)
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
		conn := mongo.NewConnection(client, time.Duration(1))
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
		conn := mongo.NewConnection(client, time.Duration(1))
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
				fn, ok := args.Get(1).(func(mongoDb.SessionContext) (any, error))
				require.True(t, ok)
				count, err := fn(mongoDb.NewSessionContext(context.Background(), mongoDb.SessionFromContext(context.Background())))
				require.Error(t, err)
				assert.Equal(t, int64(0), count)
			}).
			Return(int64(0), assert.AnError)

		collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongoDb.UpdateResult{ModifiedCount: 1}, nil).Once()
		collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
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
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)
		defer test.AssertNoLog(t)()
		err := repo.UpdateRetainedQuotes(context.Background(), retainedQuotes)
		client.AssertExpectations(t)
		session.AssertExpectations(t)
		require.ErrorContains(t, err, "mismatch on updated documents. Expected 2, updated 1")
	})
}

func TestPegoutMongoRepository_GetPegoutCreationData(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	t.Run("read pegout creation data properly", func(t *testing.T) {
		const (
			expectedLog = "READ interaction with db: {FeeRate:1.55 FeePercentage:2.41 GasPrice:123 FixedFee:456}"
			hash        = "27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f"
		)
		client, collection := getClientAndCollectionMocks(mongo.PegoutCreationDataCollection)
		collection.EXPECT().FindOne(mock.Anything, bson.D{primitive.E{Key: "hash", Value: hash}}).
			Return(mongoDb.NewSingleResultFromDocument(testPegoutCreationData, nil, nil)).Once()
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		defer assertDbInteractionLog(t, expectedLog)()
		result := repo.GetPegoutCreationData(context.Background(), hash)
		collection.AssertExpectations(t)
		assert.Equal(t, testPegoutCreationData, result)
	})
	t.Run("return zero value on invalid hash", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegoutCreationDataCollection)
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		result := repo.GetPegoutCreationData(context.Background(), test.AnyString)
		collection.AssertNotCalled(t, "FindOne")
		assert.Equal(t, quote.PegoutCreationDataZeroValue(), result)
	})
	t.Run("return zero value on db error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegoutCreationDataCollection)
		collection.EXPECT().FindOne(mock.Anything, mock.Anything).
			Return(mongoDb.NewSingleResultFromDocument(nil, nil, nil)).Once()
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		result := repo.GetPegoutCreationData(context.Background(), test.AnyHash)
		collection.AssertExpectations(t)
		assert.Equal(t, quote.PegoutCreationDataZeroValue(), result)
	})
}

func TestPegoutMongoRepository_GetQuotes(t *testing.T) {
	t.Run("Get quotes with hash filters and timestamp filters", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		pegoutCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(pegoutCollection)

		hashList := []string{"27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f"}
		expectedQuotes := []quote.PegoutQuote{testPegoutQuote}
		pegoutCollection.On("Find", mock.Anything, mock.MatchedBy(func(filter bson.M) bool {
			return true
		}), mock.Anything).Return(mongoDb.NewCursorFromDocuments([]any{testPegoutQuote}, nil, nil))
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)

		startDateTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		endDateTime := time.Date(2025, 1, 1, 23, 59, 59, 0, time.UTC)

		result, err := repo.GetQuotesByHashesAndDate(context.Background(), hashList, startDateTime, endDateTime)

		require.NoError(t, err)
		assert.Equal(t, expectedQuotes, result)

		pegoutCollection.AssertExpectations(t)
		pegoutCollection.AssertExpectations(t)
	})

	t.Run("error reading quotes from DB", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegoutQuoteCollection)

		collection.On("Find", mock.Anything, mock.Anything).Return(nil, mongoDb.ErrNoDocuments).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)

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
func TestPegoutMongoRepository_ListQuotesByDateRange(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// Test data setup
	testHash1 := "8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819"
	testHash2 := "27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f"

	testStoredQuote1 := mongo.StoredPegoutQuote{
		PegoutQuote: testPegoutQuote,
		Hash:        testHash1,
	}
	testStoredQuote2 := mongo.StoredPegoutQuote{
		PegoutQuote: testPegoutQuote,
		Hash:        testHash2,
	}

	testRetainedQuote1 := testRetainedPegoutQuote
	testRetainedQuote1.QuoteHash = testHash1
	testRetainedQuote1.Signature = "first_signature"

	startDate := time.Date(2023, 9, 25, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 9, 26, 23, 59, 59, 0, time.UTC)

	t.Run("Successfully list quotes with pagination and retained quotes", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		pegoutCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(pegoutCollection).Times(1)
		db.EXPECT().Collection(mongo.RetainedPegoutQuoteCollection).Return(retainedCollection).Times(1)

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		pegoutCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(mongoDb.NewCursorFromDocuments([]any{testStoredQuote1, testStoredQuote2}, nil, nil)).Once()

		retainedFilter := bson.D{{Key: "quote_hash", Value: bson.D{{Key: "$in", Value: []string{testHash1, testHash2}}}}}
		retainedCollection.On("Find", mock.Anything, retainedFilter).
			Return(mongoDb.NewCursorFromDocuments([]any{testRetainedQuote1, testRetainedPegoutQuote}, nil, nil)).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)

		defer assertDbInteractionLog(t, "READ interaction with db: 2")()

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, 2, count)
		require.Len(t, result, 2)
		assert.Equal(t, testPegoutQuote, result[0].Quote)
		assert.Equal(t, testRetainedQuote1, result[0].RetainedQuote)
		assert.Equal(t, testPegoutQuote, result[1].Quote)
		assert.Equal(t, testRetainedPegoutQuote, result[1].RetainedQuote)

		pegoutCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
	})

	t.Run("Successfully list quotes without pagination", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		pegoutCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(pegoutCollection).Times(1)
		db.EXPECT().Collection(mongo.RetainedPegoutQuoteCollection).Return(retainedCollection).Times(1)

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		pegoutCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(mongoDb.NewCursorFromDocuments([]any{testStoredQuote1}, nil, nil)).Once()

		retainedFilter := bson.D{{Key: "quote_hash", Value: bson.D{{Key: "$in", Value: []string{testHash1}}}}}
		retainedCollection.On("Find", mock.Anything, retainedFilter).
			Return(mongoDb.NewCursorFromDocuments([]any{}, nil, nil)).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)

		defer assertDbInteractionLog(t, "READ interaction with db: 1")()

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 0, 0)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
		require.Len(t, result, 1)
		assert.Equal(t, testPegoutQuote, result[0].Quote)
		assert.Equal(t, quote.RetainedPegoutQuote{}, result[0].RetainedQuote)

		pegoutCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
	})

	t.Run("Successfully return empty result", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		pegoutCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(pegoutCollection).Times(1)

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		pegoutCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(mongoDb.NewCursorFromDocuments([]any{}, nil, nil)).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)

		defer assertDbInteractionLog(t, "READ interaction with db: []")()

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, 0, count)
		require.Empty(t, result)

		pegoutCollection.AssertExpectations(t)
	})

	t.Run("Successfully list quotes without retained quotes", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		pegoutCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(pegoutCollection).Times(1)
		db.EXPECT().Collection(mongo.RetainedPegoutQuoteCollection).Return(retainedCollection).Times(1)

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		pegoutCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(mongoDb.NewCursorFromDocuments([]any{testStoredQuote1}, nil, nil)).Once()

		retainedFilter := bson.D{{Key: "quote_hash", Value: bson.D{{Key: "$in", Value: []string{testHash1}}}}}
		retainedCollection.On("Find", mock.Anything, retainedFilter).
			Return(mongoDb.NewCursorFromDocuments([]any{}, nil, nil)).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)

		defer assertDbInteractionLog(t, "READ interaction with db: 1")()

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
		require.Len(t, result, 1)
		assert.Equal(t, testPegoutQuote, result[0].Quote)
		assert.Equal(t, quote.RetainedPegoutQuote{}, result[0].RetainedQuote)

		pegoutCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
	})

	t.Run("Error when fetching quotes from database", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		pegoutCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(pegoutCollection).Times(1)

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		pegoutCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(nil, assert.AnError).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 1, 10)

		require.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		assert.Equal(t, 0, count)
		assert.Nil(t, result)

		pegoutCollection.AssertExpectations(t)
	})

	t.Run("Error when fetching retained quotes", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		pegoutCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(pegoutCollection).Times(1)
		db.EXPECT().Collection(mongo.RetainedPegoutQuoteCollection).Return(retainedCollection).Times(1)

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		pegoutCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(mongoDb.NewCursorFromDocuments([]any{testStoredQuote1}, nil, nil)).Once()

		retainedFilter := bson.D{{Key: "quote_hash", Value: bson.D{{Key: "$in", Value: []string{testHash1}}}}}
		retainedCollection.On("Find", mock.Anything, retainedFilter).
			Return(nil, assert.AnError).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 1, 10)

		require.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		assert.Equal(t, 1, count)
		require.Len(t, result, 1)

		pegoutCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
	})

	t.Run("Successfully handle pagination edge cases", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		pegoutCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(pegoutCollection).Times(1)
		db.EXPECT().Collection(mongo.RetainedPegoutQuoteCollection).Return(retainedCollection).Times(1)

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		pegoutCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(mongoDb.NewCursorFromDocuments([]any{testStoredQuote1}, nil, nil)).Once()

		retainedFilter := bson.D{{Key: "quote_hash", Value: bson.D{{Key: "$in", Value: []string{testHash1}}}}}
		retainedCollection.On("Find", mock.Anything, retainedFilter).
			Return(mongoDb.NewCursorFromDocuments([]any{}, nil, nil)).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)

		defer assertDbInteractionLog(t, "READ interaction with db: 1")()

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 1, 1)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
		require.Len(t, result, 1)

		pegoutCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
	})

	t.Run("Should fill zero values for retained pegout quotes with missing gas fields", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		pegoutCollection := &mocks.CollectionBindingMock{}
		retainedCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(pegoutCollection).Times(1)
		db.EXPECT().Collection(mongo.RetainedPegoutQuoteCollection).Return(retainedCollection).Times(1)

		// Create old database record with missing gas fields (represented as BSON document)
		oldRetainedDocument := bson.D{
			{Key: "quote_hash", Value: testHash1},
			{Key: "deposit_address", Value: "test_deposit_address"},
			{Key: "signature", Value: "test_signature"},
			{Key: "required_liquidity", Value: uint64(1000000)},
			{Key: "state", Value: "WaitingForDeposit"},
			{Key: "bridge_refund_gas_used", Value: uint64(21000)},
			{Key: "refund_pegout_gas_used", Value: uint64(21000)},
			{Key: "owner_account_address", Value: "0x123"},
			// Note: BridgeRefundGasPrice, RefundPegoutGasPrice, and SendPegoutBtcFee are missing (nil)
		}

		expectedFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}}}

		pegoutCollection.On("Find", mock.Anything, expectedFilter, mock.Anything).
			Return(mongoDb.NewCursorFromDocuments([]any{testStoredQuote1}, nil, nil)).Once()

		retainedFilter := bson.D{{Key: "quote_hash", Value: bson.D{{Key: "$in", Value: []string{testHash1}}}}}
		retainedCollection.On("Find", mock.Anything, retainedFilter).
			Return(mongoDb.NewCursorFromDocuments([]any{oldRetainedDocument}, nil, nil)).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPegoutMongoRepository(conn)

		defer assertDbInteractionLog(t, "READ interaction with db: 1")()

		result, count, err := repo.ListQuotesByDateRange(context.Background(), startDate, endDate, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
		require.Len(t, result, 1)
		assert.Equal(t, testPegoutQuote, result[0].Quote)

		// Verify that FillZeroValues() was applied - gas prices should be zero Wei instead of nil
		assert.NotNil(t, result[0].RetainedQuote.BridgeRefundGasPrice)
		assert.NotNil(t, result[0].RetainedQuote.RefundPegoutGasPrice)
		assert.NotNil(t, result[0].RetainedQuote.SendPegoutBtcFee)
		assert.Equal(t, entities.NewWei(0), result[0].RetainedQuote.BridgeRefundGasPrice)
		assert.Equal(t, entities.NewWei(0), result[0].RetainedQuote.RefundPegoutGasPrice)
		assert.Equal(t, entities.NewWei(0), result[0].RetainedQuote.SendPegoutBtcFee)

		pegoutCollection.AssertExpectations(t)
		retainedCollection.AssertExpectations(t)
	})
}

// nolint:funlen
func TestPegoutMongoRepository_GetRetainedQuotesForAddress(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	const address = "0xc2A630c053D12D63d32b025082f6Ba268db18300"

	t.Run("Get retained pegout quotes for address with specific state", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		const expectedLog = "READ interaction with db: [{QuoteHash:27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f DepositAddress:mkE1WWdiu5VgjfugomDk8GxV6JdEEEJR9s Signature:5c9eab91c753355f87c19d09ea88b2fd02773981e513bc2821fed5ceba0d452a0a3d21e2252cb35348ce5c6803117e3abb62837beb8f5866a375ce66587d004b1c RequiredLiquidity:55 State:WaitingForDeposit UserRskTxHash:0x6b2e1e4daf8cf00c5c3534b72cdeec3526e8a38f70c11e44888b6e4ae1ee7d38 LpBtcTxHash:6ac3779dc33ad52f3409cbb909bcd458745995496a2a3954406206f6e5d4cb0e RefundPegoutTxHash:0x8e773a2826e73f8e5792304379a7e46dff38f17089c6d344335e03537b31c2bc BridgeRefundTxHash:0x4f3f6f0664a732e4c907971e75c1e3fd8671461dcb53f566660432fc47255d8b BridgeRefundGasUsed:21000 BridgeRefundGasPrice:20000000000 RefundPegoutGasUsed:22000 RefundPegoutGasPrice:25000000000 SendPegoutBtcFee:15000 BtcReleaseTxHash:0xd8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb OwnerAccountAddress:0xc2A630c053D12D63d32b025082f6Ba268db18300}]"

		mockQuote := testRetainedPegoutQuote
		mockQuote.State = quote.PegoutStateWaitingForDeposit
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
			stateValues, ok := stateFilter[0].Value.([]quote.PegoutState)
			assert.True(t, ok)
			assert.Len(t, stateValues, 1)
			assert.Contains(t, stateValues, quote.PegoutStateWaitingForDeposit)

			return true
		})).Return(mongoDb.NewCursorFromDocuments([]any{mockQuote}, nil, nil)).Once()

		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetRetainedQuotesForAddress(context.Background(), address, quote.PegoutStateWaitingForDeposit)

		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, mockQuote, result[0])
	})

	t.Run("Get retained pegout quotes for address with multiple specific states", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		const expectedLog = "READ interaction with db: [{QuoteHash:27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f DepositAddress:mkE1WWdiu5VgjfugomDk8GxV6JdEEEJR9s Signature:5c9eab91c753355f87c19d09ea88b2fd02773981e513bc2821fed5ceba0d452a0a3d21e2252cb35348ce5c6803117e3abb62837beb8f5866a375ce66587d004b1c RequiredLiquidity:55 State:SendPegoutSucceeded UserRskTxHash:0x6b2e1e4daf8cf00c5c3534b72cdeec3526e8a38f70c11e44888b6e4ae1ee7d38 LpBtcTxHash:6ac3779dc33ad52f3409cbb909bcd458745995496a2a3954406206f6e5d4cb0e RefundPegoutTxHash:0x8e773a2826e73f8e5792304379a7e46dff38f17089c6d344335e03537b31c2bc BridgeRefundTxHash:0x4f3f6f0664a732e4c907971e75c1e3fd8671461dcb53f566660432fc47255d8b BridgeRefundGasUsed:21000 BridgeRefundGasPrice:20000000000 RefundPegoutGasUsed:22000 RefundPegoutGasPrice:25000000000 SendPegoutBtcFee:15000 BtcReleaseTxHash:0xd8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb OwnerAccountAddress:0xc2A630c053D12D63d32b025082f6Ba268db18300} {QuoteHash:second DepositAddress:mkE1WWdiu5VgjfugomDk8GxV6JdEEEJR9s Signature:123 RequiredLiquidity:777 State:SendPegoutFailed UserRskTxHash:0x6b2e1e4daf8cf00c5c3534b72cdeec3526e8a38f70c11e44888b6e4ae1ee7d38 LpBtcTxHash:6ac3779dc33ad52f3409cbb909bcd458745995496a2a3954406206f6e5d4cb0e RefundPegoutTxHash:0x8e773a2826e73f8e5792304379a7e46dff38f17089c6d344335e03537b31c2bc BridgeRefundTxHash:0x4f3f6f0664a732e4c907971e75c1e3fd8671461dcb53f566660432fc47255d8b BridgeRefundGasUsed:21000 BridgeRefundGasPrice:20000000000 RefundPegoutGasUsed:22000 RefundPegoutGasPrice:25000000000 SendPegoutBtcFee:15000 BtcReleaseTxHash:0xd8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb OwnerAccountAddress:0xc2A630c053D12D63d32b025082f6Ba268db18300}]"

		firstQuote := testRetainedPegoutQuote
		firstQuote.State = quote.PegoutStateSendPegoutSucceeded
		firstQuote.OwnerAccountAddress = address

		secondQuote := testRetainedPegoutQuote
		secondQuote.QuoteHash = "second"
		secondQuote.Signature = "123"
		secondQuote.RequiredLiquidity = entities.NewWei(777)
		secondQuote.State = quote.PegoutStateSendPegoutFailed
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
			stateValues, ok := stateFilter[0].Value.([]quote.PegoutState)
			assert.True(t, ok)
			assert.Len(t, stateValues, 2)
			assert.Contains(t, stateValues, quote.PegoutStateSendPegoutSucceeded)
			assert.Contains(t, stateValues, quote.PegoutStateSendPegoutFailed)

			return true
		})).Return(mongoDb.NewCursorFromDocuments([]any{firstQuote, secondQuote}, nil, nil)).Once()

		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetRetainedQuotesForAddress(context.Background(), address, quote.PegoutStateSendPegoutSucceeded, quote.PegoutStateSendPegoutFailed)

		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, firstQuote, result[0])
		assert.Equal(t, secondQuote, result[1])
	})

	t.Run("Empty result with no matching quotes", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
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
			stateValues, ok := stateFilter[0].Value.([]quote.PegoutState)
			assert.True(t, ok)
			assert.Len(t, stateValues, 2)
			assert.Contains(t, stateValues, quote.PegoutStateWaitingForDeposit)
			assert.Contains(t, stateValues, quote.PegoutStateWaitingForDepositConfirmations)

			return true
		})).Return(mongoDb.NewCursorFromDocuments([]any{}, nil, nil)).Once()

		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetRetainedQuotesForAddress(context.Background(), address, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations)

		collection.AssertExpectations(t)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Db error when getting retained pegout quotes for address", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))

		collection.On("Find", mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()

		result, err := repo.GetRetainedQuotesForAddress(context.Background(), address)

		collection.AssertExpectations(t)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("FillZeroValues is applied to retained pegout quotes with missing gas fields", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))

		// Mock strategy for similar tests did not work for mock limitations on unmarshalling into a struct
		firstOldDocument := bson.D{
			{Key: "quote_hash", Value: "address_first"},
			{Key: "deposit_address", Value: testRetainedPegoutQuote.DepositAddress},
			{Key: "signature", Value: testRetainedPegoutQuote.Signature},
			{Key: "required_liquidity", Value: testRetainedPegoutQuote.RequiredLiquidity.String()},
			{Key: "state", Value: quote.PegoutStateSendPegoutSucceeded},
			{Key: "user_rsk_tx_hash", Value: testRetainedPegoutQuote.UserRskTxHash},
			{Key: "lp_btc_tx_hash", Value: testRetainedPegoutQuote.LpBtcTxHash},
			{Key: "refund_pegout_tx_hash", Value: testRetainedPegoutQuote.RefundPegoutTxHash},
			{Key: "bridge_refund_tx_hash", Value: testRetainedPegoutQuote.BridgeRefundTxHash},
			{Key: "bridge_refund_gas_used", Value: testRetainedPegoutQuote.BridgeRefundGasUsed},
			{Key: "refund_pegout_gas_used", Value: testRetainedPegoutQuote.RefundPegoutGasUsed},
			{Key: "btc_release_tx_hash", Value: testRetainedPegoutQuote.BtcReleaseTxHash},
			{Key: "owner_account_address", Value: address},
			// NOTE: bridge_refund_gas_price, refund_pegout_gas_price, and send_pegout_btc_fee are MISSING
		}

		secondOldDocument := bson.D{
			{Key: "quote_hash", Value: "address_second"},
			{Key: "deposit_address", Value: testRetainedPegoutQuote.DepositAddress},
			{Key: "signature", Value: "address_signature"},
			{Key: "required_liquidity", Value: entities.NewWei(600).String()},
			{Key: "state", Value: quote.PegoutStateSendPegoutFailed},
			{Key: "user_rsk_tx_hash", Value: testRetainedPegoutQuote.UserRskTxHash},
			{Key: "lp_btc_tx_hash", Value: testRetainedPegoutQuote.LpBtcTxHash},
			{Key: "refund_pegout_tx_hash", Value: testRetainedPegoutQuote.RefundPegoutTxHash},
			{Key: "bridge_refund_tx_hash", Value: testRetainedPegoutQuote.BridgeRefundTxHash},
			{Key: "bridge_refund_gas_used", Value: uint64(35000)},
			{Key: "refund_pegout_gas_used", Value: uint64(40000)},
			{Key: "btc_release_tx_hash", Value: testRetainedPegoutQuote.BtcReleaseTxHash},
			{Key: "owner_account_address", Value: address},
			// NOTE: bridge_refund_gas_price, refund_pegout_gas_price, and send_pegout_btc_fee are MISSING
		}

		cursor, err := mongoDb.NewCursorFromDocuments([]any{firstOldDocument, secondOldDocument}, nil, nil)
		require.NoError(t, err)

		collection.On("Find", mock.Anything, mock.Anything).
			Return(cursor, nil).Once()

		result, err := repo.GetRetainedQuotesForAddress(context.Background(), address, quote.PegoutStateSendPegoutSucceeded, quote.PegoutStateSendPegoutFailed)

		collection.AssertExpectations(t)
		require.NoError(t, err)
		require.Len(t, result, 2)

		// Verify normalization applied to first document
		assert.NotNil(t, result[0].BridgeRefundGasPrice, "BridgeRefundGasPrice should not be nil after normalization")
		assert.NotNil(t, result[0].RefundPegoutGasPrice, "RefundPegoutGasPrice should not be nil after normalization")
		assert.NotNil(t, result[0].SendPegoutBtcFee, "SendPegoutBtcFee should not be nil after normalization")
		assert.Equal(t, entities.NewWei(0), result[0].BridgeRefundGasPrice)
		assert.Equal(t, entities.NewWei(0), result[0].RefundPegoutGasPrice)
		assert.Equal(t, entities.NewWei(0), result[0].SendPegoutBtcFee)
		assert.Equal(t, "address_first", result[0].QuoteHash)

		// Verify normalization applied to second document
		assert.NotNil(t, result[1].BridgeRefundGasPrice, "BridgeRefundGasPrice should not be nil after normalization")
		assert.NotNil(t, result[1].RefundPegoutGasPrice, "RefundPegoutGasPrice should not be nil after normalization")
		assert.NotNil(t, result[1].SendPegoutBtcFee, "SendPegoutBtcFee should not be nil after normalization")
		assert.Equal(t, entities.NewWei(0), result[1].BridgeRefundGasPrice)
		assert.Equal(t, entities.NewWei(0), result[1].RefundPegoutGasPrice)
		assert.Equal(t, entities.NewWei(0), result[1].SendPegoutBtcFee)
		assert.Equal(t, "address_second", result[1].QuoteHash)
		assert.Equal(t, uint64(35000), result[1].BridgeRefundGasUsed)
		assert.Equal(t, uint64(40000), result[1].RefundPegoutGasUsed)
	})
}

//nolint:funlen
func TestPegoutMongoRepository_GetRetainedQuotesInBatch(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	batch := rootstock.BatchPegOut{
		TransactionHash:    testRetainedPegoutQuote.BtcReleaseTxHash,
		BlockHash:          test.AnyString,
		BlockNumber:        1,
		BtcTxHash:          test.AnyString,
		ReleaseRskTxHashes: []string{test.AnyHash},
	}
	secondQuote := testRetainedPegoutQuote
	secondQuote.QuoteHash = "other hash"
	expectedQuotes := []quote.RetainedPegoutQuote{testRetainedPegoutQuote, secondQuote}
	t.Run("should return quotes of transactions present in batch", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))
		const expectedLog = "READ interaction with db: [{QuoteHash:27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f DepositAddress:mkE1WWdiu5VgjfugomDk8GxV6JdEEEJR9s Signature:5c9eab91c753355f87c19d09ea88b2fd02773981e513bc2821fed5ceba0d452a0a3d21e2252cb35348ce5c6803117e3abb62837beb8f5866a375ce66587d004b1c RequiredLiquidity:55 State:WaitingForDepositConfirmations UserRskTxHash:0x6b2e1e4daf8cf00c5c3534b72cdeec3526e8a38f70c11e44888b6e4ae1ee7d38 LpBtcTxHash:6ac3779dc33ad52f3409cbb909bcd458745995496a2a3954406206f6e5d4cb0e RefundPegoutTxHash:0x8e773a2826e73f8e5792304379a7e46dff38f17089c6d344335e03537b31c2bc BridgeRefundTxHash:0x4f3f6f0664a732e4c907971e75c1e3fd8671461dcb53f566660432fc47255d8b BridgeRefundGasUsed:21000 BridgeRefundGasPrice:20000000000 RefundPegoutGasUsed:22000 RefundPegoutGasPrice:25000000000 SendPegoutBtcFee:15000 BtcReleaseTxHash:0xd8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb OwnerAccountAddress:0x233845a26a4dA08E16218e7B401501D048670674} {QuoteHash:other hash DepositAddress:mkE1WWdiu5VgjfugomDk8GxV6JdEEEJR9s Signature:5c9eab91c753355f87c19d09ea88b2fd02773981e513bc2821fed5ceba0d452a0a3d21e2252cb35348ce5c6803117e3abb62837beb8f5866a375ce66587d004b1c RequiredLiquidity:55 State:WaitingForDepositConfirmations UserRskTxHash:0x6b2e1e4daf8cf00c5c3534b72cdeec3526e8a38f70c11e44888b6e4ae1ee7d38 LpBtcTxHash:6ac3779dc33ad52f3409cbb909bcd458745995496a2a3954406206f6e5d4cb0e RefundPegoutTxHash:0x8e773a2826e73f8e5792304379a7e46dff38f17089c6d344335e03537b31c2bc BridgeRefundTxHash:0x4f3f6f0664a732e4c907971e75c1e3fd8671461dcb53f566660432fc47255d8b BridgeRefundGasUsed:21000 BridgeRefundGasPrice:20000000000 RefundPegoutGasUsed:22000 RefundPegoutGasPrice:25000000000 SendPegoutBtcFee:15000 BtcReleaseTxHash:0xd8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb OwnerAccountAddress:0x233845a26a4dA08E16218e7B401501D048670674}]"

		collection.On("Find", mock.Anything,
			bson.D{
				primitive.E{Key: "state", Value: quote.PegoutStateBridgeTxSucceeded},
				primitive.E{Key: "bridge_refund_tx_hash", Value: bson.D{
					primitive.E{Key: "$in", Value: batch.ReleaseRskTxHashes},
				}},
			},
		).Return(mongoDb.NewCursorFromDocuments([]any{testRetainedPegoutQuote, secondQuote}, nil, nil)).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		retainedQuotes, err := repo.GetRetainedQuotesInBatch(context.Background(), batch)
		require.NoError(t, err)
		assert.Equal(t, expectedQuotes, retainedQuotes)
	})
	t.Run("should return empty slice when no quotes found", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))

		collection.On("Find", mock.Anything,
			bson.D{
				primitive.E{Key: "state", Value: quote.PegoutStateBridgeTxSucceeded},
				primitive.E{Key: "bridge_refund_tx_hash", Value: bson.D{
					primitive.E{Key: "$in", Value: batch.ReleaseRskTxHashes},
				}},
			},
		).Return(mongoDb.NewCursorFromDocuments([]any{}, nil, nil)).Once()
		retainedQuotes, err := repo.GetRetainedQuotesInBatch(context.Background(), batch)
		require.NoError(t, err)
		assert.Empty(t, retainedQuotes)
	})
	t.Run("should handle error reading from database", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))

		collection.On("Find", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		retainedQuotes, err := repo.GetRetainedQuotesInBatch(context.Background(), batch)
		require.Error(t, err)
		assert.Empty(t, retainedQuotes)
	})
	t.Run("FillZeroValues is applied to retained pegout quotes with missing gas fields", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.RetainedPegoutQuoteCollection)
		repo := mongo.NewPegoutMongoRepository(mongo.NewConnection(client, time.Duration(1)))

		// Mock strategy for similar tests did not work for mock limitations on unmarshalling into a struct
		firstOldDocument := bson.D{
			{Key: "quote_hash", Value: "batch_first"},
			{Key: "deposit_address", Value: testRetainedPegoutQuote.DepositAddress},
			{Key: "signature", Value: testRetainedPegoutQuote.Signature},
			{Key: "required_liquidity", Value: testRetainedPegoutQuote.RequiredLiquidity.String()},
			{Key: "state", Value: testRetainedPegoutQuote.State},
			{Key: "user_rsk_tx_hash", Value: testRetainedPegoutQuote.UserRskTxHash},
			{Key: "lp_btc_tx_hash", Value: testRetainedPegoutQuote.LpBtcTxHash},
			{Key: "refund_pegout_tx_hash", Value: testRetainedPegoutQuote.RefundPegoutTxHash},
			{Key: "bridge_refund_tx_hash", Value: testRetainedPegoutQuote.BridgeRefundTxHash},
			{Key: "bridge_refund_gas_used", Value: testRetainedPegoutQuote.BridgeRefundGasUsed},
			{Key: "refund_pegout_gas_used", Value: testRetainedPegoutQuote.RefundPegoutGasUsed},
			{Key: "btc_release_tx_hash", Value: batch.TransactionHash},
			{Key: "owner_account_address", Value: testRetainedPegoutQuote.OwnerAccountAddress},
			// NOTE: bridge_refund_gas_price, refund_pegout_gas_price, and send_pegout_btc_fee are MISSING
		}

		secondOldDocument := bson.D{
			{Key: "quote_hash", Value: "batch_second"},
			{Key: "deposit_address", Value: testRetainedPegoutQuote.DepositAddress},
			{Key: "signature", Value: "batch_signature"},
			{Key: "required_liquidity", Value: entities.NewWei(800).String()},
			{Key: "state", Value: quote.PegoutStateSendPegoutSucceeded},
			{Key: "user_rsk_tx_hash", Value: testRetainedPegoutQuote.UserRskTxHash},
			{Key: "lp_btc_tx_hash", Value: testRetainedPegoutQuote.LpBtcTxHash},
			{Key: "refund_pegout_tx_hash", Value: testRetainedPegoutQuote.RefundPegoutTxHash},
			{Key: "bridge_refund_tx_hash", Value: testRetainedPegoutQuote.BridgeRefundTxHash},
			{Key: "bridge_refund_gas_used", Value: uint64(45000)},
			{Key: "refund_pegout_gas_used", Value: uint64(50000)},
			{Key: "btc_release_tx_hash", Value: batch.TransactionHash},
			{Key: "owner_account_address", Value: testRetainedPegoutQuote.OwnerAccountAddress},
			// NOTE: bridge_refund_gas_price, refund_pegout_gas_price, and send_pegout_btc_fee are MISSING
		}

		cursor, err := mongoDb.NewCursorFromDocuments([]any{firstOldDocument, secondOldDocument}, nil, nil)
		require.NoError(t, err)

		collection.On("Find", mock.Anything, mock.Anything).Return(cursor, nil).Once()

		result, err := repo.GetRetainedQuotesInBatch(context.Background(), batch)

		collection.AssertExpectations(t)
		require.NoError(t, err)
		require.Len(t, result, 2)

		// Verify normalization applied to first document
		assert.NotNil(t, result[0].BridgeRefundGasPrice, "BridgeRefundGasPrice should not be nil after normalization")
		assert.NotNil(t, result[0].RefundPegoutGasPrice, "RefundPegoutGasPrice should not be nil after normalization")
		assert.NotNil(t, result[0].SendPegoutBtcFee, "SendPegoutBtcFee should not be nil after normalization")
		assert.Equal(t, entities.NewWei(0), result[0].BridgeRefundGasPrice)
		assert.Equal(t, entities.NewWei(0), result[0].RefundPegoutGasPrice)
		assert.Equal(t, entities.NewWei(0), result[0].SendPegoutBtcFee)
		assert.Equal(t, "batch_first", result[0].QuoteHash)

		// Verify normalization applied to second document
		assert.NotNil(t, result[1].BridgeRefundGasPrice, "BridgeRefundGasPrice should not be nil after normalization")
		assert.NotNil(t, result[1].RefundPegoutGasPrice, "RefundPegoutGasPrice should not be nil after normalization")
		assert.NotNil(t, result[1].SendPegoutBtcFee, "SendPegoutBtcFee should not be nil after normalization")
		assert.Equal(t, entities.NewWei(0), result[1].BridgeRefundGasPrice)
		assert.Equal(t, entities.NewWei(0), result[1].RefundPegoutGasPrice)
		assert.Equal(t, entities.NewWei(0), result[1].SendPegoutBtcFee)
		assert.Equal(t, "batch_second", result[1].QuoteHash)
		assert.Equal(t, uint64(45000), result[1].BridgeRefundGasUsed)
		assert.Equal(t, uint64(50000), result[1].RefundPegoutGasUsed)
	})
}

// nolint: cyclop, funlen
func validatePegoutPipelineStructure(pipeline mongoDb.Pipeline, states []quote.PegoutState, startDate, endDate time.Time) bool {
	// Verify the pipeline structure
	if len(pipeline) != 4 {
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

	// Stage 4: $match by state - verify states are correct
	matchStateStage := pipeline[3]
	if len(matchStateStage) == 0 || matchStateStage[0].Key != "$match" {
		return false
	}
	matchStateFilter, ok := matchStateStage[0].Value.(bson.M)
	if !ok {
		return false
	}
	stateFilter, ok := matchStateFilter["retained.state"].(bson.M)
	if !ok {
		return false
	}
	inStates, ok := stateFilter["$in"].([]quote.PegoutState)
	if !ok {
		return false
	}
	// Verify the states match (must have same length and contain same elements)
	if len(inStates) != len(states) {
		return false
	}
	stateMap := make(map[quote.PegoutState]bool)
	for _, s := range states {
		stateMap[s] = true
	}
	for _, s := range inStates {
		if !stateMap[s] {
			return false
		}
	}

	return true
}

func buildPegoutAggregationDocument(q quote.PegoutQuote, retained quote.RetainedPegoutQuote) bson.D {
	return bson.D{
		{Key: "lbc_address", Value: q.LbcAddress},
		{Key: "lp_rsk_address", Value: q.LpRskAddress},
		{Key: "btc_refund_address", Value: q.BtcRefundAddress},
		{Key: "rsk_refund_address", Value: q.RskRefundAddress},
		{Key: "lp_btc_address", Value: q.LpBtcAddress},
		{Key: "call_fee", Value: q.CallFee.String()},
		{Key: "penalty_fee", Value: q.PenaltyFee.String()},
		{Key: "nonce", Value: q.Nonce},
		{Key: "deposit_address", Value: q.DepositAddress},
		{Key: "value", Value: q.Value.String()},
		{Key: "agreement_timestamp", Value: q.AgreementTimestamp},
		{Key: "deposit_date_limit", Value: q.DepositDateLimit},
		{Key: "deposit_confirmations", Value: q.DepositConfirmations},
		{Key: "transfer_confirmations", Value: q.TransferConfirmations},
		{Key: "transfer_time", Value: q.TransferTime},
		{Key: "expire_date", Value: q.ExpireDate},
		{Key: "expire_block", Value: q.ExpireBlock},
		{Key: "gas_fee", Value: q.GasFee.String()},
		{Key: "product_fee_amount", Value: q.ProductFeeAmount.String()},
		{Key: "hash", Value: retained.QuoteHash},
		{Key: "retained", Value: bson.D{
			{Key: "quote_hash", Value: retained.QuoteHash},
			{Key: "deposit_address", Value: retained.DepositAddress},
			{Key: "signature", Value: retained.Signature},
			{Key: "required_liquidity", Value: retained.RequiredLiquidity.String()},
			{Key: "state", Value: retained.State},
			{Key: "user_rsk_tx_hash", Value: retained.UserRskTxHash},
			{Key: "lp_btc_tx_hash", Value: retained.LpBtcTxHash},
			{Key: "refund_pegout_tx_hash", Value: retained.RefundPegoutTxHash},
			{Key: "bridge_refund_tx_hash", Value: retained.BridgeRefundTxHash},
			{Key: "bridge_refund_gas_used", Value: retained.BridgeRefundGasUsed},
			{Key: "bridge_refund_gas_price", Value: retained.BridgeRefundGasPrice.String()},
			{Key: "refund_pegout_gas_used", Value: retained.RefundPegoutGasUsed},
			{Key: "refund_pegout_gas_price", Value: retained.RefundPegoutGasPrice.String()},
			{Key: "send_pegout_btc_fee", Value: retained.SendPegoutBtcFee.String()},
			{Key: "btc_release_tx_hash", Value: retained.BtcReleaseTxHash},
			{Key: "owner_account_address", Value: retained.OwnerAccountAddress},
		}},
	}
}

// nolint: funlen
func TestPegoutMongoRepository_GetQuotesWithRetainedByStateAndDate_Success(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	states := []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded, quote.PegoutStateBtcReleased}

	client, db := getClientAndDatabaseMocks()
	collection := &mocks.CollectionBindingMock{}
	db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(collection)

	// Create test documents that simulate aggregation pipeline results
	firstQuote := testPegoutQuote
	firstQuote.AgreementTimestamp = uint32(startDate.Add(24 * time.Hour).Unix())
	firstRetained := testRetainedPegoutQuote
	firstRetained.QuoteHash = "test-hash-1"
	firstRetained.State = quote.PegoutStateRefundPegOutSucceeded

	secondQuote := testPegoutQuote
	secondQuote.AgreementTimestamp = uint32(startDate.Add(48 * time.Hour).Unix())
	secondRetained := testRetainedPegoutQuote
	secondRetained.QuoteHash = "test-hash-2"
	secondRetained.State = quote.PegoutStateBridgeTxSucceeded

	thirdQuote := testPegoutQuote
	thirdQuote.AgreementTimestamp = uint32(startDate.Add(72 * time.Hour).Unix())
	thirdRetained := testRetainedPegoutQuote
	thirdRetained.QuoteHash = "test-hash-3"
	thirdRetained.State = quote.PegoutStateBtcReleased

	// Simulate aggregation result documents
	doc1 := buildPegoutAggregationDocument(firstQuote, firstRetained)
	doc2 := buildPegoutAggregationDocument(secondQuote, secondRetained)
	doc3 := buildPegoutAggregationDocument(thirdQuote, thirdRetained)

	cursor, err := mongoDb.NewCursorFromDocuments([]any{doc1, doc2, doc3}, nil, nil)
	require.NoError(t, err)

	collection.On("Aggregate", mock.Anything, mock.MatchedBy(func(pipeline mongoDb.Pipeline) bool {
		return validatePegoutPipelineStructure(pipeline, states, startDate, endDate)
	})).Return(cursor, nil).Once()

	conn := mongo.NewConnection(client, time.Duration(1))
	repo := mongo.NewPegoutMongoRepository(conn)

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

func TestPegoutMongoRepository_GetQuotesWithRetainedByStateAndDate_EmptyResult(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	states := []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded}

	client, db := getClientAndDatabaseMocks()
	collection := &mocks.CollectionBindingMock{}
	db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(collection)

	cursor, err := mongoDb.NewCursorFromDocuments([]any{}, nil, nil)
	require.NoError(t, err)

	collection.On("Aggregate", mock.Anything, mock.Anything).Return(cursor, nil).Once()

	conn := mongo.NewConnection(client, time.Duration(1))
	repo := mongo.NewPegoutMongoRepository(conn)

	result, err := repo.GetQuotesWithRetainedByStateAndDate(context.Background(), states, startDate, endDate)

	collection.AssertExpectations(t)
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestPegoutMongoRepository_GetQuotesWithRetainedByStateAndDate_AggregationError(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	states := []quote.PegoutState{quote.PegoutStateBridgeTxSucceeded}

	client, db := getClientAndDatabaseMocks()
	collection := &mocks.CollectionBindingMock{}
	db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(collection)

	expectedError := assert.AnError
	collection.On("Aggregate", mock.Anything, mock.Anything).Return(nil, expectedError).Once()

	conn := mongo.NewConnection(client, time.Duration(1))
	repo := mongo.NewPegoutMongoRepository(conn)

	result, err := repo.GetQuotesWithRetainedByStateAndDate(context.Background(), states, startDate, endDate)

	collection.AssertExpectations(t)
	require.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
}

func TestPegoutMongoRepository_GetQuotesWithRetainedByStateAndDate_DecodeError(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	states := []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded}

	client, db := getClientAndDatabaseMocks()
	collection := &mocks.CollectionBindingMock{}
	db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(collection)

	// Create a document with invalid data that will fail decode
	invalidDoc := bson.D{
		{Key: "call_fee", Value: "invalid"}, // Invalid type for Wei
	}

	cursor, err := mongoDb.NewCursorFromDocuments([]any{invalidDoc}, nil, nil)
	require.NoError(t, err)

	collection.On("Aggregate", mock.Anything, mock.Anything).Return(cursor, nil).Once()

	conn := mongo.NewConnection(client, time.Duration(1))
	repo := mongo.NewPegoutMongoRepository(conn)

	result, err := repo.GetQuotesWithRetainedByStateAndDate(context.Background(), states, startDate, endDate)

	collection.AssertExpectations(t)
	require.Error(t, err)
	assert.Nil(t, result)
}

// nolint: funlen
func TestPegoutMongoRepository_GetQuotesWithRetainedByStateAndDate_ZeroValueNormalization(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	states := []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded}

	client, db := getClientAndDatabaseMocks()
	collection := &mocks.CollectionBindingMock{}
	db.EXPECT().Collection(mongo.PegoutQuoteCollection).Return(collection)

	testQuote := testPegoutQuote
	testQuote.AgreementTimestamp = uint32(startDate.Add(24 * time.Hour).Unix())

	// Create document with missing gas price fields
	doc := bson.D{
		{Key: "lbc_address", Value: testQuote.LbcAddress},
		{Key: "lp_rsk_address", Value: testQuote.LpRskAddress},
		{Key: "btc_refund_address", Value: testQuote.BtcRefundAddress},
		{Key: "rsk_refund_address", Value: testQuote.RskRefundAddress},
		{Key: "lp_btc_address", Value: testQuote.LpBtcAddress},
		{Key: "call_fee", Value: testQuote.CallFee.String()},
		{Key: "penalty_fee", Value: testQuote.PenaltyFee.String()},
		{Key: "nonce", Value: testQuote.Nonce},
		{Key: "deposit_address", Value: testQuote.DepositAddress},
		{Key: "value", Value: testQuote.Value.String()},
		{Key: "agreement_timestamp", Value: testQuote.AgreementTimestamp},
		{Key: "deposit_date_limit", Value: testQuote.DepositDateLimit},
		{Key: "deposit_confirmations", Value: testQuote.DepositConfirmations},
		{Key: "transfer_confirmations", Value: testQuote.TransferConfirmations},
		{Key: "transfer_time", Value: testQuote.TransferTime},
		{Key: "expire_date", Value: testQuote.ExpireDate},
		{Key: "expire_block", Value: testQuote.ExpireBlock},
		{Key: "gas_fee", Value: testQuote.GasFee.String()},
		{Key: "product_fee_amount", Value: testQuote.ProductFeeAmount.String()},
		{Key: "hash", Value: "test-hash"},
		{Key: "retained", Value: bson.D{
			{Key: "quote_hash", Value: "test-hash"},
			{Key: "deposit_address", Value: testRetainedPegoutQuote.DepositAddress},
			{Key: "signature", Value: testRetainedPegoutQuote.Signature},
			{Key: "required_liquidity", Value: testRetainedPegoutQuote.RequiredLiquidity.String()},
			{Key: "state", Value: quote.PegoutStateRefundPegOutSucceeded},
			{Key: "user_rsk_tx_hash", Value: testRetainedPegoutQuote.UserRskTxHash},
			{Key: "lp_btc_tx_hash", Value: testRetainedPegoutQuote.LpBtcTxHash},
			{Key: "refund_pegout_tx_hash", Value: testRetainedPegoutQuote.RefundPegoutTxHash},
			{Key: "bridge_refund_tx_hash", Value: testRetainedPegoutQuote.BridgeRefundTxHash},
			{Key: "bridge_refund_gas_used", Value: uint64(21000)},
			{Key: "refund_pegout_gas_used", Value: uint64(22000)},
			{Key: "btc_release_tx_hash", Value: testRetainedPegoutQuote.BtcReleaseTxHash},
			{Key: "owner_account_address", Value: testRetainedPegoutQuote.OwnerAccountAddress},
			// NOTE: bridge_refund_gas_price, refund_pegout_gas_price, and send_pegout_btc_fee are MISSING
		}},
	}

	cursor, err := mongoDb.NewCursorFromDocuments([]any{doc}, nil, nil)
	require.NoError(t, err)

	collection.On("Aggregate", mock.Anything, mock.Anything).Return(cursor, nil).Once()

	conn := mongo.NewConnection(client, time.Duration(1))
	repo := mongo.NewPegoutMongoRepository(conn)

	result, err := repo.GetQuotesWithRetainedByStateAndDate(context.Background(), states, startDate, endDate)

	collection.AssertExpectations(t)
	require.NoError(t, err)
	require.Len(t, result, 1)

	// Verify normalization applied
	assert.NotNil(t, result[0].RetainedQuote.BridgeRefundGasPrice, "BridgeRefundGasPrice should not be nil after normalization")
	assert.NotNil(t, result[0].RetainedQuote.RefundPegoutGasPrice, "RefundPegoutGasPrice should not be nil after normalization")
	assert.NotNil(t, result[0].RetainedQuote.SendPegoutBtcFee, "SendPegoutBtcFee should not be nil after normalization")
	assert.Equal(t, entities.NewWei(0), result[0].RetainedQuote.BridgeRefundGasPrice)
	assert.Equal(t, entities.NewWei(0), result[0].RetainedQuote.RefundPegoutGasPrice)
	assert.Equal(t, entities.NewWei(0), result[0].RetainedQuote.SendPegoutBtcFee)
}
