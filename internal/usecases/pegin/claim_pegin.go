package pegin

import (
	"context"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type ClaimPegInUseCase struct {
	claims         rootstock.PegInClaimRepository
	contracts      blockchain.RskContracts
	rpc            blockchain.Rpc
	peginProvider  liquidity_provider.PeginLiquidityProvider
	rskWalletMutex sync.Locker
	eventBus       entities.EventBus
}

func NewClaimPegInUseCase(
	claims rootstock.PegInClaimRepository,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	peginProvider liquidity_provider.PeginLiquidityProvider,
	rskWalletMutex sync.Locker,
	eventBus entities.EventBus,
) *ClaimPegInUseCase {
	return &ClaimPegInUseCase{
		claims:         claims,
		contracts:      contracts,
		rpc:            rpc,
		peginProvider:  peginProvider,
		rskWalletMutex: rskWalletMutex,
		eventBus:       eventBus,
	}
}

type claimRequest struct {
	entry       rootstock.PegInWatch
	depositTxID string
	existing    *rootstock.PegInClaim
	amount      *entities.Wei
	fee         *entities.Wei
	params      blockchain.RequestPegInParams
}

func (useCase *ClaimPegInUseCase) Run(
	ctx context.Context,
	entry rootstock.PegInWatch,
	depositTxID string,
) error {
	existing, ok, err := useCase.runnableClaim(ctx, entry.RskAddress, depositTxID)
	if err != nil || !ok {
		return err
	}
	request, ready, err := useCase.prepareClaim(ctx, existing, entry, depositTxID)
	if err != nil || !ready {
		return err
	}

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	enough, err := useCase.hasWalletLiquidity(ctx, request)
	if err != nil || !enough {
		return err
	}
	return useCase.armAndSubmit(ctx, request)
}

func (useCase *ClaimPegInUseCase) runnableClaim(
	ctx context.Context,
	rskAddress, depositTxID string,
) (*rootstock.PegInClaim, bool, error) {
	existing, err := useCase.claims.Get(ctx, rskAddress, depositTxID)
	if err != nil {
		return nil, false, useCase.unavailable(err)
	}
	if existing != nil && (existing.IsTerminal() || existing.TxHash != "") {
		return nil, false, nil
	}
	if existing != nil && existing.State == rootstock.PegInClaimSubmitting {
		log.Error(LogPegInClaimSubmittingEmptyTxHash(existing.RskAddress, existing.DepositTxID))
		return nil, false, nil
	}
	return existing, true, nil
}

func (useCase *ClaimPegInUseCase) prepareClaim(
	ctx context.Context,
	existing *rootstock.PegInClaim,
	entry rootstock.PegInWatch,
	depositTxID string,
) (request claimRequest, ready bool, err error) {
	tx, err := useCase.rpc.Btc.GetTransactionInfo(depositTxID)
	if err != nil {
		return claimRequest{}, false, useCase.unavailable(err)
	}
	amount := tx.FirstOutputToAddress(entry.BtcAddress)
	if amount.Cmp(entities.NewWei(0)) <= 0 {
		return claimRequest{}, false, nil
	}
	claimable, err := useCase.isClaimable(tx, amount)
	if err != nil || !claimable {
		return claimRequest{}, false, err
	}
	fee, err := useCase.contracts.FlyoverConfigurations.CalculatePegInFee(amount)
	if err != nil {
		return claimRequest{}, false, useCase.unavailable(err)
	}
	if amount.Cmp(fee) < 0 {
		return claimRequest{}, false, nil
	}
	request = claimRequest{entry: entry, depositTxID: depositTxID, existing: existing, amount: amount, fee: fee}
	request.params, ready, err = useCase.dryRunRequestParams(ctx, request)
	if err != nil || !ready {
		return claimRequest{}, false, err
	}
	return request, true, nil
}

func (useCase *ClaimPegInUseCase) isClaimable(
	tx blockchain.BitcoinTransactionInformation,
	amount *entities.Wei,
) (claimable bool, err error) {
	requiredConfirmations, err := useCase.contracts.FlyoverConfigurations.GetRequiredPegInBtcConfirmations(amount)
	if err != nil {
		return false, useCase.unavailable(err)
	}
	if tx.Confirmations < requiredConfirmations {
		return false, nil
	}
	if err = usecases.CheckPauseLevel(useCase.contracts.PauseRegistry, blockchain.PauseLevelSoft); err != nil {
		if errors.Is(err, blockchain.ContractPausedError) {
			return false, nil
		}
		return false, useCase.unavailable(err)
	}
	return true, nil
}

func (useCase *ClaimPegInUseCase) dryRunRequestParams(
	ctx context.Context,
	request claimRequest,
) (params blockchain.RequestPegInParams, ready bool, err error) {
	rawTx, err := useCase.rpc.Btc.GetRawTransaction(request.depositTxID)
	if err != nil {
		return blockchain.RequestPegInParams{}, false, useCase.unavailable(err)
	}
	if err = blockchain.RejectWitnessSerializedTx(rawTx); err != nil {
		return blockchain.RequestPegInParams{}, false, usecases.WrapUseCaseError(usecases.ClaimPegInId, err)
	}
	params, err = usecases.BuildRequestPegInParams(useCase.rpc.Btc, usecases.RequestPegInInput{
		RskAddress:  request.entry.RskAddress,
		DepositTxID: request.depositTxID,
		Amount:      request.amount,
		Fee:         request.fee,
		RawTx:       rawTx,
	})
	if err != nil {
		return blockchain.RequestPegInParams{}, false, useCase.unavailable(err)
	}
	// Gas estimation returns the raw revert, so a claim another caller already won would look
	// like an outage. The dry run decodes it into ErrPegInAlreadyProcessed.
	err = useCase.contracts.PegIn.SimulateRequestPegIn(params)
	switch {
	case err == nil:
		return params, true, nil
	case errors.Is(err, blockchain.ErrPegInAlreadyProcessed):
		return blockchain.RequestPegInParams{}, false, useCase.persistAlreadyProcessed(ctx, request)
	case errors.Is(err, blockchain.ErrPegInBelowMinimum):
		log.Debug(LogPegInClaimBelowMinimum(request.entry.RskAddress, request.depositTxID))
		return blockchain.RequestPegInParams{}, false, nil
	default:
		return blockchain.RequestPegInParams{}, false, useCase.unavailable(err)
	}
}

func (useCase *ClaimPegInUseCase) hasWalletLiquidity(ctx context.Context, request claimRequest) (enough bool, err error) {
	estimatedGas, err := useCase.contracts.PegIn.EstimateRequestPegInGas(request.params)
	if err != nil {
		return false, useCase.unavailable(err)
	}
	gasPrice, err := useCase.rpc.Rsk.GasPrice(ctx)
	if err != nil {
		return false, useCase.unavailable(err)
	}
	gasCost := new(entities.Wei).Mul(gasPrice, entities.NewUWei(estimatedGas))
	payable, err := rootstock.CalculatePegInClaimPayableValue(request.amount, request.fee)
	if err != nil {
		return false, usecases.WrapUseCaseError(usecases.ClaimPegInId, err)
	}
	required := new(entities.Wei).Add(payable, gasCost)
	// The wallet mutex is held until RequestPegIn returns the receipt, so this balance
	// already includes earlier claims.
	available, err := useCase.peginProvider.AvailablePeginWalletLiquidity(ctx)
	if err != nil {
		return false, useCase.unavailable(err)
	}
	if available.Cmp(required) < 0 {
		log.Debug(LogPegInClaimInsufficientWalletLiquidity(request.entry.RskAddress, request.depositTxID, available, required))
		return false, nil
	}
	return true, nil
}

func (useCase *ClaimPegInUseCase) armAndSubmit(ctx context.Context, request claimRequest) error {
	claim := rootstock.NewCandidatePegInClaim(request.entry, request.depositTxID, request.existing)
	stored, alreadySubmitted, err := useCase.save(ctx, claim)
	if err != nil {
		return useCase.unavailable(err)
	}
	if alreadySubmitted {
		return nil
	}
	stored.State = rootstock.PegInClaimSubmitting
	stored.TxHash = ""
	stored.UpdatedAt = time.Now().UTC()
	if err = useCase.claims.Update(ctx, stored); err != nil {
		return useCase.unavailable(err)
	}
	return useCase.submit(ctx, stored, request.params)
}

func (useCase *ClaimPegInUseCase) submit(
	ctx context.Context,
	claim rootstock.PegInClaim,
	params blockchain.RequestPegInParams,
) error {
	result, submitErr := useCase.contracts.PegIn.RequestPegIn(params)
	if result.Receipt.TransactionHash == "" && submitErr != nil {
		return useCase.classifySubmitError(ctx, claim, submitErr)
	}
	claim.TxHash = result.Receipt.TransactionHash
	if submitErr == nil {
		claim.PegInID = hex.EncodeToString(result.Event.PegInId[:])
	}
	// The claim stays submitting until SettlePegInClaimUseCase marks it claimed after maxReorgDepth.
	if persistErr := useCase.persistSubmission(ctx, &claim); persistErr != nil {
		return useCase.unavailable(errors.Join(persistErr, submitErr))
	}
	return nil
}

func (useCase *ClaimPegInUseCase) persistSubmission(ctx context.Context, claim *rootstock.PegInClaim) error {
	claim.UpdatedAt = time.Now().UTC()
	// RequestPegIn already waited on awaitTx with a detached context, so the
	// caller deadline can expire after broadcast. Bound this write with the
	// repository timeout instead of that deadline.
	persistCtx := context.WithoutCancel(ctx)
	return useCase.claims.Update(persistCtx, *claim)
}

func (useCase *ClaimPegInUseCase) persistAlreadyProcessed(ctx context.Context, request claimRequest) error {
	if request.existing != nil && request.existing.IsTerminal() {
		return nil
	}
	claim := rootstock.NewCandidatePegInClaim(request.entry, request.depositTxID, request.existing)
	stored, alreadySubmitted, err := useCase.save(ctx, claim)
	if err != nil {
		return useCase.unavailable(err)
	}
	if alreadySubmitted && stored.IsTerminal() {
		return nil
	}
	return useCase.classifySubmitError(ctx, stored, blockchain.ErrPegInAlreadyProcessed)
}

func (useCase *ClaimPegInUseCase) classifySubmitError(
	ctx context.Context,
	claim rootstock.PegInClaim,
	submitErr error,
) error {
	claim.UpdatedAt = time.Now().UTC()
	if errors.Is(submitErr, blockchain.ErrPegInAlreadyProcessed) {
		claim.State = rootstock.PegInClaimRaceLost
		if err := useCase.claims.Update(ctx, claim); err != nil {
			return useCase.unavailable(err)
		}
		useCase.eventBus.Publish(rootstock.NewPegInClaimCompletedEvent(claim))
		return nil
	}
	claim.State = rootstock.PegInClaimRetryableFailure
	if err := useCase.claims.Update(ctx, claim); err != nil {
		return useCase.unavailable(errors.Join(submitErr, err))
	}
	return usecases.WrapUseCaseError(usecases.ClaimPegInId, submitErr)
}

func (useCase *ClaimPegInUseCase) save(
	ctx context.Context,
	claim rootstock.PegInClaim,
) (stored rootstock.PegInClaim, alreadySubmitted bool, err error) {
	existing, err := useCase.claims.Get(ctx, claim.RskAddress, claim.DepositTxID)
	if err != nil {
		return rootstock.PegInClaim{}, false, err
	}
	if existing == nil {
		return useCase.insertClaim(ctx, claim)
	}
	return useCase.refreshClaim(ctx, claim, existing)
}

func (useCase *ClaimPegInUseCase) insertClaim(
	ctx context.Context,
	claim rootstock.PegInClaim,
) (stored rootstock.PegInClaim, alreadySubmitted bool, err error) {
	if err = useCase.claims.Insert(ctx, claim); err == nil {
		return claim, false, nil
	}
	if !errors.Is(err, rootstock.ErrPegInClaimAlreadyExists) {
		return rootstock.PegInClaim{}, false, err
	}
	existing, err := useCase.claims.Get(ctx, claim.RskAddress, claim.DepositTxID)
	if err != nil {
		return rootstock.PegInClaim{}, false, err
	}
	if existing == nil {
		return rootstock.PegInClaim{}, false, rootstock.ErrPegInClaimNotFound
	}
	return useCase.refreshClaim(ctx, claim, existing)
}

func (useCase *ClaimPegInUseCase) refreshClaim(
	ctx context.Context,
	claim rootstock.PegInClaim,
	existing *rootstock.PegInClaim,
) (stored rootstock.PegInClaim, alreadySubmitted bool, err error) {
	if existing.IsTerminal() || existing.State == rootstock.PegInClaimSubmitting {
		return *existing, true, nil
	}
	claim.CreatedAt = existing.CreatedAt
	if existing.TxHash != "" && claim.TxHash == "" {
		return *existing, true, nil
	}
	if err = useCase.claims.Update(ctx, claim); err != nil {
		return rootstock.PegInClaim{}, false, err
	}
	return claim, false, nil
}

func (useCase *ClaimPegInUseCase) unavailable(err error) error {
	if err == nil {
		return nil
	}
	return usecases.WrapUseCaseError(usecases.ClaimPegInId, errors.Join(err, usecases.InfrastructureUnavailableError))
}
