package registry

import (
	"context"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/alerting"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	log "github.com/sirupsen/logrus"
)

func NewAlertSender(ctx context.Context, env environment.Environment) entities.AlertSender {
	awsConfiguration, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("error loading alert sender: ", err)
	}
	sesClient := ses.NewFromConfig(awsConfiguration)
	return alerting.NewSesAlertSender(sesClient, env.Provider.AlertSenderEmail)
}
