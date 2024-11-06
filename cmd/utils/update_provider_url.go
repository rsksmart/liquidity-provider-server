package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/rsksmart/liquidity-provider-server/cmd/utils/defaults"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/wallet"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment/secrets"
	"golang.org/x/term"
	"net/url"
	"os"
	"syscall"
)

type UpdateProviderScriptInput struct {
	ProviderName                string `validate:"required"`
	ProviderUrl                 string `validate:"required,http_url"`
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

type PasswordReader = func(int) ([]byte, error)

type UpdateProviderArgs struct {
	Name    string
	url     *url.URL
	network string
}

func NewUpdateProviderArgs(name string, rawUrl string, network string) (UpdateProviderArgs, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return UpdateProviderArgs{}, err
	}
	return UpdateProviderArgs{Name: name, url: parsedUrl, network: network}, nil
}

func (args UpdateProviderArgs) Validate() error {
	if args.Name == "" {
		return errors.New("empty name")
	}
	if args.url.Scheme == "" || args.url.Host == "" {
		return errors.New("invalid url")
	}
	if args.network != "regtest" && args.url.Scheme != "https" {
		return errors.New("invalid url, not using https")
	}
	return nil
}

func (args UpdateProviderArgs) Url() string {
	return fmt.Sprintf("%s://%s", args.url.Scheme, args.url.Host)
}

func main() {
	scriptInput := new(UpdateProviderScriptInput)
	ReadUpdateProviderScriptInput(scriptInput)
	env, err := ParseUpdateProviderScriptInput(scriptInput, term.ReadPassword)
	if err != nil {
		ExitWithError(2, "Error reading input", err)
	}

	ctx := context.Background()

	rskClient, err := bootstrap.Rootstock(ctx, env.Rsk)
	if err != nil {
		ExitWithError(2, "Error connecting to RSK node", err)
	}
	rskWallet, err := GetWallet(ctx, env, rskClient)
	if err != nil {
		ExitWithError(2, "Error accessing to wallet", err)
	}

	updateArgs, err := NewUpdateProviderArgs(scriptInput.ProviderName, scriptInput.ProviderUrl, scriptInput.Network)
	if err != nil {
		ExitWithError(2, "Invalid provider information", err)
	}
	err = updateArgs.Validate()
	if err != nil {
		ExitWithError(2, "Invalid provider information", err)
	}
	err = ExecuteUpdateProvider(ctx, env, rskWallet, rskClient, updateArgs)
	if err != nil {
		ExitWithError(2, "Error on transaction execution", err)
	}
}

func ReadUpdateProviderScriptInput(scriptInput *UpdateProviderScriptInput) {
	flag.StringVar(&scriptInput.Network, "network", "", "The network to execute the script. Must be one of the following: regtest, testnet, mainnet")
	flag.StringVar(&scriptInput.ProviderName, "provider-name", "", "The liquidity provider name to display")
	flag.StringVar(&scriptInput.ProviderUrl, "provider-url", "", "The URL of the liquidity provider to be accessible by the users")

	flag.StringVar(&scriptInput.AwsLocalEndpoint, "aws-endpoint", "http://localhost:4566", "AWS endpoint for localstack")
	flag.StringVar(&scriptInput.SecretSource, "secret-src", "", "The source of the secrets to execute the transaction. Must be one of the following: env, aws")
	flag.StringVar(&scriptInput.RskEndpoint, "rsk-endpoint", "", "The URL of the RSK RPC server. E.g. http://localhost:4444")
	flag.StringVar(&scriptInput.CustomLbcAddress, "lbc-address", "", "Custom address of the liquidity bridge contract. If not provided will use the network default.")

	flag.StringVar(&scriptInput.KeystoreFile, "keystore-file", "", "Path to the keystore file. Only required if the secret source is env")
	flag.StringVar(&scriptInput.EncryptedJsonSecret, "keystore-secret", "", "Name of the secret storing the keystore. Only required if the secret source is aws")
	flag.StringVar(&scriptInput.EncryptedJsonPasswordSecret, "password-secret", "", "Name of the secret storing the keystore password. Only required if the secret source is aws")
}

