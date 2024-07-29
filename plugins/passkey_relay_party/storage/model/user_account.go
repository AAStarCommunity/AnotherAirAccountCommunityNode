package model

import (
	"another_node/internal/community/storage"
)

type UserAccount struct {
	storage.BaseData
	Email      string `json:"email" gorm:"type:varchar(255);unique_index"`
	InitCode   string `json:"init_code" gorm:"type:text;column:init_code"`
	Address    string `json:"address" gorm:"column:address"`
	EoaAddress string `json:"eoa_address" gorm:"column:eoa_address"`
	Chain      string `json:"chain" gorm:"column:chain"`
}

func (UserAccount) TableName() string {
	return "use_accounts"
}
