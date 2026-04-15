package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// The wrapper structs defined in this class are meant to ease the mocking of the mongo driver structs

type SessionBinding interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) (interface{}, error),
		opts ...options.Lister[options.TransactionOptions]) (interface{}, error)
	EndSession(context.Context)
}

type DbClientBinding interface {
	Database(name string, opts ...options.Lister[options.DatabaseOptions]) DbBinding
	Disconnect(ctx context.Context) error
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	StartSession(opts ...options.Lister[options.SessionOptions]) (SessionBinding, error)
}

type ClientWrapper struct {
	*mongo.Client
	db DbBinding
}

func NewClientWrapper(client *mongo.Client) DbClientBinding {
	return &ClientWrapper{Client: client}
}

func (c *ClientWrapper) Database(name string, opts ...options.Lister[options.DatabaseOptions]) DbBinding {
	if c.db == nil {
		c.db = NewDatabaseWrapper(c.Client.Database(name, opts...))
	}
	return c.db
}

func (c *ClientWrapper) StartSession(opts ...options.Lister[options.SessionOptions]) (SessionBinding, error) {
	return c.Client.StartSession(opts...)
}

type DbBinding interface {
	Collection(name string, opts ...options.Lister[options.CollectionOptions]) CollectionBinding
}

type DatabaseWrapper struct {
	*mongo.Database
}

func NewDatabaseWrapper(db *mongo.Database) DbBinding {
	return &DatabaseWrapper{Database: db}
}

func (d *DatabaseWrapper) Collection(name string, opts ...options.Lister[options.CollectionOptions]) CollectionBinding {
	return d.Database.Collection(name, opts...)
}

type CollectionBinding interface {
	InsertOne(ctx context.Context, document any, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents any, opts ...options.Lister[options.InsertManyOptions]) (*mongo.InsertManyResult, error)
	FindOne(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult
	Find(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error)
	UpdateOne(ctx context.Context, filter any, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter any, update any, opts ...options.Lister[options.UpdateManyOptions]) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter any, opts ...options.Lister[options.DeleteManyOptions]) (*mongo.DeleteResult, error)
	ReplaceOne(ctx context.Context, filter any, replacement any, opts ...options.Lister[options.ReplaceOptions]) (*mongo.UpdateResult, error)
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...options.Lister[options.BulkWriteOptions]) (*mongo.BulkWriteResult, error)
	Aggregate(ctx context.Context, pipeline any, opts ...options.Lister[options.AggregateOptions]) (*mongo.Cursor, error)
}
