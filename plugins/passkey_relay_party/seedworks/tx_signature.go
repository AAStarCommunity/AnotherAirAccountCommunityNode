package seedworks

import "another_node/internal/seedworks"

type TxSignature struct {
	Ticket       string          `json:"ticket"`
	Email        string          `json:"-"`
	Origin       string          `json:"origin"`
	TxData       string          `json:"txdata"`
	Network      seedworks.Chain `json:"network"`
	NetworkAlias string          `json:"network_alias"`
}

type TxSignatureResult struct {
	Code    int    `json:"code"`
	TxData  string `json:"txdata"`
	Sign    string `json:"sign"`
	Address string `json:"address"`
}
