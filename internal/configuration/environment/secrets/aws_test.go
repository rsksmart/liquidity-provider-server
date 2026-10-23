package secrets_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment/secrets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAwsSecretsLoader_LoadFireBlocksSecrets(t *testing.T) {
	t.Run("should return error when api key secret name is missing", func(t *testing.T) {
		loader, err := secrets.NewAwsSecretsLoader(context.Background(), environment.Environment{
			Fireblocks: environment.FireblocksEnv{SecretKeySecret: "secret-key"},
		})
		require.NoError(t, err)
		result, err := loader.LoadFireBlocksSecrets(context.Background())
		require.ErrorContains(t, err, "missing fireblocks api key or secret key secret")
		assert.Empty(t, result)
	})
	t.Run("should return error when secret key secret name is missing", func(t *testing.T) {
		loader, err := secrets.NewAwsSecretsLoader(context.Background(), environment.Environment{
			Fireblocks: environment.FireblocksEnv{ApiKeySecret: "api-key"},
		})
		require.NoError(t, err)
		result, err := loader.LoadFireBlocksSecrets(context.Background())
		require.ErrorContains(t, err, "missing fireblocks api key or secret key secret")
		assert.Empty(t, result)
	})
}
