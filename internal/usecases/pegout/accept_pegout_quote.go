package pegout

import (
	"context"
	"sync"

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
	pegoutLp                 liquidity_provider.PegoutLiquidityProvider
	eventBus                 entities.EventBus
	pegoutLiquidityMutex     sync.Locker
	trustedAccountRepository liquidity_provider.TrustedAccountRepository
	hashFunction             entities.HashFunction
}

func NewAcceptQuoteUseCase(
	quoteRepository quote.PegoutQuoteRepository,
	contracts blockchain.RskContracts,
	lp liquidity_provider.LiquidityProvider,
	pegoutLp liquidity_provider.PegoutLiquidityProvider,
	eventBus entities.EventBus,
	pegoutLiquidityMutex sync.Locker,
	trustedAccountRepository liquidity_provider.TrustedAccountRepository,
	hashFunction entities.HashFunction,
) *AcceptQuoteUseCase {
	return &AcceptQuoteUseCase{
		quoteRepository:          quoteRepository,
		contracts:                contracts,
		lp:                       lp,
		pegoutLp:                 pegoutLp,
		eventBus:                 eventBus,
		pegoutLiquidityMutex:     pegoutLiquidityMutex,
		trustedAccountRepository: trustedAccountRepository,
		hashFunction:             hashFunction,
	}
}

func (useCase *AcceptQuoteUseCase) Run(ctx context.Context, quoteHash, signature string) (quote.AcceptedQuote, error) {
	var err error
	var pegoutQuote quote.PegoutQuote
	var retainedQuote *quote.RetainedPegoutQuote
	var quoteSignature string
	var requiredLiquidity *entities.Wei
	var trustedAccount liquidity_provider.TrustedAccountDetails

	if err = usecases.CheckPauseState(useCase.contracts.PegOut); err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPegoutQuoteId, err)
	}

	if pegoutQuote, err = useCase.getQuote(ctx, quoteHash); err != nil {
		return quote.AcceptedQuote{}, err
	}

	if trustedAccount, err = useCase.handleTrustedAccountSignature(ctx, signature, pegoutQuote); err != nil {
		return quote.AcceptedQuote{}, err
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

	if requiredLiquidity, err = useCase.calculateAndCheckLiquidity(ctx, pegoutQuote); err != nil {
		return quote.AcceptedQuote{}, err
	}

	if quoteSignature, err = useCase.lp.SignPegoutQuote(ctx, quoteHash); err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPegoutQuoteId, err)
	}

	retainedQuote = &quote.RetainedPegoutQuote{
		QuoteHash:           quoteHash,
		DepositAddress:      useCase.contracts.PegOut.GetAddress(),
		Signature:           quoteSignature,
		RequiredLiquidity:   requiredLiquidity,
		State:               quote.PegoutStateWaitingForDeposit,
		OwnerAccountAddress: trustedAccount.Address,
	}
	creationData := useCase.quoteRepository.GetPegoutCreationData(ctx, quoteHash)
	if err = useCase.publishQuote(ctx, pegoutQuote, retainedQuote, creationData); err != nil {
		return quote.AcceptedQuote{}, err
	}

	return quote.AcceptedQuote{
		Signature:      retainedQuote.Signature,
		DepositAddress: retainedQuote.DepositAddress,
	}, nil
}

func (useCase *AcceptQuoteUseCase) getQuote(ctx context.Context, quoteHash string) (quote.PegoutQuote, error) {
	var pegoutQuote *quote.PegoutQuote
	var err error
	errorArgs := usecases.NewErrorArgs()

	if pegoutQuote, err = useCase.quoteRepository.GetQuote(ctx, quoteHash); err != nil {
		return quote.PegoutQuote{}, usecases.WrapUseCaseError(usecases.AcceptPegoutQuoteId, err)
	} else if pegoutQuote == nil {
		errorArgs["quoteHash"] = quoteHash
		return quote.PegoutQuote{}, usecases.WrapUseCaseErrorArgs(usecases.AcceptPegoutQuoteId, usecases.QuoteNotFoundError, errorArgs)
	}
	if pegoutQuote.IsExpired() {
		errorArgs["quoteHash"] = quoteHash
		return quote.PegoutQuote{}, usecases.WrapUseCaseErrorArgs(usecases.AcceptPegoutQuoteId, usecases.ExpiredQuoteError, errorArgs)
	}

	return *pegoutQuote, nil
}

