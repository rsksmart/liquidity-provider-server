package mongo

import (
	"context"
	"errors"
	"fmt"
	"regexp"
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
	PegoutQuoteCollection         = "pegoutQuote"
	RetainedPegoutQuoteCollection = "retainedPegoutQuote"
	DepositEventsCollection       = "depositEvents"
	PegoutCreationDataCollection  = "pegoutQuoteCreationData"
)

type StoredPegoutQuote struct {
	quote.PegoutQuote `bson:",inline"`
	Hash              string `json:"hash" bson:"hash"`
}

type StoredPegoutCreationData struct {
	quote.PegoutCreationData `bson:",inline"`
	Hash                     string `json:"hash" bson:"hash"`
}

type pegoutMongoRepository struct {
	conn *Connection
}

func NewPegoutMongoRepository(conn *Connection) quote.PegoutQuoteRepository {
	return &pegoutMongoRepository{conn: conn}
}

func (repo *pegoutMongoRepository) InsertQuote(ctx context.Context, createdQuote quote.CreatedPegoutQuote) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	collection := repo.conn.Collection(PegoutQuoteCollection)
	storedQuote := StoredPegoutQuote{
		PegoutQuote: createdQuote.Quote,
		Hash:        createdQuote.Hash,
	}
	_, err := collection.InsertOne(dbCtx, storedQuote)
	if err != nil {
		return err
	}
	logDbInteraction(Insert, storedQuote)

	collection = repo.conn.Collection(PegoutCreationDataCollection)
	storedCreationData := StoredPegoutCreationData{
		PegoutCreationData: createdQuote.CreationData,
		Hash:               createdQuote.Hash,
	}
	_, err = collection.InsertOne(dbCtx, storedCreationData)
	if err != nil {
		return err
	}
	logDbInteraction(Insert, storedQuote)
	return nil
}

func (repo *pegoutMongoRepository) GetPegoutCreationData(ctx context.Context, hash string) quote.PegoutCreationData {
	var result StoredPegoutCreationData
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	if err := quote.ValidateQuoteHash(hash); err != nil {
		log.Error("Invalid hash. Returning default pegout creation data")
		return quote.PegoutCreationDataZeroValue()
	}

	collection := repo.conn.Collection(PegoutCreationDataCollection)
	filter := bson.D{primitive.E{Key: "hash", Value: hash}}

	err := collection.FindOne(dbCtx, filter).Decode(&result)
	if err != nil {
		log.Error("Hash not found. Returning default pegout creation data")
		return quote.PegoutCreationDataZeroValue()
	}
	logDbInteraction(Read, result.PegoutCreationData)
	return result.PegoutCreationData
}

func (repo *pegoutMongoRepository) GetQuote(ctx context.Context, hash string) (*quote.PegoutQuote, error) {
	var result StoredPegoutQuote
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	if err := quote.ValidateQuoteHash(hash); err != nil {
		return nil, err
	}

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

func (repo *pegoutMongoRepository) GetQuotesByHashesAndDate(
	ctx context.Context,
	hashes []string,
	startDate, endDate time.Time,
) ([]quote.PegoutQuote, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	for _, hash := range hashes {
		if err := quote.ValidateQuoteHash(hash); err != nil {
			return nil, err
		}
	}

	collection := repo.conn.Collection(PegoutQuoteCollection)
	quotesReturn := make([]quote.PegoutQuote, 0)
	filter := bson.M{
		"hash": bson.M{"$in": hashes},
		"agreement_timestamp": bson.M{
			"$gte": startDate.Unix(),
			"$lte": endDate.Unix(),
		},
	}

	cursor, err := collection.Find(dbCtx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var result StoredPegoutQuote
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		quotesReturn = append(quotesReturn, result.PegoutQuote)
	}
	logDbInteraction(Read, quotesReturn)
	return quotesReturn, nil
}

func (repo *pegoutMongoRepository) GetRetainedQuote(ctx context.Context, hash string) (*quote.RetainedPegoutQuote, error) {
	var result quote.RetainedPegoutQuote
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	if err := quote.ValidateQuoteHash(hash); err != nil {
		return nil, err
	}

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
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
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
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	sanitizedAddress := regexp.QuoteMeta(address)
	filter := bson.M{"from": bson.M{"$regex": sanitizedAddress, "$options": "i"}}
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
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
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
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
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
			updateResult, updateErr := collection.UpdateOne(dbCtx, filter, updateStatement)
			if updateErr != nil {
				return int64(0), updateErr
			}
			count += updateResult.ModifiedCount
		}
		return count, nil
	})
	if err != nil {
		return err
	}
	parsedResult, ok := result.(int64)
	if !ok {
		return errors.New("unexpected result type")
	} else if parsedResult != int64(len(retainedQuotes)) {
		return fmt.Errorf("mismatch on updated documents. Expected %d, updated %d", len(retainedQuotes), result)
	}
	logDbInteraction(Update, retainedQuotes)
	return nil
}

