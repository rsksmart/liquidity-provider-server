package pegin

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type StatusUseCase struct {
	quoteRepository quote.PeginQuoteRepository
}

func NewStatusUseCase(quoteRepository quote.PeginQuoteRepository) *StatusUseCase {
	return &StatusUseCase{quoteRepository: quoteRepository}
}

func (useCase *StatusUseCase) Run(ctx context.Context, quoteHash string) (quote.WatchedPeginQuote, error) {
	logger := log.WithField("quote_hash", usecases.SafeLogStr(quoteHash))
	peginQuote, err := useCase.quoteRepository.GetQuote(ctx, quoteHash)
	if err != nil {
		logger.WithError(err).Error("QuoteStatus: repository load failed")
		return quote.WatchedPeginQuote{}, usecases.WrapUseCaseError(usecases.PeginQuoteStatusId, err)
	} else if peginQuote == nil {
		logger.Debug("QuoteStatus: quote not found")
		return quote.WatchedPeginQuote{}, usecases.WrapUseCaseError(usecases.PeginQuoteStatusId, usecases.QuoteNotFoundError)
	}
	retainedQuote, err := useCase.quoteRepository.GetRetainedQuote(ctx, quoteHash)
	if err != nil {
		logger.WithError(err).Error("QuoteStatus: repository load failed")
		return quote.WatchedPeginQuote{}, usecases.WrapUseCaseError(usecases.PeginQuoteStatusId, err)
	} else if retainedQuote == nil {
		logger.Debug("QuoteStatus: quote not yet accepted")
		return quote.WatchedPeginQuote{}, usecases.WrapUseCaseError(usecases.PeginQuoteStatusId, usecases.QuoteNotAcceptedError)
	}
	creationData := useCase.quoteRepository.GetPeginCreationData(ctx, quoteHash)

	watchedQuote := quote.NewWatchedPeginQuote(*peginQuote, *retainedQuote, creationData)
	return watchedQuote, nil
}
