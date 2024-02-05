package environment

import "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"

func ConfigurationFromEnv(env Environment) *dataproviders.Configuration {
	return &dataproviders.Configuration{
		RskConfig: dataproviders.RskConfig{
			ChainId:       env.Rsk.ChainId,
			Account:       uint64(env.Rsk.AccountNumber),
			Confirmations: env.Rsk.Confirmations,
		},
		BtcConfig: dataproviders.BitcoinConfig{
			BtcAddress:    env.Provider.BtcAddress,
			Confirmations: env.Btc.Confirmations,
		},
		PeginConfig: dataproviders.PeginConfig{
			TimeForDeposit:      env.Pegin.TimeForDeposit,
			CallTime:            env.Pegin.CallTime,
			PenaltyFee:          env.Pegin.PenaltyFee,
			CallFee:             env.Pegin.CallFee,
			MinTransactionValue: env.Pegin.MinTransactionValue,
			MaxTransactionValue: env.Pegin.MaxTransactionValue,
		},
		PegoutConfig: dataproviders.PegoutConfig{
			TimeForDeposit:      env.Pegout.TimeForDeposit,
			CallTime:            env.Pegout.CallTime,
			PenaltyFee:          env.Pegout.PenaltyFee,
			CallFee:             env.Pegout.CallFee,
			MinTransactionValue: env.Pegout.MinTransactionValue,
			MaxTransactionValue: env.Pegout.MaxTransactionValue,
			ExpireBlocks:        env.Pegout.ExpireBlocks,
		},
	}
}
