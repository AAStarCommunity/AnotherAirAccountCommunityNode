package storage

import (
	"another_node/plugins/passkey_relay_party/seedworks"
	"errors"
)

type InMemory struct {
	users map[string]*seedworks.User
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

	return nil, seedworks.UserNotFoundError{}
}
