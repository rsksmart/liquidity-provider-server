package secrets

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
)

type walletSecretLayout struct {
	HotWallet  json.RawMessage `json:"hotWallet"`
	ColdWallet struct {
		Type          string          `json:"type"`
		Configuration json.RawMessage `json:"configuration"`
	} `json:"coldWallet"`
}

type SecretLoader interface {
	LoadDerivativeSecrets(ctx context.Context) (DerivativeWalletSecrets, error)
	LoadFireBlocksSecrets(ctx context.Context) (FireBlocksWalletSecrets, error)
}

type ColdWalletConfiguration struct {
	Type          string
	Configuration json.RawMessage
}

type DerivativeWalletSecrets struct {
	ColdWalletConfiguration
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
