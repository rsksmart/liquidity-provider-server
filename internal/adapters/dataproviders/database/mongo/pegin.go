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

func (repo *peginMongoRepository) ListQuotesByDateRange(ctx context.Context, startDate, endDate time.Time) (quote.PeginQuoteResult, error) {
	result := quote.PeginQuoteResult{
		Quotes:         make([]quote.PeginQuote, 0),
		RetainedQuotes: make([]quote.RetainedPeginQuote, 0),
	}
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	startTimestamp := uint32(startDate.Unix())
	endTimestamp := uint32(endDate.Unix())
	quoteCollection := repo.conn.Collection(PeginQuoteCollection)
	quoteFilter := bson.D{
		primitive.E{
			Key: "agreement_timestamp",
			Value: bson.D{
				primitive.E{Key: "$gte", Value: startTimestamp},
				primitive.E{Key: "$lte", Value: endTimestamp},
			},
		},
	}
	quoteCursor, err := quoteCollection.Find(dbCtx, quoteFilter)
	if err != nil {
		return result, err
	}
	defer quoteCursor.Close(dbCtx)
	var storedQuotes []StoredPeginQuote
	if err = quoteCursor.All(dbCtx, &storedQuotes); err != nil {
		return result, err
	}
	quoteHashes := make([]string, 0, len(storedQuotes))
	for _, stored := range storedQuotes {
		result.Quotes = append(result.Quotes, stored.PeginQuote)
		quoteHashes = append(quoteHashes, stored.Hash)
	}
	if len(quoteHashes) > 0 {
		retainedCollection := repo.conn.Collection(RetainedPeginQuoteCollection)
		retainedFilter := bson.D{
			primitive.E{
				Key: "quote_hash",
				Value: bson.D{
					primitive.E{Key: "$in", Value: quoteHashes},
				},
			},
		}
		retainedCursor, err := retainedCollection.Find(dbCtx, retainedFilter)
		if err != nil {
			return result, err
		}
		defer retainedCursor.Close(dbCtx)
		if err = retainedCursor.All(dbCtx, &result.RetainedQuotes); err != nil {
			return result, err
		}
	}
	logDbInteraction(Read, result)
	return result, nil
}
