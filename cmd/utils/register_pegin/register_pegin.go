package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rsksmart/liquidity-provider-server/cmd/utils/scripts"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/btc_bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"golang.org/x/term"
	"os"
)

type RegisterPegInScriptInput struct {
	scripts.BaseInput
	InputFilePath  string `validate:"required,filepath"`
	BtcRcpHost     string `validate:"required"`
	BtcRpcUser     string
	BtcRpcPassword string
}

type ParsedRegisterPegInInput struct {
	Quote     quote.PeginQuote
	Signature []byte
	BtcTxHash string
}

func (input RegisterPegInScriptInput) ToEnv(pwdReader scripts.PasswordReader) (environment.Environment, error) {
	const dummyCredential = "none"

	env, err := input.BaseInput.ToEnv(pwdReader)
	if err != nil {
		return environment.Environment{}, err
	}

	env.Btc.Endpoint = input.BtcRcpHost

	// if we provide empty credentials the client library will try to authenticate with a cookie, instead
	// if we provide dummy credentials, it will be able to connect to public RPC services if needed
	if input.BtcRpcUser != "" {
		env.Btc.Username = input.BtcRpcUser
	} else {
		env.Btc.Username = dummyCredential
	}
	if input.BtcRpcPassword != "" {
		env.Btc.Password = input.BtcRpcPassword
	} else {
		env.Btc.Password = dummyCredential
	}
	return env, nil
}

func main() {
	const errorCode = 2
	var err error
	scripts.SetUsageMessage(
		"This script is used to register a PegIn transaction in the Liquidity Bridge Contract." +
			" It can be used to refund a Peg In if something went wrong during the Flyover Protocol process." +
			" In order to execute such refund, an input file with the details of the quote, the LP signature of the quote," +
			" and the hash of the Bitcoin transaction to be registered must be provided.",
	)
	defer scripts.EnableSecureBuffers()()
	ctx := context.Background()

	scriptInput := new(RegisterPegInScriptInput)
	ReadRegisterPegInScriptInput(scriptInput)

	parsedInput, err := ParseRegisterPegInScriptInput(flag.Parse, scriptInput, os.ReadFile)
	if err != nil {
		scripts.ExitWithError(errorCode, "Error parsing input", err)
	}

	env, err := scriptInput.ToEnv(term.ReadPassword)
	if err != nil {
		scripts.ExitWithError(errorCode, "Error parsing input", err)
	}

	btcClient, err := btc_bootstrap.Bitcoin(env.Btc)
	if err != nil {
		scripts.ExitWithError(errorCode, "Error connecting to BTC node", err)
	}

	lbc, err := scripts.CreateLiquidityBridgeContract(ctx, bootstrap.Rootstock, env, environment.DefaultTimeouts())
	if err != nil {
		scripts.ExitWithError(errorCode, "Error accessing the Liquidity Bridge Contract", err)
	}

	txHash, err := ExecuteRegisterPegIn(bitcoin.NewBitcoindRpc(btcClient), lbc, parsedInput)
	if err != nil {
		scripts.ExitWithError(errorCode, "Error executing register PegIn", err)
	}
	fmt.Println("Register PegIn executed successfully. Transaction hash: ", txHash)
}

func ReadRegisterPegInScriptInput(scriptInput *RegisterPegInScriptInput) {
	scripts.ReadBaseInput(&scriptInput.BaseInput)
	flag.StringVar(&scriptInput.InputFilePath, "input-file", "", "The input file containing the PegIn quote, the provider signature and the bitcoin transaction hash to register")
	flag.StringVar(&scriptInput.BtcRcpHost, "btc-rpc-host", "", "The host of the Bitcoin RPC interface. It shouldn't include the scheme. E.g. localhost:5555")
	flag.StringVar(&scriptInput.BtcRpcUser, "btc-rpc-user", "", "The Bitcoin RPC username. You can skip it if you're using a public service.")
	flag.StringVar(&scriptInput.BtcRpcPassword, "btc-rpc-password", "", "The Bitcoin RPC password. You can skip it if you're using a public service.")
}

func ParseRegisterPegInScriptInput(parse scripts.ParseFunc, scriptInput *RegisterPegInScriptInput, reader scripts.FileReader) (ParsedRegisterPegInInput, error) {
	var rawInput struct {
		Quote     pkg.PeginQuoteDTO `json:"quote"`
		Signature string            `json:"signature"`
		BtcTxHash string            `json:"btcTxHash"`
	}

	parse()
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(scriptInput); err != nil {
		return ParsedRegisterPegInInput{}, fmt.Errorf("invalid input: %w", err)
	}

	fileBytes, err := reader(scriptInput.InputFilePath)
	if err != nil {
		return ParsedRegisterPegInInput{}, err
	}

	if err = json.Unmarshal(fileBytes, &rawInput); err != nil {
		return ParsedRegisterPegInInput{}, err
	}

	signatureBytes, err := hex.DecodeString(rawInput.Signature)
	if err != nil {
		return ParsedRegisterPegInInput{}, fmt.Errorf("invalid signature: %w", err)
	}
	return ParsedRegisterPegInInput{
		Quote:     pkg.FromPeginQuoteDTO(rawInput.Quote),
		Signature: signatureBytes,
		BtcTxHash: rawInput.BtcTxHash,
	}, nil
}

func ExecuteRegisterPegIn(
	btcRpc blockchain.BitcoinNetwork,
	lbc blockchain.LiquidityBridgeContract,
	parsedInput ParsedRegisterPegInInput,
) (string, error) {
	var pmt, rawTx []byte
	var err error

	if pmt, err = btcRpc.GetPartialMerkleTree(parsedInput.BtcTxHash); err != nil {
		return "", err
	}
	if rawTx, err = btcRpc.GetRawTransaction(parsedInput.BtcTxHash); err != nil {
		return "", err
	}
	blockInfo, err := btcRpc.GetTransactionBlockInfo(parsedInput.BtcTxHash)
	if err != nil {
		return "", err
	}

	return lbc.RegisterPegin(blockchain.RegisterPeginParams{
		QuoteSignature:        parsedInput.Signature,
		BitcoinRawTransaction: rawTx,
		PartialMerkleTree:     pmt,
		BlockHeight:           blockInfo.Height,
		Quote:                 parsedInput.Quote,
	})
}
