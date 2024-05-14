package secrets

import (
	"context"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

type EnvSecretsLoader struct {
	env environment.Environment
}

func NewEnvSecretsLoader(environment environment.Environment) SecretLoader {
	return &EnvSecretsLoader{env: environment}
}

func (loader *EnvSecretsLoader) LoadDerivativeSecrets(ctx context.Context) (DerivativeWalletSecrets, error) {
	if loader.env.Rsk.KeystoreFile == "" || loader.env.Rsk.KeystorePassword == "" {
		return DerivativeWalletSecrets{}, errors.New("missing keystore file or password")
	}

	keystoreFile, err := os.Open(loader.env.Rsk.KeystoreFile)
	if err != nil {
		return DerivativeWalletSecrets{}, fmt.Errorf("error opening keystore file: %w", err)
	}

	defer func(file *os.File) {
		if closingErr := file.Close(); closingErr != nil {
			log.Error("Error closing keystore file:", closingErr)
		}
	}(keystoreFile)

	keystoreBytes, err := io.ReadAll(keystoreFile)
	if err != nil {
		return DerivativeWalletSecrets{}, fmt.Errorf("error reading keystore file: %w", err)
	}

	return DerivativeWalletSecrets{
		EncryptedJson:         string(keystoreBytes),
		EncryptedJsonPassword: loader.env.Rsk.KeystorePassword,
	}, nil
}

func (loader *EnvSecretsLoader) LoadFireBlocksSecrets(ctx context.Context) (FireBlocksWalletSecrets, error) {
	// TODO complete with fireblocks integration
	panic("feature unavailable")
}
