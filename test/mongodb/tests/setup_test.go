//go:build integration

package mongodb_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test/mongodb/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mongoDriver "go.mongodb.org/mongo-driver/v2/mongo"
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

var allCollections = func() []string {
	names := make([]string, len(utils.FixtureCollections))
	for i, fc := range utils.FixtureCollections {
		names[i] = fc.Collection
	}
	return names
}()

func TestMain(m *testing.M) {
	host := utils.EnvOr("MONGODB_HOST", "localhost")
	port, err := utils.EnvOrUint("MONGODB_PORT", 27018)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid MONGODB_PORT: %v\n", err)
		os.Exit(1)
	}
	username := utils.EnvOr("MONGODB_USER", "test")
	password := utils.EnvOr("MONGODB_PASSWORD", "test")

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, name := range allCollections {
		coll := rawCollection(name)
		_, err := coll.DeleteMany(ctx, utils.EmptyFilter())
		if err != nil {
			t.Fatalf("failed to clean collection %s: %v", name, err)
		}
	}
}

func rawCollection(name string) mongo.CollectionBinding {
	return conn.Collection(name)
}

func restoreFixtures(t *testing.T) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fixturesDir := fixturesPath()

	for _, fc := range utils.FixtureCollections {
		filePath := filepath.Join(fixturesDir, fc.FileName)
		data, err := os.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			t.Fatalf("Failed to read fixture file %s: %v", fc.FileName, err)
		}

		var rawDocs []json.RawMessage
		if err := json.Unmarshal(data, &rawDocs); err != nil {
			t.Fatalf("Failed to unmarshal fixture file %s: %v", fc.FileName, err)
		}

		if len(rawDocs) == 0 {
			continue
		}

		docsInterface := make([]any, len(rawDocs))
		for i, raw := range rawDocs {
			doc, err := utils.ExtJSONToDocument(raw)
			if err != nil {
				t.Fatalf("Failed to parse extended json in %s: %v", fc.FileName, err)
			}
			docsInterface[i] = doc
		}

		coll := rawCollection(fc.Collection)
		if _, err := coll.DeleteMany(ctx, utils.EmptyFilter()); err != nil {
			t.Fatalf("Failed to clear collection %s before restoring fixtures: %v", fc.Collection, err)
		}
		_, err = coll.InsertMany(ctx, docsInterface)
		if err != nil {
			t.Fatalf("Failed to insert fixtures for %s: %v", fc.Collection, err)
		}
	}
}

func fixturesPath() string {
	return utils.FixturesPath()
}

func assertWeiEqual(t testing.TB, want, got *entities.Wei) {
	t.Helper()
	require.NotNil(t, want)
	require.NotNil(t, got)
	assert.Zero(t, want.Cmp(got))
}
