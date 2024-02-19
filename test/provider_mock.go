package test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/stretchr/testify/mock"
)

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

func (m *ProviderMock) BtcAddress() string {
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

func (m *ProviderMock) ValidateAmountForPegin(amount *entities.Wei) error {
	args := m.Called(amount)
	return args.Error(0)
}

func (m *ProviderMock) ValidateAmountForPegout(amount *entities.Wei) error {
	args := m.Called(amount)
	return args.Error(0)
}

func (m *ProviderMock) PenaltyFeePegin() *entities.Wei {
	args := m.Called()
	return args.Get(0).(*entities.Wei)
}

func (m *ProviderMock) PenaltyFeePegout() *entities.Wei {
	args := m.Called()
	return args.Get(0).(*entities.Wei)
}

func (m *ProviderMock) GetRootstockConfirmationsForValue(value *entities.Wei) uint16 {
	args := m.Called(value)
	return args.Get(0).(uint16)
}

func (m *ProviderMock) GetBitcoinConfirmationsForValue(value *entities.Wei) uint16 {
	args := m.Called(value)
	return args.Get(0).(uint16)
}

func (m *ProviderMock) TimeForDepositPegin() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

func (m *ProviderMock) CallTime() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

func (m *ProviderMock) TimeForDepositPegout() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

func (m *ProviderMock) ExpireBlocksPegout() uint64 {
	args := m.Called()
	return args.Get(0).(uint64)
}

func (m *ProviderMock) SignQuote(quoteHash string) (string, error) {
	args := m.Called(quoteHash)
	return args.String(0), args.Error(1)
}
