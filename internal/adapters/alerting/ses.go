package alerting

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	log "github.com/sirupsen/logrus"
)

type SesAlertSender struct {
	sesClient *ses.Client
	from      string
}

func NewSesAlertSender(sesClient *ses.Client, from string) entities.AlertSender {
	return &SesAlertSender{sesClient: sesClient, from: from}
}

func (sender *SesAlertSender) SendAlert(ctx context.Context, subject, body, recipient string) error {
	log.Info("Sending alert to liquidity provider")
	result, err := sender.sesClient.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &sesTypes.Destination{
			ToAddresses: []string{recipient},
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
