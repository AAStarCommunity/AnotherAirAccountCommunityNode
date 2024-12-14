package seedworks

import (
	"another_node/internal/seedworks"

	"github.com/go-webauthn/webauthn/protocol"
)

type TxSignature struct {
	Ticket       string                                  `json:"ticket"`
	Account      string                                  `json:"-"`
	Origin       string                                  `json:"origin"`
	TxData       string                                  `json:"txdata"`
	Network      seedworks.Chain                         `json:"network"`
	NetworkAlias string                                  `json:"network_alias"`
	CA           *protocol.ParsedCredentialAssertionData `json:"-"`
	CAPublicKey  []byte                                  `json:"-"`
}

type ValidatorResult struct {
	Message    []string `json:"message"`
	PublicKeys []string `json:"pubkeys"`
}

// DvtResult aggregate distributed validators signatures by BLS
type DvtResult struct {
	Signatures []string          `json:"signatures"`
	Validator  []ValidatorResult `json:"validator"`
}

type TxSignatureResult struct {
	Code    int        `json:"code"`
	DVT     *DvtResult `json:"dvt"`
	TxData  string     `json:"txdata"`
	Sign    string     `json:"sign"`
	Address string     `json:"address"`
}
