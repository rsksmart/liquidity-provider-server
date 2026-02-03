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
}

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
	}
}

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
		currentLiquidity.Btc,
		currentLiquidity.Rbtc,
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

	return NewTransferToColdWalletResult(btcResult, rskResult), nil
}

func (useCase *TransferExcessToColdWalletUseCase) getAndValidateConfiguration(ctx context.Context) (*liquidity_provider.GeneralConfiguration, *liquidity_provider.StateConfiguration, error) {
	if err := useCase.validateColdWallet(); err != nil {
		return nil, nil, err
	}

	generalConfig := useCase.generalProvider.GeneralConfiguration(ctx)
	if generalConfig.MaxLiquidity == nil || generalConfig.MaxLiquidity.Cmp(entities.NewWei(0)) == 0 {
		return nil, nil, NoMaxLiquidityConfiguredError
	}

	if err := generalConfig.ExcessTolerance.Validate(); err != nil {
		return nil, nil, err
	}

	stateConfig, err := useCase.getStateConfiguration(ctx)
	if err != nil {
		return nil, nil, err
	}

	return &generalConfig, stateConfig, nil
}

func (useCase *TransferExcessToColdWalletUseCase) calculateExcessForBothNetworks(
	generalConfig *liquidity_provider.GeneralConfiguration,
	stateConfig *liquidity_provider.StateConfiguration,
	currentBtcLiquidity *entities.Wei,
	currentRbtcLiquidity *entities.Wei,
) (*excessCalculationResult, error) {
	targetPerNetwork, err := new(entities.Wei).Div(generalConfig.MaxLiquidity, entities.NewWei(2))
	if err != nil {
		return nil, err
	}
	threshold := useCase.calculateThreshold(targetPerNetwork, generalConfig.ExcessTolerance)

	btcLiquidityExcess, btcIsTimeForced := useCase.calculateExcessWithTimeForcing(
		targetPerNetwork,
		threshold,
		currentBtcLiquidity,
		stateConfig.LastBtcToColdWalletTransfer,
	)
	rbtcLiquidityExcess, rbtcIsTimeForced := useCase.calculateExcessWithTimeForcing(
		targetPerNetwork,
		threshold,
		currentRbtcLiquidity,
		stateConfig.LastRbtcToColdWalletTransfer,
	)

	return &excessCalculationResult{
		BtcExcess:        btcLiquidityExcess,
		BtcIsTimeForced:  btcIsTimeForced,
		RbtcExcess:       rbtcLiquidityExcess,
		RbtcIsTimeForced: rbtcIsTimeForced,
	}, nil
}

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

func (useCase *TransferExcessToColdWalletUseCase) validateColdWallet() error {
	btcAddress := useCase.coldWallet.GetBtcAddress()
	rskAddress := useCase.coldWallet.GetRskAddress()

	if btcAddress == "" || rskAddress == "" {
		return NoColdWalletConfiguredError
	}

	return nil
}

func (useCase *TransferExcessToColdWalletUseCase) getCurrentLiquidity(ctx context.Context) (*currentLiquidityResult, error) {
	// AvailablePegoutLiquidity returns BTC balance (from BTC wallet)
	btcCurrentLiquidity, err := useCase.pegoutProvider.AvailablePegoutLiquidity(ctx)
	if err != nil {
		return nil, err
	}

	// AvailablePeginLiquidity returns RBTC balance (from RSK wallet + PegIn contract)
	rskCurrentLiquidity, err := useCase.peginProvider.AvailablePeginLiquidity(ctx)
	if err != nil {
		return nil, err
	}

	return &currentLiquidityResult{
		Btc:  btcCurrentLiquidity,
		Rbtc: rskCurrentLiquidity,
	}, nil
}

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

func (useCase *TransferExcessToColdWalletUseCase) getStateConfiguration(ctx context.Context) (*liquidity_provider.StateConfiguration, error) {
	stateConfig := useCase.generalProvider.StateConfiguration(ctx)

	// Check if configuration is valid (empty/nil fields mean validation failed or not found)
	if stateConfig.LastBtcToColdWalletTransfer == nil {
		return nil, NoTransferHistoryConfiguredError
	}
	if stateConfig.LastRbtcToColdWalletTransfer == nil {
		return nil, NoTransferHistoryConfiguredError
	}

	return &stateConfig, nil
}

func (useCase *TransferExcessToColdWalletUseCase) calculateExcessWithTimeForcing(
	target *entities.Wei,
	threshold *entities.Wei,
	currentLiquidity *entities.Wei,
	lastTransferTime *time.Time,
) (*entities.Wei, bool) {
	// First, calculate normal excess (respecting threshold)
	excess := useCase.calculateExcess(target, currentLiquidity, threshold)

	if excess.Cmp(entities.NewWei(0)) > 0 {
		return excess, false // Transfer due to threshold, not time forcing
	}

	if useCase.shouldForceTransferDueToTime(lastTransferTime) {
		// Calculate excess respecting target
		forcedExcess := useCase.calculateExcess(target, currentLiquidity, target)
		return forcedExcess, true
	}

	return entities.NewWei(0), false
}

func (useCase *TransferExcessToColdWalletUseCase) shouldForceTransferDueToTime(lastTransferTime *time.Time) bool {
	timeSinceLastTransfer := time.Since(*lastTransferTime)
	timeRequiredToForceTransfer := time.Duration(useCase.forceTransferAfterSeconds) * time.Second
	return timeSinceLastTransfer >= timeRequiredToForceTransfer
}
