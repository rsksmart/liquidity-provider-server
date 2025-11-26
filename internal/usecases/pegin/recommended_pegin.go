package pegin

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

type RecommendedPeginUseCase struct {
	peginProvider       liquidity_provider.PeginLiquidityProvider
	contracts           blockchain.RskContracts
	rpc                 blockchain.Rpc
	feeCollectorAddress string
	scale               int64
}

func NewRecommendedPeginUseCase(
	peginProvider liquidity_provider.PeginLiquidityProvider,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	feeCollectorAddress string,
	scale int64,
) *RecommendedPeginUseCase {
	return &RecommendedPeginUseCase{
		peginProvider:       peginProvider,
		contracts:           contracts,
		rpc:                 rpc,
		feeCollectorAddress: feeCollectorAddress,
		scale:               scale,
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

	if !blockchain.IsRskAddress(destinationAddress) {
		destinationAddress = blockchain.RskZeroAddress
	}

	// Percentage fees
	scaledProductFee, err := useCase.getScaledProductFeePercentage()
	if err != nil {
		return usecases.RecommendedOperationResult{}, usecases.WrapUseCaseError(usecases.RecommendedPeginId, err)
	}
	scaledCallFeePercentage := useCase.getScaledCallFeePercentage(config)

	// Fixed fees
	gasFeeEstimation, err := useCase.getGasFee(ctx, destinationAddress, data, userBalance, scaledProductFee)
	if err != nil {
		return usecases.RecommendedOperationResult{}, usecases.WrapUseCaseError(usecases.RecommendedPeginId, err)
	}
	fixedCallFeeEstimation := useCase.getFixedCallFee(config)

	// Result calculation
	totalPercentages := big.NewInt(0)
	totalPercentages.Add(scaledProductFee, scaledCallFeePercentage)
	totalPercentages.Add(totalPercentages, big.NewInt(useCase.scale))

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
		EstimatedProductFee: entities.NewBigWei(
			useCase.estimateProductFee(result, scaledProductFee),
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

func (useCase *RecommendedPeginUseCase) getScaledProductFeePercentage() (*big.Int, error) {
	// should be already scaled in the contract
	uintProductFee, err := useCase.contracts.PegIn.DaoFeePercentage()
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetUint64(uintProductFee), nil
}

func (useCase *RecommendedPeginUseCase) getGasFee(
	ctx context.Context,
	destinationAddress string,
	data []byte,
	amount *entities.Wei,
	scaledFeePercentage *big.Int,
) (*big.Int, error) {
	daoFeeAmount := new(big.Int).Quo(
		new(big.Int).Mul(amount.AsBigInt(), scaledFeePercentage),
		big.NewInt(useCase.scale),
	)
	daoGasAmount, err := useCase.rpc.Rsk.EstimateGas(
		ctx,
		useCase.feeCollectorAddress,
		entities.NewBigWei(daoFeeAmount),
		make([]byte, 0),
	)
	if err != nil {
		return nil, err
	}

	estimatedCallGas, err := useCase.rpc.Rsk.EstimateGas(ctx, destinationAddress, amount, data)
	if err != nil {
		return nil, err
	}

	gasPrice, err := useCase.rpc.Rsk.GasPrice(ctx)
	if err != nil {
		return nil, err
	}
	daoGasFee := new(big.Int).Mul(daoGasAmount.AsBigInt(), gasPrice.AsBigInt())
	callGasFee := new(big.Int).Mul(estimatedCallGas.AsBigInt(), gasPrice.AsBigInt())

	return new(big.Int).Add(callGasFee, daoGasFee), nil
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
) error {
	var err error

	if err = config.ValidateAmount(entities.NewBigWei(result)); err != nil {
		err = fmt.Errorf("recommended amount %s is out of range: %w", entities.NewBigWei(result).String(), err)
		return usecases.WrapUseCaseError(usecases.RecommendedPeginId, err)
	}

	if err = useCase.peginProvider.HasPeginLiquidity(ctx, entities.NewBigWei(result)); err != nil {
		return usecases.WrapUseCaseError(usecases.RecommendedPeginId, usecases.NoLiquidityError)
	}

	if err = usecases.ValidateMinLockValue(usecases.RecommendedPeginId, useCase.contracts.Bridge, entities.NewBigWei(result)); err != nil {
		err = fmt.Errorf("recommended amount %s is below the minimum lock value: %w", entities.NewBigWei(result).String(), err)
		return err
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

func (useCase *RecommendedPeginUseCase) estimateProductFee(
	result *big.Int,
	scaledProductFee *big.Int,
) *big.Int {
	return new(big.Int).Quo(
		new(big.Int).Mul(result, scaledProductFee),
		big.NewInt(useCase.scale),
	)
}
