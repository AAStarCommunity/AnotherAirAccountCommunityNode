package seedwork

import (
	"sync"
	"time"

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

func newSessionCache(data *webauthn.SessionData) *sessionCache {
	p := &sessionCache{
		data:    data,
		expires: 120,
	}
	p.countdown()
	return p
}

type sessionCache struct {
	data     *webauthn.SessionData
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
				cache.data = nil
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
				if session.data == nil || session.obsolete {
					delete(store.sessions, id)
				}
			}
			store.locker.Unlock()
		}
	}()
}

func (store *SessionStore) Get(id string) *webauthn.SessionData {
	store.locker.RLock()
	defer store.locker.RUnlock()
	if session, ok := store.sessions[id]; ok {
		if session.obsolete {
			delete(store.sessions, id)
			return nil
		}
		return session.data
	}
	return nil
}

func (store *SessionStore) Set(id string, session *webauthn.SessionData) {
	store.locker.Lock()
	defer store.locker.Unlock()
	store.sessions[id] = newSessionCache(session)
}
