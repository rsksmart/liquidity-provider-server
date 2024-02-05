package pegout

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"sync"
)

type AcceptQuoteUseCase struct {
	quoteRepository      quote.PegoutQuoteRepository
	lbc                  blockchain.LiquidityBridgeContract
	lp                   entities.LiquidityProvider
	pegoutLp             entities.PegoutLiquidityProvider
	eventBus             entities.EventBus
	pegoutLiquidityMutex *sync.Mutex
}

func NewAcceptQuoteUseCase(
	quoteRepository quote.PegoutQuoteRepository,
	lbc blockchain.LiquidityBridgeContract,
	lp entities.LiquidityProvider,
	pegoutLp entities.PegoutLiquidityProvider,
	eventBus entities.EventBus,
	pegoutLiquidityMutex *sync.Mutex,
) *AcceptQuoteUseCase {
	return &AcceptQuoteUseCase{
		quoteRepository:      quoteRepository,
		lbc:                  lbc,
		lp:                   lp,
		pegoutLp:             pegoutLp,
		eventBus:             eventBus,
		pegoutLiquidityMutex: pegoutLiquidityMutex,
	}
}

func (useCase *AcceptQuoteUseCase) Run(ctx context.Context, quoteHash string) (quote.AcceptedQuote, error) {
	var err error
	errorArgs := usecases.NewErrorArgs()
	var pegoutQuote *quote.PegoutQuote
	var retainedQuote *quote.RetainedPegoutQuote
	var quoteSignature string

	requiredLiquidity := new(entities.Wei)

	if pegoutQuote, err = useCase.quoteRepository.GetQuote(ctx, quoteHash); err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPegoutQuoteId, err)
	} else if pegoutQuote == nil {
		errorArgs["quoteHash"] = quoteHash
		return quote.AcceptedQuote{}, usecases.WrapUseCaseErrorArgs(usecases.AcceptPegoutQuoteId, usecases.QuoteNotFoundError, errorArgs)
	}

	if pegoutQuote.IsExpired() {
		errorArgs["quoteHash"] = quoteHash
		return quote.AcceptedQuote{}, usecases.WrapUseCaseErrorArgs(usecases.AcceptPegoutQuoteId, usecases.ExpiredQuoteError, errorArgs)
	}

	useCase.pegoutLiquidityMutex.Lock()
	defer useCase.pegoutLiquidityMutex.Unlock()

	if retainedQuote, err = useCase.quoteRepository.GetRetainedQuote(ctx, quoteHash); err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPegoutQuoteId, err)
	} else if retainedQuote != nil {
		return quote.AcceptedQuote{
			Signature:      retainedQuote.Signature,
			DepositAddress: retainedQuote.DepositAddress,
		}, nil
	}

	if requiredLiquidity, err = useCase.calculateAndCheckLiquidity(ctx, *pegoutQuote); err != nil {
		return quote.AcceptedQuote{}, err
	}

	if quoteSignature, err = useCase.lp.SignQuote(quoteHash); err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPegoutQuoteId, err)
	}

	retainedQuote = &quote.RetainedPegoutQuote{
		QuoteHash:         quoteHash,
		DepositAddress:    useCase.lbc.GetAddress(),
		Signature:         quoteSignature,
		RequiredLiquidity: requiredLiquidity,
		State:             quote.PegoutStateWaitingForDeposit,
	}

	if err = entities.ValidateStruct(retainedQuote); err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPegoutQuoteId, err)
	}
	if err = useCase.quoteRepository.InsertRetainedQuote(ctx, *retainedQuote); err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPegoutQuoteId, err)
	}

	useCase.eventBus.Publish(quote.AcceptedPegoutQuoteEvent{
		Event:         entities.NewBaseEvent(quote.AcceptedPegoutQuoteEventId),
		Quote:         *pegoutQuote,
		RetainedQuote: *retainedQuote,
	})

	return quote.AcceptedQuote{
		Signature:      retainedQuote.Signature,
		DepositAddress: retainedQuote.DepositAddress,
	}, nil
}

func (useCase *AcceptQuoteUseCase) calculateAndCheckLiquidity(ctx context.Context, pegoutQuote quote.PegoutQuote) (*entities.Wei, error) {
	var err error
	requiredLiquidity := new(entities.Wei)
	errorArgs := usecases.NewErrorArgs()

	requiredLiquidity.Add(pegoutQuote.Value, pegoutQuote.GasFee)
	if err = useCase.pegoutLp.HasPegoutLiquidity(ctx, requiredLiquidity); err != nil {
		errorArgs["amount"] = requiredLiquidity.String()
		return nil, usecases.WrapUseCaseErrorArgs(usecases.AcceptPegoutQuoteId, usecases.NoLiquidityError, errorArgs)
	}
	return requiredLiquidity, nil
}
