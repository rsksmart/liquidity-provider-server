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
	indexes := []struct {
		collection string
		field      string
	}{
		{collection: DepositEventsCollection, field: "tx_hash"},
		{collection: TrustedAccountCollection, field: "address"},
		{collection: BatchPegOutEventsCollection, field: "transaction_hash"},
	}
	for _, idx := range indexes {
		if err := createUniqueIndex(ctx, db, idx.collection, idx.field); err != nil {
			return fmt.Errorf("error creating unique index on %s.%s: %w", idx.collection, idx.field, err)
		}
		log.Infof("Created unique index on %s.%s", idx.collection, idx.field)
	}

	return nil
}

func createUniqueIndex(ctx context.Context, db *mongo.Database, collectionName, field string) error {
	_, err := db.Collection(collectionName).Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: field, Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	return err
}
