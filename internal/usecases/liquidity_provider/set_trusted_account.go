package liquidity_provider

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type SetTrustedAccountUseCase struct {
	trustedAccountRepository liquidity_provider.TrustedAccountRepository
	signer                   entities.Signer
	hashFunc                 entities.HashFunction
}

func NewSetTrustedAccountUseCase(
	trustedAccountRepository liquidity_provider.TrustedAccountRepository,
	signer entities.Signer,
	hashFunc entities.HashFunction,
) *SetTrustedAccountUseCase {
	return &SetTrustedAccountUseCase{
		trustedAccountRepository: trustedAccountRepository,
		signer:                   signer,
		hashFunc:                 hashFunc,
	}
}

func (useCase *SetTrustedAccountUseCase) Run(ctx context.Context, account liquidity_provider.TrustedAccountDetails) error {
	signedAccount, err := usecases.SignConfiguration(usecases.SetTrustedAccountId, useCase.signer, useCase.hashFunc, account)
	if err != nil {
		return err
	}
	err = useCase.trustedAccountRepository.UpdateTrustedAccount(ctx, signedAccount)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.SetTrustedAccountId, err)
	}
	return nil
}
