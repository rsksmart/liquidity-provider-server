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
	rpcCallRetryMax   = 3
	rpcCallRetrySleep = 1 * time.Minute
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

func awaitTx(client RpcClientBinding, miningTimeout time.Duration, logName string, txCall func() (*geth.Transaction, error)) (r *geth.Receipt, e error) {
	return AwaitTxWithCtx(client, miningTimeout, logName, context.Background(), txCall)
}

func AwaitTxWithCtx(client RpcClientBinding, miningTimeout time.Duration, logName string, ctx context.Context, txCall func() (*geth.Transaction, error)) (*geth.Receipt, error) {
	var tx *geth.Transaction
	var err error

	log.Infof("Executing %s transaction...", logName)
	deadline, ok := ctx.Deadline()
	if ok {
		log.Debugf("Waiting for transaction to be mined until %v...", deadline)
	}
	tx, err = txCall()
	if err != nil {
		return nil, err
	} else if tx == nil {
		return nil, errors.New("invalid transaction")
	}

	ctx, cancel := context.WithTimeout(ctx, miningTimeout)
	defer cancel()

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil || receipt == nil {
		log.Infof("Error waiting for transaction %s (%s) to be mined: %v", logName, tx.Hash().String(), err)
		return nil, err
	}

	if receipt.Status == 1 {
		log.Infof("Transaction %s (%s) executed successfully", logName, tx.Hash().String())
	} else {
		log.Infof("Transaction %s (%s) reverted", logName, tx.Hash().String())
	}
	return receipt, nil
}
