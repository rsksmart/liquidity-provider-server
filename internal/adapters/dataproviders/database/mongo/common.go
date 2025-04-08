package mongo

import (
	"context"
	"fmt"
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

type QuoteResult[Q any, R QuoteHashProvider] struct {
	Quotes           []Q
	RetainedQuotes   []R
	QuoteHashToIndex map[string]int
	Error            error
}

type QuoteQuery struct {
	Ctx                context.Context
	Conn               *Connection
	StartDate          time.Time
	EndDate            time.Time
	QuoteCollection    string
	RetainedCollection string
}

func ListQuotesByDateRange[Q any, R QuoteHashProvider](
	query QuoteQuery,
	mapper func(bson.D) Q,
) QuoteResult[Q, R] {
	dbCtx, cancel := context.WithTimeout(query.Ctx, query.Conn.timeout)
	defer cancel()
	quotes, quoteHashes, err := fetchQuotesByDateRange(dbCtx, query.Conn, query.StartDate, query.EndDate, query.QuoteCollection, mapper)
	if err != nil {
		return QuoteResult[Q, R]{Error: err}
	}
	quoteHashToIndex := make(map[string]int, len(quoteHashes))
	for i, hash := range quoteHashes {
		if hash != "" {
			quoteHashToIndex[hash] = i
		}
	}
	retainedQuotes, additionalHashes, err := fetchRetainedQuotes[R](dbCtx, query.Conn, query.StartDate, query.EndDate, query.RetainedCollection, quoteHashes)
	if err != nil {
		return QuoteResult[Q, R]{Error: err}
	}
	if len(additionalHashes) > 0 {
		additionalQuotes, additionalHashIndices, err := fetchAdditionalQuotes(dbCtx, query.Conn, query.QuoteCollection, additionalHashes, mapper)
		if err != nil {
			log.Errorf("Error processing additional quotes: %v", err)
		} else {
			baseIndex := len(quotes)
			for i, hash := range additionalHashIndices {
				if hash != "" {
					quoteHashToIndex[hash] = baseIndex + i
				}
			}
			quotes = append(quotes, additionalQuotes...)
		}
	}
	logDbInteraction(Read, fmt.Sprintf("Found %d quotes and %d retained quotes in date range",
		len(quotes), len(retainedQuotes)))
	return QuoteResult[Q, R]{
		Quotes:           quotes,
		RetainedQuotes:   retainedQuotes,
		QuoteHashToIndex: quoteHashToIndex,
		Error:            nil,
	}
}

func fetchQuotesByDateRange[Q any](
	ctx context.Context,
	conn *Connection,
	startDate, endDate time.Time,
	collectionName string,
	mapper func(bson.D) Q,
) ([]Q, []string, error) {
	quoteFilter := bson.D{
		{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}},
	}
	var storedQuotes []bson.D
	quoteCursor, err := conn.Collection(collectionName).Find(ctx, quoteFilter)
	if err != nil {
		return nil, nil, err
	}
	if err = quoteCursor.All(ctx, &storedQuotes); err != nil {
		return nil, nil, err
	}
	quoteHashes := make([]string, 0, len(storedQuotes))
	quotes := make([]Q, 0, len(storedQuotes))
	for _, stored := range storedQuotes {
		quoteObj := mapper(stored)
		quotes = append(quotes, quoteObj)
		hashValue, ok := getStringValueFromBSON(stored, "hash")
		if ok {
			quoteHashes = append(quoteHashes, hashValue)
		}
	}
	return quotes, quoteHashes, nil
}

func getStringValueFromBSON(doc bson.D, key string) (string, bool) {
	data, err := bson.Marshal(doc)
	if err != nil {
		return "", false
	}
	rawValue := bson.Raw(data).Lookup(key)
	return rawValue.StringValueOK()
}

type QuoteHashProvider interface {
	GetQuoteHash() string
}

func fetchRetainedQuotes[R QuoteHashProvider](
	ctx context.Context,
	conn *Connection,
	startDate, endDate time.Time,
	collectionName string,
	existingQuoteHashes []string,
) ([]R, []string, error) {
	retainedFilter := createRetainedFilter(startDate, endDate, existingQuoteHashes)
	var retainedQuotes []R
	retainedCursor, err := conn.Collection(collectionName).Find(ctx, retainedFilter)
	if err != nil {
		return nil, nil, err
	}
	if err = retainedCursor.All(ctx, &retainedQuotes); err != nil {
		return nil, nil, err
	}
	additionalHashes := findAdditionalQuoteHashes(retainedQuotes, existingQuoteHashes)
	return retainedQuotes, additionalHashes, nil
}

func createRetainedFilter(startDate, endDate time.Time, quoteHashes []string) bson.D {
	return bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "quote_hash", Value: bson.D{
				{Key: "$in", Value: quoteHashes},
			}}},
			bson.D{
				{Key: "created_at", Value: bson.D{
					{Key: "$gte", Value: startDate.Unix()},
					{Key: "$lte", Value: endDate.Unix()},
				}},
			},
		}},
	}
}

func findAdditionalQuoteHashes[R QuoteHashProvider](retainedQuotes []R, existingQuoteHashes []string) []string {
	existingMap := make(map[string]bool, len(existingQuoteHashes))
	for _, hash := range existingQuoteHashes {
		existingMap[hash] = true
	}
	additionalMap := make(map[string]bool)
	for i := range retainedQuotes {
		hash := retainedQuotes[i].GetQuoteHash()
		if !existingMap[hash] {
			additionalMap[hash] = true
		}
	}
	additionalHashes := make([]string, 0, len(additionalMap))
	for hash := range additionalMap {
		additionalHashes = append(additionalHashes, hash)
	}
	return additionalHashes
}

func fetchAdditionalQuotes[Q any](
	ctx context.Context,
	conn *Connection,
	collectionName string,
	hashes []string,
	mapper func(bson.D) Q,
) ([]Q, []string, error) {
	quoteFilter := bson.D{
		{Key: "hash", Value: bson.D{
			{Key: "$in", Value: hashes},
		}},
	}
	var storedQuotes []bson.D
	quoteCursor, err := conn.Collection(collectionName).Find(ctx, quoteFilter)
	if err != nil {
		return nil, nil, err
	}
	if err = quoteCursor.All(ctx, &storedQuotes); err != nil {
		return nil, nil, err
	}
	quotes := make([]Q, 0, len(storedQuotes))
	resultHashes := make([]string, 0, len(storedQuotes))
	for _, stored := range storedQuotes {
		quoteObj := mapper(stored)
		quotes = append(quotes, quoteObj)
		hashValue, ok := getStringValueFromBSON(stored, "hash")
		if ok {
			resultHashes = append(resultHashes, hashValue)
		} else {
			resultHashes = append(resultHashes, "")
		}
	}
	return quotes, resultHashes, nil
}
