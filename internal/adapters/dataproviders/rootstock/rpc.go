package rootstock

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

// newAccountGasCost fixed gas amount to add to the estimation if the destination address is a new account
const newAccountGasCost = 25000

type rskjRpcServer struct {
	client      RpcClientBinding
	retryParams RetryParams
}

func NewRskjRpcServer(client *RskClient, retryParams RetryParams) blockchain.RootstockRpcServer {
	return &rskjRpcServer{client: client.client, retryParams: retryParams}
}

func (rpc *rskjRpcServer) GetBalance(ctx context.Context, address string) (*entities.Wei, error) {
	var destination common.Address
	var err error

	if err = ParseAddress(&destination, address); err != nil {
		return nil, err
	}

	result, err := rskRetry(rpc.retryParams.Retries, rpc.retryParams.Sleep,
		func() (*big.Int, error) {
			return rpc.client.BalanceAt(ctx, destination, nil)
		})

	if err != nil {
		return nil, err
	} else {
		return entities.NewBigWei(result), nil
	}
}

func (rpc *rskjRpcServer) EstimateGas(ctx context.Context, address string, value *entities.Wei, data []byte) (*entities.Wei, error) {
	var destination common.Address
	var additionalGas uint64
	var newAccount bool
	var err error

	if err = ParseAddress(&destination, address); err != nil {
		return nil, err
	}

	if newAccount, err = rpc.isNewAccount(ctx, destination); err != nil {
		return nil, err
	} else if newAccount {
		additionalGas = newAccountGasCost
	}

	tx := ethereum.CallMsg{
		To:    &destination,
		Data:  data,
		Value: value.AsBigInt(),
	}
	result, err := rskRetry(rpc.retryParams.Retries, rpc.retryParams.Sleep,
		func() (uint64, error) {
			return rpc.client.EstimateGas(ctx, tx)
		})
	if err != nil {
		return nil, err
	} else {
		return new(entities.Wei).Add(entities.NewUWei(result), entities.NewUWei(additionalGas)), nil
	}
}

func (rpc *rskjRpcServer) GasPrice(ctx context.Context) (*entities.Wei, error) {
	result, err := rskRetry(rpc.retryParams.Retries, rpc.retryParams.Sleep,
		func() (*big.Int, error) {
			return rpc.client.SuggestGasPrice(ctx)
		})
	if err != nil {
		return nil, err
	} else {
		return entities.NewBigWei(result), err
	}
}

func (rpc *rskjRpcServer) GetHeight(ctx context.Context) (uint64, error) {
	return rskRetry(rpc.retryParams.Retries, rpc.retryParams.Sleep,
		func() (uint64, error) {
			return rpc.client.BlockNumber(ctx)
		})
}

func (rpc *rskjRpcServer) GetTransactionReceipt(ctx context.Context, hash string) (blockchain.TransactionReceipt, error) {
	_, err := hex.DecodeString(strings.TrimPrefix(hash, "0x"))
	if err != nil {
		return blockchain.TransactionReceipt{}, errors.New("invalid transaction hash")
	}

	receipt, err := rskRetry(rpc.retryParams.Retries, rpc.retryParams.Sleep,
		func() (*types.Receipt, error) {
			return rpc.client.TransactionReceipt(ctx, common.HexToHash(hash))
		})
	if err != nil {
		return blockchain.TransactionReceipt{}, err
	}

	tx, err := rskRetry(rpc.retryParams.Retries, rpc.retryParams.Sleep,
		func() (*types.Transaction, error) {
			transaction, _, rpcError := rpc.client.TransactionByHash(ctx, common.HexToHash(hash))
			return transaction, rpcError
		})
	if err != nil {
		return blockchain.TransactionReceipt{}, err
	}
	return ParseReceipt(tx, receipt)
}

func (rpc *rskjRpcServer) isNewAccount(ctx context.Context, address common.Address) (bool, error) {
	var (
		err     error
		code    []byte
		balance *big.Int
		nonce   uint64
	)
	code, err = rskRetry(rpc.retryParams.Retries, rpc.retryParams.Sleep,
		func() ([]byte, error) {
			return rpc.client.CodeAt(ctx, address, nil)
		})
	if err != nil {
		return false, err
	}

	balance, err = rskRetry(rpc.retryParams.Retries, rpc.retryParams.Sleep,
		func() (*big.Int, error) {
			return rpc.client.BalanceAt(ctx, address, nil)
		})
	if err != nil {
		return false, err
	}

	nonce, err = rskRetry(rpc.retryParams.Retries, rpc.retryParams.Sleep,
		func() (uint64, error) {
			return rpc.client.NonceAt(ctx, address, nil)
		})
	if err != nil {
		return false, err
	}

	return len(code) == 0 && balance.Cmp(common.Big0) == 0 && nonce == 0, nil
}

func (rpc *rskjRpcServer) GetBlockByHash(ctx context.Context, hash string) (blockchain.BlockInfo, error) {
	if _, err := hex.DecodeString(strings.TrimPrefix(hash, "0x")); err != nil {
		return blockchain.BlockInfo{}, errors.New("invalid block hash")
	}
	result, err := rskRetry(rpc.retryParams.Retries, rpc.retryParams.Sleep,
		func() (*types.Block, error) {
			return rpc.client.BlockByHash(ctx, common.HexToHash(hash))
		})
	if err != nil {
		return blockchain.BlockInfo{}, err
	}

	return blockchain.BlockInfo{
		Hash:      result.Hash().String(),
		Number:    result.NumberU64(),
		Timestamp: time.Unix(int64(result.Time()), 0),
		Nonce:     result.Nonce(),
	}, nil
}

func (rpc *rskjRpcServer) GetBlockByNumber(ctx context.Context, blockNumber *big.Int) (blockchain.BlockInfo, error) {
	result, err := rskRetry(rpc.retryParams.Retries, rpc.retryParams.Sleep,
		func() (*types.Block, error) {
			return rpc.client.BlockByNumber(ctx, blockNumber)
		})
	if err != nil {
		return blockchain.BlockInfo{}, err
	}
	return blockchain.BlockInfo{
		Hash:      result.Hash().String(),
		Number:    result.NumberU64(),
		Timestamp: time.Unix(int64(result.Time()), 0),
		Nonce:     result.Nonce(),
	}, nil
}

func (rpc *rskjRpcServer) ChainId(ctx context.Context) (uint64, error) {
	result, err := rskRetry(rpc.retryParams.Retries, rpc.retryParams.Sleep,
		func() (uint64, error) {
			if chainId, rpcErr := rpc.client.ChainID(ctx); rpcErr != nil {
				return 0, rpcErr
			} else {
				return chainId.Uint64(), nil
			}
		})
	if err != nil {
		return 0, err
	}
	return result, nil
}
