package dvt

import (
	"another_node/plugins/dvt/passkey"
	"another_node/plugins/dvt/seedworks"
	"another_node/plugins/dvt/signature"

	"github.com/go-webauthn/webauthn/protocol"
)

func Signature(ca *protocol.ParsedCredentialAssertionData, publicKey *string) (blsSignature []byte, blsPublickey []byte, err error) {
	if ok, err := passkey.SignatureVerify(
		ca.Raw.AssertionResponse.AuthenticatorData,
		ca.Raw.AssertionResponse.ClientDataJSON, publicKey,
		ca.Response.Signature); !ok || err != nil {
		return nil, nil, func() error {
			if err != nil {
				return err
			}
			return seedworks.ErrSignatureVerifyFailed{}
		}()
	} else {
		return signature.Bls()
	}
}
