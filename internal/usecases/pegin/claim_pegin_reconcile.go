package pegin

import (
	"context"
	"errors"
	"strings"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	log "github.com/sirupsen/logrus"
)

func (useCase *ClaimPegInUseCase) reconcile(ctx context.Context, claim rootstock.PegInClaim) error {
	if claim.TxHash == "" {
		logUnrecoverableClaimReceipt(claim)
		return nil
	}
	receipt, err := useCase.rpc.Rsk.GetTransactionReceipt(ctx, claim.TxHash)
	if err != nil {
		if errors.Is(err, blockchain.ErrTransactionReceiptNotFound) {
			logUnrecoverableClaimReceipt(claim)
			return nil
		}
		return useCase.wrap(err)
	}
	if !useCase.receiptCanonical(ctx, receipt) {
		return nil
	}
	if receipt.Status == 0 {
		return useCase.identifyFailed(ctx, claim)
	}
	if receipt.Status != blockchain.SuccessfulTxStatus {
		return nil
	}
	event, unpackErr := useCase.contracts.PegIn.UnpackPegInRequested(receipt)
	if !hasPegInRequested(event, unpackErr) {
		log.Errorf(
			"ClaimPegIn: status-1 receipt %s for %s/%s is missing PegInRequested; follow incident-recovery; not resubmitting: %v",
			claim.TxHash,
			claim.RskAddress,
			claim.DepositTxID,
			unpackErr,
		)
		return nil
	}
	return useCase.finalizeSuccess(ctx, claim, blockchain.RequestPegInResult{
		Receipt: receipt,
		Event:   event,
	})
}

func (useCase *ClaimPegInUseCase) receiptCanonical(ctx context.Context, receipt blockchain.TransactionReceipt) bool {
	for _, eventLog := range receipt.Logs {
		if eventLog.Removed {
			return false
		}
	}
	if receipt.BlockHash == "" {
		return false
	}
	block, err := useCase.rpc.Rsk.GetBlockByHash(ctx, receipt.BlockHash)
	if err != nil {
		return false
	}
	return strings.EqualFold(block.Hash, receipt.BlockHash)
}

func (useCase *ClaimPegInUseCase) identifyFailed(ctx context.Context, claim rootstock.PegInClaim) error {
	tx, err := useCase.rpc.Btc.GetTransactionInfo(claim.DepositTxID)
	if err != nil {
		return useCase.wrap(err)
	}
	amount := tx.FirstOutputToAddress(claim.BtcAddress)
	fee, err := useCase.contracts.FlyoverConfigurations.CalculatePegInFee(amount)
	if err != nil {
		return useCase.wrap(err)
	}
	params, err := useCase.buildRequestParams(claim.RskAddress, claim.DepositTxID, amount, fee)
	if err != nil {
		return useCase.wrap(err)
	}
	identifyErr := useCase.contracts.PegIn.IdentifyRequestPegIn(params)
	if identifyErr == nil {
		return useCase.failRetryable(ctx, claim, errors.New("status-0 receipt; preflight no longer reverts; not resubmitting"))
	}
	return useCase.classifySubmitError(ctx, claim, identifyErr)
}

func logUnrecoverableClaimReceipt(claim rootstock.PegInClaim) {
	log.Errorf(
		"ClaimPegIn: submitting tx %s for %s/%s has no recoverable receipt; follow incident-recovery; not resubmitting",
		claim.TxHash,
		claim.RskAddress,
		claim.DepositTxID,
	)
}

func hasPegInRequested(event blockchain.PegInRequestedEvent, err error) bool {
	return err == nil
}
