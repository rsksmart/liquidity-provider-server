package mongo

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, connectTimeout time.Duration, username, password, host string, port uint) (*mongo.Client, error) {
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
	_, depositErr := db.Collection(DepositEventsCollection).Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "tx_hash", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if depositErr != nil {
		return depositErr
	}
	_, trustedAccountErr := db.Collection(TrustedAccountCollection).Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "address", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	return trustedAccountErr
}
