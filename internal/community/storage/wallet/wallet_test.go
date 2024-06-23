package wallet_storage

import (
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
	reflect.DeepEqual(a, wallet)
}
