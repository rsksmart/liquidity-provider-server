package blockchain

// Pause levels reported by IPauseRegistry. A soft pause still allows a send; a hard pause does not.
const (
	PauseLevelNone uint8 = 0
	PauseLevelSoft uint8 = 1
	PauseLevelHard uint8 = 2
)

// PauseRegistryContract is a read-only port over IPauseRegistry.
// Pause writes (setPauseLevel) are not exposed: the liquidity provider server only queries pause state.
type PauseRegistryContract interface {
	GetAddress() string
	PauseLevel() (uint8, error)
	PauseStatus() (PauseStatus, error)
}
