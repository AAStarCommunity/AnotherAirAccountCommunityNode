package seedwork

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

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

func (store *SessionStore) NewSession(origin, email string) (*protocol.CredentialCreation, error) {
	user := newUser(email)
	webauthn, _ := newWebAuthn(origin)
	if opt, session, err := webauthn.BeginRegistration(user); err != nil {
		return nil, err
	} else {
		store.set(string(session.UserID), webauthn, session, origin, user)
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
	for k := range store.sessions {
		return store.sessions[k]
	}
	return nil
}

func (store *SessionStore) set(id string, webauthn *webauthn.WebAuthn, session *webauthn.SessionData, origin string, user *User) {
	store.locker.Lock()
	defer store.locker.Unlock()

	cache := &sessionCache{
		Data:     *session,
		WebAuthn: *webauthn,
		User:     *user,
		expires:  120,
	}
	cache.countdown()
	store.sessions[id] = cache
}

// TODO: origin should includes hostname and protocol://host
func newWebAuthn(origin string) (*webauthn.WebAuthn, error) {
	wconfig := &webauthn.Config{
		RPDisplayName: origin,
		RPID:          origin,                                    // Generally the FQDN for your site
		RPOrigins:     []string{origin, "http://localhost:3000"}, // The origin URLs allowed for WebAuthn requests
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
