package common_util

import (
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestEthereumSignHexStr(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	//publicKey := privateKey.Public()

	res, err := EthereumSignHexStr("0x123456", privateKey)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("res: %s\n", res)
}
