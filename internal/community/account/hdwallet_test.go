package account

import (
	"testing"
)

func TestNewHdWallet(t *testing.T) {
	hierarchicalPath := HierarchicalPath(HierarchicalPath_ETH)

	wallet, err := NewHdWallet(hierarchicalPath)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if wallet == nil {
		t.Error("expected wallet to be created, but got nil")
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
