package registry

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/alerting"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
)

func NewAlertSender(ctx context.Context, env environment.Environment) alerts.AlertSender {
	return alerting.NewLogAlertSender()
}
