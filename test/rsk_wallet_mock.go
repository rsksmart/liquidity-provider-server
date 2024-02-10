package test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/mock"
)

type RskWalletMock struct {
	blockchain.RootstockWallet
	mock.Mock
}

func (m *RskWalletMock) SendRbtc(ctx context.Context, config blockchain.TransactionConfig, toAddress string) (string, error) {
	args := m.Called(ctx, config, toAddress)
	return args.String(0), args.Error(1)
}
