package rootstock

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	emergencyPause "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/emergency_pause"
	pauseRegistry "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pause_registry"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

type pauseRegistryContractImpl struct {
	client          RpcClientBinding
	pausable        *bind.BoundContract
	retryParams     RetryParams
	emergencyPause  *emergencyPause.EmergencyPauseContract
	registryBinding *pauseRegistry.PauseRegistryContract
}

func NewPauseRegistryContractImpl(
	client *RskClient,
	pausable *bind.BoundContract,
	retryParams RetryParams,
	emergencyPauseBinding *emergencyPause.EmergencyPauseContract,
	registryBinding *pauseRegistry.PauseRegistryContract,
) blockchain.PauseRegistryContract {
	return &pauseRegistryContractImpl{
		client:          client.client,
		pausable:        pausable,
		retryParams:     retryParams,
		emergencyPause:  emergencyPauseBinding,
		registryBinding: registryBinding,
	}
}

func (registry *pauseRegistryContractImpl) PauseLevel() (blockchain.PauseLevel, error) {
	opts := &bind.CallOpts{}
	level, err := rskRetry(registry.retryParams.Retries, registry.retryParams.Sleep,
		func() (uint8, error) {
			registryCallData, dataErr := registry.emergencyPause.TryPackPauseRegistry()
			if dataErr != nil {
				return 0, dataErr
			}
			address, addressErr := bind.Call(registry.pausable, opts, registryCallData, registry.emergencyPause.UnpackPauseRegistry)
			if addressErr != nil {
				return 0, addressErr
			}
			levelCallData, dataErr := registry.registryBinding.TryPackPauseLevel()
			if dataErr != nil {
				return 0, dataErr
			}
			return bind.Call(registry.registryBinding.Instance(registry.client, address), opts, levelCallData, registry.registryBinding.UnpackPauseLevel)
		})
	if err != nil {
		return blockchain.PauseLevelNone, err
	}
	return blockchain.PauseLevel(level), nil
}
