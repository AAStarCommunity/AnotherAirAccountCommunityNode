package storage

import (
	"reflect"
	"testing"
)

func TestGetAllMembers(t *testing.T) {

	member1 := &Member{
		HashedAccount:   "hashedAccount1",
		RpcAddress:      "rpcAddress1",
		RpcPort:         12345,
		PublicKey:       "publicKey1",
		PrivateKeyVault: func() *string { s := "privateKeyVault1"; return &s }(),
	}
	member2 := &Member{
		HashedAccount:   "hashedAccount2",
		RpcAddress:      "rpcAddress2",
		RpcPort:         uint16(54321),
		PublicKey:       "publicKey2",
		PrivateKeyVault: nil,
	}
	UpsertMember(member1.HashedAccount, member1.PublicKey, *member1.PrivateKeyVault, member1.RpcAddress, member1.RpcPort, &member1.Version)
	UpsertMember(member2.HashedAccount, member2.PublicKey, "", member2.RpcAddress, member2.RpcPort, &member2.Version)

	// Call the GetAllMembers function
	members, err := GetAllMembers()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Check the returned members
	expectedMembers := []Member{*member1, *member2}
	if !reflect.DeepEqual(members, expectedMembers) {
		t.Errorf("Expected %v, but got %v", expectedMembers, members)
	}
}
