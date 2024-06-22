package member_storage

import (
	"another_node/internal/community/storage"
	"os"
	"reflect"
	"testing"
)

func TestGetAllMembers(t *testing.T) {
	os.Setenv("UnitTest", "1")
	defer func() {
		os.Unsetenv("UnitTest")
		storage.Close()
	}()

	member1 := &Member{
		HashedAccount:   "hashedAccount1",
		RpcAddress:      "rpcAddress1",
		RpcPort:         12345,
		PublicKey:       "publicKey1",
		PrivateKeyVault: "privateKeyVault1",
	}
	member2 := &Member{
		HashedAccount:   "hashedAccount2",
		RpcAddress:      "rpcAddress2",
		RpcPort:         uint16(54321),
		PublicKey:       "publicKey2",
		PrivateKeyVault: "",
	}

	UpsertMember(member1.HashedAccount, member1.PublicKey, member1.PrivateKeyVault, member1.RpcAddress, member1.RpcPort, &member1.Version)
	UpsertMember(member2.HashedAccount, member2.PublicKey, "", member2.RpcAddress, member2.RpcPort, &member2.Version)

	// Call the GetAllMembers function
	members := GetMembers(0, ^uint32(0))

	// Check the returned members
	expectedMembers := Members{*member1, *member2}
	if !reflect.DeepEqual(members, expectedMembers) {
		t.Errorf("Expected %v, but got %v", expectedMembers, members)
	}
}
