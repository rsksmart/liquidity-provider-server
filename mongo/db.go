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
	GetQuote(quoteHash string) (*types.Quote, error) // returns nil if not found
	RetainQuote(entry *types.RetainedQuote) error
	GetRetainedQuotes(filter []types.RQState) ([]*types.RetainedQuote, error)
	GetRetainedQuote(hash string) (*types.RetainedQuote, error) // returns nil if not found
	DeleteExpiredQuotes(expTimestamp int64) error
	UpdateRetainedQuoteState(hash string, oldState types.RQState, newState types.RQState) error
	GetLockedLiquidity() (*types.Wei, error)
	InsertPegOutQuote(id string, q *pegout.Quote, derivationAddress string) error
	GetPegOutQuote(quoteHash string) (*pegout.Quote, error)
	RetainPegOutQuote(entry *pegout.RetainedQuote) error
	GetRetainedPegOutQuote(hash string) (*pegout.RetainedQuote, error)
	UpdateRetainedPegOutQuoteState(hash string, oldState types.RQState, newState types.RQState) error
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

type PegoutQuote struct {
	Hash              string        `bson:"quotehash,omitempty"`
	DerivationAddress string        `bson:"derivationAddress,omitempty"`
	Quote             *pegout.Quote `bson:"quote,omitempty"`
}

type RetainedPeginQuote struct {
	QuoteHash   string        `json:"quoteHash" db:"quote_hash"`
	DepositAddr string        `json:"depositAddr" db:"deposit_addr"`
	Signature   string        `json:"signature" db:"signature"`
	ReqLiq      string        `json:"reqLiq" db:"req_liq"`
	State       types.RQState `json:"state" db:"state"`
}

func Connect(host string, user string, password string) (*DB, error) {
	log.Debug("Connecting to MongoDB")

	clientOptions := options.Client().
		ApplyURI("mongodb://" + user + ":" + password + "@" + host + ":27017/admin")
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

	quote := types.Quote{
		AgreementTimestamp: result.AgreementTimestamp,
		BTCRefundAddr:      result.BTCRefundAddr,
		CallFee:            types.NewWei(callFee),
		CallOnRegister:     result.CallOnRegister,
		CallTime:           result.CallTime,
		Confirmations:      result.Confirmations,
		ContractAddr:       result.ContractAddr,
		Data:               result.Data,
		FedBTCAddr:         result.FedBTCAddr,
		GasLimit:           result.GasLimit,
		LBCAddr:            result.LBCAddr,
		LPBTCAddr:          result.LPBTCAddr,
		LPRSKAddr:          result.LPRSKAddr,
		Nonce:              result.Nonce,
		PenaltyFee:         types.NewWei(penaltyFee),
		RSKRefundAddr:      result.RSKRefundAddr,
		TimeForDeposit:     result.TimeForDeposit,
		Value:              types.NewWei(value),
	}

	return &quote, nil
}

func (db *DB) RetainQuote(entry *types.RetainedQuote) error {
	log.Debug("inserting retained quote mongo DB:", entry.QuoteHash, "; DepositAddr: ", entry.DepositAddr, "; Signature: ", entry.Signature, "; ReqLiq: ", entry.ReqLiq)
	coll := db.db.Database("flyover").Collection("retainedPeginQuote")

	quoteToRetain := RetainedPeginQuote{
		DepositAddr: entry.DepositAddr,
		QuoteHash:   entry.QuoteHash,
		ReqLiq:      entry.ReqLiq.String(),
		Signature:   entry.Signature,
		State:       entry.State,
	}

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
	rows.All(context.TODO(), &retainedQuotes)

	defer rows.Close(context.TODO())
	for rows.Next(context.TODO()) {
		var rq RetainedPeginQuote
		if err = rows.Decode(&rq); err != nil {
			return nil, err
		}

		reqLiq, err := strconv.ParseInt(rq.ReqLiq, 10, 64)
		if err != nil {
			return nil, err
		}
		rqToReturn := types.RetainedQuote{
			DepositAddr: rq.DepositAddr,
			QuoteHash:   rq.QuoteHash,
			ReqLiq:      types.NewWei(reqLiq),
			Signature:   rq.Signature,
			State:       rq.State,
		}

		retainedQuotes = append(retainedQuotes, &rqToReturn)
	}

	log.Debug("Retained Quotes: ", retainedQuotes)

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

	reqLiq, err := strconv.ParseInt(result.ReqLiq, 10, 64)
	if err != nil {
		return nil, err
	}

	rqToReturn := types.RetainedQuote{
		DepositAddr: result.DepositAddr,
		QuoteHash:   result.QuoteHash,
		ReqLiq:      types.NewWei(reqLiq),
		Signature:   result.Signature,
		State:       result.State,
	}

	return &rqToReturn, nil

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

func (db *DB) InsertPegOutQuote(id string, q *pegout.Quote, derivationAddress string) error {
	log.Debug("inserting pegout_quote{", id, "}", ": ", q)
	coll := db.db.Database("flyover").Collection("pegoutQuote")

	quoteToInsert := &PegoutQuote{
		Hash:              id,
		DerivationAddress: derivationAddress,
		Quote:             q,
	}

	_, err := coll.InsertOne(context.TODO(), quoteToInsert)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetPegOutQuote(quoteHash string) (*pegout.Quote, error) {
	log.Debug("retrieving pegout quote: ", quoteHash)
	coll := db.db.Database("flyover").Collection("pegoutQuote")
	filter := bson.D{primitive.E{Key: "quotehash", Value: quoteHash}}
	var result PegoutQuote
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return result.Quote, nil
}

func (db *DB) RetainPegOutQuote(entry *pegout.RetainedQuote) error {
	log.Debug("inserting retained quote:", entry.QuoteHash, "; DepositAddr: ", entry.DepositAddr, "; Signature: ", entry.Signature, "; ReqLiq: ", entry.ReqLiq)

	coll := db.db.Database("flyover").Collection("retainedPegoutQuote")

	_, err := coll.InsertOne(context.TODO(), entry)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetRetainedPegOutQuote(hash string) (*pegout.RetainedQuote, error) {
	log.Debug("getting retained quote: ", hash)

	coll := db.db.Database("flyover").Collection("retainedPegoutQuote")
	filter := bson.D{primitive.E{Key: "quoteHash", Value: hash}}

	var result *pegout.RetainedQuote
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (db *DB) UpdateRetainedPegOutQuoteState(hash string, oldState types.RQState, newState types.RQState) error {
	log.Debugf("updating state from %v to %v for retained quote: %v", oldState, newState, hash)

	coll := db.db.Database("flyover").Collection("retainedPegoutQuote")
	filter := bson.D{primitive.E{Key: "quoteHash", Value: hash}, primitive.E{Key: "state", Value: oldState}}
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
