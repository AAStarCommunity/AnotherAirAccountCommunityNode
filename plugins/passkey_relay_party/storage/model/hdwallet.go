package model

import "another_node/internal/community/storage"

type HdWallet struct {
	storage.BaseData
	AirAccountID uint   `json:"-" gorm:"column:airaccount_id"`
	WalletVault  string `json:"wallet_vault" gorm:"type:text"`
}

func (HdWallet) TableName() string {
	return "hdwallets"
}
