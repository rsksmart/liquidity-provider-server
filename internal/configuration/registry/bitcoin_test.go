package registry_test

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewBitcoinRegistry(t *testing.T) {
	t.Run("should return a new bitcoin registry", func(t *testing.T) {
		walletFactoryMock := new(mocks.AbstractFactoryMock)
		client := new(mocks.ClientAdapterMock)
		connection := bitcoin.NewConnection(&chaincfg.TestNet3Params, client)
		monitoringWallet := new(mocks.BitcoinWalletMock)
		paymentWallet := new(mocks.BitcoinWalletMock)
		walletFactoryMock.On("BitcoinMonitoringWallet", bitcoin.PeginWalletId).Return(monitoringWallet, nil)
		walletFactoryMock.On("BitcoinPaymentWallet", bitcoin.DerivativeWalletId).Return(paymentWallet, nil)
		btcRegistry, err := registry.NewBitcoinRegistry(walletFactoryMock, connection)
		require.NoError(t, err)
		assert.Equal(t, monitoringWallet, btcRegistry.MonitoringWallet)
		assert.Equal(t, paymentWallet, btcRegistry.PaymentWallet)
		assert.Equal(t, connection, btcRegistry.RpcConnection)
		walletFactoryMock.AssertExpectations(t)
	})
	t.Run("should return an error when monitoring wallet creation fails", func(t *testing.T) {
		walletFactoryMock := new(mocks.AbstractFactoryMock)
		client := new(mocks.ClientAdapterMock)
		connection := bitcoin.NewConnection(&chaincfg.TestNet3Params, client)
		walletFactoryMock.On("BitcoinPaymentWallet", bitcoin.DerivativeWalletId).Return(nil, assert.AnError)
		btcRegistry, err := registry.NewBitcoinRegistry(walletFactoryMock, connection)
		require.Error(t, err)
		assert.Nil(t, btcRegistry)
		walletFactoryMock.AssertExpectations(t)
	})
	t.Run("should return an error when payment wallet creation fails", func(t *testing.T) {
		walletFactoryMock := new(mocks.AbstractFactoryMock)
		connection := bitcoin.NewConnection(&chaincfg.TestNet3Params, new(mocks.ClientAdapterMock))
		walletFactoryMock.On("BitcoinPaymentWallet", bitcoin.DerivativeWalletId).Return(new(mocks.BitcoinWalletMock), nil)
		walletFactoryMock.On("BitcoinMonitoringWallet", bitcoin.PeginWalletId).Return(nil, assert.AnError)
		btcRegistry, err := registry.NewBitcoinRegistry(walletFactoryMock, connection)
		require.Error(t, err)
		assert.Nil(t, btcRegistry)
		walletFactoryMock.AssertExpectations(t)
	})
}
