package seedworks

type Registration struct {
	Origin  string `json:"origin"`
	Email   string `json:"email"`
	Captcha string `json:"captcha"`
}
