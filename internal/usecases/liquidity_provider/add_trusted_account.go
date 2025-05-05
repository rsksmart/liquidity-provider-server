package liquidity_provider

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type AddTrustedAccountUseCase struct {
	trustedAccountRepository liquidity_provider.TrustedAccountRepository
	signer                   entities.Signer
	hashFunc                 entities.HashFunction
}

func NewAddTrustedAccountUseCase(
	trustedAccountRepository liquidity_provider.TrustedAccountRepository,
	signer entities.Signer,
	hashFunc entities.HashFunction,
) *AddTrustedAccountUseCase {
	return &AddTrustedAccountUseCase{
		trustedAccountRepository: trustedAccountRepository,
		signer:                   signer,
		hashFunc:                 hashFunc,
	}
}

func (useCase *AddTrustedAccountUseCase) Run(ctx context.Context, account liquidity_provider.TrustedAccountDetails) error {
	signedAccount, err := usecases.SignConfiguration(usecases.AddTrustedAccountId, useCase.signer, useCase.hashFunc, account)
	if err != nil {
		return err
	}
	err = useCase.trustedAccountRepository.AddTrustedAccount(ctx, signedAccount)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.AddTrustedAccountId, err)
	}
	return nil
}
