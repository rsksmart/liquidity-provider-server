package registry

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

type Messaging struct {
	Rpc         blockchain.Rpc
	EventBus    entities.EventBus
	AlertSender entities.AlertSender
	RskExtraRpc []blockchain.RootstockRpcServer
	BtcExtraRpc []blockchain.BitcoinNetwork
}

type ExternalClients struct {
	RskExternalClients []*rootstock.RskClient
	BtcExternalClients []*bitcoin.Connection
}

func NewMessagingRegistry(
	ctx context.Context,
	env environment.Environment,
	rskClient *rootstock.RskClient,
	btcConn *bitcoin.Connection,
	externalClients ExternalClients,
) *Messaging {
	rskExtraRpcs := make([]blockchain.RootstockRpcServer, len(externalClients.RskExternalClients))
	for i, client := range externalClients.RskExternalClients {
		rskExtraRpcs[i] = rootstock.NewRskjRpcServer(client, rootstock.DefaultRetryParams)
	}
	btcExtraRpcs := make([]blockchain.BitcoinNetwork, len(externalClients.BtcExternalClients))
	for i, client := range externalClients.BtcExternalClients {
		btcExtraRpcs[i] = bitcoin.NewBitcoindRpc(client)
	}
	return &Messaging{
		Rpc: blockchain.Rpc{
			Btc: bitcoin.NewBitcoindRpc(btcConn),
			Rsk: rootstock.NewRskjRpcServer(rskClient, rootstock.DefaultRetryParams),
		},
		EventBus:    NewEventBus(),
		AlertSender: NewAlertSender(ctx, env),
		RskExtraRpc: rskExtraRpcs,
		BtcExtraRpc: btcExtraRpcs,
	}
}
