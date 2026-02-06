package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rsksmart/liquidity-provider-server/cmd/utils/scripts"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"golang.org/x/term"
)

type ResignUtilsInput struct {
	scripts.BaseInput
}

func main() {
	const errorCode = 2
	scripts.SetUsageMessage(
		"This script is used to resign from the Liquidity Provider system.",
	)
	defer scripts.EnableSecureBuffers()()

	scriptInput := new(ResignUtilsInput)
	ReadResignUtilsInput(scriptInput)

	if err := ParseResignUtilsInput(flag.Parse, scriptInput); err != nil {
		scripts.ExitWithError(errorCode, "Error parsing input", err)
	}

	env, err := scriptInput.ToEnv(term.ReadPassword)
	if err != nil {
		scripts.ExitWithError(errorCode, "Error parsing input", err)
	}

	ctx := context.Background()
	lbc, err := scripts.CreateLiquidityBridgeContract(ctx, bootstrap.Rootstock, env, environment.DefaultTimeouts())
	if err != nil {
		scripts.ExitWithError(errorCode, "Error accessing Liquidity Bridge Contract", err)
	}

	if err = ExecuteResign(lbc); err != nil {
		scripts.ExitWithError(errorCode, "Error executing resign", err)
	}
	fmt.Println("Resign executed successfully.")
}

func ReadResignUtilsInput(scriptInput *ResignUtilsInput) {
	scripts.ReadBaseInput(&scriptInput.BaseInput)
}

func ParseResignUtilsInput(parseFunc scripts.ParseFunc, input *ResignUtilsInput) error {
	parseFunc()
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(input); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}
	return nil
}

func ExecuteResign(collateralContract interface {
	ProviderResign() error
}) error {
	return collateralContract.ProviderResign()
}