func ParseUpdateProviderScriptInput(scriptInput *UpdateProviderScriptInput, pwdReader PasswordReader) (environment.Environment, error) {
	var env environment.Environment
	flag.Parse()
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(scriptInput)
	if err != nil {
		return environment.Environment{}, fmt.Errorf("invalid input: %w", err)
	}

	if scriptInput.SecretSource == "env" {
		var password []byte
		fmt.Println("Insert keystore password:")
		if password, err = pwdReader(syscall.Stdin); err != nil {
			return environment.Environment{}, fmt.Errorf("error reading password: %w", err)
		}
		scriptInput.KeystorePassword = string(password)
	}

	rskEnvDefaults, err := defaults.GetRsk(scriptInput.Network)
	if err != nil {
		return environment.Environment{}, fmt.Errorf("invalid input: %w", err)
	}

	var lbcAddress string
	if scriptInput.CustomLbcAddress != "" {
		lbcAddress = scriptInput.CustomLbcAddress
	} else {
		lbcAddress = rskEnvDefaults.LbcAddress
	}

	env.LpsStage = scriptInput.Network
	env.AwsLocalEndpoint = scriptInput.AwsLocalEndpoint
	env.SecretSource = scriptInput.SecretSource
	env.WalletManagement = "native"
	env.Rsk = environment.RskEnv{
		Endpoint:                    scriptInput.RskEndpoint,
		ChainId:                     rskEnvDefaults.ChainId,
		LbcAddress:                  lbcAddress,
		BridgeAddress:               rskEnvDefaults.BridgeAddress,
		AccountNumber:               rskEnvDefaults.AccountNumber,
		EncryptedJsonSecret:         scriptInput.EncryptedJsonSecret,
		EncryptedJsonPasswordSecret: scriptInput.EncryptedJsonPasswordSecret,
		KeystoreFile:                scriptInput.KeystoreFile,
		KeystorePassword:            scriptInput.KeystorePassword,
	}
	env.Btc = environment.BtcEnv{Network: scriptInput.Network}
	return env, nil
}

func GetWallet(
	ctx context.Context,
	env environment.Environment,
	rskClient *rootstock.RskClient,
) (rootstock.RskSignerWallet, error) {
	secretLoader, err := secrets.GetSecretLoader(ctx, env)
	if err != nil {
		return nil, err
	}
	walletFactory, err := wallet.NewFactory(env, wallet.FactoryCreationArgs{
		Ctx: ctx, Env: env, SecretLoader: secretLoader, RskClient: rskClient,
	})
	if err != nil {
		return nil, err
	}
	return walletFactory.RskWallet()
}

func ExecuteUpdateProvider(
	ctx context.Context,
	env environment.Environment,
	rskWallet rootstock.RskSignerWallet,
	rskClient *rootstock.RskClient,
	args UpdateProviderArgs,
) error {
	lbc, err := bindings.NewLiquidityBridgeContract(common.HexToAddress(env.Rsk.LbcAddress), rskClient.Rpc())
	if err != nil {
		return err
	}
	opts := &bind.TransactOpts{From: rskWallet.Address(), Signer: rskWallet.Sign}
	tx, err := lbc.UpdateProvider(opts, args.Name, args.Url())
	if err != nil {
		return err
	}

	receipt, err := bind.WaitMined(ctx, rskClient.Rpc(), tx)
	if err != nil {
		return err
	}

	if receipt.Status == 1 {
		fmt.Println("Provider information updated successfully. Transaction hash: ", receipt.TxHash.Hex())
		return nil
	} else {
		return fmt.Errorf("transaction %s failed", receipt.TxHash.Hex())
	}
}

func ExitWithError(code int, message string, err error) {
	fmt.Println(fmt.Sprintf("%s: %s", message, err.Error()))
	os.Exit(code)
}
