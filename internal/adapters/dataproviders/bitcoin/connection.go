package bitcoin

import (
	"context"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin/btcclient"
	log "github.com/sirupsen/logrus"
)

type Connection struct {
	NetworkParams *chaincfg.Params
	client        btcclient.ClientAdapter
	WalletId      string
}

func NewWalletConnection(networkParams *chaincfg.Params, client btcclient.ClientAdapter, walletId string) *Connection {
	return &Connection{NetworkParams: networkParams, client: client, WalletId: walletId}
}

func NewConnection(networkParams *chaincfg.Params, client btcclient.ClientAdapter) *Connection {
	return &Connection{NetworkParams: networkParams, client: client}
}

func (c *Connection) Shutdown(endChannel chan<- bool) {
	c.client.Disconnect()
	endChannel <- true
	log.Debug("Disconnected from BTC node")
}

func (c *Connection) CheckConnection(ctx context.Context) bool {
	err := c.client.Ping()
	if err != nil {
		log.Error("Error checking BTC node connection: ", err)
	}
	return err == nil
}
