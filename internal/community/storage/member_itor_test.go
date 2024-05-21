package storage

import (
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestGetAllMembers(t *testing.T) {

	dir := os.TempDir()
	uuid := uuid.New().String()
	os.Setenv("storage", dir+"/testing.dat/"+uuid)
	defer func() {
		os.Unsetenv("storage")
		os.RemoveAll(dir + "/testing.dat")
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
	expectedMembers := []Member{*member1, *member2}
	if !reflect.DeepEqual(members, expectedMembers) {
		t.Errorf("Expected %v, but got %v", expectedMembers, members)
	}
}
