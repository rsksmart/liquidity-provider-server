package bitcoin_test

import (
	"context"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConnection(t *testing.T) {
	networkParams := &chaincfg.Params{}
	client := &mocks.ClientAdapterMock{}

	regularConnection := bitcoin.NewConnection(networkParams, client)
	assert.NotNil(t, regularConnection)
	assert.Equal(t, networkParams, regularConnection.NetworkParams)
	assert.Empty(t, regularConnection.WalletId)

	walletConnection := bitcoin.NewWalletConnection(networkParams, client, test.AnyString)
	assert.NotNil(t, walletConnection)
	assert.Equal(t, networkParams, walletConnection.NetworkParams)
	assert.Equal(t, test.AnyString, walletConnection.WalletId)
}

func TestConnection_CheckConnection(t *testing.T) {
	networkParams := &chaincfg.Params{}
	client := &mocks.ClientAdapterMock{}
	client.On("Ping").Return(assert.AnError).Once()
	client.On("Ping").Return(nil).Once()
	conn := bitcoin.NewConnection(networkParams, client)
	conn.CheckConnection(context.Background())
	conn.CheckConnection(context.Background())
	client.AssertExpectations(t)
}

func TestConnection_Shutdown(t *testing.T) {
	endChannel := make(chan bool)
	client := &mocks.ClientAdapterMock{}
	client.On("Disconnect").Once()
	conn := bitcoin.NewConnection(&chaincfg.Params{}, client)
	go conn.Shutdown(endChannel)
	<-endChannel
	client.AssertExpectations(t)
}
