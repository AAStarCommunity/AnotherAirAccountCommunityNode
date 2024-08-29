package model

import "another_node/internal/community/storage"

type Passkey struct {
	storage.BaseData
	AirAccountId int64  `json:"airaccount_id" gorm:"column:airaccount_id"`
	PasskeyId    string `json:"passkey_id" gorm:"type:varchar(128)"`
	PasskeyRawId string `json:"passkey_rawid" gorm:"type:varchar(128)"`
	PublicKey    string `json:"public_key" gorm:"type:varchar(512)"`
	Algorithm    string `json:"algorithm" gorm:"type:varchar(64)"`
	Origin       string `json:"origin" gorm:"type:varchar(255)"`
}

func (Passkey) TableName() string {
	return "passkeys"
}
