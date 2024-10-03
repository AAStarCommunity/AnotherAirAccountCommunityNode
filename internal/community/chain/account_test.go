package chain

import (
	"another_node/conf"
	"another_node/internal/community/account"
	"another_node/internal/seedworks"
	"fmt"
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

	p := make([]account.HierarchicalPath, 2)
	for i := 0; i < 2; i++ {
		p[i] = account.HierarchicalPath(fmt.Sprintf(account.HierarchicalPath_ETH_FMT, i))
	}
	wallets, err := account.NewHdWallet(p...)
	if err != nil {
		t.Errorf("Failed to create account: %v", err)
	}

	for _, w := range wallets {
		address, initCode, err := CreateSmartAccount(&w, seedworks.OptimismSepolia)
		if err != nil {
			t.Errorf("Failed to create account: %v", err)
		}
		if address == "" {
			t.Error("Expected account to be created, but got empty string")
		}
		fmt.Printf("address: %v\n", address)
		fmt.Printf("initCode: %v\n", initCode)
	}
}
