package wallet

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment/secrets"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

type AbstractFactory interface {
	BitcoinMonitoringWallet(walletId string) (blockchain.BitcoinWallet, error)
	BitcoinPaymentWallet(walletId string) (blockchain.BitcoinWallet, error)
	RskWallet() (rootstock.RskSignerWallet, error)
}

type FactoryCreationArgs struct {
	Ctx          context.Context
	Env          environment.Environment
	SecretLoader secrets.SecretLoader
	RskClient    *rootstock.RskClient
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
