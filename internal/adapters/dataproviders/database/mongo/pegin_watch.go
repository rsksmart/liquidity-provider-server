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
	PegInWatchCollection                     = "peginWatch"
	peginAddressRegistryCheckpointDocumentID = "checkpoint"
)

type peginWatchMongoRepository struct {
	conn *Connection
}

type peginAddressRegistryCheckpointDocument struct {
	Id                 string    `bson:"_id"`
	LocalRoot          *[32]byte `bson:"local_root,omitempty"`
	LastProcessedBlock *uint64   `bson:"last_processed_block,omitempty"`
}

func NewPegInWatchMongoRepository(conn *Connection) rootstock.PegInWatchRepositorySet {
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

	var checkpointDocument peginAddressRegistryCheckpointDocument
	err := repo.conn.Collection(PegInWatchCollection).
		FindOne(dbCtx, bson.M{"_id": peginAddressRegistryCheckpointDocumentID}).
		Decode(&checkpointDocument)
	if errors.Is(err, mongoDb.ErrNoDocuments) {
		return rootstock.PegInWatchCheckpoint{}, false, nil
	}
	if err != nil {
		return rootstock.PegInWatchCheckpoint{}, false, err
	}
	if checkpointDocument.LocalRoot == nil ||
		checkpointDocument.LastProcessedBlock == nil {
		return rootstock.PegInWatchCheckpoint{}, false, nil
	}
	return rootstock.PegInWatchCheckpoint{
		LocalRoot:          *checkpointDocument.LocalRoot,
		LastProcessedBlock: *checkpointDocument.LastProcessedBlock,
	}, true, nil
}

func (repo *peginWatchMongoRepository) SetCheckpoint(
	ctx context.Context,
	checkpoint rootstock.PegInWatchCheckpoint,
) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	result, err := repo.conn.Collection(PegInWatchCollection).UpdateOne(
		dbCtx,
		bson.M{"_id": peginAddressRegistryCheckpointDocumentID},
		bson.M{"$set": bson.M{
			"local_root":           checkpoint.LocalRoot,
			"last_processed_block": checkpoint.LastProcessedBlock,
		}},
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		return err
	}
	if result == nil || result.MatchedCount+result.UpsertedCount != 1 {
		return errors.New("pegin address registry checkpoint was not persisted")
	}
	return nil
}

func (repo *peginWatchMongoRepository) DeleteCheckpoint(ctx context.Context) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	_, err := repo.conn.Collection(PegInWatchCollection).DeleteOne(
		dbCtx,
		bson.M{"_id": peginAddressRegistryCheckpointDocumentID},
	)
	return err
}

func rskAddressIdentity(rskAddress string) bson.M {
	return bson.M{"rsk_address": rskAddress}
}
