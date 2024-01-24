package mongo

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	dbTimeout = 3 * time.Second
	dbName    = "flyover"
)

type dbInteraction string

const (
	read   dbInteraction = "READ"
	insert dbInteraction = "INSERT"
	update dbInteraction = "UPDATE"
	upsert dbInteraction = "UPSERT"
	delete dbInteraction = "DELETE"
)

func logDbInteraction(interaction dbInteraction, value any) {
	switch interaction {
	case insert, update, upsert:
		log.Infof("%s interaction with db: %#v\n", interaction, value)
	case read:
		log.Debugf("%s interaction with db: %#v\n", interaction, value)
	case delete:
		log.Debugf("%s interaction with db: %v\n", interaction, value)
	default:
		log.Debug("Unknown DB interaction")
	}
}

type Connection struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewConnection(client *mongo.Client) *Connection {
	db := client.Database(dbName)
	return &Connection{client: client, db: db}
}

func (c *Connection) GetDb() *mongo.Database {
	return c.db
}

func (c *Connection) Collection(collection string) *mongo.Collection {
	return c.db.Collection(collection)
}

func (c *Connection) Shutdown(closeChannel chan<- bool) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	err := c.client.Disconnect(ctx)

	cancel()
	closeChannel <- true

	if err != nil {
		log.Error("Error disconnecting from MongoDB: ", err)
	} else {
		log.Debug("Disconnected from MongoDB")
	}
}

func (c *Connection) CheckConnection(ctx context.Context) bool {
	err := c.client.Ping(ctx, nil)
	if err != nil {
		log.Error("Error checking database connection: ", err)
	}
	return err == nil
}
