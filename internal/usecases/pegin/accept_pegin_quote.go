package pegin

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type AcceptQuoteUseCase struct {
	quoteRepository          quote.PeginQuoteRepository
	contracts                blockchain.RskContracts
	rpc                      blockchain.Rpc
	lp                       liquidity_provider.LiquidityProvider
	peginLp                  liquidity_provider.PeginLiquidityProvider
	eventBus                 entities.EventBus
	peginLiquidityMutex      sync.Locker
	trustedAccountRepository liquidity_provider.TrustedAccountRepository
	hashFunction             entities.HashFunction
}

func NewAcceptQuoteUseCase(
	quoteRepository quote.PeginQuoteRepository,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	lp liquidity_provider.LiquidityProvider,
	peginLp liquidity_provider.PeginLiquidityProvider,
	eventBus entities.EventBus,
	peginLiquidityMutex sync.Locker,
	trustedAccountRepository liquidity_provider.TrustedAccountRepository,
	hashFunction entities.HashFunction,
) *AcceptQuoteUseCase {
	return &AcceptQuoteUseCase{
		quoteRepository:          quoteRepository,
		contracts:                contracts,
		rpc:                      rpc,
		lp:                       lp,
		peginLp:                  peginLp,
		eventBus:                 eventBus,
		peginLiquidityMutex:      peginLiquidityMutex,
		trustedAccountRepository: trustedAccountRepository,
		hashFunction:             hashFunction,
	}
}

func (useCase *AcceptQuoteUseCase) Run(ctx context.Context, quoteHash, signature string) (quote.AcceptedQuote, error) {
	logger := log.WithField("quote_hash", usecases.SafeLogStr(quoteHash))
	logger.WithField("has_signature", signature != "").Debug("Accepting pegin quote")

	peginQuote, err := useCase.loadValidQuote(ctx, quoteHash)
	if err != nil {
		return quote.AcceptedQuote{}, err
	}

	trustedAccount, err := useCase.getTrustedAccount(ctx, signature, *peginQuote)
	if err != nil && !errors.Is(err, liquidity_provider.NoSignatureError) {
		return quote.AcceptedQuote{}, err
	}

	useCase.peginLiquidityMutex.Lock()
	defer useCase.peginLiquidityMutex.Unlock()

	if err = useCase.validateTrustedAccountIfFound(ctx, logger, trustedAccount, *peginQuote, err); err != nil {
		return quote.AcceptedQuote{}, err
	}

	existing, err := useCase.quoteRepository.GetRetainedQuote(ctx, quoteHash)
	if err != nil {
		return quote.AcceptedQuote{}, usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	}
	if existing != nil {
		logger.WithField("deposit_address", existing.DepositAddress).
			Info("Accept pegin: returning cached signature")
		return quote.AcceptedQuote{
			Signature:      existing.Signature,
			DepositAddress: existing.DepositAddress,
		}, nil
	}

	retainedQuote, err := useCase.buildRetainedQuote(ctx, quoteHash, peginQuote, trustedAccount.Address)
	if err != nil {
		return quote.AcceptedQuote{}, err
	}
	if err = useCase.persistAndPublish(ctx, quoteHash, peginQuote, retainedQuote); err != nil {
		return quote.AcceptedQuote{}, err
	}

	return quote.AcceptedQuote{
		Signature:      retainedQuote.Signature,
		DepositAddress: retainedQuote.DepositAddress,
	}, nil
}

func (useCase *AcceptQuoteUseCase) persistAndPublish(ctx context.Context, quoteHash string, peginQuote *quote.PeginQuote, retainedQuote *quote.RetainedPeginQuote) error {
	logger := log.WithField("quote_hash", usecases.SafeLogStr(quoteHash))
	if err := useCase.quoteRepository.InsertRetainedQuote(ctx, *retainedQuote); err != nil {
		logger.WithError(err).Error("Accept pegin: failed to persist retained quote")
		return usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	}

	creationData := useCase.quoteRepository.GetPeginCreationData(ctx, quoteHash)

	useCase.eventBus.Publish(quote.AcceptedPeginQuoteEvent{
		Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
		Quote:         *peginQuote,
		RetainedQuote: *retainedQuote,
		CreationData:  creationData,
	})

	logger.WithFields(log.Fields{
		"deposit_address":    retainedQuote.DepositAddress,
		"required_liquidity": retainedQuote.RequiredLiquidity.String(),
		"owner":              retainedQuote.OwnerAccountAddress,
	}).Info("Accepted pegin quote")

	return nil
}

