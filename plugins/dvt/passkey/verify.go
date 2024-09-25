package passkey

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/protocol/webauthncose"
)

func SignatureVerify(authData, clientDataJson protocol.URLEncodedBase64, publicKey *string, signature []byte) (bool, error) {
	return verifySignature(authData, clientDataJson, publicKey, signature)
}

func extractAssertionData(ctx *http.Request) (authData, clientDataJson protocol.URLEncodedBase64, signature []byte, err error) {
	p, err := protocol.ParseCredentialRequestResponse(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	authData = p.Raw.AssertionResponse.AuthenticatorData
	clientDataJson = p.Raw.AssertionResponse.ClientDataJSON
	signature = p.Response.Signature
	return authData, clientDataJson, signature, nil
}

func verifySignature(authData, clientDataJson protocol.URLEncodedBase64, publicKey *string, signature []byte) (bool, error) {
	pubKeyBytes, err := base64.RawURLEncoding.DecodeString(*publicKey)
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
