package environment

import (
	"context"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	log "github.com/sirupsen/logrus"
)

type ApplicationSecrets struct {
	BtcWalletPassword     string
	EncryptedJson         string
	EncryptedJsonPassword string
}

func LoadSecrets(ctx context.Context, environment Environment) *ApplicationSecrets {
	return loadFromSecretsManager(ctx, environment)
}

func loadFromSecretsManager(ctx context.Context, env Environment) *ApplicationSecrets {
	awsConfiguration, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("error loading secrets: ", err)
	}
	sm := secretsmanager.NewFromConfig(awsConfiguration)
	encryptedJson, err := sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: &env.Rsk.EncryptedJsonSecret})
	if err != nil {
		log.Fatal("error loading encrypted json: ", err)
	}

	jsonPassword, err := sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: &env.Rsk.EncryptedJsonPasswordSecret})
	if err != nil {
		log.Fatal("error loading encrypted json password: ", err)
	}

	btcWalletPassword, err := sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: &env.Btc.WalletPasswordSecret})
	if err != nil {
		log.Fatal("error loading btc wallet password: ", err)
	}

	return &ApplicationSecrets{
		BtcWalletPassword:     *btcWalletPassword.SecretString,
		EncryptedJson:         *encryptedJson.SecretString,
		EncryptedJsonPassword: *jsonPassword.SecretString,
	}
}
