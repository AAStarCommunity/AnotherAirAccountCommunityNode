package model

import "another_node/internal/community/storage"

type AirAccount struct {
	storage.BaseData
	Email            string            `json:"email" gorm:"type:varchar(255)"`
	Facebook         string            `json:"facebook" gorm:"type:varchar(255)"`
	Twitter          string            `json:"twitter" gorm:"type:varchar(255)"`
	HdWallet         HdWallet          `json:"hdwallet" gorm:"foreignKey:AirAccountID"`
	Passkeys         []Passkey         `json:"passkeys" gorm:"foreignKey:AirAccountID"`
	AirAccountChains []AirAccountChain `json:"airaccount_chains" gorm:"foreignKey:AirAccountID"`
}

func (AirAccount) TableName() string {
	return "airaccounts"
}
