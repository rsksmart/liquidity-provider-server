package registry_test

import (
	"context"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewLiquidityProvider(t *testing.T) {
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

	messagingRegistry := registry.NewMessagingRegistry(context.Background(), environment.Environment{}, rskClient, connection, registry.ExternalClients{})

	lp := registry.NewLiquidityProvider(dbRegistry, rskRegistry, btcRegistry, messagingRegistry)
	require.NotNil(t, lp)
	assert.IsType(t, &dataproviders.LocalLiquidityProvider{}, lp)
	walletFactoryMock.AssertExpectations(t)
	client.AssertExpectations(t)
}
