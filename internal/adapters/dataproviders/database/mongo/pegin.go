package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	PeginQuoteCollection         = "peginQuote"
	RetainedPeginQuoteCollection = "retainedPeginQuote"
	PeginCreationDataCollection  = "peginQuoteCreationData"
)

type peginMongoRepository struct {
	conn *Connection
}

func NewPeginMongoRepository(conn *Connection) quote.PeginQuoteRepository {
	return &peginMongoRepository{conn: conn}
}

type StoredPeginQuote struct {
	quote.PeginQuote `bson:",inline"`
	Hash             string `json:"hash" bson:"hash"`
}

type StoredPeginCreationData struct {
	quote.PeginCreationData `bson:",inline"`
	Hash                    string `json:"hash" bson:"hash"`
}

func (repo *peginMongoRepository) InsertQuote(ctx context.Context, createdQuote quote.CreatedPeginQuote) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	collection := repo.conn.Collection(PeginQuoteCollection)
	storedQuote := StoredPeginQuote{
		PeginQuote: createdQuote.Quote,
		Hash:       createdQuote.Hash,
	}
	_, err := collection.InsertOne(dbCtx, storedQuote)
	if err != nil {
		return err
	}
	logDbInteraction(Insert, storedQuote)

	collection = repo.conn.Collection(PeginCreationDataCollection)
	storedCreationData := StoredPeginCreationData{
		PeginCreationData: createdQuote.CreationData,
		Hash:              createdQuote.Hash,
	}
	_, err = collection.InsertOne(dbCtx, storedCreationData)
	if err != nil {
		return err
	}
	logDbInteraction(Insert, storedCreationData)
	return nil
}

func (repo *peginMongoRepository) GetPeginCreationData(ctx context.Context, hash string) quote.PeginCreationData {
	var result StoredPeginCreationData
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	if err := quote.ValidateQuoteHash(hash); err != nil {
		log.Error("Invalid hash. Returning default pegin creation data")
		return quote.PeginCreationDataZeroValue()
	}

	collection := repo.conn.Collection(PeginCreationDataCollection)
	filter := bson.D{primitive.E{Key: "hash", Value: hash}}

	err := collection.FindOne(dbCtx, filter).Decode(&result)
	if err != nil {
		log.Error("Hash not found. Returning default pegin creation data")
		return quote.PeginCreationDataZeroValue()
	}
	logDbInteraction(Read, result.PeginCreationData)
	return result.PeginCreationData
}

func (repo *peginMongoRepository) GetQuote(ctx context.Context, hash string) (*quote.PeginQuote, error) {
	var result StoredPeginQuote
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	if err := quote.ValidateQuoteHash(hash); err != nil {
		return nil, err
	}

	collection := repo.conn.Collection(PeginQuoteCollection)
	filter := bson.D{primitive.E{Key: "hash", Value: hash}}

	err := collection.FindOne(dbCtx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	logDbInteraction(Read, result.PeginQuote)
	return &result.PeginQuote, nil
}

func (repo *peginMongoRepository) GetQuotesByHashesAndDate(
	ctx context.Context,
	hashes []string,
	startDate,
	endDate time.Time,
) ([]quote.PeginQuote, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	for _, hash := range hashes {
		if err := quote.ValidateQuoteHash(hash); err != nil {
			return nil, err
		}
	}

	collection := repo.conn.Collection(PeginQuoteCollection)
	filter := bson.M{
		"hash": bson.M{"$in": hashes},
		"agreement_timestamp": bson.M{
			"$gte": startDate.Unix(),
			"$lte": endDate.Unix(),
		},
	}

	quotesReturn := make([]quote.PeginQuote, 0)

	cursor, err := collection.Find(dbCtx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var result StoredPeginQuote
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		quotesReturn = append(quotesReturn, result.PeginQuote)
	}
	logDbInteraction(Read, quotesReturn)
	return quotesReturn, nil
}

func (repo *peginMongoRepository) GetRetainedQuote(ctx context.Context, hash string) (*quote.RetainedPeginQuote, error) {
	var result quote.RetainedPeginQuote
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	if err := quote.ValidateQuoteHash(hash); err != nil {
		return nil, err
	}

	collection := repo.conn.Collection(RetainedPeginQuoteCollection)
	filter := bson.D{primitive.E{Key: "quote_hash", Value: hash}}

	err := collection.FindOne(dbCtx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	result.FillZeroValues()
	logDbInteraction(Read, result)
	return &result, nil
}

func (repo *peginMongoRepository) InsertRetainedQuote(ctx context.Context, retainedQuote quote.RetainedPeginQuote) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	collection := repo.conn.Collection(RetainedPeginQuoteCollection)
	_, err := collection.InsertOne(dbCtx, retainedQuote)
	if err != nil {
		return err
	} else {
		logDbInteraction(Insert, retainedQuote)
		return nil
	}
}

func (repo *peginMongoRepository) UpdateRetainedQuote(ctx context.Context, retainedQuote quote.RetainedPeginQuote) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	collection := repo.conn.Collection(RetainedPeginQuoteCollection)
	filter := bson.D{primitive.E{Key: "quote_hash", Value: retainedQuote.QuoteHash}}
	updateStatement := bson.D{primitive.E{Key: "$set", Value: retainedQuote}}

	result, err := collection.UpdateOne(dbCtx, filter, updateStatement)
	if err != nil {
		return err
	} else if result.ModifiedCount == 0 {
		return usecases.QuoteNotFoundError
	} else if result.ModifiedCount > 1 {
		return errors.New("multiple documents updated")
	}
	logDbInteraction(Update, retainedQuote)
	return nil
}

