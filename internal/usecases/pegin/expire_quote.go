package pegin

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type ExpiredPeginQuoteUseCase struct {
	peginRepository quote.PeginQuoteRepository
}

func NewExpiredPeginQuoteUseCase(peginRepository quote.PeginQuoteRepository) *ExpiredPeginQuoteUseCase {
	return &ExpiredPeginQuoteUseCase{peginRepository: peginRepository}
}

func (useCase *ExpiredPeginQuoteUseCase) Run(ctx context.Context, peginQuote quote.RetainedPeginQuote) error {
	peginQuote.State = quote.PeginStateTimeForDepositElapsed
	err := useCase.peginRepository.UpdateRetainedQuote(ctx, peginQuote)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ExpiredPeginQuoteId, err)
	}
	return nil
}
