package chain

import (
	"another_node/conf"
	"another_node/internal/community/account"
	"another_node/internal/global_const"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	if testing.Short() {
		return
	}
	configPath := "../../../conf/network_config.json"
	conf.InitNetworkConfig(configPath)
	w, err := account.NewHdWallet(account.HierarchicalPath(account.HierarchicalPath_ETH))
	if err != nil {
		t.Errorf("Failed to create account: %v", err)
	}

	address, err := CreateSmartAccount(w, global_const.OptimismSepolia)
	if err != nil {
		t.Errorf("Failed to create account: %v", err)
	}
	if address == "" {
		t.Error("Expected account to be created, but got empty string")
	}
	t.Logf("address: %v", address)

	// test code
}