func (repo *peginMongoRepository) GetRetainedQuoteByState(ctx context.Context, states ...quote.PeginState) ([]quote.RetainedPeginQuote, error) {
	result := make([]quote.RetainedPeginQuote, 0)
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	collection := repo.conn.Collection(RetainedPeginQuoteCollection)
	query := bson.D{primitive.E{Key: "state", Value: bson.D{primitive.E{Key: "$in", Value: states}}}}
	rows, err := collection.Find(dbCtx, query)
	if err != nil {
		return nil, err
	}
	if err = rows.All(ctx, &result); err != nil {
		return nil, err
	}
	for i := range result {
		result[i].FillZeroValues()
	}
	logDbInteraction(Read, result)
	return result, nil
}

func (repo *peginMongoRepository) GetQuotesByState(ctx context.Context, states ...quote.PeginState) ([]quote.PeginQuote, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	retainedQuotes, err := repo.GetRetainedQuoteByState(ctx, states...)
	if err != nil {
		return nil, err
	}

	if len(retainedQuotes) == 0 {
		result := make([]quote.PeginQuote, 0)
		logDbInteraction(Read, result)
		return result, nil
	}

	quoteHashes := make([]string, len(retainedQuotes))
	for i, rq := range retainedQuotes {
		quoteHashes[i] = rq.QuoteHash
	}

	storedQuotes, err := repo.fetchQuotesByHashes(dbCtx, quoteHashes)
	if err != nil {
		return nil, err
	}

	result := make([]quote.PeginQuote, len(storedQuotes))
	for i, stored := range storedQuotes {
		result[i] = stored.PeginQuote
	}

	logDbInteraction(Read, len(result))
	return result, nil
}

