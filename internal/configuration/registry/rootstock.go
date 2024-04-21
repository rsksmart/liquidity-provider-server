package registry

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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

	ethClient, ok := client.Rpc().(*ethclient.Client)
	if !ok {
		return nil, errors.New("invalid RSK client type, expected *ethclient.Client to build the registry")
	}

	bridge, err := bindings.NewRskBridge(bridgeAddress, ethClient)
	if err != nil {
		return nil, err
	}

	lbc, err := bindings.NewLiquidityBridgeContract(lbcAddress, ethClient)
	if err != nil {
		return nil, err
	}
	wallet := rootstock.NewRskWalletImpl(client, account, env.ChainId)

	return &Rootstock{
		Contracts: blockchain.RskContracts{
			Bridge: rootstock.NewRskBridgeImpl(
				rootstock.RskBridgeConfig{
					Address:               env.BridgeAddress,
					RequiredConfirmations: env.BridgeRequiredConfirmations,
					IrisActivationHeight:  env.IrisActivationHeight,
					ErpKeys:               env.ErpKeys,
				},
				bridge,
				client,
				bitcoinConn.NetworkParams,
				rootstock.DefaultRetryParams,
			),
			Lbc: rootstock.NewLiquidityBridgeContractImpl(
				client,
				env.LbcAddress,
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
