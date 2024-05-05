package storage

import (
	"another_node/conf"
	"reflect"
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func TestGetAllMembers(t *testing.T) {

	member1 := &Member{
		HashedAccount:   "hashedAccount1",
		RpcAddress:      "rpcAddress1",
		RpcPort:         12345,
		PublicKey:       "publicKey1",
		PrivateKeyVault: nil,
	}
	member2 := &Member{
		HashedAccount:   "hashedAccount2",
		RpcAddress:      "rpcAddress2",
		RpcPort:         67890,
		PublicKey:       "publicKey2",
		PrivateKeyVault: nil,
	}
	data1 := member1.Marshal()
	data2 := member2.Marshal()

	func() {
		// os.Setenv("UnitTest", "1")
		stor, _ := conf.GetStorage()
		db, _ := leveldb.Open(stor, &opt.Options{})
		defer func() {
			stor.Close()
			db.Close()
		}()

		if err := db.Put([]byte(memberKey(member1)), data1, nil); err != nil {
			t.Fatal(err)
		}
		if err := db.Put([]byte(memberKey(member2)), data2, nil); err != nil {
			t.Fatal(err)
		}
	}()

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
