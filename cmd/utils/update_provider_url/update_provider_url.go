package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rsksmart/liquidity-provider-server/cmd/utils/scripts"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"golang.org/x/term"
	"net/url"
)

type UpdateProviderScriptInput struct {
	scripts.BaseInput
	ProviderName string `validate:"required"`
	ProviderUrl  string `validate:"required,http_url"`
}

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
	const errorCode = 2
	scripts.SetUsageMessage("This script is used to update the provider information displayed in the Liquidity Bridge Contract when the discovery function is executed.")
	scriptInput := new(UpdateProviderScriptInput)
	ReadUpdateProviderScriptInput(scriptInput)

	err := ParseUpdateProviderScriptInput(scriptInput)
	if err != nil {
		scripts.ExitWithError(errorCode, "Error parsing input", err)
	}
	env, err := scriptInput.ToEnv(term.ReadPassword)
	if err != nil {
		scripts.ExitWithError(errorCode, "Error parsing input", err)
	}

	ctx := context.Background()
	lbc, err := scripts.CreateLiquidityBridgeContract(ctx, bootstrap.Rootstock, env)
	if err != nil {
		scripts.ExitWithError(errorCode, "Error accessing the Liquidity Bridge Contract", err)
	}

	updateArgs, err := NewUpdateProviderArgs(scriptInput.ProviderName, scriptInput.ProviderUrl, scriptInput.Network)
	if err != nil {
		scripts.ExitWithError(errorCode, "Invalid provider information", err)
	} else if err = updateArgs.Validate(); err != nil {
		scripts.ExitWithError(errorCode, "Invalid provider information", err)
	}

	txHash, err := ExecuteUpdateProvider(lbc, updateArgs)
	if err != nil {
		scripts.ExitWithError(errorCode, "Error executing update provider", err)
	}
	fmt.Println("Update provider executed successfully. Transaction hash: ", txHash)
}

func ParseUpdateProviderScriptInput(input *UpdateProviderScriptInput) error {
	flag.Parse()
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(input); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}
	return nil
}

func ReadUpdateProviderScriptInput(scriptInput *UpdateProviderScriptInput) {
	scripts.ReadBaseInput(&scriptInput.BaseInput)
	flag.StringVar(&scriptInput.ProviderName, "provider-name", "", "The liquidity provider name to display")
	flag.StringVar(&scriptInput.ProviderUrl, "provider-url", "", "The URL of the liquidity provider to be accessible by the users")
}

func ExecuteUpdateProvider(
	lbc blockchain.LiquidityBridgeContract,
	args UpdateProviderArgs,
) (string, error) {
	return lbc.UpdateProvider(args.Name, args.Url())
}
