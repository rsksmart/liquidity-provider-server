package liquidity_provider

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/cold_wallet"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

const (
	SimpleTransferGasLimit = 21000
)

var (
	NoColdWalletConfiguredError      = errors.New("cold wallet not configured")
	NoMaxLiquidityConfiguredError    = errors.New("max liquidity not configured")
	TransferNotEconomicalError       = errors.New("transfer not economical - amount is less than minimum fee multiplier")
	NoTransferHistoryConfiguredError = errors.New("no transfer history configured - state configuration must be initialized")
)

type TransferStatus string

const (
	TransferStatusSuccess              TransferStatus = "success"
	TransferStatusSkippedNoExcess      TransferStatus = "skipped_no_excess"
	TransferStatusSkippedNotEconomical TransferStatus = "skipped_not_economical"
	TransferStatusFailed               TransferStatus = "failed"
)

type NetworkTransferResult struct {
	Status  TransferStatus
	TxHash  string
	Amount  *entities.Wei
	Fee     *entities.Wei
	Error   error
	Message string
}

type TransferToColdWalletResult struct {
	BtcResult NetworkTransferResult
	RskResult NetworkTransferResult
}

func NewTransferToColdWalletResult(btcResult, rskResult NetworkTransferResult) *TransferToColdWalletResult {
	return &TransferToColdWalletResult{
		BtcResult: btcResult,
		RskResult: rskResult,
	}
}

type excessCalculationResult struct {
	BtcExcess        *entities.Wei
	BtcIsTimeForced  bool
	RbtcExcess       *entities.Wei
	RbtcIsTimeForced bool
}

type internalTxResult struct {
	TxHash string
	Amount *entities.Wei
	Fee    *entities.Wei
}

type TransferExcessToColdWalletUseCase struct {
	peginProvider                liquidity_provider.PeginLiquidityProvider
	pegoutProvider               liquidity_provider.PegoutLiquidityProvider
	generalProvider              liquidity_provider.LiquidityProvider
	lpRepository                 liquidity_provider.LiquidityProviderRepository
	coldWallet                   cold_wallet.ColdWallet
	btcWallet                    blockchain.BitcoinWallet
	rskWallet                    blockchain.RootstockWallet
	rpc                          blockchain.Rpc
	btcWalletMutex               sync.Locker
	rskWalletMutex               sync.Locker
	btcMinTransferFeeMultiplier  uint64
	rbtcMinTransferFeeMultiplier uint64
	forceTransferAfterSeconds    uint64
	eventBus                     entities.EventBus
	signer                       entities.Signer
	hashFunc                     entities.HashFunction
}

// NewTransferExcessToColdWalletUseCase creates a new use case for transferring excess liquidity to the cold wallet.
func NewTransferExcessToColdWalletUseCase(
	peginProvider liquidity_provider.PeginLiquidityProvider,
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	generalProvider liquidity_provider.LiquidityProvider,
	lpRepository liquidity_provider.LiquidityProviderRepository,
	coldWallet cold_wallet.ColdWallet,
	btcWallet blockchain.BitcoinWallet,
	rskWallet blockchain.RootstockWallet,
	rpc blockchain.Rpc,
	btcWalletMutex sync.Locker,
	rskWalletMutex sync.Locker,
	btcMinTransferFeeMultiplier uint64,
	rbtcMinTransferFeeMultiplier uint64,
	forceTransferAfterSeconds uint64,
	eventBus entities.EventBus,
	signer entities.Signer,
	hashFunc entities.HashFunction,
) *TransferExcessToColdWalletUseCase {
	return &TransferExcessToColdWalletUseCase{
		peginProvider:                peginProvider,
		pegoutProvider:               pegoutProvider,
		generalProvider:              generalProvider,
		lpRepository:                 lpRepository,
		coldWallet:                   coldWallet,
		btcWallet:                    btcWallet,
		rskWallet:                    rskWallet,
		rpc:                          rpc,
		btcWalletMutex:               btcWalletMutex,
		rskWalletMutex:               rskWalletMutex,
		btcMinTransferFeeMultiplier:  btcMinTransferFeeMultiplier,
		rbtcMinTransferFeeMultiplier: rbtcMinTransferFeeMultiplier,
		forceTransferAfterSeconds:    forceTransferAfterSeconds,
		eventBus:                     eventBus,
		signer:                       signer,
		hashFunc:                     hashFunc,
	}
}

