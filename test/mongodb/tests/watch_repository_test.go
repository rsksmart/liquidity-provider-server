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
	collection := mongoClient.Database(mongoAdapter.DbName).Collection(mongoAdapter.PegInWatchCollection)
	repository := mongoAdapter.NewPegInWatchMongoRepository(conn)

	assertRskAddressIndexIsUnique(t, ctx, collection)

	event := rootstock.PegInWatch{
		TxHash:      fmt.Sprintf("0xfly2514-%d", time.Now().UnixNano()),
		LogIndex:    7,
		BlockNumber: 4321,
		RskAddress:  fmt.Sprintf("0xfly2514rsk-%d", time.Now().UnixNano()),
		State:       rootstock.PegInWatchDiscovered,
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
	imported.State = rootstock.PegInWatchImported
	imported.BtcAddress = "n1BE7ioVukYS2GC88hT2K6cUvRiKwMwio7"
	require.NoError(t, repository.Update(ctx, *imported))

	require.NoError(t, repository.Upsert(ctx, event))
	stored, err := repository.Get(ctx, event.RskAddress)
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, rootstock.PegInWatchImported, stored.State)
	assert.Equal(t, imported.BtcAddress, stored.BtcAddress)
	assert.Equal(t, event.TxHash, stored.TxHash)
}

func TestVerifiedRegistrationsRoundTripAgainstMongo(t *testing.T) {
	ctx := context.Background()
	watches := mongoClient.Database(mongoAdapter.DbName).Collection(mongoAdapter.PegInWatchCollection)
	checkpoints := mongoClient.Database(mongoAdapter.DbName).Collection(mongoAdapter.PegInWatchCheckpointCollection)
	repository := mongoAdapter.NewPegInWatchMongoRepository(conn)
	require.NoError(t, checkpoints.Drop(ctx))

	watch := rootstock.NewPegInWatch("0xround-trip", 3, 10, fmt.Sprintf("0xroundtrip-%d", time.Now().UnixNano()), "0xregistrant", [32]byte{1})
	t.Cleanup(func() {
		_, err := watches.DeleteMany(ctx, bson.M{"rsk_address": watch.RskAddress})
		assert.NoError(t, err)
		assert.NoError(t, checkpoints.Drop(ctx))
	})

	_, found, err := repository.GetCheckpoint(ctx)
	require.NoError(t, err)
	assert.False(t, found)

	first := rootstock.PegInWatchCheckpoint{Root: [32]byte{1}, Block: 10}
	require.NoError(t, repository.SaveVerifiedRegistrations(ctx, []rootstock.PegInWatch{watch}, first))
	stored, found, err := repository.GetCheckpoint(ctx)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, first, stored)
	storedWatch, err := repository.Get(ctx, watch.RskAddress)
	require.NoError(t, err)
	require.NotNil(t, storedWatch)
	assert.Equal(t, watch.TxHash, storedWatch.TxHash)
	assert.Equal(t, rootstock.PegInWatchDiscovered, storedWatch.State)

	replay := watch
	replay.TxHash = "0xreplay"
	second := rootstock.PegInWatchCheckpoint{Root: [32]byte{2}, Block: 20}
	require.NoError(t, repository.SaveVerifiedRegistrations(ctx, []rootstock.PegInWatch{replay}, second))
	stored, _, err = repository.GetCheckpoint(ctx)
	require.NoError(t, err)
	assert.Equal(t, second, stored)
	storedWatch, err = repository.Get(ctx, watch.RskAddress)
	require.NoError(t, err)
	require.NotNil(t, storedWatch)
	assert.Equal(t, watch.TxHash, storedWatch.TxHash)
	documents, err := checkpoints.CountDocuments(ctx, bson.M{})
	require.NoError(t, err)
	assert.Equal(t, int64(1), documents)
}

