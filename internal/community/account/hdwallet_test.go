package account

import (
	"testing"
)

func TestNewHdWallet(t *testing.T) {
	// hierarchicalPath := HierarchicalPath(HierarchicalPath_ETH)
	paths := []HierarchicalPath{
		"m/44'/60'/0'/0/0",
		"m/44'/60'/0'/0/1",
		"m/44'/60'/0'/0/2",
	}
	wallets, err := NewHdWallet(paths...)
	for _, wallet := range wallets {

		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if wallet.Mnemonic == "" {
			t.Error("expected mnemonic to be set, but got empty string")
		}

		if wallet.Address == "" {
			t.Error("expected address to be set, but got empty string")
		}

		if len(wallet.PrivateKey) == 0 {
			t.Error("expected privateKey to be set, but got empty slice")
		}
	}
}
