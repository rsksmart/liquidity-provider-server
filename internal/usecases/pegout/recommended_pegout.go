package pegout

import (
	"context"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"math/big"
)

type RecommendedPegoutUseCase struct {
	pegoutProvider liquidity_provider.PegoutLiquidityProvider
	contracts      blockchain.RskContracts
	rpc            blockchain.Rpc
	btcWallet      blockchain.BitcoinWallet
	scale          int64
}

func NewRecommendedPegoutUseCase(
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	btcWallet blockchain.BitcoinWallet,
	scale int64,
) *RecommendedPegoutUseCase {
	return &RecommendedPegoutUseCase{
		pegoutProvider: pegoutProvider,
		contracts:      contracts,
		rpc:            rpc,
		btcWallet:      btcWallet,
		scale:          scale,
	}
}

func (useCase *RecommendedPegoutUseCase) Run(
	ctx context.Context,
	userBalance *entities.Wei,
	destinationType blockchain.BtcAddressType,
) (usecases.RecommendedOperationResult, error) {
	const defaultBtcAddressType = "p2pkh"
	if destinationType == "" {
		destinationType = defaultBtcAddressType
	}
	config := useCase.pegoutProvider.PegoutConfiguration(ctx)
	result := new(big.Int).Set(userBalance.AsBigInt())

	// Percentage fees
	scaledCallFeePercentage := useCase.getScaledCallFeePercentage(config)

	// Fixed fees
	gasFeeEstimation, err := useCase.getGasFee(destinationType, userBalance)
	if err != nil {
		return usecases.RecommendedOperationResult{}, usecases.WrapUseCaseError(usecases.RecommendedPegoutId, err)
	}
	fixedCallFeeEstimation := useCase.getFixedCallFee(config)

	// Result calculation
	totalPercentages := big.NewInt(0)
	totalPercentages.Add(big.NewInt(useCase.scale), scaledCallFeePercentage)

	result.Sub(result, gasFeeEstimation)
	result.Sub(result, fixedCallFeeEstimation)
	scaledBase := new(big.Int).Mul(big.NewInt(useCase.scale), result)
	result.Quo(scaledBase, totalPercentages)

	if err = useCase.validateRecommendedValue(ctx, config, result); err != nil {
		return usecases.RecommendedOperationResult{}, err
	}

	return usecases.RecommendedOperationResult{
		RecommendedQuoteValue: entities.NewBigWei(result),
		EstimatedGasFee:       entities.NewBigWei(gasFeeEstimation),
		EstimatedCallFee: entities.NewBigWei(
			useCase.estimateCallFee(result, fixedCallFeeEstimation, scaledCallFeePercentage),
		),
	}, nil
}

func (useCase *RecommendedPegoutUseCase) getScaledCallFeePercentage(
	config liquidity_provider.PegoutConfiguration,
) *big.Int {
	var scaledPercentageFee = new(big.Int)
	new(big.Float).Mul(
		new(big.Float).Quo(
			config.FeePercentage.Native(),
			utils.NewBigFloat64(100).Native(),
		),
		big.NewFloat(float64(useCase.scale)),
	).Int(scaledPercentageFee)

	return scaledPercentageFee
}

func (useCase *RecommendedPegoutUseCase) getGasFee(
	destinationType blockchain.BtcAddressType,
	amount *entities.Wei,
) (*big.Int, error) {
	address, err := useCase.rpc.Btc.GetZeroAddress(destinationType)
	if err != nil {
		return nil, err
	}
	estimation, err := useCase.btcWallet.EstimateTxFees(address, amount)
	if err != nil {
		return nil, err
	}
	return estimation.Value.AsBigInt(), nil
}

func (useCase *RecommendedPegoutUseCase) getFixedCallFee(
	config liquidity_provider.PegoutConfiguration,
) *big.Int {
	return config.FixedFee.AsBigInt()
}

func (useCase *RecommendedPegoutUseCase) estimateCallFee(
	result *big.Int,
	fixedCallFee *big.Int,
	scaledCallFeePercentage *big.Int,
) *big.Int {
	return new(big.Int).Add(
		fixedCallFee,
		new(big.Int).Quo(
			new(big.Int).Mul(result, scaledCallFeePercentage),
			big.NewInt(useCase.scale),
		),
	)
}

func (useCase *RecommendedPegoutUseCase) validateRecommendedValue(
	ctx context.Context,
	config liquidity_provider.PegoutConfiguration,
	result *big.Int,
) error {
	var err error

	if err = config.ValidateAmount(entities.NewBigWei(result)); err != nil {
		err = fmt.Errorf("recommended amount %s is out of range: %w", entities.NewBigWei(result).String(), err)
		return usecases.WrapUseCaseError(usecases.RecommendedPegoutId, err)
	}

	if err = useCase.pegoutProvider.HasPegoutLiquidity(ctx, entities.NewBigWei(result)); err != nil {
		return usecases.WrapUseCaseError(usecases.RecommendedPegoutId, usecases.NoLiquidityError)
	}

	if err = usecases.ValidateMinLockValue(usecases.RecommendedPegoutId, useCase.contracts.Bridge, entities.NewBigWei(result)); err != nil {
		err = fmt.Errorf("recommended amount %s is below the minimum lock value: %w", entities.NewBigWei(result).String(), err)
		return err
	}
	return nil
}
