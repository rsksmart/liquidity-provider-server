package pegout

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type StatusUseCase struct {
	quoteRepository quote.PegoutQuoteRepository
}

func NewStatusUseCase(quoteRepository quote.PegoutQuoteRepository) *StatusUseCase {
	return &StatusUseCase{quoteRepository: quoteRepository}
}

func (useCase *StatusUseCase) Run(ctx context.Context, quoteHash string) (quote.WatchedPegoutQuote, error) {
	pegoutQuote, err := useCase.quoteRepository.GetQuote(ctx, quoteHash)
	if err != nil {
		return quote.WatchedPegoutQuote{}, usecases.WrapUseCaseError(usecases.PegoutQuoteStatusId, err)
	} else if pegoutQuote == nil {
		return quote.WatchedPegoutQuote{}, usecases.WrapUseCaseError(usecases.PegoutQuoteStatusId, usecases.QuoteNotFoundError)
	}
	retainedQuote, err := useCase.quoteRepository.GetRetainedQuote(ctx, quoteHash)
	if err != nil {
		return quote.WatchedPegoutQuote{}, usecases.WrapUseCaseError(usecases.PegoutQuoteStatusId, err)
	} else if retainedQuote == nil {
		return quote.WatchedPegoutQuote{}, usecases.WrapUseCaseError(usecases.PegoutQuoteStatusId, usecases.QuoteNotAcceptedError)
	}

	return quote.WatchedPegoutQuote{
		PegoutQuote:   *pegoutQuote,
		RetainedQuote: *retainedQuote,
	}, nil
}
