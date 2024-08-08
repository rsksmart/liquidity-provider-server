package liquidity_provider

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type LoginUseCase struct {
	lpRepository           liquidity_provider.LiquidityProviderRepository
	defaultPasswordChannel <-chan entities.Event
	defaultCredentials     *liquidity_provider.HashedCredentials
}

func NewLoginUseCase(
	lpRepository liquidity_provider.LiquidityProviderRepository,
	eventBus entities.EventBus,
) *LoginUseCase {
	evenChannel := eventBus.Subscribe(liquidity_provider.DefaultCredentialsSetEventId)
	return &LoginUseCase{lpRepository: lpRepository, defaultPasswordChannel: evenChannel}
}

func (useCase *LoginUseCase) Run(ctx context.Context, credentials liquidity_provider.Credentials) error {
	if err := ValidateCredentials(ctx, useCase, credentials); err != nil {
		return usecases.WrapUseCaseError(usecases.LoginId, err)
	}
	return nil
}

func (useCase *LoginUseCase) LiquidityProviderRepository() liquidity_provider.LiquidityProviderRepository {
	return useCase.lpRepository
}

func (useCase *LoginUseCase) GetDefaultCredentialsChannel() <-chan entities.Event {
	return useCase.defaultPasswordChannel
}

func (useCase *LoginUseCase) SetDefaultCredentials(credentials *liquidity_provider.HashedCredentials) {
	useCase.defaultCredentials = credentials
}

func (useCase *LoginUseCase) DefaultCredentials() *liquidity_provider.HashedCredentials {
	return useCase.defaultCredentials
}
