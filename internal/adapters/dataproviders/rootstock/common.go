package rootstock

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	rpcCallRetryMax     = 3
	rpcCallRetrySleep   = 1 * time.Minute
	txMiningWaitTimeout = 2 * time.Minute
)

var DefaultRetryParams = RetryParams{
	Retries: rpcCallRetryMax,
	Sleep:   rpcCallRetrySleep,
}

type RskClient struct {
	client RpcClientBinding
}

type RetryParams struct {
	Retries uint
	Sleep   time.Duration
}

func NewRskClient(client RpcClientBinding) *RskClient {
	return &RskClient{client: client}
}

func (c *RskClient) Rpc() RpcClientBinding {
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
	entities.Signer
	Address() common.Address
	Sign(common.Address, *types.Transaction) (*types.Transaction, error)
}

type RskSignerWallet interface {
	blockchain.RootstockWallet
	TransactionSigner
}

func ParseAddress(address *common.Address, textAddress string) error {
	if !common.IsHexAddress(textAddress) {
		return blockchain.InvalidAddressError
	}
	*address = common.HexToAddress(textAddress)
	return nil
}

func rskRetry[R any](retries uint, retrySleep time.Duration, call func() (R, error)) (R, error) {
	var result R
	var err error
	var i uint

	if retries == 0 {
		return call()
	}

	for i = 0; i < retries; i++ {
		result, err = call()
		if err == nil {
			return result, nil
		}
		time.Sleep(retrySleep)
	}
	return result, err
}

func awaitTx(client RpcClientBinding, logName string, txCall func() (*geth.Transaction, error)) (r *geth.Receipt, e error) {
	return awaitTxWithCtx(client, logName, context.Background(), txCall)
}

func awaitTxWithCtx(client RpcClientBinding, logName string, ctx context.Context, txCall func() (*geth.Transaction, error)) (r *geth.Receipt, e error) {
	var tx *geth.Transaction
	var err error

	log.Infof("Executing %s transaction...", logName)
	tx, err = txCall()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, txMiningWaitTimeout)
	defer func() {
		cancel()
		if r != nil && r.Status == 1 {
			log.Infof("Transaction %s (%s) executed successfully\n", logName, tx.Hash().String())
		} else if tx != nil {
			log.Infof("Transaction %s (%s) failed\n", logName, tx.Hash().String())
		} else {
			log.Info("Transaction failed")
		}
	}()
	if tx != nil {
		return bind.WaitMined(ctx, client, tx)
	}
	return nil, errors.New("invalid transaction")
}
