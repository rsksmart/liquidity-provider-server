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
	PegInAddressRegistryWatchCollection      = "peginAddressRegistryWatch"
	peginAddressRegistryCheckpointDocumentID = "checkpoint"
)

type peginAddressRegistryWatchMongoRepository struct {
	conn *Connection
}

type peginAddressRegistryCheckpointDocument struct {
	Id                 string    `bson:"_id"`
	LocalRoot          *[32]byte `bson:"local_root,omitempty"`
	LastProcessedBlock *uint64   `bson:"last_processed_block,omitempty"`
}

func NewPegInAddressRegistryWatchMongoRepository(conn *Connection) rootstock.PegInAddressRegistryWatchRepositorySet {
	return &peginAddressRegistryWatchMongoRepository{conn: conn}
}

func (repo *peginAddressRegistryWatchMongoRepository) Upsert(
	ctx context.Context,
	entry rootstock.PegInAddressRegistryWatchEntry,
) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	filter := rskAddressIdentity(entry.RskAddress)
	update := bson.M{"$setOnInsert": entry}
	_, err := repo.conn.Collection(PegInAddressRegistryWatchCollection).UpdateOne(
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

func (repo *peginAddressRegistryWatchMongoRepository) Get(
	ctx context.Context,
	rskAddress string,
) (*rootstock.PegInAddressRegistryWatchEntry, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	var entry rootstock.PegInAddressRegistryWatchEntry
	err := repo.conn.Collection(PegInAddressRegistryWatchCollection).
		FindOne(dbCtx, rskAddressIdentity(rskAddress)).
		Decode(&entry)
	if errors.Is(err, mongoDb.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (repo *peginAddressRegistryWatchMongoRepository) List(
	ctx context.Context,
) ([]rootstock.PegInAddressRegistryWatchEntry, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	cursor, err := repo.conn.Collection(PegInAddressRegistryWatchCollection).Find(
		dbCtx,
		bson.M{"rsk_address": bson.M{"$exists": true}},
		options.Find().SetSort(bson.D{{Key: "block_number", Value: 1}, {Key: "log_index", Value: 1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(dbCtx)

	entries := make([]rootstock.PegInAddressRegistryWatchEntry, 0)
	if err = cursor.All(dbCtx, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (repo *peginAddressRegistryWatchMongoRepository) Update(
	ctx context.Context,
	entry rootstock.PegInAddressRegistryWatchEntry,
) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	result, err := repo.conn.Collection(PegInAddressRegistryWatchCollection).ReplaceOne(
		dbCtx,
		rskAddressIdentity(entry.RskAddress),
		entry,
	)
	if err != nil {
		return err
	}
	if result.MatchedCount != 1 {
		return errors.New("pegin address registry watch entry not found")
	}
	return nil
}

func (repo *peginAddressRegistryWatchMongoRepository) GetCheckpoint(
	ctx context.Context,
) (rootstock.PegInAddressRegistryWatchCheckpoint, bool, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	var checkpointDocument peginAddressRegistryCheckpointDocument
	err := repo.conn.Collection(PegInAddressRegistryWatchCollection).
		FindOne(dbCtx, bson.M{"_id": peginAddressRegistryCheckpointDocumentID}).
		Decode(&checkpointDocument)
	if errors.Is(err, mongoDb.ErrNoDocuments) {
		return rootstock.PegInAddressRegistryWatchCheckpoint{}, false, nil
	}
	if err != nil {
		return rootstock.PegInAddressRegistryWatchCheckpoint{}, false, err
	}
	if checkpointDocument.LocalRoot == nil ||
		checkpointDocument.LastProcessedBlock == nil {
		return rootstock.PegInAddressRegistryWatchCheckpoint{}, false, nil
	}
	return rootstock.PegInAddressRegistryWatchCheckpoint{
		LocalRoot:          *checkpointDocument.LocalRoot,
		LastProcessedBlock: *checkpointDocument.LastProcessedBlock,
	}, true, nil
}

func (repo *peginAddressRegistryWatchMongoRepository) SetCheckpoint(
	ctx context.Context,
	checkpoint rootstock.PegInAddressRegistryWatchCheckpoint,
) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	result, err := repo.conn.Collection(PegInAddressRegistryWatchCollection).UpdateOne(
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

func (repo *peginAddressRegistryWatchMongoRepository) DeleteCheckpoint(ctx context.Context) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	_, err := repo.conn.Collection(PegInAddressRegistryWatchCollection).DeleteOne(
		dbCtx,
		bson.M{"_id": peginAddressRegistryCheckpointDocumentID},
	)
	return err
}

func rskAddressIdentity(rskAddress string) bson.M {
	return bson.M{"rsk_address": rskAddress}
}
