package liquidity_provider

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
)

type GetTrustedAccountUseCase struct {
	trustedAccountRepository liquidity_provider.TrustedAccountRepository
	hashFunction             entities.HashFunction
}

func NewGetTrustedAccountUseCase(
	trustedAccountRepository liquidity_provider.TrustedAccountRepository,
	hashFunction entities.HashFunction,
) *GetTrustedAccountUseCase {
	return &GetTrustedAccountUseCase{
		trustedAccountRepository: trustedAccountRepository,
		hashFunction:             hashFunction,
	}
}

func (useCase *GetTrustedAccountUseCase) Run(ctx context.Context, address string) (*entities.Signed[liquidity_provider.TrustedAccountDetails], error) {
	signedAccount, err := useCase.trustedAccountRepository.GetTrustedAccount(ctx, address)
	if err != nil {
		return nil, err
	}
	if signedAccount == nil {
		return nil, liquidity_provider.ErrTrustedAccountNotFound
	}
	if err := signedAccount.CheckIntegrity(useCase.hashFunction); err != nil {
		return nil, liquidity_provider.ErrTamperedTrustedAccount
	}
	return signedAccount, nil
}
