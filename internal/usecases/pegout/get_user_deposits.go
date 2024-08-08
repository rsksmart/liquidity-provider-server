package pegout

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetUserDepositsUseCase struct {
	quoteRepository quote.PegoutQuoteRepository
}

func NewGetUserDepositsUseCase(quoteRepository quote.PegoutQuoteRepository) *GetUserDepositsUseCase {
	return &GetUserDepositsUseCase{quoteRepository: quoteRepository}
}

func (useCase *GetUserDepositsUseCase) Run(ctx context.Context, address string) ([]quote.PegoutDeposit, error) {
	var err error
	var deposits []quote.PegoutDeposit
	if deposits, err = useCase.quoteRepository.ListPegoutDepositsByAddress(ctx, address); err != nil {
		return deposits, usecases.WrapUseCaseError(usecases.GetUserQuotesId, err)
	}
	return deposits, nil
}
