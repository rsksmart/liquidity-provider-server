package mongo_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	mongoDb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestBatchPegOutMongoRepository_UpsertBatch(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	batch := rootstock.BatchPegOut{
		TransactionHash:    test.AnyHash,
		BlockHash:          test.AnyRskAddress,
		BlockNumber:        5,
		BtcTxHash:          test.AnyString,
		ReleaseRskTxHashes: []string{"0x1234", "0xabcd"},
	}
	t.Run("should upsert batch successfully", func(t *testing.T) {
		expectedLog := "UPSERT interaction with db: {TransactionHash:d8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb BlockHash:cd22 BlockNumber:790 BtcTxHash:ab11 ReleaseRskTxHashes:[0x1122 0x3344]}"
		client, collection := getClientAndCollectionMocks(mongo.BatchPegOutEventsCollection)
		newBatch := batch
		newBatch.BtcTxHash = "ab11"
		newBatch.BlockHash = "cd22"
		newBatch.BlockNumber = 790
		newBatch.ReleaseRskTxHashes = []string{"0x1122", "0x3344"}
		collection.On("ReplaceOne", mock.Anything,
			bson.M{"transaction_hash": batch.TransactionHash},
			newBatch,
			options.Replace().SetUpsert(true),
		).Return(&mongoDb.UpdateResult{ModifiedCount: 1}, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewBatchPegOutMongoRepository(conn)
		defer test.AssertLogContains(t, expectedLog)()
		err := repo.UpsertBatch(context.Background(), newBatch)
		collection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("should handle error when upserting batch", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.BatchPegOutEventsCollection)
		collection.On("ReplaceOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewBatchPegOutMongoRepository(conn)
		err := repo.UpsertBatch(context.Background(), batch)
		collection.AssertExpectations(t)
		require.Error(t, err)
	})
	t.Run("should handle error when upserting more than one batch", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.BatchPegOutEventsCollection)
		collection.On("ReplaceOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(&mongoDb.UpdateResult{ModifiedCount: 2}, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewBatchPegOutMongoRepository(conn)
		err := repo.UpsertBatch(context.Background(), batch)
		collection.AssertExpectations(t)
		require.ErrorContains(t, err, "multiple batch pegouts updated")
	})
}
