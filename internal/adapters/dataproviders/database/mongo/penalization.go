package mongo

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	PenalizedEventCollection = "penalizedEvent"
)

type penalizedEventMongoRepository struct {
	conn *Connection
}

func NewPenalizedEventRepository(conn *Connection) penalization.PenalizedEventRepository {
	return &penalizedEventMongoRepository{conn: conn}
}

func (repo *penalizedEventMongoRepository) InsertPenalization(ctx context.Context, PenalizedEvent penalization.PenalizedEvent) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()
	collection := repo.conn.Collection(PenalizedEventCollection)
	_, err := collection.InsertOne(dbCtx, PenalizedEvent)
	if err != nil {
		return err
	} else {
		logDbInteraction(Insert, PenalizedEvent)
		return nil
	}
}

func (repo *penalizedEventMongoRepository) GetPenalizationsByQuoteHashes(ctx context.Context, quoteHashes []string) ([]penalization.PenalizedEvent, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	for _, hash := range quoteHashes {
		if err := quote.ValidateQuoteHash(hash); err != nil {
			return nil, err
		}
	}

	collection := repo.conn.Collection(PenalizedEventCollection)
	filter := bson.M{
		"quote_hash": bson.M{"$in": quoteHashes},
	}

	penalizations := make([]penalization.PenalizedEvent, 0)

	cursor, err := collection.Find(dbCtx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var result penalization.PenalizedEvent
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		penalizations = append(penalizations, result)
	}
	logDbInteraction(Read, penalizations)
	return penalizations, nil
}
