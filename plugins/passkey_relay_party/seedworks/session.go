package seedworks

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

func GetSessionKey(reg Registration) string {
	return reg.Origin + ":" + reg.Email
}

type SessionStore struct {
	sessions map[string]*sessionCache
	locker   sync.RWMutex
}

func NewInMemorySessionStore() *SessionStore {
	store := &SessionStore{
		sessions: make(map[string]*sessionCache),
	}
	store.recycle()
	return store
}

func (store *SessionStore) NewRegSession(origin, email string) (*protocol.CredentialCreation, error) {
	user := newUser(email)
	webauthn, _ := newWebAuthn(origin)
	sessionKey := GetSessionKey(Registration{Origin: origin, Email: email})
	if opt, session, err := webauthn.BeginRegistration(user); err != nil {
		return nil, err
	} else {
		store.set(sessionKey, webauthn, session, user)
		return opt, nil
	}
}

func (store *SessionStore) NewAuthSession(origin, email string) (*protocol.CredentialAssertion, error) {
	user := newUser(email)
	webauthn, _ := newWebAuthn(origin)
	sessionKey := GetSessionKey(Registration{Origin: origin, Email: email})
	if opt, session, err := webauthn.BeginLogin(user); err != nil {
		return nil, err
	} else {
		store.set(sessionKey, webauthn, session, user)
		return opt, nil
	}
}

func (store *SessionStore) Get(id string) *sessionCache {
	store.locker.RLock()
	defer store.locker.RUnlock()
	if session, ok := store.sessions[id]; ok {
		if session.obsolete {
			delete(store.sessions, id)
			return nil
		}
		return session
	}
	return nil
}

func (store *SessionStore) set(key string, webauthn *webauthn.WebAuthn, session *webauthn.SessionData, user *User) {
	store.locker.Lock()
	defer store.locker.Unlock()

	cache := &sessionCache{
		Data:     *session,
		WebAuthn: *webauthn,
		User:     *user,
		expires:  120,
	}
	cache.countdown()
	store.sessions[key] = cache
}

func newWebAuthn(origin string) (*webauthn.WebAuthn, error) {
	u, err := url.Parse(origin)
	if err != nil {
		return nil, err
	}
	hostname := u.Hostname()
	wconfig := &webauthn.Config{
		RPDisplayName: origin,
		RPID:          hostname,                   // Generally the FQDN for your site
		RPOrigins:     []string{origin, hostname}, // The origin URLs allowed for WebAuthn requests
	}

	if webAuthn, err := webauthn.New(wconfig); err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		return webAuthn, nil
	}
}

type sessionCache struct {
	Data     webauthn.SessionData
	WebAuthn webauthn.WebAuthn
	User     User
	expires  int8
	obsolete bool
}

func (cache *sessionCache) countdown() {
	go func() {
		for {
			ch := time.After(time.Second)
			<-ch
			cache.expires--

			if cache.expires < 0 {
				cache.obsolete = true
				break
			}
		}
	}()
}

func (store *SessionStore) recycle() {
	go func() {
		for {
			ch := time.After(15 * time.Second)
			<-ch

			store.locker.Lock()
			for id, session := range store.sessions {
				if session.obsolete {
					delete(store.sessions, id)
				}
			}
			store.locker.Unlock()
		}
	}()
}
