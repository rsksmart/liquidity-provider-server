package liquidity_provider

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type DeleteTrustedAccountUseCase struct {
	trustedAccountRepository liquidity_provider.TrustedAccountRepository
}

func NewDeleteTrustedAccountUseCase(
	trustedAccountRepository liquidity_provider.TrustedAccountRepository,
) *DeleteTrustedAccountUseCase {
	return &DeleteTrustedAccountUseCase{
		trustedAccountRepository: trustedAccountRepository,
	}
}

func (useCase *DeleteTrustedAccountUseCase) Run(ctx context.Context, address string) error {
	err := useCase.trustedAccountRepository.DeleteTrustedAccount(ctx, address)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.SetTrustedAccountId, err)
	}
	return nil
}
