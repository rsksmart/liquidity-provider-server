package environment

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	log "github.com/sirupsen/logrus"
)

func GetAwsConfig(ctx context.Context, env Environment) (aws.Config, error) {
	if env.LpsStage != "regtest" {
		return config.LoadDefaultConfig(ctx)
	}

	log.Debugf("Running in regtest mode. Using localstack endpoint (%s)", env.AwsLocalEndpoint)
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           env.AwsLocalEndpoint,
			SigningRegion: region,
		}, nil
	})
	return config.LoadDefaultConfig(ctx, config.WithEndpointResolverWithOptions(customResolver))
}
