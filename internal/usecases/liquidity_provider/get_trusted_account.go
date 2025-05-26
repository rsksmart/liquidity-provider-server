package liquidity_provider

import (
	"context"
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
)

type GetTrustedAccountUseCase struct {
	trustedAccountRepository liquidity_provider.TrustedAccountRepository
	hashFunction             entities.HashFunction
	signer                   entities.Signer
}

func NewGetTrustedAccountUseCase(
	trustedAccountRepository liquidity_provider.TrustedAccountRepository,
	hashFunction entities.HashFunction,
	signer entities.Signer,
) *GetTrustedAccountUseCase {
	return &GetTrustedAccountUseCase{
		trustedAccountRepository: trustedAccountRepository,
		hashFunction:             hashFunction,
		signer:                   signer,
	}
}

func (useCase *GetTrustedAccountUseCase) Run(ctx context.Context, address string) (*entities.Signed[liquidity_provider.TrustedAccountDetails], error) {
	readFunction := func() (*entities.Signed[liquidity_provider.TrustedAccountDetails], error) {
		signedAccount, err := useCase.trustedAccountRepository.GetTrustedAccount(ctx, address)
		if err != nil {
			return nil, err
		}
		if signedAccount == nil {
			return nil, liquidity_provider.TrustedAccountNotFoundError
		}
		return signedAccount, nil
	}
	signedAccount, err := liquidity_provider.ValidateConfiguration(
		useCase.signer,
		readFunction,
		useCase.hashFunction,
	)
	if errors.Is(err, liquidity_provider.TrustedAccountNotFoundError) {
		return nil, err
	}
	if errors.Is(err, liquidity_provider.ConfigurationNotFoundError) {
		return nil, liquidity_provider.TrustedAccountNotFoundError
	}
	if err != nil {
		return nil, liquidity_provider.TamperedTrustedAccountError
	}
	return signedAccount, nil
}
