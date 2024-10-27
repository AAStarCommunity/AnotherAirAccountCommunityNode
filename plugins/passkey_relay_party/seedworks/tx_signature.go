package seedworks

import (
	"another_node/internal/seedworks"

	"github.com/go-webauthn/webauthn/protocol"
)

type TxSignature struct {
	Ticket       string                                  `json:"ticket"`
	Email        string                                  `json:"-"`
	Origin       string                                  `json:"origin"`
	TxData       string                                  `json:"txdata"`
	Network      seedworks.Chain                         `json:"network"`
	NetworkAlias string                                  `json:"network_alias"`
	CA           *protocol.ParsedCredentialAssertionData `json:"-"`
	CAPublicKey  []byte                                  `json:"-"`
}

type TxSignatureResult struct {
	Code      int    `json:"code"`
	TxData    string `json:"txdata"`
	Sign      string `json:"sign"`
	BlsSign   string `json:"bls_sign"`
	BlsPubKey string `json:"bls_pubkey"`
	Address   string `json:"address"`
	BlsSchema string `json:"bls_schema"`
}
