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

const eventIdentityIndex = "tx_hash_1_log_index_1"

// The three repository semantics the in-memory watch store models, pinned against a real server:
// the unique index on the event identity builds, a replayed event inserts no second document, and a
// replayed event does not overwrite the progress already recorded for it. The package TestMain
// opens the connection through the production mongo.Connect path, so the index under test is the
// one a real boot would build.
func TestWatchSetKeepsOneDocumentPerEventAgainstMongo(t *testing.T) {
	ctx := context.Background()
	collection := mongoClient.Database(mongoAdapter.DbName).Collection(mongoAdapter.PegInAddressRegistryWatchCollection)
	repository := mongoAdapter.NewPegInAddressRegistryWatchMongoRepository(conn)

	assertEventIdentityIndexIsUnique(t, ctx, collection)

	event := rootstock.PegInAddressRegistryWatch{
		TxHash:      fmt.Sprintf("0xfly2514-%d", time.Now().UnixNano()),
		LogIndex:    7,
		BlockNumber: 4321,
		RskAddress:  "0xfly2514rsk",
		State:       rootstock.PegInAddressRegistryWatchDiscovered,
		CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
		UpdatedAt:   time.Now().UTC().Truncate(time.Millisecond),
	}
	identity := bson.M{"tx_hash": event.TxHash, "log_index": event.LogIndex}
	t.Cleanup(func() {
		_, err := collection.DeleteMany(ctx, identity)
		assert.NoError(t, err)
	})

	require.NoError(t, repository.Upsert(ctx, event))
	require.NoError(t, repository.Upsert(ctx, event))
	documents, err := collection.CountDocuments(ctx, identity)
	require.NoError(t, err)
	assert.Equal(t, int64(1), documents, "a replayed event must not insert a second document")

	imported, err := repository.Get(ctx, event.TxHash, event.LogIndex)
	require.NoError(t, err)
	require.NotNil(t, imported)
	imported.State = rootstock.PegInAddressRegistryWatchImported
	imported.BtcAddress = "n1BE7ioVukYS2GC88hT2K6cUvRiKwMwio7"
	imported.DepositTxID = "0xdeposit"
	require.NoError(t, repository.Update(ctx, *imported))

	// The overlap window re-delivers the raw event after the entry has progressed, which is the
	// case that would regress the entry if the upsert were a replace rather than a set-on-insert.
	require.NoError(t, repository.Upsert(ctx, event))
	stored, err := repository.Get(ctx, event.TxHash, event.LogIndex)
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, rootstock.PegInAddressRegistryWatchImported, stored.State)
	assert.Equal(t, imported.BtcAddress, stored.BtcAddress)
	assert.Equal(t, imported.DepositTxID, stored.DepositTxID)
}

func assertEventIdentityIndexIsUnique(t *testing.T, ctx context.Context, collection *mongoDriver.Collection) {
	t.Helper()
	cursor, err := collection.Indexes().List(ctx)
	require.NoError(t, err)
	var indexes []bson.M
	require.NoError(t, cursor.All(ctx, &indexes))
	for _, index := range indexes {
		if index["name"] == eventIdentityIndex {
			assert.Equal(t, true, index["unique"], "the event identity index must be unique")
			return
		}
	}
	require.Failf(t, "missing index", "%s was not built on %s", eventIdentityIndex, collection.Name())
}
