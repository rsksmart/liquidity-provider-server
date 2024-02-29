package environment

import "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"

func ConfigurationFromEnv(env Environment) *dataproviders.Configuration {
	return &dataproviders.Configuration{
		RskConfig: dataproviders.RskConfig{
			ChainId: env.Rsk.ChainId,
			Account: uint64(env.Rsk.AccountNumber),
		},
		BtcConfig: dataproviders.BitcoinConfig{
			BtcAddress: env.Provider.BtcAddress,
		},
		PeginConfig:  dataproviders.PeginConfig{},
		PegoutConfig: dataproviders.PegoutConfig{},
	}
}
