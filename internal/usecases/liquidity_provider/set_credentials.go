package liquidity_provider

import (
	"context"
	"encoding/hex"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

const credentialSaltSize = 32

type SetCredentialsUseCase struct {
	lpRepository           liquidity_provider.LiquidityProviderRepository
	signer                 entities.Signer
	hashFunc               entities.HashFunction
	defaultPassword        string
	defaultPasswordChannel <-chan entities.Event
}

func NewSetCredentialsUseCase(
	lpRepository liquidity_provider.LiquidityProviderRepository,
	signer entities.Signer,
	hashFunc entities.HashFunction,
	eventBus entities.EventBus,
) *SetCredentialsUseCase {
	defaultPasswordChannel := eventBus.Subscribe(liquidity_provider.DefaultCredentialsSetEventId)
	return &SetCredentialsUseCase{
		lpRepository:           lpRepository,
		signer:                 signer,
		hashFunc:               hashFunc,
		defaultPasswordChannel: defaultPasswordChannel,
	}
}

func (useCase *SetCredentialsUseCase) Run(ctx context.Context, oldCredentials, newCredentials liquidity_provider.Credentials) error {
	if err := ValidateCredentials(ctx, useCase, oldCredentials); err != nil {
		return usecases.WrapUseCaseError(usecases.ChangeCredentialsId, err)
	}

	rules := utils.DefaultPasswordValidationRuleset()
	if err := utils.CheckPasswordComplexity(newCredentials.Password, rules...); err != nil {
		return usecases.WrapUseCaseError(usecases.ChangeCredentialsId, err)
	}

	hashedUsername, usernameSalt, err := utils.HashAndSaltArgon2(newCredentials.Username, credentialSaltSize)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ChangeCredentialsId, err)
	}
	hashedPassword, passwordSalt, err := utils.HashAndSaltArgon2(newCredentials.Password, credentialSaltSize)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ChangeCredentialsId, err)
	}

	hashedCredentials := liquidity_provider.HashedCredentials{
		HashedUsername: hex.EncodeToString(hashedUsername),
		HashedPassword: hex.EncodeToString(hashedPassword),
		UsernameSalt:   usernameSalt,
		PasswordSalt:   passwordSalt,
	}
	signedCredentials, err := usecases.SignConfiguration(usecases.ChangeCredentialsId, useCase.signer, useCase.hashFunc, hashedCredentials)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ChangeCredentialsId, err)
	}
	if err = useCase.lpRepository.UpsertCredentials(ctx, signedCredentials); err != nil {
		return usecases.WrapUseCaseError(usecases.ChangeCredentialsId, err)
	}
	return nil
}

func (useCase *SetCredentialsUseCase) LiquidityProviderRepository() liquidity_provider.LiquidityProviderRepository {
	return useCase.lpRepository
}

func (useCase *SetCredentialsUseCase) GetDefaultPasswordChannel() <-chan entities.Event {
	return useCase.defaultPasswordChannel
}

func (useCase *SetCredentialsUseCase) SetDefaultPassword(password string) {
	useCase.defaultPassword = password
}

func (useCase *SetCredentialsUseCase) DefaultPassword() string {
	return useCase.defaultPassword
}
