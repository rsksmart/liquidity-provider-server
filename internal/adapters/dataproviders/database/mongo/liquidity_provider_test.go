package mongo_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

var peginTestConfig = &entities.Signed[liquidity_provider.PeginConfiguration]{
	Value: liquidity_provider.PeginConfiguration{
		TimeForDeposit: 1,
		CallTime:       2,
		PenaltyFee:     entities.NewWei(3),
		CallFee:        entities.NewWei(4),
		MaxValue:       entities.NewWei(5),
		MinValue:       entities.NewWei(6),
	},
	Signature: "pegin signature",
	Hash:      "pegin hash",
}

var pegoutTestConfig = &entities.Signed[liquidity_provider.PegoutConfiguration]{
	Value: liquidity_provider.PegoutConfiguration{
		TimeForDeposit:       1,
		ExpireTime:           2,
		PenaltyFee:           entities.NewWei(3),
		CallFee:              entities.NewWei(4),
		MaxValue:             entities.NewWei(5),
		MinValue:             entities.NewWei(6),
		ExpireBlocks:         7,
		BridgeTransactionMin: entities.NewWei(8),
	},
	Signature: "pegout signature",
	Hash:      "pegout hash",
}

var generalTestConfig = &entities.Signed[liquidity_provider.GeneralConfiguration]{
	Value: liquidity_provider.GeneralConfiguration{
		RskConfirmations: liquidity_provider.ConfirmationsPerAmount{
			1: 2,
			3: 4,
		},
		BtcConfirmations: liquidity_provider.ConfirmationsPerAmount{
			5: 6,
			7: 8,
		},
	},
	Signature: "general signature",
	Hash:      "general hash",
}

var testCredentials = &entities.Signed[liquidity_provider.HashedCredentials]{
	Value: liquidity_provider.HashedCredentials{
		HashedUsername: "username",
		HashedPassword: "password",
		UsernameSalt:   "username salt",
		PasswordSalt:   "password salt",
	},
	Signature: "credentials signature",
	Hash:      "credentials hash",
}

