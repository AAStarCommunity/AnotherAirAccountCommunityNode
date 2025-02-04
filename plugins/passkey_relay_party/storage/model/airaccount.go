package model

import "another_node/internal/community/storage"

type AirAccount struct {
	storage.BaseData
	WebAuthnID       string            `json:"webauthn_id" gorm:"column:webauthn_id;type:varchar(255)"` // the same value when first registered
	Email            string            `json:"email" gorm:"type:varchar(255)"`
	Facebook         string            `json:"facebook" gorm:"type:varchar(255)"`
	Twitter          string            `json:"twitter" gorm:"type:varchar(255)"`
	EoaAddress       string            `json:"eoa_address" gorm:"type:varchar(255)"`
	ZuzaluCityID     string            `json:"zuzalu_city_id" gorm:"column:zuzalu_city_id;type:varchar(255)"`
	HdWallet         []HdWallet        `json:"hdwallet" gorm:"foreignKey:AirAccountID"`
	Passkeys         []Passkey         `json:"passkeys" gorm:"foreignKey:AirAccountID"`
	AirAccountChains []AirAccountChain `json:"airaccount_chains" gorm:"foreignKey:AirAccountID"`
}

func (AirAccount) TableName() string {
	return "airaccounts"
}
