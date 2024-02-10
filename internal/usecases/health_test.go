package usecases_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	u "github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type serviceMock struct {
	mock.Mock
	entities.Service
	Fail bool
}

func (m *serviceMock) CheckConnection(ctx context.Context) bool {
	return !m.Fail
}

func TestHealthUseCase_Run(t *testing.T) {
	ctx := context.Background()
	rsk := &serviceMock{}
	btc := &serviceMock{}
	db := &serviceMock{}

	useCase := u.NewHealthUseCase(rsk, btc, db)

	rsk.Fail = true
	btc.Fail = false
	db.Fail = false
	badHealthRsk := useCase.Run(ctx)

	rsk.Fail = false
	btc.Fail = true
	db.Fail = false
	badHealthBtc := useCase.Run(ctx)

	rsk.Fail = false
	btc.Fail = false
	db.Fail = true
	badHealthDb := useCase.Run(ctx)

	rsk.Fail = false
	btc.Fail = false
	db.Fail = false
	healthOk := useCase.Run(ctx)

	assert.Equal(t, u.HealthStatus{
		Status: u.SvcStatusDegraded,
		Services: u.Services{
			Rsk: u.SvcStatusUnreachable,
			Btc: u.SvcStatusOk,
			Db:  u.SvcStatusOk,
		},
	}, badHealthRsk)
	assert.Equal(t, u.HealthStatus{
		Status: u.SvcStatusDegraded,
		Services: u.Services{
			Rsk: u.SvcStatusOk,
			Btc: u.SvcStatusUnreachable,
			Db:  u.SvcStatusOk,
		},
	}, badHealthBtc)
	assert.Equal(t, u.HealthStatus{
		Status: u.SvcStatusDegraded,
		Services: u.Services{
			Rsk: u.SvcStatusOk,
			Btc: u.SvcStatusOk,
			Db:  u.SvcStatusUnreachable,
		},
	}, badHealthDb)
	assert.Equal(t, u.HealthStatus{
		Status: u.SvcStatusOk,
		Services: u.Services{
			Rsk: u.SvcStatusOk,
			Btc: u.SvcStatusOk,
			Db:  u.SvcStatusOk,
		},
	}, healthOk)
}
