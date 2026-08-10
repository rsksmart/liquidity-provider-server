package rootstock

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_address_registry"
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

func (registry *peginAddressRegistryContractImpl) GetPegInAddress(rskAddr string) (blockchain.PegInAddress, error) {
	var parsedAddress common.Address
	if err := ParseAddress(&parsedAddress, rskAddr); err != nil {
		return blockchain.PegInAddress{}, err
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
		return blockchain.PegInAddress{}, err
	}
	return blockchain.PegInAddress{
		Payload:  result.Payload,
		Encoding: blockchain.PegInAddressRegistryEncoding(result.Encoding),
	}, nil
}

func (registry *peginAddressRegistryContractImpl) GetPegInAddresses(rskAddrs []string) (blockchain.PegInAddressBatch, error) {
	parsedAddresses := make([]common.Address, len(rskAddrs))
	for i, rskAddr := range rskAddrs {
		if err := ParseAddress(&parsedAddresses[i], rskAddr); err != nil {
			return blockchain.PegInAddressBatch{}, err
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
		return blockchain.PegInAddressBatch{}, err
	}
	return blockchain.PegInAddressBatch{
		Payloads: result.Payloads,
		Encoding: blockchain.PegInAddressRegistryEncoding(result.Encoding),
	}, nil
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
		RegistrationBlock: result.RegistrationBlock.Uint64(),
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

func (registry *peginAddressRegistryContractImpl) GetRegisteredBtcTransactionHash(
	ctx context.Context,
	registrationTxHash string,
	rskAddr string,
) (string, error) {
	hashBytes, err := hex.DecodeString(strings.TrimPrefix(registrationTxHash, "0x"))
	if err != nil || len(hashBytes) != common.HashLength {
		return "", errors.New("invalid registration transaction hash")
	}
	var expectedRskAddress common.Address
	if err = ParseAddress(&expectedRskAddress, rskAddr); err != nil {
		return "", err
	}
	var registryAddress common.Address
	if err = ParseAddress(&registryAddress, registry.address); err != nil {
		return "", fmt.Errorf("invalid PegIn address registry contract address: %w", err)
	}

	transaction, err := registry.getRegistrationTransaction(ctx, common.BytesToHash(hashBytes), registryAddress)
	if err != nil {
		return "", err
	}
	serializedBtcTransaction, err := registry.registeredBtcTransactionBytes(transaction, expectedRskAddress)
	if err != nil {
		return "", err
	}
	return registeredBtcTransactionHash(serializedBtcTransaction)
}

func (registry *peginAddressRegistryContractImpl) getRegistrationTransaction(
	ctx context.Context,
	transactionHash common.Hash,
	registryAddress common.Address,
) (*types.Transaction, error) {
	transaction, err := rskRetry(registry.retryParams.Retries, registry.retryParams.Sleep,
		func() (*types.Transaction, error) {
			tx, pending, rpcErr := registry.client.TransactionByHash(ctx, transactionHash)
			if rpcErr != nil {
				return nil, rpcErr
			}
			if pending {
				return nil, errors.New("registration transaction is still pending")
			}
			return tx, nil
		})
	if err != nil {
		return nil, err
	}
	if transaction == nil || transaction.To() == nil || *transaction.To() != registryAddress {
		return nil, errors.New("registration transaction does not target the PegIn address registry")
	}
	return transaction, nil
}

func (registry *peginAddressRegistryContractImpl) registeredBtcTransactionBytes(
	transaction *types.Transaction,
	expectedRskAddress common.Address,
) ([]byte, error) {
	method, ok := registry.abis.PegInAddressRegistry.Methods["registerAddress"]
	if !ok {
		return nil, errors.New("registerAddress method missing from PegIn address registry ABI")
	}
	data := transaction.Data()
	if len(data) < len(method.ID) || !bytes.Equal(data[:len(method.ID)], method.ID) {
		return nil, errors.New("registration transaction does not call registerAddress")
	}
	arguments, err := method.Inputs.Unpack(data[len(method.ID):])
	if err != nil {
		return nil, fmt.Errorf("decode registerAddress input: %w", err)
	}
	if len(arguments) < 2 {
		return nil, errors.New("registerAddress input has insufficient arguments")
	}
	registeredRskAddress, ok := arguments[0].(common.Address)
	if !ok || registeredRskAddress != expectedRskAddress {
		return nil, errors.New("registerAddress input RSK address does not match event")
	}
	serializedBtcTransaction, ok := arguments[1].([]byte)
	if !ok {
		return nil, errors.New("registerAddress input has invalid BTC transaction")
	}
	return serializedBtcTransaction, nil
}

func registeredBtcTransactionHash(serializedBtcTransaction []byte) (string, error) {
	var btcTransaction wire.MsgTx
	reader := bytes.NewReader(serializedBtcTransaction)
	if err := btcTransaction.Deserialize(reader); err != nil {
		return "", fmt.Errorf("deserialize registered BTC transaction: %w", err)
	}
	if reader.Len() != 0 {
		return "", errors.New("registered BTC transaction has trailing bytes")
	}
	return btcTransaction.TxHash().String(), nil
}

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
			LogIndex:         lbcEvent.Raw.Index,
		})
	}
	if err = iterator.Error(); err != nil {
		return nil, err
	}

	return result, nil
}
