package seedworks

import consts "another_node/internal/seedworks"

type CreateAARequest struct {
	Network consts.Chain `json:"network"`
	Alias   string       `json:"alias"`
}
