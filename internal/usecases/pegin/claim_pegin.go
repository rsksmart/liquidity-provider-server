package pegin

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

var DepositOutputNotFoundError = errors.New("the deposit has no output paying the watched address")

var errPegInClaimSkipped = errors.New("pegin claim skipped")

var skippedDryRunErrors = []error{
	blockchain.PegInBelowMinimumError,
	blockchain.DepositBelowBridgeMinimumError,
	blockchain.InsufficientConfirmationsError,
	blockchain.UnableToGetConfirmationsError,
	blockchain.ContractPausedError,
}

type ClaimPegInUseCase struct {
	claims              rootstock.PegInClaimRepository
	contracts           blockchain.RskContracts
	rpc                 blockchain.Rpc
	rskWallet           blockchain.RootstockWallet
	hashFunction        entities.HashFunction
	peginLiquidityMutex sync.Locker
	rskWalletMutex      sync.Locker
}

func NewClaimPegInUseCase(
	claims rootstock.PegInClaimRepository,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	rskWallet blockchain.RootstockWallet,
	hashFunction entities.HashFunction,
	peginLiquidityMutex sync.Locker,
	rskWalletMutex sync.Locker,
) *ClaimPegInUseCase {
	return &ClaimPegInUseCase{
		claims:              claims,
		contracts:           contracts,
		rpc:                 rpc,
		rskWallet:           rskWallet,
		hashFunction:        hashFunction,
		peginLiquidityMutex: peginLiquidityMutex,
		rskWalletMutex:      rskWalletMutex,
	}
}

type pegInClaimAmounts struct {
	amount *entities.Wei
	fee    *entities.Wei
	gasFee *entities.Wei
	value  *entities.Wei
}

type pegInClaimRequest struct {
	pegInClaimAmounts
	watch       rootstock.PegInWatch
	depositTxID string
	pegInID     string
	params      blockchain.RequestPegInParams
	createdAt   time.Time
}

func (request pegInClaimRequest) newClaim(state rootstock.PegInClaimState) rootstock.PegInClaim {
	return rootstock.PegInClaim{
		RskAddress:  request.watch.RskAddress,
		DepositTxID: request.depositTxID,
		BtcAddress:  request.watch.BtcAddress,
		PegInID:     request.pegInID,
		State:       state,
		Amount:      request.amount,
		Fee:         request.fee,
		GasFee:      request.gasFee,
		Value:       request.value,
		CreatedAt:   request.createdAt,
		UpdatedAt:   time.Now().UTC(),
	}
}

func (useCase *ClaimPegInUseCase) Run(ctx context.Context, watch rootstock.PegInWatch, depositTxID string) error {
	request, err := useCase.prepareClaim(ctx, watch, depositTxID)
	if errors.Is(err, blockchain.PegInAlreadyProcessedError) {
		err = useCase.addRaceLostClaim(ctx, request)
	} else if err == nil {
		err = useCase.submitClaim(ctx, request)
	}
	if err == nil || errors.Is(err, errPegInClaimSkipped) {
		return nil
	}
	return usecases.WrapUseCaseError(usecases.ClaimPegInId, err)
}

func (useCase *ClaimPegInUseCase) prepareClaim(
	ctx context.Context,
	watch rootstock.PegInWatch,
	depositTxID string,
) (pegInClaimRequest, error) {
	request := pegInClaimRequest{watch: watch, depositTxID: depositTxID}
	var err error
	if request.createdAt, err = useCase.checkStoredClaim(ctx, watch.RskAddress, depositTxID); err != nil {
		return pegInClaimRequest{}, err
	}
	amount, err := useCase.confirmedDepositAmount(watch, depositTxID)
	if err != nil {
		return pegInClaimRequest{}, err
	}
	if err = useCase.checkNotHardPaused(watch.RskAddress, depositTxID); err != nil {
		return pegInClaimRequest{}, err
	}
	rawTx, err := useCase.depositTransaction(depositTxID)
	if err != nil {
		return pegInClaimRequest{}, err
	}
	if request.pegInClaimAmounts, err = useCase.claimAmounts(watch.RskAddress, depositTxID, amount); err != nil {
		return pegInClaimRequest{}, err
	}
	if request.params, err = useCase.requestParams(watch.RskAddress, depositTxID, rawTx); err != nil {
		return pegInClaimRequest{}, err
	}
	if request.pegInID, err = useCase.pegInID(watch.RskAddress, depositTxID); err != nil {
		return pegInClaimRequest{}, err
	}
	return request, useCase.dryRun(request)
}

func (useCase *ClaimPegInUseCase) checkStoredClaim(
	ctx context.Context,
	rskAddress, depositTxID string,
) (createdAt time.Time, err error) {
	stored, err := useCase.claims.Get(ctx, rskAddress, depositTxID)
	if err != nil {
		return time.Time{}, err
	}
	if stored == nil {
		return time.Now().UTC(), nil
	}
	if stored.State == rootstock.PegInClaimRetryableFailure {
		return stored.CreatedAt, nil
	}
	log.Debug(LogPegInClaimAlreadyStored(rskAddress, depositTxID, stored.State))
	return time.Time{}, errPegInClaimSkipped
}

