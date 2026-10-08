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
		return rootstock.ErrPegInClaimNotFound
	}
	logDbInteraction(Update, claim)
	return nil
}

func (repo *peginClaimMongoRepository) ListByStates(
	ctx context.Context,
	states ...rootstock.PegInClaimState,
) ([]rootstock.PegInClaim, error) {
	dbCtx, cancel := context.WithTimeout(ctx, repo.conn.timeout)
	defer cancel()

	cursor, err := repo.conn.Collection(PegInClaimCollection).Find(dbCtx, pegInClaimStatesFilter(states))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(dbCtx)

	claims := make([]rootstock.PegInClaim, 0)
	if err = cursor.All(dbCtx, &claims); err != nil {
		return nil, err
	}
	logDbInteraction(Read, claims)
	return claims, nil
}

func pegInClaimIdentity(rskAddress, depositTxID string) bson.M {
	return bson.M{
		"rsk_address":  rskAddress,
		"deposit_txid": depositTxID,
	}
}

func pegInClaimStatesFilter(states []rootstock.PegInClaimState) bson.M {
	if len(states) == 0 {
		return bson.M{}
	}
	return bson.M{"state": bson.M{"$in": states}}
}
