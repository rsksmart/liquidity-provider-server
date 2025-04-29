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
	signedAccount, err := usecases.SignTrustedAccount(usecases.SetTrustedAccountId, useCase.signer, useCase.hashFunc, account)
	if err != nil {
		return err
	}
	account.Signature = signedAccount.Signature
	account.Hash = signedAccount.Hash
	err = useCase.trustedAccountRepository.UpdateTrustedAccount(ctx, account)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.SetTrustedAccountId, err)
	}
	return nil
}

func (useCase *SetTrustedAccountUseCase) Add(ctx context.Context, account liquidity_provider.TrustedAccountDetails) error {
	signedAccount, err := usecases.SignTrustedAccount(usecases.SetTrustedAccountId, useCase.signer, useCase.hashFunc, account)
	if err != nil {
		return err
	}
	account.Signature = signedAccount.Signature
	account.Hash = signedAccount.Hash
	err = useCase.trustedAccountRepository.AddTrustedAccount(ctx, account)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.SetTrustedAccountId, err)
	}
	return nil
}

func (useCase *SetTrustedAccountUseCase) Delete(ctx context.Context, address string) error {
	err := useCase.trustedAccountRepository.DeleteTrustedAccount(ctx, address)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.SetTrustedAccountId, err)
	}
	return nil
}
