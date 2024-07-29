package model

import (
	"another_node/internal/community/storage"
)

type ChallengeType string

const (
	Email ChallengeType = "email"
)

type CaptchaChallenge struct {
	storage.BaseData
	Type   ChallengeType `gorm:"type:varchar(255);column:type"`
	Object string        `gorm:"type:varchar(255);column:object;index:idx_object"`
	Code   string        `gorm:"type:varchar(16);column:code"`
}

func (CaptchaChallenge) TableName() string {
	return "challenge_captchas"
}
