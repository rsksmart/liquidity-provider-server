package liquidity_provider

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
)

type GetTrustedAccountsUseCase struct {
	trustedAccountRepository liquidity_provider.TrustedAccountRepository
}

func NewGetTrustedAccountsUseCase(
	trustedAccountRepository liquidity_provider.TrustedAccountRepository,
) *GetTrustedAccountsUseCase {
	return &GetTrustedAccountsUseCase{
		trustedAccountRepository: trustedAccountRepository,
	}
}

func (useCase *GetTrustedAccountsUseCase) Run(ctx context.Context) ([]liquidity_provider.TrustedAccountDetails, error) {
	return useCase.trustedAccountRepository.GetAllTrustedAccounts(ctx)
}