// Run is the main entry point of the use case. It acquires wallet locks, validates configuration, calculates
// excess liquidity for both BTC and RBTC networks, executes the transfers, and persists the updated state.
func (useCase *TransferExcessToColdWalletUseCase) Run(ctx context.Context) (*TransferToColdWalletResult, error) {
	useCase.btcWalletMutex.Lock()
	defer useCase.btcWalletMutex.Unlock()

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	generalConfig, stateConfig, err := useCase.getAndValidateConfiguration(ctx)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.TransferExcessToColdWalletId, err)
	}

	currentLiquidity, err := useCase.getCurrentLiquidity(ctx)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.TransferExcessToColdWalletId, err)
	}

	excessResult, err := useCase.calculateExcessForBothNetworks(
		generalConfig,
		stateConfig,
		currentLiquidity,
	)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.TransferExcessToColdWalletId, err)
	}

	btcResult := useCase.handleBtcTransfer(excessResult.BtcExcess, excessResult.BtcIsTimeForced)
	if btcResult.Status == TransferStatusFailed {
		result := NewTransferToColdWalletResult(btcResult, NetworkTransferResult{})
		return result, usecases.WrapUseCaseError(usecases.TransferExcessToColdWalletId, btcResult.Error)
	}

	rskResult := useCase.handleRskTransfer(ctx, excessResult.RbtcExcess, excessResult.RbtcIsTimeForced)
	if rskResult.Status == TransferStatusFailed {
		result := NewTransferToColdWalletResult(btcResult, rskResult)
		return result, usecases.WrapUseCaseError(usecases.TransferExcessToColdWalletId, rskResult.Error)
	}

	result := NewTransferToColdWalletResult(btcResult, rskResult)
	if btcResult.Status == TransferStatusSuccess || rskResult.Status == TransferStatusSuccess {
		useCase.persistStateConfigAfterTransfers(ctx, stateConfig, btcResult, rskResult)
	}
	return result, nil
}

// persistStateConfigAfterTransfers updates the last transfer timestamps for each network that had a successful
// transfer, signs the updated configuration, and persists it to the database.
func (useCase *TransferExcessToColdWalletUseCase) persistStateConfigAfterTransfers(ctx context.Context, stateConfig liquidity_provider.StateConfiguration, btcResult, rskResult NetworkTransferResult) {
	updatedStateConfig := stateConfig
	now := time.Now().UTC().Unix()
	if btcResult.Status == TransferStatusSuccess {
		updatedStateConfig.LastBtcToColdWalletTransfer = now
	}
	if rskResult.Status == TransferStatusSuccess {
		updatedStateConfig.LastRbtcToColdWalletTransfer = now
	}
	newSigned, err := usecases.SignConfiguration(usecases.TransferExcessToColdWalletId, useCase.signer, useCase.hashFunc, updatedStateConfig)
	if err != nil {
		log.Errorf("TransferExcessToColdWallet: failed to sign state configuration: %v", err)
		return
	}
	if err := useCase.lpRepository.UpsertStateConfiguration(ctx, newSigned); err != nil {
		log.Errorf("TransferExcessToColdWallet: failed to persist state configuration: %v", err)
		return
	}
}

// getAndValidateConfiguration reads the general and state configurations through the LiquidityProvider interface
// (which validates signatures) and checks that the cold wallet, max liquidity, and excess tolerance are properly set.
func (useCase *TransferExcessToColdWalletUseCase) getAndValidateConfiguration(ctx context.Context) (liquidity_provider.GeneralConfiguration, liquidity_provider.StateConfiguration, error) {
	var zeroGeneral liquidity_provider.GeneralConfiguration
	var zeroState liquidity_provider.StateConfiguration

	if err := useCase.validateColdWallet(); err != nil {
		return zeroGeneral, zeroState, err
	}

	generalConfig := useCase.generalProvider.GeneralConfiguration(ctx)
	if generalConfig.MaxLiquidity == nil || generalConfig.MaxLiquidity.Cmp(entities.NewWei(0)) == 0 {
		return zeroGeneral, zeroState, NoMaxLiquidityConfiguredError
	}

	if err := generalConfig.ExcessTolerance.Validate(); err != nil {
		return zeroGeneral, zeroState, err
	}

	stateConfig, err := useCase.getStateConfiguration(ctx)
	if err != nil {
		return zeroGeneral, zeroState, err
	}

	return generalConfig, stateConfig, nil
}

