package account

import (
	"testing"
)

func TestNewHdWallet(t *testing.T) {
	hierarchicalPath := HierarchicalPath(HierarchicalPath_Main_ETH_TestNet)

	wallet, err := NewHdWallet(hierarchicalPath)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if wallet == nil {
		t.Error("expected wallet to be created, but got nil")
		return
	}

	if wallet.mnemonic == "" {
		t.Error("expected mnemonic to be set, but got empty string")
	}

	if wallet.address == "" {
		t.Error("expected address to be set, but got empty string")
	}

	if len(wallet.privateKey) == 0 {
		t.Error("expected privateKey to be set, but got empty slice")
	}
}
