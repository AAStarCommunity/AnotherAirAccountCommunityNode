package seedworks

import (
	"another_node/internal/community/account"
	"another_node/internal/seedworks"
	"reflect"
	"testing"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

func TestUserMarshal(t *testing.T) {
	u := newUser("mail@domain.io")

	u.AddCredential(&webauthn.Credential{
		ID:        []byte("id"),
		PublicKey: []byte("publicKey"),
		Transport: []protocol.AuthenticatorTransport{
			protocol.USB,
			protocol.Hybrid,
		},
		Flags: webauthn.CredentialFlags{
			UserPresent:  true,
			UserVerified: false,
		},
		Authenticator: webauthn.Authenticator{
			AAGUID:       []byte("AAGUID"),
			SignCount:    1,
			CloneWarning: false,
			Attachment:   protocol.Platform,
		},
	})
	u.AddCredential(&webauthn.Credential{
		ID:        []byte("id2"),
		PublicKey: []byte("publicKey2"),
		Transport: []protocol.AuthenticatorTransport{
			protocol.USB,
		},
		Flags: webauthn.CredentialFlags{
			UserPresent:  false,
			UserVerified: true,
		},
		Authenticator: webauthn.Authenticator{
			AAGUID:       []byte("AAGUID2"),
			SignCount:    3,
			CloneWarning: true,
			Attachment:   protocol.CrossPlatform,
		},
	})

	hdWallet, err := account.NewHdWallet(account.HierarchicalPath_ETH)
	if err != nil {
		t.Fatal(err)
	}
	u.SetWallet(hdWallet, "abc", seedworks.ArbitrumNova)

	if s, err := u.Marshal(); err != nil {
		t.Fatal(err)
	} else {
		raw := string(s)

		if u2, err := UnmarshalUser(&raw); err != nil {
			t.Fatal(err)
		} else {
			if !reflect.DeepEqual(u, u2) {
				t.Fatal("not equal")
			}
		}
	}
}
