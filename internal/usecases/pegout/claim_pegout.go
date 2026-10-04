package pegout

import (
	"context"
	"encoding/hex"
	"errors"
	"strings"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type ClaimOutcome uint8

const (
	ClaimOutcomeSkipped ClaimOutcome = iota
	ClaimOutcomeClaimed
	ClaimOutcomeClosed
)

var errRequestClosed = errors.New("peg-out request is no longer open")

type ClaimPegOutUseCase struct {
	contracts       blockchain.RskContracts
	rpc             blockchain.Rpc
	btcWallet       blockchain.BitcoinWallet
	lp              liquidity_provider.LiquidityProvider
	quoteRepository quote.PegoutQuoteRepository
	eventBus        entities.EventBus
	rskWalletMutex  sync.Locker
}

func NewClaimPegOutUseCase(
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	btcWallet blockchain.BitcoinWallet,
	lp liquidity_provider.LiquidityProvider,
	quoteRepository quote.PegoutQuoteRepository,
	eventBus entities.EventBus,
	rskWalletMutex sync.Locker,
) *ClaimPegOutUseCase {
	return &ClaimPegOutUseCase{
		contracts:       contracts,
		rpc:             rpc,
		btcWallet:       btcWallet,
		lp:              lp,
		quoteRepository: quoteRepository,
		eventBus:        eventBus,
		rskWalletMutex:  rskWalletMutex,
	}
}

func (useCase *ClaimPegOutUseCase) Run(ctx context.Context, candidate blockchain.PegOutRequested) (ClaimOutcome, error) {
	claimed, err := useCase.run(ctx, candidate)
	switch {
	case errors.Is(err, errRequestClosed):
		return ClaimOutcomeClosed, nil
	case err != nil:
		return ClaimOutcomeSkipped, err
	case claimed:
		return ClaimOutcomeClaimed, nil
	default:
		return ClaimOutcomeSkipped, nil
	}
}

func (useCase *ClaimPegOutUseCase) run(ctx context.Context, candidate blockchain.PegOutRequested) (bool, error) {
	if useCase.contracts.PegOutEscrow == nil {
		return false, nil
	}
	if err := usecases.CheckPauseState(useCase.contracts.PegOut); err != nil {
		return false, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	if skip, err := useCase.shouldSkipClaim(ctx, candidate.RequestHash); err != nil || skip {
		return false, err
	}
	pegoutQuote, signature, skip, err := useCase.prepareClaim(ctx, candidate.RequestHash)
	if err != nil || skip {
		return false, err
	}
	return useCase.performClaim(ctx, candidate.RequestHash, pegoutQuote, signature)
}

func (useCase *ClaimPegOutUseCase) shouldSkipClaim(ctx context.Context, requestHash string) (bool, error) {
	if skip, err := useCase.checkAlreadyClaimed(ctx, requestHash); err != nil || skip {
		return skip, err
	}
	if err := useCase.checkRequestedState(requestHash); err != nil {
		return false, err
	}
	return useCase.checkRestriction(ctx, requestHash)
}

func (useCase *ClaimPegOutUseCase) prepareClaim(
	ctx context.Context,
	requestHash string,
) (quote.PegoutQuote, []byte, bool, error) {
	pegoutQuote, err := useCase.loadEncodedQuote(requestHash)
	if err != nil {
		return quote.PegoutQuote{}, nil, false, err
	}
	skip, err := useCase.checkCapacity(ctx, requestHash, pegoutQuote)
	if err != nil || skip {
		return quote.PegoutQuote{}, nil, skip, err
	}
	signature, err := useCase.signQuote(pegoutQuote)
	if err != nil {
		return quote.PegoutQuote{}, nil, false, err
	}
	claimGas, err := useCase.estimateClaimGas(requestHash, signature)
	if err != nil {
		return quote.PegoutQuote{}, nil, false, err
	}
	skip, err = useCase.checkProfitability(ctx, requestHash, pegoutQuote, claimGas)
	if err != nil || skip {
		return quote.PegoutQuote{}, nil, skip, err
	}
	return pegoutQuote, signature, false, nil
}

func (useCase *ClaimPegOutUseCase) checkAlreadyClaimed(ctx context.Context, requestHash string) (bool, error) {
	retainedQuote, err := useCase.quoteRepository.GetRetainedQuote(ctx, requestHash)
	if err != nil {
		return false, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	if retainedQuote == nil {
		return false, nil
	}
	log.Debug(LogClaimPegoutAlreadyClaimed(requestHash))
	return true, nil
}

func (useCase *ClaimPegOutUseCase) checkRequestedState(requestHash string) error {
	state, err := useCase.contracts.PegOutEscrow.GetPegOutState(requestHash)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	if state == blockchain.EscrowedPegOutStateRequested {
		return nil
	}
	log.Info(LogClaimPegoutLostRace(requestHash))
	return errRequestClosed
}

func (useCase *ClaimPegOutUseCase) checkRestriction(ctx context.Context, requestHash string) (bool, error) {
	restrictedUntil, err := useCase.contracts.PegOutEscrow.RestrictedUntil(useCase.lp.RskAddress())
	if err != nil {
		return false, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	height, err := useCase.rpc.Rsk.GetHeight(ctx)
	if err != nil {
		return false, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	if restrictedUntil == 0 || height >= restrictedUntil {
		return false, nil
	}
	log.Debug(LogClaimPegoutRestrictedSkip(requestHash, restrictedUntil))
	return true, nil
}

func (useCase *ClaimPegOutUseCase) loadEncodedQuote(requestHash string) (quote.PegoutQuote, error) {
	pegoutQuote, err := useCase.contracts.PegOutEscrow.GetPegOutQuote(requestHash)
	if err != nil {
		return quote.PegoutQuote{}, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	if pegoutQuote.DepositAddress, err = useCase.encodeHexAddress(pegoutQuote.DepositAddress); err != nil {
		return quote.PegoutQuote{}, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	if pegoutQuote.BtcRefundAddress, err = useCase.encodeHexAddress(pegoutQuote.BtcRefundAddress); err != nil {
		return quote.PegoutQuote{}, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	if pegoutQuote.LpBtcAddress, err = useCase.encodeHexAddress(pegoutQuote.LpBtcAddress); err != nil {
		return quote.PegoutQuote{}, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	return pegoutQuote, nil
}

func (useCase *ClaimPegOutUseCase) encodeHexAddress(hexAddress string) (string, error) {
	addressBytes, err := hex.DecodeString(strings.TrimPrefix(hexAddress, "0x"))
	if err != nil {
		return "", err
	}
	return useCase.rpc.Btc.EncodeAddress(addressBytes)
}

func (useCase *ClaimPegOutUseCase) checkCapacity(
	ctx context.Context,
	requestHash string,
	pegoutQuote quote.PegoutQuote,
) (bool, error) {
	available, err := useCase.availableLiveLiquidity(ctx)
	if err != nil {
		return false, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	if available.Cmp(pegoutQuote.Value) >= 0 {
		return false, nil
	}
	log.Debug(LogClaimPegoutCapacitySkip(requestHash))
	return true, nil
}

func (useCase *ClaimPegOutUseCase) availableLiveLiquidity(ctx context.Context) (*entities.Wei, error) {
	balance, err := useCase.btcWallet.GetBalance()
	if err != nil {
		return nil, err
	}
	inFlight, err := useCase.quoteRepository.GetRetainedQuoteByState(
		ctx,
		quote.PegoutStateClaimPending,
		quote.PegoutStateClaimed,
		quote.PegoutStateWaitingForDeposit,
		quote.PegoutStateWaitingForDepositConfirmations,
	)
	if err != nil {
		return nil, err
	}
	locked := entities.NewWei(0)
	for _, retained := range inFlight {
		if retained.RequiredLiquidity != nil {
			locked.Add(locked, retained.RequiredLiquidity)
		}
	}
	if balance.Cmp(locked) < 0 {
		return entities.NewWei(0), nil
	}
	return new(entities.Wei).Sub(balance, locked), nil
}

func (useCase *ClaimPegOutUseCase) signQuote(pegoutQuote quote.PegoutQuote) ([]byte, error) {
	hash, err := useCase.contracts.PegOut.HashPegoutQuoteEIP712(pegoutQuote)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	signatureBytes, err := useCase.lp.GetSigner().SignBytes(hash[:])
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	signatureBytes[len(signatureBytes)-1] += 27
	return signatureBytes, nil
}

func (useCase *ClaimPegOutUseCase) estimateClaimGas(requestHash string, signature []byte) (*entities.Wei, error) {
	claimGas, err := useCase.contracts.PegOutEscrow.EstimateClaimPegOut(requestHash, signature)
	if err == nil {
		return claimGas, nil
	}
	state, stateErr := useCase.contracts.PegOutEscrow.GetPegOutState(requestHash)
	if stateErr != nil {
		return nil, usecases.WrapUseCaseError(usecases.ClaimPegoutId, errors.Join(err, stateErr))
	}
	if state != blockchain.EscrowedPegOutStateRequested {
		log.Info(LogClaimPegoutLostRace(requestHash))
		return nil, errRequestClosed
	}
	return nil, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
}

func (useCase *ClaimPegOutUseCase) checkProfitability(
	ctx context.Context,
	requestHash string,
	pegoutQuote quote.PegoutQuote,
	claimGas *entities.Wei,
) (bool, error) {
	btcFeeEstimation, err := useCase.btcWallet.EstimateTxFees(pegoutQuote.DepositAddress, pegoutQuote.Value)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "insufficient funds") {
			log.Debug(LogClaimPegoutCapacitySkip(requestHash))
			return true, nil
		}
		return false, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	gasPrice, err := useCase.rpc.Rsk.GasPrice(ctx)
	if err != nil {
		return false, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	totalGas := new(entities.Wei).Add(claimGas, entities.NewUWei(refundPegoutGasLimit))
	rskCost := new(entities.Wei).Mul(totalGas, gasPrice)
	totalCost := new(entities.Wei).Add(rskCost, btcFeeEstimation.Value)
	if pegoutQuote.CallFee.Cmp(totalCost) > 0 {
		return false, nil
	}
	log.Debug(LogClaimPegoutProfitabilitySkip(requestHash))
	return true, nil
}

func (useCase *ClaimPegOutUseCase) performClaim(
	ctx context.Context,
	requestHash string,
	pegoutQuote quote.PegoutQuote,
	signature []byte,
) (bool, error) {
	completedQuote, quoteHash, err := useCase.completeQuote(pegoutQuote)
	if err != nil {
		return false, err
	}

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	skip, err := useCase.checkClaimRecorded(ctx, requestHash, quoteHash)
	if err != nil || skip {
		return false, err
	}
	// Once the claim mines the LP owes the BTC, so it is recorded before it is sent.
	retainedQuote, err := useCase.persistPendingClaim(ctx, quoteHash, completedQuote, signature)
	if err != nil {
		return false, err
	}

	txConfig := blockchain.NewTransactionConfig(nil, 0, nil)
	receipt, err := useCase.contracts.PegOutEscrow.ClaimPegOut(txConfig, requestHash, signature)
	if err != nil {
		return useCase.handleClaimError(ctx, requestHash, retainedQuote, completedQuote, err)
	}
	// A successful claimPegOut re-keys the peg-out to quoteHash. If promoting fails,
	// ReconcilePendingClaims promotes the pending record later.
	if err = useCase.promoteClaim(ctx, retainedQuote, completedQuote, receipt.TransactionHash); err != nil {
		return false, err
	}
	log.Info(LogClaimPegoutSuccess(requestHash, quoteHash, receipt.TransactionHash))
	return true, nil
}

// checkClaimRecorded skips a request this LP already has a claim record for, e.g. a claim
// still being mined after a timeout, so it is never sent twice.
func (useCase *ClaimPegOutUseCase) checkClaimRecorded(ctx context.Context, requestHash, quoteHash string) (bool, error) {
	retainedQuote, err := useCase.quoteRepository.GetRetainedQuote(ctx, quoteHash)
	if err != nil {
		return false, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	if retainedQuote == nil {
		return false, nil
	}
	log.Info(LogClaimPegoutAlreadyClaimed(requestHash))
	return true, nil
}

// completeQuote returns the quote as claimPegOut stores it (lpRskAddress set to the claiming LP)
// and its hashPegOutQuote, the id the escrow, PegOutContract and the OP_RETURN use after the claim.
func (useCase *ClaimPegOutUseCase) completeQuote(pegoutQuote quote.PegoutQuote) (quote.PegoutQuote, string, error) {
	completedQuote := pegoutQuote
	completedQuote.LpRskAddress = useCase.lp.RskAddress()
	quoteHash, err := useCase.contracts.PegOut.HashPegoutQuote(completedQuote)
	if err != nil {
		return quote.PegoutQuote{}, "", usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	return completedQuote, quoteHash, nil
}

func (useCase *ClaimPegOutUseCase) handleClaimError(
	ctx context.Context,
	requestHash string,
	retainedQuote quote.RetainedPegoutQuote,
	completedQuote quote.PegoutQuote,
	claimErr error,
) (bool, error) {
	if errors.Is(claimErr, context.DeadlineExceeded) {
		// The claim may still be mined, ReconcilePendingClaims resolves the pending record.
		return false, usecases.WrapUseCaseError(usecases.ClaimPegoutId, claimErr)
	}
	// Read requestHash first: once it is no longer REQUESTED, the quoteHash state is final.
	requestState, err := useCase.contracts.PegOutEscrow.GetPegOutState(requestHash)
	if err != nil {
		return false, usecases.WrapUseCaseError(usecases.ClaimPegoutId, errors.Join(claimErr, err))
	}
	if requestState != blockchain.EscrowedPegOutStateRequested {
		quoteState, stateErr := useCase.contracts.PegOutEscrow.GetPegOutState(retainedQuote.QuoteHash)
		if stateErr != nil {
			return false, usecases.WrapUseCaseError(usecases.ClaimPegoutId, errors.Join(claimErr, stateErr))
		}
		if quoteState == blockchain.EscrowedPegOutStateClaimed {
			// The claim mined although the call failed, e.g. reading its receipt.
			if err = useCase.promoteClaim(ctx, retainedQuote, completedQuote, ""); err != nil {
				return false, errors.Join(claimErr, err)
			}
			log.Info(LogClaimPegoutSuccess(requestHash, retainedQuote.QuoteHash, ""))
			return true, nil
		}
	}
	if err = useCase.deletePendingClaim(ctx, retainedQuote.QuoteHash); err != nil {
		return false, errors.Join(usecases.WrapUseCaseError(usecases.ClaimPegoutId, claimErr), err)
	}
	if requestState != blockchain.EscrowedPegOutStateRequested {
		log.Info(LogClaimPegoutLostRace(requestHash))
		return false, errRequestClosed
	}
	return false, usecases.WrapUseCaseError(usecases.ClaimPegoutId, claimErr)
}

func (useCase *ClaimPegOutUseCase) persistPendingClaim(
	ctx context.Context,
	quoteHash string,
	pegoutQuote quote.PegoutQuote,
	signature []byte,
) (quote.RetainedPegoutQuote, error) {
	retainedQuote := quote.RetainedPegoutQuote{
		QuoteHash:         quoteHash,
		DepositAddress:    useCase.contracts.PegOut.GetAddress(),
		Signature:         hex.EncodeToString(signature),
		RequiredLiquidity: pegoutQuote.Value.Copy(),
		State:             quote.PegoutStateClaimPending,
		RemainingToRefund: pegoutQuote.Total(),
	}
	if err := entities.ValidateStruct(retainedQuote); err != nil {
		return quote.RetainedPegoutQuote{}, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	createdQuote := quote.CreatedPegoutQuote{
		Hash:         quoteHash,
		Quote:        pegoutQuote,
		CreationData: quote.PegoutCreationDataZeroValue(),
	}
	if err := useCase.quoteRepository.InsertQuote(ctx, createdQuote); err != nil {
		return quote.RetainedPegoutQuote{}, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	if err := useCase.quoteRepository.InsertRetainedQuote(ctx, retainedQuote); err != nil {
		err = usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
		return quote.RetainedPegoutQuote{}, errors.Join(err, useCase.deletePendingClaim(ctx, quoteHash))
	}
	return retainedQuote, nil
}

func (useCase *ClaimPegOutUseCase) promoteClaim(
	ctx context.Context,
	retainedQuote quote.RetainedPegoutQuote,
	pegoutQuote quote.PegoutQuote,
	claimTxHash string,
) error {
	retainedQuote.State = quote.PegoutStateClaimed
	retainedQuote.UserRskTxHash = claimTxHash
	if err := useCase.quoteRepository.UpdateRetainedQuote(ctx, retainedQuote); err != nil {
		return usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	useCase.eventBus.Publish(quote.ClaimedPegoutQuoteEvent{
		Event:         entities.NewBaseEvent(quote.ClaimedPegoutQuoteEventId),
		Quote:         pegoutQuote,
		RetainedQuote: retainedQuote,
	})
	return nil
}

func (useCase *ClaimPegOutUseCase) deletePendingClaim(ctx context.Context, quoteHash string) error {
	if _, err := useCase.quoteRepository.DeleteQuotes(ctx, []string{quoteHash}); err != nil {
		return usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	return nil
}

// ReconcilePendingClaims resolves pending claims whose outcome was not recorded, e.g. after a
// restart or a mining timeout: a mined claim is promoted, one that can no longer mine is deleted.
func (useCase *ClaimPegOutUseCase) ReconcilePendingClaims(ctx context.Context) error {
	if useCase.contracts.PegOutEscrow == nil {
		return nil
	}
	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	pending, err := useCase.quoteRepository.GetRetainedQuoteByState(ctx, quote.PegoutStateClaimPending)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	} else if len(pending) == 0 {
		return nil
	}
	// Read the time before any state: if a claim is still unmined after a block past its
	// depositDateLimit, claimPegOut can no longer succeed.
	block, err := useCase.rpc.Rsk.GetBlockByNumber(ctx, nil)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	now := uint64(block.Timestamp.Unix())
	errs := make([]error, 0)
	for _, retainedQuote := range pending {
		errs = append(errs, useCase.reconcilePendingClaim(ctx, retainedQuote, now))
	}
	return errors.Join(errs...)
}

func (useCase *ClaimPegOutUseCase) reconcilePendingClaim(ctx context.Context, retainedQuote quote.RetainedPegoutQuote, now uint64) error {
	pegoutQuote, err := useCase.quoteRepository.GetQuote(ctx, retainedQuote.QuoteHash)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	} else if pegoutQuote == nil {
		return usecases.WrapUseCaseError(usecases.ClaimPegoutId, usecases.QuoteNotFoundError)
	}
	state, err := useCase.contracts.PegOutEscrow.GetPegOutState(retainedQuote.QuoteHash)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	switch {
	case state == blockchain.EscrowedPegOutStateClaimed:
		if err = useCase.promoteClaim(ctx, retainedQuote, *pegoutQuote, ""); err != nil {
			return err
		}
		log.Info(LogClaimPegoutPendingPromoted(retainedQuote.QuoteHash))
	case state != blockchain.EscrowedPegOutStateNone:
		log.Error(LogClaimPegoutPendingSettled(retainedQuote.QuoteHash, state))
		return useCase.deletePendingClaim(ctx, retainedQuote.QuoteHash)
	case now > uint64(pegoutQuote.DepositDateLimit):
		log.Info(LogClaimPegoutPendingDropped(retainedQuote.QuoteHash))
		return useCase.deletePendingClaim(ctx, retainedQuote.QuoteHash)
	}
	return nil
}
