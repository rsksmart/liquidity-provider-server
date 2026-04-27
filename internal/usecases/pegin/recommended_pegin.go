package pegin

import (
	"context"
	"fmt"
	"math/big"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type RecommendedPeginUseCase struct {
	peginProvider liquidity_provider.PeginLiquidityProvider
	contracts     blockchain.RskContracts
	rpc           blockchain.Rpc
	scale         int64
}

func NewRecommendedPeginUseCase(
	peginProvider liquidity_provider.PeginLiquidityProvider,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	scale int64,
) *RecommendedPeginUseCase {
	return &RecommendedPeginUseCase{
		peginProvider: peginProvider,
		contracts:     contracts,
		rpc:           rpc,
		scale:         scale,
	}
}

func (useCase *RecommendedPeginUseCase) Run(
	ctx context.Context,
	userBalance *entities.Wei,
	destinationAddress string,
	data []byte,
) (usecases.RecommendedOperationResult, error) {
	config := useCase.peginProvider.PeginConfiguration(ctx)
	result := new(big.Int).Set(userBalance.AsBigInt())

	if err := config.ValidateAmount(userBalance); err != nil {
		err = fmt.Errorf("provided amount %s is out of range: %w", userBalance.String(), err)
		return usecases.RecommendedOperationResult{}, usecases.WrapUseCaseError(usecases.RecommendedPeginId, err)
	}

	if !blockchain.IsRskAddress(destinationAddress) {
		destinationAddress = blockchain.RskZeroAddress
	}

	// Percentage fees
	scaledCallFeePercentage := useCase.getScaledCallFeePercentage(config)

	// Fixed fees
	gasFeeEstimation, err := useCase.getGasFee(ctx, destinationAddress, data, userBalance)
	if err != nil {
		return usecases.RecommendedOperationResult{}, usecases.WrapUseCaseError(usecases.RecommendedPeginId, err)
	}
	fixedCallFeeEstimation := useCase.getFixedCallFee(config)

	// Result calculation
	totalPercentages := big.NewInt(0)
	totalPercentages.Add(big.NewInt(useCase.scale), scaledCallFeePercentage)

	result.Sub(result, gasFeeEstimation)
	result.Sub(result, fixedCallFeeEstimation)
	scaledBase := new(big.Int).Mul(big.NewInt(useCase.scale), result)
	result.Quo(scaledBase, totalPercentages)

	// Compute a best-effort minimum raw input that should produce result >= MinValue after fees.
	// It is an estimate because gas is re-estimated against the original user amount, not the
	// suggested one, so the follow-up request may yield a slightly different fee.
	// requiredNet = ceil(MinValue * totalPercentages / scale)
	// minimumAcceptable = gasFee + fixedFee + requiredNet
	minValue := config.MinValue.AsBigInt()
	requiredNet := new(big.Int).Mul(minValue, totalPercentages)
	requiredNet.Add(requiredNet, new(big.Int).Sub(big.NewInt(useCase.scale), big.NewInt(1)))
	requiredNet.Quo(requiredNet, big.NewInt(useCase.scale))
	minimumAcceptable := new(big.Int).Add(gasFeeEstimation, fixedCallFeeEstimation)
	minimumAcceptable.Add(minimumAcceptable, requiredNet)

	if err = useCase.validateRecommendedValue(ctx, config, result, minimumAcceptable); err != nil {
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

func (useCase *RecommendedPeginUseCase) getScaledCallFeePercentage(
	config liquidity_provider.PeginConfiguration,
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

func (useCase *RecommendedPeginUseCase) getGasFee(
	ctx context.Context,
	destinationAddress string,
	data []byte,
	amount *entities.Wei,
) (*big.Int, error) {
	estimatedCallGas, err := useCase.rpc.Rsk.EstimateGas(ctx, destinationAddress, amount, data)
	if err != nil {
		return nil, err
	}

	gasPrice, err := useCase.rpc.Rsk.GasPrice(ctx)
	if err != nil {
		return nil, err
	}
	callGasFee := new(big.Int).Mul(estimatedCallGas.AsBigInt(), gasPrice.AsBigInt())

	return callGasFee, nil
}

func (useCase *RecommendedPeginUseCase) getFixedCallFee(
	config liquidity_provider.PeginConfiguration,
) *big.Int {
	return config.FixedFee.AsBigInt()
}

func (useCase *RecommendedPeginUseCase) validateRecommendedValue(
	ctx context.Context,
	config liquidity_provider.PeginConfiguration,
	result *big.Int,
	minimumAcceptable *big.Int,
) error {
	minValue := config.MinValue.AsBigInt()
	if result.Cmp(minValue) < 0 && minimumAcceptable.Cmp(config.MaxValue.AsBigInt()) > 0 {
		return usecases.WrapUseCaseError(usecases.RecommendedPeginId,
			fmt.Errorf(
				"suggested minimum amount %s exceeds provider max %s with current fee estimates: %w",
				entities.NewBigWei(minimumAcceptable).String(),
				config.MaxValue.String(),
				liquidity_provider.AmountOutOfRangeError,
			),
		)
	} else if result.Cmp(minValue) < 0 {
		return usecases.WrapUseCaseError(usecases.RecommendedPeginId,
			usecases.NewEffectiveAmountTooLowError(entities.NewBigWei(result), config.MinValue.Copy(), entities.NewBigWei(minimumAcceptable)),
		)
	}

	if err := useCase.peginProvider.HasPeginLiquidity(ctx, entities.NewBigWei(result)); err != nil {
		return usecases.WrapUseCaseError(usecases.RecommendedPeginId, usecases.NoLiquidityError)
	}

	if err := usecases.ValidateMinLockValue(usecases.RecommendedPeginId, useCase.contracts.Bridge, entities.NewBigWei(result)); err != nil {
		return fmt.Errorf("recommended amount %s is below the minimum lock value: %w", entities.NewBigWei(result).String(), err)
	}
	return nil
}

func (useCase *RecommendedPeginUseCase) estimateCallFee(
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
