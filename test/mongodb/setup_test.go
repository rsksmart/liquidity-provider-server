//go:build integration

package mongodb_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

var (
	peginRepo   quote.PeginQuoteRepository
	pegoutRepo  quote.PegoutQuoteRepository
	lpRepo      liquidity_provider.LiquidityProviderRepository
	trustedRepo liquidity_provider.TrustedAccountRepository
	penaltyRepo penalization.PenalizedEventRepository
	batchRepo   rootstock.BatchPegOutRepository
	conn        *mongo.Connection
	mongoClient *mongoDriver.Client
	testDbName  = mongo.DbName
)

var allCollections = []string{
	mongo.PeginQuoteCollection,
	mongo.RetainedPeginQuoteCollection,
	mongo.PeginCreationDataCollection,
	mongo.PegoutQuoteCollection,
	mongo.RetainedPegoutQuoteCollection,
	mongo.PegoutCreationDataCollection,
	mongo.DepositEventsCollection,
	mongo.LiquidityProviderCollection,
	mongo.TrustedAccountCollection,
	mongo.PenalizedEventCollection,
	mongo.BatchPegOutEventsCollection,
}

func TestMain(m *testing.M) {
	host := envOr("MONGODB_HOST", "localhost")
	port := envOrUint("MONGODB_PORT", 27018)
	username := envOr("MONGODB_USER", "test")
	password := envOr("MONGODB_PASSWORD", "test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, 10*time.Second, username, password, host, port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to MongoDB: %v\n", err)
		os.Exit(1)
	}
	mongoClient = client

	wrapper := mongo.NewClientWrapper(client)
	conn = mongo.NewConnection(wrapper, 10*time.Second)

	db := registry.NewDatabaseRegistry(conn)
	peginRepo = db.PeginRepository
	pegoutRepo = db.PegoutRepository
	lpRepo = db.LiquidityProviderRepository
	trustedRepo = db.TrustedAccountRepository
	penaltyRepo = db.PenalizedEventRepository
	batchRepo = db.BatchPegOutRepository

	code := m.Run()

	closeChannel := make(chan bool, 1)
	conn.Shutdown(closeChannel)
	<-closeChannel

	os.Exit(code)
}

func cleanCollections(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for _, name := range allCollections {
			coll := rawCollection(name)
			_, err := coll.DeleteMany(ctx, emptyFilter())
			if err != nil {
				// Cleanup failures can contaminate later tests; treat them as test failures.
				t.Errorf("failed to clean collection %s: %v", name, err)
			}
		}
	})
}

func rawCollection(name string) mongo.CollectionBinding {
	return conn.Collection(name)
}

func restoreFixtures(t *testing.T) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fixturesDir := fixturesPath()

	collectionFileMap := map[string]string{
		mongo.PeginQuoteCollection:          "peginQuote.json",
		mongo.RetainedPeginQuoteCollection:  "retainedPeginQuote.json",
		mongo.PeginCreationDataCollection:   "peginQuoteCreationData.json",
		mongo.PegoutQuoteCollection:         "pegoutQuote.json",
		mongo.RetainedPegoutQuoteCollection: "retainedPegoutQuote.json",
		mongo.PegoutCreationDataCollection:  "pegoutQuoteCreationData.json",
		mongo.DepositEventsCollection:       "depositEvents.json",
		mongo.LiquidityProviderCollection:   "liquidityProvider.json",
		mongo.TrustedAccountCollection:      "trustedAccounts.json",
		mongo.PenalizedEventCollection:      "penalizedEvent.json",
		mongo.BatchPegOutEventsCollection:   "batchPegOutEvents.json",
	}

	for collName, fileName := range collectionFileMap {
		filePath := filepath.Join(fixturesDir, fileName)
		data, err := os.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			t.Fatalf("Failed to read fixture file %s: %v", fileName, err)
		}

		var rawDocs []json.RawMessage
		if err := json.Unmarshal(data, &rawDocs); err != nil {
			t.Fatalf("Failed to unmarshal fixture file %s: %v", fileName, err)
		}

		if len(rawDocs) == 0 {
			continue
		}

		docsInterface := make([]any, len(rawDocs))
		for i, raw := range rawDocs {
			doc, err := extJSONToDocument(raw)
			if err != nil {
				t.Fatalf("Failed to parse extended json in %s: %v", fileName, err)
			}
			docsInterface[i] = doc
		}

		coll := rawCollection(collName)
		if _, err := coll.DeleteMany(ctx, emptyFilter()); err != nil {
			t.Fatalf("Failed to clear collection %s before restoring fixtures: %v", collName, err)
		}
		_, err = coll.InsertMany(ctx, docsInterface)
		if err != nil {
			t.Fatalf("Failed to insert fixtures for %s: %v", collName, err)
		}
	}
}

func fixturesPath() string {
	if v := os.Getenv("FIXTURES_DIR"); v != "" {
		return v
	}
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join("test", "mongodb", "fixtures")
	}
	return filepath.Join(filepath.Dir(filename), "fixtures")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envOrUint(key string, fallback uint) uint {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.ParseUint(v, 10, 32)
	if err != nil {
		return fallback
	}
	return uint(parsed)
}
