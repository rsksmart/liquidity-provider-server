package pegout

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	refundPegoutGasLimit = 2500000
	// BridgeConversionGasLimit see https://dev.rootstock.io/rsk/rbtc/conversion/networks/
	bridgeConversionGasLimit = 100000
	// BridgeConversionGasPrice see https://dev.rootstock.io/rsk/rbtc/conversion/networks/
	bridgeConversionGasPrice = 60000000
)

type RefundPegoutUseCase struct {
	quoteRepository quote.PegoutQuoteRepository
	lbc             blockchain.LiquidityBridgeContract
	eventBus        entities.EventBus
	btc             blockchain.BitcoinNetwork
	rskWallet       blockchain.RootstockWallet
	bridge          blockchain.RootstockBridge
	rskWalletMutex  *sync.Mutex
}

func NewRefundPegoutUseCase(
	quoteRepository quote.PegoutQuoteRepository,
	lbc blockchain.LiquidityBridgeContract,
	eventBus entities.EventBus,
	btc blockchain.BitcoinNetwork,
	rskWallet blockchain.RootstockWallet,
	bridge blockchain.RootstockBridge,
	rskWalletMutex *sync.Mutex,
) *RefundPegoutUseCase {
	return &RefundPegoutUseCase{
		quoteRepository: quoteRepository,
		lbc:             lbc,
		eventBus:        eventBus,
		btc:             btc,
		rskWallet:       rskWallet,
		bridge:          bridge,
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

	if params, err = useCase.buildRefundPegoutParams(retainedQuote); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, err, true)
	}
	txConfig := blockchain.NewTransactionConfig(nil, refundPegoutGasLimit, nil)

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	if err = useCase.performRefundPegout(ctx, retainedQuote, txConfig, params); err != nil {
		return err
	}

	if _, sendRbtcError := useCase.sendRbtcToBridge(ctx, *pegoutQuote); err != nil {
		err = errors.Join(err, sendRbtcError)
	} else if updateError := useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); updateError != nil {
		err = errors.Join(err, updateError)
	}

	if err != nil {
		return usecases.WrapUseCaseErrorArgs(usecases.RefundPegoutId, err, usecases.ErrorArg("quoteHash", retainedQuote.QuoteHash))
	}
	return nil
}

func (useCase *RefundPegoutUseCase) publishErrorEvent(ctx context.Context, retainedQuote quote.RetainedPegoutQuote, err error, recoverable bool) error {
	wrappedError := usecases.WrapUseCaseErrorArgs(usecases.RefundPegoutId, err, usecases.ErrorArg("quoteHash", retainedQuote.QuoteHash))
	if !recoverable {
		if err = useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); err != nil {
			wrappedError = errors.Join(wrappedError, err)
		}
		retainedQuote.State = quote.PegoutStateRefundPegOutFailed
		useCase.eventBus.Publish(quote.PegoutQuoteCompletedEvent{
			Event:         entities.NewBaseEvent(quote.PegoutQuoteCompletedEventId),
			RetainedQuote: retainedQuote,
			Error:         wrappedError,
		})
	}
	return wrappedError
}

func (useCase *RefundPegoutUseCase) buildRefundPegoutParams(retainedQuote quote.RetainedPegoutQuote) (blockchain.RefundPegoutParams, error) {
	var merkleBranch blockchain.MerkleBranch
	var block blockchain.BitcoinBlockInformation
	var err error
	var rawTx, quoteHashBytes []byte
	var quoteHashFixedBytes [32]byte

	if merkleBranch, err = useCase.btc.BuildMerkleBranch(retainedQuote.LpBtcTxHash); err != nil {
		return blockchain.RefundPegoutParams{}, err
	}

	if block, err = useCase.btc.GetTransactionBlockInfo(retainedQuote.LpBtcTxHash); err != nil {
		return blockchain.RefundPegoutParams{}, err
	}

	if rawTx, err = useCase.btc.GetRawTransaction(retainedQuote.LpBtcTxHash); err != nil {
		return blockchain.RefundPegoutParams{}, err
	}

	if quoteHashBytes, err = hex.DecodeString(retainedQuote.QuoteHash); err != nil {
		return blockchain.RefundPegoutParams{}, err
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

func (useCase *RefundPegoutUseCase) sendRbtcToBridge(ctx context.Context, pegoutQuote quote.PegoutQuote) (string, error) {
	var err error
	var txHash string
	value := new(entities.Wei)

	value.Add(pegoutQuote.Value, pegoutQuote.CallFee)
	value.Add(value, pegoutQuote.GasFee)
	config := blockchain.NewTransactionConfig(value, bridgeConversionGasLimit, entities.NewWei(bridgeConversionGasPrice))
	if txHash, err = useCase.rskWallet.SendRbtc(ctx, config, useCase.bridge.GetAddress()); err != nil {
		return "", err
	}
	log.Debugf("%s: transaction sent to the bridge successfully (%s)\n", usecases.RefundPegoutId, txHash)
	return txHash, nil
}

func (useCase *RefundPegoutUseCase) performRefundPegout(
	ctx context.Context,
	retainedQuote quote.RetainedPegoutQuote,
	txConfig blockchain.TransactionConfig,
	params blockchain.RefundPegoutParams,
) error {
	var newState quote.PegoutState
	var refundPegoutTxHash string
	var err error

	if refundPegoutTxHash, err = useCase.lbc.RefundPegout(txConfig, params); errors.Is(err, blockchain.WaitingForBridgeError) {
		return useCase.publishErrorEvent(ctx, retainedQuote, err, true)
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
	return nil
}

func (useCase *RefundPegoutUseCase) validateBtcTransaction(
	ctx context.Context,
	pegoutQuote quote.PegoutQuote,
	retainedQuote quote.RetainedPegoutQuote,
) error {
	var txInfo blockchain.BitcoinTransactionInformation
	var err error
	if txInfo, err = useCase.btc.GetTransactionInfo(retainedQuote.LpBtcTxHash); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, err, true)
	} else if txInfo.Confirmations < uint64(pegoutQuote.TransferConfirmations) {
		return useCase.publishErrorEvent(ctx, retainedQuote, usecases.NoEnoughConfirmationsError, true)
	}
	return nil
}
