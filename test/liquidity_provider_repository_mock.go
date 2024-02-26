package test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/stretchr/testify/mock"
)

type LpRepositoryMock struct {
	mock.Mock
	lp.LiquidityProviderRepository
}

func (m *LpRepositoryMock) GeneralConfiguration(ctx context.Context) (*entities.Signed[lp.GeneralConfiguration], error) {
	args := m.Called(ctx)
	return args.Get(0).(*entities.Signed[lp.GeneralConfiguration]), args.Error(1)
}

func (m *LpRepositoryMock) PeginConfiguration(ctx context.Context) (*entities.Signed[lp.PeginConfiguration], error) {
	args := m.Called(ctx)
	return args.Get(0).(*entities.Signed[lp.PeginConfiguration]), args.Error(1)
}

func (m *LpRepositoryMock) PegoutConfiguration(ctx context.Context) (*entities.Signed[lp.PegoutConfiguration], error) {
	args := m.Called(ctx)
	return args.Get(0).(*entities.Signed[lp.PegoutConfiguration]), args.Error(1)
}

func (m *LpRepositoryMock) UpsertGeneralConfiguration(ctx context.Context, config entities.Signed[lp.GeneralConfiguration]) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *LpRepositoryMock) UpsertPeginConfiguration(ctx context.Context, config entities.Signed[lp.PeginConfiguration]) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *LpRepositoryMock) UpsertPegoutConfiguration(ctx context.Context, config entities.Signed[lp.PegoutConfiguration]) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}
