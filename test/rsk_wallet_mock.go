package test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/mock"
)

type HashMock struct {
	mock.Mock
}

func (m *HashMock) Hash(bytes ...[]byte) []byte {
	args := m.Called(bytes)
	return args.Get(0).([]byte)
}

type RskWalletMock struct {
	blockchain.RootstockWallet
	entities.Signer
	mock.Mock
}

func (m *RskWalletMock) SendRbtc(ctx context.Context, config blockchain.TransactionConfig, toAddress string) (string, error) {
	args := m.Called(ctx, config, toAddress)
	return args.String(0), args.Error(1)
}

func (m *RskWalletMock) SignBytes(msg []byte) ([]byte, error) {
	args := m.Called(msg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *RskWalletMock) Validate(signature, hash string) bool {
	args := m.Called(signature, hash)
	return args.Bool(0)
}
