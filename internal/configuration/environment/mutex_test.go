package environment_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestNewApplicationMutexes(t *testing.T) {
	mutexes := environment.NewApplicationMutexes()
	t.Run("should create every mutex", func(t *testing.T) {
		assert.NotNil(t, mutexes.RskWalletMutex())
		assert.NotNil(t, mutexes.BtcWalletMutex())
		assert.NotNil(t, mutexes.PeginLiquidityMutex())
		assert.NotNil(t, mutexes.PegoutLiquidityMutex())
	})
	t.Run("every mutex should be different", func(t *testing.T) {
		mutexArray := []*sync.Mutex{
			mutexes.RskWalletMutex(),
			mutexes.BtcWalletMutex(),
			mutexes.PeginLiquidityMutex(),
			mutexes.PegoutLiquidityMutex(),
		}
		seen := make(map[*sync.Mutex]bool)
		for _, element := range mutexArray {
			if seen[element] {
				assert.Fail(t, "Found duplicate element in mutexes array")
			}
			seen[element] = true
		}
	})
}
