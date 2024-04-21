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
	HashedAccount   string `gorm:"column:hashed_account; type:varchar(1024); not null; uniqueIndex"`
	RpcAddress      string `gorm:"column:rpc_address; type:varchar(128); not null"`
	PublicKey       string `gorm:"column:public_key; type:varchar(1024)"`
	PrivateKeyVault string `gorm:"column:private_key_vault; type:varchar(1024); null"`
}

func (m *Member) TableName() string {
	return "members"
}

// UpsertMember update a member if exists and newer than old by version
func UpsertMember(hashedAccount string, publicKey, privateKey, rpcAddress *string, version int) error {
	db := conf.GetDbClient()

	return db.Transaction(func(tx *gorm.DB) error {
		var member Member
		err := tx.Where("hashed_account = ?", hashedAccount).First(&member).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Omit("updated_at", "id").Create(&Member{
				Model: Model{
					CreatedAt: time.Now(),
					Version:   uint(version),
				},
				HashedAccount:   hashedAccount,
				PublicKey:       *publicKey,
				PrivateKeyVault: *privateKey,
				RpcAddress:      *rpcAddress,
			}).Error
		} else {
			if member.Version >= uint(version) {
				return nil
			}

			if publicKey == nil {
				tx.Omit("public_key")
			}
			if privateKey == nil {
				tx.Omit("private_key_vault")
			}
			if rpcAddress == nil {
				tx.Omit("rpc_address")
			}
			err := tx.Where("id=?", member.ID).Updates(Member{
				PublicKey: func() string {
					if publicKey != nil {
						return *publicKey
					} else {
						return member.PublicKey
					}
				}(),
				PrivateKeyVault: func() string {
					if privateKey != nil {
						return *privateKey
					} else {
						return member.PrivateKeyVault
					}
				}(),
				RpcAddress: func() string {
					if rpcAddress != nil {
						return *rpcAddress
					} else {
						return member.RpcAddress
					}
				}(),
			}).Error

			return err
		}
	})
}
