package mocks

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/stretchr/testify/mock"
	"math/big"
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

func (m *BtcRpcMock) DecodeAddress(address string) ([]byte, error) {
	args := m.Called(address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *BtcRpcMock) GetPartialMerkleTree(hash string) ([]byte, error) {
	args := m.Called(hash)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *BtcRpcMock) ValidateAddress(address string) error {
	args := m.Called(address)
	return args.Error(0)
}

func (m *BtcRpcMock) GetHeight() (*big.Int, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *BtcRpcMock) GetCoinbaseInformation(txHash string) (rootstock.BtcCoinbaseTransactionInformation, error) {
	args := m.Called(txHash)
	return args.Get(0).(rootstock.BtcCoinbaseTransactionInformation), args.Error(1)
}

func (m *BtcRpcMock) NetworkName() string {
	args := m.Called()
	return args.String(0)
}

func (m *BtcRpcMock) GetBlockchainInfo() (blockchain.BitcoinBlockchainInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return blockchain.BitcoinBlockchainInfo{}, args.Error(1)
	}
	return (args.Get(0)).(blockchain.BitcoinBlockchainInfo), args.Error(1)
}

func (m *BtcRpcMock) GetZeroAddress(addressType blockchain.BtcAddressType) (string, error) {
	args := m.Called(addressType)
	return args.String(0), args.Error(1)
}
