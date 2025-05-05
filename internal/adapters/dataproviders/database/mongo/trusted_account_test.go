package mongo_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	mongoDb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testAccount = liquidity_provider.TrustedAccountDetails{
	Address:          "0x1234567890abcdef1234567890abcdef12345678",
	Name:             "Test Account",
	Btc_locking_cap:  entities.NewWei(1000000000000000000),
	Rbtc_locking_cap: entities.NewWei(2000000000000000000),
}

var signedTestAccount = entities.Signed[liquidity_provider.TrustedAccountDetails]{
	Value:     testAccount,
	Signature: "signature",
	Hash:      "hash",
}

func TestLpMongoRepository_GetTrustedAccount(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	t.Run("trusted account found successfully", func(t *testing.T) {
		const expectedLog = "READ interaction with db: &{Address:0x1234567890abcdef1234567890abcdef12345678 Name:Test Account Btc_locking_cap:1000000000000000000 Rbtc_locking_cap:2000000000000000000}"
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		filter := bson.M{"address": testAccount.Address}
		collection.On("FindOne", mock.Anything, filter).
			Return(mongoDb.NewSingleResultFromDocument(testAccount, nil, nil)).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetTrustedAccount(context.Background(), testAccount.Address)
		require.NoError(t, err)
		assert.Equal(t, &testAccount, result)
	})
	t.Run("trusted account not found", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		filter := bson.M{"address": testAccount.Address}
		collection.On("FindOne", mock.Anything, filter).
			Return(mongoDb.NewSingleResultFromDocument(liquidity_provider.TrustedAccountDetails{}, mongoDb.ErrNoDocuments, nil)).Once()
		result, err := repo.GetTrustedAccount(context.Background(), testAccount.Address)
		require.Error(t, err)
		assert.Equal(t, liquidity_provider.ErrTrustedAccountNotFound, err)
		assert.Nil(t, result)
	})
	t.Run("Db error reading trusted account", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		filter := bson.M{"address": testAccount.Address}
		collection.On("FindOne", mock.Anything, filter).
			Return(mongoDb.NewSingleResultFromDocument(nil, assert.AnError, nil)).Once()
		result, err := repo.GetTrustedAccount(context.Background(), testAccount.Address)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestLpMongoRepository_GetAllTrustedAccounts(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	t.Run("all trusted accounts found successfully", func(t *testing.T) {
		const expectedLog = "READ interaction with db: [{Address:0x1234567890abcdef1234567890abcdef12345678 Name:Test Account Btc_locking_cap:1000000000000000000 Rbtc_locking_cap:2000000000000000000}]"
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		accounts := []liquidity_provider.TrustedAccountDetails{testAccount}
		accountsAny := make([]any, len(accounts))
		for i, account := range accounts {
			accountsAny[i] = account
		}
		collection.On("Find", mock.Anything, bson.M{}).
			Return(mongoDb.NewCursorFromDocuments(accountsAny, nil, nil)).Once()

		defer assertDbInteractionLog(t, expectedLog)()
		result, err := repo.GetAllTrustedAccounts(context.Background())
		require.NoError(t, err)
		assert.Equal(t, accounts, result)
	})
	t.Run("Db error finding trusted accounts", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("Find", mock.Anything, bson.M{}).Return(nil, assert.AnError).Once()
		result, err := repo.GetAllTrustedAccounts(context.Background())
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Db error handling", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("Find", mock.Anything, bson.M{}).Return(nil, assert.AnError).Once()
		result, err := repo.GetAllTrustedAccounts(context.Background())
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestLpMongoRepository_UpdateTrustedAccount(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	t.Run("trusted account updated successfully", func(t *testing.T) {
		const expectedLog = "UPDATE interaction with db: {Value:{Address:0x1234567890abcdef1234567890abcdef12345678 Name:Test Account Btc_locking_cap:1000000000000000000 Rbtc_locking_cap:2000000000000000000} Signature:signature Hash:hash}"
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, bson.M{"address": signedTestAccount.Value.Address}).
			Return(mongoDb.NewSingleResultFromDocument(&testAccount, nil, nil)).Once()
		filter := bson.M{"address": signedTestAccount.Value.Address}
		opts := options.Update()
		update := bson.M{"$set": signedTestAccount}
		collection.On("UpdateOne", mock.Anything, filter, update, opts).Return(&mongoDb.UpdateResult{}, nil).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		err := repo.UpdateTrustedAccount(context.Background(), signedTestAccount)
		require.NoError(t, err)
	})
	t.Run("trusted account not found", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, bson.M{"address": signedTestAccount.Value.Address}).
			Return(mongoDb.NewSingleResultFromDocument(liquidity_provider.TrustedAccountDetails{}, mongoDb.ErrNoDocuments, nil)).Once()
		err := repo.UpdateTrustedAccount(context.Background(), signedTestAccount)
		require.Error(t, err)
		assert.Equal(t, liquidity_provider.ErrTrustedAccountNotFound, err)
	})
	t.Run("Db error checking existing account", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, bson.M{"address": signedTestAccount.Value.Address}).
			Return(mongoDb.NewSingleResultFromDocument(nil, assert.AnError, nil)).Once()
		err := repo.UpdateTrustedAccount(context.Background(), signedTestAccount)
		require.Error(t, err)
	})
	t.Run("Db error updating account", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, bson.M{"address": signedTestAccount.Value.Address}).
			Return(mongoDb.NewSingleResultFromDocument(&testAccount, nil, nil)).Once()
		filter := bson.M{"address": signedTestAccount.Value.Address}
		opts := options.Update()
		update := bson.M{"$set": signedTestAccount}
		collection.On("UpdateOne", mock.Anything, filter, update, opts).Return(nil, assert.AnError).Once()
		err := repo.UpdateTrustedAccount(context.Background(), signedTestAccount)
		require.Error(t, err)
	})
}

