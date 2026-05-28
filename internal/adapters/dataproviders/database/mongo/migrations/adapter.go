package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoDatabaseAdapter wraps *mongo.Database to satisfy DatabaseBinding.
type MongoDatabaseAdapter struct {
	db *mongo.Database
}

func NewMongoDatabaseAdapter(db *mongo.Database) *MongoDatabaseAdapter {
	return &MongoDatabaseAdapter{db: db}
}

func (a *MongoDatabaseAdapter) RunCommand(
	ctx context.Context, runCommand any, opts ...options.Lister[options.RunCmdOptions],
) *mongo.SingleResult {
	return a.db.RunCommand(ctx, runCommand, opts...)
}

func (a *MongoDatabaseAdapter) Collection(
	name string, opts ...options.Lister[options.CollectionOptions],
) CollectionBinding {
	return a.db.Collection(name, opts...)
}
