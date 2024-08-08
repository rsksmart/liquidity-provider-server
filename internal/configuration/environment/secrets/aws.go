package secrets

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
)

type AwsSecretsLoader struct {
	config         aws.Config
	secretsManager *secretsmanager.Client
	env            environment.Environment
}

func NewAwsSecretsLoader(ctx context.Context, env environment.Environment) (SecretLoader, error) {
	awsConfiguration, err := environment.GetAwsConfig(ctx, env)
	if err != nil {
		return nil, fmt.Errorf("error loading secrets: %w", err)
	}
	return &AwsSecretsLoader{
		config:         awsConfiguration,
		secretsManager: secretsmanager.NewFromConfig(awsConfiguration),
		env:            env,
	}, nil
}

func (loader *AwsSecretsLoader) LoadDerivativeSecrets(ctx context.Context) (DerivativeWalletSecrets, error) {
	if loader.env.Rsk.EncryptedJsonSecret == "" || loader.env.Rsk.EncryptedJsonPasswordSecret == "" {
		return DerivativeWalletSecrets{}, errors.New("missing encrypted json or password secret")
	}
	encryptedJsonInput := &secretsmanager.GetSecretValueInput{SecretId: &loader.env.Rsk.EncryptedJsonSecret}
	encryptedJson, err := loader.secretsManager.GetSecretValue(ctx, encryptedJsonInput)
	if err != nil {
		return DerivativeWalletSecrets{}, fmt.Errorf("error loading encrypted json: %w", err)
	}

	jsonPasswordInput := &secretsmanager.GetSecretValueInput{SecretId: &loader.env.Rsk.EncryptedJsonPasswordSecret}
	jsonPassword, err := loader.secretsManager.GetSecretValue(ctx, jsonPasswordInput)
	if err != nil {
		return DerivativeWalletSecrets{}, fmt.Errorf("error loading encrypted json password: %w", err)
	}

	return DerivativeWalletSecrets{
		EncryptedJson:         *encryptedJson.SecretString,
		EncryptedJsonPassword: *jsonPassword.SecretString,
	}, nil
}

func (loader *AwsSecretsLoader) LoadFireBlocksSecrets(ctx context.Context) (FireBlocksWalletSecrets, error) {
	// TODO complete with fireblocks integration
	panic("feature unavailable")
}
