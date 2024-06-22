package impl

import (
	"another_node/internal/community/wallet"
	"testing"
)

func TestNewAlchemyProvider(t *testing.T) {
	apiKey := "test-api-key"
	provider, err := NewAlchemyProvider(apiKey)
	if err != nil {
		t.Errorf("Failed to create AlchemyProvider: %v", err)
	}

	if provider.GetRpc() != "https://eth-sepolia.g.alchemy.com/v2/"+apiKey {
		t.Errorf("Expected rpc to be %s, but got %s", apiKey, provider.GetRpc())
	}
}

func TestAlchemyProvider_CreateAccount(t *testing.T) {
	if testing.Short() {
		return
	}

	apiKey := "<API-KEY>"
	provider, err := NewAlchemyProvider(apiKey)
	if err != nil {
		t.Errorf("Failed to create account: %v", err)
	}
	w, err := wallet.NewHdWallet(wallet.HierarchicalPath(wallet.HierarchicalPath_Main_ETH_TestNet))
	if err != nil {
		t.Errorf("Failed to create account: %v", err)
	}

	account, err := provider.CreateAccount(w)
	if err != nil {
		t.Errorf("Failed to create account: %v", err)
	}

	if account == "" {
		t.Error("Expected account to be created, but got empty string")
	}
}
