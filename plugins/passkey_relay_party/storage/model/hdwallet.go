package model

import "another_node/internal/community/storage"

type HdWallet struct {
	storage.BaseData
	AirAccount   AirAccount `json:"-" gorm:"foreignKey:AirAccountID;references:ID"`
	AirAccountID int64      `json:"-" gorm:"column:airaccount_id"`
	WalletVault  string     `json:"wallet_vault" gorm:"type:text"`
}

func (HdWallet) TableName() string {
	return "hdwallets"
}
