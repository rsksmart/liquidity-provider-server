package storage

import (
	"github.com/rsksmart/liquidity-provider-server/connectors"
	"github.com/rsksmart/liquidity-provider/providers"
	"github.com/rsksmart/liquidity-provider/types"
)

type LPRepository struct {
	db  DBConnector
	rsk connectors.RSKConnector
}

func NewLPRepository(db DBConnector, rsk connectors.RSKConnector) *LPRepository {
	return &LPRepository{db, rsk}
}

func (r *LPRepository) RetainQuote(rq *types.RetainedQuote) error {
	return r.db.RetainQuote(rq)
}

func (r *LPRepository) HasRetainedQuote(hash string) (bool, error) {
	rq, err := r.db.GetRetainedQuote(hash)
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
	lockedLiq, err := r.db.GetLockedLiquidity()
	if err != nil {
		return false, err
	}
	return new(types.Wei).Sub(types.NewBigWei(availableLiq), lockedLiq).Cmp(wei) >= 0, nil
}
