package storage

import (
	"github.com/rsksmart/liquidity-provider-server/connectors"
	mongoDB "github.com/rsksmart/liquidity-provider-server/mongo"
	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider-server/pegout"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
)

type LPRepository struct {
	dbMongo *mongoDB.DB
	rsk     connectors.RSKConnector
}

func NewLPRepository(dbMongo *mongoDB.DB, rsk connectors.RSKConnector) *LPRepository {
	return &LPRepository{dbMongo, rsk}
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

func (r *LPRepository) HasLiquidity(lp pegin.LiquidityProvider, wei *types.Wei) (bool, error) {
	log.Debug("Verifying if has liquidity")
	lpBalance, err := r.rsk.GetAvailableLiquidity(lp.Address())
	log.Debug("LP balance ", lpBalance.Text(10))
	if err != nil {
		return false, err
	}
	lockedLiquidity, err := r.dbMongo.GetLockedLiquidity()
	log.Debug("Locked Liquidity ", lockedLiquidity.String())
	if err != nil {
		return false, err
	}
	return new(types.Wei).Sub(types.NewBigWei(lpBalance), lockedLiquidity).Cmp(wei) >= 0, nil
}

func (r *LPRepository) RetainPegOutQuote(rq *pegout.RetainedQuote) error {
	return r.dbMongo.RetainPegOutQuote(rq)
}

func (r *LPRepository) HasRetainedPegOutQuote(hash string) (bool, error) {
	rq, err := r.dbMongo.GetRetainedPegOutQuote(hash)
	if err != nil {
		return false, err
	}
	return rq != nil, nil
}

func (r *LPRepository) HasLiquidityPegOut(lp pegout.LiquidityProvider, satoshis uint64) (bool, error) {
	return true, nil
}
