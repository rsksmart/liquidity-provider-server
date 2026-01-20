package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/rsksmart/liquidity-provider-server/cmd/utils/scripts"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
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
		"This script withdraws liquidity locked in the PegIn contract either fully or by a specified amount in wei.",
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
	peginContract, wallet, err := CreatePeginContractWithWallet(ctx, env)
	if err != nil {
		scripts.ExitWithError(errorCode, "Error accessing PegIn contract", err)
	}

	amount, err := ResolveWithdrawAmount(peginContract, wallet.Address().String(), scriptInput.All, scriptInput.Amount)
	if err != nil {
		scripts.ExitWithError(errorCode, "Invalid amount", err)
	}

	receipt, err := ExecuteWithdraw(peginContract, amount)
	if err != nil {
		scripts.ExitWithError(errorCode, "Error executing withdraw", err)
	}

	fmt.Println("Withdraw executed successfully!")
	PrintReceipt(receipt)
}

func ReadWithdrawScriptInput(scriptInput *WithdrawScriptInput) {
	scripts.ReadBaseInput(&scriptInput.BaseInput)
	flag.BoolVar(&scriptInput.All, "all", false, "Withdraw the full available balance")
	flag.StringVar(&scriptInput.Amount, "amount", "", "Withdraw a specific amount in wei")
}

func ParseWithdrawScriptInput(parseFunc scripts.ParseFunc, input *WithdrawScriptInput) error {
	parseFunc()
	if err := scripts.ApplyDefaultEnvFile(&input.BaseInput); err != nil {
		return fmt.Errorf("error loading default env file: %w", err)
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(input); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}
	if input.All == (input.Amount != "") {
		return errors.New("select exactly one action: --all or --amount <amount>")
	}
	return nil
}

func ParseWeiAmount(rawAmount string) (*entities.Wei, error) {
	value, ok := new(big.Int).SetString(rawAmount, 10)
	if !ok {
		return nil, errors.New("amount must be a base-10 integer in wei")
	}
	if value.Sign() <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	return entities.NewBigWei(value), nil
}

func ResolveWithdrawAmount(
	peginContract interface {
		GetBalance(address string) (*entities.Wei, error)
	},
	walletAddress string,
	withdrawAll bool,
	rawAmount string,
) (*entities.Wei, error) {
	if withdrawAll {
		balance, err := peginContract.GetBalance(walletAddress)
		if err != nil {
			return nil, fmt.Errorf("error getting available balance: %w", err)
		}
		if balance.AsBigInt().Sign() <= 0 {
			return nil, errors.New("balance is zero")
		}
		return balance, nil
	}
	return ParseWeiAmount(rawAmount)
}

func ExecuteWithdraw(
	peginContract interface {
		Withdraw(amount *entities.Wei) (blockchain.TransactionReceipt, error)
	},
	amount *entities.Wei,
) (blockchain.TransactionReceipt, error) {
	return peginContract.Withdraw(amount)
}

func CreatePeginContractWithWallet(
	ctx context.Context,
	env environment.Environment,
) (blockchain.PeginContract, rootstock.RskSignerWallet, error) {
	rskClient, err := bootstrap.Rootstock(ctx, env)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to RSK node: %w", err)
	}
	rskWallet, err := scripts.GetWallet(ctx, env, environment.DefaultTimeouts(), rskClient)
	if err != nil {
		return nil, nil, fmt.Errorf("error accessing to wallet: %w", err)
	}
	peginBinding, err := bindings.NewIPegIn(common.HexToAddress(env.Rsk.PeginContractAddress), rskClient.Rpc())
	if err != nil {
		return nil, nil, err
	}
	return rootstock.NewPeginContractImpl(
		rskClient,
		env.Rsk.PeginContractAddress,
		rootstock.NewPeginContractAdapter(peginBinding),
		rskWallet,
		rootstock.RetryParams{Retries: 0, Sleep: 0},
		environment.DefaultTimeouts().MiningWait.Seconds(),
		rootstock.MustLoadFlyoverABIs(),
	), rskWallet, nil
}

func PrintReceipt(receipt blockchain.TransactionReceipt) {
	fmt.Printf("Transaction Hash: %s\n", receipt.TransactionHash)
	fmt.Printf("Block Number: %d\n", receipt.BlockNumber)
	fmt.Printf("Gas Used: %s\n", receipt.GasUsed.String())
	fmt.Printf("Gas Price: %s wei\n", receipt.GasPrice.String())
	fmt.Printf("From: %s\n", receipt.From)
	fmt.Printf("To: %s\n", receipt.To)
	fmt.Printf("Value: %s wei\n", receipt.Value.String())
}
