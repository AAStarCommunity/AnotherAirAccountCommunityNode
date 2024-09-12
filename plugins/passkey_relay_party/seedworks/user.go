package seedworks

import (
	"another_node/internal/community/account"
	"another_node/internal/community/chain"
	consts "another_node/internal/seedworks"
	"another_node/plugins/passkey_relay_party/storage/model"
	"bytes"
	"crypto/ecdsa"
	"encoding/json"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-webauthn/webauthn/webauthn"
)

type userChain struct {
	InitCode string
	AA_Addr  string
	Alias    string
}

type User struct {
	id             []byte
	credentials    []webauthn.Credential
	email          string
	wallet         *account.HdWallet
	chainAddresses map[consts.Chain]userChain
}

func (user *User) GetEOA() string {
	return user.wallet.Address
}

func (user *User) TryCreateAA(network consts.Chain, alias string) (err error) {
	var w *account.HdWallet
	if user.wallet == nil || len(user.wallet.PrivateKey) == 0 {
		if w, err = account.NewHdWallet(account.HierarchicalPath_ETH); err != nil {
			return
		} else {
			user.wallet = w
		}
	} else {
		w = user.wallet
	}

	if _, aaAddr := user.GetChainAddresses(network, alias); aaAddr != nil {
		return nil
	}

	aa_address, initCode, err := chain.CreateSmartAccount(w, network)
	if err != nil {
		return err
	}
	user.SetAAWallet(&initCode, &aa_address, network)
	return nil
}

func newUser(email string) *User {
	return &User{
		id:             []byte(email),
		email:          email,
		chainAddresses: make(map[consts.Chain]userChain),
	}
}

func MappingUser(airaccount *model.AirAccount, getFromVault func() (string, error)) (*User, error) {
	user := &User{
		id:             []byte(airaccount.Email),
		email:          airaccount.Email,
		credentials:    make([]webauthn.Credential, 0),
		chainAddresses: make(map[consts.Chain]userChain),
	}

	if hdwalletStr, err := getFromVault(); err != nil {
		return nil, err
	} else {
		if hdwalletStr != "" {
			var hdwallet account.HdWallet
			if err := json.Unmarshal([]byte(hdwalletStr), &hdwallet); err != nil {
				return nil, err
			} else {
				user.wallet = &hdwallet
			}
		} else {
			return nil, &ErrWalletNotFound{}
		}
	}

	for i := range airaccount.Passkeys {
		passkey := airaccount.Passkeys[i]
		var cred webauthn.Credential
		json.Unmarshal([]byte(passkey.Rawdata), &cred)
		user.credentials = append(user.credentials, cred)
	}

	for i := range airaccount.AirAccountChains {
		chain := airaccount.AirAccountChains[i]
		user.chainAddresses[consts.Chain(chain.ChainName)] = userChain{
			InitCode: chain.InitCode,
			AA_Addr:  chain.AA_Address,
			Alias:    chain.Alias,
		}
	}
	return user, nil
}

var _ webauthn.User = (*User)(nil)

func (user *User) GetEmail() string {
	email, _, _ := user.GetAccounts()
	return email
}

func (user *User) GetAccounts() (email, facebook, twitter string) {
	email = user.email
	facebook = ""
	twitter = ""
	return
}

func (user *User) GetChains() map[consts.Chain]userChain {
	return user.chainAddresses
}

func (user *User) GetChainAddresses(chain consts.Chain, _ string) (initCode, aaAddr *string) {
	if len(chain) == 0 {
		return nil, nil
	}

	if chainAddr, ok := user.chainAddresses[chain]; ok {
		initCode = &chainAddr.InitCode
		aaAddr = &chainAddr.AA_Addr
		return
	} else {
		return nil, nil
	}
}

func (user *User) GetPrivateKeyStr() string {
	return user.wallet.PrivateKey
}
func (user *User) GetPrivateKeyEcdsa() (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(user.GetPrivateKeyStr())
}

func (user *User) WalletMarshal() ([]byte, error) {
	return json.Marshal(user.wallet)
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
		if bytes.Equal(c.ID, cred.ID) {
			user.credentials[i] = *cred
		}
	}
}

func (user *User) SetAAWallet(init_code, aa_address *string, network consts.Chain) {
	user.chainAddresses[network] = userChain{
		InitCode: *init_code,
		AA_Addr:  *aa_address,
	}
}
