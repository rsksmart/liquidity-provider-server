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
	normalized, err := usecases.NormalizeTrustedAccountAddress(usecases.DeleteTrustedAccountId, address)
	if err != nil {
		return err
	}
	err = useCase.trustedAccountRepository.DeleteTrustedAccount(ctx, normalized)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.DeleteTrustedAccountId, err)
	}
	return nil
}
