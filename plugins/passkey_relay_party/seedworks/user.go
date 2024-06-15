package seedwork

import "github.com/go-webauthn/webauthn/webauthn"

type User struct {
	id          []byte
	credentials []webauthn.Credential
	name        string
	displayName string
}

func NewUser(id, name, displayName string) *User {
	return &User{
		id:          []byte(id),
		name:        name,
		displayName: displayName,
	}
}

var _ webauthn.User = (*User)(nil)

func (user *User) WebAuthnID() []byte {
	return user.id
}

func (user *User) WebAuthnName() string {
	return user.name
}

func (user *User) WebAuthnDisplayName() string {
	return user.displayName
}

func (user *User) WebAuthnCredentials() []webauthn.Credential {
	return user.credentials
}

// WebAuthnIcon is a deprecated option.
// Deprecated: this has been removed from the specification recommendation. Suggest a blank string.
func (user *User) WebAuthnIcon() string {
	return ""
}
