package rootstock

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_address_registry"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
)

type peginAddressRegistryContractImpl struct {
	client      RpcClientBinding
	address     string
	contract    *bind.BoundContract
	retryParams RetryParams
	abis        *FlyoverABIs
	binding     *bindings.PegInAddressRegistryContract
}

// NewPegInAddressRegistryContractImpl builds the read-only adapter for the frozen
// IPegInAddressRegistry ABI. Registration (registerAddress) is deliberately not adapted
// here: writing registrations belongs to a separate on-chain watcher process, not the
// liquidity provider server.
func NewPegInAddressRegistryContractImpl(
	client *RskClient,
	address string,
	contract *bind.BoundContract,
	retryParams RetryParams,
	binding *bindings.PegInAddressRegistryContract,
	abis *FlyoverABIs,
) blockchain.PegInAddressRegistryContract {
	return &peginAddressRegistryContractImpl{
		client:      client.client,
		address:     address,
		contract:    contract,
		retryParams: retryParams,
		binding:     binding,
		abis:        abis,
	}
}

func (registry *peginAddressRegistryContractImpl) GetAddress() string {
	return registry.address
}

func (registry *peginAddressRegistryContractImpl) GetPegInAddress(rskAddr string) ([]byte, blockchain.PegInAddressRegistryEncoding, error) {
	var parsedAddress common.Address
	if err := ParseAddress(&parsedAddress, rskAddr); err != nil {
		return nil, 0, err
	}
	opts := &bind.CallOpts{}
	result, err := rskRetry(registry.retryParams.Retries, registry.retryParams.Sleep,
		func() (bindings.GetPegInAddressOutput, error) {
			callData, dataErr := registry.binding.TryPackGetPegInAddress(parsedAddress)
			if dataErr != nil {
				return bindings.GetPegInAddressOutput{}, dataErr
			}
			return bind.Call(registry.contract, opts, callData, registry.binding.UnpackGetPegInAddress)
		})
	if err != nil {
		return nil, 0, err
	}
	return result.Payload, blockchain.PegInAddressRegistryEncoding(result.Encoding), nil
}

func (registry *peginAddressRegistryContractImpl) GetPegInAddresses(rskAddrs []string) ([][]byte, blockchain.PegInAddressRegistryEncoding, error) {
	parsedAddresses := make([]common.Address, len(rskAddrs))
	for i, rskAddr := range rskAddrs {
		if err := ParseAddress(&parsedAddresses[i], rskAddr); err != nil {
			return nil, 0, err
		}
	}
	opts := &bind.CallOpts{}
	result, err := rskRetry(registry.retryParams.Retries, registry.retryParams.Sleep,
		func() (bindings.GetPegInAddressesOutput, error) {
			callData, dataErr := registry.binding.TryPackGetPegInAddresses(parsedAddresses)
			if dataErr != nil {
				return bindings.GetPegInAddressesOutput{}, dataErr
			}
			return bind.Call(registry.contract, opts, callData, registry.binding.UnpackGetPegInAddresses)
		})
	if err != nil {
		return nil, 0, err
	}
	return result.Payloads, blockchain.PegInAddressRegistryEncoding(result.Encoding), nil
}

func (registry *peginAddressRegistryContractImpl) IsRegistered(rskAddr string) (bool, error) {
	var parsedAddress common.Address
	if err := ParseAddress(&parsedAddress, rskAddr); err != nil {
		return false, err
	}
	opts := &bind.CallOpts{}
	return rskRetry(registry.retryParams.Retries, registry.retryParams.Sleep,
		func() (bool, error) {
			callData, dataErr := registry.binding.TryPackIsRegistered(parsedAddress)
			if dataErr != nil {
				return false, dataErr
			}
			return bind.Call(registry.contract, opts, callData, registry.binding.UnpackIsRegistered)
		})
}

func (registry *peginAddressRegistryContractImpl) GetRegistration(rskAddr string) (blockchain.PegInRegistration, error) {
	var parsedAddress common.Address
	if err := ParseAddress(&parsedAddress, rskAddr); err != nil {
		return blockchain.PegInRegistration{}, err
	}
	opts := &bind.CallOpts{}
	result, err := rskRetry(registry.retryParams.Retries, registry.retryParams.Sleep,
		func() (bindings.IPegInAddressRegistryRegistration, error) {
			callData, dataErr := registry.binding.TryPackGetRegistration(parsedAddress)
			if dataErr != nil {
				return bindings.IPegInAddressRegistryRegistration{}, dataErr
			}
			return bind.Call(registry.contract, opts, callData, registry.binding.UnpackGetRegistration)
		})
	if err != nil {
		return blockchain.PegInRegistration{}, err
	}
	return blockchain.PegInRegistration{
		Registrant:        result.Registrant.String(),
		RegistrationBlock: entities.NewBigWei(result.RegistrationBlock),
	}, nil
}

func (registry *peginAddressRegistryContractImpl) GetRegistrationRoot() ([32]byte, error) {
	opts := &bind.CallOpts{}
	return rskRetry(registry.retryParams.Retries, registry.retryParams.Sleep,
		func() ([32]byte, error) {
			callData, dataErr := registry.binding.TryPackGetRegistrationRoot()
			if dataErr != nil {
				return [32]byte{}, dataErr
			}
			return bind.Call(registry.contract, opts, callData, registry.binding.UnpackGetRegistrationRoot)
		})
}

// GetAddressRegisteredEvents replicates pegoutContractImpl.GetDepositEvents' filter-and-decode
// pattern for the registry's AddressRegistered event. It only reads and returns matching
// events; watching for new ones and reacting to them is the caller's responsibility.
func (registry *peginAddressRegistryContractImpl) GetAddressRegisteredEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]blockchain.AddressRegistered, error) {
	var lbcEvent *bindings.PegInAddressRegistryContractAddressRegistered
	result := make([]blockchain.AddressRegistered, 0)

	iterator, err := bind.FilterEvents(
		registry.contract,
		&bind.FilterOpts{
			Start:   fromBlock,
			End:     toBlock,
			Context: ctx,
		},
		registry.binding.UnpackAddressRegisteredEvent,
	)

	defer func() {
		if iterator == nil {
			return
		}
		if iteratorError := iterator.Close(); iteratorError != nil {
			log.Error("Error closing AddressRegistered event iterator: ", iteratorError)
		}
	}()
	if err != nil || iterator == nil {
		return nil, err
	}

	for iterator.Next() {
		lbcEvent = iterator.Value()
		result = append(result, blockchain.AddressRegistered{
			RskAddress:       lbcEvent.RskAddr.String(),
			Registrant:       lbcEvent.Registrant.String(),
			RegistrationRoot: lbcEvent.RegistrationRoot,
			TxHash:           lbcEvent.Raw.TxHash.String(),
			BlockNumber:      lbcEvent.Raw.BlockNumber,
		})
	}
	if err = iterator.Error(); err != nil {
		return nil, err
	}

	return result, nil
}
