package environment

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	log "github.com/sirupsen/logrus"
)

func GetAwsConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx)
}

func AwsLocalEndpoint(env Environment) string {
	if env.LpsStage != "regtest" {
		return ""
	}
	log.Debugf("Running in regtest mode. Using localstack endpoint (%s)", env.AwsLocalEndpoint)
	return env.AwsLocalEndpoint
}
