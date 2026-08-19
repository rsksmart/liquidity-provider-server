package mongo

import (
	"context"
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDb "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	PegOutEscrowWatchCollection = "pegoutEscrowWatch"
	pegOutEscrowWatchCursorId   = "scanCursor"
)

type pegOutEscrowWatchMongoRepository struct {
	conn *Connection
}

type pegOutEscrowWatchCursor struct {
	Id               string `bson:"_id"`
	LastScannedBlock uint64 `bson:"last_scanned_block"`
}

type storedPegOutEscrowCandidate struct {
	Id                         string `bson:"_id"`
	blockchain.PegOutRequested `bson:",inline"`
}

func NewPegOutEscrowWatchMongoRepository(conn *Connection) blockchain.PegOutEscrowWatchRepository {
	return &pegOutEscrowWatchMongoRepository{conn: conn}
}

func (repo *pegOutEscrowWatchMongoRepository) GetCheckpoint(ctx context.Context) (uint64, bool, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	var cursor pegOutEscrowWatchCursor
	err := repo.conn.Collection(PegOutEscrowWatchCollection).
		FindOne(dbCtx, bson.M{"_id": pegOutEscrowWatchCursorId}).
		Decode(&cursor)
	if errors.Is(err, mongoDb.ErrNoDocuments) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	logDbInteraction(Read, cursor)
	return cursor.LastScannedBlock, true, nil
}

func (repo *pegOutEscrowWatchMongoRepository) SetCheckpoint(ctx context.Context, lastScannedBlock uint64) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	cursor := pegOutEscrowWatchCursor{
		Id:               pegOutEscrowWatchCursorId,
		LastScannedBlock: lastScannedBlock,
	}
	_, err := repo.conn.Collection(PegOutEscrowWatchCollection).UpdateOne(
		dbCtx,
		bson.M{"_id": pegOutEscrowWatchCursorId},
		bson.M{"$set": bson.M{"last_scanned_block": lastScannedBlock}},
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		return err
	}
	logDbInteraction(Upsert, cursor)
	return nil
}

func (repo *pegOutEscrowWatchMongoRepository) UpsertCandidate(ctx context.Context, candidate blockchain.PegOutRequested) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	stored := storedPegOutEscrowCandidate{
		Id:              candidate.RequestHash,
		PegOutRequested: candidate,
	}
	_, err := repo.conn.Collection(PegOutEscrowWatchCollection).ReplaceOne(
		dbCtx,
		bson.M{"_id": candidate.RequestHash},
		stored,
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		return err
	}
	logDbInteraction(Upsert, stored)
	return nil
}

func (repo *pegOutEscrowWatchMongoRepository) DeleteCandidate(ctx context.Context, requestHash string) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	_, err := repo.conn.Collection(PegOutEscrowWatchCollection).DeleteOne(dbCtx, bson.M{"_id": requestHash})
	if err != nil {
		return err
	}
	logDbInteraction(Delete, requestHash)
	return nil
}

func (repo *pegOutEscrowWatchMongoRepository) ListCandidates(ctx context.Context) ([]blockchain.PegOutRequested, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	cursor, err := repo.conn.Collection(PegOutEscrowWatchCollection).Find(
		dbCtx,
		bson.M{"request_hash": bson.M{"$exists": true}},
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = cursor.Close(dbCtx) }()

	stored := make([]storedPegOutEscrowCandidate, 0)
	if err = cursor.All(dbCtx, &stored); err != nil {
		return nil, err
	}
	candidates := make([]blockchain.PegOutRequested, 0, len(stored))
	for _, item := range stored {
		candidates = append(candidates, item.PegOutRequested)
	}
	logDbInteraction(Read, candidates)
	return candidates, nil
}
