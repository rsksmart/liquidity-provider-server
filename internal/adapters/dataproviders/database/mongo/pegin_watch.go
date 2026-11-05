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
	PegInWatchCollection           = "peginWatch"
	PegInWatchCheckpointCollection = "peginWatchCheckpoint"
	peginWatchCursorId             = "scanCursor"
	peginWatchCheckpointId         = "verified"
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

	filter := rskAddressIdentity(watch.RskAddress)
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
	rskAddress string,
) (*rootstock.PegInWatch, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	var watch rootstock.PegInWatch
	err := repo.conn.Collection(PegInWatchCollection).
		FindOne(dbCtx, rskAddressIdentity(rskAddress)).
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
		bson.M{"rsk_address": bson.M{"$exists": true}},
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
		rskAddressIdentity(watch.RskAddress),
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

func (repo *peginWatchMongoRepository) GetCheckpoint(
	ctx context.Context,
) (rootstock.PegInWatchCheckpoint, bool, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	var checkpoint rootstock.PegInWatchCheckpoint
	err := repo.conn.Collection(PegInWatchCheckpointCollection).
		FindOne(dbCtx, bson.M{"_id": peginWatchCheckpointId}).
		Decode(&checkpoint)
	if errors.Is(err, mongoDb.ErrNoDocuments) {
		return rootstock.PegInWatchCheckpoint{}, false, nil
	}
	if err != nil {
		return rootstock.PegInWatchCheckpoint{}, false, err
	}
	return checkpoint, true, nil
}

func (repo *peginWatchMongoRepository) SaveVerifiedRegistrations(
	ctx context.Context,
	watches []rootstock.PegInWatch,
	checkpoint rootstock.PegInWatchCheckpoint,
) error {
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

	_, err = session.WithTransaction(dbCtx, func(sessionContext context.Context) (any, error) {
		watchCollection := repo.conn.Collection(PegInWatchCollection)
		for _, watch := range watches {
			if _, upsertErr := watchCollection.UpdateOne(
				sessionContext,
				rskAddressIdentity(watch.RskAddress),
				bson.M{"$setOnInsert": watch},
				options.UpdateOne().SetUpsert(true),
			); upsertErr != nil {
				return nil, upsertErr
			}
		}
		_, checkpointErr := repo.conn.Collection(PegInWatchCheckpointCollection).UpdateOne(
			sessionContext,
			bson.M{"_id": peginWatchCheckpointId},
			bson.M{"$set": checkpoint},
			options.UpdateOne().SetUpsert(true),
		)
		return nil, checkpointErr
	})
	return err
}

func rskAddressIdentity(rskAddress string) bson.M {
	return bson.M{"rsk_address": rskAddress}
}
