package liquidity_provider

import (
	"context"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type SetGeneralConfigUseCase struct {
	lpRepository   liquidity_provider.LiquidityProviderRepository
	peginProvider  liquidity_provider.PeginLiquidityProvider
	pegoutProvider liquidity_provider.PegoutLiquidityProvider
	signer         entities.Signer
	hashFunc       entities.HashFunction
}

func NewSetGeneralConfigUseCase(
	lpRepository liquidity_provider.LiquidityProviderRepository,
	peginProvider liquidity_provider.PeginLiquidityProvider,
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	signer entities.Signer,
	hashFunc entities.HashFunction,
) *SetGeneralConfigUseCase {
	return &SetGeneralConfigUseCase{
		lpRepository:   lpRepository,
		peginProvider:  peginProvider,
		pegoutProvider: pegoutProvider,
		signer:         signer,
		hashFunc:       hashFunc,
	}
}

func (useCase *SetGeneralConfigUseCase) Run(ctx context.Context, config liquidity_provider.GeneralConfiguration) error {
	if err := usecases.ValidateConfirmations(usecases.SetGeneralConfigId, config.RskConfirmations); err != nil {
		return err
	}
	if err := usecases.ValidateConfirmations(usecases.SetGeneralConfigId, config.BtcConfirmations); err != nil {
		return err
	}
	if err := useCase.validateMaxLiquidity(ctx, config.MaxLiquidity); err != nil {
		return err
	}
	if err := usecases.ValidateReimbursementWindowBlocks(usecases.SetGeneralConfigId, config.ReimbursementWindowBlocks); err != nil {
		return err
	}

	signedConfig, err := usecases.SignConfiguration(usecases.SetGeneralConfigId, useCase.signer, useCase.hashFunc, config)
	if err != nil {
		return err
	}
	err = useCase.lpRepository.UpsertGeneralConfiguration(ctx, signedConfig)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.SetGeneralConfigId, err)
	}
	return nil
}

func (useCase *SetGeneralConfigUseCase) validateMaxLiquidity(ctx context.Context, maxLiquidity *entities.Wei) error {
	if err := usecases.ValidatePositiveWeiValues(usecases.SetGeneralConfigId, maxLiquidity); err != nil {
		return err
	}

	peginConfiguration := useCase.peginProvider.PeginConfiguration(ctx)
	pegoutConfiguration := useCase.pegoutProvider.PegoutConfiguration(ctx)

	combinedMinimum := new(entities.Wei).Add(peginConfiguration.MinValue, pegoutConfiguration.MinValue)
	if maxLiquidity.Cmp(combinedMinimum) <= 0 {
		return usecases.WrapUseCaseError(
			usecases.SetGeneralConfigId,
			fmt.Errorf(
				"%w: maxLiquidity %s is not enough to support one minimum pegin and pegout (%s)",
				liquidity_provider.AmountOutOfRangeError,
				maxLiquidity.String(),
				combinedMinimum.String(),
			),
		)
	}
	return nil
}
