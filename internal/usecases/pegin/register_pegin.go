package pegin

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

type RegisterPeginUseCase struct {
	lbc             blockchain.LiquidityBridgeContract
	quoteRepository quote.PeginQuoteRepository
	eventBus        entities.EventBus
	bridge          blockchain.RootstockBridge
	btc             blockchain.BitcoinNetwork
	rskWalletMutex  *sync.Mutex
}

func NewRegisterPeginUseCase(
	lbc blockchain.LiquidityBridgeContract,
	quoteRepository quote.PeginQuoteRepository,
	eventBus entities.EventBus,
	bridge blockchain.RootstockBridge,
	btc blockchain.BitcoinNetwork,
	rskWalletMutex *sync.Mutex,
) *RegisterPeginUseCase {
	return &RegisterPeginUseCase{
		lbc:             lbc,
		quoteRepository: quoteRepository,
		eventBus:        eventBus,
		bridge:          bridge,
		btc:             btc,
		rskWalletMutex:  rskWalletMutex,
	}
}

func (useCase *RegisterPeginUseCase) Run(ctx context.Context, retainedQuote quote.RetainedPeginQuote) error {
	var err error
	var peginQuote *quote.PeginQuote
	var params blockchain.RegisterPeginParams
	var newState quote.PeginState
	var registerPeginTxHash string

	if retainedQuote.State != quote.PeginStateCallForUserSucceeded {
		return useCase.publishErrorEvent(ctx, retainedQuote, usecases.WrongStateError, true)
	}

	if peginQuote, err = useCase.quoteRepository.GetQuote(ctx, retainedQuote.QuoteHash); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, err, true)
	} else if peginQuote == nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, usecases.QuoteNotFoundError, false)
	}

	if err = useCase.validateTransaction(ctx, retainedQuote); err != nil {
		return err
	}

	if params, err = useCase.buildRegisterPeginParams(*peginQuote, retainedQuote); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, err, true)
	}

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()
	if registerPeginTxHash, err = useCase.lbc.RegisterPegin(params); errors.Is(err, blockchain.WaitingForBridgeError) {
		return useCase.publishErrorEvent(ctx, retainedQuote, err, true)
	} else if err != nil {
		newState = quote.PeginStateRegisterPegInFailed
	} else {
		newState = quote.PeginStateRegisterPegInSucceeded
	}

	retainedQuote.State = newState
	retainedQuote.RegisterPeginTxHash = registerPeginTxHash
	useCase.eventBus.Publish(quote.RegisterPeginCompletedEvent{
		Event:         entities.NewBaseEvent(quote.RegisterPeginCompletedEventId),
		RetainedQuote: retainedQuote,
		Error:         err,
	})

	if updateError := useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); updateError != nil {
		err = errors.Join(err, updateError)
	}
	if err != nil {
		err = errors.Join(err, usecases.NonRecoverableError)
		return usecases.WrapUseCaseErrorArgs(usecases.RegisterPeginId, err, usecases.ErrorArg("quoteHash", retainedQuote.QuoteHash))
	}
	return nil
}

func (useCase *RegisterPeginUseCase) publishErrorEvent(ctx context.Context, retainedQuote quote.RetainedPeginQuote, err error, recoverable bool) error {
	errorArgs := usecases.NewErrorArgs()
	errorArgs["quoteHash"] = retainedQuote.QuoteHash
	errorArgs["btcTx"] = retainedQuote.UserBtcTxHash
	wrappedError := usecases.WrapUseCaseErrorArgs(usecases.RegisterPeginId, err, errorArgs)
	if !recoverable {
		retainedQuote.State = quote.PeginStateRegisterPegInFailed
		if err = useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); err != nil {
			wrappedError = errors.Join(wrappedError, err, usecases.NonRecoverableError)
		}
		useCase.eventBus.Publish(quote.RegisterPeginCompletedEvent{
			Event:         entities.NewBaseEvent(quote.RegisterPeginCompletedEventId),
			RetainedQuote: retainedQuote,
			Error:         wrappedError,
		})
	}
	return wrappedError
}

func (useCase *RegisterPeginUseCase) buildRegisterPeginParams(peginQuote quote.PeginQuote, retainedQuote quote.RetainedPeginQuote) (blockchain.RegisterPeginParams, error) {
	var quoteSignature, rawBtcTx, pmt []byte
	var block blockchain.BitcoinBlockInformation
	var err error

	if quoteSignature, err = hex.DecodeString(retainedQuote.Signature); err != nil {
		return blockchain.RegisterPeginParams{}, err
	}

	if rawBtcTx, err = useCase.btc.GetRawTransaction(retainedQuote.UserBtcTxHash); err != nil {
		return blockchain.RegisterPeginParams{}, err
	}

	if pmt, err = useCase.btc.GetPartialMerkleTree(retainedQuote.UserBtcTxHash); err != nil {
		return blockchain.RegisterPeginParams{}, err
	}

	if block, err = useCase.btc.GetTransactionBlockInfo(retainedQuote.UserBtcTxHash); err != nil {
		return blockchain.RegisterPeginParams{}, err
	}

	return blockchain.RegisterPeginParams{
		QuoteSignature:        quoteSignature,
		BitcoinRawTransaction: rawBtcTx,
		PartialMerkleTree:     pmt,
		BlockHeight:           block.Height,
		Quote:                 peginQuote,
	}, nil
}

func (useCase *RegisterPeginUseCase) validateTransaction(ctx context.Context, retainedQuote quote.RetainedPeginQuote) error {
	var txInfo blockchain.BitcoinTransactionInformation
	var err error
	if txInfo, err = useCase.btc.GetTransactionInfo(retainedQuote.UserBtcTxHash); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, err, true)
	} else if txInfo.Confirmations < useCase.bridge.GetRequiredTxConfirmations() {
		return useCase.publishErrorEvent(ctx, retainedQuote, usecases.NoEnoughConfirmationsError, true)
	}
	return nil
}
