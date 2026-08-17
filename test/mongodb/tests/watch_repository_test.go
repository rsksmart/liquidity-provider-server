//go:build integration

package mongodb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	mongoAdapter "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDriver "go.mongodb.org/mongo-driver/v2/mongo"
)

const rskAddressIndex = "rsk_address_1"

func TestWatchSetKeepsOneDocumentPerRskAddressAgainstMongo(t *testing.T) {
	ctx := context.Background()
	collection := mongoClient.Database(mongoAdapter.DbName).Collection(mongoAdapter.PegInAddressRegistryWatchCollection)
	repository := mongoAdapter.NewPegInAddressRegistryWatchMongoRepository(conn)

	assertRskAddressIndexIsUnique(t, ctx, collection)

	event := rootstock.PegInAddressRegistryWatch{
		TxHash:      fmt.Sprintf("0xfly2514-%d", time.Now().UnixNano()),
		LogIndex:    7,
		BlockNumber: 4321,
		RskAddress:  fmt.Sprintf("0xfly2514rsk-%d", time.Now().UnixNano()),
		State:       rootstock.PegInAddressRegistryWatchDiscovered,
		CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
		UpdatedAt:   time.Now().UTC().Truncate(time.Millisecond),
	}
	identity := bson.M{"rsk_address": event.RskAddress}
	t.Cleanup(func() {
		_, err := collection.DeleteMany(ctx, identity)
		assert.NoError(t, err)
	})

	require.NoError(t, repository.Upsert(ctx, event))
	replay := event
	replay.TxHash = "0xreplay"
	replay.LogIndex = 8
	require.NoError(t, repository.Upsert(ctx, replay))
	documents, err := collection.CountDocuments(ctx, identity)
	require.NoError(t, err)
	assert.Equal(t, int64(1), documents, "a replayed rsk_address must not insert a second document")

	imported, err := repository.Get(ctx, event.RskAddress)
	require.NoError(t, err)
	require.NotNil(t, imported)
	imported.State = rootstock.PegInAddressRegistryWatchImported
	imported.BtcAddress = "n1BE7ioVukYS2GC88hT2K6cUvRiKwMwio7"
	require.NoError(t, repository.Update(ctx, *imported))

	require.NoError(t, repository.Upsert(ctx, event))
	stored, err := repository.Get(ctx, event.RskAddress)
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, rootstock.PegInAddressRegistryWatchImported, stored.State)
	assert.Equal(t, imported.BtcAddress, stored.BtcAddress)
	assert.Equal(t, event.TxHash, stored.TxHash)
}

func assertRskAddressIndexIsUnique(t *testing.T, ctx context.Context, collection *mongoDriver.Collection) {
	t.Helper()
	cursor, err := collection.Indexes().List(ctx)
	require.NoError(t, err)
	var indexes []bson.M
	require.NoError(t, cursor.All(ctx, &indexes))
	for _, index := range indexes {
		if index["name"] == rskAddressIndex {
			assert.Equal(t, true, index["unique"], "the rsk_address index must be unique")
			return
		}
	}
	require.Failf(t, "missing index", "%s was not built on %s", rskAddressIndex, collection.Name())
}
