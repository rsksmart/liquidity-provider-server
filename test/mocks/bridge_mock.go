package mocks

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/stretchr/testify/mock"
)

type BridgeMock struct {
	rootstock.Bridge
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

func (m *BridgeMock) GetFlyoverDerivationAddress(args rootstock.FlyoverDerivationArgs) (rootstock.FlyoverDerivation, error) {
	a := m.Called(args)
	return a.Get(0).(rootstock.FlyoverDerivation), a.Error(1)
}

func (m *BridgeMock) FetchFederationInfo() (rootstock.FederationInfo, error) {
	args := m.Called()
	return args.Get(0).(rootstock.FederationInfo), args.Error(1)
}

func (m *BridgeMock) GetRequiredTxConfirmations() uint64 {
	return m.Called().Get(0).(uint64)
}

func (m *BridgeMock) RegisterBtcCoinbaseTransaction(registrationParams rootstock.BtcCoinbaseTransactionInformation) (string, error) {
	args := m.Called(registrationParams)
	return args.String(0), args.Error(1)
}

func (m *BridgeMock) GetBatchPegOutCreatedEvent(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]rootstock.BatchPegOut, error) {
	args := m.Called(ctx, fromBlock, toBlock)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]rootstock.BatchPegOut), args.Error(1)
}
