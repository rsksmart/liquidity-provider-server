package registry

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

type Bitcoin struct {
	Wallet     blockchain.BitcoinWallet
	Connection *bitcoin.Connection
}

func NewBitcoinRegistry(env environment.BtcEnv, secrets environment.ApplicationSecrets, connection *bitcoin.Connection) (*Bitcoin, error) {
	wallet := bitcoin.NewBitcoindWallet(connection, env.BtcAddress, env.FixedTxFeeRate, env.WalletEncrypted, secrets.BtcWalletPassword)
	if err := wallet.Unlock(); err != nil {
		return nil, err
	}
	return &Bitcoin{
		Wallet:     wallet,
		Connection: connection,
	}, nil
}
