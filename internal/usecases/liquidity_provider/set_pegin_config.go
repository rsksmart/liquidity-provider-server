package liquidity_provider

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type SetPeginConfigUseCase struct {
	lpRepository liquidity_provider.LiquidityProviderRepository
	signer       entities.Signer
	hashFunc     entities.HashFunction
	contracts    blockchain.RskContracts
}

func NewSetPeginConfigUseCase(
	lpRepository liquidity_provider.LiquidityProviderRepository,
	signer entities.Signer,
	hashFunc entities.HashFunction,
	contracts blockchain.RskContracts,
) *SetPeginConfigUseCase {
	return &SetPeginConfigUseCase{lpRepository: lpRepository, signer: signer, hashFunc: hashFunc, contracts: contracts}
}

func (useCase *SetPeginConfigUseCase) Run(ctx context.Context, config liquidity_provider.PeginConfiguration) error {
	if err := usecases.ValidatePositiveWeiValues(
		usecases.SetPeginConfigId,
		config.PenaltyFee,
		config.FixedFee,
		config.MaxValue,
		config.MinValue,
	); err != nil {
		return err
	}

	if err := usecases.ValidateMinLockValue(usecases.SetPeginConfigId, useCase.contracts.Bridge, config.MinValue); err != nil {
		return err
	}
	signedConfig, err := usecases.SignConfiguration(usecases.SetPeginConfigId, useCase.signer, useCase.hashFunc, config)
	if err != nil {
		return err
	}
	if err = useCase.lpRepository.UpsertPeginConfiguration(ctx, signedConfig); err != nil {
		return usecases.WrapUseCaseError(usecases.SetPeginConfigId, err)
	}
	return nil
}
