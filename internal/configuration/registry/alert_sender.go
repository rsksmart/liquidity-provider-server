package registry

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/alerting"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	log "github.com/sirupsen/logrus"
)

func NewAlertSender(ctx context.Context, env environment.Environment) entities.AlertSender {
	awsConfiguration, err := environment.GetAwsConfig(ctx, env)
	if err != nil {
		log.Fatal("error loading alert sender: ", err)
	}
	sesClient := ses.NewFromConfig(awsConfiguration)
	return alerting.NewSesAlertSender(sesClient, env.Provider.AlertSenderEmail)
}
