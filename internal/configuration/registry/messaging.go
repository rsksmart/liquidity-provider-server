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
}

func NewMessagingRegistry(
	ctx context.Context,
	env environment.Environment,
	rskClient *rootstock.RskClient,
	btcConn *bitcoin.Connection,
) *Messaging {
	return &Messaging{
		Rpc: blockchain.Rpc{
			Btc: bitcoin.NewBitcoindRpc(btcConn),
			Rsk: rootstock.NewRskjRpcServer(rskClient),
		},
		EventBus:    NewEventBus(),
		AlertSender: NewAlertSender(ctx, env),
	}
}
