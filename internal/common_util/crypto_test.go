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

func TestSignMessage(t *testing.T) {
	hash := "0x3c49667bcbfe10315d3656f50523b2685d584610398ac0ef230046cbcea70e8a"
	pk := "f0f5aec309e1fc3c26f828bb9eabc407457284bae8d012e26558038a85003eff"
	pke, _ := crypto.HexToECDSA(pk)
	sign, _ := EthereumSignHexStr(hash, pke)
	t.Logf("sign: %s", sign)
}
