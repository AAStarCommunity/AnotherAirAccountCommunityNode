package seedworks

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

func GetSessionKey(reg *Registration) string {
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

func (store *SessionStore) NewRegSession(reg *Registration) (*protocol.CredentialCreation, error) {
	user := newUser(reg.Email)
	wan, _ := newWebAuthn(reg.Origin)
	sessionKey := GetSessionKey(reg)

	authSelect := protocol.AuthenticatorSelection{
		AuthenticatorAttachment: protocol.Platform,
		RequireResidentKey:      protocol.ResidentKeyNotRequired(),
		UserVerification:        protocol.VerificationRequired,
	}

	if opt, session, err := wan.BeginRegistration(user, webauthn.WithAuthenticatorSelection(authSelect)); err != nil {
		return nil, err
	} else {
		store.set(sessionKey, wan, session, user)
		return opt, nil
	}
}

func (store *SessionStore) FinishRegSession(reg *Registration, ctx *gin.Context) (*User, error) {
	key := GetSessionKey(reg)
	if session := store.Get(key); session == nil {
		return nil, fmt.Errorf("%s: not found", reg.Email)
	} else {
		if cred, err := session.WebAuthn.FinishRegistration(&session.User, session.Data, ctx.Request); err == nil {
			session.User.AddCredential(cred)
			store.remove(key)
			return &session.User, nil
		} else {
			return nil, err
		}
	}
}

func (store *SessionStore) NewAuthSession(user *User, signIn *SiginIn) (*protocol.CredentialAssertion, error) {
	if user == nil || signIn == nil {
		return nil, fmt.Errorf("user or signIn is nil")
	}

	webauthn, _ := newWebAuthn(signIn.Origin)
	sessionKey := GetSessionKey(&signIn.Registration)
	if opt, session, err := webauthn.BeginLogin(user); err != nil {
		return nil, err
	} else {
		store.set(sessionKey, webauthn, session, user)
		return opt, nil
	}
}

func (store *SessionStore) FinishAuthSession(signIn *SiginIn, ctx *gin.Context) (*User, *webauthn.Credential, error) {
	key := GetSessionKey(&signIn.Registration)
	if session := store.Get(key); session == nil {
		return nil, nil, fmt.Errorf("%s: not found", signIn.Email)
	} else {
		if cred, err := session.WebAuthn.FinishLogin(&session.User, session.Data, ctx.Request); err == nil {
			store.remove(key)
			session.User.UpdateCredential(cred)
			return &session.User, cred, nil
		} else {
			return nil, nil, err
		}
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

func (store *SessionStore) remove(id string) {
	store.locker.Lock()
	defer store.locker.Unlock()
	delete(store.sessions, id)
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
