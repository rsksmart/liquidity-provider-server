package pegin

import (
	"context"
	"encoding/hex"
	"errors"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type SettlePegInClaimUseCase struct {
	claims        rootstock.PegInClaimRepository
	contracts     blockchain.RskContracts
	rpc           blockchain.Rpc
	maxReorgDepth uint64
}

func NewSettlePegInClaimUseCase(
	claims rootstock.PegInClaimRepository,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	maxReorgDepth uint64,
) *SettlePegInClaimUseCase {
	return &SettlePegInClaimUseCase{
		claims:        claims,
		contracts:     contracts,
		rpc:           rpc,
		maxReorgDepth: maxReorgDepth,
	}
}

func (useCase *SettlePegInClaimUseCase) Run(ctx context.Context, claim rootstock.PegInClaim) error {
	if claim.TxHash == "" {
		log.Error(LogPegInClaimEmptyTxHash(claim.RskAddress, claim.DepositTxID))
		return nil
	}
	receipt, err := useCase.rpc.Rsk.GetTransactionReceipt(ctx, claim.TxHash)
	if errors.Is(err, blockchain.ErrTransactionReceiptNotFound) {
		log.Error(LogPegInClaimMissingReceipt(claim.TxHash, claim.RskAddress, claim.DepositTxID))
		return nil
	}
	if err != nil {
		return useCase.unavailable(err)
	}
	if !receiptOnCanonicalChain(receipt) {
		return nil
	}
	height, err := useCase.rpc.Rsk.GetHeight(ctx)
	if err != nil {
		return useCase.unavailable(err)
	}
	if !receipt.IsFinal(height, useCase.maxReorgDepth) {
		return nil
	}
	if receipt.Status != blockchain.SuccessfulTxStatus {
		return useCase.identifyFailed(ctx, claim)
	}
	return useCase.finalizeSuccess(ctx, claim, receipt)
}

// receiptOnCanonicalChain needs no block lookup: rskj eth_getTransactionReceipt returns a receipt
// only when its block is the main-chain block at that height (ReceiptStore.getInMainChain).
func receiptOnCanonicalChain(receipt blockchain.TransactionReceipt) bool {
	for _, eventLog := range receipt.Logs {
		if eventLog.Removed {
			return false
		}
	}
	return receipt.BlockHash != "" && receipt.BlockNumber != 0
}

func (useCase *SettlePegInClaimUseCase) finalizeSuccess(
	ctx context.Context,
	claim rootstock.PegInClaim,
	receipt blockchain.TransactionReceipt,
) error {
	event, unpackErr := useCase.contracts.PegIn.UnpackPegInRequested(receipt)
	if unpackErr != nil {
		// Keep a PegInID the submit path already stored.
		log.Error(LogPegInClaimMissingEvent(claim.TxHash, claim.RskAddress, claim.DepositTxID, unpackErr))
	} else {
		claim.PegInID = hex.EncodeToString(event.PegInId[:])
	}
	claim.State = rootstock.PegInClaimClaimed
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
		// Unreachable with the current adapters: GetRawTransaction serializes without witness.
		return usecases.WrapUseCaseError(usecases.SettlePegInClaimId, errors.Join(err, usecases.NonRecoverableError))
	}
	params, err := usecases.BuildRequestPegInParams(useCase.rpc.Btc, usecases.RequestPegInInput{
		RskAddress:  claim.RskAddress,
		DepositTxID: claim.DepositTxID,
		Amount:      amount,
		Fee:         fee,
		RawTx:       rawTx,
	})
	if err != nil {
		return useCase.unavailable(err)
	}
	simulateErr := useCase.contracts.PegIn.SimulateRequestPegIn(params)
	if simulateErr == nil {
		return useCase.failRetryable(ctx, claim)
	}
	return useCase.classifySubmitError(ctx, claim, simulateErr)
}

func (useCase *SettlePegInClaimUseCase) classifySubmitError(
	ctx context.Context,
	claim rootstock.PegInClaim,
	submitErr error,
) error {
	if errors.Is(submitErr, blockchain.ErrPegInAlreadyProcessed) {
		claim.State = rootstock.PegInClaimRaceLost
		claim.UpdatedAt = time.Now().UTC()
		if err := useCase.claims.Update(ctx, claim); err != nil {
			return useCase.unavailable(err)
		}
		return nil
	}
	if errors.Is(submitErr, blockchain.ErrAddressNotRegistered) ||
		errors.Is(submitErr, blockchain.ErrDepositOutputNotFound) ||
		errors.Is(submitErr, blockchain.ErrInsufficientConfirmations) ||
		errors.Is(submitErr, blockchain.ErrIncorrectFronting) ||
		errors.Is(submitErr, blockchain.ErrPegInBelowMinimum) {
		claim.State = rootstock.PegInClaimRetryableFailure
		claim.TxHash = ""
		claim.UpdatedAt = time.Now().UTC()
		if err := useCase.claims.Update(ctx, claim); err != nil {
			return useCase.unavailable(errors.Join(submitErr, err))
		}
		return usecases.WrapUseCaseError(usecases.SettlePegInClaimId, errors.Join(submitErr, usecases.NonRecoverableError))
	}
	return useCase.unavailable(submitErr)
}

func (useCase *SettlePegInClaimUseCase) failRetryable(
	ctx context.Context,
	claim rootstock.PegInClaim,
) error {
	claim.State = rootstock.PegInClaimRetryableFailure
	claim.TxHash = ""
	claim.UpdatedAt = time.Now().UTC()
	if err := useCase.claims.Update(ctx, claim); err != nil {
		return useCase.unavailable(errors.Join(blockchain.TxFailedError, err))
	}
	return usecases.WrapUseCaseError(usecases.SettlePegInClaimId, errors.Join(blockchain.TxFailedError, usecases.NonRecoverableError))
}

func (useCase *SettlePegInClaimUseCase) unavailable(err error) error {
	if err == nil {
		return nil
	}
	return usecases.WrapUseCaseError(usecases.SettlePegInClaimId, errors.Join(err, usecases.InfrastructureUnavailableError))
}
