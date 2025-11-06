package liquidity_provider

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type UpdateTrustedAccountUseCase struct {
	trustedAccountRepository liquidity_provider.TrustedAccountRepository
	signer                   entities.Signer
	hashFunc                 entities.HashFunction
}

func NewUpdateTrustedAccountUseCase(
	trustedAccountRepository liquidity_provider.TrustedAccountRepository,
	signer entities.Signer,
	hashFunc entities.HashFunction,
) *UpdateTrustedAccountUseCase {
	return &UpdateTrustedAccountUseCase{
		trustedAccountRepository: trustedAccountRepository,
		signer:                   signer,
		hashFunc:                 hashFunc,
	}
}

func (useCase *UpdateTrustedAccountUseCase) Run(ctx context.Context, account liquidity_provider.TrustedAccountDetails) error {
	signedAccount, err := usecases.SignConfiguration(usecases.UpdateTrustedAccountId, useCase.signer, useCase.hashFunc, account)
	if err != nil {
		return err
	}
	err = useCase.trustedAccountRepository.UpdateTrustedAccount(ctx, signedAccount)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.UpdateTrustedAccountId, err)
	}
	return nil
}
