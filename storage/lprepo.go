package storage

import (
	"github.com/rsksmart/liquidity-provider-server/connectors"
	mongoDB "github.com/rsksmart/liquidity-provider-server/mongo"
	"github.com/rsksmart/liquidity-provider-server/pegout"
	"github.com/rsksmart/liquidity-provider/providers"
	"github.com/rsksmart/liquidity-provider/types"
)

type LPRepository struct {
	db      DBConnector
	dbMongo *mongoDB.DB
	rsk     connectors.RSKConnector
}

func NewLPRepository(db DBConnector, dbMongo *mongoDB.DB, rsk connectors.RSKConnector) *LPRepository {
	return &LPRepository{db, dbMongo, rsk}
}

func (r *LPRepository) RetainQuote(rq *types.RetainedQuote) error {
	return r.dbMongo.RetainQuote(rq)
}

func (r *LPRepository) HasRetainedQuote(hash string) (bool, error) {
	rq, err := r.dbMongo.GetRetainedQuote(hash)
	if err != nil {
		return false, err
	}
	return rq != nil, nil
}

func (r *LPRepository) HasLiquidity(lp providers.LiquidityProvider, wei *types.Wei) (bool, error) {
	availableLiq, err := r.rsk.GetAvailableLiquidity(lp.Address())
	if err != nil {
		return false, err
	}
	lockedLiq, err := r.dbMongo.GetLockedLiquidity()
	if err != nil {
		return false, err
	}
	return new(types.Wei).Sub(types.NewBigWei(availableLiq), lockedLiq).Cmp(wei) >= 0, nil
}

func (r *LPRepository) RetainPegOutQuote(rq *pegout.RetainedQuote) error {
	return r.db.RetainPegOutQuote(rq)
}

func (r *LPRepository) HasRetainedPegOutQuote(hash string) (bool, error) {
	rq, err := r.db.GetRetainedPegOutQuote(hash)
	if err != nil {
		return false, err
	}
	return rq != nil, nil
}

func (r *LPRepository) HasLiquidityPegOut(lp pegout.LiquidityProvider, satoshis uint64) (bool, error) {
	return true, nil
}
