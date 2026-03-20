package liquidity_provider

import (
	"context"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type SetLiquidityRatioUseCase struct {
	generalProvider liquidity_provider.LiquidityProvider
	lpRepository    liquidity_provider.LiquidityProviderRepository
	signer          entities.Signer
	hashFunc        entities.HashFunction
}

func NewSetLiquidityRatioUseCase(
	generalProvider liquidity_provider.LiquidityProvider,
	lpRepository liquidity_provider.LiquidityProviderRepository,
	signer entities.Signer,
	hashFunc entities.HashFunction,
) *SetLiquidityRatioUseCase {
	return &SetLiquidityRatioUseCase{
		generalProvider: generalProvider,
		lpRepository:    lpRepository,
		signer:          signer,
		hashFunc:        hashFunc,
	}
}

func (useCase *SetLiquidityRatioUseCase) Run(ctx context.Context, btcPercentage uint64) error {
	stateConfig, err := useCase.generalProvider.StateConfiguration(ctx)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.SetLiquidityRatioId, err)
	}

	if stateConfig.BtcLiquidityTargetPercentage == btcPercentage {
		return nil
	}

	stateConfig.BtcLiquidityTargetPercentage = btcPercentage
	stateConfig.RatioCooldownEndTimestamp = time.Now().Unix() + CooldownAfterRatioChange

	signedConfig, err := usecases.SignConfiguration(usecases.SetLiquidityRatioId, useCase.signer, useCase.hashFunc, stateConfig)
	if err != nil {
		return err
	}

	if err := useCase.lpRepository.UpsertStateConfiguration(ctx, signedConfig); err != nil {
		return usecases.WrapUseCaseError(usecases.SetLiquidityRatioId, err)
	}

	return nil
}
