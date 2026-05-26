package liquidity_provider

import (
	"context"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type NetworkImpactType string

const (
	NetworkImpactExcess          NetworkImpactType = "excess"
	NetworkImpactDeficit         NetworkImpactType = "deficit"
	NetworkImpactWithinTolerance NetworkImpactType = "withinTolerance"
)

type NetworkImpactDetail struct {
	Type   NetworkImpactType
	Amount *entities.Wei
}

type LiquidityRatioDetail struct {
	BtcPercentage           uint64
	RbtcPercentage          uint64
	MaxLiquidity            *entities.Wei
	BtcTarget               *entities.Wei
	BtcThreshold            *entities.Wei
	RbtcTarget              *entities.Wei
	RbtcThreshold           *entities.Wei
	BtcCurrentBalance       *entities.Wei
	RbtcCurrentBalance      *entities.Wei
	BtcImpact               NetworkImpactDetail
	RbtcImpact              NetworkImpactDetail
	CooldownActive          bool
	CooldownEndTimestamp    int64
	CooldownDurationSeconds int64
	IsPreview               bool
}

type GetLiquidityRatioUseCase struct {
	generalProvider liquidity_provider.LiquidityProvider
	peginProvider   liquidity_provider.PeginLiquidityProvider
	pegoutProvider  liquidity_provider.PegoutLiquidityProvider
}

func NewGetLiquidityRatioUseCase(
	generalProvider liquidity_provider.LiquidityProvider,
	peginProvider liquidity_provider.PeginLiquidityProvider,
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
) *GetLiquidityRatioUseCase {
	return &GetLiquidityRatioUseCase{
		generalProvider: generalProvider,
		peginProvider:   peginProvider,
		pegoutProvider:  pegoutProvider,
	}
}

func (useCase *GetLiquidityRatioUseCase) Run(ctx context.Context, proposedBtcPercentage uint64) (LiquidityRatioDetail, error) {
	generalConfig := useCase.generalProvider.GeneralConfiguration(ctx)
	stateConfig, err := useCase.generalProvider.StateConfiguration(ctx)
	if err != nil {
		return LiquidityRatioDetail{}, usecases.WrapUseCaseError(usecases.GetLiquidityRatioId, err)
	}

	btcBalance, err := useCase.pegoutProvider.AvailablePegoutLiquidity(ctx)
	if err != nil {
		return LiquidityRatioDetail{}, usecases.WrapUseCaseError(usecases.GetLiquidityRatioId, err)
	}

	rbtcBalance, err := useCase.peginProvider.AvailablePeginLiquidity(ctx)
	if err != nil {
		return LiquidityRatioDetail{}, usecases.WrapUseCaseError(usecases.GetLiquidityRatioId, err)
	}

	isPreview := proposedBtcPercentage > 0
	btcPercentage := stateConfig.BtcLiquidityTargetPercentage
	if isPreview {
		btcPercentage = proposedBtcPercentage
	}
	rbtcPercentage := 100 - btcPercentage

	btcTarget, err := useCase.calculateTarget(generalConfig.MaxLiquidity, btcPercentage)
	if err != nil {
		return LiquidityRatioDetail{}, usecases.WrapUseCaseError(usecases.GetLiquidityRatioId, err)
	}

	rbtcTarget, err := useCase.calculateTarget(generalConfig.MaxLiquidity, rbtcPercentage)
	if err != nil {
		return LiquidityRatioDetail{}, usecases.WrapUseCaseError(usecases.GetLiquidityRatioId, err)
	}

	btcThreshold := generalConfig.ExcessTolerance.ComputeThreshold(btcTarget)
	rbtcThreshold := generalConfig.ExcessTolerance.ComputeThreshold(rbtcTarget)

	return LiquidityRatioDetail{
		BtcPercentage:           btcPercentage,
		RbtcPercentage:          rbtcPercentage,
		MaxLiquidity:            generalConfig.MaxLiquidity,
		BtcTarget:               btcTarget,
		BtcThreshold:            btcThreshold,
		RbtcTarget:              rbtcTarget,
		RbtcThreshold:           rbtcThreshold,
		BtcCurrentBalance:       btcBalance,
		RbtcCurrentBalance:      rbtcBalance,
		BtcImpact:               useCase.calculateImpact(btcBalance, btcTarget, btcThreshold),
		RbtcImpact:              useCase.calculateImpact(rbtcBalance, rbtcTarget, rbtcThreshold),
		CooldownActive:          time.Now().Unix() < stateConfig.RatioCooldownEndTimestamp,
		CooldownEndTimestamp:    stateConfig.RatioCooldownEndTimestamp,
		CooldownDurationSeconds: CooldownAfterRatioChange,
		IsPreview:               isPreview,
	}, nil
}

func (useCase *GetLiquidityRatioUseCase) calculateTarget(maxLiquidity *entities.Wei, percentage uint64) (*entities.Wei, error) {
	return new(entities.Wei).Div(
		new(entities.Wei).Mul(maxLiquidity, entities.NewUWei(percentage)),
		entities.NewUWei(100),
	)
}

func (useCase *GetLiquidityRatioUseCase) calculateImpact(balance, target, threshold *entities.Wei) NetworkImpactDetail {
	if balance.Cmp(target) < 0 {
		return NetworkImpactDetail{
			Type:   NetworkImpactDeficit,
			Amount: new(entities.Wei).Sub(target, balance),
		}
	}
	if balance.Cmp(threshold) > 0 {
		return NetworkImpactDetail{
			Type:   NetworkImpactExcess,
			Amount: new(entities.Wei).Sub(balance, target),
		}
	}
	return NetworkImpactDetail{
		Type:   NetworkImpactWithinTolerance,
		Amount: entities.NewWei(0),
	}
}
