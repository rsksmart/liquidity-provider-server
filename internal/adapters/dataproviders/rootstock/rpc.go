package rootstock

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"math/big"
	"strings"
)

const newAccountGasCost = 25000

type rskjRpcServer struct {
	client *ethclient.Client
}

func NewRskjRpcServer(client *RskClient) blockchain.RootstockRpcServer {
	return &rskjRpcServer{client: client.client}
}

func (rpc *rskjRpcServer) GetBalance(ctx context.Context, address string) (*entities.Wei, error) {
	var destination common.Address
	var err error

	if err = ParseAddress(&destination, address); err != nil {
		return nil, err
	}

	result, err := rskRetry(func() (*big.Int, error) {
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
	result, err := rskRetry(func() (uint64, error) {
		return rpc.client.EstimateGas(ctx, tx)
	})
	if err != nil {
		return nil, err
	} else {
		return entities.NewUWei(result + additionalGas), nil
	}
}

func (rpc *rskjRpcServer) GasPrice(ctx context.Context) (*entities.Wei, error) {
	result, err := rskRetry(func() (*big.Int, error) {
		return rpc.client.SuggestGasPrice(ctx)
	})
	if err != nil {
		return nil, err
	} else {
		return entities.NewBigWei(result), err
	}
}

func (rpc *rskjRpcServer) GetHeight(ctx context.Context) (uint64, error) {
	return rskRetry(func() (uint64, error) {
		return rpc.client.BlockNumber(ctx)
	})
}

func (rpc *rskjRpcServer) GetTransactionReceipt(ctx context.Context, hash string) (blockchain.TransactionReceipt, error) {
	var from common.Address

	_, err := hex.DecodeString(strings.TrimPrefix(hash, "0x"))
	if err != nil {
		return blockchain.TransactionReceipt{}, errors.New("invalid transaction hash")
	}

	receipt, err := rskRetry(func() (*types.Receipt, error) {
		return rpc.client.TransactionReceipt(ctx, common.HexToHash(hash))
	})
	if err != nil {
		return blockchain.TransactionReceipt{}, err
	}

	tx, err := rskRetry(func() (*types.Transaction, error) {
		transaction, _, rpcError := rpc.client.TransactionByHash(ctx, common.HexToHash(hash))
		return transaction, rpcError
	})
	if err != nil {
		return blockchain.TransactionReceipt{}, err
	}

	gasUsed := new(big.Int)
	gasUsed.SetUint64(receipt.GasUsed)
	cumulativeGasUsed := new(big.Int)
	cumulativeGasUsed.SetUint64(receipt.CumulativeGasUsed)
	from, err = types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
	if err != nil {
		from, _ = types.Sender(types.HomesteadSigner{}, tx)
	}
	return blockchain.TransactionReceipt{
		TransactionHash:   receipt.TxHash.String(),
		BlockHash:         receipt.BlockHash.String(),
		BlockNumber:       receipt.BlockNumber.Uint64(),
		From:              from.String(),
		To:                tx.To().String(),
		CumulativeGasUsed: cumulativeGasUsed,
		GasUsed:           gasUsed,
		Value:             entities.NewBigWei(tx.Value()),
	}, nil
}

func (rpc *rskjRpcServer) isNewAccount(ctx context.Context, address common.Address) (bool, error) {
	var (
		err     error
		code    []byte
		balance *big.Int
		nonce   uint64
	)
	code, err = rskRetry(func() ([]byte, error) {
		return rpc.client.CodeAt(ctx, address, nil)
	})
	if err != nil {
		return false, err
	}

	balance, err = rskRetry(func() (*big.Int, error) {
		return rpc.client.BalanceAt(ctx, address, nil)
	})
	if err != nil {
		return false, err
	}

	nonce, err = rskRetry(func() (uint64, error) {
		return rpc.client.NonceAt(ctx, address, nil)
	})
	if err != nil {
		return false, err
	}

	return len(code) == 0 && balance.Cmp(common.Big0) == 0 && nonce == 0, nil
}
