package mocks

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/mock"
)

type BtcWalletMock struct {
	mock.Mock
	blockchain.BitcoinWallet
}

func (m *BtcWalletMock) EstimateTxFees(toAddress string, value *entities.Wei) (blockchain.BtcFeeEstimation, error) {
	args := m.Called(toAddress, value)
	return args.Get(0).(blockchain.BtcFeeEstimation), args.Error(1)
}

func (m *BtcWalletMock) GetBalance() (*entities.Wei, error) {
	args := m.Called()
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func (m *BtcWalletMock) SendWithOpReturn(address string, value *entities.Wei, opReturnContent []byte) (blockchain.BitcoinTransactionResult, error) {
	args := m.Called(address, value, opReturnContent)
	return args.Get(0).(blockchain.BitcoinTransactionResult), args.Error(1)
}

func (m *BtcWalletMock) ImportAddress(address string) error {
	args := m.Called(address)
	return args.Error(0)
}

func (m *BtcWalletMock) GetTransactions(address string) ([]blockchain.BitcoinTransactionInformation, error) {
	args := m.Called(address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]blockchain.BitcoinTransactionInformation), args.Error(1)
}
