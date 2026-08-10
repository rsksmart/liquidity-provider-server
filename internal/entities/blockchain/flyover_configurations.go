package blockchain

import "github.com/rsksmart/liquidity-provider-server/internal/entities"

type FlyoverConfigurationsContract interface {
	GetAddress() string
	CalculatePegInFee(amount *entities.Wei) (*entities.Wei, error)
	GetRequiredPegInBtcConfirmations(amount *entities.Wei) (uint64, error)
}
