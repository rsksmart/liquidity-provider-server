package cookies_test

import (
	"encoding/hex"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	key1String = "01fbac02d66202e8468d2a4f1deba4fa5c2491f592e0e22e32fe1e6acac25923"
	key2String = "02fbac02d66202e8468d2a4f1deba4fa5c2491f592e0e22e32fe1e6acac25923"
	cookieName = "cookie"
)

func TestUniqueSessionStore_New(t *testing.T) {
	var (
		cookie         *http.Cookie
		firstSessionId string
	)
	k1, err := hex.DecodeString(key1String)
	require.NoError(t, err)
	k2, err := hex.DecodeString(key2String)
	require.NoError(t, err)
	store := cookies.NewUniqueSessionStore(cookieName, k1, k2)
	t.Run("should return an error if the session name is different", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		session, err := store.New(req, test.AnyString)
		assertDummySession(t, session)
		require.ErrorContains(t, err, "UniqueSessionStore is expecting cookie session name and received any value")
	})
	t.Run("should return an new session the first time", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		session, err := store.New(req, cookieName)
		require.NoError(t, err)
		assert.NotEmpty(t, session)
		assert.True(t, session.IsNew)

		// get cookie for next test
		err = store.Save(req, res, session)
		require.NoError(t, err)
		// nolint:bodyclose
		cookie = res.Result().Cookies()[0]
		require.NoError(t, res.Result().Body.Close())
		firstSessionId = session.ID
	})
	t.Run("should return an existing session the second time", func(t *testing.T) {
		t.Run("should handle error decoding cookie", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.AddCookie(&http.Cookie{Name: cookieName, Value: "-"})
			session, err := store.New(req, cookieName)
			assertDummySession(t, session)
			require.Error(t, err)
		})
		t.Run("should return error when session not recognized", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			otherId, err := utils.GetRandomBytes(32)
			require.NoError(t, err)
			otherCookie, err := securecookie.EncodeMulti(cookieName, hex.EncodeToString(otherId), securecookie.CodecsFromPairs(k1, k2)...)
			require.NoError(t, err)
			req.AddCookie(&http.Cookie{Name: cookieName, Value: otherCookie})
			session, err := store.New(req, cookieName)
			assertDummySession(t, session)
			require.ErrorContains(t, err, "session not recognized")
		})
		t.Run("should return existing session successfully", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.AddCookie(cookie)
			session, err := store.New(req, cookieName)
			require.NoError(t, err)
			assert.NotNil(t, session)
			assert.Equal(t, firstSessionId, session.ID)
			assert.False(t, session.IsNew)
		})
	})
}

func TestUniqueSessionStore_Get(t *testing.T) {
	k1, err := hex.DecodeString(key1String)
	require.NoError(t, err)
	k2, err := hex.DecodeString(key2String)
	require.NoError(t, err)
	store := cookies.NewUniqueSessionStore(cookieName, k1, k2)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	newSession, err := store.New(req, cookieName)
	require.NoError(t, err)
	existingSession, err := store.Get(req, cookieName)
	require.NoError(t, err)
	assert.Equal(t, newSession.ID, existingSession.ID)
}

func TestUniqueSessionStore_Save(t *testing.T) {
	var session *sessions.Session
	var err error
	k1, err := hex.DecodeString(key1String)
	require.NoError(t, err)
	k2, err := hex.DecodeString(key2String)
	require.NoError(t, err)
	store := cookies.NewUniqueSessionStore(cookieName, k1, k2)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	session, err = store.New(req, cookieName)
	require.NoError(t, err)
	assert.NotEmpty(t, session)
	assert.True(t, session.IsNew)

	t.Run("should save session", func(t *testing.T) {
		err = store.Save(req, res, session)
		require.NoError(t, err)
		// nolint:bodyclose
		req.AddCookie(res.Result().Cookies()[0])
		require.NoError(t, res.Result().Body.Close())

		session, err = store.Get(req, cookieName)
		require.NoError(t, err)
		assert.False(t, session.IsNew)
		assert.NotEmpty(t, session.ID)
	})
	t.Run("should remove session if max age is less than or equal to 0", func(t *testing.T) {
		res = httptest.NewRecorder()
		session, err = store.Get(req, cookieName)
		session.Options.MaxAge = -1
		err = store.Save(req, res, session)
		require.NoError(t, err)
		// nolint:bodyclose
		assert.Empty(t, res.Result().Cookies()[0].Value)
		require.NoError(t, res.Result().Body.Close())
	})
}

func TestUniqueSessionStore_CloseUniqueSession(t *testing.T) {
	k1, err := hex.DecodeString(key1String)
	require.NoError(t, err)
	k2, err := hex.DecodeString(key2String)
	require.NoError(t, err)
	store := cookies.NewUniqueSessionStore(cookieName, k1, k2)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	session, err := store.New(req, cookieName)
	require.NoError(t, err)
	assert.True(t, session.IsNew)
	err = store.Save(req, res, session)
	require.NoError(t, err)

	// nolint:bodyclose
	req.AddCookie(res.Result().Cookies()[0])
	require.NoError(t, res.Result().Body.Close())
	session, err = store.Get(req, cookieName)
	require.NoError(t, err)
	require.False(t, session.IsNew)
	require.NotEmpty(t, session.ID)

	t.Run("should close the session", func(t *testing.T) {
		req = httptest.NewRequest(http.MethodGet, "/", nil)
		// nolint:bodyclose
		req.AddCookie(res.Result().Cookies()[0])
		require.NoError(t, res.Result().Body.Close())
		res = httptest.NewRecorder()
		err = store.CloseUniqueSession(req, res)
		require.NoError(t, err)

		session, err = store.Get(req, cookieName)
		require.NoError(t, err)
		assert.Empty(t, session.ID)
	})
}

func assertDummySession(t *testing.T, session *sessions.Session) {
	assert.NotNil(t, session)
	assert.Empty(t, session.Options)
	assert.False(t, session.IsNew)
	assert.Empty(t, session.Values)
	assert.Empty(t, session.ID)
}
