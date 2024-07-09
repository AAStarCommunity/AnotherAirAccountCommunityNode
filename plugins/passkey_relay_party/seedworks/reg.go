package seedworks

import "another_node/internal/global_const"

type Registration struct {
	Origin  string               `json:"origin"`
	Email   string               `json:"email"`
	Network global_const.Network `json:"network"`
}
