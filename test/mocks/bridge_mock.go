package mocks

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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func (m *BridgeMock) GetAddress() string {
	args := m.Called()
	return args.String(0)
}

func (m *BridgeMock) GetFedAddress() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *BridgeMock) GetFlyoverDerivationAddress(args blockchain.FlyoverDerivationArgs) (blockchain.FlyoverDerivation, error) {
	a := m.Called(args)
	return a.Get(0).(blockchain.FlyoverDerivation), a.Error(1)
}

func (m *BridgeMock) FetchFederationInfo() (blockchain.FederationInfo, error) {
	args := m.Called()
	return args.Get(0).(blockchain.FederationInfo), args.Error(1)
}

func (m *BridgeMock) GetRequiredTxConfirmations() uint64 {
	return m.Called().Get(0).(uint64)
}

func (m *BridgeMock) RegisterBtcCoinbaseTransaction(registrationParams blockchain.BtcCoinbaseTransactionInformation) (string, error) {
	args := m.Called(registrationParams)
	return args.String(0), args.Error(1)
}
