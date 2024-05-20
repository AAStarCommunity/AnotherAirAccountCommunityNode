package storage

import (
	"another_node/conf"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	hashedAccountCap    = 128
	rpcAddressCap       = 128
	publicKeyCap        = 512
	privateKeyVaultCap  = 512
	memberMarshalHeader = byte(0x01)
)

// marshalTotalCap is the total cap of marshaled member
// 1 byte for header
// 128 bytes for hashed account
// 128 bytes for rpc address
// 2 bytes for rpc port
// 512 bytes for public key
// 512 bytes for private key vault
// 4 bytes for version
var marshalTotalCap = 1 + hashedAccountCap + rpcAddressCap + 2 + publicKeyCap + privateKeyVaultCap + 4

// Member represent a web2 account
type Member struct {
	HashedAccount   string
	RpcAddress      string
	RpcPort         uint16
	PublicKey       string
	PrivateKeyVault *string
	Version         uint32
}

// uintToBytes convert uint to bytes in little endian
func uintToBytes[T uint16 | uint32 | int64](n T) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, n)
	ret := buf.Bytes()
	return ret
}

func (m *Member) Marshal() []byte {

	hashedAccountBytes := []byte(m.HashedAccount)
	if len(hashedAccountBytes) > hashedAccountCap {
		return nil
	}

	rpcAddressBytes := []byte(m.RpcAddress)
	if len(rpcAddressBytes) > rpcAddressCap {
		return nil
	}

	rpcPortBytes := uintToBytes(m.RpcPort)

	publicKeyBytes := []byte(m.PublicKey)
	if len(publicKeyBytes) > publicKeyCap {
		return nil
	}

	privateKeyVaultBytes := []byte{}
	if m.PrivateKeyVault != nil {
		privateKeyVaultBytes = []byte(*m.PrivateKeyVault)
		if len(privateKeyVaultBytes) > privateKeyVaultCap {
			return nil
		}
	}

	versionBytes := uintToBytes(m.Version)

	ret := make([]byte, marshalTotalCap)
	offset := 0
	copy(ret, []byte{memberMarshalHeader})
	offset += 1
	copy(ret[offset:offset+hashedAccountCap], hashedAccountBytes)
	offset += hashedAccountCap
	copy(ret[offset:offset+rpcAddressCap], rpcAddressBytes)
	offset += rpcAddressCap
	copy(ret[offset:offset+2], rpcPortBytes)
	offset += 2
	copy(ret[offset:offset+publicKeyCap], publicKeyBytes)
	offset += publicKeyCap
	copy(ret[offset:offset+privateKeyVaultCap], privateKeyVaultBytes)
	offset += privateKeyVaultCap
	copy(ret[offset:offset+4], versionBytes)

	return ret
}

func compareAndUpdateMember(oldMember, newMember *Member) *Member {
	if oldMember.Version >= newMember.Version {
		return oldMember
	}

	if len(newMember.PublicKey) == 0 {
		newMember.PublicKey = oldMember.PublicKey
	}

	if newMember.PrivateKeyVault == nil {
		newMember.PrivateKeyVault = oldMember.PrivateKeyVault
	}

	if len(newMember.RpcAddress) == 0 || newMember.RpcPort == 0 {
		newMember.RpcAddress = oldMember.RpcAddress
		newMember.RpcPort = oldMember.RpcPort
	}

	return newMember
}

const MemberPrefix = "member:"

func memberKey(member *Member) string {
	return fmt.Sprintf("%s%s", MemberPrefix, member.HashedAccount)
}

