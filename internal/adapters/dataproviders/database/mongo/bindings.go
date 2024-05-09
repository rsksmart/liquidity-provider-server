package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// The wrapper structs defined in this class are meant to ease the mocking of the mongo driver structs

type DbClientBinding interface {
	Database(name string, opts ...*options.DatabaseOptions) DbBinding
	Disconnect(ctx context.Context) error
	Ping(ctx context.Context, rp *readpref.ReadPref) error
}

type ClientWrapper struct {
	*mongo.Client
	db DbBinding
}

func NewClientWrapper(client *mongo.Client) DbClientBinding {
	return &ClientWrapper{Client: client}
}

func (c *ClientWrapper) Database(name string, opts ...*options.DatabaseOptions) DbBinding {
	if c.db == nil {
		c.db = NewDatabaseWrapper(c.Client.Database(name, opts...))
	}
	return c.db
}

type DbBinding interface {
	Collection(name string, opts ...*options.CollectionOptions) CollectionBinding
}

type DatabaseWrapper struct {
	*mongo.Database
}

func NewDatabaseWrapper(db *mongo.Database) DbBinding {
	return &DatabaseWrapper{Database: db}
}

func (d *DatabaseWrapper) Collection(name string, opts ...*options.CollectionOptions) CollectionBinding {
	return d.Database.Collection(name, opts...)
}

type CollectionBinding interface {
	InsertOne(ctx context.Context, document any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []any, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult
	Find(ctx context.Context, filter any, opts ...*options.FindOptions) (*mongo.Cursor, error)
	UpdateOne(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	ReplaceOne(ctx context.Context, filter any, replacement any, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
}
