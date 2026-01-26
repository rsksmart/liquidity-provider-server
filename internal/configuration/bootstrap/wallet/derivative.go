package wallet

import (
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/account"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/btc_bootstrap"
	cold_wallet_bootstrap "github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/cold_wallet"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment/secrets"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/cold_wallet"
	log "github.com/sirupsen/logrus"
)

type DerivativeWalletFactory struct {
	rskAccount       *account.RskAccount
	env              environment.Environment
	rskClient        *rootstock.RskClient
	timeouts         environment.ApplicationTimeouts
	coldWalletConfig secrets.ColdWalletConfiguration
}

func NewDerivativeFactory(args FactoryCreationArgs) (AbstractFactory, error) {
	applicationSecrets, err := args.SecretLoader.LoadDerivativeSecrets(args.Ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting application secrets: %w", err)
	}
	rskAccount, err := bootstrap.RootstockAccount(args.Env.Rsk, args.Env.Btc, applicationSecrets)
	if err != nil {
		return nil, fmt.Errorf("error connecting to RSK account: %w", err)
	}
	log.Debug("Connected to RSK account")
	return &DerivativeWalletFactory{
		rskAccount:       rskAccount,
		env:              args.Env,
		rskClient:        args.RskClient,
		timeouts:         args.Timeouts,
		coldWalletConfig: applicationSecrets.ColdWalletConfiguration,
	}, nil
}

func (factory *DerivativeWalletFactory) BitcoinMonitoringWallet(walletId string) (blockchain.BitcoinWallet, error) {
	walletConnection, err := btc_bootstrap.BitcoinWallet(factory.env.Btc, walletId)
	if err != nil {
		return nil, fmt.Errorf("error creating BTC monitoring connection: %w", err)
	}
	wallet, err := bitcoin.NewWatchOnlyWallet(walletConnection)
	if err != nil {
		return nil, err
	}
	log.Debug("Connected to BTC node wallet for monitoring")
	return wallet, nil
}

func (factory *DerivativeWalletFactory) BitcoinPaymentWallet(walletId string) (blockchain.BitcoinWallet, error) {
	walletConnection, err := btc_bootstrap.BitcoinWallet(factory.env.Btc, walletId)
	if err != nil {
		return nil, fmt.Errorf("error creating BTC payment connection: %w", err)
	}
	derivative, err := bitcoin.NewDerivativeWallet(walletConnection, factory.rskAccount)
	if err != nil {
		return nil, err
	}
	log.Debug("Connected to BTC node wallet for payments")
	return derivative, nil
}

func (factory *DerivativeWalletFactory) RskWallet() (rootstock.RskSignerWallet, error) {
	wallet := rootstock.NewRskWalletImpl(factory.rskClient, factory.rskAccount, factory.env.Rsk.ChainId, factory.timeouts.MiningWait.Seconds())
	return wallet, nil
}

func (factory *DerivativeWalletFactory) ColdWallet(rpc blockchain.Rpc) (cold_wallet.ColdWallet, error) {
	coldWallet, err := cold_wallet_bootstrap.Create(rpc, factory.coldWalletConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating cold wallet: %w", err)
	}
	if err = coldWallet.Init(); err != nil {
		return nil, fmt.Errorf("error initializing cold wallet: %w", err)
	}
	return coldWallet, nil
}
