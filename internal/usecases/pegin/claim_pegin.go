package pegin

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	pegin "github.com/rsksmart/liquidity-provider-server/internal/entities/pegin"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

// ClaimPegInUseCase implements the LP side of the commit-first peg-in (DoS-removal redesign,
// EPICs E4/E5). It does NOT reserve liquidity ahead of a user commitment: it acts only after a
// confirmed BTC deposit to a registered address is observed, checking liquidity against the LP's
// own wallet at claim time and fronting RBTC competitively.
//
// The flow has two on-chain steps, mirroring the legacy callForUser/registerPegIn split:
//   - Request: front RBTC (amount - protocol fee) via PegInContract.requestPegIn. First-mined-wins,
//     so a revert because another LP already claimed is reported as benign (ClaimedByOther).
//   - Resolve: once the BTC tx reaches the required confirmations, settle with the Bridge via
//     PegInContract.resolvePegIn, recovering the fronted RBTC plus the fee.
type ClaimPegInUseCase struct {
	contracts      blockchain.RskContracts
	rpc            blockchain.Rpc
	lp             liquidity_provider.LiquidityProvider
	peginLp        liquidity_provider.PeginLiquidityProvider
	rskWalletMutex sync.Locker
}

func NewClaimPegInUseCase(
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	lp liquidity_provider.LiquidityProvider,
	peginLp liquidity_provider.PeginLiquidityProvider,
	rskWalletMutex sync.Locker,
) *ClaimPegInUseCase {
	return &ClaimPegInUseCase{
		contracts:      contracts,
		rpc:            rpc,
		lp:             lp,
		peginLp:        peginLp,
		rskWalletMutex: rskWalletMutex,
	}
}

// Request fronts RBTC to claim a confirmed deposit. It returns the next state for the watched
// address: Requested on a won claim, ClaimedByOther when another LP got there first.
func (useCase *ClaimPegInUseCase) Request(ctx context.Context, deposit pegin.ConfirmedDeposit) (pegin.WatchedAddressState, error) {
	if useCase.contracts.FlyoverConfigurations == nil || useCase.contracts.PegInAddressRegistry == nil {
		return pegin.WatchedAddressStateWaitingForDeposit,
			usecases.WrapUseCaseError(usecases.RequestPegInId, errors.New("commit-first peg-in contracts are not configured"))
	}
	if err := usecases.CheckPauseState(useCase.contracts.PegIn); err != nil {
		return pegin.WatchedAddressStateWaitingForDeposit, usecases.WrapUseCaseError(usecases.RequestPegInId, err)
	}

	// The address must still be registered (it could have been affected by a federation change).
	registered, err := useCase.contracts.PegInAddressRegistry.IsRegistered(deposit.RskAddress)
	if err != nil {
		return pegin.WatchedAddressStateWaitingForDeposit, usecases.WrapUseCaseError(usecases.RequestPegInId, err)
	}
	if !registered {
		return pegin.WatchedAddressStateWaitingForDeposit,
			usecases.WrapUseCaseError(usecases.RequestPegInId, blockchain.AddressNotRegisteredError)
	}

	fee, err := useCase.contracts.FlyoverConfigurations.CalculatePegInFee(deposit.Amount)
	if err != nil {
		return pegin.WatchedAddressStateWaitingForDeposit, usecases.WrapUseCaseError(usecases.RequestPegInId, err)
	}
	// The LP fronts amount - fee from its own wallet.
	valueToFront := new(entities.Wei).Sub(deposit.Amount, fee)
	if valueToFront.Cmp(entities.NewWei(0)) <= 0 {
		return pegin.WatchedAddressStateWaitingForDeposit,
			usecases.WrapUseCaseError(usecases.RequestPegInId, fmt.Errorf("deposit amount %s does not cover the peg-in fee %s", deposit.Amount.String(), fee.String()))
	}

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	// Liquidity is checked against the LP's own wallet AT CLAIM TIME — never reserved earlier.
	if err = useCase.checkOwnLiquidity(ctx, valueToFront); err != nil {
		return pegin.WatchedAddressStateWaitingForDeposit, usecases.WrapUseCaseError(usecases.RequestPegInId, err)
	}

	params := blockchain.RequestPegInParams{
		RskAddress: deposit.RskAddress,
		Amount:     deposit.Amount,
		BtcTxHash:  deposit.BtcTxHashRaw,
		OpReturn:   deposit.OpReturn,
		Value:      valueToFront,
	}
	receipt, err := useCase.contracts.PegIn.RequestPegIn(params)
	if errors.Is(err, blockchain.AlreadyClaimedError) {
		log.Infof("ClaimPegIn: peg-in for %s (btcTx %s) already claimed by another provider", deposit.RskAddress, deposit.BtcTxHash)
		return pegin.WatchedAddressStateClaimedByOther, nil
	}
	if err != nil {
		return pegin.WatchedAddressStateWaitingForDeposit,
			usecases.WrapUseCaseErrorArgs(usecases.RequestPegInId, err, usecases.ErrorArg("rskAddress", deposit.RskAddress))
	}
	log.Infof("ClaimPegIn: requested peg-in for %s, fronted %s wei (tx %s)", deposit.RskAddress, valueToFront.String(), receipt.TransactionHash)
	return pegin.WatchedAddressStateRequested, nil
}

