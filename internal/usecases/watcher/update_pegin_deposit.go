package watcher

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type UpdatePeginDepositUseCase struct {
	peginRepository quote.PeginQuoteRepository
}

func NewUpdatePeginDepositUseCase(peginRepository quote.PeginQuoteRepository) *UpdatePeginDepositUseCase {
	return &UpdatePeginDepositUseCase{peginRepository: peginRepository}
}

func (useCase *UpdatePeginDepositUseCase) Run(
	ctx context.Context,
	watchedQuote quote.WatchedPeginQuote,
	block blockchain.BitcoinBlockInformation,
	deposit blockchain.BitcoinTransactionInformation,
) (quote.WatchedPeginQuote, error) {
	sentAmount := deposit.AmountToAddress(watchedQuote.RetainedQuote.DepositAddress)
	enoughAmount := sentAmount.Cmp(watchedQuote.PeginQuote.Total()) >= 0
	sentBeforeExpire := block.Time.Before(watchedQuote.PeginQuote.ExpireTime())

	if !enoughAmount || !sentBeforeExpire {
		return quote.WatchedPeginQuote{}, usecases.WrapUseCaseError(usecases.UpdatePeginDepositId, errors.New("invalid bitcoin transaction for quote"))
	} else if watchedQuote.RetainedQuote.State != quote.PeginStateWaitingForDeposit {
		return quote.WatchedPeginQuote{}, usecases.WrapUseCaseError(usecases.UpdatePeginDepositId, usecases.IllegalQuoteStateError)
	}

	watchedQuote.RetainedQuote.State = quote.PeginStateWaitingForDepositConfirmations
	watchedQuote.RetainedQuote.UserBtcTxHash = deposit.Hash

	if err := useCase.peginRepository.UpdateRetainedQuote(ctx, watchedQuote.RetainedQuote); err != nil {
		return quote.WatchedPeginQuote{}, usecases.WrapUseCaseError(usecases.UpdatePeginDepositId, err)
	}
	return watchedQuote, nil
}
