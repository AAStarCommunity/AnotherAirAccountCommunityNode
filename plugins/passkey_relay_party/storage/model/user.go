package model

import (
	"another_node/internal/community/account"
	"another_node/internal/community/storage"
	"another_node/plugins/passkey_relay_party/conf"
	"another_node/plugins/passkey_relay_party/seedworks"
	"encoding/json"
	"time"
)

type User struct {
	storage.BaseData
	seedworkUser *seedworks.User   `gorm:"-"`
	hdWallet     *account.HdWallet `gorm:"-"`
	Email        string            `json:"email" gorm:"type:varchar(255);unique_index"`
	LastLoginAt  *time.Time        `json:"last_login_at" gorm:"column:last_login_at"`
	rawdata      string            `json:"-" gorm:"type:tinytext;column:rawdata"`
}

func (u *User) TableName() string {
	return "passkey_users"
}

func NewUser(email string, user *seedworks.User, hdwallet *account.HdWallet) *User {
	vault := conf.Get().VaultSecret
	obj := struct {
		U *seedworks.User   `json:"u"`
		H *account.HdWallet `json:"h"`
	}{
		U: user,
		H: hdwallet,
	}

	se, _ := json.Marshal(obj)
	rawdata, _ := seedworks.Encrypt([]byte(vault), se)

	return &User{
		Email:        email,
		seedworkUser: user,
		hdWallet:     hdwallet,
		rawdata:      string(rawdata),
	}
}
