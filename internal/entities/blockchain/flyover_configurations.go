package blockchain

import "github.com/rsksmart/liquidity-provider-server/internal/entities"

// ConfirmationTier is one FlyoverConfigurations amount/confirmations bucket.
type ConfirmationTier struct {
	MaxAmount     *entities.Wei
	Confirmations uint64
}

// PegOutConfiguration is the snapshotted peg-out parameter set from FlyoverConfigurations.
type PegOutConfiguration struct {
	FixedFee          *entities.Wei
	PercentageFee     uint64
	MinAmount         *entities.Wei
	MaxAmount         *entities.Wei
	ConfirmationTiers []ConfirmationTier
	PenaltyFee        *entities.Wei
	ClaimWindow       uint64
	ClaimWindowBlocks uint64
	CallTime          uint64
	ExpireTime        uint64
	ExpireBlocks      uint64
	MaxMinerFee       *entities.Wei
}

type FlyoverConfigurationsContract interface {
	GetAddress() string
	CalculatePegInFee(amount *entities.Wei) (*entities.Wei, error)
	GetRequiredPegInBtcConfirmations(amount *entities.Wei) (uint64, error)
	CalculatePegOutFee(amount *entities.Wei) (*entities.Wei, error)
	GetRequiredPegOutBtcConfirmations(amount *entities.Wei) (uint64, error)
	GetPegOutConfiguration() (PegOutConfiguration, error)
}
