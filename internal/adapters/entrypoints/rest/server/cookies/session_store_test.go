package cookies_test

import (
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	aesKeyString  = "01fbac02d66202e8468d2a4f1deba4fa5c2491f592e0e22e32fe1e6acac25923"
	aesCookieName = "lp-session"
)

func TestNewAESGCMSessionStore_KeyLength(t *testing.T) {
	t.Run("rejects 16-byte key", func(t *testing.T) {
		_, err := cookies.NewAESGCMSessionStore(aesCookieName, make([]byte, 16), cookies.SessionMaxSeconds, false)
		require.Error(t, err)
	})
	t.Run("rejects 24-byte key", func(t *testing.T) {
		_, err := cookies.NewAESGCMSessionStore(aesCookieName, make([]byte, 24), cookies.SessionMaxSeconds, false)
		require.Error(t, err)
	})
	t.Run("accepts 32-byte key", func(t *testing.T) {
		key, err := hex.DecodeString(aesKeyString)
		require.NoError(t, err)
		store, err := cookies.NewAESGCMSessionStore(aesCookieName, key, cookies.SessionMaxSeconds, false)
		require.NoError(t, err)
		assert.NotNil(t, store)
	})
}

func TestAESGCMSessionStore_RoundTrip(t *testing.T) {
	store := newTestStore(t)
	cookie := createSession(t, store)

	err := store.Validate(reqWithCookie(cookie))
	require.NoError(t, err)
}

func TestAESGCMSessionStore_TamperRejection(t *testing.T) {
	store := newTestStore(t)
	cookie := createSession(t, store)

	t.Run("flipped byte in GCM tag", func(t *testing.T) {
		tampered := tamperCookieValue(t, cookie.Value, len(cookie.Value)-1)
		req := reqWithCookie(&http.Cookie{Name: cookie.Name, Value: tampered})
		require.ErrorIs(t, store.Validate(req), cookies.ErrSessionNotRecognized)
	})

	t.Run("flipped byte inside ciphertext", func(t *testing.T) {
		tampered := tamperCookieValue(t, cookie.Value, 20)
		req := reqWithCookie(&http.Cookie{Name: cookie.Name, Value: tampered})
		require.ErrorIs(t, store.Validate(req), cookies.ErrSessionNotRecognized)
	})

	t.Run("Refresh rejects tampered cookie", func(t *testing.T) {
		tampered := tamperCookieValue(t, cookie.Value, 10)
		req := reqWithCookie(&http.Cookie{Name: cookie.Name, Value: tampered})
		rec := httptest.NewRecorder()
		require.ErrorIs(t, store.Refresh(rec, req), cookies.ErrSessionNotRecognized)
		assert.Empty(t, rec.Result().Cookies())
	})
}

func TestAESGCMSessionStore_SingleSession(t *testing.T) {
	store := newTestStore(t)
	cookieA := createSession(t, store)
	cookieB := createSession(t, store)

	require.ErrorIs(t, store.Validate(reqWithCookie(cookieA)), cookies.ErrSessionNotRecognized)
	require.NoError(t, store.Validate(reqWithCookie(cookieB)))

	t.Run("Refresh with stale cookie does not issue Set-Cookie", func(t *testing.T) {
		rec := httptest.NewRecorder()
		err := store.Refresh(rec, reqWithCookie(cookieA))
		require.ErrorIs(t, err, cookies.ErrSessionNotRecognized)
		assert.Empty(t, rec.Result().Cookies())
	})

	t.Run("Close with stale cookie does not clear active session", func(t *testing.T) {
		rec := httptest.NewRecorder()
		require.NoError(t, store.Close(rec, reqWithCookie(cookieA)))
		require.NoError(t, store.Validate(reqWithCookie(cookieB)))
	})
}

func TestAESGCMSessionStore_NoCookie(t *testing.T) {
	store := newTestStore(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	require.ErrorIs(t, store.Validate(req), cookies.ErrSessionNotRecognized)
}

func TestAESGCMSessionStore_CloseAndRefresh(t *testing.T) {
	store := newTestStore(t)
	cookie := createSession(t, store)
	req := reqWithCookie(cookie)

	t.Run("Refresh re-issues a valid cookie", func(t *testing.T) {
		rec := httptest.NewRecorder()
		require.NoError(t, store.Refresh(rec, req))
		refreshed := rec.Result().Cookies()[0]
		require.NoError(t, store.Validate(reqWithCookie(refreshed)))
	})

	t.Run("Close clears active session", func(t *testing.T) {
		rec := httptest.NewRecorder()
		require.NoError(t, store.Close(rec, req))
		require.ErrorIs(t, store.Validate(req), cookies.ErrSessionNotRecognized)
	})
}

func newTestStore(t *testing.T) *cookies.AESGCMSessionStore {
	t.Helper()
	key, err := hex.DecodeString(aesKeyString)
	require.NoError(t, err)
	store, err := cookies.NewAESGCMSessionStore(aesCookieName, key, cookies.SessionMaxSeconds, false)
	require.NoError(t, err)
	return store
}

func createSession(t *testing.T, store *cookies.AESGCMSessionStore) *http.Cookie {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	require.NoError(t, store.Create(rec, req))
	result := rec.Result()
	require.NoError(t, result.Body.Close())
	cookie := result.Cookies()[0]
	assert.Equal(t, cookies.SessionMaxSeconds, cookie.MaxAge)
	assert.Equal(t, "/", cookie.Path)
	assert.True(t, cookie.HttpOnly)
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
	assert.Empty(t, cookie.Domain, "host-only cookie scope: Domain must be omitted")
	return cookie
}

func reqWithCookie(cookie *http.Cookie) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)
	return req
}

func tamperCookieValue(t *testing.T, value string, byteIndex int) string {
	t.Helper()
	raw, err := base64.RawURLEncoding.DecodeString(value)
	require.NoError(t, err)
	if byteIndex >= len(raw) {
		byteIndex = len(raw) - 1
	}
	raw[byteIndex] ^= 0xff
	return base64.RawURLEncoding.EncodeToString(raw)
}
