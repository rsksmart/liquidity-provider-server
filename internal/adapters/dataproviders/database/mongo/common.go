package mongo

import (
	"context"
	"fmt"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
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

func ListQuotesByDateRange[S any, Q any, R any]( //nolint:funlen,cyclop
	ctx context.Context,
	conn *Connection,
	startDate, endDate time.Time,
	quoteCollection, retainedCollection string,
	extractor func(S) (string, Q),
) ([]Q, []R, error) {
	dbCtx, cancel := context.WithTimeout(ctx, conn.timeout)
	defer cancel()

	quoteFilter := bson.D{
		{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}},
	}

	var storedQuotes []S
	quoteCursor, err := conn.Collection(quoteCollection).Find(dbCtx, quoteFilter)
	if err != nil {
		return nil, nil, err
	}

	if err = quoteCursor.All(dbCtx, &storedQuotes); err != nil {
		return nil, nil, err
	}
	
	quoteHashes := make([]string, 0, len(storedQuotes))
	quotes := make([]Q, 0, len(storedQuotes))

	for _, stored := range storedQuotes {
		hash, quoteObj := extractor(stored)
		quoteHashes = append(quoteHashes, hash)
		quotes = append(quotes, quoteObj)
	}

	retainedFilter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "quote_hash", Value: bson.D{
				{Key: "$in", Value: quoteHashes},
			}}},
			bson.D{
				{Key: "updated_at", Value: bson.D{
					{Key: "$gte", Value: startDate.Unix()},
					{Key: "$lte", Value: endDate.Unix()},
				}},
			},
		}},
	}

	var retainedQuotes []R
	retainedCursor, err := conn.Collection(retainedCollection).Find(dbCtx, retainedFilter)
	if err != nil {
		return nil, nil, err
	}

	if err = retainedCursor.All(dbCtx, &retainedQuotes); err != nil {
		return nil, nil, err
	}

	additionalHashes := make(map[string]bool)
	for _, retained := range retainedQuotes {
		v := reflect.ValueOf(retained)
		hashField := v.FieldByName("QuoteHash")
		if hashField.IsValid() && hashField.Kind() == reflect.String {
			hash := hashField.String()			
			found := false
			for _, existingHash := range quoteHashes {
				if existingHash == hash {
					found = true
					break
				}
			}
			if !found {
				additionalHashes[hash] = true
			}
		}
	}
	if len(additionalHashes) > 0 { //nolint:nestif
		additionalHashesList := make([]string, 0, len(additionalHashes))
		for hash := range additionalHashes {
			additionalHashesList = append(additionalHashesList, hash)
		}
		additionalFilter := bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "$in", Value: additionalHashesList},
			}},
		}
		var additionalStoredQuotes []S
		additionalCursor, err := conn.Collection(quoteCollection).Find(dbCtx, additionalFilter)
		if err != nil {
			log.Errorf("Error fetching additional quotes: %v", err)
		} else {
			if err = additionalCursor.All(dbCtx, &additionalStoredQuotes); err != nil {
				log.Errorf("Error processing additional quotes: %v", err)
			} else {
				for _, stored := range additionalStoredQuotes {
					hash, quoteObj := extractor(stored)
					quoteHashes = append(quoteHashes, hash)
					quotes = append(quotes, quoteObj)
				}
			}
		}
	}

	logDbInteraction(Read, fmt.Sprintf("Found %d quotes and %d retained quotes in date range",
		len(quotes), len(retainedQuotes)))

	return quotes, retainedQuotes, nil
}
