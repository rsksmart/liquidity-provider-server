package pegin

import (
	"context"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"math/big"
	"sync"
)

type CallForUserUseCase struct {
	lbc             blockchain.LiquidityBridgeContract
	quoteRepository quote.PeginQuoteRepository
	btc             blockchain.BitcoinNetwork
	peginProvider   entities.LiquidityProvider
	eventBus        entities.EventBus
	rskWalletMutex  *sync.Mutex
}

func NewCallForUserUseCase(
	lbc blockchain.LiquidityBridgeContract,
	quoteRepository quote.PeginQuoteRepository,
	btc blockchain.BitcoinNetwork,
	peginProvider entities.LiquidityProvider,
	eventBus entities.EventBus,
	rskWalletMutex *sync.Mutex,
) *CallForUserUseCase {
	return &CallForUserUseCase{
		lbc:             lbc,
		quoteRepository: quoteRepository,
		btc:             btc,
		peginProvider:   peginProvider,
		eventBus:        eventBus,
		rskWalletMutex:  rskWalletMutex,
	}
}

func (useCase *CallForUserUseCase) Run(ctx context.Context, bitcoinTx string, retainedQuote quote.RetainedPeginQuote) error {
	balance := new(entities.Wei)
	valueToSend := new(entities.Wei)
	var txInfo blockchain.BitcoinTransactionInformation
	var peginQuote *quote.PeginQuote
	var quoteState quote.PeginState
	var callForUserTx string
	var txConfirmations big.Int
	var err error

	if retainedQuote.State != quote.PeginStateWaitingForDeposit {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PeginQuote{}, err, true)
	}

	if peginQuote, err = useCase.quoteRepository.GetQuote(ctx, retainedQuote.QuoteHash); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PeginQuote{}, err, true)
	} else if peginQuote == nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PeginQuote{}, usecases.QuoteNotFoundError, false)
	}

	if peginQuote.IsExpired() {
		return useCase.publishErrorEvent(ctx, retainedQuote, *peginQuote, usecases.ExpiredQuoteError, false)
	}

	if txInfo, err = useCase.btc.GetTransactionInfo(bitcoinTx); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, *peginQuote, err, true)
	}
	txConfirmations.SetUint64(txInfo.Confirmations)
	if txConfirmations.Cmp(big.NewInt(int64(peginQuote.Confirmations))) < 0 {
		return useCase.publishErrorEvent(ctx, retainedQuote, *peginQuote, usecases.NoEnoughConfirmationsError, true)
	}

	sentAmount := txInfo.AmountToAddress(retainedQuote.DepositAddress)
	if sentAmount.Cmp(peginQuote.Total()) < 0 {
		retainedQuote.UserBtcTxHash = bitcoinTx
		return useCase.publishErrorEvent(
			ctx,
			retainedQuote,
			*peginQuote,
			fmt.Errorf("insufficient amount %v < %v", sentAmount, peginQuote.Total()),
			false,
		)
	}

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	if balance, err = useCase.lbc.GetBalance(useCase.peginProvider.RskAddress()); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, *peginQuote, err, true)
	}

	if balance.Cmp(peginQuote.Value) < 0 { // lbc balance is not sufficient, calc delta to transfer
		valueToSend.Sub(peginQuote.Value, balance)
	}

	config := blockchain.NewTransactionConfig(valueToSend, uint64(peginQuote.GasLimit+CallForUserExtraGas), nil)
	if callForUserTx, err = useCase.lbc.CallForUser(config, *peginQuote); err != nil {
		quoteState = quote.PeginStateCallForUserFailed
	} else {
		quoteState = quote.PeginStateCallForUserSucceeded
	}

	retainedQuote.CallForUserTxHash = callForUserTx
	retainedQuote.UserBtcTxHash = bitcoinTx
	retainedQuote.State = quoteState
	useCase.eventBus.Publish(quote.CallForUserCompletedEvent{
		Event:         entities.NewBaseEvent(quote.CallForUserCompletedEventId),
		PeginQuote:    *peginQuote,
		RetainedQuote: retainedQuote,
		Error:         err,
	})

	if updateError := useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); updateError != nil {
		err = errors.Join(err, updateError)
	}
	if err != nil {
		err = errors.Join(err, usecases.NonRecoverableError)
		return usecases.WrapUseCaseErrorArgs(usecases.CallForUserId, err, usecases.ErrorArg("quoteHash", retainedQuote.QuoteHash))
	} else {
		return nil
	}
}

func (useCase *CallForUserUseCase) publishErrorEvent(
	ctx context.Context,
	retainedQuote quote.RetainedPeginQuote,
	peginQuote quote.PeginQuote,
	err error,
	recoverable bool,
) error {
	wrappedError := usecases.WrapUseCaseErrorArgs(usecases.CallForUserId, err, usecases.ErrorArg("quoteHash", retainedQuote.QuoteHash))
	if !recoverable {
		retainedQuote.State = quote.PeginStateCallForUserFailed
		if err = useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); err != nil {
			wrappedError = errors.Join(wrappedError, err, usecases.NonRecoverableError)
		}
		useCase.eventBus.Publish(quote.CallForUserCompletedEvent{
			Event:         entities.NewBaseEvent(quote.CallForUserCompletedEventId),
			RetainedQuote: retainedQuote,
			PeginQuote:    peginQuote,
			Error:         wrappedError,
		})

	}
	return wrappedError
}
