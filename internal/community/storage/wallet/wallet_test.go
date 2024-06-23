package wallet_storage

import (
	"another_node/internal/community/storage"
	"os"
	"reflect"
	"testing"
)

func TestWallet_Marshal_And_Unmarshal(t *testing.T) {
	wallet := Wallet{
		mnemonic:   "test mnemonic",
		Address:    "test address",
		privateKey: "test private key",
		AAAdress:   "test aa address",
	}
	a := wallet
	marshaledWallet := wallet.Marshal()
	err := wallet.Unmarshal(marshaledWallet)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(a, wallet) {
		t.Log("TestWallet_Marshal_And_Unmarshal passed")
	} else {
		t.Error("TestWallet_Marshal_And_Unmarshal failed")
	}
}

func TestUpsert_And_Get(t *testing.T) {
	os.Setenv("UnitTest", "1")
	defer func() {
		os.Unsetenv("UnitTest")
		storage.Close()
	}()

	account := "test account"
	account2 := "test account 2"

	wallet1 := Wallet{
		mnemonic:   "test mnemonic",
		Address:    "test address",
		privateKey: "test private key",
		AAAdress:   "test aa address",
	}

	wallet2 := Wallet{
		mnemonic:   "test mnemonic 2",
		Address:    "test address 2",
		privateKey: "test private key 2",
		AAAdress:   "test aa address 2",
	}

	wallet3 := Wallet{
		mnemonic:   "test mnemonic 3",
		Address:    "test address 3",
		privateKey: "test private key 3",
		AAAdress:   "test aa address 3",
	}

	if err := UpsertWallet(&account, &wallet1); err != nil {
		t.Error("UpsertWallet failed")
	}

	UpsertWallet(&account2, &wallet3)

	if err := UpsertWallet(&account, &wallet2); err != nil {
		t.Error("UpsertWallet failed")
	}

	if wallet, err := TryFindWallet(account); err != nil {
		t.Error("GetWallet failed")
	} else {
		if len(wallet) != 2 {
			t.Error("TestUpsert_And_Get TryFindWallet failed")
		}
		if reflect.DeepEqual(wallet[0], wallet1) && reflect.DeepEqual(wallet[1], wallet2) {
			t.Log("TestUpsert_And_Get passed")
		} else {
			t.Error("TestUpsert_And_Get failed")
		}
	}
}
