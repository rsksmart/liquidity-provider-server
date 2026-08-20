package pegin

import (
	"context"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type rskAccount interface {
	RskAddress() string
}

type ClaimPegInUseCase struct {
	claims         rootstock.PegInClaimRepository
	contracts      blockchain.RskContracts
	rpc            blockchain.Rpc
	account        rskAccount
	rskWalletMutex sync.Locker
	maxReorgDepth  uint64
}

func NewClaimPegInUseCase(
	claims rootstock.PegInClaimRepository,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	account rskAccount,
	rskWalletMutex sync.Locker,
	maxReorgDepth uint64,
) *ClaimPegInUseCase {
	return &ClaimPegInUseCase{
		claims:         claims,
		contracts:      contracts,
		rpc:            rpc,
		account:        account,
		rskWalletMutex: rskWalletMutex,
		maxReorgDepth:  maxReorgDepth,
	}
}

func (useCase *ClaimPegInUseCase) Run(
	ctx context.Context,
	entry rootstock.PegInWatch,
	depositTxID string,
) error {
	existing, err := useCase.claims.Get(ctx, entry.RskAddress, depositTxID)
	if err != nil {
		return useCase.wrap(err)
	}
	if isClaimNoOp(existing) {
		return nil
	}
	if existing != nil && existing.TxHash != "" {
		return useCase.reconcile(ctx, *existing)
	}

	tx, err := useCase.rpc.Btc.GetTransactionInfo(depositTxID)
	if err != nil {
		return useCase.wrap(err)
	}
	amount := tx.FirstOutputToAddress(entry.BtcAddress)
	if amount.Cmp(entities.NewWei(0)) <= 0 {
		return nil
	}
	fee, params, err := useCase.evaluateGates(ctx, existing, entry, depositTxID, tx, amount)
	if err != nil || fee == nil {
		return err
	}
	claim := newCandidateClaim(entry, depositTxID, payableValue(amount, fee), existing)
	if err = useCase.save(ctx, claim); err != nil {
		return useCase.wrap(err)
	}
	return useCase.submit(ctx, claim, params)
}

func (useCase *ClaimPegInUseCase) ReconcileSubmitting(ctx context.Context) error {
	claims, err := useCase.claims.ListByStates(ctx, rootstock.PegInClaimSubmitting)
	if err != nil {
		return useCase.wrap(err)
	}
	for _, claim := range claims {
		if err = useCase.reconcile(ctx, claim); err != nil {
			return err
		}
	}
	return nil
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
		return nil, blockchain.RequestPegInParams{}, useCase.wrap(err)
	}
	if tx.Confirmations < requiredConfirmations {
		return nil, blockchain.RequestPegInParams{}, useCase.releaseReserve(ctx, existing)
	}
	level, err := useCase.contracts.PauseRegistry.PauseLevel()
	if err != nil {
		return nil, blockchain.RequestPegInParams{}, useCase.wrap(err)
	}
	if level >= blockchain.PauseLevelHard {
		return nil, blockchain.RequestPegInParams{}, useCase.releaseReserve(ctx, existing)
	}
	fee, err := useCase.contracts.FlyoverConfigurations.CalculatePegInFee(amount)
	if err != nil {
		return nil, blockchain.RequestPegInParams{}, useCase.wrap(err)
	}
	inFlight, err := useCase.inFlightReserved(ctx, entry.RskAddress, depositTxID)
	if err != nil {
		return nil, blockchain.RequestPegInParams{}, useCase.wrap(err)
	}
	params, enough, err := useCase.checkClaimSpendable(ctx, existing, entry, depositTxID, amount, fee, inFlight)
	if err != nil || !enough {
		return nil, blockchain.RequestPegInParams{}, err
	}
	return fee, params, nil
}

