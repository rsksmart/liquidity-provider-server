package dataproviders

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

type Configuration struct {
	RskConfig    RskConfig
	BtcConfig    BitcoinConfig
	PeginConfig  PeginConfig
	PegoutConfig PegoutConfig
}

type RskConfig struct {
	ChainId       uint64
	Account       uint64
	Confirmations map[int]uint16
}

type BitcoinConfig struct {
	BtcAddress    string
	Confirmations map[int]uint16
}

type PeginConfig struct {
	TimeForDeposit      uint32
	CallTime            uint32
	PenaltyFee          *entities.Wei
	CallFee             *entities.Wei
	MinTransactionValue *entities.Wei
	MaxTransactionValue *entities.Wei
}

type PegoutConfig struct {
	TimeForDeposit      uint32
	CallTime            uint32
	PenaltyFee          *entities.Wei
	CallFee             *entities.Wei
	MinTransactionValue *entities.Wei
	MaxTransactionValue *entities.Wei
	ExpireBlocks        uint32
}
