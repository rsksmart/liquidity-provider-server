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
	pegoutProvider      liquidity_provider.PegoutLiquidityProvider
	contracts           blockchain.RskContracts
	rpc                 blockchain.Rpc
	btcWallet           blockchain.BitcoinWallet
	scale               int64
	feeCollectorAddress string
}

func NewRecommendedPegoutUseCase(
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	btcWallet blockchain.BitcoinWallet,
	scale int64,
	feeCollectorAddress string,
) *RecommendedPegoutUseCase {
	return &RecommendedPegoutUseCase{
		pegoutProvider:      pegoutProvider,
		contracts:           contracts,
		rpc:                 rpc,
		btcWallet:           btcWallet,
		scale:               scale,
		feeCollectorAddress: feeCollectorAddress,
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
	scaledProductFee, err := useCase.getScaledProductFeePercentage()
	if err != nil {
		return usecases.RecommendedOperationResult{}, usecases.WrapUseCaseError(usecases.RecommendedPegoutId, err)
	}
	scaledCallFeePercentage := useCase.getScaledCallFeePercentage(config)

	// Fixed fees
	gasFeeEstimation, err := useCase.getGasFee(ctx, destinationType, userBalance, scaledProductFee)
	if err != nil {
		return usecases.RecommendedOperationResult{}, usecases.WrapUseCaseError(usecases.RecommendedPegoutId, err)
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

func (useCase *RecommendedPegoutUseCase) getScaledProductFeePercentage() (*big.Int, error) {
	// should be already scaled in the contract
	uintProductFee, err := useCase.contracts.FeeCollector.DaoFeePercentage()
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetUint64(uintProductFee), nil
}

func (useCase *RecommendedPegoutUseCase) getGasFee(
	ctx context.Context,
	destinationType blockchain.BtcAddressType,
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
	gasPrice, err := useCase.rpc.Rsk.GasPrice(ctx)
	if err != nil {
		return nil, err
	}
	daoGasFee := new(big.Int).Mul(daoGasAmount.AsBigInt(), gasPrice.AsBigInt())

	address, err := useCase.rpc.Btc.GetZeroAddress(destinationType)
	if err != nil {
		return nil, err
	}
	estimation, err := useCase.btcWallet.EstimateTxFees(address, amount)
	if err != nil {
		return nil, err
	}
	return new(big.Int).Add(estimation.Value.AsBigInt(), daoGasFee), nil
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

func (useCase *RecommendedPegoutUseCase) estimateProductFee(
	result *big.Int,
	scaledProductFee *big.Int,
) *big.Int {
	return new(big.Int).Quo(
		new(big.Int).Mul(result, scaledProductFee),
		big.NewInt(useCase.scale),
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
