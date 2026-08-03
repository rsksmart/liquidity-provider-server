package blockchain

import "github.com/rsksmart/liquidity-provider-server/internal/entities"

// FlyoverConfigurationsContract is a read-only port over the frozen IFlyoverConfigurations ABI.
// Ticket AC exposes exactly CalculatePegInFee and GetRequiredPegInBtcConfirmations;
// getPegInConfiguration/queueChange/applyChange are out of scope for this port.
type FlyoverConfigurationsContract interface {
	GetAddress() string
	CalculatePegInFee(amount *entities.Wei) (*entities.Wei, error)
	GetRequiredPegInBtcConfirmations(amount *entities.Wei) (uint64, error)
}
