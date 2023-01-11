package mongoDB

import (
	"context"
	"fmt"
	"time"

	"github.com/rsksmart/liquidity-provider-server/pegout"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConnector interface {
	CheckConnection() error
	Close() error
	InsertQuote(id string, q *types.Quote) error
	InsertPegOutQuote(id string, q *pegout.Quote, derivationAddress string) error
	GetQuote(quoteHash string) (*types.Quote, error) // returns nil if not found
	GetPegOutQuote(quoteHash string) (*pegout.Quote, error)
	DeleteExpiredQuotes(expTimestamp int64) error
	RetainQuote(entry *types.RetainedQuote) error
	RetainPegOutQuote(entry *pegout.RetainedQuote) error
	GetRetainedQuotes(filter []types.RQState) ([]*types.RetainedQuote, error)
	GetRetainedQuote(hash string) (*types.RetainedQuote, error) // returns nil if not found
	GetRetainedPegOutQuote(hash string) (*pegout.RetainedQuote, error)
	UpdateRetainedQuoteState(hash string, oldState types.RQState, newState types.RQState) error
	UpdateRetainedPegOutQuoteState(hash string, oldState types.RQState, newState types.RQState) error
	GetLockedLiquidity() (*types.Wei, error)
}

type DB struct {
	db *mongo.Client
}

type PeginQuote struct {
	Hash       string       `bson:"hash,omitempty"`
	Expiration uint32       `bson:"expiration,omitempty"`
	Quote      *types.Quote `bson:"quote,omitempty"`
}

type QuoteHash struct {
	QuoteHash string `db:"quote_hash"`
}

type UpdateQuoteState struct {
	QuoteHash string        `db:"quote_hash"`
	OldState  types.RQState `db:"old_state"`
	NewState  types.RQState `db:"new_state"`
}

func Connect() (*DB, error) {
	log.Debug("Connecting to MongoDB")
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("**REMOVED**/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	log.Debug("Connected to MongoDB... ")

	return &DB{client}, nil
}

func (db *DB) CheckConnection() error {
	return db.db.Ping(context.TODO(), nil)
}

func (db *DB) Close() error {
	log.Debug("closing connection to mongoDB DB")
	err := db.db.Disconnect(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) InsertQuote(id string, q *types.Quote) error {
	log.Debug("inserting quote{", id, "}", ": ", q)
	coll := db.db.Database("flyover").Collection("peginQuote")

	quoteToinsert := &PeginQuote{
		id,
		q.AgreementTimestamp + q.TimeForDeposit,
		q,
	}

	_, err := coll.InsertOne(context.TODO(), quoteToinsert)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetQuote(quoteHash string) (*types.Quote, error) {
	log.Debug("retrieving quote: ", quoteHash)

	coll := db.db.Database("flyover").Collection("peginQuote")
	filter := bson.D{primitive.E{Key: "hash", Value: quoteHash}}
	var result PeginQuote
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	quote := result.Quote
	switch err {
	case nil:
		return quote, nil
	case mongo.ErrNoDocuments:
		return nil, nil
	default:
		return nil, err
	}
}

func (db *DB) RetainQuote(entry *types.RetainedQuote) error {
	log.Debug("inserting retained quote mongo DB:", entry.QuoteHash, "; DepositAddr: ", entry.DepositAddr, "; Signature: ", entry.Signature, "; ReqLiq: ", entry.ReqLiq)
	coll := db.db.Database("flyover").Collection("retainedPeginQuote")

	_, err := coll.InsertOne(context.TODO(), entry)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetRetainedQuotes(filter []types.RQState) ([]*types.RetainedQuote, error) {
	log.Debug("retrieving retained quotes MongoDB")
	coll := db.db.Database("flyover").Collection("retainedPeginQuote")
	query := bson.D{primitive.E{Key: "state", Value: bson.D{primitive.E{Key: "$in", Value: filter}}}}
	rows, err := coll.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	var retainedQuotes []*types.RetainedQuote
	err = rows.All(context.TODO(), &retainedQuotes)

	if err != nil {
		return nil, err
	}

	return retainedQuotes, nil
}

func (db *DB) GetRetainedQuote(hash string) (*types.RetainedQuote, error) {
	log.Debug("getting retained quote mongo: ", hash)
	coll := db.db.Database("flyover").Collection("retainedPeginQuote")
	filter := bson.D{primitive.E{Key: "quotehash", Value: hash}}

	var result *types.RetainedQuote
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	switch err {
	case nil:
		return result, nil
	case mongo.ErrNoDocuments:
		return nil, nil
	default:
		return nil, err
	}
}

func (db *DB) DeleteExpiredQuotes(expTimestamp int64) error {
	log.Debug("deleting expired quotes...")
	coll := db.db.Database("flyover").Collection("peginQuote")
	filter := bson.D{primitive.E{Key: "expiration", Value: bson.D{primitive.E{Key: "$lt", Value: expTimestamp}}}}
	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount > 0 {
		log.Infof("deleted %v expired quote(s)", result.DeletedCount)
	} else {
		log.Debug("no expired quotes found; nothing to delete")
	}

	return nil
}

func (db *DB) UpdateRetainedQuoteState(hash string, oldState types.RQState, newState types.RQState) error {
	log.Debugf("updating state from %v to %v for retained quote mongo: %v", oldState, newState, hash)

	coll := db.db.Database("flyover").Collection("retainedPeginQuote")
	filter := bson.D{primitive.E{Key: "quotehash", Value: hash}, primitive.E{Key: "state", Value: oldState}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "state", Value: newState}}}}
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount != 1 {
		return fmt.Errorf("error updating retained quote mongoBD: %v; oldState: %v; newState: %v", hash, oldState, newState)
	}

	return nil
}

func (db *DB) GetLockedLiquidity() (*types.Wei, error) {
	log.Debug("retrieving locked liquidity")

	coll := db.db.Database("flyover").Collection("retainedPeginQuote")
	stateFilter := []types.RQState{types.RQStateWaitingForDeposit, types.RQStateCallForUserFailed}
	filter := bson.D{primitive.E{Key: "state", Value: bson.D{primitive.E{Key: "$in", Value: stateFilter}}}}
	rows, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	lockedLiq := types.NewWei(0)
	for rows.Next(context.TODO()) {
		reqLiq := new(types.Wei)

		err = rows.Decode(&reqLiq)
		if err != nil {
			return nil, err
		}

		lockedLiq.Add(lockedLiq, reqLiq)
	}

	log.Debug("Loked Liquidity: ", lockedLiq.String())

	return lockedLiq, nil
}