// calculateExcessForBothNetworks computes the excess liquidity for BTC and RBTC by comparing current balances
// against the target plus the configured tolerance threshold, applying time-forcing logic
// when the last transfer exceeds the configured interval.
func (useCase *TransferExcessToColdWalletUseCase) calculateExcessForBothNetworks(
	generalConfig liquidity_provider.GeneralConfiguration,
	stateConfig liquidity_provider.StateConfiguration,
	currentLiquidity currentLiquidityResult,
) (excessCalculationResult, error) {
	targetPerNetwork, err := new(entities.Wei).Div(generalConfig.MaxLiquidity, entities.NewWei(2))
	if err != nil {
		return excessCalculationResult{}, err
	}
	threshold := useCase.calculateThreshold(targetPerNetwork, generalConfig.ExcessTolerance)

	btcLiquidityExcess, btcIsTimeForced := useCase.calculateExcessWithTimeForcing(
		targetPerNetwork,
		threshold,
		currentLiquidity.Btc,
		stateConfig.LastBtcToColdWalletTransfer,
	)
	rbtcLiquidityExcess, rbtcIsTimeForced := useCase.calculateExcessWithTimeForcing(
		targetPerNetwork,
		threshold,
		currentLiquidity.Rbtc,
		stateConfig.LastRbtcToColdWalletTransfer,
	)

	return excessCalculationResult{
		BtcExcess:        btcLiquidityExcess,
		BtcIsTimeForced:  btcIsTimeForced,
		RbtcExcess:       rbtcLiquidityExcess,
		RbtcIsTimeForced: rbtcIsTimeForced,
	}, nil
}

// publishBtcTransferEvent emits a BTC cold wallet transfer event to the event bus, distinguishing
// between threshold-triggered and time-forced transfers.
func (useCase *TransferExcessToColdWalletUseCase) publishBtcTransferEvent(txResult *internalTxResult, isTimeForcingTransfer bool) {
	if isTimeForcingTransfer {
		useCase.eventBus.Publish(cold_wallet.BtcTransferredDueToTimeForcingEvent{
			Event:  entities.NewBaseEvent(cold_wallet.BtcTransferredDueToTimeForcingEventId),
			Amount: txResult.Amount,
			TxHash: txResult.TxHash,
			Fee:    txResult.Fee,
		})
	} else {
		useCase.eventBus.Publish(cold_wallet.BtcTransferredDueToThresholdEvent{
			Event:  entities.NewBaseEvent(cold_wallet.BtcTransferredDueToThresholdEventId),
			Amount: txResult.Amount,
			TxHash: txResult.TxHash,
			Fee:    txResult.Fee,
		})
	}
}

// handleBtcTransfer orchestrates a single BTC transfer: skips if no excess, executes the transfer,
// checks if the amount is economical, and publishes the corresponding event on success.
func (useCase *TransferExcessToColdWalletUseCase) handleBtcTransfer(excess *entities.Wei, isTimeForcingTransfer bool) NetworkTransferResult {
	if excess.Cmp(entities.NewWei(0)) <= 0 {
		return NetworkTransferResult{
			Status:  TransferStatusSkippedNoExcess,
			Message: "No BTC excess to transfer",
		}
	}

	txResult, err := useCase.executeBtcTransfer(excess)
	if err == nil {
		useCase.publishBtcTransferEvent(txResult, isTimeForcingTransfer)

		return NetworkTransferResult{
			Status:  TransferStatusSuccess,
			TxHash:  txResult.TxHash,
			Amount:  txResult.Amount,
			Fee:     txResult.Fee,
			Message: "BTC transfer successful",
		}
	}

	if errors.Is(err, TransferNotEconomicalError) {
		return NetworkTransferResult{
			Status:  TransferStatusSkippedNotEconomical,
			Message: "BTC transfer skipped - amount not economical (less than minimum fee multiplier)",
		}
	}

	return NetworkTransferResult{
		Status:  TransferStatusFailed,
		Error:   err,
		Message: "BTC transfer failed",
	}
}

// publishRskTransferEvent emits an RBTC cold wallet transfer event to the event bus, distinguishing
// between threshold-triggered and time-forced transfers.
func (useCase *TransferExcessToColdWalletUseCase) publishRskTransferEvent(txResult *internalTxResult, isTimeForcingTransfer bool) {
	if isTimeForcingTransfer {
		useCase.eventBus.Publish(cold_wallet.RbtcTransferredDueToTimeForcingEvent{
			Event:  entities.NewBaseEvent(cold_wallet.RbtcTransferredDueToTimeForcingEventId),
			Amount: txResult.Amount,
			TxHash: txResult.TxHash,
			Fee:    txResult.Fee,
		})
	} else {
		useCase.eventBus.Publish(cold_wallet.RbtcTransferredDueToThresholdEvent{
			Event:  entities.NewBaseEvent(cold_wallet.RbtcTransferredDueToThresholdEventId),
			Amount: txResult.Amount,
			TxHash: txResult.TxHash,
			Fee:    txResult.Fee,
		})
	}
}

