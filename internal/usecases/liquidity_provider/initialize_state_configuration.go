package liquidity_provider

import (
	"context"
	"errors"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type InitializeStateConfigurationUseCase struct {
	provider     liquidity_provider.LiquidityProvider
	lpRepository liquidity_provider.LiquidityProviderRepository
	signer       entities.Signer
	hashFunc     entities.HashFunction
}

func NewInitializeStateConfigurationUseCase(
	provider liquidity_provider.LiquidityProvider,
	lpRepository liquidity_provider.LiquidityProviderRepository,
	signer entities.Signer,
	hashFunc entities.HashFunction,
) *InitializeStateConfigurationUseCase {
	return &InitializeStateConfigurationUseCase{
		provider:     provider,
		lpRepository: lpRepository,
		signer:       signer,
		hashFunc:     hashFunc,
	}
}

func (useCase *InitializeStateConfigurationUseCase) Run(ctx context.Context) error {
	stateConfig, err := useCase.provider.StateConfiguration(ctx)
	if err != nil && !errors.Is(err, liquidity_provider.ConfigurationNotFoundError) {
		return usecases.WrapUseCaseError(usecases.InitializeStateConfigurationId, err)
	}

	modified := false
	now := time.Now().UTC().Unix()

	if stateConfig.LastBtcToColdWalletTransfer == 0 {
		log.Info("Initializing LastBtcToColdWalletTransfer with current timestamp")
		stateConfig.LastBtcToColdWalletTransfer = now
		modified = true
	}

	if stateConfig.LastRbtcToColdWalletTransfer == 0 {
		log.Info("Initializing LastRbtcToColdWalletTransfer with current timestamp")
		stateConfig.LastRbtcToColdWalletTransfer = now
		modified = true
	}

	if !modified {
		log.Debug("State configuration already fully initialized")
		return nil
	}

	signedConfig, err := usecases.SignConfiguration(usecases.InitializeStateConfigurationId, useCase.signer, useCase.hashFunc, stateConfig)
	if err != nil {
		return err
	}

	if err := useCase.lpRepository.UpsertStateConfiguration(ctx, signedConfig); err != nil {
		return usecases.WrapUseCaseError(usecases.InitializeStateConfigurationId, err)
	}

	log.Info("State configuration initialized successfully")
	return nil
}
