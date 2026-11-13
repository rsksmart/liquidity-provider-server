package mongo_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDb "go.mongodb.org/mongo-driver/v2/mongo"
)

func newTestPegInClaim() rootstock.PegInClaim {
	now := time.Now().UTC().Truncate(time.Millisecond)
	return rootstock.PegInClaim{
		RskAddress:  "0x79568c2989232dca1840087d73d403602364c0d4",
		DepositTxID: "9c2d9d5ab8b7e3fe1c2f2e70ad0d77f0e3c5a3e8f1b7a4e6d2c0f9e8b7a6c5d4",
		BtcAddress:  "2N4bCrVMPu6ixHvxRS7qXcpRkdFh6WW4wEz",
		PegInID:     "6c0c1c2ba8a9cd4f0a3f0e15cc1d0a1b2c3d4e5f60718293a4b5c6d7e8f90a1b",
		State:       rootstock.PegInClaimSubmitting,
		Amount:      entities.NewWei(5000000000000000),
		Fee:         entities.NewWei(100000000000000),
		GasFee:      entities.NewWei(20000000000000),
		Value:       entities.NewWei(4880000000000000),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TestPegInClaimMongoRepository_Add(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	claim := newTestPegInClaim()
	retryableFilter := bson.M{
		"rsk_address":  claim.RskAddress,
		"deposit_txid": claim.DepositTxID,
		"state":        rootstock.PegInClaimRetryableFailure,
	}

	t.Run("stores the claim, replacing only a retryable failure", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().ReplaceOne(mock.Anything, retryableFilter, claim, withUpsert()).
			Return(&mongoDb.UpdateResult{UpsertedCount: 1}, nil).Once()
		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))

		defer assertDbInteractionLog(t, mongo.Upsert)()
		require.NoError(t, repo.Add(context.Background(), claim))
		collection.AssertExpectations(t)
	})

	t.Run("returns PegInClaimExistsError when a claim that is not a retryable failure is stored", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().ReplaceOne(mock.Anything, retryableFilter, claim, withUpsert()).
			Return(nil, mongoDb.WriteException{WriteErrors: []mongoDb.WriteError{{Code: 11000}}}).Once()
		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))

		require.ErrorIs(t, repo.Add(context.Background(), claim), rootstock.PegInClaimExistsError)
		collection.AssertExpectations(t)
	})

	t.Run("returns the ReplaceOne error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().ReplaceOne(mock.Anything, retryableFilter, claim, withUpsert()).
			Return(nil, assert.AnError).Once()
		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))

		require.ErrorIs(t, repo.Add(context.Background(), claim), assert.AnError)
		collection.AssertExpectations(t)
	})
}

func TestPegInClaimMongoRepository_Get(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	claim := newTestPegInClaim()
	identity := bson.M{"rsk_address": claim.RskAddress, "deposit_txid": claim.DepositTxID}

	t.Run("returns the claim stored for the deposit", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().FindOne(mock.Anything, identity).
			Return(mongoDb.NewSingleResultFromDocument(claim, nil, nil)).Once()
		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))

		defer assertDbInteractionLog(t, mongo.Read)()
		result, err := repo.Get(context.Background(), claim.RskAddress, claim.DepositTxID)
		require.NoError(t, err)
		assert.Equal(t, &claim, result)
		collection.AssertExpectations(t)
	})

	t.Run("returns nil when the deposit has no claim", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().FindOne(mock.Anything, identity).
			Return(mongoDb.NewSingleResultFromDocument(rootstock.PegInClaim{}, mongoDb.ErrNoDocuments, nil)).Once()
		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))

		result, err := repo.Get(context.Background(), claim.RskAddress, claim.DepositTxID)
		require.NoError(t, err)
		assert.Nil(t, result)
		collection.AssertExpectations(t)
	})

	t.Run("returns the FindOne error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().FindOne(mock.Anything, identity).
			Return(mongoDb.NewSingleResultFromDocument(rootstock.PegInClaim{}, assert.AnError, nil)).Once()
		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))

		result, err := repo.Get(context.Background(), claim.RskAddress, claim.DepositTxID)
		require.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, result)
		collection.AssertExpectations(t)
	})
}

func TestPegInClaimMongoRepository_Update(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	claim := newTestPegInClaim()
	claim.State = rootstock.PegInClaimClaimed
	claim.TxHash = "0x3f4c5d6e7f8091a2b3c4d5e6f708192a3b4c5d6e7f8091a2b3c4d5e6f708192a"
	claim.CallSuccess = true
	identity := bson.M{"rsk_address": claim.RskAddress, "deposit_txid": claim.DepositTxID}

	t.Run("replaces the stored claim", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().ReplaceOne(mock.Anything, identity, claim).
			Return(&mongoDb.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil).Once()
		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))

		defer assertDbInteractionLog(t, mongo.Update)()
		require.NoError(t, repo.Update(context.Background(), claim))
		collection.AssertExpectations(t)
	})

	t.Run("returns PegInClaimNotFoundError when no claim matches", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().ReplaceOne(mock.Anything, identity, claim).
			Return(&mongoDb.UpdateResult{MatchedCount: 0}, nil).Once()
		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))

		require.ErrorIs(t, repo.Update(context.Background(), claim), rootstock.PegInClaimNotFoundError)
		collection.AssertExpectations(t)
	})

	t.Run("returns the ReplaceOne error", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PegInClaimCollection)
		collection.EXPECT().ReplaceOne(mock.Anything, identity, claim).
			Return(nil, assert.AnError).Once()
		repo := mongo.NewPegInClaimMongoRepository(mongo.NewConnection(client, time.Second))

		require.ErrorIs(t, repo.Update(context.Background(), claim), assert.AnError)
		collection.AssertExpectations(t)
	})
}
