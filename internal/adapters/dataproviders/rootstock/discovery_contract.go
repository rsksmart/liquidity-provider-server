package rootstock

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"math/big"
	"strings"
	"time"
)

type discoveryContractImpl struct {
	client        RpcClientBinding
	address       string
	contract      DiscoveryBinding
	signer        TransactionSigner
	retryParams   RetryParams
	miningTimeout time.Duration
	abis          *FlyoverABIs
}

func NewDiscoveryContractImpl(
	client *RskClient,
	address string,
	contract DiscoveryBinding,
	signer TransactionSigner,
	retryParams RetryParams,
	miningTimeout time.Duration,
	abis *FlyoverABIs,
) blockchain.DiscoveryContract {
	return &discoveryContractImpl{
		client:        client.client,
		address:       address,
		contract:      contract,
		signer:        signer,
		retryParams:   retryParams,
		miningTimeout: miningTimeout,
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
				return discovery.contract.SetProviderStatus(opts, parsedId, newStatus)
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
	const lbcProviderNotRegisteredError = "ProviderNotRegistered"

	if !common.IsHexAddress(address) {
		return liquidity_provider.RegisteredLiquidityProvider{}, blockchain.InvalidAddressError
	}

	opts := &bind.CallOpts{}
	provider, revert := discovery.contract.GetProvider(opts, common.HexToAddress(address))
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
			return discovery.contract.GetProviders(opts)
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
		return discovery.contract.UpdateProvider(opts, name, url)
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
				return discovery.contract.Register(opts, params.Name, params.ApiBaseUrl, params.Status, parsedProviderType)
			})
		})

	if err != nil {
		return 0, fmt.Errorf("error registering provider: %w", err)
	} else if receipt == nil || receipt.Status == 0 || len(receipt.Logs) == 0 {
		return 0, errors.New("error registering provider: incomplete receipt")
	}

	registerEvent, err := discovery.contract.ParseRegister(*receipt.Logs[0])
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
			return discovery.contract.IsOperational(opts, parsedProviderType, parsedAddress)
		})
}

func (discovery *discoveryContractImpl) toContractProviderType(providerType liquidity_provider.ProviderType) (uint8, error) {
	if !providerType.IsValid() {
		return 0, liquidity_provider.InvalidProviderTypeError
	}
	return uint8(providerType), nil
}
