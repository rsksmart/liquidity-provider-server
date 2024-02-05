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
	peginQuoteCollection         = "peginQuote"
	retainedPeginQuoteCollection = "retainedPeginQuote"
)

type peginMongoRepository struct {
	conn *Connection
}

func NewPeginMongoRepository(conn *Connection) quote.PeginQuoteRepository {
	return &peginMongoRepository{conn: conn}
}

type storedPeginQuote struct {
	quote.PeginQuote `bson:",inline"`
	Hash             string `json:"hash" bson:"hash"`
}

func (repo *peginMongoRepository) InsertQuote(ctx context.Context, hash string, peginQuote quote.PeginQuote) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	collection := repo.conn.Collection(peginQuoteCollection)
	storedQuote := storedPeginQuote{
		PeginQuote: peginQuote,
		Hash:       hash,
	}
	_, err := collection.InsertOne(dbCtx, storedQuote)
	if err != nil {
		return err
	} else {
		logDbInteraction(insert, storedQuote)
		return nil
	}
}

func (repo *peginMongoRepository) GetQuote(ctx context.Context, hash string) (*quote.PeginQuote, error) {
	var result storedPeginQuote
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	collection := repo.conn.Collection(peginQuoteCollection)
	filter := bson.D{primitive.E{Key: "hash", Value: hash}}

	err := collection.FindOne(dbCtx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	logDbInteraction(read, result.PeginQuote)
	return &result.PeginQuote, nil
}

func (repo *peginMongoRepository) GetRetainedQuote(ctx context.Context, hash string) (*quote.RetainedPeginQuote, error) {
	var result quote.RetainedPeginQuote
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	collection := repo.conn.Collection(retainedPeginQuoteCollection)
	filter := bson.D{primitive.E{Key: "quote_hash", Value: hash}}

	err := collection.FindOne(dbCtx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	logDbInteraction(read, result)
	return &result, nil
}

func (repo *peginMongoRepository) InsertRetainedQuote(ctx context.Context, retainedQuote quote.RetainedPeginQuote) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	collection := repo.conn.Collection(retainedPeginQuoteCollection)
	_, err := collection.InsertOne(dbCtx, retainedQuote)
	if err != nil {
		return err
	} else {
		logDbInteraction(insert, retainedQuote)
		return nil
	}
}

func (repo *peginMongoRepository) UpdateRetainedQuote(ctx context.Context, retainedQuote quote.RetainedPeginQuote) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	collection := repo.conn.Collection(retainedPeginQuoteCollection)
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
	logDbInteraction(update, retainedQuote)
	return nil
}

func (repo *peginMongoRepository) GetRetainedQuoteByState(ctx context.Context, states ...quote.PeginState) ([]quote.RetainedPeginQuote, error) {
	result := make([]quote.RetainedPeginQuote, 0)
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	collection := repo.conn.Collection(retainedPeginQuoteCollection)
	query := bson.D{primitive.E{Key: "state", Value: bson.D{primitive.E{Key: "$in", Value: states}}}}
	rows, err := collection.Find(dbCtx, query)
	if err != nil {
		return nil, err
	}
	if err = rows.All(ctx, &result); err != nil {
		return nil, err
	}
	logDbInteraction(read, result)
	return result, nil
}

func (repo *peginMongoRepository) DeleteQuotes(ctx context.Context, quotes []string) (uint, error) {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout*2)
	defer cancel()

	filter := bson.D{primitive.E{Key: "quote_hash", Value: bson.D{primitive.E{Key: "$in", Value: quotes}}}}
	peginResult, err := repo.conn.Collection(peginQuoteCollection).DeleteMany(dbCtx, filter)
	if err != nil {
		return 0, err
	}
	retainedResult, err := repo.conn.Collection(retainedPeginQuoteCollection).DeleteMany(dbCtx, filter)
	if err != nil {
		return 0, err
	} else if peginResult.DeletedCount != retainedResult.DeletedCount {
		return 0, errors.New("pegin quote collections didn't match")
	}
	logDbInteraction(delete, fmt.Sprintf("removed %d records from %s collection", peginResult.DeletedCount, pegoutQuoteCollection))
	logDbInteraction(delete, fmt.Sprintf("removed %d records from %s collection", retainedResult.DeletedCount, retainedPegoutQuoteCollection))
	return uint(peginResult.DeletedCount + retainedResult.DeletedCount), nil
}
