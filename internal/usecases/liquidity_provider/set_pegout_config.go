package liquidity_provider

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type SetPegoutConfigUseCase struct {
	lpRepository liquidity_provider.LiquidityProviderRepository
	signer       entities.Signer
	hashFunc     entities.HashFunction
}

func NewSetPegoutConfigUseCase(
	lpRepository liquidity_provider.LiquidityProviderRepository,
	signer entities.Signer,
	hashFunc entities.HashFunction,
) *SetPegoutConfigUseCase {
	return &SetPegoutConfigUseCase{lpRepository: lpRepository, signer: signer, hashFunc: hashFunc}
}

func (useCase *SetPegoutConfigUseCase) Run(ctx context.Context, config liquidity_provider.PegoutConfiguration) error {
	signedConfig, err := usecases.SignConfiguration(usecases.SetPegoutConfigId, useCase.signer, useCase.hashFunc, config)
	if err != nil {
		return err
	}
	err = useCase.lpRepository.UpsertPegoutConfiguration(ctx, signedConfig)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.SetPegoutConfigId, err)
	}
	return nil
}
