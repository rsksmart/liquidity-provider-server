package pegout

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type DepositParser = func(receipt blockchain.TransactionReceipt) (blockchain.ParsedLog[quote.PegoutDeposit], error)

type SendPegoutUseCase struct {
	btcWallet       blockchain.BitcoinWallet
	quoteRepository quote.PegoutQuoteRepository
	rpc             blockchain.Rpc
	eventBus        entities.EventBus
	contracts       blockchain.RskContracts
	btcWalletMutex  sync.Locker
	depositParser   DepositParser
}

func NewSendPegoutUseCase(
	btcWallet blockchain.BitcoinWallet,
	quoteRepository quote.PegoutQuoteRepository,
	rpc blockchain.Rpc,
	eventBus entities.EventBus,
	contracts blockchain.RskContracts,
	btcWalletMutex sync.Locker,
	depositParser DepositParser,
) *SendPegoutUseCase {
	return &SendPegoutUseCase{
		btcWallet:       btcWallet,
		quoteRepository: quoteRepository,
		rpc:             rpc,
		eventBus:        eventBus,
		contracts:       contracts,
		btcWalletMutex:  btcWalletMutex,
		depositParser:   depositParser,
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
		wrappedError = errors.Join(wrappedError, usecases.NonRecoverableError)
		if err = useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); err != nil {
			wrappedError = errors.Join(wrappedError, err)
		}
		useCase.eventBus.Publish(quote.PegoutBtcSentToUserEvent{
			Event:         entities.NewBaseEvent(quote.PegoutBtcSentEventId),
			PegoutQuote:   pegoutQuote,
			RetainedQuote: retainedQuote,
			Error:         wrappedError,
			CreationData:  quote.PegoutCreationDataZeroValue(),
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
	var block blockchain.BlockInfo
	var completed bool

	if chainHeight, err = useCase.rpc.Rsk.GetHeight(ctx); err != nil {
		return blockchain.TransactionReceipt{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, true)
	}

	if receipt, err = useCase.rpc.Rsk.GetTransactionReceipt(ctx, retainedQuote.UserRskTxHash); err != nil {
		return blockchain.TransactionReceipt{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, true)
	} else if chainHeight-receipt.BlockNumber < uint64(pegoutQuote.DepositConfirmations) {
		return blockchain.TransactionReceipt{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, usecases.NoEnoughConfirmationsError, true)
	} else if err = useCase.validateDepositEvent(receipt, &retainedQuote, pegoutQuote); err != nil {
		retainedQuote.UserRskTxHash = receipt.TransactionHash
		return blockchain.TransactionReceipt{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, false)
	} else if block, err = useCase.rpc.Rsk.GetBlockByHash(ctx, receipt.BlockHash); err != nil {
		return blockchain.TransactionReceipt{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, true)
	} else if pegoutQuote.ExpireTime().Before(block.Timestamp) || uint64(pegoutQuote.ExpireBlock) <= block.Number {
		return blockchain.TransactionReceipt{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, usecases.ExpiredQuoteError, false)
	}

	if completed, err = useCase.contracts.Lbc.IsPegOutQuoteCompleted(retainedQuote.QuoteHash); err != nil {
		return blockchain.TransactionReceipt{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, true)
	} else if completed {
		return blockchain.TransactionReceipt{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, fmt.Errorf("quote %s was already completed", retainedQuote.QuoteHash), false)
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

	quoteHashBytes, err := hex.DecodeString(retainedQuote.QuoteHash)
	if err != nil {
		retainedQuote.UserRskTxHash = receipt.TransactionHash
		return quote.RetainedPegoutQuote{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, false)
	}

	var txResult blockchain.BitcoinTransactionResult
	if txResult, err = useCase.btcWallet.SendWithOpReturn(pegoutQuote.DepositAddress, pegoutQuote.Value, quoteHashBytes); err != nil {
		newState = quote.PegoutStateSendPegoutFailed
	} else {
		newState = quote.PegoutStateSendPegoutSucceeded
	}

	creationData := useCase.quoteRepository.GetPegoutCreationData(ctx, retainedQuote.QuoteHash)

	retainedQuote.LpBtcTxHash = txResult.Hash
	if txResult.Fee != nil {
		retainedQuote.SendPegoutBtcFee = txResult.Fee
	}
	retainedQuote.State = newState
	useCase.eventBus.Publish(quote.PegoutBtcSentToUserEvent{
		Event:         entities.NewBaseEvent(quote.PegoutBtcSentEventId),
		PegoutQuote:   *pegoutQuote,
		RetainedQuote: retainedQuote,
		CreationData:  creationData,
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

func (useCase *SendPegoutUseCase) validateDepositEvent(
	receipt blockchain.TransactionReceipt,
	retainedQuote *quote.RetainedPegoutQuote,
	pegoutQuote *quote.PegoutQuote,
) error {
	depositEvent, err := useCase.depositParser(receipt)
	if err != nil {
		return err
	} else if !strings.EqualFold(depositEvent.RawLog.Address, pegoutQuote.LbcAddress) {
		return errors.New("invalid LBC address")
	} else if !utils.CompareIgnore0x(depositEvent.Log.QuoteHash, retainedQuote.QuoteHash) {
		return errors.New("deposit belongs to other quote")
	} else if depositEvent.Log.Amount.Cmp(pegoutQuote.Total()) < 0 {
		retainedQuote.UserRskTxHash = receipt.TransactionHash
		return usecases.InsufficientAmountError
	}
	return nil
}
