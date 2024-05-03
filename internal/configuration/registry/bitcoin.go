package registry

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/account"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

type Bitcoin struct {
	MonitoringWallet           blockchain.BitcoinWallet
	PaymentWallet              blockchain.BitcoinWallet
	MonitoringWalletConnection *bitcoin.Connection
	PaymentWalletConnection    *bitcoin.Connection
}

func NewBitcoinRegistry(
	env environment.BtcEnv,
	secrets environment.ApplicationSecrets,
	monitoringWalletConnection *bitcoin.Connection,
	paymentWalletConnection *bitcoin.Connection,
	rskAccount *account.RskAccount,
) (*Bitcoin, error) {
	if monitoringWalletConnection.WalletId == "" {
		return nil, errors.New("monitoringWalletConnection must be a wallet connection to the RPC server")
	}
	if paymentWalletConnection.WalletId == "" {
		return nil, errors.New("paymentWalletConnection must be a wallet connection to the RPC server")
	}
	bitcoind := bitcoin.NewBitcoindWallet(monitoringWalletConnection, env.BtcAddress, env.FixedTxFeeRate, env.WalletEncrypted, secrets.BtcWalletPassword)
	if err := bitcoind.Unlock(); err != nil {
		return nil, err
	}
	derivative, err := bitcoin.NewDerivativeWallet(paymentWalletConnection, rskAccount)
	if err != nil {
		return nil, err
	}
	return &Bitcoin{
		MonitoringWallet:           bitcoind,
		PaymentWallet:              derivative,
		MonitoringWalletConnection: monitoringWalletConnection,
		PaymentWalletConnection:    paymentWalletConnection,
	}, nil
}
