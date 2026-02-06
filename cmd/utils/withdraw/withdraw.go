package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/big"

	"github.com/go-playground/validator/v10"
	"github.com/rsksmart/liquidity-provider-server/cmd/utils/scripts"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"golang.org/x/term"
)

type WithdrawScriptInput struct {
	scripts.BaseInput
	All    bool
	Amount string
}

func main() {
	const errorCode = 2
	scripts.SetUsageMessage(
		"This script withdraws funds used for pegins from the Liquidity Bridge Contract.",
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
	rskClient, err := bootstrap.Rootstock(ctx, env)
	if err != nil {
		scripts.ExitWithError(errorCode, "Error connecting to RSK node", err)
	}

	rskWallet, err := scripts.GetWallet(ctx, env, environment.DefaultTimeouts(), rskClient)
	if err != nil {
		scripts.ExitWithError(errorCode, "Error accessing wallet", err)
	}

	lbc, err := scripts.CreateLiquidityBridgeContract(ctx, bootstrap.Rootstock, env, environment.DefaultTimeouts())
	if err != nil {
		scripts.ExitWithError(errorCode, "Error accessing Liquidity Bridge Contract", err)
	}

	if err = ExecuteWithdraw(lbc, rskWallet.Address().String(), scriptInput.All, scriptInput.Amount); err != nil {
		scripts.ExitWithError(errorCode, "Error executing withdraw", err)
	}

	fmt.Println("Withdraw executed successfully!")
}

func ReadWithdrawScriptInput(scriptInput *WithdrawScriptInput) {
	scripts.ReadBaseInput(&scriptInput.BaseInput)
	flag.BoolVar(&scriptInput.All, "all", false, "Withdraw all available funds")
	flag.StringVar(&scriptInput.Amount, "amount", "", "Amount to withdraw (in wei)")
}

func ParseWithdrawScriptInput(parseFunc scripts.ParseFunc, input *WithdrawScriptInput) error {
	parseFunc()

	// Validate that either --all or --amount is provided, but not both
	if input.All && input.Amount != "" {
		return errors.New("cannot use both --all and --amount flags")
	}
	if !input.All && input.Amount == "" {
		return errors.New("must provide either --all or --amount flag")
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(input); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}
	return nil
}

func ExecuteWithdraw(lbc interface {
	GetBalance(address string) (*entities.Wei, error)
	Withdraw(amount *entities.Wei) error
}, address string, all bool, amountStr string) error {
	var amount *entities.Wei

	if all {
		// Get the current balance
		balance, err := lbc.GetBalance(address)
		if err != nil {
			return fmt.Errorf("failed to get balance: %w", err)
		}
		amount = balance
		fmt.Printf("Withdrawing all funds: %s wei\n", amount.String())
	} else {
		// Parse the provided amount
		bigInt := new(big.Int)
		_, ok := bigInt.SetString(amountStr, 10)
		if !ok {
			return fmt.Errorf("invalid amount format: %s", amountStr)
		}
		amount = entities.NewBigWei(bigInt)
		fmt.Printf("Withdrawing: %s wei\n", amount.String())
	}

	return lbc.Withdraw(amount)
}
