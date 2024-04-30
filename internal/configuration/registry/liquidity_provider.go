package registry

import "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"

func NewLiquidityProvider(
	databaseRegistry *Database,
	rskRegistry *Rootstock,
	btcRegistry *Bitcoin,
	messaging *Messaging,
) *dataproviders.LocalLiquidityProvider {
	return dataproviders.NewLocalLiquidityProvider(
		databaseRegistry.PeginRepository,
		databaseRegistry.PegoutRepository,
		databaseRegistry.LiquidityProviderRepository,
		messaging.Rpc,
		rskRegistry.Wallet,
		btcRegistry.PaymentWallet,
		rskRegistry.Contracts,
	)
}
