package secrets

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManagerStorage[secretType any] struct {
	client *secretsmanager.Client
}

func NewSecretsManagerStorage[secretType any](config aws.Config) SecretStorage[secretType] {
	secretsManager := secretsmanager.NewFromConfig(config)
	return &SecretsManagerStorage[secretType]{client: secretsManager}
}

func (secretsManager *SecretsManagerStorage[secretType]) SaveJsonSecret(name string, secret *secretType) error {
	if secretBytes, err := json.Marshal(secret); err != nil {
		return err
	} else {
		return secretsManager.SaveTextSecret(name, string(secretBytes))
	}
}

func (secretsManager *SecretsManagerStorage[secretType]) SaveTextSecret(name, secret string) error {
	ctx, cancel := context.WithTimeout(context.Background(), WriteTimeout)
	defer cancel()

	_, err := secretsManager.client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         &name,
		SecretString: &secret,
	})

	return err
}

func (secretsManager *SecretsManagerStorage[secretType]) GetJsonSecret(name string) (*secretType, error) {
	secretString, err := secretsManager.GetTextSecret(name)
	if err != nil {
		return nil, err
	}

	var secret secretType
	if err = json.Unmarshal([]byte(secretString), &secret); err != nil {
		return nil, err
	}
	return &secret, nil
}

func (secretsManager *SecretsManagerStorage[secretType]) GetTextSecret(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ReadTimeout)
	defer cancel()

	value, err := secretsManager.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: &name})
	if err != nil {
		return "", err
	}
	return *value.SecretString, nil
}

func (secretsManager *SecretsManagerStorage[secretType]) DeleteSecret(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), WriteTimeout)
	defer cancel()

	_, err := secretsManager.client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{SecretId: &name})
	return err
}
