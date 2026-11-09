package blockchain

type PauseLevel uint8

const (
	PauseLevelNone PauseLevel = iota
	PauseLevelSoft
	PauseLevelHard
)

type PauseRegistryContract interface {
	PauseLevel() (PauseLevel, error)
}
