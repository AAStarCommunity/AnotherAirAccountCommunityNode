package model

import "another_node/internal/community/storage"

type AirAccountChain struct {
	storage.BaseData
	AirAccountId int64  `json:"airaccount_id" gorm:"column:airaccount_id"`
	InitCode     string `json:"init_code" gorm:"type:text"`
	AA_Address   string `json:"aa_address" gorm:"type:varchar(255)"`
	EOA_Address  string `json:"eoa_address" gorm:"type:varchar(255)"`
	Chain        string `json:"chain" gorm:"type:varchar(64)"`
	WalletVault  string `json:"wallet_vault" gorm:"type:text"`
}

func (AirAccountChain) TableName() string {
	return "airaccount_chains"
}
