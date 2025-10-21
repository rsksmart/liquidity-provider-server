package defaults

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
)

func GetRsk(network string) (environment.RskEnv, error) {
	switch network {
	case "regtest":
		return environment.RskEnv{
			ChainId: 33,
			// TODO add addresses when deplpoyed
			PeginContractAddress:        "0x8901a2bbf639bfd21a97004ba4d7ae2bd00b8da8",
			PegoutContractAddress:       "0x8901a2bbf639bfd21a97004ba4d7ae2bd00b8da8",
			DiscoveryAddress:            "0x8901a2bbf639bfd21a97004ba4d7ae2bd00b8da8",
			CollateralManagementAddress: "0x8901a2bbf639bfd21a97004ba4d7ae2bd00b8da8",
			BridgeAddress:               "0x0000000000000000000000000000000001000006",
			AccountNumber:               0,
		}, nil
	case "testnet":
		return environment.RskEnv{
			ChainId: 31,
			// TODO add addresses when deplpoyed
			PeginContractAddress:        "0x8901a2bbf639bfd21a97004ba4d7ae2bd00b8da8",
			PegoutContractAddress:       "0x8901a2bbf639bfd21a97004ba4d7ae2bd00b8da8",
			DiscoveryAddress:            "0x8901a2bbf639bfd21a97004ba4d7ae2bd00b8da8",
			CollateralManagementAddress: "0x8901a2bbf639bfd21a97004ba4d7ae2bd00b8da8",
			BridgeAddress:               "0x0000000000000000000000000000000001000006",
			AccountNumber:               0,
		}, nil
	case "mainnet":
		return environment.RskEnv{
			ChainId: 30,
			// TODO add addresses when deplpoyed
			PeginContractAddress:        "0x8901a2bbf639bfd21a97004ba4d7ae2bd00b8da8",
			PegoutContractAddress:       "0x8901a2bbf639bfd21a97004ba4d7ae2bd00b8da8",
			DiscoveryAddress:            "0x8901a2bbf639bfd21a97004ba4d7ae2bd00b8da8",
			CollateralManagementAddress: "0x8901a2bbf639bfd21a97004ba4d7ae2bd00b8da8",
			BridgeAddress:               "0x0000000000000000000000000000000001000006",
			AccountNumber:               0,
		}, nil
	default:
		return environment.RskEnv{}, errors.New("invalid network")
	}
}
