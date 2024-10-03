package model

import "another_node/internal/community/storage"

type AirAccountChain struct {
	storage.BaseData
	AirAccountID int64       `json:"-" gorm:"column:airaccount_id"`
	AirAccount   AirAccount `json:"-" gorm:"foreignKey:AirAccountID;references:ID"`
	InitCode     string     `json:"init_code" gorm:"type:text"`
	AA_Address   string     `json:"aa_address" gorm:"type:varchar(255)"`
	ChainName    string     `json:"chain_name" gorm:"type:varchar(50);column:chain_name"`
	Alias        string     `json:"alias" gorm:"type:varchar(50)"`
	FromWalletID int64       `json:"-" gorm:"column:hdwallet_id"`
	FromWallet   HdWallet   `json:"-" gorm:"foreignKey:FromWalletID;references:ID"`
}

func (AirAccountChain) TableName() string {
	return "airaccount_chains"
}
