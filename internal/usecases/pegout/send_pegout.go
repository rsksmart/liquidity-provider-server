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
	btc             blockchain.BitcoinNetwork
	quoteRepository quote.PegoutQuoteRepository
	rsk             blockchain.RootstockRpcServer
	eventBus        entities.EventBus
	btcWalletMutex  *sync.Mutex
}

func NewSendPegoutUseCase(
	btcWallet blockchain.BitcoinWallet,
	btc blockchain.BitcoinNetwork,
	quoteRepository quote.PegoutQuoteRepository,
	rsk blockchain.RootstockRpcServer,
	eventBus entities.EventBus,
	btcWalletMutex *sync.Mutex,
) *SendPegoutUseCase {
	return &SendPegoutUseCase{
		btcWallet:       btcWallet,
		btc:             btc,
		quoteRepository: quoteRepository,
		rsk:             rsk,
		eventBus:        eventBus,
		btcWalletMutex:  btcWalletMutex,
	}
}

func (useCase *SendPegoutUseCase) Run(ctx context.Context, retainedQuote quote.RetainedPegoutQuote) error {
	var err error
	var pegoutQuote *quote.PegoutQuote
	var chainHeight uint64
	var receipt blockchain.TransactionReceipt
	var txHash string
	var newState quote.PegoutState

	balance := new(entities.Wei)

	if retainedQuote.State != quote.PegoutStateWaitingForDepositConfirmations {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, usecases.WrongStateError, true)
	} else if retainedQuote.UserRskTxHash == "" {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, errors.New("user rsk tx hash not provided"), true)
	}

	if pegoutQuote, err = useCase.quoteRepository.GetQuote(ctx, retainedQuote.QuoteHash); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, err, true)
	} else if pegoutQuote == nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, usecases.QuoteNotFoundError, false)
	}

	if pegoutQuote.IsExpired() {
		return useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, usecases.ExpiredQuoteError, false)
	}

	if chainHeight, err = useCase.rsk.GetHeight(ctx); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, true)
	}

	if receipt, err = useCase.rsk.GetTransactionReceipt(ctx, retainedQuote.UserRskTxHash); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, true)
	} else if chainHeight-receipt.BlockNumber < uint64(pegoutQuote.DepositConfirmations) {
		return useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, usecases.NoEnoughConfirmationsError, true)
	} else if receipt.Value.Cmp(pegoutQuote.Total()) < 0 {
		retainedQuote.UserRskTxHash = receipt.TransactionHash
		return useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, usecases.InsufficientAmountError, false)
	}

	useCase.btcWalletMutex.Lock()
	defer useCase.btcWalletMutex.Unlock()

	if balance, err = useCase.btcWallet.GetBalance(); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, true)
	} else if balance.Cmp(pegoutQuote.Value) < 0 {
		return useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, usecases.NoLiquidityError, true)
	}

	quoteHashBytes, err := hex.DecodeString(retainedQuote.QuoteHash)
	if err != nil {
		retainedQuote.UserRskTxHash = receipt.TransactionHash
		return useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, false)
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

	if updateError := useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); updateError != nil {
		err = errors.Join(err, updateError)
	}
	if err != nil {
		return usecases.WrapUseCaseErrorArgs(usecases.SendPegoutId, err, usecases.ErrorArg("quoteHash", retainedQuote.QuoteHash))
	} else {
		return nil
	}
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
		if err = useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); err != nil {
			wrappedError = errors.Join(wrappedError, err)
		}
		retainedQuote.State = quote.PegoutStateSendPegoutFailed
		useCase.eventBus.Publish(quote.PegoutBtcSentToUserEvent{
			Event:         entities.NewBaseEvent(quote.PegoutBtcSentEventId),
			PegoutQuote:   pegoutQuote,
			RetainedQuote: retainedQuote,
			Error:         wrappedError,
		})
	}
	return wrappedError
}
