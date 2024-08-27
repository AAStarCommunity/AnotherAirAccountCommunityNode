package model

import "another_node/internal/community/storage"

// PasskeyAccount keeps user passkey information
// like public-key, origin, algorithm, etc.
type PasskeyAccount struct {
	storage.BaseData
	UserId          int64  `json:"user_id" gorm:"type:integer"` // foreign key of user.id
	Origin          string `json:"origin" gorm:"type:varchar(255)"`
	Algorithm       string `json:"algorithm" gorm:"type:varchar(255)"` // split by ',' if multiple, support -7, -8, -256
	PublicKey       string `json:"public_key" gorm:"type:text"`
	CredentialId    string `json:"credential_id" gorm:"type:varchar(255)"`
	CredentialRawId string `json:"credential_raw_id" gorm:"type:varchar(255)"`
}

func (PasskeyAccount) TableName() string {
	return "passkey_accounts"
}
