package regtest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	mongoURI = "mongodb://root:root@127.0.0.1:27017/admin?directConnection=true"
)

type MongoStack struct {
	client *mongoDriver.Client
	db     *mongoDriver.Database
}

func OpenMongo(t *testing.T) *MongoStack {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongoDriver.Connect(options.Client().ApplyURI(mongoURI))
	require.NoError(t, err)
	require.NoError(t, client.Ping(ctx, nil))
	t.Cleanup(func() {
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer disconnectCancel()
		require.NoError(t, client.Disconnect(disconnectCtx))
	})
	return &MongoStack{
		client: client,
		db:     client.Database(mongo.DbName),
	}
}

func (stack *MongoStack) GetWatch(t *testing.T, user common.Address) *rootstock.PegInWatch {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var watch rootstock.PegInWatch
	err := stack.db.Collection(mongo.PegInWatchCollection).FindOne(ctx, bson.M{
		"rsk_address": user.Hex(),
	}).Decode(&watch)
	if errors.Is(err, mongoDriver.ErrNoDocuments) {
		return nil
	}
	require.NoError(t, err)
	return &watch
}

func (stack *MongoStack) ListWatches(t *testing.T) []rootstock.PegInWatch {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := stack.db.Collection(mongo.PegInWatchCollection).Find(ctx, bson.M{})
	require.NoError(t, err)
	defer cursor.Close(ctx)
	var watches []rootstock.PegInWatch
	require.NoError(t, cursor.All(ctx, &watches))
	return watches
}

func (stack *MongoStack) CountWatches(t *testing.T, user common.Address) int {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := stack.db.Collection(mongo.PegInWatchCollection).CountDocuments(ctx, bson.M{
		"rsk_address": user.Hex(),
	})
	require.NoError(t, err)
	return int(count)
}

func (stack *MongoStack) GetClaim(t *testing.T, user common.Address, depositTxID string) *rootstock.PegInClaim {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var claim rootstock.PegInClaim
	err := stack.db.Collection(mongo.PegInClaimCollection).FindOne(ctx, bson.M{
		"rsk_address":  user.Hex(),
		"deposit_txid": depositTxID,
	}).Decode(&claim)
	if errors.Is(err, mongoDriver.ErrNoDocuments) {
		return nil
	}
	require.NoError(t, err)
	return &claim
}

func (stack *MongoStack) InsertSubmittingEmptyTxHash(
	t *testing.T,
	user common.Address,
	depositTxID string,
	btcAddr string,
) {
	t.Helper()
	now := time.Now().UTC()
	claim := rootstock.PegInClaim{
		RskAddress:  user.Hex(),
		DepositTxID: depositTxID,
		BtcAddress:  btcAddr,
		State:       rootstock.PegInClaimSubmitting,
		TxHash:      "",
		PegInID:     "",
		ReservedWei: entities.NewWei(0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := stack.db.Collection(mongo.PegInClaimCollection).InsertOne(ctx, claim)
	require.NoError(t, err)
}