func (useCase *AcceptQuoteUseCase) handleTrustedAccountSignature(ctx context.Context, signature string, pegoutQuote quote.PegoutQuote) (liquidity_provider.TrustedAccountDetails, error) {
	if signature == "" {
		return liquidity_provider.TrustedAccountDetails{}, nil
	}
	trustedAccount, err := useCase.getTrustedAccount(ctx, pegoutQuote, useCase.lp.GetSigner(), signature)
	if err != nil {
		return liquidity_provider.TrustedAccountDetails{}, err
	}
	if err = useCase.checkLockingCap(ctx, trustedAccount, pegoutQuote); err != nil {
		return liquidity_provider.TrustedAccountDetails{}, err
	}
	return trustedAccount, nil
}

func (useCase *AcceptQuoteUseCase) getTrustedAccount(ctx context.Context, pegoutQuote quote.PegoutQuote, signer entities.Signer, signature string) (liquidity_provider.TrustedAccountDetails, error) {
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
	if err != nil {
		return liquidity_provider.TrustedAccountDetails{}, liquidity_provider.TamperedTrustedAccountError
	}
	return trustedAccount.Value, nil
}

func (useCase *AcceptQuoteUseCase) checkLockingCap(ctx context.Context, trustedAccount liquidity_provider.TrustedAccountDetails, pegoutQuote quote.PegoutQuote) error {
	errorArgs := usecases.NewErrorArgs()

	activeQuotesStates := []quote.PegoutState{
		quote.PegoutStateWaitingForDeposit,
		quote.PegoutStateWaitingForDepositConfirmations,
	}

	// Get all retained quotes for this trusted account
	quotes, err := useCase.quoteRepository.GetRetainedQuotesForAddress(ctx, trustedAccount.Address, activeQuotesStates...)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	}

	// Sum the total value of the quotes
	totalLocked := entities.NewUWei(0)
	for _, quote := range quotes {
		totalLocked = new(entities.Wei).Add(totalLocked, quote.RequiredLiquidity)
	}

	// Add the value of the new quote and gas fee
	totalWithNewQuote := new(entities.Wei).Add(totalLocked, pegoutQuote.Value)
	totalWithNewQuote = new(entities.Wei).Add(totalWithNewQuote, pegoutQuote.GasFee)

	// Check if the sum exceeds the locking cap
	if totalWithNewQuote.Cmp(trustedAccount.BtcLockingCap) > 0 {
		errorArgs["address"] = trustedAccount.Address
		errorArgs["currentLocked"] = totalLocked.String()
		errorArgs["lockingCap"] = trustedAccount.BtcLockingCap.String()
		return usecases.WrapUseCaseErrorArgs(
			usecases.AcceptPegoutQuoteId,
			usecases.LockingCapExceededError,
			errorArgs,
		)
	}

	return nil
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

func (useCase *AcceptQuoteUseCase) publishQuote(
	ctx context.Context,
	pegoutQuote quote.PegoutQuote,
	retainedQuote *quote.RetainedPegoutQuote,
	creationData quote.PegoutCreationData,
) error {
	var err error
	if err = entities.ValidateStruct(retainedQuote); err != nil {
		return usecases.WrapUseCaseError(usecases.AcceptPegoutQuoteId, err)
	}
	if err = useCase.quoteRepository.InsertRetainedQuote(ctx, *retainedQuote); err != nil {
		return usecases.WrapUseCaseError(usecases.AcceptPegoutQuoteId, err)
	}

	useCase.eventBus.Publish(quote.AcceptedPegoutQuoteEvent{
		Event:         entities.NewBaseEvent(quote.AcceptedPegoutQuoteEventId),
		Quote:         pegoutQuote,
		RetainedQuote: *retainedQuote,
		CreationData:  creationData,
	})
	return nil
}
