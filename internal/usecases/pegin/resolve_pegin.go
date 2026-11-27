package pegin

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

var PegInClaimNotResolvableError = errors.New("pegin claim can't be resolved in its state")

var errPegInResolveSkipped = errors.New("pegin resolve skipped")

type ResolvePegInUseCase struct {
	claims         rootstock.PegInClaimRepository
	contracts      blockchain.RskContracts
	rpc            blockchain.Rpc
	rskWalletMutex sync.Locker
	miningWait     time.Duration
}

func NewResolvePegInUseCase(
	claims rootstock.PegInClaimRepository,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	rskWalletMutex sync.Locker,
	miningWait time.Duration,
) *ResolvePegInUseCase {
	return &ResolvePegInUseCase{
		claims:         claims,
		contracts:      contracts,
		rpc:            rpc,
		rskWalletMutex: rskWalletMutex,
		miningWait:     miningWait,
	}
}

func (useCase *ResolvePegInUseCase) Run(ctx context.Context, claim rootstock.PegInClaim) error {
	params, err := useCase.prepareResolve(claim)
	if err == nil {
		err = useCase.submitResolve(ctx, claim, params)
	} else if errors.Is(err, blockchain.PegInAlreadyProcessedError) || errors.Is(err, blockchain.BridgeRejectedPegInError) {
		err = useCase.updateResolvedClaim(ctx, claim, blockchain.TransactionReceipt{}, err)
	}
	if err == nil || errors.Is(err, errPegInResolveSkipped) {
		return nil
	}
	return usecases.WrapUseCaseError(usecases.ResolvePegInId, err)
}

func (useCase *ResolvePegInUseCase) prepareResolve(claim rootstock.PegInClaim) (blockchain.ResolvePegInParams, error) {
	if err := useCase.checkResolvable(claim); err != nil {
		return blockchain.ResolvePegInParams{}, err
	}
	if err := useCase.checkBridgeDepth(claim); err != nil {
		return blockchain.ResolvePegInParams{}, err
	}
	params, err := useCase.resolveParams(claim)
	if err != nil {
		return blockchain.ResolvePegInParams{}, err
	}
	return params, useCase.dryRun(claim, params)
}

func (useCase *ResolvePegInUseCase) checkResolvable(claim rootstock.PegInClaim) error {
	switch claim.State {
	case rootstock.PegInClaimClaimed:
		return nil
	case rootstock.PegInClaimResolving:
		return useCase.checkMiningWait(claim)
	default:
		return fmt.Errorf(
			"%w: deposit %s for %s is %s",
			PegInClaimNotResolvableError, claim.DepositTxID, claim.RskAddress, claim.State,
		)
	}
}

func (useCase *ResolvePegInUseCase) checkMiningWait(claim rootstock.PegInClaim) error {
	if time.Since(claim.UpdatedAt) < useCase.miningWait {
		log.Debug(LogPegInResolvePending(claim.RskAddress, claim.DepositTxID, claim.ResolveTxHash))
		return errPegInResolveSkipped
	}
	return nil
}

func (useCase *ResolvePegInUseCase) checkBridgeDepth(claim rootstock.PegInClaim) error {
	tx, err := useCase.rpc.Btc.GetTransactionInfo(claim.DepositTxID)
	if err != nil {
		return err
	}
	required := useCase.contracts.Bridge.GetRequiredTxConfirmations()
	if tx.Confirmations < required {
		log.Debug(LogPegInResolveWaitingForConfirmations(claim.RskAddress, claim.DepositTxID, tx.Confirmations, required))
		return errPegInResolveSkipped
	}
	return nil
}

func (useCase *ResolvePegInUseCase) resolveParams(claim rootstock.PegInClaim) (blockchain.ResolvePegInParams, error) {
	rawTx, err := useCase.rpc.Btc.GetRawTransaction(claim.DepositTxID)
	if err != nil {
		return blockchain.ResolvePegInParams{}, err
	}
	partialMerkleTree, err := useCase.rpc.Btc.GetPartialMerkleTree(claim.DepositTxID)
	if err != nil {
		return blockchain.ResolvePegInParams{}, err
	}
	block, err := useCase.rpc.Btc.GetTransactionBlockInfo(claim.DepositTxID)
	if err != nil {
		return blockchain.ResolvePegInParams{}, err
	}
	return blockchain.ResolvePegInParams{
		RskAddress:        claim.RskAddress,
		BtcRawTx:          rawTx,
		PartialMerkleTree: partialMerkleTree,
		BlockHeight:       block.Height,
	}, nil
}

