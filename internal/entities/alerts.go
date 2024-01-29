package entities

import "context"

type AlertSender interface {
	SendAlert(ctx context.Context, subject, body string, recipient []string) error
}
