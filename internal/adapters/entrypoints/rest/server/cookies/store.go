package cookies

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
)

const sessionIDBytes = 32

var ErrSessionNotRecognized = errors.New("session not recognized")

// SessionStore manages the single management login session.
type SessionStore interface {
	Create(w http.ResponseWriter, r *http.Request) error
	Validate(r *http.Request) error
	Refresh(w http.ResponseWriter, r *http.Request) error
	Close(w http.ResponseWriter, r *http.Request) error
}

// UniqueSessionStore keeps one login alive at a time: the active session ID lives in
// memory behind a mutex (this is what blocks concurrent logins) and is sealed into a
// cookie with AES-256-GCM. Creating a new session overwrites the active ID, invalidating
// any previously issued cookie.
type UniqueSessionStore struct {
	name     string
	gcm      cipher.AEAD
	maxAge   int
	secure   bool
	mu       sync.Mutex
	activeID []byte
}

func NewUniqueSessionStore(name string, key []byte, maxAge int, secure bool) (*UniqueSessionStore, error) {
	if len(key) != KeysBytesLength {
		return nil, fmt.Errorf("session key must be %d bytes for AES-256, got %d", KeysBytesLength, len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating GCM: %w", err)
	}
	return &UniqueSessionStore{name: name, gcm: gcm, maxAge: maxAge, secure: secure}, nil
}

func (s *UniqueSessionStore) Create(w http.ResponseWriter, r *http.Request) error {
	id, err := utils.GetRandomBytes(sessionIDBytes)
	if err != nil {
		return err
	}
	value, err := s.seal(id)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.activeID = id
	s.mu.Unlock()
	s.setCookie(w, value, s.maxAge)
	return nil
}

// Validate is read-only (no Set-Cookie): used by the UI handler to compute loggedIn.
func (s *UniqueSessionStore) Validate(r *http.Request) error {
	id, err := s.requestID(r)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.activeID == nil || subtle.ConstantTimeCompare(id, s.activeID) != 1 {
		return ErrSessionNotRecognized
	}
	return nil
}

// Refresh validates the REQUEST's cookie against the active ID and, only if it still
// matches, re-seals THAT SAME id into a fresh cookie (sliding window).
func (s *UniqueSessionStore) Refresh(w http.ResponseWriter, r *http.Request) error {
	id, err := s.requestID(r)
	if err != nil {
		return err
	}
	s.mu.Lock()
	active := s.activeID != nil && subtle.ConstantTimeCompare(id, s.activeID) == 1
	s.mu.Unlock()
	if !active {
		return ErrSessionNotRecognized
	}
	value, err := s.seal(id)
	if err != nil {
		return err
	}
	s.setCookie(w, value, s.maxAge)
	return nil
}

// Close always expires the requester's own cookie, but only clears the in-memory active ID
// if the request's cookie is that same session.
func (s *UniqueSessionStore) Close(w http.ResponseWriter, r *http.Request) error {
	s.setCookie(w, "", -1)
	if id := s.openRequestID(r); id != nil {
		s.mu.Lock()
		if s.activeID != nil && subtle.ConstantTimeCompare(id, s.activeID) == 1 {
			s.activeID = nil
		}
		s.mu.Unlock()
	}
	return nil
}

func (s *UniqueSessionStore) openRequestID(r *http.Request) []byte {
	id, err := s.requestID(r)
	if err != nil {
		return nil
	}
	return id
}

func (s *UniqueSessionStore) requestID(r *http.Request) ([]byte, error) {
	cookie, err := r.Cookie(s.name)
	if err != nil {
		return nil, ErrSessionNotRecognized
	}
	id, err := s.open(cookie.Value)
	if err != nil {
		return nil, ErrSessionNotRecognized
	}
	return id, nil
}

func (s *UniqueSessionStore) seal(id []byte) (string, error) {
	nonce, err := utils.GetRandomBytes(int64(s.gcm.NonceSize()))
	if err != nil {
		return "", err
	}
	sealed := s.gcm.Seal(nonce, nonce, id, nil)
	return base64.RawURLEncoding.EncodeToString(sealed), nil
}

func (s *UniqueSessionStore) open(value string) ([]byte, error) {
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	nonceSize := s.gcm.NonceSize()
	if len(raw) < nonceSize {
		return nil, errors.New("sealed cookie too short")
	}
	return s.gcm.Open(nil, raw[:nonceSize], raw[nonceSize:], nil)
}

func (s *UniqueSessionStore) setCookie(w http.ResponseWriter, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		Secure:   s.secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
