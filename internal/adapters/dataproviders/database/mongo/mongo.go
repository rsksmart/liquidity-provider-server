package mongo

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const connectTimeout = 10 * time.Second

func Connect(ctx context.Context, username, password, host string, port uint) (*mongo.Client, error) {
	var err error
	var client *mongo.Client
	log.Info("Connecting to MongoDB")
	clientOptions := options.Client().ApplyURI(
		fmt.Sprintf(
			"mongodb://%s:%s@%s:%d/admin",
			username, password, host, port,
		),
	)

	ctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	if client, err = mongo.Connect(ctx, clientOptions); err != nil {
		return nil, err
	}
	db := client.Database(DbName)
	if err = createIndexes(ctx, db); err == nil {
		return client, nil
	} else {
		return nil, err
	}
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection(DepositEventsCollection).Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "tx_hash", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	return err
}
