package liquidity_provider

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
)

type GetTrustedAccountsUseCase struct {
	trustedAccountRepository liquidity_provider.TrustedAccountRepository
	hashFunction             entities.HashFunction
	signer                   entities.Signer
}

func NewGetTrustedAccountsUseCase(
	trustedAccountRepository liquidity_provider.TrustedAccountRepository,
	hashFunction entities.HashFunction,
	signer entities.Signer,
) *GetTrustedAccountsUseCase {
	return &GetTrustedAccountsUseCase{
		trustedAccountRepository: trustedAccountRepository,
		hashFunction:             hashFunction,
		signer:                   signer,
	}
}

func (useCase *GetTrustedAccountsUseCase) Run(ctx context.Context) ([]entities.Signed[liquidity_provider.TrustedAccountDetails], error) {
	signedAccounts, err := useCase.trustedAccountRepository.GetAllTrustedAccounts(ctx)
	if err != nil {
		return nil, err
	}
	validatedAccounts := make([]entities.Signed[liquidity_provider.TrustedAccountDetails], 0, len(signedAccounts))
	for i := range signedAccounts {
		readFunction := func() (*entities.Signed[liquidity_provider.TrustedAccountDetails], error) {
			return &signedAccounts[i], nil
		}
		validatedAccount, err := liquidity_provider.ValidateConfiguration(
			"trusted account",
			useCase.signer,
			readFunction,
			useCase.hashFunction,
		)
		if err != nil {
			return nil, liquidity_provider.ErrTamperedTrustedAccount
		}
		validatedAccounts = append(validatedAccounts, *validatedAccount)
	}
	return validatedAccounts, nil
}