// UpsertMember update a member if exists and newer than old by version
func UpsertMember(hashedAccount, publicKey, privateKey, rpcAddress string, rpcPort uint16, version *uint32) error {
	if ins, err := Open(); err != nil {
		return err
	} else {
		newMember := &Member{
			HashedAccount: hashedAccount,
			RpcAddress:    rpcAddress,
			RpcPort:       rpcPort,
			PublicKey:     publicKey,
			PrivateKeyVault: func() *string {
				if len(privateKey) == 0 {
					return nil
				} else {
					return &privateKey
				}
			}(),
			Version: uint32(*version),
		}

		defer ins.Close()
		db := ins.Instance

		if oldMemberByte, err := db.Get([]byte(memberKey(newMember)), nil); err != nil {
			if errors.Is(err, leveldb.ErrNotFound) {
				if err := db.Put([]byte(memberKey(newMember)), newMember.Marshal(), nil); err != nil {
					return err
				} else {
					return nil
				}
			}
			return err
		} else {
			if oldMember, err := UnmarshalMember(oldMemberByte); err != nil {
				return err
			} else {
				newMember = compareAndUpdateMember(oldMember, newMember)
				if err := db.Put([]byte(memberKey(newMember)), newMember.Marshal(), nil); err != nil {
					return err
				} else {
					return nil
				}
			}
		}
	}
}

// TryFindMember find a member by hashed account
func TryFindMember(hashedAccount string) (*Member, error) {
	if stor, err := conf.GetStorage(); err != nil {
		return nil, err
	} else {
		if db, err := leveldb.Open(stor, nil); err != nil {
			return nil, err
		} else {
			defer func() {
				stor.Close()
				db.Close()
			}()
			if member, err := db.Get([]byte(MemberPrefix+hashedAccount), nil); err != nil {
				if errors.Is(err, leveldb.ErrNotFound) {
					return nil, nil
				}
				return nil, err
			} else {
				return UnmarshalMember(member)
			}
		}
	}
}

func MarshalMembers(m []Member) []byte {
	ret := []byte{}
	for _, member := range m {
		b := member.Marshal()
		sz := uintToBytes(uint16(len(b)))
		ret = append(ret, sz...)
		ret = append(ret, b...)
	}
	return ret
}

func UnmarshalMember(data []byte) (*Member, error) {
	if len(data) < 1 {
		return nil, errors.New("invalid data")
	}

	if data[0] != memberMarshalHeader {
		return nil, errors.New("invalid header")
	}

	if len(data) < 1+hashedAccountCap+rpcAddressCap+2+publicKeyCap+privateKeyVaultCap+4 {
		return nil, errors.New("invalid data length")
	}

	m := &Member{
		HashedAccount: strings.Trim(string(data[1:1+hashedAccountCap]), "\x00"),
		RpcAddress:    strings.Trim(string(data[1+hashedAccountCap:1+hashedAccountCap+rpcAddressCap]), "\x00"),
		RpcPort:       binary.LittleEndian.Uint16(data[1+hashedAccountCap+rpcAddressCap : 1+hashedAccountCap+rpcAddressCap+2]),
		PublicKey:     strings.Trim(string(data[1+hashedAccountCap+rpcAddressCap+2:1+hashedAccountCap+rpcAddressCap+2+publicKeyCap]), "\x00"),
		PrivateKeyVault: func() *string {
			if len(data[1+hashedAccountCap+rpcAddressCap+2+publicKeyCap:]) == 0 {
				return nil
			}
			privateKeyVault := strings.Trim(string(data[1+hashedAccountCap+rpcAddressCap+2+publicKeyCap:1+hashedAccountCap+rpcAddressCap+2+publicKeyCap+privateKeyVaultCap]), "\x00")
			if len(privateKeyVault) == 0 {
				return nil
			} else {
				return &privateKeyVault
			}
		}(),
		Version: binary.LittleEndian.Uint32(data[1+hashedAccountCap+rpcAddressCap+2+publicKeyCap+privateKeyVaultCap:]),
	}

	return m, nil
}

func UnmarshalMembers(b []byte) []Member {
	ret := []Member{}
	for len(b) > 0 {
		sz := binary.LittleEndian.Uint16(b[:2])
		b = b[2:]
		m, _ := UnmarshalMember(b[:sz])
		ret = append(ret, *m)
		b = b[sz:]
	}
	return ret
}

func InitRemoteMember(members []Member) {
	for _, member := range members {
		if err := UpsertMember(member.HashedAccount, member.PublicKey, "", member.RpcAddress, member.RpcPort, &member.Version); err != nil {
			fmt.Print("Failed to init remote member: ", err)
		}
	}
}
