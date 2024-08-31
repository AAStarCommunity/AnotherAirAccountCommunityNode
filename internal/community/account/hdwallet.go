package account

import (
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

type HierarchicalPath string

// HierarchicalPath_ETH is the default path for Mainet / Eth / TestNet
const HierarchicalPath_ETH HierarchicalPath = "m/44'/60'/0'/0/0"

type HdWallet struct {
	Mnemonic   string `json:"mnemonic"`
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

func NewHdWallet(hierarchicalPath HierarchicalPath) (*HdWallet, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	path := hdwallet.MustParseDerivationPath(string(hierarchicalPath))
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, err
	}

	privateKey, err := wallet.PrivateKeyHex(account)
	if err != nil {
		return nil, err
	}

	return &HdWallet{
		Mnemonic:   mnemonic,
		Address:    account.Address.Hex(),
		PrivateKey: privateKey,
	}, nil
}
