package seedworks

import (
	"fmt"
	"math/rand"
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

func (store *SessionStore) NewRegSession(reg *Registration) (*protocol.CredentialCreation, error) {
	user := newUser(reg.Email)
	wan, _ := newWebAuthn(reg.Origin)
	sessionKey := GetSessionKey(reg.Origin, reg.Email)

	authSelect := protocol.AuthenticatorSelection{
		AuthenticatorAttachment: protocol.Platform,
		RequireResidentKey:      protocol.ResidentKeyNotRequired(),
		UserVerification:        protocol.VerificationRequired,
	}

	challenge := make([]byte, 36)
	if opt, session, err := wan.BeginRegistration(user,
		func() protocol.URLEncodedBase64 {
			rand.Read(challenge)
			return protocol.URLEncodedBase64(challenge)
		}(),
		webauthn.WithAuthenticatorSelection(authSelect),
	); err != nil {
		return nil, err
	} else {
		store.set(sessionKey, wan, session, user)
		return opt, nil
	}
}

func (store *SessionStore) FinishRegSession(reg *FinishRegistration, ctx *gin.Context) (*User, error) {
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

func (store *SessionStore) FinishAuthSession(signIn *SiginIn, ctx *gin.Context) (*User, *webauthn.Credential, error) {
	key := GetSessionKey(signIn.Origin, signIn.Email)
	if session := store.Get(key); session == nil {
		return nil, nil, fmt.Errorf("%s: not found", signIn.Email)
	} else {
		if cred, err := session.WebAuthn.FinishLogin(&session.User, session.Data, ctx.Request); err == nil {
			store.Remove(key)
			session.User.UpdateCredential(cred)
			return &session.User, cred, nil
		} else {
			return nil, nil, err
		}
	}
}

func (store *SessionStore) NewTxSession(user *User, txSignature *TxSignature) (*protocol.CredentialAssertion, error) {
	if user == nil || txSignature == nil {
		return nil, fmt.Errorf("user or signIn is nil")
	}

	webAuthn, _ := newWebAuthn(txSignature.Origin)
	sessionKey := GetSessionKey(txSignature.Origin, txSignature.Email, txSignature.Nonce)
	if opt, session, err := webAuthn.BeginLogin(user,
		func(opt *protocol.PublicKeyCredentialRequestOptions) {
			// opt.Challenge, _ = CreateChallenge(txSignature) // TODO: rewrite the challenge algorithm
			if opt.Extensions == nil {
				opt.Extensions = make(map[string]interface{})
			}
			opt.Extensions["txdata"] = txSignature.TxData
			opt.Extensions["nonce"] = txSignature.Nonce
		},
	); err != nil {
		return nil, err
	} else {
		store.set(sessionKey, webAuthn, session, user)
		return opt, nil
	}
}

func (store *SessionStore) FinishSignSession(paymentSign *TxSignature, ctx *gin.Context) (*User, error) {
	key := GetSessionKey(paymentSign.Origin, paymentSign.Email, paymentSign.Nonce)
	if session := store.Get(key); session == nil {
		return nil, fmt.Errorf("%s: not found", paymentSign.Email)
	} else {
		if _, err := session.WebAuthn.FinishLogin(&session.User, session.Data, ctx.Request); err == nil {
			store.Remove(key)
			paymentSign.TxData = session.Data.Extensions["txdata"].(string)
			if paymentSign.Nonce != session.Data.Extensions["nonce"].(string) {
				return nil, fmt.Errorf("nonce not match")
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
		User:     *user,
		expires:  120,
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
