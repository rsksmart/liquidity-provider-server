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

// NewFlyoverConfigurationsContractImpl builds the read-only adapter for the frozen
// IFlyoverConfigurations ABI. Only the fee/confirmation reads are exposed; getPegInConfiguration
// and the time-locked admin writes (queueChange/applyChange) are out of scope for this adapter.
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
