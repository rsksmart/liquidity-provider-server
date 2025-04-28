package scripts

import (
	"context"
	"flag"
	"fmt"
	"github.com/awnumar/memguard"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/wallet"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment/secrets"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"os"
)

type ParseFunc = func()
type PasswordReader = func(int) ([]byte, error)
type FileReader = func(string) ([]byte, error)
type RskClientFactory = func(context.Context, environment.RskEnv) (*rootstock.RskClient, error)

func ExitWithError(code int, message string, err error) {
	fmt.Printf("%s: %s\n", message, err.Error())
	memguard.Purge() // We add it here because exit doesn't execute deferred functions
	memguard.SafeExit(code)
}

func GetWallet(
	ctx context.Context,
	env environment.Environment,
	timeouts environment.ApplicationTimeouts,
	rskClient *rootstock.RskClient,
) (rootstock.RskSignerWallet, error) {
	secretLoader, err := secrets.GetSecretLoader(ctx, env)
	if err != nil {
		return nil, err
	}
	walletFactory, err := wallet.NewFactory(env, wallet.FactoryCreationArgs{
		Ctx: ctx, Env: env, SecretLoader: secretLoader, RskClient: rskClient, Timeouts: timeouts,
	})
	if err != nil {
		return nil, err
	}
	return walletFactory.RskWallet()
}

func CreateLiquidityBridgeContract(
	ctx context.Context,
	factory RskClientFactory,
	env environment.Environment,
	timeouts environment.ApplicationTimeouts,
) (blockchain.LiquidityBridgeContract, error) {
	rskClient, err := factory(ctx, env.Rsk)
	if err != nil {
		return nil, fmt.Errorf("error connecting to RSK node: %w", err)
	}
	rskWallet, err := GetWallet(ctx, env, timeouts, rskClient)
	if err != nil {
		return nil, fmt.Errorf("error accessing to wallet: %w", err)
	}
	lbc, err := bindings.NewLiquidityBridgeContract(common.HexToAddress(env.Rsk.LbcAddress), rskClient.Rpc())
	if err != nil {
		return nil, err
	}
	return rootstock.NewLiquidityBridgeContractImpl(
		rskClient,
		env.Rsk.LbcAddress,
		rootstock.NewLbcAdapter(lbc),
		rskWallet,
		rootstock.RetryParams{Retries: 0, Sleep: 0},
		environment.DefaultTimeouts().MiningWait.Seconds(),
	), nil
}

func SetUsageMessage(msg string) {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s\n", msg)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

// EnableSecureBuffers is a function that sets up secure memory buffers by calling memguard.CatchInterrupt().
// This function returns another function that must be called at the end of the program to purge the secure memory buffers.
// Every LPS script that interacts with a wallet must call this function at the beginning of the script.
// An example of the correct way of calling this function is:
//
//	func main() {
//		defer scripts.EnableSecureBuffers()()
//		// Your script logic here
//	}
func EnableSecureBuffers() func() {
	memguard.CatchInterrupt()
	fmt.Println("Secure buffers enabled")
	return func() {
		memguard.Purge()
		fmt.Println("Sensitive buffers were destroyed successfully")
	}
}