func TestLpMongoRepository_GetPeginConfiguration(t *testing.T) {
	filter := bson.D{primitive.E{Key: "name", Value: mongo.ConfigurationName("pegin")}}
	log.SetLevel(log.DebugLevel)
	t.Run("pegin configuration read successfully", func(t *testing.T) {
		const expectedLog = "READ interaction with db: {Value:{TimeForDeposit:1 CallTime:2 PenaltyFee:3 CallFee:4 MaxValue:5 MinValue:6} Signature:pegin signature Hash:pegin hash}"
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, filter).
			Return(mongoDb.NewSingleResultFromDocument(peginTestConfig, nil, nil)).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetPeginConfiguration(context.Background())
		require.NoError(t, err)
		assert.Equal(t, peginTestConfig, result)
	})
	t.Run("pegin configuration not found", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, filter).
			Return(
				mongoDb.NewSingleResultFromDocument(entities.Signed[liquidity_provider.PeginConfiguration]{}, mongoDb.ErrNoDocuments, nil),
			).Once()
		result, err := repo.GetPeginConfiguration(context.Background())
		require.NoError(t, err)
		assert.Nil(t, result)
	})
	t.Run("Db error reading pegin configuration", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, filter).Return(mongoDb.NewSingleResultFromDocument(nil, assert.AnError, nil)).Once()
		result, err := repo.GetPeginConfiguration(context.Background())
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestLpMongoRepository_GetPegoutConfiguration(t *testing.T) {
	filter := bson.D{primitive.E{Key: "name", Value: mongo.ConfigurationName("pegout")}}
	log.SetLevel(log.DebugLevel)
	t.Run("pegout configuration read successfully", func(t *testing.T) {
		const expectedLog = "READ interaction with db: {Value:{TimeForDeposit:1 ExpireTime:2 PenaltyFee:3 CallFee:4 MaxValue:5 MinValue:6 ExpireBlocks:7 BridgeTransactionMin:8} Signature:pegout signature Hash:pegout hash}"
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, filter).
			Return(mongoDb.NewSingleResultFromDocument(pegoutTestConfig, nil, nil)).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetPegoutConfiguration(context.Background())
		require.NoError(t, err)
		assert.Equal(t, pegoutTestConfig, result)
	})
	t.Run("pegout configuration not found", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, filter).
			Return(
				mongoDb.NewSingleResultFromDocument(entities.Signed[liquidity_provider.PegoutConfiguration]{}, mongoDb.ErrNoDocuments, nil),
			).Once()
		result, err := repo.GetPegoutConfiguration(context.Background())
		require.NoError(t, err)
		assert.Nil(t, result)
	})
	t.Run("Db error reading pegout configuration", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, filter).Return(mongoDb.NewSingleResultFromDocument(nil, assert.AnError, nil)).Once()
		result, err := repo.GetPegoutConfiguration(context.Background())
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestLpMongoRepository_GetGeneralConfiguration(t *testing.T) {
	filter := bson.D{primitive.E{Key: "name", Value: mongo.ConfigurationName("general")}}
	log.SetLevel(log.DebugLevel)
	t.Run("general configuration read successfully", func(t *testing.T) {
		const expectedLog = "READ interaction with db: {Value:{RskConfirmations:map[1:2 3:4] BtcConfirmations:map[5:6 7:8] PublicLiquidityCheck:false} Signature:general signature Hash:general hash}"
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, filter).
			Return(mongoDb.NewSingleResultFromDocument(generalTestConfig, nil, nil)).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetGeneralConfiguration(context.Background())
		require.NoError(t, err)
		assert.Equal(t, generalTestConfig, result)
	})
	t.Run("general configuration not found", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, filter).
			Return(
				mongoDb.NewSingleResultFromDocument(entities.Signed[liquidity_provider.GeneralConfiguration]{}, mongoDb.ErrNoDocuments, nil),
			).Once()
		result, err := repo.GetGeneralConfiguration(context.Background())
		require.NoError(t, err)
		assert.Nil(t, result)
	})
	t.Run("Db error reading general configuration", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, filter).Return(mongoDb.NewSingleResultFromDocument(nil, assert.AnError, nil)).Once()
		result, err := repo.GetGeneralConfiguration(context.Background())
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestLpMongoRepository_GetCredentials(t *testing.T) {
	filter := bson.D{primitive.E{Key: "name", Value: mongo.ConfigurationName("credentials")}}
	log.SetLevel(log.DebugLevel)
	t.Run("credentials read successfully", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, filter).
			Return(mongoDb.NewSingleResultFromDocument(testCredentials, nil, nil)).Once()
		defer test.AssertNoLog(t)()
		result, err := repo.GetCredentials(context.Background())
		require.NoError(t, err)
		assert.Equal(t, testCredentials, result)
	})
	t.Run("credentials not found", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, filter).
			Return(
				mongoDb.NewSingleResultFromDocument(entities.Signed[liquidity_provider.HashedCredentials]{}, mongoDb.ErrNoDocuments, nil),
			).Once()
		result, err := repo.GetCredentials(context.Background())
		require.NoError(t, err)
		assert.Nil(t, result)
	})
	t.Run("Db error reading credentials", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, filter).Return(mongoDb.NewSingleResultFromDocument(nil, assert.AnError, nil)).Once()
		result, err := repo.GetCredentials(context.Background())
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestLpMongoRepository_UpsertPeginConfiguration(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	configName := mongo.ConfigurationName("pegin")
	filter := bson.D{primitive.E{Key: "name", Value: configName}}
	t.Run("pegin configuration upserted successfully", func(t *testing.T) {
		const expectedLog = "INSERT interaction with db: {Signed:{Value:{TimeForDeposit:1 CallTime:2 PenaltyFee:3 CallFee:4 MaxValue:5 MinValue:6} Signature:pegin signature Hash:pegin hash} Name:pegin}"
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("ReplaceOne", mock.Anything, filter, mongo.StoredConfiguration[liquidity_provider.PeginConfiguration]{
			Signed: *peginTestConfig,
			Name:   configName,
		}, options.Replace().SetUpsert(true)).
			Return(nil, nil).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		err := repo.UpsertPeginConfiguration(context.Background(), *peginTestConfig)
		require.NoError(t, err)
	})
	t.Run("Db error upserting pegin configuration", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("ReplaceOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()
		err := repo.UpsertPeginConfiguration(context.Background(), *peginTestConfig)
		require.Error(t, err)
	})
}

