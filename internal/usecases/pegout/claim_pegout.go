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

func (useCase *ClaimPegOutUseCase) Run(ctx context.Context, candidate blockchain.PegOutRequested) (bool, error) {
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
	if skip, err := useCase.checkRequestedState(requestHash); err != nil || skip {
		return skip, err
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
	if claimGas == nil {
		return quote.PegoutQuote{}, nil, true, nil
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

func (useCase *ClaimPegOutUseCase) checkRequestedState(requestHash string) (bool, error) {
	state, err := useCase.contracts.PegOutEscrow.GetPegOutState(requestHash)
	if err != nil {
		return false, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	if state == blockchain.EscrowedPegOutStateRequested {
		return false, nil
	}
	log.Info(LogClaimPegoutLostRace(requestHash))
	return true, nil
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
		quote.PegoutStateClaimed,
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
		return nil, nil
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
	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	txConfig := blockchain.NewTransactionConfig(nil, 0, nil)
	receipt, err := useCase.contracts.PegOutEscrow.ClaimPegOut(txConfig, requestHash, signature)
	if err != nil {
		return false, useCase.handleClaimError(requestHash, err)
	}

	state, err := useCase.contracts.PegOutEscrow.GetPegOutState(requestHash)
	if err != nil {
		return false, usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	if state != blockchain.EscrowedPegOutStateClaimed {
		log.Info(LogClaimPegoutLostRace(requestHash))
		return false, nil
	}
	if err = useCase.persistClaim(ctx, requestHash, pegoutQuote, signature, receipt.TransactionHash); err != nil {
		return false, err
	}
	return true, nil
}

func (useCase *ClaimPegOutUseCase) handleClaimError(requestHash string, claimErr error) error {
	state, err := useCase.contracts.PegOutEscrow.GetPegOutState(requestHash)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ClaimPegoutId, errors.Join(claimErr, err))
	}
	if state != blockchain.EscrowedPegOutStateRequested {
		log.Info(LogClaimPegoutLostRace(requestHash))
		return nil
	}
	return usecases.WrapUseCaseError(usecases.ClaimPegoutId, claimErr)
}

func (useCase *ClaimPegOutUseCase) persistClaim(
	ctx context.Context,
	requestHash string,
	pegoutQuote quote.PegoutQuote,
	signature []byte,
	claimTxHash string,
) error {
	retainedQuote := quote.RetainedPegoutQuote{
		QuoteHash:         requestHash,
		DepositAddress:    useCase.contracts.PegOut.GetAddress(),
		Signature:         hex.EncodeToString(signature),
		RequiredLiquidity: pegoutQuote.Value.Copy(),
		State:             quote.PegoutStateClaimed,
		UserRskTxHash:     claimTxHash,
		RemainingToRefund: pegoutQuote.Total(),
	}
	if err := entities.ValidateStruct(retainedQuote); err != nil {
		return usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	createdQuote := quote.CreatedPegoutQuote{
		Hash:         requestHash,
		Quote:        pegoutQuote,
		CreationData: quote.PegoutCreationDataZeroValue(),
	}
	if err := useCase.quoteRepository.InsertQuote(ctx, createdQuote); err != nil {
		return usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	if err := useCase.quoteRepository.InsertRetainedQuote(ctx, retainedQuote); err != nil {
		return usecases.WrapUseCaseError(usecases.ClaimPegoutId, err)
	}
	useCase.eventBus.Publish(quote.ClaimedPegoutQuoteEvent{
		Event:         entities.NewBaseEvent(quote.ClaimedPegoutQuoteEventId),
		Quote:         pegoutQuote,
		RetainedQuote: retainedQuote,
	})
	log.Info(LogClaimPegoutSuccess(requestHash, claimTxHash))
	return nil
}