func (useCase *ResolvePegInUseCase) dryRun(claim rootstock.PegInClaim, params blockchain.ResolvePegInParams) error {
	err := useCase.contracts.PegIn.SimulateResolvePegIn(params)
	if errors.Is(err, blockchain.WaitingForBridgeError) || errors.Is(err, blockchain.ContractPausedError) {
		log.Debug(LogPegInResolveNotResolvable(claim.RskAddress, claim.DepositTxID, err))
		return errPegInResolveSkipped
	}
	return err
}

func (useCase *ResolvePegInUseCase) submitResolve(
	ctx context.Context,
	claim rootstock.PegInClaim,
	params blockchain.ResolvePegInParams,
) error {
	resending := claim.State == rootstock.PegInClaimResolving
	claim.State = rootstock.PegInClaimResolving
	claim.UpdatedAt = time.Now().UTC()
	if err := useCase.claims.Update(ctx, claim); err != nil {
		return err
	}
	if resending {
		log.Info(LogPegInResolveResent(claim.RskAddress, claim.DepositTxID, claim.ResolveTxHash))
	}
	useCase.rskWalletMutex.Lock()
	receipt, sendErr := useCase.contracts.PegIn.ResolvePegIn(params)
	useCase.rskWalletMutex.Unlock()
	return useCase.updateResolvedClaim(context.WithoutCancel(ctx), claim, receipt, sendErr)
}

func (useCase *ResolvePegInUseCase) updateResolvedClaim(
	ctx context.Context,
	claim rootstock.PegInClaim,
	receipt blockchain.TransactionReceipt,
	resolveErr error,
) error {
	var event blockchain.PegInResolvedEvent
	var eventErr error
	switch {
	case resolveErr == nil:
		claim.State = rootstock.PegInClaimResolved
		claim.ResolveTxHash = receipt.TransactionHash
		event, eventErr = useCase.contracts.PegIn.PegInResolvedEvent(receipt)
		claim.ClaimerPayout = event.ClaimerPayout
		claim.RegistrantFee = event.RegistrantFee
	case errors.Is(resolveErr, blockchain.PegInAlreadyProcessedError):
		claim.State = rootstock.PegInClaimResolved
		claim.ResolveTxHash = ""
	case errors.Is(resolveErr, blockchain.BridgeRejectedPegInError):
		claim.State = rootstock.PegInClaimResolveFailed
	default:
		claim.ResolveTxHash = cmp.Or(receipt.TransactionHash, claim.ResolveTxHash)
	}
	claim.UpdatedAt = time.Now().UTC()
	if err := useCase.claims.Update(ctx, claim); err != nil {
		return errors.Join(resolveErr, err)
	}
	switch {
	case claim.State == rootstock.PegInClaimResolveFailed:
		return errors.Join(resolveErr, usecases.NonRecoverableError)
	case claim.State == rootstock.PegInClaimResolving:
		return resolveErr
	case resolveErr == nil:
		useCase.logResolved(claim, event, eventErr)
	default:
		log.Info(LogPegInAlreadyResolved(claim.RskAddress, claim.DepositTxID))
	}
	return nil
}

func (useCase *ResolvePegInUseCase) logResolved(claim rootstock.PegInClaim, event blockchain.PegInResolvedEvent, eventErr error) {
	log.Info(LogPegInResolved(claim.RskAddress, claim.DepositTxID, claim.ResolveTxHash, claim.PegInID))
	if eventErr != nil {
		log.Error(LogPegInResolveEventUnreadable(claim.DepositTxID, claim.ResolveTxHash, eventErr))
		return
	}
	if event.PegInID != claim.PegInID {
		log.Error(LogPegInResolveIdMismatch(claim.DepositTxID, claim.PegInID, event.PegInID))
	}
}
