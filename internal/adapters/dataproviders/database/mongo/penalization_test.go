package mongo_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	mongoDb "go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"testing"
	"time"
)

var testPenalization = penalization.PenalizedEvent{
	LiquidityProvider: "0x0000000000000000000000000000000000000000",
	QuoteHash:         "8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819",
	Penalty:           entities.NewWei(100),
}

func TestPenalizedEventMongoRepository_InsertPenalization(t *testing.T) {
	t.Run("Insert penalization successfully", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PenalizedEventCollection)
		collection.On("InsertOne", mock.Anything, mock.MatchedBy(func(q penalization.PenalizedEvent) bool {
			return q.QuoteHash == testPenalization.QuoteHash && reflect.TypeFor[penalization.PenalizedEvent]().NumField() == test.CountNonZeroValues(q)
		})).Return(nil, nil).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPenalizedEventRepository(conn)
		err := repo.InsertPenalization(context.Background(), testPenalization)
		collection.AssertExpectations(t)
		require.NoError(t, err)
	})
	t.Run("Db error inserting penalization", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PenalizedEventCollection)
		collection.On("InsertOne", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPenalizedEventRepository(conn)
		err := repo.InsertPenalization(context.Background(), testPenalization)
		collection.AssertExpectations(t)
		require.Error(t, err)
	})
}

func TestLpMongoRepository_GetPenalizationsByQuoteHashes(t *testing.T) {
	t.Run("Get penalizations by quote hashes", func(t *testing.T) {
		client, db := getClientAndDatabaseMocks()
		penalizationCollection := &mocks.CollectionBindingMock{}

		db.EXPECT().Collection(mongo.PenalizedEventCollection).Return(penalizationCollection)

		hashList := []string{"27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f"}
		expectedPenalizations := []penalization.PenalizedEvent{testPenalization}
		penalizationCollection.On("Find", mock.Anything, mock.MatchedBy(func(filter bson.M) bool {
			return true
		}), mock.Anything).Return(mongoDb.NewCursorFromDocuments([]any{testPenalization}, nil, nil))
		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPenalizedEventRepository(conn)

		result, err := repo.GetPenalizationsByQuoteHashes(context.Background(), hashList)

		require.NoError(t, err)
		assert.Equal(t, expectedPenalizations, result)

		penalizationCollection.AssertExpectations(t)
		penalizationCollection.AssertExpectations(t)
	})

	t.Run("error reading quotes from DB", func(t *testing.T) {
		client, collection := getClientAndCollectionMocks(mongo.PenalizedEventCollection)

		collection.On("Find", mock.Anything, mock.Anything).Return(nil, mongoDb.ErrNoDocuments).Once()

		conn := mongo.NewConnection(client, time.Duration(1))
		repo := mongo.NewPenalizedEventRepository(conn)

		hashList := []string{"27d70ec2bc2c3154dc9a5b53b118a755441b22bc1c8ccde967ed33609970c25f"}
		quotes, err := repo.GetPenalizationsByQuoteHashes(context.Background(), hashList)
		require.Error(t, err)
		assert.Equal(t, "mongo: no documents in result", err.Error())
		assert.Nil(t, quotes)
	})
}
