package mongo

import (
	"context"
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDb "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const PegInClaimCollection = "peginClaims"

type peginClaimMongoRepository struct {
	conn *Connection
}

func NewPegInClaimMongoRepository(conn *Connection) rootstock.PegInClaimRepository {
	return &peginClaimMongoRepository{conn: conn}
}

func (repo *peginClaimMongoRepository) Add(ctx context.Context, claim rootstock.PegInClaim) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	filter := bson.M{
		"rsk_address":  claim.RskAddress,
		"deposit_txid": claim.DepositTxID,
		"state":        rootstock.PegInClaimRetryableFailure,
	}
	_, err := repo.conn.Collection(PegInClaimCollection).ReplaceOne(
		dbCtx,
		filter,
		claim,
		options.Replace().SetUpsert(true),
	)
	if mongoDb.IsDuplicateKeyError(err) {
		return rootstock.PegInClaimExistsError
	}
	if err != nil {
		return err
	}
	logDbInteraction(Upsert, claim)
	return nil
}

func (repo *peginClaimMongoRepository) Get(
	ctx context.Context,
	rskAddress, depositTxID string,
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

func (repo *peginClaimMongoRepository) Update(ctx context.Context, claim rootstock.PegInClaim) error {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	result, err := repo.conn.Collection(PegInClaimCollection).ReplaceOne(
		dbCtx,
		pegInClaimIdentity(claim.RskAddress, claim.DepositTxID),
		claim,
	)
	if err != nil {
		return err
	}
	if result.MatchedCount != 1 {
		return rootstock.PegInClaimNotFoundError
	}
	logDbInteraction(Update, claim)
	return nil
}

func pegInClaimIdentity(rskAddress, depositTxID string) bson.M {
	return bson.M{
		"rsk_address":  rskAddress,
		"deposit_txid": depositTxID,
	}
}
