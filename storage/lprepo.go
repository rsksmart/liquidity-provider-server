package storage

import (
	"github.com/rsksmart/liquidity-provider-server/connectors"
	mongoDB "github.com/rsksmart/liquidity-provider-server/mongo"
	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider-server/pegout"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
	"math/big"
)

type LPRepository struct {
	dbMongo mongoDB.DBConnector
	rsk     connectors.RSKConnector
	btc     connectors.BTCConnector
}

func NewLPRepository(dbMongo mongoDB.DBConnector, rsk connectors.RSKConnector, btc connectors.BTCConnector) *LPRepository {
	return &LPRepository{dbMongo, rsk, btc}
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

func (r *LPRepository) HasLiquidityPegOut(satoshis uint64) (bool, error) {
	log.Debug("Verifying if has liquidity")

	lpBalance, err := r.btc.GetAvailableLiquidity()
	log.Debugf("LP balance %v satoshis\n", lpBalance)
	if err != nil {
		return false, err
	}
	lockedLiquidity, err := r.dbMongo.GetLockedLiquidityPegOut()
	log.Debugf("Locked Liquidity %d satoshis\n", lockedLiquidity)
	if err != nil {
		return false, err
	}

	return new(big.Int).Sub(lpBalance, big.NewInt(int64(lockedLiquidity))).Cmp(big.NewInt(int64(satoshis))) >= 0, nil
}
