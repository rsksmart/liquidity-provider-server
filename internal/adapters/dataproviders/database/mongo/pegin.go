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
	logDbInteraction(Read, result)
	return result, nil
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

func (repo *peginMongoRepository) ListQuotesByDateRange(ctx context.Context, startDate, endDate time.Time) ([]quote.PeginQuoteWithRetained, error) {
	result := make([]quote.PeginQuoteWithRetained, 0)
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	quoteFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
		{Key: "$gte", Value: startDate.Unix()},
		{Key: "$lte", Value: endDate.Unix()},
	}}}
	findOpts := options.Find().SetSort(bson.D{{Key: "agreement_timestamp", Value: 1}})
	quoteCursor, err := repo.conn.Collection(PeginQuoteCollection).Find(dbCtx, quoteFilter, findOpts)
	if err != nil {
		return nil, err
	}
	defer func() {
		if quoteCursor != nil {
			if err := quoteCursor.Close(dbCtx); err != nil {
				log.Error("Error closing quote cursor: ", err)
			}
		}
	}()
	var storedQuotes []StoredPeginQuote
	if err = quoteCursor.All(dbCtx, &storedQuotes); err != nil {
		return nil, err
	}
	if len(storedQuotes) == 0 {
		logDbInteraction(Read, result)
		return result, nil
	}
	hashToIndex := make(map[string]int, len(storedQuotes))
	quoteHashes := make([]string, len(storedQuotes))
	result = make([]quote.PeginQuoteWithRetained, len(storedQuotes))
	for i, stored := range storedQuotes {
		quoteHashes[i] = stored.Hash
		hashToIndex[stored.Hash] = i
		result[i] = quote.PeginQuoteWithRetained{
			Quote:         stored.PeginQuote,
			RetainedQuote: nil,
		}
	}
	retainedCursor, err := repo.conn.Collection(RetainedPeginQuoteCollection).Find(
		dbCtx,
		bson.D{{Key: "quote_hash", Value: bson.D{{Key: "$in", Value: quoteHashes}}}},
	)
	if err != nil {
		return result, err
	}
	defer func() {
		if retainedCursor != nil {
			if err := retainedCursor.Close(dbCtx); err != nil {
				log.Error("Error closing retained cursor: ", err)
			}
		}
	}()
	var retainedQuote quote.RetainedPeginQuote
	for retainedCursor.Next(dbCtx) {
		if err := retainedCursor.Decode(&retainedQuote); err != nil {
			return result, err
		}
		if idx, exists := hashToIndex[retainedQuote.QuoteHash]; exists {
			quoteCopy := retainedQuote
			result[idx].RetainedQuote = &quoteCopy
		}
	}
	if err := retainedCursor.Err(); err != nil {
		return result, err
	}
	logDbInteraction(Read, len(result))
	return result, nil
}
