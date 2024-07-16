package storage

import (
	"another_node/plugins/passkey_relay_party/seedworks"
)

type Db interface {
	Save(user *seedworks.User, allowUpdate bool) error
	Find(email string) (*seedworks.User, error)
	SaveChallenge(email, captcha string) error
	Challenge(email, captcha string) bool
	SaveAccounts(user *seedworks.User, initCode, addr, eoaAddr, network string) error
	GetAccounts(email, network string) (initCode, addr, eoaAddr string, err error)
}
