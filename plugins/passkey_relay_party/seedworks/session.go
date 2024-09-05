package seedworks

import (
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

func GetSessionKey(origin, id string, ext ...string) string {
	m := ""
	if len(ext) > 0 {
		m = ":" + strings.Join(ext, ":")
	}
	return origin + ":" + id + m
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

func (store *SessionStore) NewRegSession(reg *RegistrationByEmail) (*protocol.CredentialCreation, error) {
	user := newUser(reg.Email)
	wan, _ := newWebAuthn(reg.Origin)
	sessionKey := GetSessionKey(reg.Origin, reg.Email)

	authSelect := protocol.AuthenticatorSelection{
		AuthenticatorAttachment: protocol.Platform,
		RequireResidentKey:      protocol.ResidentKeyNotRequired(),
		UserVerification:        protocol.VerificationRequired,
	}

	if opt, session, err := wan.BeginRegistration(user,
		webauthn.WithAuthenticatorSelection(authSelect),
	); err != nil {
		return nil, err
	} else {
		store.set(sessionKey, wan, session, user)
		return opt, nil
	}
}

func (store *SessionStore) FinishRegSession(reg *FinishRegistrationByEmail, ctx *gin.Context) (*User, error) {
	key := GetSessionKey(reg.Origin, reg.Email)
	if session := store.Get(key); session == nil {
		return nil, fmt.Errorf("%s: not found", reg.Email)
	} else {
		if cred, err := session.WebAuthn.FinishRegistration(&session.User, session.Data, ctx.Request); err == nil {
			session.User.AddCredential(cred)
			store.Remove(key)
			return &session.User, nil
		} else {
			return nil, err
		}
	}
}

const session_key_discoverlogin string = "DiscoverLogin"

func (store *SessionStore) NewDiscoverableAuthSession(user *User, signIn *SiginIn) (*protocol.CredentialAssertion, error) {
	webauthn, _ := newWebAuthn(signIn.Origin)
	sessionKey := GetSessionKey(signIn.Origin, session_key_discoverlogin)
	if opt, session, err := webauthn.BeginDiscoverableLogin(); err != nil {
		return nil, err
	} else {
		store.set(sessionKey, webauthn, session, user)
		return opt, nil
	}
}

func (store *SessionStore) FinishDiscoverableAuthSession(signIn *SiginIn, ctx *gin.Context) (*User, error) {
	key := GetSessionKey(signIn.Origin, session_key_discoverlogin)
	if session := store.Get(key); session == nil {
		return nil, fmt.Errorf("not found")
	} else {
		defer store.Remove(key)
		if _, err := session.WebAuthn.FinishDiscoverableLogin(func(rawID, userHandle []byte) (user webauthn.User, err error) {
			rawIDStr := base64.URLEncoding.EncodeToString(rawID)
			_ = rawIDStr
			return &session.User, nil
		}, session.Data, ctx.Request); err == nil {
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
	sessionKey := GetSessionKey(signIn.Origin, signIn.Email)
	if opt, session, err := webauthn.BeginLogin(user); err != nil {
		return nil, err
	} else {
		store.set(sessionKey, webauthn, session, user)
		return opt, nil
	}
}

func (store *SessionStore) FinishAuthSession(signIn *SiginIn, ctx *gin.Context) (*User, error) {
	key := GetSessionKey(signIn.Origin, signIn.Email)
	if session := store.Get(key); session == nil {
		return nil, fmt.Errorf("%s: not found", signIn.Email)
	} else {
		defer store.Remove(key)
		if _, err := session.WebAuthn.FinishLogin(&session.User, session.Data, ctx.Request); err == nil {
			return &session.User, nil
		} else {
			return nil, err
		}
	}
}

func (store *SessionStore) NewTxSession(user *User, txSignature *TxSignature) (*protocol.CredentialAssertion, error) {
	if user == nil || txSignature == nil {
		return nil, fmt.Errorf("user or signIn is nil")
	}

	webAuthn, _ := newWebAuthn(txSignature.Origin)
	sessionKey := GetSessionKey(txSignature.Origin, txSignature.Email, txSignature.Ticket)
	if opt, session, err := webAuthn.BeginLogin(user,
		func(opt *protocol.PublicKeyCredentialRequestOptions) {
			opt.Challenge = protocol.URLEncodedBase64(txSignature.TxData)
			if opt.Extensions == nil {
				opt.Extensions = make(map[string]interface{})
			}
			opt.Extensions["ticket"] = txSignature.Ticket
		},
	); err != nil {
		return nil, err
	} else {
		store.set(sessionKey, webAuthn, session, user)
		return opt, nil
	}
}

func (store *SessionStore) FinishSignSession(paymentSign *TxSignature, ctx *gin.Context) (*User, error) {
	key := GetSessionKey(paymentSign.Origin, paymentSign.Email, paymentSign.Ticket)
	if session := store.Get(key); session == nil {
		return nil, fmt.Errorf("%s: not found", paymentSign.Email)
	} else {
		if _, err := session.WebAuthn.FinishLogin(&session.User, session.Data, ctx.Request); err == nil {
			store.Remove(key)
			if paymentSign.Ticket != session.Data.Extensions["ticket"].(string) {
				return nil, fmt.Errorf("ticket not match")
			}
			if txData, err := base64.RawURLEncoding.DecodeString(session.Data.Challenge); err != nil {
				return nil, err
			} else {
				paymentSign.TxData = string(txData)
			}
			return &session.User, nil
		} else {
			return nil, err
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

func (store *SessionStore) Remove(id string) {
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
		User: func() User {
			if user == nil {
				return User{}
			} else {
				return *user
			}
		}(),
		expires: 120,
	}
	cache.countdown()
	store.sessions[key] = cache
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
