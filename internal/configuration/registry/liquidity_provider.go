package registry

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/wallet"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/cold_wallet"
)

type LiquidityProvider struct {
	LiquidityProvider *dataproviders.LocalLiquidityProvider
	ColdWallet        cold_wallet.ColdWallet
}

func NewLiquidityProviderRegistry(
	databaseRegistry *Database,
	rskRegistry *Rootstock,
	btcRegistry *Bitcoin,
	messaging *Messaging,
	walletFactory wallet.AbstractFactory,
) (*LiquidityProvider, error) {
	coldWallet, err := walletFactory.ColdWallet(messaging.Rpc)
	if err != nil {
		return nil, err
	}
	return &LiquidityProvider{
		LiquidityProvider: dataproviders.NewLocalLiquidityProvider(
			databaseRegistry.PeginRepository,
			databaseRegistry.PegoutRepository,
			databaseRegistry.LiquidityProviderRepository,
			messaging.Rpc,
			rskRegistry.Wallet,
			btcRegistry.PaymentWallet,
			rskRegistry.Contracts,
		),
		ColdWallet: coldWallet,
	}, nil
}
