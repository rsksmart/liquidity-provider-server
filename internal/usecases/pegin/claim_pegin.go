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
	account        liquidity_provider.LiquidityProvider
	rskWalletMutex sync.Locker
	maxReorgDepth  uint64
}

func NewClaimPegInUseCase(
	claims rootstock.PegInClaimRepository,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	account liquidity_provider.LiquidityProvider,
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
		return useCase.releaseReserve(ctx, existing)
	}
	fee, params, err := useCase.evaluateGates(ctx, existing, entry, depositTxID, tx, amount)
	if err != nil || fee == nil {
		return err
	}

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	payable, err := useCase.ensureSpendable(ctx, existing, entry.RskAddress, depositTxID, amount, fee, params)
	if err != nil || payable == nil {
		return err
	}
	return useCase.armAndSubmit(ctx, existing, entry, depositTxID, payable, params)
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
		return nil, blockchain.RequestPegInParams{}, useCase.releaseReserve(ctx, existing)
	}
	level, err := useCase.contracts.PauseRegistry.PauseLevel()
	if err != nil {
		return nil, blockchain.RequestPegInParams{}, useCase.unavailable(err)
	}
	if level >= blockchain.PauseLevelHard {
		return nil, blockchain.RequestPegInParams{}, useCase.releaseReserve(ctx, existing)
	}
	fee, err := useCase.contracts.FlyoverConfigurations.CalculatePegInFee(amount)
	if err != nil {
		return nil, blockchain.RequestPegInParams{}, useCase.unavailable(err)
	}
	if amount.Cmp(fee) < 0 {
		return nil, blockchain.RequestPegInParams{}, useCase.releaseReserve(ctx, existing)
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
		if releaseErr := useCase.releaseReserve(ctx, existing); releaseErr != nil {
			return nil, blockchain.RequestPegInParams{}, releaseErr
		}
		return nil, blockchain.RequestPegInParams{}, usecases.WrapUseCaseError(usecases.ClaimPegInId, err)
	}
	params, err := useCase.buildRequestParams(entry.RskAddress, depositTxID, amount, fee, rawTx)
	if err != nil {
		return nil, blockchain.RequestPegInParams{}, useCase.unavailable(err)
	}
	return fee, params, nil
}

func (useCase *ClaimPegInUseCase) ensureSpendable(
	ctx context.Context,
	existing *rootstock.PegInClaim,
	rskAddress, depositTxID string,
	amount, fee *entities.Wei,
	params blockchain.RequestPegInParams,
) (*entities.Wei, error) {
	inFlight, err := useCase.inFlightReserved(ctx, rskAddress, depositTxID)
	if err != nil {
		return nil, useCase.unavailable(err)
	}
	estimatedGas, err := useCase.contracts.PegIn.EstimateRequestPegInGas(params)
	if err != nil {
		return nil, useCase.unavailable(err)
	}
	gasPrice, err := useCase.rpc.Rsk.GasPrice(ctx)
	if err != nil {
		return nil, useCase.unavailable(err)
	}
	wallet, err := useCase.rpc.Rsk.GetBalance(ctx, useCase.account.RskAddress())
	if err != nil {
		return nil, useCase.unavailable(err)
	}
	gasCost := new(entities.Wei).Mul(gasPrice, entities.NewUWei(estimatedGas))
	payable, err := rootstock.CalculatePegInClaimPayableValue(amount, fee)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.ClaimPegInId, err)
	}
	required := new(entities.Wei).Add(payable, gasCost)
	required.Add(required, inFlight)
	if wallet.Cmp(required) < 0 {
		return nil, useCase.releaseReserve(ctx, existing)
	}
	return payable, nil
}

func (useCase *ClaimPegInUseCase) armAndSubmit(
	ctx context.Context,
	existing *rootstock.PegInClaim,
	entry rootstock.PegInWatch,
	depositTxID string,
	payable *entities.Wei,
	params blockchain.RequestPegInParams,
) error {
	claim := rootstock.NewCandidatePegInClaim(entry, depositTxID, payable, existing)
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
	if result.Receipt.TransactionHash != "" {
		if persistErr := useCase.persistTxHash(ctx, &claim, result.Receipt.TransactionHash); persistErr != nil {
			return useCase.unavailable(errors.Join(persistErr, submitErr))
		}
		if submitErr == nil {
			return useCase.finalizeSuccess(ctx, claim, result)
		}
		return nil
	}
	if submitErr != nil {
		return useCase.classifySubmitError(ctx, claim, submitErr)
	}
	return useCase.finalizeSuccess(ctx, claim, result)
}

func (useCase *ClaimPegInUseCase) persistTxHash(
	ctx context.Context,
	claim *rootstock.PegInClaim,
	txHash string,
) error {
	claim.TxHash = txHash
	claim.UpdatedAt = time.Now().UTC()
	// RequestPegIn waits with a detached context, so the caller can be canceled
	// after broadcast but before this critical write. Detach caller cancellation
	// and let the repository's DatabaseInteraction timeout bound TxHash persistence.
	persistCtx := context.WithoutCancel(ctx)
	var persistErr error
	for range 3 {
		persistErr = useCase.claims.Update(persistCtx, *claim)
		if persistErr == nil {
			break
		}
	}
	return persistErr
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

func (useCase *ClaimPegInUseCase) finalizeSuccess(
	ctx context.Context,
	claim rootstock.PegInClaim,
	result blockchain.RequestPegInResult,
) error {
	height, err := useCase.rpc.Rsk.GetHeight(ctx)
	if err != nil {
		return useCase.unavailable(err)
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
		return useCase.unavailable(err)
	}
	return nil
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
		if (claim.RskAddress != rskAddress || claim.DepositTxID != depositTxID) && claim.ReservedWei != nil {
			total.Add(total, claim.ReservedWei)
		}
	}
	return total, nil
}

func (useCase *ClaimPegInUseCase) releaseReserve(ctx context.Context, existing *rootstock.PegInClaim) error {
	if existing == nil || existing.ReservedWei == nil || existing.ReservedWei.Cmp(entities.NewWei(0)) == 0 {
		return nil
	}
	current, err := useCase.claims.Get(ctx, existing.RskAddress, existing.DepositTxID)
	if err != nil {
		return useCase.unavailable(err)
	}
	if !reservableClaim(current) {
		return nil
	}
	current.ReservedWei = entities.NewWei(0)
	current.UpdatedAt = time.Now().UTC()
	return useCase.unavailable(useCase.claims.Update(ctx, *current))
}

func reservableClaim(current *rootstock.PegInClaim) bool {
	return current != nil &&
		!current.IsTerminal() &&
		current.State != rootstock.PegInClaimSubmitting &&
		current.TxHash == "" &&
		current.ReservedWei != nil &&
		current.ReservedWei.Cmp(entities.NewWei(0)) > 0
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
