package model

import "another_node/internal/community/storage"

type PasskeySession struct {
	storage.BaseData
}

func (PasskeySession) TableName() string {
	return "passkey_sessions"
}
