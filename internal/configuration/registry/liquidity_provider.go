package registry

import "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"

func NewLiquidityProvider(
	config *dataproviders.Configuration,
	databaseRegistry *Database,
	rskRegistry *Rootstock,
	btcRegistry *Bitcoin,
) *dataproviders.LocalLiquidityProvider {
	return dataproviders.NewLocalLiquidityProvider(
		config,
		databaseRegistry.PeginRepository,
		databaseRegistry.PegoutRepository,
		rskRegistry.RpcServer,
		rskRegistry.Wallet,
		btcRegistry.Wallet,
		rskRegistry.Lbc,
	)
}
