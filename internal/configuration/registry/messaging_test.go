package registry_test

import (
	"context"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMessagingRegistry(t *testing.T) {
	client := new(mocks.ClientAdapterMock)
	btcConnection := bitcoin.NewConnection(&chaincfg.TestNet3Params, client)
	rskConnBinging := new(mocks.RpcClientBindingMock)
	rskClient := rootstock.NewRskClient(rskConnBinging)
	messagingRegistry := registry.NewMessagingRegistry(context.Background(), environment.Environment{}, rskClient, btcConnection, registry.ExternalRpc{})
	assert.NotNil(t, messagingRegistry)
	assert.NotEmpty(t, messagingRegistry.Rpc)
	assert.NotNil(t, messagingRegistry.Rpc.Rsk)
	assert.NotNil(t, messagingRegistry.Rpc.Btc)
	assert.NotNil(t, messagingRegistry.EventBus)
	assert.NotNil(t, messagingRegistry.AlertSender)
}