func (useCase *ClaimPegInUseCase) confirmedDepositAmount(watch rootstock.PegInWatch, depositTxID string) (*entities.Wei, error) {
	tx, err := useCase.rpc.Btc.GetTransactionInfo(depositTxID)
	if err != nil {
		return nil, err
	}
	outputs := tx.UTXOsToAddress(watch.BtcAddress)
	if len(outputs) == 0 {
		return nil, fmt.Errorf("%w: %s", DepositOutputNotFoundError, depositTxID)
	}
	amount := outputs[0]
	required, err := useCase.contracts.FlyoverConfigurations.GetRequiredPegInBtcConfirmations(amount)
	if err != nil {
		return nil, err
	}
	if tx.Confirmations < required {
		log.Debug(LogPegInClaimWaitingForConfirmations(watch.RskAddress, depositTxID, tx.Confirmations, required))
		return nil, errPegInClaimSkipped
	}
	return amount, nil
}

func (useCase *ClaimPegInUseCase) checkNotHardPaused(rskAddress, depositTxID string) error {
	level, err := useCase.contracts.PauseRegistry.PauseLevel()
	if err != nil {
		return err
	}
	if level == blockchain.PauseLevelHard {
		log.Debug(LogPegInClaimHardPaused(rskAddress, depositTxID))
		return errPegInClaimSkipped
	}
	return nil
}

func (useCase *ClaimPegInUseCase) depositTransaction(depositTxID string) ([]byte, error) {
	rawTx, err := useCase.rpc.Btc.GetRawTransaction(depositTxID)
	if err != nil {
		return nil, err
	}
	if len(rawTx) > MaxPeginDepositTxSize {
		return nil, fmt.Errorf("%w: deposit %s has %d bytes", DataCapExceededError, depositTxID, len(rawTx))
	}
	return rawTx, nil
}

func (useCase *ClaimPegInUseCase) claimAmounts(
	rskAddress, depositTxID string,
	amount *entities.Wei,
) (pegInClaimAmounts, error) {
	configuration, err := useCase.contracts.FlyoverConfigurations.GetPegInConfiguration()
	if err != nil {
		return pegInClaimAmounts{}, err
	}
	if amount.Cmp(configuration.MinAmount) < 0 {
		log.Debug(LogPegInClaimBelowMinimum(rskAddress, depositTxID, amount, configuration.MinAmount))
		return pegInClaimAmounts{}, errPegInClaimSkipped
	}
	cappedAmount := amount
	if amount.Cmp(configuration.MaxAmount) > 0 {
		cappedAmount = configuration.MaxAmount
	}
	fee, err := useCase.contracts.FlyoverConfigurations.CalculatePegInFee(cappedAmount)
	if err != nil {
		return pegInClaimAmounts{}, err
	}
	deduction := new(entities.Wei).Add(fee, configuration.ClaimGasFee)
	if cappedAmount.Cmp(deduction) < 0 {
		log.Debug(LogPegInClaimDeductionAboveAmount(rskAddress, depositTxID, cappedAmount, deduction))
		return pegInClaimAmounts{}, errPegInClaimSkipped
	}
	return pegInClaimAmounts{
		amount: amount,
		fee:    fee,
		gasFee: configuration.ClaimGasFee,
		value:  new(entities.Wei).Sub(cappedAmount, deduction),
	}, nil
}

func (useCase *ClaimPegInUseCase) requestParams(
	rskAddress, depositTxID string,
	rawTx []byte,
) (blockchain.RequestPegInParams, error) {
	block, err := useCase.rpc.Btc.GetTransactionBlockInfo(depositTxID)
	if err != nil {
		return blockchain.RequestPegInParams{}, err
	}
	merkleBranch, err := useCase.rpc.Btc.BuildMerkleBranch(depositTxID)
	if err != nil {
		return blockchain.RequestPegInParams{}, err
	}
	return blockchain.RequestPegInParams{
		RskAddress:         rskAddress,
		BtcRawTx:           rawTx,
		BtcBlockHeaderHash: block.Hash,
		MerkleBranchPath:   merkleBranch.Path,
		MerkleBranchHashes: merkleBranch.Hashes,
	}, nil
}

func (useCase *ClaimPegInUseCase) pegInID(rskAddress, depositTxID string) (string, error) {
	rskAddressBytes, err := blockchain.DecodeStringTrimPrefix(rskAddress)
	if err != nil {
		return "", err
	}
	depositTxIDBytes, err := hex.DecodeString(depositTxID)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(useCase.hashFunction(rskAddressBytes, depositTxIDBytes)), nil
}

