package bitcoin

import (
	"context"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin/btcclient"
	log "github.com/sirupsen/logrus"
)

// Connection is a struct that holds the connection to a Bitcoin node
type Connection struct {
	NetworkParams *chaincfg.Params
	client        btcclient.ClientAdapter
	WalletId      string
}

// NewWalletConnection creates a new Connection with a walletId. This connection will use the walletId
// to make the requests to an endpoint referring to that specific wallet on the RPC server.
// E.g. http://127.0.0.1:8332/wallet/<walletId>. Using this constructor has the same effect as using the
// -rpcwallet option in bitcoin-cli. Prefer this constructor if the Connection will be used by a
// blockchain.BitcoinWallet implementation
func NewWalletConnection(networkParams *chaincfg.Params, client btcclient.ClientAdapter, walletId string) *Connection {
	return &Connection{NetworkParams: networkParams, client: client, WalletId: walletId}
}

// NewConnection creates a new Connection with no walletId. This connection will make requests to the default
// endpoint of the RPC server. Prefer this constructor if the Connection will be used by a structure that only
// needs to use non-wallet related RPC methods
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
