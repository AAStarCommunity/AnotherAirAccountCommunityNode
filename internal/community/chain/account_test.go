package chain

import (
	"another_node/conf"
	"another_node/internal/community/account"
	"another_node/internal/seedworks"
	"os"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	os.Setenv("Env", "dev")
	defer os.Unsetenv("Env")

	conf.Environment.Name = "dev"

	if testing.Short() {
		return
	}

	os.Chdir("../../../")

	w, err := account.NewHdWallet(account.HierarchicalPath(account.HierarchicalPath_ETH))
	if err != nil {
		t.Errorf("Failed to create account: %v", err)
	}

	address, initCode, err := CreateSmartAccount(w, seedworks.OptimismSepolia)
	if err != nil {
		t.Errorf("Failed to create account: %v", err)
	}
	if address == "" {
		t.Error("Expected account to be created, but got empty string")
	}
	t.Logf("address: %v", address)
	t.Logf("initCode: %v", initCode)

	// test code
}
