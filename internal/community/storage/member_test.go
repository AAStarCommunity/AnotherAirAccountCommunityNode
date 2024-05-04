package storage

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	privateKey := "privateKey"
	member := &Member{
		HashedAccount:   "hashedAccount",
		RpcAddress:      "rpcAddress",
		RpcPort:         12345,
		PublicKey:       "publicKey",
		PrivateKeyVault: &privateKey,
	}

	data := member.Marshal()

	hashedAccount := fmt.Sprintf("%-*s", hashedAccountCapacity, "hashedAccount")
	rpcAddress := fmt.Sprintf("%-*s", rpcAddressCapacity, "rpcAddress")
	rpcPort := fmt.Sprintf("%-*d", rpcPortCapacity, 12345)
	publicKey := fmt.Sprintf("%-*s", publicKeyCapacity, "publicKey")
	privateKeyVault := fmt.Sprintf("%-*s", privateKeyVaultCapacity, "privateKey")

	expected := []byte(hashedAccount + rpcAddress + rpcPort + publicKey + privateKeyVault)

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Expected %v, but got %v", expected, data)
	}
}

func TestUnmarshal(t *testing.T) {
	privateKey := "privateKey"

	hashedAccount := fmt.Sprintf("%-*s", hashedAccountCapacity, "hashedAccount")
	rpcAddress := fmt.Sprintf("%-*s", rpcAddressCapacity, "rpcAddress")
	rpcPort := fmt.Sprintf("%-*d", rpcPortCapacity, 12345)
	publicKey := fmt.Sprintf("%-*s", publicKeyCapacity, "publicKey")
	privateKeyVault := fmt.Sprintf("%-*s", privateKeyVaultCapacity, privateKey)

	data := []byte(hashedAccount + rpcAddress + rpcPort + publicKey + privateKeyVault)
	member, err := Unmarshal(data)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	expected := &Member{
		HashedAccount:   "hashedAccount",
		RpcAddress:      "rpcAddress",
		RpcPort:         12345,
		PublicKey:       "publicKey",
		PrivateKeyVault: &privateKey,
	}
	if !reflect.DeepEqual(member, expected) {
		t.Errorf("Expected %v, but got %v", expected, member)
	}
}
