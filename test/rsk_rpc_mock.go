package test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/mock"
)

type RskRpcMock struct {
	mock.Mock
	blockchain.RootstockRpcServer
}

func (m *RskRpcMock) GetTransactionReceipt(ctx context.Context, hash string) (blockchain.TransactionReceipt, error) {
	args := m.Called(ctx, hash)
	return args.Get(0).(blockchain.TransactionReceipt), args.Error(1)
}

func (m *RskRpcMock) GasPrice(ctx context.Context) (*entities.Wei, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func (m *RskRpcMock) GetBalance(ctx context.Context, address string) (*entities.Wei, error) {
	args := m.Called(ctx, address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func (m *RskRpcMock) GetHeight(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}
func (m *RskRpcMock) EstimateGas(ctx context.Context, addr string, value *entities.Wei, data []byte) (*entities.Wei, error) {
	args := m.Called(ctx, addr, value, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Wei), args.Error(1)
}
