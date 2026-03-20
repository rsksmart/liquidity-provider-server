package pkg

import (
	"math/big"

	lp "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
)

type SetLiquidityRatioRequest struct {
	BtcPercentage uint64 `json:"btcPercentage" validate:"required,gte=10,lte=90"`
}

type NetworkImpact struct {
	Type   string   `json:"type"`
	Amount *big.Int `json:"amount"`
}

type LiquidityRatioResponse struct {
	BtcPercentage           uint64        `json:"btcPercentage"`
	RbtcPercentage          uint64        `json:"rbtcPercentage"`
	MaxLiquidity            *big.Int      `json:"maxLiquidity"`
	BtcTarget               *big.Int      `json:"btcTarget"`
	BtcThreshold            *big.Int      `json:"btcThreshold"`
	RbtcTarget              *big.Int      `json:"rbtcTarget"`
	RbtcThreshold           *big.Int      `json:"rbtcThreshold"`
	BtcCurrentBalance       *big.Int      `json:"btcCurrentBalance"`
	RbtcCurrentBalance      *big.Int      `json:"rbtcCurrentBalance"`
	BtcImpact               NetworkImpact `json:"btcImpact"`
	RbtcImpact              NetworkImpact `json:"rbtcImpact"`
	CooldownActive          bool          `json:"cooldownActive"`
	CooldownEndTimestamp    int64         `json:"cooldownEndTimestamp"`
	CooldownDurationSeconds int64         `json:"cooldownDurationSeconds"`
	IsPreview               bool          `json:"isPreview"`
}

func ToLiquidityRatioResponse(detail lp.LiquidityRatioDetail) LiquidityRatioResponse {
	return LiquidityRatioResponse{
		BtcPercentage:           detail.BtcPercentage,
		RbtcPercentage:          detail.RbtcPercentage,
		MaxLiquidity:            detail.MaxLiquidity.AsBigInt(),
		BtcTarget:               detail.BtcTarget.AsBigInt(),
		BtcThreshold:            detail.BtcThreshold.AsBigInt(),
		RbtcTarget:              detail.RbtcTarget.AsBigInt(),
		RbtcThreshold:           detail.RbtcThreshold.AsBigInt(),
		BtcCurrentBalance:       detail.BtcCurrentBalance.AsBigInt(),
		RbtcCurrentBalance:      detail.RbtcCurrentBalance.AsBigInt(),
		BtcImpact:               toNetworkImpact(detail.BtcImpact),
		RbtcImpact:              toNetworkImpact(detail.RbtcImpact),
		CooldownActive:          detail.CooldownActive,
		CooldownEndTimestamp:    detail.CooldownEndTimestamp,
		CooldownDurationSeconds: detail.CooldownDurationSeconds,
		IsPreview:               detail.IsPreview,
	}
}

func toNetworkImpact(detail lp.NetworkImpactDetail) NetworkImpact {
	return NetworkImpact{
		Type:   string(detail.Type),
		Amount: detail.Amount.AsBigInt(),
	}
}
