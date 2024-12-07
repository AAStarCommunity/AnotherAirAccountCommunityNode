package seedworks

type AccountType string

const (
	EOA AccountType = "EOA"
)

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

type RegistrationByAccount struct {
	Origin  string      `json:"origin"`
	Type    AccountType `json:"type"`
	Account string      `json:"account"`
}
