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

// DvtResult aggregate distributed validators signatures by BLS
type DvtResult struct {
	MessagePt  [2]string   `json:"message"`
	Signatures []string    `json:"signatures"`
	PublicKeys [][4]string `json:"pubkeys"`
}

type TxSignatureResult struct {
	Code    int        `json:"code"`
	DVT     *DvtResult `json:"dvt"`
	TxData  string     `json:"txdata"`
	Sign    string     `json:"sign"`
	Address string     `json:"address"`
}
