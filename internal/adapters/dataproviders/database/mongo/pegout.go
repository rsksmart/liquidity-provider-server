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
	PegoutQuoteCollection         = "pegoutQuote"
	RetainedPegoutQuoteCollection = "retainedPegoutQuote"
	DepositEventsCollection       = "depositEvents"
)

type StoredPegoutQuote struct {
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
	collection := repo.conn.Collection(PegoutQuoteCollection)
	storedQuote := StoredPegoutQuote{
		PegoutQuote: pegoutQuote,
		Hash:        hash,
	}
	_, err := collection.InsertOne(dbCtx, storedQuote)
	if err != nil {
		return err
	} else {
		logDbInteraction(Insert, storedQuote)
		return nil
	}
}

func (repo *pegoutMongoRepository) GetQuote(ctx context.Context, hash string) (*quote.PegoutQuote, error) {
	var result StoredPegoutQuote
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	collection := repo.conn.Collection(PegoutQuoteCollection)
	filter := bson.D{primitive.E{Key: "hash", Value: hash}}

	err := collection.FindOne(dbCtx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	logDbInteraction(Read, result.PegoutQuote)
	return &result.PegoutQuote, nil
}

func (repo *pegoutMongoRepository) GetRetainedQuote(ctx context.Context, hash string) (*quote.RetainedPegoutQuote, error) {
	var result quote.RetainedPegoutQuote
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	collection := repo.conn.Collection(RetainedPegoutQuoteCollection)
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

func (repo *pegoutMongoRepository) InsertRetainedQuote(ctx context.Context, retainedQuote quote.RetainedPegoutQuote) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	collection := repo.conn.Collection(RetainedPegoutQuoteCollection)
	_, err := collection.InsertOne(dbCtx, retainedQuote)
	if err != nil {
		return err
	} else {
		logDbInteraction(Insert, retainedQuote)
		return nil
	}
}

func (repo *pegoutMongoRepository) ListPegoutDepositsByAddress(ctx context.Context, address string) ([]quote.PegoutDeposit, error) {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	filter := bson.M{"from": bson.M{"$regex": address, "$options": "i"}}
	sort := options.Find().SetSort(bson.M{"timestamp": -1})
	cursor, err := repo.conn.Collection(DepositEventsCollection).Find(dbCtx, filter, sort)
	if err != nil {
		return make([]quote.PegoutDeposit, 0), err
	}

	var documents []quote.PegoutDeposit
	if err = cursor.All(ctx, &documents); err != nil {
		return make([]quote.PegoutDeposit, 0), err
	}
	logDbInteraction(Read, fmt.Sprintf("%d pegout deposits", len(documents)))
	return documents, nil
}

func (repo *pegoutMongoRepository) UpdateRetainedQuote(ctx context.Context, retainedQuote quote.RetainedPegoutQuote) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	collection := repo.conn.Collection(RetainedPegoutQuoteCollection)
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

func (repo *pegoutMongoRepository) UpdateRetainedQuotes(ctx context.Context, retainedQuotes []quote.RetainedPegoutQuote) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	session, err := repo.conn.client.StartSession()
	defer func() {
		if session != nil {
			session.EndSession(dbCtx)
		}
	}()
	if err != nil {
		return err
	}
	collection := repo.conn.Collection(RetainedPegoutQuoteCollection)
	result, err := session.WithTransaction(dbCtx, func(sessionContext mongo.SessionContext) (any, error) {
		var count int64 = 0
		for _, retainedQuote := range retainedQuotes {
			filter := bson.D{primitive.E{Key: "quote_hash", Value: retainedQuote.QuoteHash}}
			updateStatement := bson.D{primitive.E{Key: "$set", Value: retainedQuote}}
			result, updateErr := collection.UpdateOne(dbCtx, filter, updateStatement)
			if updateErr != nil {
				return int64(0), updateErr
			}
			count += result.ModifiedCount
		}
		return count, nil
	})
	if err != nil {
		return err
	} else if result.(int64) != int64(len(retainedQuotes)) {
		return fmt.Errorf("mismatch on updated documents. Expected %d, updated %d", len(retainedQuotes), result)
	}
	logDbInteraction(Update, retainedQuotes)
	return nil
}

func (repo *pegoutMongoRepository) GetRetainedQuoteByState(ctx context.Context, states ...quote.PegoutState) ([]quote.RetainedPegoutQuote, error) {
	result := make([]quote.RetainedPegoutQuote, 0)
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	collection := repo.conn.Collection(RetainedPegoutQuoteCollection)
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

func (repo *pegoutMongoRepository) DeleteQuotes(ctx context.Context, quotes []string) (uint, error) {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout*2)
	defer cancel()

	filter := bson.D{primitive.E{Key: "quote_hash", Value: bson.D{primitive.E{Key: "$in", Value: quotes}}}}
	pegoutResult, err := repo.conn.Collection(PegoutQuoteCollection).DeleteMany(dbCtx, filter)
	if err != nil {
		return 0, err
	}
	retainedResult, err := repo.conn.Collection(RetainedPegoutQuoteCollection).DeleteMany(dbCtx, filter)
	if err != nil {
		return 0, err
	}
	logDbInteraction(Delete, fmt.Sprintf("removed %d records from %s collection", pegoutResult.DeletedCount, PegoutQuoteCollection))
	logDbInteraction(Delete, fmt.Sprintf("removed %d records from %s collection", retainedResult.DeletedCount, RetainedPegoutQuoteCollection))
	if pegoutResult.DeletedCount != retainedResult.DeletedCount {
		return 0, errors.New("pegout quote collections didn't match")
	}
	return uint(pegoutResult.DeletedCount + retainedResult.DeletedCount), nil
}

func (repo *pegoutMongoRepository) UpsertPegoutDeposit(ctx context.Context, deposit quote.PegoutDeposit) error {
	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	result, err := repo.conn.Collection(DepositEventsCollection).ReplaceOne(
		dbCtx,
		bson.M{"tx_hash": deposit.TxHash},
		deposit,
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		return err
	} else if result.ModifiedCount > 1 {
		return errors.New("multiple deposits updated")
	}
	logDbInteraction(Upsert, deposit)
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

	_, err := repo.conn.Collection(DepositEventsCollection).BulkWrite(
		dbCtx,
		documents,
	)
	if err == nil {
		logDbInteraction(Upsert, deposits)
	}
	return err
}
