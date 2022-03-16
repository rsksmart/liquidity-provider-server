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

func (R *RSKClientMock) ChainID(ctx context.Context) (*big.Int, error) {
	args := R.Called(ctx)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (R *RSKClientMock) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	args := R.Called(ctx, msg)
	return args.Get(0).(uint64), args.Error(1)
}

func (R *RSKClientMock) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	args := R.Called(ctx)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (R *RSKClientMock) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	args := R.Called(ctx, account, blockNumber)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (R *RSKClientMock) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	args := R.Called(ctx, account, blockNumber)
	return args.Get(0).(uint64), args.Error(1)
}

func (R *RSKClientMock) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	args := R.Called(ctx, account, blockNumber)
	return args.Get(0).([]byte), args.Error(1)
}

func (R *RSKClientMock) TransactionReceipt(ctx context.Context, txHash common.Hash) (*gethTypes.Receipt, error) {
	args := R.Called(ctx, txHash)
	return args.Get(0).(*gethTypes.Receipt), args.Error(1)
}

func (R *RSKClientMock) Close() {
	R.Called()
}