func TestVerifiedRegistrationsAreAtomicAgainstMongo(t *testing.T) {
	ctx := context.Background()
	database := mongoClient.Database(mongoAdapter.DbName)
	watches := database.Collection(mongoAdapter.PegInWatchCollection)
	checkpoints := database.Collection(mongoAdapter.PegInWatchCheckpointCollection)
	repository := mongoAdapter.NewPegInWatchMongoRepository(conn)
	require.NoError(t, checkpoints.Drop(ctx))

	suffix := time.Now().UnixNano()
	kept := rootstock.NewPegInWatch("0xkept", 1, 10, fmt.Sprintf("0xkept-%d", suffix), "0xregistrant", [32]byte{1})
	rolledBack := rootstock.NewPegInWatch("0xrolled-back", 1, 20, fmt.Sprintf("0xrolledback-%d", suffix), "0xregistrant", [32]byte{2})
	t.Cleanup(func() {
		_, err := watches.DeleteMany(ctx, bson.M{"rsk_address": bson.M{"$in": []string{kept.RskAddress, rolledBack.RskAddress}}})
		assert.NoError(t, err)
		assert.NoError(t, checkpoints.Drop(ctx))
	})

	previous := rootstock.PegInWatchCheckpoint{Root: [32]byte{1}, Block: 10}
	require.NoError(t, repository.SaveVerifiedRegistrations(ctx, []rootstock.PegInWatch{kept}, previous))
	require.NoError(t, database.RunCommand(ctx, bson.D{
		{Key: "collMod", Value: mongoAdapter.PegInWatchCheckpointCollection},
		{Key: "validator", Value: bson.M{"block": bson.M{"$lt": 0}}},
	}).Err())

	err := repository.SaveVerifiedRegistrations(
		ctx,
		[]rootstock.PegInWatch{rolledBack},
		rootstock.PegInWatchCheckpoint{Root: [32]byte{2}, Block: 20},
	)

	require.Error(t, err)
	missing, err := repository.Get(ctx, rolledBack.RskAddress)
	require.NoError(t, err)
	assert.Nil(t, missing, "a watch written before the failed checkpoint write must be rolled back")
	stored, found, err := repository.GetCheckpoint(ctx)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, previous, stored)
	keptWatch, err := repository.Get(ctx, kept.RskAddress)
	require.NoError(t, err)
	assert.NotNil(t, keptWatch)
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

func TestWatchRowsAreIndexedByBlockAndLogIndexAgainstMongo(t *testing.T) {
	ctx := context.Background()
	collection := mongoClient.Database(mongoAdapter.DbName).Collection(mongoAdapter.PegInWatchCollection)

	cursor, err := collection.Indexes().List(ctx)
	require.NoError(t, err)
	var indexes []struct {
		Name string `bson:"name"`
		Key  bson.D `bson:"key"`
	}
	require.NoError(t, cursor.All(ctx, &indexes))

	for _, index := range indexes {
		if index.Name == "block_number_1_log_index_1" {
			require.Len(t, index.Key, 2)
			assert.Equal(t, "block_number", index.Key[0].Key)
			assert.Equal(t, "log_index", index.Key[1].Key)
			assert.EqualValues(t, 1, index.Key[0].Value)
			assert.EqualValues(t, 1, index.Key[1].Value)
			return
		}
	}
	require.Failf(t, "missing index", "block_number_1_log_index_1 was not built on %s", collection.Name())
}

func replaceFixtureBlock() uint64 {
	return 900_000_000 + uint64(time.Now().UnixNano()%1_000_000)*100
}

