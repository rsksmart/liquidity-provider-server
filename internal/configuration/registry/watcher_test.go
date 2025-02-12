package registry_test

import (
	"context"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func TestNewWatcherRegistry(t *testing.T) {
	t.Run("Watcher registry constructor should initialize every watcher", func(t *testing.T) {
		env := environment.Environment{
			Rsk: environment.RskEnv{LbcAddress: "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA8", BridgeAddress: "0x0000000000000000000000000000000001000006"},
			Btc: environment.BtcEnv{Network: "testnet"},
		}

		client := &mocks.DbClientBindingMock{}
		client.On("Database", mongo.DbName).Return(&mocks.DbBindingMock{})
		conn := mongo.NewConnection(client, time.Duration(1))
		dbRegistry := registry.NewDatabaseRegistry(conn)

		walletFactoryMock := new(mocks.AbstractFactoryMock)
		walletFactoryMock.On("RskWallet").Return(new(mocks.RskSignerWalletMock), nil)
		rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
		rskRegistry, err := registry.NewRootstockRegistry(env, rskClient, walletFactoryMock, environment.DefaultTimeouts())
		require.NoError(t, err)

		connection := bitcoin.NewConnection(&chaincfg.TestNet3Params, new(mocks.ClientAdapterMock))
		walletFactoryMock.On("BitcoinMonitoringWallet", bitcoin.PeginWalletId).Return(new(mocks.BitcoinWalletMock), nil)
		walletFactoryMock.On("BitcoinPaymentWallet", bitcoin.DerivativeWalletId).Return(new(mocks.BitcoinWalletMock), nil)
		btcRegistry, err := registry.NewBitcoinRegistry(walletFactoryMock, connection)
		require.NoError(t, err)

		messagingRegistry := registry.NewMessagingRegistry(context.Background(), environment.Environment{}, rskClient, connection)
		lp := registry.NewLiquidityProvider(dbRegistry, rskRegistry, btcRegistry, messagingRegistry)
		mutexes := environment.NewApplicationMutexes()
		useCaseRegistry := registry.NewUseCaseRegistry(env, rskRegistry, btcRegistry, dbRegistry, lp, messagingRegistry, mutexes)

		watcherRegistry := registry.NewWatcherRegistry(env, useCaseRegistry, rskRegistry, btcRegistry, lp, messagingRegistry, watcher.NewApplicationTickers(), environment.DefaultTimeouts())

		require.NotNil(t, watcherRegistry)
		value := reflect.ValueOf(watcherRegistry).Elem()
		for i := 0; i < value.NumField(); i++ {
			if value.Field(i).IsNil() {
				t.Errorf("Field %s of watcher registry is nil", value.Type().Field(i).Name)
			}
		}
	})
}
