package secrets

import (
	"context"
	"encoding/json"
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
	if loader.env.Rsk.WalletFile == "" || loader.env.Rsk.KeystorePassword == "" {
		return DerivativeWalletSecrets{}, errors.New("missing keystore file or password")
	}

	walletFile, err := os.Open(loader.env.Rsk.WalletFile)
	if err != nil {
		return DerivativeWalletSecrets{}, fmt.Errorf("error opening wallet file: %w", err)
	}

	defer func(file *os.File) {
		if closingErr := file.Close(); closingErr != nil {
			log.Error("Error closing wallet file:", closingErr)
		}
	}(walletFile)

	walletFileBytes, err := io.ReadAll(walletFile)
	if err != nil {
		return DerivativeWalletSecrets{}, fmt.Errorf("error reading wallet file: %w", err)
	}
	var parsedWalletSecret walletSecretLayout
	if err = json.Unmarshal(walletFileBytes, &parsedWalletSecret); err != nil {
		return DerivativeWalletSecrets{}, errors.New("error parsing wallet file")
	}

	return DerivativeWalletSecrets{
		ColdWalletConfiguration: ColdWalletConfiguration{
			Type:          parsedWalletSecret.ColdWallet.Type,
			Configuration: parsedWalletSecret.ColdWallet.Configuration,
		},
		EncryptedJson:         string(parsedWalletSecret.HotWallet),
		EncryptedJsonPassword: loader.env.Rsk.KeystorePassword,
	}, nil
}

func (loader *EnvSecretsLoader) LoadFireBlocksSecrets(ctx context.Context) (FireBlocksWalletSecrets, error) {
	// TODO complete with fireblocks integration
	panic("feature unavailable")
}
