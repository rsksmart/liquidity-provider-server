package mongo

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	DbName = "flyover"
)

type DbInteraction string

const (
	Read   DbInteraction = "READ"
	Insert DbInteraction = "INSERT"
	Update DbInteraction = "UPDATE"
	Upsert DbInteraction = "UPSERT"
	Delete DbInteraction = "DELETE"
)

func logDbInteraction(interaction DbInteraction, value any) {
	const msgTemplate = "%s interaction with db: %+v"
	switch interaction {
	case Insert, Update, Upsert:
		log.Infof(msgTemplate, interaction, value)
	case Read:
		log.Debugf(msgTemplate, interaction, value)
	case Delete:
		log.Debugf(msgTemplate, interaction, value)
	default:
		log.Debug("Unknown DB interaction")
	}
}

type Connection struct {
	client  DbClientBinding
	db      DbBinding
	timeout time.Duration
}

func NewConnection(client DbClientBinding, timeout time.Duration) *Connection {
	db := client.Database(DbName)
	return &Connection{client: client, db: db, timeout: timeout}
}

func (c *Connection) GetDb() DbBinding {
	return c.db
}

func (c *Connection) Collection(collection string) CollectionBinding {
	return c.db.Collection(collection)
}

func (c *Connection) Shutdown(closeChannel chan<- bool) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
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
