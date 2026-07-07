package pegout

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type StatusUseCase struct {
	quoteRepository quote.PegoutQuoteRepository
}

func NewStatusUseCase(quoteRepository quote.PegoutQuoteRepository) *StatusUseCase {
	return &StatusUseCase{quoteRepository: quoteRepository}
}

func (useCase *StatusUseCase) Run(ctx context.Context, quoteHash string) (quote.WatchedPegoutQuote, error) {
	logger := log.WithField("quote_hash", usecases.SafeLogStr(quoteHash))
	pegoutQuote, err := useCase.quoteRepository.GetQuote(ctx, quoteHash)
	if err != nil {
		logger.WithError(err).Error("QuoteStatus: repository load failed")
		return quote.WatchedPegoutQuote{}, usecases.WrapUseCaseError(usecases.PegoutQuoteStatusId, err)
	} else if pegoutQuote == nil {
		logger.Debug("QuoteStatus: quote not found")
		return quote.WatchedPegoutQuote{}, usecases.WrapUseCaseError(usecases.PegoutQuoteStatusId, usecases.QuoteNotFoundError)
	}
	retainedQuote, err := useCase.quoteRepository.GetRetainedQuote(ctx, quoteHash)
	if err != nil {
		logger.WithError(err).Error("QuoteStatus: repository load failed")
		return quote.WatchedPegoutQuote{}, usecases.WrapUseCaseError(usecases.PegoutQuoteStatusId, err)
	} else if retainedQuote == nil {
		logger.Debug("QuoteStatus: quote not yet accepted")
		return quote.WatchedPegoutQuote{}, usecases.WrapUseCaseError(usecases.PegoutQuoteStatusId, usecases.QuoteNotAcceptedError)
	}
	creationData := useCase.quoteRepository.GetPegoutCreationData(ctx, quoteHash)

	watchedQuote := quote.NewWatchedPegoutQuote(*pegoutQuote, *retainedQuote, creationData)
	return watchedQuote, nil
}
