package registry_test

import (
	"context"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewLiquidityProvider(t *testing.T) {
	env := environment.Environment{
		Rsk: environment.RskEnv{
			DiscoveryAddress:            "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA8",
			CollateralManagementAddress: "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA7",
			PeginContractAddress:        "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA6",
			PegoutContractAddress:       "0x8901a2Bbf639bFD21A97004BA4D7aE2BD00B8DA5",
			BridgeAddress:               "0x0000000000000000000000000000000001000006",
		},
		Btc: environment.BtcEnv{Network: "testnet"},
	}

	client := &mocks.DbClientBindingMock{}
	client.On("Database", mongo.DbName).Return(&mocks.DbBindingMock{})
	conn := mongo.NewConnection(client, time.Duration(1))
	dbRegistry := registry.NewDatabaseRegistry(conn)

	walletFactoryMock := new(mocks.AbstractFactoryMock)
	walletMock := new(mocks.RskSignerWalletMock)
	walletMock.EXPECT().Address().Return(common.HexToAddress(test.AnyRskAddress))
	walletFactoryMock.On("RskWallet").Return(walletMock, nil)
	rskClient := rootstock.NewRskClient(new(mocks.RpcClientBindingMock))
	rskRegistry, err := registry.NewRootstockRegistry(env, rskClient, walletFactoryMock, environment.DefaultTimeouts())
	require.NoError(t, err)

	connection := bitcoin.NewConnection(&chaincfg.TestNet3Params, new(mocks.ClientAdapterMock))
	walletFactoryMock.On("BitcoinMonitoringWallet", bitcoin.PeginWalletId).Return(new(mocks.BitcoinWalletMock), nil)
	walletFactoryMock.On("BitcoinPaymentWallet", bitcoin.DerivativeWalletId).Return(new(mocks.BitcoinWalletMock), nil)
	btcRegistry, err := registry.NewBitcoinRegistry(walletFactoryMock, connection)
	require.NoError(t, err)

	messagingRegistry := registry.NewMessagingRegistry(context.Background(), environment.Environment{}, rskClient, connection, registry.ExternalRpc{})

	lp := registry.NewLiquidityProvider(dbRegistry, rskRegistry, btcRegistry, messagingRegistry)
	require.NotNil(t, lp)
	assert.IsType(t, &dataproviders.LocalLiquidityProvider{}, lp)
	walletFactoryMock.AssertExpectations(t)
	client.AssertExpectations(t)
}
