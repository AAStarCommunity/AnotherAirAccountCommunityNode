package model

import "another_node/internal/community/storage"

type AirAccount struct {
	storage.BaseData
	Email    string `json:"email" gorm:"type:varchar(255)"`
	Facebook string `json:"facebook" gorm:"type:varchar(255)"`
	Twitter  string `json:"twitter" gorm:"type:varchar(255)"`
}

func (AirAccount) TableName() string {
	return "airaccounts"
}
