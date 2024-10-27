package dvt

import (
	"crypto/sha256"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/protocol/webauthncose"
)

type TxVerifier struct {
	publicKey []byte
}

func NewTxVerifier(publicKey []byte) *TxVerifier {
	return &TxVerifier{publicKey: publicKey}
}

func (v *TxVerifier) Verify(authData, clientDataJson protocol.URLEncodedBase64, signature []byte) (bool, error) {
	pubKey, err := webauthncose.ParsePublicKey(v.publicKey)
	if err != nil {
		return false, err
	}
	clientDataHash := sha256.Sum256([]byte(clientDataJson))
	// sigData = authData + sha256(clientData)
	sigData := append([]byte(authData), clientDataHash[:]...)
	return webauthncose.VerifySignature(pubKey, sigData, signature)
}