func (useCase *ClaimPegInUseCase) dryRun(request pegInClaimRequest) error {
	err := useCase.contracts.PegIn.SimulateRequestPegIn(
		blockchain.NewTransactionConfig(request.value, 0, nil),
		request.params,
	)
	for _, skipped := range skippedDryRunErrors {
		if errors.Is(err, skipped) {
			log.Debug(LogPegInClaimNotClaimable(request.watch.RskAddress, request.depositTxID, err))
			return errPegInClaimSkipped
		}
	}
	return err
}

func (useCase *ClaimPegInUseCase) addRaceLostClaim(ctx context.Context, request pegInClaimRequest) error {
	log.Info(LogPegInClaimRaceLost(request.watch.RskAddress, request.depositTxID))
	err := useCase.claims.Add(ctx, request.newClaim(rootstock.PegInClaimRaceLost))
	if errors.Is(err, rootstock.PegInClaimExistsError) {
		return nil
	}
	return err
}

func (useCase *ClaimPegInUseCase) submitClaim(ctx context.Context, request pegInClaimRequest) error {
	useCase.peginLiquidityMutex.Lock()
	defer useCase.peginLiquidityMutex.Unlock()

	txConfig, err := useCase.fundedTransactionConfig(ctx, request)
	if errors.Is(err, blockchain.PegInAlreadyProcessedError) {
		return useCase.addRaceLostClaim(ctx, request)
	} else if err != nil {
		return err
	}
	claim := request.newClaim(rootstock.PegInClaimSubmitting)
	if err = useCase.claims.Add(ctx, claim); errors.Is(err, rootstock.PegInClaimExistsError) {
		log.Debug(LogPegInClaimStoredMeanwhile(claim.RskAddress, claim.DepositTxID))
		return errPegInClaimSkipped
	} else if err != nil {
		return err
	}

	useCase.rskWalletMutex.Lock()
	result, sendErr := useCase.contracts.PegIn.RequestPegIn(txConfig, request.params)
	useCase.rskWalletMutex.Unlock()
	return useCase.updateSentClaim(ctx, claim, result, sendErr)
}

func (useCase *ClaimPegInUseCase) fundedTransactionConfig(
	ctx context.Context,
	request pegInClaimRequest,
) (blockchain.TransactionConfig, error) {
	gasLimit, err := useCase.contracts.PegIn.EstimateRequestPegInGas(
		blockchain.NewTransactionConfig(request.value, 0, nil),
		request.params,
	)
	if err != nil {
		return blockchain.TransactionConfig{}, err
	}
	gasPrice, err := useCase.rpc.Rsk.GasPrice(ctx)
	if err != nil {
		return blockchain.TransactionConfig{}, err
	}
	balance, err := useCase.rskWallet.GetBalance(ctx)
	if err != nil {
		return blockchain.TransactionConfig{}, err
	}
	gasCost := new(entities.Wei).Mul(gasPrice, entities.NewUWei(gasLimit))
	required := new(entities.Wei).Add(request.value, gasCost)
	if balance.Cmp(required) < 0 {
		log.Warn(LogPegInClaimInsufficientBalance(request.watch.RskAddress, request.depositTxID, balance, required))
		return blockchain.TransactionConfig{}, errPegInClaimSkipped
	}
	return blockchain.NewTransactionConfig(request.value, gasLimit, gasPrice), nil
}

func (useCase *ClaimPegInUseCase) updateSentClaim(
	ctx context.Context,
	claim rootstock.PegInClaim,
	result blockchain.RequestPegInResult,
	sendErr error,
) error {
	claim.TxHash = result.Receipt.TransactionHash
	claim.UpdatedAt = time.Now().UTC()
	switch {
	case sendErr == nil:
		claim.State = rootstock.PegInClaimClaimed
		claim.CallSuccess = result.Event.CallSuccess
	case errors.Is(sendErr, blockchain.PegInAlreadyProcessedError):
		claim.State = rootstock.PegInClaimRaceLost
	case claim.TxHash == "" || errors.Is(sendErr, blockchain.TxFailedError):
		claim.State = rootstock.PegInClaimRetryableFailure
		claim.TxHash = ""
	}
	if err := useCase.claims.Update(ctx, claim); err != nil {
		return errors.Join(sendErr, err)
	}
	switch claim.State {
	case rootstock.PegInClaimClaimed:
		useCase.logClaimed(claim, result.Event)
		return nil
	case rootstock.PegInClaimRaceLost:
		log.Info(LogPegInClaimRaceLost(claim.RskAddress, claim.DepositTxID))
		return nil
	default:
		return sendErr
	}
}

func (useCase *ClaimPegInUseCase) logClaimed(claim rootstock.PegInClaim, event blockchain.PegInRequestedEvent) {
	log.Info(LogPegInClaimed(claim.RskAddress, claim.DepositTxID, claim.TxHash, claim.PegInID))
	if event.PegInID != claim.PegInID {
		log.Error(LogPegInClaimIdMismatch(claim.DepositTxID, claim.PegInID, event.PegInID))
	}
	if !event.CallSuccess {
		log.Warn(LogPegInClaimUserPaymentFailed(claim.RskAddress, claim.DepositTxID, claim.TxHash))
	}
}
