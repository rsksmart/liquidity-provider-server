package liquidity_provider

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/cold_wallet"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
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
	rskWalletMutex               sync.Locker
	btcMinTransferFeeMultiplier  uint64
	rbtcMinTransferFeeMultiplier uint64
	forceTransferAfterSeconds    uint64
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
	rskWalletMutex sync.Locker,
	btcMinTransferFeeMultiplier uint64,
	rbtcMinTransferFeeMultiplier uint64,
	forceTransferAfterSeconds uint64,
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
		rskWalletMutex:               rskWalletMutex,
		btcMinTransferFeeMultiplier:  btcMinTransferFeeMultiplier,
		rbtcMinTransferFeeMultiplier: rbtcMinTransferFeeMultiplier,
		forceTransferAfterSeconds:    forceTransferAfterSeconds,
	}
}

func (useCase *TransferExcessToColdWalletUseCase) Run(ctx context.Context) (*TransferToColdWalletResult, error) {
	generalConfig, stateConfig, err := useCase.getAndValidateConfiguration(ctx)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.TransferExcessToColdWalletId, err)
	}

	currentBtcLiquidity, currentRbtcLiquidity, err := useCase.getCurrentLiquidity(ctx)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.TransferExcessToColdWalletId, err)
	}

	btcLiquidityExcess, rbtcLiquidityExcess, err := useCase.calculateExcessForBothNetworks(
		generalConfig,
		stateConfig,
		currentBtcLiquidity,
		currentRbtcLiquidity,
	)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.TransferExcessToColdWalletId, err)
	}

	result := &TransferToColdWalletResult{}

	result.BtcResult = useCase.executeBtcTransfer(btcLiquidityExcess)
	if result.BtcResult.Status == TransferStatusFailed {
		return result, usecases.WrapUseCaseError(usecases.TransferExcessToColdWalletId, result.BtcResult.Error)
	}

	result.RskResult = useCase.executeRskTransfer(ctx, rbtcLiquidityExcess)
	if result.RskResult.Status == TransferStatusFailed {
		return result, usecases.WrapUseCaseError(usecases.TransferExcessToColdWalletId, result.RskResult.Error)
	}

	return result, nil
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

	if stateConfig.LastBtcToColdWalletTransfer == nil {
		return nil, nil, NoTransferHistoryConfiguredError
	}
	if stateConfig.LastRbtcToColdWalletTransfer == nil {
		return nil, nil, NoTransferHistoryConfiguredError
	}

	return &generalConfig, stateConfig, nil
}

func (useCase *TransferExcessToColdWalletUseCase) calculateExcessForBothNetworks(
	generalConfig *liquidity_provider.GeneralConfiguration,
	stateConfig *liquidity_provider.StateConfiguration,
	currentBtcLiquidity *entities.Wei,
	currentRbtcLiquidity *entities.Wei,
) (*entities.Wei, *entities.Wei, error) {
	targetPerNetwork, err := new(entities.Wei).Div(generalConfig.MaxLiquidity, entities.NewWei(2))
	if err != nil {
		return nil, nil, err
	}
	threshold := useCase.calculateThreshold(targetPerNetwork, generalConfig.ExcessTolerance)

	btcLiquidityExcess := useCase.calculateExcessWithTimeForcing(
		targetPerNetwork,
		threshold,
		currentBtcLiquidity,
		stateConfig.LastBtcToColdWalletTransfer,
	)
	rbtcLiquidityExcess := useCase.calculateExcessWithTimeForcing(
		targetPerNetwork,
		threshold,
		currentRbtcLiquidity,
		stateConfig.LastRbtcToColdWalletTransfer,
	)

	return btcLiquidityExcess, rbtcLiquidityExcess, nil
}

