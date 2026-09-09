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
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type DepositParser = func(receipt blockchain.TransactionReceipt, quoteHash string, lbcAddress string) (blockchain.ParsedLog[quote.PegoutDeposit], error)

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
	if err := usecases.CheckPauseState(useCase.contracts.PegOut); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, err, true)
	}
	if err := useCase.validateRetainedQuote(ctx, retainedQuote); err != nil {
		return err
	}
	if useCase.isEscrowClaimPath(retainedQuote) {
		return useCase.runEscrowClaim(ctx, retainedQuote)
	}
	return useCase.runLegacyDeposit(ctx, retainedQuote)
}

func (useCase *SendPegoutUseCase) isEscrowClaimPath(retainedQuote quote.RetainedPegoutQuote) bool {
	return useCase.contracts.PegOutEscrow != nil && retainedQuote.State == quote.PegoutStateClaimed
}

func (useCase *SendPegoutUseCase) runEscrowClaim(ctx context.Context, retainedQuote quote.RetainedPegoutQuote) error {
	pegoutQuote, err := useCase.getEscrowQuote(ctx, retainedQuote)
	if err != nil {
		return err
	}
	if err = useCase.validateEscrowClaimed(ctx, retainedQuote, pegoutQuote); err != nil {
		return err
	}

	useCase.btcWalletMutex.Lock()
	defer useCase.btcWalletMutex.Unlock()
	if err = useCase.revalidateRetainedQuote(ctx, retainedQuote); err != nil {
		return err
	}
	if err = useCase.validateBalance(ctx, retainedQuote, pegoutQuote); err != nil {
		return err
	}

	retainedQuote, err = useCase.performSendPegout(ctx, retainedQuote, pegoutQuote, blockchain.TransactionReceipt{
		TransactionHash: retainedQuote.UserRskTxHash,
	})
	return useCase.processSendPegoutResult(ctx, retainedQuote, err)
}

func (useCase *SendPegoutUseCase) runLegacyDeposit(ctx context.Context, retainedQuote quote.RetainedPegoutQuote) error {
	pegoutQuote, err := useCase.getQuote(ctx, retainedQuote)
	if err != nil {
		return err
	}
	receipt, err := useCase.validateQuote(ctx, retainedQuote, pegoutQuote)
	if err != nil {
		return err
	}

	useCase.btcWalletMutex.Lock()
	defer useCase.btcWalletMutex.Unlock()
	if err = useCase.revalidateRetainedQuote(ctx, retainedQuote); err != nil {
		return err
	}
	if err = useCase.validateBalance(ctx, retainedQuote, pegoutQuote); err != nil {
		return err
	}

	retainedQuote, err = useCase.performSendPegout(ctx, retainedQuote, pegoutQuote, receipt)
	return useCase.processSendPegoutResult(ctx, retainedQuote, err)
}

func (useCase *SendPegoutUseCase) processSendPegoutResult(ctx context.Context, retainedQuote quote.RetainedPegoutQuote, err error) error {
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

func (useCase *SendPegoutUseCase) getQuote(ctx context.Context, retainedQuote quote.RetainedPegoutQuote) (*quote.PegoutQuote, error) {
	pegoutQuote, err := useCase.quoteRepository.GetQuote(ctx, retainedQuote.QuoteHash)
	if err != nil {
		return nil, useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, err, true)
	} else if pegoutQuote == nil {
		return nil, useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, usecases.QuoteNotFoundError, false)
	}
	return pegoutQuote, nil
}

func (useCase *SendPegoutUseCase) getEscrowQuote(ctx context.Context, retainedQuote quote.RetainedPegoutQuote) (*quote.PegoutQuote, error) {
	pegoutQuote, err := useCase.contracts.PegOutEscrow.GetPegOutQuote(retainedQuote.QuoteHash)
	if err != nil {
		return nil, useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, err, true)
	}
	if pegoutQuote.DepositAddress, err = useCase.encodeHexAddress(pegoutQuote.DepositAddress); err != nil {
		return nil, useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, err, false)
	}
	if pegoutQuote.BtcRefundAddress, err = useCase.encodeHexAddress(pegoutQuote.BtcRefundAddress); err != nil {
		return nil, useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, err, false)
	}
	if pegoutQuote.LpBtcAddress, err = useCase.encodeHexAddress(pegoutQuote.LpBtcAddress); err != nil {
		return nil, useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, err, false)
	}
	return &pegoutQuote, nil
}

func (useCase *SendPegoutUseCase) encodeHexAddress(hexAddress string) (string, error) {
	addressBytes, err := hex.DecodeString(strings.TrimPrefix(hexAddress, "0x"))
	if err != nil {
		return "", err
	}
	return useCase.rpc.Btc.EncodeAddress(addressBytes)
}

