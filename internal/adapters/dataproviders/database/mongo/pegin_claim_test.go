package mongo_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDb "go.mongodb.org/mongo-driver/v2/mongo"
)

func samplePegInClaim(state rootstock.PegInClaimState) rootstock.PegInClaim {
	now := time.Now().UTC().Truncate(time.Millisecond)
	return rootstock.PegInClaim{
		RskAddress:  "0xabcd",
		DepositTxID: "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899",
		BtcAddress:  "bcrt1qexample",
		State:       state,
		ReservedWei: entities.NewWei(1_000_000_000_000_000_000),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func pegInClaimIdentity(claim rootstock.PegInClaim) bson.M {
	return bson.M{
		"rsk_address":  claim.RskAddress,
		"deposit_txid": claim.DepositTxID,
	}
}

func TestPegInClaimMongoRepository(t *testing.T) {
	claim := samplePegInClaim(rootstock.PegInClaimCandidate)
	identity := pegInClaimIdentity(claim)

	t.Run("inserts a claim", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().InsertOne(mock.Anything, claim).
			Return(&mongoDb.InsertOneResult{InsertedID: identity}, nil).Once()

		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
		require.NoError(t, repo.Insert(context.Background(), claim))
		collection.AssertExpectations(t)
	})

	t.Run("returns a unique-key conflict on duplicate insert", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().InsertOne(mock.Anything, claim).
			Return(nil, mongoDb.WriteException{WriteErrors: []mongoDb.WriteError{{Code: 11000}}}).Once()

		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
		err := repo.Insert(context.Background(), claim)
		require.ErrorIs(t, err, rootstock.ErrPegInClaimAlreadyExists)
		collection.AssertExpectations(t)
	})

	t.Run("returns a write error on insert", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().InsertOne(mock.Anything, claim).Return(nil, assert.AnError).Once()

		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
		err := repo.Insert(context.Background(), claim)
		require.Error(t, err)
		require.NotErrorIs(t, err, rootstock.ErrPegInClaimAlreadyExists)
		collection.AssertExpectations(t)
	})
}

func TestPegInClaimMongoRepository_Get(t *testing.T) {
	claim := samplePegInClaim(rootstock.PegInClaimCandidate)
	identity := pegInClaimIdentity(claim)

	t.Run("loads a claim by unique key", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().FindOne(mock.Anything, identity).
			Return(mongoDb.NewSingleResultFromDocument(claim, nil, nil)).Once()

		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
		result, err := repo.Get(context.Background(), claim.RskAddress, claim.DepositTxID)
		require.NoError(t, err)
		assert.Equal(t, &claim, result)
		collection.AssertExpectations(t)
	})

	t.Run("returns nil when the unique key is missing", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().FindOne(mock.Anything, identity).
			Return(mongoDb.NewSingleResultFromDocument(rootstock.PegInClaim{}, mongoDb.ErrNoDocuments, nil)).Once()

		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
		result, err := repo.Get(context.Background(), claim.RskAddress, claim.DepositTxID)
		require.NoError(t, err)
		assert.Nil(t, result)
		collection.AssertExpectations(t)
	})

	t.Run("returns a read error on get", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().FindOne(mock.Anything, identity).
			Return(mongoDb.NewSingleResultFromDocument(rootstock.PegInClaim{}, assert.AnError, nil)).Once()

		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
		result, err := repo.Get(context.Background(), claim.RskAddress, claim.DepositTxID)
		require.Error(t, err)
		assert.Nil(t, result)
		collection.AssertExpectations(t)
	})
}

func TestPegInClaimMongoRepository_Update(t *testing.T) {
	claim := samplePegInClaim(rootstock.PegInClaimSubmitting)
	claim.TxHash = "0xabc"
	identity := pegInClaimIdentity(claim)

	t.Run("replaces a matching claim", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().ReplaceOne(mock.Anything, identity, claim).
			Return(&mongoDb.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil).Once()

		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
		require.NoError(t, repo.Update(context.Background(), claim))
		collection.AssertExpectations(t)
	})

	t.Run("returns a write error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().ReplaceOne(mock.Anything, mock.Anything, claim).
			Return(nil, assert.AnError).Once()

		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
		err := repo.Update(context.Background(), claim)
		require.Error(t, err)
		require.NotErrorIs(t, err, rootstock.ErrPegInClaimNotFound)
		collection.AssertExpectations(t)
	})

	t.Run("returns not found when no document matches", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().ReplaceOne(mock.Anything, mock.Anything, claim).
			Return(&mongoDb.UpdateResult{MatchedCount: 0}, nil).Once()

		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
		err := repo.Update(context.Background(), claim)
		require.ErrorIs(t, err, rootstock.ErrPegInClaimNotFound)
		collection.AssertExpectations(t)
	})
}

func TestPegInClaimMongoRepository_ListByStates(t *testing.T) {
	claim := samplePegInClaim(rootstock.PegInClaimCandidate)

	t.Run("filters by states", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().Find(
			mock.Anything,
			bson.M{"state": bson.M{"$in": []rootstock.PegInClaimState{
				rootstock.PegInClaimCandidate,
				rootstock.PegInClaimSubmitting,
			}}},
		).Return(mongoDb.NewCursorFromDocuments([]any{claim}, nil, nil)).Once()

		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
		result, err := repo.ListByStates(
			context.Background(),
			rootstock.PegInClaimCandidate,
			rootstock.PegInClaimSubmitting,
		)
		require.NoError(t, err)
		assert.Equal(t, []rootstock.PegInClaim{claim}, result)
		collection.AssertExpectations(t)
	})

	t.Run("returns a find error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().Find(mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
		result, err := repo.ListByStates(context.Background(), rootstock.PegInClaimSubmitting)
		require.Error(t, err)
		assert.Nil(t, result)
		collection.AssertExpectations(t)
	})
}