func (useCase *TransferExcessToColdWalletUseCase) executeBtcTransfer(excess *entities.Wei) NetworkTransferResult {
	if excess.Cmp(entities.NewWei(0)) <= 0 {
		return NetworkTransferResult{
			Status:  TransferStatusSkippedNoExcess,
			Message: "No BTC excess to transfer",
		}
	}

	txResult, err := useCase.transferBtcExcess(excess)
	if err == nil {
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

func (useCase *TransferExcessToColdWalletUseCase) executeRskTransfer(ctx context.Context, excess *entities.Wei) NetworkTransferResult {
	if excess.Cmp(entities.NewWei(0)) <= 0 {
		return NetworkTransferResult{
			Status:  TransferStatusSkippedNoExcess,
			Message: "No RSK excess to transfer",
		}
	}

	txResult, err := useCase.transferRskExcess(ctx, excess)
	if err == nil {
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

func (useCase *TransferExcessToColdWalletUseCase) validateColdWallet() error {
	if useCase.coldWallet == nil {
		return NoColdWalletConfiguredError
	}

	btcAddress := useCase.coldWallet.GetBtcAddress()
	if btcAddress == "" {
		return NoColdWalletConfiguredError
	}

	rskAddress := useCase.coldWallet.GetRskAddress()
	if rskAddress == "" {
		return NoColdWalletConfiguredError
	}

	return nil
}

func (useCase *TransferExcessToColdWalletUseCase) getCurrentLiquidity(ctx context.Context) (*entities.Wei, *entities.Wei, error) {
	// AvailablePegoutLiquidity returns BTC balance (from BTC wallet)
	btcUsableLiquidity, err := useCase.pegoutProvider.AvailablePegoutLiquidity(ctx)
	if err != nil {
		return nil, nil, err
	}

	// AvailablePeginLiquidity returns RBTC balance (from RSK wallet + PegIn contract)
	rskUsableLiquidity, err := useCase.peginProvider.AvailablePeginLiquidity(ctx)
	if err != nil {
		return nil, nil, err
	}

	return btcUsableLiquidity, rskUsableLiquidity, nil
}

func (useCase *TransferExcessToColdWalletUseCase) calculateThreshold(
	target *entities.Wei,
	tolerance liquidity_provider.ExcessTolerance,
) *entities.Wei {
	if tolerance.IsFixed {
		return new(entities.Wei).Add(target, tolerance.FixedValue)
	}

	hundred := big.NewFloat(100)
	one := big.NewFloat(1)
	multiplier := new(big.Float).Add(
		one,
		new(big.Float).Quo(tolerance.PercentageValue.Native(), hundred),
	)
	targetFloat := new(big.Float).SetInt(target.AsBigInt())
	thresholdFloat := new(big.Float).Mul(targetFloat, multiplier)
	thresholdBigInt, _ := thresholdFloat.Int(nil)
	return entities.NewBigWei(thresholdBigInt)
}

func (useCase *TransferExcessToColdWalletUseCase) calculateExcess(
	target *entities.Wei,
	currentLiquidity *entities.Wei,
	threshold *entities.Wei,
) *entities.Wei {
	compareValue := target
	if threshold != nil {
		compareValue = threshold
	}

	if currentLiquidity.Cmp(compareValue) > 0 {
		return new(entities.Wei).Sub(currentLiquidity, target)
	}
	return entities.NewWei(0)
}

func (useCase *TransferExcessToColdWalletUseCase) transferBtcExcess(amount *entities.Wei) (*internalTxResult, error) {
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

func (useCase *TransferExcessToColdWalletUseCase) transferRskExcess(ctx context.Context, amount *entities.Wei) (*internalTxResult, error) {
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

	// Lock wallet for transaction
	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

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
	signedConfig, err := useCase.lpRepository.GetStateConfiguration(ctx)
	if err != nil {
		return nil, err
	}
	if signedConfig == nil {
		return nil, NoTransferHistoryConfiguredError
	}
	return &signedConfig.Value, nil
}

func (useCase *TransferExcessToColdWalletUseCase) calculateExcessWithTimeForcing(
	target *entities.Wei,
	threshold *entities.Wei,
	currentLiquidity *entities.Wei,
	lastTransferTime *time.Time,
) *entities.Wei {
	// First, calculate normal excess (respecting threshold)
	excess := useCase.calculateExcess(target, currentLiquidity, threshold)

	if excess.Cmp(entities.NewWei(0)) > 0 {
		return excess
	}

	// No excess above threshold, but check if we should force transfer due to time
	if useCase.shouldForceTransferDueToTime(lastTransferTime) {
		// Calculate excess without threshold (pass nil)
		forcedExcess := useCase.calculateExcess(target, currentLiquidity, nil)
		return forcedExcess
	}

	return entities.NewWei(0)
}

func (useCase *TransferExcessToColdWalletUseCase) shouldForceTransferDueToTime(lastTransferTime *time.Time) bool {
	secondsSinceLastTransfer := time.Since(*lastTransferTime).Seconds()
	return secondsSinceLastTransfer >= float64(useCase.forceTransferAfterSeconds)
}