func (useCase *AcceptQuoteUseCase) loadValidQuote(
	ctx context.Context, quoteHash string,
) (*quote.PeginQuote, error) {
	logger := log.WithField("quote_hash", usecases.SafeLogStr(quoteHash))
	if err := usecases.CheckPauseState(useCase.contracts.PegIn); err != nil {
		logger.WithError(err).Warn("Accept pegin rejected: contract paused")
		return nil, usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	}

	peginQuote, err := useCase.quoteRepository.GetQuote(ctx, quoteHash)
	if err != nil {
		logger.WithError(err).Error("Accept pegin: failed to load quote")
		return nil, usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	}
	if peginQuote == nil {
		logger.Warn("Accept pegin rejected: quote not found")
		return nil, usecases.WrapUseCaseErrorArgs(
			usecases.AcceptPeginQuoteId, usecases.QuoteNotFoundError,
			usecases.ErrorArg("quoteHash", quoteHash),
		)
	}

	if peginQuote.IsExpired() {
		logger.WithField("expire_time", peginQuote.ExpireTime()).
			Warn("Accept pegin rejected: quote expired")
		return nil, usecases.WrapUseCaseErrorArgs(
			usecases.AcceptPeginQuoteId, usecases.ExpiredQuoteError,
			usecases.ErrorArg("quoteHash", quoteHash),
		)
	}

	return peginQuote, nil
}

func (useCase *AcceptQuoteUseCase) validateTrustedAccountIfFound(
	ctx context.Context,
	logger *log.Entry,
	trustedAccount liquidity_provider.TrustedAccountDetails,
	peginQuote quote.PeginQuote,
	lookupErr error,
) error {
	if errors.Is(lookupErr, liquidity_provider.NoSignatureError) {
		return nil
	}
	return useCase.checkLockingCap(ctx, logger, trustedAccount, peginQuote)
}

func (useCase *AcceptQuoteUseCase) getTrustedAccount(ctx context.Context, signature string, peginQuote quote.PeginQuote) (liquidity_provider.TrustedAccountDetails, error) {
	if signature == "" {
		return liquidity_provider.TrustedAccountDetails{}, liquidity_provider.NoSignatureError
	}
	trustedAccount, err := useCase.recoverTrustedAccount(ctx, peginQuote, useCase.lp.GetSigner(), signature)
	if err != nil {
		return liquidity_provider.TrustedAccountDetails{}, err
	}
	return trustedAccount, nil
}

