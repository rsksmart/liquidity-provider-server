package defaults

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
)

func GetRsk(network string) (environment.RskEnv, error) {
	switch network {
	case "regtest":
		return environment.RskEnv{
			ChainId:       33,
			LbcAddress:    "0x8901a2bbf639bfd21a97004ba4d7ae2bd00b8da8",
			BridgeAddress: "0x0000000000000000000000000000000001000006",
			AccountNumber: 0,
		}, nil
	case "testnet":
		return environment.RskEnv{
			ChainId:       31,
			LbcAddress:    "0xc2a630c053d12d63d32b025082f6ba268db18300",
			BridgeAddress: "0x0000000000000000000000000000000001000006",
			AccountNumber: 0,
		}, nil
	case "mainnet":
		return environment.RskEnv{
			ChainId:       30,
			LbcAddress:    "0xaa9caf1e3967600578727f975f283446a3da6612",
			BridgeAddress: "0x0000000000000000000000000000000001000006",
			AccountNumber: 0,
		}, nil
	default:
		return environment.RskEnv{}, errors.New("invalid network")
	}
}
