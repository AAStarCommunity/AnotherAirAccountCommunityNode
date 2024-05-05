package storage

import (
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
		PrivateKeyVault *string
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
				PrivateKeyVault: nil,
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
				PrivateKeyVault: func() *string { s := "privateKeyVault"; return &s }(),
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
				PrivateKeyVault: nil,
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
				PrivateKeyVault: func() *string { s := "privateKeyVault"; return &s }(),
				Version:         123,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal(tt.args.data)
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
		PrivateKeyVault: func() *string { s := "privateKeyVault"; return &s }(),
		Version:         22222,
	}
	marshal := member.Marshal()
	unmarshal, _ := Unmarshal(marshal)
	if reflect.DeepEqual(member, unmarshal) {
		t.Log("TestMashalToUnmashalThenMarhsalCompare passed")
	} else {
		t.Error("TestMashalToUnmashalThenMarhsalCompare failed")
	}

	member.PrivateKeyVault = nil
	marshal = member.Marshal()
	unmarshal, _ = Unmarshal(marshal)
	if reflect.DeepEqual(member, unmarshal) {
		t.Log("TestMashalToUnmashalThenMarhsalCompare passed")
	} else {
		t.Error("TestMashalToUnmashalThenMarhsalCompare failed")
	}
}

func TestMemberIndex(t *testing.T) {
	key1 := memberIndexKey()
	key2 := memberIndexKey()
	if len(key1) > 1 && len(key2) > 1 && key1 != key2 {
		t.Log("TestMemberIndex passed")
	} else {
		t.Error("TestMemberIndex failed")
	}
}

func TestUpsertMember(t *testing.T) {
	member := &Member{
		HashedAccount:   "HelloWorld",
		RpcAddress:      "test",
		RpcPort:         165,
		PublicKey:       "Abc Def",
		PrivateKeyVault: func() *string { s := "privateKeyVault"; return &s }(),
		Version:         22222,
	}
	err := UpsertMember(member.HashedAccount, member.PublicKey, *member.PrivateKeyVault, member.RpcAddress, member.RpcPort, &member.Version)
	if err != nil {
		t.Error("TestUpsertMember failed")
	}
}
