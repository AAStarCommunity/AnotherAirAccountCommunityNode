package impl

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
)

func TestSepoliaCalling(t *testing.T) {
	apiKey := "test-api-key"
	caller, err := NewInfuraCaller(InfuraSepolia, apiKey)
	if err != nil {
		t.Errorf("Failed to create InfuraCaller: %v", err)
	}

	err = caller.Call(func(client *ethclient.Client) interface{} {
		return "string"
	})

	if err != nil {
		t.Errorf("Failed to call the client: %v", err)
	}

	err = caller.Call(func(client *ethclient.Client) interface{} {
		return nil
	})

	if err != nil {
		t.Errorf("Failed to call the client: %v", err)
	}

	err = caller.Call(func(client *ethclient.Client) interface{} {
		return errors.New("error")
	})
	if err == nil {
		t.Errorf("Failed to call the client: %v", err)
	}
}
