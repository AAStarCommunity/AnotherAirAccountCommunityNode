package impl

import (
	"another_node/internal/community/account"
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
	w, err := account.NewHdWallet(account.HierarchicalPath(account.HierarchicalPath_ETH))
	if err != nil {
		t.Errorf("Failed to create account: %v", err)
	}

	for i := range w {
		account, err := provider.CreateAccount(&w[i])
		if err != nil {
			t.Errorf("Failed to create account: %v", err)
		}

		if account == "" {
			t.Error("Expected account to be created, but got empty string")
		}
	}
}
