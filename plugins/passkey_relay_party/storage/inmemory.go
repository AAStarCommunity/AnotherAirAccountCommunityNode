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
	accounts map[string]interface{}
}

var _ Db = (*InMemory)(nil)

func NewInMemory() *InMemory {
	return &InMemory{
		users: make(map[string]*seedworks.User),
	}
}

func (db *InMemory) Save(user *seedworks.User, allowUpdate bool) error {
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
		if v.code == code && time.Since(v.time) < 10*time.Minute {
			return true
		}
	}
	return false
}

func (db *InMemory) SaveAccounts(user *seedworks.User, initCode, addr, eoaAddr, chain string) error {
	if len(db.accounts) == 0 {
		db.accounts = make(map[string]interface{})
	}

	db.accounts[user.GetEmail()+":"+chain] = &struct {
		InitCode   string
		Address    string
		EoaAddress string
		Chain      string
	}{
		InitCode:   initCode,
		Address:    addr,
		EoaAddress: eoaAddr,
	}
	return nil
}

func (db *InMemory) GetAccounts(email, chain string) (initCode, addr, eoaAddr string, err error) {
	if v, ok := db.accounts[email+":"+chain]; ok {
		if a, ok := v.(*struct {
			InitCode   string
			Address    string
			EoaAddress string
		}); ok {
			return a.InitCode, a.Address, a.EoaAddress, nil
		}
	}
	return "", "", "", errors.New("account not found")
}
