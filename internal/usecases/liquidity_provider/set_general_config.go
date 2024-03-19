package liquidity_provider

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type SetGeneralConfigUseCase struct {
	lpRepository liquidity_provider.LiquidityProviderRepository
	signer       entities.Signer
	hashFunc     entities.HashFunction
}

func NewSetGeneralConfigUseCase(
	lpRepository liquidity_provider.LiquidityProviderRepository,
	signer entities.Signer,
	hashFunc entities.HashFunction,
) *SetGeneralConfigUseCase {
	return &SetGeneralConfigUseCase{lpRepository: lpRepository, signer: signer, hashFunc: hashFunc}
}

func (useCase *SetGeneralConfigUseCase) Run(ctx context.Context, config liquidity_provider.GeneralConfiguration) error {
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
