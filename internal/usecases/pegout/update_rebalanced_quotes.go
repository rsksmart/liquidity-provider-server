package pegout

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type UpdateBtcReleaseUseCase struct {
	pegoutRepository quote.PegoutQuoteRepository
	batchRepository  rootstock.BatchPegOutRepository
	eventBus         entities.EventBus
}

func NewUpdateBtcReleaseUseCase(
	pegoutRepository quote.PegoutQuoteRepository,
	batchRepository rootstock.BatchPegOutRepository,
	eventBus entities.EventBus,
) *UpdateBtcReleaseUseCase {
	return &UpdateBtcReleaseUseCase{
		pegoutRepository: pegoutRepository,
		batchRepository:  batchRepository,
		eventBus:         eventBus,
	}
}

func (useCase *UpdateBtcReleaseUseCase) Run(ctx context.Context, batch rootstock.BatchPegOut) (uint, error) {
	retainedQuotes, err := useCase.pegoutRepository.GetRetainedQuotesInBatch(ctx, batch)
	totalQuotes := uint(len(retainedQuotes))
	if err != nil {
		return 0, usecases.WrapUseCaseError(usecases.UpdateBtcReleaseId, err)
	} else if totalQuotes == 0 {
		return 0, nil
	}

	quoteHashes := make([]string, totalQuotes)
	for i := range retainedQuotes {
		quoteHashes[i] = retainedQuotes[i].QuoteHash
		retainedQuotes[i].BtcReleaseTxHash = batch.TransactionHash
		retainedQuotes[i].State = quote.PegoutStateBtcReleased
	}

	err = useCase.pegoutRepository.UpdateRetainedQuotes(ctx, retainedQuotes)
	if err != nil {
		return 0, usecases.WrapUseCaseError(usecases.UpdateBtcReleaseId, err)
	}

	err = useCase.batchRepository.UpsertBatch(ctx, batch)
	if err != nil {
		return totalQuotes, usecases.WrapUseCaseError(usecases.UpdateBtcReleaseId, err)
	}

	useCase.eventBus.Publish(rootstock.BatchPegOutUpdatedEvent{
		Event:       entities.NewBaseEvent(rootstock.BatchPegOutUpdatedEventId),
		QuoteHashes: quoteHashes,
		BatchPegOut: batch,
	})

	return totalQuotes, nil
}
