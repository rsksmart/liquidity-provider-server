package rootstock

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	rpcCallRetryMax     = 3
	rpcCallRetrySleep   = 2 * time.Minute
	txMiningWaitTimeout = 3 * time.Minute
)

type RskAccount struct {
	Account  *accounts.Account
	Keystore *keystore.KeyStore
}

type RskClient struct {
	client *ethclient.Client
}

func NewRskClient(client *ethclient.Client) *RskClient {
	return &RskClient{client: client}
}

func (c *RskClient) Rpc() *ethclient.Client {
	return c.client
}

func (c *RskClient) Shutdown(endChannel chan<- bool) {
	c.client.Close()
	endChannel <- true
	log.Debug("Disconnected from RSK node")
}

func (c *RskClient) CheckConnection(ctx context.Context) bool {
	_, err := c.client.ChainID(ctx)
	if err != nil {
		log.Error("Error checking RSK node connection: ", err)
	}
	return err == nil
}

type TransactionSigner interface {
	Address() common.Address
	Sign(common.Address, *types.Transaction) (*types.Transaction, error)
	SignBytes(msg []byte) ([]byte, error)
}

func ParseAddress(address *common.Address, textAddress string) error {
	if !common.IsHexAddress(textAddress) {
		return blockchain.InvalidAddressError
	}
	*address = common.HexToAddress(textAddress)
	return nil
}

func rskRetry[R any](call func() (R, error)) (R, error) {
	var result R
	var err error
	for i := 0; i < rpcCallRetryMax; i++ {
		result, err = call()
		if err == nil {
			return result, nil
		}
		time.Sleep(rpcCallRetrySleep)
	}
	return result, err
}

func awaitTx(client *ethclient.Client, logName string, txCall func() (*geth.Transaction, error)) (r *geth.Receipt, e error) {
	var tx *geth.Transaction
	var err error

	log.Infof("Executing %s transaction...\n", logName)
	tx, err = txCall()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), txMiningWaitTimeout)
	defer func() {
		cancel()
		if r.Status == 1 {
			log.Infof("Transaction %s (%s) executed successfully\n", logName, tx.Hash().String())
		} else {
			log.Infof("Transaction %s (%s) failed\n", logName, tx.Hash().String())
		}
	}()
	return bind.WaitMined(ctx, client, tx)
}
