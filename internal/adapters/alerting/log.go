package alerting

import (
	"context"
	"errors"
	"strings"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
	log "github.com/sirupsen/logrus"
)

type LogAlertSender struct{}

func NewLogAlertSender() alerts.AlertSender {
	return &LogAlertSender{}
}

func (sender *LogAlertSender) SendAlert(ctx context.Context, subject, body string, recipient []string) error {
	if strings.TrimSpace(subject) == "" {
		return errors.New("alert subject cannot be empty")
	}

	recipients := strings.Join(recipient, ", ")
	log.Infof("Alert! - Subject: %s | Recipients: %s | Body: %s", subject, recipients, body)

	return nil
}
