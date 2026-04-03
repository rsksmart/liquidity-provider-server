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

// peginQuoteWithRetainedAggDoc decodes pegin quote aggregation rows after $lookup and $unwind on "retained"
// (e.g. ListQuotesByDateRange, GetQuotesWithRetainedByStateAndDate).
type peginQuoteWithRetainedAggDoc struct {
	StoredPeginQuote `bson:",inline"`
	Retained         quote.RetainedPeginQuote `bson:"retained"`
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

	dateMatch := repo.listQuotesByDateRangeDateMatch(startDate, endDate)
	total, err := repo.countQuotesWithRetainedInDateRange(dbCtx, dateMatch)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		empty := make([]quote.PeginQuoteWithRetained, 0)
		logDbInteraction(Read, empty)
		return empty, 0, nil
	}

	dataPipeline := append(repo.listQuotesByDateRangePipelinePrefix(dateMatch),
		bson.D{{Key: "$unwind", Value: "$retained"}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "agreement_timestamp", Value: SortAscending}}}},
	)
	if page > 0 && perPage > 0 {
		skip := int64((page - 1) * perPage)
		dataPipeline = append(dataPipeline,
			bson.D{{Key: "$skip", Value: skip}},
			bson.D{{Key: "$limit", Value: int64(perPage)}},
		)
	}

	dataCursor, err := repo.conn.Collection(PeginQuoteCollection).Aggregate(dbCtx, dataPipeline)
	if err != nil {
		return nil, 0, err
	}
	defer dataCursor.Close(dbCtx)

	var result []quote.PeginQuoteWithRetained
	for dataCursor.Next(dbCtx) {
		var doc peginQuoteWithRetainedAggDoc
		if err := dataCursor.Decode(&doc); err != nil {
			return nil, 0, err
		}
		doc.Retained.FillZeroValues()
		result = append(result, quote.PeginQuoteWithRetained{
			Quote:         doc.PeginQuote,
			RetainedQuote: doc.Retained,
		})
	}
	if err := dataCursor.Err(); err != nil {
		return nil, 0, err
	}

	logDbInteraction(Read, len(result))
	return result, total, nil
}

func (repo *peginMongoRepository) listQuotesByDateRangeDateMatch(startDate, endDate time.Time) bson.D {
	return bson.D{{Key: "agreement_timestamp", Value: bson.D{
		{Key: "$gte", Value: startDate.Unix()},
		{Key: "$lte", Value: endDate.Unix()},
	}}}
}

func (repo *peginMongoRepository) listQuotesByDateRangeLookupStage() bson.D {
	return bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: RetainedPeginQuoteCollection},
		{Key: "localField", Value: "hash"},
		{Key: "foreignField", Value: "quote_hash"},
		{Key: "as", Value: "retained"},
	}}}
}

func (repo *peginMongoRepository) listQuotesByDateRangeHasRetainedMatch() bson.D {
	return bson.D{{Key: "$match", Value: bson.D{
		{Key: "retained.0", Value: bson.D{{Key: "$exists", Value: true}}},
	}}}
}

// listQuotesByDateRangePipelinePrefix is shared by the count and data aggregations so both apply the same filters.
func (repo *peginMongoRepository) listQuotesByDateRangePipelinePrefix(dateMatch bson.D) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: dateMatch}},
		repo.listQuotesByDateRangeLookupStage(),
		repo.listQuotesByDateRangeHasRetainedMatch(),
	}
}

func (repo *peginMongoRepository) countQuotesWithRetainedInDateRange(ctx context.Context, dateMatch bson.D) (int, error) {
	pipeline := append(repo.listQuotesByDateRangePipelinePrefix(dateMatch),
		bson.D{{Key: "$count", Value: "total"}},
	)
	cursor, err := repo.conn.Collection(PeginQuoteCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var out struct {
		Total int `bson:"total"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&out); err != nil {
			return 0, err
		}
	}
	if err := cursor.Err(); err != nil {
		return 0, err
	}
	return out.Total, nil
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

// GetQuotesWithRetainedByStateAndDate retrieves pegin quotes filtered by state and date range,
// optionally joined with their retained data.
//
// IMPORTANT: This method may return quotes WITHOUT retained data (non-accepted quotes).
// The aggregation pipeline includes quotes that have no matching RetainedPeginQuote
// record. For these quotes, the RetainedQuote field will be a zero-valued struct:
//   - QuoteHash: "" (empty string)
//   - DepositAddress: ""
//   - State: ""
//   - All numeric fields: 0
//   - All Wei pointers: set to NewWei(0) by FillZeroValues()
//
// The states parameter filters by RetainedQuote state when retained data exists, or
// includes quotes without retained data regardless of the provided states.
//
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
		// Stage 4: Filter by state (or include quotes without retained data)
		{{Key: "$match", Value: bson.M{
			"$or": []bson.M{
				{"retained.state": bson.M{"$in": states}},
				{"retained.state": nil}, // Include quotes without retained quote (non-accepted)
			},
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
		var doc peginQuoteWithRetainedAggDoc
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
