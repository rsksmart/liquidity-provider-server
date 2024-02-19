package test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/mock"
)

type FeeCollectorMock struct {
	blockchain.FeeCollector
	mock.Mock
}

func (m *FeeCollectorMock) DaoFeePercentage() (uint64, error) {
	args := m.Called()
	return args.Get(0).(uint64), args.Error(1)
}
