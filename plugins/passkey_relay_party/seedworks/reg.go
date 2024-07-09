package seedworks

import "another_node/internal/global_const"

type Registration struct {
	Origin  string `json:"origin"`
	Email   string `json:"email"`
	Captcha string `json:"captcha"`
	Network global_const.Network `json:"network"`
}
