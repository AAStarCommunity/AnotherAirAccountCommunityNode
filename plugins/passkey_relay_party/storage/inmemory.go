package storage

import (
	consts "another_node/internal/seedworks"
	"another_node/plugins/passkey_relay_party/seedworks"
	"another_node/plugins/passkey_relay_party/storage/model"
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
	accounts map[string]interface{}
}

var _ Db = (*InMemory)(nil)

func NewInMemory() *InMemory {
	return &InMemory{
		users: make(map[string]*seedworks.User),
	}
}

func save(db *InMemory, user *seedworks.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	db.users[string(user.WebAuthnID())] = user

	return nil
}

func (db *InMemory) FindUser(email string) (*seedworks.User, error) {
	if user, ok := db.users[email]; ok {
		return user, nil
	}

	return nil, seedworks.ErrUserNotFound{}
}

func (db *InMemory) SaveChallenge(_ model.ChallengeType, email, code string) error {
	db.captchas[email] = captcha{
		code: code,
		time: time.Now(),
	}
	return nil
}
func (db *InMemory) Challenge(_ model.ChallengeType, email, code string) bool {
	if v, ok := db.captchas[email]; ok {
		if v.code == code && time.Since(v.time) < 10*time.Minute {
			return true
		}
	}
	return false
}

func (db *InMemory) SaveAccounts(user *seedworks.User, chain consts.Chain) error {
	save(db, user)
	if len(db.accounts) == 0 {
		db.accounts = make(map[string]interface{})
	}

	initCode, addr, eoaAddr := user.GetChainAddresses(chain)

	db.accounts[user.GetEmail()+":"+string(chain)] = &struct {
		InitCode   string
		Address    string
		EoaAddress string
		Chain      string
	}{
		InitCode:   *initCode,
		Address:    *addr,
		EoaAddress: *eoaAddr,
	}
	return nil
}
