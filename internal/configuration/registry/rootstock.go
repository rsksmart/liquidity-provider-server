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

func NewRootstockRegistry(env environment.Environment, client *rootstock.RskClient, walletFactory wallet.AbstractFactory) (*Rootstock, error) {
	var bridgeAddress, lbcAddress common.Address
	var err error

	if err = rootstock.ParseAddress(&lbcAddress, env.Rsk.LbcAddress); err != nil {
		return nil, err
	} else if err = rootstock.ParseAddress(&bridgeAddress, env.Rsk.BridgeAddress); err != nil {
		return nil, err
	}

	bridge, err := bindings.NewRskBridge(bridgeAddress, client.Rpc())
	if err != nil {
		return nil, err
	}

	lbc, err := bindings.NewLiquidityBridgeContract(lbcAddress, client.Rpc())
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
				bridge,
				client,
				btcParams,
				rootstock.DefaultRetryParams,
				wallet,
			),
			Lbc: rootstock.NewLiquidityBridgeContractImpl(
				client,
				env.Rsk.LbcAddress,
				rootstock.NewLbcAdapter(lbc),
				wallet,
				rootstock.DefaultRetryParams,
			),
			FeeCollector: rootstock.NewFeeCollectorImpl(rootstock.NewLbcAdapter(lbc), rootstock.DefaultRetryParams),
		},
		Wallet: wallet,
		Client: client,
	}, nil
}
