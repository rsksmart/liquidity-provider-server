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
}

func NewInitializeStateConfigurationUseCase(
	lpRepository liquidity_provider.LiquidityProviderRepository,
) *InitializeStateConfigurationUseCase {
	return &InitializeStateConfigurationUseCase{lpRepository: lpRepository}
}

func (useCase *InitializeStateConfigurationUseCase) Run(ctx context.Context) error {
	stateConfig, err := useCase.lpRepository.GetStateConfiguration(ctx)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.InitializeStateConfigurationId, err)
	}
	if stateConfig != nil {
		log.Debug("State configuration already initialized")
		return nil
	}

	// If it doesn't exist, create it with current timestamps
	log.Info("Initializing state configuration with current timestamps...")
	now := time.Now()
	newStateConfig := liquidity_provider.StateConfiguration{
		LastBtcToColdWalletTransfer:  &now,
		LastRbtcToColdWalletTransfer: &now,
	}

	signedConfig := entities.Signed[liquidity_provider.StateConfiguration]{
		Value:     newStateConfig,
		Signature: "",
	}

	if err := useCase.lpRepository.UpsertStateConfiguration(ctx, signedConfig); err != nil {
		return usecases.WrapUseCaseError(usecases.InitializeStateConfigurationId, err)
	}

	log.Info("State configuration initialized successfully")
	return nil
}
