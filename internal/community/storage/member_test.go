package storage

import (
	"fmt"
	"reflect"
	"strconv"
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
		Version:         1,
	}

	data := member.Marshal()

	hashedAccount := fmt.Sprintf("%-*s", hashedAccountCapacity, "hashedAccount")
	rpcAddress := fmt.Sprintf("%-*s", rpcAddressCapacity, "rpcAddress")
	rpcPort := fmt.Sprintf("%-*d", rpcPortCapacity, 12345)
	publicKey := fmt.Sprintf("%-*s", publicKeyCapacity, "publicKey")
	privateKeyVault := fmt.Sprintf("%-*s", privateKeyVaultCapacity, privateKey)

	expected := append([]byte{memberMarshalHeader}, []byte(hashedAccount+rpcAddress+rpcPort+publicKey+privateKeyVault+strconv.Itoa(int(member.Version)))...)

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
	version := "100"

	data := append([]byte{1}, []byte(hashedAccount+rpcAddress+rpcPort+publicKey+privateKeyVault+version)...)
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
		Version:         100,
	}
	if !reflect.DeepEqual(member, expected) {
		t.Errorf("Expected %v, but got %v", expected, member)
	}
}
