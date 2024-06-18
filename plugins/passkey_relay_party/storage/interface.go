package storage

import "another_node/plugins/passkey_relay_party/seedworks"

type Db interface {
	Save(user *seedworks.User) error
	Find(email string) (*seedworks.User, error)
}
