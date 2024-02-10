package test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/mock"
)

type BtcRpcMock struct {
	blockchain.BitcoinNetwork
	mock.Mock
}

func (m *BtcRpcMock) GetTransactionInfo(hash string) (blockchain.BitcoinTransactionInformation, error) {
	args := m.Called(hash)
	return args.Get(0).(blockchain.BitcoinTransactionInformation), args.Error(1)
}

func (m *BtcRpcMock) GetRawTransaction(hash string) ([]byte, error) {
	args := m.Called(hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *BtcRpcMock) BuildMerkleBranch(txHash string) (blockchain.MerkleBranch, error) {
	args := m.Called(txHash)
	return args.Get(0).(blockchain.MerkleBranch), args.Error(1)
}

func (m *BtcRpcMock) GetTransactionBlockInfo(txHash string) (blockchain.BitcoinBlockInformation, error) {
	args := m.Called(txHash)
	return args.Get(0).(blockchain.BitcoinBlockInformation), args.Error(1)
}
