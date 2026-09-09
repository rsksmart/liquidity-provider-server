package pegout

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type AcceptQuoteUseCase struct {
	quoteRepository quote.PegoutQuoteRepository
	contracts       blockchain.RskContracts
	lp              liquidity_provider.LiquidityProvider
}

func NewAcceptQuoteUseCase(
	quoteRepository quote.PegoutQuoteRepository,
	contracts blockchain.RskContracts,
	lp liquidity_provider.LiquidityProvider,
) *AcceptQuoteUseCase {
	return &AcceptQuoteUseCase{
		quoteRepository: quoteRepository,
		contracts:       contracts,
		lp:              lp,
	}
}

func (useCase *AcceptQuoteUseCase) Run(ctx context.Context, quoteHash, _ string) (quote.AcceptedQuote, error) {
	if err := usecases.CheckPauseState(useCase.contracts.PegOut); err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPegoutQuoteId, err)
	}

	if _, err := useCase.getQuote(ctx, quoteHash); err != nil {
		return quote.AcceptedQuote{}, err
	}

	retainedQuote, err := useCase.quoteRepository.GetRetainedQuote(ctx, quoteHash)
	if err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPegoutQuoteId, err)
	}
	if retainedQuote != nil {
		return quote.AcceptedQuote{
			Signature:      retainedQuote.Signature,
			DepositAddress: retainedQuote.DepositAddress,
		}, nil
	}

	quoteSignature, err := useCase.lp.SignPegoutQuote(ctx, quoteHash)
	if err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPegoutQuoteId, err)
	}

	return quote.AcceptedQuote{
		Signature:      quoteSignature,
		DepositAddress: useCase.contracts.PegOut.GetAddress(),
	}, nil
}

func (useCase *AcceptQuoteUseCase) getQuote(ctx context.Context, quoteHash string) (quote.PegoutQuote, error) {
	errorArgs := usecases.NewErrorArgs()

	pegoutQuote, err := useCase.quoteRepository.GetQuote(ctx, quoteHash)
	if err != nil {
		return quote.PegoutQuote{}, usecases.WrapUseCaseError(usecases.AcceptPegoutQuoteId, err)
	}
	if pegoutQuote == nil {
		errorArgs["quoteHash"] = quoteHash
		return quote.PegoutQuote{}, usecases.WrapUseCaseErrorArgs(usecases.AcceptPegoutQuoteId, usecases.QuoteNotFoundError, errorArgs)
	}
	if pegoutQuote.IsExpired() {
		errorArgs["quoteHash"] = quoteHash
		return quote.PegoutQuote{}, usecases.WrapUseCaseErrorArgs(usecases.AcceptPegoutQuoteId, usecases.ExpiredQuoteError, errorArgs)
	}

	return *pegoutQuote, nil
}
