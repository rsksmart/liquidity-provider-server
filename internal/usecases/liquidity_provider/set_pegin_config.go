package liquidity_provider

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type SetPeginConfigUseCase struct {
	lpRepository liquidity_provider.LiquidityProviderRepository
	signer       entities.Signer
	hashFunc     entities.HashFunction
}

func NewSetPeginConfigUseCase(
	lpRepository liquidity_provider.LiquidityProviderRepository,
	signer entities.Signer,
	hashFunc entities.HashFunction,
) *SetPeginConfigUseCase {
	return &SetPeginConfigUseCase{lpRepository: lpRepository, signer: signer, hashFunc: hashFunc}
}

func (useCase *SetPeginConfigUseCase) Run(ctx context.Context, config liquidity_provider.PeginConfiguration) error {
	signedConfig, err := usecases.SignConfiguration(usecases.SetPeginConfigId, useCase.signer, useCase.hashFunc, config)
	if err != nil {
		return err
	}
	return useCase.lpRepository.UpsertPeginConfiguration(ctx, signedConfig)
}
