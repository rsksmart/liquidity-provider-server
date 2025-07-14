package alerting

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
	log "github.com/sirupsen/logrus"
)

type sesClient interface {
	SendEmail(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error)
}

type SesAlertSender struct {
	sesClient sesClient
	from      string
}

func NewSesAlertSender(sesClient sesClient, from string) alerts.AlertSender {
	return &SesAlertSender{sesClient: sesClient, from: from}
}

func (sender *SesAlertSender) SendAlert(ctx context.Context, subject, body string, recipient []string) error {
	log.Info("Sending alert to liquidity provider")
	result, err := sender.sesClient.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &sesTypes.Destination{
			ToAddresses: recipient,
		},
		Message: &sesTypes.Message{
			Body: &sesTypes.Body{
				Text: &sesTypes.Content{Data: &body},
			},
			Subject: &sesTypes.Content{Data: &subject},
		},
		Source: &sender.from,
	})
	if err != nil {
		return err
	}
	log.Info("Alert sent")
	log.Debugf("Alert sent with ID: %s\n", *result.MessageId)
	return nil
}
