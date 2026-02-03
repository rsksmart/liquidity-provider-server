package mocks

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/stretchr/testify/mock"
)

type ProviderMock struct {
	liquidity_provider.LiquidityProvider
	liquidity_provider.PeginLiquidityProvider
	liquidity_provider.PegoutLiquidityProvider
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

func (m *ProviderMock) SignQuote(quoteHash string) (string, error) {
	args := m.Called(quoteHash)
	return args.String(0), args.Error(1)
}

func (m *ProviderMock) GeneralConfiguration(ctx context.Context) liquidity_provider.GeneralConfiguration {
	args := m.Called(ctx)
	return args.Get(0).(liquidity_provider.GeneralConfiguration)
}

func (m *ProviderMock) PeginConfiguration(ctx context.Context) liquidity_provider.PeginConfiguration {
	args := m.Called(ctx)
	return args.Get(0).(liquidity_provider.PeginConfiguration)
}

func (m *ProviderMock) PegoutConfiguration(ctx context.Context) liquidity_provider.PegoutConfiguration {
	args := m.Called(ctx)
	return args.Get(0).(liquidity_provider.PegoutConfiguration)
}

func (m *ProviderMock) StateConfiguration(ctx context.Context) liquidity_provider.StateConfiguration {
	args := m.Called(ctx)
	return args.Get(0).(liquidity_provider.StateConfiguration)
}

func (m *ProviderMock) AvailablePeginLiquidity(ctx context.Context) (*entities.Wei, error) {
	args := m.Called(ctx)
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func (m *ProviderMock) AvailablePegoutLiquidity(ctx context.Context) (*entities.Wei, error) {
	args := m.Called(ctx)
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func (m *ProviderMock) GetSigner() entities.Signer {
	args := m.Called()
	return args.Get(0).(entities.Signer)
}
