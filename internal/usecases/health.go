package usecases

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

const (
	svcStatusOk          = "ok"
	svcStatusDegraded    = "degraded"
	svcStatusUnreachable = "unreachable"
)

type HealthUseCase struct {
	rsk entities.Service
	btc entities.Service
	db  entities.Service
}

func NewHealthUseCase(rsk entities.Service, btc entities.Service, db entities.Service) *HealthUseCase {
	return &HealthUseCase{rsk: rsk, btc: btc, db: db}
}

type HealthStatus struct {
	Status   string
	Services struct {
		Db  string
		Rsk string
		Btc string
	}
}

func (useCase *HealthUseCase) Run(ctx context.Context) HealthStatus {
	lpsSvcStatus := svcStatusOk
	dbSvcStatus := svcStatusOk
	rskSvcStatus := svcStatusOk
	btcSvcStatus := svcStatusOk

	if !useCase.db.CheckConnection(ctx) {
		dbSvcStatus = svcStatusUnreachable
		lpsSvcStatus = svcStatusDegraded
	}
	if !useCase.btc.CheckConnection(ctx) {
		btcSvcStatus = svcStatusUnreachable
		lpsSvcStatus = svcStatusDegraded
	}
	if !useCase.rsk.CheckConnection(ctx) {
		rskSvcStatus = svcStatusUnreachable
		lpsSvcStatus = svcStatusDegraded
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
