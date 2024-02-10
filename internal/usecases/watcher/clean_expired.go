package watcher

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
)

type CleanExpiredQuotesUseCase struct {
	peginRepository  quote.PeginQuoteRepository
	pegoutRepository quote.PegoutQuoteRepository
}

func NewCleanExpiredQuotesUseCase(peginRepository quote.PeginQuoteRepository, pegoutRepository quote.PegoutQuoteRepository) *CleanExpiredQuotesUseCase {
	return &CleanExpiredQuotesUseCase{peginRepository: peginRepository, pegoutRepository: pegoutRepository}
}

func (useCase *CleanExpiredQuotesUseCase) Run(ctx context.Context) ([]string, error) {
	var peginQuotes []quote.RetainedPeginQuote
	var pegoutQuotes []quote.RetainedPegoutQuote
	var err error

	peginHashes := make([]string, 0)
	pegoutHashes := make([]string, 0)

	if peginQuotes, err = useCase.peginRepository.GetRetainedQuoteByState(ctx, quote.PeginStateTimeForDepositElapsed); err != nil {
		return nil, err
	}
	for _, peginQuote := range peginQuotes {
		peginHashes = append(peginHashes, peginQuote.QuoteHash)
	}

	if pegoutQuotes, err = useCase.pegoutRepository.GetRetainedQuoteByState(ctx, quote.PegoutStateTimeForDepositElapsed); err != nil {
		return nil, err
	}
	for _, pegoutQuote := range pegoutQuotes {
		pegoutHashes = append(pegoutHashes, pegoutQuote.QuoteHash)
	}

	if _, err = useCase.peginRepository.DeleteQuotes(ctx, peginHashes); err != nil {
		return nil, err
	}
	if _, err = useCase.pegoutRepository.DeleteQuotes(ctx, pegoutHashes); err != nil {
		return nil, err
	}

	peginHashes = append(peginHashes, pegoutHashes...)
	return peginHashes, nil
}
