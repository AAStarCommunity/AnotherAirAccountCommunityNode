package storage

import "another_node/plugins/passkey_relay_party/seedworks"

type PgsqlStorage struct {
}

var _ Db = (*PgsqlStorage)(nil)

func NewPgsqlStorage() *PgsqlStorage {
	return &PgsqlStorage{}
}

func (db *PgsqlStorage) Save(user *seedworks.User) error {
	return nil
}

func (db *PgsqlStorage) Find(email string) (*seedworks.User, error) {
	return nil, nil
}

func (db *PgsqlStorage) SaveChallenge(email, captcha string) error {
	return nil
}

func (db *PgsqlStorage) Challenge(email, captcha string) bool {
	return true
}