// handleRskTransfer orchestrates a single RBTC transfer: skips if no excess, executes the transfer,
// checks if the amount is economical, and publishes the corresponding event on success.
func (useCase *TransferExcessToColdWalletUseCase) handleRskTransfer(ctx context.Context, excess *entities.Wei, isTimeForcingTransfer bool) NetworkTransferResult {
	if excess.Cmp(entities.NewWei(0)) <= 0 {
		return NetworkTransferResult{
			Status:  TransferStatusSkippedNoExcess,
			Message: "No RSK excess to transfer",
		}
	}

	txResult, err := useCase.executeRskTransfer(ctx, excess)
	if err == nil {
		useCase.publishRskTransferEvent(txResult, isTimeForcingTransfer)

		return NetworkTransferResult{
			Status:  TransferStatusSuccess,
			TxHash:  txResult.TxHash,
			Amount:  txResult.Amount,
			Fee:     txResult.Fee,
			Message: "RSK transfer successful",
		}
	}

	if errors.Is(err, TransferNotEconomicalError) {
		return NetworkTransferResult{
			Status:  TransferStatusSkippedNotEconomical,
			Message: "RSK transfer skipped - amount not economical (less than minimum fee multiplier)",
		}
	}

	return NetworkTransferResult{
		Status:  TransferStatusFailed,
		Error:   err,
		Message: "RSK transfer failed",
	}
}

type currentLiquidityResult struct {
	Btc  *entities.Wei
	Rbtc *entities.Wei
}

// validateColdWallet checks that both BTC and RSK cold wallet addresses are configured.
func (useCase *TransferExcessToColdWalletUseCase) validateColdWallet() error {
	btcAddress := useCase.coldWallet.GetBtcAddress()
	rskAddress := useCase.coldWallet.GetRskAddress()

	if btcAddress == "" || rskAddress == "" {
		return NoColdWalletConfiguredError
	}

	return nil
}

// getCurrentLiquidity retrieves the available BTC and RBTC balances from the pegout and pegin providers respectively.
func (useCase *TransferExcessToColdWalletUseCase) getCurrentLiquidity(ctx context.Context) (currentLiquidityResult, error) {
	btcCurrentLiquidity, err := useCase.pegoutProvider.AvailablePegoutLiquidity(ctx)
	if err != nil {
		return currentLiquidityResult{}, err
	}

	rskCurrentLiquidity, err := useCase.peginProvider.AvailablePeginLiquidity(ctx)
	if err != nil {
		return currentLiquidityResult{}, err
	}

	return currentLiquidityResult{
		Btc:  btcCurrentLiquidity,
		Rbtc: rskCurrentLiquidity,
	}, nil
}

// calculateThreshold computes the balance above which a transfer is triggered, by adding the configured
// tolerance (either a fixed amount or a percentage) to the per-network target.
func (useCase *TransferExcessToColdWalletUseCase) calculateThreshold(
	target *entities.Wei,
	tolerance liquidity_provider.ExcessTolerance,
) *entities.Wei {
	if tolerance.IsFixed {
		return new(entities.Wei).Add(target, tolerance.FixedValue)
	}

	thresholdBigInt := utils.ApplyPercentageIncrease(target.AsBigInt(), tolerance.PercentageValue.Native())
	return entities.NewBigWei(thresholdBigInt)
}

// calculateExcess returns currentLiquidity minus target if currentLiquidity exceeds compareValue, otherwise zero.
func (useCase *TransferExcessToColdWalletUseCase) calculateExcess(
	target *entities.Wei,
	currentLiquidity *entities.Wei,
	compareValue *entities.Wei,
) *entities.Wei {
	if currentLiquidity.Cmp(compareValue) > 0 {
		return new(entities.Wei).Sub(currentLiquidity, target)
	}
	return entities.NewWei(0)
}

