package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo/migrations"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(ctx context.Context, connectTimeout time.Duration, username, password, host string, port uint, runMigrations bool) (*mongo.Client, error) {
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

	if client, err = mongo.Connect(clientOptions); err != nil {
		return nil, err
	}
	db := client.Database(DbName)
	if runMigrations {
		if err = migrations.NewRunner(migrations.NewMongoDatabaseAdapter(db)).RunAll(ctx); err != nil {
			return nil, err
		}
	}
	if err = createIndexes(ctx, db); err != nil {
		return nil, err
	}
	return client, nil
}

// indexCreator applies an index model to a collection, so index definitions
// can be asserted without a live MongoDB.
type indexCreator interface {
	CreateIndex(ctx context.Context, collection string, model mongo.IndexModel) error
}

type databaseIndexCreator struct {
	db *mongo.Database
}

func (creator databaseIndexCreator) CreateIndex(ctx context.Context, collection string, model mongo.IndexModel) error {
	_, err := creator.db.Collection(collection).Indexes().CreateOne(ctx, model)
	return err
}

func pegInAddressRegistryWatchRskAddressIndex() mongo.IndexModel {
	return mongo.IndexModel{
		Keys:    bson.D{{Key: "rsk_address", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
}

func ensurePegInAddressRegistryWatchRskAddressIndex(ctx context.Context, creator indexCreator) error {
	if err := creator.CreateIndex(ctx, PegInAddressRegistryWatchCollection, pegInAddressRegistryWatchRskAddressIndex()); err != nil {
		return fmt.Errorf("error creating unique index on %s rsk_address: %w", PegInAddressRegistryWatchCollection, err)
	}
	return nil
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	if err := ensurePegInAddressRegistryWatchRskAddressIndex(ctx, databaseIndexCreator{db: db}); err != nil {
		return err
	}

	uniqueIndexes := []struct {
		collection string
		field      string
	}{
		{collection: DepositEventsCollection, field: "tx_hash"},
		{collection: TrustedAccountCollection, field: "address"},
		{collection: BatchPegOutEventsCollection, field: "transaction_hash"},
	}
	for _, idx := range uniqueIndexes {
		if err := createUniqueIndex(ctx, db, idx.collection, idx.field); err != nil {
			return fmt.Errorf("error creating unique index on %s.%s: %w", idx.collection, idx.field, err)
		}
		log.Infof("Created unique index on %s.%s", idx.collection, idx.field)
	}
	nonUniqueIndexes := []struct {
		collection string
		field      string
	}{
		{collection: RetainedPegoutQuoteCollection, field: "bridge_rebalances.tx_hash"},
		{collection: RetainedPeginQuoteCollection, field: "state"},
		{collection: RetainedPegoutQuoteCollection, field: "state"},
		{collection: PegInAddressRegistryWatchCollection, field: "state"},
		// agreement_timestamp is a Unix-seconds quote-creation time used by reports as a
		// range filter — two quotes issued in the same second is a legitimate insert.
		{collection: PeginQuoteCollection, field: "agreement_timestamp"},
		{collection: PegoutQuoteCollection, field: "agreement_timestamp"},
		// quote_hash is semantically one-to-one with a quote, but InsertRetainedQuote does
		// not enforce uniqueness; a unique build would fail and block startup on any pre-
		// existing duplicate, with no read-perf gain.
		{collection: RetainedPeginQuoteCollection, field: "quote_hash"},
		{collection: RetainedPegoutQuoteCollection, field: "quote_hash"},
		// hash is the quote's content hash and effectively its PK, but InsertQuote performs a
		// plain insert with no duplicate check; a unique build would block startup on any
		// pre-existing duplicate (legacy data, replayed insert) with no read-perf gain.
		{collection: PeginQuoteCollection, field: "hash"},
		{collection: PegoutQuoteCollection, field: "hash"},
	}
	for _, idx := range nonUniqueIndexes {
		if err := createIndex(ctx, db, idx.collection, idx.field); err != nil {
			return fmt.Errorf("error creating index on %s.%s: %w", idx.collection, idx.field, err)
		}
		log.Infof("Created index on %s.%s", idx.collection, idx.field)
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

func createIndex(ctx context.Context, db *mongo.Database, collectionName, field string) error {
	_, err := db.Collection(collectionName).Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{{Key: field, Value: 1}},
		},
	)
	return err
}
