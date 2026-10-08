package pegin

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type ResolvePegInUseCase struct {
	pegInClaimRepository rootstock.PegInClaimRepository
	contracts            blockchain.RskContracts
	btcRpc               blockchain.BitcoinNetwork
	rskWalletMutex       sync.Locker
}

func NewResolvePegInUseCase(
	pegInClaimRepository rootstock.PegInClaimRepository,
	contracts blockchain.RskContracts,
	btcRpc blockchain.BitcoinNetwork,
	rskWalletMutex sync.Locker,
) *ResolvePegInUseCase {
	return &ResolvePegInUseCase{
		pegInClaimRepository: pegInClaimRepository,
		contracts:            contracts,
		btcRpc:               btcRpc,
		rskWalletMutex:       rskWalletMutex,
	}
}

func (useCase *ResolvePegInUseCase) Run(ctx context.Context, claim rootstock.PegInClaim) error {
	var params blockchain.ResolvePegInParams
	var err error

	if !claim.IsResolvable() {
		return useCase.wrapError(claim, errors.Join(usecases.WrongStateError, usecases.NonRecoverableError))
	}

	if err = usecases.CheckPauseLevel(useCase.contracts.PauseRegistry, blockchain.PauseLevelHard); errors.Is(err, blockchain.ContractPausedError) {
		return useCase.wrapError(claim, err)
	} else if err != nil {
		return useCase.wrapError(claim, usecases.JoinInfrastructureUnavailable(err))
	}

	if err = useCase.validateConfirmations(claim); err != nil {
		return useCase.wrapError(claim, err)
	}

	if params, err = useCase.buildResolvePegInParams(claim); err != nil {
		return useCase.wrapError(claim, usecases.JoinInfrastructureUnavailable(err))
	}

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	return useCase.performResolvePegIn(ctx, params, claim)
}

func (useCase *ResolvePegInUseCase) validateConfirmations(claim rootstock.PegInClaim) error {
	txInfo, err := useCase.btcRpc.GetTransactionInfo(claim.DepositTxID)
	if err != nil {
		return usecases.JoinInfrastructureUnavailable(err)
	}
	if txInfo.Confirmations < useCase.contracts.Bridge.GetRequiredTxConfirmations() {
		return usecases.NoEnoughConfirmationsError
	}
	return nil
}

func (useCase *ResolvePegInUseCase) buildResolvePegInParams(claim rootstock.PegInClaim) (blockchain.ResolvePegInParams, error) {
	var rawBtcTx, pmt []byte
	var block blockchain.BitcoinBlockInformation
	var err error

	if rawBtcTx, err = useCase.btcRpc.GetRawTransaction(claim.DepositTxID); err != nil {
		return blockchain.ResolvePegInParams{}, err
	}

	if pmt, err = useCase.btcRpc.GetPartialMerkleTree(claim.DepositTxID); err != nil {
		return blockchain.ResolvePegInParams{}, err
	}

	if block, err = useCase.btcRpc.GetTransactionBlockInfo(claim.DepositTxID); err != nil {
		return blockchain.ResolvePegInParams{}, err
	}

	return blockchain.ResolvePegInParams{
		RskAddress:            claim.RskAddress,
		BitcoinRawTransaction: rawBtcTx,
		PartialMerkleTree:     pmt,
		BlockHeight:           block.Height,
	}, nil
}

func (useCase *ResolvePegInUseCase) performResolvePegIn(
	ctx context.Context,
	params blockchain.ResolvePegInParams,
	claim rootstock.PegInClaim,
) error {
	result, resolveErr := useCase.contracts.PegIn.ResolvePegIn(params)
	claim, write := useCase.applyResolveResult(claim, result, resolveErr)
	if !write {
		return useCase.wrapError(claim, resolveErr)
	}

	claim.UpdatedAt = time.Now().UTC()
	if err := useCase.pegInClaimRepository.Update(ctx, claim); err != nil {
		return useCase.wrapError(claim, errors.Join(resolveErr, usecases.JoinInfrastructureUnavailable(err)))
	}
	switch claim.State {
	case rootstock.PegInClaimResolved:
		return nil
	case rootstock.PegInClaimResolveFailed:
		return useCase.wrapError(claim, errors.Join(resolveErr, usecases.NonRecoverableError))
	default:
		return useCase.wrapError(claim, resolveErr)
	}
}

// applyResolveResult logs the outcome and reports whether the claim must be written
func (useCase *ResolvePegInUseCase) applyResolveResult(
	claim rootstock.PegInClaim,
	result blockchain.ResolvePegInResult,
	resolveErr error,
) (rootstock.PegInClaim, bool) {
	sent := result.Receipt.TransactionHash != ""
	if sent {
		claim.ResolveTxHash = result.Receipt.TransactionHash
		claim.ResolveGasUsed = result.Receipt.GasUsed.Uint64()
		claim.ResolveGasPrice = result.Receipt.GasPrice
	}

	switch {
	case resolveErr == nil:
		claim.State = rootstock.PegInClaimResolved
		claim.ClaimerPayout = result.ClaimerPayout
		claim.RegistrantFee = result.RegistrantFee
		log.Info(LogPegInResolved(claim))
	case errors.Is(resolveErr, blockchain.ErrPegInAlreadyProcessed):
		claim.State = rootstock.PegInClaimResolved
		log.Info(LogPegInResolveAlreadyProcessed(claim))
	case errors.Is(resolveErr, blockchain.WaitingForBridgeError):
		log.Warn(LogPegInResolveWaitingForBridge(claim))
		return claim, false
	case errors.Is(resolveErr, blockchain.ErrBridgeRejectedPegIn), errors.Is(resolveErr, blockchain.ErrPegInNotClaimed):
		claim.State = rootstock.PegInClaimResolveFailed
		log.Error(LogPegInResolveFailed(claim, resolveErr))
	case sent:
		log.Error(LogPegInResolveTxNotResolved(claim, resolveErr))
	default:
		return claim, false
	}
	return claim, true
}

func (useCase *ResolvePegInUseCase) wrapError(claim rootstock.PegInClaim, err error) error {
	args := usecases.NewErrorArgs()
	args["rskAddress"] = claim.RskAddress
	args["depositTxId"] = claim.DepositTxID
	return usecases.WrapUseCaseErrorArgs(usecases.ResolvePegInId, err, args)
}
