package liquidity_provider

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
)

type GetTrustedAccountsUseCase struct {
	trustedAccountRepository liquidity_provider.TrustedAccountRepository
	hashFunction             entities.HashFunction
}

func NewGetTrustedAccountsUseCase(
	trustedAccountRepository liquidity_provider.TrustedAccountRepository,
	hashFunction entities.HashFunction,
) *GetTrustedAccountsUseCase {
	return &GetTrustedAccountsUseCase{
		trustedAccountRepository: trustedAccountRepository,
		hashFunction:             hashFunction,
	}
}

func (useCase *GetTrustedAccountsUseCase) Run(ctx context.Context) ([]entities.Signed[liquidity_provider.TrustedAccountDetails], error) {
	signedAccounts, err := useCase.trustedAccountRepository.GetAllTrustedAccounts(ctx)
	if err != nil {
		return nil, err
	}
	for _, account := range signedAccounts {
		if err := account.CheckIntegrity(useCase.hashFunction); err != nil {
			return nil, liquidity_provider.ErrTamperedTrustedAccount
		}
	}
	return signedAccounts, nil
}
