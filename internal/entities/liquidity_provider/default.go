package liquidity_provider

import "github.com/rsksmart/liquidity-provider-server/internal/entities"

const (
	PeginTimeForDeposit = 3600
	PeginCallTime       = 7200
	PeginPenaltyFee     = 10000000000000
	PeginCallFee        = 100000000000000
	PeginMinValue       = 5000000000000000
	PeginMaxValue       = 100000000000000000
)

const (
	PegoutTimeForDeposit       = 3600
	PegoutCallTime             = 7200
	PegoutPenaltyFee           = 10000000000000
	PegoutCallFee              = 100000000000000
	PegoutMinValue             = 5000000000000000
	PegoutMaxValue             = 100000000000000000
	PegoutExpireBlocks         = 500
	PegoutBridgeTransactionMin = 15000000000000000
)

func DefaultRskConfirmationsPerAmount() ConfirmationsPerAmount {
	return ConfirmationsPerAmount{
		100000000000000000:  40,
		400000000000000000:  120,
		2000000000000000000: 200,
		4000000000000000000: 400,
		8000000000000000000: 800,
	}
}

func DefaultBtcConfirmationsPerAmount() ConfirmationsPerAmount {
	return ConfirmationsPerAmount{
		100000000000000000:  2,
		400000000000000000:  6,
		2000000000000000000: 10,
		4000000000000000000: 20,
		8000000000000000000: 40,
	}
}

func DefaultPeginConfiguration() PeginConfiguration {
	return PeginConfiguration{
		TimeForDeposit: PeginTimeForDeposit,
		CallTime:       PeginCallTime,
		PenaltyFee:     entities.NewWei(PeginPenaltyFee),
		CallFee:        entities.NewWei(PeginCallFee),
		MaxValue:       entities.NewWei(PeginMaxValue),
		MinValue:       entities.NewWei(PeginMinValue),
	}
}

func DefaultPegoutConfiguration() PegoutConfiguration {
	return PegoutConfiguration{
		TimeForDeposit:       PegoutTimeForDeposit,
		CallTime:             PegoutCallTime,
		PenaltyFee:           entities.NewWei(PegoutPenaltyFee),
		CallFee:              entities.NewWei(PegoutCallFee),
		MaxValue:             entities.NewWei(PegoutMaxValue),
		MinValue:             entities.NewWei(PegoutMinValue),
		ExpireBlocks:         PegoutExpireBlocks,
		BridgeTransactionMin: entities.NewWei(PegoutBridgeTransactionMin),
	}
}

func DefaultGeneralConfiguration() GeneralConfiguration {
	return GeneralConfiguration{
		RskConfirmations:     DefaultRskConfirmationsPerAmount(),
		BtcConfirmations:     DefaultBtcConfirmationsPerAmount(),
		PublicLiquidityCheck: false,
	}
}
