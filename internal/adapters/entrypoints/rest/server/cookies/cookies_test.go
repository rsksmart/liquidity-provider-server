package cookies

import (
	"encoding/hex"
	"sync"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSessionCookieStore(t *testing.T) {
	t.Cleanup(resetSessionCookieStoreForTest)

	t.Run("return error if encryption key is invalid", func(t *testing.T) {
		env := environment.ManagementEnv{
			SessionEncryptionKey: "invalid",
			UseHttps:             false,
		}
		_, err := GetSessionCookieStore(env)
		require.Error(t, err)
	})
	t.Run("always return the same store", func(t *testing.T) {
		resetSessionCookieStoreForTest()
		env := environment.ManagementEnv{
			SessionEncryptionKey: hex.EncodeToString(make([]byte, 32)),
			UseHttps:             false,
		}
		stores := make([]*UniqueSessionStore, 0, 10)
		errs := make([]error, 0, 10)
		wg := sync.WaitGroup{}
		mutex := sync.Mutex{}
		wg.Add(10)
		for i := 0; i < 10; i++ {
			go func() {
				defer wg.Done()
				store, err := GetSessionCookieStore(env)
				mutex.Lock()
				stores = append(stores, store)
				errs = append(errs, err)
				mutex.Unlock()
			}()
		}
		wg.Wait()
		for _, err := range errs {
			require.NoError(t, err)
		}
		for i := 1; i < 10; i++ {
			assert.Same(t, stores[0], stores[i])
		}
	})
}
