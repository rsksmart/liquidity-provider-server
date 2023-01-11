package mongoDB

import (
	"context"
	"fmt"
	"strconv"
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
	Hash       string `bson:"hash,omitempty"`
	Expiration uint32 `bson:"expiration,omitempty"`

	FedBTCAddr         string `bson:"fedBTCAddr,omitempty"`
	LBCAddr            string `bson:"lbcAddr,omitempty"`
	LPRSKAddr          string `bson:"lpRSKAddr,omitempty"`
	BTCRefundAddr      string `bson:"btcRefundAddr,omitempty"`
	RSKRefundAddr      string `bson:"rskRefundAddr,omitempty"`
	LPBTCAddr          string `bson:"lpBTCAddr,omitempty"`
	CallFee            string `bson:"callFee,omitempty"`
	PenaltyFee         string `bson:"penaltyFee,omitempty"`
	ContractAddr       string `bson:"contractAddr,omitempty"`
	Data               string `bson:"data,omitempty"`
	GasLimit           uint32 `bson:"gasLimit,omitempty"`
	Nonce              int64  `bson:"nonce,omitempty"`
	Value              string `bson:"value,omitempty"`
	AgreementTimestamp uint32 `bson:"agreementTimestamp,omitempty"`
	TimeForDeposit     uint32 `bson:"timeForDeposit,omitempty"`
	CallTime           uint32 `bson:"callTime,omitempty"`
	Confirmations      uint16 `bson:"confirmations,omitempty"`
	CallOnRegister     bool   `bson:"callOnRegister,omitempty"`
}

type RetainedPeginQuote struct {
	QuoteHash   string        `bson:"quoteHash,omitempty"`
	DepositAddr string        `bson:"depositAddr,omitempty"`
	Signature   string        `bson:"signature,omitempty"`
	ReqLiq      string        `bson:"reqLiq,omitempty"`
	State       types.RQState `bson:"state,omitempty"`
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
		Hash:               id,
		Expiration:         q.AgreementTimestamp + q.TimeForDeposit,
		FedBTCAddr:         q.FedBTCAddr,
		LBCAddr:            q.LBCAddr,
		LPRSKAddr:          q.LPRSKAddr,
		BTCRefundAddr:      q.BTCRefundAddr,
		RSKRefundAddr:      q.RSKRefundAddr,
		LPBTCAddr:          q.LPBTCAddr,
		CallFee:            q.CallFee.String(),
		PenaltyFee:         q.PenaltyFee.String(),
		ContractAddr:       q.ContractAddr,
		Data:               q.Data,
		GasLimit:           q.GasLimit,
		Nonce:              q.Nonce,
		Value:              q.Value.String(),
		AgreementTimestamp: q.AgreementTimestamp,
		TimeForDeposit:     q.TimeForDeposit,
		CallTime:           q.CallTime,
		Confirmations:      q.Confirmations,
		CallOnRegister:     q.CallOnRegister,
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

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	var quote *types.Quote
	callFee, err := strconv.ParseInt(result.CallFee, 10, 64)
	if err != nil {
		return nil, err
	}
	penaltyFee, err := strconv.ParseInt(result.PenaltyFee, 10, 64)
	if err != nil {
		return nil, err
	}
	value, err := strconv.ParseInt(result.Value, 10, 64)
	if err != nil {
		return nil, err
	}

	quote.AgreementTimestamp = result.AgreementTimestamp
	quote.BTCRefundAddr = result.BTCRefundAddr
	quote.CallFee = types.NewWei(callFee)
	quote.CallOnRegister = result.CallOnRegister
	quote.CallTime = result.CallTime
	quote.Confirmations = result.Confirmations
	quote.ContractAddr = result.ContractAddr
	quote.Data = result.Data
	quote.FedBTCAddr = result.FedBTCAddr
	quote.GasLimit = result.GasLimit
	quote.LBCAddr = result.LBCAddr
	quote.LPBTCAddr = result.LPBTCAddr
	quote.LPRSKAddr = result.LPRSKAddr
	quote.Nonce = result.Nonce
	quote.PenaltyFee = types.NewWei(penaltyFee)
	quote.RSKRefundAddr = result.RSKRefundAddr
	quote.TimeForDeposit = result.TimeForDeposit
	quote.Value = types.NewWei(value)

	return quote, nil
}

func (db *DB) RetainQuote(entry *types.RetainedQuote) error {
	log.Debug("inserting retained quote mongo DB:", entry.QuoteHash, "; DepositAddr: ", entry.DepositAddr, "; Signature: ", entry.Signature, "; ReqLiq: ", entry.ReqLiq)
	coll := db.db.Database("flyover").Collection("retainedPeginQuote")

	var quoteToRetain RetainedPeginQuote
	quoteToRetain.DepositAddr = entry.DepositAddr
	quoteToRetain.QuoteHash = entry.QuoteHash
	quoteToRetain.ReqLiq = entry.ReqLiq.String()
	quoteToRetain.Signature = entry.Signature
	quoteToRetain.State = entry.State

	_, err := coll.InsertOne(context.TODO(), quoteToRetain)

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

	defer rows.Close(context.TODO())
	for rows.Next(context.TODO()) {
		var rq RetainedPeginQuote
		var rqToReturn *types.RetainedQuote
		if err = rows.Decode(&rq); err != nil {
			return nil, err
		}

		reqLiq, err := strconv.ParseInt(rq.ReqLiq, 10, 64)
		if err != nil {
			return nil, err
		}

		rqToReturn.DepositAddr = rq.DepositAddr
		rqToReturn.QuoteHash = rq.QuoteHash
		rqToReturn.ReqLiq = types.NewWei(reqLiq)
		rqToReturn.Signature = rq.Signature
		rqToReturn.State = rq.State

		retainedQuotes = append(retainedQuotes, rqToReturn)
	}

	if err != nil {
		return nil, err
	}

	return retainedQuotes, nil
}

func (db *DB) GetRetainedQuote(hash string) (*types.RetainedQuote, error) {
	log.Debug("getting retained quote mongo: ", hash)
	coll := db.db.Database("flyover").Collection("retainedPeginQuote")
	filter := bson.D{primitive.E{Key: "quotehash", Value: hash}}

	var result RetainedPeginQuote
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	var rqToReturn *types.RetainedQuote

	reqLiq, err := strconv.ParseInt(result.ReqLiq, 10, 64)
	if err != nil {
		return nil, err
	}

	rqToReturn.DepositAddr = result.DepositAddr
	rqToReturn.QuoteHash = result.QuoteHash
	rqToReturn.ReqLiq = types.NewWei(reqLiq)
	rqToReturn.Signature = result.Signature
	rqToReturn.State = result.State

	return rqToReturn, nil

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

	var lockedLiq = types.NewWei(0)
	for rows.Next(context.TODO()) {
		var reqLiqString string
		err = rows.Decode(&reqLiqString)
		if err != nil {
			return nil, err
		}
		reqLiqInt, err := strconv.ParseInt(reqLiqString, 10, 64)

		if err != nil {
			return nil, err
		}

		reqLiq := types.NewWei(reqLiqInt)

		lockedLiq.Add(lockedLiq, reqLiq)
	}

	log.Debug("Loked Liquidity: ", lockedLiq.String())

	return lockedLiq, nil
}
