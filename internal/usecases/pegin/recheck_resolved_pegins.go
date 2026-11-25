package pegin

import (
	"context"
	"errors"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

const ResolvedPegInRecheckWindow = time.Hour

type RecheckResolvedPegInsUseCase struct {
	claims    rootstock.PegInClaimRepository
	contracts blockchain.RskContracts
	rpc       blockchain.Rpc
}

func NewRecheckResolvedPegInsUseCase(
	claims rootstock.PegInClaimRepository,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
) *RecheckResolvedPegInsUseCase {
	return &RecheckResolvedPegInsUseCase{
		claims:    claims,
		contracts: contracts,
		rpc:       rpc,
	}
}

func (useCase *RecheckResolvedPegInsUseCase) Run(ctx context.Context) error {
	since := time.Now().UTC().Add(-ResolvedPegInRecheckWindow)
	claims, err := useCase.claims.ListByStateUpdatedSince(ctx, rootstock.PegInClaimResolved, since)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.RecheckResolvedPegInsId, err)
	}
	for _, claim := range claims {
		if err = useCase.recheck(ctx, claim); err != nil {
			log.Error(LogPegInResolveRecheckError(claim.RskAddress, claim.DepositTxID, err))
		}
	}
	return nil
}

func (useCase *RecheckResolvedPegInsUseCase) recheck(ctx context.Context, claim rootstock.PegInClaim) error {
	params, err := resolvePegInParams(useCase.rpc.Btc, claim)
	if err != nil {
		return err
	}
	err = useCase.contracts.PegIn.SimulateResolvePegIn(params)
	switch {
	case errors.Is(err, blockchain.PegInAlreadyProcessedError):
		return nil
	case errors.Is(err, blockchain.ContractPausedError):
		log.Debug(LogPegInResolveRecheckSkipped(claim.RskAddress, claim.DepositTxID, err))
		return nil
	case err == nil, errors.Is(err, blockchain.WaitingForBridgeError), errors.Is(err, blockchain.BridgeRejectedPegInError):
		return useCase.reopen(ctx, claim)
	default:
		return err
	}
}

func (useCase *RecheckResolvedPegInsUseCase) reopen(ctx context.Context, claim rootstock.PegInClaim) error {
	claim.State = rootstock.PegInClaimResolving
	claim.ClaimerPayout = nil
	claim.RegistrantFee = nil
	claim.UpdatedAt = time.Now().UTC()
	if err := useCase.claims.Update(ctx, claim); err != nil {
		return err
	}
	log.Warn(LogPegInResolveReopened(claim.RskAddress, claim.DepositTxID, claim.ResolveTxHash))
	return nil
}
