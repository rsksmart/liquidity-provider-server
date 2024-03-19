package registry

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	environment2 "github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

type Bitcoin struct {
	RpcServer  blockchain.BitcoinNetwork
	Wallet     blockchain.BitcoinWallet
	Connection *bitcoin.Connection
}

func NewBitcoinRegistry(env environment2.BtcEnv, secrets environment2.ApplicationSecrets, connection *bitcoin.Connection) (*Bitcoin, error) {
	wallet := bitcoin.NewBitcoindWallet(connection, env.FixedTxFeeRate, env.WalletEncrypted, secrets.BtcWalletPassword)
	if err := wallet.Unlock(); err != nil {
		return nil, err
	}
	return &Bitcoin{
		RpcServer:  bitcoin.NewBitcoindRpc(connection),
		Wallet:     wallet,
		Connection: connection,
	}, nil
}
