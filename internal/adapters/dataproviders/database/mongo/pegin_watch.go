package mongo

import (
	"context"
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDb "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	PegInWatchCollection = "peginWatch"
	peginWatchCursorId   = "scanCursor"
)

type peginWatchMongoRepository struct {
	conn *Connection
}

type peginWatchCursor struct {
	Id               string `bson:"_id"`
	LastScannedBlock uint64 `bson:"last_scanned_block"`
}

func NewPegInWatchMongoRepository(conn *Connection) rootstock.PegInWatchRepository {
	return &peginWatchMongoRepository{conn: conn}
}

func (repo *peginWatchMongoRepository) Upsert(
	ctx context.Context,
	watch rootstock.PegInWatch,
) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	filter := eventIdentity(watch.TxHash, watch.LogIndex)
	update := bson.M{"$setOnInsert": watch}
	_, err := repo.conn.Collection(PegInWatchCollection).UpdateOne(
		dbCtx,
		filter,
		update,
		options.UpdateOne().SetUpsert(true),
	)
	if mongoDb.IsDuplicateKeyError(err) {
		return nil
	}
	return err
}

func (repo *peginWatchMongoRepository) Get(
	ctx context.Context,
	txHash string,
	logIndex uint,
) (*rootstock.PegInWatch, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	var watch rootstock.PegInWatch
	err := repo.conn.Collection(PegInWatchCollection).
		FindOne(dbCtx, eventIdentity(txHash, logIndex)).
		Decode(&watch)
	if errors.Is(err, mongoDb.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &watch, nil
}

func (repo *peginWatchMongoRepository) List(
	ctx context.Context,
) ([]rootstock.PegInWatch, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	cursor, err := repo.conn.Collection(PegInWatchCollection).Find(
		dbCtx,
		bson.M{"tx_hash": bson.M{"$exists": true}},
		options.Find().SetSort(bson.D{{Key: "block_number", Value: 1}, {Key: "log_index", Value: 1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(dbCtx)

	watches := make([]rootstock.PegInWatch, 0)
	if err = cursor.All(dbCtx, &watches); err != nil {
		return nil, err
	}
	return watches, nil
}

func (repo *peginWatchMongoRepository) Update(
	ctx context.Context,
	watch rootstock.PegInWatch,
) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	result, err := repo.conn.Collection(PegInWatchCollection).ReplaceOne(
		dbCtx,
		eventIdentity(watch.TxHash, watch.LogIndex),
		watch,
	)
	if err != nil {
		return err
	}
	if result.MatchedCount != 1 {
		return errors.New("pegin watch not found")
	}
	return nil
}

func (repo *peginWatchMongoRepository) GetCursor(
	ctx context.Context,
) (uint64, bool, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	var cursor peginWatchCursor
	err := repo.conn.Collection(PegInWatchCollection).
		FindOne(dbCtx, bson.M{"_id": peginWatchCursorId}).
		Decode(&cursor)
	if errors.Is(err, mongoDb.ErrNoDocuments) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return cursor.LastScannedBlock, true, nil
}

func (repo *peginWatchMongoRepository) SetCursor(
	ctx context.Context,
	lastScannedBlock uint64,
) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	_, err := repo.conn.Collection(PegInWatchCollection).UpdateOne(
		dbCtx,
		bson.M{"_id": peginWatchCursorId},
		bson.M{"$set": bson.M{"last_scanned_block": lastScannedBlock}},
		options.UpdateOne().SetUpsert(true),
	)
	return err
}

func eventIdentity(txHash string, logIndex uint) bson.M {
	return bson.M{"tx_hash": txHash, "log_index": logIndex}
}
