package dvt

import (
	"another_node/plugins/dvt/seedworks"
	"another_node/plugins/dvt/signature"

	"github.com/go-webauthn/webauthn/protocol"
)

func Signature(ca *protocol.ParsedCredentialAssertionData, publicKey []byte) (blsSignature []byte, blsPublickey []byte, err error) {
	if ca == nil || publicKey == nil {
		return nil, nil, seedworks.ErrSignatureVerifyFailed{}
	}

	verifier := NewTxVerifier(publicKey)

	if ok, err := verifier.Verify(
		ca.Raw.AssertionResponse.AuthenticatorData,
		ca.Raw.AssertionResponse.ClientDataJSON,
		ca.Response.Signature); !ok || err != nil {
		return nil, nil, func() error {
			if err != nil {
				return err
			}
			return seedworks.ErrSignatureVerifyFailed{}
		}()
	} else {
		return signature.BlsTss(1, 1, ca.Response.Signature)
	}
}
