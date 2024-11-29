package scripts

import (
	"flag"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/cmd/utils/defaults"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"syscall"
)

type BaseInput struct {
	Network                     string `validate:"required,oneof=regtest testnet mainnet"`
	RskEndpoint                 string `validate:"required,http_url"`
	CustomLbcAddress            string `validate:"omitempty,eth_addr"`
	AwsLocalEndpoint            string `validate:"http_url"`
	SecretSource                string `validate:"required,oneof=aws env"`
	EncryptedJsonSecret         string
	EncryptedJsonPasswordSecret string
	KeystoreFile                string `validate:"omitempty,filepath"`
	KeystorePassword            string
}

func ReadBaseInput(scriptInput *BaseInput) {
	flag.StringVar(&scriptInput.Network, "network", "", "The network to execute the script. Must be one of the following: regtest, testnet, mainnet")

	flag.StringVar(&scriptInput.AwsLocalEndpoint, "aws-endpoint", "http://localhost:4566", "AWS endpoint for localstack")
	flag.StringVar(&scriptInput.SecretSource, "secret-src", "", "The source of the secrets to execute the transaction. Must be one of the following: env, aws")
	flag.StringVar(&scriptInput.RskEndpoint, "rsk-endpoint", "", "The URL of the RSK RPC server. E.g. http://localhost:4444")
	flag.StringVar(&scriptInput.CustomLbcAddress, "lbc-address", "", "Custom address of the liquidity bridge contract. If not provided will use the network default.")

	flag.StringVar(&scriptInput.KeystoreFile, "keystore-file", "", "Path to the keystore file. Only required if the secret source is env")
	flag.StringVar(&scriptInput.EncryptedJsonSecret, "keystore-secret", "", "Name of the secret storing the keystore. Only required if the secret source is aws")
	flag.StringVar(&scriptInput.EncryptedJsonPasswordSecret, "password-secret", "", "Name of the secret storing the keystore password. Only required if the secret source is aws")
}

func (input BaseInput) ToEnv(pwdReader PasswordReader) (environment.Environment, error) {
	var env environment.Environment
	var err error

	if input.SecretSource == "env" {
		var password []byte
		fmt.Println("Insert keystore password:")
		if password, err = pwdReader(syscall.Stdin); err != nil {
			return environment.Environment{}, fmt.Errorf("error reading password: %w", err)
		}
		input.KeystorePassword = string(password)
	}

	rskEnvDefaults, err := defaults.GetRsk(input.Network)
	if err != nil {
		return environment.Environment{}, fmt.Errorf("invalid input: %w", err)
	}

	var lbcAddress string
	if input.CustomLbcAddress != "" {
		lbcAddress = input.CustomLbcAddress
	} else {
		lbcAddress = rskEnvDefaults.LbcAddress
	}

	env.LpsStage = input.Network
	env.AwsLocalEndpoint = input.AwsLocalEndpoint
	env.SecretSource = input.SecretSource
	env.WalletManagement = "native"
	env.Rsk = environment.RskEnv{
		Endpoint:                    input.RskEndpoint,
		ChainId:                     rskEnvDefaults.ChainId,
		LbcAddress:                  lbcAddress,
		BridgeAddress:               rskEnvDefaults.BridgeAddress,
		AccountNumber:               rskEnvDefaults.AccountNumber,
		EncryptedJsonSecret:         input.EncryptedJsonSecret,
		EncryptedJsonPasswordSecret: input.EncryptedJsonPasswordSecret,
		KeystoreFile:                input.KeystoreFile,
		KeystorePassword:            input.KeystorePassword,
	}
	env.Btc = environment.BtcEnv{Network: input.Network}
	return env, nil
}
