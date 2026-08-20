package blockchain

// PauseRegistryContract is a read-only port over IPauseRegistry.
// Pause writes (setPauseLevel) are not exposed: the liquidity provider server only queries pause state.
type PauseRegistryContract interface {
	GetAddress() string
	PauseLevel() (uint8, error)
	PauseStatus() (PauseStatus, error)
}
