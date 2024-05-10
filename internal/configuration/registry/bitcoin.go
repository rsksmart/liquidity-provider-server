package registry

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/wallet"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

type Bitcoin struct {
	MonitoringWallet blockchain.BitcoinWallet
	PaymentWallet    blockchain.BitcoinWallet
	RpcConnection    *bitcoin.Connection
}

func NewBitcoinRegistry(
	walletFactory wallet.AbstractFactory,
	connection *bitcoin.Connection,
) (*Bitcoin, error) {
	paymentWallet, err := walletFactory.BitcoinPaymentWallet(bitcoin.DerivativeWalletId)
	if err != nil {
		return nil, err
	}
	peginWatchOnly, err := walletFactory.BitcoinMonitoringWallet(bitcoin.PeginWalletId)
	if err != nil {
		return nil, err
	}
	return &Bitcoin{
		MonitoringWallet: peginWatchOnly,
		PaymentWallet:    paymentWallet,
		RpcConnection:    connection,
	}, nil
}
