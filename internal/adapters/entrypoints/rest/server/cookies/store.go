package cookies

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"net/http"
	"sync"
)

// UniqueSessionStore is a custom implementation of the sessions.Store interface. The rationale to implement this is that
// existing implementations don't provide any way to prevent concurrent logins and the LPS management session should be unique
type UniqueSessionStore struct {
	sessions.CookieStore
	session      *sessions.Session
	name         string
	sessionMutex *sync.Mutex
}

func NewUniqueSessionStore(uniqueSessionName string, keyPairs ...[]byte) *UniqueSessionStore {
	store := &UniqueSessionStore{
		CookieStore:  *sessions.NewCookieStore(keyPairs...),
		name:         uniqueSessionName,
		sessionMutex: &sync.Mutex{},
	}
	return store
}

func (s *UniqueSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *UniqueSessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	if name != s.name {
		return s.dummySession(name), fmt.Errorf("UniqueSessionStore is expecting %s session name and received %s", s.name, name)
	}

	if s.session != nil {
		return s.getExistingSession(r, name)
	}

	session := sessions.NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true
	return session, nil
}

func (s *UniqueSessionStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	const idSize = 32
	s.sessionMutex.Lock()
	defer s.sessionMutex.Unlock()
	// Delete if max-age is <= 0
	if session.Options.MaxAge <= 0 {
		s.session = nil
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}

	if session.ID == "" {
		idBytes, err := utils.GetRandomBytes(idSize)
		if err != nil {
			return err
		}
		session.ID = hex.EncodeToString(idBytes)
	}

	s.session = session
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

func (s *UniqueSessionStore) CloseUniqueSession(r *http.Request, w http.ResponseWriter) error {
	if s.session == nil {
		return nil
	}
	s.session.Options.MaxAge = -1
	if err := s.session.Save(r, w); err != nil {
		return err
	}
	return nil
}

func (s *UniqueSessionStore) getExistingSession(r *http.Request, name string) (*sessions.Session, error) {
	var err error
	var cookie *http.Cookie
	var sessionId string

	cookie, err = r.Cookie(name)
	if err != nil {
		return s.dummySession(name), err
	}

	err = securecookie.DecodeMulti(name, cookie.Value, &sessionId, s.Codecs...)
	if err != nil {
		return s.dummySession(name), err
	}
	if sessionId == s.session.ID {
		s.session.IsNew = false
		return s.session, nil
	} else {
		return s.dummySession(name), errors.New("session not recognized")
	}
}

// dummySession some parts of the gorilla sessions library expect a session even if an error is returned from a function
// such is the case of New function of sessions.Store interface. In order to maintain compatibility with the API and avoid
// nil pointer errors, this function should be used only to return a dummy session together with an error
func (s *UniqueSessionStore) dummySession(name string) *sessions.Session {
	return sessions.NewSession(s, name)
}
