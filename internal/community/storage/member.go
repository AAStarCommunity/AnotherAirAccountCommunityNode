package storage

import (
	"another_node/conf"
	"errors"
	"time"

	"gorm.io/gorm"
)

// Member represent a web2 account
type Member struct {
	Model
	HashedAccount   string  `gorm:"column:hashed_account; type:varchar(1024); not null; uniqueIndex"`
	RpcAddress      string  `gorm:"column:rpc_address; type:varchar(128); not null"`
	RpcPort         int     `gorm:"column:rpc_port; type:int; not null; default:0"`
	PublicKey       string  `gorm:"column:public_key; type:varchar(1024)"`
	PrivateKeyVault *string `gorm:"column:private_key_vault; type:varchar(1024); null"`
}

func (m *Member) TableName() string {
	return "members"
}

// UpsertMember update a member if exists and newer than old by version
func UpsertMember(hashedAccount, publicKey, privateKey, rpcAddress string, rpcPort int, version *int) error {
	db := conf.GetDbClient()

	return db.Transaction(func(tx *gorm.DB) error {
		var member Member
		err := tx.Where("hashed_account = ?", hashedAccount).First(&member).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Omit("updated_at", "id").Create(&Member{
				Model: Model{
					CreatedAt: time.Now(),
					Version:   uint(*version),
				},
				HashedAccount: hashedAccount,
				PublicKey:     publicKey,
				PrivateKeyVault: func() *string {
					if len(privateKey) == 0 {
						return nil
					} else {
						return &privateKey
					}
				}(),
				RpcAddress: rpcAddress,
				RpcPort:    rpcPort,
			}).Error
		} else {
			if member.Version >= uint(*version) {
				*version = int(member.Version)
				return nil
			}

			if publicKey == nil {
				tx.Omit("public_key")
			}
			if privateKey == nil {
				tx.Omit("private_key_vault")
			}
			if rpcAddress == nil || rpcPort == nil {
				tx.Omit("rpc_address")
				tx.Omit("rpc_port")
			}
			err := tx.Where("id=?", member.ID).Updates(Member{
				PublicKey: func() string {
					if publicKey != nil {
						return *publicKey
					} else {
						return member.PublicKey
					}
				}(),
				PrivateKeyVault: privateKey,
				RpcAddress: func() string {
					if rpcAddress != nil {
						return *rpcAddress
					} else {
						return member.RpcAddress
					}
				}(),
				RpcPort: func() int {
					if rpcPort != nil {
						return *rpcPort
					} else {
						return member.RpcPort
					}
				}(),
			}).Error

			return err
		}
	})
}

// TryFindMember find a member by hashed account
func TryFindMember(hashedAccount string) (*Member, error) {
	db := conf.GetDbClient()

	var member Member
	err := db.Where("hashed_account = ?", hashedAccount).First(&member).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &member, err
}
