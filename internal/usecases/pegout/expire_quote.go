package pegout

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type ExpiredPegoutQuoteUseCase struct {
	pegoutRepository quote.PegoutQuoteRepository
}

func NewExpiredPegoutQuoteUseCase(pegoutRepository quote.PegoutQuoteRepository) *ExpiredPegoutQuoteUseCase {
	return &ExpiredPegoutQuoteUseCase{pegoutRepository: pegoutRepository}
}

func (useCase *ExpiredPegoutQuoteUseCase) Run(ctx context.Context, pegoutQuote quote.RetainedPegoutQuote) error {
	pegoutQuote.State = quote.PegoutStateTimeForDepositElapsed
	err := useCase.pegoutRepository.UpdateRetainedQuote(ctx, pegoutQuote)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ExpiredPegoutQuoteId, err)
	}
	return nil
}
