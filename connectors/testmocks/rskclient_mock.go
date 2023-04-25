package testmocks

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"math/big"
)

type RSKClientMock struct {
	mock.Mock
}

func (mock *RSKClientMock) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	args := mock.Called(ctx, account)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *RSKClientMock) SendTransaction(ctx context.Context, tx *gethTypes.Transaction) error {
	args := mock.Called(ctx, tx)
	return args.Error(0)
}

func (mock *RSKClientMock) ChainID(ctx context.Context) (*big.Int, error) {
	arg := mock.Called(ctx)
	return arg.Get(0).(*big.Int), arg.Error(1)
}

func (mock *RSKClientMock) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	args := mock.Called(ctx, msg)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *RSKClientMock) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	args := mock.Called(ctx)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (mock *RSKClientMock) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	args := mock.Called(ctx, account, blockNumber)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (mock *RSKClientMock) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	args := mock.Called(ctx, account, blockNumber)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *RSKClientMock) TransactionReceipt(ctx context.Context, txHash common.Hash) (*gethTypes.Receipt, error) {
	args := mock.Called(ctx, txHash)
	return args.Get(0).(*gethTypes.Receipt), args.Error(1)
}

func (mock *RSKClientMock) BlockNumber(ctx context.Context) (uint64, error) {
	args := mock.Called(ctx)
	return uint64(args.Int(0)), args.Error(1)
}

func (mock *RSKClientMock) Close() {
	//dummy implementation for mock
}

func (mock *RSKClientMock) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	args := mock.Called(ctx, account, blockNumber)
	return args.Get(0).([]byte), args.Error(1)
}