// executeBtcTransfer estimates the BTC transaction fee, verifies the transfer is economical
// (amount >= fee * multiplier), and sends the funds to the cold wallet.
func (useCase *TransferExcessToColdWalletUseCase) executeBtcTransfer(amount *entities.Wei) (*internalTxResult, error) {
	coldBtcAddress := useCase.coldWallet.GetBtcAddress()

	feeEstimation, err := useCase.btcWallet.EstimateTxFees(coldBtcAddress, amount)
	if err != nil {
		return nil, err
	}

	// Check if transfer is economical: amount >= fee * multiplier
	minWorthwhileAmount := new(entities.Wei).Mul(feeEstimation.Value, entities.NewUWei(useCase.btcMinTransferFeeMultiplier))
	if amount.Cmp(minWorthwhileAmount) < 0 {
		return nil, TransferNotEconomicalError
	}

	txResult, err := useCase.btcWallet.Send(coldBtcAddress, amount)
	if err != nil {
		return nil, err
	}

	return &internalTxResult{
		TxHash: txResult.Hash,
		Amount: amount,
		Fee:    feeEstimation.Value,
	}, nil
}

// executeRskTransfer estimates gas costs, subtracts them from the amount, verifies the transfer is economical
// (net amount >= gasCost * multiplier), and sends the RBTC to the cold wallet.
func (useCase *TransferExcessToColdWalletUseCase) executeRskTransfer(ctx context.Context, amount *entities.Wei) (*internalTxResult, error) {
	coldRskAddress := useCase.coldWallet.GetRskAddress()

	gasPrice, err := useCase.rpc.Rsk.GasPrice(ctx)
	if err != nil {
		return nil, err
	}

	gasCost := new(entities.Wei).Mul(gasPrice, entities.NewWei(SimpleTransferGasLimit))

	// Subtract gas cost from the amount to transfer
	amountToTransfer := new(entities.Wei).Sub(amount, gasCost)

	// Check if transfer is economical: amountToTransfer >= gasCost * multiplier
	minWorthwhileAmount := new(entities.Wei).Mul(gasCost, entities.NewUWei(useCase.rbtcMinTransferFeeMultiplier))
	if amountToTransfer.Cmp(minWorthwhileAmount) < 0 {
		return nil, TransferNotEconomicalError
	}

	config := blockchain.NewTransactionConfig(amountToTransfer, SimpleTransferGasLimit, gasPrice)

	receipt, err := useCase.rskWallet.SendRbtc(ctx, config, coldRskAddress)
	if err != nil {
		return nil, err
	}

	actualFee := new(entities.Wei).Mul(receipt.GasPrice, entities.NewUWei(receipt.GasUsed.Uint64()))

	return &internalTxResult{
		TxHash: receipt.TransactionHash,
		Amount: amountToTransfer,
		Fee:    actualFee,
	}, nil
}

// getStateConfiguration reads the state configuration through the LiquidityProvider interface and validates
// that both last-transfer timestamps have been initialized (non-zero).
func (useCase *TransferExcessToColdWalletUseCase) getStateConfiguration(ctx context.Context) (liquidity_provider.StateConfiguration, error) {
	stateConfig, err := useCase.generalProvider.StateConfiguration(ctx)
	if err != nil {
		return liquidity_provider.StateConfiguration{}, err
	}

	if stateConfig.LastBtcToColdWalletTransfer == 0 || stateConfig.LastRbtcToColdWalletTransfer == 0 {
		return liquidity_provider.StateConfiguration{}, NoTransferHistoryConfiguredError
	}

	return stateConfig, nil
}

// calculateExcessWithTimeForcing first checks if liquidity exceeds the threshold (normal excess). If not,
// it checks whether enough time has elapsed since the last transfer to force a transfer at the target level.
// Returns the excess amount and whether the transfer was time-forced.
func (useCase *TransferExcessToColdWalletUseCase) calculateExcessWithTimeForcing(
	target *entities.Wei,
	threshold *entities.Wei,
	currentLiquidity *entities.Wei,
	lastTransferUnix int64,
) (*entities.Wei, bool) {
	// First, calculate normal excess (respecting threshold)
	excess := useCase.calculateExcess(target, currentLiquidity, threshold)

	if excess.Cmp(entities.NewWei(0)) > 0 {
		return excess, false // Transfer due to threshold, not time forcing
	}

	if useCase.shouldForceTransferDueToTime(lastTransferUnix) {
		// Calculate excess respecting target
		forcedExcess := useCase.calculateExcess(target, currentLiquidity, target)
		return forcedExcess, true
	}

	return entities.NewWei(0), false
}

// shouldForceTransferDueToTime returns true if the elapsed time since the last transfer exceeds the
// configured forceTransferAfterSeconds interval.
func (useCase *TransferExcessToColdWalletUseCase) shouldForceTransferDueToTime(lastTransferUnix int64) bool {
	secondsSinceLastTransfer := time.Now().Unix() - lastTransferUnix
	return secondsSinceLastTransfer >= int64(useCase.forceTransferAfterSeconds)
}
