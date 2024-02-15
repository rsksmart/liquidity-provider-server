package usecases

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

const (
	SvcStatusOk          = "ok"
	SvcStatusDegraded    = "degraded"
	SvcStatusUnreachable = "unreachable"
)

type HealthUseCase struct {
	rsk entities.Service
	btc entities.Service
	db  entities.Service
}

func NewHealthUseCase(rsk entities.Service, btc entities.Service, db entities.Service) *HealthUseCase {
	return &HealthUseCase{rsk: rsk, btc: btc, db: db}
}

type Services struct {
	Db  string
	Rsk string
	Btc string
}

type HealthStatus struct {
	Status   string
	Services Services
}

func (useCase *HealthUseCase) Run(ctx context.Context) HealthStatus {
	lpsSvcStatus := SvcStatusOk
	dbSvcStatus := SvcStatusOk
	rskSvcStatus := SvcStatusOk
	btcSvcStatus := SvcStatusOk

	if !useCase.db.CheckConnection(ctx) {
		dbSvcStatus = SvcStatusUnreachable
		lpsSvcStatus = SvcStatusDegraded
	}
	if !useCase.btc.CheckConnection(ctx) {
		btcSvcStatus = SvcStatusUnreachable
		lpsSvcStatus = SvcStatusDegraded
	}
	if !useCase.rsk.CheckConnection(ctx) {
		rskSvcStatus = SvcStatusUnreachable
		lpsSvcStatus = SvcStatusDegraded
	}

	return HealthStatus{
		Status: lpsSvcStatus,
		Services: struct {
			Db  string
			Rsk string
			Btc string
		}{
			Rsk: rskSvcStatus,
			Btc: btcSvcStatus,
			Db:  dbSvcStatus,
		},
	}
}
