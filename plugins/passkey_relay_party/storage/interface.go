package storage

import (
	consts "another_node/internal/seedworks"
	"another_node/plugins/passkey_relay_party/seedworks"
	"another_node/plugins/passkey_relay_party/storage/model"
)

type Db interface {
	FindUser(email string) (*seedworks.User, error)
	SaveChallenge(captchaType model.ChallengeType, challenger, captcha string) error
	Challenge(captchaType model.ChallengeType, challenger, captcha string) bool
	SaveAccounts(user *seedworks.User, network consts.Chain) error
	GetAccountsByEmail(email, network string) (initCode, addr, eoaAddr string, err error)
}
