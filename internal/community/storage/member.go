package storage

import (
	"another_node/conf"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

// Member represent a web2 account
type Member struct {
	HashedAccount   string
	RpcAddress      string
	RpcPort         int
	PublicKey       string
	PrivateKeyVault *string
	Version         uint
}

const (
	hashedAccountCapacity   = 128
	rpcAddressCapacity      = 128
	rpcPortCapacity         = 5
	publicKeyCapacity       = 1024
	privateKeyVaultCapacity = 2048
)

func (m *Member) Marshal() []byte {
	hashedAccount := fmt.Sprintf("%-*s", hashedAccountCapacity, m.HashedAccount)
	if len(hashedAccount) > hashedAccountCapacity {
		hashedAccount = hashedAccount[:hashedAccountCapacity]
	}

	rpcAddress := fmt.Sprintf("%-*s", rpcAddressCapacity, m.RpcAddress)
	if len(rpcAddress) > rpcAddressCapacity {
		rpcAddress = rpcAddress[:rpcAddressCapacity]
	}

	rpcPort := fmt.Sprintf("%-*d", rpcPortCapacity, m.RpcPort)
	if len(rpcPort) > rpcPortCapacity {
		rpcPort = rpcPort[:rpcPortCapacity]
	}

	publicKey := fmt.Sprintf("%-*s", publicKeyCapacity, m.PublicKey)
	if len(publicKey) > publicKeyCapacity {
		publicKey = publicKey[:publicKeyCapacity]
	}

	privateKeyVault := fmt.Sprintf("%-*s", privateKeyVaultCapacity, *m.PrivateKeyVault)
	if len(privateKeyVault) > privateKeyVaultCapacity {
		privateKeyVault = privateKeyVault[:privateKeyVaultCapacity]
	}

	result := hashedAccount + rpcAddress + rpcPort + publicKey + privateKeyVault
	return []byte(result)
}

func Unmarshal(data []byte) (*Member, error) {
	if len(data) < (hashedAccountCapacity + rpcAddressCapacity + rpcPortCapacity + publicKeyCapacity + privateKeyVaultCapacity) {
		return nil, errors.New("data is too short to unmarshal into Member")
	}

	hashedAccount := strings.TrimSpace(string(data[:hashedAccountCapacity]))
	rpcAddress := strings.TrimSpace(string(data[hashedAccountCapacity : hashedAccountCapacity+rpcAddressCapacity]))
	rpcPort, err := strconv.Atoi(strings.TrimSpace(string(data[hashedAccountCapacity+rpcAddressCapacity : hashedAccountCapacity+rpcAddressCapacity+rpcPortCapacity])))
	if err != nil {
		return nil, err
	}
	publicKey := strings.TrimSpace(string(data[hashedAccountCapacity+rpcAddressCapacity+rpcPortCapacity : hashedAccountCapacity+rpcAddressCapacity+rpcPortCapacity+publicKeyCapacity]))
	privateKeyVault := strings.TrimSpace(string(data[hashedAccountCapacity+rpcAddressCapacity+rpcPortCapacity+publicKeyCapacity:]))

	return &Member{
		HashedAccount:   hashedAccount,
		RpcAddress:      rpcAddress,
		RpcPort:         rpcPort,
		PublicKey:       publicKey,
		PrivateKeyVault: &privateKeyVault,
	}, nil
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

// UpsertMember update a member if exists and newer than old by version
func UpsertMember(hashedAccount, publicKey, privateKey, rpcAddress string, rpcPort int, version *int) error {
	if stor, err := conf.GetStorage(); err != nil {
		return err
	} else {
		if db, err := leveldb.Open(stor, nil); err != nil {
			return err
		} else {
			defer db.Close()

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
				Version: uint(*version),
			}
			if oldMemberByte, err := db.Get([]byte(hashedAccount), nil); err != nil {
				if errors.Is(err, leveldb.ErrNotFound) {
					return db.Put([]byte(hashedAccount), newMember.Marshal(), nil)
				}
				return err
			} else {
				if oldMember, err := Unmarshal(oldMemberByte); err != nil {
					return err
				} else {
					newMember = compareAndUpdateMember(oldMember, newMember)

					return db.Put([]byte(hashedAccount), newMember.Marshal(), nil)
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
			defer db.Close()
			if member, err := db.Get([]byte(hashedAccount), nil); err != nil {
				if errors.Is(err, leveldb.ErrNotFound) {
					return nil, nil
				}
				return nil, err
			} else {
				return Unmarshal(member)
			}
		}
	}
}
