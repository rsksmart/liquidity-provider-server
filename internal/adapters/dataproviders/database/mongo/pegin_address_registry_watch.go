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
	PegInAddressRegistryWatchCollection = "peginAddressRegistryWatch"
	peginAddressRegistryCursorId        = "scanCursor"
)

type peginAddressRegistryWatchMongoRepository struct {
	conn *Connection
}

type peginAddressRegistryCursor struct {
	Id               string `bson:"_id"`
	LastScannedBlock uint64 `bson:"last_scanned_block"`
}

func NewPegInAddressRegistryWatchMongoRepository(conn *Connection) rootstock.PegInAddressRegistryWatchRepository {
	return &peginAddressRegistryWatchMongoRepository{conn: conn}
}

func (repo *peginAddressRegistryWatchMongoRepository) Upsert(
	ctx context.Context,
	watch rootstock.PegInAddressRegistryWatch,
) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	filter := rskAddressIdentity(watch.RskAddress)
	update := bson.M{"$setOnInsert": watch}
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
) (*rootstock.PegInAddressRegistryWatch, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	var watch rootstock.PegInAddressRegistryWatch
	err := repo.conn.Collection(PegInAddressRegistryWatchCollection).
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

func (repo *peginAddressRegistryWatchMongoRepository) List(
	ctx context.Context,
) ([]rootstock.PegInAddressRegistryWatch, error) {
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

	watches := make([]rootstock.PegInAddressRegistryWatch, 0)
	if err = cursor.All(dbCtx, &watches); err != nil {
		return nil, err
	}
	return watches, nil
}

func (repo *peginAddressRegistryWatchMongoRepository) Update(
	ctx context.Context,
	watch rootstock.PegInAddressRegistryWatch,
) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	result, err := repo.conn.Collection(PegInAddressRegistryWatchCollection).ReplaceOne(
		dbCtx,
		rskAddressIdentity(watch.RskAddress),
		watch,
	)
	if err != nil {
		return err
	}
	if result.MatchedCount != 1 {
		return errors.New("pegin address registry watch not found")
	}
	return nil
}

func (repo *peginAddressRegistryWatchMongoRepository) GetCursor(
	ctx context.Context,
) (uint64, bool, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	var cursor peginAddressRegistryCursor
	err := repo.conn.Collection(PegInAddressRegistryWatchCollection).
		FindOne(dbCtx, bson.M{"_id": peginAddressRegistryCursorId}).
		Decode(&cursor)
	if errors.Is(err, mongoDb.ErrNoDocuments) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return cursor.LastScannedBlock, true, nil
}

func (repo *peginAddressRegistryWatchMongoRepository) SetCursor(
	ctx context.Context,
	lastScannedBlock uint64,
) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	_, err := repo.conn.Collection(PegInAddressRegistryWatchCollection).UpdateOne(
		dbCtx,
		bson.M{"_id": peginAddressRegistryCursorId},
		bson.M{"$set": bson.M{"last_scanned_block": lastScannedBlock}},
		options.UpdateOne().SetUpsert(true),
	)
	return err
}

func rskAddressIdentity(rskAddress string) bson.M {
	return bson.M{"rsk_address": rskAddress}
}
