package rootstock

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	registry "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
)

// peginAddressRegistryImpl is the adapter for the new PegInAddressRegistry contract (EPIC E2).
// It is read-only from the LPS perspective: LPS reads derived deposit addresses and consumes the
// AddressRegistered event stream to drive commit-first discovery (EPIC E5). Registration itself is
// performed by users / watchtowers, not by LPS.
type peginAddressRegistryImpl struct {
	client      RpcClientBinding
	address     string
	contract    *bind.BoundContract
	retryParams RetryParams
	binding     *registry.PegInAddressRegistry
}

func NewPegInAddressRegistryImpl(
	client *RskClient,
	address string,
	contract *bind.BoundContract,
	retryParams RetryParams,
	binding *registry.PegInAddressRegistry,
) blockchain.PegInAddressRegistryContract {
	return &peginAddressRegistryImpl{
		client:      client.client,
		address:     address,
		contract:    contract,
		retryParams: retryParams,
		binding:     binding,
	}
}

func (r *peginAddressRegistryImpl) GetAddress() string {
	return r.address
}

func (r *peginAddressRegistryImpl) IsRegistered(rskAddress string) (bool, error) {
	var parsed common.Address
	if err := ParseAddress(&parsed, rskAddress); err != nil {
		return false, err
	}
	return rskRetry(r.retryParams.Retries, r.retryParams.Sleep, func() (bool, error) {
		callData, dataErr := r.binding.TryPackIsRegistered(parsed)
		if dataErr != nil {
			return false, dataErr
		}
		return bind.Call(r.contract, &bind.CallOpts{}, callData, r.binding.UnpackIsRegistered)
	})
}

func (r *peginAddressRegistryImpl) GetPegInAddress(rskAddress string) (blockchain.PegInDepositAddress, error) {
	var parsed common.Address
	if err := ParseAddress(&parsed, rskAddress); err != nil {
		return blockchain.PegInDepositAddress{}, err
	}
	out, err := rskRetry(r.retryParams.Retries, r.retryParams.Sleep, func() (registry.GetPegInAddressOutput, error) {
		callData, dataErr := r.binding.TryPackGetPegInAddress(parsed)
		if dataErr != nil {
			return registry.GetPegInAddressOutput{}, dataErr
		}
		return bind.Call(r.contract, &bind.CallOpts{}, callData, r.binding.UnpackGetPegInAddress)
	})
	if err != nil {
		return blockchain.PegInDepositAddress{}, err
	}
	return blockchain.PegInDepositAddress{
		RskAddress:      rskAddress,
		ScriptOrAddress: out.Arg0,
		Encoding:        blockchain.PegInAddressEncoding(out.Arg1),
	}, nil
}

func (r *peginAddressRegistryImpl) GetPegInAddresses(rskAddresses []string) ([]blockchain.PegInDepositAddress, error) {
	parsed := make([]common.Address, 0, len(rskAddresses))
	for _, a := range rskAddresses {
		var addr common.Address
		if err := ParseAddress(&addr, a); err != nil {
			return nil, err
		}
		parsed = append(parsed, addr)
	}
	out, err := rskRetry(r.retryParams.Retries, r.retryParams.Sleep, func() (registry.GetPegInAddressesOutput, error) {
		callData, dataErr := r.binding.TryPackGetPegInAddresses(parsed)
		if dataErr != nil {
			return registry.GetPegInAddressesOutput{}, dataErr
		}
		return bind.Call(r.contract, &bind.CallOpts{}, callData, r.binding.UnpackGetPegInAddresses)
	})
	if err != nil {
		return nil, err
	}
	result := make([]blockchain.PegInDepositAddress, 0, len(out.DerivationAddresses))
	for i, raw := range out.DerivationAddresses {
		rsk := ""
		if i < len(rskAddresses) {
			rsk = rskAddresses[i]
		}
		result = append(result, blockchain.PegInDepositAddress{
			RskAddress:      rsk,
			ScriptOrAddress: raw,
			Encoding:        blockchain.PegInAddressEncoding(out.Encoding),
		})
	}
	return result, nil
}

func (r *peginAddressRegistryImpl) GetRegistrationBlock(rskAddress string) (uint64, error) {
	var parsed common.Address
	if err := ParseAddress(&parsed, rskAddress); err != nil {
		return 0, err
	}
	block, err := rskRetry(r.retryParams.Retries, r.retryParams.Sleep, func() (*big.Int, error) {
		callData, dataErr := r.binding.TryPackGetRegistrationBlock(parsed)
		if dataErr != nil {
			return nil, dataErr
		}
		return bind.Call(r.contract, &bind.CallOpts{}, callData, r.binding.UnpackGetRegistrationBlock)
	})
	if err != nil {
		return 0, err
	}
	return block.Uint64(), nil
}

func (r *peginAddressRegistryImpl) GetRegistrationRoot() ([32]byte, error) {
	return rskRetry(r.retryParams.Retries, r.retryParams.Sleep, func() ([32]byte, error) {
		callData, dataErr := r.binding.TryPackGetRegistrationRoot()
		if dataErr != nil {
			return [32]byte{}, dataErr
		}
		return bind.Call(r.contract, &bind.CallOpts{}, callData, r.binding.UnpackGetRegistrationRoot)
	})
}

// GetAddressRegisteredEvents polls the AddressRegistered logs over the given block range using
// the higher-level bind.FilterEvents helper (the same pattern used by GetDepositEvents on the
// pegout contract). It is the discovery mechanism for the commit-first peg-in path (EPIC E5).
func (r *peginAddressRegistryImpl) GetAddressRegisteredEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]blockchain.AddressRegisteredEvent, error) {
	result := make([]blockchain.AddressRegisteredEvent, 0)
	iterator, err := bind.FilterEvents(
		r.contract,
		&bind.FilterOpts{Start: fromBlock, End: toBlock, Context: ctx},
		r.binding.UnpackAddressRegisteredEvent,
	)
	defer func() {
		if iterator != nil {
			if closeErr := iterator.Close(); closeErr != nil {
				log.Error("error closing AddressRegistered iterator: ", closeErr)
			}
		}
	}()
	if err != nil || iterator == nil {
		return nil, err
	}
	for iterator.Next() {
		ev := iterator.Value()
		result = append(result, blockchain.AddressRegisteredEvent{
			RskAddress:       strings.ToLower(ev.Addr.String()),
			RegistrationRoot: ev.RegistrationRoot,
			BlockNumber:      ev.Raw.BlockNumber,
			TxHash:           ev.Raw.TxHash.String(),
		})
	}
	if err = iterator.Error(); err != nil {
		return nil, fmt.Errorf("error iterating AddressRegistered events: %w", err)
	}
	return result, nil
}
