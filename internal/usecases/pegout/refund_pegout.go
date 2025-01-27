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

const (
	// refundPegoutGasLimit Fixed gas limit for refundPegout function, should change only if the function does
	refundPegoutGasLimit = 2500000
)

type RefundPegoutUseCase struct {
	quoteRepository quote.PegoutQuoteRepository
	contracts       blockchain.RskContracts
	eventBus        entities.EventBus
	rpc             blockchain.Rpc
	rskWalletMutex  sync.Locker
}

func NewRefundPegoutUseCase(
	quoteRepository quote.PegoutQuoteRepository,
	contracts blockchain.RskContracts,
	eventBus entities.EventBus,
	rpc blockchain.Rpc,
	rskWalletMutex sync.Locker,
) *RefundPegoutUseCase {
	return &RefundPegoutUseCase{
		quoteRepository: quoteRepository,
		contracts:       contracts,
		eventBus:        eventBus,
		rpc:             rpc,
		rskWalletMutex:  rskWalletMutex,
	}
}

func (useCase *RefundPegoutUseCase) Run(ctx context.Context, retainedQuote quote.RetainedPegoutQuote) error {
	var params blockchain.RefundPegoutParams
	var pegoutQuote *quote.PegoutQuote
	var err error

	if retainedQuote.State != quote.PegoutStateSendPegoutSucceeded {
		return useCase.publishErrorEvent(ctx, retainedQuote, usecases.WrongStateError, true)
	}

	if pegoutQuote, err = useCase.quoteRepository.GetQuote(ctx, retainedQuote.QuoteHash); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, err, true)
	} else if pegoutQuote == nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, usecases.QuoteNotFoundError, false)
	}

	if err = useCase.validateBtcTransaction(ctx, *pegoutQuote, retainedQuote); err != nil {
		return err
	}

	if params, err = useCase.buildRefundPegoutParams(ctx, retainedQuote); err != nil {
		return err
	}
	txConfig := blockchain.NewTransactionConfig(nil, refundPegoutGasLimit, nil)

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	if _, err = useCase.performRefundPegout(ctx, retainedQuote, txConfig, params); err != nil {
		return err
	}
	return nil
}

func (useCase *RefundPegoutUseCase) publishErrorEvent(ctx context.Context, retainedQuote quote.RetainedPegoutQuote, err error, recoverable bool) error {
	wrappedError := usecases.WrapUseCaseErrorArgs(usecases.RefundPegoutId, err, usecases.ErrorArg("quoteHash", retainedQuote.QuoteHash))
	if !recoverable {
		retainedQuote.State = quote.PegoutStateRefundPegOutFailed
		wrappedError = errors.Join(wrappedError, usecases.NonRecoverableError)
		if err = useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); err != nil {
			wrappedError = errors.Join(wrappedError, err)
		}
		useCase.eventBus.Publish(quote.PegoutQuoteCompletedEvent{
			Event:         entities.NewBaseEvent(quote.PegoutQuoteCompletedEventId),
			RetainedQuote: retainedQuote,
			Error:         wrappedError,
		})
	}
	return wrappedError
}

func (useCase *RefundPegoutUseCase) buildRefundPegoutParams(ctx context.Context, retainedQuote quote.RetainedPegoutQuote) (blockchain.RefundPegoutParams, error) {
	var merkleBranch blockchain.MerkleBranch
	var block blockchain.BitcoinBlockInformation
	var err error
	var rawTx, quoteHashBytes []byte
	var quoteHashFixedBytes [32]byte

	if merkleBranch, err = useCase.rpc.Btc.BuildMerkleBranch(retainedQuote.LpBtcTxHash); err != nil {
		return blockchain.RefundPegoutParams{}, useCase.publishErrorEvent(ctx, retainedQuote, err, true)
	}

	if block, err = useCase.rpc.Btc.GetTransactionBlockInfo(retainedQuote.LpBtcTxHash); err != nil {
		return blockchain.RefundPegoutParams{}, useCase.publishErrorEvent(ctx, retainedQuote, err, true)
	}

	if rawTx, err = useCase.rpc.Btc.GetRawTransaction(retainedQuote.LpBtcTxHash); err != nil {
		return blockchain.RefundPegoutParams{}, useCase.publishErrorEvent(ctx, retainedQuote, err, true)
	}

	if quoteHashBytes, err = hex.DecodeString(retainedQuote.QuoteHash); err != nil {
		return blockchain.RefundPegoutParams{}, useCase.publishErrorEvent(ctx, retainedQuote, err, false)
	}
	copy(quoteHashFixedBytes[:], quoteHashBytes)

	return blockchain.RefundPegoutParams{
		QuoteHash:          quoteHashFixedBytes,
		BtcRawTx:           rawTx,
		BtcBlockHeaderHash: block.Hash,
		MerkleBranchPath:   merkleBranch.Path,
		MerkleBranchHashes: merkleBranch.Hashes,
	}, nil
}

func (useCase *RefundPegoutUseCase) performRefundPegout(
	ctx context.Context,
	retainedQuote quote.RetainedPegoutQuote,
	txConfig blockchain.TransactionConfig,
	params blockchain.RefundPegoutParams,
) (quote.RetainedPegoutQuote, error) {
	var newState quote.PegoutState
	var refundPegoutTxHash string
	var err, updateError error

	if refundPegoutTxHash, err = useCase.contracts.Lbc.RefundPegout(txConfig, params); errors.Is(err, blockchain.WaitingForBridgeError) {
		return quote.RetainedPegoutQuote{}, useCase.publishErrorEvent(ctx, retainedQuote, err, true)
	} else if err != nil {
		newState = quote.PegoutStateRefundPegOutFailed
	} else {
		newState = quote.PegoutStateRefundPegOutSucceeded
	}

	retainedQuote.State = newState
	retainedQuote.RefundPegoutTxHash = refundPegoutTxHash
	useCase.eventBus.Publish(quote.PegoutQuoteCompletedEvent{
		Event:         entities.NewBaseEvent(quote.PegoutQuoteCompletedEventId),
		RetainedQuote: retainedQuote,
		Error:         err,
	})

	if updateError = useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); updateError != nil {
		err = errors.Join(err, updateError)
	}
	if err != nil {
		err = errors.Join(err, usecases.NonRecoverableError)
		return quote.RetainedPegoutQuote{}, usecases.WrapUseCaseErrorArgs(usecases.RefundPegoutId, err, usecases.ErrorArg("quoteHash", retainedQuote.QuoteHash))
	}
	return retainedQuote, nil
}

func (useCase *RefundPegoutUseCase) validateBtcTransaction(
	ctx context.Context,
	pegoutQuote quote.PegoutQuote,
	retainedQuote quote.RetainedPegoutQuote,
) error {
	var txInfo blockchain.BitcoinTransactionInformation
	var err error
	if txInfo, err = useCase.rpc.Btc.GetTransactionInfo(retainedQuote.LpBtcTxHash); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, err, true)
	} else if txInfo.Confirmations < uint64(pegoutQuote.TransferConfirmations) {
		return useCase.publishErrorEvent(ctx, retainedQuote, usecases.NoEnoughConfirmationsError, true)
	}
	return nil
}