func (useCase *SendPegoutUseCase) validateEscrowClaimed(
	ctx context.Context,
	retainedQuote quote.RetainedPegoutQuote,
	pegoutQuote *quote.PegoutQuote,
) error {
	state, err := useCase.contracts.PegOutEscrow.GetPegOutState(retainedQuote.QuoteHash)
	if err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, true)
	}
	if state != blockchain.EscrowedPegOutStateClaimed {
		return useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, usecases.WrongStateError, true)
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

	if completed, err = useCase.contracts.PegOut.IsPegOutQuoteCompleted(retainedQuote.QuoteHash); err != nil {
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

	requestHashBytes, err := hex.DecodeString(retainedQuote.QuoteHash)
	if err != nil {
		retainedQuote.UserRskTxHash = receipt.TransactionHash
		return quote.RetainedPegoutQuote{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, false)
	}

	if err = useCase.validatePegoutTransaction(retainedQuote, pegoutQuote, requestHashBytes); err != nil {
		retainedQuote.UserRskTxHash = receipt.TransactionHash
		errorStr := err.Error()
		isNonRecoverable := strings.Contains(errorStr, "(non-recoverable)") ||
			strings.Contains(errorStr, "reverted with:")
		return quote.RetainedPegoutQuote{}, useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, !isNonRecoverable)
	}

	var txResult blockchain.BitcoinTransactionResult
	if txResult, err = useCase.btcWallet.SendWithOpReturn(pegoutQuote.DepositAddress, pegoutQuote.Value, requestHashBytes); err != nil {
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

func (useCase *SendPegoutUseCase) validatePegoutTransaction(
	retainedQuote quote.RetainedPegoutQuote,
	pegoutQuote *quote.PegoutQuote,
	requestHashBytes []byte,
) error {
	rawTx, err := useCase.btcWallet.CreateUnfundedTransactionWithOpReturn(
		pegoutQuote.DepositAddress,
		pegoutQuote.Value,
		requestHashBytes,
	)
	if err != nil {
		return fmt.Errorf("failed to create unfunded transaction (non-recoverable): %w", err)
	}

	if err = useCase.contracts.PegOut.ValidatePegout(retainedQuote.QuoteHash, rawTx); err != nil {
		return fmt.Errorf("transaction validation failed: %w", err)
	}

	return nil
}

func (useCase *SendPegoutUseCase) validateBalance(
	ctx context.Context,
	retainedQuote quote.RetainedPegoutQuote,
	pegoutQuote *quote.PegoutQuote,
) error {
	var err error
	var balance *entities.Wei

	if balance, err = useCase.btcWallet.GetBalance(); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, err, true)
	}
	requiredBalance := new(entities.Wei)
	requiredBalance.Add(pegoutQuote.Value, pegoutQuote.GasFee)
	if balance.Cmp(requiredBalance) < 0 {
		return useCase.publishErrorEvent(ctx, retainedQuote, *pegoutQuote, usecases.NoLiquidityError, true)
	}
	return nil
}

func (useCase *SendPegoutUseCase) validateRetainedQuote(ctx context.Context, retainedQuote quote.RetainedPegoutQuote) error {
	validState := retainedQuote.State == quote.PegoutStateWaitingForDepositConfirmations ||
		retainedQuote.State == quote.PegoutStateClaimed
	if !validState {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, usecases.WrongStateError, true)
	} else if retainedQuote.UserRskTxHash == "" {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, errors.New("user rsk tx hash not provided"), true)
	}
	return nil
}

func (useCase *SendPegoutUseCase) revalidateRetainedQuote(ctx context.Context, retainedQuote quote.RetainedPegoutQuote) error {
	if dbQuote, err := useCase.quoteRepository.GetRetainedQuote(ctx, retainedQuote.QuoteHash); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, err, true)
	} else if dbQuote == nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, usecases.QuoteNotFoundError, false)
	} else if dbQuote.State != retainedQuote.State || dbQuote.UserRskTxHash != retainedQuote.UserRskTxHash {
		return useCase.publishErrorEvent(ctx, retainedQuote, quote.PegoutQuote{}, usecases.WrongStateError, true)
	}
	return nil
}

func (useCase *SendPegoutUseCase) validateDepositEvent(
	receipt blockchain.TransactionReceipt,
	retainedQuote *quote.RetainedPegoutQuote,
	pegoutQuote *quote.PegoutQuote,
) error {
	depositEvent, err := useCase.depositParser(receipt, retainedQuote.QuoteHash, pegoutQuote.LbcAddress)
	if err != nil {
		return err
	} else if depositEvent.Log.Amount.Cmp(pegoutQuote.Total()) < 0 {
		retainedQuote.UserRskTxHash = receipt.TransactionHash
		return usecases.InsufficientAmountError
	}
	return nil
}
