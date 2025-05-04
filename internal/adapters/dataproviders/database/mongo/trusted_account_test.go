package mongo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	mongoDb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testTrustedAccount = liquidity_provider.TrustedAccount{
	Address:        "0x1234567890abcdef",
	Name:           "Test Account",
	BtcLockingCap:  entities.NewWei(1000000),
	RbtcLockingCap: entities.NewWei(2000000),
}

var testTrustedAccount2 = liquidity_provider.TrustedAccount{
	Address:        "0xabcdef1234567890",
	Name:           "Test Account 2",
	BtcLockingCap:  entities.NewWei(3000000),
	RbtcLockingCap: entities.NewWei(4000000),
}

func TestTrustedAccountMongoRepository_GetTrustedAccount(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)

	t.Run("account found successfully", func(t *testing.T) {
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))

		filter := bson.M{"address": testTrustedAccount.Address}
		collection.On("FindOne", mock.Anything, filter).
			Return(mongoDb.NewSingleResultFromDocument(testTrustedAccount, nil, nil)).Once()

		result, err := repo.GetTrustedAccount(context.Background(), testTrustedAccount.Address)
		require.NoError(t, err)
		assert.Equal(t, &testTrustedAccount, result)
		collection.AssertExpectations(t)
	})

	t.Run("account not found", func(t *testing.T) {
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))

		filter := bson.M{"address": testTrustedAccount.Address}
		collection.On("FindOne", mock.Anything, filter).
			Return(mongoDb.NewSingleResultFromDocument(liquidity_provider.TrustedAccount{}, mongoDb.ErrNoDocuments, nil)).Once()

		result, err := repo.GetTrustedAccount(context.Background(), testTrustedAccount.Address)
		require.Equal(t, liquidity_provider.TrustedAccountNotFoundError, err)
		assert.Nil(t, result)
		collection.AssertExpectations(t)
	})

	t.Run("database error", func(t *testing.T) {
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))

		filter := bson.M{"address": testTrustedAccount.Address}
		collection.On("FindOne", mock.Anything, filter).
			Return(mongoDb.NewSingleResultFromDocument(nil, errors.New("document is nil"), nil)).Once()

		result, err := repo.GetTrustedAccount(context.Background(), testTrustedAccount.Address)
		require.Error(t, err)
		assert.Equal(t, "document is nil", err.Error())
		assert.Nil(t, result)
		collection.AssertExpectations(t)
	})
}

func TestTrustedAccountMongoRepository_AddTrustedAccount(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)

	t.Run("account added successfully", func(t *testing.T) {
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))

		findFilter := bson.M{"address": testTrustedAccount.Address}
		collection.On("FindOne", mock.Anything, findFilter).
			Return(mongoDb.NewSingleResultFromDocument(liquidity_provider.TrustedAccount{}, mongoDb.ErrNoDocuments, nil)).Once()

		collection.On("InsertOne", mock.Anything, testTrustedAccount).
			Return(&mongoDb.InsertOneResult{InsertedID: testTrustedAccount.Address}, nil).Once()

		err := repo.AddTrustedAccount(context.Background(), testTrustedAccount)
		require.NoError(t, err)
		collection.AssertExpectations(t)
	})

	t.Run("account already exists", func(t *testing.T) {
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))

		findFilter := bson.M{"address": testTrustedAccount.Address}
		collection.On("FindOne", mock.Anything, findFilter).
			Return(mongoDb.NewSingleResultFromDocument(testTrustedAccount, nil, nil)).Once()

		err := repo.AddTrustedAccount(context.Background(), testTrustedAccount)
		require.Error(t, err)
		assert.Equal(t, liquidity_provider.DuplicateAddressError, err)
		collection.AssertExpectations(t)
	})

	t.Run("database error during check", func(t *testing.T) {
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))

		dbError := errors.New("document is nil")
		findFilter := bson.M{"address": testTrustedAccount.Address}
		collection.On("FindOne", mock.Anything, findFilter).
			Return(mongoDb.NewSingleResultFromDocument(nil, dbError, nil)).Once()

		err := repo.AddTrustedAccount(context.Background(), testTrustedAccount)
		require.Error(t, err)
		assert.Equal(t, dbError, err)
		collection.AssertExpectations(t)
	})

	t.Run("database error during insert", func(t *testing.T) {
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))

		findFilter := bson.M{"address": testTrustedAccount.Address}
		collection.On("FindOne", mock.Anything, findFilter).
			Return(mongoDb.NewSingleResultFromDocument(liquidity_provider.TrustedAccount{}, mongoDb.ErrNoDocuments, nil)).Once()

		dbError := errors.New("insert error")
		collection.On("InsertOne", mock.Anything, testTrustedAccount).
			Return(nil, dbError).Once()

		err := repo.AddTrustedAccount(context.Background(), testTrustedAccount)
		require.Error(t, err)
		assert.Equal(t, dbError, err)
		collection.AssertExpectations(t)
	})
}

