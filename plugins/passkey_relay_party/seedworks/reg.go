package seedworks

import consts "another_node/internal/seedworks"

type RegistrationByEmailPrepare struct {
	Email string `json:"email"`
}

type RegistrationByEmail struct {
	RegistrationByEmailPrepare
	Origin  string `json:"origin"`
	Captcha string `json:"captcha"`
}

type FinishRegistrationByEmail struct {
	RegistrationByEmailPrepare
	Origin  string       `json:"origin"`
	Network consts.Chain `json:"network"`
	Alias   string       `json:"alias"`
}
