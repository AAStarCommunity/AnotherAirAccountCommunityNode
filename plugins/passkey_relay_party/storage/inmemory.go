package storage

import (
	"another_node/plugins/passkey_relay_party/seedworks"
	"errors"
	"time"
)

type captcha struct {
	code string
	time time.Time
}

type InMemory struct {
	users    map[string]*seedworks.User
	captchas map[string]captcha
}

var _ Db = (*InMemory)(nil)

func NewInMemory() *InMemory {
	return &InMemory{
		users: make(map[string]*seedworks.User),
	}
}

func (db *InMemory) Save(user *seedworks.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	db.users[string(user.WebAuthnID())] = user

	return nil
}

func (db *InMemory) Find(email string) (*seedworks.User, error) {
	if user, ok := db.users[email]; ok {
		return user, nil
	}

	return nil, seedworks.ErrUserNotFound{}
}

func (db *InMemory) SaveChallenge(email, code string) error {
	db.captchas[email] = captcha{
		code: code,
		time: time.Now(),
	}
	return nil
}
func (db *InMemory) Challenge(email, code string) bool {
	if v, ok := db.captchas[email]; ok {
		if v.code == code && time.Now().Sub(v.time) < 10*time.Minute {
			return true
		}
	}
	return false
}
