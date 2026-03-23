package pkg_test

import (
	"math/big"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lp "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
)

func TestToLiquidityRatioResponse(t *testing.T) {
	detail := lp.LiquidityRatioDetail{
		BtcPercentage:      60,
		RbtcPercentage:     40,
		MaxLiquidity:       entities.NewBigWei(big.NewInt(10_000_000_000)),
		BtcTarget:          entities.NewBigWei(big.NewInt(6_000_000_000)),
		BtcThreshold:       entities.NewBigWei(big.NewInt(6_500_000_000)),
		RbtcTarget:         entities.NewBigWei(big.NewInt(4_000_000_000)),
		RbtcThreshold:      entities.NewBigWei(big.NewInt(4_500_000_000)),
		BtcCurrentBalance:  entities.NewBigWei(big.NewInt(5_000_000_000)),
		RbtcCurrentBalance: entities.NewBigWei(big.NewInt(5_000_000_000)),
		BtcImpact: lp.NetworkImpactDetail{
			Type:   lp.NetworkImpactDeficit,
			Amount: entities.NewBigWei(big.NewInt(1_000_000_000)),
		},
		RbtcImpact: lp.NetworkImpactDetail{
			Type:   lp.NetworkImpactExcess,
			Amount: entities.NewBigWei(big.NewInt(1_000_000_000)),
		},
		CooldownActive:          true,
		CooldownEndTimestamp:    1700000000,
		CooldownDurationSeconds: 3600,
		IsPreview:               true,
	}

	response := pkg.ToLiquidityRatioResponse(detail)

	assert.Equal(t, uint64(60), response.BtcPercentage)
	assert.Equal(t, uint64(40), response.RbtcPercentage)
	assert.Equal(t, "10000000000", response.MaxLiquidity.String())
	assert.Equal(t, "6000000000", response.BtcTarget.String())
	assert.Equal(t, "6500000000", response.BtcThreshold.String())
	assert.Equal(t, "4000000000", response.RbtcTarget.String())
	assert.Equal(t, "4500000000", response.RbtcThreshold.String())
	assert.Equal(t, "5000000000", response.BtcCurrentBalance.String())
	assert.Equal(t, "5000000000", response.RbtcCurrentBalance.String())
	assert.Equal(t, "deficit", response.BtcImpact.Type)
	assert.Equal(t, "1000000000", response.BtcImpact.Amount.String())
	assert.Equal(t, "excess", response.RbtcImpact.Type)
	assert.Equal(t, "1000000000", response.RbtcImpact.Amount.String())
	assert.True(t, response.CooldownActive)
	assert.Equal(t, int64(1700000000), response.CooldownEndTimestamp)
	assert.Equal(t, int64(3600), response.CooldownDurationSeconds)
	assert.True(t, response.IsPreview)
	test.AssertNonZeroValues(t, response)
}
