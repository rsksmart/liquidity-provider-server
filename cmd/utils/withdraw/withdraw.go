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

type WithdrawScriptInput struct {
	scripts.BaseInput
}

func main() {
	const errorCode = 2
	scripts.SetUsageMessage(
		"This script withdraws all collateral from the Liquidity Bridge Contract.",
	)
	defer scripts.EnableSecureBuffers()()

	scriptInput := new(WithdrawScriptInput)
	ReadWithdrawScriptInput(scriptInput)

	if err := ParseWithdrawScriptInput(flag.Parse, scriptInput); err != nil {
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

	if err = ExecuteWithdrawCollateral(lbc); err != nil {
		scripts.ExitWithError(errorCode, "Error executing withdraw collateral", err)
	}

	fmt.Println("Withdraw collateral executed successfully!")
}

func ReadWithdrawScriptInput(scriptInput *WithdrawScriptInput) {
	scripts.ReadBaseInput(&scriptInput.BaseInput)
}

func ParseWithdrawScriptInput(parseFunc scripts.ParseFunc, input *WithdrawScriptInput) error {
	parseFunc()
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(input); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}
	return nil
}

func ExecuteWithdrawCollateral(lbc interface {
	WithdrawCollateral() error
}) error {
	return lbc.WithdrawCollateral()
}
