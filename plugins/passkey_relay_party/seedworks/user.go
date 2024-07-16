package seedworks

import (
	"another_node/internal/community/account"
	consts "another_node/internal/seedworks"
	"encoding/json"

	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	id             []byte
	credentials    []webauthn.Credential
	email          string
	wallet         *account.HdWallet
	chainAddresses map[consts.Chain]string
}

func newUser(email string) *User {
	return &User{
		id:             []byte(email),
		email:          email,
		chainAddresses: make(map[consts.Chain]string),
	}
}

type marshalHDWallet struct {
	Mnemonic   string `json:"mnemonic"`
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

type marshalUser struct {
	Id          []byte                  `json:"id"`
	Email       string                  `json:"email"`
	Credentials []webauthn.Credential   `json:"credentials"`
	Wallet      marshalHDWallet         `json:"wallet"`
	AddressMap  map[consts.Chain]string `json:"addressMap"`
}

func UnmarshalUser(rawdata *string) (*User, error) {
	m := marshalUser{}

	if err := json.Unmarshal([]byte(*rawdata), &m); err != nil {
		return nil, err
	}

	return &User{
		id:             m.Id,
		email:          m.Email,
		credentials:    m.Credentials,
		wallet:         account.RecoverHdWallet(&m.Wallet.Mnemonic, &m.Wallet.Address, &m.Wallet.PrivateKey),
		chainAddresses: m.AddressMap,
	}, nil
}

var _ webauthn.User = (*User)(nil)

func (user *User) GetEmail() string {
	return user.email
}
func (user *User) GetPrivateKey() string {
	return user.wallet.PrivateKey()
}

func (user *User) Marshal() ([]byte, error) {
	return json.Marshal(marshalUser{
		Id:          user.id,
		Email:       user.email,
		Credentials: user.credentials,
		Wallet: marshalHDWallet{
			Mnemonic:   user.wallet.Mnemonic(),
			Address:    user.wallet.Address(),
			PrivateKey: user.wallet.PrivateKey(),
		},
		AddressMap: user.chainAddresses,
	})
}

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

func (user *User) SetWallet(wallet *account.HdWallet, address string, network consts.Chain) {
	user.wallet = wallet
	user.chainAddresses[network] = address
}
