package seedworks

import "github.com/go-webauthn/webauthn/protocol"

// CreateChallenge custom define the challenge code
func CreateChallenge(tx *TxSignature) (protocol.URLEncodedBase64, error) {
	return protocol.CreateChallenge()
}
