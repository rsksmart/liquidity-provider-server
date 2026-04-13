package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/test/mongodb/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type mongoConfig struct {
	host     string
	port     uint
	username string
	password string
}

func loadConfigFromEnv() (mongoConfig, error) {
	port, err := utils.EnvOrUint("MONGODB_PORT", 27018)
	if err != nil {
		return mongoConfig{}, fmt.Errorf("invalid MONGODB_PORT: %w", err)
	}
	return mongoConfig{
		host:     utils.EnvOr("MONGODB_HOST", "localhost"),
		port:     port,
		username: utils.EnvOr("MONGODB_USER", "test"),
		password: utils.EnvOr("MONGODB_PASSWORD", "test"),
	}, nil
}

func connectMongo(ctx context.Context, cfg mongoConfig) (*mongodriver.Client, error) {
	client, err := mongo.Connect(ctx, defaultMongoTimeout, cfg.username, cfg.password, cfg.host, cfg.port)
	if err != nil {
		return nil, fmt.Errorf("connect to mongo: %w", err)
	}
	return client, nil
}

func disconnectMongo(client *mongodriver.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultMongoTimeout)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "disconnect mongo:", err)
	}
}

func resetCollections(ctx context.Context, client *mongodriver.Client) error {
	db := client.Database(mongo.DbName)
	for _, c := range fixtureCollections {
		if _, err := db.Collection(c.Collection).DeleteMany(ctx, bson.M{}); err != nil {
			return fmt.Errorf("reset collection %s: %w", c.Collection, err)
		}
	}
	return nil
}
