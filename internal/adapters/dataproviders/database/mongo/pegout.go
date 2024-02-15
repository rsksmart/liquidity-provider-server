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
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	pegoutQuoteCollection         = "pegoutQuote"
	retainedPegoutQuoteCollection = "retainedPegoutQuote"
	depositEventsCollection       = "depositEvents"
)

type storedPegoutQuote struct {
	quote.PegoutQuote `bson:",inline"`
	Hash              string `json:"hash" bson:"hash"`
}

type pegoutMongoRepository struct {
	conn *Connection
}

func NewPegoutMongoRepository(conn *Connection) quote.PegoutQuoteRepository {
	return &pegoutMongoRepository{conn: conn}
}

func (repo *pegoutMongoRepository) InsertQuote(ctx context.Context, hash string, pegoutQuote quote.PegoutQuote) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	collection := repo.conn.Collection(pegoutQuoteCollection)
	storedQuote := storedPegoutQuote{
		PegoutQuote: pegoutQuote,
		Hash:        hash,
	}
	_, err := collection.InsertOne(dbCtx, storedQuote)
	if err != nil {
		return err
	} else {
		logDbInteraction(insert, storedQuote)
		return nil
	}
}

func (repo *pegoutMongoRepository) GetQuote(ctx context.Context, hash string) (*quote.PegoutQuote, error) {
	var result storedPegoutQuote
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	collection := repo.conn.Collection(pegoutQuoteCollection)
	filter := bson.D{primitive.E{Key: "hash", Value: hash}}

	err := collection.FindOne(dbCtx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	logDbInteraction(read, result.PegoutQuote)
	return &result.PegoutQuote, nil
}

func (repo *pegoutMongoRepository) GetRetainedQuote(ctx context.Context, hash string) (*quote.RetainedPegoutQuote, error) {
	var result quote.RetainedPegoutQuote
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	collection := repo.conn.Collection(retainedPegoutQuoteCollection)
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

func (repo *pegoutMongoRepository) InsertRetainedQuote(ctx context.Context, retainedQuote quote.RetainedPegoutQuote) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	collection := repo.conn.Collection(retainedPegoutQuoteCollection)
	_, err := collection.InsertOne(dbCtx, retainedQuote)
	if err != nil {
		return err
	} else {
		logDbInteraction(insert, retainedQuote)
		return nil
	}
}

func (repo *pegoutMongoRepository) ListPegoutDepositsByAddress(ctx context.Context, address string) ([]quote.PegoutDeposit, error) {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	filter := bson.M{"from": address}
	sort := options.Find().SetSort(bson.M{"timestamp": -1})
	cursor, err := repo.conn.Collection(depositEventsCollection).Find(dbCtx, filter, sort)
	if err != nil {
		return make([]quote.PegoutDeposit, 0), err
	}

	var documents []quote.PegoutDeposit
	if err = cursor.All(ctx, &documents); err != nil {
		return make([]quote.PegoutDeposit, 0), err
	}
	logDbInteraction(read, fmt.Sprintf("%d pegout deposits", len(documents)))
	return documents, nil
}

func (repo *pegoutMongoRepository) UpdateRetainedQuote(ctx context.Context, retainedQuote quote.RetainedPegoutQuote) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	collection := repo.conn.Collection(retainedPegoutQuoteCollection)
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

func (repo *pegoutMongoRepository) GetRetainedQuoteByState(ctx context.Context, states ...quote.PegoutState) ([]quote.RetainedPegoutQuote, error) {
	result := make([]quote.RetainedPegoutQuote, 0)
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	collection := repo.conn.Collection(retainedPegoutQuoteCollection)
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

func (repo *pegoutMongoRepository) DeleteQuotes(ctx context.Context, quotes []string) (uint, error) {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout*2)
	defer cancel()

	filter := bson.D{primitive.E{Key: "quote_hash", Value: bson.D{primitive.E{Key: "$in", Value: quotes}}}}
	pegoutResult, err := repo.conn.Collection(pegoutQuoteCollection).DeleteMany(dbCtx, filter)
	if err != nil {
		return 0, err
	}
	retainedResult, err := repo.conn.Collection(retainedPegoutQuoteCollection).DeleteMany(dbCtx, filter)
	if err != nil {
		return 0, err
	} else if pegoutResult.DeletedCount != retainedResult.DeletedCount {
		return 0, errors.New("pegout quote collections didn't match")
	}
	logDbInteraction(delete, fmt.Sprintf("removed %d records from %s collection", pegoutResult.DeletedCount, pegoutQuoteCollection))
	logDbInteraction(delete, fmt.Sprintf("removed %d records from %s collection", retainedResult.DeletedCount, retainedPegoutQuoteCollection))
	return uint(pegoutResult.DeletedCount + retainedResult.DeletedCount), nil
}

func (repo *pegoutMongoRepository) UpsertPegoutDeposit(ctx context.Context, deposit quote.PegoutDeposit) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	_, err := repo.conn.Collection(depositEventsCollection).ReplaceOne(
		dbCtx,
		bson.M{"tx_hash": deposit.TxHash},
		deposit,
		options.Replace().SetUpsert(true),
	)
	if err == nil {
		logDbInteraction(upsert, deposit)
	}
	return err
}

func (repo *pegoutMongoRepository) UpsertPegoutDeposits(ctx context.Context, deposits []quote.PegoutDeposit) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	if len(deposits) == 0 {
		return nil
	}

	documents := make([]mongo.WriteModel, 0)
	for _, deposit := range deposits {
		filter := bson.M{"tx_hash": deposit.TxHash}
		replaceModel := mongo.NewReplaceOneModel()
		replaceModel.SetFilter(filter)
		replaceModel.SetReplacement(deposit)
		replaceModel.SetUpsert(true)

		documents = append(documents, replaceModel)
	}

	_, err := repo.conn.Collection(depositEventsCollection).BulkWrite(
		dbCtx,
		documents,
	)
	if err == nil {
		logDbInteraction(upsert, deposits)
	}
	return err
}
