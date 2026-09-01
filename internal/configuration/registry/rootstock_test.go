package registry_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// newRskWalletFactoryMock returns a wallet factory mock wired to a signer wallet mock that
// resolves to a fixed address.
func newRskWalletFactoryMock() (*mocks.AbstractFactoryMock, *mocks.RskSignerWalletMock) {
	rskWalletMock := new(mocks.RskSignerWalletMock)
	rskWalletMock.On("Address").Return(common.HexToAddress(test.AnyRskAddress))
	walletFactoryMock := new(mocks.AbstractFactoryMock)
	walletFactoryMock.On("RskWallet").Return(rskWalletMock, nil)
	return walletFactoryMock, rskWalletMock
}

func newRskClientWithGenesisRegistry(t *testing.T, registryAddress string, deploymentBlock uint64) *rootstock.RskClient {
	t.Helper()
	rpc := mocks.NewRpcClientBindingMock(t)
	address := common.HexToAddress(registryAddress)
	rpc.On("CodeAt", mock.Anything, address, new(big.Int).SetUint64(deploymentBlock)).Return([]byte{0x60}, nil)
	if deploymentBlock > 0 {
		rpc.On("CodeAt", mock.Anything, address, new(big.Int).SetUint64(deploymentBlock-1)).Return([]byte{}, nil)
	}
	rpc.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(make([]byte, 32), nil)
	rpc.On("FilterLogs", mock.Anything, mock.Anything).Return([]geth.Log{}, nil)
	return rootstock.NewRskClient(rpc)
}

// nolint:funlen
func TestNewRootstockRegistry(t *testing.T) {
	testEnv := environment.Environment{
		Rsk: environment.RskEnv{
			DiscoveryAddress:             "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA8",
			CollateralManagementAddress:  "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA7",
			PeginContractAddress:         "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA6",
			PegoutContractAddress:        "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA5",
			PegInAddressRegistryAddress:  "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA4",
			PauseRegistryAddress:         "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA2",
			FlyoverConfigurationsAddress: "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA3",
			BridgeAddress:                "0x0000000000000000000000000000000001000006",
		},
		Btc: environment.BtcEnv{Network: "testnet"},
	}
	t.Run("should create a new Rootstock registry", func(t *testing.T) {
		env := testEnv
		walletFactoryMock, rskWalletMock := newRskWalletFactoryMock()
		rskClient := newRskClientWithGenesisRegistry(t, env.Rsk.PegInAddressRegistryAddress, 0)
		rskRegistry, err := registry.NewRootstockRegistry(context.Background(), env, rskClient, walletFactoryMock, environment.DefaultTimeouts())
		require.NoError(t, err)
		require.NotNil(t, rskRegistry)
		require.NotEmpty(t, rskRegistry.Contracts)
		require.NotNil(t, rskRegistry.Contracts.Discovery)
		require.NotNil(t, rskRegistry.Contracts.CollateralManagement)
		require.NotNil(t, rskRegistry.Contracts.PegIn)
		require.NotNil(t, rskRegistry.Contracts.PegOut)
		require.NotNil(t, rskRegistry.Contracts.Bridge)
		require.NotNil(t, rskRegistry.Contracts.PegInAddressRegistry)
		require.NotNil(t, rskRegistry.Contracts.PauseRegistry)
		require.NotNil(t, rskRegistry.Contracts.FlyoverConfigurations)
		require.Equal(t, env.Rsk.PegInAddressRegistryAddress, rskRegistry.Contracts.PegInAddressRegistry.GetAddress())
		require.Equal(t, env.Rsk.PauseRegistryAddress, rskRegistry.Contracts.PauseRegistry.GetAddress())
		require.Equal(t, env.Rsk.FlyoverConfigurationsAddress, rskRegistry.Contracts.FlyoverConfigurations.GetAddress())
		require.Equal(t, rskWalletMock, rskRegistry.Wallet)
		require.Equal(t, rskClient, rskRegistry.Client)
	})
	t.Run("should return an error when the pegin address registry address is invalid", func(t *testing.T) {
		env := testEnv
		env.Rsk.PegInAddressRegistryAddress = test.AnyString
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(context.Background(), env, rskClient, new(mocks.AbstractFactoryMock), environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the pause registry address is invalid", func(t *testing.T) {
		env := testEnv
		env.Rsk.PauseRegistryAddress = test.AnyString
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(context.Background(), env, rskClient, new(mocks.AbstractFactoryMock), environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the flyover configurations address is invalid", func(t *testing.T) {
		env := testEnv
		env.Rsk.FlyoverConfigurationsAddress = test.AnyString
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(context.Background(), env, rskClient, new(mocks.AbstractFactoryMock), environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the discovery contract address is invalid", func(t *testing.T) {
		env := testEnv
		env.Rsk.DiscoveryAddress = test.AnyString
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(context.Background(), env, rskClient, new(mocks.AbstractFactoryMock), environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the pegin contract address is invalid", func(t *testing.T) {
		env := testEnv
		env.Rsk.PeginContractAddress = test.AnyString
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(context.Background(), env, rskClient, new(mocks.AbstractFactoryMock), environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the pegout contract address is invalid", func(t *testing.T) {
		env := testEnv
		env.Rsk.PegoutContractAddress = test.AnyString
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(context.Background(), env, rskClient, new(mocks.AbstractFactoryMock), environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the collateral management address is invalid", func(t *testing.T) {
		env := testEnv
		env.Rsk.CollateralManagementAddress = test.AnyString
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(context.Background(), env, rskClient, new(mocks.AbstractFactoryMock), environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the bridge address is invalid", func(t *testing.T) {
		env := testEnv
		env.Rsk.BridgeAddress = test.AnyString
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(context.Background(), env, rskClient, new(mocks.AbstractFactoryMock), environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the wallet factory fails", func(t *testing.T) {
		env := testEnv
		walletFactoryMock := new(mocks.AbstractFactoryMock)
		walletFactoryMock.On("RskWallet").Return(nil, assert.AnError)
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(context.Background(), env, rskClient, walletFactoryMock, environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the BTC network params cannot be retrieved", func(t *testing.T) {
		env := testEnv
		env.Btc.Network = test.AnyString
		walletFactoryMock := new(mocks.AbstractFactoryMock)
		walletFactoryMock.On("RskWallet").Return(new(mocks.RskSignerWalletMock), nil)
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(context.Background(), env, rskClient, walletFactoryMock, environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the watcher start block is not the deployment block", func(t *testing.T) {
		env := testEnv
		env.Pegin.AddressRegistryWatcherStartBlock = 100
		env.Pegin.AddressRegistryWatcherPageSize = 10
		rpc := mocks.NewRpcClientBindingMock(t)
		address := common.HexToAddress(env.Rsk.PegInAddressRegistryAddress)
		rpc.EXPECT().CodeAt(mock.Anything, address, big.NewInt(100)).Return(nil, nil).Once()
		walletFactoryMock, _ := newRskWalletFactoryMock()
		rskClient := rootstock.NewRskClient(rpc)
		rskRegistry, err := registry.NewRootstockRegistry(context.Background(), env, rskClient, walletFactoryMock, environment.DefaultTimeouts())
		require.ErrorContains(t, err, "configured start block 100 is not the PegIn address registry deployment block")
		require.Nil(t, rskRegistry)
	})
}
