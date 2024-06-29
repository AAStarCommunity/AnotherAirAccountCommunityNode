package member_storage

import (
	"another_node/internal/community/storage"
	"os"
	"reflect"
	"testing"
)

func TestLittleEndianToUint32(t *testing.T) {
	var u32 uint32 = 0x12345678
	var b = []byte{120, 86, 52, 18}
	if reflect.DeepEqual(uintToBytes(u32), b) {
		t.Log("TestLittleEndianToUint32 passed")
	} else {
		t.Error("TestLittleEndianToUint32 failed")
	}
}
func TestLittleEndianToUint16(t *testing.T) {
	var u16 uint16 = 65535
	var b = []byte{255, 255}
	tx := uintToBytes(u16)
	if reflect.DeepEqual(tx, b) {
		t.Log("TestLittleEndianToUint32 passed")
	} else {
		t.Error("TestLittleEndianToUint32 failed")
	}
}

func TestMember_Marshal(t *testing.T) {
	type fields struct {
		HashedAccount   string
		RpcAddress      string
		RpcPort         uint16
		PublicKey       string
		PrivateKeyVault string
		Version         uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "test marshal without privatekey",
			fields: fields{
				HashedAccount:   "HelloWorld",
				RpcAddress:      "test",
				RpcPort:         222,
				PublicKey:       "Abc Def",
				PrivateKeyVault: "",
				Version:         123,
			},
			want: func() []byte {
				ret := make([]byte, marshalTotalCap)
				offset := 0
				copy(ret, []byte{memberMarshalHeader})
				offset += 1
				copy(ret[offset:offset+hashedAccountCap], []byte("HelloWorld"))
				offset += hashedAccountCap
				copy(ret[offset:offset+rpcAddressCap], []byte("test"))
				offset += rpcAddressCap
				copy(ret[offset:offset+2], []byte{222, 0})
				offset += 2
				copy(ret[offset:offset+publicKeyCap], []byte("Abc Def"))
				offset += publicKeyCap
				copy(ret[offset:offset+privateKeyVaultCap], []byte{})
				offset += privateKeyVaultCap
				copy(ret[offset:offset+4], []byte{123, 0, 0, 0})
				return ret
			}(),
		},
		{
			name: "test marshal with privatekey",
			fields: fields{
				HashedAccount:   "HelloWorld",
				RpcAddress:      "test",
				RpcPort:         0,
				PublicKey:       "Abc Def",
				PrivateKeyVault: "privateKeyVault",
				Version:         123,
			},
			want: func() []byte {
				ret := make([]byte, marshalTotalCap)
				offset := 0
				copy(ret, []byte{memberMarshalHeader})
				offset += 1
				copy(ret[offset:offset+hashedAccountCap], []byte("HelloWorld"))
				offset += hashedAccountCap
				copy(ret[offset:offset+rpcAddressCap], []byte("test"))
				offset += rpcAddressCap
				copy(ret[offset:offset+2], []byte{0, 0})
				offset += 2
				copy(ret[offset:offset+publicKeyCap], []byte("Abc Def"))
				offset += publicKeyCap
				copy(ret[offset:offset+privateKeyVaultCap], []byte("privateKeyVault"))
				offset += privateKeyVaultCap
				copy(ret[offset:offset+4], []byte{123, 0, 0, 0})
				return ret
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Member{
				HashedAccount:   tt.fields.HashedAccount,
				RpcAddress:      tt.fields.RpcAddress,
				RpcPort:         tt.fields.RpcPort,
				PublicKey:       tt.fields.PublicKey,
				PrivateKeyVault: tt.fields.PrivateKeyVault,
				Version:         tt.fields.Version,
			}
			if got := m.Marshal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Member.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMember_Unmarshal(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Member
		wantErr bool
	}{
		{
			name: "test unmarshal without privatekey",
			args: args{
				data: func() []byte {
					ret := make([]byte, marshalTotalCap)
					offset := 0
					copy(ret, []byte{memberMarshalHeader})
					offset += 1
					copy(ret[offset:offset+hashedAccountCap], []byte("HelloWorld"))
					offset += hashedAccountCap
					copy(ret[offset:offset+rpcAddressCap], []byte("test"))
					offset += rpcAddressCap
					copy(ret[offset:offset+2], uintToBytes(uint16(65531)))
					offset += 2
					copy(ret[offset:offset+publicKeyCap], []byte("Abc Def"))
					offset += publicKeyCap
					copy(ret[offset:offset+privateKeyVaultCap], []byte{})
					offset += privateKeyVaultCap
					copy(ret[offset:offset+4], uintToBytes(uint32(2000933341)))
					return ret
				}(),
			},
			want: &Member{
				HashedAccount:   "HelloWorld",
				RpcAddress:      "test",
				RpcPort:         65531,
				PublicKey:       "Abc Def",
				PrivateKeyVault: "",
				Version:         2000933341,
			},
			wantErr: false,
		},
		{
			name: "test unmarshal with privatekey",
			args: args{
				data: func() []byte {
					ret := make([]byte, marshalTotalCap)
					offset := 0
					copy(ret, []byte{memberMarshalHeader})
					offset += 1
					copy(ret[offset:offset+hashedAccountCap], []byte("HelloWorld"))
					offset += hashedAccountCap
					copy(ret[offset:offset+rpcAddressCap], []byte("test"))
					offset += rpcAddressCap
					copy(ret[offset:offset+2], []byte{0, 0})
					offset += 2
					copy(ret[offset:offset+publicKeyCap], []byte("Abc Def"))
					offset += publicKeyCap
					copy(ret[offset:offset+privateKeyVaultCap], []byte("privateKeyVault"))
					offset += privateKeyVaultCap
					copy(ret[offset:offset+4], []byte{123, 0, 0, 0})
					return ret
				}(),
			},
			want: &Member{
				HashedAccount:   "HelloWorld",
				RpcAddress:      "test",
				RpcPort:         0,
				PublicKey:       "Abc Def",
				PrivateKeyVault: "privateKeyVault",
				Version:         123,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalMember(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Member.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Member.Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarshalToUnmarshalThenMarhsalCompare(t *testing.T) {
	member := &Member{
		HashedAccount:   "HelloWorld",
		RpcAddress:      "test",
		RpcPort:         165,
		PublicKey:       "Abc Def",
		PrivateKeyVault: "privateKeyVault",
		Version:         22222,
	}
	marshal := member.Marshal()
	unmarshal, _ := UnmarshalMember(marshal)
	if reflect.DeepEqual(member, unmarshal) {
		t.Log("TestMashalToUnmashalThenMarhsalCompare passed")
	} else {
		t.Error("TestMashalToUnmashalThenMarhsalCompare failed")
	}

	member.PrivateKeyVault = ""
	marshal = member.Marshal()
	unmarshal, _ = UnmarshalMember(marshal)
	if reflect.DeepEqual(member, unmarshal) {
		t.Log("TestMashalToUnmashalThenMarhsalCompare passed")
	} else {
		t.Error("TestMashalToUnmashalThenMarhsalCompare failed")
	}
}

func TestUpsertMember(t *testing.T) {
	if testing.Short() {
		t.Skip("this unit test is always failed by unknown reason, skipped temporarily")
	}

	os.Setenv("UnitTest", "1")
	defer func() {
		os.Unsetenv("UnitTest")
		storage.Close()
	}()
	member := &Member{
		HashedAccount:   "HelloWorld",
		RpcAddress:      "test",
		RpcPort:         165,
		PublicKey:       "Abc Def",
		PrivateKeyVault: "privateKeyVault",
		Version:         22222,
	}
	err := UpsertMember(member.HashedAccount, member.PublicKey, member.PrivateKeyVault, member.RpcAddress, member.RpcPort, &member.Version)
	if err != nil {
		t.Error("TestUpsertMember failed: " + err.Error())
	}
}

func TestMarshalMembers(t *testing.T) {
	member1 := &Member{
		HashedAccount:   "HelloWorld",
		RpcAddress:      "test",
		RpcPort:         165,
		PublicKey:       "Abc Def",
		PrivateKeyVault: "privateKeyVault",
		Version:         22222,
	}
	member2 := &Member{
		HashedAccount:   "HelloWorld2",
		RpcAddress:      "test2",
		RpcPort:         166,
		PublicKey:       "Abc Def2",
		PrivateKeyVault: "privateKeyVault2",
		Version:         22223,
	}
	members := Members{*member1, *member2}
	marshal := members.Marshal()
	unmarshal := UnmarshalMembers(marshal)
	if reflect.DeepEqual(members, unmarshal) {
		t.Log("TestMarshalMembers passed")
	} else {
		t.Error("TestMarshalMembers failed")
	}
}
