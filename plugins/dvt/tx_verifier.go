package dvt

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/protocol/webauthncose"
)

type TxVerifier struct {
	publicKey *string
}

func NewTxVerifier(publicKey *string) *TxVerifier {
	return &TxVerifier{publicKey: publicKey}
}

func (v *TxVerifier) Verify(authData, clientDataJson protocol.URLEncodedBase64, signature []byte) (bool, error) {
	pubKeyBytes, err := base64.RawURLEncoding.DecodeString(*v.publicKey)
	if err != nil {
		return false, err
	}
	pubKey, err := webauthncose.ParsePublicKey(pubKeyBytes)
	if err != nil {
		return false, err
	}
	clientDataHash := sha256.Sum256([]byte(clientDataJson))
	// sigData = authData + sha256(clientData)
	sigData := append([]byte(authData), clientDataHash[:]...)
	return webauthncose.VerifySignature(pubKey, sigData, signature)
}
