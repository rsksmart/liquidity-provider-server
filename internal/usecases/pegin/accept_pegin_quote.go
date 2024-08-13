package pegin

import (
	"context"
	"encoding/hex"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"sync"
)

type AcceptQuoteUseCase struct {
	quoteRepository     quote.PeginQuoteRepository
	contracts           blockchain.RskContracts
	rpc                 blockchain.Rpc
	lp                  liquidity_provider.LiquidityProvider
	peginLp             liquidity_provider.PeginLiquidityProvider
	eventBus            entities.EventBus
	peginLiquidityMutex sync.Locker
}

func NewAcceptQuoteUseCase(
	quoteRepository quote.PeginQuoteRepository,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	lp liquidity_provider.LiquidityProvider,
	peginLp liquidity_provider.PeginLiquidityProvider,
	eventBus entities.EventBus,
	peginLiquidityMutex sync.Locker,
) *AcceptQuoteUseCase {
	return &AcceptQuoteUseCase{
		quoteRepository:     quoteRepository,
		contracts:           contracts,
		rpc:                 rpc,
		lp:                  lp,
		peginLp:             peginLp,
		eventBus:            eventBus,
		peginLiquidityMutex: peginLiquidityMutex,
	}
}

func (useCase *AcceptQuoteUseCase) Run(ctx context.Context, quoteHash string) (quote.AcceptedQuote, error) {
	var err error
	errorArgs := usecases.NewErrorArgs()
	var peginQuote *quote.PeginQuote
	var retainedQuote *quote.RetainedPeginQuote

	if peginQuote, err = useCase.quoteRepository.GetQuote(ctx, quoteHash); err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	} else if peginQuote == nil {
		errorArgs["quoteHash"] = quoteHash
		return quote.AcceptedQuote{}, usecases.WrapUseCaseErrorArgs(usecases.AcceptPeginQuoteId, usecases.QuoteNotFoundError, errorArgs)
	}

	if peginQuote.IsExpired() {
		errorArgs["quoteHash"] = quoteHash
		return quote.AcceptedQuote{}, usecases.WrapUseCaseErrorArgs(usecases.AcceptPeginQuoteId, usecases.ExpiredQuoteError, errorArgs)
	}

	useCase.peginLiquidityMutex.Lock()
	defer useCase.peginLiquidityMutex.Unlock()

	if retainedQuote, err = useCase.quoteRepository.GetRetainedQuote(ctx, quoteHash); err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	} else if retainedQuote != nil {
		return quote.AcceptedQuote{
			Signature:      retainedQuote.Signature,
			DepositAddress: retainedQuote.DepositAddress,
		}, nil
	}

	if retainedQuote, err = useCase.buildRetainedQuote(ctx, quoteHash, peginQuote); err != nil {
		return quote.AcceptedQuote{}, err
	}

	if err = entities.ValidateStruct(retainedQuote); err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	}
	if err = useCase.quoteRepository.InsertRetainedQuote(ctx, *retainedQuote); err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	}

	useCase.eventBus.Publish(quote.AcceptedPeginQuoteEvent{
		Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
		Quote:         *peginQuote,
		RetainedQuote: *retainedQuote,
	})

	return quote.AcceptedQuote{
		Signature:      retainedQuote.Signature,
		DepositAddress: retainedQuote.DepositAddress,
	}, nil
}

func (useCase *AcceptQuoteUseCase) calculateDerivationAddress(quoteHashBytes []byte, peginQuote quote.PeginQuote) (blockchain.FlyoverDerivation, error) {
	var err error
	errorArgs := usecases.NewErrorArgs()
	var fedInfo blockchain.FederationInfo
	var userBtcAddress, lpBtcAddress, lbcAddress []byte

	if userBtcAddress, err = useCase.rpc.Btc.DecodeAddress(peginQuote.BtcRefundAddress); err != nil {
		errorArgs["btcAddress"] = peginQuote.BtcRefundAddress
		return blockchain.FlyoverDerivation{}, usecases.WrapUseCaseErrorArgs(usecases.AcceptPeginQuoteId, err, errorArgs)
	} else if lpBtcAddress, err = useCase.rpc.Btc.DecodeAddress(peginQuote.LpBtcAddress); err != nil {
		errorArgs["btcAddress"] = peginQuote.LpBtcAddress
		return blockchain.FlyoverDerivation{}, usecases.WrapUseCaseErrorArgs(usecases.AcceptPeginQuoteId, err, errorArgs)
	} else if lbcAddress, err = blockchain.DecodeStringTrimPrefix(peginQuote.LbcAddress); err != nil {
		errorArgs["rskAddress"] = peginQuote.LbcAddress
		return blockchain.FlyoverDerivation{}, usecases.WrapUseCaseErrorArgs(usecases.AcceptPeginQuoteId, err, errorArgs)
	}

	if fedInfo, err = useCase.contracts.Bridge.FetchFederationInfo(); err != nil {
		return blockchain.FlyoverDerivation{}, usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	}
	return useCase.contracts.Bridge.GetFlyoverDerivationAddress(blockchain.FlyoverDerivationArgs{
		FedInfo:              fedInfo,
		LbcAdress:            lbcAddress,
		UserBtcRefundAddress: userBtcAddress,
		LpBtcAddress:         lpBtcAddress,
		QuoteHash:            quoteHashBytes,
	})
}

func (useCase *AcceptQuoteUseCase) calculateAndCheckLiquidity(ctx context.Context, peginQuote quote.PeginQuote) (*entities.Wei, error) {
	var err error
	var gasPrice *entities.Wei
	errorArgs := usecases.NewErrorArgs()

	gasLimit := new(entities.Wei).Add(
		entities.NewUWei(uint64(peginQuote.GasLimit)),
		entities.NewUWei(CallForUserExtraGas),
	)
	if gasPrice, err = useCase.rpc.Rsk.GasPrice(ctx); err != nil {
		return nil, usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	}
	gasCost := new(entities.Wei).Mul(gasLimit, gasPrice)
	requiredLiquidity := new(entities.Wei).Add(gasCost, peginQuote.Value)

	if err = useCase.peginLp.HasPeginLiquidity(ctx, requiredLiquidity); err != nil {
		errorArgs["amount"] = requiredLiquidity.String()
		return nil, usecases.WrapUseCaseErrorArgs(usecases.AcceptPeginQuoteId, usecases.NoLiquidityError, errorArgs)
	}
	return requiredLiquidity, nil
}

func (useCase *AcceptQuoteUseCase) buildRetainedQuote(ctx context.Context, quoteHash string, peginQuote *quote.PeginQuote) (*quote.RetainedPeginQuote, error) {
	var derivation blockchain.FlyoverDerivation
	var requiredLiquidity *entities.Wei
	var quoteHashBytes []byte
	var quoteSignature string
	var err error

	if quoteHashBytes, err = hex.DecodeString(quoteHash); err != nil {
		return nil, usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	}
	if derivation, err = useCase.calculateDerivationAddress(quoteHashBytes, *peginQuote); err != nil {
		return nil, err
	}
	if requiredLiquidity, err = useCase.calculateAndCheckLiquidity(ctx, *peginQuote); err != nil {
		return nil, err
	}
	if quoteSignature, err = useCase.lp.SignQuote(quoteHash); err != nil {
		return nil, usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	}

	return &quote.RetainedPeginQuote{
		QuoteHash:         quoteHash,
		DepositAddress:    derivation.Address,
		Signature:         quoteSignature,
		RequiredLiquidity: requiredLiquidity,
		State:             quote.PeginStateWaitingForDeposit,
	}, nil
}
