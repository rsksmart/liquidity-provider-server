package cookies

import (
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"net/http"
	"sync"
)

const (
	ManagementSessionCookieName = "lp-session"
	CsrfCookieName              = "lps-csrf-cookie"
	KeysBytesLength             = 32
	SessionMaxSeconds           = 60 * 30
)

var storeOnce sync.Once
var cookieStore sessions.Store

func GetSessionCookieStore(env environment.ManagementEnv) (sessions.Store, error) {
	var authKey, encryptionKey []byte
	var err error
	if cookieStore != nil {
		return cookieStore, nil
	}

	authKey, err = utils.DecodeKey(env.SessionAuthKey, KeysBytesLength)
	if err != nil {
		err = fmt.Errorf("error decoding session auth key: %w", err)
		return nil, err
	}
	encryptionKey, err = utils.DecodeKey(env.SessionEncryptionKey, KeysBytesLength)
	if err != nil {
		err = fmt.Errorf("error decoding session encryption key: %w", err)
		return nil, err
	}

	storeOnce.Do(func() {
		cookieStore = NewUniqueSessionStore(ManagementSessionCookieName, authKey, encryptionKey)
	})
	return cookieStore, err
}

type CreateSessionArgs struct {
	Store   sessions.Store
	Env     environment.ManagementEnv
	Request *http.Request
	Writer  http.ResponseWriter
}

func CreateManagementSession(args *CreateSessionArgs) error {
	session, err := args.Store.Get(args.Request, ManagementSessionCookieName)
	if err != nil {
		return err
	}
	session.Options.Domain = args.Request.URL.Host
	session.Options.MaxAge = SessionMaxSeconds
	session.Options.Path = "/"
	session.Options.Secure = args.Env.UseHttps
	session.Options.HttpOnly = true
	session.Options.SameSite = http.SameSiteStrictMode
	return session.Save(args.Request, args.Writer)
}

type CloseSessionArgs struct {
	Store   sessions.Store
	Request *http.Request
	Writer  http.ResponseWriter
}

func CloseManagementSession(args *CloseSessionArgs) error {
	store, ok := args.Store.(*UniqueSessionStore)
	if !ok {
		return errors.New("closing a unique session is only supported by UniqueSessionStore")
	}
	return store.CloseUniqueSession(args.Request, args.Writer)
}
