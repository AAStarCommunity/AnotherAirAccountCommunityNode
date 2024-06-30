package seedworks

import (
	"another_node/internal/community/account"

	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	id          []byte
	credentials []webauthn.Credential
	email       string
	wallet      *account.HdWallet
	address     string
}

func newUser(email string) *User {
	return &User{
		id:    []byte(email),
		email: email,
	}
}

var _ webauthn.User = (*User)(nil)

func (user *User) WebAuthnID() []byte {
	return user.id
}

func (user *User) WebAuthnName() string {
	return user.email
}

func (user *User) WebAuthnDisplayName() string {
	return user.email
}

func (user *User) WebAuthnCredentials() []webauthn.Credential {
	return user.credentials
}

// WebAuthnIcon is a deprecated option.
// Deprecated: this has been removed from the specification recommendation. Suggest a blank string.
func (user *User) WebAuthnIcon() string {
	return ""
}

func (user *User) AddCredential(cred *webauthn.Credential) {
	user.credentials = append(user.credentials, *cred)
}
func (user *User) UpdateCredential(cred *webauthn.Credential) {
	for i, c := range user.credentials {
		if string(c.ID) == string(cred.ID) {
			user.credentials[i] = *cred
		}
	}
}

func (user *User) SetWallet(wallet *account.HdWallet, address string) {
	user.wallet = wallet
	user.address = address
}
