package rootstock

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/discovery"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"math/big"
	"strings"
	"time"
)

type discoveryContractImpl struct {
	client        RpcClientBinding
	address       string
	contract      *bind.BoundContract
	signer        TransactionSigner
	retryParams   RetryParams
	miningTimeout time.Duration
	abis          *FlyoverABIs
	binding       *bindings.FlyoverDiscovery
}

func NewDiscoveryContractImpl(
	client *RskClient,
	address string,
	contract *bind.BoundContract,
	signer TransactionSigner,
	retryParams RetryParams,
	miningTimeout time.Duration,
	binding *bindings.FlyoverDiscovery,
	abis *FlyoverABIs,
) blockchain.DiscoveryContract {
	return &discoveryContractImpl{
		client:        client.client,
		address:       address,
		contract:      contract,
		signer:        signer,
		retryParams:   retryParams,
		miningTimeout: miningTimeout,
		binding:       binding,
		abis:          abis,
	}
}

func (discovery *discoveryContractImpl) GetAddress() string {
	return discovery.address
}

func (discovery *discoveryContractImpl) SetProviderStatus(id uint64, newStatus bool) error {
	parsedId := new(big.Int)
	parsedId.SetUint64(id)
	opts := &bind.TransactOpts{
		From:   discovery.signer.Address(),
		Signer: discovery.signer.Sign,
	}

	receipt, err := rskRetry(discovery.retryParams.Retries, discovery.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(discovery.client, discovery.miningTimeout, "SetProviderStatus", func() (*geth.Transaction, error) {
				callData, dataErr := discovery.binding.TryPackSetProviderStatus(parsedId, newStatus)
				if dataErr != nil {
					return nil, dataErr
				}
				return bind.Transact(discovery.contract, opts, callData)
			})
		})

	if err != nil {
		return err
	} else if receipt == nil || receipt.Status == 0 {
		return errors.New("setProviderStatus transaction failed")
	}
	return nil
}

func (discovery *discoveryContractImpl) GetProvider(address string) (liquidity_provider.RegisteredLiquidityProvider, error) {
	var providerType liquidity_provider.ProviderType

	if !common.IsHexAddress(address) {
		return liquidity_provider.RegisteredLiquidityProvider{}, blockchain.InvalidAddressError
	}

	opts := &bind.CallOpts{}
	callData, dataErr := discovery.binding.TryPackGetProvider(common.HexToAddress(address))
	if dataErr != nil {
		return liquidity_provider.RegisteredLiquidityProvider{}, dataErr
	}
	provider, revert := bind.Call(discovery.contract, opts, callData, discovery.binding.UnpackGetProvider)
	parsedRevert, err := ParseRevertReason(discovery.abis.Flyover, revert)
	if err != nil && parsedRevert == nil {
		return liquidity_provider.RegisteredLiquidityProvider{}, fmt.Errorf("error parsing getProvider result: %w", err)
	} else if parsedRevert != nil && strings.EqualFold(lbcProviderNotRegisteredError, parsedRevert.Name) {
		return liquidity_provider.RegisteredLiquidityProvider{}, liquidity_provider.ProviderNotFoundError
	} else if parsedRevert != nil {
		return liquidity_provider.RegisteredLiquidityProvider{}, fmt.Errorf("getProvider reverted with: %s", parsedRevert.Name)
	}

	providerType = liquidity_provider.ProviderType(provider.ProviderType)
	if !providerType.IsValid() {
		return liquidity_provider.RegisteredLiquidityProvider{}, liquidity_provider.InvalidProviderTypeError
	}
	return liquidity_provider.RegisteredLiquidityProvider{
		Id:           provider.Id.Uint64(),
		Address:      provider.ProviderAddress.String(),
		Name:         provider.Name,
		ApiBaseUrl:   provider.ApiBaseUrl,
		Status:       provider.Status,
		ProviderType: providerType,
	}, nil
}

func (discovery *discoveryContractImpl) GetProviders() ([]liquidity_provider.RegisteredLiquidityProvider, error) {
	var providerType liquidity_provider.ProviderType
	var providers []bindings.FlyoverLiquidityProvider

	opts := &bind.CallOpts{}
	providers, err := rskRetry(discovery.retryParams.Retries, discovery.retryParams.Sleep,
		func() ([]bindings.FlyoverLiquidityProvider, error) {
			callData, dataErr := discovery.binding.TryPackGetProviders()
			if dataErr != nil {
				return nil, dataErr
			}
			return bind.Call(discovery.contract, opts, callData, discovery.binding.UnpackGetProviders)
		})
	if err != nil {
		return nil, err
	}
	parsedProviders := make([]liquidity_provider.RegisteredLiquidityProvider, 0)
	for _, provider := range providers {
		providerType = liquidity_provider.ProviderType(provider.ProviderType)
		if !providerType.IsValid() {
			return nil, liquidity_provider.InvalidProviderTypeError
		}
		parsedProviders = append(parsedProviders, liquidity_provider.RegisteredLiquidityProvider{
			Id:           provider.Id.Uint64(),
			Address:      provider.ProviderAddress.String(),
			Name:         provider.Name,
			ApiBaseUrl:   provider.ApiBaseUrl,
			Status:       provider.Status,
			ProviderType: providerType,
		})
	}
	return parsedProviders, nil
}

