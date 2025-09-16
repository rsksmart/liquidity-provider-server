package registry

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/wallet"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

type Rootstock struct {
	Contracts blockchain.RskContracts
	Wallet    rootstock.RskSignerWallet
	Client    *rootstock.RskClient
}

func NewRootstockRegistry(env environment.Environment, client *rootstock.RskClient, walletFactory wallet.AbstractFactory, timeouts environment.ApplicationTimeouts) (*Rootstock, error) {
	var (
		err                         error
		bridgeAddress               common.Address
		peginContractAddress        common.Address
		pegoutContractAddress       common.Address
		collateralManagementAddress common.Address
		discoveryAddress            common.Address
	)

	if err = rootstock.ParseAddress(&peginContractAddress, env.Rsk.PeginContractAddress); err != nil {
		return nil, err
	} else if err = rootstock.ParseAddress(&pegoutContractAddress, env.Rsk.PegoutContractAddress); err != nil {
		return nil, err
	} else if err = rootstock.ParseAddress(&collateralManagementAddress, env.Rsk.CollateralManagementAddress); err != nil {
		return nil, err
	} else if err = rootstock.ParseAddress(&discoveryAddress, env.Rsk.DiscoveryAddress); err != nil {
		return nil, err
	} else if err = rootstock.ParseAddress(&bridgeAddress, env.Rsk.BridgeAddress); err != nil {
		return nil, err
	}

	bridge, err := bindings.NewRskBridge(bridgeAddress, client.Rpc())
	if err != nil {
		return nil, err
	}

	peginContract, err := bindings.NewIPegIn(peginContractAddress, client.Rpc())
	if err != nil {
		return nil, err
	}
	pegoutContract, err := bindings.NewIPegOut(pegoutContractAddress, client.Rpc())
	if err != nil {
		return nil, err
	}
	collateralManagement, err := bindings.NewICollateralManagement(collateralManagementAddress, client.Rpc())
	if err != nil {
		return nil, err
	}
	discovery, err := bindings.NewIFlyoverDiscovery(discoveryAddress, client.Rpc())
	if err != nil {
		return nil, err
	}
	wallet, err := walletFactory.RskWallet()
	if err != nil {
		return nil, err
	}

	btcParams, err := env.Btc.GetNetworkParams()
	if err != nil {
		return nil, err
	}

	return &Rootstock{
		Contracts: blockchain.RskContracts{
			Bridge: rootstock.NewRskBridgeImpl(
				rootstock.RskBridgeConfig{
					Address:               env.Rsk.BridgeAddress,
					RequiredConfirmations: env.Rsk.BridgeRequiredConfirmations,
					IrisActivationHeight:  env.Rsk.IrisActivationHeight,
					ErpKeys:               env.Rsk.ErpKeys,
				},
				rootstock.NewRskBridgeAdapter(bridge),
				client,
				btcParams,
				rootstock.DefaultRetryParams,
				wallet,
				timeouts.MiningWait.Seconds(),
			),
			PegIn: rootstock.NewPeginContractImpl(
				client,
				env.Rsk.PeginContractAddress,
				rootstock.NewPeginContractAdapter(peginContract),
				wallet,
				rootstock.DefaultRetryParams,
				timeouts.MiningWait.Seconds(),
			),
			PegOut: rootstock.NewPegoutContractImpl(
				client,
				env.Rsk.PegoutContractAddress,
				rootstock.NewPegoutContractAdapter(pegoutContract),
				wallet,
				rootstock.DefaultRetryParams,
				timeouts.MiningWait.Seconds(),
			),
			CollateralManagement: rootstock.NewCollateralManagementContractImpl(
				client,
				wallet.Address().String(),
				env.Rsk.CollateralManagementAddress,
				rootstock.NewCollateralManagementAdapter(collateralManagement),
				wallet,
				rootstock.DefaultRetryParams,
				timeouts.MiningWait.Seconds(),
			),
			Discovery: rootstock.NewDiscoveryContractImpl(
				client,
				env.Rsk.DiscoveryAddress,
				discovery,
				wallet,
				rootstock.DefaultRetryParams,
				timeouts.MiningWait.Seconds(),
			),
		},
		Wallet: wallet,
		Client: client,
	}, nil
}