func (repo *peginMongoRepository) fetchQuotesByHashes(ctx context.Context, quoteHashes []string) ([]StoredPeginQuote, error) {
	collection := repo.conn.Collection(PeginQuoteCollection)
	filter := bson.D{{Key: "hash", Value: bson.D{{Key: "$in", Value: quoteHashes}}}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var storedQuotes []StoredPeginQuote
	if err = cursor.All(ctx, &storedQuotes); err != nil {
		return nil, err
	}

	return storedQuotes, nil
}

func (repo *peginMongoRepository) DeleteQuotes(ctx context.Context, quotes []string) (uint, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout*2)
	defer cancel()

	quoteFilter := bson.D{primitive.E{Key: "hash", Value: bson.D{primitive.E{Key: "$in", Value: quotes}}}}
	retainedFilter := bson.D{primitive.E{Key: "quote_hash", Value: bson.D{primitive.E{Key: "$in", Value: quotes}}}}
	peginResult, err := repo.conn.Collection(PeginQuoteCollection).DeleteMany(dbCtx, quoteFilter)
	if err != nil {
		return 0, err
	}
	retainedResult, err := repo.conn.Collection(RetainedPeginQuoteCollection).DeleteMany(dbCtx, retainedFilter)
	if err != nil {
		return 0, err
	}
	creationDataResult, err := repo.conn.Collection(PeginCreationDataCollection).DeleteMany(dbCtx, quoteFilter)
	if err != nil {
		return 0, err
	}
	const msgTemplate = "removed %d records from %s collection"
	logDbInteraction(Delete, fmt.Sprintf(msgTemplate, peginResult.DeletedCount, PeginQuoteCollection))
	logDbInteraction(Delete, fmt.Sprintf(msgTemplate, retainedResult.DeletedCount, RetainedPeginQuoteCollection))
	logDbInteraction(Delete, fmt.Sprintf(msgTemplate, creationDataResult.DeletedCount, PegoutCreationDataCollection))
	// creation data doesn't count for mismatch because not all the quotes have it
	if peginResult.DeletedCount != retainedResult.DeletedCount {
		return 0, errors.New("pegin quote collections didn't match")
	}
	return uint(peginResult.DeletedCount + retainedResult.DeletedCount + creationDataResult.DeletedCount), nil
}

func (repo *peginMongoRepository) ListQuotesByDateRange(ctx context.Context, startDate, endDate time.Time, page, perPage int) ([]quote.PeginQuoteWithRetained, int, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	// Fetch quotes with pagination
	storedQuotes, err := repo.fetchQuotesByDateRange(dbCtx, startDate, endDate, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	if len(storedQuotes) == 0 {
		result := make([]quote.PeginQuoteWithRetained, 0)
		logDbInteraction(Read, result)
		return result, 0, nil
	}

	// Build initial result structure with quotes
	result, quoteHashes := repo.buildQuoteResults(storedQuotes)

	// Fetch and merge retained quotes
	if err := repo.mergeRetainedQuotes(dbCtx, result, quoteHashes); err != nil {
		return result, len(result), err
	}

	logDbInteraction(Read, len(result))
	return result, len(result), nil
}

func (repo *peginMongoRepository) fetchQuotesByDateRange(ctx context.Context, startDate, endDate time.Time, page, perPage int) ([]StoredPeginQuote, error) {
	quoteFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
		{Key: "$gte", Value: startDate.Unix()},
		{Key: "$lte", Value: endDate.Unix()},
	}}}

	findOpts := repo.buildFindOptions(page, perPage)

	quoteCursor, err := repo.conn.Collection(PeginQuoteCollection).Find(ctx, quoteFilter, findOpts)
	if err != nil {
		return nil, err
	}

	var storedQuotes []StoredPeginQuote
	if err = quoteCursor.All(ctx, &storedQuotes); err != nil {
		return nil, err
	}

	return storedQuotes, nil
}

func (repo *peginMongoRepository) buildFindOptions(page, perPage int) *options.FindOptions {
	findOpts := options.Find().SetSort(bson.D{{Key: "agreement_timestamp", Value: SortAscending}})

	// Apply pagination if page and perPage are provided (not 0)
	// When page=0 and perPage=0, return all results
	if page > 0 && perPage > 0 {
		skip := (page - 1) * perPage
		findOpts.SetSkip(int64(skip)).SetLimit(int64(perPage))
	}

	return findOpts
}

func (repo *peginMongoRepository) buildQuoteResults(storedQuotes []StoredPeginQuote) ([]quote.PeginQuoteWithRetained, []string) {
	result := make([]quote.PeginQuoteWithRetained, len(storedQuotes))
	quoteHashes := make([]string, len(storedQuotes))

	for i, stored := range storedQuotes {
		quoteHashes[i] = stored.Hash
		result[i] = quote.PeginQuoteWithRetained{
			Quote:         stored.PeginQuote,
			RetainedQuote: quote.RetainedPeginQuote{},
		}
	}

	return result, quoteHashes
}

