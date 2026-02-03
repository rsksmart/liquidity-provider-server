package liquidity_provider

import (
	"context"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type InitializeStateConfigurationUseCase struct {
	lpRepository liquidity_provider.LiquidityProviderRepository
	signer       entities.Signer
	hashFunc     entities.HashFunction
}

func NewInitializeStateConfigurationUseCase(
	lpRepository liquidity_provider.LiquidityProviderRepository,
	signer entities.Signer,
	hashFunc entities.HashFunction,
) *InitializeStateConfigurationUseCase {
	return &InitializeStateConfigurationUseCase{
		lpRepository: lpRepository,
		signer:       signer,
		hashFunc:     hashFunc,
	}
}

func (useCase *InitializeStateConfigurationUseCase) Run(ctx context.Context) error {
	signedStateConfig, err := useCase.lpRepository.GetStateConfiguration(ctx)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.InitializeStateConfigurationId, err)
	}

	var stateConfig liquidity_provider.StateConfiguration
	if signedStateConfig != nil {
		stateConfig = signedStateConfig.Value
	}

	modified := false
	now := time.Now()

	// BTC field
	if stateConfig.LastBtcToColdWalletTransfer == nil {
		log.Info("Initializing LastBtcToColdWalletTransfer with current timestamp")
		stateConfig.LastBtcToColdWalletTransfer = &now
		modified = true
	}

	// RBTC field
	if stateConfig.LastRbtcToColdWalletTransfer == nil {
		log.Info("Initializing LastRbtcToColdWalletTransfer with current timestamp")
		stateConfig.LastRbtcToColdWalletTransfer = &now
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
