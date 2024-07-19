package common_util

import (
	"another_node/internal/community/account"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestEthereumSignHexStr(t *testing.T) {
	//publicKey := privateKey.Public()
	hdWallet, err := account.NewHdWallet(account.HierarchicalPath_ETH)
	if err != nil {
		t.Fatal(err)
	}
	privateKeyStr := hdWallet.PrivateKey()
	address := hdWallet.Address()
	t.Logf("address: %s", address)
	t.Logf("privateKeyStr: %s", privateKeyStr)
	privateKeyECDSA, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		t.Fatal(err)
	}
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	addressAno := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	t.Logf("addressAno: %s", addressAno)
	//sign, err := EthereumSignHexStr(private
}
