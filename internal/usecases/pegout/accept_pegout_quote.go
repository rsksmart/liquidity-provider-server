package pegout

import (
	"context"
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type AcceptQuoteUseCase struct {
	quoteRepository          quote.PegoutQuoteRepository
	contracts                blockchain.RskContracts
	lp                       liquidity_provider.LiquidityProvider
	trustedAccountRepository liquidity_provider.TrustedAccountRepository
	hashFunction             entities.HashFunction
}

func NewAcceptQuoteUseCase(
	quoteRepository quote.PegoutQuoteRepository,
	contracts blockchain.RskContracts,
	lp liquidity_provider.LiquidityProvider,
	trustedAccountRepository liquidity_provider.TrustedAccountRepository,
	hashFunction entities.HashFunction,
) *AcceptQuoteUseCase {
	return &AcceptQuoteUseCase{
		quoteRepository:          quoteRepository,
		contracts:                contracts,
		lp:                       lp,
		trustedAccountRepository: trustedAccountRepository,
		hashFunction:             hashFunction,
	}
}

func (useCase *AcceptQuoteUseCase) Run(ctx context.Context, quoteHash, signature string) (quote.AcceptedQuote, error) {
	if err := usecases.CheckPauseState(useCase.contracts.PegOut); err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPegoutQuoteId, err)
	}

	pegoutQuote, err := useCase.getQuote(ctx, quoteHash)
	if err != nil {
		return quote.AcceptedQuote{}, err
	}

	if _, err = useCase.getTrustedAccount(ctx, signature, pegoutQuote); err != nil && !errors.Is(err, liquidity_provider.NoSignatureError) {
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

func (useCase *AcceptQuoteUseCase) getTrustedAccount(ctx context.Context, signature string, pegoutQuote quote.PegoutQuote) (liquidity_provider.TrustedAccountDetails, error) {
	if signature == "" {
		return liquidity_provider.TrustedAccountDetails{}, liquidity_provider.NoSignatureError
	}
	trustedAccount, err := useCase.recoverTrustedAccount(ctx, pegoutQuote, useCase.lp.GetSigner(), signature)
	if err != nil {
		return liquidity_provider.TrustedAccountDetails{}, err
	}
	return trustedAccount, nil
}

func (useCase *AcceptQuoteUseCase) recoverTrustedAccount(ctx context.Context, pegoutQuote quote.PegoutQuote, signer entities.Signer, signature string) (liquidity_provider.TrustedAccountDetails, error) {
	address, err := usecases.RecoverSignerAddress(signature, func() ([]byte, error) {
		if hash, err := useCase.contracts.PegOut.HashPegoutQuoteEIP712(pegoutQuote); err != nil {
			return nil, err
		} else {
			return hash[:], nil
		}
	})
	if err != nil {
		return liquidity_provider.TrustedAccountDetails{}, err
	}

	trustedAccount, err := liquidity_provider.ValidateConfiguration(signer, useCase.hashFunction, func() (*entities.Signed[liquidity_provider.TrustedAccountDetails], error) {
		return useCase.trustedAccountRepository.GetTrustedAccount(ctx, address)
	})
	if err != nil && errors.Is(err, liquidity_provider.TrustedAccountNotFoundError) {
		return liquidity_provider.TrustedAccountDetails{}, err
	} else if err != nil {
		return liquidity_provider.TrustedAccountDetails{}, liquidity_provider.TamperedTrustedAccountError
	}
	return trustedAccount.Value, nil
}