func (repo *pegoutMongoRepository) GetRetainedQuoteByState(ctx context.Context, states ...quote.PegoutState) ([]quote.RetainedPegoutQuote, error) {
	result := make([]quote.RetainedPegoutQuote, 0)
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
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
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout*2)
	defer cancel()

	quoteFilter := bson.D{primitive.E{Key: "hash", Value: bson.D{primitive.E{Key: "$in", Value: quotes}}}}
	retainedFilter := bson.D{primitive.E{Key: "quote_hash", Value: bson.D{primitive.E{Key: "$in", Value: quotes}}}}
	pegoutResult, err := repo.conn.Collection(PegoutQuoteCollection).DeleteMany(dbCtx, quoteFilter)
	if err != nil {
		return 0, err
	}
	retainedResult, err := repo.conn.Collection(RetainedPegoutQuoteCollection).DeleteMany(dbCtx, retainedFilter)
	if err != nil {
		return 0, err
	}
	creationDataResult, err := repo.conn.Collection(PegoutCreationDataCollection).DeleteMany(dbCtx, quoteFilter)
	if err != nil {
		return 0, err
	}
	const msgTemplate = "removed %d records from %s collection"
	logDbInteraction(Delete, fmt.Sprintf(msgTemplate, pegoutResult.DeletedCount, PegoutQuoteCollection))
	logDbInteraction(Delete, fmt.Sprintf(msgTemplate, retainedResult.DeletedCount, RetainedPegoutQuoteCollection))
	logDbInteraction(Delete, fmt.Sprintf(msgTemplate, creationDataResult.DeletedCount, PegoutCreationDataCollection))
	// creation data doesn't count for mismatch because not all the quotes have it
	if pegoutResult.DeletedCount != retainedResult.DeletedCount {
		return 0, errors.New("pegout quote collections didn't match")
	}
	return uint(pegoutResult.DeletedCount + retainedResult.DeletedCount + creationDataResult.DeletedCount), nil
}

func (repo *pegoutMongoRepository) UpsertPegoutDeposit(ctx context.Context, deposit quote.PegoutDeposit) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
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
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
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

func (repo *pegoutMongoRepository) ListQuotesByDateRange(ctx context.Context, startDate, endDate time.Time) ([]quote.PegoutQuoteWithRetained, error) {
	result := make([]quote.PegoutQuoteWithRetained, 0)
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	quoteFilter := bson.D{{Key: "agreement_timestamp", Value: bson.D{
		{Key: "$gte", Value: startDate.Unix()},
		{Key: "$lte", Value: endDate.Unix()},
	}}}
	findOpts := options.Find().SetSort(bson.D{{Key: "agreement_timestamp", Value: 1}})
	quoteCursor, err := repo.conn.Collection(PegoutQuoteCollection).Find(dbCtx, quoteFilter, findOpts)
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
	var storedQuotes []StoredPegoutQuote
	if err = quoteCursor.All(dbCtx, &storedQuotes); err != nil {
		return nil, err
	}
	if len(storedQuotes) == 0 {
		logDbInteraction(Read, result)
		return result, nil
	}
	hashToIndex := make(map[string]int, len(storedQuotes))
	quoteHashes := make([]string, len(storedQuotes))
	result = make([]quote.PegoutQuoteWithRetained, len(storedQuotes))
	for i, stored := range storedQuotes {
		quoteHashes[i] = stored.Hash
		hashToIndex[stored.Hash] = i
		result[i] = quote.PegoutQuoteWithRetained{
			Quote:         stored.PegoutQuote,
			RetainedQuote: quote.RetainedPegoutQuote{},
		}
	}
	retainedCursor, err := repo.conn.Collection(RetainedPegoutQuoteCollection).Find(
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
	var retainedQuote quote.RetainedPegoutQuote
	for retainedCursor.Next(dbCtx) {
		if err := retainedCursor.Decode(&retainedQuote); err != nil {
			return result, err
		}
		if idx, exists := hashToIndex[retainedQuote.QuoteHash]; exists {
			result[idx].RetainedQuote = retainedQuote
		}
	}
	if err := retainedCursor.Err(); err != nil {
		return result, err
	}
	logDbInteraction(Read, len(result))
	return result, nil
}
