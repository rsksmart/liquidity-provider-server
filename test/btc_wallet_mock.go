package test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/mock"
)

type BtcWalletMock struct {
	mock.Mock
	blockchain.BitcoinWallet
}

func (m *BtcWalletMock) EstimateTxFees(toAddress string, value *entities.Wei) (*entities.Wei, error) {
	args := m.Called(toAddress, value)
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func (m *BtcWalletMock) GetBalance() (*entities.Wei, error) {
	args := m.Called()
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func (m *BtcWalletMock) SendWithOpReturn(address string, value *entities.Wei, opReturnContent []byte) (string, error) {
	args := m.Called(address, value, opReturnContent)
	return args.String(0), args.Error(1)
}
