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
	peginWatchCheckpointId         = "verified"
)

type peginWatchMongoRepository struct {
	conn *Connection
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

func (repo *peginWatchMongoRepository) GetRootAt(ctx context.Context, block uint64) ([32]byte, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	var watch rootstock.PegInWatch
	err := repo.conn.Collection(PegInWatchCollection).
		FindOne(
			dbCtx,
			bson.M{"block_number": bson.M{"$lte": block}},
			options.FindOne().SetSort(bson.D{{Key: "block_number", Value: -1}, {Key: "log_index", Value: -1}}),
		).
		Decode(&watch)
	if errors.Is(err, mongoDb.ErrNoDocuments) {
		return [32]byte{}, nil
	}
	if err != nil {
		return [32]byte{}, err
	}
	return watch.RegistrationRoot, nil
}

func (repo *peginWatchMongoRepository) ReplaceVerifiedRegistrations(
	ctx context.Context,
	fromBlock uint64,
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
		if _, bulkErr := repo.conn.Collection(PegInWatchCollection).BulkWrite(
			sessionContext,
			replaceRegistrationModels(fromBlock, watches),
			options.BulkWrite().SetOrdered(true),
		); bulkErr != nil {
			return nil, bulkErr
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

func replaceRegistrationModels(fromBlock uint64, watches []rootstock.PegInWatch) []mongoDb.WriteModel {
	addresses := make([]string, 0, len(watches))
	models := make([]mongoDb.WriteModel, 1, len(watches)+1)
	for _, watch := range watches {
		addresses = append(addresses, watch.RskAddress)
		models = append(models, mongoDb.NewUpdateOneModel().
			SetFilter(rskAddressIdentity(watch.RskAddress)).
			SetUpdate(bson.M{
				"$set": bson.M{
					"tx_hash":           watch.TxHash,
					"log_index":         watch.LogIndex,
					"block_number":      watch.BlockNumber,
					"registrant":        watch.Registrant,
					"registration_root": watch.RegistrationRoot,
					"updated_at":        watch.UpdatedAt,
				},
				"$setOnInsert": bson.M{
					"encoding":     watch.Encoding,
					"btc_address":  watch.BtcAddress,
					"state":        watch.State,
					"last_seen_at": watch.LastSeenAt,
					"last_error":   watch.LastError,
					"created_at":   watch.CreatedAt,
				},
			}).
			SetUpsert(true))
	}
	models[0] = mongoDb.NewDeleteManyModel().SetFilter(bson.M{
		"block_number": bson.M{"$gte": fromBlock},
		"rsk_address":  bson.M{"$nin": addresses},
	})
	return models
}

func rskAddressIdentity(rskAddress string) bson.M {
	return bson.M{"rsk_address": rskAddress}
}
