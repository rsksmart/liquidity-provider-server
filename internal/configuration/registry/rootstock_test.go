package registry_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// nolint:funlen
func TestNewRootstockRegistry(t *testing.T) {
	testEnv := environment.Environment{
		Rsk: environment.RskEnv{
			DiscoveryAddress:            "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA8",
			CollateralManagementAddress: "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA7",
			PeginContractAddress:        "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA6",
			PegoutContractAddress:       "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA5",
			BridgeAddress:               "0x0000000000000000000000000000000001000006",
		},
		Btc: environment.BtcEnv{Network: "testnet"},
	}
	t.Run("should create a new Rootstock registry", func(t *testing.T) {
		env := testEnv
		walletFactoryMock := new(mocks.AbstractFactoryMock)
		rskWalletMock := new(mocks.RskSignerWalletMock)
		rskWalletMock.On("Address").Return(common.HexToAddress(test.AnyRskAddress))
		walletFactoryMock.On("RskWallet").Return(rskWalletMock, nil)
		rskConnBinding := new(mocks.RpcClientBindingMock)
		rskClient := rootstock.NewRskClient(rskConnBinding)
		rskRegistry, err := registry.NewRootstockRegistry(env, rskClient, walletFactoryMock, environment.DefaultTimeouts())
		require.NoError(t, err)
		require.NotNil(t, rskRegistry)
		require.NotEmpty(t, rskRegistry.Contracts)
		require.NotNil(t, rskRegistry.Contracts.Discovery)
		require.NotNil(t, rskRegistry.Contracts.CollateralManagement)
		require.NotNil(t, rskRegistry.Contracts.PegIn)
		require.NotNil(t, rskRegistry.Contracts.PegOut)
		require.NotNil(t, rskRegistry.Contracts.Bridge)
		require.Equal(t, rskWalletMock, rskRegistry.Wallet)
		require.Equal(t, rskClient, rskRegistry.Client)
	})
	t.Run("should return an error when the discovery contract address is invalid", func(t *testing.T) {
		env := testEnv
		env.Rsk.DiscoveryAddress = test.AnyString
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(env, rskClient, new(mocks.AbstractFactoryMock), environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the pegin contract address is invalid", func(t *testing.T) {
		env := testEnv
		env.Rsk.PeginContractAddress = test.AnyString
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(env, rskClient, new(mocks.AbstractFactoryMock), environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the pegout contract address is invalid", func(t *testing.T) {
		env := testEnv
		env.Rsk.PegoutContractAddress = test.AnyString
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(env, rskClient, new(mocks.AbstractFactoryMock), environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the collateral management address is invalid", func(t *testing.T) {
		env := testEnv
		env.Rsk.CollateralManagementAddress = test.AnyString
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(env, rskClient, new(mocks.AbstractFactoryMock), environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the bridge address is invalid", func(t *testing.T) {
		env := testEnv
		env.Rsk.BridgeAddress = test.AnyString
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(env, rskClient, new(mocks.AbstractFactoryMock), environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the wallet factory fails", func(t *testing.T) {
		env := testEnv
		walletFactoryMock := new(mocks.AbstractFactoryMock)
		walletFactoryMock.On("RskWallet").Return(nil, assert.AnError)
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(env, rskClient, walletFactoryMock, environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the BTC network params cannot be retrieved", func(t *testing.T) {
		env := testEnv
		env.Btc.Network = test.AnyString
		walletFactoryMock := new(mocks.AbstractFactoryMock)
		walletFactoryMock.On("RskWallet").Return(new(mocks.RskSignerWalletMock), nil)
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(env, rskClient, walletFactoryMock, environment.DefaultTimeouts())
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
}