func (repo *peginMongoRepository) mergeRetainedQuotes(ctx context.Context, result []quote.PeginQuoteWithRetained, quoteHashes []string) error {
	retainedQuotes, err := repo.fetchRetainedQuotes(ctx, quoteHashes)
	if err != nil {
		return err
	}

	// Create hash to index mapping for efficient lookup
	hashToIndex := make(map[string]int, len(quoteHashes))
	for i, hash := range quoteHashes {
		hashToIndex[hash] = i
	}

	// Merge retained quotes into result
	for _, retainedQuote := range retainedQuotes {
		if idx, exists := hashToIndex[retainedQuote.QuoteHash]; exists {
			result[idx].RetainedQuote = retainedQuote
		}
	}

	return nil
}

func (repo *peginMongoRepository) fetchRetainedQuotes(ctx context.Context, quoteHashes []string) ([]quote.RetainedPeginQuote, error) {
	retainedCursor, err := repo.conn.Collection(RetainedPeginQuoteCollection).Find(
		ctx,
		bson.D{{Key: "quote_hash", Value: bson.D{{Key: "$in", Value: quoteHashes}}}},
	)
	if err != nil {
		return nil, err
	}

	var retainedQuotes []quote.RetainedPeginQuote
	if err = retainedCursor.All(ctx, &retainedQuotes); err != nil {
		return nil, err
	}

	for i := range retainedQuotes {
		retainedQuotes[i].FillZeroValues()
	}

	return retainedQuotes, nil
}

func (repo *peginMongoRepository) GetRetainedQuotesForAddress(ctx context.Context, address string, states ...quote.PeginState) ([]quote.RetainedPeginQuote, error) {
	result := make([]quote.RetainedPeginQuote, 0)
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	collection := repo.conn.Collection(RetainedPeginQuoteCollection)
	filter := bson.D{
		primitive.E{Key: "owner_account_address", Value: address},
		primitive.E{Key: "state", Value: bson.D{
			primitive.E{Key: "$in", Value: states},
		}},
	}

	rows, err := collection.Find(dbCtx, filter)
	if err != nil {
		return nil, err
	}
	if err = rows.All(ctx, &result); err != nil {
		return nil, err
	}
	for i := range result {
		result[i].FillZeroValues()
	}
	logDbInteraction(Read, result)
	return result, nil
}

// TODO: add pagination to this method
func (repo *peginMongoRepository) GetQuotesWithRetainedByStateAndDate(ctx context.Context, states []quote.PeginState, startDate, endDate time.Time) ([]quote.PeginQuoteWithRetained, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	collection := repo.conn.Collection(PeginQuoteCollection)
	const maxRecordLimit = 500

	// Single aggregation pipeline that:
	// 1. Filters by date (most selective filter)
	// 2. Joins with retained quotes
	// 3. Filters by state
	// 4. Limits results to protect server
	pipeline := mongo.Pipeline{
		// Stage 1: Filter by date (most selective)
		{{Key: "$match", Value: bson.M{
			"agreement_timestamp": bson.M{
				"$gte": startDate.Unix(),
				"$lte": endDate.Unix(),
			},
		}}},
		// Stage 2: Lookup retained quotes
		{{Key: "$lookup", Value: bson.M{
			"from":         RetainedPeginQuoteCollection,
			"localField":   "hash",
			"foreignField": "quote_hash",
			"as":           "retained",
		}}},
		// Stage 3: Unwind retained array (should have 0 or 1 element)
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$retained",
			"preserveNullAndEmptyArrays": true, // Keep all quotes, even without retained data
		}}},
		// Stage 4: Filter by state
		{{Key: "$match", Value: bson.M{
			"retained.state": bson.M{"$in": states},
		}}},
		// Stage 5: Limit to maxRecordLimit + 1 to detect if we exceeded
		{{Key: "$limit", Value: maxRecordLimit + 1}},
	}

	cursor, err := collection.Aggregate(dbCtx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(dbCtx)

	result := make([]quote.PeginQuoteWithRetained, 0)
	for cursor.Next(dbCtx) {
		var doc struct {
			StoredPeginQuote `bson:",inline"`
			Retained         quote.RetainedPeginQuote `bson:"retained"`
		}

		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		doc.Retained.FillZeroValues()
		result = append(result, quote.PeginQuoteWithRetained{
			Quote:         doc.PeginQuote,
			RetainedQuote: doc.Retained,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	if len(result) > maxRecordLimit {
		return nil, errors.New("dataset too large, please try a shorter time range")
	}
	logDbInteraction(Read, len(result))
	return result, nil
}
