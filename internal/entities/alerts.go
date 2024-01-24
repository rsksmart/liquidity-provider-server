package entities

import "context"

type AlertSender interface {
	SendAlert(ctx context.Context, subject, body, recipient string) error
}
