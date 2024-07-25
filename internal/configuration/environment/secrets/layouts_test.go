package secrets_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment/secrets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetSecretLoader(t *testing.T) {
	t.Run("should create aws secret loader", func(t *testing.T) {
		loader, err := secrets.GetSecretLoader(context.Background(), environment.Environment{SecretSource: "aws"})
		require.NoError(t, err)
		assert.IsType(t, &secrets.AwsSecretsLoader{}, loader)
	})
	t.Run("should create env secret loader", func(t *testing.T) {
		loader, err := secrets.GetSecretLoader(context.Background(), environment.Environment{SecretSource: "env"})
		require.NoError(t, err)
		assert.IsType(t, &secrets.EnvSecretsLoader{}, loader)
	})
	t.Run("should return error for unknown secret source", func(t *testing.T) {
		loader, err := secrets.GetSecretLoader(context.Background(), environment.Environment{SecretSource: "gcp"})
		require.Error(t, err)
		assert.Nil(t, loader)
	})
}
