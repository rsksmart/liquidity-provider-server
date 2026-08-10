package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDb "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// capturingIndexCreator records every index application so the tests can assert on the
// exact collection and model the production code sends through the seam.
type capturingIndexCreator struct {
	collections []string
	models      []mongoDb.IndexModel
	err         error
}

func (c *capturingIndexCreator) CreateIndex(_ context.Context, collection string, model mongoDb.IndexModel) error {
	c.collections = append(c.collections, collection)
	c.models = append(c.models, model)
	return c.err
}

func resolveIndexOptions(t *testing.T, model mongoDb.IndexModel) *options.IndexOptions {
	t.Helper()
	require.NotNil(t, model.Options)
	resolved := &options.IndexOptions{}
	for _, apply := range model.Options.List() {
		require.NoError(t, apply(resolved))
	}
	return resolved
}

func TestPegInAddressRegistryWatchEventIdentityIndex(t *testing.T) {
	model := pegInAddressRegistryWatchEventIdentityIndex()

	keys, ok := model.Keys.(bson.D)
	require.True(t, ok, "Keys should be bson.D")
	assert.Equal(t, bson.D{
		{Key: "tx_hash", Value: 1},
		{Key: "log_index", Value: 1},
	}, keys)

	resolved := resolveIndexOptions(t, model)
	require.NotNil(t, resolved.Unique)
	assert.True(t, *resolved.Unique)
}

func TestEnsurePegInAddressRegistryWatchEventIdentityIndex(t *testing.T) {
	t.Run("applies the unique event identity index to the watch collection", func(t *testing.T) {
		creator := &capturingIndexCreator{}

		require.NoError(t, ensurePegInAddressRegistryWatchEventIdentityIndex(context.Background(), creator))

		require.Equal(t, []string{PegInAddressRegistryWatchCollection}, creator.collections)
		require.Len(t, creator.models, 1)
		assert.Equal(t, bson.D{
			{Key: "tx_hash", Value: 1},
			{Key: "log_index", Value: 1},
		}, creator.models[0].Keys)
		resolved := resolveIndexOptions(t, creator.models[0])
		require.NotNil(t, resolved.Unique)
		assert.True(t, *resolved.Unique)
	})
	t.Run("preserves the underlying error with startup context", func(t *testing.T) {
		creator := &capturingIndexCreator{err: assert.AnError}

		err := ensurePegInAddressRegistryWatchEventIdentityIndex(context.Background(), creator)

		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		assert.Contains(t, err.Error(), "error creating unique index on "+PegInAddressRegistryWatchCollection+" event identity")
	})
}
