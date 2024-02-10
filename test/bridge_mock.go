package test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/mock"
)

type BridgeMock struct {
	blockchain.RootstockBridge
	mock.Mock
}

func (m *BridgeMock) GetMinimumLockTxValue() (*entities.Wei, error) {
	args := m.Called()
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func (m *BridgeMock) GetAddress() string {
	args := m.Called()
	return args.String(0)
}
