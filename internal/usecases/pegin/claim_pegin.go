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
}

func NewClaimPegInUseCase(
	claims rootstock.PegInClaimRepository,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	peginProvider liquidity_provider.PeginLiquidityProvider,
	rskWalletMutex sync.Locker,
) *ClaimPegInUseCase {
	return &ClaimPegInUseCase{
		claims:         claims,
		contracts:      contracts,
		rpc:            rpc,
		peginProvider:  peginProvider,
		rskWalletMutex: rskWalletMutex,
	}
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

	tx, err := useCase.rpc.Btc.GetTransactionInfo(depositTxID)
	if err != nil {
		return useCase.unavailable(err)
	}
	amount := tx.FirstOutputToAddress(entry.BtcAddress)
	if amount.Cmp(entities.NewWei(0)) <= 0 {
		return nil
	}
	fee, params, err := useCase.evaluateGates(ctx, existing, entry, depositTxID, tx, amount)
	if err != nil || fee == nil {
		return err
	}

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	enough, err := useCase.ensureSpendable(ctx, entry.RskAddress, depositTxID, amount, fee, params)
	if err != nil || !enough {
		return err
	}
	return useCase.armAndSubmit(ctx, existing, entry, depositTxID, params)
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

func (useCase *ClaimPegInUseCase) evaluateGates(
	ctx context.Context,
	existing *rootstock.PegInClaim,
	entry rootstock.PegInWatch,
	depositTxID string,
	tx blockchain.BitcoinTransactionInformation,
	amount *entities.Wei,
) (*entities.Wei, blockchain.RequestPegInParams, error) {
	requiredConfirmations, err := useCase.contracts.FlyoverConfigurations.GetRequiredPegInBtcConfirmations(amount)
	if err != nil {
		return nil, blockchain.RequestPegInParams{}, useCase.unavailable(err)
	}
	if tx.Confirmations < requiredConfirmations {
		return nil, blockchain.RequestPegInParams{}, nil
	}
	if err = usecases.CheckPauseLevel(useCase.contracts.PauseRegistry, blockchain.PauseLevelSoft); err != nil {
		if errors.Is(err, blockchain.ContractPausedError) {
			return nil, blockchain.RequestPegInParams{}, nil
		}
		return nil, blockchain.RequestPegInParams{}, useCase.unavailable(err)
	}
	fee, err := useCase.contracts.FlyoverConfigurations.CalculatePegInFee(amount)
	if err != nil {
		return nil, blockchain.RequestPegInParams{}, useCase.unavailable(err)
	}
	if amount.Cmp(fee) < 0 {
		return nil, blockchain.RequestPegInParams{}, nil
	}
	return useCase.requestParamsIfAccepted(ctx, existing, entry, depositTxID, amount, fee)
}

func (useCase *ClaimPegInUseCase) requestParamsIfAccepted(
	ctx context.Context,
	existing *rootstock.PegInClaim,
	entry rootstock.PegInWatch,
	depositTxID string,
	amount, fee *entities.Wei,
) (*entities.Wei, blockchain.RequestPegInParams, error) {
	rawTx, err := useCase.rpc.Btc.GetRawTransaction(depositTxID)
	if err != nil {
		return nil, blockchain.RequestPegInParams{}, useCase.unavailable(err)
	}
	if err = blockchain.RejectWitnessSerializedTx(rawTx); err != nil {
		return nil, blockchain.RequestPegInParams{}, usecases.WrapUseCaseError(usecases.ClaimPegInId, err)
	}
	params, err := useCase.buildRequestParams(entry.RskAddress, depositTxID, amount, fee, rawTx)
	if err != nil {
		return nil, blockchain.RequestPegInParams{}, useCase.unavailable(err)
	}
	// Gas estimation returns the raw revert, so a claim another caller already won would look
	// like an outage. The dry run decodes it into ErrPegInAlreadyProcessed.
	if err = useCase.contracts.PegIn.SimulateRequestPegIn(params); err != nil {
		if errors.Is(err, blockchain.ErrPegInAlreadyProcessed) {
			return nil, blockchain.RequestPegInParams{}, useCase.persistAlreadyProcessed(ctx, existing, entry, depositTxID)
		}
		if errors.Is(err, blockchain.ErrPegInBelowMinimum) {
			log.Debug(LogPegInClaimBelowMinimum(entry.RskAddress, depositTxID))
			return nil, blockchain.RequestPegInParams{}, nil
		}
		return nil, blockchain.RequestPegInParams{}, useCase.unavailable(err)
	}
	return fee, params, nil
}

func (useCase *ClaimPegInUseCase) ensureSpendable(
	ctx context.Context,
	rskAddress, depositTxID string,
	amount, fee *entities.Wei,
	params blockchain.RequestPegInParams,
) (bool, error) {
	estimatedGas, err := useCase.contracts.PegIn.EstimateRequestPegInGas(params)
	if err != nil {
		return false, useCase.unavailable(err)
	}
	gasPrice, err := useCase.rpc.Rsk.GasPrice(ctx)
	if err != nil {
		return false, useCase.unavailable(err)
	}
	gasCost := new(entities.Wei).Mul(gasPrice, entities.NewUWei(estimatedGas))
	payable, err := rootstock.CalculatePegInClaimPayableValue(amount, fee)
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
		log.Debug(LogPegInClaimInsufficientWalletLiquidity(rskAddress, depositTxID, available, required))
		return false, nil
	}
	return true, nil
}

func (useCase *ClaimPegInUseCase) armAndSubmit(
	ctx context.Context,
	existing *rootstock.PegInClaim,
	entry rootstock.PegInWatch,
	depositTxID string,
	params blockchain.RequestPegInParams,
) error {
	claim := rootstock.NewCandidatePegInClaim(entry, depositTxID, existing)
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
	return useCase.submit(ctx, stored, params)
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

func (useCase *ClaimPegInUseCase) buildRequestParams(
	rskAddress string,
	depositTxID string,
	amount *entities.Wei,
	fee *entities.Wei,
	rawTx []byte,
) (blockchain.RequestPegInParams, error) {
	block, err := useCase.rpc.Btc.GetTransactionBlockInfo(depositTxID)
	if err != nil {
		return blockchain.RequestPegInParams{}, err
	}
	merkle, err := useCase.rpc.Btc.BuildMerkleBranch(depositTxID)
	if err != nil {
		return blockchain.RequestPegInParams{}, err
	}
	return blockchain.RequestPegInParams{
		RskAddress:         rskAddress,
		BitcoinRawTx:       rawTx,
		BtcBlockHash:       block.Hash,
		MerkleBranchPath:   merkle.Path,
		MerkleBranchHashes: merkle.Hashes,
		Amount:             amount,
		Fee:                fee,
	}, nil
}

func (useCase *ClaimPegInUseCase) persistAlreadyProcessed(
	ctx context.Context,
	existing *rootstock.PegInClaim,
	entry rootstock.PegInWatch,
	depositTxID string,
) error {
	if existing != nil && existing.IsTerminal() {
		return nil
	}
	claim := rootstock.NewCandidatePegInClaim(entry, depositTxID, existing)
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
