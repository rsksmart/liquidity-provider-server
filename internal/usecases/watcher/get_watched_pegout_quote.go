package watcher

import (
	"context"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetWatchedPegoutQuoteUseCase struct {
	pegoutRepository quote.PegoutQuoteRepository
}

func NewGetWatchedPegoutQuoteUseCase(pegoutRepository quote.PegoutQuoteRepository) *GetWatchedPegoutQuoteUseCase {
	return &GetWatchedPegoutQuoteUseCase{pegoutRepository: pegoutRepository}
}

func (useCase *GetWatchedPegoutQuoteUseCase) Run(ctx context.Context, states ...quote.PegoutState) ([]quote.WatchedPegoutQuote, error) {
	result := make([]quote.WatchedPegoutQuote, 0)
	for _, state := range states {
		switch state {
		case
			quote.PegoutStateWaitingForDepositConfirmations,
			quote.PegoutStateWaitingForDeposit,
			quote.PegoutStateSendPegoutSucceeded,
			quote.PegoutStateRefundPegOutSucceeded:
			if watchedQuotes, err := useCase.getWatchedQuotes(ctx, state); err == nil {
				result = append(result, watchedQuotes...)
			} else {
				return nil, usecases.WrapUseCaseError(usecases.GetWatchedPegoutQuoteId, err)
			}
		default:
			return nil, usecases.WrapUseCaseError(usecases.GetWatchedPegoutQuoteId, fmt.Errorf("illegal state %s", state))
		}
	}
	return result, nil
}

func (useCase *GetWatchedPegoutQuoteUseCase) getWatchedQuotes(ctx context.Context, state quote.PegoutState) ([]quote.WatchedPegoutQuote, error) {
	var retainedQuotes []quote.RetainedPegoutQuote
	watchedQuote := make([]quote.WatchedPegoutQuote, 0)
	var pegoutQuote *quote.PegoutQuote
	var err error
	if retainedQuotes, err = useCase.pegoutRepository.GetRetainedQuoteByState(ctx, state); err != nil {
		return nil, usecases.WrapUseCaseError(usecases.GetWatchedPegoutQuoteId, err)
	}

	for _, retainedQuote := range retainedQuotes {
		if pegoutQuote, err = useCase.pegoutRepository.GetQuote(ctx, retainedQuote.QuoteHash); err != nil {
			return nil, usecases.WrapUseCaseError(usecases.GetWatchedPegoutQuoteId, err)
		}
		watchedQuote = append(
			watchedQuote,
			quote.NewWatchedPegoutQuote(*pegoutQuote, retainedQuote),
		)
	}
	return watchedQuote, nil
}
