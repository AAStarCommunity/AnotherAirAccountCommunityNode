package model

import "another_node/internal/community/storage"

type Passkey struct {
	storage.BaseData
	AirAccountID uint       `json:"-" gorm:"column:airaccount_id"`
	AirAccount   AirAccount `json:"-" gorm:"foreignKey:AirAccountID;references:ID"`
	CredentialId string     `json:"credential_id" gorm:"type:varchar(128)"`
	PublicKey    string     `json:"public_key" gorm:"type:varchar(512)"`
	Algorithm    string     `json:"algorithm" gorm:"type:varchar(64)"`
	Origin       string     `json:"origin" gorm:"type:varchar(255)"`
	Rawdata      string     `json:"-" gorm:"type:text;column:rawdata"`
}

func (Passkey) TableName() string {
	return "passkeys"
}