func (useCase *AcceptQuoteUseCase) recoverTrustedAccount(ctx context.Context, peginQuote quote.PeginQuote, signer entities.Signer, signature string) (liquidity_provider.TrustedAccountDetails, error) {
	address, err := usecases.RecoverSignerAddress(signature, func() ([]byte, error) {
		if hash, err := useCase.contracts.PegIn.HashPeginQuoteEIP712(peginQuote); err != nil {
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

func (useCase *AcceptQuoteUseCase) checkLockingCap(ctx context.Context, logger *log.Entry, trustedAccount liquidity_provider.TrustedAccountDetails, peginQuote quote.PeginQuote) error {
	errorArgs := usecases.NewErrorArgs()

	activeQuotesStates := []quote.PeginState{
		quote.PeginStateWaitingForDeposit,
		quote.PeginStateWaitingForDepositConfirmations,
	}

	// Get all retained quotes for this trusted account
	quotes, err := useCase.quoteRepository.GetRetainedQuotesForAddress(ctx, trustedAccount.Address, activeQuotesStates...)
	if err != nil {
		logger.WithError(err).Error("Accept pegin: failed to load retained quotes")
		return usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	}

	// Sum the total value of the quotes
	totalLocked := entities.NewUWei(0)
	for _, quote := range quotes {
		totalLocked = new(entities.Wei).Add(totalLocked, quote.RequiredLiquidity)
	}

	newQuoteValue := new(entities.Wei).Add(peginQuote.Value, peginQuote.GasFee)
	totalWithNewQuote := new(entities.Wei).Add(totalLocked, newQuoteValue)

	// Check if the sum exceeds the locking cap
	if totalWithNewQuote.Cmp(trustedAccount.RbtcLockingCap) > 0 {
		logger.WithFields(log.Fields{
			"trusted_account": usecases.SafeLogStr(trustedAccount.Address),
			"current_locked":  totalLocked.String(),
			"locking_cap":     trustedAccount.RbtcLockingCap.String(),
			"new_quote_value": newQuoteValue.String(),
		}).Warn("Accept pegin rejected: locking cap exceeded")
		errorArgs["address"] = trustedAccount.Address
		errorArgs["currentLocked"] = totalLocked.String()
		errorArgs["lockingCap"] = trustedAccount.RbtcLockingCap.String()
		return usecases.WrapUseCaseErrorArgs(
			usecases.AcceptPeginQuoteId,
			usecases.LockingCapExceededError,
			errorArgs,
		)
	}

	return nil
}

func (useCase *AcceptQuoteUseCase) calculateDerivationAddress(quoteHashBytes []byte, peginQuote quote.PeginQuote) (rootstock.FlyoverDerivation, error) {
	var err error
	errorArgs := usecases.NewErrorArgs()
	var fedInfo rootstock.FederationInfo
	var userBtcAddress, lpBtcAddress, lbcAddress []byte

	if userBtcAddress, err = useCase.rpc.Btc.DecodeAddress(peginQuote.BtcRefundAddress); err != nil {
		errorArgs["btcAddress"] = peginQuote.BtcRefundAddress
		return rootstock.FlyoverDerivation{}, usecases.WrapUseCaseErrorArgs(usecases.AcceptPeginQuoteId, err, errorArgs)
	} else if lpBtcAddress, err = useCase.rpc.Btc.DecodeAddress(peginQuote.LpBtcAddress); err != nil {
		errorArgs["btcAddress"] = peginQuote.LpBtcAddress
		return rootstock.FlyoverDerivation{}, usecases.WrapUseCaseErrorArgs(usecases.AcceptPeginQuoteId, err, errorArgs)
	} else if lbcAddress, err = blockchain.DecodeStringTrimPrefix(peginQuote.LbcAddress); err != nil {
		errorArgs["rskAddress"] = peginQuote.LbcAddress
		return rootstock.FlyoverDerivation{}, usecases.WrapUseCaseErrorArgs(usecases.AcceptPeginQuoteId, err, errorArgs)
	}

	if fedInfo, err = useCase.contracts.Bridge.FetchFederationInfo(); err != nil {
		return rootstock.FlyoverDerivation{}, usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	}
	return useCase.contracts.Bridge.GetFlyoverDerivationAddress(rootstock.FlyoverDerivationArgs{
		FedInfo:              fedInfo,
		LbcAdress:            lbcAddress,
		UserBtcRefundAddress: userBtcAddress,
		LpBtcAddress:         lpBtcAddress,
		QuoteHash:            quoteHashBytes,
	})
}

func (useCase *AcceptQuoteUseCase) calculateAndCheckLiquidity(ctx context.Context, quoteHash string, peginQuote quote.PeginQuote) (*entities.Wei, error) {
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
		log.WithFields(log.Fields{
			"quote_hash": usecases.SafeLogStr(quoteHash),
			"required":   requiredLiquidity.String(),
		}).Warn("Accept pegin rejected: insufficient liquidity")
		errorArgs["amount"] = requiredLiquidity.String()
		return nil, usecases.WrapUseCaseErrorArgs(usecases.AcceptPeginQuoteId, usecases.NoLiquidityError, errorArgs)
	}
	return requiredLiquidity, nil
}

func (useCase *AcceptQuoteUseCase) buildRetainedQuote(ctx context.Context, quoteHash string, peginQuote *quote.PeginQuote, owner string) (*quote.RetainedPeginQuote, error) {
	var derivation rootstock.FlyoverDerivation
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
	if requiredLiquidity, err = useCase.calculateAndCheckLiquidity(ctx, quoteHash, *peginQuote); err != nil {
		return nil, err
	}
	if quoteSignature, err = useCase.lp.SignPeginQuote(ctx, quoteHash); err != nil {
		log.WithField("quote_hash", usecases.SafeLogStr(quoteHash)).WithError(err).
			Error("Accept pegin: failed to sign quote")
		return nil, usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	}

	retainedQuote := &quote.RetainedPeginQuote{
		QuoteHash:           quoteHash,
		DepositAddress:      derivation.Address,
		Signature:           quoteSignature,
		RequiredLiquidity:   requiredLiquidity,
		State:               quote.PeginStateWaitingForDeposit,
		OwnerAccountAddress: owner,
	}
	if err = entities.ValidateStruct(retainedQuote); err != nil {
		return nil, usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, err)
	}
	return retainedQuote, nil
}
