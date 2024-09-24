package plugin_passkey_relay_party

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/protocol/webauthncose"
)

func VerifySignature(data *protocol.ParsedCredentialAssertionData, publicKey *string, signature []byte) (bool, error) {
	pubKeyBytes, err := base64.RawURLEncoding.DecodeString(*publicKey)
	if err != nil {
		return false, err
	}
	authData := data.Raw.AssertionResponse.AuthenticatorData
	pubKey, err := webauthncose.ParsePublicKey(pubKeyBytes)
	if err != nil {
		return false, err
	}
	sig := data.Response.Signature
	clientDataHash := sha256.Sum256(data.Raw.AssertionResponse.ClientDataJSON)
	// sigData = authData + sha256(clientData)
	sigData := append(authData, clientDataHash[:]...)
	return webauthncose.VerifySignature(pubKey, sigData, sig)
}
