package mongo

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	BatchPegOutEventsCollection = "batchPegOutEvents"
)

type batchPegOutMongoRepository struct {
	conn *Connection
}

func NewBatchPegOutMongoRepository(conn *Connection) rootstock.BatchPegOutRepository {
	return &batchPegOutMongoRepository{conn: conn}
}

func (repo *batchPegOutMongoRepository) UpsertBatch(ctx context.Context, batch rootstock.BatchPegOut) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	result, err := repo.conn.Collection(BatchPegOutEventsCollection).ReplaceOne(
		dbCtx,
		bson.M{"transaction_hash": batch.TransactionHash},
		batch,
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		return err
	} else if result.ModifiedCount > 1 {
		return errors.New("multiple batch pegouts updated")
	}
	logDbInteraction(Upsert, batch)
	return err
}
