package model

import (
	"another_node/internal/community/storage"
	"time"
)

type User struct {
	storage.BaseData
	Email       string     `json:"email" gorm:"type:varchar(255);unique_index"`
	LastLoginAt *time.Time `json:"last_login_at" gorm:"column:last_login_at"`
	Rawdata     string     `json:"-" gorm:"type:tinytext;column:rawdata"`
}

func (User) TableName() string {
	return "passkey_users"
}
