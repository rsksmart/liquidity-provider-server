package rootstock

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pause_registry"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

type pauseRegistryContractImpl struct {
	client      RpcClientBinding
	address     string
	contract    *bind.BoundContract
	retryParams RetryParams
	abis        *FlyoverABIs
	binding     *bindings.PauseRegistryContract
}

func NewPauseRegistryContractImpl(
	client *RskClient,
	address string,
	contract *bind.BoundContract,
	retryParams RetryParams,
	binding *bindings.PauseRegistryContract,
	abis *FlyoverABIs,
) blockchain.PauseRegistryContract {
	return &pauseRegistryContractImpl{
		client:      client.client,
		address:     address,
		contract:    contract,
		retryParams: retryParams,
		binding:     binding,
		abis:        abis,
	}
}

func (registry *pauseRegistryContractImpl) GetAddress() string {
	return registry.address
}

func (registry *pauseRegistryContractImpl) PauseLevel() (uint8, error) {
	opts := &bind.CallOpts{}
	return rskRetry(registry.retryParams.Retries, registry.retryParams.Sleep,
		func() (uint8, error) {
			callData, dataErr := registry.binding.TryPackPauseLevel()
			if dataErr != nil {
				return 0, dataErr
			}
			return bind.Call(registry.contract, opts, callData, registry.binding.UnpackPauseLevel)
		})
}

func (registry *pauseRegistryContractImpl) PauseStatus() (blockchain.PauseStatus, error) {
	opts := &bind.CallOpts{}
	result, err := rskRetry(registry.retryParams.Retries, registry.retryParams.Sleep,
		func() (bindings.PauseStatusOutput, error) {
			callData, dataErr := registry.binding.TryPackPauseStatus()
			if dataErr != nil {
				return bindings.PauseStatusOutput{}, dataErr
			}
			return bind.Call(registry.contract, opts, callData, registry.binding.UnpackPauseStatus)
		})
	if err != nil {
		return blockchain.PauseStatus{}, err
	}
	return blockchain.PauseStatus{
		IsPaused: result.IsPaused,
		Reason:   result.Reason,
		Since:    result.Since,
	}, nil
}
