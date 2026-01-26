package wallet

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment/secrets"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/cold_wallet"
)

type AbstractFactory interface {
	BitcoinMonitoringWallet(walletId string) (blockchain.BitcoinWallet, error)
	BitcoinPaymentWallet(walletId string) (blockchain.BitcoinWallet, error)
	RskWallet() (rootstock.RskSignerWallet, error)
	ColdWallet(rpc blockchain.Rpc) (cold_wallet.ColdWallet, error)
}

type FactoryCreationArgs struct {
	Ctx          context.Context
	Env          environment.Environment
	SecretLoader secrets.SecretLoader
	RskClient    *rootstock.RskClient
	Timeouts     environment.ApplicationTimeouts
}

func NewFactory(env environment.Environment, args FactoryCreationArgs) (AbstractFactory, error) {
	switch env.WalletManagement {
	case "native":
		return NewDerivativeFactory(args)
	case "fireblocks":
		return NewFireBlocksFactory(args)
	default:
		return nil, errors.New("unknown wallet management scheme")
	}
}
