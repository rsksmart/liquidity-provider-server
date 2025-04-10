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

type ListQuotesResult[Q any, R QuoteHashProvider] struct {
	Quotes           []Q
	RetainedQuotes   []R
	quoteHashToIndex map[string]int
}

func (r *ListQuotesResult[Q, R]) GetQuoteByHash(hash string) (Q, bool) {
	if idx, ok := r.quoteHashToIndex[hash]; ok && idx < len(r.Quotes) {
		return r.Quotes[idx], true
	}
	var zero Q
	return zero, false
}

type QuoteQuery struct {
	Ctx                context.Context
	Conn               *Connection
	StartDate          time.Time
	EndDate            time.Time
	QuoteCollection    string
	RetainedCollection string
}

func extractHash(stored bson.D) (string, error) {
	var doc map[string]interface{}
	bsonBytes, err := bson.Marshal(stored)
	if err != nil {
		return "", err
	}
	if err := bson.Unmarshal(bsonBytes, &doc); err != nil {
		return "", err
	}
	hash, ok := doc["hash"].(string)
	if !ok || hash == "" {
		return "", nil
	}
	return hash, nil
}

func fetchQuotesByDateRange[Q any](
	ctx context.Context,
	conn *Connection,
	startDate, endDate time.Time,
	collectionName string,
	mapper func(bson.D) Q,
) ([]Q, []string, map[string]int, error) {
	quoteFilter := bson.D{
		{Key: "agreement_timestamp", Value: bson.D{
			{Key: "$gte", Value: startDate.Unix()},
			{Key: "$lte", Value: endDate.Unix()},
		}},
	}
	quoteCursor, err := conn.Collection(collectionName).Find(ctx, quoteFilter)
	if err != nil {
		return nil, nil, nil, err
	}
	var storedQuotes []bson.D
	if err = quoteCursor.All(ctx, &storedQuotes); err != nil {
		return nil, nil, nil, err
	}
	quoteHashToIndex := make(map[string]int)
	quoteHashes := make([]string, 0, len(storedQuotes))
	quotes := make([]Q, 0, len(storedQuotes))
	for i, stored := range storedQuotes {
		quotes = append(quotes, mapper(stored))
		hashVal, extractErr := extractHash(stored)
		if extractErr == nil && hashVal != "" {
			quoteHashToIndex[hashVal] = i
			quoteHashes = append(quoteHashes, hashVal)
		} else if extractErr != nil {
			log.Errorf("Error extracting hash: %v", extractErr)
		}
	}
	return quotes, quoteHashes, quoteHashToIndex, nil
}

func fetchRetainedQuotesByFilter[R any](
	ctx context.Context,
	conn *Connection,
	collectionName string,
	filter bson.D,
) ([]R, error) {
	retainedCursor, err := conn.Collection(collectionName).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var retainedQuotes []R
	if err = retainedCursor.All(ctx, &retainedQuotes); err != nil {
		return nil, err
	}
	return retainedQuotes, nil
}

func ListQuotesByDateRange[Q any, R QuoteHashProvider](
	query QuoteQuery,
	mapper func(bson.D) Q,
) (ListQuotesResult[Q, R], error) {
	dbCtx, cancel := context.WithTimeout(query.Ctx, query.Conn.timeout)
	defer cancel()
	result := ListQuotesResult[Q, R]{
		quoteHashToIndex: make(map[string]int),
	}
	quotes, quoteHashes, hashToIndex, err := fetchQuotesByDateRange(
		dbCtx, query.Conn, query.StartDate, query.EndDate,
		query.QuoteCollection, mapper,
	)
	if err != nil {
		return result, err
	}
	result.Quotes = quotes
	result.quoteHashToIndex = hashToIndex
	retainedFilter := createRetainedFilter(query.StartDate, query.EndDate, quoteHashes)
	retainedQuotes, err := fetchRetainedQuotesByFilter[R](
		dbCtx, query.Conn, query.RetainedCollection, retainedFilter,
	)
	if err != nil {
		result.Quotes = nil
		return result, err
	}
	result.RetainedQuotes = retainedQuotes
	additionalHashes := findAdditionalQuoteHashes(result.RetainedQuotes, quoteHashes)
	processAdditionalQuotes(dbCtx, query.Conn, query.QuoteCollection, additionalHashes, mapper, &result)
	logDbInteraction(Read, fmt.Sprintf("Found %d quotes and %d retained quotes in date range",
		len(result.Quotes), len(result.RetainedQuotes)))
	return result, nil
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
	existingMap := make(map[string]bool)
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

func processAdditionalQuotes[Q any, R QuoteHashProvider](
	ctx context.Context,
	conn *Connection,
	collectionName string,
	additionalHashes []string,
	mapper func(bson.D) Q,
	result *ListQuotesResult[Q, R],
) {
	if len(additionalHashes) == 0 {
		return
	}
	additionalFilter := bson.D{
		{Key: "hash", Value: bson.D{
			{Key: "$in", Value: additionalHashes},
		}},
	}
	additionalCursor, err := conn.Collection(collectionName).Find(ctx, additionalFilter)
	if err != nil {
		log.Errorf("Error fetching additional quotes: %v", err)
		return
	}
	var additionalStoredQuotes []bson.D
	if err = additionalCursor.All(ctx, &additionalStoredQuotes); err != nil {
		log.Errorf("Error reading additional quotes: %v", err)
		return
	}
	baseIndex := len(result.Quotes)
	for i, stored := range additionalStoredQuotes {
		result.Quotes = append(result.Quotes, mapper(stored))
		hash, err := extractHash(stored)
		if err == nil && hash != "" {
			result.quoteHashToIndex[hash] = baseIndex + i
		} else if err != nil {
			log.Errorf("Error extracting hash: %v", err)
		}
	}
}

type QuoteHashProvider interface {
	GetQuoteHash() string
}
