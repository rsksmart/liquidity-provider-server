package mocks

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/stretchr/testify/mock"
)

type LbcMock struct {
	blockchain.LiquidityBridgeContract
	mock.Mock
}

func (m *LbcMock) GetAddress() string {
	args := m.Called()
	return args.String(0)
}

func (m *LbcMock) GetMinimumCollateral() (*entities.Wei, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func (m *LbcMock) GetCollateral(providerAddress string) (*entities.Wei, error) {
	args := m.Called(providerAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func (m *LbcMock) GetPegoutCollateral(providerAddress string) (*entities.Wei, error) {
	args := m.Called(providerAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func (m *LbcMock) IsOperationalPegin(providerAddress string) (bool, error) {
	args := m.Called(providerAddress)
	return args.Bool(0), args.Error(1)
}

func (m *LbcMock) IsOperationalPegout(providerAddress string) (bool, error) {
	args := m.Called(providerAddress)
	return args.Bool(0), args.Error(1)
}

func (m *LbcMock) AddCollateral(collateral *entities.Wei) error {
	args := m.Called(collateral)
	return args.Error(0)
}

func (m *LbcMock) AddPegoutCollateral(collateral *entities.Wei) error {
	args := m.Called(collateral)
	return args.Error(0)
}

func (m *LbcMock) ProviderResign() error {
	args := m.Called()
	return args.Error(0)
}

func (m *LbcMock) RegisterProvider(txConfig blockchain.TransactionConfig, params blockchain.ProviderRegistrationParams) (int64, error) {
	args := m.Called(txConfig, params)
	return args.Get(0).(int64), args.Error(1)
}

func (m *LbcMock) GetProviders() ([]liquidity_provider.RegisteredLiquidityProvider, error) {
	args := m.Called()
	return args.Get(0).([]liquidity_provider.RegisteredLiquidityProvider), args.Error(1)
}

func (m *LbcMock) GetProvider(address string) (liquidity_provider.RegisteredLiquidityProvider, error) {
	args := m.Called(address)
	return args.Get(0).(liquidity_provider.RegisteredLiquidityProvider), args.Error(1)
}

func (m *LbcMock) SetProviderStatus(id uint64, status bool) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *LbcMock) GetPenalizedEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]penalization.PenalizedEvent, error) {
	args := m.Called(ctx, fromBlock, toBlock)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]penalization.PenalizedEvent), args.Error(1)
}

func (m *LbcMock) HashPeginQuote(peginQuote quote.PeginQuote) (string, error) {
	args := m.Called(peginQuote)
	return args.String(0), args.Error(1)
}

func (m *LbcMock) HashPegoutQuote(pegoutQuote quote.PegoutQuote) (string, error) {
	args := m.Called(pegoutQuote)
	return args.String(0), args.Error(1)
}

func (m *LbcMock) WithdrawPegoutCollateral() error {
	args := m.Called()
	return args.Error(0)
}

func (m *LbcMock) WithdrawCollateral() error {
	args := m.Called()
	return args.Error(0)
}

func (m *LbcMock) GetDepositEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]quote.PegoutDeposit, error) {
	args := m.Called(ctx, fromBlock, toBlock)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]quote.PegoutDeposit), args.Error(1)
}

func (m *LbcMock) RefundPegout(txConfig blockchain.TransactionConfig, params blockchain.RefundPegoutParams) (string, error) {
	args := m.Called(txConfig, params)
	return args.String(0), args.Error(1)
}

func (m *LbcMock) GetBalance(address string) (*entities.Wei, error) {
	args := m.Called(address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func (m *LbcMock) CallForUser(txConfig blockchain.TransactionConfig, peginQuote quote.PeginQuote) (string, error) {
	args := m.Called(txConfig, peginQuote)
	return args.String(0), args.Error(1)
}

func (m *LbcMock) RegisterPegin(params blockchain.RegisterPeginParams) (string, error) {
	args := m.Called(params)
	return args.String(0), args.Error(1)
}

func (m *LbcMock) IsPegOutQuoteCompleted(quoteHash string) (bool, error) {
	args := m.Called(quoteHash)
	return args.Bool(0), args.Error(1)
}

func (m *LbcMock) UpdateProvider(name, url string) (string, error) {
	args := m.Called(name, url)
	return args.String(0), args.Error(1)
}

func (m *LbcMock) RefundUserPegOut(quoteHash string) (string, error) {
	args := m.Called(quoteHash)
	return args.String(0), args.Error(1)
}