func (useCase *ClaimPegInUseCase) checkClaimSpendable(
	ctx context.Context,
	existing *rootstock.PegInClaim,
	entry rootstock.PegInAddressRegistryWatchEntry,
	depositTxID string,
	amount, fee, inFlight *entities.Wei,
) (blockchain.RequestPegInParams, bool, error) {
	params, err := useCase.buildRequestParams(entry.RskAddress, depositTxID, amount, fee)
	if err != nil {
		return blockchain.RequestPegInParams{}, false, useCase.wrap(err)
	}
	wallet, err := useCase.rpc.Rsk.GetBalance(ctx, useCase.account.RskAddress())
	if err != nil {
		return blockchain.RequestPegInParams{}, false, useCase.wrap(err)
	}
	estimatedGas, err := useCase.contracts.PegIn.EstimateRequestPegInGas(params)
	if err != nil {
		return blockchain.RequestPegInParams{}, false, useCase.wrap(err)
	}
	gasPrice, err := useCase.rpc.Rsk.GasPrice(ctx)
	if err != nil {
		return blockchain.RequestPegInParams{}, false, useCase.wrap(err)
	}
	gasCost := new(entities.Wei).Mul(gasPrice, entities.NewUWei(estimatedGas))
	required := new(entities.Wei).Add(payableValue(amount, fee), gasCost)
	required.Add(required, inFlight)
	if wallet.Cmp(required) < 0 {
		return blockchain.RequestPegInParams{}, false, useCase.releaseReserve(ctx, existing)
	}
	return params, true, nil
}

func payableValue(amount, fee *entities.Wei) *entities.Wei {
	required := new(entities.Wei).Sub(amount.Copy(), fee.Copy())
	if required.Cmp(entities.NewWei(0)) < 0 {
		return entities.NewWei(0)
	}
	return required
}

func (useCase *ClaimPegInUseCase) submit(
	ctx context.Context,
	claim rootstock.PegInClaim,
	params blockchain.RequestPegInParams,
) error {
	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	result, submitErr := useCase.contracts.PegIn.RequestPegIn(params)
	if result.Receipt.TransactionHash != "" {
		claim.TxHash = result.Receipt.TransactionHash
		claim.State = rootstock.PegInClaimSubmitting
		claim.UpdatedAt = time.Now().UTC()
		if err := useCase.claims.Update(ctx, claim); err != nil {
			return useCase.wrap(err)
		}
	}
	if submitErr != nil {
		return useCase.classifySubmitError(ctx, claim, submitErr)
	}
	return useCase.finalizeSuccess(ctx, claim, result)
}

