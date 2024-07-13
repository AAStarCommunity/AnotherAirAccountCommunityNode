package seedworks

import consts "another_node/internal/seedworks"

type RegistrationPrepare struct {
	Email string `json:"email"`
}

type Registration struct {
	RegistrationPrepare
	Origin  string         `json:"origin"`
	Captcha string         `json:"captcha"`
	Network consts.Network `json:"network"`
}