func (discovery *discoveryContractImpl) UpdateProvider(name, url string) (string, error) {
	opts := &bind.TransactOpts{From: discovery.signer.Address(), Signer: discovery.signer.Sign}
	receipt, err := awaitTx(discovery.client, discovery.miningTimeout, "UpdateProvider", func() (*geth.Transaction, error) {
		callData, dataErr := discovery.binding.TryPackUpdateProvider(name, url)
		if dataErr != nil {
			return nil, dataErr
		}
		return bind.Transact(discovery.contract, opts, callData)
	})

	if err != nil {
		return "", fmt.Errorf("update provider error: %w", err)
	} else if receipt == nil {
		return "", errors.New("update provider error: incomplete receipt")
	} else if receipt.Status == 0 {
		txHash := receipt.TxHash.String()
		return txHash, fmt.Errorf("update provider error: transaction reverted (%s)", txHash)
	}
	return receipt.TxHash.String(), nil
}

func (discovery *discoveryContractImpl) RegisterProvider(txConfig blockchain.TransactionConfig, params blockchain.ProviderRegistrationParams) (int64, error) {
	var err error

	opts := &bind.TransactOpts{
		Value:  txConfig.Value.AsBigInt(),
		From:   discovery.signer.Address(),
		Signer: discovery.signer.Sign,
	}

	parsedProviderType, err := discovery.toContractProviderType(params.Type)
	if err != nil {
		return 0, err
	}
	receipt, err := rskRetry(discovery.retryParams.Retries, discovery.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(discovery.client, discovery.miningTimeout, "Register", func() (*geth.Transaction, error) {
				callData, dataErr := discovery.binding.TryPackRegister(params.Name, params.ApiBaseUrl, params.Status, parsedProviderType)
				if dataErr != nil {
					return nil, dataErr
				}
				return bind.Transact(discovery.contract, opts, callData)
			})
		})

	if err != nil {
		return 0, fmt.Errorf("error registering provider: %w", err)
	} else if receipt == nil || receipt.Status == 0 || len(receipt.Logs) == 0 {
		return 0, errors.New("error registering provider: incomplete receipt")
	}

	registerEvent, err := discovery.binding.UnpackRegisterEvent(receipt.Logs[0])
	if err != nil {
		return 0, fmt.Errorf("error parsing register event: %w", err)
	}
	return registerEvent.Id.Int64(), nil
}

func (discovery *discoveryContractImpl) IsOperational(providerType liquidity_provider.ProviderType, address string) (bool, error) {
	var parsedAddress common.Address
	var err error
	opts := &bind.CallOpts{}

	if err = ParseAddress(&parsedAddress, address); err != nil {
		return false, err
	}

	parsedProviderType, err := discovery.toContractProviderType(providerType)
	if err != nil {
		return false, err
	}

	return rskRetry(discovery.retryParams.Retries, discovery.retryParams.Sleep,
		func() (bool, error) {
			callData, dataErr := discovery.binding.TryPackIsOperational(parsedProviderType, parsedAddress)
			if dataErr != nil {
				return false, dataErr
			}
			result, revert := bind.Call(discovery.contract, opts, callData, discovery.binding.UnpackIsOperational)
			parsedRevert, parseErr := ParseRevertReason(discovery.abis.Flyover, revert)
			if parseErr != nil && parsedRevert == nil {
				return false, fmt.Errorf("error parsing IsOperational result: %w", err)
			} else if parsedRevert != nil && strings.EqualFold(lbcProviderNotRegisteredError, parsedRevert.Name) {
				return false, nil
			} else if parsedRevert != nil {
				return false, fmt.Errorf("IsOperational reverted with: %s", parsedRevert.Name)
			}
			return result, revert
		})
}

func (discovery *discoveryContractImpl) PausedStatus() (blockchain.PauseStatus, error) {
	opts := new(bind.CallOpts)
	result, err := rskRetry(
		discovery.retryParams.Retries,
		discovery.retryParams.Sleep,
		func() (bindings.PauseStatusOutput, error) {
			callData, dataErr := discovery.binding.TryPackPauseStatus()
			if dataErr != nil {
				return bindings.PauseStatusOutput{}, dataErr
			}
			return bind.Call(discovery.contract, opts, callData, discovery.binding.UnpackPauseStatus)
		},
	)
	if err != nil {
		return blockchain.PauseStatus{}, err
	}
	return blockchain.PauseStatus{
		IsPaused: result.IsPaused,
		Reason:   result.Reason,
		Since:    result.Since,
	}, nil
}

func (discovery *discoveryContractImpl) toContractProviderType(providerType liquidity_provider.ProviderType) (uint8, error) {
	if !providerType.IsValid() {
		return 0, liquidity_provider.InvalidProviderTypeError
	}
	return uint8(providerType), nil
}
