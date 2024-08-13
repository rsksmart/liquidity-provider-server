package cookies_test

import (
	"encoding/hex"
	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestSession(t *testing.T) {
	testGetSessionCookieStore(t)
	cookie := testCreateManagementSession(t)
	testCloseManagementSession(t, cookie)
}

func testGetSessionCookieStore(t *testing.T) {
	t.Run("return error if auth key is invalid", func(t *testing.T) {
		env := environment.Environment{Management: environment.ManagementEnv{
			EnableManagementApi:  false,
			SessionAuthKey:       "invalid",
			SessionEncryptionKey: hex.EncodeToString(make([]byte, 32)),
			SessionTokenAuthKey:  hex.EncodeToString(make([]byte, 32)),
			UseHttps:             false,
		}}
		_, err := cookies.GetSessionCookieStore(env.Management)
		require.Error(t, err)
	})
	t.Run("return error if encryption key is invalid", func(t *testing.T) {
		env := environment.Environment{Management: environment.ManagementEnv{
			EnableManagementApi:  false,
			SessionAuthKey:       hex.EncodeToString(make([]byte, 32)),
			SessionEncryptionKey: "invalid",
			SessionTokenAuthKey:  hex.EncodeToString(make([]byte, 32)),
			UseHttps:             false,
		}}
		_, err := cookies.GetSessionCookieStore(env.Management)
		require.Error(t, err)
	})
	t.Run("always return the same store", func(t *testing.T) {
		env := environment.Environment{Management: environment.ManagementEnv{
			EnableManagementApi:  false,
			SessionAuthKey:       hex.EncodeToString(make([]byte, 32)),
			SessionEncryptionKey: hex.EncodeToString(make([]byte, 32)),
			SessionTokenAuthKey:  hex.EncodeToString(make([]byte, 32)),
			UseHttps:             false,
		}}
		stores := make([]sessions.Store, 0)
		wg := sync.WaitGroup{}
		mutex := sync.Mutex{}
		wg.Add(10)
		for i := 0; i < 10; i++ {
			func() {
				defer wg.Done()
				defer mutex.Unlock()
				mutex.Lock()
				store, err := cookies.GetSessionCookieStore(env.Management)
				require.NoError(t, err)
				stores = append(stores, store)
			}()
		}
		wg.Wait()
		for i := 1; i < 10; i++ {
			assert.Same(t, stores[0], stores[i])
		}
	})
}

func testCreateManagementSession(t *testing.T) *http.Cookie {
	var cookie *http.Cookie
	t.Run("should create a new session", func(t *testing.T) {
		env := environment.Environment{Management: environment.ManagementEnv{}}
		store, err := cookies.GetSessionCookieStore(env.Management)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		args := &cookies.CreateSessionArgs{
			Store:   store,
			Env:     environment.ManagementEnv{},
			Request: req,
			Writer:  response,
		}
		err = cookies.CreateManagementSession(args)
		require.NoError(t, err)
		session, err := store.Get(req, cookies.ManagementSessionCookieName)
		require.NoError(t, err)
		assert.Equal(t, cookies.SessionMaxSeconds, session.Options.MaxAge)
		assert.Equal(t, "/", session.Options.Path)
		assert.True(t, session.Options.HttpOnly)
		assert.Equal(t, http.SameSiteStrictMode, session.Options.SameSite)
		assert.Equal(t, req.URL.Host, session.Options.Domain)
		assert.Equal(t, env.Management.UseHttps, session.Options.Secure)
		// nolint:bodyclose
		cookie = response.Result().Cookies()[0]
		require.NoError(t, response.Result().Body.Close())
	})

	t.Run("should return error if session cannot be saved", func(t *testing.T) {
		env := environment.Environment{Management: environment.ManagementEnv{}}
		store, err := cookies.GetSessionCookieStore(env.Management)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		args := &cookies.CreateSessionArgs{
			Store:   store,
			Env:     environment.ManagementEnv{},
			Request: req,
			Writer:  nil,
		}
		err = cookies.CreateManagementSession(args)
		require.Error(t, err)
	})
	return cookie
}

func testCloseManagementSession(t *testing.T, cookie *http.Cookie) {
	t.Run("should close the session", func(t *testing.T) {
		// validate session is open
		env := environment.Environment{Management: environment.ManagementEnv{}}
		store, err := cookies.GetSessionCookieStore(env.Management)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookie)
		session, err := store.Get(req, cookies.ManagementSessionCookieName)
		require.NoError(t, err)
		assert.False(t, session.IsNew)

		// close the session
		res := httptest.NewRecorder()
		args := &cookies.CloseSessionArgs{
			Store:   store,
			Request: req,
			Writer:  res,
		}
		require.NoError(t, cookies.CloseManagementSession(args))

		// validate session is closed
		req = httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookie)
		session, err = store.Get(req, cookies.ManagementSessionCookieName)
		require.NoError(t, err)
		assert.True(t, session.IsNew)
	})
	t.Run("should return error if the store doesn't have the correct type", func(t *testing.T) {
		err := cookies.CloseManagementSession(&cookies.CloseSessionArgs{
			Store: sessions.NewCookieStore(),
		})
		require.ErrorContains(t, err, "closing a unique session is only supported by UniqueSessionStore")
	})
}