func TestPegInClaimMongoRepository_InsertWriteError(t *testing.T) {
	claim := rootstock.PegInClaim{RskAddress: "0xabcd", DepositTxID: "aabb"}
	client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
	collection.EXPECT().InsertOne(mock.Anything, claim).Return(nil, assert.AnError).Once()

	repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
	err := repo.Insert(context.Background(), claim)
	require.Error(t, err)
	require.NotErrorIs(t, err, rootstock.ErrPegInClaimAlreadyExists)
	collection.AssertExpectations(t)
}

func TestPegInClaimMongoRepository_GetMissing(t *testing.T) {
	identity := bson.M{"rsk_address": "0xabcd", "deposit_txid": "aabb"}
	client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
	collection.EXPECT().FindOne(mock.Anything, identity).
		Return(mongoDb.NewSingleResultFromDocument(rootstock.PegInClaim{}, mongoDb.ErrNoDocuments, nil)).Once()

	repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
	result, err := repo.Get(context.Background(), "0xabcd", "aabb")
	require.NoError(t, err)
	assert.Nil(t, result)
	collection.AssertExpectations(t)
}

func TestPegInClaimMongoRepository_GetError(t *testing.T) {
	identity := bson.M{"rsk_address": "0xabcd", "deposit_txid": "aabb"}
	client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
	collection.EXPECT().FindOne(mock.Anything, identity).
		Return(mongoDb.NewSingleResultFromDocument(rootstock.PegInClaim{}, assert.AnError, nil)).Once()

	repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
	result, err := repo.Get(context.Background(), "0xabcd", "aabb")
	require.Error(t, err)
	assert.Nil(t, result)
	collection.AssertExpectations(t)
}

func TestPegInClaimMongoRepository_Update(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond)
	claim := rootstock.PegInClaim{
		RskAddress:  "0xabcd",
		DepositTxID: "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899",
		BtcAddress:  "bcrt1qexample",
		State:       rootstock.PegInClaimSubmitting,
		TxHash:      "0xabc",
		ReservedWei: entities.NewWei(1_000_000_000_000_000_000),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	identity := bson.M{
		"rsk_address":  claim.RskAddress,
		"deposit_txid": claim.DepositTxID,
	}
	client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
	collection.EXPECT().ReplaceOne(mock.Anything, identity, claim).
		Return(&mongoDb.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil).Once()

	repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
	require.NoError(t, repo.Update(context.Background(), claim))
	collection.AssertExpectations(t)
}

func TestPegInClaimMongoRepository_UpdateError(t *testing.T) {
	claim := rootstock.PegInClaim{RskAddress: "0xabcd", DepositTxID: "aabb"}
	client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
	collection.EXPECT().ReplaceOne(mock.Anything, mock.Anything, claim).
		Return(nil, assert.AnError).Once()

	repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
	err := repo.Update(context.Background(), claim)
	require.Error(t, err)
	require.NotErrorIs(t, err, rootstock.ErrPegInClaimNotFound)
	collection.AssertExpectations(t)
}

func TestPegInClaimMongoRepository_UpdateNotFound(t *testing.T) {
	claim := rootstock.PegInClaim{RskAddress: "0xabcd", DepositTxID: "aabb"}
	client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
	collection.EXPECT().ReplaceOne(mock.Anything, mock.Anything, claim).
		Return(&mongoDb.UpdateResult{MatchedCount: 0}, nil).Once()

	repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
	err := repo.Update(context.Background(), claim)
	require.ErrorIs(t, err, rootstock.ErrPegInClaimNotFound)
	collection.AssertExpectations(t)
}

func TestPegInClaimMongoRepository_ListByStates(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond)
	claim := rootstock.PegInClaim{
		RskAddress:  "0xabcd",
		DepositTxID: "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899",
		BtcAddress:  "bcrt1qexample",
		State:       rootstock.PegInClaimCandidate,
		ReservedWei: entities.NewWei(1_000_000_000_000_000_000),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
	collection.EXPECT().Find(
		mock.Anything,
		bson.M{"state": bson.M{"$in": []rootstock.PegInClaimState{
			rootstock.PegInClaimCandidate,
			rootstock.PegInClaimSubmitting,
		}}},
	).Return(mongoDb.NewCursorFromDocuments([]any{claim}, nil, nil)).Once()

	repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
	result, err := repo.ListByStates(
		context.Background(),
		rootstock.PegInClaimCandidate,
		rootstock.PegInClaimSubmitting,
	)
	require.NoError(t, err)
	assert.Equal(t, []rootstock.PegInClaim{claim}, result)
	collection.AssertExpectations(t)
}

func TestPegInClaimMongoRepository_ListByStatesFindError(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
	collection.EXPECT().Find(mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

	repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
	result, err := repo.ListByStates(context.Background(), rootstock.PegInClaimSubmitting)
	require.Error(t, err)
	assert.Nil(t, result)
	collection.AssertExpectations(t)
}

func TestPegInClaimMongoRepository_ListByStatesEmptyFilter(t *testing.T) {
	client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
	collection.EXPECT().Find(mock.Anything, bson.M{}).
		Return(mongoDb.NewCursorFromDocuments([]any{}, nil, nil)).Once()

	repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
	result, err := repo.ListByStates(context.Background())
	require.NoError(t, err)
	assert.Empty(t, result)
	collection.AssertExpectations(t)
}