func TestTrustedAccountMongoRepository_UpdateTrustedAccount(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)

	t.Run("account exists and update successful", func(t *testing.T) {
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))

		findFilter := bson.M{"address": testTrustedAccount.Address}
		collection.On("FindOne", mock.Anything, findFilter).
			Return(mongoDb.NewSingleResultFromDocument(testTrustedAccount, nil, nil)).Once()

		updateFilter := bson.M{"address": testTrustedAccount.Address}
		updateSet := bson.M{"$set": testTrustedAccount}
		collection.On("UpdateOne", mock.Anything, updateFilter, updateSet, mock.Anything).
			Return(&mongoDb.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil).Once()

		err := repo.UpdateTrustedAccount(context.Background(), testTrustedAccount)
		require.NoError(t, err)
		collection.AssertExpectations(t)
	})

	t.Run("account not found but upsert successful", func(t *testing.T) {
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))

		findFilter := bson.M{"address": testTrustedAccount.Address}
		collection.On("FindOne", mock.Anything, findFilter).
			Return(mongoDb.NewSingleResultFromDocument(liquidity_provider.TrustedAccount{}, liquidity_provider.TrustedAccountNotFoundError, nil)).Once()

		updateFilter := bson.M{"address": testTrustedAccount.Address}
		updateSet := bson.M{"$set": testTrustedAccount}
		opts := options.Update().SetUpsert(true)
		collection.On("UpdateOne", mock.Anything, updateFilter, updateSet, opts).
			Return(&mongoDb.UpdateResult{UpsertedCount: 1}, nil).Once()

		err := repo.UpdateTrustedAccount(context.Background(), testTrustedAccount)
		require.NoError(t, err)
		collection.AssertExpectations(t)
	})

	t.Run("database error during get", func(t *testing.T) {
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))

		docNilError := errors.New("document is nil")
		findFilter := bson.M{"address": testTrustedAccount.Address}
		collection.On("FindOne", mock.Anything, findFilter).
			Return(mongoDb.NewSingleResultFromDocument(nil, docNilError, nil)).Once()

		err := repo.UpdateTrustedAccount(context.Background(), testTrustedAccount)
		require.Error(t, err)
		assert.Equal(t, docNilError.Error(), err.Error())
		collection.AssertExpectations(t)
	})

	t.Run("database error during update", func(t *testing.T) {
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))

		findFilter := bson.M{"address": testTrustedAccount.Address}
		collection.On("FindOne", mock.Anything, findFilter).
			Return(mongoDb.NewSingleResultFromDocument(testTrustedAccount, nil, nil)).Once()

		dbError := errors.New("database error")
		updateFilter := bson.M{"address": testTrustedAccount.Address}
		updateSet := bson.M{"$set": testTrustedAccount}
		opts := options.Update().SetUpsert(true)
		collection.On("UpdateOne", mock.Anything, updateFilter, updateSet, opts).
			Return(nil, dbError).Once()

		err := repo.UpdateTrustedAccount(context.Background(), testTrustedAccount)
		require.Error(t, err)
		assert.Equal(t, dbError, err)
		collection.AssertExpectations(t)
	})
}

func TestTrustedAccountMongoRepository_DeleteTrustedAccount(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.TrustedAccountCollection)

	t.Run("account deleted successfully", func(t *testing.T) {
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))

		filter := bson.M{"address": testTrustedAccount.Address}
		collection.On("DeleteOne", mock.Anything, filter).
			Return(&mongoDb.DeleteResult{DeletedCount: 1}, nil).Once()

		err := repo.DeleteTrustedAccount(context.Background(), testTrustedAccount.Address)
		require.NoError(t, err)
		collection.AssertExpectations(t)
	})

	t.Run("account not found", func(t *testing.T) {
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))

		filter := bson.M{"address": testTrustedAccount.Address}
		collection.On("DeleteOne", mock.Anything, filter).
			Return(&mongoDb.DeleteResult{DeletedCount: 0}, nil).Once()

		err := repo.DeleteTrustedAccount(context.Background(), testTrustedAccount.Address)
		require.Error(t, err)
		assert.Equal(t, liquidity_provider.TrustedAccountNotFoundError, err)
		collection.AssertExpectations(t)
	})

	t.Run("database error", func(t *testing.T) {
		repo := mongo.NewTrustedAccountRepository(mongo.NewConnection(client, time.Duration(1)))

		dbError := errors.New("database error")
		filter := bson.M{"address": testTrustedAccount.Address}
		collection.On("DeleteOne", mock.Anything, filter).
			Return(nil, dbError).Once()

		err := repo.DeleteTrustedAccount(context.Background(), testTrustedAccount.Address)
		require.Error(t, err)
		assert.Equal(t, dbError, err)
		collection.AssertExpectations(t)
	})
}
