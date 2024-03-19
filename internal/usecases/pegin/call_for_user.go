package pegin

import (
	"context"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"math/big"
	"sync"
)

type CallForUserUseCase struct {
	contracts       blockchain.RskContracts
	quoteRepository quote.PeginQuoteRepository
	rpc             blockchain.Rpc
	peginProvider   liquidity_provider.LiquidityProvider
	eventBus        entities.EventBus
	rskWalletMutex  sync.Locker
}

func NewCallForUserUseCase(
	contracts blockchain.RskContracts,
	quoteRepository quote.PeginQuoteRepository,
	rpc blockchain.Rpc,
	peginProvider liquidity_provider.LiquidityProvider,
	eventBus entities.EventBus,
	rskWalletMutex sync.Locker,
) *CallForUserUseCase {
	return &CallForUserUseCase{
		contracts:       contracts,
		quoteRepository: quoteRepository,
		rpc:             rpc,
		peginProvider:   peginProvider,
		eventBus:        eventBus,
		rskWalletMutex:  rskWalletMutex,
	}
}

func (useCase *CallForUserUseCase) Run(ctx context.Context, bitcoinTx string, retainedQuote quote.RetainedPeginQuote) error {
	var valueToSend *entities.Wei
	var peginQuote *quote.PeginQuote
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

	if err = useCase.validateBitcoinTx(ctx, bitcoinTx, peginQuote, retainedQuote); err != nil {
		return err
	}

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	if valueToSend, err = useCase.calculateValueToSend(ctx, *peginQuote, retainedQuote); err != nil {
		return err
	}

	retainedQuote, err = useCase.performCallForUser(bitcoinTx, valueToSend, peginQuote, retainedQuote)

	if updateError := useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); updateError != nil {
		err = errors.Join(err, updateError)
	}
	if err != nil {
		err = errors.Join(err, usecases.NonRecoverableError)
		return usecases.WrapUseCaseErrorArgs(usecases.CallForUserId, err, usecases.ErrorArg("quoteHash", retainedQuote.QuoteHash))
	}
	return nil
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

func (useCase *CallForUserUseCase) calculateValueToSend(
	ctx context.Context,
	peginQuote quote.PeginQuote,
	retainedQuote quote.RetainedPeginQuote,
) (*entities.Wei, error) {
	var contractBalance, networkBalance *entities.Wei
	var err error

	if contractBalance, err = useCase.contracts.Lbc.GetBalance(useCase.peginProvider.RskAddress()); err != nil {
		return nil, useCase.publishErrorEvent(ctx, retainedQuote, peginQuote, err, true)
	}

	valueToSend := entities.NewWei(0)
	if contractBalance.Cmp(peginQuote.Value) < 0 { // lbc balance is not sufficient, calc delta to transfer
		valueToSend.Sub(peginQuote.Value, contractBalance)
	} else {
		return valueToSend, nil
	}

	if networkBalance, err = useCase.rpc.Rsk.GetBalance(ctx, useCase.peginProvider.RskAddress()); err != nil {
		return nil, useCase.publishErrorEvent(ctx, retainedQuote, peginQuote, err, true)
	} else if networkBalance.Cmp(valueToSend) < 0 {
		return nil, useCase.publishErrorEvent(ctx, retainedQuote, peginQuote, usecases.NoLiquidityError, true)
	}
	return valueToSend, nil
}

func (useCase *CallForUserUseCase) performCallForUser(
	bitcoinTx string,
	valueToSend *entities.Wei,
	peginQuote *quote.PeginQuote,
	retainedQuote quote.RetainedPeginQuote,
) (quote.RetainedPeginQuote, error) {
	var quoteState quote.PeginState
	var callForUserTx string
	var err error

	config := blockchain.NewTransactionConfig(valueToSend, uint64(peginQuote.GasLimit+CallForUserExtraGas), nil)
	if callForUserTx, err = useCase.contracts.Lbc.CallForUser(config, *peginQuote); err != nil {
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
	return retainedQuote, err
}

func (useCase *CallForUserUseCase) validateBitcoinTx(
	ctx context.Context,
	bitcoinTx string,
	peginQuote *quote.PeginQuote,
	retainedQuote quote.RetainedPeginQuote,
) error {
	var txInfo blockchain.BitcoinTransactionInformation
	var txConfirmations big.Int
	var err error

	if txInfo, err = useCase.rpc.Btc.GetTransactionInfo(bitcoinTx); err != nil {
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
			fmt.Errorf("%w: %v < %v", usecases.InsufficientAmountError, sentAmount, peginQuote.Total()),
			false,
		)
	}
	return nil
}
