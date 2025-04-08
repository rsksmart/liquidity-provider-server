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

// QuoteResult holds the result of a quote listing operation
type QuoteResult[Q any, R QuoteHashProvider] struct {
	Quotes         []Q
	RetainedQuotes []R
	Error          error
}

// QuoteQuery represents parameters for querying quotes
type QuoteQuery struct {
	Ctx                context.Context
	Conn               *Connection
	StartDate          time.Time
	EndDate            time.Time
	QuoteCollection    string
	RetainedCollection string
}

// ListQuotesByDateRange retrieves quotes and retained quotes within a date range
func ListQuotesByDateRange[Q any, R QuoteHashProvider](
	query QuoteQuery,
	mapper func(bson.D) Q,
) QuoteResult[Q, R] {
	dbCtx, cancel := context.WithTimeout(query.Ctx, query.Conn.timeout)
	defer cancel()

	// Step 1: Fetch quotes by date range
	quotes, quoteHashes, err := fetchQuotesByDateRange(dbCtx, query.Conn, query.StartDate, query.EndDate, query.QuoteCollection, mapper)
	if err != nil {
		return QuoteResult[Q, R]{Error: err}
	}

	// Step 2: Fetch retained quotes
	retainedQuotes, additionalHashes, err := fetchRetainedQuotes[R](dbCtx, query.Conn, query.StartDate, query.EndDate, query.RetainedCollection, quoteHashes)
	if err != nil {
		return QuoteResult[Q, R]{Error: err}
	}

	// Step 3: Fetch any additional quotes referenced by retained quotes
	if len(additionalHashes) > 0 {
		additionalQuotes, err := fetchAdditionalQuotes(dbCtx, query.Conn, query.QuoteCollection, additionalHashes, mapper)
		if err != nil {
			log.Errorf("Error processing additional quotes: %v", err)
		} else {
			quotes = append(quotes, additionalQuotes...)
		}
	}

	logDbInteraction(Read, fmt.Sprintf("Found %d quotes and %d retained quotes in date range",
		len(quotes), len(retainedQuotes)))

	return QuoteResult[Q, R]{
		Quotes:         quotes,
		RetainedQuotes: retainedQuotes,
		Error:          nil,
	}
}

// fetchQuotesByDateRange retrieves quotes from the database within the specified date range
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

		// Extract hash from the BSON document
		hashValue, ok := getStringValueFromBSON(stored, "hash")
		if ok {
			quoteHashes = append(quoteHashes, hashValue)
		}
	}

	return quotes, quoteHashes, nil
}

// getStringValueFromBSON extracts a string value from a BSON document by key
func getStringValueFromBSON(doc bson.D, key string) (string, bool) {
	// Convert bson.D to bson.Raw for direct lookup
	data, err := bson.Marshal(doc)
	if err != nil {
		return "", false
	}

	// Use the Raw.Lookup method to find the value
	rawValue := bson.Raw(data).Lookup(key)

	// Extract string value if possible
	return rawValue.StringValueOK()
}

// QuoteHashProvider defines an interface for objects that can provide a hash
type QuoteHashProvider interface {
	GetQuoteHash() string
}

// fetchRetainedQuotes retrieves retained quotes and identifies any additional quote hashes to fetch
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

// createRetainedFilter creates a filter for retained quotes
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

// findAdditionalQuoteHashes identifies quote hashes in retained quotes that are not in existingQuoteHashes
func findAdditionalQuoteHashes[R QuoteHashProvider](retainedQuotes []R, existingQuoteHashes []string) []string {
	// Create a map for faster lookup
	existingMap := make(map[string]bool, len(existingQuoteHashes))
	for _, hash := range existingQuoteHashes {
		existingMap[hash] = true
	}

	// Store additional hashes
	additionalMap := make(map[string]bool)

	for i := range retainedQuotes {
		hash := retainedQuotes[i].GetQuoteHash()
		if !existingMap[hash] {
			additionalMap[hash] = true
		}
	}

	// Convert the map to a slice
	additionalHashes := make([]string, 0, len(additionalMap))
	for hash := range additionalMap {
		additionalHashes = append(additionalHashes, hash)
	}

	return additionalHashes
}

// fetchAdditionalQuotes retrieves quotes by their hash
func fetchAdditionalQuotes[Q any](
	ctx context.Context,
	conn *Connection,
	collectionName string,
	hashes []string,
	mapper func(bson.D) Q,
) ([]Q, error) {
	quoteFilter := bson.D{
		{Key: "hash", Value: bson.D{
			{Key: "$in", Value: hashes},
		}},
	}

	var storedQuotes []bson.D
	quoteCursor, err := conn.Collection(collectionName).Find(ctx, quoteFilter)
	if err != nil {
		return nil, err
	}

	if err = quoteCursor.All(ctx, &storedQuotes); err != nil {
		return nil, err
	}

	quotes := make([]Q, 0, len(storedQuotes))
	for _, stored := range storedQuotes {
		quoteObj := mapper(stored)
		quotes = append(quotes, quoteObj)
	}

	return quotes, nil
}