//nolint:funlen // One scenario checks every kind of row a rebuild touches.
func TestReplacedRegistrationsKeepImportStateAgainstMongo(t *testing.T) {
	ctx := context.Background()
	watches := mongoClient.Database(mongoAdapter.DbName).Collection(mongoAdapter.PegInWatchCollection)
	checkpoints := mongoClient.Database(mongoAdapter.DbName).Collection(mongoAdapter.PegInWatchCheckpointCollection)
	repository := mongoAdapter.NewPegInWatchMongoRepository(conn)
	require.NoError(t, checkpoints.Drop(ctx))

	base := replaceFixtureBlock()
	suffix := time.Now().UnixNano()
	below := rootstock.NewPegInWatch("0xbelow", 1, base-1, fmt.Sprintf("0xbelow-%d", suffix), "0xr", [32]byte{1})
	moved := rootstock.NewPegInWatch("0xmoved", 1, base+1, fmt.Sprintf("0xmoved-%d", suffix), "0xr", [32]byte{2})
	dropped := rootstock.NewPegInWatch("0xdropped", 1, base+5, fmt.Sprintf("0xdropped-%d", suffix), "0xr", [32]byte{3})
	added := rootstock.NewPegInWatch("0xadded", 2, base+7, fmt.Sprintf("0xadded-%d", suffix), "0xr", [32]byte{5})
	t.Cleanup(func() {
		_, err := watches.DeleteMany(ctx, bson.M{"block_number": bson.M{"$gte": base - 1, "$lte": base + 10}})
		assert.NoError(t, err)
		assert.NoError(t, checkpoints.Drop(ctx))
	})
	require.NoError(t, repository.SaveVerifiedRegistrations(
		ctx,
		[]rootstock.PegInWatch{below, moved, dropped},
		rootstock.PegInWatchCheckpoint{Root: [32]byte{3}, Block: base + 5},
	))
	moved.MarkImported()
	moved.BtcAddress = "n1BE7ioVukYS2GC88hT2K6cUvRiKwMwio7"
	require.NoError(t, repository.Update(ctx, moved))

	rebuiltMoved := rootstock.NewPegInWatch("0xmoved-again", 0, base+7, moved.RskAddress, "0xr", [32]byte{4})
	rebuilt := rootstock.PegInWatchCheckpoint{Root: [32]byte{5}, Block: base + 8}
	require.NoError(t, repository.ReplaceVerifiedRegistrations(ctx, base, []rootstock.PegInWatch{rebuiltMoved, added}, rebuilt))

	storedMoved, err := repository.Get(ctx, moved.RskAddress)
	require.NoError(t, err)
	require.NotNil(t, storedMoved)
	assert.Equal(t, rootstock.PegInWatchImported, storedMoved.State)
	assert.Equal(t, moved.BtcAddress, storedMoved.BtcAddress)
	assert.Equal(t, "0xmoved-again", storedMoved.TxHash)
	assert.Equal(t, base+7, storedMoved.BlockNumber)
	assert.Equal(t, [32]byte{4}, storedMoved.RegistrationRoot)
	storedAdded, err := repository.Get(ctx, added.RskAddress)
	require.NoError(t, err)
	require.NotNil(t, storedAdded)
	assert.Equal(t, rootstock.PegInWatchDiscovered, storedAdded.State)
	assert.Equal(t, added.RskAddress, storedAdded.RskAddress)
	storedDropped, err := repository.Get(ctx, dropped.RskAddress)
	require.NoError(t, err)
	assert.Nil(t, storedDropped)
	storedBelow, err := repository.Get(ctx, below.RskAddress)
	require.NoError(t, err)
	assert.NotNil(t, storedBelow)

	checkpoint, found, err := repository.GetCheckpoint(ctx)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, rebuilt, checkpoint)
	rootBefore, err := repository.GetRootAt(ctx, base)
	require.NoError(t, err)
	assert.Equal(t, below.RegistrationRoot, rootBefore)
	rootAfter, err := repository.GetRootAt(ctx, base+8)
	require.NoError(t, err)
	assert.Equal(t, added.RegistrationRoot, rootAfter)
}

