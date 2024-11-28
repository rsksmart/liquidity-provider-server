package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rsksmart/liquidity-provider-server/cmd/utils/scripts"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"golang.org/x/term"
)

type RefundUserPegOutScriptInput struct {
	scripts.BaseInput        // Embedding BaseInput
	QuoteHashBytes    string `validate:"required,hexadecimal"`
}

type PasswordReader = func(int) ([]byte, error)

func main() {
	scripts.SetUsageMessage(
		"This script is used to execute a refund for a PegOut transaction in the Liquidity Bridge Contract." +
			" It is intended for use when the final user does not receive their funds." +
			" To perform this refund, you must provide the hash of the quote agreed for the service.",
	)
	scriptInput := new(RefundUserPegOutScriptInput)
	ReadRefundUserPegOutScriptInput(scriptInput)
	env, err := ParseRefundUserPegOutScriptInput(flag.Parse, scriptInput, term.ReadPassword)
	if err != nil {
		scripts.ExitWithError(2, "Error reading input", err)
	}

	ctx := context.Background()
	lbc, err := scripts.CreateLiquidityBridgeContract(ctx, bootstrap.Rootstock, env)
	if err != nil {
		scripts.ExitWithError(2, "Error accessing the Liquidity Bridge Contract", err)
	}

	txHash, err := ExecuteRefundUserPegOut(lbc, scriptInput.QuoteHashBytes)
	if err != nil {
		scripts.ExitWithError(2, "Error on transaction execution", err)
	}
	fmt.Println("Refund user peg out executed successfully. Transaction hash: ", txHash)
}

func ReadRefundUserPegOutScriptInput(scriptInput *RefundUserPegOutScriptInput) {
	flag.StringVar(&scriptInput.Network, "network", "", "The network to execute the script. Must be one of the following: regtest, testnet, mainnet")
	flag.StringVar(&scriptInput.QuoteHashBytes, "quote-hash", "", "The quote hash to refund the user peg out")

	flag.StringVar(&scriptInput.AwsLocalEndpoint, "aws-endpoint", "http://localhost:4566", "AWS endpoint for localstack")
	flag.StringVar(&scriptInput.SecretSource, "secret-src", "", "The source of the secrets to execute the transaction. Must be one of the following: env, aws")
	flag.StringVar(&scriptInput.RskEndpoint, "rsk-endpoint", "", "The URL of the RSK RPC server. E.g. http://localhost:4444")
	flag.StringVar(&scriptInput.CustomLbcAddress, "lbc-address", "", "Custom address of the liquidity bridge contract. If not provided will use the network default.")

	flag.StringVar(&scriptInput.KeystoreFile, "keystore-file", "", "Path to the keystore file. Only required if the secret source is env")
	flag.StringVar(&scriptInput.EncryptedJsonSecret, "keystore-secret", "", "Name of the secret storing the keystore. Only required if the secret source is aws")
	flag.StringVar(&scriptInput.EncryptedJsonPasswordSecret, "password-secret", "", "Name of the secret storing the keystore password. Only required if the secret source is aws")
}

func ParseRefundUserPegOutScriptInput(parse scripts.ParseFunc, scriptInput *RefundUserPegOutScriptInput, pwdReader PasswordReader) (environment.Environment, error) {
	parse()
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(scriptInput)
	if err != nil {
		return environment.Environment{}, fmt.Errorf("invalid input: %w", err)
	}

	return scriptInput.BaseInput.ToEnv(pwdReader)
}

func ExecuteRefundUserPegOut(lbc blockchain.LiquidityBridgeContract, quoteHash string) (string, error) {
	return lbc.RefundUserPegOut(quoteHash)
}
