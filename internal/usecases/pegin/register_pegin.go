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
	contracts       blockchain.RskContracts
	quoteRepository quote.PeginQuoteRepository
	eventBus        entities.EventBus
	rpc             blockchain.Rpc
	rskWalletMutex  sync.Locker
}

func NewRegisterPeginUseCase(
	contracts blockchain.RskContracts,
	quoteRepository quote.PeginQuoteRepository,
	eventBus entities.EventBus,
	rpc blockchain.Rpc,
	rskWalletMutex sync.Locker,
) *RegisterPeginUseCase {
	return &RegisterPeginUseCase{
		contracts:       contracts,
		quoteRepository: quoteRepository,
		eventBus:        eventBus,
		rpc:             rpc,
		rskWalletMutex:  rskWalletMutex,
	}
}

func (useCase *RegisterPeginUseCase) Run(ctx context.Context, retainedQuote quote.RetainedPeginQuote) error {
	var err error
	var peginQuote *quote.PeginQuote
	var params blockchain.RegisterPeginParams
	var userBtcTx blockchain.BitcoinTransactionInformation

	if retainedQuote.State != quote.PeginStateCallForUserSucceeded {
		return useCase.publishErrorEvent(ctx, retainedQuote, usecases.WrongStateError, true)
	}

	if peginQuote, err = useCase.quoteRepository.GetQuote(ctx, retainedQuote.QuoteHash); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, err, true)
	} else if peginQuote == nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, usecases.QuoteNotFoundError, false)
	}

	if userBtcTx, err = useCase.getUserBtcTransactionIfValid(ctx, retainedQuote); err != nil {
		return err
	}

	if params, err = useCase.buildRegisterPeginParams(*peginQuote, retainedQuote); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, err, true)
	}

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	if err = usecases.RegisterCoinbaseTransaction(useCase.rpc.Btc, useCase.contracts.Bridge, userBtcTx); err != nil {
		return useCase.publishErrorEvent(ctx, retainedQuote, err, errors.Is(err, blockchain.WaitingForBridgeError))
	}

	return useCase.performRegisterPegin(ctx, params, retainedQuote)
}

func (useCase *RegisterPeginUseCase) publishErrorEvent(ctx context.Context, retainedQuote quote.RetainedPeginQuote, err error, recoverable bool) error {
	errorArgs := usecases.NewErrorArgs()
	errorArgs["quoteHash"] = retainedQuote.QuoteHash
	errorArgs["btcTx"] = retainedQuote.UserBtcTxHash
	wrappedError := usecases.WrapUseCaseErrorArgs(usecases.RegisterPeginId, err, errorArgs)
	if !recoverable {
		retainedQuote.State = quote.PeginStateRegisterPegInFailed
		wrappedError = errors.Join(wrappedError, usecases.NonRecoverableError)
		if err = useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); err != nil {
			wrappedError = errors.Join(wrappedError, err)
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

	if rawBtcTx, err = useCase.rpc.Btc.GetRawTransaction(retainedQuote.UserBtcTxHash); err != nil {
		return blockchain.RegisterPeginParams{}, err
	}

	if pmt, err = useCase.rpc.Btc.GetPartialMerkleTree(retainedQuote.UserBtcTxHash); err != nil {
		return blockchain.RegisterPeginParams{}, err
	}

	if block, err = useCase.rpc.Btc.GetTransactionBlockInfo(retainedQuote.UserBtcTxHash); err != nil {
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

func (useCase *RegisterPeginUseCase) getUserBtcTransactionIfValid(ctx context.Context, retainedQuote quote.RetainedPeginQuote) (blockchain.BitcoinTransactionInformation, error) {
	var txInfo blockchain.BitcoinTransactionInformation
	var err error
	if txInfo, err = useCase.rpc.Btc.GetTransactionInfo(retainedQuote.UserBtcTxHash); err != nil {
		return blockchain.BitcoinTransactionInformation{}, useCase.publishErrorEvent(ctx, retainedQuote, err, true)
	} else if txInfo.Confirmations < useCase.contracts.Bridge.GetRequiredTxConfirmations() {
		return blockchain.BitcoinTransactionInformation{}, useCase.publishErrorEvent(ctx, retainedQuote, usecases.NoEnoughConfirmationsError, true)
	}
	return txInfo, nil
}

func (useCase *RegisterPeginUseCase) performRegisterPegin(ctx context.Context, params blockchain.RegisterPeginParams, retainedQuote quote.RetainedPeginQuote) error {
	var registerPeginTxHash string
	var newState quote.PeginState
	var err error

	if registerPeginTxHash, err = useCase.contracts.Lbc.RegisterPegin(params); errors.Is(err, blockchain.WaitingForBridgeError) {
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
