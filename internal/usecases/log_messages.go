package usecases

import "fmt"

func LogPauseLevelBlocks(registryAddress string, level uint8) string {
	return fmt.Sprintf("PauseRegistry %s is at pause level %d", registryAddress, level)
}
