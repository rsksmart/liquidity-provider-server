package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rsksmart/liquidity-provider-server/cmd/utils/scripts"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"golang.org/x/term"
)

type ResignUtilsInput struct {
	scripts.BaseInput
	Resign             bool
	WithdrawCollateral bool
}

func main() {
	const errorCode = 2
	scripts.SetUsageMessage(
		"This script is used to resign from the Liquidity Provider system and withdraw collateral from the Liquidity Bridge Contract.",
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

	if scriptInput.Resign {
		if err = ExecuteResign(lbc); err != nil {
			scripts.ExitWithError(errorCode, "Error executing resign", err)
		}
		fmt.Println("Resign executed successfully.")
		return
	}

	if err = ExecuteWithdrawCollateral(lbc); err != nil {
		if errors.Is(err, usecases.ProviderNotResignedError) {
			scripts.ExitWithError(errorCode, "Withdraw collateral rejected", err)
		}
		scripts.ExitWithError(errorCode, "Error executing withdraw collateral", err)
	}
	fmt.Println("Withdraw collateral executed successfully.")
}

func ReadResignUtilsInput(scriptInput *ResignUtilsInput) {
	scripts.ReadBaseInput(&scriptInput.BaseInput)
	flag.BoolVar(&scriptInput.Resign, "resign", false, "Execute the resign operation")
	flag.BoolVar(&scriptInput.WithdrawCollateral, "withdraw-collateral", false, "Withdraw collateral after resignation delay")
}

func ParseResignUtilsInput(parseFunc scripts.ParseFunc, input *ResignUtilsInput) error {
	parseFunc()
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(input); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	if input.Resign == input.WithdrawCollateral {
		return errors.New("select exactly one action: --resign or --withdraw-collateral")
	}
	return nil
}

func ExecuteResign(collateralContract interface {
	ProviderResign() error
}) error {
	return collateralContract.ProviderResign()
}

func ExecuteWithdrawCollateral(collateralContract interface {
	WithdrawCollateral() error
}) error {
	return collateralContract.WithdrawCollateral()
}
