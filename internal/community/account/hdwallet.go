package account

import (
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

type HierarchicalPath string

// HierarchicalPath_ETH is the default path for Mainet / Eth / TestNet
const HierarchicalPath_ETH HierarchicalPath = "m/44'/60'/0'/0/0"
const HierarchicalPath_ETH_FMT string = "m/44'/60'/0'/0/%d"

type HdWallet struct {
	Id         int64  `json:"-"`
	Mnemonic   string `json:"mnemonic"`
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

func newWallet() (*hdwallet.Wallet, *string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, nil, err
	}
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, nil, err
	}

	return wallet, &mnemonic, nil
}

func NewHdWallet(hierarchicalPath ...HierarchicalPath) ([]HdWallet, error) {
	if wallet, mnemonic, err := newWallet(); err != nil {
		return nil, err
	} else {
		hdwallets := make([]HdWallet, 0)
		for p := range hierarchicalPath {
			path := hdwallet.MustParseDerivationPath(string(hierarchicalPath[p]))
			account, err := wallet.Derive(path, false)
			if err != nil {
				return nil, err
			}

			privateKey, err := wallet.PrivateKeyHex(account)
			if err != nil {
				return nil, err
			}

			hdwallets = append(hdwallets, HdWallet{
				Mnemonic:   *mnemonic,
				Address:    account.Address.Hex(),
				PrivateKey: privateKey,
			})
		}
		return hdwallets, nil
	}
}