// Resolve settles a previously requested peg-in with the Bridge. WaitingForBridgeError is returned
// (wrapped) when there are not yet enough confirmations so the watcher can retry; AlreadyClaimedError
// maps to ClaimedByOther.
func (useCase *ClaimPegInUseCase) Resolve(ctx context.Context, deposit pegin.ConfirmedDeposit) (pegin.WatchedAddressState, error) {
	if useCase.contracts.FlyoverConfigurations == nil {
		return pegin.WatchedAddressStateRequested,
			usecases.WrapUseCaseError(usecases.ResolvePegInId, errors.New("commit-first peg-in contracts are not configured"))
	}
	if err := usecases.CheckPauseState(useCase.contracts.PegIn); err != nil {
		return pegin.WatchedAddressStateRequested, usecases.WrapUseCaseError(usecases.ResolvePegInId, err)
	}

	requiredConfirmations, err := useCase.contracts.FlyoverConfigurations.GetRequiredPegInConfirmations(deposit.Amount)
	if err != nil {
		return pegin.WatchedAddressStateRequested, usecases.WrapUseCaseError(usecases.ResolvePegInId, err)
	}
	if deposit.Confirmations < requiredConfirmations {
		return pegin.WatchedAddressStateRequested,
			usecases.WrapUseCaseError(usecases.ResolvePegInId, blockchain.WaitingForBridgeError)
	}

	rawTx, err := useCase.rpc.Btc.GetRawTransaction(deposit.BtcTxHash)
	if err != nil {
		return pegin.WatchedAddressStateRequested, usecases.WrapUseCaseError(usecases.ResolvePegInId, err)
	}
	pmt, err := useCase.rpc.Btc.GetPartialMerkleTree(deposit.BtcTxHash)
	if err != nil {
		return pegin.WatchedAddressStateRequested, usecases.WrapUseCaseError(usecases.ResolvePegInId, err)
	}
	blockInfo, err := useCase.rpc.Btc.GetTransactionBlockInfo(deposit.BtcTxHash)
	if err != nil {
		return pegin.WatchedAddressStateRequested, usecases.WrapUseCaseError(usecases.ResolvePegInId, err)
	}

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	params := blockchain.ResolvePegInParams{
		RskAddress:        deposit.RskAddress,
		BtcTxHash:         deposit.BtcTxHashRaw,
		BtcRawTransaction: rawTx,
		PartialMerkleTree: pmt,
		Height:            new(big.Int).Set(blockInfo.Height),
		// The registrant (e.g. a watchtower) is paid its fee from the contract; for the PoC the
		// claiming LP is the registrant.
		Registrant: useCase.lp.RskAddress(),
	}
	receipt, err := useCase.contracts.PegIn.ResolvePegIn(params)
	if errors.Is(err, blockchain.WaitingForBridgeError) {
		return pegin.WatchedAddressStateRequested, usecases.WrapUseCaseError(usecases.ResolvePegInId, err)
	}
	if errors.Is(err, blockchain.AlreadyClaimedError) {
		return pegin.WatchedAddressStateClaimedByOther, nil
	}
	if err != nil {
		return pegin.WatchedAddressStateRequested,
			usecases.WrapUseCaseErrorArgs(usecases.ResolvePegInId, err, usecases.ErrorArg("rskAddress", deposit.RskAddress))
	}
	log.Infof("ClaimPegIn: resolved peg-in for %s (tx %s)", deposit.RskAddress, receipt.TransactionHash)
	return pegin.WatchedAddressStateResolved, nil
}

// checkOwnLiquidity verifies the LP's network balance covers the value to front, with no reservation.
func (useCase *ClaimPegInUseCase) checkOwnLiquidity(ctx context.Context, valueToFront *entities.Wei) error {
	balance, err := useCase.rpc.Rsk.GetBalance(ctx, useCase.lp.RskAddress())
	if err != nil {
		return err
	}
	if balance.Cmp(valueToFront) < 0 {
		return fmt.Errorf("%w: have %s wei, need %s wei to front", usecases.NoLiquidityError, balance.String(), valueToFront.String())
	}
	return nil
}
