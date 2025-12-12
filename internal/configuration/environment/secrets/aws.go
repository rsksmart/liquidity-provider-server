package secrets

import (
	"context"
	"encoding/json"
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
	if loader.env.Rsk.WalletSecret == "" || loader.env.Rsk.PasswordSecret == "" {
		return DerivativeWalletSecrets{}, errors.New("missing encrypted json or password secret")
	}
	walletInput := &secretsmanager.GetSecretValueInput{SecretId: &loader.env.Rsk.WalletSecret}
	walletSecret, err := loader.secretsManager.GetSecretValue(ctx, walletInput)
	if err != nil {
		return DerivativeWalletSecrets{}, fmt.Errorf("error loading encrypted json: %w", err)
	}
	var parsedWalletSecret walletSecretLayout
	if err = json.Unmarshal([]byte(*walletSecret.SecretString), &parsedWalletSecret); err != nil {
		return DerivativeWalletSecrets{}, errors.New("error parsing wallet secret")
	}

	jsonPasswordInput := &secretsmanager.GetSecretValueInput{SecretId: &loader.env.Rsk.PasswordSecret}
	jsonPassword, err := loader.secretsManager.GetSecretValue(ctx, jsonPasswordInput)
	if err != nil {
		return DerivativeWalletSecrets{}, fmt.Errorf("error loading encrypted json password: %w", err)
	}

	return DerivativeWalletSecrets{
		ColdWalletConfiguration: ColdWalletConfiguration{
			Type:          parsedWalletSecret.ColdWallet.Type,
			Configuration: parsedWalletSecret.ColdWallet.Configuration,
		},
		EncryptedJson:         string(parsedWalletSecret.HotWallet),
		EncryptedJsonPassword: *jsonPassword.SecretString,
	}, nil
}

func (loader *AwsSecretsLoader) LoadFireBlocksSecrets(ctx context.Context) (FireBlocksWalletSecrets, error) {
	// TODO complete with fireblocks integration
	panic("feature unavailable")
}
