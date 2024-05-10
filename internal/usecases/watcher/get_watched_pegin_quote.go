package watcher

import (
	"context"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetWatchedPeginQuoteUseCase struct {
	peginRepository quote.PeginQuoteRepository
}

func NewGetWatchedPeginQuoteUseCase(peginRepository quote.PeginQuoteRepository) *GetWatchedPeginQuoteUseCase {
	return &GetWatchedPeginQuoteUseCase{peginRepository: peginRepository}
}

func (useCase *GetWatchedPeginQuoteUseCase) Run(ctx context.Context, state quote.PeginState) ([]quote.WatchedPeginQuote, error) {
	switch state {
	case quote.PeginStateWaitingForDeposit, quote.PeginStateCallForUserSucceeded:
		return useCase.getWatchedQuotes(ctx, state)
	default:
		return nil, fmt.Errorf("GetWatchedPeginQuoteUseCase: illegal state %s", state)
	}
}

func (useCase *GetWatchedPeginQuoteUseCase) getWatchedQuotes(ctx context.Context, state quote.PeginState) ([]quote.WatchedPeginQuote, error) {
	var retainedQuotes []quote.RetainedPeginQuote
	watchedQuotes := make([]quote.WatchedPeginQuote, 0)
	var peginQuote *quote.PeginQuote
	var err error
	if retainedQuotes, err = useCase.peginRepository.GetRetainedQuoteByState(ctx, state); err != nil {
		return nil, usecases.WrapUseCaseError(usecases.GetWatchedPeginQuoteId, err)
	}

	for _, retainedQuote := range retainedQuotes {
		if peginQuote, err = useCase.peginRepository.GetQuote(ctx, retainedQuote.QuoteHash); err != nil {
			return nil, usecases.WrapUseCaseError(usecases.GetWatchedPeginQuoteId, err)
		}
		watchedQuotes = append(
			watchedQuotes,
			quote.NewWatchedPeginQuote(*peginQuote, retainedQuote),
		)
	}
	return watchedQuotes, nil
}
