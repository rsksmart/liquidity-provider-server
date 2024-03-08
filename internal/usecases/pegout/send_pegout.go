package pegout

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"sync"
)

type SendPegoutUseCase struct {
	btcWallet       blockchain.BitcoinWallet
	quoteRepository quote.PegoutQuoteRepository
	rpc             blockchain.Rpc
	eventBus        entities.EventBus
	btcWalletMutex  sync.Locker
}

func NewSendPegoutUseCase(
	btcWallet blockchain.BitcoinWallet,
	quoteRepository quote.PegoutQuoteRepository,
	rpc blockchain.Rpc,
	eventBus entities.EventBus,
	btcWalletMutex sync.Locker,
) *SendPegoutUseCase {
	return &SendPegoutUseCase{
		btcWallet:       btcWallet,
		quoteRepository: quoteRepository,
		rpc:             rpc,
		eventBus:        eventBus,
		btcWalletMutex:  btcWalletMutex,
	}
}

func (useCase *SendPegoutUseCase) Run(ctx context.Context, retainedQuote quote.RetainedPegoutQuote) error {
	var err error
	var pegoutQuote *quote.PegoutQuote
	var receipt blockchain.TransactionReceipt

	if err = useCase.validateRetainedQuote(ctx, retainedQuote); err != nil {
		return err
	}

	if pegoutQuote, err = useCase.quoteRepository.GetQuote(ctx, retainedQuote.QuoteHash); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, err, true)
	} else if pegoutQuote == nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, usecases.QuoteNotFoundError, false)
	}

	if receipt, err = useCase.validateQuote(ctx, retainedQuote, pegoutQuote); err != nil {
		return err
	}

	useCase.btcWalletMutex.Lock()
	defer useCase.btcWalletMutex.Unlock()

	if err = useCase.validateBalance(ctx, retainedQuote, pegoutQuote); err != nil {
		return err
	}

	retainedQuote, err = useCase.performSendPegout(ctx, retainedQuote, pegoutQuote, receipt)
	// if the error is not nil and the state is not SendPegoutFailed,
	// means that an error happened before sending the tx
	if err != nil && retainedQuote.State != quote.PegoutStateSendPegoutFailed {
		return err
	}

	if updateError := useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); updateError != nil {
		err = errors.Join(err, updateError)
	}
	if err != nil {
		return usecases.WrapUseCaseErrorArgs(usecases.SendPegoutId, err, usecases.ErrorArg("quoteHash", retainedQuote.QuoteHash))
	}
	return nil
}

func (useCase *SendPegoutUseCase) publishErrorEvent(
	ctx context.Context,
	retainedQuote quote.RetainedPegoutQuote,
	pegoutQuote quote.PegoutQuote,
	err error,
	recoverable bool,
) error {
	wrappedError := usecases.WrapUseCaseErrorArgs(usecases.SendPegoutId, err, usecases.ErrorArg("quoteHash", retainedQuote.QuoteHash))
	if !recoverable {
		retainedQuote.State = quote.PegoutStateSendPegoutFailed
		if err = useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); err != nil {
			wrappedError = errors.Join(wrappedError, err)
		}
		useCase.eventBus.Publish(quote.PegoutBtcSentToUserEvent{
			Event:         entities.NewBaseEvent(quote.PegoutBtcSentEventId),
			PegoutQuote:   pegoutQuote,
			RetainedQuote: retainedQuote,
			Error:         wrappedError,
		})
	}
	return wrappedError
}

func (useCase *SendPegoutUseCase) validateQuote(
	ctx context.Context,
	retainedQuote quote.RetainedPegoutQuote,
	pegoutQuote *quote.PegoutQuote,
) (blockchain.TransactionReceipt, error) {
	var err error
	var chainHeight uint64
	var receipt blockchain.TransactionReceipt

	if pegoutQuote.IsExpired() {
		return blockchain.TransactionReceipt{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, usecases.ExpiredQuoteError, false)
	}

	if chainHeight, err = useCase.rpc.Rsk.GetHeight(ctx); err != nil {
		return blockchain.TransactionReceipt{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, true)
	}

	if receipt, err = useCase.rpc.Rsk.GetTransactionReceipt(ctx, retainedQuote.UserRskTxHash); err != nil {
		return blockchain.TransactionReceipt{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, true)
	} else if chainHeight-receipt.BlockNumber < uint64(pegoutQuote.DepositConfirmations) {
		return blockchain.TransactionReceipt{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, usecases.NoEnoughConfirmationsError, true)
	} else if receipt.Value.Cmp(pegoutQuote.Total()) < 0 {
		retainedQuote.UserRskTxHash = receipt.TransactionHash
		return blockchain.TransactionReceipt{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, usecases.InsufficientAmountError, false)
	}
	return receipt, nil
}

func (useCase *SendPegoutUseCase) performSendPegout(
	ctx context.Context,
	retainedQuote quote.RetainedPegoutQuote,
	pegoutQuote *quote.PegoutQuote,
	receipt blockchain.TransactionReceipt,
) (quote.RetainedPegoutQuote, error) {
	var err error
	var newState quote.PegoutState
	var txHash string

	quoteHashBytes, err := hex.DecodeString(retainedQuote.QuoteHash)
	if err != nil {
		retainedQuote.UserRskTxHash = receipt.TransactionHash
		return quote.RetainedPegoutQuote{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, false)
	}

	if txHash, err = useCase.btcWallet.SendWithOpReturn(pegoutQuote.DepositAddress, pegoutQuote.Value, quoteHashBytes); err != nil {
		newState = quote.PegoutStateSendPegoutFailed
	} else {
		newState = quote.PegoutStateSendPegoutSucceeded
	}

	retainedQuote.LpBtcTxHash = txHash
	retainedQuote.State = newState
	useCase.eventBus.Publish(quote.PegoutBtcSentToUserEvent{
		Event:         entities.NewBaseEvent(quote.PegoutBtcSentEventId),
		PegoutQuote:   *pegoutQuote,
		RetainedQuote: retainedQuote,
		Error:         err,
	})
	return retainedQuote, err
}

func (useCase *SendPegoutUseCase) validateBalance(
	ctx context.Context,
	retainedQuote quote.RetainedPegoutQuote,
	pegoutQuote *quote.PegoutQuote,
) error {
	var err error
	var balance *entities.Wei

	requiredBalance := new(entities.Wei)
	if balance, err = useCase.btcWallet.GetBalance(); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, true)
	}
	requiredBalance = new(entities.Wei)
	requiredBalance.Add(pegoutQuote.Value, pegoutQuote.GasFee)
	if balance.Cmp(requiredBalance) < 0 {
		return useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, usecases.NoLiquidityError, true)
	}
	return nil
}

func (useCase *SendPegoutUseCase) validateRetainedQuote(ctx context.Context, retainedQuote quote.RetainedPegoutQuote) error {
	if retainedQuote.State != quote.PegoutStateWaitingForDepositConfirmations {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, usecases.WrongStateError, true)
	} else if retainedQuote.UserRskTxHash == "" {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, errors.New("user rsk tx hash not provided"), true)
	}
	return nil
}
