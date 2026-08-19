package rootstock

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/flyover_configurations"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

type flyoverConfigurationsContractImpl struct {
	client      RpcClientBinding
	address     string
	contract    *bind.BoundContract
	retryParams RetryParams
	abis        *FlyoverABIs
	binding     *bindings.FlyoverConfigurationsContract
}

func NewFlyoverConfigurationsContractImpl(
	client *RskClient,
	address string,
	contract *bind.BoundContract,
	retryParams RetryParams,
	binding *bindings.FlyoverConfigurationsContract,
	abis *FlyoverABIs,
) blockchain.FlyoverConfigurationsContract {
	return &flyoverConfigurationsContractImpl{
		client:      client.client,
		address:     address,
		contract:    contract,
		retryParams: retryParams,
		binding:     binding,
		abis:        abis,
	}
}

func (configurations *flyoverConfigurationsContractImpl) GetAddress() string {
	return configurations.address
}

func (configurations *flyoverConfigurationsContractImpl) CalculatePegInFee(amount *entities.Wei) (*entities.Wei, error) {
	opts := &bind.CallOpts{}
	result, err := rskRetry(configurations.retryParams.Retries, configurations.retryParams.Sleep,
		func() (*big.Int, error) {
			callData, dataErr := configurations.binding.TryPackCalculatePegInFee(amount.AsBigInt())
			if dataErr != nil {
				return nil, dataErr
			}
			return bind.Call(configurations.contract, opts, callData, configurations.binding.UnpackCalculatePegInFee)
		})
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(result), nil
}

func (configurations *flyoverConfigurationsContractImpl) GetRequiredPegInBtcConfirmations(amount *entities.Wei) (uint64, error) {
	opts := &bind.CallOpts{}
	result, err := rskRetry(configurations.retryParams.Retries, configurations.retryParams.Sleep,
		func() (*big.Int, error) {
			callData, dataErr := configurations.binding.TryPackGetRequiredPegInBtcConfirmations(amount.AsBigInt())
			if dataErr != nil {
				return nil, dataErr
			}
			return bind.Call(configurations.contract, opts, callData, configurations.binding.UnpackGetRequiredPegInBtcConfirmations)
		})
	if err != nil {
		return 0, err
	}
	return result.Uint64(), nil
}

func (configurations *flyoverConfigurationsContractImpl) CalculatePegOutFee(amount *entities.Wei) (*entities.Wei, error) {
	opts := &bind.CallOpts{}
	result, err := rskRetry(configurations.retryParams.Retries, configurations.retryParams.Sleep,
		func() (*big.Int, error) {
			callData, dataErr := configurations.binding.TryPackCalculatePegOutFee(amount.AsBigInt())
			if dataErr != nil {
				return nil, dataErr
			}
			return bind.Call(configurations.contract, opts, callData, configurations.binding.UnpackCalculatePegOutFee)
		})
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(result), nil
}

func (configurations *flyoverConfigurationsContractImpl) GetRequiredPegOutBtcConfirmations(amount *entities.Wei) (uint64, error) {
	opts := &bind.CallOpts{}
	result, err := rskRetry(configurations.retryParams.Retries, configurations.retryParams.Sleep,
		func() (*big.Int, error) {
			callData, dataErr := configurations.binding.TryPackGetRequiredPegOutBtcConfirmations(amount.AsBigInt())
			if dataErr != nil {
				return nil, dataErr
			}
			return bind.Call(configurations.contract, opts, callData, configurations.binding.UnpackGetRequiredPegOutBtcConfirmations)
		})
	if err != nil {
		return 0, err
	}
	return result.Uint64(), nil
}

func (configurations *flyoverConfigurationsContractImpl) GetPegOutConfiguration() (blockchain.PegOutConfiguration, error) {
	opts := &bind.CallOpts{}
	result, err := rskRetry(configurations.retryParams.Retries, configurations.retryParams.Sleep,
		func() (bindings.IFlyoverConfigurationsPegOutConfiguration, error) {
			callData, dataErr := configurations.binding.TryPackGetPegOutConfiguration()
			if dataErr != nil {
				return bindings.IFlyoverConfigurationsPegOutConfiguration{}, dataErr
			}
			return bind.Call(configurations.contract, opts, callData, configurations.binding.UnpackGetPegOutConfiguration)
		})
	if err != nil {
		return blockchain.PegOutConfiguration{}, err
	}
	return mapPegOutConfiguration(result), nil
}

func mapPegOutConfiguration(raw bindings.IFlyoverConfigurationsPegOutConfiguration) blockchain.PegOutConfiguration {
	tiers := make([]blockchain.ConfirmationTier, 0, len(raw.ConfirmationTiers))
	for _, tier := range raw.ConfirmationTiers {
		tiers = append(tiers, blockchain.ConfirmationTier{
			MaxAmount:     bigIntToWei(tier.MaxAmount),
			Confirmations: bigIntToUint64(tier.Confirmations),
		})
	}
	return blockchain.PegOutConfiguration{
		FixedFee:          bigIntToWei(raw.FixedFee),
		PercentageFee:     bigIntToUint64(raw.PercentageFee),
		MinAmount:         bigIntToWei(raw.MinAmount),
		MaxAmount:         bigIntToWei(raw.MaxAmount),
		ConfirmationTiers: tiers,
		PenaltyFee:        bigIntToWei(raw.PenaltyFee),
		ClaimWindow:       bigIntToUint64(raw.ClaimWindow),
		ClaimWindowBlocks: bigIntToUint64(raw.ClaimWindowBlocks),
		CallTime:          bigIntToUint64(raw.CallTime),
		ExpireTime:        bigIntToUint64(raw.ExpireTime),
		ExpireBlocks:      bigIntToUint64(raw.ExpireBlocks),
		MaxMinerFee:       bigIntToWei(raw.MaxMinerFee),
	}
}
