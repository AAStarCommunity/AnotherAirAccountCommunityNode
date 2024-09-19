package seedworks

import (
	"another_node/internal/community/account"
	"another_node/internal/community/chain"
	consts "another_node/internal/seedworks"
	"another_node/plugins/passkey_relay_party/storage/model"
	"crypto/ecdsa"
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-webauthn/webauthn/webauthn"
)

type UserChain struct {
	InitCode       string
	AA_Addr        string
	Alias          string
	Name           consts.Chain
	FromHdWalletId int64
}

type User struct {
	id             []byte
	credentials    []webauthn.Credential
	email          string
	wallets        []account.HdWallet
	chainAddresses map[string]UserChain
}

func (user *User) GetPrimaryEOA() string {
	if user.wallets == nil {
		return ""
	}
	for i := range user.wallets {
		if user.wallets[i].Primary {
			return user.wallets[i].Address
		}
	}
	return ""
}

func (user *User) TryCreateAA(network consts.Chain, alias string) (err error) {
	var w *account.HdWallet
	for i := range user.wallets {
		used := false
		for _, j := range user.GetChains() {
			if user.wallets[i].Id == j.FromHdWalletId && j.Name == network {
				used = true
				break
			}
		}
		if !used {
			w = &user.wallets[i]
			break
		}
	}

	if w == nil {
		return &ErrNoAvailableWallet{}
	}

	if _, aaAddr := user.GetChainAddresses(network, alias); aaAddr != nil {
		return nil
	}

	aa_address, initCode, err := chain.CreateSmartAccount(w, network)
	if err != nil {
		return err
	}

	user.SetAAWallet(w, &initCode, &aa_address, alias, network)

	return
}

func newUser(email string) *User {
	return &User{
		id:             []byte(email),
		email:          email,
		chainAddresses: make(map[string]UserChain),
	}
}

func MappingUser(airaccount *model.AirAccount, getFromVault func(vault *string) (string, error)) (*User, error) {
	user := &User{
		id:             []byte(airaccount.Email),
		email:          airaccount.Email,
		wallets:        make([]account.HdWallet, 0),
		credentials:    make([]webauthn.Credential, 0),
		chainAddresses: make(map[string]UserChain),
	}

	for i := range airaccount.HdWallet {
		if hdwalletStr, err := getFromVault(&airaccount.HdWallet[i].WalletVault); err != nil {
			return nil, err
		} else {
			if hdwalletStr != "" {
				var hdwallet account.HdWallet
				if err := json.Unmarshal([]byte(hdwalletStr), &hdwallet); err != nil {
					return nil, err
				} else {
					hdwallet.Id = airaccount.HdWallet[i].ID
					hdwallet.Primary = airaccount.HdWallet[i].Primary
					user.wallets = append(user.wallets, hdwallet)
				}
			} else {
				return nil, &ErrWalletNotFound{}
			}
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
		user.chainAddresses[chain.ChainName+":"+chain.Alias] = UserChain{
			InitCode:       chain.InitCode,
			AA_Addr:        chain.AA_Address,
			Name:           consts.Chain(chain.ChainName),
			Alias:          chain.Alias,
			FromHdWalletId: chain.FromWalletID,
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

func (user *User) GetChains() map[string]UserChain {
	return user.chainAddresses
}

func (user *User) GetWallets() []account.HdWallet {
	return user.wallets
}

func (user *User) GetChainAddresses(chain consts.Chain, alias string) (initCode, aaAddr *string) {
	if len(chain) == 0 {
		return nil, nil
	}

	if chainAddr, ok := user.chainAddresses[string(chain)+":"+alias]; ok {
		if strings.EqualFold(chainAddr.Alias, alias) {
			initCode = &chainAddr.InitCode
			aaAddr = &chainAddr.AA_Addr
			return
		}
	}
	return
}

func (user *User) GetPrivateKeyEcdsa() (*ecdsa.PrivateKey, error) {
	for i := range user.wallets {
		if user.wallets[i].Primary {
			return crypto.HexToECDSA(user.wallets[i].PrivateKey)
		}
	}
	panic("no primary wallet")
}

func (user *User) WalletMarshal() ([][]byte, error) {
	rlt := make([][]byte, 0)
	for i := range user.wallets {
		if wallet, err := json.Marshal(user.wallets[i]); err != nil {
			return nil, err
		} else {
			rlt = append(rlt, wallet)
		}
	}
	return rlt, nil
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

func (user *User) SetAAWallet(wallet *account.HdWallet, init_code, aa_address *string, alias string, network consts.Chain) {
	user.chainAddresses[string(network)+":"+alias] = UserChain{
		InitCode:       *init_code,
		AA_Addr:        *aa_address,
		Alias:          alias,
		Name:           network,
		FromHdWalletId: wallet.Id,
	}
}
