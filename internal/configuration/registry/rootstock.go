package registry

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

type Rootstock struct {
	Contracts blockchain.RskContracts
	Wallet    *rootstock.RskWalletImpl
	Client    *rootstock.RskClient
}

func NewRootstockRegistry(env environment.RskEnv, client *rootstock.RskClient, account *rootstock.RskAccount, bitcoinConn *bitcoin.Connection) (*Rootstock, error) {
	var bridgeAddress, lbcAddress common.Address
	var err error

	if err = rootstock.ParseAddress(&lbcAddress, env.LbcAddress); err != nil {
		return nil, err
	} else if err = rootstock.ParseAddress(&bridgeAddress, env.BridgeAddress); err != nil {
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
	wallet := rootstock.NewRskWalletImpl(client, account, env.ChainId)

	return &Rootstock{
		Contracts: blockchain.RskContracts{
			Bridge: rootstock.NewRskBridgeImpl(
				env.BridgeAddress,
				env.BridgeRequiredConfirmations,
				env.IrisActivationHeight,
				env.ErpKeys,
				bridge,
				client,
				bitcoinConn.NetworkParams,
			),
			Lbc:          rootstock.NewLiquidityBridgeContractImpl(client, env.LbcAddress, lbc, wallet),
			FeeCollector: rootstock.NewFeeCollectorImpl(lbc),
		},
		Wallet: wallet,
		Client: client,
	}, nil
}
