package test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/mock"
)

type AlertSenderMock struct {
	entities.AlertSender
	mock.Mock
}

func (m *AlertSenderMock) SendAlert(ctx context.Context, subject, body string, recipients []string) error {
	args := m.Called(ctx, subject, body, recipients)
	return args.Error(0)
}

type BridgeMock struct {
	blockchain.RootstockBridge
	mock.Mock
}

func (m *BridgeMock) GetMinimumLockTxValue() (*entities.Wei, error) {
	args := m.Called()
	return args.Get(0).(*entities.Wei), args.Error(1)
}

type ProviderMock struct {
	entities.LiquidityProvider
	entities.PeginLiquidityProvider
	entities.PegoutLiquidityProvider
	mock.Mock
}

func (m *ProviderMock) RskAddress() string {
	args := m.Called()
	return args.String(0)
}

func (m *ProviderMock) HasPeginLiquidity(ctx context.Context, amount *entities.Wei) error {
	args := m.Called(ctx, amount)
	return args.Error(0)
}

func (m *ProviderMock) HasPegoutLiquidity(ctx context.Context, amount *entities.Wei) error {
	args := m.Called(ctx, amount)
	return args.Error(0)
}

func (m *ProviderMock) CallFeePegin() *entities.Wei {
	args := m.Called()
	return args.Get(0).(*entities.Wei)
}

func (m *ProviderMock) MinPegin() *entities.Wei {
	args := m.Called()
	return args.Get(0).(*entities.Wei)
}

func (m *ProviderMock) MaxPegin() *entities.Wei {
	args := m.Called()
	return args.Get(0).(*entities.Wei)
}

func (m *ProviderMock) MaxPeginConfirmations() uint16 {
	args := m.Called()
	return args.Get(0).(uint16)
}

func (m *ProviderMock) CallFeePegout() *entities.Wei {
	args := m.Called()
	return args.Get(0).(*entities.Wei)
}

func (m *ProviderMock) MinPegout() *entities.Wei {
	args := m.Called()
	return args.Get(0).(*entities.Wei)
}

func (m *ProviderMock) MaxPegout() *entities.Wei {
	args := m.Called()
	return args.Get(0).(*entities.Wei)
}

func (m *ProviderMock) MaxPegoutConfirmations() uint16 {
	args := m.Called()
	return args.Get(0).(uint16)
}

type LbcMock struct {
	blockchain.LiquidityBridgeContract
	mock.Mock
}

func (m *LbcMock) GetMinimumCollateral() (*entities.Wei, error) {
	args := m.Called()
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func (m *LbcMock) GetCollateral(providerAddress string) (*entities.Wei, error) {
	args := m.Called(providerAddress)
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func (m *LbcMock) GetPegoutCollateral(providerAddress string) (*entities.Wei, error) {
	args := m.Called(providerAddress)
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

func (m *LbcMock) GetProviders() ([]entities.RegisteredLiquidityProvider, error) {
	args := m.Called()
	return args.Get(0).([]entities.RegisteredLiquidityProvider), args.Error(1)
}

func (m *LbcMock) SetProviderStatus(id uint64, status bool) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *LbcMock) GetPeginPunishmentEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]entities.PunishmentEvent, error) {
	args := m.Called(ctx, fromBlock, toBlock)
	return args.Get(0).([]entities.PunishmentEvent), args.Error(1)
}
