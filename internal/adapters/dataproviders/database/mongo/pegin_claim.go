package mongo

import (
	"context"
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDb "go.mongodb.org/mongo-driver/v2/mongo"
)

const PegInClaimCollection = "peginClaims"

type peginClaimMongoRepository struct {
	conn *Connection
}

func NewPegInClaimMongoRepository(conn *Connection) rootstock.PegInClaimRepository {
	return &peginClaimMongoRepository{conn: conn}
}

func (repo *peginClaimMongoRepository) Insert(ctx context.Context, claim rootstock.PegInClaim) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	_, err := repo.conn.Collection(PegInClaimCollection).InsertOne(dbCtx, claim)
	if mongoDb.IsDuplicateKeyError(err) {
		return rootstock.ErrPegInClaimAlreadyExists
	}
	if err != nil {
		return err
	}
	logDbInteraction(Insert, claim)
	return nil
}

func (repo *peginClaimMongoRepository) Get(
	ctx context.Context,
	rskAddress string,
	depositTxID string,
) (*rootstock.PegInClaim, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	var claim rootstock.PegInClaim
	err := repo.conn.Collection(PegInClaimCollection).
		FindOne(dbCtx, pegInClaimIdentity(rskAddress, depositTxID)).
		Decode(&claim)
	if errors.Is(err, mongoDb.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	logDbInteraction(Read, claim)
	return &claim, nil
}

func pegInClaimIdentity(rskAddress, depositTxID string) bson.M {
	return bson.M{
		"rsk_address":  rskAddress,
		"deposit_txid": depositTxID,
	}
}
