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
	cookieStore *UniqueSessionStore
	storeErr    error
)

func GetSessionCookieStore(env environment.ManagementEnv) (*UniqueSessionStore, error) {
	key, err := utils.DecodeKey(env.SessionEncryptionKey, KeysBytesLength)
	if err != nil {
		return nil, fmt.Errorf("error decoding session encryption key: %w", err)
	}
	storeOnce.Do(func() {
		cookieStore, storeErr = NewUniqueSessionStore(ManagementSessionCookieName, key, SessionMaxSeconds, env.UseHttps)
	})
	return cookieStore, storeErr
}
