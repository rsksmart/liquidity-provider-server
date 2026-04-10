package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
)

const (
	defaultMongoTimeout   = 10 * time.Second
	defaultCommandTimeout = 60 * time.Second
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := loadConfigFromEnv()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultCommandTimeout)
	defer cancel()

	client, err := connectMongo(ctx, cfg)
	if err != nil {
		return err
	}
	defer disconnectMongo(client)

	db := registry.NewDatabaseRegistry(
		mongo.NewConnection(mongo.NewClientWrapper(client), defaultMongoTimeout),
	)

	if err = resetCollections(ctx, client); err != nil {
		return fmt.Errorf("reset collections: %w", err)
	}
	if err = writeRepresentativeData(ctx, db); err != nil {
		return fmt.Errorf("write representative data: %w", err)
	}
	fixturesDir, err := exportFixturesAsExtJSON(ctx, client)
	if err != nil {
		return fmt.Errorf("export fixtures: %w", err)
	}

	fmt.Println("\nDone. Fixtures written to", fixturesDir)
	return nil
}
