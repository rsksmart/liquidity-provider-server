package secrets

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
)

type SecretLoader interface {
	LoadDerivativeSecrets(ctx context.Context) (DerivativeWalletSecrets, error)
	LoadFireBlocksSecrets(ctx context.Context) (FireBlocksWalletSecrets, error)
}

type DerivativeWalletSecrets struct {
	EncryptedJson         string
	EncryptedJsonPassword string
}

type FireBlocksWalletSecrets struct {
	// TODO complete with fireblocks integration
	// TBD
}

func GetSecretLoader(ctx context.Context, environment environment.Environment) (SecretLoader, error) {
	switch environment.SecretSource {
	case "aws":
		return NewAwsSecretsLoader(ctx, environment)
	case "env":
		return NewEnvSecretsLoader(environment), nil
	default:
		return nil, errors.New("unknown secret source")
	}
}
