package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	PeginQuoteCollection         = "peginQuote"
	RetainedPeginQuoteCollection = "retainedPeginQuote"
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

func (repo *peginMongoRepository) InsertQuote(ctx context.Context, hash string, peginQuote quote.PeginQuote) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	collection := repo.conn.Collection(PeginQuoteCollection)
	storedQuote := StoredPeginQuote{
		PeginQuote: peginQuote,
		Hash:       hash,
	}
	_, err := collection.InsertOne(dbCtx, storedQuote)
	if err != nil {
		return err
	} else {
		logDbInteraction(Insert, storedQuote)
		return nil
	}
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

func (repo *peginMongoRepository) GetQuotes(ctx context.Context, hashes []string) ([]quote.PeginQuote, error) {
	var result StoredPeginQuote
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	for _, hash := range hashes {
		if err := quote.ValidateQuoteHash(hash); err != nil {
			return nil, err
		}
	}

	collection := repo.conn.Collection(PeginQuoteCollection)
	filter := bson.M{"hash": bson.M{"$in": hashes}}

	quotesReturn := make([]quote.PeginQuote, 0)

	cursor, err := collection.Find(dbCtx, filter)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		err := cursor.Decode(&result)
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
	logDbInteraction(Delete, fmt.Sprintf("removed %d records from %s collection", peginResult.DeletedCount, PeginQuoteCollection))
	logDbInteraction(Delete, fmt.Sprintf("removed %d records from %s collection", retainedResult.DeletedCount, RetainedPeginQuoteCollection))
	if peginResult.DeletedCount != retainedResult.DeletedCount {
		return 0, errors.New("pegin quote collections didn't match")
	}
	return uint(peginResult.DeletedCount + retainedResult.DeletedCount), nil
}