func (useCase *ClaimPegInUseCase) buildRequestParams(
	rskAddress string,
	depositTxID string,
	amount *entities.Wei,
	fee *entities.Wei,
) (blockchain.RequestPegInParams, error) {
	rawTx, err := useCase.rpc.Btc.GetRawTransaction(depositTxID)
	if err != nil {
		return blockchain.RequestPegInParams{}, err
	}
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

func (useCase *ClaimPegInUseCase) classifySubmitError(
	ctx context.Context,
	claim rootstock.PegInClaim,
	submitErr error,
) error {
	claim.ReservedWei = entities.NewWei(0)
	claim.UpdatedAt = time.Now().UTC()
	if errors.Is(submitErr, blockchain.ErrPegInAlreadyProcessed) {
		claim.State = rootstock.PegInClaimRaceLost
		if err := useCase.claims.Update(ctx, claim); err != nil {
			return useCase.wrap(err)
		}
		return nil
	}
	claim.State = rootstock.PegInClaimRetryableFailure
	if err := useCase.claims.Update(ctx, claim); err != nil {
		return useCase.wrap(errors.Join(submitErr, err))
	}
	return useCase.wrap(submitErr)
}

func (useCase *ClaimPegInUseCase) finalizeSuccess(
	ctx context.Context,
	claim rootstock.PegInClaim,
	result blockchain.RequestPegInResult,
) error {
	if result.Receipt.Status != blockchain.SuccessfulTxStatus {
		return useCase.failRetryable(ctx, claim, errors.New("requestPegIn receipt is not successful"))
	}
	height, err := useCase.rpc.Rsk.GetHeight(ctx)
	if err != nil {
		return useCase.wrap(err)
	}
	claim.PegInID = hex.EncodeToString(result.Event.PegInId[:])
	claim.UpdatedAt = time.Now().UTC()
	if result.Receipt.BlockNumber+useCase.maxReorgDepth <= height {
		claim.State = rootstock.PegInClaimClaimed
		claim.ReservedWei = entities.NewWei(0)
	} else {
		claim.State = rootstock.PegInClaimSubmitting
	}
	if err = useCase.claims.Update(ctx, claim); err != nil {
		return useCase.wrap(err)
	}
	return nil
}

func (useCase *ClaimPegInUseCase) failRetryable(
	ctx context.Context,
	claim rootstock.PegInClaim,
	cause error,
) error {
	claim.State = rootstock.PegInClaimRetryableFailure
	claim.ReservedWei = entities.NewWei(0)
	claim.UpdatedAt = time.Now().UTC()
	if err := useCase.save(ctx, claim); err != nil {
		return useCase.wrap(errors.Join(cause, err))
	}
	return useCase.wrap(cause)
}

// inFlightReserved sums ReservedWei on other candidate and submitting rows,
// excluding the current (rskAddress, depositTxID).
func (useCase *ClaimPegInUseCase) inFlightReserved(
	ctx context.Context,
	rskAddress string,
	depositTxID string,
) (*entities.Wei, error) {
	claims, err := useCase.claims.ListByStates(
		ctx,
		rootstock.PegInClaimCandidate,
		rootstock.PegInClaimSubmitting,
	)
	if err != nil {
		return nil, err
	}
	total := entities.NewWei(0)
	for _, claim := range claims {
		if claim.RskAddress == rskAddress && claim.DepositTxID == depositTxID {
			continue
		}
		if claim.ReservedWei != nil {
			total.Add(total, claim.ReservedWei)
		}
	}
	return total, nil
}

func (useCase *ClaimPegInUseCase) releaseReserve(ctx context.Context, existing *rootstock.PegInClaim) error {
	if existing == nil || existing.ReservedWei == nil || existing.ReservedWei.Cmp(entities.NewWei(0)) == 0 {
		return nil
	}
	existing.ReservedWei = entities.NewWei(0)
	existing.UpdatedAt = time.Now().UTC()
	return useCase.wrap(useCase.claims.Update(ctx, *existing))
}

func (useCase *ClaimPegInUseCase) save(ctx context.Context, claim rootstock.PegInClaim) error {
	existing, err := useCase.claims.Get(ctx, claim.RskAddress, claim.DepositTxID)
	if err != nil {
		return err
	}
	if existing == nil {
		if err = useCase.claims.Insert(ctx, claim); err == nil {
			return nil
		}
		if !errors.Is(err, rootstock.ErrPegInClaimAlreadyExists) {
			return err
		}
	}
	return useCase.claims.Update(ctx, claim)
}

func (useCase *ClaimPegInUseCase) wrap(err error) error {
	if err == nil {
		return nil
	}
	return usecases.WrapUseCaseError(usecases.ClaimPegInId, err)
}

func isClaimNoOp(existing *rootstock.PegInClaim) bool {
	if existing == nil {
		return false
	}
	switch existing.State {
	case rootstock.PegInClaimClaimed, rootstock.PegInClaimRaceLost:
		return true
	default:
		return false
	}
}

func newCandidateClaim(
	entry rootstock.PegInWatch,
	depositTxID string,
	reserved *entities.Wei,
	existing *rootstock.PegInClaim,
) rootstock.PegInClaim {
	now := time.Now().UTC()
	claim := rootstock.PegInClaim{
		RskAddress:  entry.RskAddress,
		DepositTxID: depositTxID,
		BtcAddress:  entry.BtcAddress,
		State:       rootstock.PegInClaimCandidate,
		ReservedWei: reserved.Copy(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if existing != nil {
		claim.CreatedAt = existing.CreatedAt
	}
	return claim
}
