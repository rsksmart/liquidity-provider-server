package watcher

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type UpdatePegoutQuoteDepositUseCase struct {
	pegoutRepository quote.PegoutQuoteRepository
}

func NewUpdatePegoutQuoteDepositUseCase(pegoutRepository quote.PegoutQuoteRepository) *UpdatePegoutQuoteDepositUseCase {
	return &UpdatePegoutQuoteDepositUseCase{pegoutRepository: pegoutRepository}
}

func (useCase *UpdatePegoutQuoteDepositUseCase) Run(ctx context.Context, watchedQuote WatchedPegoutQuote, deposit quote.PegoutDeposit) (WatchedPegoutQuote, error) {
	var err error
	if !deposit.IsValidForQuote(watchedQuote.PegoutQuote) {
		return WatchedPegoutQuote{}, usecases.WrapUseCaseError(usecases.UpdatePegoutDepositId, errors.New("deposit not valid for quote"))
	} else if watchedQuote.RetainedQuote.State != quote.PegoutStateWaitingForDeposit {
		return WatchedPegoutQuote{}, usecases.WrapUseCaseError(usecases.UpdatePegoutDepositId, errors.New("illegal quote state"))
	}
	watchedQuote.RetainedQuote.State = quote.PegoutStateWaitingForDepositConfirmations
	watchedQuote.RetainedQuote.UserRskTxHash = deposit.TxHash
	if err = useCase.pegoutRepository.UpdateRetainedQuote(ctx, watchedQuote.RetainedQuote); err != nil {
		return WatchedPegoutQuote{}, usecases.WrapUseCaseError(usecases.UpdatePegoutDepositId, err)
	}
	if err = useCase.pegoutRepository.UpsertPegoutDeposit(ctx, deposit); err != nil {
		return WatchedPegoutQuote{}, usecases.WrapUseCaseError(usecases.UpdatePegoutDepositId, err)
	}
	return watchedQuote, nil
}
