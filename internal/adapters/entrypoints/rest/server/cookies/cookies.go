package cookies

import (
	"fmt"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
)

const (
	ManagementSessionCookieName = "lp-session"
	KeysBytesLength             = 32
	SessionMaxSeconds           = 60 * 30
)

var (
	storeOnce   sync.Once
	cookieStore *AESGCMSessionStore
	storeErr    error
)

func resetSessionCookieStoreForTest() {
	storeOnce = sync.Once{}
	cookieStore = nil
	storeErr = nil
}

func GetSessionCookieStore(env environment.ManagementEnv) (*AESGCMSessionStore, error) {
	storeOnce.Do(func() {
		key, err := utils.DecodeKey(env.SessionEncryptionKey, KeysBytesLength)
		if err != nil {
			storeErr = fmt.Errorf("error decoding session encryption key: %w", err)
			return
		}
		cookieStore, storeErr = NewAESGCMSessionStore(ManagementSessionCookieName, key, SessionMaxSeconds, env.UseHttps)
	})
	return cookieStore, storeErr
}
