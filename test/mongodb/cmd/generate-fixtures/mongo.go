package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
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
	port, err := envOrUint("MONGODB_PORT", 27018)
	if err != nil {
		return mongoConfig{}, fmt.Errorf("invalid MONGODB_PORT: %w", err)
	}
	return mongoConfig{
		host:     envOr("MONGODB_HOST", "localhost"),
		port:     port,
		username: envOr("MONGODB_USER", "test"),
		password: envOr("MONGODB_PASSWORD", "test"),
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
		if _, err := db.Collection(c.collection).DeleteMany(ctx, bson.M{}); err != nil {
			return fmt.Errorf("reset collection %s: %w", c.collection, err)
		}
	}
	return nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envOrUint(key string, fallback uint) (uint, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}
