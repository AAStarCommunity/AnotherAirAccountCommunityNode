package seedworks

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
	Origin string `json:"origin"`
}
