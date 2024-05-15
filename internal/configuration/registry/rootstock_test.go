package registry_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewRootstockRegistry(t *testing.T) {
	testEnv := environment.Environment{
		Rsk: environment.RskEnv{LbcAddress: "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA8", BridgeAddress: "0x0000000000000000000000000000000001000006"},
		Btc: environment.BtcEnv{Network: "testnet"},
	}
	t.Run("should create a new Rootstock registry", func(t *testing.T) {
		env := testEnv
		walletFactoryMock := new(mocks.AbstractFactoryMock)
		rskWalletMock := new(mocks.RskSignerWalletMock)
		walletFactoryMock.On("RskWallet").Return(rskWalletMock, nil)
		rskConnBinding := new(mocks.RpcClientBindingMock)
		rskClient := rootstock.NewRskClient(rskConnBinding)
		rskRegistry, err := registry.NewRootstockRegistry(env, rskClient, walletFactoryMock)
		require.NoError(t, err)
		require.NotNil(t, rskRegistry)
		require.NotEmpty(t, rskRegistry.Contracts)
		require.NotNil(t, rskRegistry.Contracts.Lbc)
		require.NotNil(t, rskRegistry.Contracts.Bridge)
		require.NotNil(t, rskRegistry.Contracts.FeeCollector)
		require.Equal(t, rskWalletMock, rskRegistry.Wallet)
		require.Equal(t, rskClient, rskRegistry.Client)
	})
	t.Run("should return an error when the LBC address is invalid", func(t *testing.T) {
		env := testEnv
		env.Rsk.LbcAddress = test.AnyString
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(env, rskClient, new(mocks.AbstractFactoryMock))
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the bridge address is invalid", func(t *testing.T) {
		env := testEnv
		env.Rsk.BridgeAddress = test.AnyString
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(env, rskClient, new(mocks.AbstractFactoryMock))
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the wallet factory fails", func(t *testing.T) {
		env := testEnv
		walletFactoryMock := new(mocks.AbstractFactoryMock)
		walletFactoryMock.On("RskWallet").Return(nil, assert.AnError)
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(env, rskClient, walletFactoryMock)
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
	t.Run("should return an error when the BTC network params cannot be retrieved", func(t *testing.T) {
		env := testEnv
		env.Btc.Network = test.AnyString
		walletFactoryMock := new(mocks.AbstractFactoryMock)
		walletFactoryMock.On("RskWallet").Return(new(mocks.RskSignerWalletMock), nil)
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(env, rskClient, walletFactoryMock)
		require.Error(t, err)
		require.Nil(t, rskRegistry)
	})
}