func TestReplacedRegistrationsAreAtomicAgainstMongo(t *testing.T) {
	ctx := context.Background()
	database := mongoClient.Database(mongoAdapter.DbName)
	watches := database.Collection(mongoAdapter.PegInWatchCollection)
	checkpoints := database.Collection(mongoAdapter.PegInWatchCheckpointCollection)
	repository := mongoAdapter.NewPegInWatchMongoRepository(conn)
	require.NoError(t, checkpoints.Drop(ctx))

	base := replaceFixtureBlock()
	suffix := time.Now().UnixNano()
	stale := rootstock.NewPegInWatch("0xstale", 1, base+1, fmt.Sprintf("0xstale-%d", suffix), "0xr", [32]byte{1})
	fresh := rootstock.NewPegInWatch("0xfresh", 1, base+2, fmt.Sprintf("0xfresh-%d", suffix), "0xr", [32]byte{2})
	t.Cleanup(func() {
		_, err := watches.DeleteMany(ctx, bson.M{"block_number": bson.M{"$gte": base, "$lte": base + 10}})
		assert.NoError(t, err)
		assert.NoError(t, checkpoints.Drop(ctx))
	})
	previous := rootstock.PegInWatchCheckpoint{Root: [32]byte{1}, Block: base + 1}
	require.NoError(t, repository.SaveVerifiedRegistrations(ctx, []rootstock.PegInWatch{stale}, previous))
	require.NoError(t, database.RunCommand(ctx, bson.D{
		{Key: "collMod", Value: mongoAdapter.PegInWatchCheckpointCollection},
		{Key: "validator", Value: bson.M{"block": bson.M{"$lt": 0}}},
	}).Err())

	err := repository.ReplaceVerifiedRegistrations(
		ctx,
		base,
		[]rootstock.PegInWatch{fresh},
		rootstock.PegInWatchCheckpoint{Root: [32]byte{2}, Block: base + 2},
	)

	require.Error(t, err)
	storedStale, err := repository.Get(ctx, stale.RskAddress)
	require.NoError(t, err)
	assert.NotNil(t, storedStale, "a row removed before the failed checkpoint write must be restored")
	storedFresh, err := repository.Get(ctx, fresh.RskAddress)
	require.NoError(t, err)
	assert.Nil(t, storedFresh, "a row written before the failed checkpoint write must be rolled back")
	stored, found, err := repository.GetCheckpoint(ctx)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, previous, stored)
}

func TestReplacedRegistrationsHandleManyAddressesAgainstMongo(t *testing.T) {
	const addresses = 5000
	ctx := context.Background()
	watches := mongoClient.Database(mongoAdapter.DbName).Collection(mongoAdapter.PegInWatchCollection)
	checkpoints := mongoClient.Database(mongoAdapter.DbName).Collection(mongoAdapter.PegInWatchCheckpointCollection)
	repository := mongoAdapter.NewPegInWatchMongoRepository(conn)
	require.NoError(t, checkpoints.Drop(ctx))

	base := replaceFixtureBlock()
	suffix := time.Now().UnixNano()
	rebuilt := make([]rootstock.PegInWatch, 0, addresses)
	for index := range addresses {
		rebuilt = append(rebuilt, rootstock.NewPegInWatch(
			fmt.Sprintf("0xmany-%d", index), uint(index), base, fmt.Sprintf("0xmany-%d-%d", suffix, index), "0xr", [32]byte{1},
		))
	}
	t.Cleanup(func() {
		_, err := watches.DeleteMany(ctx, bson.M{"block_number": base})
		assert.NoError(t, err)
		assert.NoError(t, checkpoints.Drop(ctx))
	})

	require.NoError(t, repository.ReplaceVerifiedRegistrations(ctx, base, rebuilt, rootstock.PegInWatchCheckpoint{Root: [32]byte{1}, Block: base}))

	documents, err := watches.CountDocuments(ctx, bson.M{"block_number": base})
	require.NoError(t, err)
	assert.Equal(t, int64(addresses), documents)
}