func TestLpMongoRepository_UpsertPegoutConfiguration(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	configName := mongo.ConfigurationName("pegout")
	filter := bson.D{primitive.E{Key: "name", Value: configName}}
	t.Run("pegout configuration upserted successfully", func(t *testing.T) {
		const expectedLog = "INSERT interaction with db: {Signed:{Value:{TimeForDeposit:1 ExpireTime:2 PenaltyFee:3 CallFee:4 MaxValue:5 MinValue:6 ExpireBlocks:7 BridgeTransactionMin:8} Signature:pegout signature Hash:pegout hash} Name:pegout}"
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("ReplaceOne", mock.Anything, filter, mongo.StoredConfiguration[liquidity_provider.PegoutConfiguration]{
			Signed: *pegoutTestConfig,
			Name:   configName,
		}, options.Replace().SetUpsert(true)).
			Return(nil, nil).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		err := repo.UpsertPegoutConfiguration(context.Background(), *pegoutTestConfig)
		require.NoError(t, err)
	})
	t.Run("Db error upserting pegout configuration", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("ReplaceOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()
		err := repo.UpsertPegoutConfiguration(context.Background(), *pegoutTestConfig)
		require.Error(t, err)
	})
}

func TestLpMongoRepository_UpsertGeneralConfiguration(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	configName := mongo.ConfigurationName("general")
	filter := bson.D{primitive.E{Key: "name", Value: configName}}
	t.Run("general configuration upserted successfully", func(t *testing.T) {
		const expectedLog = "INSERT interaction with db: {Signed:{Value:{RskConfirmations:map[1:2 3:4] BtcConfirmations:map[5:6 7:8] PublicLiquidityCheck:false} Signature:general signature Hash:general hash} Name:general}"
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("ReplaceOne", mock.Anything, filter, mongo.StoredConfiguration[liquidity_provider.GeneralConfiguration]{
			Signed: *generalTestConfig,
			Name:   configName,
		}, options.Replace().SetUpsert(true)).
			Return(nil, nil).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		err := repo.UpsertGeneralConfiguration(context.Background(), *generalTestConfig)
		require.NoError(t, err)
	})
	t.Run("Db error upserting general configuration", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("ReplaceOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()
		err := repo.UpsertGeneralConfiguration(context.Background(), *generalTestConfig)
		require.Error(t, err)
	})
}

func TestLpMongoRepository_UpsertCredentials(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	configName := mongo.ConfigurationName("credentials")
	filter := bson.D{primitive.E{Key: "name", Value: configName}}
	t.Run("credentials upserted successfully", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("ReplaceOne", mock.Anything, filter, mongo.StoredConfiguration[liquidity_provider.HashedCredentials]{
			Signed: *testCredentials,
			Name:   configName,
		}, options.Replace().SetUpsert(true)).
			Return(nil, nil).Once()
		defer test.AssertNoLog(t)()
		err := repo.UpsertCredentials(context.Background(), *testCredentials)
		require.NoError(t, err)
	})
	t.Run("Db error upserting credentials", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.LiquidityProviderCollection)
		repo := mongo.NewLiquidityProviderRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("ReplaceOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()
		err := repo.UpsertCredentials(context.Background(), *testCredentials)
		require.Error(t, err)
	})
}
