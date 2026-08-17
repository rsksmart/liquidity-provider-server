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

func TestPegInClaimMongoRepository(t *testing.T) {
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
	identity := bson.M{
		"rsk_address":  claim.RskAddress,
		"deposit_txid": claim.DepositTxID,
	}

	t.Run("inserts a claim", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().InsertOne(mock.Anything, claim).
			Return(&mongoDb.InsertOneResult{InsertedID: identity}, nil).Once()

		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
		require.NoError(t, repo.Insert(context.Background(), claim))
		collection.AssertExpectations(t)
	})

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

	t.Run("returns a unique-key conflict on duplicate insert", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().InsertOne(mock.Anything, claim).
			Return(nil, mongoDb.WriteException{WriteErrors: []mongoDb.WriteError{{Code: 11000}}}).Once()

		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))
		err := repo.Insert(context.Background(), claim)
		require.ErrorIs(t, err, rootstock.ErrPegInClaimAlreadyExists)
		collection.AssertExpectations(t)
	})
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