func TestLpMongoRepository_AddTrustedAccount(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	t.Run("trusted account added successfully", func(t *testing.T) {
		const expectedLog = "INSERT interaction with db: {Value:{Address:0x1234567890abcdef1234567890abcdef12345678 Name:Test Account Btc_locking_cap:1000000000000000000 Rbtc_locking_cap:2000000000000000000} Signature:signature Hash:hash}"
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, bson.M{"address": signedTestAccount.Value.Address}).
			Return(mongoDb.NewSingleResultFromDocument(liquidity_provider.TrustedAccountDetails{}, mongoDb.ErrNoDocuments, nil)).Once()
		collection.On("InsertOne", mock.Anything, signedTestAccount).Return(&mongoDb.InsertOneResult{}, nil).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		err := repo.AddTrustedAccount(context.Background(), signedTestAccount)
		require.NoError(t, err)
	})
	t.Run("trusted account already exists", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, bson.M{"address": signedTestAccount.Value.Address}).
			Return(mongoDb.NewSingleResultFromDocument(&testAccount, nil, nil)).Once()
		err := repo.AddTrustedAccount(context.Background(), signedTestAccount)
		require.Error(t, err)
		assert.Equal(t, liquidity_provider.ErrDuplicateTrustedAccount, err)
	})
	t.Run("Db error checking existing account", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, bson.M{"address": signedTestAccount.Value.Address}).
			Return(mongoDb.NewSingleResultFromDocument(nil, assert.AnError, nil)).Once()
		err := repo.AddTrustedAccount(context.Background(), signedTestAccount)
		require.Error(t, err)
		collection.AssertNotCalled(t, "InsertOne")
	})
	t.Run("Db error inserting account", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		collection.On("FindOne", mock.Anything, bson.M{"address": signedTestAccount.Value.Address}).
			Return(mongoDb.NewSingleResultFromDocument(liquidity_provider.TrustedAccountDetails{}, mongoDb.ErrNoDocuments, nil)).Once()
		collection.On("InsertOne", mock.Anything, signedTestAccount).Return(nil, assert.AnError).Once()
		err := repo.AddTrustedAccount(context.Background(), signedTestAccount)
		require.Error(t, err)
	})
}

func TestLpMongoRepository_DeleteTrustedAccount(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	t.Run("trusted account deleted successfully", func(t *testing.T) {
		const expectedLog = "DELETE interaction with db: map[address:0x1234567890abcdef1234567890abcdef12345678]"
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		filter := bson.M{"address": testAccount.Address}
		result := &mongoDb.DeleteResult{DeletedCount: 1}
		collection.On("DeleteOne", mock.Anything, filter).Return(result, nil).Once()
		defer assertDbInteractionLog(t, expectedLog)()
		err := repo.DeleteTrustedAccount(context.Background(), testAccount.Address)
		require.NoError(t, err)
	})
	t.Run("trusted account not found", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		filter := bson.M{"address": testAccount.Address}
		result := &mongoDb.DeleteResult{DeletedCount: 0}
		collection.On("DeleteOne", mock.Anything, filter).Return(result, nil).Once()
		err := repo.DeleteTrustedAccount(context.Background(), testAccount.Address)
		require.Error(t, err)
		assert.Equal(t, liquidity_provider.ErrTrustedAccountNotFound, err)
	})
	t.Run("Db error deleting account", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))
		filter := bson.M{"address": testAccount.Address}
		collection.On("DeleteOne", mock.Anything, filter).Return(nil, assert.AnError).Once()
		err := repo.DeleteTrustedAccount(context.Background(), testAccount.Address)
		require.Error(t, err)
	})
}
