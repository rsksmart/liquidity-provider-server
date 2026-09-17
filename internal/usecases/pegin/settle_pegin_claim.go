package pegin

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

var ErrStatus0ReceiptStillCallable = errors.New("status-0 receipt; preflight no longer reverts; not resubmitting")

type SettlePegInClaimUseCase struct {
	claims    rootstock.PegInClaimRepository
	contracts blockchain.RskContracts
	rpc       blockchain.Rpc
}

func NewSettlePegInClaimUseCase(
	claims rootstock.PegInClaimRepository,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
) *SettlePegInClaimUseCase {
	return &SettlePegInClaimUseCase{
		claims:    claims,
		contracts: contracts,
		rpc:       rpc,
	}
}

func (useCase *SettlePegInClaimUseCase) Run(ctx context.Context, claim rootstock.PegInClaim) error {
	if claim.TxHash == "" {
		log.Error(LogPegInClaimEmptyTxHash(claim.RskAddress, claim.DepositTxID))
		return nil
	}
	receipt, err := useCase.rpc.Rsk.GetTransactionReceipt(ctx, claim.TxHash)
	if err != nil && !errors.Is(err, blockchain.TxFailedError) {
		if errors.Is(err, blockchain.ErrTransactionReceiptNotFound) {
			log.Error(LogPegInClaimMissingReceipt(claim.TxHash, claim.RskAddress, claim.DepositTxID))
			return nil
		}
		return useCase.unavailable(err)
	}
	onChain, chainErr := useCase.receiptOnCanonicalChain(ctx, receipt)
	if chainErr != nil {
		return chainErr
	}
	if !onChain {
		return nil
	}
	if errors.Is(err, blockchain.TxFailedError) {
		return useCase.identifyFailed(ctx, claim)
	}
	return useCase.finalizeSuccess(ctx, claim, receipt)
}

func (useCase *SettlePegInClaimUseCase) receiptOnCanonicalChain(
	ctx context.Context,
	receipt blockchain.TransactionReceipt,
) (onChain bool, err error) {
	for _, eventLog := range receipt.Logs {
		if eventLog.Removed {
			return false, nil
		}
	}
	if receipt.BlockHash == "" || receipt.BlockNumber == 0 {
		return false, nil
	}
	block, err := useCase.rpc.Rsk.GetBlockByNumber(ctx, big.NewInt(int64(receipt.BlockNumber)))
	if err != nil {
		return false, useCase.unavailable(err)
	}
	if !strings.EqualFold(block.Hash, receipt.BlockHash) {
		return false, nil
	}
	return true, nil
}

func (useCase *SettlePegInClaimUseCase) finalizeSuccess(
	ctx context.Context,
	claim rootstock.PegInClaim,
	receipt blockchain.TransactionReceipt,
) error {
	event, unpackErr := useCase.contracts.PegIn.UnpackPegInRequested(receipt)
	if unpackErr != nil {
		log.Error(LogPegInClaimMissingEvent(claim.TxHash, claim.RskAddress, claim.DepositTxID, unpackErr))
		claim.PegInID = ""
	} else {
		claim.PegInID = hex.EncodeToString(event.PegInId[:])
	}
	claim.State = rootstock.PegInClaimClaimed
	claim.ReservedWei = entities.NewWei(0)
	claim.UpdatedAt = time.Now().UTC()
	if err := useCase.claims.Update(ctx, claim); err != nil {
		return useCase.unavailable(err)
	}
	return nil
}

func (useCase *SettlePegInClaimUseCase) identifyFailed(ctx context.Context, claim rootstock.PegInClaim) error {
	tx, err := useCase.rpc.Btc.GetTransactionInfo(claim.DepositTxID)
	if err != nil {
		return useCase.unavailable(err)
	}
	amount := tx.FirstOutputToAddress(claim.BtcAddress)
	fee, err := useCase.contracts.FlyoverConfigurations.CalculatePegInFee(amount)
	if err != nil {
		return useCase.unavailable(err)
	}
	rawTx, err := useCase.rpc.Btc.GetRawTransaction(claim.DepositTxID)
	if err != nil {
		return useCase.unavailable(err)
	}
	if err = blockchain.RejectWitnessSerializedTx(rawTx); err != nil {
		if releaseErr := useCase.releaseReserve(ctx, &claim); releaseErr != nil {
			return releaseErr
		}
		return usecases.WrapUseCaseError(usecases.SettlePegInClaimId, err)
	}
	params, err := useCase.buildRequestParams(claim.RskAddress, claim.DepositTxID, amount, fee, rawTx)
	if err != nil {
		return useCase.unavailable(err)
	}
	simulateErr := useCase.contracts.PegIn.SimulateRequestPegIn(params)
	if simulateErr == nil {
		return useCase.failRetryable(ctx, claim)
	}
	if errors.Is(simulateErr, blockchain.ErrPegInAlreadyProcessed) {
		return useCase.classifySubmitError(ctx, claim, simulateErr)
	}
	return useCase.unavailable(simulateErr)
}

func (useCase *SettlePegInClaimUseCase) classifySubmitError(
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
	return usecases.WrapUseCaseError(usecases.SettlePegInClaimId, submitErr)
}

func (useCase *SettlePegInClaimUseCase) failRetryable(
	ctx context.Context,
	claim rootstock.PegInClaim,
) error {
	claim.State = rootstock.PegInClaimRetryableFailure
	claim.ReservedWei = entities.NewWei(0)
	claim.UpdatedAt = time.Now().UTC()
	if err := useCase.claims.Update(ctx, claim); err != nil {
		return useCase.unavailable(errors.Join(ErrStatus0ReceiptStillCallable, err))
	}
	return usecases.WrapUseCaseError(usecases.SettlePegInClaimId, ErrStatus0ReceiptStillCallable)
}

func (useCase *SettlePegInClaimUseCase) buildRequestParams(
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

func (useCase *SettlePegInClaimUseCase) releaseReserve(ctx context.Context, existing *rootstock.PegInClaim) error {
	if existing == nil || existing.ReservedWei == nil || existing.ReservedWei.Cmp(entities.NewWei(0)) == 0 {
		return nil
	}
	existing.ReservedWei = entities.NewWei(0)
	existing.UpdatedAt = time.Now().UTC()
	return useCase.unavailable(useCase.claims.Update(ctx, *existing))
}

func (useCase *SettlePegInClaimUseCase) unavailable(err error) error {
	if err == nil {
		return nil
	}
	return usecases.WrapUseCaseError(usecases.SettlePegInClaimId, errors.Join(err, usecases.InfrastructureUnavailableError))
}
